# Android 模擬器（headless smoke 用）。
#
# 以本專案的 Android 建置環境為基底，補三樣東西：
#
#   1. emulator 在 headless Linux 需要的 X11／PulseAudio 動態庫。
#   2. `emulator` 本體與一份 system image。
#   3. **一顆建好的 AVD**——留在映像裡，smoke 才不必每次重建，
#      也才不會因為建立參數漂移而每次跑在不同的機器上。
#
# ⚠ system image 用 `android-34;google_apis;x86_64`：這台只有 x86_64，
# arm64 的 image 要靠模擬指令集，慢到不適合當驗收迴圈。
# **API 34 不等於 targetSdk 35 的實機行為**，只是這一版 SDK 有現成的 x86_64
# image；行為差異要靠里程碑 H 的實機驗收擋，模擬器擋不到。
FROM wolong-go-android:20260820

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

ARG SYSIMG="system-images;android-34;google_apis;x86_64"
RUN yes | "$ANDROID_HOME/cmdline-tools/latest/bin/sdkmanager" --licenses >/dev/null \
    && "$ANDROID_HOME/cmdline-tools/latest/bin/sdkmanager" "emulator" "$SYSIMG"

# AVD 建在共用位置，並開放寫入——smoke 以 `-u $(id -u)` 跑，
# 而模擬器開機時會寫回 AVD 目錄。
ENV ANDROID_AVD_HOME=/opt/avd
RUN mkdir -p /opt/avd \
    && echo no | "$ANDROID_HOME/cmdline-tools/latest/bin/avdmanager" create avd \
        --name wolong --package "$SYSIMG" --device pixel_5 --path /opt/avd/wolong.avd \
    && printf 'hw.lcd.density=440\nhw.keyboard=yes\ndisk.dataPartition.size=2G\n' \
        >> /opt/avd/wolong.avd/config.ini \
    && chmod -R 0777 /opt/avd

ENV PATH=$ANDROID_HOME/emulator:$ANDROID_HOME/platform-tools:$PATH
WORKDIR /src
