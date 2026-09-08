#!/usr/bin/env bash
# 容器內渲染；沿用實錄素材並明示歷史版本，不冒充完整同狀態對拍。
set -euo pipefail
[[ -f /.dockerenv ]]
VERSION=${1:?版本}
GUI=${2:?本版正常操作截圖目錄}
[[ "$VERSION" =~ ^v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}$ ]]
OUT=/out
FONT=/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc
WORK=$(mktemp -d /tmp/wolong-promo-render-XXXXXX)
trap 'rm -rf "$WORK"' EXIT
cd "$WORK"
card() {
    local name=$1 duration=$2 title=$3 subtitle=$4
    printf '%s' "$title" > "$name-title.txt"
    printf '%s' "$subtitle" > "$name-subtitle.txt"
    ffmpeg -nostdin -hide_banner -loglevel error -f lavfi -i "color=c=0x122638:s=1280x720:r=30:d=$duration" \
      -vf "drawtext=fontfile=$FONT:textfile=$name-title.txt:fontsize=56:fontcolor=0xead9aa:x=(w-text_w)/2:y=265,drawtext=fontfile=$FONT:textfile=$name-subtitle.txt:fontsize=28:fontcolor=white:x=(w-text_w)/2:y=370" \
      -an -c:v libx264 -threads 2 -preset fast -crf 19 -pix_fmt yuv420p "$name.mp4"
}
clip() {
    local name=$1 duration=$2 source=$3 start=$4 label=$5
    printf '%s' "$label" > "$name-label.txt"
    ffmpeg -nostdin -hide_banner -loglevel error -ss "$start" -i "$source" -t "$duration" \
      -vf "scale=1280:720:flags=neighbor,setsar=1,fps=30,drawbox=x=0:y=0:w=iw:h=54:color=0x122638:t=fill,drawtext=fontfile=$FONT:textfile=$name-label.txt:fontsize=26:fontcolor=white:x=(w-text_w)/2:y=12" \
      -an -c:v libx264 -threads 2 -preset fast -crf 19 -pix_fmt yuv420p "$name.mp4"
}
still() {
    local name=$1 source=$2 label=$3
    printf '%s' "$label" > "$name-label.txt"
    ffmpeg -nostdin -hide_banner -loglevel error -loop 1 -framerate 30 -i "$source" -t 6 \
      -vf "scale=960:600:flags=neighbor,pad=1280:720:160:60:color=0x122638,drawtext=fontfile=$FONT:textfile=$name-label.txt:fontsize=26:fontcolor=white:x=(w-text_w)/2:y=15" \
      -an -c:v libx264 -threads 2 -preset fast -crf 19 -pix_fmt yuv420p "$name.mp4"
}
card 00 4 '臥龍傳・桌面重製' "$VERSION  Linux / Windows / macOS"
clip 01 12 /src/dist/promo/wolong-remake-trailer.mp4 8 '大地圖與政略｜歷史 remake 實錄（2026-08）'
still 02 "$GUI/defaults.png" '新版 AppImage 實拍｜戰後結果頁預設關閉'
still 03 "$GUI/preferences-after-restart.png" '新版 AppImage 實拍｜偏好設定可保留至下次啟動'
clip 04 14 /src/dist/promo/wolong-remake-trailer.mp4 24 '野戰與攻城｜歷史 remake 實錄，非本版逐拍驗收'
clip 05 8 /src/dist/promo/wolong-remake-dosv-realmachine.mp4 44 '原版／remake 歷史展示｜不同局面，不作同狀態對拍'
card 06 6 '保存經典，持續修正' '原版資料請自備｜逐拍戰況仍有差異｜Windows / macOS 待人工驗收'
for name in 00 01 02 03 04 05 06; do printf "file '%s.mp4'\n" "$name"; done > concat.txt
ffmpeg -nostdin -hide_banner -loglevel error -f concat -safe 0 -i concat.txt \
  -i /src/workplace/promo-live/original-audio/original-adlib.wav -t 56 \
  -map 0:v -map 1:a -c:v copy -af 'atrim=0:56,afade=t=in:d=1,afade=t=out:st=54:d=2,loudnorm=I=-18:TP=-1.5:LRA=11' \
  -c:a aac -b:a 192k -ar 44100 -ac 2 -movflags +faststart "$OUT/wolong-remake-promo-$VERSION.mp4"
ffprobe -v error -show_format -show_streams -of json "$OUT/wolong-remake-promo-$VERSION.mp4" > "$OUT/ffprobe.json"
ffmpeg -nostdin -hide_banner -i "$OUT/wolong-remake-promo-$VERSION.mp4" \
  -vf 'blackdetect=d=0.2:pix_th=0.10,freezedetect=n=-50dB:d=2' -af 'volumedetect,silencedetect=n=-50dB:d=2' -f null - 2> "$OUT/media-check.log"
ffmpeg -nostdin -hide_banner -loglevel error -i "$OUT/wolong-remake-promo-$VERSION.mp4" \
  -vf 'fps=1/4,scale=320:180,tile=4x4' -frames:v 1 "$OUT/contact-sheet.png"
ffmpeg -nostdin -hide_banner -loglevel error -i "$OUT/wolong-remake-promo-$VERSION.mp4" \
  -lavfi 'showspectrumpic=s=1280x480:legend=1' -frames:v 1 "$OUT/audio-spectrum.png"
sha256sum "$OUT/wolong-remake-promo-$VERSION.mp4" > "$OUT/video.sha256"
