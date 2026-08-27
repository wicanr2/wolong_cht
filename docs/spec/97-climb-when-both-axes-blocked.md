# 97 — 登城的觸發：X 與 Y 都走不動就試 Z，不必先走到目標格

**狀態：CONFORMED。** remake 多加了「X、Y 都已經等於目標」這個條件，
於是被擋住的登城部隊永遠不會嘗試往上爬。順帶修掉同一條路上的第二個錯：
**命令 3 的目標 Z 取的是腳下的地面層，而原版取的是高平面那一層（牆頂）**
——牆腳有地面的那幾格會算出「目前 Z ＝ 目標 Z」，純 Z 移動一次都不會試。
⚠ 兩處都改完之後，據點 82 那條 fixture 的**行為沒有變**：那一隊的 Y 落在
28–35，而那張圖的門在 Y 6／20／41／54，兵根本走不到門格
（[`../playtest/52`](../playtest/52-siege-timeseries-parity.md) §6）。

- 日期：2026-08-27
- 出處：`sub_1AF69`（`0001AF69`–`0001AFFC`，本輪由
  `tools/ida.sh script dosv tools/ida_range.py` dump）
- 推論等級：**confirmed**（三軸的分派鏈逐條讀出來）

## 1. 原版做什麼

```asm
0001AF69  mov ah, 1                    ; ah ＝「這一幀完全沒動」
0001AF6E  al = [si+6]                  ; 目前 X
0001AF71  cmp al, [si+10h]             ; 目標 X
0001AF74  ja  loc_1AFDD                ;   → sub_1B047（−X）
0001AF76  jb  loc_1AFE5                ;   → sub_1B069（+X）
0001AF78  al = [si+8]                  ; 目前 Y
0001AF7B  cmp al, [si+11h]             ; 目標 Y
0001AF7E  ja  loc_1AFED                ;   → sub_1B08B（−Y）
0001AF80  jb  loc_1AFF5                ;   → sub_1B0AF（+Y）
0001AF82  cmp byte ptr [si+4], 12h     ; ★ 兵種 ≤ 0x12（大將與騎馬）不做 Z
0001AF86  jbe short loc_1AF92
0001AF88  al = [si+0Ah]                ; 目前 Z
0001AF8B  cmp al, [si+12h]             ; 目標 Z
0001AF8E  jb  loc_1AFFD                ;   → sub_1B0D3（爬上去）
0001AF90  ja  loc_1B005                ;   → sub_1B116（爬下來）
0001AF92  and ah, ah / jz loc_1AFD0    ; 動過就結束
0001AF96  call sub_1B00D               ; 完全沒動 → 取下一個繞路點
```

每一個軸的處理長這樣：

```asm
0001AFDD  call sub_1B047               ; 走 −X
0001AFE0  mov  ah, 0                   ; 標記「試過」
0001AFE2  jb   short loc_1AF78         ; ★ 走**失敗**（CF=1）才往下試 Y
0001AFE4  retn                         ;   走成功 → 這一幀就結束
```

⭐ **落到 Z 那一段的條件只有一個：X 與 Y 這一幀都沒走成。**
`0001AF78`（Y）與 `0001AF82`（Z）是**同一條 `jb` 鏈的下一站**——
既可以從「X 走失敗」跳過來，也可以從「X 已經等於目標」直接落下來。
**原版沒有要求 X、Y 先等於目標。**

而 `sub_1B0D3` 自己會擋：腳下那一格的圖塊 < `0xF0`（不是門）就失敗。
所以「只有站在門格上才爬得上去」這一條在被呼叫的那一支，不在呼叫端。

## 2. remake 現況

```go
if !moved && s.X == s.StepX && s.Y == s.StepY && s.Z != s.StepZ {
    if b.tryClimb(side, k) { moved = true }
}
```

`s.X == s.StepX && s.Y == s.StepY` 是 remake 自己加的。後果：

腳本在攻城段對兵種 2／3 下命令 3（登城），目標是
**(索引第二欄, 自己的 Y)**（[`../re/11`](../re/11-tactical-battle.md) §5.8i）。
兵往那個 X 走，中途**站在門格上**時被前面的城壁擋住——
這時原版會試 Z 並爬上去，remake 因為「X 還沒到目標」而不試，
於是繼續撞牆。`tryMove` 把每一次撞擊算成一點耐久，
**城壁被磨穿，門強度條一路亮著**。

原版七張影格的門強度條**一次都沒出現**
（[`../playtest/52`](../playtest/52-siege-timeseries-parity.md) §2）。

## 3. 演算法

```
X 差 → 走 X；走成功就結束這一幀
Y 差 → 走 Y；走成功就結束這一幀
兵種 > 0x12 且 Z 差 → 爬（成不成由腳下的圖塊決定）
三軸都沒動 → 取下一個繞路點
```

## 4. remake 實作

| 項目 | 位置 |
|---|---|
| 規則層 | `internal/rules/tactical/soldier.go` 的 `moveToward`：Z 那一段的守衛只留 `!moved && s.Z != s.StepZ` |
| 目標 Z | 同一支檔案的 `doScaleWall`：改讀 `GroundLevel(x, y, PlaneHigh)`（原版 `loc_1AB39` 的 `bh \|= 10h` ＋ `al = es:[bx] & 7`）|
| 差異 | 無 |

## 5. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestClimbTriedWhenBlockedBeforeReachingGoal`（`internal/rules/tactical`）。做過正對照——**改回舊守衛會紅** |
| 對原版 | [`../playtest/52`](../playtest/52-siege-timeseries-parity.md)：門強度條要跟原版一樣不出現 |

## 6. 未解

| 項目 | 現況 |
|---|---|

<!-- 缺口：無 -->
