# 39 — 行軍指示的三選一：戰鬥指揮／委任／解體

**狀態：CONFORMED。** 原版的流程、三個分支的寫入值與抵達時的狀態機都有
機器碼出處，remake 已實作並有單測與截圖。**版面也接回原版了**（§3.7）。

- 日期：2026-08-16（§3.7 於 2026-09-06 補）
- 出處：[`../re/45`](../re/45-corps-command-mode.md) §1–§2（流程與寫入值）、
  [`../re/64`](../re/64-corps-arrival-state-machine.md)（抵達分派、解體的五個動作）、
  [`../re/27`](../re/27-list-row-fields.md) §7（委任位元）、
  [`../re/09`](../re/09-combat.md) §2（委任怎麼影響遭遇）
- 推論等級：**confirmed（靜態）**；§3.5 的選單位置是**強證據**
  （夾制與讀取端 confirmed，兩個全域的語意由用法推）

## 1. 原版怎麼做

### 1.1 流程（`sub_17FDB`）

```
TALK #3「請指示行軍目標之據點。」
選據點（sub_1703C）；右鍵取消 → 整條流程結束
TALK #21「向{2}移動下。請下達戰鬥指示。」
選單（sub_1804E → sub_193E9，cx=0x4Ch）：
    項目數 ＝ 2；**目標據點就是自己的首都時才給 3**
    字串是 TALK #76 的三行：「　戰鬥指揮　」「　委　　任　」「　解　　體　」
    右鍵取消 → 回去重選據點
寫入：
    選 0 戰鬥指揮 → and [si], 0FBh（清委任位元）、[si+23h] = 0
    選 1 委任     → or  [si], 4  （設委任位元）、[si+23h] = 0
    選 ≥2 解體    → [si+23h] = 0x0B
    [si+0Bh] = 1            ; 計時器 ＝ 1，下一個 tick 就動
    [si+20h] = 目標據點
```

⭐ **「解體」不是隨時可選**：只有把目標指向自己的首都時才出現第三項——
解散要把兵還回預備兵池，而池子在首都。

### 1.2 抵達時（`sub_14325` 的分派）

| Stage | 抵達時做什麼 |
|---:|---|
| 0（戰鬥指揮／委任都寫 0）| 到了**首都**且兵力 < 600 點（6,000 人）→ 轉 Stage 9 |
| 9 | **補兵**：六槽退回池 → 重新分配 → 轉 Stage 3 → 重畫 |
| 11 | **解體**：目標校正成首都；到了就解散 |

解散的五個動作（`sub_14651`）：勢力軍團數 −1、六槽兵員退回預備兵池、
軍團記錄歸零、主將 `+0x17` 職務歸 0、大地圖佔用圖 −1。

完整的 12 筆分派表與 Stage 8／10 見 [`../re/64`](../re/64-corps-arrival-state-machine.md)。

### 1.3 委任只影響遭遇，不影響行軍

`sub_14E5C`（野戰）與 `sub_14ED7`（攻城）**分兩條路**：玩家是攻方走
`test byte ptr [si], 4`、是守方走 `test byte ptr [di], 4`，
設起來就退回自動判定，不跳「戰鬥指揮／委任」選單
（[`../re/09`](../re/09-combat.md) §2）。**判準是「玩家那一方」，與攻守無關。**

## 2. remake 要做什麼

| 項目 | 作法 |
|---|---|
| 選單 | `cmd/wlgame/marchmode.go`：選完目的地後跳三選一，字串取自 TALK **#76**（`sub_1804E` 的 `cx = 0x4Ch` 就是這個索引）；**第三項只在目的地 ＝ 玩家首都時出現**。版面見 §3.7 |
| 訊息 #21 | 左下角的狀態列框（[`140`](140-status-message-box.md)），與選單是兩個獨立視窗 |
| 戰鬥指揮 | `state.SetMarchMode(i, MarchCommand)`：`Delegated = false`、`Stage = 0` |
| 委任 | `MarchDelegate`：`Delegated = true`、`Stage = 0` |
| 解體 | `MarchDisband`：`Stage = 11`（`state.StageDisband`）；目標不是首都時回錯誤 |
| 下令的共同尾巴 | `Timer = 1`、`Ordered = 目的地`（`March` 已經寫了）|
| 抵達分派 | `internal/state/corpsorder.go` 的 `arriveCorps(i)`，由 `tickOneCorps` 呼叫：Stage 0 → 首都且未滿編就轉 9；Stage 9 → 補兵並轉 3；Stage 11 → 目標校正成首都，已在首都就解散 |
| 「到了」的判準 | `Node == TargetNode` **而且** 座標也到了。remake 的 `Node` 在踩到據點座標時就更新（中繼據點也算），只看它會把「經過目標據點」誤判成抵達 |
| 補兵 | 走既有的 `distributeReserves`（與編成畫面同一支）——**原版也是同一支** |
| 解散 | `World.disbandCorps(i)`：四個動作照 §1.2。第五個（佔用圖 −1）remake 不需要——佔用是每 tick 由位置推導的 |
| 遭遇 | `wantsTactical` 加一道：**玩家那一方委任中就不跳選單**。原版兩條路各檢查各自那一方（玩家是攻方看 `[si]`、是守方看 `[di]`），所以與攻守無關 |
| 取消 | 選單右鍵／ESC ＝ 回去重選據點（原版是回去重選，不是整條取消）|

**不做的**：Stage 8（等士氣）與 Stage 10（回首都補兵的路上）——
它們的寫入端都在 AI 那一側（`sub_14466`／`sub_1474A`／`sub_143AF`），
而非玩家的 Stage 0–3 決策鏈還沒解（[`../re/64`](../re/64-corps-arrival-state-machine.md) §6）。
remake 的 AI 軍團維持現有行為。

## 3. 驗證

| 方式 | 內容 |
|---|---|
| 單元測試 | `TestSetMarchModeWritesRawValues`、`TestDisbandOnlyOfferedAtCapital` |
| 單元測試 | `TestDisbandReturnsMenAndFreesLeader`：軍團數 −1、預備兵增加、`Alive=false`、主將 `Posted=false` |
| 單元測試 | `TestArriveAtCapitalDisbands`、`TestArriveAtCapitalResupplies`、`TestFullCorpsDoesNotResupply` |
| 單元測試 | `TestDelegatedPlayerSideSkipsEncounter`：四種攻守／委任組合 |
| 單元測試 | `TestMarchModeAnchorClamps`（`cmd/wlgame`）：兩道夾制的邊界值，且夾住之後右／下緣剛好 640／400 |
| 單元測試 | `TestMarchModeBoxMatchesOriginal`：兩項 `112×48`、三項 `112×64`，與 §3.6 量到的一致 |
| 截圖 | `-open-march-mode`：三個選項是 TALK #76 的原文（第三項因為目標是首都而出現），訊息 #21 在左下角的狀態列框。**加 `-march-to N` 指到別的據點就只有兩項** |
| 對原版 ✅ | [`../playtest/79`](../playtest/79-march-menu-original-layout.md)：**三項 `112×64`、兩項 `112×48`、狀態列框 `256×80`，三塊都是逐像素 0 px** |

## 3.5 選單開在游標旁邊，而且會被夾住不出畫面（2026-08-21）

`sub_1804E` 只是替 `sub_193E9` 算位置（`tools/ida_dump.py`，
`KI.EXE` SHA-256 `fffeba985231cda4…`）：

```asm
0001804F  ah = 18h − al                  ; al ＝ 項數（2 或 3）→ ah = 0x16／0x15
00018053  bx = cs:word_19898
00018058  bl > ah → bl = ah              ; ★ Y 夾住，選單不會掉出畫面底部
0001805E  dx = cs:word_19896
00018063  dl > 21h → dl = 21h            ; ★ X 夾住（33 格）
0001806A  dh = bl
0001806C  cx = 4Ch / ah = 1 / call sub_193E9
```

**`dl` ＝ X、`dh` ＝ Y，兩個都取自全域，再各夾一次上限**——
上限跟著項數變（多一項就少一格），所以那是「不讓選單超出畫面」的夾制。

⭐ **兩個都是「游標相對鏡頭的格座標」**，寫入端在 `sub_11F7F`
（confirmed，2026-09-02）：

```asm
00012020  mov     ds:9886h, cx        ; 游標 X − 鏡頭 X，夾在 0..27Fh（640 px）
00012024  mov     ds:9888h, dx        ; 游標 Y − 鏡頭 Y，夾在 0..18Fh（400 px）
00012028  mov     ax, cx / shr ax,1 ×4
00012032  mov     ds:9896h, ax        ; ★ ÷16 → 格
0001203F  mov     ds:9898h, ax        ; ★ 同上，Y
```

同一支後面還把鏡頭本身 ÷16 寫進 `ds:988Eh`／`ds:9890h`，再相加得到
游標的**絕對**格座標 `ds:989Ah`／`ds:989Ch`（後者再減 2 ＝ 橫幅那兩列）。

⚠ **交叉參考看不到這些寫入**：運算元是 `ds:9896h` 這種絕對定址，而 DS 是
函式開頭 `mov ax, cs / mov ds, ax` 設的，IDA 沒把段值傳播進來。
`tools/ida_var_writers.py` 因此報「寫 0」，`tools/ida_disp_users.py` 掃位移也是 0 處
——**兩支工具的盲區重疊在同一種寫法上**，成因與判法記在
[`../re/47`](../re/47-main-screen-window-registry.md) §5.1。

### 3.6 ⭐ 框的大小是夾制值自己說出來的（2026-09-05）

實測（dosgolem，游標停在 (320,160)）：

| 項目數 | 量到的框 |
|---:|---|
| 2（目標＝宛，不是首都）| `(320, 160, 112, 48)` |
| 3（目標＝許昌，是首都）| `(320, 160, 112, 64)` |

⇒ **框是 112 px 寬、`(項目數 + 1) × 16` px 高，左上角就是游標所在的格。**
兩個夾制值把這件事自己證完：

| 夾制 | 換算 | 意思 |
|---|---|---|
| 欄 ≤ `0x21`（33）| 33 × 16 ＋ **112** ＝ **640** | 右緣剛好貼齊畫面右邊 |
| 列 ≤ `0x18 − n` | (24 − n) × 16 ＋ **(n+1) × 16** ＝ **400** | 下緣剛好貼齊畫面下邊 |

⭐ 與指令列那三張彈出選單量到的 `112×48`（兩列，
[`../playtest/60`](../playtest/60-corps-menu-parity.md)）是**同一組幾何**——
它們走的本來就是同一支 `sub_193E9`。

原版側的兩張擷取都有了（`workplace/dosgolem/tac4/m2-menu3.png` 三項、
`workplace/dosgolem/march2/e2-menu.png` 兩項），
第一次比對的紀錄在 [`../playtest/70`](../playtest/70-dosgolem-tactical-commands.md) §4——
那一輪量到的兩項差異已於 §3.7 改掉。

### 3.7 版面接回原版（2026-09-06）

§3.6 量到的兩項差異（選單開在固定位置、訊息併進選單框當標題）都改掉了。

**一、選單走 `sub_193E9` 的幾何**，與指令列那三張共用一份實作
（[`126`](126-command-popup-menus.md)）：

```
框寬 ＝ (第一列的全形字數 + 1) × 16 ＝ (6 + 1) × 16 ＝ 112
框高 ＝ (項目數 + 1) × 16
```

TALK #76 的三行都是 6 個全形字（「　戰鬥指揮　」補到等寬），所以
`legacyChoiceRect` 算出來就是 §3.6 量到的 `112 × (n+1)×16`。
⚠ **標籤要走 `talkmenu.MenuLabels` 不是 `talkLines`**——後者會 `TrimRight`
掉行尾的全形空白，框會窄 16 px（[`124`](124-menu-highlight-xor.md)）。

**二、位置由游標算，兩道夾制照抄 `sub_1804E`**：

```
欄 ＝ min(游標X ÷ 16, 0x21)        ; 33 × 16 + 112 ＝ 640
列 ＝ min(游標Y ÷ 16, 0x18 − 項目數) ; (24−n) × 16 + (n+1) × 16 ＝ 400
選單左上角 ＝ (欄 × 16, 列 × 16)
```

⭐ **兩個上限都是「右／下緣剛好貼齊畫面」**，所以夾制值自己就把框的大小
說出來了（§3.6）。remake 讀的是滑鼠座標；原版讀的是 `word_19896`／
`word_19898`，那兩個全域由 `sub_11F7F` 每圈寫入，值同樣是**游標相對鏡頭的
格座標**（§3.5）——單位與座標框一致。

**三、訊息 #21 改走左下角的狀態列框**（[`140`](140-status-message-box.md)）：
`(0, 320, 256, 80)`、肖像 `0x93`，與選單是兩個獨立視窗。
選單收掉時一起清（原版是 `sub_18853(cx = 0FFFFh)`）。
框裡的 `\2`（據點名）畫色 `0x0B` 並補到三個全形字——五支 marker handler
都是 `al = 3`（[`../re/79`](../re/79-talk-marker-handlers.md) §2）。

⚠ **remake 的選目的地仍然是一覽表**，原版是在大地圖上用游標點一格
（`sub_1703C`，[`../re/85`](../re/85-march-target-hit-test.md)）。
這是既有的 remake 差異（[`../re/22`](../re/22-strategy-command-tree.md) §6），
本輪沒有改；影響到的只有「開選單時游標在哪」，夾制之後仍然是合法位置。

## 4. 未解

| 項目 | 現況 | 下手點 |
|---|---|---|
| ~~`word_19896`／`word_19898` 的寫入端~~ | **已解**（confirmed，2026-09-02）：`sub_11F7F` 的 `00012032`／`0001203F`，值是游標相對鏡頭的格座標。見 §3.5 | — |
| ~~remake 的選單位置與訊息窗~~ | **已改**（2026-09-06，§3.7）：選單走 `legacyChoiceRect` 的幾何 ＋ `sub_1804E` 的兩道夾制，訊息改走左下角的狀態列框（[`140`](140-status-message-box.md)）| — |
| `sub_193E9` 內部的列高與配色 | 只解出外框幾何，內部（`loc_19409`）沒逐行讀 | 反白的畫法已有 `docs/spec/124`，列高可由框高 ÷(n+1) 推但沒驗 |

<!-- 缺口：無 -->
