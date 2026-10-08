#!/usr/bin/env bash
# 十一個月結世界更新 C 函式；spec/209。
set -euo pipefail
WORLD_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORLD_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
WORLD_OUT="$WORLD_ROOT/workplace/matching-decompilation/c-world"
WORLD_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$WORLD_GOLEM/oracle";test -f "$WORLD_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$WORLD_OUT";test -O "$WORLD_OUT"
WORLD_IMAGE_ID="$(docker image inspect "$WORLD_IMAGE" --format '{{.Id}}')"
WORLD_CONTAINER="wolong-c-world-$$-$RANDOM"
trap 'docker rm -f "$WORLD_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1200 docker run --rm --init --name "$WORLD_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$WORLD_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$WORLD_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$WORLD_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$WORLD_OUT,dst=/output" \
  --env "WORLD_MODE=${WOLONG_WORLD_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$WORLD_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/world.go.work
    export GOWORK=/tmp/world.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/world_update.c /repo/tools/c_recovery/world_update.h /repo/tools/c_recovery/world_update_fixture.h /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_world.go /repo/internal/rules/economy/economy.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$WORLD_MODE" in
      smoke|go-probe)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/world-probe /repo/tools/c_recovery_world.go
        if [[ "$WORLD_MODE" == smoke ]]; then
          /output/world-probe -smoke -skip-go -out /output/results/smoke.json
        else
          /output/world-probe -go-probe -out /output/results/go-probe.json
        fi
        ;;
      full)
        go test -count=1 -p 2 ./internal/rules/economy
        for optimize in O2 O0; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/world-$optimize" /repo/tools/c_recovery_world.go
          "/output/world-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:growth 2:score 3:writer 4:governor 5:diplomat 6:disaster 7:marker; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_WORLD_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_world.go
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
if [[ "${WOLONG_WORLD_MODE:-full}" == full ]]; then
  test -s "$WORLD_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$WORLD_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$WORLD_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_world_verify.py --repo /repo --output /output
fi
