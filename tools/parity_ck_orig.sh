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

# ⭐ **君主召見的回應**：跑到 196 年 5 月 13 日 11 時原版會停下來等玩家
# （「孫乾　大人，主公有事召見。」），而**它問的是右鍵**——
# `int 33h AX=5, BX=1`，左鍵一次都不會被讀到（docs/playtest/118 §2）。
#
# 序列照 `docs/playtest/118` §3，只是選第 1 列而不是第 3 列：
# `press` 直接按下反白那一列 ＝「無條件同意」（使用者裁定 2026-09-10
# 「都固定 yes」；remake 側對應 `rng_pace.go` 的 `-answer accept`，那是預設）。
#
# ⚠ 三件事：
#   · **只有 `steps` 超過召見點才插**，否則右鍵會落在大地圖上開選單。
#   · **回應期間遊戲是停的**（時鐘不動）而指令照跑，所以 budget 要多給
#     那 6M；召見之後接著跑 `steps − SUMMON_STEPS` 道，`拍 × 24509`
#     的估計因此不受影響。
#   · 90 天內**只有這一次召見**（`docs/playtest/118` §5）。
SUMMON_STEPS=${WOLONG_SUMMON_STEPS:-147054000}
# 召見在 `SUMMON_STEPS` **之前**就發生了，中間那一段遊戲是停的。
# 由實測校準（`docs/playtest/119` §46.21）：跑到 `SUMMON_STEPS` 時
# 遊戲時鐘停在 5/13 11 時、據點游標 144。
SUMMON_WASTE=${WOLONG_SUMMON_WASTE:-0}
SUMMON_REPLY="rclick:287,199;steps:500000;rclick:287,199;steps:500000"
SUMMON_REPLY="$SUMMON_REPLY;press;steps:2000000"
for _ in 1 2 3 4 5 6 7 8; do
    SUMMON_REPLY="$SUMMON_REPLY;rclick:287,199;steps:500000"
done

for t in "$@"; do
    steps=$((t * STEPS_PER_TICK))
    if (( steps >= SUMMON_STEPS )); then
        # ⚠ **跑到 `SUMMON_STEPS` 時遊戲已經卡住一陣子了**，那一段的指令
        # 沒有推進時鐘。`SUMMON_WASTE` 是實測出來的浪費量，要補回去，
        # 否則「拍 × 24509」會低估，`parity_ck.sh` 的游標反推（±96）鎖錯相位。
        RUN="steps:$SUMMON_STEPS;$SUMMON_REPLY"
        RUN="$RUN;steps:$((steps - SUMMON_STEPS + SUMMON_WASTE))"
        echo "=== 拍 $t（steps:$steps，含召見回應）==="
    else
        RUN="steps:$steps"
        echo "=== 拍 $t（steps:$steps）==="
    fi
    WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-liubei13 \
    WOLONG_DOSGOLEM_TIMEOUT=${WOLONG_DOSGOLEM_TIMEOUT:-60m} \
    WOLONG_DOSGOLEM_BUDGET=${WOLONG_DOSGOLEM_BUDGET:-$((steps + 30000000))} \
        tools/dosgolem.sh "$CK" \
        "wait;click:320,200;click:300,151;$RUN;clock;$PEEKS" \
        > "$CK/ck$t.log" 2>&1
    grep "遊戲時鐘" "$CK/ck$t.log" | tail -1
    tools/py.sh tools/orig_snapshot.py "$CK/ck$t.log" \
        workplace/orig/dosv/SAVE.DAT "$CK/orig-ck$t.DAT" | tail -1
    od -An -tu1 -j $((0x2E)) -N2 "$CK/orig-ck$t.DAT" |
        awk '{printf "  據點游標 %d\n", ($1 + $2 * 256) / 32}'
done
