// wlsim 是無頭的世界模擬器。
//
// 它把 internal/rules 的三個套件（clock／economy／general）接成一條
// 可以跑很久的迴圈，用**長期行為**去驗證那些從機器碼讀出來的公式。
// 這不是遊戲——沒有畫面、沒有輸入、不會存檔。
//
//	tools/go.sh run ./cmd/wlsim -orig workplace/orig/dosv/SINARIO.DAT
//	tools/go.sh run ./cmd/wlsim -scenario 1 -years 20 -tax 50
//
// 為什麼需要它：規則層每個套件都有單元測試，但單元測試只能驗
// 「一次呼叫的結果對不對」。經濟是複利模型（docs/mechanics/40 §4），
// 錯誤要跑幾十個月才看得出來——這支程式就是那個 loop。
//
// ⚠ 這支程式**不 import Ebiten**。Ebiten 在 init 期就要求顯示器，
// 跟它放同一個 binary 在容器裡跑不起來（cmd/wlshot 踩過同樣的坑）。
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/rules/rng"
	"github.com/wicanr2/wolong_cht/internal/state"
)

func main() {
	path := flag.String("orig", "workplace/orig/dosv/SINARIO.DAT",
		"原版劇本檔（不隨本專案散布，請自備）")
	scenario := flag.Int("scenario", 0, "劇本編號 0–3")
	player := flag.Int("player", 0, "玩家所仕的勢力編號")
	years := flag.Int("years", 10, "要跑幾年")
	tax := flag.Int("tax", -1, "覆寫稅率（-1 ＝ 用劇本預設）")
	seed := flag.Int("seed", 1, "亂數種子（0–255，原版是從即時時鐘播的）")
	check := flag.Bool("check", false,
		"每個 tick 檢查資料模型的不變量，第一個違反就停（見 internal/state/invariant.go）")
	every := flag.Int("every", 12, "每幾個月印一列")
	flag.Parse()

	w, err := state.LoadScenario(*path, *scenario)
	if err != nil {
		log.Fatal(err)
	}
	w.Player = *player
	if *tax >= 0 {
		w.TaxRate, w.NextTaxRate = *tax, *tax
	}

	fmt.Printf("劇本 %d　起始 %d年%d月%d日　勢力 %d 個　稅率 %d%%　種子 %d\n",
		*scenario+1, w.Clock.Year, w.Clock.Month, w.Clock.Day,
		len(w.AliveFactions()), w.TaxRate, *seed)
	fmt.Printf("玩家所仕勢力 %d（君主 %s）\n\n", w.Player, big5(w.LordName(w.Player)))

	rng := rng.NewFixed(*seed)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "年月\t勢力數\t玩家據點\t玩家資金\t玩家預備兵\t平均生產力\t玩家上昇值\tAI上昇值\tAI低於−36\t火災\t暴動\t暴風雨")

	var fires, riots, storms int
	months := 0
	total := *years * 12

	ticks := 0
	for months < total {
		ev := w.Tick(rng)
		ticks++
		// ⭐ 不變量檢查：驗的是「規則組合起來對不對」，
		// 不是單條公式對不對（單元測試已經釘住單條了）。
		if *check {
			if v := w.CheckInvariants(); len(v) > 0 {
				tw.Flush()
				fmt.Printf("\n✗ 第 %d 個 tick（%d 年 %d 月 %d 日）違反 %d 條：\n",
					ticks, w.Clock.Year, w.Clock.Month, w.Clock.Day, len(v))
				for i, x := range v {
					if i >= 5 {
						fmt.Printf("  …還有 %d 條\n", len(v)-5)
						break
					}
					fmt.Printf("  %s\n", x)
				}
				os.Exit(1)
			}
		}
		if !ev.Settled {
			continue
		}
		months++
		for _, d := range ev.Disaster {
			switch d {
			case economy.Fire:
				fires++
			case economy.Riot:
				riots++
			}
		}
		if ev.Storm != nil {
			storms++
		}
		for _, i := range ev.Eliminated {
			fmt.Fprintf(tw, "%d/%d\t— 勢力 %d（%s）滅亡 —\n",
				w.Clock.Year, w.Clock.Month, i, big5(w.LordName(i)))
		}
		if months%*every != 0 {
			continue
		}

		p := w.Factions[w.Player]
		prodSum, growSum, owned := 0, 0, 0
		// ⭐ AI 的據點要**分開統計**。先前只印玩家的，於是實作內政官之後
		// 看到「平均上昇值 +100」就以為問題解決了，但暴動只降了 64%
		// 而且速率沒隨時間下降——那代表沒被拉起來的是別人的城。
		// **看不到的那一半才是問題所在。**
		aiSum, aiN, aiRiskly := 0, 0, 0
		for _, c := range w.Cities {
			switch {
			case c.Owner == w.Player:
				prodSum += c.Production
				growSum += c.Growth
				owned++
			case c.Owner >= 0 && c.Owner < 22:
				aiSum += c.Growth
				aiN++
				// 暴動的免疫門檻：上昇值存值 ≥ 64 ＝ 實際值 ≥ −36
				// （docs/re/07 §17、internal/rules/economy/disaster.go）。
				if c.Growth < -36 {
					aiRiskly++
				}
			}
		}
		avgP, avgG := 0, 0
		if owned > 0 {
			avgP, avgG = prodSum/owned, growSum/owned
		}
		avgAI := 0
		if aiN > 0 {
			avgAI = aiSum / aiN
		}
		fmt.Fprintf(tw, "%d/%d\t%d\t%d\t%d\t%d/%d/%d\t%d\t%+d\t%+d\t%d\t%d\t%d\t%d\n",
			w.Clock.Year, w.Clock.Month, len(w.AliveFactions()),
			p.Cities, p.Funds,
			p.Reserves[economy.Cavalry], p.Reserves[economy.Archer], p.Reserves[economy.Infantry],
			avgP, avgG, avgAI, aiRiskly, fires, riots, storms)
	}
	tw.Flush()

	fmt.Printf("\n跑了 %d 個月（%d tick）。火災 %d 次、暴動 %d 次、暴風雨 %d 次。\n",
		months, months*30*clock.TicksPerDay, fires, riots, storms)
	p := w.Factions[w.Player]
	fmt.Printf("玩家勢力：據點 %d　資金 %d　君主好戰等級 %d　侵攻可持續 %v\n",
		p.Cities, p.Funds, p.Aggression,
		diplomacy.CanSustainInvasion(p.Funds, p.Cities))
}

// big5 把 internal/state 保留的原始位元組轉成可以印的字串。
//
// 規則層刻意不做編碼轉換（存原始 byte 才能 round-trip 回原版檔案），
// 所以轉換發生在最外層。
func big5(s string) string {
	if s == "" {
		return "?"
	}
	return text.Decode([]byte(s), text.Big5)
}
