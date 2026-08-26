package text

import (
	"encoding/json"
	"fmt"
	"os"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// jsonTable 是 tools/talkdat.py export／correct 使用的文字表格式。
// JSON 是呈現層的可版控來源；raw TALK.DAT 仍由 Parse 保留，兩者不要混為一談。
type jsonTable struct {
	Encoding string     `json:"encoding"`
	Messages [][]string `json:"messages"`
}

// LoadJSON 載入一份以 UTF-8 儲存、每行已經排版好的訊息表。
//
// 這條路徑只給 remake 的呈現層使用。它會把 marker 重新切成 Part，並用
// 原版對應編碼編回 Raw，因此 talkLines 不需要另造一套 marker parser。
// raw TALK.DAT 的 Parse → Bytes round-trip 不經過這裡，也不會被覆寫。
func LoadJSON(path string, enc Encoding) (*Table, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("talk json: 讀不到 %s：%w", path, err)
	}
	return parseJSON(raw, path, enc)
}

// ParseJSON 與 LoadJSON 相同，但吃已經讀進來的位元組——
// 語系包會內嵌進執行檔（docs/spec/86 §2），那條路沒有檔案可以開。
func ParseJSON(raw []byte, enc Encoding) (*Table, error) {
	return parseJSON(raw, "<embedded>", enc)
}

func parseJSON(raw []byte, path string, enc Encoding) (*Table, error) {
	var in jsonTable
	if err := json.Unmarshal(raw, &in); err != nil {
		return nil, fmt.Errorf("talk json: %s 不是有效 JSON：%w", path, err)
	}
	if len(in.Messages) != MessageCount {
		return nil, fmt.Errorf("talk json: %s 有 %d 則訊息，預期 %d",
			path, len(in.Messages), MessageCount)
	}

	t := &Table{enc: enc}
	for i, lines := range in.Messages {
		msg, err := messageFromJSONLines(lines, enc)
		if err != nil {
			return nil, fmt.Errorf("talk json: %s #%d：%w", path, i, err)
		}
		t.Messages[i] = msg
	}
	return t, nil
}

func messageFromJSONLines(lines []string, enc Encoding) (Message, error) {
	msg := Message{Lines: make([]Line, len(lines))}
	for i, line := range lines {
		parts, err := partsFromJSONLine(line, enc)
		if err != nil {
			return Message{}, fmt.Errorf("第 %d 行：%w", i, err)
		}
		msg.Lines[i] = Line{Parts: parts}
	}
	return msg, nil
}

func partsFromJSONLine(line string, enc Encoding) ([]Part, error) {
	runes := []rune(line)
	parts := make([]Part, 0, 2)
	var raw []rune
	flush := func() error {
		if len(raw) == 0 {
			return nil
		}
		encoded, err := encodeText(string(raw), enc)
		if err != nil {
			return err
		}
		parts = append(parts, Part{Raw: encoded})
		raw = raw[:0]
		return nil
	}

	for i := 0; i < len(runes); {
		if runes[i] == '{' && i+2 < len(runes) && runes[i+2] == '}' &&
			runes[i+1] <= 0x7f {
			if err := flush(); err != nil {
				return nil, err
			}
			parts = append(parts, Part{Marker: byte(runes[i+1])})
			i += 3
			continue
		}
		raw = append(raw, runes[i])
		i++
	}
	if err := flush(); err != nil {
		return nil, err
	}
	return parts, nil
}

func encodeText(s string, enc Encoding) ([]byte, error) {
	if !utf8.ValidString(s) {
		return nil, fmt.Errorf("文字不是有效 UTF-8")
	}
	if enc == UTF8 {
		return []byte(s), nil
	}
	var encoder interface {
		Bytes([]byte) ([]byte, error)
	}
	switch enc {
	case Big5:
		encoder = traditionalchinese.Big5.NewEncoder()
	case ShiftJIS:
		encoder = japanese.ShiftJIS.NewEncoder()
	default:
		return nil, fmt.Errorf("不支援的編碼 %d", enc)
	}
	encoded, err := encoder.Bytes([]byte(s))
	if err != nil {
		return nil, fmt.Errorf("無法以原版編碼表示 %q：%w", s, err)
	}
	return encoded, nil
}
