//go:build ignore

// 受控戰況比較；從 dosgolem 匯出的執行期資料與 RNG 建立同一戰場。
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wicanr2/wolong_cht/internal/assets/battle"
	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/assets/world"
	"github.com/wicanr2/wolong_cht/internal/battlesetup"
	"github.com/wicanr2/wolong_cht/internal/rules/combat"
	"github.com/wicanr2/wolong_cht/internal/rules/march"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func main() {
	original := flag.String("original", "", "dosgolem 輸出目錄")
	out := flag.String("out", "", "本次 remake 收據")
	assets := flag.String("assets", "workplace/orig/dosv", "原版素材")
	field := flag.Bool("field", false, "野戰案例")
	attack := flag.Bool("attack", false, "下攻擊令並等待自然勝負")
	noOrders := flag.Bool("no-orders", false, "進戰術畫面後不下令；不是行軍委任")
	commandTick := flag.Int("command-tick", -1, "原版退卻輸入節拍；預設從 command-events.json 讀取")
	flag.Parse()
	// 原版收據明示野戰時，不容許默認攻城造成假同狀態比較。
	if data, err := os.ReadFile(filepath.Join(*original, "result.json")); err == nil {
		var identity struct {
			Fixture string `json:"fixture"`
		}
		if err := json.Unmarshal(data, &identity); err != nil {
			panic(err)
		}
		if strings.Contains(identity.Fixture, "SAVE-FIELD.DAT") && !*field {
			panic("原版收據為野戰；必須指定 -field，禁止以預設攻城模式重播")
		}
	}
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	var events []struct {
		Tick int `json:"tick"`
	}
	if *noOrders {
		*commandTick = -1
	} else if *commandTick < 0 {
		data, err := os.ReadFile(filepath.Join(*original, "command-events.json"))
		must(err)
		must(json.Unmarshal(data, &events))
		if len(events) == 0 {
			panic("原版沒有退卻命令取樣")
		}
		*commandTick = events[0].Tick
	} else {
		events = append(events, struct {
			Tick int `json:"tick"`
		}{*commandTick})
	}
	must(os.MkdirAll(*out, 0755))
	write := func(name string, data any) {
		raw, err := json.MarshalIndent(data, "", "  ")
		must(err)
		must(os.WriteFile(filepath.Join(*out, name+".json"), raw, 0644))
	}
	w, err := state.LoadScenario(filepath.Join(*original, "entry-SAVE.DAT"), 0)
	must(err)
	raw, err := os.ReadFile(filepath.Join(*original, "spawn-rng.bin"))
	must(err)
	r, ok := rng.FromRaw(raw)
	if !ok {
		panic("亂數長度錯誤")
	}
	lib, err := library.Load(*assets)
	must(err)
	xy := make([][2]int, len(w.Cities))
	for i, c := range w.Cities {
		xy[i] = [2]int{c.X, c.Y}
	}
	edges, err := world.RoadEdges(lib.World, xy)
	must(err)
	w.SetRoads(march.New(len(w.Cities), world.MarchEdges(edges, xy)))
	provider, setup, err := battlesetup.Load(battlesetup.Options{Dir: *assets, World: w, Map: lib.World, Warn: func(s string) { panic(s) }})
	must(err)
	w.SetTactical(setup)
	// 原版 OpenSiege 的 DI 明示城 82，攻方本身仍在野外；只在初始化期間提供同一個城參數。
	before := w.Corps[35].Node
	mode := combat.Siege
	if *field {
		mode = combat.Field
	} else {
		w.Corps[35].Node = 82
	}
	must(w.StageBattle(35, 39, mode, r))
	w.Corps[35].Node = before
	b := w.PendingBattle().Battle
	// 正式 newBattleView 同樣以共享 RNG 初始化實際戰場旗幟，不可任意跳過固定次數。
	fieldNumber := provider.FieldNumber(w.PendingBattle().Node, !*field)
	tiles := provider.Library().Tiles(fieldNumber)
	if provider.Rotate() {
		tiles = battle.Rotate180(tiles)
	}
	write("banners", provider.Library().BannersFor(fieldNumber, tiles, r.Next))
	duel := tactical.DuelInput{FieldNumber: fieldNumber}
	for side, corps := range [2]int{35, 39} {
		g := w.Generals[w.Leader(corps)]
		duel.Martial[side], duel.CommandStat[side] = g.Martial, g.Command
	}
	b.SetDuelInput(duel)
	snapshot := func(name string) {
		var units []map[string]any
		for role := 0; role < 2; role++ {
			side := b.PlayerSide
			if role == 1 {
				side = 1 - side
			}
			for k, u := range b.Sides[side].Soldiers {
				units = append(units, map[string]any{"Side": role, "Squad": k / 8, "Slot": k % 8, "X": u.X, "Y": u.Y, "Stamina": u.HP, "Order": u.Cmd, "NewOrder": u.Next, "Alive": u.Alive, "Power": u.Power, "state": u})
			}
		}
		c, s := r.State()
		write(name, map[string]any{"tick": b.Frame, "units": units, "sides": b.Sides, "rng_prefix": hex.EncodeToString([]byte{c, s}), "opening_active": b.OpeningActive(), "done": b.Done, "winner": b.Winner, "corps": map[int]any{35: w.Corps[35], 39: w.Corps[39]}})
	}
	snapshot("initialized")
	command := tactical.Retreat
	if *attack {
		command = tactical.Attack
	}
	openingRecorded := false
	var inputReceipts []map[string]any
	for !b.Done && b.Frame < 12000 {
		if !openingRecorded && !b.OpeningActive() {
			snapshot("opening-ready")
			openingRecorded = true
		}
		for i, event := range events {
			if event.Tick != b.Frame {
				continue
			}
			accepted := !b.OpeningActive()
			if accepted {
				b.OrderSelected(b.PlayerSide, command)
			}
			inputReceipts = append(inputReceipts, map[string]any{"tick": b.Frame, "accepted": accepted, "opening_active": b.OpeningActive()})
			if i == 0 {
				snapshot("command-state")
			}
		}
		if b.Frame < 1000 || b.Frame%50 == 0 {
			snapshot(fmt.Sprintf("tick-%04d", b.Frame))
		}
		if b.Frame == 50 {
			snapshot("tick50-sample")
		}
		b.Step()
	}
	write("command", map[string]any{"requested_tick": *commandTick, "actual_tick": *commandTick, "field_number": fieldNumber, "events": inputReceipts, "timing": "與原版同拍送入；開場阻擋照正式 UI 丟棄，不延後重送"})
	if !b.Done {
		panic("戰鬥超過 12000 拍")
	}
	snapshot("before-settlement")
	event := w.ResolvePending(r)
	snapshot("world-return")
	write("result", map[string]any{"completed": true, "event": event, "battle": b.Result(), "log": b.Log, "classification": "受控入口；共用規則與結算，非正常 GUI"})
}
