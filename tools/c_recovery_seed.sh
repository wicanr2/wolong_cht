#!/usr/bin/env bash
# RNG 播種的受控 RTC 對照；契約 docs/spec/202-c-rng-seed.md。
set -euo pipefail
SEED_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SEED_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
SEED_OUT="$SEED_ROOT/workplace/matching-decompilation/c-seed"
SEED_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$SEED_GOLEM/oracle"
test -f "$SEED_ROOT/workplace/orig/dosv/KI.EXE"
test -f "$SEED_ROOT/tools/c_recovery/seed.c"
mkdir -p "$SEED_OUT"
test -O "$SEED_OUT"
SEED_IMAGE_ID="$(docker image inspect "$SEED_IMAGE" --format '{{.Id}}')"
SEED_CONTAINER="wolong-c-seed-$$-$RANDOM"
trap 'docker rm -f "$SEED_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 600 docker run --rm --init --name "$SEED_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$SEED_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$SEED_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$SEED_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$SEED_OUT,dst=/output" \
  --workdir /repo --entrypoint /bin/bash "$SEED_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off
    export GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/seed.go.work
    export GOWORK=/tmp/seed.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/seed.c /repo/tools/c_recovery/seed.h /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_seed.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt
    gcc --version >> /output/results/tool-versions.txt
    for optimize in O2 O0; do
      CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/seed-$optimize" /repo/tools/c_recovery_seed.go
      "/output/seed-$optimize" -out "/output/results/$optimize.json"
    done
    for mutation in 1 2 3; do
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_SEED_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_seed.go
      set +e
      "/output/mutant-$mutation" -max 2 -out "/output/results/mutant-$mutation.json"
      status=$?
      set -e
      test "$status" -eq 1
      test -s "/output/results/mutant-$mutation.json"
    done
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
timeout --kill-after=10 45 docker run --rm --init --network none \
  --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$SEED_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$SEED_OUT,dst=/output" \
  --entrypoint python python:3.13-bookworm /repo/tools/c_recovery_seed_verify.py \
  --output /output --repo /repo
