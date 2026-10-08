#!/usr/bin/env bash
# spec/225／226：矩形、選取與計量，原版／C／正式 Go 長度。
set -euo pipefail
RECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RECT_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
RECT_OUT="$RECT_ROOT/workplace/matching-decompilation/c-rect"
test -d "$RECT_GOLEM/internal/machine";test -f "$RECT_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$RECT_OUT";test -O "$RECT_OUT"
RECT_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
RECT_CONTAINER="wolong-c-rect-$$-$RANDOM"
trap 'docker rm -f "$RECT_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$RECT_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$RECT_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$RECT_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$RECT_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$RECT_OUT,dst=/output" --env "RECT_MODE=${WOLONG_RECT_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$RECT_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/rect-build /output/results
    cp /repo/tools/c_recovery_rect.go /tmp/rect-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/rect-build/bus.go
    printf "module github.com/wicanr2/dosgolem/wolongcrect\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/rect-build/go.mod
    # Extract the exact production function, including nested closure, without another formula.
    awk "/^func battleSideBarLengths\(/ {capture=1} capture {print} capture && /^}/ {exit}" /repo/cmd/wlgame/battlelayout.go > /output/results/go-helper-source.txt
    test -s /output/results/go-helper-source.txt
    printf "//go:build matching_rect\npackage main\n" > /tmp/rect-build/bars.go
    awk "/^[[:space:]]*battleSideBarMaxLen[[:space:]]*=/ {print \"const \" \$0;exit}" /repo/cmd/wlgame/battlelayout.go > /output/results/go-helper-constant.txt
    test -s /output/results/go-helper-constant.txt
    cat /output/results/go-helper-constant.txt >> /tmp/rect-build/bars.go
    cat /output/results/go-helper-source.txt >> /tmp/rect-build/bars.go
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery_rect.go /repo/tools/c_recovery_vga_bus.go /repo/cmd/wlgame/battlelayout.go >> /output/results/c-source.sha256
    RECT_DIGEST="$(sha256sum /output/results/c-source.sha256)";RECT_DIGEST="${RECT_DIGEST%% *}"
    RECT_DEFINE="-DKI_RECT_SOURCE_DIGEST=0x${RECT_DIGEST:0:16}"
    printf "%s\n" "$RECT_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/rect-build
    if [[ "$RECT_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $RECT_DEFINE" go build -tags matching_rect,matching_vga -p 2 -o /output/rect-smoke .
      /output/rect-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $RECT_DEFINE" go build -tags matching_rect,matching_vga -p 2 -o "/output/rect-$optimize" .
        go version -m "/output/rect-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/rect-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:rectangle 2:rectangle 3:rectangle 4:rectangle 5:rectangle 6:rectangle 7:gauge 8:gauge 9:waiting 10:bar 11:rectangle 12:selection; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $RECT_DEFINE -DKI_RECT_MUTATION=$mutation" go build -tags matching_rect,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_RECT_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$RECT_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$RECT_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_rect_verify.py --repo /repo --output /output
fi
