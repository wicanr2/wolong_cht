#!/usr/bin/env bash
# spec/233：完整 DOS 讀檔、BGM cue 與 allocation。
set -euo pipefail
RESOURCE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RESOURCE_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
RESOURCE_OUT="$RESOURCE_ROOT/workplace/matching-decompilation/c-resource"
test -d "$RESOURCE_GOLEM/internal/machine";test -f "$RESOURCE_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$RESOURCE_OUT";test -O "$RESOURCE_OUT"
RESOURCE_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
RESOURCE_CONTAINER="wolong-c-resource-$$-$RANDOM"
trap 'docker rm -f "$RESOURCE_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$RESOURCE_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$RESOURCE_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$RESOURCE_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$RESOURCE_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$RESOURCE_OUT,dst=/output" --env "RESOURCE_MODE=${WOLONG_RESOURCE_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$RESOURCE_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/resource-build /output/results
    cp /repo/tools/c_recovery_resource.go /tmp/resource-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/resource-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/resource-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/resource-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcresource\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/resource-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/resource_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_resource_generate.py /repo/tools/c_recovery_resource.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    RESOURCE_DIGEST="$(sha256sum /output/results/c-source.sha256)";RESOURCE_DIGEST="${RESOURCE_DIGEST%% *}"
    RESOURCE_DEFINE="-DKI_RESOURCE_SOURCE_DIGEST=0x${RESOURCE_DIGEST:0:16}"
    printf "%s\n" "$RESOURCE_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/resource-build
    if [[ "$RESOURCE_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $RESOURCE_DEFINE" go build -tags matching_resource,matching_input,matching_glyph,matching_vga -p 2 -o /output/resource-smoke .
      /output/resource-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $RESOURCE_DEFINE" go build -tags matching_resource,matching_input,matching_glyph,matching_vga -p 2 -o "/output/resource-$optimize" .
        go version -m "/output/resource-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/resource-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:stream 2:stream 3:stream 4:mmap-caller 5:cue 6:cue 7:cue 8:cue 9:allocation 10:allocation 11:stop 12:stream; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $RESOURCE_DEFINE -DKI_RESOURCE_MUTATION=$mutation" go build -tags matching_resource,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_RESOURCE_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$RESOURCE_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$RESOURCE_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_resource_verify.py --repo /repo --output /output
fi
