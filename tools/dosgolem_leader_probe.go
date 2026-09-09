//go:build ignore

// 在自然委任遭遇的原版將領修正入口做狀態對照；不修改機器碼或強制勝負。
package main

import (
	"crypto/sha256"
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
	root := flag.String("root", "", "受控委任存檔素材根目錄")
	out := flag.String("out", "", "收據 JSON")
	flag.Parse()
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	o, err := wolong.Load(filepath.Join(*root, "KI.EXE"), *root)
	must(err)
	defer o.Close()
	must(o.RunUntil(wolong.Booted(), oracle.Budget(40_000_000)))
	must(o.Click(320, 200))
	must(o.Click(300, 152))
	must(o.RunUntil(oracle.At(o.IDA(0x152E7)), oracle.Budget(400_000_000)))
	start := o.Save()
	regs := o.Regs()
	stats := oracle.Far(regs.DS, regs.BX+0x4251)
	var rows []any
	for _, pair := range [][2]uint8{{0, 0}, {1, 0}, {12, 6}, {15, 11}, {15, 15}, {4, 12}, {0, 15}} {
		seen := map[bool]bool{}
		for counter := 0; counter < 256; counter++ {
			o.Restore(start)
			o.WriteBytes(stats, pair[:])
			o.WriteU8(o.IDA(0x1ECFC), uint8(counter))
			rngBefore := hex.EncodeToString(o.Bytes(o.IDA(0x1ECFC), 258))
			must(o.RunUntil(oracle.At(o.IDA(0x15308)), oracle.Budget(100_000)))
			r := o.Regs()
			mixed := pair[0] < pair[1] || r.AX&3 == 0
			if seen[mixed] {
				continue
			}
			seen[mixed] = true
			rows = append(rows, map[string]any{"martial": pair[0], "command": pair[1],
				"counter": counter, "rng_before": rngBefore,
				"rng_after": hex.EncodeToString(o.Bytes(o.IDA(0x1ECFC), 258)),
				"mixed":     mixed, "leader_value": r.CX >> 8, "registers": r})
			if pair[0] < pair[1] || len(seen) == 2 {
				break
			}
		}
		if pair[0] >= pair[1] && len(seen) != 2 {
			panic("未覆蓋兩個 RNG 分支")
		}
	}
	exe, err := os.ReadFile(filepath.Join(*root, "KI.EXE"))
	must(err)
	data, err := json.MarshalIndent(map[string]any{"classification": "自然委任入口的窄狀態對照；不是完整玩家路徑", "tool": "dosgolem",
		"address_space": "IDA DOS/V linear", "entry": "000152E7", "end": "00015308",
		"exe_sha256": fmt.Sprintf("%x", sha256.Sum256(exe)), "original_registers": regs, "cases": rows}, "", "  ")
	must(err)
	must(os.WriteFile(*out, data, 0644))
}
