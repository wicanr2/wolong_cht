// Package scenario 把劇本／存檔區塊在 `state.World` 與 JSON 之間互轉。
//
// 這一層存在的理由：編輯器需要一個人看得懂、diff 得出來、又能無損寫回的
// 中間格式（docs/spec/28）。
//
// ⭐ **JSON 不是完整劇本，是已解欄位的投影。** 匯入一定要有原始檔：
// 寫回沿用 `state` 的改寫策略（從原始 bytes 出發，只蓋已解欄位），
// 未解區域一個 byte 都不動。拿 JSON 從零重建會產生一個看起來對、
// 實際少東西的檔案。
package scenario

import (
	"fmt"

	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/wicanr2/wolong_cht/internal/assets/text"
	"github.com/wicanr2/wolong_cht/internal/rules/clock"
	"github.com/wicanr2/wolong_cht/internal/rules/economy"
	"github.com/wicanr2/wolong_cht/internal/state"
)

// nameBytes 是姓名欄的長度：6 bytes ＝ 3 個全形字（docs/formats/08）。
const nameBytes = 6

// Meta 記來源，讓一份 JSON 自己說得出它是從哪裡出來的。
type Meta struct {
	Title     string // UTF-8 的劇本標題（區塊 +0x40）
	Source    string // 來源檔名
	Block     int    // 區塊編號 0–3
	SourceSHA string // 來源檔的 SHA-256
}

// City 是 state.City 的 JSON 版：**只有名字換成 UTF-8**，其餘欄位原樣提升。
//
// 淺層的 Name 會遮蔽內嵌的 state.City.Name，所以編碼出來只有一個 Name。
type City struct {
	state.City
	Name string
}

// General 同理。
type General struct {
	state.General
	Name  string
	Alias string
}

// Scenario 是一個區塊的 JSON 形態。
type Scenario struct {
	Meta     Meta
	Clock    clock.Clock
	Player   int
	Trust    int
	TaxRate  int
	NextTax  int
	Recruit  [economy.NumTroopTypes]int
	NextRec  [economy.NumTroopTypes]int
	Factions []state.Faction
	Cities   []City
	Generals []General
}

// FromWorld 把一個 World 投影成 JSON 結構。
func FromWorld(w *state.World, meta Meta) *Scenario {
	meta.Title = decode(w.Title)
	s := &Scenario{
		Meta:    meta,
		Clock:   w.Clock,
		Player:  w.Player,
		Trust:   w.Trust,
		TaxRate: w.TaxRate,
		NextTax: w.NextTaxRate,
		Recruit: w.RecruitCap,
		NextRec: w.NextRecruitCap,
	}
	s.Factions = append(s.Factions, w.Factions[:]...)
	for _, c := range w.Cities {
		s.Cities = append(s.Cities, City{City: c, Name: decode(c.Name)})
	}
	for _, g := range w.Generals {
		s.Generals = append(s.Generals, General{
			General: g, Name: decode(g.Name), Alias: decode(g.Alias)})
	}
	return s
}

// ApplyTo 把 JSON 的內容套回 World。**World 要是從原始檔載入的**——
// 未解區域靠它保留。
func (s *Scenario) ApplyTo(w *state.World) error {
	if n := len(s.Factions); n != len(w.Factions) {
		return fmt.Errorf("scenario: 勢力 %d 筆，預期 %d", n, len(w.Factions))
	}
	if n := len(s.Cities); n != len(w.Cities) {
		return fmt.Errorf("scenario: 據點 %d 筆，預期 %d", n, len(w.Cities))
	}
	if n := len(s.Generals); n != len(w.Generals) {
		return fmt.Errorf("scenario: 武將 %d 筆，預期 %d", n, len(w.Generals))
	}
	w.Clock = s.Clock
	w.Player = s.Player
	w.Trust = s.Trust
	w.TaxRate = s.TaxRate
	w.NextTaxRate = s.NextTax
	w.RecruitCap = s.Recruit
	w.NextRecruitCap = s.NextRec
	copy(w.Factions[:], s.Factions)

	for i, c := range s.Cities {
		name, err := keepOrEncode(w.Cities[i].Name, c.Name)
		if err != nil {
			return fmt.Errorf("scenario: 據點 %d 的名字：%w", i, err)
		}
		w.Cities[i] = c.City
		w.Cities[i].Name = name
	}
	for i, g := range s.Generals {
		name, err := keepOrEncode(w.Generals[i].Name, g.Name)
		if err != nil {
			return fmt.Errorf("scenario: 武將 %d 的名字：%w", i, err)
		}
		alias, err := keepOrEncode(w.Generals[i].Alias, g.Alias)
		if err != nil {
			return fmt.Errorf("scenario: 武將 %d 的別號：%w", i, err)
		}
		w.Generals[i] = g.General
		w.Generals[i].Name = name
		w.Generals[i].Alias = alias
	}
	return nil
}

func decode(raw string) string { return text.Decode([]byte(raw), text.Big5) }

// keepOrEncode 決定姓名欄要寫什麼。
//
// ⭐ **名字沒被改動就原樣保留**，不重新編碼一次。理由是 round-trip：
// `text.Decode` 會把尾端的全形空白與 NUL 修掉，再編碼回去補的是全形空白，
// 而原始檔那兩個 byte 可能是 NUL——一來一回就不是同一份 bytes 了。
//
// 真的改了才編碼，並**補全形空白到 6 bytes**，否則新名字比舊的短時，
// 尾巴會留著上一個名字的後半（`state` 的寫回是 copy，不清欄位）。
func keepOrEncode(oldRaw, newUTF8 string) (string, error) {
	if decode(oldRaw) == newUTF8 {
		return oldRaw, nil
	}
	b, err := traditionalchinese.Big5.NewEncoder().Bytes([]byte(newUTF8))
	if err != nil {
		return "", fmt.Errorf("%q 轉不成 Big5：%w", newUTF8, err)
	}
	if len(b) > nameBytes {
		return "", fmt.Errorf("%q 編碼後 %d bytes，超過欄位的 %d",
			newUTF8, len(b), nameBytes)
	}
	for len(b) < nameBytes {
		b = append(b, 0xA1, 0x40) // 全形空白
	}
	return string(b), nil
}
