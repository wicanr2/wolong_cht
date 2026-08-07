#!/usr/bin/env bash
# 在 headless DOSBox 裡跑原版，送鍵並截圖。用來建立**視覺 oracle**。
#
#   tools/dosbox.sh dosv "wait:5;shot:01-boot;key:Return;wait:2;shot:02"
#
# timeline 步驟：wait:秒 / key:KEYSYM / type:字串 / shot:名稱
#
# 為什麼一定要固定 cycles：這款是即時制。`cycles=auto` 之下同一串按鍵
# 每次跑到的遊戲內時間點都不同，截圖對不起來、bug 重現不了。
# 設定出自 msdostest（見 dosbox/README.md）。
#
# ⚠ 只跑 dosv。PC-98 版要 DOSBox-X 的 machine=pc98，本 image 的
# DOSBox 0.74 不支援。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_DOSBOX_IMAGE:-dosbox-run:latest}"

VER="${1:?用法: tools/dosbox.sh <dosv> [timeline]}"
TIMELINE="${2:-wait:8;shot:default}"

GAME="$REPO_ROOT/workplace/dosbox/$VER"
SHOTS="$REPO_ROOT/workplace/dosbox/shots"
mkdir -p "$SHOTS"

# workplace/orig 是唯讀的，複製一份可寫副本給 DOSBox 用。
# 已存在就不覆蓋 —— 存檔 diff 實驗做到一半不該被洗掉。
if [ ! -d "$GAME" ]; then
    echo "[dosbox] 建立可寫副本 $GAME"
    mkdir -p "$GAME"
    cp "$REPO_ROOT/workplace/orig/$VER/"* "$GAME/"
    chmod -R u+w "$GAME"
fi

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$GAME:/game" -v "$SHOTS:/shots" \
    -u "$(id -u):$(id -g)" -e HOME=/tmp \
    --entrypoint bash \
    "$IMAGE" -c "
set -e
export XDG_RUNTIME_DIR=/tmp
cat > /tmp/dosbox.conf <<'EOF'
[sdl]
autolock=false
[dosbox]
machine=vgaonly
memsize=16
[cpu]
core=normal
cputype=486
cycles=fixed 20000
[render]
aspect=false
scaler=none
[dos]
xms=true
ems=true
umb=true
[autoexec]
mount c /game
c:
EOF
Xvfb :99 -screen 0 1280x1024x24 >/tmp/xvfb.log 2>&1 &
XVFB=\$!
for i in \$(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
DISPLAY=:99 dosbox -conf /tmp/dosbox.conf -noconsole >/tmp/dosbox.log 2>&1 &
DB=\$!
sleep 4
WID=\$(timeout 20 env DISPLAY=:99 xdotool search --sync --onlyvisible --class -- DOSBox 2>/dev/null | head -1) || true
if [ -z \"\$WID\" ]; then
    echo '!!! DOSBox 視窗沒出現'; cat /tmp/dosbox.log || true
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
        shot) DISPLAY=:99 import -window \$WID \"/shots/\$arg.png\" && echo \"  截圖 \$arg.png\" ;;
        *) echo \"未知步驟 \$step\" >&2 ;;
    esac
done
echo '--- dosbox log ---'; tail -20 /tmp/dosbox.log || true
kill -9 \$DB \$XVFB 2>/dev/null || true
"
ls -la "$SHOTS"
