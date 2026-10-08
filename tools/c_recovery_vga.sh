#!/usr/bin/env bash
# Eight original VGA entries with real plane/latch adapter; spec/223.
set -euo pipefail
VGA_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VGA_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
VGA_OUT="$VGA_ROOT/workplace/matching-decompilation/c-vga"
test -d "$VGA_GOLEM/internal/machine";test -f "$VGA_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$VGA_OUT";test -O "$VGA_OUT"
VGA_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
VGA_CONTAINER="wolong-c-vga-$$-$RANDOM"
trap 'docker rm -f "$VGA_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1800 docker run --rm --init --name "$VGA_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$VGA_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$VGA_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$VGA_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$VGA_OUT,dst=/output" --env "VGA_MODE=${WOLONG_VGA_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$VGA_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/vga-build /output/results
    cp /repo/tools/c_recovery_vga.go /tmp/vga-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/vga-build/bus.go
    printf "module github.com/wicanr2/dosgolem/wolongcvga\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/vga-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/vga.c /repo/tools/c_recovery/vga.h /repo/tools/c_recovery/vga_fixture.h /repo/tools/c_recovery_vga.go /repo/tools/c_recovery_vga_bus.go > /output/results/c-source.sha256
    VGA_SOURCE_DIGEST="$(sha256sum /output/results/c-source.sha256)";VGA_SOURCE_DIGEST="${VGA_SOURCE_DIGEST%% *}"
    VGA_SOURCE_DEFINE="-DKI_VGA_SOURCE_DIGEST=0x${VGA_SOURCE_DIGEST:0:16}"
    printf "%s\n" "$VGA_SOURCE_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/vga-build
    if [[ "$VGA_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $VGA_SOURCE_DEFINE" go build -tags matching_vga -p 2 -o /output/vga-probe .
      go version -m /output/vga-probe > /output/results/buildinfo-smoke.txt
      /output/vga-probe -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $VGA_SOURCE_DEFINE" go build -tags matching_vga -p 2 -o "/output/vga-$optimize" .
        go version -m "/output/vga-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/vga-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:leaf 2:draw 3:leaf 4:draw 5:draw 6:save 7:wrapper 8:draw; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $VGA_SOURCE_DEFINE -DKI_VGA_MUTATION=$mutation" go build -tags matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_VGA_MODE:-full}" == full ]]; then
  timeout 45 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$VGA_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$VGA_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_vga_verify.py --repo /repo --output /output
fi
