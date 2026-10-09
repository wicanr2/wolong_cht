#!/usr/bin/env bash
# spec/248：據點輪轉、動畫物件與原提示控制。
set -euo pipefail
TK_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TK_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
TK_OUT="$TK_ROOT/workplace/matching-decompilation/c-tick"
test -d "$TK_GOLEM/internal/machine";test -f "$TK_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$TK_OUT";test -O "$TK_OUT"
case "${WOLONG_TICK_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
TK_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
TK_CONTAINER="wolong-c-tick-$$-$RANDOM"
trap 'docker rm -f "$TK_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$TK_ROOT,dst=/repo,readonly" --mount "type=bind,src=$TK_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$TK_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$TK_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$TK_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$TK_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$TK_OUT,dst=/output" --env "TK_MODE=${WOLONG_TICK_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$TK_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/tick-build /output/results
    cp /repo/tools/c_recovery_tick.go /tmp/tick-build/tick.go
    cp /repo/tools/c_recovery_tick_data.go /tmp/tick-build/data.go
    cp /repo/tools/c_recovery_interaction_data.go /tmp/tick-build/interaction-data.go
    cp /repo/tools/c_recovery_main_data.go /tmp/tick-build/main-data.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/tick-build/strategy-data.go
    cp /repo/tools/c_recovery_tick_platform.go /tmp/tick-build/tick-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/tick-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/tick-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/tick-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongctick\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/tick-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$TK_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/tick-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/tick-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/tick_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_tick_generate.py /repo/tools/c_recovery_tick.go /repo/tools/c_recovery_tick_data.go /repo/tools/c_recovery_interaction_data.go /repo/tools/c_recovery_main_data.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_tick_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$TK_MODE" == controls ]]; then
      cmp /tmp/tick-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/tick-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    TK_DIGEST="$(sha256sum /output/results/c-source.sha256)";TK_DIGEST="${TK_DIGEST%% *}"
    TK_DEFINE="-DKI_TICK_SOURCE_DIGEST=0x${TK_DIGEST:0:16}"
    printf "%s\n" "$TK_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/tick-build
    if [[ "$TK_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $TK_DEFINE" go build -tags matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/tick-smoke .
      /output/tick-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$TK_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $TK_DEFINE" go build -tags matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/tick-$optimize" .
        go version -m "/output/tick-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/tick-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      : > /output/results/mutant-names.tsv
      for entry in 1:objects:object-active-threshold 2:objects:object-slot-count 3:objects:object-mobile-boundary 4:objects:jitter-roll-mask 5:objects:velocity-lower-clamp 6:objects:jitter-step-threshold 7:city:city-cursor-wrap 8:city:city-owner-write-offset 9:neighbors:neighbor-neutral-branch 10:neighbors:neighbor-alliance-threshold 11:threat:status-flag-mask 12:growth:growth-default-rate 13:growth:growth-upper-limit 14:disaster:disaster-multiplier-byte 15:threat:original-distance-field 16:alert:alert-speaker-enable; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $TK_DEFINE -DKI_TICK_MUTATION=$mutation" go build -tags matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
        go version -m "/output/mutant-$mutation" > "/output/results/buildinfo-mutant-$mutation.txt"
        rm -f "/output/results/mutant-$mutation.json"
        set +e
        "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json" > "/output/results/mutant-$mutation.log" 2>&1
        status=$?
        set -e
        printf "%s\n" "$status" > "/output/results/mutant-$mutation.exit"
        cat "/output/results/mutant-$mutation.log"
        test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
      done
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
if [[ "${WOLONG_TICK_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$TK_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$TK_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_tick_verify.py --repo /repo --output /output
fi
