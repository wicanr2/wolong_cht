package tactical

// 尋路請求佇列（`docs/re/80`、`docs/spec/120`）。
//
// ⭐ **原版的兵平常不尋路。** 每一幀直接比目前座標與這一步的目標，
// 三軸各試一步；**走不動才把自己排進這條佇列**。佇列由
// `sub_1AED2` 在逐兵迴圈**之前**消化，而且**每幀最多兩筆**——
// 所以「重算」不是每個兵各自的計時器，是一份**全域預算**。

// pathQueueSize 是佇列格數。原版是 `and di, 0FFh` 的環狀陣列
// （`CS:0xD352` 起 `100h` byte ＝ 128 筆 × 2 byte，`docs/re/80` §2）。
//
// ⭐ 128 不是隨便取的：入隊端用兵記錄 `+0x00` bit 4 去重，
// 在飛的請求最多 96 筆（兵的上限），**結構上塞不滿**。
const pathQueueSize = 128

// pathQueuePerFrame 是每幀消化幾筆（`sub_1AED2` 的 `mov cx, 2`）。
const pathQueuePerFrame = 2

// pathRequest 指到一個兵。原版存的是兵記錄的段內位址，同一件事。
type pathRequest struct{ side, k int }

// pathQueue 是環狀 FIFO。`head` ＝ `word_1D350`、`tail` ＝ `word_1D34E`。
type pathQueue struct {
	buf        [pathQueueSize]pathRequest
	head, tail int
	n          int
}

func (q *pathQueue) push(r pathRequest) bool {
	if q.n >= pathQueueSize {
		return false
	}
	q.buf[q.tail] = r
	q.tail = (q.tail + 1) % pathQueueSize
	q.n++
	return true
}

func (q *pathQueue) pop() (pathRequest, bool) {
	if q.n == 0 {
		return pathRequest{}, false
	}
	r := q.buf[q.head]
	q.head = (q.head + 1) % pathQueueSize
	q.n--
	return r, true
}

// requestPath 把一個兵排進佇列（原版 `sub_1C653`）。
//
// ⚠ **去重靠旗標不靠掃描**：原版 `test byte ptr [si], 10h / jnz` 一條指令
// 就擋掉重複排隊，所以同一個兵在佇列裡永遠只有一筆。
func (b *Battle) requestPath(side, k int) {
	s := &b.Sides[side].Soldiers[k]
	if s.PathQueued || !s.Alive {
		return
	}
	if !b.paths.push(pathRequest{side: side, k: k}) {
		// 結構上到不了（見 pathQueueSize）。真的滿了就這一次不排，
		// 下一次走不動還會再來——**不要 panic，也不要擠掉別人**。
		return
	}
	s.PathQueued = true
}

// drainPathQueue 消化佇列（原版 `sub_1AED2`）。
//
// ⭐ **呼叫點在逐兵迴圈之前**（`sub_1ADC8` 的順序），所以這一幀排進去的
// 請求最快也要下一幀才算得到。
func (b *Battle) drainPathQueue() {
	for i := 0; i < pathQueuePerFrame; i++ {
		r, ok := b.paths.pop()
		if !ok {
			return
		}
		s := &b.Sides[r.side].Soldiers[r.k]
		s.PathQueued = false
		if !s.Alive {
			continue // 排隊期間陣亡：位子讓出來，這一筆作廢
		}
		b.computePath(r.side, r.k)
	}
}
