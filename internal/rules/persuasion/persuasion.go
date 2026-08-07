// Package persuasion 實作「進言」與「說得」。
//
// 這是本作的招牌機制：**玩家扮演軍師，不是君主**。所有戰略指令都不是
// 直接執行，而是向君主提議；君主可能拒絕，然後玩家挑理由說服他。
//
// 規格出自日文原版說明書 3.9 節（PDF p.16–18），機制說明見
// docs/mechanics/70-ai.md §1。好戰等級的欄位（勢力記錄 +0x28）
// 已在機器碼裡定案，見 docs/formats/08 §1.5。
package persuasion

// Command 是三種需要說服的進言。
type Command int

const (
	Hostility Command = iota // 敵對提案
	CeaseFire                // 停戰提案
	Cooperate                // 協力要請
)

func (c Command) String() string {
	return [...]string{"敵對提案", "停戰提案", "協力要請"}[c]
}

// Reason 是說服時可以提出的理由。說明書 3.9 一共列了九個，
// 外加一個「進言撤回」。
type Reason int

const (
	FriendshipBad     Reason = iota // 交友關係惡い
	FriendshipGood                  // 交友關係良い
	WeAreStronger                   // 我が国有利
	EnemyIsStronger                 // 敵国強大
	EnemyInvading                   // 敵が他国侵攻中
	WeAreDefending                  // 我が国防戦中
	EnemyExhausted                  // 敵が疲弊中
	WeAreExhausted                  // 我が国疲弊
	AllyIsStronger                  // 協力国強大
	InvaderIsStronger               // 侵攻国強大
	Withdraw                        // 進言撤回
	numReasons
)

func (r Reason) String() string {
	return [...]string{
		"交友關係惡", "交友關係良", "我國有利", "敵國強大",
		"敵在侵攻他國", "我國防戰中", "敵已疲弊", "我國疲弊",
		"協力國強大", "侵攻國強大", "進言撤回",
	}[r]
}

// Situation 是判斷各理由是否「符合狀況」需要的局勢。
//
// 欄位都對應已解出的資料：國力用**據點數**（說明書 3.9 明說
// 「国力の基準は拠点数からチェックされ」），疲弊用**資金 < 0**
// （機器碼用 24 位值的符號位判斷，docs/formats/08 §1.5）。
type Situation struct {
	// Aggression 是自家君主的好戰等級（0–15，勢力記錄 +0x28）。
	// **不顯示給玩家**——玩家只能靠被拒絕的次數去推。
	Aggression int

	OurCities, TheirCities int // 國力 ＝ 據點數
	AllyCities             int // 協力勢力的據點數
	InvaderCities          int // 侵攻勢力的據點數

	OurFunds, TheirFunds int // 疲弊 ＝ 資金 < 0

	// Friendship 是自家君主看對方的交友值（0–100）。
	Friendship int

	TheyInvadeThirdParty bool // 對方正在侵攻第三方
	TheyInvadeUs         bool // 對方正在侵攻我方
}

// badFriendshipGate 是「交友關係惡」成立的交友值上限。
//
// 說明書 3.9 舉了兩個例子：
//
//	例えば、呂布などの場合は特に悪くない状態でも説得可能ですし、
//	劉備などの場合は逆にかなり悪くならないと説得できません。
//
// **好戰等級越高，門檻越高**——呂布（15）在關係「不特別差」時就買帳，
// 劉備（4）要關係「相當差」才買帳。
//
// ⚠ 第一版我把方向寫反了（寫成好戰等級越高門檻越低），
// 照說明書那兩個例子寫的測試當場抓到。**實際換算還沒反組譯出來**，
// 這條線性式是 remake 的暫定值，只保證方向對。
func badFriendshipGate(aggression int) int { return 10 + aggression*2 }

// 國力比較的修正量。說明書：「好戦レベルが高い場合は多少拠点数が
// 少なくとも有利とみなします」。
//
// ⚠ 同樣是暫定值——只保證方向對。
func powerBonus(aggression int) int { return aggression / 3 }

// Applies 回報某個理由在這個局勢下是否「符合狀況」。
//
// **選到不符合的理由 → 說服失敗、信賴度下降。**
// 所以這個函式的正確性直接決定玩家會不會被冤枉扣分。
func (s Situation) Applies(r Reason) bool {
	switch r {
	case FriendshipBad:
		return s.Friendship < badFriendshipGate(s.Aggression)
	case FriendshipGood:
		return s.Friendship >= badFriendshipGate(s.Aggression)
	case WeAreStronger:
		return s.OurCities+powerBonus(s.Aggression) > s.TheirCities
	case EnemyIsStronger:
		return s.TheirCities > s.OurCities+powerBonus(s.Aggression)
	case EnemyInvading:
		return s.TheyInvadeThirdParty
	case WeAreDefending:
		return s.TheyInvadeUs
	case EnemyExhausted:
		return s.TheirFunds < 0
	case WeAreExhausted:
		return s.OurFunds < 0
	case AllyIsStronger:
		return s.AllyCities > s.OurCities
	case InvaderIsStronger:
		return s.InvaderCities > s.OurCities
	case Withdraw:
		return true // 隨時可用
	}
	return false
}

// 每個指令會用到的理由池（說明書 3.9 的分類）。
// 畫面上一次只給五個，從對應的池裡挑。
var pools = map[Command][]Reason{
	Hostility: {FriendshipBad, WeAreStronger, EnemyInvading, EnemyExhausted, WeAreDefending},
	CeaseFire: {EnemyIsStronger, WeAreExhausted, FriendshipGood, InvaderIsStronger, WeAreDefending},
	Cooperate: {FriendshipGood, AllyIsStronger, InvaderIsStronger, WeAreExhausted, EnemyInvading},
}

// Options 回傳這個指令會顯示的五個理由，外加「進言撤回」。
//
// 說明書：「常に 5 つの項目が選択肢として用意されています」。
// **不符合狀況的理由也會出現在選項裡**——那正是這個系統的難處：
// 玩家要自己判斷哪些成立，選錯就扣信賴度。
func Options(c Command) []Reason {
	return append(append([]Reason{}, pools[c]...), Withdraw)
}

// Outcome 是一次說服動作的結果。
type Outcome int

const (
	Continue  Outcome = iota // 理由成立，但君主還沒點頭
	Agreed                   // 君主同意了
	Failed                   // 選到不符合狀況的理由 → 信賴度下降
	Withdrawn                // 用進言撤回收手 → 信賴度不變
)

// 信賴度的增減量。
//
// ⚠ **實際數值還沒反組譯出來**（說明書只說「下降」「上昇」
// 「大幅に上昇」）。這裡的值是 remake 的暫定值，
// 只保證「外交成功 ≫ 進言成功 > 0 > 失敗」的相對關係。
const (
	TrustOnSuccess       = 2
	TrustOnFailure       = -5
	TrustOnDiplomaticWin = 20 // 停戰／協力外交成功後另外給
)

// Session 是一次進行中的說服。
type Session struct {
	Command   Command
	Situation Situation

	// used 記住已經講過的理由，同一個理由不能重複用。
	used [numReasons]bool
	// need 是還要講幾個成立的理由君主才會點頭。
	need int
}

// requiredReasons 是君主點頭前要聽幾個成立的理由。
//
// 好戰等級決定他對哪一類指令買不買帳（說明書 3.9）：
// 好戰的君主容易被說服去打人、不容易被說服停戰；消極的相反。
//
// ⚠ 換算是 remake 的暫定值，只保證方向。
func requiredReasons(c Command, aggression int) int {
	n := 2
	switch c {
	case Hostility:
		n += (15 - aggression) / 5 // 好戰的君主要的理由少
	case CeaseFire, Cooperate:
		n += aggression / 5 // 好戰的君主要的理由多
	}
	if n < 1 {
		n = 1
	}
	return n
}

// Begin 開始一次說服。
func Begin(c Command, s Situation) *Session {
	return &Session{Command: c, Situation: s, need: requiredReasons(c, s.Aggression)}
}

// Remaining 回傳君主還想再聽幾個理由。純粹給測試與除錯用——
// **這個數字不該顯示給玩家**，猜君主的脾氣正是本作的玩法。
func (s *Session) Remaining() int { return s.need }

// Offer 提出一個理由。
//
// 說明書 3.9：
//
//	この中から状況に合うものだけを君主が合意するまで選んでいきます。
//	このとき合わないものを選ぶと、説得失敗で信頼度が下がります。
//	また、成功すると上がります。
//	状況に合うものを総て選択しても君主が納得しない場合は、
//	進言撤回でキャンセルする事が出来ます。この場合は信頼度は変化しません。
//
// 回傳結果與信賴度的變化量。
func (s *Session) Offer(r Reason) (Outcome, int) {
	if r == Withdraw {
		return Withdrawn, 0
	}
	// 重複提同一個理由不算數，也不罰——原版沒寫，
	// 但把它當成「選到不符合的」會讓誤點的代價過重。
	if s.used[r] {
		return Continue, 0
	}
	s.used[r] = true

	if !s.Situation.Applies(r) {
		return Failed, TrustOnFailure
	}
	s.need--
	if s.need <= 0 {
		return Agreed, TrustOnSuccess
	}
	return Continue, 0
}

// Exhausted 回報「符合狀況的理由都講完了，君主還是不點頭」。
//
// 這時玩家應該用進言撤回收手（信賴度不變），
// 而不是硬選一個不成立的理由（會扣信賴度）。
func (s *Session) Exhausted() bool {
	for _, r := range pools[s.Command] {
		if !s.used[r] && s.Situation.Applies(r) {
			return false
		}
	}
	return true
}
