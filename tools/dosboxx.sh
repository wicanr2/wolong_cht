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
#   noopen     同 mount，但跳過 ROGO.EXE 的開場（省約兩分鐘）。
#   floppy     照原樣掛五片磁片開機。**遊戲會報「ファイルが異常です」**
#              ——它預期檔案都在同一顆硬碟上，跨片找不到。
MODE="${WOLONG_PC98_MODE:-mount}"
AUTOLOCK="${WOLONG_AUTOLOCK:-true}"
# 滑鼠相關的三個旋鈕（docs/playtest/04、05）。預設值是目前已知「不會動」的組合，
# 留著是為了讓舊 timeline 重跑得到一樣的結果；要試別的組合就覆蓋這三個。
SDL_AUTOLOCK="${WOLONG_SDL_AUTOLOCK:-false}"
MOUSE_EMU="${WOLONG_MOUSE_EMU:-never}"
# LOG_MOUSE=true 會打開 DOSBox-X 自己的 mouse log。**這是分辨
# 「事件沒到 DOSBox」與「事件到了但 guest 不理」的唯一直接證據。**
LOG_MOUSE="${WOLONG_LOG_MOUSE:-false}"
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
elif [ "$MODE" = "noopen" ]; then
    # 跳過開場。YNSHELL.COM 只是個啟動器——它的字串表就四個檔名：
    # YNFONT.COM／YNSOUND.COM／YNMOUSE.COM／ROGO.EXE。
    # 自己載那三個 TSR 再直接跑 KI.EXE，就少掉 ROGO.EXE 的兩分鐘開場。
    # ⚠ YNMOUSE.COM 一定要載，不然遊戲收不到滑鼠（KI.EXE 會印
    #    「ERROR: Mouse driver not install ?」）。
    AUTOEXEC="mount c /game
c:
YNFONT.COM
YNSOUND.COM
YNMOUSE.COM
KI.EXE"
else
    AUTOEXEC="imgmount 0 /fdi/Garyou_A.fdi /fdi/Garyou_B.fdi /fdi/Garyou_C.fdi /fdi/Garyou_D.fdi /fdi/Garyou_E.fdi -t floppy
boot -l a"
fi

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$FDI:/fdi" -v "$GAME:/game" -v "$SHOTS:/shots" \
    -u "$(id -u):$(id -g)" -e HOME=/tmp -e AUTOLOCK="$AUTOLOCK" \
    "$IMAGE" bash -c "
set -e
cat > /tmp/dosbox-x.conf <<'EOF'
[sdl]
# 這兩個值由 WOLONG_SDL_AUTOLOCK／WOLONG_MOUSE_EMU 決定，見 docs/playtest/04、05。
autolock=$SDL_AUTOLOCK
mouse_emulation=$MOUSE_EMU
usesystemcursor=false
[log]
mouse=$LOG_MOUSE
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
# 視窗在 root 上的絕對位置。DOSBox 讀的是**絕對指標位置**，
# 用 xdotool mousemove --window 送相對座標它收不到 ——
# 症狀是點擊完全沒反應，看起來像「遊戲不理滑鼠」。
eval \$(DISPLAY=:99 xdotool getwindowgeometry --shell \$WID)
echo \"視窗位置 (\$X,\$Y) 大小 \${WIDTH}x\${HEIGHT}\"
echo \"鍵盤焦點視窗 = \$(DISPLAY=:99 xdotool getwindowfocus 2>&1 || true)（DOSBox-X 視窗 = \$WID）\"
# probe 只印狀態、不動任何東西：指標在哪、指標底下是哪個視窗。
# 用來分辨「host 的指標根本沒進到視窗矩形」與「進去了但 DOSBox 不理」。
# ⚠ 一定要用 ( ) 當函式體（子 shell）—— getmouselocation --shell 也吐 X=／Y=，
# 用 { } 會把上面那組視窗座標蓋掉，後面的點擊就全部歪掉。
probe() (
    eval \$(DISPLAY=:99 xdotool getmouselocation --shell 2>/dev/null)
    echo \"    指標 root(\$X,\$Y)　底下的視窗=\$WINDOW\"
)
# 所有 xdotool 都走這個包裝，失敗時把**完整命令列**印出來。
# 上一輪有一行 mousemove: option --window requires an argument
# 無從歸屬到任何一步——裸呼叫的 stderr 認不出是誰發的。
xd() {
    if ! DISPLAY=:99 xdotool \"\$@\" 2>/tmp/xd.err; then
        echo \"    ⚠ 這一步失敗了：xdotool \$*\"; sed 's/^/      /' /tmp/xd.err
    elif [ -s /tmp/xd.err ]; then
        echo \"    ⚠ xdotool \$* 有 stderr：\"; sed 's/^/      /' /tmp/xd.err
    fi
}
# ⭐ 游標的位置感測器。純紅 #FF0000 在整張畫面**只有游標在用**
# （實測 640×400 裡剛好 56 個像素，全部落在游標框內），所以認色就夠，
# 不必做樣板比對。回傳游標紅色區塊的左上角，也就是箭頭尖端。
curpos() (
    DISPLAY=:99 import -window \$WID /tmp/cur.png 2>/dev/null || true
    # +opaque 把「不是純紅」的像素全部塗黑，剩下的只有游標；
    # %@ 直接給那塊非背景區域的邊界框（如 12x12+349+215）。
    # 比掃 txt: 快 20 倍（38 ms vs 768 ms），16 步的閉迴路差半分鐘。
    bb=\$(convert /tmp/cur.png -fill black +opaque '#FF0000' -format '%@' info:- 2>/dev/null)
    case \"\$bb\" in
        [3-9]x*|[1-9][0-9]x*) echo \"\${bb#*+}\" | tr '+' ' ' ;;
        *) echo '-1 -1' ;;   # 找不到游標（畫面全黑 → trim 退化成 1x1）
    esac
)

# gotoxy x y —— **閉迴路**把遊戲內游標移到指定位置。
# 為什麼不能開迴路送一個大位移：PC-98 是匯流排滑鼠，每次回報的位移是
# 帶號 8 bit，大跳躍會被截掉（實測 clip(-2000)+clip(600) = -128+127 = -1，
# 送 600 px 只走了 1 px）。所以一次最多送 60 px，看著截圖逼近。
gotoxy() {
    local tx=\$1 ty=\$2 cx cy dx dy ax ay i ok=0
    for i in \$(seq 1 16); do
        set -- \$(curpos); cx=\$1; cy=\$2
        [ \"\$cx\" = -1 ] && { echo \"    ⚠ 畫面上找不到游標\"; return 1; }
        dx=\$((tx-cx)); dy=\$((ty-cy)); ax=\${dx#-}; ay=\${dy#-}
        if [ \$ax -le 2 ] && [ \$ay -le 2 ]; then
            echo \"    游標到 (\$cx,\$cy)，走了 \$i 步\"; ok=1; break
        fi
        [ \$dx -gt 60 ] && dx=60; [ \$dx -lt -60 ] && dx=-60
        [ \$dy -gt 60 ] && dy=60; [ \$dy -lt -60 ] && dy=-60
        xd mousemove_relative -- \$dx \$dy
        sleep 0.15
    done
    [ \$ok = 0 ] && echo \"    ⚠ 16 步內沒逼近到 (\$tx,\$ty)，停在 (\$cx,\$cy)\"
    true
}
IFS=';' read -ra STEPS <<< '$TIMELINE'
for step in \"\${STEPS[@]}\"; do
    kind=\${step%%:*}; arg=\${step#*:}
    case \"\$kind\" in
        wait) sleep \"\$arg\" ;;
        key)  xd key --window \$WID --clearmodifiers \"\$arg\"; sleep 0.4 ;;
        type) xd type --window \$WID --delay 60 \"\$arg\"; sleep 0.4 ;;
        # xkey／xtype 走 XTEST（不帶 --window），送的是**真的**按鍵事件而不是
        # XSendEvent 合成事件。前提是視窗握有鍵盤焦點——實測是（docs/playtest/05）。
        # DOSBox-X 自己的快捷鍵（如 ctrl+F10 擷取滑鼠）走 mapper，
        # 合成事件不一定吃得到，這條路比較穩。
        xkey)  xd key --clearmodifiers \"\$arg\"; sleep 0.4 ;;
        xtype) xd type --delay 60 \"\$arg\"; sleep 0.4 ;;
        click|rclick)
               # ⚠ autolock 擷取滑鼠之後，DOSBox 收的是**相對位移**，
               # 絕對定位（xdotool mousemove 到某個 root 座標）會漂 ——
               # 症狀是游標位置與送的座標對不上，點到隔壁的選項。
               # 標準做法：先往左上角灌一個大位移把游標「歸零」（會被夾住），
               # 再走精確的相對位移過去。
               cx=\${arg%%,*}; cy=\${arg##*,}
               btn=1; [ \"\$kind\" = rclick ] && btn=3
               if [ \"\$AUTOLOCK\" = false ]; then
                   xd mousemove \$((X+cx)) \$((Y+cy))
               else
                   xd mousemove_relative -- -2000 -2000
                   sleep 0.1
                   xd mousemove_relative -- \"\$cx\" \"\$cy\"
               fi
               sleep 0.2
               xd mousedown \$btn
               sleep 0.15
               xd mouseup \$btn
               sleep 0.5 ;;
        move)  cx=\${arg%%,*}; cy=\${arg##*,}
               # 與 click 同一套規則：AUTOLOCK=false 走絕對，否則走「歸零再相對」。
               # （舊版這裡永遠走相對，跟 click 不一致——docs/playtest/05）
               if [ \"\$AUTOLOCK\" = false ]; then
                   xd mousemove \$((X+cx)) \$((Y+cy))
               else
                   xd mousemove_relative -- -2000 -2000
                   sleep 0.1
                   xd mousemove_relative -- \"\$cx\" \"\$cy\"
               fi
               sleep 0.3; probe ;;
        # 找不到游標時 gotoxy 回 1，而整段腳本是 set -e —— 不接 || true 會整個中止。
        goto)    gotoxy \${arg%%,*} \${arg##*,} || true ;;
        clickat) gotoxy \${arg%%,*} \${arg##*,} || true
                 sleep 0.2; xd mousedown 1; sleep 0.15; xd mouseup 1; sleep 0.6 ;;
        # press —— **原地**按一下，不移動也不找游標。
        # 兩段式選取（說明書 3.8）的第二下就是這個：第一下反白之後
        # **遊戲會把游標藏起來**，clickat 會卡在「找不到游標」。
        press)   xd mousedown 1; sleep 0.15; xd mouseup 1; sleep 0.6 ;;
        rpress)  xd mousedown 3; sleep 0.15; xd mouseup 3; sleep 0.6 ;;
        probe) probe ;;
        # settle:起,上限 —— 從「起」秒開始每 2 秒截一張，連續三張一模一樣就往下走。
        # ⭐ 為什麼不用固定 wait：開場長度**每次不一樣**（實測 104 s 與 >114 s），
        # 用 wait 會讓整段測試跑在過場動畫裡，而那時候既沒有游標也沒有選單——
        # 測試不會失敗，它會**什麼都沒測到**。過場一直在動、選單完全靜止，
        # 所以「靜止」本身就是可靠的到站訊號。
        # until:md5,上限 —— 每 2 秒截一張，等到畫面的 md5 等於指定值才往下走。
        # settle（等靜止）不夠：**過場動畫每張插圖會定格好幾秒**，
        # 連續三張相同會在開場中途就誤判到站（實測 96 s 那次）。
        # 認畫面本身沒有這個問題——選單畫面的 md5 跨執行是固定的。
        until)
               want=\${arg%%,*}; lim=\${arg##*,}; t=0
               while [ \$t -lt \$lim ]; do
                   DISPLAY=:99 import -window \$WID /tmp/until.png 2>/dev/null || true
                   cur=\$(md5sum /tmp/until.png 2>/dev/null | cut -d' ' -f1)
                   [ \"\$cur\" = \"\$want\" ] && { echo \"    \${t}s 到站\"; break; }
                   sleep 2; t=\$((t+2))
               done
               if [ \$t -ge \$lim ]; then
                   echo \"    ⚠ 到 \${lim}s 還沒等到指定畫面，最後看到的是 \$cur\"
                   DISPLAY=:99 import -window \$WID /shots/until-timeout.png 2>/dev/null || true
               fi
               true ;;
        settle)
               st=\${arg%%,*}; lim=\${arg##*,}
               sleep \"\$st\"; t=\$st; prev=''; same=0
               while [ \$t -lt \$lim ]; do
                   DISPLAY=:99 import -window \$WID /tmp/settle.png 2>/dev/null || true
                   cur=\$(md5sum /tmp/settle.png 2>/dev/null | cut -d' ' -f1)
                   if [ -n \"\$cur\" ] && [ \"\$cur\" = \"\$prev\" ]; then
                       same=\$((same+1))
                       [ \$same -ge 2 ] && { echo \"    畫面在 \${t}s 靜止（連續三張相同）\"; break; }
                   else
                       same=0
                   fi
                   prev=\$cur; sleep 2; t=\$((t+2))
               done
               [ \$t -ge \$lim ] && echo \"    ⚠ 到 \${lim}s 畫面仍在變，沒等到靜止\"
               true ;;
        shot) DISPLAY=:99 import -window \$WID \"/shots/\$arg.png\" && echo \"  截圖 \$arg.png\" ;;
    esac
done
echo '--- DOSBox-X 看到的滑鼠事件 ---'
grep -i -E 'mouse|MOUSE' /tmp/dbx.log | head -40 || echo '（一行都沒有）'
echo '--- dosbox-x log 尾段 ---'; tail -15 /tmp/dbx.log || true
kill -9 \$DB \$XVFB 2>/dev/null || true
"
