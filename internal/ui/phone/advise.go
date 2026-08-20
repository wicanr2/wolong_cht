package phone

import (
	"github.com/wicanr2/wolong_cht/internal/rules/persuasion"
	"github.com/wicanr2/wolong_cht/internal/ui/talkmenu"
)

// 進言：玩家是軍師不是君主，戰略指令是**向君主提議**，君主可能拒絕
//（CLAUDE.md §3.2）。判定全在 `internal/rules/persuasion`，
// 對白的 TALK 索引算式也在那裡——這一層只做「點哪裡、顯示什麼」。

// adviseStage 是進言流程的階段。與桌面版同一組概念，但少了兩個框的排版。
type adviseStage int

const (
	adviseIdle       adviseStage = iota
	advisePickAlly               // 請求協助：先選協力方
	advisePickTarget             // 選交涉對象
	advisePersuade               // 君主要理由，進說服迴圈
	adviseVerdict                // 君主已定案，只等關掉
)

// adviseCommands 是前三項對應的說服指令。第四、五項不走說服迴圈。
var adviseCommands = []persuasion.Command{
	persuasion.Hostility, persuasion.CeaseFire, persuasion.Cooperate,
}

const (
	adviseRelocateRow = 3 // 遷　都
	adviseSortieRow   = 4 // 請求君主出陣
)

// adviseFallbackNames 只在讀不到 `TALK.DAT` 時用。正常路徑走 #77。
var adviseFallbackNames = []string{
	"敵對提案", "停戰提案", "請求協助", "遷　都", "請求君主出陣",
}

// advise 是進言流程的狀態。
type advise struct {
	stage adviseStage
	cmd   persuasion.Command
	row   int
	ally  int
	tgt   int

	sess *persuasion.Session

	// said 是畫面上要顯示的對白，最新的在最後。
	// **兩個角色合成一串**：手機的版面放不下原版的上下兩個框，
	// 改成一列一句並標明是誰講的（remake 差異，docs/mobile/android-ux.md §7）。
	said []adviseLine
}

type adviseLine struct {
	lord bool
	text string
}

// adviseLabels 是進言的五項，取自 `TALK.DAT` #77。
func (s *Session) adviseLabels() []string {
	if s.lib == nil {
		return append([]string(nil), adviseFallbackNames...)
	}
	return talkmenu.Labels(s.lib.Talk, persuasion.TalkMenu, nil, adviseFallbackNames)
}

// adviseHint 是每一項右邊那一欄：**現在點下去會發生什麼**。
//
// ⭐ 遷都要先在地圖上點好目的地。手機沒有「先開選單再選城」的餘裕
//（選單一開就蓋住地圖），所以順序反過來：先選城、再開進言。
func (s *Session) adviseHint(row int) string {
	switch row {
	case adviseRelocateRow:
		if s.selected < 0 {
			return "先在地圖上點目的地"
		}
		return "遷往 " + big5(s.world.Cities[s.selected].Name)
	case adviseSortieRow:
		return "君主親自出陣"
	default:
		return "選對象勢力"
	}
}

// AdviseStage 回報進言流程走到哪。給驗收與輸入層用。
func (s *Session) AdviseStage() adviseStage { return s.advise.stage }

// AdviseLines 是目前畫面上的對白。
func (s *Session) AdviseLines() []adviseLine { return s.advise.said }

// AdviseChoices 是這一刻要玩家選的那幾列，沒有選項時回 nil。
func (s *Session) AdviseChoices() []string {
	switch s.advise.stage {
	case advisePickAlly, advisePickTarget:
		rows := make([]string, 0, len(s.world.Factions))
		for _, i := range s.adviseFactions() {
			rows = append(rows, big5(s.world.LordName(i)))
		}
		return rows
	case advisePersuade:
		if s.advise.sess == nil {
			return nil
		}
		return s.reasonLabels(s.advise.cmd)
	}
	return nil
}

// adviseFactions 是可以當交涉對象的勢力編號，**不含自己**。
func (s *Session) adviseFactions() []int {
	out := make([]int, 0, len(s.world.Factions))
	for i := range s.world.Factions {
		if s.world.Factions[i].Alive && i != s.world.Player {
			out = append(out, i)
		}
	}
	return out
}

func (s *Session) reasonLabels(c persuasion.Command) []string {
	opts := persuasion.Options(c)
	fallback := make([]string, len(opts))
	for i, r := range opts {
		fallback[i] = r.String()
	}
	if s.lib == nil {
		return fallback
	}
	return talkmenu.Labels(s.lib.Talk, persuasion.TalkReasonBase(c),
		s.adviseVars(), fallback)
}

// adviseVars 是進言那幾則的變數：`{3}` 交涉對象、`{4}` 軍師（玩家）、
// `{6}` 排版標記（空字串，原版 handler 只調 X 不輸出字元）。
func (s *Session) adviseVars() map[byte]string {
	vars := map[byte]string{'6': ""}
	if t := s.advise.tgt; t >= 0 && t < len(s.world.Factions) {
		vars['3'] = big5(s.world.LordName(t))
	}
	if p := s.world.Player; p >= 0 && p < len(s.world.Factions) {
		if a := s.world.Factions[p].Advisor; a >= 0 && a < len(s.world.Generals) {
			vars['4'] = big5(s.world.Generals[a].Name)
		}
	}
	return vars
}

// PickAdvise 點了進言的第 row 項。
func (s *Session) PickAdvise(row int) {
	s.advise = advise{row: row, ally: -1, tgt: -1}
	switch row {
	case adviseRelocateRow:
		if s.selected < 0 {
			return // 還沒選目的地，什麼都不做（提示已經印在那一列上）
		}
		ok := s.world.AdviseRelocateAccepted(s.selected)
		s.sayVerdict(persuasion.TalkRelocateBase, ok)
		if ok {
			s.world.AdviseRelocate(s.selected)
		}
	case adviseSortieRow:
		ok := s.world.AdviseSortieAccepted()
		s.sayVerdict(persuasion.TalkSortieBase, ok)
		if ok {
			s.world.AdviseSortie()
		}
	default:
		if row < 0 || row >= len(adviseCommands) {
			return
		}
		s.advise.cmd = adviseCommands[row]
		if s.advise.cmd == persuasion.Cooperate {
			s.advise.stage = advisePickAlly
			return
		}
		s.advise.stage = advisePickTarget
	}
}

// PickAdviseChoice 點了選項清單的第 i 列。
func (s *Session) PickAdviseChoice(i int) {
	switch s.advise.stage {
	case advisePickAlly:
		ids := s.adviseFactions()
		if i < 0 || i >= len(ids) {
			return
		}
		s.advise.ally = ids[i]
		s.advise.stage = advisePickTarget
	case advisePickTarget:
		ids := s.adviseFactions()
		if i < 0 || i >= len(ids) {
			return
		}
		s.advise.tgt = ids[i]
		s.beginPersuasion()
	case advisePersuade:
		opts := persuasion.Options(s.advise.cmd)
		if i < 0 || i >= len(opts) {
			return
		}
		s.offerReason(opts[i])
	}
}

// CloseAdvise 關掉進言流程。
func (s *Session) CloseAdvise() { s.advise = advise{ally: -1, tgt: -1} }

// sayVerdict 演 `sub_13B08` 的三句：君主開場、軍師、君主定案。
//
// ⚠ **看到君主說話不代表提議被接受**——原版無論通不通過都跳同一組，
// 差別只在第三句（docs/mechanics/70-ai.md）。
func (s *Session) sayVerdict(base int, accepted bool) {
	s.advise.stage = adviseVerdict
	s.advise.said = nil
	v := s.talkVariant()
	s.say(true, base+v)
	s.say(false, base+3)
	reply := base + 4
	if !accepted {
		reply += 3
	}
	s.say(true, reply+v)
}

func (s *Session) beginPersuasion() {
	sit := s.world.PersuasionSituation(s.advise.cmd, s.advise.tgt, s.advise.ally)
	base := persuasion.TalkBase(s.advise.cmd)
	v := s.talkVariant()
	s.advise.stage = advisePersuade
	s.advise.said = nil
	s.say(true, base+v)
	s.say(false, base+3)

	queued := false
	if s.advise.cmd == persuasion.Hostility {
		queued = s.world.HasQueuedDeclaration(s.world.Player, s.advise.tgt)
	}
	switch reaction := persuasion.FirstReaction(s.advise.cmd, sit, queued); reaction {
	case persuasion.AskReason:
		s.advise.sess = persuasion.Begin(s.advise.cmd, sit)
	default:
		s.advise.sess = nil
		s.world.AdjustTrust(persuasion.ReactionTrustDelta(reaction))
		s.say(true, persuasion.TalkReplyIndex(base, reaction, v))
		if reaction == persuasion.Agree {
			s.commitAdvice()
		}
		s.advise.stage = adviseVerdict
	}
}

func (s *Session) offerReason(r persuasion.Reason) {
	if s.advise.sess == nil {
		return
	}
	repeat := s.advise.sess.Offered(r)
	out, dt := s.advise.sess.Offer(r)
	s.world.AdjustTrust(dt)
	base := persuasion.TalkReasonBase(s.advise.cmd)
	slot := persuasion.TalkReasonSlot(s.advise.cmd, r)
	v := s.talkVariant()

	s.say(false, base+slot+1)
	s.say(true, persuasion.TalkReasonReply(base, slot, out, repeat, v))

	switch out {
	case persuasion.Agreed:
		s.commitAdvice()
		s.advise.stage = adviseVerdict
		s.advise.sess = nil
	case persuasion.Failed, persuasion.Withdrawn:
		s.advise.stage = adviseVerdict
		s.advise.sess = nil
	}
}

// commitAdvice 把「君主同意了」接到規則層。
func (s *Session) commitAdvice() bool {
	var ok bool
	switch s.advise.cmd {
	case persuasion.Hostility:
		ok = s.world.ApplyPlayerHostility(s.advise.tgt)
	case persuasion.CeaseFire:
		ok = s.world.QueuePlayerCeasefire(s.advise.tgt)
	case persuasion.Cooperate:
		ok = s.world.QueuePlayerCooperation(s.advise.ally, s.advise.tgt)
	}
	if !ok {
		// remake 專屬的守門句：君主點頭之後規則層才發現條件變了。
		// 原版沒有這條路徑，所以沒有對應的原文（docs/spec/44 §6）。
		s.advise.said = append(s.advise.said,
			adviseLine{lord: true, text: "局勢已變，這項進言沒有成立。"})
	}
	return ok
}

// say 把一則 TALK 加進對白。讀不到就不加——**寧可少一句，
// 也不要把索引當文字印出去**。
func (s *Session) say(lord bool, index int) {
	if s.lib == nil {
		return
	}
	lines, ok := s.lib.Talk.Lines(index, s.adviseVars())
	if !ok {
		return
	}
	for _, l := range lines {
		if l == "" {
			continue
		}
		s.advise.said = append(s.advise.said, adviseLine{lord: lord, text: l})
	}
}

// talkVariant 是君主的**說話型**（`sub_13C99` 會把它加進索引）。
func (s *Session) talkVariant() int {
	p := s.world.Player
	if p < 0 || p >= len(s.world.Factions) {
		return 0
	}
	lord := s.world.Factions[p].Lord
	if lord < 0 || lord >= len(s.world.Generals) {
		return 0
	}
	return s.world.Generals[lord].TalkVariant
}
