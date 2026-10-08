#!/usr/bin/env bash
# spec/230：原始 INPUT、七標記、肖像快取與 DOS 讀檔。
set -euo pipefail
INPUT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INPUT_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
INPUT_OUT="$INPUT_ROOT/workplace/matching-decompilation/c-input"
test -d "$INPUT_GOLEM/internal/machine";test -f "$INPUT_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$INPUT_OUT";test -O "$INPUT_OUT"
INPUT_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
INPUT_CONTAINER="wolong-c-input-$$-$RANDOM"
trap 'docker rm -f "$INPUT_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$INPUT_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$INPUT_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$INPUT_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$INPUT_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$INPUT_OUT,dst=/output" --env "INPUT_MODE=${WOLONG_INPUT_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$INPUT_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/input-build /output/results
    cp /repo/tools/c_recovery_input.go /tmp/input-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/input-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/input-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/input-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcinput\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/input-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/input_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_input_generate.py /repo/tools/c_recovery_input.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    INPUT_DIGEST="$(sha256sum /output/results/c-source.sha256)";INPUT_DIGEST="${INPUT_DIGEST%% *}"
    INPUT_DEFINE="-DKI_INPUT_SOURCE_DIGEST=0x${INPUT_DIGEST:0:16}"
    printf "%s\n" "$INPUT_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/input-build
    if [[ "$INPUT_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $INPUT_DEFINE" go build -tags matching_input,matching_glyph,matching_vga -p 2 -o /output/input-smoke .
      /output/input-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $INPUT_DEFINE" go build -tags matching_input,matching_glyph,matching_vga -p 2 -o "/output/input-$optimize" .
        go version -m "/output/input-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/input-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:poll 2:wait 3:wait 4:wait 5:cursor 6:cursor 7:edges 8:edges 9:cursor 10:services 11:cursor 12:message; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $INPUT_DEFINE -DKI_INPUT_MUTATION=$mutation" go build -tags matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_INPUT_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$INPUT_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$INPUT_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_input_verify.py --repo /repo --output /output
fi
