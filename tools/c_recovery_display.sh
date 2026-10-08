#!/usr/bin/env bash
# spec/225／226：矩形、選取與計量，原版／C／正式 Go 長度。
set -euo pipefail
DISPLAY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DISPLAY_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
DISPLAY_OUT="$DISPLAY_ROOT/workplace/matching-decompilation/c-display"
test -d "$DISPLAY_GOLEM/internal/machine";test -f "$DISPLAY_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$DISPLAY_OUT";test -O "$DISPLAY_OUT"
DISPLAY_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
DISPLAY_CONTAINER="wolong-c-display-$$-$RANDOM"
trap 'docker rm -f "$DISPLAY_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$DISPLAY_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$DISPLAY_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$DISPLAY_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$DISPLAY_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$DISPLAY_OUT,dst=/output" --env "DISPLAY_MODE=${WOLONG_DISPLAY_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$DISPLAY_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/display-build /output/results
    cp /repo/tools/c_recovery_display.go /tmp/display-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/display-build/bus.go
    printf "module github.com/wicanr2/dosgolem/wolongcdisplay\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/display-build/go.mod
    # Extract the exact production function, including nested closure, without another formula.
    awk "/^func battleSideBarLengths\(/ {capture=1} capture {print} capture && /^}/ {exit}" /repo/cmd/wlgame/battlelayout.go > /output/results/go-helper-source.txt
    test -s /output/results/go-helper-source.txt
    printf "//go:build matching_display\npackage main\n" > /tmp/display-build/bars.go
    awk "/^[[:space:]]*battleSideBarMaxLen[[:space:]]*=/ {print \"const \" \$0;exit}" /repo/cmd/wlgame/battlelayout.go > /output/results/go-helper-constant.txt
    test -s /output/results/go-helper-constant.txt
    cat /output/results/go-helper-constant.txt >> /tmp/display-build/bars.go
    cat /output/results/go-helper-source.txt >> /tmp/display-build/bars.go
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/display_generated.inc /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_display.go /repo/tools/c_recovery_vga_bus.go /repo/cmd/wlgame/battlelayout.go >> /output/results/c-source.sha256
    DISPLAY_DIGEST="$(sha256sum /output/results/c-source.sha256)";DISPLAY_DIGEST="${DISPLAY_DIGEST%% *}"
    DISPLAY_DEFINE="-DKI_DISPLAY_SOURCE_DIGEST=0x${DISPLAY_DIGEST:0:16}"
    printf "%s\n" "$DISPLAY_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/display-build
    if [[ "$DISPLAY_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $DISPLAY_DEFINE" go build -tags matching_display,matching_vga -p 2 -o /output/display-smoke .
      /output/display-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $DISPLAY_DEFINE" go build -tags matching_display,matching_vga -p 2 -o "/output/display-$optimize" .
        go version -m "/output/display-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/display-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:scene 2:opcode 3:table-swap 4:opcode 5:text 6:text 7:texture 8:axis 9:line-box 10:register; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $DISPLAY_DEFINE -DKI_DISPLAY_MUTATION=$mutation" go build -tags matching_display,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_DISPLAY_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$DISPLAY_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$DISPLAY_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_display_verify.py --repo /repo --output /output
fi
