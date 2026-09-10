# 175 — 撞上敵人不是立刻開打：`+0x03` 先倒數 12 次（96 拍）

**狀態：CONFORMED。** 兩步都實作了，同局面對拍量到 96 拍：
remake 的軍團 50 停在 `(302,158)`——與原版停的那一格相同——倒數
12 → 1 每 8 拍減 1，96 拍之後才開打（§4）。
原版軍團要踏進「敵方軍團佔著的格子」或「別人的據點」時會**停住**，
設 `+0x00` 位元 5，`+0x03` 從 `0Ch` 開始倒數，每個軍團巡迴週期減 1，
**減到 1 的那一次才開打**。remake 一走到就結算，因此早 96 拍。

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_12708`（IDA 線性 00012708）、
  `sub_12831`（00012831）、`sub_12880`（00012880）、`sub_1264A`（0001264A）、
  `sub_14A7B`（00014A7B）、`sub_14ADE`（00014ADE）、`sub_12977`（00012977）、
  `sub_12B3C`（00012B3C）、`sub_125A3`（000125A3）
- 推論等級：confirmed（四支常式逐條讀出，且實測的 96 拍與 12 × 8 吻合）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §32（原版）、§33（接上之後）
- 相關：[`43`](43-rout-on-blocked-return.md)（`+0x03` 的另一個用途）、
  [`174`](174-relief-dispatch-does-not-move-marching-corps.md)

## 1. 原版做什麼

### 1.1 下一格有人就不走

`sub_12708` 寫座標之前先看佔用圖：

```asm
        cmp     byte ptr [di], 0      ; di ＝ 佔用圖上「要踏進去的那一格」
        jz      short loc_1273C       ; 空的 → 照常走
        mov     ds, cs:word_10D52
        mov     dx, di
        xor     ah, ah
        call    sub_12831
        jb      short loc_1273C       ; CF=1 → 還是走
        retn                          ; ⭐ CF=0 → **這一拍不動**
```

### 1.2 `sub_12831`：找出佔那一格的是誰

```asm
sub_12831:
        mov     bx, [si+0Eh]
        mov     di, 2240h             ; 軍團表
        mov     cx, 7Fh               ; 127 支
loc_1283D:
        cmp     byte ptr [di], 80h
        jb      short loc_1284C
        cmp     ax, [di+12h]          ; Y 相同？
        jnz     short loc_1284C
        cmp     dx, [di+10h]          ; X 相同？
        jz      short loc_12853
loc_1284C:
        add     di, 40h
        loop    loc_1283D
        jmp     short loc_1287B       ; 找不到 → stc（放行）
loc_12853:
        mov     cl, [si+1]
        cmp     cl, [di+1]
        jz      short loc_1287B       ; 自己人 → stc（放行）
        or      byte ptr [si], 20h    ; ⭐ 位元 5 ＝「前面卡著敵人」
        cmp     byte ptr [si+3], 1
        ja      short loc_1286C       ; > 1 → 還在倒數
        jz      short loc_12873       ; ＝ 1 → ⭐ 開打
        mov     byte ptr [si+3], 0Ch  ; ⭐ 第一次撞上：倒數設 12
        jmp     short loc_12876
loc_1286C:
        mov     al, 3
        call    sub_102F5             ; 倒數期間的音效
        jmp     short loc_12876
loc_12873:
        call    sub_14A7B             ; 開打
loc_12876:
        clc                           ; ⭐ 不放行
```

### 1.3 倒數在軍團巡迴裡減

`sub_125A3` 每次巡到一支軍團就先走 `sub_1264A`：

```asm
sub_1264A:
        test    byte ptr [si], 20h
        jnz     short loc_12658
        mov     byte ptr [si+3], 0    ; 沒卡著 → 歸零
        mov     byte ptr [si+21h], 0
        retn
loc_12658:
        dec     byte ptr [si+3]       ; ⭐ 卡著 → 每個巡迴週期減 1
        jnz     short locret_12661
        mov     byte ptr [si+3], 1    ; 減到 0 就停在 1，等下一次撞上時開打
```

⭐ **`sub_1264A` 由 `sub_125A3` 呼叫，而巡迴一圈是 8 拍**
（[`168`](168-corps-and-hour-cursors.md)），所以倒數 12 次 ＝ **96 拍**。
實測完全吻合：原版軍團 19 在拍 2,024 停住（旗標 `C5` → `F5`），
一直到拍 2,120 才開打（[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §32）。

### 1.35 據點那一條走 `sub_12880`，同一套倒數

`sub_12708` 還有第二個入口：位元 0 設著（已經走上路徑）而下一格的地形值
落在 `0CEh`–`0DDh`（**據點自己的圖塊**）時走 `sub_12880`：

```asm
sub_12880:
        mov     di, [si+0Eh]          ; 連結記錄
        cmp     byte ptr [si+0Ah], 4  ; 步進 ＝ 4（正向）？
        jnz     short loc_12892
        mov     di, es:[di+8]         ; 連結記錄 +08 ＝ B 端節點
        jmp     short loc_12896
loc_12892:
        mov     di, es:[di+6]         ; +06 ＝ A 端節點
loc_12896:
        shl     di, 1 / shl di, 1     ; → 據點記錄位址
        mov     al, [si+1]
        cmp     [di+841h], al         ; 那座據點是自己的嗎
        jz      short loc_128C7       ; 是 → stc（走進去）
        or      byte ptr [si], 20h    ; ⭐ 同一個位元 5
        cmp     byte ptr [si+3], 1
        ja      short loc_128B4       ; > 1 → 音效
        jz      short loc_128BB       ; ＝ 1 → ⭐ 攻城／接收
        mov     byte ptr [si+3], 0Ch  ; ⭐ 同一個 12
        jmp     short loc_128C2
loc_128BB:
        add     di, 840h
        call    sub_14ADE
loc_128C2:
        clc                           ; 不放行
```

⇒ **兩個入口共用位元 5 與 `+0x03`，倒數都是 12。**
差別只在倒數結束時呼叫誰：撞軍團走 `sub_14A7B`（野戰），
撞據點走 `sub_14ADE`（攻城／接收）。

⭐ 實測的那一場（軍團 19 → 據點 122）走的是**據點這一條**：
軍團停在最後一個路徑點 `(302,158)`，下一格 `(303,158)` 是 122 的城門格。

### 1.4 開打時把兩邊的旗標與倒數都清掉

```asm
sub_14A7B:
        …
        call    sub_14C72             ; 開戰前置；CF=1 就不打
        mov     di, bx
        and     byte ptr [si], 0DFh   ; 攻方清位元 5
        mov     byte ptr [si+3], 0
        and     byte ptr [di], 0DFh   ; 守方也清
        mov     byte ptr [di+3], 0
        call    sub_14E5C             ; 戰鬥判定，ah ＝ 結果
        …                             ; 贏的一方 sub_1291A（接收據點）
```

### 1.5 `+0x03` 是一個 byte、兩種用途

| 什麼時候 | 判準 | 用途 | 初值 |
|---|---|---|---|
| 敗走中 | `+0x00` ＝ `08h`（**已經不存在**，`sub_12977` 整個 byte 寫 8）| 敗走倒數 | `30h` ＝ 48 |
| 對峙中 | `+0x00` 位元 5，而且軍團**還活著**（≥ `80h`）| 對峙倒數 | `0Ch` ＝ 12 |
| 其他 | 兩者都不成立 | 恆為 **0**（`sub_1264A` 每個巡迴週期歸零）| — |

**兩種狀態互斥**：敗走的軍團旗標 < `80h`，對峙的軍團旗標 ≥ `C0h`。
原版檢查點的 44 筆活著的軍團記錄，`+0x03` **全部是 0**。

⚠ 對峙倒數還兼**動畫相位**：`sub_12B3C`（大地圖畫軍團）拿
`[si+3] & 3` ＋ `[si+21h] << 2` 當索引取圖。所以那 12 拍不只是等待，
畫面上是兩軍對砍的動畫。

## 2. remake 差在哪

`resolveContact` 一踏進敵人那一格就結算，**沒有對峙這一段**。
後果不是「早一點打完」，是**整條時間軸從那一刻起錯開 96 拍**——
同局面對拍的第一個分歧就是它（拍 2,023 對 2,119）。

⚠ 還有一個副作用：原版對峙期間軍團**停在原地**，所以那 96 拍裡
它佔的格子、據點的 `+0x18`、威脅量都跟著不同。

## 3. remake 要改什麼

分兩步做。**第一步（已完成）：把 `+0x03` 的兩種用途分開。**

| 項目 | 位置 |
|---|---|
| 狀態層 | `internal/state/corps.go`：`Corps.Standoff` ← `+0x00` 位元 5，進 `modelledCorpsBits` |
| 狀態層 | `internal/state/corps.go`：`RoutTimer` 改名 `Countdown`，一個欄位對一個 byte |
| 狀態層 | `internal/state/corps.go`：軍團巡迴每次先跑 `tickStandoff`（＝ `sub_1264A`）|
| 狀態層 | `internal/state/corps.go`：`saveCorps` 對**活著的**軍團也寫 `+0x03` |

⭐ **一個 byte 就用一個欄位。** 拆成 `RoutTimer` 與 `StandoffTimer` 兩個
Go 欄位會多出「兩個欄位、一個 byte」的失步風險，而原版本來就是靠
`+0x00` 分辨——照抄那個結構，載入、寫回、讀取三處都不可能對不上。
分開的是**用途**（哪個狀態讀它），不是儲存。

**第二步（已完成）：接上觸發。**

| 項目 | 位置 |
|---|---|
| 下一格 | `internal/state/corps.go`：`nextCell` ＝ `sub_12708` 進來時 `es:[bx]`／`es:[bx+2]` 指的那一筆（指標已在 `sub_126FF` 加過步進）|
| 兩個閘 | `internal/state/corps.go`：`blockerAt` ＝ `sub_12831`（佔用圖）＋ `sub_12880`（據點圖塊）|
| 倒數與結算 | `internal/state/corps.go`：`standoffBlocks`，由 `tickOneCorps` 在 `step` **之前**呼叫 |
| 結算入口 | `internal/state/corps.go`：`siegeAt` ＝ `sub_14ADE`、`fieldAt` ＝ `sub_14A7B`。舊的 `resolveContact`（移動**之後**才檢查同格）整支換掉 |
| 目標據點 | `fight`／`resolveCorpsBattle`／`fightGarrison`／`capture`／`afterBattle`／`beginTactical` 都加一個 `node` 參數 |
| 退卻判準 | `internal/state/aimarch.go`：`retreatOrPerish` 的「站在自家據點上」改用 `onCity`（§3.1）|

### 3.1 ⚠ 連帶要修的：`Node` 不等於「站在據點上」

對峙的軍團**還沒踏進去**，所以結算的時候：

- **目標據點找不到**。`fightGarrison`／`capture` 原本用 `w.Corps[att].Node`
  當那一場的據點——攻方對峙時 `Node` 還留在前一站，靠它會找錯。
  所以整條結算鏈改成由呼叫端把 `node` 傳下去。
  ⭐ 值在改之前之後**相同**（攻城時攻方踏進去後 `Node` 正是目標據點），
  改的是**來源**：從「猜」變成「傳」。
- ⭐ **`retreatOrPerish` 會把攻方讀成「站在自家城裡」而不退。**
  原版的判準是節點欄 `+0x0E < 800h`（[`46`](46-post-battle-retreat.md) §2），
  行軍中那一欄放的是連結記錄位址，一定 ≥ `800h`；而 remake 的 `Node`
  **在行軍中留著出發那一站**。以前攻方踏進敵城才結算，`Node` 已經換成
  敵城所以剛好躲過；改成城外對峙之後就踩到了。判準改成
  `onCity`（`Node` 是據點**而且**座標就在那個據點上）。

### 3.2 ⚠ 沒照抄的一個閘：位元 0

`sub_12880` 前面有 `test byte ptr [si], 1`（已經走上路徑才問據點）。
**remake 不照抄**，理由與判定順序倒置同源（§3.3）：原版一個據點佔
`0CEh`–`0DDh` 一整段圖塊，軍團站在自家城裡時下一格往往還在自己的據點
圖塊上，位元 0 是用來擋掉那一步的；remake 的據點只佔一個點，
出城第一步永遠踏在道路格上，沒有這個情況。
照抄反而會漏掉直線退路——缺道路圖時 `OnPath` 一次都不會設
（`step` 只在走格子路徑時設它），攻城那條就整個消失。

### 3.3 判定順序照舊倒過來

`blockerAt` **先問據點再問佔用圖**，與原版相反。理由是地圖模型不同，
在 [`../re/09`](../re/09-combat.md) §2 已經記過：本專案的據點是一個點，
守軍就站在中心，照抄順序會永遠打成野戰。

## 4. 怎麼驗

第一步（`+0x03` 分開）：

1. 五個檢查點的軍團表維持**逐 byte 相同**——原版活著的軍團 `+0x03`
   全部是 0，所以「活著也寫 `+0x03`」不能讓任何一格變。
2. 單元測試：兩種倒數**不會同時非零**；載入→寫回 `+0x03` 與位元 5 無損；
   對峙中每個巡迴週期減 1、減到 0 停在 1；不在對峙就歸零。

第二步（接上觸發）——**都量到了**（[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §33）：

| 判準 | 原版 | 接上之前 | 接上之後 |
|---|---:|---:|---:|
| 第一場 AI 自動判定戰鬥 | 拍 **2,119** | 拍 2,025（早 94 拍）| 拍 **2,121**（晚 2 拍）|
| 第二場 | 拍 2,313 | 拍 2,049 | 拍 2,311 |
| 逐拍取數不一致 | — | 369 / 5,480（6.7%）| **332 / 5,480（6.1%）** |

⭐ **對峙期間的軍團記錄**：remake 的軍團 50 停在 `(302,158)`——與原版
停的那一格**相同**——`+0x03` 在拍 2,030／2,050／2,100 是 11／8／2，
每 8 拍減 1，拍 2,121 開打。

⚠ 旗標 remake 是 `E5`、原版是 `F5`，差的是**位元 4**（§5 的未解項，
remake 沒建模）。位元 5 與倒數都對得上。

⚠ 剩下的 2 拍差來自**撞上的時機**本來就晚 2 拍，不是倒數本身——
兩邊都是「撞上之後整整 96 拍」。

五個檢查點的軍團表逐 byte 與接上之前**完全相同**（拍 200／700／1,200／
1,600／1,900 各 4／1／1／1／25 個 byte）——第一次對峙發生在拍 2,025
之後，所以這一段時間軸不受影響。

單元測試（`internal/state/standoff_test.go`）：撞上不移動、
第 12 個巡迴週期才結算、據點那條打的是攻城、擋路的走了就散、
自己人不擋。四個突變（一撞上就打／照樣走進去／自己人也擋／不問據點）
各自被擋下。

## 5. 未解

| 項目 | 現況 |
|---|---|
| `+0x00` 位元 4 | `sub_12B3C` 設、`sub_12BA8` 清，而 `sub_12BA8` 接著呼叫 `sub_19656`／`sub_196ED`（繪圖）。**像是「這一格要重畫」的髒旗標**，不是規則狀態——待確認 |
| `+0x21` | `sub_1264A` 在沒卡住時一併歸零；`sub_12B3C` 拿 `<< 2` 與 `+0x03 & 3` 合成圖塊索引。**像是對峙動畫的第二個維度**，語意未讀 |
| `sub_102F5(al=3)` | 對峙期間每個週期呼叫一次，推測是音效 |
