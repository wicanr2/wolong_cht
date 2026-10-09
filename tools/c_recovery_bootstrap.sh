#!/usr/bin/env bash
# spec/246：獨立地圖與評分啟動前綴。
set -euo pipefail
BS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BS_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
BS_OUT="$BS_ROOT/workplace/matching-decompilation/c-bootstrap"
test -d "$BS_GOLEM/internal/machine";test -f "$BS_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$BS_OUT";test -O "$BS_OUT"
case "${WOLONG_BOOTSTRAP_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
BS_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
BS_CONTAINER="wolong-c-bootstrap-$$-$RANDOM"
trap 'docker rm -f "$BS_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$BS_ROOT,dst=/repo,readonly" --mount "type=bind,src=$BS_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$BS_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$BS_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$BS_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$BS_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$BS_OUT,dst=/output" --env "BS_MODE=${WOLONG_BOOTSTRAP_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$BS_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/bootstrap-build /output/results
    cp /repo/tools/c_recovery_bootstrap.go /tmp/bootstrap-build/bootstrap.go
    cp /repo/tools/c_recovery_bootstrap_data.go /tmp/bootstrap-build/data.go
    cp /repo/tools/c_recovery_main_data.go /tmp/bootstrap-build/main-data.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/bootstrap-build/strategy-data.go
    cp /repo/tools/c_recovery_bootstrap_platform.go /tmp/bootstrap-build/bootstrap-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/bootstrap-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/bootstrap-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/bootstrap-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcbootstrap\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/bootstrap-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$BS_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/bootstrap-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/bootstrap-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/bootstrap_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_bootstrap_generate.py /repo/tools/c_recovery_bootstrap.go /repo/tools/c_recovery_bootstrap_data.go /repo/tools/c_recovery_main_data.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_bootstrap_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$BS_MODE" == controls ]]; then
      cmp /tmp/bootstrap-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/bootstrap-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    BS_DIGEST="$(sha256sum /output/results/c-source.sha256)";BS_DIGEST="${BS_DIGEST%% *}"
    BS_DEFINE="-DKI_BOOTSTRAP_SOURCE_DIGEST=0x${BS_DIGEST:0:16}"
    printf "%s\n" "$BS_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/bootstrap-build
    if [[ "$BS_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $BS_DEFINE" go build -tags matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/bootstrap-smoke .
      /output/bootstrap-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$BS_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $BS_DEFINE" go build -tags matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/bootstrap-$optimize" .
        go version -m "/output/bootstrap-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/bootstrap-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      : > /output/results/mutant-names.tsv
      for entry in 1:city:center-divisor 2:city:neutral-owner-colour 3:city:own-fringe-colour 4:fringe:fringe-upper-bound 5:army:active-threshold 6:army:descriptor-segment 7:army:occupancy-increment 8:bulk:last-corps-slot 9:score:score-date-ds-result 10:prefix:saved-prefix-sp; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $BS_DEFINE -DKI_BOOTSTRAP_MUTATION=$mutation" go build -tags matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_BOOTSTRAP_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$BS_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$BS_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_bootstrap_verify.py --repo /repo --output /output
fi
