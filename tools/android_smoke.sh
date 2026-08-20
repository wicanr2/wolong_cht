#!/usr/bin/env bash
# Android smoke：起模擬器 → 裝 APK → 推原版資料 → 抓指紋與截圖。
#
#   tools/android_smoke.sh [輸出目錄]
#
# 驗的是三件事（docs/mobile/android-plan.md §5–§6）：
#
#   A. **核心真的在 Android 上跑**——logcat 抓得到 `WOLONG_FP`，
#      而且與桌面同幀的指紋**一模一樣**（docs/spec/69）。
#   G. 安裝、啟動、橫向鎖定不崩。
#   UI. 截一張圖，肉眼看版面（模擬器驗不到真實 DPI 的可讀性）。
#
# ⚠ **smoke 期間不可以送任何觸控**。指紋只在「完全沒有輸入」時兩邊可比，
# 點一下畫面就會改到選取狀態以外的東西（例如指令列），比出來的差異沒有意義。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${1:-$REPO_ROOT/workplace/android}"
IMAGE="${WOLONG_EMULATOR_IMAGE:-wolong-android-emulator:20260820}"
APK="$REPO_ROOT/android/app/build/outputs/apk/debug/app-debug.apk"
PKG=com.wicanr2.wolong

[ -f "$APK" ] || { echo "[smoke] 找不到 APK，先跑 tools/android_build.sh" >&2; exit 1; }
[ -e /dev/kvm ] || { echo "[smoke] 沒有 /dev/kvm，模擬器會慢到不能用" >&2; exit 1; }
mkdir -p "$OUT_DIR"

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --device /dev/kvm --memory 6g --cpus 2 --pids-limit 512 \
    -v "$REPO_ROOT:/src:ro" -v "$OUT_DIR:/out" \
    -e HOME=/tmp -e FP_FRAMES="${WOLONG_FP_FRAMES:-1,60,180}" \
    -w /src "$IMAGE" bash -c '
set -euo pipefail
export PATH=$ANDROID_HOME/emulator:$ANDROID_HOME/platform-tools:$PATH
PKG=com.wicanr2.wolong

echo "── 起模擬器 ──"
emulator -avd wolong -no-window -no-audio -no-boot-anim -no-snapshot \
    -gpu swiftshader_indirect -camera-back none -camera-front none \
    >/tmp/emulator.log 2>&1 &

adb wait-for-device
# `sys.boot_completed` 之前 pm 還沒起來，install 會失敗。
for i in $(seq 1 180); do
    [ "$(adb shell getprop sys.boot_completed 2>/dev/null | tr -d "\r")" = "1" ] && break
    sleep 2
done
adb shell settings put global window_animation_scale 0 || true
# ⚠ 關掉 ANR／崩潰對話框。機器忙的時候 System UI 會被判定沒回應，
# 跳出來的對話框會搶走焦點，Ebiten 隨即因為 **GL context 遺失**自殺
#（logcat 只留一行 `The application was killed due to context loss`）。
# 那是驗收環境的雜訊，不是 app 的缺陷。
adb shell settings put global hide_error_dialogs 1 || true
# ⚠ **螢幕不能睡著**。模擬器在這台機器上只跑到個位數 fps，
# 跑滿 600 幀要好幾分鐘，中途螢幕一關 surface 就沒了，Ebiten 隨即自殺
#（`The application was killed due to context loss`）。
# 那看起來像 app 崩潰，其實是驗收環境把它的畫布收走了。
adb shell svc power stayon true || true
adb shell settings put secure lockscreen.disabled 1 || true
adb shell input keyevent KEYCODE_WAKEUP || true

# ⚠ `sys.boot_completed` **不代表 /sdcard 可用**。外部儲存是 FUSE 掛的，
# 開機旗標亮起之後還要一段時間才掛好，這段期間對它下 mkdir 會拿到
# `Transport endpoint is not connected`——那不是權限問題，是還沒掛上。
# 機器忙的時候這段特別長，所以要等到真的讀得到為止。
for i in $(seq 1 120); do
    adb shell "ls /sdcard/Android >/dev/null 2>&1" && break
    sleep 2
done

# ⚠ 關掉安裝時的完整性驗證。機器忙的時候它會逾時，
# 症狀是 `INSTALL_FAILED_VERIFICATION_FAILURE: Integrity verification timed out`
# ——那是**驗收環境太慢**，不是 APK 有問題。
adb shell settings put global package_verifier_enable 0 || true
adb shell settings put global verifier_verify_adb_installs 0 || true

echo "── 安裝 ──"
for i in 1 2 3; do
    adb install -r -g /src/android/app/build/outputs/apk/debug/app-debug.apk && break
    echo "（第 $i 次安裝失敗，重試）"
    sleep 10
done

# ⭐ 資料要進 **app 的內部目錄**，不是 /sdcard。
# `/sdcard/Android/data/<pkg>/` 看起來剛好是給這個 app 的，但 Android 11 以上
# 它是 FUSE 掛的：adb 寫進去的檔案 app 讀不到（permission denied），
# 而目錄還列得出來——症狀像「檔案在那裡但壞掉」。
# debuggable 的 APK 可以用 `run-as` 以 app 身分複製，這條路徑乾淨。
STAGE=/data/local/tmp/wolong
adb shell rm -rf "$STAGE"
adb shell mkdir -p "$STAGE/orig" "$STAGE/eten"

echo "── 推原版資料（唯讀來源，不進 APK）──"
push_flat() {
    # 目的地要在最後一個參數，所以借 sh 的 $0 把它帶到 xargs 的另一端。
    find "$1" -maxdepth 1 -type f -print0 |
        xargs -0 -r sh -c "adb push \"\$@\" \"\$0\" >/dev/null" "$2"
}
push_flat /src/workplace/orig/dosv "$STAGE/orig"
push_flat /src/workplace/eten "$STAGE/eten"

# ⚠ `run-as` 的工作目錄**不是** app 的私有目錄（是 `/`，唯讀），
# 而且 `adb shell` 會把參數重新拼成一行——多行的 `sh -c` 會被拆散。
# 兩件事合起來的症狀是 `mkdir: files: Read-only file system`。
# 所以：整段包成一行、路徑全部寫絕對路徑。
FILES=/data/data/$PKG/files
adb shell "run-as $PKG sh -c \"mkdir -p $FILES/orig $FILES/eten && cp $STAGE/orig/* $FILES/orig/ && cp $STAGE/eten/* $FILES/eten/\""
echo "orig 檔數：$(adb shell "run-as $PKG ls $FILES/orig" | wc -l)"

echo "── 重新啟動，抓指紋 ──"
adb shell am force-stop "$PKG"
adb logcat -c
adb shell am start --es fp_frames "$FP_FRAMES" -n "$PKG/.MainActivity" >/dev/null
# 指紋的幀數 ≒ 幀數 ÷ 60 秒，但模擬器不保證跑滿 60 fps，機器忙的時候差很多。
# 等到最後一個幀號出現為止，不要用固定秒數賭它跑到了。
LAST=${FP_FRAMES##*,}
for i in $(seq 1 60); do
    adb logcat -d | grep -q "WOLONG_FP frame=$LAST" && break
    sleep 5
done
adb logcat -d > /out/logcat.txt
adb exec-out screencap -p > /out/screen.png
grep -E "WOLONG_FP|FATAL|AndroidRuntime" /out/logcat.txt | tail -20 || true
adb emu kill >/dev/null 2>&1 || true
'
echo "輸出：$OUT_DIR/logcat.txt、$OUT_DIR/screen.png"
