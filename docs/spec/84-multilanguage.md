# 84 — 多語系：簡體中文與英文支援

**狀態：DRAFT → 第一期 CONFORMED（簡體全鏈＋英文骨架，2026-08-26）。**

- 日期：2026-08-26
- 出處：這是 **remake 差異**（現代化外殼），不對應任何原版機制。
  架構政策出自 `CLAUDE.md` §6（繁中母本、文字進語系檔、排版層抽離）；
  簡體／英文這一輪啟動是使用者裁定（2026-08-26）。
  母本文字的來源仍是松崗版 `TALK.DAT`（[`../formats/01`](../formats/01-talk-dat.md)）
  ＋ 60 筆校訂（`translations/corrections.json`）。

## 1. 一個語系由四樣東西組成

| 元件 | zh-Hant（母本） | zh-Hans | en |
|---|---|---|---|
| talk 表（1,022 則） | `translations/talk-dosv-corrected.json`（cp950） | `translations/talk-zh-hans.json`（utf-8，OpenCC t2s 機轉初稿） | 未譯——**fallback 母本**（第二期人工翻譯） |
| UI 詞表（Go 內字串） | 原文即母本 | **字級 t2s 表 runtime 轉換**（`translations/t2s-chars.json`，3,130 字，掛在 textdraw 的 rune 替換層） | 逐句詞表（第二期抽字串後接 `uitext.Convert`；本輪顯示繁中） |
| 人名地名（`SINARIO.DAT` raw Big5） | `big5()` 解碼 | 解碼後過字級 t2s 表 | 第二期（需要羅馬化／譯名 glossary）——先顯示繁中 |
| 全形字型 | 倚天 `STDFONT.15`＋`SPCFONT.15`（16×15，Big5 索引） | **`HZK16`**（GB2312 16×16，UCDOS 格式，使用者自備） | 半形 `ASCFONT.15` 既有 |

設計取捨：

- **talk 走逐則語系檔、UI 走字級轉換**。talk 是成句文本，字級轉換的
  一簡對多繁誤差（乾／干、後／后）會露在劇情裡，所以用 OpenCC 的
  詞級轉換離線產檔，之後逐則人工校訂；UI 詞是短標籤，字級表的錯誤率
  低且可隨時以 `ui-zh-hans.json` 覆寫單則（覆寫優先於字級表）。
- **字型不隨專案散布**（同倚天的處理）：HZK16 是 UCDOS 的商業字型，
  由使用者放進 `-font-dir`；缺檔時照既有規則畫空心方框，遊戲照跑。
- **fallback 一律回母本繁中**，不顯示空字串——缺譯要看得見。

## 2. 改動

| 層 | 內容 |
|---|---|
| `internal/assets/text` | `Encoding` 增 `UTF8`（`Decode`／`encodeText` passthrough）；`Table` 記住自己的編碼，`Lines()` 不再寫死 Big5（`Parse` 預設 Big5，零值相容） |
| `internal/assets/cjk` | `LoadHZK16(path)`：GB2312 區位索引（`((qu−0xA1)×94＋(wei−0xA1))×32`），16×16、每列 2 byte、MSB-first。自我檢查：「啊」（B0A1）第一個字非空。`Font.Glyph` 依來源分派；**字高 16 比倚天多 1px**，版面常數不動，多的一列畫在行距裡 |
| `internal/ui/uitext` | 新套件：語系代號解析、逐句覆寫詞表與字級表的載入（`Load`／`Convert`／`RuneMap`）。**字級表實際掛在 `textdraw.Drawer.SetRuneMap`**——那是「同一個字選哪個字形」的層，一次涵蓋 Go 內 literal、人名與 talk fallback；簡體語系包的文字再過表是恆等。逐句覆寫詞表（en 與簡體例外詞）由 `Convert` 提供，第二期抽字串時接上 |
| `cmd/wlgame` | `-lang` flag（`zh-hant` 預設／`zh-hans`／`en`）：選 talk 檔預設值、載入詞表與字級表、字型載入順序（zh-hans 先試 `HZK16` 再退倚天） |
| `tools/langpack.py` | 離線產生器（需網路容器跑 OpenCC）：`talk-zh-hans.json`＋`t2s-chars.json`。**產出進版控**，重跑可重生 |

## 3. 分期

| 期 | 內容 | 狀態 |
|---|---|---|
| 1 | 上表全部＋簡體端到端可玩（talk 機轉初稿＋HZK16＋人名轉換）；英文＝flag＋UI 詞表機制＋talk fallback | **本輪** |
| 2 | 簡體 talk 逐則人工校訂（機轉標 DRAFT）；英文 talk 全量翻譯＋人名 glossary；UI 字串全量抽出 | 未排 |
| 3 | 手機版接線；英文比例字寬排版（現為半形等寬 8px）；行寬 guard 對非全形語系重算 | 未排 |

## 4. 驗證

- 單測：`text` UTF8 表載入與 `Lines()`；`cjk` HZK16 索引自檢與 Glyph；
  `uitext` 覆寫優先序與 fallback。
- `tools/langpack.py --selftest`：marker `{N}` 逐則保留、則數 1,022、
  簡轉結果不含殘留繁體常用字抽樣。
- 實跑（2026-08-26）：`-lang zh-hans -open-talk-index 70` 截圖，
  訊息框顯示「□昌□生了暴□雨。」——語系包生效（`{2}` 展開成「许昌」、
  字級表把人名轉成簡體），簡體字在無 HZK16 時按政策畫缺字方框、
  繁簡同形字正常。有 HZK16 時全字顯示（字型使用者自備，未隨測）。

## 5. 未解

| 缺口 | 下手點 |
|---|---|
| 簡體 talk 是機轉初稿，未逐則校訂 | 第二期：沿用 `corrections.json` 的覆寫格式建 `corrections-zh-hans.json` |
| 英文 talk／人名未譯 | 第二期：glossary 先行（三國人名有既定英譯慣例，Wade-Giles vs Pinyin 要先裁定） |
| 英文半形排版未重算 | 行寬 guard 與訊息框 22 格是全形格假設；en 全量翻譯前一併處理 |
