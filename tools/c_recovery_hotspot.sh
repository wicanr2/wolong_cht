#!/usr/bin/env bash
# 熱區／真實 pixel query 與數值接線；spec/222。
set -euo pipefail
HOTSPOT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOTSPOT_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
HOTSPOT_OUT="$HOTSPOT_ROOT/workplace/matching-decompilation/c-hotspot"
HOTSPOT_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$HOTSPOT_GOLEM/oracle";test -f "$HOTSPOT_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$HOTSPOT_OUT";test -O "$HOTSPOT_OUT"
HOTSPOT_IMAGE_ID="$(docker image inspect "$HOTSPOT_IMAGE" --format '{{.Id}}')"
HOTSPOT_MODULES="$HOTSPOT_ROOT/workplace/matching-decompilation/c-economy/state-modcache"
test -d "$HOTSPOT_MODULES"
HOTSPOT_CONTAINER="wolong-c-hotspot-$$-$RANDOM"
trap 'docker rm -f "$HOTSPOT_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1800 docker run --rm --init --name "$HOTSPOT_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$HOTSPOT_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$HOTSPOT_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$HOTSPOT_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$HOTSPOT_OUT,dst=/output" \
  --mount "type=bind,src=$HOTSPOT_MODULES,dst=/modcache,readonly" \
  --env "HOTSPOT_MODE=${WOLONG_HOTSPOT_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$HOTSPOT_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOMODCACHE=/modcache GOPROXY=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/hotspot.go.work
    export GOWORK=/tmp/hotspot.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/hotspot.c /repo/tools/c_recovery/hotspot.h /repo/tools/c_recovery/hotspot_fixture.h /repo/tools/c_recovery/numeric.c /repo/tools/c_recovery/numeric.h /repo/tools/c_recovery/numeric_fixture.h /repo/tools/c_recovery/modal.c /repo/tools/c_recovery/modal.h /repo/tools/c_recovery/modal_fixture.h /repo/tools/c_recovery/events.c /repo/tools/c_recovery/events.h /repo/tools/c_recovery/hourly.c /repo/tools/c_recovery/hourly.h /repo/tools/c_recovery/hourly_fixture.h /repo/tools/c_recovery/politics.c /repo/tools/c_recovery/politics.h /repo/tools/c_recovery/world_update.c /repo/tools/c_recovery/world_update.h /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_hotspot.go > /output/results/c-source.sha256
    find /repo/internal/state /repo/internal/rules -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/go-source.sha256
    HOTSPOT_SOURCE_DIGEST="$(sha256sum /output/results/c-source.sha256)"
    HOTSPOT_SOURCE_DIGEST="${HOTSPOT_SOURCE_DIGEST%% *}"
    HOTSPOT_SOURCE_DEFINE="-DKI_HOTSPOT_SOURCE_DIGEST=0x${HOTSPOT_SOURCE_DIGEST:0:16}"
    printf "%s\n" "$HOTSPOT_SOURCE_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$HOTSPOT_MODE" in
      smoke)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $HOTSPOT_SOURCE_DEFINE" go build -p 2 -o /output/hotspot-probe /repo/tools/c_recovery_hotspot.go
        go version -m /output/hotspot-probe > /output/results/buildinfo-smoke.txt
        /output/hotspot-probe -smoke -out /output/results/smoke.json
        ;;
      full)
        for optimize in O2 O0; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $HOTSPOT_SOURCE_DEFINE" go build -p 2 -o "/output/hotspot-$optimize" /repo/tools/c_recovery_hotspot.go
          go version -m "/output/hotspot-$optimize" > "/output/results/buildinfo-$optimize.txt"
          "/output/hotspot-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:init 2:rectangle 3:rectangle 4:rectangle 5:query-pixels 6:wrapper 7:wrapper 8:wrapper; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $HOTSPOT_SOURCE_DEFINE -DKI_HOTSPOT_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_hotspot.go
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
if [[ "${WOLONG_HOTSPOT_MODE:-full}" == full ]]; then
  test -s "$HOTSPOT_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$HOTSPOT_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$HOTSPOT_OUT,dst=/output" \
  --mount "type=bind,src=$HOTSPOT_MODULES,dst=/modcache,readonly" \
    python:3.13-bookworm python /repo/tools/c_recovery_hotspot_verify.py --repo /repo --output /output
fi
