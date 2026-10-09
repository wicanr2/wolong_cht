#!/usr/bin/env bash
# spec/251：容器內的有界原生工作；主機背景啟動後由本腳本保存退出收據。
set -euo pipefail

engagement_payload() {
    set -euo pipefail
    export GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
    mkdir -p /tmp/engagement-build /output/results
    cp /repo/tools/c_recovery_engagement.go /tmp/engagement-build/engagement.go
    cp /repo/tools/c_recovery_engagement_data.go /tmp/engagement-build/data.go
    cp /repo/tools/c_recovery_engagement_control_data.go /tmp/engagement-build/control-data.go
    cp /repo/tools/c_recovery_engagement_ui_data.go /tmp/engagement-build/ui-data.go
    cp /repo/tools/c_recovery_engagement_front_data.go /tmp/engagement-build/front-data.go
    cp /repo/tools/c_recovery_engagement_extra_data.go /tmp/engagement-build/extra-data.go
    cp /repo/tools/c_recovery_engagement_combat_data.go /tmp/engagement-build/combat-data.go
    cp /repo/tools/c_recovery_engagement_save.go /tmp/engagement-build/save.go
    cp /repo/tools/c_recovery_outcome_data.go /tmp/engagement-build/outcome-data.go
    cp /repo/tools/c_recovery_route_data.go /tmp/engagement-build/route-data.go
    cp /repo/tools/c_recovery_interaction_data.go /tmp/engagement-build/interaction-data.go
    cp /repo/tools/c_recovery_main_data.go /tmp/engagement-build/main-data.go
    cp /repo/tools/c_recovery_strategy_data.go /tmp/engagement-build/strategy-data.go
    cp /repo/tools/c_recovery_engagement_platform.go /tmp/engagement-build/engagement-platform.go
    cp /repo/tools/c_recovery_vga_bus.go /tmp/engagement-build/bus.go
    cp /repo/tools/c_recovery_glyph_platform.go /tmp/engagement-build/platform.go
    cp /repo/tools/c_recovery_input_platform.go /tmp/engagement-build/input-platform.go
    printf "module github.com/wicanr2/dosgolem/wolongcengagement\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n" > /tmp/engagement-build/go.mod
    git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
    if [[ "$EG_MODE" == controls || "$EG_MODE" == mutants ]]; then
      cp /output/results/c-source.sha256 /tmp/engagement-source-before.sha256
      cp /output/results/golem-source-before.sha256 /tmp/engagement-golem-before.sha256
      test -s /output/results/compiled-source-digest.txt
      cp /output/results/compiled-source-digest.txt /tmp/engagement-compiled-before.txt
      if [[ "$EG_MODE" == controls ]]; then
        test -s /output/results/O0.json;test -s /output/results/O2.json
      fi
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
    find /repo/tools/c_recovery -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" -o -name "*.inc" \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256
    sha256sum /repo/tools/c_recovery/engagement_generated.inc /repo/tools/c_recovery_glyph_platform.go /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_engagement_generate.py /repo/tools/c_recovery_engagement.go /repo/tools/c_recovery_engagement_data.go /repo/tools/c_recovery_engagement_control_data.go /repo/tools/c_recovery_engagement_ui_data.go /repo/tools/c_recovery_engagement_front_data.go /repo/tools/c_recovery_engagement_extra_data.go /repo/tools/c_recovery_engagement_combat_data.go /repo/tools/c_recovery_engagement_save.go /repo/tools/c_recovery_outcome_data.go /repo/tools/c_recovery_route_data.go /repo/tools/c_recovery_interaction_data.go /repo/tools/c_recovery_main_data.go /repo/tools/c_recovery_strategy_data.go /repo/tools/c_recovery_engagement_platform.go /repo/tools/c_recovery_vga_bus.go /repo/tools/c_recovery_input_platform.go >> /output/results/c-source.sha256
    if [[ "$EG_MODE" == controls || "$EG_MODE" == mutants ]]; then
      cmp /tmp/engagement-source-before.sha256 /output/results/c-source.sha256
      cmp /tmp/engagement-golem-before.sha256 /output/results/golem-source-before.sha256
    fi
    EG_DIGEST="$(sha256sum /output/results/c-source.sha256)";EG_DIGEST="${EG_DIGEST%% *}"
    EG_DEFINE="-DKI_ENGAGEMENT_SOURCE_DIGEST=0x${EG_DIGEST:0:16}"
    printf "%s\n" "$EG_DIGEST" > /output/results/compiled-source-digest.txt
    if [[ "$EG_MODE" == controls || "$EG_MODE" == mutants ]]; then
      cmp /tmp/engagement-compiled-before.txt /output/results/compiled-source-digest.txt
    fi
    go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
    cd /tmp/engagement-build
    if [[ "$EG_MODE" == smoke ]]; then
      CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $EG_DEFINE" go build -tags matching_engagement,matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o /output/engagement-smoke .
      smoke_args=(-smoke -out /output/results/smoke.json)
      if [[ -n "$EG_GROUP" ]]; then
        smoke_args=(-smoke -group "$EG_GROUP" -out "/output/results/probe-$EG_GROUP.json")
      fi
      /output/engagement-smoke "${smoke_args[@]}"
    else
      if [[ "$EG_MODE" != controls && "$EG_MODE" != mutants ]]; then
       optimizations=(O2 O0)
       case "$EG_MODE" in normal-o2) optimizations=(O2) ;; normal-o0) optimizations=(O0) ;; esac
       for optimize in "${optimizations[@]}"; do
        CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $EG_DEFINE" go build -tags matching_engagement,matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/engagement-$optimize" .
        go version -m "/output/engagement-$optimize" > "/output/results/buildinfo-$optimize.txt"
        "/output/engagement-$optimize" -out "/output/results/$optimize.json"
       done
      fi
      if [[ "$EG_MODE" == full || "$EG_MODE" == controls || "$EG_MODE" == mutants ]]; then
      : > /output/results/mutant-names.tsv
      for entry in 1:field:standoff-mask 2:siege:standoff-mask 3:terrain:terrain-limit 4:patch:timer-live-opcode 5:patch:unit-x-cutoff 6:patch:unit-y-cutoff 7:patch:unit-dl-live 8:patch:target-live-compare 9:patch:target-live-opcode 10:patch:path-live-opcode 11:script:script-dispatch 12:script:direction-dispatch 13:unit:unit-command 14:unit:unit-order 15:settings:settings-dispatch 16:direction:direction-table 17:projection:signed-shift 18:grid:grid-clear-count 19:settings:six-button-count 20:direction:negative-x; do
        mutation="${entry%%:*}";rest="${entry#*:}";group="${rest%%:*}";name="${rest#*:}"
        printf "%s\t%s\t%s\n" "$mutation" "$group" "$name" >> /output/results/mutant-names.tsv
        CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $EG_DEFINE -DKI_ENGAGEMENT_MUTATION=$mutation" go build -tags matching_engagement,matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga -p 2 -o "/output/mutant-$mutation" .
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
    fi
    find /golem/internal -type f -name "*.go" -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
    cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
}

if [[ "${1:-}" == --payload ]]; then
  shift
  test "$#" -eq 0
  engagement_payload
  exit 0
fi

test "$#" -eq 0
[[ "${EG_JOB_ID:-}" =~ ^[a-zA-Z0-9_.-]+$ ]]
[[ "${EG_TIMEOUT:-}" =~ ^[1-9][0-9]*$ ]]
test -d /output/results;test -O /output/results
EG_JOB_LOG="/output/results/$EG_JOB_ID.job.log"
EG_JOB_EXIT="/output/results/$EG_JOB_ID.job.exit"
test ! -e "$EG_JOB_LOG";test ! -e "$EG_JOB_EXIT"
umask 022
: > "$EG_JOB_LOG"
engagement_record_exit() {
  EG_JOB_STATUS=$?
  trap - EXIT
  printf "%s\n" "$EG_JOB_STATUS" > "$EG_JOB_EXIT.tmp"
  mv "$EG_JOB_EXIT.tmp" "$EG_JOB_EXIT"
}
trap engagement_record_exit EXIT
trap 'exit 143' TERM
trap 'exit 130' INT
set +e
timeout --kill-after=10 "$EG_TIMEOUT" /bin/bash "$0" --payload 2>&1 | tee -a "$EG_JOB_LOG"
EG_PIPE_STATUS=("${PIPESTATUS[@]}")
set -e
EG_JOB_STATUS="${EG_PIPE_STATUS[0]}"
if [[ "$EG_JOB_STATUS" -eq 0 && "${EG_PIPE_STATUS[1]}" -ne 0 ]]; then
  EG_JOB_STATUS="${EG_PIPE_STATUS[1]}"
fi
exit "$EG_JOB_STATUS"
