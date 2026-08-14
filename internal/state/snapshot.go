package state

import "fmt"

// Snapshot 是 remake 存檔要保存、而**原版記錄裡沒有**的執行期狀態。
//
// 這些欄位在原版的 22,208 B 區塊裡沒有位置，所以走原版格式存檔一定會掉。
// 其中 CityCursor 尤其要命：它決定 AI 下一個處理哪個據點
// （`docs/spec/10`），掉了就等於每次讀檔都把 AI 的掃描進度重設。
//
// 規格見 `docs/spec/20-save-format.md` §2.4。**所有欄位都是 remake 差異**，
// 匯出回原版格式時會遺失。
type Snapshot struct {
	// Player 是玩家扮演的勢力。原版把它存在區塊裡（`playerOffset`），
	// 但 remake 允許 −1（旁觀），所以一併帶著走。
	Player int `json:"player"`

	// CityCursor 是據點整備的輪轉游標（原版 `word_10D1E`）。
	CityCursor int `json:"city_cursor"`

	// EventCursor／EventDelay 是事件佇列的執行期游標（`word_10D20`／`byte_131AD`）。
	EventCursor int   `json:"event_cursor"`
	EventDelay  uint8 `json:"event_delay"`

	// Routes 是每支軍團的行軍路線。**原版沒有這個欄位**——原版玩家只能
	// 對相鄰據點下令，不需要多段路由（`internal/state/corps.go`）。
	// 用 map 而不是陣列：存檔裡多半只有幾支軍團在走，陣列會塞滿 null。
	Routes map[int][][2]int `json:"routes,omitempty"`

	// StrategicAI／ApproximateEvent10 是執行期開關，不是世界狀態，
	// 但讀檔後要能回到同一個行為。
	StrategicAI        bool `json:"strategic_ai"`
	ApproximateEvent10 bool `json:"approximate_event10"`
}

// RawBlock 回傳**載入當時**那 22,208 B 原始區塊的副本。
//
// ⚠ 這不是「目前狀態存成原版格式」——那是 Bytes()。
// 兩者的差別正是 `docs/spec/20` §2.3 的權威分工：
// RawBlock 是未解區域的保存錨點，Bytes 是套用已解欄位之後的產物。
func (w *World) RawBlock() []byte {
	return append([]byte(nil), w.raw...)
}

// TakeSnapshot 取出執行期狀態。
func (w *World) TakeSnapshot() Snapshot {
	s := Snapshot{
		Player:             w.Player,
		CityCursor:         w.cityCursor,
		EventCursor:        w.eventCursor,
		EventDelay:         w.eventDelay,
		StrategicAI:        w.strategicAI,
		ApproximateEvent10: w.approximateEvent10,
	}
	for i := range w.routes {
		if len(w.routes[i]) == 0 {
			continue
		}
		if s.Routes == nil {
			s.Routes = make(map[int][][2]int, 4)
		}
		s.Routes[i] = append([][2]int(nil), w.routes[i]...)
	}
	return s
}

// Restore 把執行期狀態放回去。索引超界一律報錯而不是默默夾住——
// 這種檔案要嘛是版本不合、要嘛是被手改壞了，兩種都該讓呼叫端知道。
func (w *World) Restore(s Snapshot) error {
	if s.Player < -1 || s.Player >= numFactions {
		return fmt.Errorf("玩家勢力 %d 超出 −1–%d", s.Player, numFactions-1)
	}
	if s.CityCursor < 0 || s.CityCursor >= numCities {
		return fmt.Errorf("據點游標 %d 超出 0–%d", s.CityCursor, numCities-1)
	}
	if s.EventCursor < 0 || s.EventCursor >= eventQueueEntries {
		return fmt.Errorf("事件游標 %d 超出 0–%d", s.EventCursor, eventQueueEntries-1)
	}
	for i := range s.Routes {
		if i < 0 || i >= numCorps {
			return fmt.Errorf("行軍路線的軍團編號 %d 超出 0–%d", i, numCorps-1)
		}
	}
	w.Player = s.Player
	w.cityCursor = s.CityCursor
	w.eventCursor = s.EventCursor
	w.eventDelay = s.EventDelay
	w.strategicAI = s.StrategicAI
	w.approximateEvent10 = s.ApproximateEvent10
	for i := range w.routes {
		w.routes[i] = nil
	}
	for i, r := range s.Routes {
		w.routes[i] = append([][2]int(nil), r...)
	}
	return nil
}

// LoadBlock 從一個 22,208 B 的區塊建出 World。
//
// LoadScenario 是「從檔案的第 N 個區塊」，這一支是「從手上這段 bytes」——
// 原生存檔把區塊帶在自己裡面（`docs/spec/20` §2.2），需要這個入口。
func LoadBlock(b []byte) (*World, error) {
	if len(b) != blockSize {
		return nil, fmt.Errorf("區塊大小 %d，預期 %d", len(b), blockSize)
	}
	return loadBlock(b), nil
}
