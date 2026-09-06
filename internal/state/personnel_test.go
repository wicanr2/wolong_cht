package state

import "testing"

// 職務 `+0x17` 是 0–4 的值，不是布林（docs/spec/143）。
// ⭐ 四個劇本的初始值全是 0，所以**只比劇本存檔驗不出壓縮**——
// 這裡把非零值灌進原始 bytes 再 round-trip，才擋得住「2／3／4 被壓成 1」。
func TestGeneralDutyRoundTripsAllValues(t *testing.T) {
	w := load(t, 0)
	want := []byte{0, 1, 2, 3, 4}
	for i := range want {
		off := generalBase + i*generalSize
		w.raw[off+0x17] = want[i]
	}
	w2 := loadBlock(w.raw)
	for i, v := range want {
		if got := w2.Generals[i].Duty; got != int(v) {
			t.Fatalf("武將 %d 載入的職務 = %d，want %d", i, got, v)
		}
	}
	out := w2.Bytes()
	for i, v := range want {
		off := generalBase + i*generalSize
		if out[off+0x17] != v {
			t.Fatalf("武將 %d 寫回的 +0x17 = %#02x，want %#02x", i, out[off+0x17], v)
		}
	}
}

// Posted() 是「有職務」，四個非零值都算。
func TestPostedIsAnyNonZeroDuty(t *testing.T) {
	for duty := 0; duty <= 4; duty++ {
		g := General{Duty: duty}
		if got, want := g.Posted(), duty != DutyNone; got != want {
			t.Fatalf("職務 %d 的 Posted() = %v，want %v", duty, got, want)
		}
	}
}

// 任命寫兩格（職務 ＋ 據點 +0x19），解任清三格（職務、經費、據點）。
// **經費歸零**是原版 `sub_16B4F` 的 `mov byte [bx+1Ah], 0`（docs/re/25 §3.1）。
func TestAssignAndDismissGovernor(t *testing.T) {
	w := load(t, 0)
	city, who := 0, 3
	w.Cities[city].Governor = noGovernorSlot
	w.Generals[who].Duty = DutyNone
	w.Generals[who].Budget = 40

	if !w.AssignGovernor(city, who) {
		t.Fatal("任命失敗")
	}
	if got := w.Generals[who].Duty; got != DutyGovernor {
		t.Fatalf("任命後職務 = %d，want %d", got, DutyGovernor)
	}
	if got := w.Cities[city].Governor; got != who {
		t.Fatalf("據點 +0x19 = %d，want %d", got, who)
	}
	// 那裡已經有人就不覆蓋（原版由呼叫端跳 TALK #52 擋掉）。
	if w.AssignGovernor(city, 4) {
		t.Fatal("已有內政官還任命成功")
	}

	if got := w.DismissGovernor(city); got != who {
		t.Fatalf("解任回傳 %d，want %d", got, who)
	}
	if got := w.Generals[who].Duty; got != DutyNone {
		t.Fatalf("解任後職務 = %d，want 0", got)
	}
	if got := w.Generals[who].Budget; got != 0 {
		t.Fatalf("解任後經費 = %d，want 0（原版沒收）", got)
	}
	if got := w.Cities[city].Governor; got != noGovernorSlot {
		t.Fatalf("解任後據點 +0x19 = %#02x，want %#02x", got, noGovernorSlot)
	}
	// 本來就沒人 → 回 −1，而且不會誤傷別人。
	if got := w.DismissGovernor(city); got != noGovernor {
		t.Fatalf("空缺解任回傳 %d，want %d", got, noGovernor)
	}
}

func TestAssignAndDismissDiplomat(t *testing.T) {
	w := load(t, 0)
	faction, who := 1, 5
	w.Factions[faction].Diplomat = noFaction
	w.Generals[who].Duty = DutyNone
	w.Generals[who].Budget = 25

	if !w.AssignDiplomat(faction, who) {
		t.Fatal("派駐失敗")
	}
	if got := w.Generals[who].Duty; got != DutyDiplomat {
		t.Fatalf("派駐後職務 = %d，want %d", got, DutyDiplomat)
	}
	if w.AssignDiplomat(faction, 6) {
		t.Fatal("已有外交官還派駐成功")
	}
	if got := w.DismissDiplomat(faction); got != who {
		t.Fatalf("召回回傳 %d，want %d", got, who)
	}
	if w.Generals[who].Duty != DutyNone || w.Generals[who].Budget != 0 {
		t.Fatalf("召回後職務 = %d、經費 = %d，want 0／0",
			w.Generals[who].Duty, w.Generals[who].Budget)
	}
	if got := w.DismissDiplomat(faction); got != noGovernor {
		t.Fatalf("空缺召回回傳 %d，want %d", got, noGovernor)
	}
}

// 選了軍師，那個人就從武將表消失（docs/spec/144）：
// 原版 `loc_11AF8` 寫的是記錄 `+0x00`，不是 `+0x17` —— 存在旗標一起沒了，
// 所以每張清單的第一個條件 `[si] >= 0x80` 一次擋掉他。
func TestTakeAdvisorRemovesHimFromTheGeneralTable(t *testing.T) {
	w := load(t, 0)
	f := &w.Factions[0]
	who := f.Advisor
	if who < 0 || who >= len(w.Generals) || who == NoAdvisor {
		t.Skipf("劇本一的勢力 0 沒有預設軍師（%d）", who)
	}
	before := f.Generals
	if !w.Generals[who].Alive {
		t.Fatalf("軍師 %d 在劇本檔裡就不存在，樣本不對", who)
	}

	w.TakeAdvisor(0, who)

	if f.Generals != before-1 {
		t.Errorf("武將數 = %d，want %d", f.Generals, before-1)
	}
	if w.Generals[who].Alive {
		t.Error("軍師還在武將表上")
	}
	// 原版寫的是整個 `+0x00`，四個旗標一起沒。
	out := w.Bytes()
	off := generalBase + who*generalSize
	if got := out[off] &^ 0x01; got != 0 {
		t.Errorf("寫回的 +0x00 = %#02x，除了未解的 bit 0 之外應該全 0", out[off])
	}
}

// 32 筆物件記錄要 byte-for-byte round-trip（docs/spec/146 §2）。
// ⭐ 後 16 筆是常駐的雲，劇本檔裡就有值——先前 remake 不讀這一段，
// 靠「沒碰過就不會壞」蒙混；現在會寫回去了，就得真的比。
func TestMapObjectsRoundTrip(t *testing.T) {
	for idx := 0; idx < 4; idx++ {
		w := load(t, idx)
		out := w.Bytes()
		for i := 0; i < 32*16; i++ {
			off := mapObjectBase + i
			if out[off] != w.raw[off] {
				t.Fatalf("劇本 %d 物件區 +%#04x：%#02x != %#02x",
					idx+1, i, out[off], w.raw[off])
			}
		}
		// 正對照：那一段不是全 0，否則這支測試等於沒比。
		clouds := 0
		for slot := 16; slot < 32; slot++ {
			if w.raw[mapObjectBase+slot*16]&0x80 != 0 {
				clouds++
			}
		}
		if clouds != 16 {
			t.Fatalf("劇本 %d 只有 %d 朵雲，want 16", idx+1, clouds)
		}
	}
}
