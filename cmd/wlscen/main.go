// wlscen 把劇本／存檔區塊在二進位與 JSON 之間互轉（docs/spec/28）。
//
//	wlscen export    -in SINARIO.DAT -block 0 -out s0.json
//	wlscen import    -in SINARIO.DAT -block 0 -json s0.json -out new.DAT
//	wlscen roundtrip -in SINARIO.DAT
//
// ⭐ **匯入一定要有原始檔**：JSON 只含已解欄位，未解區域靠 `-in` 原樣保留。
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/wolong_cht/internal/scenario"
	"github.com/wicanr2/wolong_cht/internal/state"
)

const numBlocks = 4

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	in := fs.String("in", "", "來源 SINARIO.DAT／SAVE.DAT（必填）")
	out := fs.String("out", "", "輸出路徑")
	jsonPath := fs.String("json", "", "import 用的 JSON")
	block := fs.Int("block", 0, "區塊編號 0–3")
	indent := fs.Bool("indent", true, "JSON 縮排（給人看的預設開）")
	if err := fs.Parse(os.Args[2:]); err != nil {
		os.Exit(2)
	}
	if *in == "" {
		usage()
	}

	var err error
	switch cmd {
	case "export":
		err = doExport(*in, *out, *block, *indent)
	case "import":
		err = doImport(*in, *jsonPath, *out, *block)
	case "roundtrip":
		err = doRoundTrip(*in)
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `用法：
  wlscen export    -in <劇本檔> -block <0-3> -out <out.json>
  wlscen import    -in <劇本檔> -block <0-3> -json <in.json> -out <out.DAT>
  wlscen roundtrip -in <劇本檔>

匯入一定要有原始檔：JSON 只含已解欄位，未解區域靠 -in 原樣保留。`)
	os.Exit(2)
}

func fileSHA(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

// load 讀一個區塊並回傳 World ＋ 來源雜湊。
func load(in string, block int) (*state.World, string, error) {
	w, err := state.LoadScenario(in, block)
	if err != nil {
		return nil, "", err
	}
	sum, err := fileSHA(in)
	if err != nil {
		return nil, "", err
	}
	return w, sum, nil
}

func doExport(in, out string, block int, indent bool) error {
	w, sum, err := load(in, block)
	if err != nil {
		return err
	}
	s := scenario.FromWorld(w, scenario.Meta{
		Source: filepath.Base(in), Block: block, SourceSHA: sum})

	var data []byte
	if indent {
		data, err = json.MarshalIndent(s, "", "  ")
	} else {
		data, err = json.Marshal(s)
	}
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if out == "" || out == "-" {
		_, err = os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(out, data, 0o644)
}

func doImport(in, jsonPath, out string, block int) error {
	if jsonPath == "" || out == "" {
		return fmt.Errorf("import 需要 -json 與 -out")
	}
	// ⚠ 不覆蓋來源。原版資產一律唯讀（CLAUDE.md §9）。
	if sameFile(in, out) {
		return fmt.Errorf("-out 不能與 -in 相同：%s", out)
	}
	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	var s scenario.Scenario
	if err := json.Unmarshal(raw, &s); err != nil {
		return err
	}
	w, _, err := load(in, block)
	if err != nil {
		return err
	}
	if err := s.ApplyTo(w); err != nil {
		return err
	}
	return w.SaveInto(in, out, block)
}

// doRoundTrip 是驗收：四個區塊各自 export → import，輸出必須與輸入相同。
func doRoundTrip(in string) error {
	orig, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	for block := 0; block < numBlocks; block++ {
		w, sum, err := load(in, block)
		if err != nil {
			return err
		}
		before := w.Bytes()

		data, err := json.Marshal(scenario.FromWorld(w, scenario.Meta{
			Source: filepath.Base(in), Block: block, SourceSHA: sum}))
		if err != nil {
			return err
		}
		var s scenario.Scenario
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w2, _, err := load(in, block)
		if err != nil {
			return err
		}
		if err := s.ApplyTo(w2); err != nil {
			return err
		}
		after := w2.Bytes()
		if diff := firstDiff(before, after); diff >= 0 {
			return fmt.Errorf("區塊 %d 在 +0x%X 不一致：%02X → %02X",
				block, diff, before[diff], after[diff])
		}
		fmt.Printf("區塊 %d：%d bytes 完全相同 ✓（標題 %q）\n",
			block, len(after), scenario.FromWorld(w, scenario.Meta{}).Meta.Title)
	}
	fmt.Printf("來源 %s 共 %d bytes，四個區塊 round-trip 全過\n", in, len(orig))
	return nil
}

func firstDiff(a, b []byte) int {
	if len(a) != len(b) {
		return 0
	}
	for i := range a {
		if a[i] != b[i] {
			return i
		}
	}
	return -1
}

func sameFile(a, b string) bool {
	ai, err := os.Stat(a)
	if err != nil {
		return false
	}
	bi, err := os.Stat(b)
	if err != nil {
		return false
	}
	return os.SameFile(ai, bi)
}
