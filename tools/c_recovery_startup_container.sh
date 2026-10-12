#!/usr/bin/env bash
# spec/253：只在Docker內執行，持久保存逾時與退出收據。
set -euo pipefail
startup_payload() {
 export GOPATH=/tmp/gopath GOCACHE=/output/gocache GOPROXY=off GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=1 GOMAXPROCS=2 CGO_LDFLAGS="-Wl,--gc-sections"
 mkdir -p /tmp/startup-build /output/results
 cp /repo/tools/c_recovery_startup.go /tmp/startup-build/startup.go
 cp /repo/tools/c_recovery_startup_data.go /tmp/startup-build/startup-data.go
 for source in startup_platform vga_bus glyph_platform; do
  cp "/repo/tools/c_recovery_$source.go" "/tmp/startup-build/$source.go"
 done
 printf 'module github.com/wicanr2/dosgolem/wolongcstartup\n\ngo 1.26.7\nrequire github.com/wicanr2/dosgolem v0.0.0\nreplace github.com/wicanr2/dosgolem => /golem\n' > /tmp/startup-build/go.mod
 git -C /golem rev-parse HEAD > /output/results/golem-revision.txt
 find /golem/internal -type f -name '*.go' -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-before.sha256
 find /repo/tools/c_recovery -maxdepth 1 -type f \( -name '*.c' -o -name '*.h' -o -name '*.inc' \) -print0 | sort -z | xargs -0 sha256sum > /output/results/c-source.sha256.tmp
 sha256sum /repo/tools/c_recovery_startup_generate.py /repo/tools/c_recovery_display_generate.py /repo/tools/c_recovery_startup.go /repo/tools/c_recovery_startup_data.go >> /output/results/c-source.sha256.tmp
 for source in startup_platform vga_bus glyph_platform; do
  sha256sum "/repo/tools/c_recovery_$source.go" >> /output/results/c-source.sha256.tmp
 done
 if [[ "$SU_MODE" == mutants ]]; then
  cmp /output/results/c-source.sha256 /output/results/c-source.sha256.tmp
 fi
 mv /output/results/c-source.sha256.tmp /output/results/c-source.sha256
 cmp /repo/tools/c_recovery_startup.go /tmp/startup-build/startup.go
 cmp /repo/tools/c_recovery_startup_data.go /tmp/startup-build/startup-data.go
 for source in startup_platform vga_bus glyph_platform; do
  cmp "/repo/tools/c_recovery_$source.go" "/tmp/startup-build/$source.go"
 done
 # 保留實際編譯來源；後續研究改檔時，歷史收據仍可回查原始文字。
 SU_SNAPSHOT="/output/results/$SU_JOB_ID.sources.tar.gz"
 test ! -e "$SU_SNAPSHOT"
 mkdir -p /tmp/startup-source-snapshot
 awk '{sub("^/repo/", "", $2);print $2}' /output/results/c-source.sha256 > /tmp/startup-source-files.txt
 tar -C /repo -cf - -T /tmp/startup-source-files.txt | tar -C /tmp/startup-source-snapshot -xf -
 (
  cd /tmp/startup-source-snapshot
  sed 's#  /repo/#  #' /output/results/c-source.sha256 | sha256sum -c --quiet
 )
 cp /output/results/c-source.sha256 /tmp/startup-source-snapshot/manifest.sha256
 cp /tmp/startup-build/go.mod /tmp/startup-source-snapshot/module-go.mod
 tar -C /tmp/startup-source-snapshot -czf "$SU_SNAPSHOT" .
 sha256sum "$SU_SNAPSHOT" > "$SU_SNAPSHOT.sha256"
 SU_DIGEST="$(sha256sum /output/results/c-source.sha256)";SU_DIGEST="${SU_DIGEST%% *}"
 printf '%s\n' "$SU_DIGEST" > /output/results/compiled-source-digest.txt
 SU_DEFINE="-DKI_STARTUP_SOURCE_DIGEST=0x${SU_DIGEST:0:16}"
 go version > /output/results/tool-versions.txt;gcc --version >> /output/results/tool-versions.txt
 cd /tmp/startup-build
 SU_TAGS=matching_startup,matching_glyph,matching_vga
 if [[ "$SU_MODE" == smoke ]]; then
  CGO_CFLAGS="-O2 -ffunction-sections -fdata-sections -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $SU_DEFINE" go build -tags "$SU_TAGS" -p 2 -o /output/startup-smoke .
  SU_SMOKE_ARGS=(-smoke -out /output/results/smoke.json)
  if [[ -n "$SU_GROUP" ]]; then SU_SMOKE_ARGS+=(-group "$SU_GROUP");fi
  /output/startup-smoke "${SU_SMOKE_ARGS[@]}"
  go version -m /output/startup-smoke > /output/results/buildinfo-smoke.txt
 else
  if [[ "$SU_MODE" != mutants ]]; then
   SU_OPTIMIZATIONS=(O2 O0)
   case "$SU_MODE" in normal-o2) SU_OPTIMIZATIONS=(O2) ;; normal-o0) SU_OPTIMIZATIONS=(O0) ;; esac
   for optimize in "${SU_OPTIMIZATIONS[@]}"; do
    CGO_CFLAGS="-$optimize -ffunction-sections -fdata-sections -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $SU_DEFINE" go build -tags "$SU_TAGS" -p 2 -o "/output/startup-$optimize" .
    go version -m "/output/startup-$optimize" > "/output/results/buildinfo-$optimize.txt"
    "/output/startup-$optimize" -out "/output/results/$optimize.json"
   done
  fi
  if [[ "$SU_MODE" == full || "$SU_MODE" == mutants ]]; then
   test -s /repo/tools/c_recovery_startup_mutants.tsv
   cp /repo/tools/c_recovery_startup_mutants.tsv /output/results/mutant-names.tsv
   while IFS=$'\t' read -r mutation group name; do
    [[ "$mutation" =~ ^[1-9][0-9]*$ ]] || exit 2
    CGO_CFLAGS="-O0 -ffunction-sections -fdata-sections -std=c11 -Werror=implicit-function-declaration -D_GNU_SOURCE $SU_DEFINE -DKI_STARTUP_MUTATION=$mutation" go build -tags "$SU_TAGS" -p 2 -o "/output/mutant-$mutation" .
    go version -m "/output/mutant-$mutation" > "/output/results/buildinfo-mutant-$mutation.txt"
    rm -f "/output/results/mutant-$mutation.json"
    set +e
    "/output/mutant-$mutation" -group "$group" -out "/output/results/mutant-$mutation.json" > "/output/results/mutant-$mutation.log" 2>&1
    status=$?
    set -e
    printf '%s\n' "$status" > "/output/results/mutant-$mutation.exit"
    cat "/output/results/mutant-$mutation.log"
    test "$status" -eq 1;test -s "/output/results/mutant-$mutation.json"
   done < /repo/tools/c_recovery_startup_mutants.tsv
  fi
 fi
 find /golem/internal -type f -name '*.go' -print0 | sort -z | xargs -0 sha256sum > /output/results/golem-source-after.sha256
 cmp /output/results/golem-source-before.sha256 /output/results/golem-source-after.sha256
}
if [[ "${1:-}" == --payload ]]; then shift;test "$#" -eq 0;startup_payload;exit 0;fi
test "$#" -eq 0
[[ "${SU_JOB_ID:-}" =~ ^[a-zA-Z0-9_.-]+$ ]];[[ "${SU_TIMEOUT:-}" =~ ^[1-9][0-9]*$ ]]
test -O /output/results
SU_LOG="/output/results/$SU_JOB_ID.job.log";SU_EXIT="/output/results/$SU_JOB_ID.job.exit"
test ! -e "$SU_LOG";test ! -e "$SU_EXIT"
record_exit() { SU_STATUS=$?;trap - EXIT;printf '%s\n' "$SU_STATUS" > "$SU_EXIT.tmp";mv "$SU_EXIT.tmp" "$SU_EXIT"; }
trap record_exit EXIT
trap 'exit 143' TERM
trap 'exit 130' INT
set +e
timeout --kill-after=10 "$SU_TIMEOUT" /bin/bash "$0" --payload 2>&1 | tee "$SU_LOG"
SU_PIPE_STATUS=("${PIPESTATUS[@]}")
set -e
SU_STATUS="${SU_PIPE_STATUS[0]}"
if [[ "$SU_STATUS" -eq 0 && "${SU_PIPE_STATUS[1]}" -ne 0 ]];then SU_STATUS="${SU_PIPE_STATUS[1]}";fi
exit "$SU_STATUS"
