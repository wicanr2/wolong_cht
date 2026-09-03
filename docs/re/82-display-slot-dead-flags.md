# 82 — 戰場顯示格的旗標 bit 5／bit 6：兩版都是死碼

**狀態：confirmed（靜態，兩版各自驗證）。** 顯示格 `+0` 的 bit 5 只有
**一支**設定端與**一支**清除端，而這兩支小常式在**松崗 DOS/V 版與
PC-98 日文原版都沒有任何呼叫端**。所以 `sub_1DE95` 收尾那道
「`dl & 0x20` 成立就對五個鄰格各跑一次 `ax = 0`」——[`../spec/58`](../spec/58-display-slot-depth-range.md)
掛了半個月的未解項——**在兩版原版上都永遠不會執行**。
remake 沒實作那條路不是缺口。bit 6 同理：唯一的寫入端就在同一支死碼裡。

- 日期：2026-09-03
- 原始輸入：`workplace/orig/dosv/KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；
  `workplace/orig/pc98/KI.EXE`，SHA-256
  `061917f9f3f5c03e29397a9c636d546052128a99b8c8ce31ded0e84cf2a481e8`
- IDA database：`workplace/ida/dosv/KI.EXE.i64` SHA-256
  `65736f11b0b28a5b3a6db9e1a3d205cc24f0eaebc82b508ee0d7d283f6240572`（739 支函式）；
  `workplace/ida/pc98/KI.EXE.i64` SHA-256
  `6b89ddc239310153d594f3f617e129497dba905f1e66544c44edde2618b67324`（731 支函式）
- 工具：`tools/ida_bitflag_users.py`、`tools/ida_callers.py`、
  `tools/ida_dump.py`、`tools/ida_range.py`、`tools/ida_var_writers.py`
- 位址：各自版本的 IDA 線性位址，**兩版不互相外推**（`CLAUDE.md` §7 第 9 條）
- 相關：[`11`](11-tactical-battle.md) §5.13b（顯示格 32 B 的版面）、
  [`68`](68-t3-frontier-functions.md) §2.1（`+1`／`+3` 的分工）、
  [`../spec/58`](../spec/58-display-slot-depth-range.md)（深度範圍與帶高）

## 1. 顯示格 `+0` 那個 byte 上有五個位元

| 位元 | 誰設 | 誰清 | 語意 |
|---|---|---|---|
| **bit 3**（`0x08`）| `sub_1D9D1`（`ah = 0`）| `sub_1D9D1`（`ah ≠ 0`，同時 `or 80h`）| **這一格被浮動視窗蓋住，不要重畫**。呼叫端是戰場上的視窗族：門強度條 `sub_1C4A6`、對白框 `sub_1C39C`／`sub_1C3F0`／`sub_1C407`。`sub_1DDB4` 進到一格的第一件事就是 `test dl, 8 / jnz`，整格跳過 |
| **bit 4**（`0x10`）| `sub_1DDB4` 自己（`or dl, 10h`）| — | **只活在暫存器裡**：高度被鄰格取代時設，寫回 `[si]` 的路上沒有這個位元。它是快路徑判斷 `test dl, 50h` 的其中一半 |
| **bit 5**（`0x20`）| `0x1D98B`（DOS/V）／`0x1DEF9`（PC-98）| `0x1D9AF`（DOS/V）／`0x1DF1D`（PC-98）| **死旗標**，見 §2 |
| **bit 6**（`0x40`）| 同上兩支的 `or [bx], 0E0h`／`or [bx], 0C0h` | `sub_1DDB4` 收尾 `and [si], 3Fh` | 同樣只由死碼寫，見 §3 |
| **bit 7**（`0x80`）| `sub_1DD22`（內容變了）、`sub_1DB34`／`sub_1DB9B`／`sub_1DA1C`／`sub_1DAAA`／`sub_1DC03`（物件進出）、`sub_1D9D1`（視窗收掉）| `sub_1DD22`（`and [si], 7Fh`）、`sub_1DDB4` 收尾 | **髒旗標**：這一格要重畫 |

## 2. bit 5 的設定端與清除端

兩支連在一起、夾在 `sub_1D971`（把整塊顯示格 `rep stosw` 清成 0）與
`sub_1D9D1`（視窗遮罩）之間，**IDA 兩支都沒認成函式**——因為沒有人呼叫它們，
而 IDA 是靠 xref 建函式的。

```asm
; DOS/V 0x1D98B（PC-98 0x1DEF9，逐條同構）
        push ds
        xchg bl, bh / shr bx,1 ×3 / add bx, dx / shl bx,1 ×5   ; 格距 0x20、列距 0x400
        mov  ds, cs:word_1E15C          ; ★ 顯示格段（PC-98 是 word_1E744）
        test byte ptr [bx], 20h
        jnz  short .done
        or   byte ptr [bx], 0E0h        ; ★ bit 5 ＋ 6 ＋ 7 一起設
.done:  pop  ds / retn

; DOS/V 0x1D9AF（PC-98 0x1DF1D）
        push ds
        （同一段位址計算）
        mov  ds, cs:word_1E15C
        and  byte ptr [bx], 0DFh        ; ★ 清 bit 5
        or   byte ptr [bx], 0C0h        ; bit 6 ＋ 7
        pop  ds / retn
```

位址計算與 `sub_1D9D1`、`sub_1DB34` 完全一致（格距 `0x20`、列距 `0x400`），
段取自 `word_1E15C` ——所以它們碰的確實是**戰場顯示格**，
不是大地圖那張 8 B 一格的顯示清單（那一族在 `sub_1D4A3`／`sub_1D5D4`／
`sub_1D66A`，格距 8、列距 320，見 [`../spec/74`](../spec/74-corps-on-world-map.md)）。
兩者的 bit 5 都叫「bit 5」但**不是同一個東西**，掃描時很容易混。

### 2.1 「沒有呼叫端」是怎麼驗的

`tools/ida_callers.py` 四層都跑過，**其中三層有正對照**：

| 層 | DOS/V `0x1D98B`／`0x1D9AF` | PC-98 `0x1DEF9`／`0x1DF1D` | 正對照 |
|---|---|---|---|
| IDA code xref | 0／0 | 0／0 | `sub_1DD22` ＝ 3、`sub_1D971` ＝ 2；PC-98 `sub_1DF3F` ＝ 5 |
| 全 segment 逐條解碼的 `call`／`jmp` | 0／0 | 0／0 | 同上，數字一致 |
| 指令立即值 ＝ 該 offset | 0／0 | 0／0 | 無（見下） |
| 資料段裡的 word ＝ 該 offset | 10／0（**全是撞號**）| 0／0 | `sub_1320C` 在 `funcs_131E8` 命中 1 筆 |

DOS/V 那 10 筆是 `mov bx, cx` 的機器碼 `8B D9` 恰好等於 `0xD98B`——
**這種撞號要逐筆看反組譯才排得掉**，只看筆數會得到相反的結論。
第三層（立即值取址）沒有正對照，所以它的 0 只算弱證據；
前兩層與第四層足以定案。

## 3. 於是 `sub_1DDB4` 的兩道判斷各少一半

`sub_1DDB4` 把**自己的旗標 OR 上四個鄰格的旗標**（`or dl, al` ×4）之後：

```asm
cmp dl, 40h / jb .skip     ; 沒有 bit 6 也沒有 bit 7 ⇒ 這一格不重畫
...
test dl, 50h / jnz .slow   ; bit 4 或 bit 6 ⇒ 不能走快路徑
```

bit 6 恆為 0，所以實際成立的是：

- `cmp dl, 40h` ＝ **「自己或任一鄰格是髒的（bit 7）」才重畫**
- `test dl, 50h` ＝ **只由 bit 4（高度被鄰格取代）決定**

remake 每幀全畫，是這個條件的超集，行為一致。

## 4. 為什麼這個結論差一點是反的

第一版掃描用 `idautils.Functions()` 逐函式走，掃出來的 bit 5 設定端
**全部落在大地圖那一族**，戰場這邊一支都沒有——看起來像
「顯示格的 bit 5 沒有人設」，而那個結論剛好也是對的**答案**，
卻是用**錯的理由**得到的（[`../../CLAUDE.md`](../../CLAUDE.md) §7 第 18 條）。
真正的設定端在 IDA 沒認成函式的區段裡，逐函式掃描結構上看不到。

兩支工具因此改掉：

| 工具 | 改了什麼 |
|---|---|
| `tools/ida_bitflag_users.py` | 新增。掃「對記憶體做位元運算而立即值含指定位元」的指令，**逐 segment 不逐函式**，所屬函式印成「（無函式）」時那一欄本身就是線索 |
| `tools/ida_callers.py` | 新增。查「誰呼叫這個位址」，四層：IDA xref ／全 segment 的 `call`／`jmp` ／指令立即值 ／資料段裡的 word。⚠ 16-bit near call 的 `op.addr` 是**段內 offset 不是 linear address**，第一版拿它直接比對 linear 目標，對已知有 3 個呼叫者的 `sub_1DD22` 也回 0 ——**那是假零，靠正對照才抓得到** |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 這兩支常式原本要做什麼 | `ax = 0` 送進 `sub_1E085`／`sub_1E0E1` ＝ 拿**子圖塊 0** 對自己與四個鄰格各貼一次。子圖塊 0 在深度迴圈裡是「空」（`and ax,ax / jz` 跳過），只有這條死路徑會真的去畫它。兩版都沒有呼叫端，所以**沒有實機可以觀察**，只能說它是開發期留下的東西 |
| 大地圖那一族的 bit 5 | 是**另一個**旗標（顯示清單一格 8 B），`sub_1D4C7` 換圖時 `or byte ptr [si], 20h` 打髒，見 [`../spec/74`](../spec/74-corps-on-world-map.md)。與本份無關，並列在這裡是為了擋掉下一次的混淆 |
