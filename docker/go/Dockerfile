# Go + Ebiten 建置／測試環境
#
# [HARD] 本專案一律用 docker 編譯，不污染主機 Go 環境。
# Ebiten 即使跑 headless 測試也要 X11 / OpenGL 的標頭檔才編得過。
FROM golang:1.25-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
        libgl1-mesa-dev \
        libxcursor-dev \
        libxi-dev \
        libxinerama-dev \
        libxrandr-dev \
        libxxf86vm-dev \
        libasound2-dev \
        pkg-config \
        # 視覺產物一律 dump PNG 肉眼比對，不接受「編譯過了」當驗收
        xvfb \
        x11-utils \
        xdotool \
    && rm -rf /var/lib/apt/lists/*

ENV GOPATH=/go
ENV GOFLAGS=-buildvcs=false
ENV CGO_ENABLED=1
WORKDIR /src
