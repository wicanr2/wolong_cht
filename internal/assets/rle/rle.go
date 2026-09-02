// Package rle 是原版的 RLE 解壓。
//
// 出處：`KI.EXE` 的 `sub_1F5E7` —— `MMAP.MAP` 專用的載入器，
// 與其他資料檔走的 `sub_1F4A2`（單純讀檔）不同支。
//
// 演算法**用「連續兩個相同的 byte」當 run 的觸發，沒有逃脫字元**：
//
//	逐 byte 複製；
//	一旦輸出的 byte 與前一個相同，下一個輸入 byte 是「再重複幾次」；
//	次數 0 表示那兩個相同的 byte 就只是字面值，回到逐 byte 模式。
//
// 所以一段 run 的總長是 2 + count，最短的 run（count=0）不省空間但也不虧
// —— 這是為了不需要逃脫字元付的代價。
package rle

import (
	"encoding/binary"
	"fmt"
)

// Decode 解壓 src。
//
// 不接受「預期長度」參數是刻意的：原版的解壓器也不知道目標長度，
// 它一路解到檔案結束。呼叫端自己截到要的長度
// （`MMAP.MAP` 解出來會比 384×256 多幾個 byte，見 docs/formats/06）。
func Decode(src []byte) []byte {
	out := make([]byte, 0, len(src)*2)
	for i := 0; i < len(src); {
		prev := src[i]
		out = append(out, prev)
		i++

		// 逐 byte 複製，直到出現與前一個相同的 byte
		matched := false
		for i < len(src) {
			cur := src[i]
			out = append(out, cur)
			i++
			if cur == prev {
				matched = true
				break
			}
			prev = cur
		}
		if !matched || i >= len(src) {
			break
		}

		count := int(src[i])
		i++
		for n := 0; n < count; n++ {
			out = append(out, prev)
		}
	}
	return out
}

// Header 是檔案前面那 4 byte：小端 u32，值就是解壓後的長度。
//
// 原版的載入器在讀第一個 byte 之前先 `LSEEK` 過它
// （`KI.EXE` 的 `sub_1F655`、`D7OPEN.EXE` 的 `sub_10E04`、`D7END.EXE` 同一段，
// 三個執行檔逐指令相同——docs/re/76 §5、docs/spec/113）。
const Header = 4

// DecodeFile 解一整個 RLE 資料檔：跳過長度頭，解剩下的，核對長度。
//
// ⭐ **長度是驗收條件，不是參考值。** 從 offset 0 解會在某一處掉相位，
// 症狀是長度差幾十個 byte、畫面整體位移（一張 320 px 寬的圖會在
// x = 160 冒出一條垂直接縫）。`MMAP.MAP` 是唯一躲過的檔：它的頭
// `00 80 01 00` 沒有相鄰重複，RLE 原樣吐出而且相位不變。
//
// Decode 本身不動——它是演算法，長度頭是容器；混在一起之後，
// 任何「拿一段裸資料試解」的呼叫端都得先湊一個假檔頭。
func DecodeFile(src []byte) ([]byte, error) {
	if len(src) < Header {
		return nil, fmt.Errorf("rle: 只有 %d B，連 %d B 的長度頭都不夠",
			len(src), Header)
	}
	want := int(binary.LittleEndian.Uint32(src[:Header]))
	out := Decode(src[Header:])
	if len(out) != want {
		return nil, fmt.Errorf("rle: 檔頭宣告 %d B，解出來 %d B——"+
			"差一個 byte 就是解錯，不是尾巴沒編進去", want, len(out))
	}
	return out, nil
}
