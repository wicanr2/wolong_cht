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
	// **用字照原版選單**（訊息 102／166／230），不是自己翻的。
	// 這是保存專案，選單文字屬於松崗版的原文。
	return [...]string{
		"外交關係惡劣", "外交關係良好", "我國較有利", "對我國較不利",
		"敵正侵攻他國", "我正在防禦戰", "敵勢力疲乏", "我國力疲乏",
		"協力國強大", "侵攻對象強大", "撤回進言",
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
// 說明書 3.9 只給了方向（呂布關係不特別差就買帳、劉備要相當差才買帳），
// **實際的式子從 `sub_16475` 讀出來了**：`ah>>1 + 5`，
// 而 `ah` 在那一行是 `好戰 × 2 + 20`，所以是 **好戰 ＋ 15**。
//
// ⚠ 舊值是 `10 + 好戰 × 2`（remake 暫定）。形狀猜對了一半——
// `好戰 × 2` 那個項確實存在，但它是**君主拒絕**的門檻
// （見 FirstReaction），不是這個理由的門檻。
// **一個猜測在錯的地方對，比全錯更難發現。**
func badFriendshipGate(aggression int) int { return aggression + 15 }

// weAreStronger 是「我國較有利」的判準（`sub_16A28`）。
//
//	我方據點數 × (好戰 ＋ 20)  >  敵方據點數 × 25
//
// **平衡點在好戰 ＝ 5**：據點數相同時，好戰 5 的君主覺得是平手，
// 6 以上才覺得自己有利。劉禪 0、劉表 1、劉備 4 都在平衡點以下。
//
// ⚠ 舊版寫成「我方據點 ＋ 好戰/3 > 敵方據點」——**加法**。
// 原版是兩邊各乘一個係數的**乘法**，差別在據點數大的時候會拉開：
// 10 vs 10 與 100 vs 100 在加法版是同一個答案，在原版不是。
func weAreStronger(ourCities, theirCities, aggression int) bool {
	return ourCities*(aggression+20) > theirCities*25
}

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
		return weAreStronger(s.OurCities, s.TheirCities, s.Aggression)
	case EnemyIsStronger:
		// 反向用同一個式子取反——不是另外定義一條，
		// 免得兩邊在邊界上同時成立或同時不成立。
		return !weAreStronger(s.OurCities, s.TheirCities, s.Aggression)
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
// 三個指令各自的理由池。**內容與順序直接抄自原版的選單訊息**
// （102／166／230），不是照說明書的分類推的。
//
//	#102 敵對提案  外交關係惡劣│我國較有利　│敵正侵攻他國│敵勢力疲乏　│撤回進言
//	#166 停戰提案  對我國較不利│我正在防禦戰│敵正侵攻他國│我國力疲乏　│撤回進言
//	#230 協力要請  外交關係良好│協力國強大　│侵攻對象強大│我正在防禦戰│撤回進言
//
// ⚠ **每個池是四個理由，不是五個。** 說明書那句「常に 5 つの項目が
// 選択肢として用意されています」的 5 是**含撤回**——先前讀成
// 「5 個理由 ＋ 撤回」，於是三個池各自多塞或錯放了一個。
// 對照之下，舊版三個池只有敵對提案接近正確：
//
//	敵對  多了「我國防戰中」
//	停戰  多了「外交關係良好」「侵攻對象強大」，少了「敵正侵攻他國」
//	協力  多了「我國力疲乏」「敵正侵攻他國」，少了「我正在防禦戰」
//
// **順序也是資料**：`sub_16475` 組出來的可用旗標，位元順序與選單一致。
var pools = map[Command][]Reason{
	Hostility: {FriendshipBad, WeAreStronger, EnemyInvading, EnemyExhausted},
	CeaseFire: {EnemyIsStronger, WeAreDefending, EnemyInvading, WeAreExhausted},
	Cooperate: {FriendshipGood, AllyIsStronger, InvaderIsStronger, WeAreDefending},
}

// Options 回傳這個指令的選單：四個理由 ＋「撤回進言」，一共五項。
//
// 說明書：「常に 5 つの項目が選択肢として用意されています」——
// **那個 5 是含撤回的**（原版選單訊息 102／166／230 各正好五行）。
//
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

// Reaction 是君主聽完進言的**第一反應**（`sub_16475`）。
//
// 說明書只講了「君主可能拒絕，然後挑理由說服他」，沒說在那之前
// 還有三種不進入說服迴圈的分支。數值照原版——
// 台詞是按 `基底 ＋ 4 ＋ 反應碼 × 3` 排的（`sub_13830`），
// 所以這幾個數字有意義，不能重排。
type Reaction int

const (
	// Refuse ＝ 0：「無法答允！別平白增加敵人。」（訊息 90–92）
	Refuse Reaction = 0
	// Agree ＝ 1：「我也有同樣的想法。立刻準備交戰！」（93–95）
	Agree Reaction = 1
	// AskReason ＝ 2：「聽你這麼說，看來是有勝算囉？」（96–98）
	// **只有這一個會進入說服迴圈**（Begin／Offer）。
	AskReason Reaction = 2
	// AlreadyAtWar ＝ 3：「你別迷糊了！不是已經在交戰狀態中了嗎！」（99–101）
	AlreadyAtWar Reaction = 3
)

// FirstReaction 回傳君主的第一反應。
//
// queued 是「事件佇列裡已經有一筆本勢力要打這個目標的宣戰事件」
// （原版 `sub_1304E` 掃 `0x000`–`0x3FF`）。為真代表**自家勢力的 AI
// 早就決定要打它了**——台詞「我也有同樣的想法」講的正是這件事，
// 不是客套話，是同一個變數的兩端。
//
// atWar 是「和平位元沒設」，也就是已經在交戰中。
//
// ⚠ **拒絕的門檻是 `好戰 × 2 ＋ 20`**，與 badFriendshipGate 的
// `好戰 ＋ 15` 是兩條不同的線。舊版把 `10 + 好戰×2` 用在後者，
// 形狀對、常數與位置都錯。
func FirstReaction(s Situation, queued, atWar bool) Reaction {
	switch {
	case queued:
		return Agree
	case atWar:
		return AlreadyAtWar
	case s.Friendship >= s.Aggression*2+20:
		return Refuse
	}
	return AskReason
}
