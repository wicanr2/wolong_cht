# DOSBox-X ＋ Native AI Bridge（動態 oracle）
#
# 與 docker/dosboxx 的差別：那一顆用 Debian 打包好的 dosbox-x，**沒有除錯器**，
# 更沒有 AI bridge。這一顆從 pmanyeh/dosbox-x 的 `ai-mcp-bridge` 分支自己編，
# 帶 `--enable-debug=heavy`（`debug_ai.cpp` 整支包在 `#if C_DEBUG` 裡）。
#
# 建置 context 是 DOSBox-X-MCP-Debugger/dosbox-src：
#   docker build -f docker/dosboxx-bridge/Dockerfile \
#       -t wolong-dosboxx-bridge ~/cht/DOSBox-X-MCP-Debugger/dosbox-src
#
# heavy debugger 需要 ncurses——少了它 configure 會安靜地退回 non-heavy，
# 而 bridge 仍然編得出來，只是斷點與單步不可用。**先驗 `--enable-debug=heavy`
# 有沒有真的生效**（build log 裡的 `Debugger. yes (heavy)`）。
FROM debian:trixie-slim AS build

RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential automake autoconf libtool pkg-config git ca-certificates \
        libncurses-dev libglu1-mesa-dev libgl-dev \
        libx11-dev libxrandr-dev libxext-dev libxcursor-dev libxinerama-dev libxi-dev \
        libpng-dev zlib1g-dev libfreetype-dev libslirp-dev \
        libfluidsynth-dev libavcodec-dev libavformat-dev libavutil-dev libswscale-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY . /src
RUN chmod +x autogen.sh build-debug configure 2>/dev/null || true
RUN ./build-debug > /tmp/build.log 2>&1 || { tail -40 /tmp/build.log; exit 1; }
# 驗收與建置分開兩層：改驗收條件不會害整包重編。
# **probe 要看 config.h，不要 grep build log 的字樣**——第一版寫成
# `grep -q "Debugger" build.log`，build 明明成功卻擋下來，
# 因為那個字樣根本不在 log 裡。**編得出來不等於編對了**：
# 少了 C_DEBUG，`debug_ai.cpp` 整支被 `#if` 掉，bridge 會安靜地不存在。
RUN test -x src/dosbox-x \
    && grep -q "define C_DEBUG 1" config.h \
    && grep -q "define C_HEAVY_DEBUG 1" config.h \
    && echo "C_DEBUG ＋ C_HEAVY_DEBUG 都在"

FROM debian:trixie-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
        libncurses6 libx11-6 libxrandr2 libxext6 libxcursor1 libxinerama1 libxi6 \
        libpng16-16t64 zlib1g libfreetype6 libslirp0 \
        libfluidsynth3 libavcodec61 libavformat61 libavutil59 libswscale8 \
        xvfb x11-apps x11-utils xdotool imagemagick procps netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*
COPY --from=build /src/src/dosbox-x /usr/local/bin/dosbox-x
WORKDIR /work
