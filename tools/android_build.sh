#!/usr/bin/env bash
# 建 Android APK：ebitenmobile bind → AAR → gradle assembleDebug。
#
#   tools/android_build.sh              # 兩個 ABI（arm64 + amd64）
#   WOLONG_ABIS=android/arm64 tools/android_build.sh
#   WOLONG_NET=none tools/android_build.sh   # 相依都在快取裡時可以斷網跑
#   WOLONG_BUNDLE_DATA=1 tools/android_build.sh  # 把原版資料內嵌進 APK
#
# ⚠ **`WOLONG_BUNDLE_DATA=1` 建出來的 APK 內含原版資料與倚天字型，
# 絕對不可外流**（docs/spec/72）。預設不內嵌——內嵌要明講，
# 而且建完就把 assets 清掉，下一次建置不會默默沿用上一次的資料。
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
BUNDLE="${WOLONG_BUNDLE_DATA:-0}"
ASSETS="$REPO_ROOT/android/app/src/main/assets"
ORIG_DIR="${WOLONG_ORIG_DIR:-$REPO_ROOT/workplace/orig/dosv}"
FONT_DIR="${WOLONG_FONT_DIR:-$REPO_ROOT/workplace/eten}"
AUDIO_DIR="${WOLONG_AUDIO_DIR:-$REPO_ROOT/workplace/audio}"

# ⭐ **每一次都先清乾淨。** 沒有這一步，上一次 `WOLONG_BUNDLE_DATA=1`
# 留下的資料會被下一次的一般建置默默包進去——而那個 APK 看起來
# 與正常的一模一樣。
rm -rf "$ASSETS"

if [ "$BUNDLE" = 1 ]; then
    [ -d "$ORIG_DIR" ] || { echo "[android_build] 找不到原版資料 $ORIG_DIR" >&2; exit 1; }
    [ -f "$ORIG_DIR/SINARIO.DAT" ] || { echo "[android_build] $ORIG_DIR 裡沒有 SINARIO.DAT" >&2; exit 1; }
    mkdir -p "$ASSETS/gamedata/orig" "$ASSETS/gamedata/eten"
    cp "$ORIG_DIR"/* "$ASSETS/gamedata/orig/"
    if [ -d "$FONT_DIR" ]; then
        cp "$FONT_DIR"/* "$ASSETS/gamedata/eten/"
    else
        echo "[android_build] ⚠ 沒有 $FONT_DIR，APK 裡不會有字型，中文會是方框" >&2
    fi
    # ⭐ 音樂與音效（docs/spec/92）。**只收 ogg**——`workplace/audio` 裡
    # ogg 旁邊躺著合成中間產物 wav，整包 239 MB 而 ogg 只有 19 MB。
    # ⚠ 這裡放進 assets 還不夠：`ImportActivity.unpackBundled()` 的子目錄
    # 陣列**就是唯一的清單**，漏了 `audio` 就是進得去、解不出來，
    # 而畫面上什麼都看不出來。
    if [ -d "$AUDIO_DIR" ] && ls "$AUDIO_DIR"/*.ogg >/dev/null 2>&1; then
        mkdir -p "$ASSETS/gamedata/audio"
        cp "$AUDIO_DIR"/*.ogg "$ASSETS/gamedata/audio/"
    else
        echo "[android_build] ⚠ 沒有 $AUDIO_DIR 的 ogg，APK 裡不會有音樂（先跑 tools/bgm2ogg.sh）" >&2
    fi
    echo "── 內嵌原版資料：$(find "$ASSETS/gamedata" -type f | wc -l) 個檔 ──"
fi

# 不論成功失敗都把 assets 清掉：原版資料不留在工作區裡。
cleanup_assets() { rm -rf "$ASSETS"; }
trap cleanup_assets EXIT

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
# ⚠ **`.so` 要 16 KB 對齊**。Android 15 起有 16 KB page size 的裝置，
# 而 Go 產出的 `libgojni.so` 預設是 4 KB（LOAD 段 align=0x1000）——
# 那種 `.so` 在 16 KB 的機器上**載不起來**，症狀是一啟動就掛，
# 而 4 KB 的機器上完全正常，所以測不出來。
# `zipalign -P 16` 驗的是 zip 那一層，**驗不到 ELF 這一層**。
/tmp/bin/ebitenmobile bind \
    -ldflags "-extldflags=-Wl,-z,max-page-size=16384" \
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

# ⭐ **對齊要驗，不能只靠旗標**。旗標打錯字、工具鏈換版、`ebitenmobile`
# 換掉傳遞方式，任何一種都會讓對齊悄悄掉回 4 KB——而 4 KB 的機器上一切正常，
# 只有 16 KB page size 的裝置會一啟動就掛。
docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 1 --pids-limit 64 \
    -v "$REPO_ROOT:/src:ro" -u "$(id -u):$(id -g)" -e HOME=/tmp \
    "$IMAGE" bash -c '
set -euo pipefail
cd /tmp && rm -rf apkcheck
unzip -o -q /src/android/app/build/outputs/apk/debug/app-debug.apk "lib/*" -d apkcheck
bad=0
for f in apkcheck/lib/*/*.so; do
    for a in $(readelf -lW "$f" | awk "/LOAD/ {print \$NF}"); do
        if [ "$a" != "0x4000" ]; then
            echo "⚠ $f 的 LOAD 段對齊是 $a，不是 16 KB（0x4000）" >&2
            bad=1
        fi
    done
done
[ "$bad" = 0 ] || { echo "16 KB page size 的裝置載不起這個 .so" >&2; exit 1; }
echo "✓ .so 的 LOAD 段都是 16 KB 對齊"
$ANDROID_HOME/build-tools/35.0.0/zipalign -c -P 16 -v 4 \
    /src/android/app/build/outputs/apk/debug/app-debug.apk >/dev/null
echo "✓ zip 內的 .so 也照 16 KB 對齊"
'

# ⭐ **內嵌與否要驗，不能只靠旗標。** 兩個方向都會安靜地出錯：
# 該內嵌卻沒進去（玩家裝上去打不開），或不該內嵌卻混進去
# （把原版資產送出門）。兩種 APK 從外觀完全分不出來。
docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 1g --cpus 1 --pids-limit 64 \
    -v "$REPO_ROOT:/src:ro" -u "$(id -u):$(id -g)" -e HOME=/tmp \
    -e BUNDLE="$BUNDLE" \
    "$IMAGE" bash -c '
set -euo pipefail
list=$(unzip -Z1 /src/android/app/build/outputs/apk/debug/app-debug.apk \
    "assets/gamedata/*" 2>/dev/null | wc -l)
if [ "$BUNDLE" = 1 ]; then
    [ "$list" -gt 60 ] || { echo "⚠ 說要內嵌，APK 裡卻只有 $list 個資料檔" >&2; exit 1; }
    echo "✓ APK 內嵌了 $list 個原版資料檔（⚠ 這個 APK 不可外流）"
else
    [ "$list" = 0 ] || { echo "⚠ 沒說要內嵌，APK 裡卻有 $list 個原版資料檔" >&2; exit 1; }
    echo "✓ APK 不含原版資產"
fi
'

ls -l "$APK"
echo "$APK"
