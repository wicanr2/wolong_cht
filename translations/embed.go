// Package translations 只做一件事：把語系檔嵌進執行檔（docs/spec/86 §2）。
//
// 放在資料目錄裡而不是複製一份到別的套件底下，是因為 `go:embed` 只吃
// **自己套件目錄底下**的檔案——複製一份就會有兩份會漂開的同名資料。
//
// 這裡嵌的全是本專案的產出（譯文與名表），**原版資產一個都沒有**；
// 原版資料走的是 `docs/spec/72` 的 `gamedata/`，deny-list 照掃。
package translations

import "embed"

// Files 是 `translations/*.json`：三個語系的訊息包、三語名表、
// UI 詞表與簡體字級表。
//
//go:embed *.json
var Files embed.FS
