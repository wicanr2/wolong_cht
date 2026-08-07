// wlshot 把原版素材解成 PNG。**不 import Ebiten** —— 所以在無頭環境跑得起來。
//
//	tools/go.sh run ./cmd/wlshot -orig workplace/orig/dosv -asset 0 -page 42 -out a.png
//	tools/go.sh run ./cmd/wlshot -orig workplace/orig/dosv -list
//	tools/go.sh run ./cmd/wlshot -orig workplace/orig/dosv -asset 0 -sheet 15 -out all.png
package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
)

func main() {
	dir := flag.String("orig", "workplace/orig/dosv", "原版素材目錄（請自備）")
	asset := flag.Int("asset", 0, "素材編號，用 -list 查")
	page := flag.Int("page", 0, "第幾張")
	bank := flag.Int("season", 0, "調色盤組 0-3 ＝ 春夏秋冬")
	sheet := flag.Int("sheet", 0, "每列幾張；>0 就輸出全部張數的總覽")
	out := flag.String("out", "", "輸出 PNG 路徑")
	list := flag.Bool("list", false, "列出素材種類就結束")
	flag.Parse()

	lib, err := library.Load(*dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range lib.Warns {
		log.Printf("⚠ %s", w)
	}
	if *list {
		for i, e := range lib.Entries {
			fmt.Printf("%d  %-22s %4d 張  %dx%d\n",
				i, e.Label, e.Count, e.Spec.Width, e.Spec.Height)
		}
		return
	}
	if *out == "" {
		log.Fatal("要給 -out")
	}

	var img image.Image
	if *sheet > 0 {
		img, err = contactSheet(lib, *asset, *bank, *sheet)
	} else {
		img, err = lib.Render(*asset, *page, *bank)
	}
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
	b := img.Bounds()
	fmt.Printf("%s: %dx%d\n", *out, b.Dx(), b.Dy())
}

// contactSheet 把一種素材的全部張數排成總覽，用於肉眼驗收。
func contactSheet(lib *library.Library, asset, bank, perRow int) (image.Image, error) {
	e := lib.Entries[asset]
	if e.Count == 0 {
		return nil, fmt.Errorf("%s 沒有可畫的張數", e.Label)
	}
	if perRow > e.Count {
		perRow = e.Count
	}
	rows := (e.Count + perRow - 1) / perRow
	dst := image.NewRGBA(image.Rect(0, 0, perRow*e.Spec.Width, rows*e.Spec.Height))
	for i := 0; i < e.Count; i++ {
		src, err := lib.Render(asset, i, bank)
		if err != nil {
			return nil, err
		}
		at := image.Pt((i%perRow)*e.Spec.Width, (i/perRow)*e.Spec.Height)
		draw.Draw(dst, src.Bounds().Add(at), src, image.Point{}, draw.Src)
	}
	return dst, nil
}
