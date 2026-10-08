#!/usr/bin/env bash
# spec/229：原始數字 raster 與兩個 code caller。
set -euo pipefail
NUMBERS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NUMBERS_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
NUMBERS_OUT="$NUMBERS_ROOT/workplace/matching-decompilation/c-numbers"
test -d "$NUMBERS_GOLEM/internal/machine";test -f "$NUMBERS_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$NUMBERS_OUT";test -O "$NUMBERS_OUT"
NUMBERS_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
NUMBERS_CONTAINER="wolong-c-numbers-$$-$RANDOM"
trap 'docker rm -f "$NUMBERS_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$NUMBERS_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$NUMBERS_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$NUMBERS_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$NUMBERS_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$NUMBERS_OUT,dst=/output" --env "NUMBERS_MODE=${WOLONG_NUMBERS_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$NUMBERS_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/numbers-build /output/results
    cp /repo/tools/c_recovery_numbers.go /tmp/numbers-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/numbers-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/numbers-build/platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcnumbers\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/numbers-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/numbers_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_numbers_generate.py /repo/tools/c_recovery_numbers.go /repo/tools/c_recovery_vga_bus.go >> /output/results/c-source.sha256
    NUMBERS_DIGEST="$(sha256sum /output/results/c-source.sha256)";NUMBERS_DIGEST="${NUMBERS_DIGEST%% *}"
    NUMBERS_DEFINE="-DKI_NUMBERS_SOURCE_DIGEST=0x${NUMBERS_DIGEST:0:16}"
    printf "%s\n" "$NUMBERS_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/numbers-build
    if [[ "$NUMBERS_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $NUMBERS_DEFINE" go build -tags matching_numbers,matching_glyph,matching_vga -p 2 -o /output/numbers-smoke .
      /output/numbers-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $NUMBERS_DEFINE" go build -tags matching_numbers,matching_glyph,matching_vga -p 2 -o "/output/numbers-$optimize" .
        go version -m "/output/numbers-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/numbers-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:number 2:number 3:number 4:digit 5:digit 6:primitive-boundary 7:digit 8:colors; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $NUMBERS_DEFINE -DKI_NUMBERS_MUTATION=$mutation" go build -tags matching_numbers,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_NUMBERS_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$NUMBERS_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$NUMBERS_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_numbers_verify.py --repo /repo --output /output
fi
