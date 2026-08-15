# 17 — 松崗 DOS/V 音源 TSR 與戰術效果碼

**狀態：INT 61h 介面、遊戲端效果碼與硬體 register 寫入已證實。
⭐ 晶片是 **OPL3（YMF262）**，見 [`57`](57-opl3-register-map.md)；
`*BGM.DAT` 的容器見 [`23`](23-bgm-resource-format.md)、事件編碼見
[`56`](56-bgm-track-events.md)。**

- 日期：2026-08-12
- 原始輸入：`workplace/orig/dosv/YNSOUND.COM`，SHA-256
  `e2c6a6a8576c4f2a96b7e3f156d7f48c9570ae03539fe9367adb78aebb364fa1`
- 關聯輸入：`KI.EXE` SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；
  `SOUND.DAT` SHA-256
  `b5624388d1bc8f6bb32aeff1d19b1da4c010d7b5a9682a249f1b807c27f2cdba`；
  `BGM.DAT` SHA-256
  `7a51c8b9a349b9e088f3796b70c268181c60bcebead70942f00e1621523dedc9`
- 工具／位址基準：IDA Pro 9.4，COM 載入後 IDA 線性位址（segment base `0x10000`）。
  `tools/ida_int61_callers.idc` 為唯讀匯出器；`objdump` 只用來交叉解碼 IDA 尚未建立
  函式邊界的 handler bytes，不能取代 IDA 的關係圖。

## 1. `YNSOUND.COM` 安裝 INT 61h

`start_0` 以 DOS 的 INT 21h 取得／設定 INT 61h 向量，將 handler 指向 COM offset
`0x0103`，之後以 DOS AH=`0x31` 常駐。這代表 `KI.EXE` 的遊戲音效不是直接呼叫一張
靜態音效表，而是先安裝的 TSR 服務。

handler 的原始 bytes 在 IDA 線性 `00010103`。經 IDA 讀取並以第二工具交叉解碼後為：

```asm
cli
push bx
mov bl, ah
shl bl, 1
xor bh, bh
add bx, 0115h
call cs:[bx]
pop bx
sti
iret
```

因此 AH 是 command index，word table 位於 COM offset `0x0115`（IDA 線性
`00010115`）。這是 **已證實** 的 dispatch 介面，不是以音色聽感推測。

## 2. 已定位的 command

| AH | handler COM offset | 已證實作用 | 證據等級 |
|---:|---:|---|---|
| `0x05` | `0x01DC` | 呼叫內部 `0x07E6` 選擇／播放遊戲效果 code | 已證實的控制流 |
| `0x06` | `0x01E2` | 保存 `DS:SI` 資源指標，並寫入 mode byte `CS:0x0A4C=6` | 已證實 |
| `0x07` | `0x01F6` | 從保存來源讀 offset words，初始化多張內部表 | 已證實為資源 layout parser；結構見 [`23`](23-bgm-resource-format.md) |
| `0x0A` | `0x02D7` | 回傳 `CS:0x099E` 的 flags 至 AL | 已證實 |

其他 table entries 存在，但本輪未把沒有資料流／呼叫端證據的 entry 命名為音樂、停止或音量。

## 3. `KI.EXE` 呼叫端與戰術效果碼

以下皆是同一份 `KI.EXE` 的 IDA linear address：

- `sub_10210` 依序以 AH=`0x00`、`0x02`、`0x03` 呼叫 INT 61h。
- `sub_10241` 先以 AH=`0x08` 停止／重設既有狀態，載入一個外部資源，再以 AH=`0x06`
  交出 `DS:SI`、AH=`0x07` 初始化；可確認這是一條資源載入路徑，不能據此宣稱完整
  `BGM.DAT` 格式已解出。
- `sub_102F5` 先以 AH=`0x0A` 查 flags，再以 AH=`0x05`、AL=效果 code 送至 TSR。

戰術的已證實 code 對照如下：

| 原始路徑 | AL code | remake 對照範圍 |
|---|---:|---|
| `sub_1AD2D` 普通投射物發射 | `0x0C` | 普通投射物發射事件 |
| `sub_1AD7F` 特殊投射物發射 | `0x0A` | 特殊投射物發射事件 |
| `sub_1B97E` 投射物命中 | `0x0B` | 投射物命中事件 |

這些是 effect code，不是已命名的可聽音效；remake 不能在未解 `SOUND.DAT` 的情況下把它們
斷言為某種特定樂器或樣本。

## 4. 硬體 register path

TSR 的 `sub_6BF` 取 channel 值後，以 `sub_890` 送出兩組 register family：

```text
register = 0xA0 + channel, data = CL
register = 0xB0 + channel, data = CH
```

`sub_890` 先把 register byte 寫到 `DX`，執行延遲讀取，再把 data byte 寫到 `DX+1`。
原始 data word `word_1097C` 預設為 `0x0220`，緊接著的 word 為 `0x0330`；啟動 `/A`
參數的結果會更新前者。

⭐ **晶片已定案：OPL3（YMF262）**，六個聲軌各佔一組 4-operator 通道，
音效走剩下的三個 2-operator 通道。決定性證據是初始化寫的 `0x104`／`0x105`
兩個暫存器（OPL2 沒有），完整推導見 [`57`](57-opl3-register-map.md) §1。

- **已證實：** address/data register pair、`A0`／`B0` channel family、
  `0x222` 是 OPL3 的第二組暫存器（不是第二顆晶片）。
- **未知：** `0x330` 的用途（MPU-401 的標準埠，但沒找到讀它的地方）。

## 5. UI 點擊音是另一條路

`KI.EXE` 的 `sub_10CDE` → `sub_1EB11` 直接操作 PC speaker port `0x61`，與 INT 61h
TSR 路徑分離。故「戰術效果」與「介面／TALK click」不能共用同一個已證實的音源實作。

## 6. remake 邊界

本輪不合成新音色，也不將未解的 `SOUND.DAT` 猜成既定格式；`BGM.DAT` 的結構雖已解
（[`23`](23-bgm-resource-format.md)），聲軌事件編碼未解，一樣不能宣稱音色 parity。若 remake 新增音效，
應先保留 `0x0A`／`0x0B`／`0x0C` 這個已證實事件層，再把實際聲音標成替代資產；只有取得
可重播的原版輸出或完整 resource parser 後，才可宣稱音色／時序 parity。

## 7. 未解

| 項目 | 現況 |
|---|---|
| `0x330` 的用途 | MPU-401 的標準埠，沒找到讀它的地方 |
| 效果碼 ↔ 聽起來像什麼 | `SOUND.DAT` 的記錄結構已解（[`57`](57-opl3-register-map.md) §6），但哪一號對應哪個動作只有 §3 的三個 |
| `INT 61h` 的四個服務號 | `ah=4`／`7`／`8` 與 `ax=09F2h`／`0C01h`，對應什麼動作要看 `YNSOUND.COM`（[`42`](42-leaf-functions.md) §7）|
