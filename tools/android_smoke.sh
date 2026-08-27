#!/usr/bin/env bash
# Android smoke：起模擬器 → 裝 APK → 推原版資料 → 抓指紋與截圖。
#
#   tools/android_smoke.sh [輸出目錄]
#   WOLONG_BUNDLED=1 tools/android_smoke.sh   # 驗內嵌資料的 APK
#
# 驗的是三件事（docs/mobile/android-plan.md §5–§6）：
#
#   A. **核心真的在 Android 上跑**——logcat 抓得到 `WOLONG_FP`，
#      而且與桌面同幀的指紋**一模一樣**（docs/spec/69）。
#   G. 安裝、啟動、橫向鎖定不崩。
#   UI. 截一張圖，肉眼看版面（模擬器驗不到真實 DPI 的可讀性）。
#
# ⭐ `WOLONG_BUNDLED=1` 換掉的是**入口那一段**（docs/spec/72 §4）：
# 內嵌資料的 APK 在乾淨安裝之後應該**自己解開並直接進遊戲**，
# 而不是停在匯入畫面。所以那一版驗的是「有沒有轉到 MainActivity」，
# 而且**不推任何資料**——推了就分不出「內嵌生效」與「我剛推進去的」。
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
    -e HOME=/tmp -e FP_FRAMES="${WOLONG_FP_FRAMES:-1,60,120}" \
    -e BUNDLED="${WOLONG_BUNDLED:-0}" \
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
# ⚠ `svc power stayon true` **擋不住這台機器上的睡眠**。實測指紋那一段
# （約 24 秒）跑到一半螢幕就關了並上鎖，logcat 留下
# `LockscreenSmartspaceController: Starting smartspace session for lockscreen`；
# 之後叫醒只會回到**桌面**，截到的是 Android 的主畫面而不是遊戲——
# 而那張圖是合法的 PNG、大小也正常，**看起來完全像截成功了**。
adb shell settings put system screen_off_timeout 1800000 || true
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

# ⭐ 先在**沒有資料**的狀態下啟動一次。
#   一般 APK：該看到匯入畫面（里程碑 B 的入口，docs/mobile/android-plan.md §3）。
#   內嵌 APK：該**自己解開並轉進遊戲**（docs/spec/72 §4）。
adb shell pm clear "$PKG" >/dev/null || true
adb logcat -c
adb shell am start -n "$PKG/.ImportActivity" >/dev/null
if [ "$BUNDLED" = 1 ]; then
    WANT=".MainActivity"
else
    WANT=".ImportActivity"
fi
# ⚠ 這台機器上第一次啟動要十幾秒才畫得出來（logcat 的 `Displayed` 是
# +14s）。截太早只會拍到系統的啟動畫面，看起來像「匯入畫面沒出現」。
for i in $(seq 1 30); do
    adb logcat -d | grep -q "Displayed $PKG/$WANT" && break
    sleep 3
done
adb logcat -d | grep -q "Displayed $PKG/$WANT" || {
    echo "⚠ 等不到 Displayed $PKG/$WANT" >&2
    adb logcat -d | tail -40 >&2
    exit 1
}
sleep 3
# ⚠ 截圖前先叫醒螢幕。這台機器上一輪 smoke 要跑好幾分鐘，
# 中途螢幕會睡著——**睡著時 screencap 拍到的是全黑**，
# 看起來像「畫面沒畫出來」，實際上 logcat 明明寫著 Displayed。
wake() {
    # ⚠ **等到真的醒著再拍。** `KEYCODE_WAKEUP` 是非同步的，送完就 screencap
    # 有機會拍到還在睡的那一瞬間。`mWakefulness` 是唯一講實話的那個值。
    # ⚠ **grep 要在裝置上跑，不要在 host 這一端接管線。**
    # 這支腳本開著 `set -o pipefail`，而 `grep -q` 命中就提早結束，
    # 上游的 `tr` 吃到 SIGPIPE（141）——**整條 pipeline 於是回非零，
    # 明明命中了卻被判成沒命中**。實測 `mWakefulness=Awake`、
    # Display State=ON，而這個檢查照樣報「叫不醒」（2026-08-22）。
    for i in $(seq 1 8); do
        adb shell input keyevent KEYCODE_WAKEUP >/dev/null 2>&1 || true
        adb shell wm dismiss-keyguard >/dev/null 2>&1 || true
        adb shell svc power stayon true >/dev/null 2>&1 || true
        sleep 2
        state=$(adb shell "dumpsys power | grep -m1 mWakefulness=" 2>/dev/null | tr -d "\r" || true)
        case "$state" in *Awake*) return 0 ;; esac
    done
    echo "⚠ 螢幕叫不醒（最後看到：$state）" >&2
}

# ⭐ 全黑的截圖要當成失敗，不能寫進輸出目錄就算數。
# 睡著時 `screencap` 回的是一張合法的 PNG，大小固定約 12 KB——
# **它看起來完全像一張成功的截圖**，只有內容是空的。
# 沒有輸出要長得像沒有輸出。
shoot() {
    want="${2:-}"
    for i in 1 2 3 4; do
        wake
        # ⭐ **前景是誰要查**。螢幕睡過一輪再叫醒會停在桌面，
        # 那張截圖不黑、大小也正常，只是拍到別的東西。
        if [ -n "$want" ] &&
           ! adb shell "dumpsys window | grep -m1 mCurrentFocus" 2>/dev/null |
             tr -d "\r" | grep -q "$want"; then
            adb shell am start -n "$want" >/dev/null 2>&1 || true
            sleep 4
        fi
        adb exec-out screencap -p > "$1"
        [ "$(wc -c < "$1")" -gt 102400 ] && return 0
        echo "（第 $i 次截到全黑，重試）" >&2
        sleep 3
    done
    echo "⚠ $1 連續四次都是全黑——這一張不是有效證據" >&2
}
adb shell dumpsys window | grep -E "mCurrentFocus|mFocusedApp" || true
shoot /out/import.png "$PKG/$WANT"
adb logcat -d > /out/logcat-import.txt
adb shell am force-stop "$PKG"

# ⭐ 資料要進 **app 的內部目錄**，不是 /sdcard。
# `/sdcard/Android/data/<pkg>/` 看起來剛好是給這個 app 的，但 Android 11 以上
# 它是 FUSE 掛的：adb 寫進去的檔案 app 讀不到（permission denied），
# 而目錄還列得出來——症狀像「檔案在那裡但壞掉」。
# debuggable 的 APK 可以用 `run-as` 以 app 身分複製，這條路徑乾淨。
STAGE=/data/local/tmp/wolong
adb shell rm -rf "$STAGE"
adb shell mkdir -p "$STAGE/orig" "$STAGE/eten"

if [ "$BUNDLED" = 1 ]; then
    # ⚠ **一個 byte 都不推。** 這一段要證明的正是「資料本來就在 APK 裡」，
    # 推了就什麼都證明不了。
    echo "── 內嵌版：不推資料，直接查 app 私有目錄 ──"
    FILES=/data/data/$PKG/files
    n=$(adb shell "run-as $PKG ls $FILES/orig 2>/dev/null" | tr -d "\r" | grep -c . || true)
    echo "APK 自己解出來的 orig 檔數：$n"
    [ "$n" -ge 60 ] || { echo "⚠ 只解出 $n 個檔，內嵌沒生效" >&2; exit 1; }
    adb shell "run-as $PKG ls $FILES/orig" | tr -d "\r" | grep -q SINARIO.DAT || {
        echo "⚠ 解出來的資料裡沒有 SINARIO.DAT" >&2; exit 1; }
    # ⭐ 音檔（docs/spec/92）。**模擬器是 `-no-audio` 起的，聽不到聲音
    # 不代表沒播**——判準是檔案解出來了，不是有沒有聲音。
    a=$(adb shell "run-as $PKG ls $FILES/audio 2>/dev/null" | tr -d "\r" | grep -c "\.ogg" || true)
    echo "APK 自己解出來的 ogg 數：$a"
    [ "$a" -ge 30 ] || echo "⚠ 只解出 $a 個 ogg，音樂不會有聲音" >&2
    # ⚠ `.part` 留下來代表改名那一步沒完成，資料可能是半套的。
    adb shell "run-as $PKG ls $FILES" | tr -d "\r" | grep -q ".part" && {
        echo "⚠ 留下了 .part，解包沒有跑完" >&2; exit 1; } || true
else

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
fi

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
shoot /out/screen.png "$PKG/.MainActivity"
grep -E "WOLONG_FP|FATAL|AndroidRuntime" /out/logcat.txt | tail -20 || true
adb emu kill >/dev/null 2>&1 || true
'
echo "輸出：$OUT_DIR/logcat.txt、$OUT_DIR/screen.png"
