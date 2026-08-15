// wlbgm 把原版的 `*BGM.DAT` 渲染成 WAV。
//
//	tools/go.sh run ./cmd/wlbgm -bgm workplace/orig/dosv/BGM.DAT -song 0 -out /tmp/s0.wav
//	tools/go.sh run ./cmd/wlbgm -bgm workplace/orig/dosv/OPENBGM.DAT -out /tmp/open.wav
//	tools/go.sh run ./cmd/wlbgm -bgm workplace/orig/dosv/BGM.DAT -song 0 -dump
//
// WAV → ogg 用 `tools/bgm2ogg.sh`（docker ffmpeg，純 Go 沒有 vorbis 編碼器）。
//
// ⚠ 輸出是**原版衍生物**：不進版控、不進發行包（`CLAUDE.md` §9）。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/audio"
)

func main() {
	bgmPath := flag.String("bgm", "workplace/orig/dosv/BGM.DAT", "`*BGM.DAT` 路徑")
	comPath := flag.String("ynsound", "", "YNSOUND.COM 路徑（預設取 -bgm 的同目錄）")
	song := flag.Int("song", 0, "第幾曲")
	seconds := flag.Float64("seconds", 60, "渲染幾秒（原版無限循環，一定要給上限）")
	rate := flag.Float64("rate", 44100, "輸出取樣率")
	peak := flag.Float64("peak", 0.7, "正規化的峰值，0 表示不正規化")
	out := flag.String("out", "", "輸出 WAV 路徑")
	dump := flag.Bool("dump", false, "只印曲目資訊，不渲染")
	flag.Parse()

	data, err := os.ReadFile(*bgmPath)
	check(err)
	if *comPath == "" {
		*comPath = filepath.Join(filepath.Dir(*bgmPath), "YNSOUND.COM")
	}
	com, err := os.ReadFile(*comPath)
	check(err)
	tbl, err := audio.LoadTables(com)
	check(err)
	songs, err := audio.Songs(data)
	check(err)

	if *dump {
		fmt.Printf("%s：%d 曲\n", *bgmPath, len(songs))
		for i := range songs {
			s := &songs[i]
			n := 0
			for s.Instrument(n) != nil {
				n++
			}
			fmt.Printf("  第 %2d 曲  %5d B  音色表 @0x%04X\n", i, len(s.Data), s.Instruments)
		}
		return
	}
	if *song < 0 || *song >= len(songs) {
		fatal(fmt.Sprintf("曲號 %d 超出範圍（0–%d）", *song, len(songs)-1))
	}
	if *out == "" {
		fatal("要給 -out")
	}

	chip := audio.NewOPL3(0)
	player := audio.NewPlayer(&songs[*song], tbl, chip)
	pcm := audio.Render(player, *seconds)
	pcm = audio.Resample(pcm, *rate)
	gain := 0.0
	if *peak > 0 {
		gain = audio.Normalize(pcm, *peak)
	}

	f, err := os.Create(*out)
	check(err)
	check(audio.WriteWAV(f, pcm))
	check(f.Close())

	// 分段 RMS 一定要印出來：「有檔案、長度正確、全靜音」是這條路的
	// 典型失敗樣子（`docs/playtest/25` §3），總 RMS 看不出來。
	fmt.Printf("%s：%.1f 秒、%.0f Hz、%d 個 tick、增益 ×%.1f\n",
		*out, float64(len(pcm.L))/pcm.Rate, pcm.Rate, pcm.TickCount, gain)
	fmt.Print("  分十段 RMS：")
	for _, v := range audio.SegmentRMS(pcm, 10) {
		fmt.Printf(" %.3f", v)
	}
	fmt.Println()
}

func check(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "wlbgm: "+msg)
	os.Exit(1)
}
