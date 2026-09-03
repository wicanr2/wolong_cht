# 122 — 音效的 TYPE 1–4：四段主衰減

**狀態：CONFORMED。** 系統選單「音　效」那一列原版有五個選項
（ＯＦＦ／TYPE 1／2／3／4），remake 先前只有開與關。
四個 TYPE 在原版是**四段音量**——每段 4 個 OPL Total Level 單位——
所以 remake 接成四段主衰減。

- 日期：2026-09-03
- 出處：[`../re/81`](../re/81-sound-type-attenuation.md)
  （`KI.EXE` 的 `sub_102D0`＝IDA `000102D0`；`YNSOUND.COM` 的
  `AH=0Bh` handler ＝記憶體 `0x02DE`、TL 計算在 `0x060F`）
- 推論等級：**confirmed（原始 bytes 交叉解碼）**
- 相關：[`29`](29-audio.md)（音訊管線）、[`13`](13-main-window-toggles.md) §5

## 1. 原版做什麼

設定值存在 `ds:0CF9h`，0–4（[`../re/55`](../re/55-system-menu-window.md) §4）。
點一下往下一個選項，到頂繞回 0。

```
0（ＯＦＦ） → INT 61h AH=8         停止並清空
1（TYPE 1） → INT 61h AH=7         初始化
              INT 61h AH=0Bh, AL=0
2–4         → INT 61h AH=0Bh, AL=(值−1)×4    ⇒ 4／8／12
```

TSR 把 `AL` 存進 `cs:[0996h]`，載波的 Total Level 就變成：

```
atten = (15 − 聲軌音量) × 4 + AL        上限 0x3F
TL    = min(音色的 TL + atten, 0x3F)    只改載波，調變器寫原值
```

OPL 的 TL 是 **0.75 dB/step**，所以一段 4 步 ≈ **3 dB**：
TYPE 1 最大聲，TYPE 4 比它小約 9 dB。

⭐ **`AH=7` 只在 0 → 1 那一步送。** 選項是環狀遞增的，要到 TYPE 2
一定先經過 TYPE 1，所以 `al ≥ 2` 那一支不必再初始化一次。

## 2. remake 的差異：主增益，不是逐載波的 TL

remake 播的是 `tools/bgm2ogg.sh` 事先算好的 OGG（[`29`](29-audio.md)），
執行期沒有 OPL 暫存器可以改，所以型別接成**播放增益**：

```
gain(n) = 10^(−(4n × 0.75) / 20)  = 10^(−0.15n)      n = TYPE − 1

TYPE 1 → 1.0000    TYPE 3 → 0.5012
TYPE 2 → 0.7079    TYPE 4 → 0.3548
```

⚠ 兩處與原版不同，都是 OGG 這條路徑帶來的：

| 項目 | 原版 | remake |
|---|---|---|
| 作用點 | 每個載波的 TL | 整個輸出的增益 |
| 飽和 | 加到 `0x3F` 就靜音，**本來就安靜的聲軌先消失** | 等比縮小，不會有聲軌先消失 |

⭐ **在 FM 裡衰減載波等於縮小那個聲部的輸出**，而四個聲部加的量相同，
所以「主增益」在數學上就是原版的效果——差別只在那道 `0x3F` 飽和。
真要連飽和一起對，得讓 `bgm2ogg` 為四個型別各算一份，
**成本（四倍的音檔）遠大於收益**，不做。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 增益 | `internal/ui/sound`：`Bank.SetLevel`／`Level`／`LevelGain` |
| 套用 | `PlayMusic`／`PlayEffect` 建 player 之後 `SetVolume`；`SetLevel` 對正在播的那一個立刻套 |
| 選單 | `cmd/wlgame/strategyhud.go`：`soundValue` 回 `TYPE 1`–`TYPE 4`；`cycleSound` 走 0–4 環狀 |
| 差異 | §2 兩條；另外「未接入」是 remake 才有的狀態（[`29`](29-audio.md) §5）|

⚠ **預設仍是 TYPE 1**（`level = 0`）——[`../playtest/39`](../playtest/39-system-window-parity.md)
的逐像素對拍就是靠這一格顯示「TYPE 1」，改預設會讓那一張退回 FAIL。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestLevelGainMatchesOPLSteps`（`internal/ui/sound`）：四段增益等於 `10^(−0.15n)`，且嚴格遞減 |
| 單元測試 | `TestSetLevelClamps`（同上）：超出 0–3 夾住，nil Bank 不炸 |
| 單元測試 | `TestCycleSoundWalksOriginalFiveOptions`（`cmd/wlgame`）：左鍵 TYPE 1→2→3→4→ＯＦＦ→TYPE 1，右鍵反向 |
| 單元測試 | `TestSoundValueDefaultsToType1`（同上）：有音檔時預設顯示「TYPE 1」——擋住 §3 那條預設值 |
| 對原版 | 版面 ✅ [`../playtest/39`](../playtest/39-system-window-parity.md)（值格逐像素相同）。**聲音本身沒有對拍**——四段的音量差沒有量過波形 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 四段的實際音量差 | 沒有錄下原版四個 TYPE 的波形量過。算式來自機器碼，聽感沒驗 |
| `AH=0Bh` 只重算三個聲部 | 原版那個迴圈是 `ah = 0、1、2`（[`../re/81`](../re/81-sound-type-attenuation.md) §5）。remake 的主增益對所有聲部一致，這一點**沒有照抄** |
