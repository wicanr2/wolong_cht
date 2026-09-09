# 173 — 軍團的 `+0x00` 位元 0、`+0x08`／`+0x09` 與 AI 編成的初值

**狀態：CONFORMED。** 四個欄位都接上了；同局面對拍的軍團表在
**拍 200／700／1,200／1,600 逐 byte 完全相同**（0 個 byte、0 支）。

- 日期：2026-09-09
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_126FF`（IDA 線性 000126FF）、
  `sub_12662`（00012662）、`sub_12708`（00012708）、`sub_127A2`（000127A2）、
  `sub_127F6`（000127F6）、`sub_12808`（00012808）、`sub_147BB`（000147BB）、
  `sub_14300`（00014300）、`sub_16F26`（00016F26）、`sub_16FD2`（00016FD2）
- 推論等級：confirmed（四個欄位都有唯一的寫入端，逐條反組譯讀出）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §28
- 相關：[`172`](172-corps-march-fields.md)（行軍五欄位）、
  [`74`](74-corps-on-world-map.md)（軍團圖塊 ＝ 勢力 × 5 ＋ 朝向）

## 1. 原版做什麼

### 1.1 `+0x00` 位元 0 ＝「這支軍團已經走上路徑」

唯一的設定端有兩處，都在「軍團拿到／推進路徑指標」的那一刻：

```asm
sub_126FF:                       ; 往前推一個路徑點
        or      byte ptr [si], 1 ; ← 設
        mov     al, [si+0Ah]     ; 步進（4 ／ 0FCh）
        cbw
        add     bx, ax           ; 路徑指標 += 步進
                                 ; 落下去就是 sub_12708（寫座標）

sub_147BB:                       ; 選路徑（三個分支各設一次）
        or      byte ptr [si], 1
```

**沒有清除端**——全庫掃「對記憶體做位元運算而立即值含位元 0」
（`tools/ida_bitflag_users.py`）只有 `or` 與 `test`，沒有 `and …, 0FEh`。
所以它是**一次性的**：軍團第一次被排上路徑就設起來，之後永遠留著。

兩個讀取端印證這個語意：

```asm
sub_12662:
        test    byte ptr [si], 2     ; 位元 1 ＝ 下一步要重算
        …
        call    sub_147BB            ; 重算
        test    byte ptr [si], 1     ; ← 已經走上路徑？
        jz      short loc_126C0      ;   沒有 → 直接把座標擺到現在這一點
        jmp     short loc_126C8      ;   有   → 先判斷是不是走到 leg 盡頭

sub_12708:
        test    byte ptr [si], 1     ; 沒走上路徑就不做地形 0CEh–0DDh 的處理
```

### 1.2 `+0x08` 朝向：四個方向 ＋ 「停著」

```asm
sub_12808:                       ; 由 sub_12804／sub_127F6 呼叫
        add     ax, [si+0Ch]     ; ax ＝ 步進，bx ＝ 下一個路徑點
        mov     bx, ax
        mov     ax, [si+10h]
        sub     ax, es:[bx]      ; 現在 X − 那一點的 X
        jz      short .yaxis
        and     ax, 8000h
        rol     ax, 1            ; 0（往 X 增）／1（往 X 減）
        mov     [si+8], al
        retn
.yaxis: mov     al, [si+12h]
        sub     al, es:[bx+2]
        jz      short .keep      ; 完全沒動就**不寫**
        and     al, 80h
        rol     al, 1
        add     al, 2            ; 2（往 Y 增）／3（往 Y 減）
        mov     [si+8], al
```

⭐ **X 有差就只看 X，X 相同才看 Y**；兩軸都沒動就保留舊值。

「停著」是 **4**，由三處寫：`sub_12662`（現在節點 ＝ 目標節點）、
`sub_127A2`（換節點到據點）、`sub_16F26`（編成）。
`sub_127F6`（leg 邊界掉頭）另外做 `xor byte ptr [si+8], 1`——
0↔1 與 2↔3 剛好是同一軸的反向。

⚠ **這個欄位有規則意義，不只是圖塊。** `sub_14300`（行軍中經過受威脅
據點就留守）第二道閘就是 `cmp byte ptr [si+8], 4 / jz 不留守`——
**停著的軍團不會被留下來**。

### 1.3 `+0x09` ＝ 勢力編號 × 5

```asm
sub_16FD2:                       ; 編成的收尾（sub_16E8F 等四處呼叫）
        inc     ch
        mov     [si+1Eh], ch     ; 間隔
        mov     bh, [si+1]       ; 勢力
        xor     bl, bl
        shr     bx, 1
        shr     bx, 1            ; bx ＝ 勢力 × 64 ＝ 勢力記錄位址
        mov     al, [bx+3Eh]     ; 勢力 +0x3E ＝ 勢力編號（＝自身索引）
        mov     ah, al
        shl     al, 1
        shl     al, 1
        add     al, ah           ; al × 5
        mov     [si+9], al
        mov     byte ptr [si+0Bh], 1   ; 計時 ← 1
```

`sub_12B2A` 把兩個欄位加起來取圖塊：`al = [si+9] + [si+8]`——
**每個勢力五張圖**（四方向 ＋ 停著），22 個勢力共 110 張，
與 [`74`](74-corps-on-world-map.md) §3 的 `world.CorpsTile` 同一條算式。

### 1.4 編成的計時是 1，不是間隔

`sub_16FD2` 收尾寫 `[si+0Bh] = 1`：**新編的軍團下一拍就輪得到**。

## 2. remake 差在哪

| 欄位 | remake 現況 | 對拍看到的 |
|---|---|---|
| `+0x00` 位元 0 | 完全沒建模；寫回時 `modelledCorpsBits` 保留存檔原值，所以**存檔時還沒走過的軍團永遠不會設它** | 軍團 39 從拍 200 起 `C1` 對 `C0` |
| `+0x08` 朝向 | 值域與算式都對（`headingTo` ＝ `sub_12808`、`HeadingStill` ＝ 4），但 **AI 編成走的是另一條路**（`autoFormCorps` 用 `var c Corps`），拿到 Go 零值 0 ＝「往 X 減」 | 軍團 50／56 從拍 1,600 起 `04` 對 `00` |
| `+0x09` | 概念在 `world.CorpsTile` 裡，但**存檔一次都沒寫** | 同上，`05` 對 `00` |
| `+0x0B` 計時 | `FormCorps` 寫 1（對），`autoFormCorps` 寫 `Interval` | 同上，`03` 對 `02` |

⭐ **三個裡有兩個是同一個成因**：AI 編成把玩家編成的初始化抄了一份，
抄漏的欄位就用 Go 的零值頂上（`CLAUDE.md` §7 第 11 條、第 6 條）。
修法是把共用的初值抽成一支，兩條路都走它。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 狀態層 | `internal/state/corps.go`：`Corps.OnPath`、`newCorpsRecord`、`saveCorps` 寫 `+0x09` |
| 狀態層 | `internal/state/strategy.go`：`autoFormCorps` 改用 `newCorpsRecord` |
| 差異 | 無（四個欄位都照原版） |

`+0x09` 是導出值（勢力 × 5），所以只在寫回時算，不進 `Corps`。
每組五張這個常數在 `internal/assets/world` 也有一份
（`world.CorpsHeadings`），狀態層不依賴資產層，因此以註解互指。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 同局面對拍的軍團表逐 byte | `tools/corps_diff.py`：拍 200／700／1,200／1,600 **0 個 byte**（改前 1／1／1／7）。拍 1,900 剩 12 個 byte／2 支，那是走位已經岔開的下游 |
| 單元測試 | `internal/state`、`cmd/wlgame` 全綠 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 軍團 `+0x00` 位元 4／5 | 有成對的設定與清除端，語意未定（[`../re/34`](../re/34-corps-status-bits.md) §2）|
| `sub_12708` 的地形 `0CEh`–`0DDh` | 位元 0 設著時才走 `sub_12880`（讀連結記錄的兩端節點）。那一段地形是什麼還沒讀 |
