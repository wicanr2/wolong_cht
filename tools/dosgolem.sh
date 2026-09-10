#!/usr/bin/env bash
# 用 dosgolem 跑原版當視覺 oracle——**不經過 DOSBox、不經過 X**。
#
#   tools/dosgolem.sh out_dir "wait;shot:z1;click:320,215;wait;shot:z2"
#
# 步驟語法與 `tools/dosv_live_capture.sh` 的時間軸同形，所以舊腳本搬得過來：
#
#   wait          畫面畫完並停住（不是「睡幾秒」）
#   steps:N       再跑 N 道指令
#   click:X,Y     點左鍵
#   rclick:X,Y    **點右鍵** ← 對話框在等的是它，不是左鍵（docs/playtest/118 §2）
#   press／rpress 原地按左／右鍵（要先 move 到位）
#
# ⛔ **`click`／`rclick` 的 settle 預設六百萬道指令**（≈ 二十幾個遊戲節拍），
#    而 `rclick` **忽略**第三個參數——只有 `click:X,Y,N` 吃得到 settle。
#    連按十次就白跑十一個遊戲日，取樣點完全失控。
#    ⇒ 要精確控制就 `move` 一次再用 `press`／`rpress` ＋ 自己給的 `steps:`。
#   move:X,Y      只移動游標
#   profile:N     取樣 N 道指令，印出 CS:IP 的熱區——**卡住時第一個要跑的**
#   sclick:X,Y    **進到遊戲之後**用這一支：畫面座標點擊（先把捲動原點歸零）
#   stap:X,Y      同上，瞬按
#   origin        印出捲動原點與遊戲算出來的畫面座標
#   hotspots      印出目前畫面的熱區圖（編號 → 像素矩形）
#   at:X,Y        印出某個像素座標上的熱區編號
#   tile:TX,TY    把游標移到大地圖的第 (TX,TY) 格 ← **選據點要用這個**
#   celltile      印出游標現在在第幾格
#   runto:LIN[,N] 跑到 CS:IP 走到 IDA 線性位址（等一場仗開打就用它）
#   ipeek:LIN:N   印出 IDA 線性位址起 N 個 byte（`cs:word_XXXX` 用這個）
#   siege:攻,據點 直接叫原版開一場指定的攻城戰（不必等 AI 打過來）
#   ticks:N       **以戰術節拍為單位**推進 N 拍 ← 走位期取樣要用這個，不要用 steps:
#   units[:all|sum|側/隊]  印場上的兵；`側/隊` 只印那一隊，`sum` 每側一行摘要
#   corps         印原版的軍團表
#   until:Y/M/D   跑到遊戲日期到某一天  ← 即時制的取樣點寫成日期
#   clock         印出目前的遊戲日期
#   shot:NAME     存一張 640×400 的 PNG（已經裁好，不必再跑 parity_crop.py）
#
# ⚠ **y 座標是遊戲座標，不是 DOSBox-X 的視窗座標。** 舊腳本要加 `-y dosbox`，
# 換算是 `遊戲 y ＝ 視窗 y × 399 ÷ 479`（分母 479 不是 480，
# 見 dosgolem 的 `apps/wolong/wolong.go`）。
#
# ⚠⚠ **進到遊戲之後畫面座標不等於滑鼠座標**：大地圖有一層捲動原點
# （`畫面 ＝ 滑鼠 − 原點`，推到負的就捲地圖），所以那時要用 `sclick`／`stap`。
# 選單畫面不走那一層，用 `click` 就好。細節見 docs/playtest/66。
#
# ⛔ **遊戲停住時先想右鍵。** 君主召見軍師那種對話框問的是
# `int 33h AX=5, BX=1`（右鍵按下計數），**左鍵一次都不會被讀到**——
# 四種左鍵操作全部無效而畫面完全正常（docs/playtest/118 §2）。
# 判準：`profile:20000000` 看熱區落在 `20000–200FF`（滑鼠分派器）。
#
# ⚠ **`wait` 在遊戲中不適用**：即時制的畫面永遠不會靜止，`wait` 會一路跑到
# 預算上限（實測跑掉兩個遊戲日）。遊戲中用 `steps:` 或 `until:`。
#
# 用別的存檔開局：另外準備一個目錄（原版檔案 ＋ 那份 SAVE.DAT），
# 再用 WOLONG_DOSGOLEM_GAMEDIR=<相對 workplace/ 的路徑> 指過去。
# **原版素材唯讀，不要就地覆蓋**（CLAUDE.md §9）。
#
# 規格：docs/spec/131-dosgolem-oracle.md，實跑紀錄：docs/playtest/65、66
set -euo pipefail
# ⭐ **教訓印在動手之前**（`docs/lessons.json`，`tools/lessons.py render` 產生）。
# 光把教訓寫進常駐文件不夠——它們要防的動作發生在任務中間，那時沒有東西會問。
_LESSONS="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/docs/lessons/parity.txt"
[[ -s "$_LESSONS" ]] && cat "$_LESSONS" >&2


REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOLEM="${WOLONG_DOSGOLEM:-$HOME/cht/dosgolem-wolong}"

if [[ ! -d "$GOLEM" ]]; then
    echo "找不到 dosgolem：$GOLEM" >&2
    echo "設 WOLONG_DOSGOLEM 指到工作副本，或 git clone https://github.com/wicanr2/dosgolem.git" >&2
    exit 1
fi

# 遊戲目錄可換：要載入某一份 SAVE.DAT 時，另外準備一個目錄
# （原版檔案 ＋ 那份存檔）再指過去。**原版素材唯讀，不要就地覆蓋**
# （CLAUDE.md §9）。路徑相對於 workplace/。
GAMEDIR="${WOLONG_DOSGOLEM_GAMEDIR:-orig/dosv}"

OUT="${1:?用法: tools/dosgolem.sh <輸出目錄> <時間軸> [dosbox]}"
TIMELINE="${2:?用法: tools/dosgolem.sh <輸出目錄> <時間軸> [dosbox]}"
YMODE="${3:-game}"

mkdir -p "$OUT"
OUT_ABS="$(cd "$OUT" && pwd)"

ARGS=(-exe "/orig/$GAMEDIR/KI.EXE" -root "/orig/$GAMEDIR" -dir /out -script "$TIMELINE")
[[ "$YMODE" == "dosbox" ]] && ARGS+=(-dosbox-y)
[[ -n "${WOLONG_DOSGOLEM_WATCH:-}" ]] && ARGS+=(-watch "$WOLONG_DOSGOLEM_WATCH")
[[ -n "${WOLONG_DOSGOLEM_BUDGET:-}" ]] && ARGS+=(-budget "$WOLONG_DOSGOLEM_BUDGET")

# 原版素材唯讀掛載；輸出目錄可寫。**本專案與 dosgolem 都不含原版檔案。**
DOSGOLEM_ORIG="$REPO_ROOT/workplace" \
DOSGOLEM_TIMEOUT="${WOLONG_DOSGOLEM_TIMEOUT:-30m}" \
DOSGOLEM_EXTRA_MOUNT="$OUT_ABS:/out" \
    "$GOLEM/tools/go.sh" run ./apps/wolong/cmd/shot "${ARGS[@]}"
