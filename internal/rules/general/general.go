// Package general 實作與武將有關的規則。
//
// 目前只有「評價」——原版 `sub_155A6`，在載入劇本時與每次月結都會重算，
// 所以它是**衍生值**，不存在存檔裡（劇本檔的 `+1Fh` 開局是 0）。
// 反組譯見 docs/re/07 §14，機制說明見 docs/mechanics/60-personnel.md。
package general

// General 是一名武將在規則層需要的部分。
// 欄位對應武將記錄的 +0Eh…+13h（docs/formats/08 §3）。
type General struct {
	Name string

	// Aptitude 是 +0Eh／+0Fh／+10h 三個欄位。
	// 原版存成 16 的倍數（0–160），讀出來時 `>>4` 變成 0–10。
	// 這裡存的是**已經除過的小值**。
	//
	// 三個欄位疑似**兵種適性**（騎馬／弓／步）：123 個武將裡有 80 人
	// 只有一欄非零、32 人兩欄、9 人三欄，形狀比較像適性而不是能力值。
	// **強證據，還沒驗**——要與戰鬥程式對照才能定案。
	Aptitude [3]int

	Martial  int // +11h 武術，1–15
	Command  int // +12h 統率，1–15
	Politics int // +13h 政治，1–15
}

// Rating 回傳武將的「評價」。
//
//	評價 = 適性₁ + 適性₂ + 適性₃ + 2 × 武術 + 2 × 統率
//
// ⭐ **武術與統率的權重相同**，所以只有兩者的和有意義——
// 這正是說明書 10.5 那句「武術と統率はその合計が同じであれば強さは同じ」。
//
// **政治完全不計入**，與說明書「政治は内政、外交に関係する」一致。
// 所以純文官（武 1 統 1 政 14）的評價會排在最後面，
// 而評價高的是呂布、趙雲、諸葛亮、關羽這一批。
//
// 原版把結果存進一個 byte（`mov [si+1Fh], ah`），所以會在 256 溢位；
// 實際值域遠低於此（劇本 1 的最高是 66），這裡不模擬溢位。
func (g General) Rating() int {
	return g.Aptitude[0] + g.Aptitude[1] + g.Aptitude[2] +
		2*g.Martial + 2*g.Command
}

// PreferForPlayerCommand 回報「玩家親自指揮時」該不該優先選這名武將。
//
// 說明書 10.5：玩家指揮時選**武術**高的，委任時選**統率**高的。
// 理由是「戦術では、プレーヤーの指揮次第で統率力をカバー出来るのに対し、
// 武術は武将個人の強さのため指揮ではどうにもならない」。
//
// 兩人評價相同時（＝武術＋統率相同），武術高的那個在玩家手上更強。
func PreferForPlayerCommand(a, b General) bool {
	if a.Rating() != b.Rating() {
		return a.Rating() > b.Rating()
	}
	return a.Martial > b.Martial
}

// PreferForDelegation 回報「委任給 AI 時」該不該優先選這名武將。
// 與 PreferForPlayerCommand 相反，看的是統率。
func PreferForDelegation(a, b General) bool {
	if a.Rating() != b.Rating() {
		return a.Rating() > b.Rating()
	}
	return a.Command > b.Command
}
