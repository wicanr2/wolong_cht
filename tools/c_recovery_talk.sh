#!/usr/bin/env bash
# spec/230：原始 TALK、七標記、肖像快取與 DOS 讀檔。
set -euo pipefail
TALK_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TALK_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
TALK_OUT="$TALK_ROOT/workplace/matching-decompilation/c-talk"
test -d "$TALK_GOLEM/internal/machine";test -f "$TALK_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$TALK_OUT";test -O "$TALK_OUT"
TALK_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
TALK_CONTAINER="wolong-c-talk-$$-$RANDOM"
trap 'docker rm -f "$TALK_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$TALK_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$TALK_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$TALK_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$TALK_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$TALK_OUT,dst=/output" --env "TALK_MODE=${WOLONG_TALK_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$TALK_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/talk-build /output/results
    cp /repo/tools/c_recovery_talk.go /tmp/talk-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/talk-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/talk-build/platform.go
    printf "module github.com/wicanr2/dosgolem/wolongctalk\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/talk-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/talk_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_talk_generate.py /repo/tools/c_recovery_talk.go /repo/tools/c_recovery_vga_bus.go >> /output/results/c-source.sha256
    TALK_DIGEST="$(sha256sum /output/results/c-source.sha256)";TALK_DIGEST="${TALK_DIGEST%% *}"
    TALK_DEFINE="-DKI_TALK_SOURCE_DIGEST=0x${TALK_DIGEST:0:16}"
    printf "%s\n" "$TALK_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/talk-build
    if [[ "$TALK_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $TALK_DEFINE" go build -tags matching_talk,matching_glyph,matching_vga -p 2 -o /output/talk-smoke .
      /output/talk-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $TALK_DEFINE" go build -tags matching_talk,matching_glyph,matching_vga -p 2 -o "/output/talk-$optimize" .
        go version -m "/output/talk-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/talk-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:corpus 2:strings 3:strings 4:strings 5:markers 6:cache 7:portrait 8:portrait 9:strings 10:corpus; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $TALK_DEFINE -DKI_TALK_MUTATION=$mutation" go build -tags matching_talk,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_TALK_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$TALK_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$TALK_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_talk_verify.py --repo /repo --output /output
fi
