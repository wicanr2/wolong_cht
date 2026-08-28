// Package namechars 解 `END_S15.DAT`：松崗版「自定」軍師命名視窗的選字表
// （docs/formats/10）。
//
// 2,621 個 Big5 字，照注音順序排（ㄅ八巴疤叭吧玻背杯…），每字 2 byte：
// **把 Big5 的高低 byte 對調之後加 0x1000** 存成 little-endian。
// 讀取端是 `KI.EXE` 的 `sub_190C0`／`sub_1928A`：
//
//	mov ax, es:[di] / xchg ah, al / sub ax, 1000h / xchg al, ah
//
// 檔名不是「結局」的一部分——`END_S13`／`S14` 是字型、`S15` 是這張表，
// 三個都只是松崗版塞在 `END_S*` 這串編號後面的資料檔。
package namechars

import (
	"fmt"
	"os"

	"golang.org/x/text/encoding/traditionalchinese"
)

// Count 是表裡的字數（5,242 B ÷ 2）。
const Count = 2621

// Table 是解出來的選字表，順序照檔案（＝注音順序）。
type Table struct {
	Runes []rune
	// Big5 是每個字的原始 Big5 碼（高 byte 在前），寫回名字緩衝區時用。
	Big5 []uint16
}

// Decode 解一份 END_S15.DAT 的內容。
func Decode(data []byte) (*Table, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("namechars: 長度 %d 不是偶數", len(data))
	}
	dec := traditionalchinese.Big5.NewDecoder()
	t := &Table{}
	for i := 0; i+1 < len(data); i += 2 {
		w := uint16(data[i]) | uint16(data[i+1])<<8 // little-endian
		sw := (w&0xFF)<<8 | w>>8                    // xchg ah, al
		v := sw - 0x1000                            // sub ax, 1000h
		// 之後再 xchg 回來存進緩衝區，緩衝區裡就是 [高 byte, 低 byte]。
		hi, lo := byte(v>>8), byte(v)
		s, err := dec.Bytes([]byte{hi, lo})
		if err != nil {
			return nil, fmt.Errorf("namechars: 第 %d 字 %02X%02X 不是 Big5", i/2, hi, lo)
		}
		rs := []rune(string(s))
		if len(rs) != 1 {
			return nil, fmt.Errorf("namechars: 第 %d 字解出 %d 個 rune", i/2, len(rs))
		}
		t.Runes = append(t.Runes, rs[0])
		t.Big5 = append(t.Big5, v)
	}
	return t, nil
}

// Load 從檔案讀。
func Load(path string) (*Table, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Decode(data)
}
