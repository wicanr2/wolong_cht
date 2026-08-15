package sound

import "testing"

// 目錄不存在時 Bank 仍然可用，而且每個方法都是 no-op。
//
// 這是呼叫端的契約：`cmd/wlgame` 到處呼叫 `g.sound.X()` 而不判斷 nil，
// **也不能因為沒有音檔就崩**（`docs/spec/29` §5）。
// ⚠ 這個測試刻意不建 audio context——無頭環境沒有音效裝置。
func TestMissingDirectoryIsUsable(t *testing.T) {
	b := Open(t.TempDir() + "/不存在")
	if b == nil {
		t.Fatal("Open 不該回傳 nil")
	}
	if b.Available() {
		t.Error("沒有音檔時 Available 應該是 false")
	}
	if len(b.Music()) != 0 {
		t.Error("沒有音檔時 Music 應該是空的")
	}
	// 這幾個呼叫要安靜地什麼都不做。
	b.PlayMusic("bgm-0")
	b.PlayEffect(12)
	b.StopMusic()
	b.SetEnabled(false)
	if err := b.Err(); err != nil {
		t.Errorf("不該有錯誤：%v", err)
	}
}

// nil 的 Bank 也要能用——`Open` 失敗時呼叫端可能拿到 nil。
func TestNilBankIsSafe(t *testing.T) {
	var b *Bank
	b.PlayMusic("bgm-0")
	b.PlayEffect(1)
	b.StopMusic()
	b.SetEnabled(true)
	if b.Available() || b.Enabled() || b.Err() != nil {
		t.Error("nil Bank 的查詢應該全部回零值")
	}
}
