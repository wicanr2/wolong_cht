#!/usr/bin/env bash
# 帶 Native AI Bridge 的 headless DOSBox-X（動態 oracle）。
#
#   tools/dosboxx_bridge.sh start      # 背景起 DOSBox-X ＋ PC-98 原版，bridge 在 9876
#   tools/dosboxx_bridge.sh health     # 對 bridge 送一次請求，確認它活著
#   tools/dosboxx_bridge.sh shot 名稱  # 截圖
#   tools/dosboxx_bridge.sh stop
#
# 為什麼是 PC-98 而不是松崗版：**PC-98 原版沒有防拷**（docs/playtest/01）。
# 兩版是同一份原始碼的兩次編譯，規則層可以在這一版上驗。
#
# ⚠ **bridge 只在 `--enable-debug=heavy` 的建置裡存在**（`debug_ai.cpp` 整支
# 包在 `#if C_DEBUG`）。Debian 打包的 dosbox-x 沒有，所以這一支用的是
# `docker/dosboxx-bridge/Dockerfile` 自己編的映像，不是 `docker/dosboxx`。
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_BRIDGE_IMAGE:-wolong-dosboxx-bridge}"
NAME="${WOLONG_BRIDGE_NAME:-wolong-bridge}"
CONF="$REPO_ROOT/dosbox/pc98.conf"
GAME="$REPO_ROOT/workplace/dosbox/pc98"
SHOTS="$REPO_ROOT/workplace/dosbox/shots"

case "${1:-health}" in
start)
    mkdir -p "$SHOTS" "$GAME"
    [ -n "$(ls -A "$GAME" 2>/dev/null)" ] || {
        cp "$REPO_ROOT/workplace/orig/pc98/"* "$GAME/"; chmod -R u+w "$GAME"; }
    docker rm -f "$NAME" >/dev/null 2>&1 || true
    # 直接跑 KI.EXE，跳過 ROGO.EXE 的兩分鐘開場（同 dosboxx.sh 的 noopen 模式）。
    # ⚠ YNMOUSE.COM 一定要載，否則 KI.EXE 印「Mouse driver not install ?」。
    # 倉庫裡的 dosbox/pc98.conf 掛的是五片 .fdi，那條路遊戲會報
    # 「ファイルが異常です」（跨片找不到檔），所以這裡自己組 conf。
    #
    # --network host：bridge 綁 127.0.0.1，MCP server 那一顆也用 host 網路才連得到。
    # ⚠ **-t 是必要的，而且 dosbox-x 的輸出不能重導向。**
    # heavy debugger 的 UI 是 ncurses，開在**控制終端**上，不是 X 視窗；
    # DOSBox-X 判斷有沒有終端是看自己的 stdout。接了 `| tee` 或
    # `> log` 之後 stdout 不再是 TTY，log 裡會出現
    # 「Debugger is not available unless you start DOSBox-X from a terminal」，
    # 而 bridge 對每個 `memory.read` 回 DEBUGGER_NOT_STOPPED——
    # **看起來像 bridge 壞了，其實只是沒有地方畫那個介面**。
    # 要看 log 用 `docker logs $NAME`，不要在容器裡重導向。
    docker run -d -t --rm --name "$NAME" \
        --log-opt max-size=10m --log-opt max-file=3 \
        --network host --memory 2g --cpus 2 --pids-limit 256 \
        -v "$GAME:/game" -v "$SHOTS:/shots" -e BREAK_START="${WOLONG_BREAK_START:-}" \
        "$IMAGE" bash -c '
            cat > /tmp/dosbox-x.conf <<EOF
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
mount c /game
c:
YNFONT.COM
YNSOUND.COM
YNMOUSE.COM
KI.EXE
EOF
            Xvfb :99 -screen 0 1024x768x24 >/dev/null 2>&1 &
            for i in $(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.2; done
            DISPLAY=:99 dosbox-x -conf /tmp/dosbox-x.conf -nomenu ${BREAK_START:-}
        ' >/dev/null
    echo "[bridge] 已啟動 $NAME；等 bridge 開 port…"
    for i in $(seq 1 30); do
        if docker exec "$NAME" bash -c 'exec 3<>/dev/tcp/127.0.0.1/9876' 2>/dev/null; then
            echo "[bridge] 9876 已就緒"; exit 0
        fi
        sleep 2
    done
    echo "[bridge] 30 次都連不上 9876。先看 workplace/dosbox/shots/dosbox.log" >&2
    exit 1
    ;;
health)
    # 正對照：沒有這一步就分不出「bridge 沒開」與「請求寫錯了」。
    # 協定是 newline-delimited JSON，方法名是 `debug.status` 不是工具名
    # （工具名 get_debug_status 是 MCP 那一層的，bridge 不認得）。
    #
    # ⚠ 判成功要看**回應內容**，不要看 nc 的 exit code。bridge 不會關連線，
    # 所以 nc 一定是被 timeout 砍掉、一定回非 0——照 exit code 判會
    # 在明明拿到正確回應時報「連不上」。
    reply="$(printf '{"id":1,"method":"debug.status"}\n' \
        | timeout 3 nc 127.0.0.1 9876 2>/dev/null | head -1 || true)"
    echo "${reply:-（沒有回應）}"
    case "$reply" in
        *'"ok":true'*) exit 0 ;;
        *) echo "[bridge] 沒有拿到正常回應——先跑 tools/dosboxx_bridge.sh start" >&2
           exit 1 ;;
    esac
    ;;
shot)
    docker exec "$NAME" bash -c "DISPLAY=:99 import -window root /shots/${2:-shot}.png"
    echo "$SHOTS/${2:-shot}.png"
    ;;
stop)
    docker rm -f "$NAME" >/dev/null 2>&1 && echo "[bridge] 已停止" ;;
*)  echo "用法: dosboxx_bridge.sh {start|health|shot 名稱|stop}" >&2; exit 2 ;;
esac
