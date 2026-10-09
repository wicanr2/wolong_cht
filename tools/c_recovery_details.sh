#!/usr/bin/env bash
# spec/242：據點資訊與軍團面板。
set -euo pipefail
DETAILS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DETAILS_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
DETAILS_OUT="$DETAILS_ROOT/workplace/matching-decompilation/c-details"
test -d "$DETAILS_GOLEM/internal/machine";test -f "$DETAILS_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$DETAILS_OUT";test -O "$DETAILS_OUT"
case "${WOLONG_DETAILS_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
DETAILS_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
DETAILS_CONTAINER="wolong-c-details-$$-$RANDOM"
trap 'docker rm -f "$DETAILS_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$DETAILS_ROOT,dst=/repo,readonly" --mount "type=bind,src=$DETAILS_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$DETAILS_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$DETAILS_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$DETAILS_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$DETAILS_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$DETAILS_OUT,dst=/output" --env "DETAILS_MODE=${WOLONG_DETAILS_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$DETAILS_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/details-build /output/results
    cp /repo/tools/c_recovery_details.go /tmp/details-build/main.go
    cp /repo/tools/c_recovery_details_data.go /tmp/details-build/data.go
    cp /repo/tools/c_recovery_details_platform.go /tmp/details-build/details-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/details-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/details-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/details-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcdetails\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/details-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$DETAILS_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/details-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/details-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/details_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_details_generate.py /repo/tools/c_recovery_details.go /repo/tools/c_recovery_details_data.go /repo/tools/c_recovery_details_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$DETAILS_MODE" == controls ]]; then
      cmp /tmp/details-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/details-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    DETAILS_DIGEST="$(sha256sum /output/results/c-source.sha256)";DETAILS_DIGEST="${DETAILS_DIGEST%% *}"
    DETAILS_DEFINE="-DKI_DETAILS_SOURCE_DIGEST=0x${DETAILS_DIGEST:0:16}"
    printf "%s\n" "$DETAILS_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/details-build
    if [[ "$DETAILS_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $DETAILS_DEFINE" go build -tags matching_details,matching_input,matching_glyph,matching_vga -p 2 -o /output/details-smoke .
      /output/details-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$DETAILS_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $DETAILS_DEFINE" go build -tags matching_details,matching_input,matching_glyph,matching_vga -p 2 -o "/output/details-$optimize" .
        go version -m "/output/details-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/details-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      for entry in 1:city-values 2:city-values 3:city-values 4:city-values 5:city-image 6:army-values 7:army-values 8:army-slots 9:army-slots 10:army-close 11:own-frame 12:city-window; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $DETAILS_DEFINE -DKI_DETAILS_MUTATION=$mutation" go build -tags matching_details,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_DETAILS_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$DETAILS_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$DETAILS_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_details_verify.py --repo /repo --output /output
fi
