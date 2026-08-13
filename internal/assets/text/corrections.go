package text

import (
	"encoding/json"
	"fmt"
	"os"
)

// correctionManifest 是可散布的最小覆蓋表。它不含完整原版訊息表；每一筆都先
// 驗證玩家自備 TALK.DAT 的原文，避免把不同版本的資料靜默改壞。
type correctionManifest struct {
	Corrections []correctionItem `json:"corrections"`
}

type correctionItem struct {
	ID  int     `json:"id"`
	CHT string  `json:"cht"`
	Fix *string `json:"fix"`
}

// ApplyCorrections 對已解析的玩家自備繁中 TALK.DAT 套用已定案校訂。
//
// 發行包不能帶出完整的 translations/talk-dosv-corrected.json，因為那是由
// 原版文字抽取的 1,022 則完整表。這條路徑只讀取公開的 corrections.json，將
// 每一項目前原文與 manifest 的 cht 欄位作 fail-closed 比對後才套用 fix。
func ApplyCorrections(table *Table, path string, enc Encoding) (*Table, error) {
	if table == nil {
		return nil, fmt.Errorf("TALK 表不可為空")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("讀不到校訂覆蓋 %s：%w", path, err)
	}
	var manifest correctionManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("校訂覆蓋不是有效 JSON：%w", err)
	}
	if len(manifest.Corrections) == 0 {
		return nil, fmt.Errorf("校訂覆蓋沒有 corrections")
	}

	result := &Table{}
	for i, msg := range table.Messages {
		result.Messages[i] = cloneMessage(msg)
	}
	byID := make(map[int]correctionItem, len(manifest.Corrections))
	for _, item := range manifest.Corrections {
		byID[item.ID] = item
	}
	for _, item := range manifest.Corrections {
		if item.Fix == nil {
			continue
		}
		if item.ID < 0 || item.ID >= len(result.Messages) {
			return nil, fmt.Errorf("校訂編號超出 TALK.DAT：#%d", item.ID)
		}
		current, err := messageText(result.Messages[item.ID], enc)
		if err != nil {
			return nil, fmt.Errorf("讀取 #%d 現況失敗：%w", item.ID, err)
		}
		expected, err := correctionSourceText(item, byID)
		if err != nil {
			return nil, err
		}
		if current != expected {
			return nil, fmt.Errorf("#%d 的現況與校訂覆蓋不符：得到 %q，預期 %q", item.ID, current, expected)
		}
		lines, err := wrapCorrection(*item.Fix, result.Messages[item.ID], enc)
		if err != nil {
			return nil, fmt.Errorf("#%d 校訂換行失敗：%w", item.ID, err)
		}
		msg, err := messageFromJSONLines(lines, enc)
		if err != nil {
			return nil, fmt.Errorf("#%d 校訂文字無法轉換：%w", item.ID, err)
		}
		result.Messages[item.ID] = msg
	}
	return result, nil
}

func correctionSourceText(item correctionItem, byID map[int]correctionItem) (string, error) {
	expected := item.CHT
	runes := []rune(expected)
	if len(runes) < 2 || runes[0] != '同' || runes[1] != ' ' {
		return expected, nil
	}
	var ref int
	if _, err := fmt.Sscanf(expected, "同 %d", &ref); err != nil || ref < 0 {
		return "", fmt.Errorf("#%d 的同文參照無效：%q", item.ID, expected)
	}
	referenced, ok := byID[ref]
	if !ok {
		return "", fmt.Errorf("#%d 的同文參照無法展開：%q", item.ID, expected)
	}
	refRunes := []rune(referenced.CHT)
	if len(refRunes) >= 2 && refRunes[0] == '同' && refRunes[1] == ' ' {
		return "", fmt.Errorf("#%d 的同文參照無法展開：%q", item.ID, expected)
	}
	return referenced.CHT, nil
}

// wrapCorrection 複製 tools/talkdat.py correct 的保守換行契約。它不是原版
// renderer 的逐像素宣稱；目的只是讓公開的最小覆蓋在 runtime 套用時與既有、
// 已驗收的 talk-dosv-corrected.json 保持相同的行／尾端空行資料。
func wrapCorrection(fix string, original Message, enc Encoding) ([]string, error) {
	originalLines := make([]string, len(original.Lines))
	for i, line := range original.Lines {
		text, err := messageText(Message{Lines: []Line{line}}, enc)
		if err != nil {
			return nil, err
		}
		originalLines[i] = text
	}
	trailing := 0
	for i := len(originalLines) - 1; i >= 0 && originalLines[i] == ""; i-- {
		trailing++
	}
	body := originalLines[:len(originalLines)-trailing]
	// 與 tools/talkdat.py correct 的 max(..., default=10) 對齊。default 只在
	// 沒有任何非空原行時才生效；不能把 10 當成下限，否則原本較窄的訊息會
	// 改變已驗收的硬換行。
	width := 0
	for _, line := range body {
		if lineWidth(line) > width {
			width = lineWidth(line)
		}
	}
	if width == 0 {
		width = 10
	}
	lines := wrapText(fix, width)
	for range trailing {
		lines = append(lines, "")
	}
	return lines, nil
}

type correctionToken struct {
	text  string
	width int
}

func correctionTokens(text string) []correctionToken {
	runes := []rune(text)
	tokens := make([]correctionToken, 0, len(runes))
	for i := 0; i < len(runes); {
		if runes[i] == '{' && i+2 < len(runes) && runes[i+2] == '}' && runes[i+1] >= '0' && runes[i+1] <= '9' {
			start, width := i, 0
			for i+2 < len(runes) && runes[i] == '{' && runes[i+2] == '}' && runes[i+1] >= '0' && runes[i+1] <= '9' {
				width += 3
				i += 3
			}
			tokens = append(tokens, correctionToken{text: string(runes[start:i]), width: width})
			continue
		}
		tokens = append(tokens, correctionToken{text: string(runes[i]), width: 1})
		i++
	}
	return tokens
}

func lineWidth(text string) int {
	width := 0
	for _, token := range correctionTokens(text) {
		width += token.width
	}
	return width
}

func wrapText(text string, width int) []string {
	if text == "" {
		return []string{""}
	}
	if width < 1 {
		width = 1
	}
	const closing = "，。！？；：、）】」』"
	lines := make([]string, 0, 2)
	var current []rune
	used := 0
	for _, token := range correctionTokens(text) {
		isClosing := len([]rune(token.text)) == 1 && containsRune(closing, []rune(token.text)[0])
		if isClosing {
			if len(current) == 0 && len(lines) > 0 {
				lines[len(lines)-1] += token.text
				continue
			}
			if len(current) > 0 && used+token.width > width {
				current = append(current, []rune(token.text)...)
				used += token.width
				continue
			}
		}
		if len(current) > 0 && used+token.width > width {
			lines = append(lines, string(current))
			current, used = nil, 0
		}
		current = append(current, []rune(token.text)...)
		used += token.width
	}
	if len(current) > 0 {
		lines = append(lines, string(current))
	}
	return lines
}

func containsRune(text string, target rune) bool {
	for _, r := range text {
		if r == target {
			return true
		}
	}
	return false
}

func cloneMessage(in Message) Message {
	out := Message{Lines: make([]Line, len(in.Lines))}
	for i, line := range in.Lines {
		out.Lines[i].Parts = make([]Part, len(line.Parts))
		for j, part := range line.Parts {
			out.Lines[i].Parts[j] = Part{Marker: part.Marker, Raw: append([]byte(nil), part.Raw...)}
		}
	}
	return out
}

func messageText(msg Message, enc Encoding) (string, error) {
	var out []byte
	for _, line := range msg.Lines {
		for _, part := range line.Parts {
			if part.Marker != 0 {
				out = append(out, '{', part.Marker, '}')
				continue
			}
			out = append(out, part.Raw...)
		}
	}
	return Decode(out, enc), nil
}
