#!/usr/bin/env bash
# 從版控指令來源冷組譯，只從自備原版匯入非指令區。
set -euo pipefail
CODE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CODE_IMAGE="sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e"
CODE_CONTAINER="wolong-code-record-$$-$RANDOM"
test -d "$CODE_ROOT/workplace/matching-decompilation/assembly"
test -f "$CODE_ROOT/workplace/orig/dosv/KI.EXE"
test -f "$CODE_ROOT/tools/c_recovery/KI.code.S"
docker image inspect "$CODE_IMAGE" >/dev/null
trap 'docker rm -f "$CODE_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 180 docker run --rm --init --name "$CODE_CONTAINER" \
  --network none --memory 1g --cpus 2 --pids-limit 128 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$CODE_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$CODE_ROOT/workplace/matching-decompilation/assembly,dst=/output" \
  --entrypoint python3 "$CODE_IMAGE" /repo/tools/matching_code_record.py build \
  --repo /repo --output /output/code-record --image-id "$CODE_IMAGE"
