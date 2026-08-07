# DOSBox-X 參考環境（原版 oracle）
#
# 為什麼不是普通 DOSBox：**`machine=pc98` 只有 DOSBox-X 有**。
# 本專案要跑 PC-98 日文原版（那一版沒有防拷，見 docs/playtest/01），
# 所以非 DOSBox-X 不可。Debian bookworm 沒有這個套件，trixie 才有。
#
# 全程 headless：Xvfb 起虛擬 framebuffer，xdotool 送鍵，import 截圖。
FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        dosbox-x \
        xvfb \
        x11-apps \
        x11-utils \
        xdotool \
        imagemagick \
        procps \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /work
