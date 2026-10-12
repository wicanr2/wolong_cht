#!/usr/bin/env bash
# spec/253：開機主排程C的有界容器工作；退出狀態留在本機研究目錄。
set -euo pipefail
SU_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SU_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
SU_OUT="$SU_ROOT/workplace/matching-decompilation/c-startup"
test -d "$SU_GOLEM/internal/machine";test -d "$SU_ROOT/workplace/orig/dosv"
test -f "$SU_ROOT/workplace/orig/dosv/KI.EXE"
test -f "$SU_ROOT/tools/c_recovery_startup_container.sh"
test -d "$SU_OUT/results";test -O "$SU_OUT";test -O "$SU_OUT/results"
SU_MODE="${WOLONG_STARTUP_MODE:-full}"
case "$SU_MODE" in full|smoke|mutants|normal-o2|normal-o0|verify) ;; *) exit 2 ;; esac
case "${WOLONG_STARTUP_DETACHED:-0}" in 0|1) ;; *) exit 2 ;; esac
SU_GROUP="${WOLONG_STARTUP_GROUP:-}"
[[ -z "$SU_GROUP" || "$SU_GROUP" =~ ^[a-z-]+$ ]] || exit 2
if [[ "$SU_MODE" == verify ]]; then
 timeout --kill-after=10 120 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$SU_ROOT,dst=/repo,readonly" --mount "type=bind,src=$SU_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$SU_OUT,dst=/output" python:3.13-bookworm \
  python /repo/tools/c_recovery_startup_verify.py --repo /repo --output /output
 exit 0
fi
SU_TIMEOUT=3600
[[ "$SU_MODE" != full ]] || SU_TIMEOUT=9600
SU_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
test "$SU_IMAGE" = sha256:e8c859f5632dcfde7b32d2012b4351728f6437930887c2f6a91ea242459e5514
SU_CONTAINER="wolong-c-startup-$SU_MODE-$$-$RANDOM"
trap 'docker rm -f "$SU_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 \
 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
 --mount "type=bind,src=$SU_ROOT,dst=/repo,readonly" --mount "type=bind,src=$SU_OUT,dst=/output" \
 python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
SU_ARGS=(--rm --init --name "$SU_CONTAINER" --network none --memory 2g --cpus 2 --pids-limit 256
 --log-opt max-size=10m --log-opt max-file=3 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation
 --mount "type=bind,src=$SU_ROOT,dst=/repo,readonly" --mount "type=bind,src=$SU_GOLEM,dst=/golem,readonly"
 --mount "type=bind,src=$SU_ROOT/workplace/orig/dosv,dst=/orig,readonly" --mount "type=bind,src=$SU_OUT,dst=/output"
 --env "SU_MODE=$SU_MODE" --env "SU_GROUP=$SU_GROUP" --env "SU_TIMEOUT=$SU_TIMEOUT" --env "SU_JOB_ID=$SU_CONTAINER" --env "WOLONG_STARTUP_TRAPRETF=${WOLONG_STARTUP_TRAPRETF:-0}"
 --workdir /repo --entrypoint /bin/bash)
if [[ "${WOLONG_STARTUP_DETACHED:-0}" == 1 ]]; then
 SU_CID="$(timeout --kill-after=10 60 docker run -d "${SU_ARGS[@]}" "$SU_IMAGE" /repo/tools/c_recovery_startup_container.sh)"
 trap - EXIT
 printf 'container_name=%s\ncontainer_id=%s\njob_log=%s/results/%s.job.log\njob_exit=%s/results/%s.job.exit\n' "$SU_CONTAINER" "$SU_CID" "$SU_OUT" "$SU_CONTAINER" "$SU_OUT" "$SU_CONTAINER"
else
 docker run "${SU_ARGS[@]}" "$SU_IMAGE" /repo/tools/c_recovery_startup_container.sh
fi
