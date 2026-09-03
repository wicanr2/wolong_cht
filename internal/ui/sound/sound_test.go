package sound

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

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

// TestSilentSkipsPlayback 釘住 docs/spec/29 §5.1：驗收模式**不出聲，
// 但狀態照舊**。
//
// ⚠ 這一支擋的不是「有沒有聲音」，是**驗收捷徑會不會改到被驗收的畫面**。
// 先前的作法是把音檔目錄清空，於是 Available() 變成 false、
// 系統選單那一格從「TYPE 1」變成「未接入」，對拍差 272 px——
// 而那 272 px 不是回歸，是驗收路徑自己造成的。
func TestSilentSkipsPlayback(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "track-01.ogg"), []byte("not really ogg"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sfx-03.ogg"), []byte("not really ogg"), 0o600); err != nil {
		t.Fatal(err)
	}
	b := Open(dir)
	if !b.Available() {
		t.Fatal("掃到了兩個 ogg，Available 應為真")
	}
	b.SetSilent(true)
	if !b.Silent() {
		t.Error("SetSilent(true) 之後 Silent() 應為真")
	}
	// ⭐ 靜音模式**不改**這兩個——系統選單那一格看的就是它們。
	if !b.Available() {
		t.Error("靜音模式把 Available 弄成 false 了，那一格會印成「未接入」")
	}
	if !b.Enabled() {
		t.Error("靜音模式不該改 Enabled")
	}
	// 檔案是壞的 ogg：真的走到解碼就會設 initErr，靜音模式應該連碰都不碰。
	b.PlayMusic("track-01")
	b.PlayEffect(3)
	if err := b.Err(); err != nil {
		t.Errorf("靜音模式仍然走到了播放層：%v", err)
	}
}

// 四段增益要對得上原版的 OPL Total Level 步進（docs/spec/122 §2）。
func TestLevelGainMatchesOPLSteps(t *testing.T) {
	// 一段 ＝ 4 個 TL 單位 × 0.75 dB ＝ 3 dB。
	want := []float64{1, 0.70794578, 0.50118723, 0.35481339}
	for n, w := range want {
		got := LevelGain(n)
		if math.Abs(got-w) > 1e-6 {
			t.Errorf("LevelGain(%d) ＝ %g，應為 %g", n, got, w)
		}
	}
	for n := 1; n < len(want); n++ {
		if LevelGain(n) >= LevelGain(n-1) {
			t.Errorf("TYPE %d 應該比 TYPE %d 小聲", n+1, n)
		}
	}
}

func TestSetLevelClamps(t *testing.T) {
	var nilBank *Bank
	nilBank.SetLevel(2) // 不能炸
	if nilBank.Level() != 0 {
		t.Error("nil Bank 的 Level 應為 0")
	}
	b := Open(t.TempDir())
	for _, c := range []struct{ in, want int }{
		{-3, 0}, {0, 0}, {3, 3}, {9, 3},
	} {
		b.SetLevel(c.in)
		if got := b.Level(); got != c.want {
			t.Errorf("SetLevel(%d) → Level() ＝ %d，應為 %d", c.in, got, c.want)
		}
	}
}
