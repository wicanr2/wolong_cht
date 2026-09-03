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
		// ⚠ **UTF-8 語系包只切 NUL，不切空白**。原版是定長欄位，
		// 尾端的空白是補位所以要切掉；英文的 `struck {2}.` 那個空白
		// 是真的空白，切掉會變成 `struckXUCHANG.`（docs/spec/84）。
		return strings.TrimRight(string(b), "\x00")
	}
	out, ok := decodeRaw(b, enc)
	if !ok {
		// 解不開就退回原始 hex，**不要靜靜吃掉**——
		// 少了字的畫面看起來像排版 bug，會被查很久。
		return "?" + strings.ToUpper(hex(b))
	}
	return strings.TrimRight(out, "\x00　 ")
}

// decodeRaw 只做編碼轉換，不砍任何東西。
func decodeRaw(b []byte, enc Encoding) (string, bool) {
	var dec interface{ Bytes([]byte) ([]byte, error) }
	switch enc {
	case Big5:
		dec = traditionalchinese.Big5.NewDecoder()
	default:
		dec = japanese.ShiftJIS.NewDecoder()
	}
	out, err := dec.Bytes(b)
	if err != nil {
		return "", false
	}
	return string(out), true
}

// DecodeKeepPad 與 Decode 相同，但**保留行尾的全形空白**。
//
// 選單的框寬由「那一行有幾個全形字」決定（`docs/spec/45` §2.2），
// 而原版把每一列補到等寬——**那些空白是版面，不是補位噪音**。
// Decode 的 TrimRight 對定長欄位（武將名）是對的，對選單是錯的：
// 「　位置確認　」被砍成 5 個字，框就少了 16 px（`docs/spec/124`）。
func DecodeKeepPad(b []byte, enc Encoding) string {
	if enc == UTF8 {
		return strings.TrimRight(string(b), "\x00")
	}
	out, ok := decodeRaw(b, enc)
	if !ok {
		return "?" + strings.ToUpper(hex(b))
	}
	return strings.TrimRight(out, "\x00")
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
