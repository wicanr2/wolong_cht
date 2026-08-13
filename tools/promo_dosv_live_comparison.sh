#!/usr/bin/env bash
# 以「松崗 DOS/V 原版實機錄影」與「remake 實機擷取」製作可重播的推廣片。
#
# 此腳本必須在具有 ffmpeg、Python 與中文字型的 Docker 容器內執行。原版來源
# 只讀掛載於 $WOLONG_PROMO_CAPTURE/original-video/original.mp4。原版錄影先固定
# 裁成 DOS/V 640×400 framebuffer，再和 remake 共用相同縮放鏈；配樂只接受
# DOSBox-X 實錄的松崗 DOS/V AdLib WAV，不再呼叫近似合成器。
#
# 預設輸出：/out/wolong-remake-dosv-live-comparison.mp4
set -euo pipefail

repo_dir=${WOLONG_PROMO_REPO:-/src}
capture_root=${WOLONG_PROMO_CAPTURE:-/capture}
font=${PROMO_FONTFILE:-/fonts/truetype/moe/MoeStandardSong.ttf}
out=${1:-/out/wolong-remake-dosv-live-comparison.mp4}
out_dir=$(dirname "$out")
original_video="$capture_root/original-video/original.mp4"
original_adlib=${WOLONG_PROMO_ADLIB:-$capture_root/original-audio/original-adlib.wav}
remake_tactical_dir=${WOLONG_PROMO_TACTICAL:-$capture_root/remake-tactical}

require_file() {
    if [ ! -f "$1" ]; then
        echo "找不到必要檔案：$1" >&2
        exit 1
    fi
}

require_dir() {
    if [ ! -d "$1" ]; then
        echo "找不到必要目錄：$1" >&2
        exit 1
    fi
}

require_file "$original_video"
require_file "$original_adlib"
require_file "$capture_root/story-timed/opening-120.png"
require_dir "$capture_root/remake-strategy"
require_dir "$capture_root/remake-normal"
require_dir "$remake_tactical_dir"
require_file "$font"

mkdir -p "$out_dir"
if [ ! -w "$out_dir" ]; then
    echo "輸出目錄不可寫：$out_dir" >&2
    exit 1
fi

tmp=$(mktemp -d /tmp/wolong-dosv-live-promo-XXXXXX)
cleanup() {
    rm -rf "$tmp"
}
trap cleanup EXIT INT TERM

write_text() {
    printf '%s' "$2" >"$1"
}

fade_filter() {
    awk -v duration="$1" 'BEGIN {
        printf "fade=t=in:st=0:d=0.22,fade=t=out:st=%.2f:d=0.22", duration - 0.22
    }'
}

full_filter() {
    local title_file=$1
    local subtitle_file=$2
    local duration=$3
    printf '%s' \
        "fps=30,scale=1280:720:force_original_aspect_ratio=decrease:flags=neighbor,pad=1280:720:(ow-iw)/2:(oh-ih)/2:color=0x07142e,drawbox=x=0:y=0:w=1280:h=64:color=0x07142e@0.84:t=fill,drawbox=x=0:y=660:w=1280:h=60:color=0x07142e@0.84:t=fill,drawtext=fontfile=$font:textfile=$title_file:fontcolor=white:fontsize=34:x=(w-text_w)/2:y=13,drawtext=fontfile=$font:textfile=$subtitle_file:fontcolor=0xb8c8e8:fontsize=23:x=(w-text_w)/2:y=678,$(fade_filter "$duration")"
}

# 使用者提供的原版錄影是 4:3 容器；下載格式目前為 956×720，但來源格式可能
# 改變。先以最近鄰正規化成 640×480，再裁出垂直置中的 DOS/V 640×400
# framebuffer，最後走和 remake 完全相同的縮放鏈。
original_full_filter() {
    local title_file=$1
    local subtitle_file=$2
    local duration=$3
    printf '%s' "scale=640:480:flags=neighbor,crop=640:400:0:40,$(full_filter "$title_file" "$subtitle_file" "$duration")"
}

make_title() {
    local name=$1
    local duration=$2
    local title=$3
    local subtitle=$4
    local title_file="$tmp/$name-title.txt"
    local subtitle_file="$tmp/$name-subtitle.txt"
    write_text "$title_file" "$title"
    write_text "$subtitle_file" "$subtitle"
    printf '產生片段：%s\n' "$name"

    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -f lavfi -i "color=c=0x07142e:s=1280x720:r=30:d=$duration" \
        -vf "drawtext=fontfile=$font:textfile=$title_file:fontcolor=white:fontsize=68:x=(w-text_w)/2:y=270,drawtext=fontfile=$font:textfile=$subtitle_file:fontcolor=0xb8c8e8:fontsize=30:x=(w-text_w)/2:y=390,$(fade_filter "$duration")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

make_full_image() {
    local name=$1
    local duration=$2
    local image=$3
    local title=$4
    local subtitle=$5
    local title_file="$tmp/$name-title.txt"
    local subtitle_file="$tmp/$name-subtitle.txt"
    write_text "$title_file" "$title"
    write_text "$subtitle_file" "$subtitle"
    printf '產生片段：%s\n' "$name"

    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -loop 1 -framerate 30 -t "$duration" -i "$image" \
        -vf "$(original_full_filter "$title_file" "$subtitle_file" "$duration")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

make_full_original() {
    local name=$1
    local duration=$2
    local start=$3
    local title=$4
    local subtitle=$5
    local title_file="$tmp/$name-title.txt"
    local subtitle_file="$tmp/$name-subtitle.txt"
    write_text "$title_file" "$title"
    write_text "$subtitle_file" "$subtitle"
    printf '產生片段：%s\n' "$name"

    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -ss "$start" -t "$duration" -i "$original_video" \
        -vf "$(full_filter "$title_file" "$subtitle_file" "$duration")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

make_full_sequence() {
    local name=$1
    local duration=$2
    local frames_dir=$3
    local prefix=$4
    local title=$5
    local subtitle=$6
    local title_file="$tmp/$name-title.txt"
    local subtitle_file="$tmp/$name-subtitle.txt"
    write_text "$title_file" "$title"
    write_text "$subtitle_file" "$subtitle"
    printf '產生片段：%s\n' "$name"

    require_file "$frames_dir/$prefix-000001.png"
    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -framerate 5 -start_number 1 -i "$frames_dir/$prefix-%06d.png" \
        -vf "fps=30,tpad=stop_mode=clone:stop_duration=$duration,trim=duration=$duration,setpts=PTS-STARTPTS,$(full_filter "$title_file" "$subtitle_file" "$duration")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

make_split() {
    local name=$1
    local duration=$2
    local original_start=$3
    local remake_dir=$4
    local remake_prefix=$5
    local title=$6
    local original_label=$7
    local remake_label=$8
    local note=$9
    local title_file="$tmp/$name-title.txt"
    local original_file="$tmp/$name-original.txt"
    local remake_file="$tmp/$name-remake.txt"
    local note_file="$tmp/$name-note.txt"
    write_text "$title_file" "$title"
    write_text "$original_file" "$original_label"
    write_text "$remake_file" "$remake_label"
    write_text "$note_file" "$note"
    printf '產生片段：%s\n' "$name"

    require_file "$remake_dir/$remake_prefix-000001.png"
    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -f lavfi -i "color=c=0x07142e:s=1280x720:r=30:d=$duration" \
        -ss "$original_start" -t "$duration" -i "$original_video" \
        -framerate 5 -start_number 1 -i "$remake_dir/$remake_prefix-%06d.png" \
        -filter_complex "[1:v]scale=640:480:flags=neighbor,crop=640:400:0:40,fps=30,scale=580:362:flags=neighbor,setsar=1,pad=580:436:0:37:color=0x07142e[original];[2:v]fps=30,tpad=stop_mode=clone:stop_duration=$duration,trim=duration=$duration,setpts=PTS-STARTPTS,scale=580:362:flags=neighbor,setsar=1,pad=580:436:0:37:color=0x07142e[remake];[0:v][original]overlay=50:150[first];[first][remake]overlay=650:150[panels];[panels]drawbox=x=0:y=0:w=1280:h=64:color=0x07142e@0.84:t=fill,drawbox=x=0:y=660:w=1280:h=60:color=0x07142e@0.84:t=fill,drawtext=fontfile=$font:textfile=$title_file:fontcolor=white:fontsize=34:x=(w-text_w)/2:y=13,drawtext=fontfile=$font:textfile=$original_file:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.88:boxborderw=8:x=62:y=103,drawtext=fontfile=$font:textfile=$remake_file:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.88:boxborderw=8:x=662:y=103,drawtext=fontfile=$font:textfile=$note_file:fontcolor=0xb8c8e8:fontsize=22:x=(w-text_w)/2:y=678,$(fade_filter "$duration")[out]" \
        -map "[out]" -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

# 總長 60 秒。所有原版影片段落明確移除音訊；戰術 remake 段落是可視化 fixture，
# 不被當成「正常自然路徑」的測試證據。
make_title opening 4 "臥龍傳 Remake" "經典再現｜松崗 DOS/V 原版實機 × remake 實機"
make_full_image original_start 5 "$capture_root/story-timed/opening-120.png" \
    "從原版實機開始" "松崗 DOS/V｜正常啟動後的新遊戲畫面"
make_split strategy_compare 7 76 "$capture_root/remake-strategy" remake-idle \
    "自然策略畫面" "松崗 DOS/V 原版實機" "Remake 實機｜固定 seed 17" \
    "同類畫面對照；非同狀態逐像素判定"
make_full_original original_strategy 6 154 \
    "原版實機｜策略地圖與時鐘" "使用者松崗 DOS 原版錄影；沿用其原版 AdLib 音軌"
make_full_sequence remake_formation 6 "$capture_root/remake-normal" remake-formed \
    "Remake 實機｜編成與行軍" "正常鍵盤路徑：編成 → 目的地 → 行軍"
make_split command_compare 5 236 "$capture_root/remake-normal" remake-destination \
    "指令與事件" "原版實機｜事件訊息" "Remake 實機｜目的地選擇" \
    "展示操作與資訊層；非同一局面"
make_full_original original_tactical 7 317 \
    "原版實機｜戰術戰場" "640×400 原始畫布；與 remake 使用相同縮放鏈"
make_full_sequence remake_tactical 5 "$remake_tactical_dir" tactical-idle \
    "Remake 實機｜戰術固定骨架" "攻城 fixture；不同戰況，只作 layout-only 比較"
make_full_original original_battle 5 398 \
    "原版實機｜戰鬥與結果" "畫面片段取自使用者提供錄影"
make_full_sequence remake_march 5 "$capture_root/remake-normal" remake-march \
    "Remake 實機｜行軍、時鐘與自然畫面" "正常鍵盤路徑錄製"
make_title closing 5 "經典再現" "原版實機錄影 × 可執行 remake｜配樂：使用者錄影中的原版 AdLib"

concat="$tmp/concat.txt"
: > "$concat"
for name in opening original_start strategy_compare original_strategy remake_formation command_compare original_tactical remake_tactical original_battle remake_march closing; do
    printf "file '%s'\n" "$tmp/$name.mp4" >> "$concat"
done

ffmpeg -hide_banner -loglevel error -y -threads 2 \
    -f concat -safe 0 -i "$concat" -stream_loop -1 -i "$original_adlib" \
    -map 0:v:0 -map 1:a:0 -t 60 \
    -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p \
    -c:a aac -b:a 160k -ar 44100 -movflags +faststart "$out"

ffprobe -v error \
    -show_entries format=duration,size:stream=codec_name,codec_type,width,height,r_frame_rate,channels,sample_rate \
    -of default=noprint_wrappers=1 "$out"
echo "實機對照推廣片：$out"
