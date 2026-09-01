#!/usr/bin/env bash
# 在 wolong-go Docker 容器中錄取 wlgame 的實際 X11 遊玩畫面。
#
# 正常模式只接受鍵盤輸入；戰術對拍可用明確標示的 fixture
# 進入攻城或野戰。輸出為 640x400 point 縮放的 PNG 逐幀序列。
#
# WOLONG_WLGAME_TIMELINE：
#   wait:秒;shot:名稱;key:鍵;record:秒,fps,前綴
set -euo pipefail

# demonwinter-go 的 Go 安裝在固定的 /usr/local/go/bin；容器預設 PATH 沒有它。
export PATH=/usr/local/go/bin:$PATH

src_dir=${WOLONG_WLGAME_SOURCE:-/src}
orig_dir=${WOLONG_WLGAME_ORIG:-/orig}
font_dir=${WOLONG_WLGAME_FONT:-/font}
out_dir=${WOLONG_WLGAME_OUT:-/out}
timeline=${WOLONG_WLGAME_TIMELINE:-"wait:2;shot:strategy"}
seed=${WOLONG_WLGAME_SEED:-17}
speed=${WOLONG_WLGAME_SPEED:-0}
mode=${WOLONG_WLGAME_MODE:-normal}
startup_wait=${WOLONG_WLGAME_STARTUP_WAIT:-2}
game_binary=${WOLONG_WLGAME_BINARY:-${WOLONG_WLGAME_PREBUILT:-}}
extra_args=${WOLONG_WLGAME_EXTRA_ARGS:-}
# ⚠ `-scenario`／`-player` 本身就是「直啟白名單」的旗標（launcher.go 的
# directStartFlagWasPassed），帶著它們會**跳過啟動殼層**。要拍新遊戲那四層
# （ＹＥＳ／ＮＯ → 劇本 → 勢力清單 → 君主卡）就把這個變數設成空字串。
scenario_args=${WOLONG_WLGAME_SCENARIO_ARGS--scenario 0 -player 0}

for dir in "$src_dir" "$orig_dir" "$font_dir"; do
    if [ ! -d "$dir" ]; then
        echo "找不到必要目錄：$dir" >&2
        exit 1
    fi
done
mkdir -p "$out_dir"
if [ ! -w "$out_dir" ]; then
    echo "輸出目錄不可寫：$out_dir" >&2
    exit 1
fi

case "$mode" in
    normal)
        mode_args=()
        ;;
    siege_fixture)
        # 僅供推廣片的戰術可視化鏡頭；不列作正常路徑 parity 證據。
        mode_args=(-open-siege)
        ;;
    field_fixture)
        # 與 siege_fixture 配對，只驗收兩種戰鬥的共用畫面骨架。
        mode_args=(-open-battle)
        ;;
    *)
        echo "不支援的 wlgame 擷取模式：$mode" >&2
        exit 1
        ;;
esac

trace_file="$out_dir/runtime-trace.txt"
trace() {
    printf '%s %s\n' "$(date +%s)" "$*" >> "$trace_file"
}
: > "$trace_file"
: > "$out_dir/manifest.txt"

XVFB=""
APP=""
cleanup() {
    if [ -n "$APP" ]; then
        kill "$APP" 2>/dev/null || true
        wait "$APP" 2>/dev/null || true
    fi
    if [ -n "$XVFB" ]; then
        kill "$XVFB" 2>/dev/null || true
        wait "$XVFB" 2>/dev/null || true
    fi
    cp /tmp/wlgame-live.log "$out_dir/wlgame.log" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

if [ -n "$game_binary" ]; then
    if [ ! -x "$game_binary" ]; then
        echo "指定的 wlgame 執行檔不存在或不可執行：$game_binary" >&2
        exit 1
    fi
    trace "use prebuilt wlgame=$game_binary"
    # AppImage 在容器裡沒有 FUSE，要先自解。`--appimage-extract-and-run` 會
    # 每次解一份到 /tmp；先解一次再直接跑解出來的 AppRun，省掉重複解壓。
    # AppImage 的魔數是 ELF 標頭第 8–10 byte 的 `41 49 02`（"AI\x02"），
    # 不是檔案裡有 "AppImage" 這個字——那個字要到 squashfs 裡才有。
    if [ "$(dd if="$game_binary" bs=1 skip=8 count=3 2>/dev/null | od -An -tx1 | tr -d ' \n')" = "414902" ]; then
        cp "$game_binary" /tmp/app.AppImage
        chmod +x /tmp/app.AppImage
        (cd /tmp && ./app.AppImage --appimage-extract >/tmp/appimage-extract.log 2>&1)
        game_binary=/tmp/squashfs-root/AppRun
        [ -x "$game_binary" ] || { echo 'AppImage 解不開' >&2; cat /tmp/appimage-extract.log >&2; exit 1; }
        trace "extracted AppImage -> $game_binary"
    fi
else
    trace "build wlgame"
    (cd "$src_dir" && go build -o /tmp/wlgame-live ./cmd/wlgame)
    game_binary=/tmp/wlgame-live
fi

Xvfb :99 -screen 0 1600x900x24 >/tmp/wlgame-live-xvfb.log 2>&1 &
XVFB=$!
for _ in $(seq 1 100); do
    if command -v xdpyinfo >/dev/null 2>&1; then
        xdpyinfo -display :99 >/dev/null 2>&1 && break
    else
        DISPLAY=:99 xdotool search --name '.*' >/dev/null 2>&1 && break
    fi
    sleep 0.1
done
if command -v xdpyinfo >/dev/null 2>&1; then
    display_ready() { xdpyinfo -display :99 >/dev/null 2>&1; }
else
    display_ready() { DISPLAY=:99 xdotool search --name '.*' >/dev/null 2>&1; }
fi
if ! display_ready; then
    cat /tmp/wlgame-live-xvfb.log >&2 || true
    echo 'Xvfb 未能啟動。' >&2
    exit 1
fi

trace "start wlgame seed=$seed speed=$speed"
# ⚠ 額外參數放在預設之後：Go 的 flag 套件同名旗標**後面的贏**，
# 所以 WOLONG_WLGAME_EXTRA_ARGS 可以覆寫 -scenario／-player／-speed。
# 不加引號是刻意的——要讓它按空白拆成多個參數。
LIBGL_ALWAYS_SOFTWARE=1 DISPLAY=:99 "$game_binary" \
    -orig "$orig_dir" -font "$font_dir" $scenario_args \
    -seed "$seed" -speed "$speed" "${mode_args[@]}" \
    $extra_args >/tmp/wlgame-live.log 2>&1 &
APP=$!
WID=""
for _ in $(seq 1 100); do
    WID=$(DISPLAY=:99 xdotool search --onlyvisible --name '臥龍傳' 2>/dev/null | head -1 || true)
    [ -n "$WID" ] && break
    sleep 0.2
done
if [ -z "$WID" ]; then
    cat /tmp/wlgame-live.log >&2 || true
    echo 'wlgame 視窗未出現。' >&2
    exit 1
fi
DISPLAY=:99 xdotool windowactivate --sync "$WID" 2>/dev/null || true
trace "found window wid=$WID"

window_geometry() {
    eval "$(DISPLAY=:99 xdotool getwindowgeometry --shell "$WID")"
}

shot() {
    local name=$1
    local root_png
    if ! kill -0 "$APP" 2>/dev/null; then
        echo 'wlgame 在擷取前已結束。' >&2
        cat /tmp/wlgame-live.log >&2 || true
        exit 1
    fi
    window_geometry
    : "${WIDTH:?未能取得 wlgame 視窗寬度}"
    : "${HEIGHT:?未能取得 wlgame 視窗高度}"
    : "${X:?未能取得 wlgame 視窗 X 座標}"
    : "${Y:?未能取得 wlgame 視窗 Y 座標}"
    root_png=$(mktemp /tmp/wlgame-live-root.XXXXXX.png)
    DISPLAY=:99 import -window root "$root_png"
    # wlgame 的 2× 視窗輸出轉回 640×400 邏輯畫布，以利和 DOS/V 直接並排。
    convert "$root_png" -crop "${WIDTH}x${HEIGHT}+${X}+${Y}" +repage \
        -filter point -resize 640x400! "$out_dir/$name.png"
    rm -f "$root_png"
    printf '%s\n' "$name.png" >> "$out_dir/manifest.txt"
    trace "shot name=$name geometry=${WIDTH}x${HEIGHT}+${X}+${Y}"
}

record() {
    local seconds=$1
    local fps=$2
    local prefix=$3
    local frames delay n file
    frames=$(awk -v seconds="$seconds" -v fps="$fps" 'BEGIN { printf "%d", seconds * fps }')
    delay=$(awk -v fps="$fps" 'BEGIN { printf "%.6f", 1 / fps }')
    if [ "$frames" -le 0 ]; then
        echo "錄影段長度或 fps 無效：$seconds,$fps" >&2
        exit 1
    fi
    for ((n = 1; n <= frames; n++)); do
        printf -v file '%s-%06d' "$prefix" "$n"
        shot "$file"
        sleep "$delay"
    done
}

{
    printf 'source=remake actual GUI capture\n'
    printf 'seed=%s\nspeed=%s\nmode=%s\n' "$seed" "$speed" "$mode"
    printf 'timeline=%s\n' "$timeline"
} > "$out_dir/capture-metadata.txt"

sleep "$startup_wait"
IFS=';' read -r -a steps <<< "$timeline"
for step in "${steps[@]}"; do
    kind=${step%%:*}
    arg=${step#*:}
    case "$kind" in
        wait)
            trace "step begin=$step"
            sleep "$arg"
            trace "step end=$step"
            ;;
        shot)
            trace "step begin=$step"
            shot "$arg"
            ;;
        key)
            trace "step begin=$step"
            DISPLAY=:99 xdotool key --window "$WID" --clearmodifiers "$arg"
            sleep 0.3
            trace "step end=$step"
            ;;
        # keyrep:鍵,次數 —— 清單要往下捲十幾列時，寫成十幾個 key: 會讓
        # 時間軸長到看不出結構。
        keyrep)
            trace "step begin=$step"
            rk=${arg%%,*}; rn=${arg##*,}
            for ((ri = 0; ri < rn; ri++)); do
                DISPLAY=:99 xdotool key --window "$WID" --clearmodifiers "$rk"
                sleep 0.12
            done
            sleep 0.2
            trace "step end=$step"
            ;;
        # click／rclick：x,y 是**邏輯 640×400 座標**，會依視窗實際大小換算。
        # 主畫面的四個視窗開關只有滑鼠路徑（原版就沒有鍵盤，docs/re/22 §2），
        # 少了這一步就拍不到「關掉視窗」的畫面。
        click|rclick)
            trace "step begin=$step"
            window_geometry
            cx=${arg%%,*}; cy=${arg##*,}
            btn=1; [ "$kind" = rclick ] && btn=3
            DISPLAY=:99 xdotool mousemove \
                "$((X + cx * WIDTH / 640))" "$((Y + cy * HEIGHT / 400))"
            sleep 0.15
            DISPLAY=:99 xdotool click "$btn"
            sleep 0.3
            trace "step end=$step"
            ;;
        record)
            trace "step begin=$step"
            IFS=',' read -r seconds fps prefix <<< "$arg"
            record "$seconds" "$fps" "$prefix"
            trace "step end=$step"
            ;;
        *)
            echo "不支援的擷取步驟：$step" >&2
            exit 1
            ;;
    esac
done
