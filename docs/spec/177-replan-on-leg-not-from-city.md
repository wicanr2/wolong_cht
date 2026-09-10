# 177 — 走在路段上重新選路只在這條邊上決定方向

**狀態：CONFORMED。** `sub_147BB` 有兩半：軍團的 `+0x0E` **≥ `800h`**（走在
某條邊上）時只決定「往這條邊的哪一端」，`< `800h`（站在據點上）才算整條路。
remake 兩種情形都走 `March`——從**出發據點**重算整條路徑，並且順手改寫
`+0x16`／`+0x18`。後果是退卻中的軍團原地不動、目標座標與原版分歧。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_147BB`（IDA 線性 000147BB，192 B）、
  `sub_12662`（00012662）、`sub_126FF`（000126FF）、`sub_144A9`（000144A9）、
  `loc_1491B`（0001491B）
- 推論等級：confirmed（`sub_147BB` 全文逐條讀出）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §38
- 相關：[`172`](172-corps-march-fields.md)（連結記錄的版面）、
  [`173`](173-corps-flag-bit0-and-sprite-fields.md)（位元 0，本篇訂正它）、
  [`46`](46-post-battle-retreat.md)（Stage 10／11 的目標校正）

## 1. 原版做什麼

`sub_12662` 每次要移動就先呼叫 `sub_147BB`。它開頭取三個欄位：

```asm
mov ax, [si+14h]        ; 行軍目標（節點 × 8）
mov bx, [si+0Eh]        ; 現在節點
mov dl, [si+1]          ; 勢力
mov es, cs:word_19874   ; 道路表的段
cmp bx, 800h
jb  short loc_1482F     ; ← 站在據點上，走「算整條路」那一半
```

### 1.1 `≥ 800h`：走在邊上，只決定方向

```asm
mov cx, es:[bx+6]       ; 這條邊的端點 A
mov bx, es:[bx+8]       ; 端點 B
cmp ax, bx / jz loc_14821   ; 目標就是 B → al = 4    （正向）
cmp ax, cx / jz loc_14828   ; 目標就是 A → al = 0FCh （反向）
… loc_1491B 算成本，再比 dx 與 bp 決定 al = 4 或 −4 …
loc_14821:  mov al, 4    / or byte ptr [si], 1
loc_14828:  mov al, 0FCh / or byte ptr [si], 1
loc_1486C:  mov [si+0Ah], al
```

⇒ **軍團走上一條邊之後，只會在這條邊上前進或後退**，走到端點才重新選邊。
三個分支都 `or byte ptr [si], 1`（設位元 0）。

`+0x0C`（路徑點位址）與 `+0x0E`（現在節點）**這一半完全不動**——
方向換了，走的還是同一批路徑點，只是 `bx += [si+0Ah]` 的正負相反。

### 1.2 `< 800h`：站在據點上，算整條路

```asm
loc_1482F:
        cmp ax, bx / jz loc_1486F      ; 目標就是腳下 → 什麼都不做
        mov cx, bx / call loc_1491B    ; 算路
        （`cx ≥ 8000h` 而且 Stage ≥ 0Ah ⇒ 敗走，見 spec/43）
        cmp al, 4 / jnz → dx = es:[bx+2] ; 否則 dx = es:[bx]
        mov [si+0Ch], dx               ; ★ 路徑點位址 ← 這條邊的頭或尾
        mov [si+0Eh], bx               ; ★ 現在節點 ← 連結記錄位址
        and byte ptr [si], 0FEh        ; ★ 清位元 0
        mov [si+0Ah], al
```

### 1.3 ⛔ 位元 0 有清除端——`docs/spec/173` §1.1 說反了

`00014869` 的 `and byte ptr [si], 0FEh` 就是它。**位元 0 不是一次性旗標**，
而是每次移動判定都會重寫的狀態：

| 移動前的 `+0x0E` | 位元 0 | 意思 |
|---|---|---|
| `< 800h`（站在據點上）| **清** | 這一步是從據點踏出去的 |
| `≥ 800h`（走在邊上）| **設** | 這一步是從路段中間走出來的 |

`sub_12662` 緊接著就靠它分流：沒設 → `loc_126C0` 直接寫座標（出城第一步）、
設 → `loc_126C8` 檢查 leg 盡頭。

> **為什麼舊斷言會成立**：`tools/ida_bitflag_users.py` 照
> 「立即值**含**指定位元」篩，而清除位元 0 的立即值 `0FEh` 恰好**不含**
> 位元 0——清除端一條都篩不到，輸出看起來完全正常。工具已改成
> 按角色分組（`or`／`bts` 含 ＝ 設、`and`／`btr` 不含 ＝ 清），
> 而且哪一組是 0 處會明講。

### 1.4 戰後退卻不改 `+0x16`／`+0x18`，Stage 校正改

兩條路要分開：

| 常式 | 寫哪些欄位 |
|---|---|
| `sub_1474A`（戰後退卻）| `+0x14`（目標 × 8）、`+0x20`（意圖）、位元 1。**座標兩格不碰** |
| `sub_144A9`／`sub_144D6`（Stage 10／11 校正成首都）| 經 `sub_14548`，`+0x14`／`+0x16`／`+0x18` **三格一起寫** |

```asm
sub_1474A:
        call    sub_1487B          ; 下一個自己的據點（bx ＝ 840h + node×32）
        jb      短 STC             ; 退不了 ⇒ 壞滅
        shr bx,1 ×2 / mov [si+14h], bx    ; ★ 只有目標節點
        shr bx,1 ×3 / mov [si+20h], bl    ; ★ 意圖
        or      byte ptr [si], 2   ; ★ 位元 1 ＝ 下一步要重算
        cmp     word ptr [si+4], 12Ch   ; 兵 ≤ 300 ⇒ Stage 10
        jbe     短 → Stage 0Ah
        cmp     al, [bx+3]         ; 退到的據點就是首都 ⇒ Stage 10
        jnz     短 → Stage 8
```

原版同局面拍 2,200 的軍團 19：`+0x14` 已經從 122 改成 129（首都），
而 `+0x16`／`+0x18` 還停在 `(304,158)` ＝ 據點 122 的座標——
**上一次寫進去的值**。這一格是分辨兩條路的判準。

### 1.5 「已經到了」要比三個欄位

`sub_14548` 回 CF=1（`sub_144A9` 據此轉 Stage 9）的條件是
`+0x10`／`+0x12`（座標）與 `+0x0E`（節點）**三個都相同**。
只看節點會把「從首都出發、還走在路上」讀成「已經到家」——
行軍中 `+0x0E` 是連結記錄位址（≥ `800h`），根本不等於 `node × 8`。

### 1.6 重算不是當場做的——位元 1 決定時機

改行軍目標的常式（`sub_1474A`／`sub_144A9`／`sub_142AB`／`sub_17FDB`…）
只設 `+0x00` 位元 1，**方向不當場換**。換向發生在下一次「輪到移動」：

```asm
sub_12662:
        …佔用圖 −1…
        test    byte ptr [si], 2     ; 位元 1 ＝ 下一步要重算
        jz      short loc_126B4
        and     byte ptr [si], 0FDh  ; ★ 清掉
        call    sub_147BB            ; ★ 重算
        jb      short loc_126F5
        test    byte ptr [si], 1     ; 位元 0 決定走哪一條出口
        jz      short loc_126C0      ;   沒設 → 直接寫座標
        jmp     short loc_126C8      ;   設了 → 先檢查 leg 盡頭
loc_126B4:                           ; 位元 1 沒設
        cmp     word ptr [si+0Eh], 800h
        jnb     short loc_126C8      ; ★ 走在邊上 ⇒ 方向根本不重算
        call    sub_147BB            ;   站在據點上才算
loc_126C0:
        mov     bx, [si+0Ch] / call sub_12708   ; 寫座標 ＝ 走一格
```

兩件事照抄：

- **重算完那一拍照樣走一格**——`sub_147BB` 之後就落到 `sub_12708`，沒有出口。
- **位元 1 沒設而且走在邊上時，`sub_147BB` 根本不會被呼叫**：方向就照
  現有的 `+0x0A` 一路走到端點。

同局面拍 2,120 的軍團 19：原版旗標 `C7`（位元 1 還設著）、步進仍是 `FC`；
remake 在改目標的當下就換成 `04`，早了一個移動週期。

## 2. remake 錯在哪

`headHomeResupply`／`arriveDisband` 用 `w.March(i, capital)` 做目標校正，
而 `March` 是**玩家下行軍指令**的入口：

1. 它寫 `c.TargetX, c.TargetY = w.Cities[node].X, w.Cities[node].Y`——
   原版的 AI 校正不碰這兩格。
2. 它用 `w.roads.CellRouteMarked(c.Node, node)` 從**出發據點**重算整條路徑。
   軍團的座標在路段中間，重算出來的路線從據點起算，兩者對不上，
   於是 `step` 走不動——軍團原地凍住。

同局面對拍（`docs/playtest/119` §38）拍 2,200 的軍團 19：

| | 步進 | 座標 | 目標 XY |
|---|---|---|---|
| 原版 | `04`（**掉頭往回**）| `(298,160)`（走了一步）| `(304,158)` |
| remake | `FC`（沒換）| `(302,158)`（**沒動**）| `(291,165)` |

## 3. 要改哪幾支

| 檔案 | 改什麼 |
|---|---|
| `internal/state/corps.go` | 新增 `replanOnLeg`：軍團在邊上時只在這條邊上選方向（§1.1）。位元 0 的語意訂正隨 `step` 改好 |
| `internal/state/aimarch.go` | 新增 `retargetAndReplan(i, node, coords)`：`coords` 決定寫不寫座標兩格（§1.4）。`retreatOrPerish` 用 `false`、`headHomeResupply` 用 `true`。Stage 9 的判準加 `onCity`（§1.5）|
| `internal/state/corpsorder.go` | `arriveDisband` 同上；補兵判準也加 `onCity`。`retarget`（`sub_14548`）改成無條件呼叫 |
| `internal/state/corps.go` | `Corps.Replan` ＝ 位元 1，進 `modelledCorpsBits`；重算移到 `tickOneCorps` 的「輪到移動」那一刻（§1.6）|

`March` 本身不動——它是玩家指令的入口，寫 `+0x16`／`+0x18` 是對的
（原版 `sub_17FDB` 也寫）。

⚠ §1.5 那四處是 `remake-only-field` 這條教訓的**第四次**：
`Corps.Node` 在行軍中留著出發那一站，拿它當「站在哪裡」永遠會在
「從 X 出發、目標也是 X」時誤判。原版沒有這個欄位——它比的是座標。

## 4. 驗證

逐拍取數的第一個分歧從拍 **2,455 推到 2,647**。檢查點
（`tools/parity_ck.sh`）：

| 檢查點 | 據點表 | 軍團表 |
|---|---|---|
| 200／700／1,200／1,600／1,900 | **0** | **0** |
| 2,120 | **0** | 2 B／2 支（都是 `+0x0B` 移動計時）|
| 2,454 | 6 B／5 座 → **1 B／1 座** | 33 B → 17 B |
| 2,500 | 58 B／39 座 → **1 B／1 座** | 34 B → 31 B |

## 5. 未解

- **移動計時器 `+0x0B` 的相位**：拍 2,120 之後軍團 19 與 72 的計時差 1–2。
  兩支都是拍 2,119 那場攻城的當事人（19 攻、72 守），而 `arriveCorps`
  三條路的共同尾巴會把計時寫 1——戰後那一輪的計時怎麼走還沒逐拍對過。
- `loc_1491B` 在「已經在邊上」那一支回傳的 `dx` 與 `bp` 比較的語意
  （`cmp dx, bp / jz / neg al`）只讀出結果，沒讀出成本函數本身——
  那一段是自我修改碼（[`../re/65`](../re/65-ai-march-decision-chain.md) §8.1）。
  remake 先用「離目標較近的端點」代替，等距時的取捨還沒有原版收據。
