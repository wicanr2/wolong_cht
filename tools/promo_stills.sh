#!/usr/bin/env bash
# 重拍推廣主片的靜態段（`tools/promo_video.sh` 的 make_image 來源）。
#
#   tools/promo_stills.sh [輸出目錄]        # 預設 docs/images
#
# ⚠ **每次重剪推廣片都要重拍這五張。** 不重拍的話，片子裡會同時出現
# 好幾個世代的 UI——2026-08-26 踩過：片中同時有純藍對話框、舊版 HUD
# 與當時的藍底龍紋（`docs/spec/88` §1.1、`docs/promo/README.md`）。
#
# ⚠ **不要借 `docs/playtest/` 的證據圖。** 那些圖的 SHA-256 記在
# `docs/re/13`／`docs/playtest/12` 當 fixture 身分，重拍會把紀錄弄壞。
# 推廣片專用的圖一律 `promo-*.png`，只有這一支會動它們。
#
# 五張的旗標寫在這裡是**這一支存在的理由**：先前只有圖進版控，
# 產生它們的命令沒有留下來，於是「重拍」得先重新猜一次旗標。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"
OUT="${1:-docs/images}"
mkdir -p "$OUT"

# 固定的開局：曹操／許昌、196年4月1日、固定亂數種子。
BASE="-direct -scenario 0 -player 0 -seed 7"

# 名稱|旗標
stills=(
"promo-talk|-open-message"
"promo-finance|-open-finance -finance-amount 0"
"promo-cityinfo|-open-cityinfo -1"
"promo-list|-open-list"
"promo-advise|-advise-sortie"
)

for s in "${stills[@]}"; do
    name=${s%%|*}; args=${s#*|}
    echo "── 拍 $name ──"
    tools/parity_shot.sh "$OUT/$name.png" $BASE $args >/dev/null
done

echo "靜態素材 → $OUT"
sha256sum "$OUT"/promo-*.png
