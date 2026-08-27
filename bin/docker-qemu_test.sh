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
chmod +x "$test_dir/mock/docker" "$test_dir/mock/gotestfmt"

for expected in 0 42; do
  case_dir="$test_dir/case-$expected"
  mkdir -p "$case_dir"
  status=0
  (
    cd "$case_dir"
    export PATH="$test_dir/mock:$PATH"
    export KRYPTO_TEST_CONTAINER_STATUS="$expected"
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
  if (( expected == 0 )) && [[ "$(wc -l < "$case_dir/qemu.csv")" -ne 2 ]]; then
    echo 'successful run did not test both variants' >&2
    exit 1
  fi
done

echo 'QEMU runner success and failure propagation tests passed'
