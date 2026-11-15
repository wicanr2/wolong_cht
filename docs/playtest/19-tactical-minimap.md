# 19 — DOS/V 戰術縮圖 raw producer 驗收

2026-11-15兵點續片勘誤：前輪保留的17拍縮圖殘差已消除。
原版逐兵擦舊位置、登錄目前位置，破牆後重畫十字與陣形線會蓋住已畫的兵點，只有稍後輪到的兵再畫回。
isoview沿既有繪圖observer保存64×64角色標記，桌面不再由存活清單每幀重補全部點。
守方319拍與攻方八張後段樣本皆0像素差異；原指令數路徑、固定seed、319拍正常呼叫及界線見[spec159](../spec/159-tactical-command-movement-order.md)
與[正式收據](../../workplace/scratch/ui-windows-20261104/cao-minidot319-integration.json)。下方先前17拍的敘述保留歷史定位，不再是現行縮圖缺陷。

## 2026-11-15縮圖受損更新勘誤

已由原版正常守汝南319個`1A065`邊界及目前同狀態重播證實：破牆／破門的當拍就更新縮圖地形。
既有DOS/V KI.EXE、IDA linear `1B824→1BB6D`的圖塊路徑沿用，輸入雜湊與執行器身分見[spec159](../spec/159-tactical-command-movement-order.md)。
`isoview.Minimap`已依目前Field與Revision更新同一張底圖，表頭／尾列保留，兵點／陣形線／十字仍由介面疊上，不重擲旗。
第266拍起1260個跨拍地形改動像素全等；六張完整畫面除95游標像素外全等，另17拍兵點擦除／順序差異維持有界。
[正式證據](../../workplace/scratch/ui-windows-20261104/cao-chrome320-evidence.json)與[整合收據](../../workplace/scratch/ui-windows-20261104/cao-chrome320-integration.json)保存新證據。
本段取代下方「縮圖建一次不更新」的當時實作限制；舊畫面與歷史數字不改。

**狀態：PASS（已證實 producer 的 remake 實作）。底圖、部隊點與陣形線都已接上
（陣形線見 `docs/spec/37`）；游標十字與小地圖上的城壁受損仍未接（位址已定位）。**

- 日期：2026-08-12（部隊點：2026-08-16）

本切片只處理松崗 DOS/V 戰術縮圖，不涵蓋 launcher、自然 HUD、戰術其他面板、按鈕
glyph、規則或事件。

## 已實作的證據鏈

IDA DOS/V `KI.EXE` 的已證實路徑為：

```text
sub_1C83E
  mapY = 0..63
  mapX = 0..63
  → sub_1C4FA
      BATTLE.MAP[mapY*64 + mapX]
      BATTLE.MDL attribute[tile]
  → sub_1C51E
      每格 2×2
      screenX = 0x1F0 + 2*mapY
      screenY = 0x50 + 2*(63-mapX)
```

remake 對應如下：

- `internal/assets/battle.Library.TileAttributes` 公開每個圖塊組的 raw 256-byte
  `BATTLE.MDL` attribute table，回傳複本。
- `internal/assets/battle.RenderTacticalMinimap` 只輸出 128×128 palette index；
  attribute 取 EGA set/reset 實際使用的低 4 bit。
- `battleView.minimap` 在 `newBattleView` 初始化時建立一次，之後畫面只重用快取。
- `SideMiniMap` 對齊 DOS/V 原點 `(496,80)`，不再將 128×128 base image 壓進
  108×96 的高度圖 fallback 內。
- 部隊點：`sub_1B240` 在 `0001B284` 依單位記錄的位址分色——`si < 0x600`
  用調色盤索引 10、否則用 3，也就是**側 0 一色、側 1 一色**
  （[`../re/60`](../re/60-tactical-sidebar.md) §7）。`drawBattleMiniMapUnits`
  照同一條座標換算把每個活著的兵畫成 2×2 點，顏色取自
  `GAMEPAL.BRG` 當季 bank 的索引 10／3。

## 還沒接的局部更新

位址都已定位（[`../re/60`](../re/60-tactical-sidebar.md) §7），實作還沒做：

| 標記 | producer | 調色盤索引 |
|---|---|---:|
| ~~陣形線（側 0 的陣形原點那一整行）~~ ✅ 已接（`cmd/wlgame/battle.go`，`docs/spec/37`）| `sub_1C5AE` | 11 |
| 游標十字（一行 ＋ 一列）| `sub_1C577` | 0 |
| 城壁受損（圖塊值 `+0x10`／`+8` 後重畫）| `sub_1B824` | 隨新圖塊 |

⚠ **「城壁受損」那一條只差在小地圖這一份。** 戰場本體的重畫已經實作了
（[`../spec/66`](../spec/66-broken-walls-repaint.md)）——缺的是
`isoview` 的縮圖在初始化時建一次之後就不再更新。

## 自動化驗收

Docker 內執行：

```text
go test ./internal/assets/battle ./cmd/wlgame -count=1
```

結果：兩個 package 均通過。測試涵蓋：

- raw tile → raw attribute → palette index；
- 2×2 填色與 128×128 尺寸；
- 原版轉置／Y 反轉座標；
- palette index → RGBA；
- `TileAttributes` 的圖塊組選取與複本隔離。

## 實機畫面

Docker／Xvfb 以 `-scenario 0 -player 0 -seed 17 -speed 1 -open-siege` 啟動，使用
wlgame 內建 `-shot` 擷取第 30 幀：

- [wlgame-tactical-minimap-siege.png](../images/wlgame-tactical-minimap-siege.png)
- 解析度：640×400
- SHA-256：`9ad40f4d6042355a1283367ee4c0d930cb001e9968459ec34f1853c84bba7438`
- 原版視覺參照的歷史位置：`workplace/promo-live/review/original-320s.png`、
  `workplace/promo-live/review/original-400s.png`。

2026-11-02 補註：這兩個歷史檔案目前缺檔。保留位置供追溯，當前驗收不引用這兩張圖。

這是 `-open-siege` 受控 tactical fixture，不能升格為完整自然流程或同狀態逐像素
parity。原版參照影片畫面與 remake fixture 的戰場狀態不同；本輪只驗證縮圖的 raw
producer、方向、尺寸與可見接線。

## 未知邊界

- `sub_1C51E` 的 EGA set/reset byte 已證實可視為低 4-bit palette index；其對應的
  戰術專用硬體調色盤語意未在本切片新增推論。
- 原版在戰鬥中破壞城壁、移動單位後是否有縮圖局部更新 routine，現有證據不足，
  保留 unknown。
- `sub_1C83E` 後續 `sub_1E3D7` 繪製的完整縮圖外框／裝飾素材不在本切片解碼；目前
  remake 只保留 128×128 base image 與最小外框。

本輪沒有修改 `WORKLIST.md`、`CONTEXT.md`，沒有建立 `HANDOFF.md`，Docker 容器已由
`--rm` 清理。
