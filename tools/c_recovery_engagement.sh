#!/usr/bin/env bash
# spec/251：完整交戰與戰術原生C。長工作可用背景容器與持久退出收據。
set -euo pipefail
EG_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EG_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
EG_OUT="$EG_ROOT/workplace/matching-decompilation/c-engagement"
test -d "$EG_GOLEM/internal/machine";test -d "$EG_ROOT/workplace/orig/dosv"
test -f "$EG_ROOT/workplace/orig/dosv/KI.EXE"
test -f "$EG_ROOT/tools/c_recovery_engagement_container.sh"
test -d "$EG_OUT";test -O "$EG_OUT"
test -d "$EG_OUT/results";test -O "$EG_OUT/results"
# 共用來源清單，mutants、normal-o2、normal-o0 須依序執行。
case "${WOLONG_ENGAGEMENT_MODE:-full}" in full|smoke|controls|mutants|normal-o2|normal-o0) ;; *) exit 2 ;; esac
case "${WOLONG_ENGAGEMENT_DETACHED:-0}" in 0|1) ;; *) exit 2 ;; esac
# 外層容器期限依工作量設定，不改原版指令預算或測試矩陣。
EG_TIMEOUT=2400
case "${WOLONG_ENGAGEMENT_MODE:-full}" in
  normal-o2|normal-o0) EG_TIMEOUT=3600 ;;
  full) EG_TIMEOUT=9600 ;;
esac
# 指定分組只供 smoke 窄診斷，其他模式忽略此環境變數。
EG_GROUP=""
if [[ "${WOLONG_ENGAGEMENT_MODE:-full}" == smoke ]]; then
  EG_GROUP="${WOLONG_ENGAGEMENT_GROUP:-}"
  [[ -z "$EG_GROUP" || "$EG_GROUP" =~ ^[a-z-]+$ ]] || exit 2
fi
EG_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
EG_CONTAINER="wolong-c-engagement-${WOLONG_ENGAGEMENT_MODE:-full}-$$-$RANDOM"
EG_JOB_BASE="$EG_OUT/results/$EG_CONTAINER"
trap 'docker rm -f "$EG_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$EG_ROOT,dst=/repo,readonly" --mount "type=bind,src=$EG_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
EG_DOCKER_ARGS=(--rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$EG_CONTAINER"
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)"
  --label wolong.task=matching-decompilation
  --mount "type=bind,src=$EG_ROOT,dst=/repo,readonly"
  --mount "type=bind,src=$EG_GOLEM,dst=/golem,readonly"
  --mount "type=bind,src=$EG_ROOT/workplace/orig/dosv,dst=/orig,readonly"
  --mount "type=bind,src=$EG_OUT,dst=/output"
  --env "EG_MODE=${WOLONG_ENGAGEMENT_MODE:-full}" --env "EG_GROUP=$EG_GROUP"
  --env "EG_TIMEOUT=$EG_TIMEOUT" --env "EG_JOB_ID=$EG_CONTAINER"
  --workdir /repo --entrypoint /bin/bash)
if [[ "${WOLONG_ENGAGEMENT_DETACHED:-0}" == 1 ]]; then
  EG_CONTAINER_ID="$(timeout --kill-after=10 60 docker run -d "${EG_DOCKER_ARGS[@]}" "$EG_IMAGE" /repo/tools/c_recovery_engagement_container.sh)"
  trap - EXIT
  printf 'container_name=%s\ncontainer_id=%s\njob_log=%s.job.log\njob_exit=%s.job.exit\n' "$EG_CONTAINER" "$EG_CONTAINER_ID" "$EG_JOB_BASE" "$EG_JOB_BASE"
  if [[ "${WOLONG_ENGAGEMENT_MODE:-full}" == full || "${WOLONG_ENGAGEMENT_MODE:-full}" == controls ]]; then
    printf 'final_verifier=pending; run tools/c_recovery_engagement_verify.py after successful native job exit\n'
  fi
  exit 0
fi
printf 'job_log=%s.job.log\njob_exit=%s.job.exit\n' "$EG_JOB_BASE" "$EG_JOB_BASE"
docker run "${EG_DOCKER_ARGS[@]}" "$EG_IMAGE" /repo/tools/c_recovery_engagement_container.sh
if [[ "${WOLONG_ENGAGEMENT_MODE:-full}" == full || "${WOLONG_ENGAGEMENT_MODE:-full}" == controls ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$EG_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$EG_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_engagement_verify.py --repo /repo --output /output
fi
