#!/usr/bin/env bash
# spec/247：世界點選、原始右鍵表與淡入。
set -euo pipefail
WI_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WI_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
WI_OUT="$WI_ROOT/workplace/matching-decompilation/c-loop"
test -d "$WI_GOLEM/internal/machine";test -f "$WI_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$WI_OUT";test -O "$WI_OUT"
case "${WOLONG_INTERACTION_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
WI_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
WI_CONTAINER="wolong-c-loop-$$-$RANDOM"
trap 'docker rm -f "$WI_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$WI_ROOT,dst=/repo,readonly" --mount "type=bind,src=$WI_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$WI_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$WI_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$WI_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$WI_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$WI_OUT,dst=/output" --env "WI_MODE=${WOLONG_INTERACTION_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$WI_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/interaction-build /output/results
    cp /repo/tools/c_recovery_interaction.go /tmp/interaction-build/interaction.go
    cp /repo/tools/c_recovery_interaction_data.go /tmp/interaction-build/data.go
    cp /repo/tools/c_recovery_main_data.go /tmp/interaction-build/main-data.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/interaction-build/strategy-data.go
    cp /repo/tools/c_recovery_interaction_platform.go /tmp/interaction-build/interaction-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/interaction-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/interaction-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/interaction-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcinteraction\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/interaction-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$WI_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/interaction-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/interaction-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/interaction_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_interaction_generate.py /repo/tools/c_recovery_interaction.go /repo/tools/c_recovery_interaction_data.go /repo/tools/c_recovery_main_data.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_interaction_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$WI_MODE" == controls ]]; then
      cmp /tmp/interaction-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/interaction-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    WI_DIGEST="$(sha256sum /output/results/c-source.sha256)";WI_DIGEST="${WI_DIGEST%% *}"
    WI_DEFINE="-DKI_INTERACTION_SOURCE_DIGEST=0x${WI_DIGEST:0:16}"
    printf "%s\n" "$WI_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/interaction-build
    if [[ "$WI_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $WI_DEFINE" go build -tags matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/interaction-smoke .
      /output/interaction-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$WI_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $WI_DEFINE" go build -tags matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/interaction-$optimize" .
        go version -m "/output/interaction-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/interaction-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      : > /output/results/mutant-names.tsv
      for entry in 1:fade:fade-final-level 2:fade:second-vblank-wait 3:popup:popup-y-clamp 4:popup:popup-x-clamp 5:point:city-tile-lower-bound 6:point:city-tile-upper-bound 7:point:popup-carry-cancel 8:right:right-index-reset 9:right:live-table-read 10:clear:clear-bit-two 11:clear:clear-bit-one 12:clear:clear-bit-zero; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $WI_DEFINE -DKI_INTERACTION_MUTATION=$mutation" go build -tags matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_INTERACTION_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$WI_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$WI_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_interaction_verify.py --repo /repo --output /output
fi
