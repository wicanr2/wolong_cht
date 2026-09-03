# 42 — 戰術以外的 47 支葉節點

**狀態：47 支全部逐行讀過。四件事因此定案：`INT 61h` 是音源 TSR 的介面、
`byte_198A6` 的位元圖完整、勢力計數的減法端找到、第三處自我修改碼。**

- 日期：2026-08-14
- 範圍：只驗松崗 DOS/V
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原始 `KI.EXE.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 工具：**IDAPython** `tools/ida_dump.py`
- 位址空間：IDA DOS/V linear address，segment base `0x10000`

[`39`](39-remaining-unread.md) 列的 90 支葉節點裡，戰術那 37 支由
[`11`](11-tactical-battle.md) §3.5 與 [`36`](36-tactical-module-map.md) §6 收掉，
這一份收其餘 47 支（1,283 bytes）。

## 1. `INT 61h`：音源 TSR 的介面

四支在啟動層的小常式全部走 `INT 61h`——**那是 `YNSOUND.COM` 的軟體中斷**
（[`17`](17-dosv-audio-tsr.md)），與 `sub_1EB11` 用的 **I/O 埠 `61h`**（PC 喇叭，
[`37`](37-graphics-and-runtime-module-map.md) §1）是**兩件不同的事**。

| 函式 | B | 呼叫 |
|---|---:|---|
| `sub_10236` | 11 | `ah=8` → `INT 61h`；接著 `ah=4` → `INT 61h` |
| `sub_102C2` | 14 | `byte_1020E = 0FFh`；`ax=09F2h` → `INT 61h` |
| `sub_102D0` | 37 | 依 `byte_10CF9`（0／1／≥2）分流，`ah=7` → `INT 61h`。**那是系統選單的音效列**：0 送 `ah=8`（停）、≥1 送 `ah=0Bh, al=(值−1)×4`（主衰減，[`81`](81-sound-type-attenuation.md)）|
| `sub_1034E` | 8 | `ax=0C01h` → `INT 61h` |

> **`INT 61h` 與 port `61h` 撞號，是這一輪最容易讀錯的地方。**
> 兩者在反組譯裡長得很像（`int 61h` vs `out 61h, al`），
> 但一個是呼叫常駐程式、一個是直接踢喇叭的閘。
> **`in`／`out` 的運算元是埠，`int` 的運算元是向量號**——看助憶碼，不要看數字。

## 2. `cs:byte_198A6` 的位元圖完整了

[`31`](31-faction-picker-screen.md) §1 解出 `sub_119CA` 靠這個位元集合決定
「哪幾塊要重畫」，但只找到設定端。清除端在這一批：

| 位元 | 設定 | 清除 | 對應的重繪 |
|---:|---|---|---|
| 0 | — | `sub_161B6`（`and 0FEh`）| `sub_1614A` 指令列 |
| 1 | — | `sub_15E4C`（`and 0FDh`）| `sub_15E1E` 狀態畫面 |
| 2 | `sub_15A3A`（`or 4`）| `sub_15AA2`（`and 0FBh`）| `sub_15A3A` 勢力一覽 |
| 3 | — | — | `sub_15FAA` |

三支清除端的形狀完全相同——**清位元，然後畫一個框把那一塊擦掉**：

```asm
sub_161B6:  and ds:98A6h, 0FEh / sub_1895D(al=0, dx=0,   bx=0,  cx=021Bh)
sub_15E4C:  and ds:98A6h, 0FDh / sub_1895D(al=0, dx=1Bh, bx=0Ah, cx=0D0Dh)
sub_15AA2:  and     byte_198A6, 0FBh 之前先 sub_1895D(al=0, dx=1Bh, bx=0, cx=0A0Dh)
```

`al = 0` 是「擦除」樣式（[`24`](24-unread-function-catalogue.md) §2.1），
幾何與各自的建構常式完全一致。

`sub_15E60`（32 B）則是**讀**它：`test cs:byte_198A6, 2 / jnz → 重畫`，
與 `sub_15AD1` 的用法相同。

## 3. 勢力計數的減法端

[`../formats/08-sinario-save.md`](../formats/08-sinario-save.md) 的勢力 `+0x14`
（軍團數）與 `+0x18`（武將數）先前只找到 `inc`：

| 函式 | B | 動作 |
|---|---:|---|
| `sub_14689` | 15 | `bh = [si+1]` → 勢力記錄；**`dec byte [bx+14h]`**（軍團數 −1）|
| `sub_12AD2` | 34 | `ah != 0FFh` 時 `bh = ah` → 勢力記錄；**`dec byte [bx+18h]`**（武將數 −1）|

`sub_14689` 的三個呼叫者是 `sub_12977`／`sub_129C3`（俘虜處理）／`sub_14651`
（解除編成）——**軍團消失的三條路都經過它**。
`sub_12AD2` 對 `ah == 0xFF` 直接跳過，因為 `0xFF` 是「在野」的哨兵。

## 4. 第三處自我修改碼：`sub_12216`

```asm
sub_12216:
  push cs:word_12239                    ; 存原值
  mov  cs:word_12239, 9090h             ; ← 蓋成兩個 NOP
  call loc_1222B
  pop  cs:word_12239                    ; 還原
  retn
```

21 B，三個呼叫者。**它把 `loc_1222B` 裡的一條指令暫時換成 NOP 再呼叫**——
等於「這一次呼叫請跳過那一步」。`loc_1222B` 是 `sub_18810`（顯示訊息）
用的等待迴圈（[`25`](25-message-variants-and-personnel.md) §2），
所以被 NOP 掉的多半是「等按鍵」那一步。

前兩處是 `loc_1A065`（[`20`](20-ida-re-coverage-audit.md) §2.2）與
`loc_10701`（[`28`](28-text-number-rendering.md) §3）。
**三處的手法都不同**：一處每輪改自己、一處把參數 patch 進立即值、
這一處是「呼叫前蓋掉、呼叫後還原」。

## 5. EGA 底層的葉節點

| 函式 | B | 動作 |
|---|---:|---|
| `sub_1F426` | 21 | GC 暫存器 `5`（Mode）寫 0 |
| `sub_1F43B` | 21 | GC 暫存器 `1`（Enable Set/Reset）寫 `ah` |
| `sub_1F999` | 23 | 同上，寫 `dl` |
| `sub_1F140` | 59 | 依 `ah` bit 0 分流的位元搬移 |
| `sub_1F17B` | 39 | 逐 byte 交換（`al` 長度、`di` 目標）|
| `sub_1F3CF` | 87 | 設 `dx = 3CFh` 後的成批寫入 |
| `sub_1FA1B`／`sub_1FB11` | 28／24 | `rep movsb` 系列的搬移 |
| `sub_1E993`／`sub_1E9A7` | 20／26 | 寫 `cs:dword_1EAE9` 等一組參數槽；`sub_1E9A7` 以 `bx × 8 + 0EAF1h` 定址 → **8 bytes 一筆的參數表** |

`sub_10612`（13 B）是最小的搬移：`ds = es = cx`、`rep movsw`。

## 6. 其餘

| 函式 | B | 動作 |
|---|---:|---|
| `start` | 67 | `INT 21h AH=51h` 取 PSP，讀 `ds:80h`（命令列長度）存進 `cs:byte_10D36` |
| `sub_11F5A` | 37 | 比對 `ds:9886h`／`9888h` 與 `988Ah`／`988Ch`，不同才發系統服務 9 → **鏡頭位置的髒檢查** |
| `sub_121B2` | 53 | `sub_12151(ax=14h, cx=0Ch)` ＋ `word_19882`／`word_19886` → 鏡頭換算 |
| `sub_12804` | 4 | `al = [si+0Ah]` ＋ `cbw` |
| `sub_11C8D` | 36 | 存全部暫存器的包裝 |
| `sub_14EB9`／`sub_14F58`／`sub_14F71` | 30／25／25 | 把 `[si+2]`／`[di+2]` 推上堆疊當訊息變數，再 `sub_10CDE`／`sub_10CE7` 發聲 |
| `sub_1533D` | 27 | `ds = 0D52h` 段後轉呼叫 `sub_155A6` |
| `sub_1562B` | 16 | `si = (dx & 0FF00h) >> 2` 之後轉呼叫 `sub_1563B` |
| `sub_15AB6` | 27 | 座標夾限：`cx -= 1B8h`、`dx -= 28h`，負的話歸零 |
| `sub_15F5D`／`sub_15F7F` | 34／43 | 讀 `[si+22h]`／`si+4` 後印數字 |
| `sub_1613B` | 15 | 擦除框（`cx=0C0Dh`）|
| `sub_16605` | 30 | `sub_1304E(al=6, dx=0FFFFh)` 查事件佇列 → 訊息 `#73` |
| `sub_17028` | 20 | `xchg al, [si]` 取出並清零，再 `test al, 8` |
| `sub_1709E` | 77 | `cs:word_1989C` × 4 的換算 |
| `sub_18AD1` | 25 | `al = es:[bx] − 0DEh`，落在 `0`–`0x13` 才處理 → **字元範圍檢查** |
| `nullsub_1`／`nullsub_2`／`nullsub_3`／`nullsub_4`／`nullsub_5`／`nullsub_6`／`nullsub_7` | 1 | **空常式**，各自佔住一張分派表的一格。IDA 用 `nullsub_N` 命名，不是 `sub_XXXXX` |

## 6.1 最後四支

| 函式 | B | 動作 |
|---|---:|---|
| `sub_1061F` | 16 | 從堆疊框取 `[bp+0Ah]`／`[bp+0Ch]`／`[bp+12h]` 當 `dx`／`bx`／`si`，`di = 10h` 後轉呼叫 `sub_10B46`——**C 呼叫慣例的橋接**（參數在堆疊上，被呼叫端要暫存器）|
| `sub_109AF` | 33 | `dx = 0D58h`／`si = 18A4h`／`ah = 0C0h`／`mul ah`／`di = 0C0h` → `sub_1E38C`。**載入 `0xC0` bytes 的資源**，檔案位移 ＝ `0xC0 × ah` |
| `sub_13D45` | 35 | `sub_187AF` 之後 `sub_189A4(al=0, dx=0, bx=2, cx=151Bh)`——**擦掉訊息框**（`al=0` 是擦除樣式）|
| `sub_1C653` | 32 | **去重的待處理佇列**：`test byte [si], 10h` 已在佇列就跳過；否則 `or byte [si], 10h` 標記，把 `si` 寫進 `cs:[di-2CAEh]`，游標 `word_1D34E` 前進 2 並 `and 0FFh` 繞回 |

`sub_1C653` 的四個呼叫者都在戰術單位更新那一叢——
**單位狀態位元 4（`0x10`）＝「已排入本輪的待處理佇列」**，
佇列是 128 筆的環形緩衝（游標 `and 0FFh`、步進 2）。

## 7. 未解

| 項目 | 現況 |
|---|---|
| `INT 61h` 的四個服務號（`ah=4`／`7`／`8`、`ax=09F2h`／`0C01h`）`[DOS/BIOS]` | 對應什麼音效動作要看 `YNSOUND.COM`（[`17`](17-dosv-audio-tsr.md)）。⚠ 原版與音效 TSR 的介面，**不擋 remake** |
| `cs:byte_198A6` 位元 3 | **全庫沒有任何一處寫它**（2026-09-02 三種寫法全掃：`byte_198A6`、`cs:byte_198A6`、`ds:98A6h`，整個 `KI.EXE.asm` 只有 19 行提到這個位址）。位元 0 在 `or ds:98A6h, 1`／`and ds:98A6h, 0FEh` 成對、位元 1 在 `or/and ds:98A6h, 2/0FDh` 成對、位元 2 在 `or/and byte_198A6, 4/0FBh` 成對，**只有位元 3 兩端皆無**；整個 byte 由 `sub_11A6E` 歸零。⇒ 它要嘛只透過 `sub_119CA` 那個逐位輪詢的迴圈被讀（[`31`](31-faction-picker-screen.md) §1），要嘛根本不會被設起來 |
| `sub_1E9A7` 的 8 bytes 參數表 | 寫入端已解（2026-09-02）：整支是 `bx = bl × 8 + 0EAF1h` 之後 `cs:[bx] = ax`／`cs:[bx+2] = dx`／`cs:[bx+4] = cx`——**一筆 8 B 但只寫前 6 B**，後兩個 byte 沒有任何寫入端。**表的內容（各筆代表什麼）仍未讀** |
| `byte_1020E`／`byte_10CF9` | 音源相關的兩個旗標 |
