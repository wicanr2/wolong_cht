//go:build ignore

// 量 remake 每個子刻取幾個亂數、由誰取——拿來跟原版的節拍對齊
// （docs/spec/160、docs/playtest/116）。
//
//	tools/go.sh run tools/rng_pace.go -save <SAVE.DAT> -ticks 324
//
// ⭐ **不改 internal/**：`World.Tick` 的 rng 是外部傳進來的介面，
// 這裡包一層記錄就好。`runtime.Caller(1)` 在包裝層取到的正是規則層的
// 呼叫點，與原版 `-watch 1ECE0` 的「來自 近=」是同一種東西。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// tracing 包住原版的產生器，記下每一次取數的呼叫點。
// **它不改變輸出**：值原封不動來自底下那一支。
type tracing struct {
	inner *rng.Rand
	seq    int
	where  []string
	values []int
}

func (t *tracing) Next() int {
	v := t.inner.Next()
	t.seq++
	at := "?"
	if _, file, line, ok := runtime.Caller(1); ok {
		at = filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	t.where = append(t.where, at)
	t.values = append(t.values, v)
	return v
}

func main() {
	save := flag.String("save", "workplace/dosgolem/root-liubei/SAVE.DAT", "存檔")
	slot := flag.Int("slot", 0, "存檔槽 0–3")
	ticks := flag.Int("ticks", 324, "要記錄幾個子刻")
	skip := flag.Int("skip", 0, "先推幾個子刻不記錄——用來對齊原版的起點")
	seed := flag.Int("seed", 1, "亂數種子（沒給 -rng-state 時用）")
	cursor := flag.Int("city-cursor", -1, "載入後把據點巡迴游標設成這個值（原版存檔 +0x2E ÷ 32）")
	mark := flag.String("mark", "", "印出含這個呼叫點的子刻位置（例：strategy.go:630）")
	rngState := flag.String("rng-state", "", "載入原版當下的亂數狀態（258 byte，docs/spec/147 §5）")
	flag.Parse()

	w, err := state.LoadScenario(*save, *slot)
	if err != nil {
		fmt.Fprintln(os.Stderr, "讀不到存檔：", err)
		os.Exit(1)
	}
	r := rng.NewFixed(*seed)
	if *rngState != "" {
		raw, err := os.ReadFile(*rngState)
		if err != nil {
			fmt.Fprintln(os.Stderr, "-rng-state：", err)
			os.Exit(1)
		}
		got, ok := rng.FromRaw(raw)
		if !ok {
			fmt.Fprintf(os.Stderr, "-rng-state：%s 是 %d byte，預期 %d\n", *rngState, len(raw), rng.RawStateLen)
			os.Exit(1)
		}
		r = got
		fmt.Printf("載入原版亂數狀態：%s\n", *rngState)
	}
	tr := &tracing{inner: r}

	if *cursor >= 0 {
		s := w.TakeSnapshot()
		s.CityCursor = *cursor
		if err := w.Restore(s); err != nil {
			fmt.Fprintln(os.Stderr, "-city-cursor：", err)
			os.Exit(1)
		}
		fmt.Printf("據點游標設成 %d\n", *cursor)
	}
	fmt.Printf("起點 %d年%d月%d日 %d時 子刻 %d\n",
		w.Clock.Year, w.Clock.Month, w.Clock.Day, w.Clock.Hour, w.Clock.Subtick)

	for i := 0; i < *skip; i++ {
		w.Tick(tr)
	}
	tr.seq, tr.where, tr.values = 0, nil, nil
	fmt.Printf("跳過 %d 子刻後：%d年%d月%d日 %d時 子刻 %d\n",
		*skip, w.Clock.Year, w.Clock.Month, w.Clock.Day, w.Clock.Hour, w.Clock.Subtick)

	markAt := []int{}
	perTick := make([]int, 0, *ticks)
	total := map[string]int{}
	prev := 0
	for i := 0; i < *ticks; i++ {
		w.Tick(tr)
		n := tr.seq - prev
		perTick = append(perTick, n)
		for _, s := range tr.where[prev:tr.seq] {
			total[s]++
			if *mark != "" && strings.Contains(s, *mark) {
				markAt = append(markAt, i+1)
			}
		}
		prev = tr.seq
	}

	// 每子刻取幾個：分布比平均值有用——原版是「每子刻 2 或 3 個」。
	dist := map[int]int{}
	for _, n := range perTick {
		dist[n]++
	}
	keys := make([]int, 0, len(dist))
	for k := range dist {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	fmt.Printf("\n%d 個子刻，共取 %d 個亂數（平均 %.2f／子刻）\n",
		*ticks, tr.seq, float64(tr.seq)/float64(*ticks))
	fmt.Println("每子刻取數的分布：")
	for _, k := range keys {
		fmt.Printf("  %d 個 × %d 子刻\n", k, dist[k])
	}

	type kv struct {
		k string
		v int
	}
	rows := make([]kv, 0, len(total))
	for k, v := range total {
		rows = append(rows, kv{k, v})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].v > rows[j].v })
	fmt.Print("\n前 10 個取到的值：")
	for i, v := range tr.values {
		if i >= 10 {
			break
		}
		fmt.Printf(" %d", v)
	}
	fmt.Println()

	fmt.Print("\n序列：")
	for i, n := range perTick {
		if i >= 60 {
			break
		}
		fmt.Printf(" %d", n)
	}
	fmt.Println()

	if *mark != "" {
		fmt.Printf("\n含 %s 的子刻（%d 次）：%v\n", *mark, len(markAt), markAt)
	}

	fmt.Println("\n呼叫點：")
	for _, r := range rows {
		fmt.Printf("  %6d  %s\n", r.v, r.k)
	}
}
