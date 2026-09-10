# 179 — 戰後重算把移動計時寫 1；「到了沒」比的是節點欄不是 `Corps.Node`

**狀態：CONFORMED。** 兩件小事，各卡住一段對拍：

1. `sub_1474A`（戰後處理）**第一行就是 `call sub_16FD2`**，而那支的收尾
   `mov byte ptr [si+0Bh], 1` 把移動計時寫成 1——勝敗都寫。
2. `sub_12662` 開頭的 `cmp bx, [si+14h]` 比的是 **`+0x0E` 對 `+0x14`**。
   行軍中 `+0x0E` 是連結記錄位址（≥ `800h`），所以「還在路上」永遠不相等。
   remake 拿 `Corps.Node` 比，而它留著出發那一站——「從 X 出發、目標也是 X」
   （退卻回首都、掉頭）**每走一格都被判成抵達**。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_16FD2`（IDA 線性 00016FD2–00017028）、
  `sub_1474A`（0001474A）、`sub_12662`（00012662）、`sub_125A3`（000125A3）
- 推論等級：confirmed（三支逐條讀出）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §39
- 相關：[`177`](177-replan-on-leg-not-from-city.md)（同一批 `+0x0E` 判準）、
  [`../re/30`](../re/30-corps-formation-ui.md) §5（`sub_16FD2` 的另外兩個公式）

## 1. `sub_16FD2` 的第三件事

`docs/re/30` §5 記了兩個公式（總兵力、移動間隔），漏了收尾那兩行：

```asm
00017006  mov bh, [si+1] / xor bl,bl / shr bx,1 ×2
0001700F  mov al, [bx+3Eh] / mov ah,al / shl al,1 ×2 / add al,ah
0001701A  mov [si+9], al               ; +0x09 ＝ 勢力欄的衍生值 × 5
0001701D  mov byte ptr [si+0Bh], 1     ; ★ 移動計時寫 1
```

**不是寫成間隔，是寫 1**——重算過的軍團**下一次輪到就走**。

四個呼叫者裡有 `sub_1474A`（戰後）、`sub_14499`（補兵）、`sub_16E8F`（AI 編成）。
所以一場戰鬥之後攻守雙方的計時器會**同步歸到 1**。

⭐ 對同局面對拍是看得見的：少了它，攻守雙方的移動節拍各自漂掉。
同局面拍 2,120 的軍團 19（攻）與 72（守），計時器原本差 2 與 1，
接上之後兩邊逐 byte 相同。

⚠ 順序也要照抄：`sub_1474A` **先**重算 `+0x04` 總兵力，**再**用
`cmp word ptr [si+4], 12Ch` 判 Stage 10——判準用的是重算後的值。

## 2. 「到了沒」的判準

原版：

```asm
sub_12662:
        mov bx, [si+0Eh]
        cmp bx, [si+14h]       ; 現在節點 == 行軍目標？
        jnz short loc_12675    ;   不是 → 移動那一半
        mov byte ptr [si+8], 4 / call sub_128F4 / call sub_14325
```

remake 的等價寫法是 `atTargetNode`：

```go
c.LinkAddr == 0 && c.Node == c.TargetNode
   && c.X == w.Cities[c.TargetNode].X && c.Y == w.Cities[c.TargetNode].Y
```

三段各有理由：

- **`LinkAddr == 0`** 對應 `+0x0E < 800h`（不在路段上）。
- **`Node == TargetNode`** 是原版那一比。
- **座標**是缺道路圖時的退路：`step` 這時退回直線逼近、`LinkAddr` 恆為 0，
  上面那道閘就失效了。比的是**目標據點的座標**，不是 `TargetX`／`TargetY`
  ——後兩格在戰後退卻時留著舊目標的值（[`177`](177-replan-on-leg-not-from-city.md) §1.4）。

⛔ 只比 `Node` 的後果不是「差一格」：退卻中的軍團每走一格都跑一次
`arriveCorps`，於是 `sub_144A9` → `sub_14548` 每一格都把三個目標欄位重寫，
`+0x16`／`+0x18` 跟著漂成首都座標。

## 3. 改了哪幾支

| 檔案 | 改什麼 |
|---|---|
| `internal/state/corps.go` | 新增 `recalcCorps`（`sub_16FD2`）與 `atTargetNode`；`tickOneCorps` 兩處判準都換過去 |
| `internal/state/aimarch.go` | `retreatOrPerish` 開頭呼叫 `recalcCorps` |

三個測試的 fixture 跟著改：`faceOff` 要讓攻方處在「還沒到目標」，
`TestWeakLoserSwitchesToHomeResupply` 要設**槽位**而不是 `Men`
（重算會蓋掉），`TestMarchIntoDefendedCityIsSiege` 的攻方 `Node`
不能等於目標。

## 4. 驗證

`tools/parity_ck.sh`：**拍 200／700／1,200／1,600／1,900／2,120／2,140／2,200
全部 0 個 byte**，拍 2,300 只剩軍團 50 的位元 4（繪圖狀態，規則層不設）。

逐拍取數不一致 210 → **178 / 5,480**。

## 5. 未解

- 拍 2,454 起據點 29 的佔用數差 2、軍團 20 的目標與意圖不同，
  往前夾逼還沒做到。
