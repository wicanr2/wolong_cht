#!/usr/bin/env bash
# 建 Android APK：ebitenmobile bind → AAR → gradle assembleDebug。
#
#   tools/android_build.sh              # 兩個 ABI（arm64 + amd64）
#   WOLONG_ABIS=android/arm64 tools/android_build.sh
#   WOLONG_NET=none tools/android_build.sh   # 相依都在快取裡時可以斷網跑
#
# ⚠ **容器要跑在 UTF-8 locale 底下**。gobind 會把 Go 的文件註解原封不動抄進
# 產生的 Java 檔，而 javac 沒有帶 `-encoding`，用的是**平台預設字集**。
# 容器預設是 POSIX／US-ASCII，於是每一個中文字都變成一行
# `unmappable character` 編譯錯誤——註解害整包建置失敗。
#
# ⚠ **第一次一定要有網路**：gradle 要抓 AGP 8.7.3 與它的相依（數百 MB），
# 之後留在具名 volume `wl-gradle` 裡。Go 的模組沿用桌面那條 `wl-gomod`。
#
# 預設只建 arm64 與 amd64：arm64 是實機、amd64 是模擬器。32 位的 arm／386
# 要的話自己加，但那會讓 bind 的時間加倍，而**這個專案沒有 32 位的驗收對象**。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${WOLONG_ANDROID_IMAGE:-wolong-go-android:20260820}"
ABIS="${WOLONG_ABIS:-android/arm64,android/amd64}"
NET="${WOLONG_NET:-host}"

if ! docker image inspect "$IMAGE" >/dev/null 2>&1; then
    echo "[android_build] 找不到 $IMAGE；先建：" >&2
    echo "  docker build --network host -t $IMAGE docker/android" >&2
    exit 1
fi

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network "$NET" --memory 8g --cpus 2 --pids-limit 512 \
    -v "$REPO_ROOT:/src" \
    -v wl-gomod:/gomod -v wl-gobuild:/gocache -v wl-gradle:/gradle \
    -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e GOCACHE=/gocache -e GOMODCACHE=/gomod \
    -e GRADLE_USER_HOME=/gradle \
    -e LANG=C.UTF-8 -e LC_ALL=C.UTF-8 \
    -e ABIS="$ABIS" \
    -w /src "$IMAGE" bash -c '
set -euo pipefail
export PATH=/usr/local/go/bin:$PATH

# ⭐ ebitenmobile 要用 `go install pkg@version` 裝，**不能 `go run`**：
# 它自己相依 `golang.org/x/tools`，而那不是遊戲的相依。`go run` 會要求
# 主模組的 `go.sum` 收下它，等於把建置工具的相依混進遊戲的相依圖。
# `@version` 形式的 install 在主模組之外解析，主模組一個字都不會動。
# 版本從 `go.mod` 讀，不寫死——工具與函式庫版本不一致會出很難查的錯。
EBITEN_VER="$(go list -m -f "{{.Version}}" github.com/hajimehoshi/ebiten/v2)"
echo "── 安裝 ebitenmobile $EBITEN_VER ──"
GOBIN=/tmp/bin go install "github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@$EBITEN_VER"

echo "── ebitenmobile bind（$ABIS）──"
mkdir -p android/app/libs
/tmp/bin/ebitenmobile bind \
    -target "$ABIS" \
    -androidapi 29 \
    -javapkg com.wicanr2.wolong.mobile \
    -o android/app/libs/wolong.aar \
    ./mobile/wolong

echo "── gradle assembleDebug ──"
cd android
/opt/gradle/bin/gradle --no-daemon --console=plain assembleDebug
'

APK="$REPO_ROOT/android/app/build/outputs/apk/debug/app-debug.apk"
ls -l "$APK"
echo "$APK"
