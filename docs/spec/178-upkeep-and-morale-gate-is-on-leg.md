# 178 — 軍費與士氣的判準是「有沒有走在路段上」

**狀態：CONFORMED。** `sub_12600` 的分岔是 `cmp word ptr [si+0Eh], 800h`——
**軍團有沒有走在兩個據點之間的那條路上**，不是它站在哪一種節點上。
remake 拿 `army.KindOf(Corps.Node)` 去問，而 `Node` 在行軍中留著出發那一站，
於是行軍中的軍團收便宜的軍費、**而且每小時回 10 點士氣**。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_12600`（IDA 線性 00012600–00012649）、
  `sub_1562B`（0001562B）
- 推論等級：confirmed（整支 74 B 逐條讀出）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §38
- 相關：[`50`](50-corps-upkeep-charges-funds.md)（軍費走哪一條扣款鏈）、
  [`177`](177-replan-on-leg-not-from-city.md)（同一個 `+0x0E` 判準的另一處）

## 1. 原版做什麼

```asm
sub_12600:
        cmp     cs:byte_10CF3, 1     ; 每天只有「一時」收
        jnz     retn
        mov     ax, [si+4]           ; 兵力
        cmp     word ptr [si+0Eh], 800h
        jb      short loc_12624
        ; ── 走在路段上 ──
        shr ax,1 / mov dx,ax / shr ax,1 / add ax,dx   ; 兵 × 3/4
        mov dh, [si+1] / call sub_1562B               ; 當場扣資金
        retn                                          ; ★ 士氣不回
loc_12624:
        ; ── 沒走在路段上（據點或野外節點都算）──
        mov cl,5 / shr ax,cl / inc ax                 ; 兵 ÷ 32 ＋ 1
        mov dh, [si+1] / call sub_1562B
        mov bh,[si+1] / bx = 勢力 × 64
        mov al, [bx+1Dh]             ; 勢力的士氣基準（+0x1D，開局 200）
        add byte ptr [si+6], 0Ah     ; ★ 士氣 +10
        cmp [si+6], al / jb → / mov [si+6], al        ; 夾到基準
```

⇒ 兩件事共用同一道閘：**軍費差 24 倍、士氣回不回**，
都看 `+0x0E` 有沒有 ≥ `800h`。

野外節點（`600h`–`7FFh`）走的是**便宜**那一邊——「野外」不是這條規則的分野。

## 2. remake 錯在哪

```go
inField := army.KindOf(c.Node) == army.FieldNode   // ⛔
```

`Corps.Node` 是 remake 自己的欄位，行軍中留著**出發那一站**（docs/spec/172）。
從據點出發的軍團，這個式子永遠回 false ⇒ 收便宜的軍費、每小時回 10 點士氣。

同局面對拍量到的：戰後退卻中的軍團 19，原版士氣停在 82，
remake 在 20 拍內回到 92。

`Upkeep`／`Recover` 的參數也跟著改名成 `onLeg`——原本叫 `inField`，
而註解把「道路上」歸到便宜那一邊，與原版相反。

## 3. 改了哪幾支

| 檔案 | 改什麼 |
|---|---|
| `internal/rules/combat/combat.go` | `Upkeep`／`Recover` 的參數 `inField` → `onLeg`，註解訂正 |
| `internal/state/corps.go` | 呼叫端改成 `onLeg := c.LinkAddr != 0` |

## 4. 驗證

`tools/parity_ck.sh` 的拍 2,140／2,200：軍團 19 的 `+0x06`（士氣）
從差 10 回到 0。

## 5. 未解

- `byte_10CF3`（「一時」的判準）remake 用 `hour == upkeepHour` 代替，
  兩者等價與否沒有逐拍證據。
