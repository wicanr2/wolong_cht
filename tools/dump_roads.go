//go:build ignore

// 把 remake 算出來的道路邊與逐格路徑輸出成 JSON，拿去跟原版執行期建好的
// 那張表逐條比（`tools/orig_roadtable.py`、docs/re/08 §7）。
//
//	tools/go.sh run tools/dump_roads.go -out workplace/parity/pace/remake-roads.json
//
// ⚠ **要比的是中段。** `RoadEdge.Path` 頭尾各有一段「城中心 → 節點格 →
// 城門格」的直線（`StubA`／`StubB`），原版的路徑點表只有城門格到城門格。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/state"
)

type edge struct {
	A, B   int      `json:"a"`
	B2     int      `json:"b"`
	StubA  int      `json:"stub_a"`
	StubB  int      `json:"stub_b"`
	Points [][2]int `json:"points"`
}

func main() {
	root := flag.String("root", "workplace/orig/dosv", "原版目錄")
	save := flag.String("save", "workplace/orig/dosv/SINARIO.DAT", "劇本或存檔（拿據點座標）")
	slot := flag.Int("slot", 0, "槽")
	out := flag.String("out", "", "輸出 JSON")
	flag.Parse()

	w, err := state.LoadScenario(*save, *slot)
	if err != nil {
		fmt.Fprintln(os.Stderr, "讀不到存檔：", err)
		os.Exit(1)
	}
	lib, err := library.Load(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "讀不到素材：", err)
		os.Exit(1)
	}
	xy := make([][2]int, len(w.Cities))
	for i := range w.Cities {
		xy[i] = [2]int{w.Cities[i].X, w.Cities[i].Y}
	}
	edges, err := world.RoadEdges(lib.World, xy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "RoadEdges：", err)
		os.Exit(1)
	}
	rows := make([]map[string]any, 0, len(edges))
	for _, e := range edges {
		pts := make([][2]int, 0, len(e.Path))
		for _, p := range e.Path {
			pts = append(pts, p)
		}
		rows = append(rows, map[string]any{
			"a": e.A, "b": e.B, "stub_a": e.StubA, "stub_b": e.StubB,
			"steps": e.Steps, "seq": e.Seq, "points": pts,
		})
	}
	fmt.Printf("remake 邊 %d 條，路徑點 %d 個\n", len(rows), func() int {
		n := 0
		for _, e := range edges {
			n += len(e.Path)
		}
		return n
	}())
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		if err := json.NewEncoder(f).Encode(rows); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("寫到 %s\n", *out)
	}
}
