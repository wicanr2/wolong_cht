#!/usr/bin/env bash
# spec/245：主迴圈非區域返回與調色盤淡出。
set -euo pipefail
MC_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MC_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
MC_OUT="$MC_ROOT/workplace/matching-decompilation/c-main"
test -d "$MC_GOLEM/internal/machine";test -f "$MC_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$MC_OUT";test -O "$MC_OUT"
case "${WOLONG_MAIN_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
MC_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
MC_CONTAINER="wolong-c-main-$$-$RANDOM"
trap 'docker rm -f "$MC_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$MC_ROOT,dst=/repo,readonly" --mount "type=bind,src=$MC_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$MC_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$MC_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$MC_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$MC_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$MC_OUT,dst=/output" --env "MC_MODE=${WOLONG_MAIN_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$MC_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/main-build /output/results
    cp /repo/tools/c_recovery_main.go /tmp/main-build/main.go
    cp /repo/tools/c_recovery_main_data.go /tmp/main-build/data.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/main-build/strategy-data.go
    cp /repo/tools/c_recovery_main_platform.go /tmp/main-build/main-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/main-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/main-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/main-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcmain\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/main-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$MC_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/main-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/main-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/main_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_main_generate.py /repo/tools/c_recovery_main.go /repo/tools/c_recovery_main_data.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_main_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$MC_MODE" == controls ]]; then
      cmp /tmp/main-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/main-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    MC_DIGEST="$(sha256sum /output/results/c-source.sha256)";MC_DIGEST="${MC_DIGEST%% *}"
    MC_DEFINE="-DKI_MAIN_SOURCE_DIGEST=0x${MC_DIGEST:0:16}"
    printf "%s\n" "$MC_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/main-build
    if [[ "$MC_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $MC_DEFINE" go build -tags matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/main-smoke .
      /output/main-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$MC_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $MC_DEFINE" go build -tags matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/main-$optimize" .
        go version -m "/output/main-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/main-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      : > /output/results/mutant-names.tsv
      for entry in 1:escape:saved-stack-segment 2:escape:saved-stack-pointer 3:escape:exit-ax-push 4:escape:exit-ax-restore 5:fade:fade-first-level 6:palette:dac-rounding 7:palette:dac-channel-order 8:trust:nonlocal-unwind; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $MC_DEFINE -DKI_MAIN_MUTATION=$mutation" go build -tags matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_MAIN_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$MC_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$MC_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_main_verify.py --repo /repo --output /output
fi
