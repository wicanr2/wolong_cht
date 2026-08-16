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
	// Trust 是原版全域信賴度 byte_10D00（0–255）。說服迴圈開始時，
	// 原版 sub_13C1E 會把它換成 1–4 級，決定還要聽幾個成立理由。
	Trust int

	// Aggression 是自家君主的好戰等級（0–15，勢力記錄 +0x28）。
	// **不顯示給玩家**——玩家只能靠被拒絕的次數去推。
	Aggression int

	// 「對方」在三個指令裡指的不是同一個勢力：
	//
	//	敵對提案  想打的目標
	//	停戰提案  想停戰的敵人
	//	請求協助  **想一起打的侵攻對象**（協力對象另外放 Ally*）
	OurCities, TheirCities int // 國力比較的據點數
	AllyCities             int // 協力對象的據點數（只有請求協助用得到）

	OurFunds, TheirFunds int // 疲弊 ＝ 資金 < 0

	// Friendship 是自家君主看「對方」的交友值，**含最高位的和平位元**
	// （`0x80` ＝ 和平）。判定式直接拿原始值比，不要先去掉那個位元——
	// 原版的門檻常數本身就帶著它。
	Friendship int
	// AllyFriendship 是看協力對象的交友值（只有請求協助用得到）。
	AllyFriendship int

	TheyInvadeThirdParty bool // 對方正在侵攻第三方
	TheyInvadeUs         bool // **對方**正在侵攻我方
	// SameFactionPicked 為真表示請求協助時，協力對象與侵攻對象選成同一家。
	SameFactionPicked bool
	// AnyoneInvadesUs 為真表示**有任何一個別的勢力**（交涉對象除外）
	// 把我方設成侵攻目標。停戰提案的「我正在防禦戰」用這個，
	// 而請求協助用的是 TheyInvadeUs ——**同一個選項，兩個條件**。
	AnyoneInvadesUs bool
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
func badFriendshipGate(aggression int) int { return peaceBit + aggression + 15 }

// goodFriendshipGate 是「交友關係良好」成立的下限（`sub_166D9`）：
//
//	交友度 ≥ 0x80 ＋ 好戰 × 4 ＋ 60
//
// ⚠ **它不是 badFriendshipGate 的反面。** 兩個門檻中間有一大段
// 「既不算差、也不算好」——好戰 5 時是 `0xA0`–`0xD0`（32–80）。
// 舊版把兩者寫成互補，於是那一段的判定必定有一邊是錯的。
func goodFriendshipGate(aggression int) int { return peaceBit + aggression*4 + 60 }

// peaceBit 是交友度最高位：設著 ＝ 和平。
const peaceBit = 0x80

// power 是 `sub_16A28` 的國力比較：兩邊各乘一個係數。
//
//	我方 ＝ 我方據點數 × (好戰 ＋ 20)
//	對方 ＝ 對方據點數 × 25
//
// **平衡點在好戰 ＝ 5**：據點數相同時，好戰 5 的君主覺得是平手。
// 劉禪 0、劉表 1、劉備 4 都在平衡點以下。
//
// ⚠ 三個指令拿它比的方向不同，而且**相等時兩邊都不成立**：
//
//	敵對「我國較有利」    我方 >  對方
//	停戰「對我國較不利」  我方 <  對方
//	協力「協力國強大」    我方 <  協力對象
//
// 舊版把後者寫成 `!weAreStronger`，於是相等時會誤判成「對我國較不利」。
func power(ourCities, theirCities, aggression int) (ours, theirs int) {
	return ourCities * (aggression + 20), theirCities * 25
}

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
	ours, theirs := power(ourCities, theirCities, aggression)
	return ours > theirs
}

// theyAreStronger 是反向。**不是 weAreStronger 取反**——相等時兩者皆偽。
func theyAreStronger(ourCities, theirCities, aggression int) bool {
	ours, theirs := power(ourCities, theirCities, aggression)
	return ours < theirs
}

// Applies 回報某個理由在這個局勢下是否「符合狀況」。
//
// **要傳指令**：同一個選項在不同指令下是**不同的條件**。
// 最明顯的是「我正在防禦戰」——停戰提案掃全部 22 個勢力
// （`sub_16577` 的 `mov cx, 16h` 迴圈），請求協助只看侵攻對象
// （`sub_166D9` 的 `cmp al, byte_10CFF`）。
//
// **選到不符合的理由 → 說服失敗、信賴度下降。**
// 所以這個函式的正確性直接決定玩家會不會被冤枉扣分。
func (s Situation) Applies(c Command, r Reason) bool {
	switch r {
	case FriendshipBad: // 敵對①
		return s.Friendship < badFriendshipGate(s.Aggression)
	case FriendshipGood: // 協力①，比的是**協力對象**
		return s.AllyFriendship >= goodFriendshipGate(s.Aggression)
	case WeAreStronger: // 敵對②
		return weAreStronger(s.OurCities, s.TheirCities, s.Aggression)
	case EnemyIsStronger: // 停戰①「對我國較不利」
		return theyAreStronger(s.OurCities, s.TheirCities, s.Aggression)
	case AllyIsStronger: // 協力②，比的是**協力對象**
		return theyAreStronger(s.OurCities, s.AllyCities, s.Aggression)
	case InvaderIsStronger: // 協力③「侵攻對象強大」
		return theyAreStronger(s.OurCities, s.TheirCities, s.Aggression)
	case EnemyInvading: // 敵對③／停戰③
		return s.TheyInvadeThirdParty
	case WeAreDefending:
		// ⚠ 同一個選項，兩個條件。
		if c == CeaseFire {
			return s.AnyoneInvadesUs
		}
		return s.TheyInvadeUs
	case EnemyExhausted: // 敵對④，**對方**資金 < 0
		return s.TheirFunds < 0
	case WeAreExhausted: // 停戰④，**我方**資金 < 0
		return s.OurFunds < 0
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

// 信賴度的增減量。數值來自 sub_13830：
// `mov al,14h` 後分別呼叫 sub_13D91／sub_13DC9；多理由成功先
// `shr al,1`，所以是 +10。原版的 byte_10D00 會在 0 與 255 飽和。
const (
	TrustOnReasonSuccess    = 10
	TrustOnImmediateSuccess = 20
	TrustOnFailure          = -20
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
// sub_13C1E 直接讀 byte_10D00，回傳放在 AH 的信賴度級別；
// sub_13BA9 每選到一個符合狀況的理由就遞減 [bp+3]（也就是該級別）。
// 這個門檻與 Command、Aggression 無關；好戰等級只參與三種指令的
// 第一反應與各理由是否成立。
func requiredReasons(trust int) int {
	switch {
	case trust >= 0xE0:
		return 1
	case trust >= 0x90:
		return 2
	case trust >= 0x20:
		return 3
	default:
		return 4
	}
}

// Begin 開始一次說服。
func Begin(c Command, s Situation) *Session {
	return &Session{Command: c, Situation: s, need: requiredReasons(s.Trust)}
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

	if !s.Situation.Applies(s.Command, r) {
		return Failed, TrustOnFailure
	}
	s.need--
	if s.need <= 0 {
		return Agreed, TrustOnReasonSuccess
	}
	return Continue, 0
}

// Offered 回報這個理由這一輪講過了沒有。
//
// 原版對「同一個理由再講一次」有專屬台詞（TALK `base+45`，
// docs/spec/44 §5），呈現層要靠這個分辨。
func (s *Session) Offered(r Reason) bool {
	if s == nil || r < 0 || r >= numReasons {
		return false
	}
	return s.used[r]
}

// Exhausted 回報「符合狀況的理由都講完了，君主還是不點頭」。
//
// 這時玩家應該用進言撤回收手（信賴度不變），
// 而不是硬選一個不成立的理由（會扣信賴度）。
func (s *Session) Exhausted() bool {
	for _, r := range pools[s.Command] {
		if !s.used[r] && s.Situation.Applies(s.Command, r) {
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

// SameFaction ＝ 4：請求協助時協力對象與侵攻對象選成同一個。
// `sub_13830` 對 ≥ 4 的碼一律顯示訊息 83「我想軍師並不是來談笑的。」
const SameFaction Reaction = 4

// ReactionTrustDelta 是 sub_13830 對第一反應碼的信賴度變化。
//
// AL=1 直接成功；AL=0、3 走失敗分支；AL=2 進入 Session，AL=4
// 是協力選到同一家，只顯示訊息而不改信賴度。
func ReactionTrustDelta(r Reaction) int {
	switch r {
	case Agree:
		return TrustOnImmediateSuccess
	case Refuse, AlreadyAtWar:
		return TrustOnFailure
	default:
		return 0
	}
}

// FirstReaction 回傳君主的第一反應。**三個指令的判定式各自不同。**
//
//	敵對 `sub_16475`   佇列裡已有 → 1；交戰中 → 3；交友度 ≥ 好戰×2+20 → 0
//	停戰 `sub_16577`   **和平中 → 3**（沒在打，停什麼）；交友度 < 好戰÷2 → 0
//	協力 `sub_166D9`   兩個選同一家 → 4；交友對象太差 → 0；
//	                   與侵攻對象和平中 → 3；被它打且國力不到它一半 → 1
//
// 好戰等級在拒絕門檻上出現**三個不同的係數**（×2+20／÷2／×4+30）。
// 看起來像可以合併，**不要合併**——那是三個獨立讀出來的常數。
//
// queued 只有敵對用得到：事件佇列裡已經有一筆本勢力要打這個目標的宣戰事件
// （原版 `sub_1304E`），代表自家 AI 早就想打它了，所以君主說
// 「我也有同樣的想法」。
func FirstReaction(c Command, s Situation, queued bool) Reaction {
	atWarWithThem := s.Friendship&peaceBit == 0
	// 去掉和平位元才是 0–127 的實值。原版一律先 `sub al, 80h`。
	friend := s.Friendship &^ peaceBit

	switch c {
	case Hostility:
		switch {
		case queued:
			return Agree
		case atWarWithThem:
			return AlreadyAtWar
		case friend >= s.Aggression*2+20:
			return Refuse
		}
	case CeaseFire:
		switch {
		case !atWarWithThem:
			return AlreadyAtWar // 「原本就沒有和\3交戰啊！」
		case friend < s.Aggression/2:
			return Refuse
		}
	case Cooperate:
		switch {
		case s.SameFactionPicked:
			return SameFaction
		case s.AllyFriendship < peaceBit+s.Aggression*4+30:
			return Refuse
		case !atWarWithThem:
			return AlreadyAtWar // 與侵攻對象還在和平
		case s.TheyInvadeUs && weakerThanHalf(s):
			return Agree
		}
	}
	return AskReason
}

// weakerThanHalf 是協力那條「直接同意」的門檻：我方國力**不到對方的一半**
// （`sub_16A28` 之後 `shr cx, 1`）。被打得夠慘，君主不必聽理由就答應。
func weakerThanHalf(s Situation) bool {
	ours, theirs := power(s.OurCities, s.TheirCities, s.Aggression)
	return ours < theirs/2
}
