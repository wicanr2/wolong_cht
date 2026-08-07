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
