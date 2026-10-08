#!/usr/bin/env bash
# C 時鐘、callee 快照與等待輸入；spec/203、spec/204。
set -euo pipefail
CLOCK_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CLOCK_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
CLOCK_OUT="$CLOCK_ROOT/workplace/matching-decompilation/c-clock"
CLOCK_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$CLOCK_GOLEM/oracle"
test -f "$CLOCK_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$CLOCK_OUT";test -O "$CLOCK_OUT"
CLOCK_IMAGE_ID="$(docker image inspect "$CLOCK_IMAGE" --format '{{.Id}}')"
CLOCK_CONTAINER="wolong-c-clock-$$-$RANDOM"
trap 'docker rm -f "$CLOCK_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 600 docker run --rm --init --name "$CLOCK_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$CLOCK_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$CLOCK_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$CLOCK_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$CLOCK_OUT,dst=/output" \
  --workdir /repo --entrypoint /bin/bash "$CLOCK_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off
    export GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/clock.go.work
    export GOWORK=/tmp/clock.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/clock.c /repo/tools/c_recovery/clock.h /repo/tools/c_recovery/clock_fixture.h /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_clock.go /repo/internal/rules/clock/clock.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    go test -count=1 -p 2 ./internal/rules/clock
    for optimize in O2 O0; do
      CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/clock-$optimize" /repo/tools/c_recovery_clock.go
      "/output/clock-$optimize" -out "/output/results/$optimize.json"
    done
    for mutation in 1 2 3 4; do
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_CLOCK_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_clock.go
      set +e
      "/output/mutant-$mutation" -max 2 -out "/output/results/mutant-$mutation.json"
      status=$?
      set -e
      test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
    done
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$CLOCK_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$CLOCK_OUT,dst=/output" \
  python:3.13-bookworm python /repo/tools/c_recovery_clock_verify.py --repo /repo --output /output
