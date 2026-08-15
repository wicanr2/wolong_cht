package audio

import (
	"math"
	"os"
	"testing"
)

const origDir = "../../workplace/orig/dosv/"

func loadOriginal(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(origDir + name)
	if err != nil {
		t.Skip("找不到原版 " + name + "，跳過")
	}
	return b
}

func loadTables(t *testing.T) *Tables {
	t.Helper()
	tbl, err := LoadTables(loadOriginal(t, "YNSOUND.COM"))
	if err != nil {
		t.Fatalf("LoadTables: %v", err)
	}
	return tbl
}

// noteHz 把查表的結果換回頻率，公式與 OPL 的規格一致。
func noteHz(tbl *Tables, lo uint8) float64 {
	b0 := tbl.B0[lo&0x7F] | 0x20
	fnum := uint16(b0&3)<<8 | uint16(tbl.A0[lo&0x0F])
	block := (b0 >> 2) & 7
	return float64(fnum) * math.Exp2(float64(block)) * NativeRate / (1 << 20)
}

// 音高表要對得上十二平均律，而且**八度 0–7 都要對**。
//
// ⚠ 這個測試的重點是「跨八度」：只驗一個八度的話，
// 表被截短成 32 筆也照樣會過（`CONTEXT.md` §6 的推翻紀錄）。
func TestPitchTableIsEqualTemperamentAcrossOctaves(t *testing.T) {
	tbl := loadTables(t)
	// 八度要一路翻倍到 7。**這一條才是抓「表被截短」的那一條**：
	// 表只讀 32 筆的話，八度 2 以上會全部退回八度 1。
	for oct := 1; oct < 8; oct++ {
		lo, hi := noteHz(tbl, uint8((oct-1)*16+1)), noteHz(tbl, uint8(oct*16+1))
		if math.Abs(hi-2*lo) > 0.01 {
			t.Errorf("八度 %d 的 C ＝ %.2f Hz，應該是八度 %d（%.2f Hz）的兩倍", oct, hi, oct-1, lo)
		}
	}
	// 半音的比值要是十二平均律。表是量化過的 ROM，容差取 1.5%。
	for oct := 0; oct < 8; oct++ {
		base := noteHz(tbl, uint8(oct*16+1))
		for n := 1; n <= 12; n++ {
			got := noteHz(tbl, uint8(oct*16+n))
			ideal := base * math.Exp2(float64(n-1)/12)
			if math.Abs(got-ideal)/ideal > 0.015 {
				t.Errorf("八度 %d 第 %d 音 ＝ %.2f Hz，十二平均律是 %.2f Hz", oct, n, got, ideal)
			}
		}
	}
	// A4 是最好記的錨點。
	if a4 := noteHz(tbl, 4*16+10); math.Abs(a4-440) > 5 {
		t.Errorf("八度 4 的 A ＝ %.2f Hz，離 440 太遠", a4)
	}
}

// 長度表要是 192 的二分序列，附點那一組正好是 0.75 倍。
func TestLengthTableIsBinaryWithDots(t *testing.T) {
	tbl := loadTables(t)
	for i := 0; i < 6; i++ {
		if int(tbl.Length[i])/2 != int(tbl.Length[i+1]) {
			t.Fatalf("長度表第 %d 筆 %d 不是前一筆的一半", i+1, tbl.Length[i+1])
		}
	}
	if got, want := int(tbl.Length[9]), int(tbl.Length[0])*3/4; got != want {
		t.Errorf("附點全音符 ＝ %d，應該是 %d", got, want)
	}
}

// 每一首用到的音色編號都要落在該曲音色表的範圍內。
//
// 這是「音色記錄是 32 bytes」的獨立檢查：大小猜錯的話筆數就對不上。
func TestInstrumentIndicesStayInsideTable(t *testing.T) {
	for _, name := range []string{"BGM.DAT", "OPENBGM.DAT", "ENDBGM.DAT", "OVERBGM.DAT"} {
		songs, err := Songs(loadOriginal(t, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		for n := range songs {
			s := &songs[n]
			count := 0
			for s.Instrument(count) != nil {
				count++
			}
			for _, start := range s.Tracks {
				if start == emptyTrack {
					continue
				}
				for i := start; i+2 <= s.Instruments; i += 2 {
					lo, hi := s.Data[i], s.Data[i+1]
					if lo >= 0x80 && (lo>>4)&7 == 2 && int(hi) >= count {
						t.Errorf("%s 第 %d 曲用了音色 %d，但表只有 %d 筆", name, n, hi, count)
					}
				}
			}
		}
	}
}

// 渲染出來的東西要有聲音，而且**每一段都要有**。
//
// 只看總 RMS 會漏掉「只有片頭有聲音」這種失敗（`docs/playtest/25` §3）。
func TestRenderProducesSoundInEverySegment(t *testing.T) {
	if testing.Short() {
		t.Skip("渲染很慢，-short 跳過")
	}
	tbl := loadTables(t)
	songs, err := Songs(loadOriginal(t, "BGM.DAT"))
	if err != nil {
		t.Fatal(err)
	}
	p := NewPlayer(&songs[0], tbl, NewOPL3(0))
	pcm := Render(p, 8)
	if pcm.TickCount < 300 {
		t.Errorf("8 秒只走了 %d 個 tick，速度換算可能錯了", pcm.TickCount)
	}
	seg := SegmentRMS(pcm, 8)
	// 第一段可能是曲子開頭的休止，從第二段開始要求有訊號。
	for i := 1; i < len(seg); i++ {
		if seg[i] < 1e-4 {
			t.Errorf("第 %d 段是靜音（RMS %.6f）", i, seg[i])
		}
	}
}

// 四種 4-operator 連接的載波集合要與原版那張遮罩表一致。
//
// 這一條是把 `docs/re/57` §1 的證據釘在程式碼裡：遮罩表是原版讀的，
// `generateFourOp` 的拓樸是照規格寫的，兩邊必須指向同一件事。
func TestCarrierMaskMatchesFourOpTopology(t *testing.T) {
	want := [4]uint8{0x08, 0x0A, 0x09, 0x0D}
	if carrierMask != want {
		t.Fatalf("載波遮罩 %v 與原版的 08 0A 09 0D 不符", carrierMask)
	}
}

// 開了 4-operator 之後，後半通道不能自己出聲。
func TestFourOpSecondChannelIsNotMixedTwice(t *testing.T) {
	c := NewOPL3(0)
	c.Write(1, 0x05, 0x01)
	c.Write(1, 0x04, 0x3F)
	for _, n := range []int{3, 4, 5, 12, 13, 14} {
		if !c.isFourOpSecond(n) {
			t.Errorf("通道 %d 應該是 4-op 的後半", n)
		}
		if ops := c.channelOperators(n); ops != nil {
			t.Errorf("通道 %d 是後半，不該有自己的 operator，卻回傳 %v", n, ops)
		}
	}
	for _, n := range []int{0, 1, 2, 9, 10, 11} {
		if ops := c.channelOperators(n); len(ops) != 4 {
			t.Errorf("通道 %d 是 4-op 的前半，應該有四個 operator，卻有 %d 個", n, len(ops))
		}
	}
	// NEW 關掉就退回 2-op。
	c.Write(1, 0x05, 0x00)
	if ops := c.channelOperators(0); len(ops) != 2 {
		t.Errorf("NEW 關掉後通道 0 應該只有兩個 operator，卻有 %d 個", len(ops))
	}
}
