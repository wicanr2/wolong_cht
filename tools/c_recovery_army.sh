#!/usr/bin/env bash
# spec/252：軍團輪轉C的有界容器工作；退出狀態留在本機研究目錄。
set -euo pipefail
AR_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AR_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
AR_OUT="$AR_ROOT/workplace/matching-decompilation/c-army"
test -d "$AR_GOLEM/internal/machine";test -d "$AR_ROOT/workplace/orig/dosv"
test -f "$AR_ROOT/workplace/orig/dosv/KI.EXE"
test -f "$AR_ROOT/tools/c_recovery_army_container.sh"
test -d "$AR_OUT/results";test -O "$AR_OUT";test -O "$AR_OUT/results"
AR_MODE="${WOLONG_ARMY_MODE:-full}"
case "$AR_MODE" in full|smoke|mutants|normal-o2|normal-o0) ;; *) exit 2 ;; esac
case "${WOLONG_ARMY_DETACHED:-0}" in 0|1) ;; *) exit 2 ;; esac
AR_GROUP="${WOLONG_ARMY_GROUP:-}"
[[ -z "$AR_GROUP" || "$AR_GROUP" =~ ^[a-z-]+$ ]] || exit 2
AR_TIMEOUT=3600
[[ "$AR_MODE" != full ]] || AR_TIMEOUT=9600
AR_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
test "$AR_IMAGE" = sha256:e8c859f5632dcfde7b32d2012b4351728f6437930887c2f6a91ea242459e5514
AR_CONTAINER="wolong-c-army-$AR_MODE-$$-$RANDOM"
trap 'docker rm -f "$AR_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 \
 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
 --mount "type=bind,src=$AR_ROOT,dst=/repo,readonly" --mount "type=bind,src=$AR_OUT,dst=/output" \
 python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
AR_ARGS=(--rm --init --name "$AR_CONTAINER" --network none --memory 2g --cpus 2 --pids-limit 256
 --log-opt max-size=10m --log-opt max-file=3 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation
 --mount "type=bind,src=$AR_ROOT,dst=/repo,readonly" --mount "type=bind,src=$AR_GOLEM,dst=/golem,readonly"
 --mount "type=bind,src=$AR_ROOT/workplace/orig/dosv,dst=/orig,readonly" --mount "type=bind,src=$AR_OUT,dst=/output"
 --env "AR_MODE=$AR_MODE" --env "AR_GROUP=$AR_GROUP" --env "AR_TIMEOUT=$AR_TIMEOUT" --env "AR_JOB_ID=$AR_CONTAINER"
 --workdir /repo --entrypoint /bin/bash)
if [[ "${WOLONG_ARMY_DETACHED:-0}" == 1 ]]; then
 AR_CID="$(timeout --kill-after=10 60 docker run -d "${AR_ARGS[@]}" "$AR_IMAGE" /repo/tools/c_recovery_army_container.sh)"
 trap - EXIT
 printf 'container_name=%s\ncontainer_id=%s\njob_log=%s/results/%s.job.log\njob_exit=%s/results/%s.job.exit\n' "$AR_CONTAINER" "$AR_CID" "$AR_OUT" "$AR_CONTAINER" "$AR_OUT" "$AR_CONTAINER"
else
 docker run "${AR_ARGS[@]}" "$AR_IMAGE" /repo/tools/c_recovery_army_container.sh
fi
