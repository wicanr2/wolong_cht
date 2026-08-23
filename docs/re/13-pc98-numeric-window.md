# 13 — DOS/V 數字輸入視窗量測與 CJK 版面決策

**狀態：以 DOS/V 為唯一畫面基準；原始座標、DOS/V 96×64 內框 blit、3×6 每格操作、
實際按鍵 glyph、DOS/V 內建硬體游標、座標選取、TALK 五行分頁、event3 場景／肖像
composite 與 pending 結束後消像已接入，並由 Docker／Xvfb 短 smoke 驗證；完整劇本
長程測試依使用者要求略過。**

- 日期：2026-08-10

本文件只處理事件 2／3／4／5 的數字輸入畫面。遊戲規則與畫面均以松崗 DOS/V
`KI.EXE`／`ICONGRF.DAT`／`TALK.DAT` 為主要來源；PC-98／Golden Box 只保留為
歷史格式與 CJK 版面交叉驗證，不取代 DOS/V 畫面 oracle。

## 1. 原始量測

證據輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.asm`，SHA-256
`FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；同批 binary
`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4；下列位址均為 IDA 線性位址，
不是檔案偏移。

| 項目 | 原始證據 | 結論 | 等級 |
|---|---|---|---|
| 畫布 | `sub_1EB6C` 設定 640×400 EGA 模式；`docs/formats/07` 亦由 80 byte 列跨距交叉確認 | 保留 640×400 邏輯畫布 | 已證實 |
| 呼叫位置 | `sub_13902`／`sub_139E8` 呼叫 `sub_17C6E` 前設 `DX=0058h`、`BX=00B8h` | **錨點由呼叫端給**；事件那四支是 `(88,184)`，財政的四支是 `(296,184)`（[`../spec/78`](../spec/78-amount-input-editor.md) §1.2）| 已證實 |
| 保存範圍 | `sub_17C6E` 先把 `DX`／`BX` 各減 `8` 呼叫 `sub_19796`，結束前由 `sub_197C3` 還原 | 保存／還原區域從 `(80,176)` 起 | 已證實 |
| 輸入格起點 | `sub_17D5F` 先 `add bx,10h`，再以目前 `DX`／`BX` 呼叫 `sub_1E3D7` | 內容首格從 `(88,200)` 起 | 已證實 |
| 格子排列 | `sub_17D5F`：`CH=3`、`CL=6`；每欄 `DX += 10h`，每列 `BX += 10h` | 3 列 × 6 欄、16×16 格 | 已證實 |
| 外框呼叫 | `sub_17D0D` 使用 `AL=51h`、`CX=507h`，並以 `DS:SI=word_10D50:0600h`、`AX=4006h` 呼叫 `sub_1FA37` | 使用數值視窗專用圖形資源，外框 blit 的輸入參數已固定 | 已證實 |
| 外框 blit 尺寸 | `sub_1FA37` 將 `AX=4006h` 傳給 `sub_1FAA2`；後者以 `DH=40h` 列、`DL=06h` 且每列複製兩次 byte | 從資源拷貝 96×64 像素區域；這是 blit 尺寸，不等同已解出裝飾內容 | 已證實 |
| 每格內容 | `sub_17D5F` 從 `CS:7D93` 讀 18 bytes；`sub_1E453` 讀回點擊座標的畫面 buffer byte；`sub_17C6E` 以 `AL-52h` 分派 | 3×6 格是可實際選取的數字／功能按鍵 | 已證實 |
| 每格 raw table | `CS:7D93h` = `59 5A 5B 5D 5E 5E / 56 57 58 52 5F 5F / 53 54 55 5C 60 60` | 3 列 6 欄；`52h..5Bh` 為數字 0..9、`5Ch` 百、`5Dh` 退位、`5Eh` 清零、`5Fh` 還原初值、`60h` 完成 | 已證實 |
| 游標／選取 | `sub_121E7` 輪詢 input API，`sub_1E453` 以 `CX/DX` 對應畫面 buffer；remake 以同一 `(88,200)` 格座標接受滑鼠點擊，並以高亮框保留目前格 | 實際滑鼠格選取與跨平台鍵盤 fallback 已接入 | 已證實 |
| DOS/V 硬體游標 | `sub_201E4`（`000201E4`）兩次呼叫 `sub_2020C`；`SI=031Bh`（`seg002:031B`）各取 32 bytes。第一層 `AX=0F00h`、第二層 `AX=0A00h` | `KI.EXE` file `0x1051B` 的 64 bytes 解為兩層 16×16 mask；`DecodeDOSVCursor` 已接入 `Library` 與 `wlgame` | 已證實 |
| DOS/V 按鍵 glyph | `sub_17D0D`（`00017D0D`）的 `word_10D50:0600h`／`AX=4006h` → `ICONGRF` 第 3 段 `0x14A0`、96×64；資源下半部與 `(88,200)` 的 3×6 格重合 | 真實資源每個 16×16 格均含靜態 glyph；有資源時不再以 vector／CJK 標籤覆蓋，缺資源才 fallback | 已證實 |

### 1.1 外框尺寸的反組譯推導

證據輸入：唯讀 DOS/V `workplace/ida/dosv/KI.EXE.asm`，SHA-256
`FFFEBA2D9E6EE947A4CF7ABF8FEF6D4B8D0FB4E6E0EC66D9D34D7B5A0D43868`；同批 binary
`KI.EXE` SHA-256 `FFFEBA985231CDA4D636E93D10F598470B1F691D00275E4AA38E285893D43868`。
工具為 `ida-pro-9.4-ver2:uidfix-v1`／IDA Pro 9.4，位址基準為 IDA 線性位址。

- `sub_17D0D`（`00017D0D`）先以 `AL=51h` 呼叫 `sub_1895D`，再設定
  `DS=CS:word_10D50`、`SI=0600h`、`AX=4006h` 呼叫 `sub_1FA37`。
- `sub_1FA37` 把 `AX` 複製到 `CX` 後呼叫 `sub_1FAA2`。`sub_1FAA2` 將 `DH` 作為列數、
  每列把 `DL` 減至零；每次迴圈讀寫兩個 byte，故 `DL=06h` 對應 12 bytes＝96 個 8-pixel
  planar 欄位，`DH=40h` 對應 64 列。
- `word_10D50` 指向的圖像內容已由 `ICONGRF.DAT` 第 3 段直接解出；其下半部含原版
  3×6 靜態按鍵 glyph。硬體游標則由 `KI.EXE` 的兩層 mask 重建；因此本項已不再以
  替代高亮框冒充 glyph。尚未宣稱的是原版自然執行畫面與 remake 的整張截圖逐像素
  對拍——這個視窗還沒去拍，密碼頁不擋（`CLAUDE.md` §4.0）。

## 2. CJK 比較參考

以下是同類 PC-98 Golden Box 作品的佈局參考來源；它們不是臥龍傳規則或美術素材
來源，也不把其框線、肖像、圖示或文字複製進本專案：

- [克萊恩英豪（Champions of Krynn）PC-98 screenshots](https://www.mobygames.com/game/833/champions-of-krynn/screenshots/pc98/)
- [幽靈騎士（Death Knights of Krynn）PC-98 screenshots](https://www.mobygames.com/game/2219/death-knights-of-krynn/screenshots/pc98/)
- [Pool of Radiance PC-98 screenshots](https://www.mobygames.com/game/502/pool-of-radiance/screenshots/pc98/)
- [Curse of the Azure Bonds PC-98 screenshots](https://www.mobygames.com/game/503/curse-of-the-azure-bonds/screenshots/pc98/)

可安全採用的跨作品模式是資訊階層，不是裝飾：主要內容、持久狀態、情境指令分區；
數字／表格欄位對齊；16×16 整數格；主題換色不改換行、選取或狀態。這與本專案
現有 `640×400`、倚天 16×15 點陣字、16 像素列高相容。

## 3. 原版格位與操作映射

`sub_17C6E` 在 `00017CA4` 先把目前 `SI` 放入 `AX` 顯示，再等待
`sub_121E7`；點擊後 `sub_1E453` 回傳格內 byte。`AL` 先減 `52h`，其值同時
作為數字參數與操作表索引：

| raw byte | 畫面格 | 動作 | 狀態層對應 |
|---|---:|---|---|
| `52h..5Bh` | 3×6 中的數字格 | 目前值 × 10 + 0..9 | `AmountAppendDigit` |
| `5Ch` | 第 3 列第 4 格 | 目前值 × 100 | `AmountAppendHundred` |
| `5Dh` | 第 1 列第 4 格 | 除以 10 | `AmountDeleteDigit` |
| `5Eh` | 第 1 列第 5／6 格 | 清零 | `AmountClear` |
| `5Fh` | 第 2 列第 5／6 格 | **設成上限**（`si = [bp+0]`）| `AmountSetMax` |
| `60h` | 第 3 列第 5／6 格 | `STC` 結束輸入、保留目前值 | `AmountFinishInput` |

`sub_17DA5` 的 `xchg AX,SI` 後以 `SI = SI×10 + AL`，所以前一版把它簡化成
「加十」是勘誤；目前 `internal/state.editAmount` 與 UI 格位已按這個資料流固定。

⭐ **`ax` 是上限，不是初值。** `sub_17C6E` 把它存進 `[bp+0]`，
而三個地方都只把那一格當上界用：`sub_17DA5` 與 `sub_17DC3` 算完各 `cmp/jbe`
夾一次，`sub_17DEC` 直接 `si = [bp+0]`。**`si` 開場是 `xor si, si`**——
呼叫端沒有任何管道給初值。六個呼叫端傳的分別是 `7530h`（事件 2／3／4／5）、
`64h`（稅率）與三個 `2710h`（募兵），
與 [`48`](48-window-display-list.md) §6 從圖庫解出來的按鍵字樣「最大」一致。
完整的參數表與 remake 實作在 [`../spec/78`](../spec/78-amount-input-editor.md)。

## 4. 本專案決策

- `cmd/wlgame` 以滑鼠 `(x,y)` 命中同一個 3×6 格；數字鍵、退格、Insert、Delete、Home
  是跨平台 fallback，全部轉成 `AmountEdit`，規則層不接 PC-98 掃描碼。
- 事件 2／3／4／5 先畫原版 `IVENTGRF`／TALK composite；只有確認指定金額後才覆蓋
  `(80,176)` 數值器，符合 `sub_13902`／`sub_139E8` 的流程順序。
- 數值格命中與 X11 點擊仍是跨平台輸入接縫；畫面有 DOS/V `KI.EXE` 解出的 16×16
  白框／紅填硬體游標，靜態 3×6 按鍵則直接來自 `ICONGRF` 96×64 資源。原始資產
  缺失時才保留通用 fallback，不把 fallback 當作 DOS/V glyph。
- pending 選單只在當幀繪製 `IVENTGRF`／肖像；完成、拒絕或取消後清除 pending，下一個
  `Draw` 先重畫地圖，再顯示後續 TALK，所以場景與肖像不會殘留在地圖上。這是本輪的
  消像 parity；它不等同原版淡出或硬體畫面 blit 的逐像素 parity。
- 中文長字串、全形標點、半形數字與混排必須在 640×400 原生解析度測試；不得用
  插值放大後的截圖判定沒有溢位。

## 5. 原生解析度驗收

完成數字視窗呈現前，至少要有下列可重跑證據：

1. `(88,184)` anchor、`(88,200)` 首格與 3×6／16×16 量測能由畫面 metadata 回查。
2. 連續五列密集繁體中文不重疊，數字欄不被全形字推離邊界。
3. 編輯數值時時間仍停止；取消／確認後狀態層結果與 `TestRawAmountEditorSemantics`
   完全一致。
4. 外框的 96×64 blit 尺寸、3×6 raw 操作表、靜態 glyph 與 DOS/V 硬體游標 mask 已固定；
   `DecodeDOSVCursor` 單測固定白色 39、紅色 56、透明 161 像素及 16 列形狀。
5. EGA／現代主題換色不改 16 像素格、欄寬、換行、選取與遊戲狀態。

因此本文件解除事件 2／3／4／5 的「功能性數值選取／TALK 分頁／消像／DOS/V 游標與
按鍵 glyph 資產」gate；仍未封口的是自然 DOS/V 執行畫面對拍、其他事件 formatter／
物件動畫與目標平台 GUI，不把這些未完成項倒推成資產未知。

## 6. PC-98 event3 fixture 與輸入鏈補證

2026-08-10 以 PC-98 `KI.EXE`（SHA-256
`061917F9F3F5C03E29397A9C636D546052128A99B8C8CE31DED0E84CF2A481E8`）及一份
暫存 `SAVE.DAT` fixture 重播事件 3。fixture 的事件字為 `0x0303`，即來源勢力 3、
目標在 Param 低 byte 指向玩家 0；來源勢力外交官欄設為有效武將 17。原版可穩定進入
前置外交通知與三選一畫面，證明這條 event producer／handler／畫面 dispatch 鏈不是
只靠靜態 patch 拼出的假路徑。

| 觀測 | 證據 | 結論 | 等級 |
|---|---|---|---|
| 三選一 caller | IDA Pro 9.4，DOS/V `sub_13B7E`（`00013B7E–00013BA9`） | 先以 `sub_19796` 畫 `(DX=0050h,BX=00B0h,CX=600Ah)`，保存初始 `AL`，再以 `CX=[BP+0]`、`DX=0B05h` 進入輸入函式 | 已證實 |
| 輸入保存／還原 | `sub_193E9`（`000193E9–00019409`）→ `sub_101B4`／`sub_101DB` | 三選一使用通用輸入／游標狀態保存，不等於 `sub_17C6E` 的數字編輯 loop | 已證實 |
| 原版選擇畫面 | `docs/images/pc98-oracle-event3-choice.png`；640×400；PNG SHA-256 `56A1A16BC6D92F75A5DB3DC49F3C961609F1EB6B908106A39136D4E3E32FDB5C` | 已取得同一 fixture 的事件 3 三選一畫面；其內容含原版肖像／動畫疊層與三列選項 | 已證實 |
| remake composite 選項 | `docs/images/wlgame-event3-choice.png`；640×400；PNG SHA-256 `CA40B865B44A6EA13ED5B4F2C0B6AB913A0BC895EF48D7A19E1825501E535151` | `IVENTGRF` page 0、玩家君主 64×64 肖像、TALK prompt、三列選項與高亮游標同畫面 | 已證實 |
| remake 數值器 | `docs/images/wlgame-event3-amount.png`；640×400；PNG SHA-256 `27A5474EBA79C92C23B24A79938CA4E1D376B9FA52C0956AE3D3359C0404609D` | Enter 進指定金額後，以 X11 點擊 `(88,200)` 格，畫面保留 3×6 raw 操作與游標高亮 | 已證實 |
| 消像／後續 TALK | 同一 fixture 的 pending→resolve→Draw 路徑；`main.go` Draw 順序 | pending 結束後不再畫 `IVENTGRF`／肖像，地圖重畫後才顯示後續 TALK；五行頁面按 Enter／Space 逐頁 | 已證實 |
| 資產／逐像素 parity | DOS/V `KI.EXE` mask、`ICONGRF` 第 3 段與 IDA blit 證據 | 96×64 內框、3×6 靜態 glyph、16×16 白框／紅填游標已由原始 bytes 解出並接線；自然 DOS/V／remake 整張畫面对拍仍缺 | 資產已完成；自然畫面对拍 P1 |

原版數字編輯函式 `sub_17C6E` 的直接輸入分派已落成：`sub_17D5F` 讀 `CS:7D93` 的
18 bytes 畫 3×6 格；`sub_1E453` 回讀點擊格；`sub_17DA5`、`sub_17DC3`、
`sub_17DDD`、`sub_17DEC`、`sub_17DF1`、`sub_17DEA` 分別對應逐位追加、追加 100、
退位、還原初值、清零、完成。remake 的格位表、`internal/state.editAmount` 與
`AmountEdit` 已用單測固定；本輪只以 DOS/V 為畫面基準，PC-98 硬體游標不作視覺依賴。

## 7. TALK 中文排版與目前 release 界線

remake 先代入已證實 marker，再以 `internal/ui/textdraw.WrapLines` 按實際點陣寬度換行：
ASCII 為 8 px、非 ASCII 為 16 px，原始 TALK 行仍是 hard boundary；結構性的最後空行
不佔頁面，中間的空行仍保留，關閉標點不會單獨落到下一列。一般訊息 modal 使用
22 個全形 cell（352 px）寬度；原版式 TALK／肖像框使用 160 px 文字區，兩者都以
`sub_18810` 的五行／16 px 分頁。事件 2／3／4／5 的 `IVENTGRF` page、speaker 肖像、
prompt、選項與消像（pending 結束後由地圖重畫）已接在同一呈現層。

因此本文件目前可宣稱的是「DOS/V 原始行邊界 + 可重播的繁中五行分頁 + event composite／
數值格選取」；完整劇本長程測試仍依要求略過，事件 6／7 次要 formatter、事件 10 與
災害／投射物的證據分別記在 `docs/re/12`、`WORKLIST.md`，不在本文件冒充畫面 parity。

## 8. 2026-08-10 DOS/V 內框資源接線

本輪以 DOS/V `KI.EXE`（SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`）的
IDA Pro 9.4 線性位址為準。先前只記錄 blit 尺寸、沒有接上資源，這一輪補上：

- `sub_17D0D`（`00017D0D`）使用 `DS:SI=word_10D50:0600h`、`AX=4006h`；
  `sub_1FA37`／`sub_1FAA2` 的實際來源是 ICONGRF 第 3 段內的 96×64 平面圖。
- 依 `sub_100DF` 的段指標配置，`word_10D50:0600h` 換算為
  `ICONGRF.DAT` 第 3 段相對 byte offset `0x14A0`；程式以
  `gfx.DOSVAmountPanel`／`DecodeAt` 固定此定位，不把它誤切成 8×8 chrome tile。
- `sub_17C6E` caller 的保存區仍從 `(80,176)`、112×80 開始；96×64 內框目的地是
  `(88,184)`；`sub_17D5F` 的首個 16×16 操作格仍是 `(88,200)`。
- `cmd/wlgame` 有 DOS/V 資源時直接繪製這張內框，缺資源才降級通用框；數值格的 raw
  table、滑鼠命中與鍵盤 fallback 不變。`sub_17D5F` 的 raw table 與 `ICONGRF` 下半部
  靜態按鍵 glyph 已分開驗證；`KI.EXE` 游標 overlay 由 `Library.DOSVCursor` 提供，
  fallback 只在資產缺失時啟用。

這項接線把「外框位置／尺寸／3×6 靜態 glyph／硬體游標來源」提升到 DOS/V 資源 parity。
**數值輸入視窗本身還沒做整張截圖的逐像素對拍**——缺的是那張截圖，
不是取得截圖的手段：松崗 DOS/V 的密碼頁四格留白直接確認就會過
（`CLAUDE.md` §4.0、[`../playtest/18`](../playtest/18-dosv-password-verification.md)），
主畫面已經這樣拍完並逐區對拍（[`../playtest/37`](../playtest/37-main-screen-parity.md)）。

## 9. 2026-08-10 DOS/V 硬體游標與按鍵 glyph 解碼

### 9.1 `KI.EXE` 內建游標

證據輸入為 DOS/V `workplace/orig/dosv/KI.EXE`，SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；反組譯來源為
`workplace/ida/dosv/KI.EXE.asm`，IDA Pro 9.4／`ida-pro-9.4-ver2:uidfix-v1`，位址基準
為 IDA 線性位址。`sub_201E4`（`seg002:01E4`）把 `SI` 設為 `031Bh`，連續呼叫
`sub_2020C` 兩次；每次 16 列、每列 word，故來源是 `seg002:031B` 的 64 bytes。
MZ header 為 `0x200` bytes，對應原始檔 `0x1051B..0x1055A`。

- 第一個 32-byte mask 由 `AX=0F00h` 的 EGA set/reset 畫成 palette index `0x0F` 白色外框。
- 第二個 32-byte mask 由 `AX=0A00h` 的 EGA set/reset 畫成 palette index `0x0A` 紅色填色。
- `sub_2020C` 讀 `[si]` 後把 `AH` 寫 `[bx]`、`AL` 寫 `[bx+1]`；decoder 按每列反轉
  source byte，並以 MSB-first 展開，不把 16-bit word 當成一般線性 bitmap。
- 解碼結果為白色 39、紅色 56、透明 161；palette-index buffer SHA-256 為
  `385c2f1949d3d1e331399316305db7d7f2fd489a0626de1c6b2b8375aadfc6fe`。

### 9.2 `ICONGRF` 靜態按鍵 glyph

`sub_17D0D`（`00017D0D`）的 `DS:SI=word_10D50:0600h`、`AX=4006h` 經
`sub_1FA37`／`sub_1FAA2` 對應 `ICONGRF.DAT` 第 3 段相對 `0x14A0` 的 96×64 資源。
資源 `(88,184)` 起始的下半部正好覆蓋 `(88,200)` 起始的 3×6、16×16 格；每格均含
大量不同於格背景的像素。`TestDOSVAmountPanelContainsStaticButtonGlyphs` 逐格檢查六欄、
三列，每格至少 128 個非背景像素，證明這些是資源內的靜態 button glyph，不是
`CS:7D93h` raw hit-test byte 本身。

### 9.3 接線與驗收

`internal/assets/gfx/cursor.go` 保留 IDA segment offset、檔案 offset 與透明／顏色常數；
`internal/assets/library.Library` 讀取 `KI.EXE` 後快取游標；`cmd/wlgame` 在真實
`amountFrame` 存在時直接繪製原始 96×64 內框，並疊上解出的 16×16 cursor。舊 vector
矩形與 CJK 按鈕標籤只在原始資源缺失時使用。Docker 內已通過：

```text
go test -p=1 -vet=off ./internal/assets/gfx ./internal/assets/library ./cmd/wlgame
```

`cmd/wlgame` 測試在有界 Xvfb 下執行；完整長程遊戲測試仍依使用者要求略過。這一節
封閉的是 DOS/V 資產解碼與接線，不宣稱 PC-98 或自然 DOS/V 整張截圖的逐像素對拍。

<!-- 缺口：無 -->

> 這份文件本身沒有未解項——內文出現的「缺口／未解」字樣指的是別處的缺口或方法論規則。
> `tools/re_open_questions.py` 靠上面那行把「真的沒有」與「抽不到」分開。
