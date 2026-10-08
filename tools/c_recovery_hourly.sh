#!/usr/bin/env bash
# 每時更新及派發核心；spec/218。
set -euo pipefail
HOURLY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOURLY_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
HOURLY_OUT="$HOURLY_ROOT/workplace/matching-decompilation/c-hourly"
HOURLY_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$HOURLY_GOLEM/oracle";test -f "$HOURLY_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$HOURLY_OUT";test -O "$HOURLY_OUT"
HOURLY_IMAGE_ID="$(docker image inspect "$HOURLY_IMAGE" --format '{{.Id}}')"
HOURLY_CONTAINER="wolong-c-hourly-$$-$RANDOM"
trap 'docker rm -f "$HOURLY_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1200 docker run --rm --init --name "$HOURLY_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$HOURLY_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$HOURLY_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$HOURLY_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$HOURLY_OUT,dst=/output" \
  --env "HOURLY_MODE=${WOLONG_HOURLY_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$HOURLY_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/hourly.go.work
    export GOWORK=/tmp/hourly.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/politics.c /repo/tools/c_recovery/politics.h /repo/tools/c_recovery/hourly.c /repo/tools/c_recovery/hourly.h /repo/tools/c_recovery/hourly_fixture.h /repo/tools/c_recovery/world_update.c /repo/tools/c_recovery/world_update.h /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_hourly.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$HOURLY_MODE" in
      smoke)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/hourly-probe /repo/tools/c_recovery_hourly.go
        /output/hourly-probe -smoke -out /output/results/smoke.json
        ;;
      full)
        for optimize in O2 O0; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/hourly-$optimize" /repo/tools/c_recovery_hourly.go
          "/output/hourly-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:expense 2:maintenance 3:dispatch 4:dispatch 5:trust 6:world 7:diplomat; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_HOURLY_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_hourly.go
          set +e
          "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json"
          status=$?
          set -e
          test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
        done
        ;;
      *) exit 2 ;;
    esac
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
if [[ "${WOLONG_HOURLY_MODE:-full}" == full ]]; then
  test -s "$HOURLY_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$HOURLY_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$HOURLY_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_hourly_verify.py --repo /repo --output /output
fi
