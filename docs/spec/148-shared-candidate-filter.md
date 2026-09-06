# 148 — 任命與編成共用同一份候選過濾（`sub_17663`）

**狀態：CONFORMED。** 內政官任命、外交官任命、編成選武將三條流程在原版
**呼叫同一支** `sub_17663`，用**同一個建表 callback**（線性 `0x176A0`），
四個條件一次講完：存在 ∧ 勢力＝玩家 ∧ 職務＝0 ∧ 不是君主。
remake 也收成一份 `candidateGenerals()`。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_17663`（`00017663`）
  與它傳給 `sub_1820E` 的 `ax = 76A0h`（callback 本體在線性 `0x176A0`，
  IDA 沒有建成函式，所以查交叉參考查不到——要從 `mov ax, 76A0h` 的立即值追）
- 推論等級：**confirmed**（三個呼叫端 ＋ callback 逐行讀完）
- 驗證：[`../playtest/92`](../playtest/92-personnel-assign-parity.md)
- 相關：[`76`](76-lord-not-in-formation.md)（君主那一關與使用者裁定的開關）、
  [`142`](142-personnel-dismiss-flow.md)（任命的迴圈）、
  [`143`](143-general-duty-field.md)（`+0x17` 職務）、
  [`144`](144-advisor-leaves-general-table.md)（軍師不在表裡）

## 1. 原版做什麼

```asm
sub_17663:                       ; 三個呼叫端：sub_16A9B（內政官任命）
    dx = 18h / bx = 58h          ;             sub_16B71（外交官任命）
    cx = 0B18h / call sub_181C0  ;             sub_16C5E（編成選武將）
    ax = 76A0h                   ; ★ 建表 callback
    bx = 76DCh / si = 7550h / di = 75C8h
    xor cl, cl                   ; ⚠ 與 sub_175FA 的差別之一：游標歸零
    cl = cs:word_198AA
    call sub_1820E
    cs:word_198AA = cl
```

建表 callback（線性 `0x176A0`）：

```asm
    call sub_187FF               ; bx ← 君主索引 × 20h
    si = 4240h                   ; 武將表起點
    bx += si / dx = bx           ; ★ dx ＝ 君主的記錄位址
    di = 0
    al = cs:byte_10CFF           ; 玩家勢力
loc_176B5:
    cmp byte [si], 80h    / jb  next   ; ① 存在（+0x00 bit7）
    cmp al, [si+1Ch]      / jnz next   ; ② 勢力 ＝ 玩家
    cmp byte [si+17h], 0  / jnz next   ; ③ ★ 職務 ＝ 0（無職）
    cmp dx, si            / jz  next   ; ④ ★ 不是君主
    [bp+di] = si / ah++ / di += 2
next:
    si += 20h / cmp si, 5240h / jnz loc_176B5
```

`sub_187FF` 從玩家的勢力記錄取君主：

```asm
    bx = cs:word_10CFD    ; 玩家的勢力記錄（docs/spec/42、130）
    bh = [bx+1]           ; 勢力 +0x01 ＝ 君主武將編號
    bl = 0                ; bx = 編號 << 8
    bx >>= 3              ; bx = 編號 × 20h
```

⭐ **四個條件就是全部**——沒有「俘虜」那一關，因為俘虜的職務是 4
（[`143`](143-general-duty-field.md)），③ 已經擋掉了；也沒有「軍師」那一關，
因為軍師整筆不在武將表裡（[`144`](144-advisor-leaves-general-table.md)），
① 已經擋掉了。**兩個看似缺少的過濾，各自由別的機制吃掉。**

### 1.1 對照：`sub_175FA`（武將一覽）只有兩個條件

指令列「武將」那格開的是**不過濾職務**的清單（callback 在線性 `0x1763C`），
所以君主、軍團長、內政官全都在裡面——[`../playtest/91`](../playtest/91-general-boast-parity.md)
的原版擷取第一列正是「曹操　…　君王」。**兩張清單本來就不一樣**，
不可以共用同一份過濾。

## 2. remake 實作

| 項目 | 位置 |
|---|---|
| 共用過濾 | `cmd/wlgame/corps.go` 的 `candidateGenerals(excludeLord bool)`：①②③④ 四個條件一次寫完 |
| 編成 | `formCandidates()` ＝ `candidateGenerals(!g.lordCorps)`——④ 由系統選單那一列決定（[`76`](76-lord-not-in-formation.md) 的使用者裁定）|
| 人事任命 | `cmd/wlgame/personnel.go` 的 `freeGenerals()` ＝ `candidateGenerals(true)`——**一律排除君主**，沒有開關 |
| 武將一覽 | `cmd/wlgame/main.go` 的 `openGeneralList()` 照 §1.1 只用 ①②，不共用 |
| 差異 | 編成那條的 ④ 預設放行（[`76`](76-lord-not-in-formation.md) §3）。人事沒有這個差異 |

⚠ **`Captor == 0xFF` 那一關拿掉了**：原版沒有，而俘虜已經被 ③ 擋掉。
留著會在「職務 0 但 `+0x1D` 有值」這種原版不會出現的狀態上多擋一個人。

## 3. 驗證

| 方式 | 證據 |
|---|---|
| 對原版 ✅ | [`../playtest/92`](../playtest/92-personnel-assign-parity.md)：內政官任命的候選清單第一列是夏侯淵，逐像素 |
| 單元測試 | `TestCandidateGeneralsMatchesOriginalFilter`：四個條件各自的正反例 |
| 單元測試 | `TestPersonnelCandidatesAlwaysExcludeLord`：`lordCorps` 開著也不影響人事 |
| 突變測試 | 把 ④ 從 `freeGenerals` 拿掉，`TestPersonnelCandidatesAlwaysExcludeLord` 要變紅 |

## 4. 未解

| 項目 | 現況 |
|---|---|
| `sub_17663` 的 `xor cl, cl` | 比 `sub_175FA` 多一行，把清單游標歸零。remake 每次開清單本來就從 0 開始，行為相同；**但那代表原版的兩張清單共用同一個游標記憶體 `word_198AA`**，切換時的殘留還沒對過 |
