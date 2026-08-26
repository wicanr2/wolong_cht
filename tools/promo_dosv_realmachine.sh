#!/usr/bin/env bash
# 以「自己錄的松崗 DOS/V 實機遊玩」與「remake 實機擷取」合成對照推廣片。
#
#   tools/promo_dosv_realmachine.sh [輸出.mp4]
#
# ⭐ 與 tools/promo_dosv_live_comparison.sh 的差別：**原版側是自己跑的**。
# 那一支的原版畫面九成來自使用者提供的 YouTube 錄影；這一支的原版畫面
# 出自 tools/dosv_capture.sh 的受控 DOSBox-X，可以照 timeline 重跑。
#
# ⚠ **只有戰術戰場那一格仍用錄影。** 原版的戰場要「編成軍團然後等 AI 來打」
# （docs/playtest/40 §1），而原版以時鐘播種、每次 RNG 不同——四次專門擷取
# 都沒等到。這一格片上標明來源，其餘全部是實機。
#
# 邊界：容器不開網路、素材唯讀掛載、只寫輸出目錄；成品不含原版執行檔、
# 資料檔或字型，也不隨遊戲包發行。
set -euo pipefail

repo_dir=${WOLONG_PROMO_REPO:-/src}
capture_root=${WOLONG_PROMO_CAPTURE:-/capture}
remake_root=${WOLONG_PROMO_REMAKE:-/remake}
font=${PROMO_FONTFILE:-/usr/share/fonts/opentype/noto/NotoSerifCJK-Bold.ttc}
out=${1:-/out/wolong-remake-dosv-realmachine.mp4}
out_dir=$(dirname "$out")

live=$capture_root/promo-dosv-main
battle=$capture_root/promo-dosv-battle
adlib=${WOLONG_PROMO_ADLIB:-$capture_root/original-audio/original-adlib.wav}
fill=${WOLONG_PROMO_FILL:-$capture_root/original-video/original.mp4}

need() { [ -e "$1" ] || { echo "找不到必要素材：$1" >&2; exit 1; }; }
for f in "$live/o1-newgame.mp4" "$live/o2-strategy.mp4" "$live/o3-corps.mp4" \
         "$live/o4-march.mp4" "$live/o5-battle.mp4" "$adlib" "$fill" "$font" \
         "$remake_root/map/f00000.png" "$remake_root/battle/f00000.png" \
         "$remake_root/siege/f00000.png" "$remake_root/newgame/f00000.png" \
         "$remake_root/form/f00000.png" "$remake_root/march/f00000.png"; do
    need "$f"
done
mkdir -p "$out_dir"

tmp=$(mktemp -d /tmp/wolong-realmachine-XXXXXX)
trap 'rm -rf "$tmp"' EXIT INT TERM

txt() { printf '%s' "$2" > "$1"; }

fade_of() { awk -v d="$1" 'BEGIN { printf "fade=t=in:st=0:d=0.25,fade=t=out:st=%.2f:d=0.25", d - 0.25 }'; }

# 一格全幅畫面的共用鏈：統一 30 fps、最近鄰放大、上下壓字條。
chrome_of() {
    local tf=$1 sf=$2 d=$3
    # ⚠ **畫面要縮進字條之間**，不要讓字條壓在遊戲畫面上。
    # 640×400 直接放到 1280×720 會填滿整個高度，上下兩條半透明字條
    # 就正好蓋住原版自己的標題列與底列——那兩列是遊戲畫面的一部分。
    # 944×590 保持 1.6 的比例，剛好落在 y=65..655 的空檔裡。
    printf '%s' \
"fps=30,scale=944:590:flags=neighbor,setsar=1,\
pad=1280:720:(ow-iw)/2:65:color=0x07142e,\
drawbox=x=0:y=0:w=1280:h=64:color=0x07142e@0.86:t=fill,\
drawbox=x=0:y=660:w=1280:h=60:color=0x07142e@0.86:t=fill,\
drawtext=fontfile=$font:textfile=$tf:fontcolor=white:fontsize=34:x=(w-text_w)/2:y=13,\
drawtext=fontfile=$font:textfile=$sf:fontcolor=0xb8c8e8:fontsize=23:x=(w-text_w)/2:y=678,\
$(fade_of "$d")"
}

card() {
    local name=$1 d=$2 title=$3 sub=$4
    txt "$tmp/$name.t" "$title"; txt "$tmp/$name.s" "$sub"
    echo "  卡：$name（${d}s）"
    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -f lavfi -i "color=c=0x07142e:s=1280x720:r=30:d=$d" \
        -vf "drawtext=fontfile=$font:textfile=$tmp/$name.t:fontcolor=white:fontsize=64:x=(w-text_w)/2:y=272,drawtext=fontfile=$font:textfile=$tmp/$name.s:fontcolor=0xb8c8e8:fontsize=28:x=(w-text_w)/2:y=392,$(fade_of "$d")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

# 原版實機：640×480 的視窗，遊戲畫面是垂直置中的 640×400。
orig_clip() {
    local name=$1 d=$2 src=$3 start=$4 title=$5 sub=$6
    txt "$tmp/$name.t" "$title"; txt "$tmp/$name.s" "$sub"
    echo "  原版實機：$name（${d}s）"
    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -ss "$start" -t "$d" -i "$src" \
        -vf "crop=640:400:0:40,$(chrome_of "$tmp/$name.t" "$tmp/$name.s" "$d")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

# 補充素材（使用者提供的錄影）：4:3 容器先正規化再裁成 640×400。
fill_clip() {
    local name=$1 d=$2 start=$3 title=$4 sub=$5
    txt "$tmp/$name.t" "$title"; txt "$tmp/$name.s" "$sub"
    echo "  補充錄影：$name（${d}s）"
    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -ss "$start" -t "$d" -i "$fill" \
        -vf "scale=640:480:flags=neighbor,crop=640:400:0:40,$(chrome_of "$tmp/$name.t" "$tmp/$name.s" "$d")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

# remake：逐幀 PNG（640×400），照 spec/71 錄的。
remake_clip() {
    local name=$1 d=$2 dir=$3 title=$4 sub=$5
    txt "$tmp/$name.t" "$title"; txt "$tmp/$name.s" "$sub"
    echo "  remake 實機：$name（${d}s）"
    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -framerate 30 -start_number 0 -i "$remake_root/$dir/f%05d.png" \
        -vf "tpad=stop_mode=clone:stop_duration=$d,trim=duration=$d,setpts=PTS-STARTPTS,$(chrome_of "$tmp/$name.t" "$tmp/$name.s" "$d")" \
        -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

# 並排：左原版實機、右 remake 實機。
# ⚠ **來源尺寸不一定是 640×480。** 受控擷取（`tools/dosv_capture.sh`）是
# 640×480 上下各 40 px 黑邊；使用者提供的錄影是 956×720。所以先一律
# `scale=640:480` 再裁 —— 少了那一步，956 寬的來源會被裁掉右邊一大半，
# 而畫面看起來只是「構圖怪」不像出錯。
split_clip() {
    local name=$1 d=$2 src=$3 start=$4 dir=$5 title=$6 lft=$7 rgt=$8 note=$9
    txt "$tmp/$name.t" "$title"; txt "$tmp/$name.l" "$lft"
    txt "$tmp/$name.r" "$rgt";  txt "$tmp/$name.n" "$note"
    echo "  並排：$name（${d}s）"
    ffmpeg -hide_banner -loglevel error -y -threads 2 \
        -f lavfi -i "color=c=0x07142e:s=1280x720:r=30:d=$d" \
        -ss "$start" -t "$d" -i "$src" \
        -framerate 30 -start_number 0 -i "$remake_root/$dir/f%05d.png" \
        -filter_complex "\
[1:v]scale=640:480:flags=neighbor,setsar=1,crop=640:400:0:40,fps=30,scale=580:362:flags=neighbor,setsar=1,pad=580:436:0:37:color=0x07142e[L];\
[2:v]fps=30,tpad=stop_mode=clone:stop_duration=$d,trim=duration=$d,setpts=PTS-STARTPTS,scale=580:362:flags=neighbor,setsar=1,pad=580:436:0:37:color=0x07142e[R];\
[0:v][L]overlay=50:150[a];[a][R]overlay=650:150[b];\
[b]drawbox=x=0:y=0:w=1280:h=64:color=0x07142e@0.86:t=fill,\
drawbox=x=0:y=660:w=1280:h=60:color=0x07142e@0.86:t=fill,\
drawtext=fontfile=$font:textfile=$tmp/$name.t:fontcolor=white:fontsize=34:x=(w-text_w)/2:y=13,\
drawtext=fontfile=$font:textfile=$tmp/$name.l:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.9:boxborderw=8:x=62:y=103,\
drawtext=fontfile=$font:textfile=$tmp/$name.r:fontcolor=white:fontsize=24:box=1:boxcolor=0x061128@0.9:boxborderw=8:x=662:y=103,\
drawtext=fontfile=$font:textfile=$tmp/$name.n:fontcolor=0xb8c8e8:fontsize=22:x=(w-text_w)/2:y=678,\
$(fade_of "$d")[out]" \
        -map "[out]" -an -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p "$tmp/$name.mp4"
}

echo "── 產生片段 ──"
card       s01 4  "臥龍傳 Remake" "松崗 DOS/V 原版實機 × remake：同一個流程並排"
split_clip s02 7  "$live/o1-newgame.mp4" 2 newgame "開新遊戲" \
           "松崗 DOS/V 原版實機" "Remake 實機｜啟動殼層" \
           "原版四層：ＹＥＳ／ＮＯ → 劇本 → 勢力清單 → 君主卡"
split_clip s03 9  "$live/o2-strategy.mp4" 1 map "大地圖對照" \
           "松崗 DOS/V 原版實機" "Remake 實機｜seed 17" \
           "同類畫面對照，非同一局面；remake 側為最高速檔錄製"
split_clip s04 8  "$live/o3-corps.mp4" 3 form "軍團編成" \
           "松崗 DOS/V 原版實機" "Remake 實機｜六個編成位置" \
           "主將／前鋒／左翼／右翼／左備／右備，配兵與士氣同一套規則"
split_clip s05 8  "$live/o5-battle.mp4" 20 map "指令與事件" \
           "原版實機｜事件訊息與時鐘" "Remake 實機｜四個常駐視窗" \
           "兩側都是實際執行畫面，非逐像素判定"
split_clip s06 8  "$live/o4-march.mp4" 2 march "行軍指示" \
           "松崗 DOS/V 原版實機" "Remake 實機｜TALK #21 ＋ 三選一" \
           "選完目的地跳同一則訊息：戰鬥指揮／委任／解體"
# ⚠ **這一段的原版側是使用者提供的錄影，不是本次實機擷取**：
# 原版的戰場要「編成軍團然後等 AI 來打」，四次專門擷取都沒等到
#（docs/promo/dosv-realmachine.md §4）。片上要寫明來源。
split_clip s07 9  "$fill" 317 battle "戰術戰場" \
           "原版（使用者提供的錄影）" "Remake 實機｜野戰" \
           "原版側非本次實機擷取；兩側不是同一場戰鬥"
remake_clip s09 15 siege "Remake 實機｜攻城" "城牆、城門與守軍；同一套戰術規則層"
card       s10 5  "經典再現" "原版實機 × 可執行 remake｜配樂為 DOSBox-X 實錄的原版 AdLib"

concat=$tmp/concat.txt; : > "$concat"
for n in s01 s02 s03 s04 s05 s06 s07 s09 s10; do
    printf "file '%s'\n" "$tmp/$n.mp4" >> "$concat"
done

echo "── 合成 ──"
ffmpeg -hide_banner -loglevel error -y -threads 2 \
    -f concat -safe 0 -i "$concat" -stream_loop -1 -i "$adlib" \
    -map 0:v:0 -map 1:a:0 -t 72 \
    -af "volume=0.9,alimiter=limit=0.92,afade=t=in:st=0:d=1,afade=t=out:st=70:d=2" \
    -c:v libx264 -preset veryfast -crf 19 -pix_fmt yuv420p \
    -c:a aac -b:a 160k -ar 44100 -movflags +faststart "$out"

ffprobe -v error \
    -show_entries format=duration,size:stream=codec_name,codec_type,width,height,r_frame_rate,nb_frames \
    -of default=noprint_wrappers=1 "$out"
echo "實機對照推廣片：$out"
