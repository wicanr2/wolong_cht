//go:build ignore

// 對拍夾具：從一份快照起跑，量 remake 每個子刻取幾個亂數、由誰取、
// 世界變成什麼樣——拿來跟原版逐拍比（docs/spec/160、docs/playtest/119）。
//
//	tools/go.sh run tools/rng_pace.go -save <SAVE.DAT> -ticks 5480 \
//	    -rng-state rng.bin -seq-out seq.txt -save-out out.DAT
//
// ⭐ **不改 internal/**：`World.Tick` 的 rng 是外部傳進來的介面，
// 這裡包一層記錄就好。`runtime.Caller(1)` 在包裝層取到的正是規則層的
// 呼叫點，與原版 `-watch 1ECE0` 的「來自 近=」是同一種東西。
//
// ⚠⚠ **夾具要把正式遊戲會做的設定全部做齊。** `state.LoadScenario` 只還原
// 存檔裡有的東西；道路圖、戰術層與政略 AI 都是**執行期注入**的，少掛一項
// 不會報錯也不會警告——`w.roads` 是 nil 時行軍走直線退路、`strategicAI`
// 是 false 時 AI 不編軍團也不做月結政略評估。兩者都讓對拍量到的是夾具，
// 而症狀長得像規則差異（docs/spec/169 §7）。
//
// ⭐ 所以這支程式**開跑前把接了什麼印出來**，缺一項就是一行看得見的字。
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

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
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
	seqOut := flag.String("seq-out", "", "把逐子刻取數序列寫到檔案（供與原版 diff）")
	traceN := flag.Int("trace", 0, "印前 N 拍的據點狀態（小樣本追蹤）")
	mark := flag.String("mark", "", "印出含這個呼叫點的子刻位置（例：strategy.go:630）")
	rngState := flag.String("rng-state", "", "載入原版當下的亂數狀態（258 byte，docs/spec/147 §5）")
	root := flag.String("root", "workplace/orig/dosv", "原版目錄（道路圖與戰術層都從這裡讀）")
	withAI := flag.Bool("ai", true, "開政略 AI（月結的宣戰／遷都評估與 AI 編軍團）")
	withTactical := flag.Bool("tactical", true, "接戰術層——玩家捲進去而且沒委任時才用得到")
	player := flag.Int("player", -1, "覆寫玩家勢力（-1 ＝ 用存檔裡的）")
	events := flag.Bool("events", false, "印每一次事件推送，對應原版的 eventwatch（docs/spec/139）")
	answer := flag.String("answer", "reject",
		"自動回應「等玩家」的視窗：reject（一律拒絕，對應「什麼都不做」的實驗）／"+
			"accept（一律接受）／stop（停下來，之後的比較沒有意義）")
	saveOut := flag.String("save-out", "", "跑完之後把世界寫成一份 SAVE.DAT（拿去跟原版同一拍的記憶體逐 byte 比）")
	corpsWatch := flag.String("corps", "",
		"逐拍印這幾支軍團（逗號分隔）——只在**任何一個追蹤欄位變了**的那一拍印")
	flag.Parse()

	w, err := state.LoadScenario(*save, *slot)
	if err != nil {
		fmt.Fprintln(os.Stderr, "讀不到存檔：", err)
		os.Exit(1)
	}
	if *player >= 0 {
		w.Player = *player
	}

	// ⭐ **規則層不讀檔案**：道路圖與戰術層由呼叫端注入，政略 AI 是執行期開關。
	// 三樣都不是「載入就有」的東西，缺哪一樣都不會報錯（docs/spec/169 §7）。
	wired := []string{}
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
	w.SetRoads(march.New(len(w.Cities), world.MarchEdges(edges, xy)))
	// ⭐ 夾具要說得出自己的狀態：還原不了的軍團會**不動**（docs/spec/172 §4.5），
	// 而不動與「原版也沒動」在取數序列上長得一樣。
	wired = append(wired, fmt.Sprintf("道路圖 %d 條邊（行軍還原 %d 支、還原不了 %d 支）",
		len(edges), w.RestoredMarches(), w.UnresolvedMarches()))

	if *withAI {
		w.EnableStrategicAI()
		wired = append(wired, "政略 AI")
	}
	if *withTactical {
		_, setup, err := battlesetup.Load(battlesetup.Options{
			Dir: *root, World: w, Map: lib.World,
			Warn: func(string) {},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "戰術層接不上：", err)
			os.Exit(1)
		}
		w.SetTactical(setup)
		wired = append(wired, "戰術層")
	}
	if *events {
		w.OnEvent = func(e state.EventTrace) {
			fmt.Printf("  ▶ %d年%d月%d日 %d時 事件%02X 發起%3d 對象%3d 第二%3d 槽%3d\n",
				e.Year, e.Month, e.Day, e.Hour, e.Code, e.Source,
				e.Param&0xFF, e.Param>>8, e.Slot)
		}
		wired = append(wired, "事件軌跡")
	}
	fmt.Printf("已接：%s；玩家勢力 %d\n", strings.Join(wired, "、"), w.Player)

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

	start := 0
	if len(w.Cities) > 0 {
		s := w.TakeSnapshot()
		start = s.CityCursor
	}
	startClock := w.Clock
	markAt := []int{}
	watched := []int{}
	for _, f := range strings.Split(*corpsWatch, ",") {
		if f = strings.TrimSpace(f); f == "" {
			continue
		}
		if n, err := strconv.Atoi(f); err == nil && n >= 0 && n < len(w.Corps) {
			watched = append(watched, n)
		}
	}
	// ⚠ 追蹤欄位要含**意圖與階段**，不只位置：AI 的分歧常常先出現在
	// 「決定去哪」而位置還沒動（docs/spec/170：從決定到踏出去跨三個週期）。
	lastCorps := map[int][6]int{}
	pendingDiplo, pendingFunding, pendingBattle := 0, 0, 0
	answeredDiplo, answeredFunding := 0, 0
	ticksRun := *ticks
	perTickWhere := make([][]string, 0, *ticks)
	perTick := make([]int, 0, *ticks)
	// ⭐ **每一拍處理的是哪一個據點**：原版側從 `sub_14194` 的 SI 拿得到
	// （`SI ÷ 32`），remake 這邊就是 `Tick` 之前的巡迴游標（先處理再前進）。
	// 兩邊都記下來，對齊就不必靠人工填位移（docs/playtest/119 §27）。
	perTickCity := make([]int, 0, *ticks)
	total := map[string]int{}
	prev := 0
	for i := 0; i < *ticks; i++ {
		perTickCity = append(perTickCity, w.TakeSnapshot().CityCursor)
		w.Tick(tr)
		n := tr.seq - prev
		perTick = append(perTick, n)
		perTickWhere = append(perTickWhere, append([]string(nil), tr.where[prev:tr.seq]...))
		for _, s := range tr.where[prev:tr.seq] {
			total[s]++
			if *mark != "" && strings.Contains(s, *mark) {
				markAt = append(markAt, i+1)
			}
		}
		if *traceN > 0 && i < *traceN {
			cs := w.Cities
			// ⭐ **原版是「先處理再前進」**（`sub_13EFD` 讀 `word_10D1E`
			// 指的那一格，處理完才 `si += 0x20`），所以第 i 拍處理的是
			// `start + i`，不是 `start + i + 1`。
			id := (start + i) % len(cs)
			c := cs[id]
			fmt.Printf("  拍 %2d 據點 %3d(si=%04X) 主 %2d 侵攻目標 %3d 佔用 %d 威脅 %3d "+
				"鄰敵 %d 冷卻 %d 受威脅 %v 具體 %v ← 取 %d 個\n",
				i+1, id, id*32, c.Owner, invasion(w, c.Owner), c.Occupancy, c.Threat,
				c.EnemyNeighbours, c.ReliefCooldown, c.Threatened, c.Specific, n)
			for k, nb := range c.Neighbours {
				if nb < 0 || nb >= len(cs) {
					fmt.Printf("      鄰 %d：（無）\n", k)
					continue
				}
				f := -1
				if c.Owner >= 0 && c.Owner < 22 && cs[nb].Owner >= 0 && cs[nb].Owner < 22 {
					f = w.Friendship[c.Owner][cs[nb].Owner].Raw()
				}
				fmt.Printf("      鄰 %d：據點 %3d 主 %2d 佔用 %d 交友度 %3d\n",
					k, nb, cs[nb].Owner, cs[nb].Occupancy, f)
			}
		}
		for _, n := range watched {
			c := w.Corps[n]
			alive := 0
			if c.Alive {
				alive = 1
			}
			cur := [6]int{c.X, c.Y, c.Node, c.TargetNode, c.Ordered, c.Stage*2 + alive}
			if cur != lastCorps[n] {
				fmt.Printf("  拍 %4d 軍團 %d xy(%d,%d) 節點 %d → %d 意圖 %d 階段 %d 朝向 %d 在 %d\n",
					i+1, n, c.X, c.Y, c.Node, c.TargetNode, c.Ordered,
					c.Stage, c.Heading, alive)
				lastCorps[n] = cur
			}
		}
		// ⭐ **夾具要說出自己卡在哪。** 規則層沒有 pending 閘，所以這兩種
		// 「等玩家」的狀態不會讓 `Tick` 停下來，但世界會停在半路——
		// 症狀是「某些拍完全沒跑內政」，而那長得像規則差異。
		// ⭐ **停在半路的拍不能拿來比。** 規則層沒有 pending 閘，所以
		// `Tick` 照樣被呼叫，但 `tickCity` 不跑——症狀是「連續幾千拍
		// 取 0 個數」，而那會把逐拍不一致的總數整個灌爆。
		//
		// 這個實驗是「玩家什麼都不做」，所以預設**一律拒絕**並繼續。
		// 回應的次數會印出來——原版在同一段沒有停，所以那個次數本身
		// 就是一個待查的差異，不能讓它靜靜地消失。
		if c := w.PendingDiplomacy(); c != nil {
			if pendingDiplo == 0 {
				pendingDiplo = i + 1
			}
			switch *answer {
			case "reject":
				w.ResolveDiplomacy(state.DiplomacyReject)
			case "accept":
				w.ResolveDiplomacy(state.DiplomacyAcceptFree)
			}
			answeredDiplo++
		}
		if c := w.PendingFunding(); c != nil {
			if pendingFunding == 0 {
				pendingFunding = i + 1
			}
			switch *answer {
			case "reject":
				w.ResolveFunding(state.FundingReject)
			case "accept":
				w.ResolveFunding(state.FundingFullAmount)
			}
			answeredFunding++
		}
		if w.PendingBattle() != nil && pendingBattle == 0 {
			pendingBattle = i + 1
		}
		// ⚠ 戰術戰鬥沒有「固定回答」可用——它要真的打完。停下來說明，
		// 不要假裝跑得動。
		if pendingBattle > 0 || (*answer == "stop" && pendingDiplo > 0) {
			ticksRun = i + 1
			break
		}
		prev = tr.seq
	}

	if *saveOut != "" {
		// ⭐ **改寫不是重建**：`SaveInto` 從來源 bytes 出發只蓋已解欄位，
		// 所以跟原版記憶體的 diff 只會落在 remake 真的有在寫的欄位上。
		if err := w.SaveInto(*save, *saveOut, *slot); err != nil {
			fmt.Fprintln(os.Stderr, "-save-out：", err)
			os.Exit(1)
		}
		fmt.Printf("世界寫到 %s\n", *saveOut)
	}

	if ticksRun < *ticks {
		fmt.Printf("⚠ 第 %d 拍停下（要 %d 拍）——出現「等玩家」的狀態，"+
			"再往下跑世界會停在半路。要硬跑加 -run-past-pending\n", ticksRun, *ticks)
	}
	fmt.Printf("終點 %d年%d月%d日 %d時 子刻 %d；活軍團 %d\n",
		w.Clock.Year, w.Clock.Month, w.Clock.Day, w.Clock.Hour, w.Clock.Subtick,
		len(w.AliveCorps()))
	if answeredDiplo > 0 {
		fmt.Printf("  外交三選一：第 %d 拍起，自動以 %q 回應 %d 次\n",
			pendingDiplo, *answer, answeredDiplo)
	}
	if answeredFunding > 0 {
		fmt.Printf("  撥款視窗：第 %d 拍起，自動以 %q 回應 %d 次\n",
			pendingFunding, *answer, answeredFunding)
	}
	if pendingBattle > 0 {
		fmt.Printf("  ⚠ 第 %d 拍起有戰術戰鬥等著——**世界從那一拍起就停在半路**\n", pendingBattle)
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
		ticksRun, tr.seq, float64(tr.seq)/float64(ticksRun))
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

	if *seqOut != "" {
		f, err := os.Create(*seqOut)
		if err != nil {
			fmt.Fprintln(os.Stderr, "-seq-out：", err)
			os.Exit(1)
		}
		// ⭐ 連**來源**一起寫：分歧要能自己說明是哪一支多取／少取，
		// 否則每次都得再跑一輪 `-mark` 去猜。
		// ⭐ **檔頭寫起點時鐘**：原版側的 `-watch` 記錄常常比快照早開始，
		// 而 dosgolem 的 `clock` 步驟會在 log 裡留下同一個時刻——
		// 對齊點因此從資料本身讀得出來，不必人工填位移
		// （docs/playtest/119 §27）。
		fmt.Fprintf(f, "# 起點 %d年%d月%d日 %d時 子刻 %d\n",
			startClock.Year, startClock.Month, startClock.Day,
			startClock.Hour, startClock.Subtick)
		// 版面：`拍 個數 據點 來源1,來源2,…`。第三欄是**這一拍處理的據點**，
		// `tools/parity_pace_diff.py` 拿它逐拍檢查兩邊沒有脫節。
		for i, n := range perTick {
			fmt.Fprintf(f, "%d %d %d %s\n", i+1, n, perTickCity[i],
				strings.Join(perTickWhere[i], ","))
		}
		f.Close()
		fmt.Printf("\n逐子刻序列寫到 %s\n", *seqOut)
	}

	if *mark != "" {
		fmt.Printf("\n含 %s 的子刻（%d 次）：%v\n", *mark, len(markAt), markAt)
	}

	fmt.Println("\n呼叫點：")
	for _, r := range rows {
		fmt.Printf("  %6d  %s\n", r.v, r.k)
	}
}

// invasion 讀勢力的侵攻目標（勢力記錄 +0x19），越界回 -1。
func invasion(w *state.World, f int) int {
	if f < 0 || f >= len(w.Factions) {
		return -1
	}
	return w.Factions[f].InvasionTarget
}
