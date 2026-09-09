//go:build ignore

// 在 dosgolem module 內建置的受控戰況取樣器；素材唯讀，不改原版記憶體。
// OpenSiege 是既有的遭遇入口 fixture，不能標成正常玩家觸發。
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/dosgolem/apps/wolong"
	"github.com/wicanr2/dosgolem/oracle"
)

func main() {
	root := flag.String("root", "/orig", "唯讀原版目錄")
	out := flag.String("out", "/out", "收據目錄")
	field := flag.Bool("field", false, "從受控野戰存檔等待正常遭遇")
	attack := flag.Bool("attack", false, "下攻擊令並等待自然勝負，不主動退卻")
	noOrders := flag.Bool("no-orders", false, "進戰術畫面後不下令；不是行軍委任")
	traceStart := flag.Bool("trace-start", false, "只取開場呼叫順序，不宣稱戰鬥完成")
	traceTicks := flag.Int("trace-ticks", 3, "開場追蹤拍數（1–20，配合 trace-start）")
	traceAll := flag.Bool("trace-all-units", false, "配合 trace-start 取全部兵及換位呼叫")
	flag.Parse()
	if *traceTicks < 1 || *traceTicks > 20 {
		panic("trace-ticks 必須介於 1 與 20")
	}
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	must(os.MkdirAll(*out, 0755))
	write := func(name string, data []byte) { must(os.WriteFile(filepath.Join(*out, name), data, 0644)) }
	jsonFile := func(name string, data any) {
		raw, err := json.MarshalIndent(data, "", "  ")
		must(err)
		write(name+".json", raw)
	}
	o, err := wolong.Load(filepath.Join(*root, "KI.EXE"), *root)
	must(err)
	defer o.Close()
	logicalTick := 0
	runto := func(lin uint32, budget uint64) { must(o.RunUntil(oracle.At(o.IDA(lin)), oracle.Budget(budget))) }
	snapshot := func(name string, units bool) {
		data := map[string]any{"instructions": o.Steps(), "clock": wolong.Clock(o), "tick": logicalTick, "raw_tick": wolong.Tick(o), "corps": wolong.CorpsTable(o), "rng": hex.EncodeToString(o.Bytes(o.IDA(0x1ECFC), 258)), "address_space": "IDA DOS/V linear", "battle_setup_10D2E": hex.EncodeToString(o.Bytes(o.IDA(0x10D2E), 8))}
		if units {
			data["units"] = wolong.Units(o, false)
			seg := o.Word(o.IDA(0x1D30E))
			data["unit_segment"] = seg
			data["unit_bytes"] = hex.EncodeToString(o.Bytes(oracle.Far(seg, 0), 0xC00))
			if *traceStart {
				data["path_segment"] = seg
				data["path_offset"] = 0x1800
				data["path_bytes"] = hex.EncodeToString(o.Bytes(oracle.Far(seg, 0x1800), 96*128))
			}
			seg = o.Word(o.IDA(0x1D30A))
			data["summary_segment"] = seg
			data["summary_bytes"] = hex.EncodeToString(o.Bytes(oracle.Far(seg, 0), 64))
			data["control_1D310"] = hex.EncodeToString(o.Bytes(o.IDA(0x1D310), 64))
		}
		jsonFile(name, data)
		must(wolong.Shot(o, filepath.Join(*out, name+".png")))
		fmt.Println(name, o.Steps(), wolong.Tick(o))
	}
	must(o.RunUntil(wolong.Booted(), oracle.Budget(40_000_000)))
	must(o.Click(320, 200))
	must(o.Click(300, 152))
	if !*field {
		must(wolong.OpenSiege(o, 35, 82))
	}
	runto(0x11B5A, 400_000_000)
	snapshot("entry", false)
	// formats/08 §0：由執行期三個已證實區塊匯出給 remake 載入；不是玩家存檔收據。
	save, err := os.ReadFile(filepath.Join(*root, "SAVE.DAT"))
	must(err)
	copy(save[:59], o.Bytes(o.IDA(0x10CF0), 59))
	copy(save[0x80:0x52C0], o.Bytes(oracle.Far(o.Word(o.IDA(0x10D52)), 0), 0x5240))
	copy(save[0x52C0:0x56C0], o.Bytes(oracle.Far(o.Word(o.IDA(0x10D56)), 0), 0x400))
	write("entry-SAVE.DAT", save)
	var calls []map[string]any
	o.OnCall(o.IDA(0x1ECE0), func(o *oracle.Oracle) {
		if len(calls) < 2048 {
			calls = append(calls, map[string]any{"tick": wolong.Tick(o), "caller_ida": fmt.Sprintf("%05X", o.ToIDA(o.NearCaller())), "si": o.Regs().SI, "state": hex.EncodeToString(o.Bytes(o.IDA(0x1ECFC), 2))})
		}
	})
	runto(0x19C45, 100_000_000)
	write("spawn-rng.bin", o.Bytes(o.IDA(0x1ECFC), 258))
	snapshot("before-spawn", true)
	runto(0x19FA0, 100_000_000)
	snapshot("initialized", true)
	jsonFile("rng-calls-initialized", calls)
	var startCalls []any
	if *traceStart {
		for _, addr := range []uint32{0x1A12A, 0x1A6FA, 0x1ADC8, 0x1AF69, 0x1A7B7, 0x1A85B, 0x1B240, 0x1B732} {
			address := addr
			o.OnCall(o.IDA(address), func(o *oracle.Oracle) {
				if *traceAll && logicalTick == 2 && address == 0x1AF69 {
					// 保留該呼叫 ES 段的碰撞層原始資料；位址基準見同筆 registers。
					write(fmt.Sprintf("collision-tick2-%04x.bin", o.Regs().SI), o.Bytes(oracle.Far(o.Regs().ES, 0), 0x8000))
				}
				limit := 512
				if *traceAll {
					limit = 10000
				}
				if len(startCalls) < limit && (*traceAll || o.Regs().SI == 0 || o.Regs().SI == 0x600 || address == 0x1A12A || address == 0x1ADC8) {
					startCalls = append(startCalls, map[string]any{"tick": logicalTick, "address": fmt.Sprintf("%05X", address), "registers": o.Regs(), "caller": fmt.Sprintf("%05X", o.ToIDA(o.NearCaller())), "unit": hex.EncodeToString(o.Bytes(oracle.Far(o.Word(o.IDA(0x1D30E)), o.Regs().SI), 32))})
				}
			})
		}
	}
	settled := false
	var commands []map[string]any
	commandAddress := uint32(0x1A8F6)
	if *attack {
		commandAddress = 0x1C1B9
	}
	o.OnCall(o.IDA(commandAddress), func(o *oracle.Oracle) {
		commands = append(commands, map[string]any{"tick": logicalTick, "raw_tick": wolong.Tick(o), "si": o.Regs().SI, "caller_ida": fmt.Sprintf("%05X", o.ToIDA(o.NearCaller()))})
	})
	o.OnCall(o.IDA(0x1A065), func(o *oracle.Oracle) {
		if settled {
			return
		}
		if logicalTick < 1000 || logicalTick%50 == 0 {
			snapshot(fmt.Sprintf("tick-%04d", logicalTick), true)
		}
		logicalTick++
	})
	if *traceStart {
		must(o.RunUntil(oracle.NewCond("指定開場拍數完成", func(*oracle.Oracle) bool { return logicalTick >= *traceTicks+1 }), oracle.Budget(100_000_000)))
		jsonFile("start-calls", startCalls)
		return
	}
	runto(0x1A426, 400_000_000)
	snapshot("opening-ready", true)
	jsonFile("rng-calls-opening", calls)
	// 正常按鈕下令全軍退卻，直到結算回大地圖；不注入勝負或改兵力。
	o.OnCall(o.IDA(0x19EBD), func(o *oracle.Oracle) {
		if !settled {
			snapshot("before-settlement", true)
			settled = true
		}
	})
	commandY := 367
	if *attack {
		commandY = 304
	}
	if !*noOrders {
		must(o.Tap(510, commandY))
		if !settled {
			must(o.Tap(510, commandY))
		}
	}
	snapshot("retreat-command", !settled)
	if !settled {
		budget := uint64(100_000_000)
		if *attack || *noOrders {
			budget = 1_500_000_000
		}
		must(o.RunUntil(oracle.NewCond("戰術結算入口", func(*oracle.Oracle) bool { return settled }), oracle.Budget(budget)))
	}
	jsonFile("rng-calls-settlement", calls)
	runto(0x1E453, 100_000_000)
	snapshot("world-return", false)
	jsonFile("command-events", commands)
	fixture := "OpenSiege(35,82)"
	if *field {
		fixture = "受控 SAVE-FIELD.DAT 載入後自然遭遇"
	}
	input := fmt.Sprintf("Tap(510,%d)，結算前至多兩次", commandY)
	if *noOrders {
		input = "進戰術畫面後完全不下令；不是行軍委任"
	}
	jsonFile("result", map[string]any{"completed": true, "fixture": fixture, "player_commands": input, "attack": *attack, "no_orders": *noOrders, "classification": "受控入口；原版結算"})
}
