#!/usr/bin/env bash
# 原版局部取數 → C 還原 → Go rng.Next；介面契約見 docs/spec/201。
set -euo pipefail
RNG_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RNG_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
RNG_OUT="$RNG_ROOT/workplace/matching-decompilation/c-rng"
RNG_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$RNG_GOLEM/oracle"
test -d "$RNG_ROOT/workplace/orig/dosv"
test -f "$RNG_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$RNG_OUT"
test -O "$RNG_OUT"
RNG_IMAGE_ID="$(docker image inspect "$RNG_IMAGE" --format '{{.Id}}')"
RNG_CONTAINER="wolong-c-rng-$$-$RANDOM"
trap 'docker rm -f "$RNG_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 300 docker run --rm --init --name "$RNG_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$RNG_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$RNG_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$RNG_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$RNG_OUT,dst=/output" \
  -e "WOLONG_C_RNG_IMAGE_ID=$RNG_IMAGE_ID" \
  --workdir /repo --entrypoint /bin/bash "$RNG_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off
    export GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/rng.go.work
    export GOWORK=/tmp/rng.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    git -C /golem status --short > /output/results/golem-worktree.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_rng.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt
    gcc --version >> /output/results/tool-versions.txt
    for optimize in O2 O0; do
      CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 \
        -o "/output/rng-$optimize" /repo/tools/c_recovery_rng.go
      "/output/rng-$optimize" -out "/output/results/$optimize.json"
    done
    for mutation in 1 2 3 4; do
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_MUTATION=$mutation" \
        go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_rng.go
      set +e
      "/output/mutant-$mutation" -max 64 -out "/output/results/mutant-$mutation.json"
      status=$?
      set -e
      test "$status" -eq 1
      test -s "/output/results/mutant-$mutation.json"
    done
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
test -f "$RNG_ROOT/tools/c_recovery_verify.py"
timeout --kill-after=10 45 docker run --rm --init --network none \
  --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$RNG_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$RNG_OUT,dst=/output" \
  --entrypoint python python:3.13-bookworm /repo/tools/c_recovery_verify.py \
  --output /output --repo /repo
