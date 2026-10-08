#!/usr/bin/env bash
# 局部 matching decompilation 試點，所有工作負載均在 Docker 內。
# 文件入口：docs/re/89-matching-decompilation-pilot.md
set -euo pipefail
MATCH_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MATCH_OUTPUT="$MATCH_ROOT/workplace/matching-decompilation"
MATCH_BUILD_IMAGE="${WOLONG_MATCH_BUILD_IMAGE:-sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e}"

if [[ -z "${WOLONG_IDA_PY_IMAGE:-}" ]]; then
  if docker image inspect ida-pro-9.4-idapython:locked-v1 >/dev/null 2>&1; then
    export WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:locked-v1
  else
    export WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1
  fi
fi
docker image inspect "$WOLONG_IDA_PY_IMAGE" >/dev/null
MATCH_BUILD_ID="$(docker image inspect "$MATCH_BUILD_IMAGE" --format '{{.Id}}')"
MATCH_CONTAINER="wolong-matching-build-$$-$RANDOM"
trap 'docker rm -f "$MATCH_CONTAINER" >/dev/null 2>&1 || true' EXIT
test -d "$MATCH_ROOT/tools"
test -f "$MATCH_ROOT/workplace/orig/dosv/KI.EXE"
test -f "$MATCH_ROOT/workplace/orig/dosv/LOGO.EXE"

"$MATCH_ROOT/tools/ida.sh" probe dosv tools/ida_matching_probe.py "$MATCH_OUTPUT/dosv"
"$MATCH_ROOT/tools/ida.sh" probe dosv tools/ida_matching_probe.py "$MATCH_OUTPUT/borland-control" LOGO.EXE
timeout --kill-after=10 120 docker run --rm --init --network none \
  --name "$MATCH_CONTAINER" \
  --memory 512m --cpus 1 --pids-limit 96 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$MATCH_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$MATCH_OUTPUT,dst=/output" \
  --entrypoint python3 "$MATCH_BUILD_ID" /repo/tools/matching_decomp.py \
  --original /repo/workplace/orig/dosv/KI.EXE \
  --probe /output/dosv/ida-probe.json --output /output/build --image-id "$MATCH_BUILD_ID"
