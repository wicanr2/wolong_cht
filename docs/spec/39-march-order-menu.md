# 39 — 行軍指示的三選一：戰鬥指揮／委任／解體

**狀態：CONFORMED。** 原版的流程、三個分支的寫入值與抵達時的狀態機都有
機器碼出處，remake 已實作並有單測與截圖。

- 日期：2026-08-16
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
| 選單 | `cmd/wlgame/marchmode.go`：選完目的地後跳三選一，字串取自 TALK **#76**（`sub_1804E` 的 `cx = 0x4Ch` 就是這個索引）；標題取 TALK **#21**；**第三項只在目的地 ＝ 玩家首都時出現** |
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
| 單元測試 | `TestMarchModeRowsAndHitTest`（`cmd/wlgame`）：三列不重疊、兩項時第三列不可點 |
| 截圖 | `-open-march-mode`：標題是 TALK #21 的原文、三個選項是 TALK #76 的原文，第三項因為目標是首都而出現 |

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

⚠ `word_19896`／`word_19898` **沒有任何直接寫入的交叉參考**
（`tools/ida_var_writers.py`：各 2 筆參考、寫 0）。讀取端只有兩處，
而**兩處都是 `sub_193E9` 的呼叫端**（`sub_11F0E` 與 `sub_1804E`）。
⇒ **強證據：那是「上一次點擊的位置」**，由指標間接寫入所以掃不到寫入端
（`CLAUDE.md` §4.1 第 4 條的形狀）。

**remake 差異**：選單畫在固定位置，不跟游標走。

## 4. 未解

| 項目 | 現況 | 下手點 |
|---|---|---|
| `word_19896`／`word_19898` 的寫入端 | 語意是強證據（見 §3.5），但**直接 xref 掃不到寫入端**，所以「就是滑鼠座標」還沒 confirmed | 找取址的那幾筆，或用 DOSBox-X bridge 下寫入中斷點 |
