#!/usr/bin/env bash
# spec/238：據點清單、排序與遷都。
set -euo pipefail
LIST_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LIST_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
LIST_OUT="$LIST_ROOT/workplace/matching-decompilation/c-list"
test -d "$LIST_GOLEM/internal/machine";test -f "$LIST_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$LIST_OUT";test -O "$LIST_OUT"
LIST_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
LIST_CONTAINER="wolong-c-list-$$-$RANDOM"
trap 'docker rm -f "$LIST_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$LIST_ROOT,dst=/repo,readonly" --mount "type=bind,src=$LIST_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --name "$LIST_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$LIST_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$LIST_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$LIST_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$LIST_OUT,dst=/output" --env "LIST_MODE=${WOLONG_LIST_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$LIST_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/list-build /output/results
    cp /repo/tools/c_recovery_list.go /tmp/list-build/main.go
    cp /repo/tools/c_recovery_list_platform.go /tmp/list-build/list-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/list-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/list-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/list-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongclist\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/list-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/list_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_list_generate.py /repo/tools/c_recovery_list.go /repo/tools/c_recovery_list_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    LIST_DIGEST="$(sha256sum /output/results/c-source.sha256)";LIST_DIGEST="${LIST_DIGEST%% *}"
    LIST_DEFINE="-DKI_LIST_SOURCE_DIGEST=0x${LIST_DIGEST:0:16}"
    printf "%s\n" "$LIST_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/list-build
    if [[ "$LIST_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $LIST_DEFINE" go build -tags matching_list,matching_input,matching_glyph,matching_vga -p 2 -o /output/list-smoke .
      /output/list-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $LIST_DEFINE" go build -tags matching_list,matching_input,matching_glyph,matching_vga -p 2 -o "/output/list-$optimize" .
        go version -m "/output/list-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/list-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:rows 2:rows 3:list-cancel 4:sort 5:sort 6:selection 7:scroll 8:scroll-map 9:relocation 10:relocation 11:relocation 12:selection; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -D_GNU_SOURCE $LIST_DEFINE -DKI_LIST_MUTATION=$mutation" go build -tags matching_list,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_LIST_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$LIST_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$LIST_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_list_verify.py --repo /repo --output /output
fi
