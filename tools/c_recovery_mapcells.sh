#!/usr/bin/env bash
# spec/234：原始世界顯示格初始化、推入與完整圖塊合成。
set -euo pipefail
MAPCELLS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MAPCELLS_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
MAPCELLS_OUT="$MAPCELLS_ROOT/workplace/matching-decompilation/c-mapcells"
test -d "$MAPCELLS_GOLEM/internal/machine"
test -f "$MAPCELLS_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$MAPCELLS_OUT";test -O "$MAPCELLS_OUT"
MAPCELLS_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
MAPCELLS_CONTAINER="wolong-c-mapcells-$$-$RANDOM"
trap 'docker rm -f "$MAPCELLS_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$MAPCELLS_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$MAPCELLS_OUT,dst=/output" \
  python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --name "$MAPCELLS_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$MAPCELLS_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$MAPCELLS_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$MAPCELLS_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$MAPCELLS_OUT,dst=/output" \
  --env "MAPCELLS_MODE=${WOLONG_MAPCELLS_MODE:-full}" --workdir /repo --entrypoint /bin/bash "$MAPCELLS_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/mapcells-build /output/results
    cp /repo/tools/c_recovery_mapcells.go /tmp/mapcells-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/mapcells-build/bus.go
    printf "module github.com/wicanr2/dosgolem/wolongcmapcells\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/mapcells-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery_mapcells.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_mapcells_generate.py /repo/tools/c_recovery_display_generate.py >> /output/results/c-source.sha256
    MAPCELLS_DIGEST="$(sha256sum /output/results/c-source.sha256)";MAPCELLS_DIGEST="${MAPCELLS_DIGEST%% *}"
    MAPCELLS_DEFINE="-DKI_MAPCELLS_SOURCE_DIGEST=0x${MAPCELLS_DIGEST:0:16}"
    printf "%s\n" "$MAPCELLS_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/mapcells-build
    if [[ "$MAPCELLS_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $MAPCELLS_DEFINE" go build -tags matching_mapcells,matching_vga -p 2 -o /output/mapcells-smoke .
      /output/mapcells-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $MAPCELLS_DEFINE" go build -tags matching_mapcells,matching_vga -p 2 -o "/output/mapcells-$optimize" .
        go version -m "/output/mapcells-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/mapcells-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:copy 2:copy 3:flags 4:flush 5:background 6:overlay 7:overlay 8:flush 9:flags 10:flags 11:push 12:push; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -D_GNU_SOURCE $MAPCELLS_DEFINE -DKI_MAPCELLS_MUTATION=$mutation" go build -tags matching_mapcells,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_MAPCELLS_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$MAPCELLS_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$MAPCELLS_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_mapcells_verify.py --repo /repo --output /output
fi
