#!/usr/bin/env bash
# spec/252：只在Docker內執行，持久保存逾時與退出收據。
set -euo pipefail
army_payload() {
 export GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2
 mkdir -p /tmp/army-build /output/results
 cp /repo/tools/c_recovery_army.go /tmp/army-build/army.go
 cp /repo/tools/c_recovery_army_data.go /tmp/army-build/army-data.go
 for source in engagement_data engagement_control_data engagement_ui_data engagement_front_data engagement_extra_data engagement_combat_data engagement_save outcome_data route_data interaction_data main_data strategy_data engagement_platform vga_bus glyph_platform input_platform; do
  cp "/repo/tools/c_recovery_$source.go" "/tmp/army-build/$source.go"
 done
 printf 'module github.com/wicanr2/dosgolem/wolongcarmy\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n' > /tmp/army-build/go.mod
 git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
 find /golem/internal -type f -name '*.go' -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
 find /repo/tools/c_recovery -maxdepth 1 -type f \( -name '*.c' -o -name '*.h' -o -name '*.inc' \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256.tmp
 sha256sum /repo/tools/c_recovery_army_generate.py /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_army.go /repo/tools/c_recovery_army_data.go >> /output/results/c-source.sha256.tmp
 for source in engagement_data engagement_control_data engagement_ui_data engagement_front_data engagement_extra_data engagement_combat_data engagement_save outcome_data route_data interaction_data main_data strategy_data engagement_platform vga_bus glyph_platform input_platform; do
  sha256sum "/repo/tools/c_recovery_$source.go" >> /output/results/c-source.sha256.tmp
 done
 if [[ "$AR_MODE" == mutants ]]; then
  cmp /output/results/c-source.sha256 /output/results/c-source.sha256.tmp
 fi
 mv /output/results/c-source.sha256.tmp /output/results/c-source.sha256
 cmp /repo/tools/c_recovery_army.go /tmp/army-build/army.go
 cmp /repo/tools/c_recovery_army_data.go /tmp/army-build/army-data.go
 for source in engagement_data engagement_control_data engagement_ui_data engagement_front_data engagement_extra_data engagement_combat_data engagement_save outcome_data route_data interaction_data main_data strategy_data engagement_platform vga_bus glyph_platform input_platform; do
  cmp "/repo/tools/c_recovery_$source.go" "/tmp/army-build/$source.go"
 done
 AR_DIGEST="$(sha256sum /output/results/c-source.sha256)";AR_DIGEST="${AR_DIGEST%% *}"
 printf '%s\n' "$AR_DIGEST" > /output/results/compiled-source-digest.txt
 AR_DEFINE="-DKI_ARMY_SOURCE_DIGEST=0x${AR_DIGEST:0:16}"
 go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
 cd /tmp/army-build
 AR_TAGS=matching_army,matching_engagement,matching_outcome,matching_route,matching_tick,matching_interaction,matching_bootstrap,matching_main,matching_strategy,matching_input,matching_glyph,matching_vga
 if [[ "$AR_MODE" == smoke ]]; then
  CGO_CFLAGS="-O2 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $AR_DEFINE" go build -tags "$AR_TAGS" -p 2 -o /output/army-smoke .
  AR_SMOKE_ARGS=(-smoke -out /output/results/smoke.json)
  if [[ -n "$AR_GROUP" ]]; then AR_SMOKE_ARGS+=(-group "$AR_GROUP");fi
  /output/army-smoke "${AR_SMOKE_ARGS[@]}"
  go version -m /output/army-smoke > /output/results/buildinfo-smoke.txt
 else
  if [[ "$AR_MODE" != mutants ]]; then
   AR_OPTIMIZATIONS=(O2 O0)
   case "$AR_MODE" in normal-o2) AR_OPTIMIZATIONS=(O2) ;; normal-o0) AR_OPTIMIZATIONS=(O0) ;; esac
   for optimize in "${AR_OPTIMIZATIONS[@]}"; do
    CGO_CFLAGS="-$optimize -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $AR_DEFINE" go build -tags "$AR_TAGS" -p 2 -o "/output/army-$optimize" .
    go version -m "/output/army-$optimize" > "/output/results/buildinfo-$optimize.txt"
    "/output/army-$optimize" -out "/output/results/$optimize.json"
   done
  fi
  if [[ "$AR_MODE" == full || "$AR_MODE" == mutants ]]; then
   test -s /repo/tools/c_recovery_army_mutants.tsv
   cp /repo/tools/c_recovery_army_mutants.tsv /output/results/mutant-names.tsv
   while IFS=$'\t' read -r mutation group name; do
    [[ "$mutation" =~ ^[1-9][0-9]*$ ]] || exit 2
    CGO_CFLAGS="-O0 -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $AR_DEFINE -DKI_ARMY_MUTATION=$mutation" go build -tags "$AR_TAGS" -p 2 -o "/output/mutant-$mutation" .
    go version -m "/output/mutant-$mutation" > "/output/results/buildinfo-mutant-$mutation.txt"
    rm -f "/output/results/mutant-$mutation.json"
    set +e
    "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json" > "/output/results/mutant-$mutation.log" 2>&1
    status=$?
    set -e
    printf '%s\n' "$status" > "/output/results/mutant-$mutation.exit"
    cat "/output/results/mutant-$mutation.log"
    test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
   done < /repo/tools/c_recovery_army_mutants.tsv
  fi
 fi
 find /golem/internal -type f -name '*.go' -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
 cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
}
if [[ "${1:-}" == --payload ]]; then shift;test "$#" -eq 0;army_payload;exit 0;fi
test "$#" -eq 0
[[ "${AR_JOB_ID:-}" =~ ^[a-zA-Z0-9_.-]+$ ]];[[ "${AR_TIMEOUT:-}" =~ ^[1-9][0-9]*$ ]]
test -O /output/results
AR_LOG="/output/results/$AR_JOB_ID.job.log";AR_EXIT="/output/results/$AR_JOB_ID.job.exit"
test ! -e "$AR_LOG";test ! -e "$AR_EXIT"
record_exit() { AR_STATUS=$?;trap - EXIT;printf '%s\n' "$AR_STATUS" > "$AR_EXIT.tmp";mv "$AR_EXIT.tmp" "$AR_EXIT"; }
trap record_exit EXIT
trap 'exit 143' TERM
trap 'exit 130' INT
set +e
timeout --kill-after=10 "$AR_TIMEOUT" /bin/bash "$0" --payload 2>&1 | tee "$AR_LOG"
AR_PIPE_STATUS=("${PIPESTATUS[@]}")
set -e
AR_STATUS="${AR_PIPE_STATUS[0]}"
if [[ "$AR_STATUS" -eq 0 && "${AR_PIPE_STATUS[1]}" -ne 0 ]];then AR_STATUS="${AR_PIPE_STATUS[1]}";fi
exit "$AR_STATUS"
