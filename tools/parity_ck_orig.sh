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
# ⛔ **不要用 `rclick`**：它的 settle 預設**六百萬道指令**（dosgolem 的
# `click` 註解），對即時制是二十幾個遊戲節拍——十次 rclick 就白跑十一天，
# 取樣點完全失控。`rclick` 還**忽略**第三個參數（只有 `click` 吃 settle）。
#
# ⇒ `move` 一次把游標放到對話框上，之後用 `rpress`（原地按右鍵）／
#   `press`（原地按左鍵）＋ 自己給的 `steps:`，每一步都可控。
#
# ⭐ **次數是量出來的，不能「多按幾次保險」**：對話關掉之後多的右鍵會在
#   大地圖上開選單並把遊戲暫停，之後的 `until:` 永遠跑不到。實測序列
#   （`workplace/parity/summon/a0`–`a7`，2026-09-10）：
#
#   | 按鍵 | 畫面 |
#   |---|---|
#   | （召見點）| 「孫乾　大人，主公有事召見。」|
#   | rpress ×2 | 推進到「曹操　的使者前來希望我國協助一事，我想聽聽你的意見。」|
#   | press | 選反白的第 1 列「為今後的外交設想，或許無條件比較好吧」＝ **無條件同意協力** |
#   | rpress ×4 | 逐段推完回覆與「已經不能再與呂布　共存了。立即固守國境。」，最後一下回到大地圖 |
#
#   判準是**時鐘**：a0–a6 都停在 5/13 11 時，第 7 下之後才開始走。
SUMMON_REPLY="move:287,199;rpress;steps:100000;rpress;steps:100000"
SUMMON_REPLY="$SUMMON_REPLY;press;steps:200000"
for _ in 1 2 3 4; do
    SUMMON_REPLY="$SUMMON_REPLY;rpress;steps:50000"
done

# ⭐ **參數含 `/` 就當遊戲日期**（`196/5/14`），走 dosgolem 的 `until:`。
#
# 召見之後 `steps` 與遊戲時間的對應整個垮掉：實測回應完之後 120 萬道指令
# 就走了 4 天 15 小時（≈ 每拍 1,050 道，正常是 24,509）——**對話關掉之後
# 遊戲不用重畫，跑得快得多**。日期取樣沒有這個問題。
for t in "$@"; do
    if [[ "$t" == */* ]]; then
        RUN="steps:$SUMMON_STEPS;$SUMMON_REPLY;until:$t"
        echo "=== $t（召見回應後跑到這一天）==="
        WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-liubei13 \
        WOLONG_DOSGOLEM_TIMEOUT=${WOLONG_DOSGOLEM_TIMEOUT:-60m} \
        WOLONG_DOSGOLEM_BUDGET=${WOLONG_DOSGOLEM_BUDGET:-200000000} \
            tools/dosgolem.sh "$CK" \
            "wait;click:320,200;click:300,151;$RUN;clock;$PEEKS" \
            > "$CK/ck${t//\//-}.log" 2>&1
        grep "遊戲時鐘" "$CK/ck${t//\//-}.log" | tail -1
        tools/py.sh tools/orig_snapshot.py "$CK/ck${t//\//-}.log" \
            workplace/orig/dosv/SAVE.DAT "$CK/orig-ck${t//\//-}.DAT" | tail -1
        od -An -tu1 -j $((0x2E)) -N2 "$CK/orig-ck${t//\//-}.DAT" |
            awk '{printf "  據點游標 %d\n", ($1 + $2 * 256) / 32}'
        continue
    fi
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
