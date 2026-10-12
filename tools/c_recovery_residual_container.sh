#!/usr/bin/env bash
# spec/254：只在Docker內執行，持久保存逾時與退出收據。
set -euo pipefail
residual_payload() {
 export GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2 CGO_LDFLAGS="-Wl,--gc-sections"
 mkdir -p /tmp/residual-build /output/results
 cp /repo/tools/c_recovery_residual.go /tmp/residual-build/residual.go
 cp /repo/tools/c_recovery_residual_data.go /tmp/residual-build/residual-data.go
 for source in residual_platform vga_bus glyph_platform; do
  cp "/repo/tools/c_recovery_$source.go" "/tmp/residual-build/$source.go"
 done
 printf 'module github.com/wicanr2/dosgolem/wolongcresidual\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n' > /tmp/residual-build/go.mod
 git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
 find /golem/internal -type f -name '*.go' -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
 find /repo/tools/c_recovery -maxdepth 1 -type f \( -name '*.c' -o -name '*.h' -o -name '*.inc' \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256.tmp
 sha256sum /repo/tools/c_recovery_residual_generate.py /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_residual.go /repo/tools/c_recovery_residual_data.go >> /output/results/c-source.sha256.tmp
 for source in residual_platform vga_bus glyph_platform; do
  sha256sum "/repo/tools/c_recovery_$source.go" >> /output/results/c-source.sha256.tmp
 done
 if [[ "$RS_MODE" == mutants ]]; then
  cmp /output/results/c-source.sha256 /output/results/c-source.sha256.tmp
 fi
 mv /output/results/c-source.sha256.tmp /output/results/c-source.sha256
 cmp /repo/tools/c_recovery_residual.go /tmp/residual-build/residual.go
 cmp /repo/tools/c_recovery_residual_data.go /tmp/residual-build/residual-data.go
 for source in residual_platform vga_bus glyph_platform; do
  cmp "/repo/tools/c_recovery_$source.go" "/tmp/residual-build/$source.go"
 done
 # 保留實際編譯來源；後續研究改檔時，歷史收據仍可回查原始文字。
 RS_SNAPSHOT="/output/results/$RS_JOB_ID.sources.tar.gz"
 test ! -e "$RS_SNAPSHOT"
 mkdir -p /tmp/residual-source-snapshot
 awk '{sub("^/repo/", "", $2);print $2}' /output/results/c-source.sha256 > /tmp/residual-source-files.txt
 tar -C /repo -cf - -T /tmp/residual-source-files.txt | tar -C /tmp/residual-source-snapshot -xf -
 (
  cd /tmp/residual-source-snapshot
  sed 's#  /repo/#  #' /output/results/c-source.sha256 | sha256sum -c --quiet
 )
 cp /output/results/c-source.sha256 /tmp/residual-source-snapshot/manifest.sha256
 cp /tmp/residual-build/go.mod /tmp/residual-source-snapshot/module-go.mod
 tar -C /tmp/residual-source-snapshot -czf "$RS_SNAPSHOT" .
 sha256sum "$RS_SNAPSHOT" > "$RS_SNAPSHOT.sha256"
 RS_DIGEST="$(sha256sum /output/results/c-source.sha256)";RS_DIGEST="${RS_DIGEST%% *}"
 printf '%s\n' "$RS_DIGEST" > /output/results/compiled-source-digest.txt
 RS_DEFINE="-DKI_RESIDUAL_SOURCE_DIGEST=0x${RS_DIGEST:0:16}"
 go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
 cd /tmp/residual-build
 RS_TAGS=matching_residual,matching_glyph,matching_vga
 if [[ "$RS_MODE" == smoke ]]; then
  CGO_CFLAGS="-O2 -ffunction-sections -fdata-sections -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $RS_DEFINE" go build -tags "$RS_TAGS" -p 2 -o /output/residual-smoke .
  RS_SMOKE_ARGS=(-smoke -out /output/results/smoke.json)
  if [[ -n "$RS_GROUP" ]]; then RS_SMOKE_ARGS+=(-group "$RS_GROUP");fi
  /output/residual-smoke "${RS_SMOKE_ARGS[@]}"
  go version -m /output/residual-smoke > /output/results/buildinfo-smoke.txt
 else
  if [[ "$RS_MODE" != mutants ]]; then
   RS_OPTIMIZATIONS=(O2 O0)
   case "$RS_MODE" in normal-o2) RS_OPTIMIZATIONS=(O2) ;; normal-o0) RS_OPTIMIZATIONS=(O0) ;; esac
   for optimize in "${RS_OPTIMIZATIONS[@]}"; do
    CGO_CFLAGS="-$optimize -ffunction-sections -fdata-sections -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $RS_DEFINE" go build -tags "$RS_TAGS" -p 2 -o "/output/residual-$optimize" .
    go version -m "/output/residual-$optimize" > "/output/results/buildinfo-$optimize.txt"
    "/output/residual-$optimize" -out "/output/results/$optimize.json"
   done
  fi
  if [[ "$RS_MODE" == full || "$RS_MODE" == mutants ]]; then
   test -s /repo/tools/c_recovery_residual_mutants.tsv
   cp /repo/tools/c_recovery_residual_mutants.tsv /output/results/mutant-names.tsv
   while IFS=$'\t' read -r mutation group name; do
    [[ "$mutation" =~ ^[1-9][0-9]*$ ]] || exit 2
    CGO_CFLAGS="-O0 -ffunction-sections -fdata-sections -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $RS_DEFINE -DKI_RESIDUAL_MUTATION=$mutation" go build -tags "$RS_TAGS" -p 2 -o "/output/mutant-$mutation" .
    go version -m "/output/mutant-$mutation" > "/output/results/buildinfo-mutant-$mutation.txt"
    rm -f "/output/results/mutant-$mutation.json"
    set +e
    "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json" > "/output/results/mutant-$mutation.log" 2>&1
    status=$?
    set -e
    printf '%s\n' "$status" > "/output/results/mutant-$mutation.exit"
    cat "/output/results/mutant-$mutation.log"
    test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
   done < /repo/tools/c_recovery_residual_mutants.tsv
  fi
 fi
 find /golem/internal -type f -name '*.go' -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
 cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
}
if [[ "${1:-}" == --payload ]]; then shift;test "$#" -eq 0;residual_payload;exit 0;fi
test "$#" -eq 0
[[ "${RS_JOB_ID:-}" =~ ^[a-zA-Z0-9_.-]+$ ]];[[ "${RS_TIMEOUT:-}" =~ ^[1-9][0-9]*$ ]]
test -O /output/results
RS_LOG="/output/results/$RS_JOB_ID.job.log";RS_EXIT="/output/results/$RS_JOB_ID.job.exit"
test ! -e "$RS_LOG";test ! -e "$RS_EXIT"
record_exit() { RS_STATUS=$?;trap - EXIT;printf '%s\n' "$RS_STATUS" > "$RS_EXIT.tmp";mv "$RS_EXIT.tmp" "$RS_EXIT"; }
trap record_exit EXIT
trap 'exit 143' TERM
trap 'exit 130' INT
set +e
timeout --kill-after=10 "$RS_TIMEOUT" /bin/bash "$0" --payload 2>&1 | tee "$RS_LOG"
RS_PIPE_STATUS=("${PIPESTATUS[@]}")
set -e
RS_STATUS="${RS_PIPE_STATUS[0]}"
if [[ "$RS_STATUS" -eq 0 && "${RS_PIPE_STATUS[1]}" -ne 0 ]];then RS_STATUS="${RS_PIPE_STATUS[1]}";fi
exit "$RS_STATUS"
