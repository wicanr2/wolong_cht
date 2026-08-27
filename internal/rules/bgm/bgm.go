// Package bgm 決定「現在該放哪一首」。
//
// 對應全部出自 `KI.EXE` 的呼叫端（docs/re/58），不是聽出來的。
// 這裡只吃純值：桌面與手機各自把自己的狀態填進 Scene，**規則只有一份**
//（CLAUDE.md §7 第 6 條：反組譯筆記重複只是浪費，程式碼重複會產生行為差異）。
package bgm

import (
	"fmt"

	"github.com/wicanr2/wolong_cht/internal/rules/battlefield"
)

// seasonMusic 是月份 → 曲號。曲 2–5 是四季配樂，跟著大地圖的月份走
//（docs/re/58 §2）。
var seasonMusic = [12]int{5, 5, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5}

// Battle 是戰術畫面的選曲條件。
//
// ⭐ 門檻用的是**戰場編號**不是「攻城／野戰」這個布林值：
// 編號 ≥ 0xD1 是中心格為山地／林地／水域的特殊戰場，原版給它另一首。
type Battle struct {
	// Field 是戰場編號（`battlefield.FieldNumber` 的結果）。
	Field int
	// PlayerAttacks 只在攻城戰用得到：玩家是攻方還是守方。
	// 原版由 `sub_14ED7` 拿玩家勢力比對據點與軍團的持有者定的。
	PlayerAttacks bool
}

// Scene 是選曲要看的全部狀態。
type Scene struct {
	// Launcher 是還停在啟動殼層。
	Launcher bool
	// Ending 是結局過場正在播。
	Ending bool
	// GameOver 是勝負已定。
	GameOver bool
	// Message 是事件訊息或進言對話開著。
	Message bool
	// Battle 非 nil 表示在戰術畫面。
	Battle *Battle
	// Month 是 1–12；**0 表示沒有世界**（還沒開局），這時不放任何曲子。
	Month int
}

// Track 回傳曲名，空字串表示「不換曲」。
//
// ⚠ 三件原版行為，remake 照做：
//   - **戰術進場先停**，設定跑完才依戰場類別挑曲（`sub_19946`）
//   - 換季那個月是**第 2 天**換曲，調色盤卻要漸變 16 天——兩者不同步
//     是原版行為（docs/re/58 §2）。remake 目前只做換曲那一半
//   - 事件與對話放曲 6，結束後回到當季那一首。原版是四支常式各自
//     呼叫曲 6、收尾再呼叫 `sub_19321` 放回當季（docs/re/58 §3）；
//     remake 這一側只要條件不成立就自然落回季節那一支，形狀一樣
func Track(s Scene) string {
	switch {
	case s.Launcher:
		// `sub_11A6E` 開機流程的第一件事就是曲 0（docs/re/58 §3）。
		return "bgm-0"
	case s.Ending:
		// ⭐ **結局這一格要排在很前面。** 排在勝負之後的話，整段結局會放成
		// `overbgm-0`（遊戲結束曲，那是 `D7OVER.EXE` 的）——結局播的時候
		// 勝負早就定了。
		return "endbgm-0"
	case s.Month == 0 && s.Battle == nil:
		return ""
	case s.GameOver:
		// `OVERBGM.DAT` 是 `D7OVER.EXE`（遊戲結束）的（docs/re/58 §6）。
		return "overbgm-0"
	case s.Battle != nil:
		return battleTrack(*s.Battle)
	case s.Message:
		// 曲 6 ＝ 事件與對話。⚠ 原版的四個呼叫端是外交對話、事件 2/3、
		// 事件 4/5 與系統服務分派（docs/re/58 §3），**remake 這一側
		// 不是一對一**：這裡用「事件訊息開著」與「進言對話開著」兩個狀態。
		return "bgm-6"
	default:
		if s.Month < 1 || s.Month > 12 {
			return ""
		}
		return fmt.Sprintf("bgm-%d", seasonMusic[s.Month-1])
	}
}

// battleTrack 依戰場編號與玩家的攻守挑曲（docs/re/58 §4）。
//
// 原版是 `sub_19946` 算的：`byte_1D34B`（戰場編號分三類）＋ 7，
// 而攻城戰那一格再看 `byte_10D35` 的 bit 6。
func battleTrack(b Battle) string {
	switch {
	case b.Field >= 0xD1:
		return "bgm-10" // 山地／林地／水域的戰場
	case b.Field >= battlefield.FieldBase:
		return "bgm-9" // 平原野戰
	case b.PlayerAttacks:
		return "bgm-7" // 攻城戰，玩家是攻方
	default:
		return "bgm-8" // 攻城戰，玩家是守方
	}
}
