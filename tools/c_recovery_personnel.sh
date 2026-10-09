#!/usr/bin/env bash
# spec/241：人事選單與任免。
set -euo pipefail
PERSONNEL_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PERSONNEL_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
PERSONNEL_OUT="$PERSONNEL_ROOT/workplace/matching-decompilation/c-personnel"
test -d "$PERSONNEL_GOLEM/internal/machine";test -f "$PERSONNEL_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$PERSONNEL_OUT";test -O "$PERSONNEL_OUT"
case "${WOLONG_PERSONNEL_MODE:-full}" in full|smoke|controls) ;; *) exit 2 ;; esac
PERSONNEL_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
PERSONNEL_CONTAINER="wolong-c-personnel-$$-$RANDOM"
trap 'docker rm -f "$PERSONNEL_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$PERSONNEL_ROOT,dst=/repo,readonly" --mount "type=bind,src=$PERSONNEL_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --name "$PERSONNEL_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$PERSONNEL_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$PERSONNEL_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$PERSONNEL_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$PERSONNEL_OUT,dst=/output" --env "PERSONNEL_MODE=${WOLONG_PERSONNEL_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$PERSONNEL_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/personnel-build /output/results
    cp /repo/tools/c_recovery_personnel.go /tmp/personnel-build/main.go
    cp /repo/tools/c_recovery_personnel_data.go /tmp/personnel-build/data.go
    cp /repo/tools/c_recovery_personnel_platform.go /tmp/personnel-build/personnel-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/personnel-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/personnel-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/personnel-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcpersonnel\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/personnel-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$PERSONNEL_MODE" == controls ]]; then
      cp /output/results/c-source.sha256 /tmp/personnel-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/personnel-golem-before.sha256
      test -s /output/results/O0.json;test -s /output/results/O2.json
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/personnel_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_personnel_generate.py /repo/tools/c_recovery_personnel.go /repo/tools/c_recovery_personnel_data.go /repo/tools/c_recovery_personnel_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$PERSONNEL_MODE" == controls ]]; then
      cmp /tmp/personnel-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/personnel-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    PERSONNEL_DIGEST="$(sha256sum /output/results/c-source.sha256)";PERSONNEL_DIGEST="${PERSONNEL_DIGEST%% *}"
    PERSONNEL_DEFINE="-DKI_PERSONNEL_SOURCE_DIGEST=0x${PERSONNEL_DIGEST:0:16}"
    printf "%s\n" "$PERSONNEL_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/personnel-build
    if [[ "$PERSONNEL_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $PERSONNEL_DEFINE" go build -tags matching_personnel,matching_input,matching_glyph,matching_vga -p 2 -o /output/personnel-smoke .
      /output/personnel-smoke -smoke -out /output/results/smoke.json
    else
      if [[ "$PERSONNEL_MODE" != controls ]]; then
       for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $PERSONNEL_DEFINE" go build -tags matching_personnel,matching_input,matching_glyph,matching_vga -p 2 -o "/output/personnel-$optimize" .
        go version -m "/output/personnel-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/personnel-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      for entry in 1:appoint 2:appoint 3:appoint 4:appoint 5:helper 6:helper 7:helper 8:helper 9:appoint 10:menu 11:appoint 12:dismiss; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $PERSONNEL_DEFINE -DKI_PERSONNEL_MUTATION=$mutation" go build -tags matching_personnel,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_PERSONNEL_MODE:-full}" != smoke ]]; then
  timeout 90 docker run --rm --init --log-opt max-size=10m --log-opt max-file=3 --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$PERSONNEL_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$PERSONNEL_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_personnel_verify.py --repo /repo --output /output
fi
