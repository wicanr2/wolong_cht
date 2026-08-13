# Android 模擬器只用於隔離 smoke；以既有 Android 工具鏈為基底，補齊
# emulator 在 headless Linux 需要的 X11／PulseAudio 動態庫。
FROM rich2-go-android:20260809

USER root
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        libasound2 \
        libdbus-1-3 \
        libdrm2 \
        libfontconfig1 \
        libgbm1 \
        libgl1 \
        libnss3 \
        libpulse0 \
        libx11-6 \
        libxcb1 \
        libxcomposite1 \
        libxext6 \
        libxfixes3 \
        libxkbfile1 \
        libxkbcommon0 \
        libxi6 \
        libxrandr2 \
        libxrender1 \
        libxtst6 \
    && rm -rf /var/lib/apt/lists/*

USER 1000:1000
