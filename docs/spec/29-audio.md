# 29 — 音樂與音效

**狀態：CONFORMED。音樂與音效都會出聲、與原版錄音比對過，
⭐ **場景對應也解出來了**（[`docs/re/58`](../re/58-bgm-scene-mapping.md)）——
四季配樂、事件曲、攻城與野戰各一首，remake 照著接。**

- 日期：2026-08-15
- 出處：[`docs/re/56`](../re/56-bgm-track-events.md)（事件與控制事件）、
  [`docs/re/57`](../re/57-opl3-register-map.md)（OPL3 暫存器、音色、音量、速度、`SOUND.DAT`）、
  [`docs/re/23`](../re/23-bgm-resource-format.md)（容器）、
  [`docs/re/17`](../re/17-dosv-audio-tsr.md)（INT 61h 介面）
- 推論等級：**confirmed**（全部由 `YNSOUND.COM` 的程式碼直接讀出，
  並以兩版 26 首曲子的音色編號範圍交叉驗證）
- 範圍：**DOS/V 版**。PC-98 側是 YM2203，這一份不涵蓋

## 1. 路線：解格式 → OPL3 合成 → ogg

使用者裁定（2026-08-15）：**不跑原版錄音，解 `BGM.DAT` 的格式。**

錄音那條路的可行性驗過了（[`docs/playtest/25`](../playtest/25-audio-capture-feasibility.md)），
但它錄的是 DOSBox 的 OPL 模擬而不是資料本身，逐曲觸發要另外解一段 RE，
而且輸出無法隨曲子的迴圈點調整。**解格式之後這些都不存在。**

| 層 | 現況 |
|---|---|
| 容器（11 曲索引、聲軌指標）| ✅ [`re/23`](../re/23-bgm-resource-format.md) |
| 音符事件 | ✅ [`re/56`](../re/56-bgm-track-events.md) §2 |
| 控制事件（音量／漸變／音色／速度／迴圈／子段／旗標）| ✅ [`re/56`](../re/56-bgm-track-events.md) §5 |
| 音色記錄（32 B）與音量換算 | ✅ [`re/57`](../re/57-opl3-register-map.md) §3–§4 |
| 時間基準（PIT 256 → 4661.65 Hz、速度分頻）| ✅ [`re/57`](../re/57-opl3-register-map.md) §5 |
| `SOUND.DAT`（19 × 16 B，含接續鏈）| ✅ [`re/57`](../re/57-opl3-register-map.md) §6 |
| OPL3 合成核心 | ✅ `internal/audio/opl3.go` |
| 事件 → 暫存器序列 | ✅ `internal/audio/bgm.go`、`effect.go` |
| 離線渲染 → ogg | ✅ `cmd/wlaudio` ＋ `tools/bgm2ogg.sh` |
| remake 播放層 | ✅ `internal/ui/sound` ＋ 系統選單第 3 列 |
| 驗證 | ✅ 與原版錄音比對（[`docs/playtest/26`](../playtest/26-bgm-render-vs-recording.md)）|

⭐ **晶片是 OPL3（YMF262）不是 OPL2。** 初始化寫了 `0x105`（NEW）與
`0x104 = 0x3F`（六對 4-operator 通道全開），這兩個暫存器 OPL2 沒有。
六個聲軌各佔一組 4-op 通道，音效走剩下的三個 2-op 通道。

## 2. 資料流

```
BGM.DAT ─┬─ 容器索引 ─→ 曲塊
         └─ 曲塊 +0x02 ─→ 音色表（32 B × N）
                +0x10 ─→ 六個聲軌指標
                              ↓
                        事件解譯（2 B 一筆）
                              ↓
                   OPL3 暫存器寫入序列（含時間戳）
                              ↓
                        OPL3 合成 → PCM
                              ↓
                             ogg
```

**中間那一層要留下來**：暫存器寫入序列（`時間, 埠, 暫存器, 值`）是
「解得對不對」與「合成得像不像」的分界。前者錯了序列就不對，
後者錯了序列對但聲音不對。**混在一起就分不出是哪一邊的問題。**

## 3. OPL3 合成核心

**從公開規格重寫成純 Go，不引入 CGO。** 跨平台建置與 Android（gomobile）
都會因為 C++ 相依而多兩套工具鏈，代價比寫核心大。

`~/cht/rich2/` 走過同一條路，那邊是 **YM3812（OPL2）**、523 行、
照規格寫而非搬運上游程式碼（行為對照 ymfm，BSD-3-Clause）。
**這裡需要的是 OPL3**，差別在：

| 項目 | OPL2 | OPL3 |
|---|---|---|
| 通道／operator | 9／18 | 18／36（兩組暫存器）|
| 4-operator 模式 | 無 | reg `0x104` 的六個位元 |
| 波形 | 4 種 | 8 種（NEW 位元開啟）|
| 輸出 | 單聲道 | 左右（`0xC0` 的 bit 4／5）|

⚠ **不要直接複製 rich2 的檔案進來。** 那是另一個 module 的私有套件，
複製會產生兩份各自漂移的實作。要嘛照規格重寫、要嘛抽成獨立套件——
**先寫再說，抽套件等第二個使用者出現。**

姓名標示與「哪些是證據、哪些是模型」的分界要寫在檔案開頭，
跟 rich2 一樣：日後對不上原版錄音時，先看是不是落在「模型」那半邊。

## 4. 離線渲染

`cmd/wlaudio`：

```
tools/bgm2ogg.sh                    # 14 首全做
tools/bgm2ogg.sh OPENBGM.DAT 0      # 只做一首
tools/go.sh run ./cmd/wlaudio -sound workplace/orig/dosv/SOUND.DAT -out workplace/audio/sfx
```

Go 這邊沒有 vorbis 編碼器，所以 WAV → ogg 走 docker ffmpeg。
**中介的 WAV 留著**：它是「合成對不對」與「編碼對不對」的分界。

**三張查表一定要從 `YNSOUND.COM` 讀進來，不准寫死在程式裡**——
那是原版資料，寫進去等於 commit 原版內容（`CLAUDE.md` §9），
而且兩版的 TSR 不見得一樣。`tools/bgmdump.py` 已經是這樣做的。

取樣率 49,716 Hz（OPL3 的原生取樣率），輸出前重取樣到 44,100。

## 5. remake 播放層

| 項目 | 作法 |
|---|---|
| 解碼 | Ebiten 的 `audio/vorbis` |
| 檔案來源 | **玩家自己從原版產生**（§6）|
| 缺檔案時 | 靜音跑，系統選單顯示「未接入」——**不要 fallback 到自製音樂** |
| 開關 | 系統選單第 3 列（`cmd/wlgame/strategyhud.go` 的 `sysRowSound`）|
| 音效 | `SOUND.DAT` 的 19 筆離線渲染成短 ogg（`sfx-NN.ogg`），接續鏈在渲染時攤平 |
| 效果碼怎麼來 | 規則層排隊（`tactical.Battle.TakeSoundEffects`），呈現層播。**碼就是 `SOUND.DAT` 的記錄編號**，不必另外對照表。已接的三個是原版證實的投射物發射／特殊發射／命中（[`docs/re/17`](../re/17-dosv-audio-tsr.md) §3）|

## 6. `[HARD]` 權利邊界

PCM、WAV 與 ogg 都是**原版衍生物**：`workplace/` 已 gitignore，
**不進版控、不進發行包**（`CLAUDE.md` §9）。
發行包裡只有渲染工具與說明，玩家拿自己的原版跑一次。

## 7. 驗證

| 方式 | 內容 |
|---|---|
| 暫存器序列 | 先驗這一層：音色編號落在表內、key-on／key-off 成對、速度換算的 tick 數與長度表一致 |
| 不是靜音 | **分段** RMS（只渲染到片頭的話後段是 0，總 RMS 看不出來）|
| 下游吃得下 | 拿 **Ebiten 自己的解碼器**跑一次，不要用別的播放器代替 |
| 對照組 | [`docs/playtest/25`](../playtest/25-audio-capture-feasibility.md) 錄到的開場曲。**只能當耳朵的參考**——那是 DOSBox 的模擬，不是硬體 |

## 8. 未解

| 項目 | 現況 |
|---|---|
| 曲 1／8／10 | 已知路徑到不了（[`re/58`](../re/58-bgm-scene-mapping.md) §5）|
| 曲 6（事件與對話）| 原版由四支對話／事件常式呼叫，remake 這一側還沒有對應的單一進入點，**沒接** |
| 換季的兩段時序 | 原版第 1 天停、第 2 天換曲，調色盤另外漸變 16 天。remake 只做了換曲那一半 |
| 迴圈點怎麼呈現 | 原版靠控制事件 `C1`／`C3` 無限循環；ogg 是有限長度，要決定渲染幾輪或另存迴圈點 |
| 全域音量偏移 | `cs:0996h` 誰設、範圍多少未解（[`re/57`](../re/57-opl3-register-map.md) §8）|
| PC-98 版 | 音源是 YM2203，暫存器路徑完全沒讀。要不要做是待裁定的問題 |
