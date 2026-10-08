#!/usr/bin/env bash
# spec/230：原始 CHOICE、七標記、肖像快取與 DOS 讀檔。
set -euo pipefail
CHOICE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHOICE_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
CHOICE_OUT="$CHOICE_ROOT/workplace/matching-decompilation/c-choice"
test -d "$CHOICE_GOLEM/internal/machine";test -f "$CHOICE_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$CHOICE_OUT";test -O "$CHOICE_OUT"
CHOICE_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
CHOICE_CONTAINER="wolong-c-choice-$$-$RANDOM"
trap 'docker rm -f "$CHOICE_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 2400 docker run --rm --init --name "$CHOICE_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$CHOICE_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$CHOICE_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$CHOICE_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$CHOICE_OUT,dst=/output" --env "CHOICE_MODE=${WOLONG_CHOICE_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$CHOICE_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/choice-build /output/results
    cp /repo/tools/c_recovery_choice.go /tmp/choice-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/choice-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/choice-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/choice-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcchoice\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/choice-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/choice_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_choice_generate.py /repo/tools/c_recovery_choice.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    CHOICE_DIGEST="$(sha256sum /output/results/c-source.sha256)";CHOICE_DIGEST="${CHOICE_DIGEST%% *}"
    CHOICE_DEFINE="-DKI_CHOICE_SOURCE_DIGEST=0x${CHOICE_DIGEST:0:16}"
    printf "%s\n" "$CHOICE_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/choice-build
    if [[ "$CHOICE_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $CHOICE_DEFINE" go build -tags matching_choice,matching_input,matching_glyph,matching_vga -p 2 -o /output/choice-smoke .
      /output/choice-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $CHOICE_DEFINE" go build -tags matching_choice,matching_input,matching_glyph,matching_vga -p 2 -o "/output/choice-$optimize" .
        go version -m "/output/choice-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/choice-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:corpus 2:helpers 3:bands 4:scroll 5:selector 6:selector 7:helpers 8:helpers 9:selector 10:scroll 11:scroll 12:live; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $CHOICE_DEFINE -DKI_CHOICE_MUTATION=$mutation" go build -tags matching_choice,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_CHOICE_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$CHOICE_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$CHOICE_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_choice_verify.py --repo /repo --output /output
fi
