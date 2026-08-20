# Android 建置環境（gomobile / ebitenmobile + Gradle）
#
# [HARD] 本專案一律用 docker 編譯。基底沿用既有的 Android 工具鏈映像，
# 只補兩樣它缺的東西：
#
#   1. **Go 1.25**。基底裝的是 1.24.13，而 `go.mod` 寫 `go 1.25.0`，
#      直接編會被 toolchain 檢查擋下。
#   2. **ebitenmobile**。Android 的 AAR 不是 `gomobile bind` 直接產的——
#      `EbitenView` 來自 Ebiten 自己的包裝器，它內部才去呼叫 gomobile。
#
# Gradle 的相依不放進映像：那是**專案的**相依，版本會跟著 `android/build.gradle`
# 走，烤進映像會讓兩者悄悄脫鉤。改用具名 volume 當 `GRADLE_USER_HOME`
#（見 tools/android_build.sh）。
FROM rich2-go-android:20260809

USER root

ARG GO_VERSION=1.25.12
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tgz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf /tmp/go.tgz \
    && rm /tmp/go.tgz \
    && /usr/local/go/bin/go version

ENV PATH=/usr/local/go/bin:/go/bin:$PATH
ENV GOPATH=/go
ENV GOFLAGS=-buildvcs=false
ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64

# gomobile 的 init 會在 NDK 底下找 toolchain；SDK/NDK 由基底提供。
ENV ANDROID_HOME=/opt/android-sdk
ENV ANDROID_NDK_HOME=/opt/android-sdk/ndk/27.2.12479018

# 讓非 root 使用者（-u $(id -u)）寫得進 GOPATH 與 gradle 的工作區。
RUN mkdir -p /go /gomod /gocache /gradle \
    && chmod -R 0777 /go /gomod /gocache /gradle

WORKDIR /src
