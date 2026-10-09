#!/usr/bin/env bash
# spec/240：玩家軍團編成介面。
set -euo pipefail
FORMATION_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FORMATION_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
FORMATION_OUT="$FORMATION_ROOT/workplace/matching-decompilation/c-formation"
test -d "$FORMATION_GOLEM/internal/machine";test -f "$FORMATION_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$FORMATION_OUT";test -O "$FORMATION_OUT"
FORMATION_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
FORMATION_CONTAINER="wolong-c-formation-$$-$RANDOM"
trap 'docker rm -f "$FORMATION_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$FORMATION_ROOT,dst=/repo,readonly" --mount "type=bind,src=$FORMATION_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$FORMATION_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$FORMATION_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$FORMATION_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$FORMATION_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$FORMATION_OUT,dst=/output" --env "FORMATION_MODE=${WOLONG_FORMATION_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$FORMATION_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/formation-build /output/results
    cp /repo/tools/c_recovery_formation.go /tmp/formation-build/main.go
    cp /repo/tools/c_recovery_formation_data.go /tmp/formation-build/data.go
    cp /repo/tools/c_recovery_formation_platform.go /tmp/formation-build/formation-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/formation-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/formation-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/formation-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcformation\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/formation-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/formation_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_formation_generate.py /repo/tools/c_recovery_formation.go /repo/tools/c_recovery_formation_data.go /repo/tools/c_recovery_formation_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    FORMATION_DIGEST="$(sha256sum /output/results/c-source.sha256)";FORMATION_DIGEST="${FORMATION_DIGEST%% *}"
    FORMATION_DEFINE="-DKI_FORMATION_SOURCE_DIGEST=0x${FORMATION_DIGEST:0:16}"
    printf "%s\n" "$FORMATION_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/formation-build
    if [[ "$FORMATION_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $FORMATION_DEFINE" go build -tags matching_formation,matching_input,matching_glyph,matching_vga -p 2 -o /output/formation-smoke .
      /output/formation-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $FORMATION_DEFINE" go build -tags matching_formation,matching_input,matching_glyph,matching_vga -p 2 -o "/output/formation-$optimize" .
        go version -m "/output/formation-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/formation-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:init 2:numbers 3:icons 4:frame 5:formation 6:switch 7:formation 8:formation 9:caller 10:formation; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -D_GNU_SOURCE $FORMATION_DEFINE -DKI_FORMATION_MUTATION=$mutation" go build -tags matching_formation,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_FORMATION_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$FORMATION_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$FORMATION_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_formation_verify.py --repo /repo --output /output
fi
