# 126 — 指令列的四個彈出選單走同一支常式

**狀態：CONFORMED。** 進言、人事、軍團、據點四格點下去都是先跳一張
`sub_193E9` 的彈出選單。remake 先前只有進言與軍團是選單，
**人事做成了一覽表、據點直接開一覽**。

- 日期：2026-09-03
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_193E9`（`000193E9`）
  與四個 handler：`sub_16224`／`sub_16265`／`sub_1628F`／`sub_162FB`；
  指令樹見 [`../re/22`](../re/22-strategy-command-tree.md) §3
- 推論等級：**confirmed**（四個 handler 的參數都是立即值）
- 相關：[`110`](110-corps-command-menu.md)（軍團那一張）、
  [`124`](124-menu-highlight-xor.md)（反白）、
  [`125`](125-menu-box-width-from-padding.md)（框寬）

## 1. 原版做什麼

```asm
mov ax, 項數 / mov cx, TALK 索引 / mov dx, 位置 / call sub_193E9
jb  取消      ; CF=1 ⇒ 使用者取消
al = 選中的第幾項
```

`dx` 的**高 byte 是列、低 byte 是欄**，一格 16 px。四張選單：

| # | 指令 | handler | `ax` | `cx` | `dx` | 位置 | 項目 |
|---:|---|---|---:|---:|---|---|---|
| 0 | 進言 | `sub_16224` | 5 | `4Dh`（77）| `400h` | (0, 64) | 敵對提案／停戰提案／請求協助／遷都／請求君主出陣 |
| 1 | 人事 | `sub_16265` | 4 | `4Eh`（78）| `403h` | **(48, 64)** | 內政官任命／解任／外交官任命／解任 |
| 4 | 軍團 | `sub_1628F` | 2 | `4Fh`（79）| `40Ch` | (192, 64) | 位置確認／行軍指示 |
| 5 | 據點 | `sub_162FB` | 2 | `52h`（82）| **`40Fh` ⇒ (240, 64)** | 首都確認／據點一覽 |

⭐ **欄號跟著那一格走**：進言在第 0 格 ⇒ 欄 0、人事第 1 格 ⇒ 欄 3、
軍團第 4 格 ⇒ 欄 12、據點第 5 格 ⇒ 欄 15。
指令格寬 48 px ＝ 3 個 16 px 的粗格，所以**欄 ＝ 指令索引 × 3**。
四張都對上，這是四筆獨立立即值互相印證的結果，不是套公式套出來的。

### 1.1 據點的兩條出口（`sub_162FB`）

```asm
and al, al / jnz 據點一覽
; 首都確認：直接算出首都的據點記錄
        mov bx, cs:word_10CFD / mov bh, [bx+3]   ; 勢力記錄 +3 ＝ 首都編號
        xor bl, bl / shr bx,1 / shr bx,1 / shr bx,1   ; 編號 × 32
        add bx, 840h                              ; 據點表
        jmp 共同尾段
據點一覽:
        call sub_12078 / mov cx, 17h / call sub_18853   ; 狀態列 #23
        call sub_17400                            ; 據點一覽，回傳選中的記錄
        call sub_120D6 / jb 收尾                   ; 取消
共同尾段:
        mov dx, [bx+8] / mov bx, [bx+0Ah]         ; 據點 +0x08／+0x0A ＝ X／Y
        mov ax, 14h / mov cx, 0Ch / call sub_12151 ; 鏡頭 ＝ (X−20, Y−12)
        call sub_11F7F / call sub_11D46
        call sub_17E1F                            ; 據點情報視窗
收尾:   mov cx, 0FFFFh / call sub_18853            ; 清狀態列
```

⭐ **兩條出口共用尾段**：不論是首都還是自己挑的，都會**把鏡頭移過去
並開情報視窗**。所以「首都確認」＝「跳過選城那一步的據點一覽」。

鏡頭的 `(20, 12)` 與軍團的位置確認、開局鏡頭是同一組立即值
（[`52`](52-main-screen-camera-and-banner-date.md)）。

### 1.2 ⭐ 選單不會被擦掉，一覽表畫在它上面（2026-09-06）

選「據點一覽」之後，**彈出選單的框還留在畫面上**——一覽表的視窗蓋住它的
下半，露出來的那一列（`首都確認`）就一直掛在清單上緣。

實測（dosgolem，`workplace/dosgolem/advise/` 的 `h1`／`h2`／`h3`）：

| 檢查 | 結果 |
|---|---|
| 開清單之後再跑 2,000 萬道指令 | 整張畫面**逐像素不變**（五區全 0）|
| 再把游標移到別處 | 一樣不變 |

⇒ **殘影是持續的，不是還沒重畫的過渡態**；而且**清單開著時時間停住**
（`banner` 一個像素都沒動）。幾何上也對得起來：選單框 `(240,64,112,48)`、
一覽表的視窗從 `y = 80` 起，露出來的正是 `y 64..79` 這一條，
量到的差異框 `x 240..351 y 64..79` 逐格相符。

這與編成視窗畫在武將一覽上面（[`../re/30`](../re/30-corps-formation-ui.md) §1）
是同一種作法：**原版不擦，只是往上面畫**。

⭐ **所以「哪一項會留框」不是選項的屬性，是「有沒有重畫地圖」的副產物。**
同一支 `sub_193E9` 的兩條出口就分得出來：

| 出口 | 之後做什麼 | 看得到殘影嗎 |
|---|---|---|
| 據點一覽 | 只開一張清單，鏡頭不動 | **看得到**（實測）|
| 首都確認 | `sub_12151` 把鏡頭跳到首都 ＋ 開情報視窗 | **看不到**——重畫地圖把它蓋掉了 |

remake 照抄的是**那個差別**：`moveCamTo()` 會順手擦掉殘留的框，
而不是逐項寫死「這一項要留、那一項不留」。

⚠ **只有據點一覽那一條有原版擷取。** 人事與軍團的四條出口也開清單，
形狀應該一樣，但**沒比過**——列在 §4。

## 2. remake 實作

四張選單長得一樣，所以**只留一份實作**（`CLAUDE.md` §7 第 6 條）：

```go
type popupMenu struct {
	talk     int              // TALK 索引
	x, y     int              // 粗格換算好的像素
	cell     naturalCommandID // 指令列要反白哪一格
	fallback []string         // 讀不到 TALK.DAT 時用
	dispatch func(*game, int) // 選中第幾項之後往哪走
}
```

| 項目 | 位置 |
|---|---|
| 型別與三張表 | `cmd/wlgame/popupmenu.go`：`corpsPopupMenu`／`cityPopupMenu`／`personnelPopupMenu` |
| 狀態 | `game.cmdMenu`（`menu` 是 nil 就是沒開——**零值安全**，不必另外記 active 旗標）。`stale` ＝ 已經選走了但框還留在畫面上（§1.2）|
| 擦掉的時機 | `moveCamTo()`（＝原版的重畫地圖）與 `syncCommandFlow()`（流程回來了）|
| 畫序 | **選單先畫、一覽表後畫**——反過來的話露出來的那一列會蓋掉清單上緣 |
| 狀態列提示 | `openCityList` 掛 TALK #23（`sub_162FB` 的 `mov cx, 17h`，[`140`](140-status-message-box.md)）|
| 輸入／繪製 | `updatePopupMenu`／`drawPopupMenu`，兩處呼叫點不變 |
| 反白 | `activeCommandCell()` 回傳 `cmdMenu.menu.cell`（[`124`](124-menu-highlight-xor.md)）|
| 首都確認 | `beginLocateCapital()`：鏡頭 ＝ 首都 −(20,12) ＋ 情報視窗，與 `openCityList` 選中之後的尾段共用一支 |
| 差異 | 數字鍵 1–5 是 remake 加的捷徑；原版只有游標與方向鍵 |

⚠ **人事從一覽表換成選單**。四個項目與順序不動
（`funcs_16279` 的四筆與 remake 既有的四條流程逐一對上），
換掉的只有那一層外殼。

## 3. 驗證

| 方式 | 證據 |
|---|---|
| 對原版 ✅ | 軍團那一張逐像素 0 px（[`../playtest/60`](../playtest/60-corps-menu-parity.md)）——**三張走同一份程式碼**，所以那一張同時驗到了框、反白與版面 |
| 單元測試 | `TestCommandMenuAnchorsFollowCommandCell`（`cmd/wlgame`）：四張的欄 ＝ 指令索引 × 3，像素位置與 `dx` 的立即值一致 |
| 單元測試 | `TestCommandMenuLabelsComeFromTalk`：三張的標籤都取自 `TALK.DAT`，每一列的全形字數與原版相同（框寬才會對）|
| 單元測試 | `TestLocateCapitalMovesCamera`：鏡頭 ＝ 首都 −(20,12) |
| 對原版 ⚠ | 據點與人事那兩張**沒有原版截圖**。`tap:30,5,5`／`tap:10,5,5` 可以照 [`../playtest/54`](../playtest/54-menu-second-row-tap.md) 的方式拍 |

## 4. 未解

| 項目 | 現況 |
|---|---|
| ~~據點／人事兩張的逐像素對拍~~ | **拍了也比了**（[`../playtest/61`](../playtest/61-city-personnel-menu-parity.md)，2026-09-03）：兩張都是四區 0 px，選單框與反白格也各 0 px。⭐ 原版側後來又用 dosgolem 重取一次並把 `banner` 也收到 0（[`../playtest/67`](../playtest/67-dosgolem-popup-menus.md)）|
| ~~「據點一覽」列的是誰的城~~ | **玩家的**（實測，2026-09-06）：196年4月20日的原版清單是濟陰／洛陽／汜水關／滎陽／官渡／函谷關／虎牢關／陝／陳留／河南**十座**，與 `-list-cities 0` 印出來的「主 0」十座逐一相同（[`../playtest/83`](../playtest/83-city-list-parity.md)）。remake 的 `playerCities` 是對的 |
| 人事／軍團那四條出口的殘影 | 只有據點一覽比過（[`../playtest/83`](../playtest/83-city-list-parity.md)）。另外四條也開清單，remake 走同一支 `dispatchPopupMenu`，**但沒有原版擷取** |
| 進言那一張還沒併進來 | `openAdvise` 有自己的一套（五項 ＋ 說服流程）。**併之前要先確認它的取消語意一樣**，這一輪沒動 |
