#!/usr/bin/env bash
# spec/224：位元對齊貼圖、真實按鈕與外框接線。
set -euo pipefail
ALIGNED_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ALIGNED_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
ALIGNED_OUT="$ALIGNED_ROOT/workplace/matching-decompilation/c-aligned"
test -d "$ALIGNED_GOLEM/internal/machine";test -f "$ALIGNED_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$ALIGNED_OUT";test -O "$ALIGNED_OUT"
ALIGNED_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
ALIGNED_CONTAINER="wolong-c-aligned-$$-$RANDOM"
trap 'docker rm -f "$ALIGNED_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1800 docker run --rm --init --name "$ALIGNED_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$ALIGNED_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$ALIGNED_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$ALIGNED_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$ALIGNED_OUT,dst=/output" --env "ALIGNED_MODE=${WOLONG_ALIGNED_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$ALIGNED_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/aligned-build /output/results
    cp /repo/tools/c_recovery_aligned.go /tmp/aligned-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/aligned-build/bus.go
    printf "module github.com/wicanr2/dosgolem/wolongcaligned\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/aligned-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery_aligned.go /repo/tools/c_recovery_vga_bus.go >> /output/results/c-source.sha256
    ALIGNED_DIGEST="$(sha256sum /output/results/c-source.sha256)";ALIGNED_DIGEST="${ALIGNED_DIGEST%% *}"
    ALIGNED_DEFINE="-DKI_ALIGNED_SOURCE_DIGEST=0x${ALIGNED_DIGEST:0:16}"
    printf "%s\n" "$ALIGNED_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/aligned-build
    if [[ "$ALIGNED_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $ALIGNED_DEFINE" go build -tags matching_aligned,matching_vga -p 2 -o /output/aligned-smoke .
      /output/aligned-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $ALIGNED_DEFINE" go build -tags matching_aligned,matching_vga -p 2 -o "/output/aligned-$optimize" .
        go version -m "/output/aligned-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/aligned-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:aligned 2:aligned 3:row 4:aligned 5:row 6:aligned 7:row 8:buttons 9:redraw 10:aligned 11:window; do
        mutation="${entry%%:*}";group="${entry#*:}";define="-DKI_ALIGNED_MUTATION=$mutation"
        if [[ "$mutation" == 11 ]]; then define="-DKI_HOTSPOT_MUTATION=7";fi
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $ALIGNED_DEFINE $define" go build -tags matching_aligned,matching_vga -p 2 -o "/output/mutant-$mutation" .
        set +e
        "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json"
        status=$?
        set -e
        test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
      done
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
if [[ "${WOLONG_ALIGNED_MODE:-full}" == full ]]; then
  timeout 45 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$ALIGNED_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$ALIGNED_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_aligned_verify.py --repo /repo --output /output
fi
