#!/usr/bin/env bash
set -Eeuo pipefail

usage() {
  echo "Usage: $0 <platform> <image>" >&2
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Error: command not found: $1" >&2
    exit 127
  }
}

sanitize() {
  # docker volume/log path safe-ish
  echo "$1" | tr '/:' '_' | tr -c 'a-zA-Z0-9_.-' '_'
}

linebuf() {
  # stdbuf가 있으면 라인버퍼링 강제, 없으면 그대로 실행
  if command -v stdbuf >/dev/null 2>&1; then
    stdbuf -oL -eL "$@"
  else
    "$@"
  fi
}

ensure_gotestfmt() {
  if ! command -v gotestfmt >/dev/null 2>&1; then
    need_cmd go
    # gotestfmt 설치 후 PATH에 GOPATH/bin이 없을 수 있어 보정
    export PATH="$(go env GOPATH)/bin:${PATH}"
    go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@v2.5.0
  fi
}

if [[ $# -ne 2 ]]; then
  usage
  exit 1
fi

need_cmd docker
need_cmd timeout
need_cmd awk
need_cmd tee

DOCKER_PLATFORM="$1"
DOCKER_IMAGE="$2"
# Include emulated builds in the outer budget; each test package still has 1h.
QEMU_TIMEOUT="${QEMU_TIMEOUT:-2h}"

ensure_gotestfmt

main() {
  echo "Starting go-test with QEMU for ${DOCKER_PLATFORM} with ${DOCKER_IMAGE}"

  DOCKER_PLATFORM_R="$(sanitize "$DOCKER_PLATFORM")"
  DOCKER_IMAGE_R="$(sanitize "$DOCKER_IMAGE")"

  cpu_cores="$(nproc)"
  cpus=$(( (cpu_cores + 3) / 4 ))
  (( cpus < 1 )) && cpus=1

  # 컨테이너 go env는 테스트마다 다시 물을 필요가 없어 1회만 조회
  CONTAINER_GOMODCACHE="$(docker run --rm --platform "$DOCKER_PLATFORM" "$DOCKER_IMAGE" go env GOMODCACHE)"
  CONTAINER_GOCACHE="$(docker run --rm --platform "$DOCKER_PLATFORM" "$DOCKER_IMAGE" go env GOCACHE)"

  echo "CONTAINER_GOMODCACHE: $CONTAINER_GOMODCACHE"
  echo "CONTAINER_GOCACHE:    $CONTAINER_GOCACHE"

  run_test "default" ""
  run_test "noasm"   "--tags=noasm"
}

run_test() {
  local variant="$1"
  local go_test_opts="$2"

  echo "Running tests with ${DOCKER_PLATFORM} (opt: ${go_test_opts})"

  local start_time start_time_human now
  start_time="$(date +%s)"
  start_time_human="$(date +"%Y-%m-%d %H:%M:%S")"
  now="$(date +"%Y-%m-%d %H-%M-%S")"

  local VOLUME_GOMODCACHE VOLUME_GOCACHE
  VOLUME_GOMODCACHE="$(sanitize "${DOCKER_IMAGE}-${DOCKER_PLATFORM}-gomodcache")"
  VOLUME_GOCACHE="$(sanitize "${DOCKER_IMAGE}-${DOCKER_PLATFORM}-${variant}-gocache")"

  echo "VOLUME_GOMODCACHE: $VOLUME_GOMODCACHE"
  echo "VOLUME_GOCACHE:    $VOLUME_GOCACHE"

  docker volume create "$VOLUME_GOMODCACHE" >/dev/null
  docker volume create "$VOLUME_GOCACHE" >/dev/null

  mkdir -p "log/${DOCKER_PLATFORM_R}/${DOCKER_IMAGE_R}" "out" || true

  local raw_log="log/${DOCKER_PLATFORM_R}/${DOCKER_IMAGE_R}/${variant}-${now}.log"
  local out_log="out/${DOCKER_IMAGE_R}-${variant}.log"

  # 컨테이너 내부에서 실행할 스크립트 (quoting 안정적)
  local inner_script
  inner_script=$(
    cat <<'EOS'
git config --global --add safe.directory /src || true
echo "===== GO BUILD ====="
go build -v ./...
echo "===== GO TEST ====="
go test ${GO_TEST_OPTS} -json -timeout 1h ./...
EOS
  )

  # docker run을 배열로 만들어 eval 제거
  local -a docker_cmd=(
    docker run --rm
    --platform "$DOCKER_PLATFORM"
    --volume "$PWD/..:/src"
    --volume "${VOLUME_GOMODCACHE}:${CONTAINER_GOMODCACHE}"
    --volume "${VOLUME_GOCACHE}:${CONTAINER_GOCACHE}"
    --cpus "$cpus"
    --memory 2g
    --env CGO_ENABLED=0
    --env "GO_TEST_OPTS=${go_test_opts}"
    --workdir /src
    "$DOCKER_IMAGE"
    /bin/sh -ceu "$inner_script"
  )

  echo "============================== Running Command =============================="
  printf '%q ' timeout "$QEMU_TIMEOUT" "${docker_cmd[@]}"
  echo
  echo "============================================================================="

  local docker_rc=0

  # 파이프라인 실패해도 rc를 기록해야 하므로 set +e로 감싸기
  set +e
  timeout "$QEMU_TIMEOUT" "${docker_cmd[@]}" 2>&1 \
    | linebuf tee >( (command -v ts >/dev/null 2>&1 && linebuf ts "%y-%m-%d %H:%M:%S" || cat) > "$raw_log" ) \
    | linebuf awk '
        /^===== GO TEST =====$/ { print; fflush(); found=1; next }
        !found { print; fflush(); next }
        { print > "/dev/fd/3"; fflush("/dev/fd/3") }
      ' 3> >(
        gotestfmt \
          | sed -u 's/\x1b\[[0-9;]*m//g' \
          | linebuf tee "$out_log"
      )
  docker_rc=${PIPESTATUS[0]}
  set -e

  # 프로세스 치환(ts/gotestfmt 파이프라인) flush/종료까지 기다림
  wait || true

  # 프롬프트가 같은 줄에 붙는 현상 방지
  printf '\n'

  if (( docker_rc != 0 )); then
    echo "Error: Command failed for ${DOCKER_PLATFORM} with ${DOCKER_IMAGE}. Check logs:"
    echo "  raw:  $raw_log"
    echo "  out:  $out_log"
  fi

  local end_time end_time_human elapsed elapsed_human
  end_time="$(date +%s)"
  end_time_human="$(date +"%Y-%m-%d %H:%M:%S")"
  elapsed=$((end_time - start_time))
  elapsed_human="$(date -u -d @"$elapsed" +"%H:%M:%S")"

  # csv에 명령도 남기고 싶다면(콤마/개행 안전하게 하려면 quoting 권장)
  local cmd_str
  cmd_str="$(printf '%q ' timeout "$QEMU_TIMEOUT" "${docker_cmd[@]}")"

  echo "${start_time_human},${end_time_human},${elapsed_human},${DOCKER_PLATFORM},${DOCKER_IMAGE},${variant},${docker_rc},${cmd_str}" >> qemu.csv
  return "$docker_rc"
}

main
