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

# ⭐ **時鐘是兩邊唯一無歧義的錨點**（CLAUDE.md §4.02）。
# 區塊 `+0x00` 日、`+0x02` 子刻、`+0x03` 時、`+0x04` 月、`+0x06` 年。
clockof() {
    od -An -tu1 -j 0 -N8 "$1" |
        awk '{printf "%d/%d/%d/%d/%d\n", $7 + $8 * 256, $5, $1, $4, $3}'
}

start=$(cursor "$SNAP")
echo "快照起始據點游標：$start"
fail=0
for t in "${ticks[@]}"; do
    orig="$CK/orig-ck$t.DAT"
    [[ -f "$orig" ]] || { echo "跳過拍 $t：找不到 $orig"; continue; }
    cur=$(cursor "$orig")
    # ⛔ **不要用標稱拍數反推。** 據點游標 192 循環，靠標稱值鎖相位只在
    #    ±96 內有效——而「停下來等玩家」會讓 `拍 × 24509` 差好幾百拍
    #    （召見的回應期間指令照跑而時鐘不動，之後又繼續跑）。
    #    時鐘沒有這個問題：讀原版檢查點的時刻，讓 remake 跑到同一刻。
    clk=$(clockof "$orig")
    tools/go.sh run tools/rng_pace.go -save "$SNAP" -rng-state "$RNG" \
        -until "$clk" -save-out "$CK/remake-$t.DAT" > "$CK/remake-$t.log" 2>&1
    n=$(sed -n 's/.*共 \([0-9]*\) 拍/\1/p' "$CK/remake-$t.log" | tail -1)
    [[ -n "$n" ]] || { echo "拍 $t：remake 跑不到 $clk"; fail=1; continue; }
    # 交叉檢查：跑到同一時刻時據點游標也該相同。
    rc=$(cursor "$CK/remake-$t.DAT")
    (( rc == cur )) || echo "  ⚠ 游標對不上：原版 $cur、remake $rc"
    c=$(tools/py.sh tools/city_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    p=$(tools/py.sh tools/corps_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    f=$(tools/py.sh tools/faction_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    g=$(tools/py.sh tools/global_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    e=$(tools/py.sh tools/event_diff.py "$orig" "$CK/remake-$t.DAT" | grep 合計)
    # ⭐ 佇列**也判**：`tools/orig_snapshot.py` 已經走 `peek:2514` 把
    #    事件佇列（`+0x52C0`）那 1,024 B 從執行期記憶體讀回來，不再是
    #    來源 `SAVE.DAT` 的模板（docs/playtest/119 §44）。⚠ 用舊快照或
    #    舊檢查點跑會在這一欄看到假差異——重取一次再判。
    printf '拍 %-5s 游標 %-4s 跑 %-5s 據點 %-18s 勢力 %-20s 軍團 %-34s 全域 %-18s 佇列 %s\n' \
        "$t" "$cur" "$n" "$c" "$f" "$p" "${g%%（*}" "${e%%，*}"
    [[ "$c$f$p$g" == *"0 個 byte／0 座"*"0 個 byte／0 個勢力"*"0 個 byte／0 支"* \
        && "$g" == "合計 0 個 byte"* && "$e" == "合計 0 筆"* ]] || fail=1
done
exit $fail
