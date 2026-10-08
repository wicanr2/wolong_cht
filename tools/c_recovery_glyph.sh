#!/usr/bin/env bash
# spec/225／226：矩形、選取與計量，原版／C／正式 Go 長度。
set -euo pipefail
GLYPH_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GLYPH_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
GLYPH_OUT="$GLYPH_ROOT/workplace/matching-decompilation/c-glyph"
test -d "$GLYPH_GOLEM/internal/machine";test -f "$GLYPH_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$GLYPH_OUT";test -O "$GLYPH_OUT"
GLYPH_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
GLYPH_CONTAINER="wolong-c-glyph-$$-$RANDOM"
trap 'docker rm -f "$GLYPH_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$GLYPH_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$GLYPH_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$GLYPH_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$GLYPH_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$GLYPH_OUT,dst=/output" --env "GLYPH_MODE=${WOLONG_GLYPH_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$GLYPH_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/glyph-build /output/results
    cp /repo/tools/c_recovery_glyph.go /tmp/glyph-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/glyph-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/glyph-build/platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcglyph\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/glyph-build/go.mod
    # Extract the exact production function, including nested closure, without another formula.
    awk "/^func battleSideBarLengths\(/ {capture=1} capture {print} capture && /^}/ {exit}" /repo/cmd/wlgame/battlelayout.go > /output/results/go-helper-source.txt
    test -s /output/results/go-helper-source.txt
    printf "//go:build matching_glyph\npackage main\n" > /tmp/glyph-build/bars.go
    awk "/^[[:space:]]*battleSideBarMaxLen[[:space:]]*=/ {print \"const \" \$0;exit}" /repo/cmd/wlgame/battlelayout.go > /output/results/go-helper-constant.txt
    test -s /output/results/go-helper-constant.txt
    cat /output/results/go-helper-constant.txt >> /tmp/glyph-build/bars.go
    cat /output/results/go-helper-source.txt >> /tmp/glyph-build/bars.go
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/glyph_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_glyph_generate.py /repo/tools/c_recovery_glyph.go /repo/tools/c_recovery_vga_bus.go /repo/cmd/wlgame/battlelayout.go >> /output/results/c-source.sha256
    GLYPH_DIGEST="$(sha256sum /output/results/c-source.sha256)";GLYPH_DIGEST="${GLYPH_DIGEST%% *}"
    GLYPH_DEFINE="-DKI_GLYPH_SOURCE_DIGEST=0x${GLYPH_DIGEST:0:16}"
    printf "%s\n" "$GLYPH_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/glyph-build
    if [[ "$GLYPH_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $GLYPH_DEFINE" go build -tags matching_glyph,matching_vga -p 2 -o /output/glyph-smoke .
      /output/glyph-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $GLYPH_DEFINE" go build -tags matching_glyph,matching_vga -p 2 -o "/output/glyph-$optimize" .
        go version -m "/output/glyph-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/glyph-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:glyph 2:glyph 3:glyph 4:glyph 5:glyph 6:glyph 7:glyph 8:text-name 9:text-name; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $GLYPH_DEFINE -DKI_GLYPH_MUTATION=$mutation" go build -tags matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_GLYPH_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$GLYPH_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$GLYPH_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_glyph_verify.py --repo /repo --output /output
fi
