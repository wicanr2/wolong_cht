// Package savefile 是 remake 的原生存檔格式。
//
// 規格：`docs/spec/20-save-format.md`。三個已裁定的取捨：
// **JSON**（可 diff、可手改造 fixture）、**不壓縮**（壓了就看不見）、
// **沒有原版檔案也載得動**（原版檔案的驗證只在啟動時做一次）。
//
// 一條不可退讓的約束：**原生存檔必須能無損匯出回原版格式**。
// 所以檔案帶著載入當時那 22,208 B 原始區塊——少了它，
// 那 7 KB 未解區與所有還沒解的欄位在匯出時只能填 0，
// 等於把「改寫不是重建」這條硬規則作廢（`CLAUDE.md` §9）。
package savefile

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/state"
)

// Format 是檔案格式識別字。改動版面時往上加，**不要沉默地換語意**。
const Format = "wolong-save/1"

// Origin 記錄這份存檔從哪一個原版區塊長出來的。
type Origin struct {
	// Source 是來源檔名（`SINARIO.DAT` 或 `SAVE.DAT`），只作紀錄。
	Source string `json:"source"`
	// Block 是劇本／存檔槽編號。
	Block int `json:"block"`
	// BlockSHA256 是 Raw 的雜湊。**它驗的是這個檔案自己有沒有壞**，
	// 不是玩家有沒有原版——後者在啟動時驗一次（決策 C）。
	BlockSHA256 string `json:"block_sha256"`
}

// File 是原生存檔的完整版面。
type File struct {
	Format string `json:"format"`
	Origin Origin `json:"origin"`

	// Raw 是載入當時的原始區塊。**保存錨點**：未解區域原樣帶著走。
	Raw []byte `json:"raw"`

	// State 是已解欄位的可讀版本。與 Raw 重複是**刻意的**——
	// 重複才有東西可以互相對（`docs/spec/20` §2.3）。
	State json.RawMessage `json:"state"`

	// Runtime 是原版記錄裡沒有的 remake 狀態。
	Runtime state.Snapshot `json:"runtime"`
}

// Encode 把世界寫成原生存檔。
func Encode(w *state.World, origin Origin) ([]byte, error) {
	if w == nil {
		return nil, fmt.Errorf("世界是 nil")
	}
	raw := w.RawBlock()
	sum := sha256.Sum256(raw)
	origin.BlockSHA256 = hex.EncodeToString(sum[:])

	// State 只是給人看與給 diff 用的鏡子，所以直接序列化 World 的公開欄位。
	// 權威仍在 Raw ＋ 已解欄位的寫回路徑（w.Bytes()）。
	view, err := json.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("序列化世界狀態：%w", err)
	}
	f := File{
		Format:  Format,
		Origin:  origin,
		Raw:     raw,
		State:   view,
		Runtime: w.TakeSnapshot(),
	}
	return json.MarshalIndent(f, "", "  ")
}

// Decode 從原生存檔建回世界。
//
// **不合就大聲失敗。** 版本不對、雜湊不符、Runtime 索引超界，
// 三種都回錯誤而不是靜靜挑一邊——那代表檔案被手改過或版本不合，
// 而默默採用會產生一個「看起來能玩但行為不對」的世界。
func Decode(data []byte) (*state.World, *File, error) {
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, nil, fmt.Errorf("解析存檔：%w", err)
	}
	if f.Format != Format {
		return nil, nil, fmt.Errorf("存檔格式 %q，這個版本只認得 %q", f.Format, Format)
	}
	sum := sha256.Sum256(f.Raw)
	if got := hex.EncodeToString(sum[:]); got != f.Origin.BlockSHA256 {
		return nil, nil, fmt.Errorf("原始區塊雜湊不符：檔案寫 %s，實際 %s",
			f.Origin.BlockSHA256, got)
	}
	w, err := state.LoadBlock(f.Raw)
	if err != nil {
		return nil, nil, err
	}
	if err := w.Restore(f.Runtime); err != nil {
		return nil, nil, fmt.Errorf("回復執行期狀態：%w", err)
	}
	return w, &f, nil
}

// ExportOriginal 把世界寫成原版格式的 22,208 B 區塊。
//
// 這就是 `w.Bytes()`：從 Raw 出發、只蓋已解欄位。留一支具名函式是為了
// 讓「匯出」在呼叫端讀得出意圖，**不要在別處另外寫一份**
// （`CLAUDE.md` §7 第 6 條）。
func ExportOriginal(w *state.World) []byte { return w.Bytes() }

// VerifyExport 檢查一份原生存檔匯出之後與預期的區塊一致。
// 給驗收流程用：**round-trip 的唯一可信驗法是 byte-for-byte 比對**。
func VerifyExport(w *state.World, want []byte) error {
	got := ExportOriginal(w)
	if !bytes.Equal(got, want) {
		for i := range got {
			if i >= len(want) {
				break
			}
			if got[i] != want[i] {
				return fmt.Errorf("匯出結果從位移 0x%04X 起不同：%02X ≠ %02X",
					i, got[i], want[i])
			}
		}
		return fmt.Errorf("匯出結果長度 %d，預期 %d", len(got), len(want))
	}
	return nil
}
