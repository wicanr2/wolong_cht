#!/usr/bin/env bash
# spec/236：完整場景入口、鏡頭與小地圖退出重畫。
set -euo pipefail
RESUME_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RESUME_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
RESUME_OUT="$RESUME_ROOT/workplace/matching-decompilation/c-resume"
test -d "$RESUME_GOLEM/internal/machine";test -f "$RESUME_ROOT/workplace/orig/dosv/KI.EXE"
test -d "$RESUME_OUT";test -O "$RESUME_OUT"
RESUME_IMAGE="$(docker image inspect golang:1.26.7-bookworm --format '{{.Id}}')"
RESUME_CONTAINER="wolong-c-resume-$$-$RANDOM"
trap 'docker rm -f "$RESUME_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout 60 docker run --rm --init --network none --memory 512m --cpus 1 --pids-limit 64 --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation --mount "type=bind,src=$RESUME_ROOT,dst=/repo,readonly" --mount "type=bind,src=$RESUME_OUT,dst=/output" python:3.13-bookworm python /repo/tools/c_recovery_mapcells_prepare.py --repo /repo --output /output
timeout --kill-after=10 2400 docker run --rm --init --name "$RESUME_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 --user "$(id -u):$(id -g)" \
  --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$RESUME_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$RESUME_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$RESUME_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$RESUME_OUT,dst=/output" --env "RESUME_MODE=${WOLONG_RESUME_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$RESUME_IMAGE" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/resume-build /output/results
    cp /repo/tools/c_recovery_resume.go /tmp/resume-build/main.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/resume-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/resume-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/resume-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcresume\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/resume-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/resume_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_resume_generate.py /repo/tools/c_recovery_resume.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    RESUME_DIGEST="$(sha256sum /output/results/c-source.sha256)";RESUME_DIGEST="${RESUME_DIGEST%% *}"
    RESUME_DEFINE="-DKI_RESUME_SOURCE_DIGEST=0x${RESUME_DIGEST:0:16}"
    printf "%s\n" "$RESUME_DIGEST" > /output/results/compiled-source-digest.txt
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/resume-build
    if [[ "$RESUME_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE $RESUME_DEFINE" go build -tags matching_resume,matching_input,matching_glyph,matching_vga -p 2 -o /output/resume-smoke .
      /output/resume-smoke -smoke -out /output/results/smoke.json
    else
      for optimize in O2 O0; do
        CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE $RESUME_DEFINE" go build -tags matching_resume,matching_input,matching_glyph,matching_vga -p 2 -o "/output/resume-$optimize" .
        go version -m "/output/resume-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/resume-$optimize" -out "/output/results/$optimize.json"
      done
      for entry in 1:camera 2:camera-helper 3:camera-helper 4:panel 5:minimap-clamp 6:minimap 7:minimap 8:minimap 9:minimap 10:scene 11:scene 12:scene; do
        mutation="${entry%%:*}";group="${entry#*:}"
        CGO_CFLAGS="-O0 -std=c11 -D_GNU_SOURCE $RESUME_DEFINE -DKI_RESUME_MUTATION=$mutation" go build -tags matching_resume,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
if [[ "${WOLONG_RESUME_MODE:-full}" == full ]]; then
  timeout 90 docker run --rm --init --network none --memory 768m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$RESUME_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$RESUME_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_resume_verify.py --repo /repo --output /output
fi
