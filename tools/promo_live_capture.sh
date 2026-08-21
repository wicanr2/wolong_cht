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
# ⚠ **戰術鏡頭要指定**。原版初值 36,14 看到的是城牆，部隊在畫面外；
# `20,15` 才對準交戰的中庭。
# ⚠ **大地圖要用最高速檔**（`-speed 0`）。低速檔下八秒只走幾天，畫面幾乎不變
#（量過：240 張裡只有 16 張不一樣）；最高速檔才看得到時鐘、季節與行軍在動。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"
LIVE="${WOLONG_PROMO_LIVE:-workplace/promo/live}"

# 名稱|張數|參數
segs=(
"map|330|-direct -scenario 0 -player 0 -seed 17 -speed 0 -open-window -2"
"battle|390|-direct -scenario 0 -player 0 -seed 17 -open-battle -battle-steps 200 -battle-cam 20,15 -tactical-speed 2"
"siege|270|-direct -scenario 0 -player 0 -seed 17 -open-siege -battle-steps 200 -battle-cam 20,15 -tactical-speed 2"
)

for s in "${segs[@]}"; do
    name=${s%%|*}; rest=${s#*|}; n=${rest%%|*}; args=${rest#*|}
    echo "── 錄 $name（$n 張）──"
    tools/wlgame_frames.sh "$LIVE/$name" "$n" $args >/dev/null
    uniq=$(sha256sum "$LIVE/$name"/f*.png | cut -d' ' -f1 | sort -u | wc -l)
    echo "   $name：$n 張，其中 $uniq 張不重複"
    # 不重複的張數是「真的在動」的下限。太低就是錄到一段靜止畫面——
    # 那正是這支腳本要取代的東西，所以在這裡擋，不要等剪完才發現。
    [ "$uniq" -ge 20 ] || { echo "⚠ $name 幾乎沒有動靜（$uniq 張）" >&2; exit 1; }
done
echo "動態素材 → $LIVE"
