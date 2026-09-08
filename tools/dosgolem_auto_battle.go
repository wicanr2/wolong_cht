//go:build ignore

// 委任野戰的受控取樣；讀存檔自然遭遇，不注入勝負。
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
	root := flag.String("root", "", "唯讀受控素材與委任存檔")
	out := flag.String("out", "", "輸出")
	flag.Parse()
	must := func(e error) {
		if e != nil {
			panic(e)
		}
	}
	must(os.MkdirAll(*out, 0755))
	write := func(n string, b []byte) { must(os.WriteFile(filepath.Join(*out, n), b, 0644)) }
	js := func(n string, v any) { b, e := json.MarshalIndent(v, "", "  "); must(e); write(n+".json", b) }
	o, e := wolong.Load(filepath.Join(*root, "KI.EXE"), *root)
	must(e)
	defer o.Close()
	must(o.RunUntil(wolong.Booted(), oracle.Budget(40_000_000)))
	must(o.Click(320, 200))
	must(o.Click(300, 152))
	must(o.RunUntil(oracle.NewCond("自動判定軍團35與39", func(o *oracle.Oracle) bool {
		r := o.Regs()
		return oracle.At(o.IDA(0x15130)).Ready(o) && r.SI == 0x2240+35*64 && r.DI == 0x2240+39*64
	}), oracle.Budget(400_000_000)))
	ret := o.NearCaller()
	snapshot := func(n string) {
		js(n, map[string]any{"instructions": o.Steps(), "registers": o.Regs(), "corps": wolong.CorpsTable(o), "rng": hex.EncodeToString(o.Bytes(o.IDA(0x1ECFC), 258)), "address_space": "IDA DOS/V linear"})
		must(wolong.Shot(o, filepath.Join(*out, n+".png")))
	}
	snapshot("before-auto")
	save, e := os.ReadFile(filepath.Join(*root, "SAVE.DAT"))
	must(e)
	copy(save[:59], o.Bytes(o.IDA(0x10CF0), 59))
	copy(save[0x80:0x52C0], o.Bytes(oracle.Far(o.Word(o.IDA(0x10D52)), 0), 0x5240))
	copy(save[0x52C0:0x56C0], o.Bytes(oracle.Far(o.Word(o.IDA(0x10D56)), 0), 0x400))
	write("entry-SAVE.DAT", save)
	write("rng.bin", o.Bytes(o.IDA(0x1ECFC), 258))
	var power []any
	o.OnCall(o.IDA(0x152D7), func(o *oracle.Oracle) {
		power = append(power, map[string]any{"stage": "進將領修正", "registers": o.Regs(), "caller": fmt.Sprintf("%05X", o.ToIDA(o.NearCaller()))})
	})
	must(o.RunUntil(oracle.At(ret), oracle.Budget(20_000_000)))
	snapshot("after-auto")
	js("power-calls", power)
	js("result", map[string]any{"completed": true, "return_ida": fmt.Sprintf("%05X", o.ToIDA(ret)), "ax": o.Regs().AX, "mode": "行軍委任自動判定；不是戰術放置"})
	must(o.RunUntil(oracle.At(o.IDA(0x1E453)), oracle.Budget(100_000_000)))
	snapshot("world-return")
}
