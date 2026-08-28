#!/usr/bin/env bash
# 錄推廣主片要用的動態段（docs/spec/71）。
#
#   tools/promo_live_capture.sh
#
# ⭐ **主片以前每一段都是靜止截圖**——畫面是 remake 真的畫的，但影片裡沒有
# 一格在動。這一支錄的是真的在跑的遊戲：時鐘走、部隊移動、兩軍交戰。
#
# ⚠ **戰場要先 `-battle-steps` 推進再錄**。開場那一幀是雙方對白，兩軍還沒接觸，
# 直接錄會得到一段幾乎不動的畫面（量過：240 張裡只有 4 張不一樣）。
# ⚠ **推太多也不行**：攻城戰推到 200 步時已經打完了（270 張裡只剩 12 張不一樣，
# 而且那 12 張是循環的動畫，不是戰況）。掃過 0／40／90／120／150／170 之後
# **120 步是尖峰**（180 張裡 52 張不一樣）。
# ⚠ **戰術鏡頭要指定，而且會隨規則改動失效**。這兩個參數綁的是「這一局在
# 這個時間點打到哪裡」，規則層一動就要重掃——症狀是錄出一段空景，
# 不會有任何錯誤訊息。2026-08-26 重掃的結果：野戰維持原設定，
# 攻城改成 `-battle-steps 120 -battle-cam 44,8`（城內，騎兵與步兵交錯）。
# ⚠ **大地圖要用最高速檔**（`-speed 0`）。低速檔下八秒只走幾天，畫面幾乎不變
#（量過：240 張裡只有 16 張不一樣）；最高速檔才看得到時鐘、季節與行軍在動。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"
LIVE="${WOLONG_PROMO_LIVE:-workplace/promo/live}"

# ⚠ **語言段是四次啟動，不是按 F9**。F9 切出來的畫面與 `-lang` 啟動逐像素相同
#（`docs/playtest/46`），但按鍵注入沒辦法逐幀重現——四次啟動同一顆種子、
# 同一個時間點，拍到的是同一個局面的四種語言。
# 這四段先前是手打的，命令沒有留下來；寫在這裡是為了讓「重錄」不必再猜一次。
#
# 名稱|張數|不重複下限|參數
#
# ⚠ **不重複下限對語言段要放寬。** 那四段的內容是「同一個局面的四種語言」，
# 會動的只有時鐘；而「130 張裡有幾張不一樣」還跟**當下的機器負載**有關
#（量到過：同一條命令在閒置的機器上 26 張、在 load 90 的機器上 17 張）。
# 三段動態段維持 20，語言段用 8——那個值仍然擋得住「錄成一張靜止畫面」，
# 但不會因為機器忙就把一段可用的素材判死。
segs=(
"map|330|20|-direct -scenario 0 -player 0 -seed 17 -speed 0 -open-window -2"
"battle|390|20|-direct -scenario 0 -player 0 -seed 17 -open-battle -battle-steps 200 -battle-cam 20,15 -tactical-speed 2"
"siege|270|20|-direct -scenario 0 -player 0 -seed 17 -open-siege -battle-steps 120 -battle-cam 44,8 -tactical-speed 2"
"lang-zh-hant|130|8|-direct -scenario 0 -player 0 -seed 17 -speed 0 -open-window -2 -lang zh-hant"
"lang-zh-hans|130|8|-direct -scenario 0 -player 0 -seed 17 -speed 0 -open-window -2 -lang zh-hans"
"lang-ja|130|8|-direct -scenario 0 -player 0 -seed 17 -speed 0 -open-window -2 -lang ja"
"lang-en|130|8|-direct -scenario 0 -player 0 -seed 17 -speed 0 -open-window -2 -lang en"
# 以下三段給原版實機對照片用（tools/promo_dosv_realmachine.sh），不進主預告。
# ⭐ **它們量到的不重複張數就是 1，而那是對的**：編成與行軍指示都是
# 非常駐視窗，開著的時候**原版的時鐘停住**（說明書 3.1，`docs/spec/13`），
# 所以畫面連日期都不會變；啟動殼層更是還沒開始跑。下限設 1 ＝ 只擋「一張都沒錄到」。
"newgame|90|1|-seed 17"
"form|90|1|-direct -scenario 0 -player 0 -seed 17 -open-form"
"march|90|1|-direct -scenario 0 -player 0 -seed 17 -open-march-mode"
)

for s in "${segs[@]}"; do
    name=${s%%|*}; rest=${s#*|}
    n=${rest%%|*}; rest=${rest#*|}
    floor=${rest%%|*}; args=${rest#*|}
    echo "── 錄 $name（$n 張）──"
    tools/wlgame_frames.sh "$LIVE/$name" "$n" $args >/dev/null
    uniq=$(sha256sum "$LIVE/$name"/f*.png | cut -d' ' -f1 | sort -u | wc -l)
    echo "   $name：$n 張，其中 $uniq 張不重複（下限 $floor）"
    # 不重複的張數是「真的在動」的下限。太低就是錄到一段靜止畫面——
    # 那正是這支腳本要取代的東西，所以在這裡擋，不要等剪完才發現。
    [ "$uniq" -ge "$floor" ] || { echo "⚠ $name 幾乎沒有動靜（$uniq 張）" >&2; exit 1; }
done
echo "動態素材 → $LIVE"
