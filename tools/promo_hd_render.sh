#!/usr/bin/env bash
# 在 ffmpeg 容器中，把本版實錄與 DOSBox 原版片頭錄音剪成一支本機推廣片。
set -euo pipefail
[[ -f /.dockerenv ]]
version=${1:?版本}
[[ "$version" =~ ^v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}$ ]]
capture=/out/capture
promo=${WOLONG_PROMO_OUTPUT:-/delivery/promo}
titlefont=/usr/share/fonts/opentype/noto/NotoSerifCJK-Bold.ttc
font=/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc
music="/repo/workplace/dosbox/audio/d7open  _000.wav"
film="$promo/wolong-remake-promo-$version.mp4"
temp=$(mktemp -d)
trap 'rm -rf "$temp"' EXIT
cd "$temp"
card() {
    printf '%s' "$3" > "$1-title"
    printf '%s' "$4" > "$1-subtitle"
    ffmpeg -nostdin -hide_banner -loglevel error -f lavfi -i "color=c=0x172832:s=1280x720:r=30:d=$2" \
        -vf "drawbox=x=480:y=342:w=320:h=2:color=0xc59457:t=fill,drawtext=fontfile=$titlefont:textfile=$1-title:fontsize=64:fontcolor=0xf5d7a2:x=(w-text_w)/2:y=238,drawtext=fontfile=$font:textfile=$1-subtitle:fontsize=28:fontcolor=white:x=(w-text_w)/2:y=382" \
        -an -c:v libx264 -threads 2 -preset fast -crf 19 -pix_fmt yuv420p "$1.mp4"
}
clip() {
    printf '%s' "$3" > "$1-label"
    ffmpeg -nostdin -hide_banner -loglevel error -ss "$5" -i "$capture/$2.mp4" -t "$4" \
      -vf "scale=1184:664:force_original_aspect_ratio=decrease:force_divisible_by=2:flags=lanczos,pad=1280:720:(ow-iw)/2:46:color=0x172832,setsar=1,fps=30,drawtext=fontfile=$font:textfile=$1-label:fontsize=23:fontcolor=0xf4dfaf:x=(w-text_w)/2:y=7" \
      -an -c:v libx264 -threads 2 -preset fast -crf 19 -pix_fmt yuv420p "$1.mp4"
}
opening=/out/capture/opening
opening_filter='format=rgb32,hqx=n=4,scale=1152:720:flags=lanczos,pad=1280:720:64:0:color=0x101820,setsar=1,fps=30'
for scene in 2 3 4; do
    ffmpeg -nostdin -hide_banner -loglevel error -filter_threads 1 -framerate 10 -i "$opening/scene$scene-%02d.png" \
        -vf "$opening_filter" -t 3 -an -c:v libx264 -threads 2 -preset fast -crf 18 -pix_fmt yuv420p "intro$scene.mp4"
done
ffmpeg -nostdin -hide_banner -loglevel error -filter_threads 1 -loop 1 -framerate 30 -i "$opening/scene6.png" \
    -vf "${opening_filter/hqx=n=4/hqx=n=2}" -t 3 -an -c:v libx264 -threads 2 -preset fast -crf 18 -pix_fmt yuv420p intro6.mp4
clip 01 map '大地圖・原版與 HD 即時切換' 9 1
clip 02 siege '椎陽攻城・雲梯與護城河・戰場視野縮放' 15 1
clip 03 field '戰術指揮・陣形與兵種動作' 8 1
card 04 3 '三國制霸之計' "Linux / Windows / macOS  $version"
for part in intro2 intro3 intro4 intro6 01 02 03 04; do printf "file '%s.mp4'\n" "$part"; done > cuts.txt
ffmpeg -y -nostdin -hide_banner -loglevel error -f concat -safe 0 -i cuts.txt -i "$music" -i "$music" -i "$music" -i "$music" \
    -filter_complex '[1:a]atrim=start=0.65:end=14.65,asetpts=PTS-STARTPTS[a];[2:a]atrim=start=0.65:end=14.65,asetpts=PTS-STARTPTS[b];[3:a]atrim=start=0.65:end=14.65,asetpts=PTS-STARTPTS[c];[4:a]atrim=start=0.65:end=14.65,asetpts=PTS-STARTPTS[d];[a][b]acrossfade=d=0.55[o];[o][c]acrossfade=d=0.55[p];[p][d]acrossfade=d=0.55,asetpts=N/SR/TB,atrim=duration=47,volume=19dB,alimiter=limit=0.85:level=false:latency=true,afade=t=in:d=0.6,afade=t=out:st=45:d=2[score]' \
    -map 0:v -map '[score]' -c:v copy -c:a aac -b:a 192k -ar 44100 -movflags +faststart "$film"
ffprobe -v error -show_format -show_streams -of json "$film" > "$promo/ffprobe.json"
ffmpeg -nostdin -hide_banner -i "$film" \
    -vf 'blackdetect=d=0.2:pix_th=0.10,freezedetect=n=-70dB:d=2' \
    -af 'volumedetect,silencedetect=n=-50dB:d=1.5' -f null - 2> "$promo/media-check.log"
ffmpeg -y -nostdin -hide_banner -loglevel error -i "$film" \
    -vf 'fps=1/4,scale=320:180,tile=4x3' -frames:v 1 "$promo/contact-sheet.png"
ffmpeg -y -nostdin -hide_banner -loglevel error -i "$film" -t 12 \
    -vf 'fps=1/3:start_time=0,scale=480:270,tile=2x2' -frames:v 1 "$promo/opening-contact-sheet.png"
cp "$capture/"*-capture.json "$promo/"
sha256sum "$music" "$capture/map.mp4" "$capture/siege.mp4" "$capture/field.mp4" \
    /repo/workplace/orig/dosv/OPENPAL.BRG /repo/workplace/orig/dosv/OPEN_S{2,3,4,6}.DAT > "$promo/source-SHA256SUMS.txt"
