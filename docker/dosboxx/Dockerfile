# DOSBox-X 參考環境（原版 oracle）
#
# 為什麼不是普通 DOSBox：**`machine=pc98` 只有 DOSBox-X 有**。
# 本專案要跑 PC-98 日文原版（那一版沒有防拷，見 docs/playtest/01），
# 所以非 DOSBox-X 不可。Debian bookworm 沒有這個套件，trixie 才有。
#
# 全程 headless：Xvfb 起虛擬 framebuffer，xdotool 送鍵，import 截圖。
#
# ffmpeg 是給推廣片用的：逐幀 `import` 一張要 280–330 ms，錄不出即時制的
# 真實速度（實測標稱 10 fps 只跑到 3.08 fps）。`x11grab` 直接錄同一個
# Xvfb，是**真實時間、固定幀率**。
FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        dosbox-x \
        xvfb \
        x11-apps \
        x11-utils \
        xdotool \
        imagemagick \
        ffmpeg \
        procps \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /work
