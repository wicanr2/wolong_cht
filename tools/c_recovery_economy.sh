#!/usr/bin/env bash
# 月結主流程與五個經濟 C 函式；spec/205。
set -euo pipefail
ECONOMY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ECONOMY_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
ECONOMY_OUT="$ECONOMY_ROOT/workplace/matching-decompilation/c-economy"
ECONOMY_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$ECONOMY_GOLEM/oracle";test -f "$ECONOMY_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$ECONOMY_OUT";test -O "$ECONOMY_OUT"
ECONOMY_IMAGE_ID="$(docker image inspect "$ECONOMY_IMAGE" --format '{{.Id}}')"
ECONOMY_CONTAINER="wolong-c-economy-$$-$RANDOM"
trap 'docker rm -f "$ECONOMY_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 600 docker run --rm --init --name "$ECONOMY_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$ECONOMY_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$ECONOMY_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$ECONOMY_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$ECONOMY_OUT,dst=/output" \
  --env "ECONOMY_MODE=${WOLONG_ECONOMY_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$ECONOMY_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off
    export GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/economy.go.work
    export GOWORK=/tmp/economy.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/economy_fixture.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_economy.go /repo/internal/rules/economy/economy.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$ECONOMY_MODE" in
      smoke)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/economy-smoke /repo/tools/c_recovery_economy.go
        /output/economy-smoke -smoke -skip-go -out /output/results/smoke.json
        ;;
      go-probe)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/economy-probe /repo/tools/c_recovery_economy.go
        /output/economy-probe -group deficit -out /output/results/go-probe.json
        ;;
      full)
        go test -count=1 -p 2 ./internal/rules/economy
        for optimize in O2 O0; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/economy-$optimize" /repo/tools/c_recovery_economy.go
          "/output/economy-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:credit 2:debit 3:reserve 4:distance 5:deficit 6:monthly; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_ECONOMY_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_economy.go
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
if [[ "${WOLONG_ECONOMY_MODE:-full}" == full ]]; then
  test -s "$ECONOMY_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$ECONOMY_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$ECONOMY_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_economy_verify.py --repo /repo --output /output
fi
