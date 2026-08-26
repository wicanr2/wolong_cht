#!/bin/sh
# 產生 60 秒 remake 推廣片。此腳本只使用 remake 自己的畫面與本專案原創合成聲，
# 應在 Docker 的 ffmpeg 容器內執行；不讀取原版影片、BGM 或 SOUND.DAT。
#
# ⭐ **大地圖與兩場戰鬥是真的在跑的遊戲**，不是截圖。素材先由
# `tools/promo_live_capture.sh` 逐幀錄出來（docs/spec/71），這裡只負責剪。
# 選單類畫面（事件視窗、校訂、存檔）本來就不動，仍用截圖。
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_dir"
out=${1:-dist/promo/wolong-remake-trailer.mp4}
font=${PROMO_FONTFILE:-/fonts/truetype/moe/MoeStandardSong.ttf}

if [ ! -f "$font" ]; then
    echo "找不到推廣片字型：$font" >&2
    exit 1
fi

mkdir -p "$(dirname -- "$out")"
tmp=$(mktemp -d /tmp/wolong-promo-XXXXXX)
cleanup() {
    rm -rf "$tmp"
}
trap cleanup EXIT INT TERM

write_text() {
    path=$1
    text=$2
    printf '%s' "$text" >"$path"
}

make_card() {
    name=$1
    duration=$2
    title=$3
    subtitle=$4
    title_file="$tmp/$name-title.txt"
    subtitle_file="$tmp/$name-subtitle.txt"
    write_text "$title_file" "$title"
    write_text "$subtitle_file" "$subtitle"
    ffmpeg -hide_banner -loglevel error -y \
        -f lavfi -i "color=c=0x07142e:s=1280x720:r=30:d=$duration" \
        -vf "drawtext=fontfile=$font:textfile=$title_file:fontcolor=white:fontsize=72:x=(w-text_w)/2:y=285,drawtext=fontfile=$font:textfile=$subtitle_file:fontcolor=0xb8c8e8:fontsize=30:x=(w-text_w)/2:y=395,fade=t=in:st=0:d=0.18,fade=t=out:st=$(awk -v d="$duration" 'BEGIN{printf "%.2f",d-0.18}'):d=0.18" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
        "$tmp/$name.mp4"
}

make_image() {
    name=$1
    duration=$2
    source=$3
    label=$4
    label_file="$tmp/$name-label.txt"
    write_text "$label_file" "$label"
    ffmpeg -hide_banner -loglevel error -y \
        -loop 1 -framerate 30 -t "$duration" -i "$repo_dir/$source" \
        -vf "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2:color=0x07142e,drawtext=fontfile=$font:textfile=$label_file:fontcolor=white:fontsize=28:box=1:boxcolor=0x061128@0.82:boxborderw=12:x=30:y=28,fade=t=in:st=0:d=0.18,fade=t=out:st=$(awk -v d="$duration" 'BEGIN{printf "%.2f",d-0.18}'):d=0.18" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
        "$tmp/$name.mp4"
}

# make_frames 把一段逐幀錄下來的 PNG 接成影片。
#
# ⚠ `-start_number 0` 不能省：檔名從 f00000 開始，ffmpeg 預設從 1 找起，
# 找不到就整段落空——而落空的症狀是「少了一段」不是「報錯」。
# ⚠ `flags=neighbor`：點陣圖放大不能內插，糊掉就看不出是原版美術了。
make_frames() {
    name=$1
    duration=$2
    dir=$3
    label=$4
    label_file="$tmp/$name-label.txt"
    write_text "$label_file" "$label"
    ffmpeg -hide_banner -loglevel error -y \
        -framerate 30 -start_number 0 -i "$repo_dir/$dir/f%05d.png" -t "$duration" \
        -vf "scale=1280:720:flags=neighbor,drawtext=fontfile=$font:textfile=$label_file:fontcolor=white:fontsize=28:box=1:boxcolor=0x061128@0.82:boxborderw=12:x=30:y=28,fade=t=in:st=0:d=0.18,fade=t=out:st=$(awk -v d="$duration" 'BEGIN{printf "%.2f",d-0.18}'):d=0.18" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
        "$tmp/$name.mp4"
}

# make_compare 把原版實機錄影與 remake 的逐幀素材並排。
#
# ⚠ **兩側不是同一局面、同一輸入**，也不是同一個時鐘速度（remake 這一段是
# 最高速檔錄的）。片上要寫清楚，不然並排本身會被讀成「逐像素一致」的宣稱。
# 逐像素的判定在 docs/playtest/37 與 40，不在推廣片裡。
make_compare() {
    name=$1
    duration=$2
    original=$3
    dir=$4
    label=$5
    label_file="$tmp/$name-label.txt"
    write_text "$label_file" "$label"
    left_file="$tmp/$name-left.txt"
    right_file="$tmp/$name-right.txt"
    # ⚠ 原版擷取是 640×480，上下各 40 px 黑邊；不裁掉的話兩側等比縮放後
    # 原版那一格會小一圈，看起來像「remake 比較大」而不是「錄影帶黑邊」。
    write_text "$left_file" "原版　松崗 DOS/V 實機"
    write_text "$right_file" "remake"
    ffmpeg -hide_banner -loglevel error -y \
        -i "$repo_dir/$original" \
        -framerate 30 -start_number 0 -i "$repo_dir/$dir/f%05d.png" \
        -filter_complex "\
[0:v]crop=640:400:0:40,scale=640:400:flags=neighbor,\
drawtext=fontfile=$font:textfile=$left_file:fontcolor=0xb8c8e8:fontsize=20:box=1:boxcolor=0x061128@0.85:boxborderw=8:x=10:y=372[l];\
[1:v]scale=640:400:flags=neighbor,\
drawtext=fontfile=$font:textfile=$right_file:fontcolor=0xb8c8e8:fontsize=20:box=1:boxcolor=0x061128@0.85:boxborderw=8:x=10:y=372[r];\
[l][r]hstack=inputs=2[stack];\
[stack]pad=1280:720:0:160:color=0x07142e,\
drawtext=fontfile=$font:textfile=$label_file:fontcolor=white:fontsize=26:box=1:boxcolor=0x061128@0.82:boxborderw=12:x=(w-text_w)/2:y=612,\
fade=t=in:st=0:d=0.18,fade=t=out:st=$(awk -v d="$duration" 'BEGIN{printf "%.2f",d-0.18}'):d=0.18[v]" \
        -map "[v]" -t "$duration" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
        "$tmp/$name.mp4"
}

live=${WOLONG_PROMO_LIVE:-workplace/promo/live}
# 原版側：自己跑的松崗 DOS/V 實機擷取（tools/dosv_capture.sh，
# 紀錄在 docs/promo/dosv-realmachine.md）。**不是 YouTube 錄影。**
original_clip=${WOLONG_PROMO_ORIGINAL:-workplace/promo-live/promo-dosv-main/o2-strategy.mp4}
if [ ! -f "$repo_dir/$original_clip" ]; then
    echo "找不到原版擷取 $original_clip —— 見 docs/promo/dosv-realmachine.md §6" >&2
    exit 1
fi
for seg in map battle siege lang-zh-hant lang-zh-hans lang-ja lang-en; do
    if [ ! -f "$repo_dir/$live/$seg/f00000.png" ]; then
        echo "找不到動態素材 $live/$seg —— 先跑 tools/promo_live_capture.sh" >&2
        exit 1
    fi
done

# 段落長度加起來要**剛好等於配樂長度**（tools/promo_score.py 的 DURATION）。
# 對不上的症狀是音樂的轉折落在畫面中間，說不出哪裡怪，只覺得鬆散。
make_card title 4 "臥龍傳 Remake" "DOS/V 對齊・四語系・跨平台"
make_frames natural 10 "$live/map" "大地圖：時鐘在走，兩軍在行軍"
# ⚠ **靜態圖一律用 `promo-*.png`，不要借 playtest 的證據圖。**
# 那些圖的 SHA-256 記在 `docs/re/13`／`docs/playtest/12` 裡當 fixture 身分，
# 重拍會把紀錄弄壞；而不重拍，片子裡就會同時出現三個世代的 UI
#（docs/spec/88 §1.1 踩過：2026-08-10 的圖與現行版本的對話框長得不一樣）。
make_image choice 2.5 docs/images/promo-talk.png "事件訊息：肖像 ＋ TALK.DAT 原文"
make_image event5 2 docs/images/promo-finance.png "財政：稅率、徵兵與數值輸入"
make_frames battle 12 "$live/battle" "戰場：45 度視角，兩軍實時交戰"
make_frames siege 8 "$live/siege" "攻城：城牆、城門與守軍"
make_image result 2.5 docs/images/promo-cityinfo.png "據點情報：生產力、上昇率、防災"
# 語言：同一顆種子、同一個時間點的四種語言。F9 在遊戲中即時切換，
# 切出來的畫面與 `-lang` 啟動**逐像素相同**（docs/playtest/46），
# 所以這四段用四次啟動錄，換得可重現。
make_frames langA 2.2 "$live/lang-zh-hant" "F9 切換語言　繁體中文"
make_frames langB 2.2 "$live/lang-zh-hans" "F9 切換語言　简体中文"
make_frames langC 2.2 "$live/lang-ja" "F9 切換語言　日本語（PC-98 原文）"
make_frames langD 2.4 "$live/lang-en" "F9 切換語言　English"
make_compare compare 10 "$original_clip" "$live/map" \
    "並排：不同局面、不同時鐘速度　逐像素判定見 docs/playtest/37・40"
make_image m7 2.5 docs/images/promo-list.png "一覽表：武將、據點、勢力、軍團"
make_image save 2.5 docs/images/promo-advise.png "進言：玩家是軍師，指令要先過君主那一關"
make_card end 7 "立即開始你的三國" "Linux・Windows・macOS・Android"

concat="$tmp/concat.txt"
: >"$concat"
for name in title natural choice event5 battle siege result langA langB langC langD compare m7 save end; do
    printf "file '%s'\n" "$tmp/$name.mp4" >>"$concat"
done

# 配樂：本專案原創合成。ffmpeg 容器裡沒有 python3，所以允許先在
# `tools/py.sh` 那邊產好再用 PROMO_SCORE 指進來。
if [ -n "${PROMO_SCORE:-}" ] && [ -f "$repo_dir/$PROMO_SCORE" ]; then
    cp "$repo_dir/$PROMO_SCORE" "$tmp/score.wav"
else
    python3 "$repo_dir/tools/promo_score.py" "$tmp/score.wav" >/dev/null
fi
ffmpeg -hide_banner -loglevel error -y \
    -f concat -safe 0 -i "$concat" -i "$tmp/score.wav" \
    -map 0:v:0 -map 1:a:0 -t 72 \
    -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
    -c:a aac -b:a 160k -ar 44100 -movflags +faststart "$out"

ffprobe -v error -show_entries format=duration:stream=codec_name,width,height,channels,sample_rate \
    -of default=noprint_wrappers=1 "$out"
echo "推廣片：$out"
