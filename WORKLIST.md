# 臥龍傳 remake 工作清單

> ⚠ **待辦看 [`CONTEXT.md`](CONTEXT.md) §7，不是這裡。**
> 本檔是**按日期的完成紀錄**：每一節記下那一天封口了什麼、當時的邊界在哪。
> 每一輪都會更新的「現在該做什麼」在 `CONTEXT.md` §7.0
> （`CLAUDE.md` §10 也指向那裡）。
>
> 兩邊都寫「唯一來源」會讓接手的人拿到舊快照——本檔只在封口時補一節，
> 而 `CONTEXT.md` §7 每一輪都動。最新一節的日期就是本檔的時效。
>
> 自 2026-08-10 起不再建立或更新 `HANDOFF.md`；歷史交接內容已濃縮到本檔，
> 深層證據請回查 `RESEARCH-LOG.md`、`CONTEXT.md` 與 `docs/`。
>
> ⚠ **日期節裡的敘述是「當時的認知」**，其中有些後來被推翻
> （例如密碼頁一度被當成 oracle 的阻擋，2026-08-12 測出**不擋**）。
> 推翻紀錄集中在 `CONTEXT.md` §6。

## 目前目標

完成《臥龍傳－三國制霸之計》的證據導向 remake：核心規則、松崗 DOS/V
呈現、繁中 TALK、存檔與正常玩家路徑接到可重播的跨平台程式。

本輪以松崗繁中版為唯一畫面／行為 oracle（`workplace/orig/dosv` 是沿用的資料夾名稱）：

- 640×400 畫面外框與數值視窗以 DOS/V 對準。
- PC-98 與其他版本只保留歷史研究，不作本輪視覺、行為或 release 驗收基準。
- 依使用者要求不跑完整長程遊戲測試；採窄 fixture、單測、Docker/Xvfb 短 smoke。
- 三平台候選包可在短驗收 gate 通過後建立；原生 Windows／macOS runtime 仍是獨立 gate，不能由交叉編譯檔頭代替。

## 已完成且可回查

- 格式／時鐘／規則／狀態／存檔 round-trip 的主要垂直切片已接入；原始
  `workplace/orig/`、`workplace/ida/` 視為唯讀。
- `wlgame` 已接真實劇本、地圖、命令、軍團、行軍、戰術遭遇、四槽存檔 overlay、
  戰鬥指揮／委任與固定種子重播。
- 事件 1、2、3、4、5、8、9、11、12、13 的已證實狀態接縫已在 `internal/state`
  接入；事件 2／3／4／5 的玩家選單與數值編輯共用同一條 modal 狀態路徑。
- `sub_17C6E` 數值核心、DOS/V `CS:7D93h` 3×6 raw 格位、16×16 命中區與
  `AmountEdit` 已接入；DOS/V `sub_17D0D` 的 ICONGRF 第 3 段 96×64 內框也已
  接到 `(88,184)`，外圍保存區仍是 `(80,176)`／112×80。`KI.EXE` `seg002:031B`
  的兩層 16×16 硬體游標 mask 與 `ICONGRF` 下半部 3×6 靜態按鍵 glyph 已解碼、
  載入並接到面板；只保留資產缺失時的 fallback。
- TALK marker 展開、原始行硬斷行、五行／16 px 分頁、肖像與 IVENT 場景、
  pending 完成後的消像已接入；M7 校訂輸出由
  `translations/talk-dosv-corrected.json` 提供。
- `MMAP.MCH` 16×16 平面圖塊、`0xA000` metadata、`CS:985Ah` type 1／2／3
  八相位查表與火災／暴動圖像接到戰略地圖；runtime object record 已依
  `sub_123FF`／`sub_12459`／`sub_12533` 接入 state，包含 16-update timer、dirty
  與 render 後相位遞增。後半區 slots 16–31 的移動分支 `sub_1248A` 已依 raw fixed-point、
  方向 byte 與邊界 wrap 接入並有定向測試。
- Docker 內已有 Go test／vet、翻譯 selftest、文件索引與 Linux/Xvfb 截圖 smoke；
  完整長程遊戲測試依使用者要求略過。

### 2026-08-21（第十輪）刪除後複驗：反向的登記不實

刪完死程式碼之後再跑一次全部機械檢查。這一輪抓到的是**反向**的登記不實
——不是文件說未解而其實已解，而是**知識只活在註解裡，文件沒登記**：

| 註解裡有、文件沒有 | 處置 |
|---|---|
| `sub_1C775` ＝ `值 >> 2`、`sub_1C78E` ＝ `值 × 3 ÷ 4`，上限 `CH=0x7C` ＝ 124 | 算式其實在 `spec/31` §2.3 與 `re/60` §5，只是**沒綁上符號**，所以 `re/24` 還把這兩支列在**未讀函式目錄**裡。兩邊都補上 |
| `sub_1A3C3` 的開戰對白 pair ＝ TALK `0x1BA`／`0x1BB` | **任何文件都沒記。** 補進 `spec/60` §3.5，並標明「上格攻方／下格守方」只是強推論（照影格位置接的），列進該份的未解表 |

⭐ **`docs/re/24` 那種「未讀函式目錄」會反過來過期**：函式被讀懂了，
但解讀寫在別份文件（甚至只寫在註解裡），目錄沒人回頭改。
`CLAUDE.md` §10 叫人「動手讀任何一支之前先查 `re/24`」——
**目錄說未讀而其實已解，會讓人白讀一遍**。

順帶把最後 5 段英文註解翻成繁中（`minimap.go`／`launcher.go`／
`battlelayout.go`／`launcher_test.go`／`events.go`）。其中 `events.go` 那段
還補上出處：`sub_13771` 讀的軍團記錄在段內 `2240h`（`re/08` §4）。

未解列數 530 → **531**（`spec/60` 新增一列真缺口）。

⚠ **數值型的斷言複驗完是乾淨的**：8 筆命中裡 6 筆是誤報
（進位換算、同一支常式在不同呼叫點的不同常數），2 筆是上面那兩條沒登記。

### 2026-08-21（第九輪）重讀 internal 註解：1,248 行機器碼斷言

三個機械訊號跑完 `internal/` 的 93 個檔：

| 訊號 | 結果 |
|---|---|
| 註解裡的**數值** vs 提到該符號的文件全文 | 574 個符號比完只剩 2 筆，**兩筆都是誤報**（`sub_1351A` 那筆是 `cmp dl,18h` 與 `cmp ah,24h` 兩個不同的比較，文件寫得比註解精確）|
| 註解裡的 **Go 識別字**找不到宣告 | 278 筆命中幾乎全是 IDA 符號（誤報）；真缺陷只有一個：`tactical.go` 寫 `Tactical` 而型別叫 `TacticalSetup` |
| 註解引用的**小節**存不存在 | 六筆，見第八輪 |

⭐ **數值層面的斷言基本上是乾淨的。** 真正的問題是另一類：
**死程式碼與指錯的指標**。

舊的 `mobileui` 套件（當時在 `internal/ui/` 底下）**沒有任何 Go 檔 import 它**——
Android 綁定 `mobile/wolong` 走的是 `internal/ui/phone`。它是 2026-08-20
換路線前那條「照原版版面」留下的原型。連帶三處錯的登記：

| 位置 | 曾經寫過 | 實際 |
|---|---|---|
| `android-plan.md` §4 架構圖 | 呈現層畫成 `mobileui` | 是 `internal/ui/phone` |
| `phone/layout.go` 檔頭 | 「平台無關的幾何在 `mobileui`」 | 那個套件沒人用 |
| `phone/layout.go` 兩處 | 座標「由 `Viewport` 換算到實際螢幕」 | **`Viewport` 只定義在那個死套件裡**，live path 上根本沒有它——縮放是 Ebiten 的 `Layout` 契約做的（`game.Layout` 回傳 `phone.LogicalW/H`）|

⭐ **這一類比數值錯更難發現**：數值錯了會有人算不出來，而「指向一個
不存在於設計裡的套件」不會有任何症狀，只會讓讀的人多花半小時。

使用者裁定（2026-08-21）：**刪掉**。`mobileui` 那兩個檔（164 行）整個移除，
`android-plan.md` §4 改成「這條鏈上只有這三層」，不留「保留供參考」
那種會長回來的但書。

順手掃全庫的其他死程式碼：**沒有第二個沒人 import 的套件**；
函式層級找到 8 支零呼叫，刪掉其中 7 支
（`amountButtonIsEdit`／`amountPanelValues`／`advName`／`plainErr`／
`fundingOfficerName`／`fundingSubjectName`／`drawDisplayEntry`，共 85 行）。

⚠ **第 8 支 `march.Swap` 不能刪**——它是 `heap.Interface` 的一員，
靠介面派發，靜態掃描看不到呼叫端。**「零呼叫」在 Go 裡不等於死碼。**

⚠ **刪 `drawDisplayEntry` 之後差點連帶刪錯**：它是 `battleDisplayEntry.kind`
唯一的**產品碼**讀者，所以 `kind` 與那個 enum 看起來也死了。實際上
`isoview_test.go` 用它過濾 raw-unit entry，釘的是 `sub_1DA1C` 的奇偶配對。
⭐ **測試是合法的消費端**——刪掉欄位會讓那條不變量無法表達。
（順帶確認繪圖層不看 `kind` 是**對的**：live path 是 `drawDisplayGrid`，
照原版 `sub_1DDB4` 統一用圖號畫，不分型別。）

順帶收掉術語漂移（`CONTEXT.md` §5 明寫「不寫城市」、能力值是「武術」）：
註解 24 行的「城市」→「據點」、9 行的「武力」→「武術」，
文件 8 處欄位標籤與機制敘述一併改。**真實地理敘述留著**
（「河北一帶的城市」講的是真的城市）。

⚠ 這三個訊號只驗得了「數值兜不兜得攏」「識別字在不在」「引用指不指得到」。
**註解裡的機制敘述有沒有講錯，只能逐支對機器碼重讀**，沒有做。

### 2026-08-21（第八輪）斷言檢查擴到原始碼註解

六輪文件稽核都沒抓到 `tools/dosboxx.sh` 的檔頭，因為 `tools/index.py` 的
斷言檢查**只走 `docs/` 與根目錄的 markdown**。那一行寫著

> 為什麼跑 PC-98 而不是松崗版：**PC-98 原版沒有防拷**

而 `CLAUDE.md` §4.0 明文寫著「選 PC-98 仍然可以，但理由只能是
『那一版的腳本比較齊』，不能是『過不去』」。

⭐ **註解裡的斷言比文件裡的更難發現**：沒有人會為了查一條結論去 grep
建置腳本的檔頭，而**照著它做的人會直接改用另一個工具**——傷害發生在
「換工具」那一刻，不會留下任何痕跡。

改掉的註解：

| 檔案 | 曾經寫過 | 實際 |
|---|---|---|
| `tools/dosboxx.sh` | 「PC-98 原版沒有防拷…不必去過密碼那一關」 | 理由改成腳本比較齊，並明寫要驗 DOS/V 就用 `dosv_live_capture.sh` |
| `path.go` | `FindPathForcing`：「那個演算法還沒解出來，解出來之後這一支要換掉」 | **同一個檔的檔頭就是 `loc_1BD46` 的移植。** 演算法不是近似，`force` 參數才是 remake 加的（理由：驗收路徑要能終止；原版不需要，破城本來就是玩家的工作）|
| `isoview.go` | 「原版是否在城壁破壞後做局部更新**尚未知**」 | 原版會更新（`sub_1B824` → `sub_1BB6D`，`spec/66`）。沒更新的是**這張縮圖**，是明知的 remake 差異 |
| `wlview/main.go` | 「原版的字型來源見 CLAUDE.md §3.6，**還沒結案**」 | **兩層都錯**：`CLAUDE.md` 沒有 §3.6，而字型懸案在 `re/29` 就結案了。那一支只是還沒接倚天字，不是缺資料 |

六筆註解引用指向不存在的小節（`re/09 §4.4`、`formats/08 §1.4`、
`spec/61 §5.1` 兩處、`formats/05 §6`、`CLAUDE.md §3.10`），全部修正。

三道新檢查，都做過正對照：

1. ③.5 的掃描範圍加上 `.go`／`.py`／`.sh`
2. ⑨ 註解裡的 `docs/xx/NN §M` 要指得到那個小節——⑧ 驗的是 markdown 的
   連結形式，註解寫的是**裸路徑**，形狀不同所以比不到
3. ⑨ 同時涵蓋 `CLAUDE.md §N`／`CONTEXT.md §N`：那兩份的章節會重編，
   而註解寫下的編號不會跟著動

⚠ 只掃了「已知會復發的斷言」與「引用指不指得到」。**註解裡還有沒有
別的過期事實，沒有系統性的辦法查**——那要逐支對機器碼重讀。

### 2026-08-21（第七輪）掃完最後四個目錄：稽核收尾

`docs/reference/`（17）、`docs/release/`（13）、`docs/mobile/`（12）、
`docs/promo/`（7）。這四個目錄的缺口大半是真的——**沒有硬體就是沒有硬體**，
Windows／macOS 實機、Android 實機、release signing、16 KB page size 裝置
這四條至今成立，一條都沒動。

改掉的五處：

| 文件 | 曾經寫過 | 實際 |
|---|---|---|
| `04-first-survey.md` | 「即時制的時間單位：一個月在真實時間是多久」 | 計時中斷是 291.3 Hz，一個遊戲日在五檔分別 0.74／1.48／2.22／2.97 秒（`re/61`、`spec/34`）|
| `yt-remake-pixel-review.md` | 「中央 raw reserve glyph 未解出原版圖形，remake 不冒充，改用自繪」 | 是 `ICONGRF` 段 3 `0x1BA0` 起的四張 24×16（天秤／馬／弓／步），`DOSVResourceIcon` 用的就是原版素材；那一區後來逐像素 PASS |
| `dosv-adlib-and-tactical-review.md` | 整節在說「下一個有效驗收是建立同攻城節點、同攻守方、同編成、同鏡頭的 capture pair」 | **那正是 `playtest/40` 做的事。** 改寫成「推廣片證不了什麼」——那是媒材的性質，不是缺口 |
| `android-ux.md` | 數值輸入器引 `docs/spec/12` | `spec/12` 是主畫面外框；數值輸入器是 `sub_17C6E`，在 `re/13` |

⭐ **`--strict` 在這裡抓到一個我自己製造的問題**：把兩份 promo 的未解列
清空之後，它們變成「提到未解卻抽不出任何一條」的盲區。
**清掉最後一條缺口時要順手宣告「沒有缺口」**（`<!-- 缺口：無 -->`），
否則下一輪會分不出「沒有缺口」與「表壞了」。

未解列數：535 → **530**（promo 7 → 2）。全六輪合計 **570 → 530**。

### 2026-08-21（第六輪）稽核 docs/playtest/：三份文件的「受阻」都已不成立

`docs/playtest/` 是**有日期的實驗紀錄**，「當時的認知」本來就該保留，
所以校正只動**宣稱現況的句子**（狀態行、未解表、受阻標記），敘述不動。

三份文件的狀態行說某件事做不到，而那件事後來都做到了：

| 文件 | 曾經寫過 | 實際 |
|---|---|---|
| `13` | 「嚴格同狀態逐像素差異仍不宣稱」 | `37`–`40` 全做了。**Windows／macOS 原生 runtime 那一條至今成立**，留著 |
| `20`／`21` | PARTIAL，「同狀態戰場」「一覽詳細層與捲軸」未完成 | 都收掉了。⭐ 更重要的是**量法已被取代**——那兩份拿的是壓縮影片的代表幀、縮到 320×200 比幾何 |
| `28` | 「攻方仍然攻不進城…**實作還沒做**」 | 登城機制同一天就解完並實作（`ground.go` ＋ 三張攻城圖的迴歸測試）|
| `30` | 「攻城仍打不進去」，未解列「攻方大多數不前進」 | **那不是缺陷**——說明書第 11 章整章在講破城要換陣形，是玩家要操作的 |
| `31` | 「戰場：同一場的逐格對拍｜沒做過」、「大地圖地形色調｜存疑」 | 前者 `40` 做了；後者是調色盤刻度差 4%（`spec/51`），改完 `map` 區逐像素相同 |
| `37` | 「橫幅數字的字模」「鏡頭水平參數 20 vs 16」 | 兩條都解了（`spec/52`；那個「謎」根本不存在，是地圖解碼少跳 4 byte）|

⭐ **`playtest/15` 與 `14` 又寫了一次密碼頁**——這次的形狀是
「不宣稱**密碼保護下**的同狀態逐像素 parity」。**它把防拷擺成修飾語而不是
主張**，所以既有的兩條檢查（動詞式「擋住／阻擋」、名詞式「不可證實／邊界」）
都比不到。`tools/index.py` 補上第三種寫法後兩處都被抓出來。
這是 `CLAUDE.md` §4.0 那條斷言的**第六、第七次**復發。

⚠ 修的時候自己也踩了一次：我寫的更正句「不是被密碼頁擋著」被檢查誤判，
因為 `NOT_BLOCKED` 認得「不擋」不認得「不是被…擋著」。
**改成規則認得的措辭**，沒有去放寬否定式的比對——放寬會讓真陽性漏掉。

未解列數：546 → **535**（playtest 69 → 58）。

⚠ `docs/re/` 剩下的未解列仍未逐列核；`docs/reference/`（17）、
`docs/release/`（13）、`docs/mobile/`（12）、`docs/promo/`（7）還沒碰。

### 2026-08-21（第五輪）稽核擴到 docs/re/：加兩道跨文件檢查

`docs/re/`（70 份、229 列）沒辦法逐列人工核，改用兩個機械訊號篩：
**未解表裡自己就寫著「已解／已讀／✅」的列**（54 列命中）、
**未解列點名的 `sub_XXXXX` 已經出現在 Go 原始碼註解裡**（59 列命中，
多數是真缺口——被移植不等於那一列的問題有答案）。

⭐ **`docs/re/11` §6 那張表 24 列裡 22 列是已解的刪除線列**，剩下兩列還
被同一張表的其他列推翻（「疲勞、士氣存在哪」——同表上面就寫著 `+0x19`
與 `+0x03`）。清掉之後真缺口剩三條。全庫共刪 17 列刪除線殘留。

改掉的還有：`re/62` 的視野框美術（`spec/55` 已 CONFORMED）、
`re/70` 的第一幕幾何（`formats/09` §2.1 已解）、
`re/06` 的 `sub_13E11`／`sub_15358` 九支子程式／勢力記錄欄位表、
`re/08` 的連結記錄佈局與「軍團 `+0x0A` ＝ 行進方向」（那是**錯的**，
`+0x0A` 是步進量、`+0x08` 才是朝向）。

**兩個結構性缺陷**，各補一道檢查：

| 缺陷 | 檢查 | 正對照 |
|---|---|---|
| `docs/spec/00-index.md` 的狀態欄與規格自己的狀態行相反（`29-audio` 記成 DRAFT、`51` 記成 READY）| `tools/index.py` ⑦：索引那一欄若出現狀態碼，必須與規格的狀態行相同 | 注入 `29-audio` ＝ DRAFT 會被擋，還原後通過 |
| 七筆引用把**行號**寫成小節號（`§1065`、`§590`、`§467`、`§150`、`§91`），而它們全部指得到檔案，所以既有的連結檢查一路放行 | `tools/index.py` ⑧：連結後面接的 `§N` 要真的存在於被連到的那份文件 | 注入 `§590` 會被擋，還原後通過 |

⭐ **`docs/re/07` 的章節編號重複**（§18／§19／§22 各兩次），所以
「`re/07` §19」這種引用是歧義的——而我這一輪自己就寫了三個。
重新編號成 1–27 之後才修得掉。**長文件的章節編號要當成識別碼維護。**

未解列數：570 → **546**（re 230→223、mechanics 46→42、formats 40→35）。

⚠ 剩下的 `docs/re/` 未解列**沒有逐列核過**，只跑了上面兩個訊號；
`docs/playtest/`（69 列）整個目錄還沒碰。

### 2026-08-21（第四輪）稽核 mechanics 與 formats：十三條「還沒解」其實已解

逐條對 code 與其他文件核實 `docs/mechanics/`（8 份）與 `docs/formats/`（9 份）。
形狀與第三輪相同：做完之後沒有回頭改登記。

已解卻還掛在「還沒解」的十三條，比較值得記的幾個：

| 曾經寫過 | 實際 |
|---|---|
| `30-combat`「繞路點是誰算出來的（真正的尋路演算法）」 | `loc_1BD46`，**uniform-cost 波前擴散**不是 BFS（每格加的是佔用成本），回溯只在轉彎記點、上限 64 個。`re/11` §5.15 早就解了，remake 也照著實作在 `internal/rules/tactical/path.go` |
| `30-combat`「士氣值存在哪」 | 軍團記錄 `+0x06`，編成時從勢力 `+0x1D` 複製（`re/08` §4，confirmed）|
| `40-economy`「預備兵維持費的單價（三兵種是否不同）」 | **三兵種同價**：`(騎馬 + 弓兵 + 步兵) ÷ 32`，先加總再除（`re/08` §2）|
| `40-economy`「北方／南方的判定邊界」 | `y < 80` ／ `y ≥ 150`（`re/07` §6）|
| `50-diplomacy`「外交官每月提升的量與金額的關係」 | `re/08` §3 全解：每小時 12.5% 機率動作、經費扣 `23 − 政治`、`rand(0..15) ≤ 政治` 才有成果、交友度 +1 上限 100 |
| `formats/05`「`MMAP.MAP` 的編碼」 | `formats/06` 整份就是它——**而 `formats/05` 自己的狀態行也這麼寫**，只有 §4 沒改 |
| `formats/07`「`64 × 62` 哪一邊是寬」 | 寬是 64。**同一份文件的下一小節** `di = Y × 64 + X` 就是答案 |
| `re/04`「`sub_20000` 的三次呼叫在設定什麼」 | `int 33h` 滑鼠驅動的分派表，設的是游標範圍 `0x17FF × 0x101F`；`playtest/38` §1 的座標換算靠的正是它 |

⭐ **反向的一條**：`CLAUDE.md` M2 寫「`TALK.DAT` 變數插入語意**全解**」，
而 `formats/01` §3 明寫 `\1`–`\4` 的文字語意「仍保留為未定案」。
查 code 證實 `cmd/wlgame` 用的是帶通用字樣的推定值（`'1': "武將"`、
`'2': "據點"`…）。**畫得出來不是證據**——里程碑表因此改成「機制全解，
文字語意是實務推定」。登記不實會往兩個方向跑，只查「說太少」會漏掉這種。

⭐ 這一輪多出一種形狀：**同一份文件內部自相矛盾**。既有的檢查 ② 比的是
「狀態行 vs 內文有沒有 confirmed 斷言」，比不到「未解表 vs 正文」。

未解列數：562 → 552（mechanics 46→42、formats 40→35）。

### 2026-08-21（第三輪）文件狀態稽核：六份「未收尾規格」裡四份是登記不實

`CONTEXT.md` §7.0 把「六份 READY 規格收尾」列成第一順位。逐份對 code 核實之後，
**真缺口只有兩件**：`spec/32` 的戰場區右鍵熱區層（`cmd/wlgame/battle.go` 完全沒有
`MouseButtonRight` 分支）、`spec/35` 的 22 勢力選擇視窗版面與熱區 `0x16`。
另外四份的缺口早就被填掉了，只是登記沒跟著改。

⭐ **源頭是 `docs/spec/00-index.md` 的「狀態」欄**。那一欄是散文不是狀態碼，
所以 `tools/index.py` 的既有檢查看不到它——它擋的是「一份文件的狀態行與自己的內文
矛盾」，擋不到「索引對某份規格的描述與那份規格自己的狀態行矛盾」。
逐列比對之後改掉十列，其中兩列連狀態碼都相反：`29-audio` 索引寫「**DRAFT**：
播放層未做」而規格是 CONFORMED、音效整條接通；`51-vga-dac` 索引寫「尚未全面套用」
而換算早就在 `internal/assets/palette` 的 `toSRGB`，主畫面因此逐像素相同。
其餘八列是「已解掉但索引還留著上一版的但書」：`22` 的頭像與滑鼠、`33` 的命令圖示
來源段、`36` 的攻方不前進、`38` 的捲軸、`30` 的結局過場、`11` 的尚未實作、
`12` 的內部排版估值、`13` 的原版執行期未驗。

改掉的斷言（全部改寫成現況，不留刪除線，推翻紀錄集中在 `CONTEXT.md` §6.1）：

| 文件 | 曾經寫過 | 實際 |
|---|---|---|
| `spec/12` §7 | 信賴度是數字不是量條 | 早就是量條，照 `sub_10AAA` 的 `(信賴度×100 + 0x9F) ÷ 0xA0`，槽 176×10 來自顯示清單 op 03 |
| `spec/12` §7 | 顯示清單「只讀了場景 0（`docs/re/48` §5）」 | `docs/re/48` 狀態行寫「十個場景的歸屬與九個 opcode 全部解出來了」；而它的 §5 是〈為什麼這件事重要〉，根本沒有那句話 |
| `spec/12` §6 | 「逐像素沒比，而模擬器上主畫面收不到點擊」 | 兩半都被推翻：`playtest/37`／`38`／`39` 逐像素對過，`playtest/38` §1 證明點得到 |
| `spec/35` §5 | 視野框的美術尺寸沒從程式碼讀到 | `spec/55` 已 CONFORMED，點陣從檔案解出來、位置在兩個鏡頭各驗一次 |
| `spec/91` | 狀態 DRAFT、§1「戰場一次都沒有對過同一場」 | `playtest/40`（2026-08-18）九區裡六區逐像素相同、`field` 0.17%；升 CONFORMED |
| `re/46` §3 | 「剩下的是視窗內部那張龍紋底圖」 | 龍紋在 `ICONGRF.DAT` 最後 128 byte，早就解出來也接上了（`formats/03` §5.5）|
| `README` | 規格 56 份／50 CONFORMED／5 READY／1 DRAFT | **59 份／56 CONFORMED／3 READY／0 DRAFT** |
| `README` | 全專案未解 553 列（spec 132）| 清掉過期列之後 **562 列**（spec 128）|
| `CLAUDE.md` §8 | M4「規格已到 `docs/spec/65`」 | 已到 `71`（另有對拍規格 `90`／`91`）|
| `CLAUDE.md` §8 | M7「排版 parity 未完成」 | `playtest/32` 已全量量過：1,022 則逐則量進訊息框，單行超寬 0 行 |

⭐ **共同的形狀：四份都不是推理錯，是做完之後沒有回頭改登記。**
所以判斷「某件事解了沒」要看**產物**（code、測試、逐像素結果），
不要看待辦表替它寫的一句話摘要——摘要只反映寫下那一刻，
而且它不會因為它引用的規格被改寫而自動更新。

⚠ **這一輪只稽核了 `CONTEXT.md` §7.0 點名的六份規格與 `docs/spec/00-index.md` 全表。**
`docs/re/`、`docs/mechanics/`、`docs/formats/` 的未解表沒有逐列對 code 核過；
以這一輪 6 份裡 4 份、索引 59 列裡 10 列的比例推，那些目錄裡應該還有同類的過期列。

### 2026-08-21（下半）16 KB 對齊、實跑推廣片、`20260821` 批次

⭐ **APK 的 `.so` 原本是 4 KB 對齊。** Android 15 起有 16 KB page size 的裝置，
Go 產出的 `libgojni.so` 預設 LOAD 段 `align=0x1000`，那種 `.so` 在 16 KB 的機器上
**載不起來**——而 4 KB 的機器上完全正常，所以測不出來。
`zipalign -c -P 16` 早就通過了，**它驗的是 zip 那一層，驗不到 ELF 這一層**。
修法是 `ebitenmobile bind` 帶 `-ldflags "-extldflags=-Wl,-z,max-page-size=16384"`，
並在 `tools/android_build.sh` 建完之後用 `readelf` 逐段檢查，不是 `0x4000` 就讓建置失敗
（[`docs/release/03`](docs/release/03-three-platform-20260821.md) §4）。

⭐ **推廣主片以前每一段都是靜止截圖。** 畫面是 remake 真的畫的，但影片裡沒有一格在動
（量得到：`freezedetect` 在 60 秒裡標出約 52 秒凍結）。現在大地圖、野戰與攻城
三段改成**逐幀錄下來的實跑畫面**（[`docs/spec/71`](docs/spec/71-promo-live-capture.md)）：
`cmd/wlgame` 加 `-frames-dir`／`-frames`，程式自己寫 PNG，不做螢幕擷取。
錄的時候踩到三個都長得像「錄成功了」的坑：戰場開場那一幀是雙方對白、兩軍還沒接觸
（240 張裡只有 4 張不同），要先 `-battle-steps 200`；原版初值的戰術鏡頭 `36,14` 對著城牆、
部隊在畫面外，要 `-battle-cam 20,15`；大地圖低速檔八秒只走幾天，要 `-speed 0`。
`tools/promo_live_capture.sh` 逐段數「不重複的張數」，低於 20 就擋下來。

**發行盤點的三個發現**：`release-manifest.json` 的 `android_experimental` 指向一個
**不存在的檔名**（手動換過 APK，清單沒跟著換）；`PROMO_FILES` 漏了 Android 推廣片；
`ANDROID-EXPERIMENTAL.md` 還寫著「觸控 shell 原型、尚未接入遊戲核心」，
而那份斷言在 2026-08-20 就不成立了。三處都改掉，APK 命名也從
`touch-prototype` 改成 `android`。

`dist-all` 重打成一致的 `wolong-remake-20260821` 批次：五個桌面包 ＋ Android APK，
`sha256sum -c` 二十筆全部相符，Linux GUI smoke 截圖補回去
（`promote` 每次都會清掉 `verification/`，要另外跑 `tools/release_smoke.sh` 再 `refresh`）。

### 2026-08-21 手機版套原版配色 ＋ 推廣片 ＋ 側欄熱區勘誤

**手機版的按鈕與面板改用原版的底色與外框**（[`docs/spec/70`](docs/spec/70-phone-chrome.md)）：
顏色查 `GAMEPAL.BRG`、外框是 `ICONGRF.DAT` 的 8×8 圖塊，兩者都走
`internal/ui/chrome`——**與桌面版同一份**。對應關係照原版的視窗種類：
狀態列與指令列是命令列（色 0 純黑、沒有龍紋）、一覽與選對象是清單視窗
（米色底黑字）、小卡與對白與擋住世界的決定是選單／情報視窗（深藍 ＋ 龍紋、
白字）、選中的一列是反白條（色 5）。⛔ 手機層不再抄任何 RGB，
`TestPhoneUsesOriginalPalette` 守著這一條。

兩處刻意不照抄，標成 remake 差異：地圖上選中的據點畫**白圈**（反白條的綠
疊在土黃城上分不出來）、戰場側欄的「我方」用白字（用綠會讀成「選中了它」）。

**Android 版推廣片**（[`docs/promo/android.md`](docs/promo/android.md)）：
48 秒、1280×720，全部是手機版自己的畫面 ＋ 原創合成配樂。
⭐ **逐幀輸出而不是錄螢幕**——X11 擷取的幀率不穩還會抖，逐幀是確定性的。

⭐ **側欄熱區勘誤**：六個命令的熱區照原版是 **48 px 寬**，不是面板那張圖的
128（`docs/re/60` §6）。算成整條寬的話點陣線的右半也會送出命令。
舊測試把 128 釘住了，一併改掉並記進 `CONTEXT.md` §6。

`docs/images/` 的 remake 截圖全部重拍（`tools/readme_shots.sh`）：
`wlgame-battle.png` 那張是 2026-08-17 15:10 截的，而側欄那個突出到陣形的
選取框同一天 19:04 就拿掉了——**圖比程式舊四小時，而過期的截圖看起來
跟正確的一模一樣**。

### 2026-08-20 Android 手機版：核心接入、四個入口、戰場與資料匯入

**里程碑 A 通過**：模擬器與桌面在同一組幀算出**完全相同的指紋**
（frame 1／60／120 ＝ `2b58e7b58f4796f9`／`5b3585cf005a6ec5`／`36eb02d379333f95`，
`World.Fingerprint`，[`docs/spec/69`](docs/spec/69-world-fingerprint.md)）。
這一條不看畫面就能驗，是整條路線最強的驗收。

手機呈現層在 `internal/ui/phone`，**只共用規則層、畫面重畫**
（[`docs/mobile/android-ux.md`](docs/mobile/android-ux.md)）：

- 主畫面：大地圖（可縮放拖曳）、狀態列、四個入口、據點小卡
- 進言五項齊全：三個外交提案走完整說服迴圈，遷都與請求出陣走一次驗收
- 一覽四張表（唯讀）、軍團（現有 ＋ 編成）、系統（速度／四槽存檔／關於）
- 戰場：45 度視角、六個編成位置、六個命令、兩軍資訊與縮圖
- 資料匯入：`ImportActivity` 是啟動入口，走系統的 SAF 選擇器

**七塊東西從 `cmd/wlgame` 抽成共用套件**（`internal/ui/isoview`、
`internal/battlesetup`、`internal/rules/speed`、`persuasion/talk.go`、
`state/persuade.go`、`assets/battle/order.go`、`text.Table.Lines` ＋
`ui/talkmenu`）。抽完桌面版的戰場畫面**byte-for-byte 沒有變**
（同一個固定局面截圖雜湊 `a450b80e…`）。

⭐ **順手抓到一個存檔的真缺陷**：`savefile.Encode` 存的是**載入當時**的
區塊，所以開檔之後的進度不會進原生存檔，而讀檔又優先讀原生檔——
存了等於沒存。既有的三個測試從來沒有推進過時鐘，所以完全看不出來。
（推翻紀錄見 `CONTEXT.md` §6。）

**還沒完成**：SAF 選資料夾之後的複製流程沒有自動驗（要驅動系統的檔案
選擇器）；實機驗收 ⛔ 沒有裝置。工具鏈與踩過的坑寫在
[`docs/mobile/android-plan.md`](docs/mobile/android-plan.md) §6。

### 2026-08-20 數量斷言複驗 ＋ 三平台重新交付 ＋ Android 路線定案

**數量斷言複驗**（使用者要求）：README 寫的「反組譯 549 條未解」**兩層都錯**。

| 錯的 | 實際 |
|---|---|
| 標籤：掛在「反組譯」那一列 | `docs/re/43` 彙總的是**全部文件**；`docs/re/` 自己只有 **230 列** |
| 數字：549 當成獨立問題數 | 那是**列數**。全專案 **553 列**，其中 9 列是索引檔對別份文件的重述 |
| 規格份數 58、DRAFT 2 | **56 份**（`TEMPLATE.md` 不是規格）、DRAFT **1** |
| 戰場 307 px 裡「211 px 消不掉」 | **299 px**（＋門破的 tick 88）；真正未歸類只剩 **8 px** |

⭐ **抽樣 11 列就抓到兩類誤判**：一類是索引在描述 `43` 這份文件本身、
只因為句子裡有「未解」兩個字；另一類是**自己剛解掉卻沒清的過期列**——
`re/59` 的結局訊息與 `D7END.EXE`、`spec/30` 的結局畫面、
`mechanics/80` 的 `END_S*` 格式，四處全部改寫。
**數字對不代表內容對**：`grep -c` 與分類小計完全吻合，誤判照樣在裡面。

`tools/re_open_questions.py` 的 §1 現在自己印出「這是列數不是獨立問題數」、
索引重述的列數與逐目錄分佈，讓引用的人看得到要扣什麼。

⚠ **只抽了約 2%**，沒有逐列核過；以那個比例外推，剩下的列裡可能還有十幾列類似。

**三平台重新交付**：`dist-all` 換成一致的 `wolong-remake-20260820` 批次
（`docs/release/02-three-platform-20260820.md`）。順帶修掉三個發行流程問題：
`release_all.sh` **沒有執行權限**（直接 exit 126）、版本字串硬寫在十幾處、
`stage` 只在 `dist/promo/` 找推廣片而影片只剩 `dist-all/promo/`——
**而且要編完三平台才會發現**。Android 附件改用它自己的建置日期命名，
因為它根本沒有重編，掛發行日期會宣稱一個沒發生過的建置。

**Android 路線定案**（使用者裁定）：只共用規則層、UI 重畫；驗證只有 Docker 模擬器。
量出決定路線的事實：`cmd/wlgame` 是 **16,587 行、53 檔、253 個方法的 `package main`**，
手機綁定 import 不到；重新設計那條路直接接 `internal/` 就好。
⚠ 另外找出一個必須先解的：`internal/*` 全用 `os.ReadFile(路徑)`，
而 Android 11+ 的 SAF 給的是 `content://`——要有 Java 側的匯入步驟。

### 2026-08-18（下半）結局、過場格式與倒地動畫

- [x] **結局全文救回來了**：燒在 `D7END.EXE` 的資料段（段內 `0x5F0`），
  松崗版 200 字 10 行、日文原版 180 字 9 行，兩版斷句不同。`docs/re/70`
- [x] **過場圖的格式**：外層就是 `MMAP.MAP` 那一種 RLE。先前「位移比檔案大」
  被讀成格式怪，而那正是「這是壓縮檔」的證據。內層 640×400 四平面、
  上下兩半各 200 列、一半之內平面優先，一幕配一組色盤。`docs/formats/09`
- [x] **結局會播了**：十二幕 ＋ 逐字文字 ＋ 十七階淡入淡出。`docs/spec/67`
- [x] **倒地動畫**：四幀，圖號 84 ＋ 側 ×90 ＋ 兵種組 ×2 ＋ 後兩幀 +1；
  大將與騎馬共用同一組。那四幀不擋路也不算場上人數。`docs/spec/68`
- [x] **打壞的城壁與門會在畫面上換掉**：繪圖層原本是進戰場那一幀的靜態副本。`docs/spec/66`
- [x] 戰場那 88 px 是**一道破掉的門**，不是「沒畫的木樁」——兩邊的門破在不同 tick。
- [x] 六個編成位置的效果：**位置本身沒有加成**，查表索引是（兵種, 場合）。
- [x] 六階外交的門檻：`sub_17A7A` 的 `dh = 14h`，**一級 20 分**，換算式收斂成一份。
- [x] 「攻方停止打城壁」不是規則層缺口：被擋的格子上**沒有實體**（是抬高的地形）。

### 2026-08-18 戰場同狀態對拍收斂 ＋ 兩個規則層缺陷

**對拍**（`docs/playtest/40`）：戰場九區裡**六區逐像素相同**
（底列、標題、我方將旗、陣形列、指令面板、`▶▶` 列），
戰場區從 19.9% 一路降到 **0.17%**（307 px），小地圖 0.04%。

| 接上的成因 | 規格 |
|---|---|
| 一帶是 **8 列**不是 16（`sub_1E0E1` 的 `mov cx, 8`）| `docs/spec/58` |
| 開場常令兩側都是「命令 0」；腳本的陣形線查表用**角色**不是側編號 | `docs/spec/59` |
| 對白框壽命 60 tick，而且**每側各一個**、兩邊可同時掛著 | `docs/spec/60` |
| 底列一格有**四樣**東西（少畫了兵種圖示；命令圖示是下令時畫一次就不更新，就位 7 要顯示成陣形 0）| `docs/spec/33` §1.6 |
| 兵的**開場體力 ＝ 軍團士氣**，`MaxHP` 100 是回復上限 | `docs/spec/61` |
| **旗也要跟著戰場翻**（`sub_19E10` 掃的是翻好之後的緩衝區）| `docs/spec/56` §3 |
| 顯示格的**圖號 0 ＝ 空**，不是第 0 張圖 | `docs/playtest/40` §12 |

**規則層的兩個缺陷**（不是畫面問題）：

- [x] **被換位的兵這一幀不動**（`sub_1ADC8` 的 `test al, 40h`）。少了它，
  兵撐得夠久就會擠在一起互相打不到——量到 48 個兵圍著 2 個打、
  95 萬 tick 零傷害。補上之後同一場 5,000 tick 內分勝負。`docs/spec/62`
- [x] **挨打有三幀硬直**，而且受擊旗標撐到硬直結束才清（它同時是換位的
  擋條件）。`docs/spec/63`
- [x] **退卻是保命不是損失**：兵離場只有一支常式，`ah=0`（退到畫面外）
  算生還、`ah=1`（倒地）不算，打完的兵力是 Σ（存活 ＋ 待機）。
  remake 兩種離場都當成戰死，打完的兵力偏低。`docs/spec/65`
- [x] **遷都的兩條播報**：自國君主下令（#518–520）；他國要**有外交官駐在
  那裡**才報得回來（#57 ＋ #521–525），沒有就一句話都不報。`docs/spec/64`

**RE 覆蓋**：四個分級全部收斂到 T1，739 支函式每一支都有 `docs/re/` 筆記
（最後兩批 `docs/re/68`／`69`）。⚠ T2 那批多半只是登記不是新理解。

**量測方法本身**（`docs/spec/90` §4.1）：參考影格是實錄，**本身會有東西**
——307 px 裡 **299 px 消不掉**：滑鼠游標 95、旗的揮舞相位 116（兩邊各自
`rand & 3` 起手）、兩邊的門破在不同 tick 88。**真正未歸類的只剩 8 px。**
到 0.2% 這個量級，「原版有、remake 沒有」不自動等於 remake 的缺陷。

**文件稽核**：`tools/index.py` 的「密碼頁擋住 oracle」檢查只走 `docs/`，
根目錄的 markdown 從來沒被掃過，於是本檔與 `VERIFICATION-MATRIX.md`、
`RESEARCH-LOG.md` 都還留著那條被推翻的斷言。檢查已擴到根目錄，
正規表示式也補上「密碼保護…不可證實／邊界」這種換句話說的寫法。

### 2026-08-12 `dist-all` 統一交付封口

- [x] 新增可重播的 `tools/release_all.sh`／`tools/release_all_fs.py`：所有建置、封裝、
  雜湊與 AppImage 步驟都在 Docker 內執行，最終只保留 `dist-all/` 的交付檔；原始
  資料、完整 TALK 表與 `.work` 建置中間檔均被排除。
- [x] `dist-all/packages/` 已集中 Linux amd64、Windows amd64、macOS Intel＋Apple Silicon
  三個桌面完整包及 Linux amd64 AppImage；Linux arm64 無頭工具包另列為伴隨包，不假稱完整
  GUI 平台。
- [x] 四支推廣片已複製至 `dist-all/promo/`。其中
  `wolong-remake-dosv-live-comparison.mp4` 是 60 秒、1280×720 的「讓經典再現」
  DOS/V／remake 實機比較片；YouTube 代表幀短片是 24 秒、1280×400，兩者不可混稱為
  同狀態逐像素 parity。
- [x] Linux 解壓包與 `APPIMAGE_EXTRACT_AND_RUN=1` AppImage 都在 Docker/Xvfb 載入玩家
  自備松崗資料與倚天字型，從同包公開 `corrections.json` 成功啟動並得到同一張 640×400
  固定種子截圖（SHA-256 `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`）。
- [x] 公開校訂覆蓋在真實松崗 `TALK.DAT` 上逐則與既有 1,022 則產生表位元組一致；完整表
  僅作測試 oracle，不進發行包。`denylist.py dist-all` 通過，根目錄 `SHA256SUMS.txt` 已驗證。
- [x] Android APK 收在 `dist-all/experimental/android/`，並清楚標為觸控 shell 原型，不計入
  三平台完整遊戲發行。
- [ ] Windows／macOS 尚缺目標作業系統的原生 GUI／輸入／音訊／字型短 smoke；交叉 ABI
  檔頭已驗證，不能取代該 gate。

## 本輪高優先工作

### P0：DOS/V 畫面外框與數值選取（內框已完成）

1. 已由 DOS/V `sub_17D0D` 接入 `ICONGRF` 第 3 段相對 `0x14A0` 的 96×64
   平面內框，目的座標為 `(88,184)`；`sub_19796` 保存區與 `(88,200)` 格位仍分離。
2. 實際數值選取已使用 `CS:7D93h` raw byte 與同一套 `AmountEdit`；DOS/V
   `KI.EXE` cursor mask（file `0x1051B`）與 `ICONGRF` 3×6 靜態 glyph 已由原始
   bytes 解碼並接線，`internal/assets/gfx` 單測固定像素統計與每格 glyph。
3. DOS/V 執行檔仍有複製保護，但使用者提供的 YouTube 錄製
   [`af6xqcicXoI`](https://www.youtube.com/watch?v=af6xqcicXoI) 已提供自然遊戲 oracle；
   原版 DOSBox 密碼頁只保留作啟動邊界，不再阻擋自然畫面結構對照。

4. `cmd/wlgame` 有真實 96×64 資源時不再疊畫 vector 矩形／CJK 按鈕標籤；解出的
   16×16 白框／紅填 cursor 只在數值面板選取狀態繪製。自然策略 HUD 已依影片／說明書
   對齊 32 px 命令列、左側 27×21 地圖與右側 208 px minimap／情報欄。

### P0：事件 6／7 次要 TALK（raw formatter 已接線）

1. 已追完 `sub_13327`、`sub_13388` → `sub_13C3D`：事件 6 是 `CX=0x48`／#72，
   事件 7 是 `CX=0x4C`／#76；條件是 `AH != 0` 且 response 不是 2／3。
2. `TalkNotice` 已保留 `Secondary`／`NoPortrait`；事件 6／7 state 單測固定雙向
   俘虜旗標、索引與順序。#76 以無肖像文字 modal 顯示。
3. 第二次呼叫沒有重建 `DI=SP`，#72 的 `\\2` 參數不能安全映射成城市；呈現層缺
   formatter 時整則 fail-closed，禁止把 `AH` 猜成城市、勢力或信賴度。最新 raw handler
   證據顯示第二次會讀 `SS:[DI]`，而此 transient stack word 不能從 World state 重建；
   因此事件 6 現在固定 `RawFormatterWord=-1`、`RawFormatterWordValid=false`，不再把
   Go 零值 `0` 偽裝成原版 payload。只有未來動態原版 trace 擷取到該 word 時才能標為有效。

### P0：事件 10（raw producer／consumer 已接）

1. `sub_131AE` dispatch 的低碼 `0x0A` 已確定進 `sub_13496`；該 handler 將
   `AL=事件字高 byte`、`AH=0xFF`、`CX=Param` 後呼叫 `sub_18810`。state consumer
   已保留 raw `Param` 為 TALK index，並在高 byte 是有效 General 時提供 `\\1`。
2. 直接檢查 `sub_12FBF`、`sub_12FB1`、`sub_1301C` 的 caller，尚未找到可證實的
   `0x0A` producer；這是負證據，不是全 binary 排除。
3. `World.QueueEvent10` 保留已證實的受控 raw producer；另新增
   `World.produceApproximateEvent10`，在月結後依玩家俘虜狀態近似排入 TALK `0x41／0x42`。
   近似規則預設開啟但可由 `SetApproximateEvent10(false)` 關閉，不冒充原版自然時序。
4. consumer 與兩種 producer 都沿 raw queue／formatter marker 接入；未知原版 TALK index
   不以泛用文案替代，近似路徑只使用已查到的俘虜文字並標為 substitute。
5. [x] `cmd/wlgame` 的 `idleClockGate` 已對應松崗繁中 `sub_11F7F` 的穩定游標條件：
   首次座標觀測、游標位移、按鈕或命令 frame 都停住世界；游標穩定且無輸入時才按
   據點／軍團／物件／時鐘順序跑 `World.TickMap`。`TestIdleClockGateRequiresStablePointerAndNoCommand`
   與 `TestIdleClockDispatchesQueuedEvent10OnHourlyCadence` 已納入 parity gate。

### P0：物件動畫時序（timer／render 已完成）

1. 已依 DOS/V IDA 線性位址實作 runtime object record：
   - `sub_123FF` `000123FF` 建立物件時寫 `[si+0C]=1`、`[si+0D]=0x10`、
     `[si+0F]=1`，並把 type 寫入 `[si+0E]`。
   - `sub_12459` `00012459` 每次主畫面更新對 active object 遞減 `[si+0C]`；
     到零時補回 `[si+0D]` 並設 dirty bit 0。
   - `sub_12533` `00012533` 只在 dirty bit 命中時把 `[si+0F]` 加一並 `& 7`，
     再依 `CS:985Ah` 查表取 MCH frame。
   - `sub_12438` `00012438` 依同一城市座標清除 active object。
2. 原本的 `disasterFrameTicks=8` fallback 已移除，改成 typed、可測試的 per-object
   `Phase`／`Timer`／`Delay` 狀態；`World.DisasterMarker` 仍不序列化。
3. `TestDisasterObjectAnimationTiming` 已驗證「建立初值 → 首次 timer cadence →
   16 次 update → dirty → render 舊 phase → phase +1 → 清除」；`wlgame` 只在可見
   map-loop Update 推進，模態視窗不穿透。
4. MCH type 1／2 的圖像接線與固定時序已完成；type 3 查表保留但事件語意未知，
   暴風雨仍依 `sub_134A6`／`sub_1237E` 的已證實範圍處理。後半區 slots 16–31 的物件移動
   `sub_1248A` 分支已依 raw fixed-point word、方向 byte、邊界與 wrap 規則接入，
   由 `TestSub124FFMatchesRawSignedByteContract`、兩個 `MovingDisasterSub1248A*`
   測試固定。

### P1：DOS/V 自然畫面对拍與目標平台 GUI

1. 使用者影片 `af6xqcicXoI`（567 秒、478×360、30 fps）已取 20／80／160／240／
   320／400／480／550 秒代表幀；80 秒幀另去黑邊／縮放成 640×400。比較確認橫幅 32 px、
   命令列 32 px、左側 27×21 地圖與右側 208 px minimap／情報欄。
2. `cmd/wlgame/strategyhud.go` 已接入這個自然 HUD；remake 固定 `seed=17`、30 幀
   輸出 [`wlgame-dosv-natural-remake.png`](docs/images/wlgame-dosv-natural-remake.png)，
   並與 [`yt-wolong-natural-80s-640x400.png`](docs/images/yt-wolong-natural-80s-640x400.png)
   做結構／色彩／欄位位置對照。
3. 影片是有損縮放來源，且 80 秒為 196 年 4 月 5 日、remake smoke 為 196 年 4 月 1 日；
   因此「影片視覺 oracle 對拍」已通過，但嚴格同狀態逐像素 diff 不宣稱。
4. Linux／Xvfb GUI smoke 通過；Windows amd64 交叉建置產生 `PE32+ x86-64`，macOS
   amd64／arm64 產生 `Mach-O`。目前 Docker 沒有 Windows／macOS 原生桌面執行環境，
   這兩項仍是 GUI 編譯／格式 gate，不寫成原生 runtime parity。

## 2026-08-10 最新增量：DOS/V 硬體游標／button glyph

- `KI.EXE` `sub_201E4`（IDA `seg002:01E4`）設定 `SI=031Bh`，兩次呼叫
  `sub_2020C`；`seg002:031B` 的 64 bytes 對應 file `0x1051B`，第一層
  `AX=0F00h` 白色、第二層 `AX=0A00h` 紅色，解成兩層 16×16 mask。
- `internal/assets/gfx.DecodeDOSVCursor` 已固定原始每列 `AH`／`AL` 反轉與 MSB-first
  展開；單測固定白 39、紅 56、透明 161 與完整 16 列形狀。
- `ICONGRF` 第 3 段相對 `0x14A0` 的 96×64 資源下半部是 3×6 個 16×16 靜態
  button glyph；單測逐格確認非背景像素。`amountFrame` 存在時直接畫資源，不再覆蓋
  vector 矩形／CJK button label；`Library.DOSVCursor` 提供原版 cursor overlay，
  缺資源才 fallback。
- Docker 全量 `go test -p=1 -vet=off ./... -count=1`、`go vet ./...`、文件
  `index.py generate/check` 與 30 幀 Xvfb 短 smoke 已通過；完整長程遊戲仍依要求略過。

## 2026-08-11 本輪封口

本輪以短 fixture、可重跑單測與 IDA 靜態證據封口，不以完整長程遊戲測試作為必要條件：

- [x] `sub_1248A`：只對 raw slot 16–31 執行移動；兩次 `sub_124FF` 的有號 byte
  漂移、方向 byte、邊界與 `-0x10..0x190`／`-0x10..0x110` wrap 已實作並測試。
- [x] 事件 2／3／4／5：外交與資助的成功、拒絕、金額邊界及 TALK 展開／硬換行／5 行
  分頁由 `TestEvent2To5TalkBranchParityGate` 表格化驗收。
- [x] 事件 9：玩家勢力／非玩家勢力／在野與空白 #409 的短 fixture 由
  `TestEvent9ShortFixtureGate` 封口；不宣稱原版長程時序。
- [x] M7：`corrections.json` 的 60 筆校訂逐筆展開，逐筆檢查寬度與最多 5 行分頁，並由
  `TestM7CorrectedTalkLayoutGate` 加 `tools/talkdat_selftest.py` 驗收。
- [x] 投射物：一般水平／垂直與特殊 `0x214`／`0x215` frame、發射姿態位元及運動規則由
  `TestProjectileParityGate` 與 tactical projectile tests 驗收。
- [x] 事件 10：已完成 IDA Pro 9.4 `.i64` 的 queue dispatcher、consumer、writer、caller
  與 data-ref 深度追查。原版自然 producer 仍是未知；這是有界的負證據結論，
  `World.QueueEvent10` 保留為受控 raw producer，已不再把未知來源列為無限追查阻塞。

### 2026-08-11 事件 10 與無輸入自動 clock 勘誤

- [x] 已證實 DOS/V 無輸入路徑：`sub_11BE0` → `sub_11F7F`（座標不變設
  `byte_198A3` bit 7）→ `sub_11CD0` → `sub_13EFD`／`sub_125A3`／`sub_12459`／
  `sub_11D8E`；因此日期、已下達目的地的軍團、據點 runtime 與 MCH 物件會自動前進。
- [x] 已證實事件 10 只是這條路徑的下游 queue consumer：每小時進入 `sub_13E11`，
  再由 `sub_131AE` 依 `byte_131AD` 節流；初始化為 7，取到一筆後重設 0x0A，
  不是每小時同步取，也不是 clock／行軍 driver。
- [x] 新增 `TestIdleClockDispatchesQueuedEvent10OnHourlyCadence`，固定 7 個每時邊界
  驗證預先注入的 `Code=0x030A`／`Param=0x42` 只在節拍邊界產生 TALK；完整證據見
  [`docs/re/16-idle-clock-event10.md`](docs/re/16-idle-clock-event10.md)。
- [x] 原版自然 `0x0A` producer 仍未知，維持限時封口；不要把 `World.QueueEvent10`
  接成 clock。正常 UI map-loop 已用 `World.TickMap` 對齊原版據點／軍團／物件／時鐘；
  同一畫面的額外規則 tick 使用不含物件的 `World.Tick`，不污染 MCH 動畫 cadence。

### 2026-08-12 追加 IDA 再審：未知項的已證實邊界

- [x] 事件 10：逐一展開所有目前可見 direct queue producer；已證實集合僅寫出
  `0x01`–`0x09`、`0x0B`–`0x0D` 的已列舉路徑，沒有 `0x0A`。這強化有限負證據，
  仍不排除間接／外部／受保護流程，原版 natural producer 保持 **unknown**。
- [x] 事件 6／7：`\\2` handler `000108DB` 已證實第二次 formatter 讀 `SS:[DI]`；
  舊有「raw word 0」已撤回，#72 以無 payload fail-closed，#76 保持無肖像次要文字。
- [x] MCH：`sub_123FF` 唯一 direct caller 是 event 12 handler；`sub_12286` 只直接排
  `0x010C`／`0x020C`。type 3 仍未知；`sub_1248A` 是 slots 16–31，不是最後一筆。
- [x] 戰術：投射物已接原版 normal／special branch、`0x214／0x215` source frame 與
  raw `[si+13]` 的 8／6 發射冷卻；完整同狀態逐像素生命週期仍非已證實。
- [x] 音源：新增 `docs/re/17-dosv-audio-tsr.md`，記錄 `YNSOUND.COM` 的 INT 61h command
  table、遊戲效果碼與 address/data 硬體寫入；精確音色格式與硬體型號仍未知。

### 2026-08-11 事件 10 近似自然 producer

- [x] 月結後新增 `produceApproximateEvent10`：只選玩家勢力目前收容的活武將，每月最多
  一筆，保留 `General.Timer` 倒數閘與固定 RNG 邊界；逃走／歸降分別排入 raw TALK
  `0x41`／`0x42`。
- [x] raw payload 維持已證實的 `(general<<8)|0x0A`／`Param`，下一個每時 queue 節拍
  才由 `sub_13496` 對應的 remake consumer 產生 `TalkNotice`；不把 producer 接成 clock。
- [x] `SetApproximateEvent10(false)` 可關閉替代規則，`TestApproximateEvent10*` 固定
  payload、逃走／歸降狀態、queue 滿／倒數邊界與 idle clock consumer 接縫。
- [x] 原版自然 producer 仍維持 **unknown**；本項完成的是可遊玩的 substitute，不提升
  原版 parity 證據等級。完整說明見 [`docs/re/15-event10-producer.md`](docs/re/15-event10-producer.md)。

### 2026-08-11 M7 人工文字與畫面抽樣封口

- [x] `tools/m7_review.py --check` 逐筆確認 60 筆校訂、TALK marker、校訂表產出一致性、
  原始硬行與保守字寬；`TestM7CorrectedTalkLayoutGate` 再以 runtime `textdraw` 實測像素寬度。
- [x] 逐筆閱讀 `translations/corrections.json` 的語意備註，確認人名／勢力／據點／金額
  marker 沒有被文字修正誤換；#321／#322／#718 群組／#751 等槽位差異保留原始證據。
- [x] 以 `#321`、`#258`、`#663`、`#718`、`#889`、`#967` 代表幀回看硬換行、五列分頁、
  寬度、標點、尾端空行與戰場命令；完整表與畫面連結見
  [`docs/playtest/14-m7-review.md`](docs/playtest/14-m7-review.md)。
- [x] 本項完成的是 60 筆已定案校訂的人工審查，不把它擴大宣稱成 1,022 則全部重譯或
  同狀態逐像素 parity。

### 2026-08-11 事件 2–5 完整 TALK 抽樣封口

- [x] `TestEvent2To5FullTalkPageSampling` 以真實校訂 TALK、原始索引、marker 展開、硬換行、
  runtime 字寬與五列上限逐頁檢查事件 2／3／4／5 的 36 個 raw TALK 頁面、18 組雙頁回應。
- [x] 已覆蓋外交的自由／有資金／拒絕／超額，以及撥款的全額／等額／低額／零額／超額／
  拒絕分支；各分支的 raw index 與下一頁 index 均由 fixture 固定，不用泛用文案代替。
- [x] 四張代表幀與完整索引、限制條件見
  [`docs/playtest/15-event2-5-talk-sampling.md`](docs/playtest/15-event2-5-talk-sampling.md)。
- [x] 本項封口的是完整分支／逐頁抽樣與版面 contract；不把受控 fixture 宣稱成完整自然劇本
  長程或原版同狀態逐像素 parity。

### 2026-08-11 事件 9 長程通知流程封口

- [x] `TestEvent9LongNaturalRoute` 以正常 `World.Tick` 跑 27 小時、每小時 9 個 subtick，
  驗證 queue delay=7 的三筆事件在第 7／17／27 小時依序取出，並確認據點／軍團／物件／時鐘
  的 idle clock 順序沒有被事件 9 取用污染。
- [x] `TestEvent9LongNotificationRoute` 逐段驗證玩家勢力釋放武將才產生 #37 modal；非玩家勢力、
  在野武將與空白 #409 均不產生錯誤通知，玩家後續釋放仍可再次通知。
- [x] 完整 fixture、raw index、通知條件與代表幀見
  [`docs/playtest/16-event9-long-route.md`](docs/playtest/16-event9-long-route.md)。
- [x] 本項封口的是可重跑的長程 queue／通知流程；完整原版自然劇本仍依使用者要求不跑，
  當時也沒有做同狀態逐像素對拍。

### 2026-08-11 推廣影片產出

- [x] 已用 remake 實機流程代表幀串接自然策略、事件 2–5 TALK、戰術、投射物、戰果、事件 9、
  M7 與存檔畫面，產出 60 秒、1280×720、H.264/AAC 影片。
- [x] 配樂由 `tools/promo_score.py` 以原創合成音色產生；未使用原版 `SOUND.DAT` 或原版 BGM，
  權利與重現命令見 [`docs/promo/README.md`](docs/promo/README.md)。
- [x] 影片輸出：`dist/promo/wolong-remake-trailer.mp4`；正式三平台可執行包與目標平台原生
  GUI smoke 仍是獨立 release gate。

### 2026-08-12 DOS/V／remake 實機動態推廣片

- [x] 依使用者指定錄製／剪輯 60 秒、1280×720、H.264/AAC 的
  [`wolong-remake-dosv-live-comparison.mp4`](dist-all/promo/wolong-remake-dosv-live-comparison.mp4)。
  原版側是使用者指定松崗 DOS/V 遊玩錄影與受控 DOSBox-X 新遊戲畫面；remake 側是正常
  鍵盤路徑的策略、編成、目的地與行軍實機擷取。
- [x] 影片明示同類畫面比較「非同狀態逐像素判定」；戰術 remake 段標為獨立 fixture，
  不被挪作自然路徑驗收證據。原版音訊已排除，唯一音軌為本專案原創合成配樂。
- [x] 完整來源、秒數、雜湊、抽樣驗收與離線重播規則見
  [`docs/promo/dosv-live-comparison.md`](docs/promo/dosv-live-comparison.md)；暫存原版錄影在驗收後
  刪除，不隨平台遊戲包或工作區長期保存。

### 2026-08-11 DOSBox／remake 可玩性專家驗證

- [x] DOS/V 原版以固定 `cycles=20000` 在 DOSBox 啟動至密碼保護頁；不繞過密碼，
  原版自然玩法保持阻擋界線。PC-98 DOSBox-X `NEW GAME` 與既有自然 oracle 截圖可用。
- [x] 目前 remake 建置以 `-seed 17`、無 `-open-*` 旗標走過編成、關閉通知、行軍、
  玩家不下命令時的日期流逝與 196/6/28 遭遇選單；逐階段證據見
  [`docs/playtest/17-expert-dosbox-remake.md`](docs/playtest/17-expert-dosbox-remake.md)。
- [x] 正常系統視窗儲存／讀取第 1 槽通過，overlay 寫入 88,832 bytes；原始 DOS/V 檔案
  唯讀掛載，沒有被覆寫。
- [x] 目前建置的戰術 GUI／`2` 號攻擊完成 debug smoke；正常無旗標戰術接續沿用
  [`docs/playtest/09-wlgame-normal-tactical-path.md`](docs/playtest/09-wlgame-normal-tactical-path.md)，
  不把 debug fixture 宣稱成原版 parity。
- PC-98 DOSBox-X 的 bus-mouse／焦點輸入重播仍不穩定，但依使用者決策，PC-98 不再是
  本輪 DOS/V remake 的畫面或 release gate；保留作未來研究，不列為收尾阻塞。

### 2026-08-11 DOS/V 自然策略骨架對齊

- [x] 以 YouTube 80 秒 DOS/V 640×400 參考幀重對常駐骨架：32 px banner、32 px 命令列、
  左側 432×336 地圖與右側 208 px 欄位維持同一座標契約。
- [x] 修正右欄 minimap／情報框的共用分隔邊：原版 16 px 紅／藍勢力色標列覆蓋共用
  8 px 邊，不再產生 remake 的雙倍分隔；情報區改為君主／首都／軍師、信賴度、黑底
  資金／預備兵結構。最新畫面為
  [`docs/images/wlgame-dosv-natural-remake-skeleton.png`](docs/images/wlgame-dosv-natural-remake-skeleton.png)。
- [x] 依新骨架重錄 60 秒推廣片與 24 秒 YouTube 對照片；raw `AE` 由 `255003` 降為
  `249178`，`RMSE` 由 `0.338208` 降為 `0.329145`。差異解讀仍不升格為同狀態逐像素 parity。

### 本輪 gate 入口

在 `wolong-go` Docker 容器內執行 [`tools/parity_gate.sh`](tools/parity_gate.sh)。
事件 10 的非破壞性 IDA 匯出器是 [`tools/ida_event10_producer.idc`](tools/ida_event10_producer.idc)，
完整證據與推論等級見 `RESEARCH-LOG.md` 的 2026-08-11 章節。

### 2026-08-11 三平台候選封裝與 Android 規劃（歷史起點）

- [x] Linux amd64、Windows amd64、macOS Intel／Apple Silicon 的候選產物已在 Docker 建置；Linux 原生 `wlgame`／`wlview`、Windows PE32+、macOS `Mach-O` 均已核對檔頭。
- [x] `dist/release-20260811/packages/` 已產出 Linux amd64、Windows amd64 與 macOS universal 三個主 `.tar.gz`，另附 Linux arm64 邏輯工具包；解包清單、每包 `SHA256SUMS.txt`、外層雜湊與 deny-list 均通過。
- [x] 以封裝 Linux `wlgame` 執行 640×400、固定 seed、30 幀 Xvfb smoke；輸出 hash 為 `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`。
- [x] Android 規劃已建立：[`docs/mobile/android-plan.md`](docs/mobile/android-plan.md)。第一版鎖定橫向、保留 640×400 邏輯畫布、用安全區／觸控抽屜／手勢轉接改善手機操作。
- 當時尚未開始 Android 觸控層、固定 SDK／NDK／`gomobile` Docker 工具鏈與 APK smoke；
  這是歷史起點，現況以下方「Android 原型與 AppImage 封口」為準，仍不宣稱 Android 完整支援。

### 2026-08-11 Android 原型、AppImage 與「經典再現」推廣片封口

- [x] Linux amd64 AppImage 已產出：`wolong-remake-linux-amd64-20260811.AppImage`（已被 8/12 的 [`dist-all/packages`](dist-all/packages) 取代）。
  AppDir 根目錄含 `.desktop`／`AppRun`，deny-list 通過，Docker／Xvfb `APPIMAGE_EXTRACT_AND_RUN=1`
  啟動與 640×400 固定 seed 截圖 smoke 通過；不含原版資料與字型。
- [x] 新增「經典再現」原版／remake 比較片：
  [`wolong-remake-classic-revival.mp4`](dist-all/promo/wolong-remake-classic-revival.mp4)，60 秒、
  1280×720、H.264/AAC。原版側使用使用者 YouTube 的代表幀，remake 側使用固定 `seed=17`
  實機代表幀；影片與 [`docs/promo/classic-revival.md`](docs/promo/classic-revival.md) 明示
  `core=normal`、`cputype=486`、`cycles=20000` 的 DOSBox 重播原則，以及不宣稱同狀態逐像素 parity。
- [x] Android 觸控 shell debug APK 已建置並在 API 35 `google_apis;x86_64` 模擬器安裝／啟動；
  1080×1920 實體畫布旋轉成 1920×1080 橫向畫面，`CONTINUE` 顯示 TALK 頁、`MENU` 開啟命令抽屜。
  三張證據圖見 [`docs/mobile/android-plan.md`](docs/mobile/android-plan.md)；完整自然時鐘、事件、
  存檔／讀檔、實機與 release signing 仍未完成。
- [x] Android 驗證容器已停止並確認沒有留下 `wolong` 專案容器；本輪沒有建立 `HANDOFF.md`，
  後續工作仍只追加在本檔。

### 2026-08-11 事件 10 remake 封口

- [x] 使用者所指的「不下命令、不移動滑鼠時，部隊依已下達指令自動跑、日期時間流逝」
  已在正常 `wlgame` 主迴圈完成，不再僅是 state fixture。游標每次移動都會暫停該 frame，
  穩定下一 frame 才恢復；modal 仍依既有規則暫停。
- [x] 月結俘虜的 `0x41`／`0x42` substitute、raw `0x0A` queue 與每時 TALK consumer
  仍完整串接。原版自然 `0x0A` writer 仍是 **unknown**，但不再阻擋這項 remake 功能；
  不把 substitute 寫成松崗原版已證實的 producer。
- [x] 松崗繁中唯一 oracle 的 30 frame Xvfb smoke hash：
  `45a68852335420dd7b22b4e240192dcd7a38fbbc62f72c8c59ec95acdc137b24`。

## 仍未宣稱完成的邊界

- YouTube 原版遊玩影片與 remake 推廣片已完成並排比較，並以 640×400 自然畫面保存
  raw pixel diff、差異圖與 24 秒研究對照片；這已封閉可見像素差異／畫面骨架的驗收，
  詳見 [`docs/promo/yt-remake-pixel-review.md`](docs/promo/yt-remake-pixel-review.md)。
  當時還沒做「同日期／同輸入／同狀態」的逐像素 parity。
- Windows／macOS 目前有交叉建置／檔頭 gate，沒有目標作業系統原生 GUI runtime；三平台候選
  封裝已產出，但原生 GUI smoke 尚未完成。本輪推廣影片已產出，不取代目標平台驗收。
- Android 目前只有已在模擬器驗證的觸控 shell debug APK；完整核心、自然 clock、事件流程、存檔／
  讀檔、實機／平板與 release signing 仍未完成，不能標成 Android release。
- 完整長程遊戲測試依使用者指示略過；這不影響上述短 fixture 的封口。

## 證據入口

| 主題 | 入口 |
|---|---|
| 總體狀態與歷史勘誤 | [`CONTEXT.md`](CONTEXT.md)、[`RESEARCH-LOG.md`](RESEARCH-LOG.md) |
| 事件 6／7 對話索引 | [`docs/re/12-diplomacy-dialogue.md`](docs/re/12-diplomacy-dialogue.md) |
| DOS/V 數值視窗 | [`docs/re/13-pc98-numeric-window.md`](docs/re/13-pc98-numeric-window.md)（檔名為歷史名稱，本輪內容須以 DOS/V 為準） |
| MCH 物件格式 | [`docs/re/14-mmap-mch-objects.md`](docs/re/14-mmap-mch-objects.md) |
| 事件原版 fixture | [`docs/playtest/11-event6-original-fixture.md`](docs/playtest/11-event6-original-fixture.md) |
| 同狀態截圖規則 | [`VERIFICATION-MATRIX.md`](VERIFICATION-MATRIX.md)、[`docs/playtest/12-event3-same-state-parity.md`](docs/playtest/12-event3-same-state-parity.md) |
| M7／事件 2–5／事件 9 抽樣 | [`docs/playtest/14-m7-review.md`](docs/playtest/14-m7-review.md)、[`docs/playtest/15-event2-5-talk-sampling.md`](docs/playtest/15-event2-5-talk-sampling.md)、[`docs/playtest/16-event9-long-route.md`](docs/playtest/16-event9-long-route.md) |
| 推廣影片與配樂 | [`docs/promo/README.md`](docs/promo/README.md)、[`docs/promo/dosv-live-comparison.md`](docs/promo/dosv-live-comparison.md)、`dist/promo/wolong-remake-trailer.mp4`、`dist/promo/wolong-remake-classic-revival.mp4`、`dist/promo/wolong-remake-dosv-live-comparison.mp4` |
| IDA DOS/V 證據 | `workplace/ida/dosv/KI.EXE.i64`、`KI.EXE.asm`、`func-sub_*.txt` |

DOS/V 證據輸入：`KI.EXE` SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`，IDA `.i64`
SHA-256 `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`，工具為
IDA Pro 9.4，位址基準為線性位址。

## 驗證與清理

所有 Go、Python、GUI、DOSBox、截圖與建置都在 Docker 內執行；主機只做 Docker 控制、
Git 狀態檢查與 `apply_patch` 編輯。每次工作結束至少執行：

```text
docker run --rm --network none ... go test -p=1 -vet=off ./... -count=1
docker run --rm --network none ... go vet ./...
docker run --rm --network none ... tools/talkdat_selftest.py
docker run --rm --network none ... tools/index.py generate
docker run --rm --network none ... tools/index.py check
git diff --check
git status --short
docker ps -a
```

一次性 container 必須 `--rm`、有界 CPU／記憶體／PID、以目前 UID/GID 寫入；原始素材
唯讀掛載。明確刪除本輪 `/tmp` fixture、截圖與 binary，不做廣域 `/tmp` 清理；不得停止
其他專案既有容器。交接時記錄沒有留下本輪持續執行 container。

## 不要重做或誤宣稱

- 不要把 PC-98 screenshot、游標或外框當 DOS/V release oracle。
- 不要把通過 remake 單測寫成原版 parity；必須附 DOS/V 位址、raw bytes、fixture 或
  同狀態畫面證據。
- 不要猜事件 6／7 次要 formatter、事件 10 producer 或未知 marker；#72 缺參數時
  必須維持 fail-closed。
- 不要把 `HANDOFF.md` 重新建立；後續接手、進度與剩餘工作只更新本檔。
- 不要把三平台候選包的交叉建置／檔頭驗證寫成 Windows／macOS 原生 GUI parity；正式發布仍
  需要目標平台短 smoke 與輸入、音訊、字型載入驗收。

## 2026-08-12 DOS/V 密碼頁驗證勘誤

- [x] 松崗 DOS/V 密碼頁的有效輸入橋接已由 `wolong-dosboxx:latest` 的
  `mouse_emulation=integration`／INT 33 640×480 建立；不是舊 `dosbox-run` 腳本的
  DOSBox 命令列假陽性。
- [x] 新的唯讀原版副本中，空白確認、`0000`、`1234` 都進入原版開場；空白與 `1234`
  的 10 秒後畫面雜湊相同。這表示「任意數字會過」可重現，且數字非必要。
- [x] 密碼頁不納入 remake，亦不再作為 DOS/V 原版自然流程／畫面採樣的阻擋理由。
  `PASS.*`／`YNFONT.EXE` 的實際比較語意、真實硬體行為、完整長程流程與同狀態逐像素
  parity 仍各自保持未驗證；詳見 `docs/playtest/18-dosv-password-verification.md`。

## 2026-08-12 一般玩家啟動／新遊戲殼層切片

- [x] 新增 `cmd/wlgame/launcher.go` 純狀態啟動殼層：`NEW GAME` 確認、四劇本選擇、由實際劇本資料篩出的合法玩家勢力／君主、確認與返回。
- [x] `-save-file` 存在時提供四槽 `LOAD DATA`；以既有四槽 overlay 解析判定可用槽，空槽 fail-closed，不讀取或修改 `SINARIO.DAT`。
- [x] 一般流程在確認後才呼叫 `startWorld`，再掛道路、AI、戰術資料、季節外框、硬體游標與相機；既有驗收入口使用明確 direct-start 白名單，不以任意旗標跳過 launcher。
- [x] 鍵盤與滑鼠共用 launcher selection state；首次靜止滑鼠位置不覆蓋鍵盤選取，實際滑鼠移動／點擊才改變選取。
- [x] `cmd/wlgame/launcher_test.go` 覆蓋新局成功、取消／返回、非法玩家、空槽拒讀與成功讀槽狀態轉移。
- [x] 正式視窗標題改為「臥龍傳－三國制霸之計」；本切片不處理戰術或自然 HUD polish，不宣稱逐像素 parity。
- [x] Docker/Xvfb 實際擷取修正後三張 640×400 launcher 畫格：`/tmp/wolong-launcher-title.png`（`60def98e0cf54726ad62794b92906017f200863241442db68aa3f536eb3b5150`）、`/tmp/wolong-launcher-scenario.png`（`9fc49e523b75e45b5a939177272a4b199fe1a36d4dd7fb67fc506b70d70eeb04`）、`/tmp/wolong-launcher-player.png`（`0b3a7f4709b538a6c9729a3bb745e3b49c523d999c8e7bfc596e1178cf5f38eb`）；確認 title／scenario 未選文字可讀，player 8 列、標題、反白與提示均未越過 panel 安全區。
- [ ] 四劇本與可用存檔槽的逐一代表畫面仍未錄製；不阻擋本小切片交付，完整長程測試依使用者要求不跑。

## 2026-08-12 Mentor polish 收斂

- [ ] DOS/V 戰術畫面已有可玩骨架與鍵鼠 dispatcher，但推廣片複驗證實戰場 viewport、上下 TALK 區、右側縮圖／狀態／命令 glyph 與底列配置尚未對齊原版；撤銷「完整骨架」宣稱。
- [x] 戰術縮圖由 `BATTLE.MAP` 與 `BATTLE.MDL` attribute 動態產生 128×128 圖，不再使用高度圖替代。
- [x] 依 `sub_1075B` 已證實公式接入開戰兩筆 `TALK.DAT`；每場只初始化一次、未知 marker fail-closed、戰場時間不因對話停止。實機證據：`docs/images/wlgame-tactical-opening-talk.png`。
- [x] 已證實兩種敗北 latch：信賴度歸零優先顯示 TALK #414；最後據點失守使用克制 fallback。研究備註不進 GUI。
- [x] 存讀檔四槽與一覽表新增滑鼠／觸控列、分頁、確認、取消；遵守原版兩段式選取且 modal 不穿透背景。
- [ ] 指揮／事件／一覽畫面的 DOS/V 幾何 parity 未完成；現有推廣片「指令與事件」左右不是同類畫面，不能作為還原證據，必須重做同類畫面對拍。
- [x] 推廣片異類「事件 vs 目的地」鏡頭已撤換為「原版系統設定 vs remake 系統設定」；戰術段已重錄為 `240:80` 主要幾何。影片視訊與音訊均為 60 秒。
- [x] 系統設定中央五列、事件左下 TALK、一覽第一層主要外框已依松崗錄影正規化座標修正。
- [ ] 戰術右欄原版命令 glyph／內框、一覽左側捲軸與選取後的前層武將詳細窗仍未完成；不得以本輪主要幾何修正宣稱完整 parity。
- [x] Docker＋Xvfb `go test -p=1 -vet=off ./... -count=1` 與 `go vet ./...` 通過；文件索引 65 份通過。
- [x] 三平台包與「經典再現」推廣片已以本輪畫面重建；`dist-all/` deny-list 與全部 SHA-256 驗證通過。實機對照主片的新版戰術段含 TALK、縮圖、六指令與動態戰場。
## 2026-08-12 攻城／兩軍遭遇共用戰術骨架

- [x] `combat.Siege` 與 `combat.Field` 的實際繪製、TALK slot、右欄、底部六命令及鍵鼠命中區均無模式分支；唯一幾何來源為 `dosvBattleLayoutFor`。
- [x] 新增 `TestFieldAndSiegeShareExactDOSVBattleChrome`，防止後續 polish 分裂成兩套畫面或命中區。
- [x] 修正 `demoBattle` 預跑 900 tick 導致野戰在 GUI 出現前結束的假差異；改為 120 tick。
- [x] Docker／Xvfb 六格實機對照完成，見 [`docs/playtest/22-field-siege-shared-layout.md`](docs/playtest/22-field-siege-shared-layout.md)；全專案 `go test -p=1 ./...` 與 `go vet ./...` 通過。
- [x] 純城兵據點攻擊保留原版自動判定，不擅自導入手動戰術畫面。
- [x] 原版指令 glyph／右欄複合面板已依 `sub_1C7F4`、`sub_1C863`、`sub_1F888`、`sub_1C6BF` 接入；舊的 2×3 文字格與綠色選取框已移除。
- [x] 三平台完整包與 AppImage 已重建至 `dist-all/`；Linux tar／AppImage Xvfb smoke 皆維持基準 SHA-256 `45a68852…b5150`，deny-list 掃描 19 個交付檔並比對 120 個原版檔後通過。
- [x] DOS/V／remake 推廣片的 40–45 秒戰術段已替換為原版 glyph／右欄面板實機畫面；繁中 overlay 以 Noto CJK 代表幀驗收，成品 60 秒、1800 幀，SHA-256 `feddd663…3961e8`。
- [ ] 兩模式仍共用同一組未完成項：戰術動畫時序與同狀態逐像素 parity。
# 2026-08-12 推廣片 AdLib／戰術骨架勘誤

- [x] 移除推廣片對 `tools/promo_score.py` 的依賴，改用使用者松崗 DOS 原版錄影中
  的實際遊戲 AdLib 音軌；來源與權利邊界見 `docs/promo/dosv-adlib-and-tactical-review.md`。
- [x] 原版錄影先正規化／裁成 640×400，再與 remake 共用最近鄰縮放鏈。
- [x] 戰術離屏 buffer 從 496×384 收斂到 480×368，消除右、下各 16 px 的遮切。
- [ ] 完成同戰況戰術 capture pair：同攻城節點、攻守方、編成、命令、鏡頭與 frame。
- [ ] 依 pair 修完雙 TALK、右欄完整狀態、底列／側欄 glyph 與選取時序。
