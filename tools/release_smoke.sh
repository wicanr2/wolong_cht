#!/usr/bin/env bash
# 對 dist-all 的 Linux 產物做 GUI smoke，並把截圖放進 dist-all/verification。
#
#   tools/release_smoke.sh [版本日期]      # 預設今天，格式 YYYYMMDD
#
# 為什麼要一支腳本：`tools/release_all.sh` 的 `promote` 會把整個 dist-all
# 換掉，**verification/ 的截圖每次都會被清掉**。手動補會忘記補；
# 有腳本才知道那幾張是怎麼來的、怎麼重來一次。
#
# 邊界：只用 2 顆 CPU、一律 --rm、原版資料唯讀掛載、只寫 dist-all/verification。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STAMP="${1:-$(date +%Y%m%d)}"
ORIG="${WOLONG_ORIG_DIR:-$REPO_ROOT/workplace/orig/dosv}"
OUT="$REPO_ROOT/dist-all/verification"
APPIMAGE="wolong-remake-linux-amd64-${STAMP}.AppImage"
TARBALL="wolong-remake-linux-amd64-${STAMP}.tar.gz"

for f in "$APPIMAGE" "$TARBALL"; do
    [ -f "$REPO_ROOT/dist-all/packages/$f" ] || { echo "找不到 dist-all/packages/$f" >&2; exit 1; }
done
mkdir -p "$OUT"

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 2g --cpus 2 --pids-limit 256 \
    -v "$REPO_ROOT/dist-all/packages:/pkg:ro" -v "$ORIG:/orig:ro" -v "$OUT:/out" \
    -u "$(id -u):$(id -g)" -e HOME=/tmp -e LIBGL_ALWAYS_SOFTWARE=1 \
    -w /tmp "${WOLONG_GO_IMAGE:-demonwinter-go}" bash -lc "
set -e
Xvfb :99 -screen 0 1600x900x24 >/tmp/xvfb.log 2>&1 &
for i in \$(seq 1 50); do xdpyinfo -display :99 >/dev/null 2>&1 && break; sleep 0.1; done
cp /pkg/${APPIMAGE} /tmp/app.AppImage && chmod +x /tmp/app.AppImage
common='-orig /orig -direct -scenario 0 -player 0 -seed 7'
DISPLAY=:99 timeout 180 /tmp/app.AppImage --appimage-extract-and-run \$common -shot /out/appimage-smoke-${STAMP}.png
DISPLAY=:99 timeout 180 /tmp/app.AppImage --appimage-extract-and-run \$common -open-ending 0 -shot /out/appimage-ending-${STAMP}.png
mkdir -p /tmp/tarpkg && tar -xzf /pkg/${TARBALL} -C /tmp/tarpkg
BIN=\$(find /tmp/tarpkg -name wlgame -type f | head -1)
DISPLAY=:99 timeout 180 \"\$BIN\" \$common -shot /out/linux-tar-smoke-${STAMP}.png
# ⭐ 內含遊戲檔案的完整包還要驗**不帶資料旗標**那條路（docs/spec/72 §3）：
# 從包的根目錄跑起來，wlgame 要自己找到旁邊的 gamedata/ 與 fonts/。
# ⚠ 這一條的失敗方式很安靜——載不到字型只會噴一行警告，畫面照常出來，
# 中文變成方框。所以要**同時**看 exit code 與那一行警告。
if [ -d /tmp/tarpkg/*/gamedata ] 2>/dev/null || ls -d /tmp/tarpkg/*/gamedata >/dev/null 2>&1; then
    cd \"\$(dirname \"\$BIN\")\"
    DISPLAY=:99 timeout 180 ./wlgame -direct -scenario 0 -player 0 -seed 7 \
        -shot /out/bundled-nodflags-${STAMP}.png 2>/tmp/bundled.log
    cat /tmp/bundled.log
    grep -q '載不到' /tmp/bundled.log && { echo '⚠ 完整包裡的資料沒被自己找到' >&2; exit 1; }
    echo '✓ 不帶任何資料旗標就跑得起來，字型也載到了'
fi
"

echo "截圖 → dist-all/verification/{appimage-smoke,appimage-ending,linux-tar-smoke}-${STAMP}.png"
echo "⚠ 截圖加進去之後 SHA256SUMS.txt 要重算：tools/py.sh tools/release_all_fs.py refresh"
