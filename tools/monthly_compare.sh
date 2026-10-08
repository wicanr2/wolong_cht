#!/usr/bin/env bash
# 原版/C/Go 月結 audit，spec/212；結果未通過也完整保留差異。
set -euo pipefail
COMPARE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPARE_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
COMPARE_OUT="$COMPARE_ROOT/workplace/matching-decompilation/monthly-compare"
COMPARE_MODULES="$COMPARE_ROOT/workplace/matching-decompilation/c-economy/state-modcache"
test -d "$COMPARE_GOLEM/oracle";test -f "$COMPARE_ROOT/workplace/orig/dosv/KI.EXE";test -d "$COMPARE_MODULES"
mkdir -p "$COMPARE_OUT";test -O "$COMPARE_OUT"
COMPARE_IMAGE_ID="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
COMPARE_CONTAINER="wolong-monthly-compare-$$-$RANDOM"
trap 'docker rm -f "$COMPARE_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 240 docker run --rm --init --name "$COMPARE_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$COMPARE_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$COMPARE_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$COMPARE_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$COMPARE_MODULES,dst=/modcache,readonly" \
  --mount "type=bind,src=$COMPARE_OUT,dst=/output" \
  --workdir /repo --entrypoint /bin/bash "$COMPARE_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOMODCACHE=/modcache GOPROXY=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/monthly.go.work
    export GOWORK=/tmp/monthly.go.work
    git -C /golem rev-parse HEAD > /output/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/golem-source-before.sha256
    find /repo/internal/state /repo/internal/rules -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/go-source.sha256
    sha256sum /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/world_update.c /repo/tools/c_recovery/politics.c /repo/tools/c_recovery/politics_fixture.h /repo/tools/monthly_compare.go > /output/c-source.sha256
    go version > /output/tool-versions.txt;gcc --version >> /output/tool-versions.txt
    CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -tags matching -mod=readonly -p 2 -o /output/monthly-compare /repo/tools/monthly_compare.go
    /output/monthly-compare -out /output/audit.json -vectors /output/vectors
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/golem-source-after.sha256
    cmp /output/golem-source-before.sha256 /output/golem-source-after.sha256
  '
timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$COMPARE_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$COMPARE_OUT,dst=/output" \
  python:3.13-bookworm python /repo/tools/monthly_compare_verify.py --repo /repo --output /output
