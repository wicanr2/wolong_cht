# 20 — DOS/V IDA 逆向覆蓋與 remake 差距審計

**狀態：REVIEWED。足以支撐可玩重製與多數規則，但不足以支撐高忠實度戰術呈現；主要缺口是原版顯示串列、相機狀態機、逐幀執行順序與同狀態動態 oracle。**

- 日期：2026-08-12
- 範圍：只審查松崗 DOS/V；不以 PC-98 作完成條件
- 原始 `KI.EXE` SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原始 `KI.EXE.i64` SHA-256：`7b7c1aa67c47f99062cfdd4439b3423302808c874f929c3ea6f75ec564034c26`
- 原始 `KI.EXE.asm` SHA-256：`4364494ffb5d285681b74ef3eb3da4d6bac072373524f4cb06e861d3e6d43668`
- 工具：IDA Pro 9.4；IDA DOS/V 線性位址
- 方法：在 `/tmp` 的資料庫副本執行 `tools/ida_re_coverage_audit.idc`，原始 `.i64` 唯讀；再把 `.i64` 關係圖、既有文件、remake 實作及驗收種類互相比對。

## 1. 結論

目前的逆向工程不是「完全不夠」，而是**深度分布失衡**：資產格式、時鐘、月結、經濟、戰略資料與不少戰術規則已有可回查的機器碼證據；戰術畫面的 runtime 合成管線卻沒有完整移植。現況足以製作可玩的重製版，不足以宣稱畫面、遮蔽、相機、動畫與操作節拍接近原版。

造成多輪 polish 仍不收斂的主因，不是單一座標或素材解碼錯誤，而是 remake 以另一套 renderer 重畫已解出的資料。只調外框、viewport、glyph 與固定截圖，無法補回原版 display-list producer／consumer 的行為。

以下百分比只是本次審計的**風險估計**，不是自動量測出的完成率：

| 子系統 | 逆向足夠度 | 判定 |
|---|---:|---|
| 圖庫、調色盤、主要資料格式 | 85–95% | 多數有整除、round-trip、硬體埠或原版畫面互證 |
| 時鐘、月結、核心戰略規則 | 75–90% | 呼叫鏈、欄位讀寫端與測試較完整 |
| 戰術規則與資料模型 | 65–80% | 腳本、陣形、移動、傷害、投射物已有深度；仍有刻意的 remake 差異 |
| 戰術相機、遮蔽與畫面合成 | 25–40% | 原版關鍵函式已定位，但 renderer 沒有按該管線實作 |
| 戰術 TALK／HUD／輸入節拍 | 40–60% | 資產與主要幾何進步很大，但仍含強推論與固定 90-frame 替代時序 |
| DOS/V 音源語意與實際播放 | 20–35% | INT 61h／register 層有證據，資源、音色、歌曲狀態機仍未知 |
| 原版／remake 同狀態動態等價 | 20–30% | 多數是 fixture smoke 或不同戰況影片，尚非逐幀差分 oracle |

## 2. IDA 資料庫本身的風險

### 2.1 函式數不是穩定證據

文件長期記錄 DOS/V 為 732 個函式；本次把雜湊為 `7b7c…34c26` 的資料庫複製到 `/tmp`，以 `idat -A` 等待自動分析後枚舉出 739 個函式、50,586 個函式 bytes 與 22,267 個 instruction heads。IDA 開啟後也會改寫該副本的雜湊。

這不應解讀成「原文件一定算錯七支」，而應解讀成：**函式數會受自動分析狀態與重新開啟影響，不能拿 732 當二進位身分或逆向完成率。** 後續輸出必須同時記錄輸入雜湊、IDA 版本、是否重新分析及輸出時間。

### 2.2 `loc_1A065` 沒有可靠函式邊界

`sub_19FA0` 在 `00019FD7` 直接 `call 0001A065`，但目前資料庫對 `0001A065` 回傳 `FUNCATTR_START/END = BADADDR`。該區又含自我修改碼：`0001A06A` 的 byte 在 runtime 被寫成 `0x74`。IDA 隨後把若干 bytes 解成不合理指令，直到後方重新同步。

既有文件知道這是自我修改碼，也人工還原了主要呼叫；但資料庫沒有一個可依賴的完整函式邊界、控制流程圖或 runtime bytes。這一段掌管每輪更新／繪製，風險高於一般未命名 helper。

### 2.3 間接分派不能靠直接 xref 封口

`sub_1A426` 在 `0001A457` 執行 `call cs:funcs_1A457[bx]`。IDA 能列出 19 個候選 handler，但一般「誰直接呼叫誰」查詢無法代表 runtime 選到哪支。相同問題也存在 `sub_1C01D` 的滑鼠 handler table。

此外，立即值形式的位址、以 segment base 加 offset 的欄位、取址後的間接寫入都可能沒有 data xref。現有 `ida_func.idc` 已警告這點，但不少完成敘述仍偏重 direct xref 的有無。

### 2.4 文件 provenance 有舊錯

`docs/re/02-palette-routine.md`、`03-image-blitter.md`、`04-mmap-entry-points.md`、`05-battle-selection.md` 與 `CLAUDE.md` 的早期表格，把 `KI.EXE` 的 `fffeba…3868` 寫成 `.i64` 雜湊；目前 `.i64` 實際是 `7b7c…34c26`。`CONTEXT.md` 的較新段落已分開記錄，但舊文件尚未全面勘誤。

這不推翻那些格式結論，卻表示「文件寫了 IDA 證據」仍需回查輸入身分，不能批次視為相同資料庫已證實。

## 3. 戰術畫面差距的直接根因

### 3.1 原版相機是狀態機，remake 是每幀跟隨大將

原版 `word_1E160`／`word_1E162` 是所有投影／顯示串列 routine 共用的 viewport origin：

- `sub_19946` 先由 `word_1D328`／`word_1D32A` 呼叫 `sub_1DC9D`。
- `sub_199F3` 把初值設為 `0x24`／`0x0E`。
- 滑鼠 handler 的 `0001C0D…0001C10…` 區段由游標座標計算新值、寫回 `word_1D328`／`word_1D32A`，並設 `byte_1D348=1`。
- `loc_1A065` 只在 dirty flag 設定時重新呼叫 `sub_1DC9D`。

remake 的 `drawBattleIso` 則每次 Draw 都呼叫 `centreOn(me.X, me.Y, me.Z)`，鏡頭持續鎖在玩家大將。這是明確的產品差異，不是原版公式的移植。即使 viewport 寬高調對，取景仍會不同。

### 3.2 原版有共用顯示／遮蔽串列，remake 沒有

原版的證據鏈是：

1. `sub_1D958` 配置／初始化 `word_1E15C`。
2. `sub_1D971`、`sub_1D9D1`、`sub_1DA1C`、`sub_1DAAA`、`sub_1DC03` 等 producer 把地形、人物、旗與效果登錄到同一份結構。
3. `sub_1DDB4` 逐列讀 `word_1E15C`，比較相鄰 cell 的 flags／depth，再選擇 `sub_1DFBB` 或 `sub_1DE95`＋`sub_1DFA6` 合成到 VGA planar VRAM。

remake 的 `battleview.go` 是固定順序：

```
地形 → 場景旗 → 投射物 → side 0 全兵 → side 1 全兵
```

兵也只按陣列順序繪製，沒有依原版 display-list 的 cell、depth、鄰格遮蔽或 sprite 高度合成。這會造成城牆前後、兵與旗、投射物和高物件的遮蔽錯誤，也是戰術畫面「骨架像、實際觀感不像」的首要原因。

### 3.3 模擬更新與畫面更新不是原版的一輪一畫

原版主迴圈是腳本、runtime 更新、滑鼠／相機、顯示串列合成緊密相連。remake 在 `updateBattle` 依 `speed` 一次執行多個 `Battle.Step()`，Ebiten `Draw` 只畫最後狀態。速度大於 1 時，中間動畫幀天然消失。

`Battle.Step()` 的順序也不是已證實的逐指令等價：目前是腳本 → 全兵 → 投射物 → 增援 → 優劣 → 攻城扣血 → 大將退卻 → 勝負；原版 `loc_1A065` 的可辨識序列包含 `sub_1A12A`、`sub_1A6FA`、`sub_1B941`、`sub_1ADC8`、`sub_1C6F6` 與 `sub_1DDB4`，而自我修改區尚未有 runtime trace。即使每個 helper 的局部公式正確，跨 helper 的一幀先後仍可能差一輪。

### 3.4 動畫相位仍有 remake 自行決定的部分

- 兵圖使用 `b.Frame+k` 產生各兵相位；原版是每筆兵記錄自己的 bit 0 在該兵更新後翻轉。
- 戰術 TALK 固定 `battleTalkDuration=90`，但原版 payload、消像及輸入／自動推進時序沒有完整 trace。
- `replanInterval=30` 明載為原版沒有的效能節流，可能改變阻塞陣線、交戰接觸與戰鬥長度。
- 城壁破壞後 minimap 是否局部更新未知；目前 remake 固定快取初始 base image。

上述差異會一起污染推廣片的視覺對拍，不能再當成單純 polish。

## 4. 為什麼現有測試沒有抓到

現有測試對「可玩、不崩、局部公式與資產索引」很有價值，但多數戰術驗收是：

- `-open-siege`／`-open-battle` fixture；
- remake 自己的幾何 contract；
- 不同戰況的 YouTube frame，只作 layout-only 比較；
- 最終截圖，沒有逐幀 state／display-list trace；
- 正常路徑 smoke，沒有同一戰場、同一雙方、同一命令、同一時刻的 DOS/V 對照。

因此測試綠只能證明 remake 內部一致，不能證明原版等價。`docs/playtest/09`、`19`、`20`、`22` 已誠實標出此邊界；問題是後續實作仍在這個證據缺口上反覆調常數。

## 5. 最短收斂方案

不應再先調 UI 常數。建議把後續工作改成以下四個 gate，依序完成：

1. **重建原版 renderer IR。** 對 `word_1E15C` 建立原始 cell entry 結構；完整解讀 `sub_1D958`、`1D971`、`1D9D1`、`1DA1C`、`1DAAA`、`1DB34`、`1DB9B`、`1DC03`、`1DC9D`、`1DD22`、`1DDB4`、`1DE95`、`1DFA6`、`1DFBB`。先輸出 deterministic display list，不碰 Ebiten。
2. **移植相機狀態機。** 保存 `word_1D328/32A/32C/32E` 對應狀態、dirty flag 與滑鼠邊界；取消每幀 `centreOn`。用原版初始 `0x24/0x0E` 及一段捲動畫面作 fixture。
3. **建立同狀態 DOS/V oracle。** 使用原版可通過密碼頁的事實，製作可重放存檔／輸入序列；每一幀記錄戰場編號、雙方編成、命令、相機原點與截圖。先鎖 60–120 幀，不需要完整遊戲測試。
4. **逐層差分。** 依序驗 terrain-only、display-list entries、composited field、HUD/TALK、整張 320×200；每層都有 machine-readable diff，禁止用不同戰況影片替代 same-state gate。

IDA 部分應同時補兩項方法：對 `loc_1A065` 使用 debugger／DOSBox trace 取得 runtime bytes 與呼叫順序；對間接 table 與 segment-relative 欄位建立受版控的位址索引，輸出時自動附 `confirmed／強推論／假說／未知`。IDA 的自動函式與 xref 只當導航，不再當完成證據。

## 6. 最終判定

**目前 IDA RE 的程度：對「玩法骨架與多數規則」大致足夠；對「DOS/V 原版外觀與逐幀手感」明顯不足。**

重製差距與逆向不足確實有直接關聯，但更精確地說，是已定位的呈現層 routine 沒有被端到端解讀／移植，反而被另一套相機、繪圖排序與動畫節拍替代。下一輪若先完成顯示串列、相機與 120 幀同狀態 oracle，戰術 parity 才會開始收斂；若繼續只看推廣片修座標，預期仍會反覆回歸。

## 7. 2026-08-12 第一個修正切片

本輪沒有把第 5 節全部誤標成完成，只先接入可由現有 IDA 證據封閉的部分：

- `battleView` 保存原版 `word_1D328/word_1D32A` 的 world origin；初值使用
  `sub_199F3` 的 `0x24/0x0E`，再由 `sub_1DC9D` 尾段公式換算成
  `word_1E160/word_1E162` 投影 origin。已移除每幀 `centreOn` 玩家大將。
- DOS/V `0001C0C6` 起的縮圖點選算式已接入 `setCameraFromMiniMap`；點縮圖會改寫
  保存的 world origin，不再只是 UI 裝飾。
- 兵圖動畫改讀每筆 `Soldier.PoseStep`，並接入 `HitGeneral` 的 `+0x02 bit 3`；不再
  使用 `b.Frame+k` 人工錯開所有兵。
- 地形、場景旗、投射物、兵與大將旗先建立同一份 `battleDisplayEntry` 中介表示，再
  依 projected row／col／Z 排序繪製；已消除固定 side 0 → side 1 的覆蓋順序。

證據等級仍分開：共用 producer／consumer 與相機公式為 **confirmed**；目前
`battleDisplayEntry` 的同格 tie-break 是**強推論**，尚未逐 byte 移植 `sub_1DE95` 的
鄰格遮蔽，因此不能宣稱原版顯示串列完整 parity。受控 `-open-siege` 截圖顯示固定原點
已把取景從城外空地移回城門／橋樑中心；這是 fixture 視覺驗收，不是同狀態 oracle。

自動驗收新增：原版初始相機換算、縮圖兩端座標、display-list 深度不受 side 順序支配、
每兵 `PoseStep` 圖號。`go test ./...` 已在 Docker／Xvfb 通過。

## 8. 2026-08-13 第二個修正切片：端到端顯示格 consumer

前一節的強推論全域排序已移除。`tools/ida_tactical_display_grid.idc` 對唯讀
`KI.EXE` 副本重建 IDA Pro 9.4 資料庫，匯出 producer／consumer 函式邊界、
xref 與完整原始指令；`sub_1DDB4` 的檔案位移 `0xDFB4` raw bytes 另行比對，
前 225 bytes 與 IDA 指令完全吻合。輸入仍是 SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`，
位址空間是 IDA DOS/V linear address。

新 renderer 已完成 32×30 cells、每格 32 B、兩 lane、七 depth；接入
`sub_1DDB4` 的 23×15 奇數欄 anchor、`sub_1DE95` 的鄰格順序、人物雙 raw-unit
producer、場景旗單 raw-unit producer，以及 `sub_1E011` 的
16×32 encoded → 32×16 VGA 四象限重排。戰術 viewport 改為原版 480×368
像素，不再先畫 240×184 再整張放大。

受控攻城 fixture 顯示城牆磚列、橋、河岸、兵列與旗的尺度已回到 DOS/V
並排截圖的量級；這是 renderer 結構驗收，仍不是同狀態自然流程逐像素 parity。
完整 `go test ./...` 已在有界 Docker／Xvfb 通過。

仍保留的邊界：原版 flags/depth 的 dirty-cell 快速路徑、EGA planar AND/OR 的
逐位元等價、逐幀相同狀態 oracle 與 renderer 效能調校。它們不再由全域排序或
錯誤 tile 尺度遮蔽，但仍應分別驗證。

## 9. 未解

§5 的四個 gate 裡，前兩個已解（§7／§8 兩個切片）。剩下的是這些：

| 項目 | 現況 |
|---|---|
| 同狀態動態 oracle | 沒有可重放的存檔／輸入序列，所以「原版等價」目前無法驗。**這是還沒做，不是做不了**——DOS/V 的密碼頁空白確認就會過（[`../playtest/18`](../playtest/18-dosv-password-verification.md)），PC-98 側連除錯器都接好了（[`../playtest/21`](../playtest/21-dosboxx-bridge-sampling.md)）|
| 逐幀執行順序 | 顯示串列與相機已重建，但整幀的呼叫順序沒有逐幀對過 |
| `loc_1A065` 的 runtime bytes | 自我修改碼，靜態影像看不到每輪的實際內容（§2.2）|
| 四層差分（terrain／display list／composited／HUD）| 沒有 machine-readable diff，目前只有 layout-only 比較 |
