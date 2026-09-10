#!/usr/bin/env bash
# 取原版側的逐欄檢查點（docs/playtest/119 §36.5）。
#
#   tools/parity_ck_orig.sh 2454 2500      # 取這兩拍
#
# 停點是 `steps:拍 × 24509`，而 24,509 只是每拍指令數的估計——**實際落點會
# 偏個幾拍**。所以檔名上的拍數只是標稱值，真正的相位由區塊 `+0x2E`
# （據點游標）決定，remake 側一律用 `tools/parity_ck.sh` 反推，不要手打拍數。
#
# ⛔ **`steps:` 停在指令級，可能落在一拍中間**，而 remake 側是跑完整拍。
# 游標對得上只代表**同一拍**，不代表**同一個位置**。平常看不出來——
# 一拍內的狀態變動小；但那一拍有戰鬥時，「拍中間」與「拍結束」會差到
# 十幾個 byte（軍團兵力、`+0x03` 倒數、子刻、軍團游標都會對不上），
# **形狀與真的規則分歧一模一樣**（2026-09-10 拍 5,020，docs/playtest/119 §46.6）。
# 判準：`tools/py.sh tools/parity_pace_diff.py … --tick N` 看那一拍有沒有
# `152F6`／`151FC`／`1521B`；有就換一個沒戰鬥的拍再取一次。
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

_LESSONS="$(pwd)/docs/lessons/parity.txt"
[[ -s "$_LESSONS" ]] && cat "$_LESSONS" >&2

CK=${WOLONG_CK_DIR:-workplace/parity/pace/ck16}
STEPS_PER_TICK=${WOLONG_STEPS_PER_TICK:-24509}
mkdir -p "$CK"

PEEKS="ipeek:10CF0:128"
for off in 0000 0800 1000 1800 2000 2800 3000 3800 4000 4800; do
    PEEKS="$PEEKS;peek:2754:$off:2048"
done
PEEKS="$PEEKS;peek:2754:5000:544"
# ⚠ 事件佇列在**第三個段**（`cs:word_10D56`，實測 0x2514），區塊的最後
# 1,024 B。少了它，檢查點那一段是從來源 `SAVE.DAT` 模板複製的
# （docs/playtest/119 §44）。
PEEKS="$PEEKS;peek:2514:0000:1024"

for t in "$@"; do
    steps=$((t * STEPS_PER_TICK))
    echo "=== 拍 $t（steps:$steps）==="
    WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-liubei13 \
    WOLONG_DOSGOLEM_TIMEOUT=${WOLONG_DOSGOLEM_TIMEOUT:-60m} \
        tools/dosgolem.sh "$CK" \
        "wait;click:320,200;click:300,151;steps:$steps;clock;$PEEKS" \
        > "$CK/ck$t.log" 2>&1
    grep "遊戲時鐘" "$CK/ck$t.log" | tail -1
    tools/py.sh tools/orig_snapshot.py "$CK/ck$t.log" \
        workplace/orig/dosv/SAVE.DAT "$CK/orig-ck$t.DAT" | tail -1
    od -An -tu1 -j $((0x2E)) -N2 "$CK/orig-ck$t.DAT" |
        awk '{printf "  據點游標 %d\n", ($1 + $2 * 256) / 32}'
done
