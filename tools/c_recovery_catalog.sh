#!/usr/bin/env bash
# spec/239：四類清單、builder 與 renderer。
set -euo pipefail
CATALOG_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CATALOG_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
CATALOG_OUT="$CATALOG_ROOT/workplace/matching-decompilation/c-catalog"
test -d "$CATALOG_GOLEM/internal/machine";test -f "$CATALOG_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$CATALOG_OUT";test -O "$CATALOG_OUT"
CATALOG_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
CATALOG_CONTAINER="wolong-c-catalog-$$-$RANDOM"
trap 'docker rm -f "$CATALOG_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$CATALOG_ROOT,dst=/repo,readonly" --mount "type=bind,src=$CATALOG_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$CATALOG_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$CATALOG_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$CATALOG_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$CATALOG_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$CATALOG_OUT,dst=/output" --env "CATALOG_MODE=${WOLONG_CATALOG_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$CATALOG_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/catalog-build /output/results
    cp /repo/tools/c_recovery_catalog.go /tmp/catalog-build/main.go
    cp /repo/tools/c_recovery_catalog_data.go /tmp/catalog-build/data.go
    cp /repo/tools/c_recovery_catalog_platform.go /tmp/catalog-build/catalog-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/catalog-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/catalog-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/catalog-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongccatalog\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/catalog-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/catalog_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_catalog_generate.py /repo/tools/c_recovery_catalog.go /repo/tools/c_recovery_catalog_data.go /repo/tools/c_recovery_catalog_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    CATALOG_DIGEST="$(sha256sum /output/results/c-source.sha256)";CATALOG_DIGEST="${CATALOG_DIGEST%% *}"
    CATALOG_DEFINE="-DKI_CATALOG_SOURCE_DIGEST=0x${CATALOG_DIGEST:0:16}"
    printf "%s\n" "$CATALOG_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/catalog-build
    if [[ "$CATALOG_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $CATALOG_DEFINE" go build -tags matching_catalog,matching_input,matching_glyph,matching_vga -p 2 -o /output/catalog-smoke .
      /output/catalog-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $CATALOG_DEFINE" go build -tags matching_catalog,matching_input,matching_glyph,matching_vga -p 2 -o "/output/catalog-$optimize" .
        go version -m "/output/catalog-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/catalog-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:builder 2:builder 3:builder 4:builder 5:caller 6:army 7:person 8:person 9:builder 10:cache 11:cache 12:caller 13:pages; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -D_GNU_SOURCE $CATALOG_DEFINE -DKI_CATALOG_MUTATION=$mutation" go build -tags matching_catalog,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_CATALOG_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$CATALOG_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$CATALOG_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_catalog_verify.py --repo /repo --output /output
fi
