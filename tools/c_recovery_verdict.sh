#!/usr/bin/env bash
# spec/237：君主出陣、自動編成與側欄。
set -euo pipefail
VERDICT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERDICT_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
VERDICT_OUT="$VERDICT_ROOT/workplace/matching-decompilation/c-verdict"
test -d "$VERDICT_GOLEM/internal/machine";test -f "$VERDICT_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$VERDICT_OUT";test -O "$VERDICT_OUT"
VERDICT_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
VERDICT_CONTAINER="wolong-c-verdict-$$-$RANDOM"
trap 'docker rm -f "$VERDICT_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$VERDICT_ROOT,dst=/repo,readonly" --mount "type=bind,src=$VERDICT_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --name "$VERDICT_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$VERDICT_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$VERDICT_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$VERDICT_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$VERDICT_OUT,dst=/output" --env "VERDICT_MODE=${WOLONG_VERDICT_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$VERDICT_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/verdict-build /output/results
    cp /repo/tools/c_recovery_verdict.go /tmp/verdict-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/verdict-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/verdict-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/verdict-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcverdict\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/verdict-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/verdict_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_verdict_generate.py /repo/tools/c_recovery_verdict.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    VERDICT_DIGEST="$(sha256sum /output/results/c-source.sha256)";VERDICT_DIGEST="${VERDICT_DIGEST%% *}"
    VERDICT_DEFINE="-DKI_VERDICT_SOURCE_DIGEST=0x${VERDICT_DIGEST:0:16}"
    printf "%s\n" "$VERDICT_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/verdict-build
    if [[ "$VERDICT_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $VERDICT_DEFINE" go build -tags matching_verdict,matching_input,matching_glyph,matching_vga -p 2 -o /output/verdict-smoke .
      /output/verdict-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $VERDICT_DEFINE" go build -tags matching_verdict,matching_input,matching_glyph,matching_vga -p 2 -o "/output/verdict-$optimize" .
        go version -m "/output/verdict-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/verdict-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:sortie 2:sortie 3:formation 4:formation 5:troops 6:troops 7:formation 8:derived 9:sortie 10:sidebar; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -D_GNU_SOURCE $VERDICT_DEFINE -DKI_VERDICT_MUTATION=$mutation" go build -tags matching_verdict,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_VERDICT_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$VERDICT_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$VERDICT_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_verdict_verify.py --repo /repo --output /output
fi
