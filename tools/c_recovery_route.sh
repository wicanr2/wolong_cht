#!/usr/bin/env bash
# spec/249：自我修改尋路、退卻與軍團潰散。
set -euo pipefail
RT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RT_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
RT_OUT="$RT_ROOT/workplace/matching-decompilation/c-route"
test -d "$RT_GOLEM/internal/machine";test -f "$RT_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$RT_OUT";test -O "$RT_OUT"
case "${WOLONG_ROUTE_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
RT_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
RT_CONTAINER="wolong-c-route-$$-$RANDOM"
trap 'docker rm -f "$RT_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$RT_ROOT,dst=/repo,readonly" --mount "type=bind,src=$RT_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$RT_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$RT_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$RT_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$RT_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$RT_OUT,dst=/output" --env "RT_MODE=${WOLONG_ROUTE_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$RT_IMAGE" -c '
    set -euo pipefail
    export GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/route-build /output/results
    cp /repo/tools/c_recovery_route.go /tmp/route-build/route.go
    cp /repo/tools/c_recovery_route_data.go /tmp/route-build/data.go
    cp /repo/tools/c_recovery_interaction_data.go /tmp/route-build/interaction-data.go
    cp /repo/tools/c_recovery_main_data.go /tmp/route-build/main-data.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/route-build/strategy-data.go
    cp /repo/tools/c_recovery_route_platform.go /tmp/route-build/route-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/route-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/route-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/route-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcroute\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/route-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$RT_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/route-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/route-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/route_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_route_generate.py /repo/tools/c_recovery_route.go /repo/tools/c_recovery_route_data.go /repo/tools/c_recovery_interaction_data.go /repo/tools/c_recovery_main_data.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_route_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$RT_MODE" == controls ]]; then
      cmp /tmp/route-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/route-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    RT_DIGEST="$(sha256sum /output/results/c-source.sha256)";RT_DIGEST="${RT_DIGEST%% *}"
    RT_DEFINE="-DKI_ROUTE_SOURCE_DIGEST=0x${RT_DIGEST:0:16}"
    printf "%s\n" "$RT_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/route-build
    if [[ "$RT_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $RT_DEFINE" go build -tags matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/route-smoke .
      /output/route-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$RT_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $RT_DEFINE" go build -tags matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/route-$optimize" .
        go version -m "/output/route-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/route-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      : > /output/results/mutant-names.tsv
      for entry in 1:search:terminal-a-live-read 2:search:terminal-b-live-read 3:search:owner-live-read 4:search:node-cost 5:search:foreign-node-cost 6:search:foreign-cost-flag 7:edge:edge-length-byte 8:edge:visited-zero-sentinel 9:search:minimum-cost-bucket 10:edge:positive-edge-step 11:replan:blocked-return-stage 12:replan:leg-direction-flag 13:collapse:capture-threshold 14:collapse:capture-occupancy 15:collapse:destroyed-army-flags 16:minimap:dirty-bit-mask; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $RT_DEFINE -DKI_ROUTE_MUTATION=$mutation" go build -tags matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_ROUTE_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$RT_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$RT_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_route_verify.py --repo /repo --output /output
fi
