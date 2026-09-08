//go:build ignore

// 於 dosgolem module 的 Docker 環境內建置；只使用公開 API 與正常滑鼠操作。
// 原始素材唯讀，SetScratch 將存檔寫入明確的可丟棄輸出目錄。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/dosgolem/apps/wolong"
	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	root := flag.String("root", "/orig", "唯讀原版素材")
	out := flag.String("out", "/out", "截圖與收據目錄")
	scratch := flag.String("scratch", "/scratch", "存檔暫存層")
	reload := flag.Bool("reload", false, "新機器讀回第二槽")
	flag.Parse()
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	must(os.MkdirAll(*out, 0755))
	must(os.MkdirAll(*scratch, 0755))
	o, err := wolong.Load(filepath.Join(*root, "KI.EXE"), *root)
	must(err)
	defer o.Close()
	o.SetScratch(*scratch)
	shot := func(name string) {
		must(wolong.Shot(o, filepath.Join(*out, name+".png")))
		fmt.Printf("%s：%s，指令 %d\n", name, wolong.Clock(o), o.Steps())
	}
	click := func(x, y int) { fmt.Printf("左鍵 %d,%d\n", x, y); must(o.Click(x, y)); must(o.Run(600000)) }
	must(o.RunUntil(wolong.Booted(), oracle.Budget(40000000)))
	shot("title")
	click(320, 200)
	shot("load-slots")
	y := 152
	if *reload {
		y = 200
	}
	click(300, y)
	shot("loaded")
	if !*reload {
		fmt.Println("畫面座標左鍵 448,16")
		must(wolong.ClickScreen(o, 448, 16))
		must(o.Run(200000))
		shot("system")
		click(376, 160)
		shot("save-slots")
		click(300, 200)
		shot("after-save")
	}
	data, err := json.MarshalIndent(o.Wrote(), "", "  ")
	must(err)
	must(os.WriteFile(filepath.Join(*out, "writes.json"), data, 0644))
}
