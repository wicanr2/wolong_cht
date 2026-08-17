#!/usr/bin/env bash
# 重拍 README 引用的 remake 截圖（`docs/images/`）。
#
#   tools/readme_shots.sh            # 全部
#   tools/readme_shots.sh advise     # 只拍名字含 advise 的那幾張
#
# ## 為什麼要有這一支
#
# 這些圖是**渲染結果**，凡是動到呈現層或素材解碼的改動都會讓它們過期，
# 而過期的截圖看起來跟正確的一模一樣。以前重拍要照著 commit message
# 裡的清單一張一張手打——**沒有腳本的清單等於下一輪要重新考古**。
#
# 每一張都固定 `-seed 17`，靠 `-open-*` fixture 進到該畫面。
# **沒有 fixture 的畫面不要放進 README**：手動按鍵序列重拍不出來。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"
filter=${1:-}

# 名稱|程式|參數
shots=(
"wlgame-dosv-natural-remake|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-window -2"
"wlgame-cht|wlgame|-direct -scenario 1 -player 0 -seed 17 -open-window -2"
"wlgame-cht-paused|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-window 3"
"wlgame-advise|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-advise -advise-menu"
"wlgame-form|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-form"
"wlgame-march|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-march-list"
"wlgame-list|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-list"
"wlgame-corps|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-corps"
"wlgame-battle|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-battle"
"wlgame-siege|wlgame|-direct -scenario 0 -player 0 -seed 17 -open-siege"
"wlgame-save-ui|wlgame|KEYS=4,s -direct -scenario 0 -player 0 -seed 17 -save-file workplace/scratch/SAVE.DAT"
"wlview-world-spring|wlview|-world -season 0 -wx 150 -wy 90"
"wlview-world-winter|wlview|-world -season 3 -wx 150 -wy 90"
)

# wlgame-ai-normal-encounter.png 刻意不在清單裡：它是 docs/playtest/08 那一次
# 「196/6/28 呂布對曹操」的歷史證據，現在的 AI 不會走到同一個狀態，
# 重拍只會得到一張空的主畫面。wlview-kyogrf.png 同理靠手動選素材。

mkdir -p workplace/scratch
for row in "${shots[@]}"; do
    name=${row%%|*}
    rest=${row#*|}
    cmd=${rest%%|*}
    args=${rest#*|}
    if [ -n "$filter" ] && [[ "$name" != *"$filter"* ]]; then
        continue
    fi
    echo "[readme-shots] $name（$cmd $args）"
    # shellcheck disable=SC2086
    WOLONG_SHOT_CMD="$cmd" tools/shot.sh "docs/images/$name.png" $args >/dev/null
done
echo "[readme-shots] 完成。用 git diff --stat docs/images 看哪幾張真的變了。"
