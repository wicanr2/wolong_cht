#!/usr/bin/env bash
# 重生 remake 側的逐欄檢查點並與原版比對（docs/playtest/119 §37）。
#
# ⭐ **拍數由據點游標反推，不用檔名上的標稱拍數。**
# 原版側的檢查點是用 `steps:拍 × 24509` 停下來的，而 24,509 只是每拍指令數的
# 估計，誤差會累積——`orig-ck700.DAT` 實際落在第 702 拍。拿 700 去跑 remake
# 會錯位兩拍，症狀是「游標前一兩座據點的上昇值與防災值差 1」，
# 看起來像規則差異，其實只是取樣點沒對齊。
#
# 游標每拍 +1、192 循環，所以它把相位釘死：同一個游標值在 ±96 拍內唯一。
#
#   tools/parity_ck.sh [檢查點...]        # 預設 200 700 1200 1600 1900
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

_LESSONS="$(pwd)/docs/lessons/parity.txt"
[[ -s "$_LESSONS" ]] && cat "$_LESSONS" >&2

CK=${WOLONG_CK_DIR:-workplace/parity/pace/ck16}
SNAP=${WOLONG_CK_SNAP:-workplace/parity/orig-at-16h.DAT}
RNG=${WOLONG_CK_RNG:-workplace/parity/liubei-rng16.bin}
CURSOR_OFF=$((0x2E))
PERIOD=192

ticks=("$@")
[[ ${#ticks[@]} -eq 0 ]] && ticks=(200 700 1200 1600 1900)

cursor() {  # 讀一份區塊的據點游標
    od -An -tu1 -j "$CURSOR_OFF" -N2 "$1" | awk '{print ($1 + $2 * 256) / 32}'
}

start=$(cursor "$SNAP")
echo "快照起始據點游標：$start"
fail=0
for t in "${ticks[@]}"; do
    orig="$CK/orig-ck$t.DAT"
    [[ -f "$orig" ]] || { echo "跳過拍 $t：找不到 $orig"; continue; }
    cur=$(cursor "$orig")
    n=$(( (cur - start + PERIOD * 100) % PERIOD ))
    while (( n < t - PERIOD / 2 )); do n=$(( n + PERIOD )); done

    tools/go.sh run tools/rng_pace.go -save "$SNAP" -rng-state "$RNG" \
        -ticks "$n" -save-out "$CK/remake-$t.DAT" >/dev/null
    c=$(tools/py.sh tools/city_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    p=$(tools/py.sh tools/corps_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    f=$(tools/py.sh tools/faction_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    g=$(tools/py.sh tools/global_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    e=$(tools/py.sh tools/event_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    # ⚠ 佇列**只印不判**：`tools/orig_snapshot.py` 只拼「全域 ＋ 四張表」，
    #    事件佇列（`+0x52C0`）那一段是從原版 `SAVE.DAT` **模板**來的，
    #    不是執行期的內容——比它等於拿模板對 remake 的活狀態
    #    （worklist `event-queue-not-peeked`）。要判就得先 peek 那一段。
    printf '拍 %-5s 游標 %-4s 跑 %-5s 據點 %-18s 勢力 %-20s 軍團 %-34s 全域 %-18s 佇列 %s（原版側是模板）\n' \
        "$t" "$cur" "$n" "$c" "$f" "$p" "${g%%（*}" "${e%%，*}"
    [[ "$c$f$p$g" == *"0 個 byte／0 座"*"0 個 byte／0 個勢力"*"0 個 byte／0 支"* \
        && "$g" == "合計 0 個 byte"* ]] || fail=1
done
exit $fail
