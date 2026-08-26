package text

import (
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/traditionalchinese"
)

// Decode 把原版的原始位元組轉成 UTF-8。
//
// 為什麼不在 Parse 就轉：`Table` 要能 round-trip 回原始檔案
// （docs/formats/01），所以解析層一律保留原始位元組，
// 只有要**顯示**的時候才呼叫這裡。
//
// ⚠ 用的是 **cp950／cp932 而不是 big5／shift_jis**。
// 兩者不是同一張表——踩過一次：`F9D8` 在 cp950 是「裏」，
// 用嚴格的 big5 解會失敗，當時被誤判成「原版有錯字」
// （CONTEXT.md「已被推翻的斷言」）。
// x/text 的 Big5 與 ShiftJIS 實作就是 cp950／cp932 的超集。
func Decode(b []byte, enc Encoding) string {
	if enc == UTF8 {
		return strings.TrimRight(string(b), "\x00　 ")
	}
	var dec interface{ Bytes([]byte) ([]byte, error) }
	switch enc {
	case Big5:
		dec = traditionalchinese.Big5.NewDecoder()
	default:
		dec = japanese.ShiftJIS.NewDecoder()
	}
	out, err := dec.Bytes(b)
	if err != nil {
		// 解不開就退回原始 hex，**不要靜靜吃掉**——
		// 少了字的畫面看起來像排版 bug，會被查很久。
		return "?" + strings.ToUpper(hex(b))
	}
	return strings.TrimRight(string(out), "\x00　 ")
}

func hex(b []byte) string {
	const d = "0123456789abcdef"
	var sb strings.Builder
	for _, c := range b {
		sb.WriteByte(d[c>>4])
		sb.WriteByte(d[c&0xF])
	}
	return sb.String()
}
