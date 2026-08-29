// Package text 解 `TALK.DAT` 訊息表。
//
// 規格：docs/formats/01-talk-dat.md（READY）
//
// 驗收標準是 byte-for-byte round-trip：解出來再組回去必須與原檔完全相同。
// 寫不回去的中文化工具是不能用的。
package text

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	// TableBytes 是偏移表的長度。
	TableBytes = 0x800
	// TableEntries 是偏移表的筆數。
	TableEntries = TableBytes / 2
	// MessageCount 是實際的訊息數。[1022] 是哨兵、[1023] 未使用。
	MessageCount = 1022
)

// Message 是一則訊息：若干行，加上這則用到的變數標記。
//
// **結尾的空行是資料不是格式** —— 有的訊息一個都沒有，有的有兩個，
// 那是對話框裡的留白。把它正規化掉就寫不回去了。
type Message struct {
	Lines []Line
}

// Line 是一行，由文字片段與變數標記交錯組成。
type Line struct {
	Parts []Part
}

// Part 要嘛是一段原始位元組（文字），要嘛是一個變數標記。
type Part struct {
	// Marker 非 0 時代表這是變數標記，值是 '1'–'7' 的 ASCII。
	// 原版用到六種，沒有 '5'。
	Marker byte
	// Raw 是文字片段的原始位元組（Big5 或 Shift-JIS，未解碼）。
	Raw []byte
}

// Table 是整份 TALK.DAT。
type Table struct {
	Messages [MessageCount]Message
	// enc 是 Raw 片段的編碼，Lines() 用它解碼。零值 Big5，
	// 與既有的 Parse 呼叫端相容；LoadJSON 依語系包設定。
	enc Encoding
}

// Parse 解一份 TALK.DAT。
//
// ⚠ table[1022] 是**最後一個位元組的位移**，不是 one-past-end。
// 照它切最後一則會少一個 byte，所以最後一則要讀到 EOF。
func Parse(data []byte) (*Table, error) {
	if len(data) < TableBytes {
		return nil, fmt.Errorf("talk: 檔案只有 %d byte，連偏移表都不夠", len(data))
	}
	offsets := make([]int, TableEntries)
	for i := range offsets {
		offsets[i] = int(binary.LittleEndian.Uint16(data[i*2:]))
	}
	t := &Table{}
	for i := 0; i < MessageCount; i++ {
		end := len(data)
		if i < MessageCount-1 {
			end = offsets[i+1]
		}
		if offsets[i] > end || end > len(data) {
			return nil, fmt.Errorf("talk: 第 %d 則的範圍 [%d,%d) 不合法",
				i, offsets[i], end)
		}
		t.Messages[i] = parseMessage(data[offsets[i]:end])
	}
	return t, nil
}

// isLead 判斷這個位元組是不是雙位元組字的第一個。
//
// 這一步不能省：Big5 與 Shift-JIS 的**第二個位元組都可能是 0x5C**（反斜線），
// 不先判斷就會把中文字切一半，掃出一堆假的變數標記。
func isLead(b byte, enc Encoding) bool {
	switch enc {
	case Big5:
		return b >= 0x81 && b <= 0xFE
	case ShiftJIS:
		return (b >= 0x81 && b <= 0x9F) || (b >= 0xE0 && b <= 0xFC)
	}
	return false
}

// Encoding 是訊息文字的編碼。
type Encoding int

const (
	Big5     Encoding = iota // 松崗 DOS/V 繁中版
	ShiftJIS                 // PC-98 日文原版
	// UTF8 是 remake 語系包（簡體／英文）的編碼：Raw 直接存 UTF-8，
	// 不再受「必須能編回原版編碼」的約束（docs/spec/84）。
	// 原版檔案的 byte-for-byte round-trip 只適用 Big5／ShiftJIS 表。
	UTF8
)

// encoding 是解析時假設的編碼。兩種編碼的雙位元組首碼範圍不同，
// 但變數標記的掃描只需要「這是不是雙位元組字的開頭」，
// 而兩者的聯集是安全的：0x81–0xFE 涵蓋 Shift-JIS 的兩段。
//
// 這裡刻意不解碼成 string —— 解碼是呈現層的事，
// 這一層只負責把結構切對，並保證能原樣組回去。
func parseMessage(raw []byte) Message {
	var msg Message
	var line Line
	var run []byte

	flushRun := func() {
		if len(run) > 0 {
			line.Parts = append(line.Parts, Part{Raw: run})
			run = nil
		}
	}
	for i := 0; i < len(raw); {
		c := raw[i]
		switch {
		case c == 0x5C && i+1 < len(raw):
			flushRun()
			line.Parts = append(line.Parts, Part{Marker: raw[i+1]})
			i += 2
		case c == 0x00:
			flushRun()
			msg.Lines = append(msg.Lines, line)
			line = Line{}
			i++
		case isLead(c, Big5) && i+1 < len(raw):
			run = append(run, raw[i], raw[i+1])
			i += 2
		default:
			run = append(run, c)
			i++
		}
	}
	flushRun()
	if len(line.Parts) > 0 { // 沒有 NUL 收尾的殘行
		msg.Lines = append(msg.Lines, line)
	}
	return msg
}

// Bytes 把一則訊息組回原始位元組。
func (m Message) Bytes() []byte {
	var out []byte
	for _, line := range m.Lines {
		for _, p := range line.Parts {
			if p.Marker != 0 {
				out = append(out, 0x5C, p.Marker)
				continue
			}
			out = append(out, p.Raw...)
		}
		out = append(out, 0x00)
	}
	return out
}

// Markers 回傳這則訊息用到的變數標記，依出現順序。
//
// 兩版比對時用得上：標記集合不一致的訊息就是譯文缺陷
// （見 docs/reference/02）。
func (m Message) Markers() []byte {
	var out []byte
	for _, line := range m.Lines {
		for _, p := range line.Parts {
			if p.Marker != 0 {
				out = append(out, p.Marker)
			}
		}
	}
	return out
}

// Bytes 把整份表組回原始位元組，用於 round-trip 驗收。
func (t *Table) Bytes() []byte {
	body := make([]byte, 0, 64*1024)
	offsets := make([]uint16, TableEntries)
	for i := 0; i < MessageCount; i++ {
		offsets[i] = uint16(TableBytes + len(body))
		body = append(body, t.Messages[i].Bytes()...)
	}
	offsets[MessageCount] = uint16(TableBytes + len(body) - 1) // 哨兵
	offsets[MessageCount+1] = 0                                // 未使用

	out := make([]byte, TableBytes, TableBytes+len(body))
	for i, v := range offsets {
		binary.LittleEndian.PutUint16(out[i*2:], v)
	}
	return append(out, body...)
}

// Lines 取出一則訊息的原始行並代入變數。
//
// ⚠ **缺變數時整則 fail-closed**（回 false），不顯示半句、也不把錯誤的索引
// 當成文字印出去。半句訊息在畫面上看起來像遊戲內容，會被當成原版行為記下來。
//
// ⚠ 最後一行的空行是 `TALK.DAT` 每則訊息的 NUL 結束行，原版 `sub_1084A`
// 讀到它就停，不畫成可見列——所以要砍掉。**中間**真正存在的空行必須保留。
//
// 桌面版與手機版共用這一份；兩邊各寫一份會長出不同的變數展開行為。
func (t *Table) Lines(index int, vars map[byte]string) ([]string, bool) {
	return t.LinesSeq(index, vars, nil)
}

// LinesSeq 與 Lines 相同，但**同一個標記出現多次時依序取 seq 的下一個值**。
//
// ⭐ 原版的 formatter 是一個共用的堆疊游標（`sub_14EB9` 的 `mov di, sp`），
// 每個標記消耗下一個參數，所以「{1}大人的兵馬，遇上{1}的兵馬了」兩個 `{1}`
// 是**兩個不同的武將**。`vars` 那張 map 一個標記只放得下一個值，
// 全庫 9 則有重複標記，其中 #29（兩個 `{1}`）與 #217（兩個 `{3}`）
// 真的要不同值，其餘 7 則重複的是 `{6}`（排版控制，空字串）。
//
// seq[m] 用完或沒給時退回 vars[m]。
func (t *Table) LinesSeq(index int, vars map[byte]string, seq map[byte][]string) ([]string, bool) {
	if t == nil || index < 0 || index >= len(t.Messages) {
		return nil, false
	}
	out := make([]string, 0, len(t.Messages[index].Lines))
	used := make(map[byte]int, len(seq))
	for _, line := range t.Messages[index].Lines {
		var b strings.Builder
		for _, part := range line.Parts {
			if part.Marker != 0 {
				value, ok := vars[part.Marker]
				if !ok {
					return nil, false
				}
				if s := seq[part.Marker]; used[part.Marker] < len(s) {
					value = s[used[part.Marker]]
				}
				used[part.Marker]++
				b.WriteString(value)
				continue
			}
			if len(part.Raw) > 0 {
				b.WriteString(Decode(part.Raw, t.enc))
			}
		}
		out = append(out, b.String())
	}
	if len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out, true
}
