#!/usr/bin/env bash
# spec/244：政略指令列與上層玩家入口。
set -euo pipefail
STRATEGY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STRATEGY_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
STRATEGY_OUT="$STRATEGY_ROOT/workplace/matching-decompilation/c-strategy"
test -d "$STRATEGY_GOLEM/internal/machine";test -f "$STRATEGY_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$STRATEGY_OUT";test -O "$STRATEGY_OUT"
case "${WOLONG_STRATEGY_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
STRATEGY_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
STRATEGY_CONTAINER="wolong-c-strategy-$$-$RANDOM"
trap 'docker rm -f "$STRATEGY_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$STRATEGY_ROOT,dst=/repo,readonly" --mount "type=bind,src=$STRATEGY_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$STRATEGY_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$STRATEGY_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$STRATEGY_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$STRATEGY_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$STRATEGY_OUT,dst=/output" --env "STRATEGY_MODE=${WOLONG_STRATEGY_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$STRATEGY_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/strategy-build /output/results
    cp /repo/tools/c_recovery_strategy.go /tmp/strategy-build/main.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/strategy-build/data.go
    cp /repo/tools/c_recovery_strategy_platform.go /tmp/strategy-build/strategy-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/strategy-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/strategy-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/strategy-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcstrategy\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/strategy-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$STRATEGY_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/strategy-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/strategy-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/strategy_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_strategy_generate.py /repo/tools/c_recovery_strategy.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_strategy_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$STRATEGY_MODE" == controls ]]; then
      cmp /tmp/strategy-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/strategy-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    STRATEGY_DIGEST="$(sha256sum /output/results/c-source.sha256)";STRATEGY_DIGEST="${STRATEGY_DIGEST%% *}"
    STRATEGY_DEFINE="-DKI_STRATEGY_SOURCE_DIGEST=0x${STRATEGY_DIGEST:0:16}"
    printf "%s\n" "$STRATEGY_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/strategy-build
    if [[ "$STRATEGY_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $STRATEGY_DEFINE" go build -tags matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/strategy-smoke .
      /output/strategy-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$STRATEGY_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $STRATEGY_DEFINE" go build -tags matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/strategy-$optimize" .
        go version -m "/output/strategy-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/strategy-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      : > /output/results/mutant-names.tsv
      for entry in 1:bar:bar-coordinate-offset 2:gate:ruler-duty-gate 3:queue-search:queue-high-byte-comparison 4:trust:trust-budget-threshold 5:reason:reason-budget-decrement 6:reason:reason-attempted-mask 7:reason:reason-success-carry 8:trust:trust-saturation 9:finance-values:signed-money-extension 10:power:power-aggression-offset 11:proposal:war-relation-threshold 12:proposal:peace-aggression-halving 13:proposal:joint-target-equality 14:general:general-aptitude-maximum 15:general:general-segment-restore 16:scene:scene-success-trust-half; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $STRATEGY_DEFINE -DKI_STRATEGY_MUTATION=$mutation" go build -tags matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
        go version -m "/output/mutant-$mutation" > "/output/results/buildinfo-mutant-$mutation.txt"
        rm -f "/output/results/mutant-$mutation.json"
        set +e
        "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json" > "/output/results/mutant-$mutation.log" 2>&1
        status=$?
        set -e
        printf "%s\n" "$status" > "/output/results/mutant-$mutation.exit"
        if [[ "$mutation" == 9 ]]; then
          test "$status" -eq 2
          test ! -e "/output/results/mutant-$mutation.json"
          grep -F "dl_div: Assertion" "/output/results/mutant-$mutation.log" > /dev/null
          grep -F "SIGABRT: abort" "/output/results/mutant-$mutation.log" > /dev/null
          printf "Mutant 9 rejected by native DIV overflow guard\n"
        else
          cat "/output/results/mutant-$mutation.log"
          test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
        fi
      done
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
if [[ "${WOLONG_STRATEGY_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$STRATEGY_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$STRATEGY_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_strategy_verify.py --repo /repo --output /output
fi
