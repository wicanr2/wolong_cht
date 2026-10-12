#!/usr/bin/env bash
# spec/254：開機主排程C的有界容器工作；退出狀態留在本機研究目錄。
set -euo pipefail
RS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RS_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
RS_OUT="$RS_ROOT/workplace/matching-decompilation/c-residual"
test -d "$RS_GOLEM/internal/machine";test -d "$RS_ROOT/workplace/orig/dosv"
test -f "$RS_ROOT/workplace/orig/dosv/KI.EXE"
test -f "$RS_ROOT/tools/c_recovery_residual_container.sh"
test -d "$RS_OUT/results";test -O "$RS_OUT";test -O "$RS_OUT/results"
RS_MODE="${WOLONG_RESIDUAL_MODE:-full}"
case "$RS_MODE" in full|smoke|mutants|normal-o2|normal-o0|verify) ;; *) exit 2 ;; esac
case "${WOLONG_RESIDUAL_DETACHED:-0}" in 0|1) ;; *) exit 2 ;; esac
RS_GROUP="${WOLONG_RESIDUAL_GROUP:-}"
[[ -z "$RS_GROUP" || "$RS_GROUP" =~ ^[a-z-]+$ ]] || exit 2
if [[ "$RS_MODE" == verify ]]; then
 timeout --kill-after=10 120 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$RS_ROOT,dst=/repo,readonly" --mount "type=bind,src=$RS_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$RS_OUT,dst=/output" python:3.13-bookworm \
  python /repo/tools/c_recovery_residual_verify.py --repo /repo --output /output
 exit 0
fi
RS_TIMEOUT=3600
[[ "$RS_MODE" != full ]] || RS_TIMEOUT=9600
RS_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
test "$RS_IMAGE" = sha256:e8c859f5632dcfde7b32d2012b4351728f6437930887c2f6a91ea242459e5514
RS_CONTAINER="wolong-c-residual-$RS_MODE-$$-$RANDOM"
trap 'docker rm -f "$RS_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 \
 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
 --mount "type=bind,src=$RS_ROOT,dst=/repo,readonly" --mount "type=bind,src=$RS_OUT,dst=/output" \
 python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
RS_ARGS=(--rm --init --name "$RS_CONTAINER" --network none --memory 2g --cpus 2 --pids-limit 256
 --log-opt max-size=10m --log-opt max-file=3 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation
 --mount "type=bind,src=$RS_ROOT,dst=/repo,readonly" --mount "type=bind,src=$RS_GOLEM,dst=/golem,readonly"
 --mount "type=bind,src=$RS_ROOT/workplace/orig/dosv,dst=/orig,readonly" --mount "type=bind,src=$RS_OUT,dst=/output"
 --env "RS_MODE=$RS_MODE" --env "RS_GROUP=$RS_GROUP" --env "RS_TIMEOUT=$RS_TIMEOUT" --env "RS_JOB_ID=$RS_CONTAINER" --env "WOLONG_RESIDUAL_TRAPRETF=${WOLONG_RESIDUAL_TRAPRETF:-0}"
 --workdir /repo --entrypoint /bin/bash)
if [[ "${WOLONG_RESIDUAL_DETACHED:-0}" == 1 ]]; then
 RS_CID="$(timeout --kill-after=10 60 docker run -d "${RS_ARGS[@]}" "$RS_IMAGE" /repo/tools/c_recovery_residual_container.sh)"
 trap - EXIT
 printf 'container_name=%s\ncontainer_id=%s\njob_log=%s/results/%s.job.log\njob_exit=%s/results/%s.job.exit\n' "$RS_CONTAINER" "$RS_CID" "$RS_OUT" "$RS_CONTAINER" "$RS_OUT" "$RS_CONTAINER"
else
 docker run "${RS_ARGS[@]}" "$RS_IMAGE" /repo/tools/c_recovery_residual_container.sh
fi
