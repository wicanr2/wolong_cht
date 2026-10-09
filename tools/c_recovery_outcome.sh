#!/usr/bin/env bash
# spec/250：自動戰鬥、退卻與據點易主。
set -euo pipefail
OC_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OC_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
OC_OUT="$OC_ROOT/workplace/matching-decompilation/c-outcome"
test -d "$OC_GOLEM/internal/machine";test -f "$OC_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$OC_OUT";test -O "$OC_OUT"
case "${WOLONG_OUTCOME_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
OC_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
OC_CONTAINER="wolong-c-outcome-$$-$RANDOM"
trap 'docker rm -f "$OC_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$OC_ROOT,dst=/repo,readonly" --mount "type=bind,src=$OC_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$OC_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$OC_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$OC_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$OC_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$OC_OUT,dst=/output" --env "OC_MODE=${WOLONG_OUTCOME_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$OC_IMAGE" -c '
    set -euo pipefail
    export GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/outcome-build /output/results
    cp /repo/tools/c_recovery_outcome.go /tmp/outcome-build/outcome.go
    cp /repo/tools/c_recovery_outcome_data.go /tmp/outcome-build/data.go
    cp /repo/tools/c_recovery_route_data.go /tmp/outcome-build/route-data.go
    cp /repo/tools/c_recovery_interaction_data.go /tmp/outcome-build/interaction-data.go
    cp /repo/tools/c_recovery_main_data.go /tmp/outcome-build/main-data.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/outcome-build/strategy-data.go
    cp /repo/tools/c_recovery_outcome_platform.go /tmp/outcome-build/outcome-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/outcome-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/outcome-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/outcome-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcoutcome\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/outcome-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$OC_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/outcome-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/outcome-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/outcome_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_outcome_generate.py /repo/tools/c_recovery_outcome.go /repo/tools/c_recovery_outcome_data.go /repo/tools/c_recovery_route_data.go /repo/tools/c_recovery_interaction_data.go /repo/tools/c_recovery_main_data.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_outcome_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$OC_MODE" == controls ]]; then
      cmp /tmp/outcome-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/outcome-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    OC_DIGEST="$(sha256sum /output/results/c-source.sha256)";OC_DIGEST="${OC_DIGEST%% *}"
    OC_DEFINE="-DKI_OUTCOME_SOURCE_DIGEST=0x${OC_DIGEST:0:16}"
    printf "%s\n" "$OC_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/outcome-build
    if [[ "$OC_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $OC_DEFINE" go build -tags matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/outcome-smoke .
      /output/outcome-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$OC_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $OC_DEFINE" go build -tags matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/outcome-$optimize" .
        go version -m "/output/outcome-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/outcome-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      : > /output/results/mutant-names.tsv
      for entry in 1:battle:first-side-coefficient 2:battle:morale-divisor 3:battle:aptitude-shift 4:battle:general-rng-mask 5:battle:battle-ratio-offset 6:battle:battle-ratio-limit 7:battle:loser-base-loss 8:battle:leader-minimum 9:retreat:retreat-men-threshold 10:retreat:own-city-stay 11:capture:old-city-count 12:capture:new-city-count 13:capture:governor-duty 14:capture:failed-capital-faction-mask 15:capture:old-corps-redirection 16:neighbors:neighbor-city-count 17:fall:surviving-faction-count 18:diplomat:diplomat-duty 19:history:city-owner-history 20:army:army-occupancy; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $OC_DEFINE -DKI_OUTCOME_MUTATION=$mutation" go build -tags matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_OUTCOME_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$OC_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$OC_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_outcome_verify.py --repo /repo --output /output
fi
