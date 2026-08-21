# 14 — `MMAP.MCH` 戰略地圖物件

**狀態：資產格式、事件 12 的火災／暴動圖形鏈與 typed 動畫／移動時序 confirmed；
type 3 的事件語意仍未知。**

- 日期：2026-08-10
- IDA 輸入：DOS/V `workplace/ida/dosv/KI.EXE`，SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA 資料庫：`KI.EXE.i64`，SHA-256
  `7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- IDA Pro：9.4，16-bit DOS/V 線性位址，segment base `0x10000`
- `MMAP.MCH`（DOS/V／PC-98 相同）SHA-256
  `b10a5b64bbffa672c1fb5cb37703ac4c14b18bf1166cc47c4e802c19aae9f8f7`

2026-08-11 補驗：`sub_12459` 的後半區（raw slots 16–31）moving slot 已接 `sub_1248A` typed slice；
raw drift word、`sub_124FF` signed-byte normalization、direction byte、storm bounds
與 wrap 由 `internal/state` 測試固定。這個補驗不替 type 3 猜事件名稱，也不把 remake
測試誤寫成原版自然畫面逐像素對拍。

## 1. 載入位置

`sub_187CC`（`000187CC`）依段位址計算：

```asm
word_19876 = base + 3000h
word_1987A = base + 4200h
```

`sub_187AF`（`000187AF`）先把 `MMAP.MDL` 載到 `word_19876`，再把
`MMAP.MCH` 載到 `word_19876 + 800h`。因此 `word_1987A` 指向
`MMAP.MCH` 的檔案 offset `0xA000`：

```text
(0x4200 - 0x3800) paragraphs × 16 = 0xA000
```

## 2. MCH 圖塊格式

`sub_1D804`（`0001D804`）以 `AH × 160` 取一張圖，先讀 16 個 mask word，
再讀 64 個 color word。因此 `MMAP.MCH` 前 `0xA000` 是 256 張 16×16
圖塊，每張：

```text
0x00–0x1F：16 列 mask（每列 2 byte）
0x20–0x9F：4 個色平面（每個 0x20 byte）
```

mask bit 為 0 的像素不寫入畫面；其餘像素由四個平面合成 0–15 的調色盤
索引。這與 `sub_1D66A` 對 `word_1D84C` 的呼叫相符，不是把 MCH 當成
一般 128-byte `MMAP.MDL` 圖塊。

## 3. 物件 metadata 與 source 矩陣

`word_1987A` 對應檔案 offset `0xA000`。`sub_12533`（`00012533`）讀：

| metadata 欄位 | 位址 | 語意 | 證據等級 |
|---|---:|---|---|
| 寬 | `word_1987A + index×4 + 0` | source tile 欄數 | 已證實 |
| 高 | `word_1987A + index×4 + 1` | source tile 列數 | 已證實 |
| 位移 | `word_1987A + index×4 + 2` | 加到 `word_1987A + 0x100` 的矩陣位移 | 已證實 |

所以：

```text
0xA000–0xA0FF：64 筆 4-byte metadata
0xA100 起    ：source tile ID 矩陣；0xFF 表示透明格
```

`loc_1D51F`（`0001D51F`）把矩陣的每一個 ID 寫到 strategic map 的物件
layer；最終 renderer 再依 MCH 圖塊格式畫出它。這也是為什麼 source 矩陣
裡會出現 `0x80`–`0xE2`：它們是 MCH 圖塊 ID，不是 palette 色號。

## 4. 事件 12 的 object type 查表

`sub_134B1`（`000134B1`）把事件字高 byte 傳給 `sub_123FF`（`000123FF`）
的 `[si+0E]`；高 byte `1`／`2` 分別建立兩種物件。`sub_12533` 用：

```asm
BL = object_type × 8 + [si+0F]
BL = CS:[BX - 67A6h]
```

16-bit wrap 後的查表位址是 `CS:985Ah`，IDA 線性位址為 `0001985A`。IDA
資料庫的原始 bytes 為：

```text
type 1：18 19 1A 1B 1C 18 19 1A
type 2：20 21 22 23 20 21 22 23
type 3：28 29 2A 2B 28 29 2A 2B
```

對應 MCH metadata 後：

| 原版 object type | MCH metadata index | 尺寸 | 目前語意 |
|---:|---|---:|---|
| 1 | `0x18`–`0x1C`，循環重用 | 16×9 tile | 事件 12 火災，已證實 |
| 2 | `0x20`–`0x23`，循環重用 | 5×5 tile | 事件 12 暴動，已證實 |
| 3 | `0x28`–`0x2B`，循環重用 | 5×5 tile | 查表已證實，事件語意未知 |

`internal/assets/world/mmapmch.go` 以這份位址／bytes 對照解碼，
`cmd/wlgame` 已將 type 1／2 接到 `drawDisasterOverlay`。物件以據點的
地圖格座標置中，與 `loc_1D51F` 的 `width/2`、`height/2` 算法一致。

## 5. runtime timer 接線與尚未升格成 parity 的部分

`sub_123FF` 初始寫入 `[si+0F]=1`、`[si+0C]=1`、`[si+0D]=0x10`；
`sub_12459` 以 timer 設定 dirty，`sub_12533` 才在該次繪製先取舊
`[si+0F]`、再加一並 `and 7`。`internal/state/disaster_objects.go` 已以
非序列化 typed runtime record 接入這三段：每一筆物件各自保存 `Phase`／`Timer`／
`Interval`／`Dirty`，`cmd/wlgame` 只在可見 map-loop Update 呼叫 advance，modal 不穿透。
`TestDisasterObjectAnimationTiming` 固定建立初值、16-update cadence、dirty render
舊 phase、八相位遞增與據點清除；因此已移除先前的固定 presentation frame clock
替代方案。`World.DisasterMarker` 仍只保存持久 marker／level，不把 runtime 欄位寫入存檔。

`sub_12459` 對 raw slots 16–31 逐筆走 `sub_1248A`；它不是只處理最後一筆物件。
此移動分支已由 typed runtime record 接線並以 fixed-point／signed-byte／wrap 測試固定；
火災／暴動圖像的 type 1／2 時序也不依賴它。

事件 11 的 `sub_134A6` → `sub_1237E` 只確認據點 `+0x15` marker 與
暴風雨範圍，沒有呼叫 `sub_123FF`；因此暴風雨仍不套用 type 1／2／3 圖形。

## 6. 2026-08-12 再審：type 3 與 `sub_1248A` 的直接證據邊界

本節使用同一份松崗 DOS/V `KI.EXE`（SHA-256 見文件開頭）、IDA Pro 9.4 與
`tools/ida_mch_deep.idc` 的唯讀匯出；位址一律是 IDA 線性位址。

- **已證實：** `sub_123FF`（`000123FF`）的唯一直接 caller 是事件 12 handler
  `sub_134B1`（`000134B1`）；`sub_134B1` 取 queue Code.high 為 object type。
- **已證實：** 月結的 `sub_12286` 兩次直接寫入 Code 為 `0x010C`、`0x020C`，故目前
  可見的原版 direct producer 只建立 type 1／2。
- **強推論：** type 3 不由目前所有可見的直接事件 12 producer 產生。這不是宣稱 type 3
  不可達：尚未被 IDA 關係圖捕獲的間接／外部 producer 仍是未知。
- **已證實：** `sub_12459` 對每個 raw `SI >= 0x2140` 的 slot 呼叫 `sub_1248A`，正好是
  slots 16–31；先前「最後一筆」的說法已撤回。

因此 remake 維持 type 3 的資產查表、拒絕替它命名或合成自然事件；已完成的
`sub_1248A` 接線只覆蓋這支函式明確可證實的移動資料流。
