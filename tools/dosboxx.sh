#!/usr/bin/env bash
# 在 headless DOSBox-X 裡跑 PC-98 日文原版，送鍵並截圖。
#
#   tools/dosboxx.sh "wait:10;shot:01-boot;key:Return;wait:3;shot:02"
#
# timeline 步驟：wait:秒 / key:KEYSYM / type:字串 / shot:名稱
#
# 為什麼跑 PC-98 而不是松崗版：**PC-98 原版沒有防拷**（docs/playtest/01）。
# 兩版是同一份原始碼的兩次編譯，所以玩法規則、AI、資料結構、
# 地圖與戰鬥畫面全部可以在這一版上驗，不必去過密碼那一關。
#
# 設定出自 msdostest（見 dosbox/README.md），差別是原設定掛 .hdi，
# 我們手上是五片 .fdi。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_DOSBOXX_IMAGE:-wolong-dosboxx}"

# 模式：
#   mount（預設）掛 workplace/orig/pc98 當 C:，用 DOSBox-X 內建 DOS 跑 YNSHELL。
#              原版是裝到硬碟的（E 片的 AUTOEXEC.BAT 寫 `CD GARYO`），
#              掛目錄等價於裝好的狀態，而且跳過 FGDOS 與換片。
#   floppy     照原樣掛五片磁片開機。**遊戲會報「ファイルが異常です」**
#              ——它預期檔案都在同一顆硬碟上，跨片找不到。
MODE="${WOLONG_PC98_MODE:-mount}"
TIMELINE="${1:-wait:15;shot:pc98-default}"
FDI="$REPO_ROOT/workplace/dosbox/pc98_fdi"
GAME="$REPO_ROOT/workplace/dosbox/pc98"
SHOTS="$REPO_ROOT/workplace/dosbox/shots"
mkdir -p "$SHOTS"

# workplace/orig 是唯讀的，磁片映像要可寫（遊戲會寫存檔）。
if [ ! -d "$FDI" ]; then
    mkdir -p "$FDI"
    cp "$REPO_ROOT/workplace/orig/pc98_fdi/"*.fdi "$FDI/"
    chmod -R u+w "$FDI"
fi
if [ ! -d "$GAME" ]; then
    echo "[dosbox-x] 建立可寫遊戲副本 $GAME"
    mkdir -p "$GAME"
    cp "$REPO_ROOT/workplace/orig/pc98/"* "$GAME/"
    chmod -R u+w "$GAME"
fi

if [ "$MODE" = "mount" ]; then
    AUTOEXEC="mount c /game
c:
YNSHELL.COM"
else
    AUTOEXEC="imgmount 0 /fdi/Garyou_A.fdi /fdi/Garyou_B.fdi /fdi/Garyou_C.fdi /fdi/Garyou_D.fdi /fdi/Garyou_E.fdi -t floppy
boot -l a"
fi

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$FDI:/fdi" -v "$GAME:/game" -v "$SHOTS:/shots" \
    -u "$(id -u):$(id -g)" -e HOME=/tmp \
    "$IMAGE" bash -c "
set -e
cat > /tmp/dosbox-x.conf <<'EOF'
[sdl]
# autolock 一定要開：關掉的話 DOSBox 不會把滑鼠事件轉給 guest，
# xdotool 的合成點擊會被吃掉，看起來像「遊戲不理滑鼠」。
autolock=true
[dosbox]
memsize=8
machine=pc98
pc-98 sound bios=true
cascade interrupt ignore in service=true
[cpu]
core=normal
cputype=486
cycles=20000
[render]
aspect=false
scaler=none
[autoexec]
$AUTOEXEC
EOF
Xvfb :99 -screen 0 1280x1024x24 >/tmp/xvfb.log 2>&1 &
XVFB=\$!
for i in \$(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
DISPLAY=:99 dosbox-x -conf /tmp/dosbox-x.conf -nomenu >/tmp/dbx.log 2>&1 &
DB=\$!
sleep 6
WID=\$(timeout 25 env DISPLAY=:99 xdotool search --sync --onlyvisible --name 'DOSBox-X' 2>/dev/null | head -1) || true
if [ -z \"\$WID\" ]; then
    WID=\$(DISPLAY=:99 xdotool search --onlyvisible '' 2>/dev/null | tail -1) || true
fi
if [ -z \"\$WID\" ]; then
    echo '!!! DOSBox-X 視窗沒出現'; tail -40 /tmp/dbx.log || true
    kill -9 \$DB \$XVFB 2>/dev/null || true; exit 1
fi
DISPLAY=:99 xdotool windowactivate --sync \$WID 2>/dev/null || true
IFS=';' read -ra STEPS <<< '$TIMELINE'
for step in \"\${STEPS[@]}\"; do
    kind=\${step%%:*}; arg=\${step#*:}
    case \"\$kind\" in
        wait) sleep \"\$arg\" ;;
        key)  DISPLAY=:99 xdotool key --window \$WID --clearmodifiers \"\$arg\"; sleep 0.4 ;;
        type) DISPLAY=:99 xdotool type --window \$WID --delay 60 \"\$arg\"; sleep 0.4 ;;
        click) # 遊戲要滑鼠（KI.EXE 有 \'Mouse driver not install ?\'）。
               # arg 格式 x,y；不給就點畫面中央。
               cx=\${arg%%,*}; cy=\${arg##*,}
               DISPLAY=:99 xdotool mousemove --window \$WID \"\$cx\" \"\$cy\" click 1
               sleep 0.5 ;;
        shot) DISPLAY=:99 import -window \$WID \"/shots/\$arg.png\" && echo \"  截圖 \$arg.png\" ;;
    esac
done
echo '--- dosbox-x log 尾段 ---'; tail -15 /tmp/dbx.log || true
kill -9 \$DB \$XVFB 2>/dev/null || true
"
