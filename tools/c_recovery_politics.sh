#!/usr/bin/env bash
# 月結政治與俘虜閉包；spec/211。
set -euo pipefail
POLITICS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLITICS_GOLEM="${WOLONG_DOSGOLEM:-/home/anr2/cht/dosgolem}"
POLITICS_OUT="$POLITICS_ROOT/workplace/matching-decompilation/c-politics"
POLITICS_IMAGE="${WOLONG_C_RNG_IMAGE:-golang:1.26.7-bookworm}"
test -d "$POLITICS_GOLEM/oracle";test -f "$POLITICS_ROOT/workplace/orig/dosv/KI.EXE"
mkdir -p "$POLITICS_OUT";test -O "$POLITICS_OUT"
POLITICS_IMAGE_ID="$(docker image inspect "$POLITICS_IMAGE" --format '{{.Id}}')"
POLITICS_CONTAINER="wolong-c-politics-$$-$RANDOM"
trap 'docker rm -f "$POLITICS_CONTAINER" >/dev/null 2>&1 || true' EXIT
timeout --kill-after=10 1200 docker run --rm --init --name "$POLITICS_CONTAINER" \
  --network none --memory 2g --cpus 2 --pids-limit 256 \
  --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
  --mount "type=bind,src=$POLITICS_ROOT,dst=/repo,readonly" \
  --mount "type=bind,src=$POLITICS_GOLEM,dst=/golem,readonly" \
  --mount "type=bind,src=$POLITICS_ROOT/workplace/orig/dosv,dst=/orig,readonly" \
  --mount "type=bind,src=$POLITICS_OUT,dst=/output" \
  --env "POLITICS_MODE=${WOLONG_POLITICS_MODE:-full}" \
  --workdir /repo --entrypoint /bin/bash "$POLITICS_IMAGE_ID" -c '
    set -euo pipefail
    export HOME=/tmp GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    printf "go 1.26.7\nuse (\n/repo\n/golem\n)\n" > /tmp/politics.go.work
    export GOWORK=/tmp/politics.go.work
    mkdir -p /output/results
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    sha256sum /repo/tools/c_recovery/politics.c /repo/tools/c_recovery/politics.h /repo/tools/c_recovery/politics_fixture.h /repo/tools/c_recovery/world_update.c /repo/tools/c_recovery/world_update.h /repo/tools/c_recovery/settlement.c /repo/tools/c_recovery/settlement.h /repo/tools/c_recovery/economy.c /repo/tools/c_recovery/economy.h /repo/tools/c_recovery/rng.c /repo/tools/c_recovery/rng.h /repo/tools/c_recovery_politics.go > /output/results/c-source.sha256
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    case "$POLITICS_MODE" in
      smoke)
        CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE" go build -p 2 -o /output/politics-probe /repo/tools/c_recovery_politics.go
        /output/politics-probe -smoke -out /output/results/smoke.json
        ;;
      full)
        for optimize in O2 O0; do
          CGO_CFLAGS="-$optimize -std=c11 -D_GNU_SOURCE" go build -p 2 -o "/output/politics-$optimize" /repo/tools/c_recovery_politics.go
          "/output/politics-$optimize" -out "/output/results/$optimize.json"
        done
        for entry in 1:writer 2:relation 3:initialize 4:assignment 5:power 6:generals 7:frontier 8:cooperation; do
          mutation="${entry%%:*}";group="${entry#*:}"
          CGO_CFLAGS="-O2 -std=c11 -D_GNU_SOURCE -DKI_POLITICS_MUTATION=$mutation" go build -p 2 -o "/output/mutant-$mutation" /repo/tools/c_recovery_politics.go
          set +e
          "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json"
          status=$?
          set -e
          test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
        done
        ;;
      *) exit 2 ;;
    esac
    find /golem/oracle /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
  '
if [[ "${WOLONG_POLITICS_MODE:-full}" == full ]]; then
  test -s "$POLITICS_OUT/ida/ida-probe.json"
  timeout --kill-after=10 45 docker run --rm --network none --memory 512m --cpus 1 --pids-limit 64 \
    --user "$(id -u):$(id -g)" --label wolong.task=matching-decompilation \
    --mount "type=bind,src=$POLITICS_ROOT,dst=/repo,readonly" \
    --mount "type=bind,src=$POLITICS_OUT,dst=/output" \
    python:3.13-bookworm python /repo/tools/c_recovery_politics_verify.py --repo /repo --output /output
fi
