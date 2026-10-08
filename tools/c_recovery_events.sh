#!/usr/bin/env bash
# 十三碼事件 handler 與政治／災害依賴；spec/219。
set -euo pipefail
EVENTS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EVENTS_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
EVENTS_OUT="$EVENTS_ROOT/workplace/matching-decompilation/c-events"
EVENTS_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$EVENTS_GOLEM/oracle";test -f "$EVENTS_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$EVENTS_OUT";test -O "$EVENTS_OUT"
EVENTS_IMAGE_ID="$(docker image inspect "$EVENTS_IMAGE" --format '{{.Id}}')"
EVENTS_CONTAINER="wolong-c-events-$$-$RANDOM"
trap 'docker rm -f "$EVENTS_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1800 docker run --rm --init --name "$EVENTS_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$EVENTS_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$EVENTS_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$EVENTS_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$EVENTS_OUT,dst=/output" \
  --env "EVENTS_MODE=${WOLONG_EVENTS_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$EVENTS_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/events.go.work
    export GOWORK=/tmp/events.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/events.c /repo/tools/c_recovery/events.h /repo/tools/c_recovery/events_fixture.h /repo/tools/c_recovery/hourly.c /repo/tools/c_recovery/hourly.h /repo/tools/c_recovery/hourly_fixture.h /repo/tools/c_recovery/politics.c /repo/tools/c_recovery/politics.h /repo/tools/c_recovery/world_update.c /repo/tools/c_recovery/world_update.h /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_events.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$EVENTS_MODE" in
      smoke)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/events-probe /repo/tools/c_recovery_events.go
        /output/events-probe -smoke -out /output/results/smoke.json
        ;;
      full)
        for optimize in O2 O0; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/events-$optimize" /repo/tools/c_recovery_events.go
          "/output/events-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:presence 2:relation 3:target 4:disaster 5:capital 6:disaster 7:cleanup 8:diplomacy; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_EVENTS_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_events.go
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
if [[ "${WOLONG_EVENTS_MODE:-full}" == full ]]; then
  test -s "$EVENTS_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$EVENTS_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$EVENTS_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_events_verify.py --repo /repo --output /output
fi
