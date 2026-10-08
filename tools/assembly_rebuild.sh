#!/usr/bin/env bash
# 組語基準重建；C 還原以 assembly/build/function-ledger.json 為定位入口。
set -euo pipefail
ASM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ASM_OUTPUT="$ASM_ROOT/workplace/matching-decompilation/assembly"
ASM_IMAGE="${WOLONG_MATCH_BUILD_IMAGE:-sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e}"
if [[ -z "${WOLONG_IDA_PY_IMAGE:-}" ]]; then
  if docker image inspect ida-pro-9.4-idapython:locked-v1 >/dev/null 2>&1; then
    export WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:locked-v1
  else
    export WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1
  fi
fi
ASM_IMAGE_ID="$(docker image inspect "$ASM_IMAGE" --format '{{.Id}}')"
ASM_CONTAINER="wolong-assembly-build-$$-$RANDOM"
trap 'docker rm -f "$ASM_CONTAINER" >/dev/null 2>&1 || true' EXIT
"$ASM_ROOT/tools/ida.sh" probe dosv tools/ida_assembly_inventory.py "$ASM_OUTPUT/dosv"
test -d "$ASM_OUTPUT"
timeout --kill-after=10 180 docker run --rm --init --name "$ASM_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 128 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$ASM_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$ASM_OUTPUT,dst=/output" \
  --entrypoint python3 "$ASM_IMAGE_ID" /repo/tools/assembly_rebuild.py \
  --original /repo/workplace/orig/dosv/KI.EXE --inventory /output/dosv/ida-probe.json \
  --output /output/build --image-id "$ASM_IMAGE_ID"
