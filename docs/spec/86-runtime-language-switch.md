# 86 — 執行期切換語言（含手機版）

**狀態：CONFORMED（2026-08-26 實作並實跑驗過）。**

- 日期：2026-08-26
- 出處：**remake 差異**，不對應原版機制（原版只有一個語言）。
  語系包的結構見 [`84`](84-multilanguage.md)，清單欄界見
  [`85`](85-latin-list-layout.md)。
- 推論等級：不涉及原版事實。

## 1. 為什麼 `-lang` 不夠

| 情境 | 為什麼要執行期切換 |
|---|---|
| **Android** | 沒有命令列。使用者裁定（2026-08-26）Android 也要能切語言，那就只能做在畫面裡 |
| 推廣片 | 要拍「同一個局面換一種語言」，重開程式就不是同一個局面了 |
| 桌面 | 想比對譯文時，重開遊戲會丟掉當下的世界狀態 |

## 2. 語系包內嵌進執行檔

`translations/*.json` 共 **370 KB**，`go:embed` 進執行檔，
**檔案系統上的同名檔優先**（開發時改檔立即生效，不必重編）。

理由：Android 的資料是使用者自己匯入的（`DataRoot()/orig`），
語系包不是原版資產、也不該叫使用者自己產。內嵌之後手機端零設定。

政策依據：`translations/` **本來就在公開 repo 裡**（`CLAUDE.md` §1
「公開產出只有引擎程式碼與譯文校訂紀錄」），內嵌不改變散布範圍。
**原版資產仍然一個都不進執行檔**——`gamedata/` 走的是另一條路
（`docs/spec/72`），deny-list 照掃。

## 3. 一個語系載入層，兩端共用

新增 `internal/ui/langpack`：把「換一個語言要動哪些東西」收在一處。

```go
type Pack struct { Lang uitext.Language; Text *uitext.Table
                   Talk *text.Table; Font textdraw.GlyphSource }
func Load(lang uitext.Language, fontDir string) (*Pack, error)
func (p *Pack) Apply(d *textdraw.Drawer)
```

| 換語言要動的 | 誰負責 |
|---|---|
| `TALK.DAT` 的 1,022 則 | `Pack.Talk`（母本語系是 nil ＝ 用原版解出來的那份）|
| UI 詞、人名地名 | `Pack.Text`（`uitext.Table`）|
| 全形字型鏈 | `Pack.Font`（`JISKAN16`／`HZK16` → 倚天）|
| 畫的那一刻 | `Pack.Apply(d)` 掛 `SetRuneMap`／`SetTranslator` |

**桌面與手機各自的 `-lang` 初始化改成呼叫這一層**，
兩邊不會再長出不同的語系行為（`CLAUDE.md` §7 第 6 條：一條規則一份實作）。

## 4. 切換的入口

**⛔ 不動原版畫面。** 系統選單是逐像素對過的（`docs/playtest/39`），
不加第五列；語言切換一律放在 **remake 自己的畫面**上：

| 端 | 入口 | 說明 |
|---|---|---|
| 桌面 | **啟動殼層**標題頁的最後一列 `LANGUAGE` | 那一頁本來就是 remake 的（`docs/spec/79`）。進去時游標停在目前語言上，選完停在原地——**下一列就是換好的樣子** |
| 桌面 | **F9 循環**：繁中 → 简体 → 日本語 → English | 遊戲中即時切，推廣片用得到；**不存進存檔** |
| 手機 | 系統面板第三頁「語言」，四列點一下就換 | `internal/ui/phone/sheet.go` 的 tab（速度／存檔／語言／關於）|

四個語言的名字**用該語言自己的寫法**寫（`langpack.Choices`，桌面與手機同一份）。
理由：換過去之後整個畫面都是那個語言，用母本的寫法列出來，
玩家反而認不出自己剛才選了什麼。

⚠ 這帶出一個字型問題：「简」不在 Big5、平假名不在 GB2312，
只掛自己那一套字型的話，選單上會有幾格方框，**偏偏那幾格正是用來認語言的**。
所以 `fontChain` 每個語系都把另外兩套字集接在鏈的後面墊底
（順序：本語系 → 倚天 Big5 → 其餘）。原版文字用不到這些碼位，
既有的逐像素對拍不受影響。同理，選中記號用 `●` 不用 `▶`——後者不在倚天 Big5 裡。

切換要重建的東西與 §3 同一份；世界狀態、時鐘、存檔全部不動——
**換的是呈現，不是遊戲**。

## 5. 手機版接語系

手機版目前完全沒有語系（`84` §6 的缺口）。三個出口與桌面相同：

| 出口 | 桌面 | 手機 |
|---|---|---|
| 訊息 | `library.LoadOptions.TalkJSON` | `Session.SetLanguage` 換 `s.lib.Talk` |
| 人名地名 | `cmd/wlgame` 的 `big5()` | `Session.Localise()`（原本的 `big5()` 全數改走它）|
| UI 詞、字型 | `Drawer.SetTranslator`／`SetFont` | `game.syncLanguage()` 每幀比對，差了就 `LangPack().Apply(td)` |

⚠ **Drawer 不在 Session 手上**（它屬於 `mobile/wolong` 的平台殼），
所以面板點下去只換得掉訊息與人名。UI 詞與字型鏈要由殼層補掛——
`syncLanguage()` 每幀比對一次語系代號，差了才動，代價可以忽略。
少了這一步的症狀是「人名換了、選單沒換」，看起來像譯文沒做完。

## 6. 驗證（2026-08-26）

| 項目 | 結果 |
|---|---|
| `langpack_test.go` | 四個語系都載得起來；`Apply` 之後逐語系正確；切回母本**完全還原** |
| `internal/ui/phone/language_test.go` | `SetLanguage` 之後**整份武將名冊**與 TALK 探針都換了語言，切回母本一字不差；語言頁四列、目前那一列剛好一個記號 |
| `cmd/wlgame/launcher_test.go` | 語言頁進得去、游標停在目前語言、選得到、ESC 退得回；有存檔時 `LANGUAGE` 仍是最後一列（不擠掉 `LOAD DATA`）|
| F9 熱鍵 | 執行期切換與 `-lang` 啟動**逐像素相同**（f9l-1/2/3 vs ref-hans/ja/en 全 0 px）|
| 啟動殼層語言頁 | 實跑截圖，四個語言各自的寫法都畫得出來（`docs/playtest/46`）|

⚠ 探針不能只看一個名字：**日文版有 271/343 個名字與繁中寫法相同**（同樣是漢字），
挑到其中一個就會誤判成「沒換」。名冊整份串起來比對才有鑑別力。

## 7. 未解

| 缺口 | 下手點 |
|---|---|
| Android 實機／模擬器還沒實地切過語言 | 面板與 `syncLanguage` 都有單測，但手機版的畫面沒拍過；下一次 Android 打包驗收時補 |
| 語言不進存檔 | 原版存檔格式沒有這一欄，塞進去會破壞 round-trip。要記住偏好得另存 remake 自己的設定檔 |
| F9 是 remake 自創的鍵 | 原版沒有這個鍵；`docs/spec/13` 的按鍵表要同步記一筆 |
