#!/usr/bin/env bash
set -Eeuo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
test_dir="$(mktemp -d)"
trap 'rm -rf -- "$test_dir"' EXIT
mkdir -p "$test_dir/mock"

# Exercise the runner without Docker, network access, or architecture emulation.
cat > "$test_dir/mock/docker" <<'MOCK'
#!/usr/bin/env bash
case "$*" in
  *"go env GOMODCACHE") echo /gomodcache ;;
  *"go env GOCACHE") echo /gocache ;;
  volume*) ;;
  *)
    echo '===== GO TEST ====='
    echo 'mock container output'
    exit "$KRYPTO_TEST_CONTAINER_STATUS"
    ;;
esac
MOCK
cat > "$test_dir/mock/gotestfmt" <<'MOCK'
#!/usr/bin/env bash
cat
MOCK
cat > "$test_dir/mock/timeout" <<'MOCK'
#!/usr/bin/env bash
printf '%s\n' "$1" >> "$KRYPTO_TEST_TIMEOUT_LOG"
if [[ "$KRYPTO_TEST_CONTAINER_STATUS" == 124 ]]; then
  exit 124
fi
shift
exec "$@"
MOCK
chmod +x "$test_dir/mock/docker" "$test_dir/mock/gotestfmt" "$test_dir/mock/timeout"

for timeout_setting in default 3h; do
  for expected in 0 42 124; do
    case_dir="$test_dir/case-$timeout_setting-$expected"
    mkdir -p "$case_dir"
    status=0
    (
      cd "$case_dir"
      export PATH="$test_dir/mock:$PATH"
      export KRYPTO_TEST_CONTAINER_STATUS="$expected"
      export KRYPTO_TEST_TIMEOUT_LOG="$case_dir/timeout.log"
      unset QEMU_TIMEOUT
      if [[ "$timeout_setting" != default ]]; then
        export QEMU_TIMEOUT="$timeout_setting"
      fi
      bash "$script_dir/docker-qemu.sh" linux/arm/v7 krypto/test
    ) > "$case_dir/output.log" 2>&1 || status=$?

    if (( status != expected )); then
      cat "$case_dir/output.log" >&2
      echo "expected exit $expected, got $status" >&2
      exit 1
    fi
    if [[ ! -s "$case_dir/qemu.csv" ]]; then
      echo 'container status was not recorded' >&2
      exit 1
    fi
    expected_runs=1
    if (( expected == 0 )); then
      expected_runs=2
    fi
    if [[ "$(wc -l < "$case_dir/qemu.csv")" -ne "$expected_runs" ]]; then
      echo "expected $expected_runs variants to be recorded" >&2
      exit 1
    fi
    expected_timeout="$timeout_setting"
    if [[ "$timeout_setting" == default ]]; then
      expected_timeout=2h
    fi
    if [[ "$(wc -l < "$case_dir/timeout.log")" -ne "$expected_runs" ]] ||
       [[ "$(sort -u "$case_dir/timeout.log")" != "$expected_timeout" ]]; then
      echo "expected timeout $expected_timeout for each variant" >&2
      exit 1
    fi
    if [[ "$(grep -c "timeout $expected_timeout docker run" "$case_dir/qemu.csv")" -ne "$expected_runs" ]]; then
      echo 'recorded command did not use the configured timeout' >&2
      exit 1
    fi
  done
done

echo 'QEMU runner success and failure propagation tests passed'
