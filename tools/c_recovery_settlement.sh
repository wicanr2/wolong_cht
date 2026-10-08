#!/usr/bin/env bash
# 五個據點結算函式與 C 月結接線；spec/207。
set -euo pipefail
SETTLE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SETTLE_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
SETTLE_OUT="$SETTLE_ROOT/workplace/matching-decompilation/c-settlement"
SETTLE_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$SETTLE_GOLEM/oracle";test -f "$SETTLE_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$SETTLE_OUT";test -O "$SETTLE_OUT"
SETTLE_IMAGE_ID="$(docker image inspect "$SETTLE_IMAGE" --format '{{.Id}}')"
SETTLE_CONTAINER="wolong-c-settle-$$-$RANDOM"
trap 'docker rm -f "$SETTLE_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1200 docker run --rm --init --name "$SETTLE_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$SETTLE_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$SETTLE_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$SETTLE_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$SETTLE_OUT,dst=/output" \
  --env "SETTLE_MODE=${WOLONG_SETTLEMENT_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$SETTLE_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off
    export GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/settle.go.work
    export GOWORK=/tmp/settle.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/settlement_fixture.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_settlement.go /repo/internal/rules/economy/economy.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$SETTLE_MODE" in
      smoke|go-probe)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/settle-probe /repo/tools/c_recovery_settlement.go
        if [[ "$SETTLE_MODE" == smoke ]]; then
          /output/settle-probe -smoke -skip-go -out /output/results/smoke.json
        else
          /output/settle-probe -go-probe -out /output/results/go-probe.json
        fi
        ;;
      full|finish)
        go test -count=1 -p 2 ./internal/rules/economy
        levels="O2 O0"
        if [[ "$SETTLE_MODE" == finish ]]; then
          test -s /output/results/O2.json
          levels="O0"
        fi
        for optimize in $levels; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/settle-$optimize" /repo/tools/c_recovery_settlement.go
          "/output/settle-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:income 2:recruit 3:player-tax 4:ai-gate 5:city-settlement; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_SETTLEMENT_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_settlement.go
          set +e
          "/output/mutant-$mutation" -group "$group" -skip-go -out "/output/results/mutant-$mutation.json"
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
if [[ "${WOLONG_SETTLEMENT_MODE:-full}" == full || "${WOLONG_SETTLEMENT_MODE:-full}" == finish ]]; then
  test -s "$SETTLE_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$SETTLE_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$SETTLE_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_settlement_verify.py --repo /repo --output /output
fi
