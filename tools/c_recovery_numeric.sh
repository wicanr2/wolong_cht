#!/usr/bin/env bash
# 數值輸入／財政與 Go 數值核心；spec/221。
set -euo pipefail
NUMERIC_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NUMERIC_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
NUMERIC_OUT="$NUMERIC_ROOT/workplace/matching-decompilation/c-numeric"
NUMERIC_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$NUMERIC_GOLEM/oracle";test -f "$NUMERIC_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$NUMERIC_OUT";test -O "$NUMERIC_OUT"
NUMERIC_IMAGE_ID="$(docker image inspect "$NUMERIC_IMAGE" --format '{{.Id}}')"
NUMERIC_MODULES="$NUMERIC_ROOT/workplace/matching-decompilation/c-economy/state-modcache"
test -d "$NUMERIC_MODULES"
NUMERIC_CONTAINER="wolong-c-numeric-$$-$RANDOM"
trap 'docker rm -f "$NUMERIC_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1800 docker run --rm --init --name "$NUMERIC_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$NUMERIC_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$NUMERIC_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$NUMERIC_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$NUMERIC_OUT,dst=/output" \
  --mount "type=bind,src=$NUMERIC_MODULES,dst=/modcache,readonly" \
  --env "NUMERIC_MODE=${WOLONG_NUMERIC_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$NUMERIC_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOMODCACHE=/modcache GOPROXY=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/numeric.go.work
    export GOWORK=/tmp/numeric.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/numeric.c /repo/tools/c_recovery/numeric.h /repo/tools/c_recovery/numeric_fixture.h /repo/tools/c_recovery/modal.c /repo/tools/c_recovery/modal.h /repo/tools/c_recovery/modal_fixture.h /repo/tools/c_recovery/events.c /repo/tools/c_recovery/events.h /repo/tools/c_recovery/hourly.c /repo/tools/c_recovery/hourly.h /repo/tools/c_recovery/hourly_fixture.h /repo/tools/c_recovery/politics.c /repo/tools/c_recovery/politics.h /repo/tools/c_recovery/world_update.c /repo/tools/c_recovery/world_update.h /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_numeric.go > /output/results/c-source.sha256
    find /repo/internal/state /repo/internal/rules -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/go-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$NUMERIC_MODE" in
      smoke)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/numeric-probe /repo/tools/c_recovery_numeric.go
        /output/numeric-probe -smoke -out /output/results/smoke.json
        ;;
      full)
        for optimize in O2 O0; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/numeric-$optimize" /repo/tools/c_recovery_numeric.go
          "/output/numeric-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:digit-full 2:hundred-full 3:delete-full 4:action-boundary 5:device 6:glyph 7:editor 8:editor 9:finance 10:finance; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_NUMERIC_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_numeric.go
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
if [[ "${WOLONG_NUMERIC_MODE:-full}" == full ]]; then
  test -s "$NUMERIC_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$NUMERIC_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$NUMERIC_OUT,dst=/output" \
  --mount "type=bind,src=$NUMERIC_MODULES,dst=/modcache,readonly" \
    python:3.13-bookworm python /repo/tools/c_recovery_numeric_verify.py --repo /repo --output /output
fi
