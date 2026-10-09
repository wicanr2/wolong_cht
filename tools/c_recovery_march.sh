#!/usr/bin/env bash
# spec/243：軍團行軍選點與分派。
set -euo pipefail
MARCH_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MARCH_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
MARCH_OUT="$MARCH_ROOT/workplace/matching-decompilation/c-march"
test -d "$MARCH_GOLEM/internal/machine";test -f "$MARCH_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$MARCH_OUT";test -O "$MARCH_OUT"
case "${WOLONG_MARCH_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
MARCH_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
MARCH_CONTAINER="wolong-c-march-$$-$RANDOM"
trap 'docker rm -f "$MARCH_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$MARCH_ROOT,dst=/repo,readonly" --mount "type=bind,src=$MARCH_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$MARCH_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$MARCH_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$MARCH_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$MARCH_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$MARCH_OUT,dst=/output" --env "MARCH_MODE=${WOLONG_MARCH_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$MARCH_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/march-build /output/results
    cp /repo/tools/c_recovery_march.go /tmp/march-build/main.go
    cp /repo/tools/c_recovery_march_data.go /tmp/march-build/data.go
    cp /repo/tools/c_recovery_march_platform.go /tmp/march-build/march-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/march-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/march-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/march-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcmarch\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/march-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$MARCH_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/march-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/march-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/march_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_march_generate.py /repo/tools/c_recovery_march.go /repo/tools/c_recovery_march_data.go /repo/tools/c_recovery_march_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$MARCH_MODE" == controls ]]; then
      cmp /tmp/march-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/march-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    MARCH_DIGEST="$(sha256sum /output/results/c-source.sha256)";MARCH_DIGEST="${MARCH_DIGEST%% *}"
    MARCH_DEFINE="-DKI_MARCH_SOURCE_DIGEST=0x${MARCH_DIGEST:0:16}"
    printf "%s\n" "$MARCH_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/march-build
    if [[ "$MARCH_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $MARCH_DEFINE" go build -tags matching_march,matching_input,matching_glyph,matching_vga -p 2 -o /output/march-smoke .
      /output/march-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$MARCH_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $MARCH_DEFINE" go build -tags matching_march,matching_input,matching_glyph,matching_vga -p 2 -o "/output/march-$optimize" .
        go version -m "/output/march-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/march-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      for entry in 1:picker 2:hit 3:command 4:command 5:menu 6:dispatch 7:handler 8:handler 9:handler 10:handler 11:dissolve 12:controller 13:hit 14:cursor 15:handler 16:route; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $MARCH_DEFINE -DKI_MARCH_MUTATION=$mutation" go build -tags matching_march,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
        go version -m "/output/mutant-$mutation" > "/output/results/buildinfo-mutant-$mutation.txt"
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
if [[ "${WOLONG_MARCH_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$MARCH_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$MARCH_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_march_verify.py --repo /repo --output /output
fi
