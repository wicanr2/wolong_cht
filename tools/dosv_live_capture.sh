#!/usr/bin/env bash
# 在一次性 DOSBox-X 容器內錄取松崗 DOS/V 的「實際顯示」畫面。
#
# 此腳本必須由 wolong-dosboxx:latest 容器內執行。原版目錄掛為 /orig:ro，
# 只會把 X11 實際畫面寫到 /out；它不修改原版資料，也不把原版檔案複製到輸出。
#
# 預設會在密碼頁按原始「確定」，接著依 WOLONG_DOSV_TIMELINE 執行：
#   wait:秒;shot:名稱;key:鍵;type:字串;click:x,y;rclick:x,y;press;
#   record:秒,fps,前綴;audio-start;audio-stop;savefile:檔名
#
# ⚠ **click／rclick 的 y 要用「遊戲座標 × 1.2」**：視窗是 640×480，
#    遊戲畫面 640×400 置中，而 INT 33 把**整個視窗**等比對映到遊戲畫面
#    （量法與證據見 tools/dosv_capture.sh 的說明）。減 40 會差十幾 px，
#    正好落進按鈕之間的空隙，看起來就跟「事件沒送到」一樣。
# ⚠ **一次按下會被兩層各吃一次。** 量到的：在指令列上送 `click;press` 之後，
#    「軍團」選單**開了，而且同一瞬間被選走第一列**（位置確認）。
#    成因是 INT 33 `AX=5` 回的 BX 是**當下的按鍵狀態**，不是計數——
#    `clickat`／`press_here` 把鍵按住 0.12 秒，這段時間內剛開起來的選單
#    再 poll 一次就看到「按著」，於是拿當時的游標位置選列
#    （`sub_193E9` 把框外的游標夾到列 0）。清單視窗不會，因為框外的點被忽略。
#    ⇒ **目前沒有可靠的辦法點到選單的第二列以後**。要繞開就別走那條路：
#    這一款的戰場可以用「編成一支軍團然後等 AI 來打」取得
#    （docs/playtest/40），不必碰「軍團 → 行軍指示」。
# record 產生前綴-000001.png … 的連續實機畫面，可交給影片容器編碼。
# audio-start／audio-stop 觸發 DOSBox-X 自身的 WAV mixer capture；這是原版
# YNSOUND.COM 經模擬硬體後的實際輸出，不是由 remake 或外部合成器重建。
set -euo pipefail

orig_dir=${WOLONG_DOSV_ORIG:-/orig}
out_dir=${WOLONG_DOSV_CAPTURE_OUT:-/out}
timeline=${WOLONG_DOSV_TIMELINE:-"wait:1;shot:after-confirm;wait:10;shot:opening-10;wait:15;shot:opening-25;wait:15;shot:opening-40"}
password_wait=${WOLONG_DOSV_PASSWORD_WAIT:-8}
password_timeout=${WOLONG_DOSV_PASSWORD_TIMEOUT:-45}
bypass_password=${WOLONG_DOSV_BYPASS_PASSWORD:-1}
debug_capture=${WOLONG_DOSV_DEBUG_CAPTURE:-0}
boot_mode=${WOLONG_DOSV_BOOT_MODE:-start}

if [ ! -d "$orig_dir" ]; then
    echo "找不到唯讀 DOS/V 原版目錄：$orig_dir" >&2
    exit 1
fi
mkdir -p "$out_dir"
if [ ! -w "$out_dir" ]; then
    echo "擷取輸出目錄不可寫：$out_dir" >&2
    exit 1
fi

case "$boot_mode" in
    start)
        boot_commands=(START.BAT)
        ;;
    logo)
        boot_commands=(YNFONT.EXE YNSOUND.COM LOGO.EXE)
        ;;
    audio-logo)
        # DOSBox-X 的 guest 內建 DX-CAPTURE 直接錄 mixer，避免依賴不同版本的
        # host 快捷鍵。LOGO.EXE 自然結束時 WAV 會被正確封口。
        boot_commands=(YNFONT.EXE "YNSOUND.COM /A" "DX-CAPTURE /A LOGO.EXE")
        ;;
    game)
        # 只使用原版自身的常駐驅動與主程式，供縮短開場後的實機策略錄影。
        # 這不是發行版捷徑，也不拿來證明完整自然啟動順序。
        boot_commands=(YNFONT.EXE YNSOUND.COM KI.EXE)
        ;;
    *)
        echo "不支援的 DOS/V 啟動模式：$boot_mode" >&2
        exit 1
        ;;
esac

trace_file="$out_dir/runtime-trace.txt"
trace() {
    printf '%s %s\n' "$(date +%s)" "$*" >> "$trace_file"
}
: > "$trace_file"
trace "start boot_mode=$boot_mode"

game_dir=$(mktemp -d /tmp/wolong-dosv-capture.XXXXXX)
log_file=/tmp/wolong-dosv-capture.log
XVFB=""
DB=""
cleanup() {
    # ⭐ 遊戲還在跑的時候送 SIGTERM，DOSBox-X 會彈一個 xmessage 確認框
    # （"Are you sure to quit anyway now?"）等一個永遠不會來的答案，
    # 於是 `wait` 掛住、**容器永遠不退出**——即使 --rm、即使工作早就做完。
    # 症狀是看起來完全正常：docker ps 一排同名容器、CPU 各 0.1%、輸出一張不少。
    # 兩道一起做：設定裡關掉確認框（quit warning=false），這裡五秒後直接砍。
    if [ -n "$DB" ]; then
        kill "$DB" 2>/dev/null || true
        for _ in $(seq 1 20); do kill -0 "$DB" 2>/dev/null || break; sleep 0.25; done
        kill -9 "$DB" 2>/dev/null || true
        wait "$DB" 2>/dev/null || true
    fi
    if [ -n "$XVFB" ]; then
        kill "$XVFB" 2>/dev/null || true
        for _ in $(seq 1 20); do kill -0 "$XVFB" 2>/dev/null || break; sleep 0.25; done
        kill -9 "$XVFB" 2>/dev/null || true
        wait "$XVFB" 2>/dev/null || true
    fi
    # 即使中途失敗也保留執行器日誌；它不是原版資產，能分辨遊戲自然轉場與
    # X11／capture 橋接失敗。
    cp "$log_file" "$out_dir/dosbox-x.log" 2>/dev/null || true
    chmod -R u+w "$game_dir" 2>/dev/null || true
    rm -rf "$game_dir"
}
trap cleanup EXIT INT TERM

# DOSBox-X 的 guest 會寫 SAVE.DAT；複製只存在於容器的明確暫存目錄。
cp -a "$orig_dir"/. "$game_dir"/
chmod -R u+w "$game_dir"
trace "copied original to container temporary game directory"

# surface 讓 ImageMagick 的 X11 擷取讀到實際 framebuffer；OpenGL 在
# headless Xvfb 下可能只回傳黑色視窗。
printf '%s\n' \
    '[sdl]' \
    'autolock=false' \
    'output=surface' \
    'mouse_emulation=integration' \
    'usesystemcursor=false' \
	    '[dosbox]' \
	    'machine=vgaonly' \
	    'memsize=16' \
	    'quit warning=false' \
	    "captures=$out_dir" \
	    'show recorded filename=false' \
    '[cpu]' \
    'core=normal' \
    'cputype=486' \
    'cycles=fixed 20000' \
    '[render]' \
    'aspect=false' \
    'scaler=none' \
    '[dos]' \
    'xms=true' \
    'ems=true' \
    'umb=true' \
    'int33=true' \
    'int33 max x=640' \
    'int33 max y=480' \
    '[autoexec]' \
    "mount c $game_dir" \
    'c:' > /tmp/wolong-dosv-capture.conf
printf '%s\n' "${boot_commands[@]}" >> /tmp/wolong-dosv-capture.conf

Xvfb :99 -screen 0 1280x1024x24 >/tmp/wolong-dosv-capture-xvfb.log 2>&1 &
XVFB=$!
for _ in $(seq 1 100); do
    xdpyinfo -display :99 >/dev/null 2>&1 && break
    sleep 0.1
done
if ! xdpyinfo -display :99 >/dev/null 2>&1; then
    echo 'Xvfb 未能啟動。' >&2
    exit 1
fi

DISPLAY=:99 SDL_AUDIODRIVER=dummy dosbox-x -conf /tmp/wolong-dosv-capture.conf -nomenu >"$log_file" 2>&1 &
DB=$!
trace "started dosbox-x pid=$DB"
WID=""
for _ in $(seq 1 100); do
    WID=$(DISPLAY=:99 xdotool search --onlyvisible --name 'DOSBox-X' 2>/dev/null | head -1 || true)
    [ -n "$WID" ] && break
    sleep 0.2
done
if [ -z "$WID" ]; then
    tail -80 "$log_file" >&2 || true
    echo 'DOSBox-X 視窗未出現。' >&2
    exit 1
fi
DISPLAY=:99 xdotool windowactivate --sync "$WID" 2>/dev/null || true
trace "found window wid=$WID"

window_geometry() {
    # DOSBox-X 的 INT 33 integration 吃 root 座標，不吃 --window 相對座標。
    eval "$(DISPLAY=:99 xdotool getwindowgeometry --shell "$WID")"
}

shot() {
    local name=$1
    local root_png
    # SDL2 surface 顯示在 Xvfb 的 root compositing 結果上；直接抓 child window
    # 會得到全黑 PNG。先抓 root，再依目前 client geometry 裁成遊戲實際畫面。
    if ! kill -0 "$DB" 2>/dev/null; then
        echo 'DOSBox-X 在擷取前已結束。' >&2
        tail -80 "$log_file" >&2 || true
        exit 1
    fi
    window_geometry
    : "${WIDTH:?未能取得 DOSBox-X 視窗寬度}"
    : "${HEIGHT:?未能取得 DOSBox-X 視窗高度}"
    : "${X:?未能取得 DOSBox-X 視窗 X 座標}"
    : "${Y:?未能取得 DOSBox-X 視窗 Y 座標}"
    root_png=$(mktemp /tmp/wolong-dosv-root.XXXXXX.png)
    DISPLAY=:99 import -window root "$root_png"
    convert "$root_png" -crop "${WIDTH}x${HEIGHT}+${X}+${Y}" +repage "$out_dir/$name.png"
    rm -f "$root_png"
    printf '%s\n' "$name.png" >> "$out_dir/manifest.txt"
    trace "shot name=$name geometry=${WIDTH}x${HEIGHT}+${X}+${Y}"
}

debug_snapshot() {
    [ "$debug_capture" = '1' ] || return 0
    DISPLAY=:99 xwininfo -root -tree > "$out_dir/x11-tree.txt" 2>&1 || true
    DISPLAY=:99 xdotool getwindowgeometry --shell "$WID" > "$out_dir/window-geometry.txt" 2>&1 || true
    DISPLAY=:99 import -window root "$out_dir/root-screen.png" || true
    cp "$log_file" "$out_dir/dosbox-x.log" 2>/dev/null || true
}

clickat() {
    local px=$1
    local py=$2
    window_geometry
    DISPLAY=:99 xdotool mousemove "$((X + px))" "$((Y + py))"
    sleep 0.15
    DISPLAY=:99 xdotool mousedown 1
    sleep 0.12
    DISPLAY=:99 xdotool mouseup 1
    sleep 0.45
}

rclickat() {
    local px=$1
    local py=$2
    window_geometry
    DISPLAY=:99 xdotool mousemove "$((X + px))" "$((Y + py))"
    sleep 0.15
    DISPLAY=:99 xdotool mousedown 3
    sleep 0.12
    DISPLAY=:99 xdotool mouseup 3
    sleep 0.45
}

press_here() {
    # 原版清單的第一下會藏起 guest 游標；第二下不應重新以游標辨識定位。
    DISPLAY=:99 xdotool mousedown 1
    sleep 0.12
    DISPLAY=:99 xdotool mouseup 1
    sleep 0.45
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

wait_for_password() {
    local deadline root_png orange
    sleep "$password_wait"
    trace "password wait started"
    deadline=$((SECONDS + password_timeout))
    while [ "$SECONDS" -lt "$deadline" ]; do
        if ! kill -0 "$DB" 2>/dev/null; then
            echo 'DOSBox-X 在等待密碼頁時已結束。' >&2
            tail -80 "$log_file" >&2 || true
            exit 1
        fi
        window_geometry
        root_png=$(mktemp /tmp/wolong-dosv-password.XXXXXX.png)
        DISPLAY=:99 import -window root "$root_png"
        # 密碼頁三枚赭色按鈕位於 client (340..540, 215..255)。以 #C38220
        # 的實像素數判定，避免開場插圖短暫靜止時就誤點。
        orange=$(convert "$root_png" -crop "200x40+$((X + 340))+$((Y + 215))" \
            -format '%c' histogram:info:- | awk '/#C38220/ { count += $1 } END { print count + 0 }')
        rm -f "$root_png"
        if [ "$orange" -ge 500 ]; then
            trace "password screen detected orange=$orange"
            return 0
        fi
        sleep 1
    done
    echo "在 ${password_timeout} 秒內未辨識到 DOS/V 密碼頁按鈕。" >&2
    tail -80 "$log_file" >&2 || true
    exit 1
}

{
    printf 'source=松崗 DOS/V 唯讀實機畫面\n'
    printf 'machine=vgaonly\ncore=normal\ncputype=486\ncycles=fixed 20000\n'
    printf 'mouse_emulation=integration\n'
    printf 'timeline=%s\n' "$timeline"
} > "$out_dir/capture-metadata.txt"
: > "$out_dir/manifest.txt"

if [ "$bypass_password" = '1' ]; then
    # docs/playtest/18 已證實：空白後按原始「確定」可進入開場。
    wait_for_password
    clickat 497 240
    trace "clicked original password confirm"
else
    sleep "$password_wait"
    trace "password bypass disabled"
fi
debug_snapshot

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
            DISPLAY=:99 xdotool key --clearmodifiers "$arg"
            sleep 0.35
            trace "step end=$step"
            ;;
        type)
            trace "step begin=$step"
            DISPLAY=:99 xdotool type --delay 60 "$arg"
            sleep 0.35
            trace "step end=$step"
            ;;
        click)
            trace "step begin=$step"
            clickat "${arg%%,*}" "${arg##*,}"
            trace "step end=$step"
            ;;
        rclick)
            # 原版的取消是右鍵（docs/re/53 §3「右鍵取消走 sub_18F7C」）。
            trace "step begin=$step"
            rclickat "${arg%%,*}" "${arg##*,}"
            trace "step end=$step"
            ;;
        press)
            trace "step begin=$step"
            press_here
            trace "step end=$step"
            ;;
	        record)
	            trace "step begin=$step"
	            IFS=',' read -r seconds fps prefix <<< "$arg"
	            record "$seconds" "$fps" "$prefix"
	            trace "step end=$step"
	            ;;
	        audio-start|audio-stop)
	            trace "step begin=$step"
	            DISPLAY=:99 xdotool key --clearmodifiers ctrl+F6
	            sleep 0.5
	            trace "step end=$step"
	            ;;
        savefile)
            # 把 guest 寫的存檔複製出去。**這是原版資產**，只會落在
            # workplace/（gitignore），不進版控（CLAUDE.md §9）。
            # 用途：讓 remake 從同一份存檔開局，做中局的同狀態對拍
            # （docs/spec/90 §2）。
            trace "step begin=$step"
            if [ -f "$game_dir/SAVE.DAT" ]; then
                cp "$game_dir/SAVE.DAT" "$out_dir/$arg"
                trace "step end=$step bytes=$(stat -c%s "$out_dir/$arg")"
            else
                echo "找不到 $game_dir/SAVE.DAT，沒有東西可複製" >&2
                exit 1
            fi
            ;;
        *)
            echo "不支援的擷取步驟：$step" >&2
            exit 1
            ;;
    esac
done

printf 'DOSBox-X 視窗=%s\n' "$WID"
tail -20 "$log_file" || true
