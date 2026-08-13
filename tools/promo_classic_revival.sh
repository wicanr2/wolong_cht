#!/bin/sh
# 製作「經典再現」原版／remake 對照推廣片。
#
# 原版只取使用者提供 YouTube 錄影的代表幀；不讀取、複製或封裝原版影片與資產。
# 對照條件來自 docs/promo/classic-revival.md，影片用途是視覺／流程展示，不是
# 同狀態逐像素 parity 證明。
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
out=${1:-$repo_dir/dist/promo/wolong-remake-classic-revival.mp4}
out_dir=$(dirname -- "$out")
review_dir=${CLASSIC_REVIEW_DIR:-$repo_dir/docs/promo}
font=${PROMO_FONTFILE:-/fonts/NotoSansTC-Regular.otf}

if [ ! -f "$font" ]; then
    font=/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf
fi
if [ ! -f "$font" ]; then
    echo "找不到推廣片字型：$font" >&2
    exit 1
fi

mkdir -p "$out_dir" "$review_dir"
tmp=$(mktemp -d /tmp/wolong-classic-revival-XXXXXX)
cleanup() {
    rm -rf "$tmp"
}
trap cleanup EXIT INT TERM

write_text() {
    path=$1
    text=$2
    printf '%s' "$text" >"$path"
}

fade_filter() {
    duration=$1
    awk -v d="$duration" 'BEGIN { printf "fade=t=in:st=0:d=0.22,fade=t=out:st=%.2f:d=0.22", d-0.22 }'
}

make_title() {
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
        -vf "drawtext=fontfile=$font:textfile=$title_file:fontcolor=white:fontsize=68:x=(w-text_w)/2:y=275,drawtext=fontfile=$font:textfile=$subtitle_file:fontcolor=0xb8c8e8:fontsize=30:x=(w-text_w)/2:y=390,$(fade_filter "$duration")" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

make_compare() {
    name=$1
    duration=$2
    original=$3
    remake=$4
    title=$5
    original_label=$6
    remake_label=$7
    note=$8
    title_file="$tmp/$name-title.txt"
    original_file="$tmp/$name-original.txt"
    remake_file="$tmp/$name-remake.txt"
    note_file="$tmp/$name-note.txt"
    write_text "$title_file" "$title"
    write_text "$original_file" "$original_label"
    write_text "$remake_file" "$remake_label"
    write_text "$note_file" "$note"

    ffmpeg -hide_banner -loglevel error -y \
        -f lavfi -i "color=c=0x07142e:s=1280x720:r=30:d=$duration" \
        -loop 1 -framerate 30 -t "$duration" -i "$repo_dir/$original" \
        -loop 1 -framerate 30 -t "$duration" -i "$repo_dir/$remake" \
        -filter_complex "[0:v]drawtext=fontfile=$font:textfile=$title_file:fontcolor=white:fontsize=38:x=(w-text_w)/2:y=26,drawtext=fontfile=$font:textfile=$note_file:fontcolor=0xb8c8e8:fontsize=22:x=(w-text_w)/2:y=650,$(fade_filter "$duration")[bg];[1:v]scale=600:375,setsar=1[orig];[2:v]scale=600:375,setsar=1[remake];[bg][orig]overlay=40:150[withorig];[withorig][remake]overlay=680:150[panels];[panels]drawtext=fontfile=$font:textfile=$original_file:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.88:boxborderw=8:x=58:y=164,drawtext=fontfile=$font:textfile=$remake_file:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.88:boxborderw=8:x=698:y=164" \
        -an -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

make_title opening 5 "臥龍傳 Remake" "經典再現｜原版遊玩錄影 × DOS/V remake"
make_compare natural 8 \
    docs/images/yt-wolong-natural-80s-640x400.png \
    docs/images/wlgame-dosv-natural-remake-skeleton.png \
    "自然策略畫面" \
    "原版參考｜YouTube 80s" \
    "Remake｜DOS/V 自然 HUD" \
    "先看畫面語彙與骨架；代表幀不是同狀態逐像素判定"
make_compare choice 8 \
    docs/images/yt-wolong-natural-160s.png \
    docs/images/wlgame-event3-choice.png \
    "事件與 TALK 選擇" \
    "原版參考｜YouTube 160s" \
    "Remake｜事件 2–5 TALK" \
    "以可重播 fixture 對照反應、數值選取與排版"
make_compare battle 8 \
    docs/images/yt-wolong-natural-240s.png \
    docs/images/wlgame-ai-battle-afterpatch.png \
    "戰術畫面" \
    "原版參考｜YouTube 240s" \
    "Remake｜戰術戰鬥" \
    "展示相同類型的戰場資訊層，不宣稱相同戰局"
make_compare projectile 8 \
    docs/images/yt-wolong-natural-320s.png \
    docs/images/wlgame-ai-battle-attack-afterpatch.png \
    "攻擊與投射物" \
    "原版參考｜YouTube 320s" \
    "Remake｜投射物時序" \
    "以 remake 實際驗收畫面呈現動作方向與結果回寫"
make_compare result 8 \
    docs/images/yt-wolong-natural-400s.png \
    docs/images/wlgame-ai-battle-result.png \
    "戰果與流程回寫" \
    "原版參考｜YouTube 400s" \
    "Remake｜戰果畫面" \
    "原版錄影作視覺參考；remake 由固定 seed 重播"
make_compare postbattle 8 \
    docs/images/yt-wolong-natural-480s.png \
    docs/images/wlgame-ai-postbattle.png \
    "戰後回到自然畫面" \
    "原版參考｜YouTube 480s" \
    "Remake｜戰後自然流程" \
    "把經典的自然／戰術切換放回可執行的跨平台骨架"
make_title closing 7 "經典再現" "原版來源：使用者提供 YouTube 錄影 af6xqcicXoI｜Remake 不含原版資產"

concat="$tmp/concat.txt"
: >"$concat"
for name in opening natural choice battle projectile result postbattle closing; do
    printf "file '%s'\n" "$tmp/$name.mp4" >>"$concat"
done

if command -v python3 >/dev/null 2>&1; then
    python3 "$repo_dir/tools/promo_score.py" "$tmp/score.wav" >/dev/null
elif [ -f "$repo_dir/dist/promo/wolong-remake-trailer-score.wav" ]; then
    # u5cht/video 影像工具容器不保證附帶 Python；沿用同一份本專案原創合成聲。
    cp "$repo_dir/dist/promo/wolong-remake-trailer-score.wav" "$tmp/score.wav"
else
    echo "找不到 promo_score.py 的 Python 執行環境或既有合成聲：$repo_dir/dist/promo/wolong-remake-trailer-score.wav" >&2
    exit 1
fi
ffmpeg -hide_banner -loglevel error -y \
    -f concat -safe 0 -i "$concat" -i "$tmp/score.wav" \
    -map 0:v:0 -map 1:a:0 -t 60 \
    -c:v libx264 -preset medium -crf 19 -pix_fmt yuv420p \
    -c:a aac -b:a 160k -ar 44100 -movflags +faststart "$out"

yt="$repo_dir/docs/images/yt-wolong-natural-80s-640x400.png"
remake="$repo_dir/docs/images/wlgame-dosv-natural-remake-skeleton.png"
convert "$yt" "$remake" -alpha off +append "$review_dir/classic-revival-natural-side-by-side.png"
convert "$yt" "$remake" -alpha off -compose difference -composite "$review_dir/classic-revival-natural-difference.png"

ffprobe -v error \
    -show_entries format=duration:stream=codec_name,width,height,channels,sample_rate \
    -of default=noprint_wrappers=1 "$out"
echo "經典再現推廣片：$out"
