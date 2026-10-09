#!/usr/bin/env bash
# spec/235：軍團／物件 producer、矩陣裁切與小地圖。
set -euo pipefail
OVERLAY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OVERLAY_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
OVERLAY_OUT="$OVERLAY_ROOT/workplace/matching-decompilation/c-overlay"
test -d "$OVERLAY_GOLEM/internal/machine"
test -f "$OVERLAY_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$OVERLAY_OUT";test -O "$OVERLAY_OUT"
OVERLAY_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
OVERLAY_CONTAINER="wolong-c-overlay-$$-$RANDOM"
trap 'docker rm -f "$OVERLAY_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$OVERLAY_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$OVERLAY_OUT,dst=/output" \
  python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --name "$OVERLAY_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$OVERLAY_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$OVERLAY_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$OVERLAY_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$OVERLAY_OUT,dst=/output" \
  --env "OVERLAY_MODE=${WOLONG_OVERLAY_MODE:-full}" --workdir /repo --entrypoint /bin/bash "$OVERLAY_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/overlay-build /output/results
    cp /repo/tools/c_recovery_overlay.go /tmp/overlay-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/overlay-build/bus.go
    printf "module github.com/wicanr2/dosgolem/wolongcoverlay\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/overlay-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery_overlay.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_overlay_generate.py /repo/tools/c_recovery_display_generate.py >> /output/results/c-source.sha256
    OVERLAY_DIGEST="$(sha256sum /output/results/c-source.sha256)";OVERLAY_DIGEST="${OVERLAY_DIGEST%% *}"
    OVERLAY_DEFINE="-DKI_OVERLAY_SOURCE_DIGEST=0x${OVERLAY_DIGEST:0:16}"
    printf "%s\n" "$OVERLAY_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/overlay-build
    if [[ "$OVERLAY_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $OVERLAY_DEFINE" go build -tags matching_overlay,matching_vga -p 2 -o /output/overlay-smoke .
      /output/overlay-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $OVERLAY_DEFINE" go build -tags matching_overlay,matching_vga -p 2 -o "/output/overlay-$optimize" .
        go version -m "/output/overlay-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/overlay-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:scan 2:point 3:point 4:army 5:objects 6:objects 7:objects 8:matrix 9:matrix 10:matrix 11:matrix 12:minimap; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -D_GNU_SOURCE $OVERLAY_DEFINE -DKI_OVERLAY_MUTATION=$mutation" go build -tags matching_overlay,matching_vga -p 2 -o "/output/mutant-$mutation" .
        go version -m "/output/mutant-$mutation" > "/output/results/buildinfo-mutant-$mutation.txt"
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
if [[ "${WOLONG_OVERLAY_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$OVERLAY_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$OVERLAY_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_overlay_verify.py --repo /repo --output /output
fi
