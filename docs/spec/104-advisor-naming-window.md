# 104 — 「自定」：軍師命名視窗

**狀態：CONFORMED（2026-08-29 實作、單測與截圖）。**

- 日期：2026-08-29
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_18FC9`（控制器）、
  `sub_190C0`（選字）、`sub_1928A`（畫表）、`sub_19223`（畫六格）、
  `sub_1905E`／`sub_1906E`（重來／繼續）、`sub_1912D`／`sub_19144`／`sub_192F5`（肖像）、
  `sub_18F7C`（取消）、`sub_18E5A`（君主卡的確定條件）；版面 [`../re/54`](../re/54-advisor-naming-window.md)；
  選字表 [`../formats/10`](../formats/10-end-s15-namechars.md)；存檔欄位同文件 §4。
- 推論等級：**全部 confirmed**（2026-08-29 補完：`funcs_19037` 逐格對過、
  肖像座標由 `sub_107D2` 的 `bx` 換算並用君主卡做正對照）。

## 1. 原版做什麼

君主卡的「自定」（熱區 `0x21`）開這個視窗（`sub_11AC3` → `sub_18FC9`）：

1. 讀 `END_S15.DAT`（2,621 個注音序的 Big5 字）進暫存記憶體。
2. `ds:5221h`（肖像編號）是 `FF` 就設成 `0x91`。
3. 畫視窗（場景 9）、六格名字（前三格軍師名、後三格別號）、肖像、本頁 96 個字。
4. 等待迴圈：每個熱區一個 handler——

⭐ **`funcs_19037` 逐格對過**（`KI.EXE.asm`）：八筆的順序與位址順序相同，
第九筆是 `start`（越界保護，選不到——`al == 0x20` 在分派之前就被分走了）。

```
funcs_19037  dw sub_1905E  ; 0x21 重來    dw sub_1906E  ; 0x22 繼續
             dw sub_1907F  ; 0x23 上一頁  dw sub_1908D  ; 0x24 下一頁
             dw sub_190A3  ; 0x25 聲母列  dw sub_190C0  ; 0x26 選字
             dw sub_1912D  ; 0x27 前 ▲   dw sub_19144  ; 0x28 後 ▼
```

| 熱區 | 動作 | 出處 |
|---|---|---|
| `0x20` 確定 | 收尾：勢力 `+0x02 = 0x7F`、擦視窗、釋放記憶體 | `sub_18FC9` 的 `mov byte ptr [bx+2], 7Fh` |
| `0x21` 重來 | 目前格寫全形空白、**退一格**、重畫 | `sub_1905E` |
| `0x22` 繼續 | 目前格寫全形空白、**跳一格**、重畫 | `sub_1906E`（與 `1905E` 同形，方向相反）|
| `0x23`／`0x24` 上下頁 | 頁位移 ∓ `0xC0`（96 字），上限 `0x13BA` | `sub_1908D` |
| `0x25` 聲母列 | 跳到那個聲母的頁——**標記字本身在表裡**（ㄅ0、ㄈ242、ㄋ569、ㄍ784、ㄐ1089、ㄑ1317、ㄓ1570、ㄕ1812、ㄙ2101、ㄨ2321）| `formats/10` §3 |
| `0x26` 選字 | (x−210)÷20、(y−238)÷20 定格，餘數 ≥16 不算；寫進目前格、游標前進（第六格不動）| `sub_190C0` |
| `0x27`／`0x28` 前▲後▼ | 肖像 −1／+1，在 0..0x92 循環，**跳過武將表裡有人用的號碼** | `sub_1912D`／`sub_19144`／`sub_192F5`（比 `[bx+1]` ＝ 武將的肖像欄）|

5. 右鍵／ESC：`sub_18F7C` 把 `+0x02` 從 `+0x3F` 抄回、`5221h = FF`、六格清成空標記 `D0A1`。

君主卡那一層的「確定」放行條件（`sub_18E5A`）：`+0x02 ≠ 0x7F` **或** 第一格不是 `D0A1`
——沒取名的自訂軍師不能開局。

## 2. 存檔

肖像在區塊 `+0x52A1`、六個字在 `+0x52A2`（`ds:0D52h` 段整段存，`formats/10` §4）。
先前 `spec/27` §5 掛的 ⛔「存不回去」**不成立**——欄位一直都在，只是不在勢力記錄裡。

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 模型（不認識 Ebiten）| `cmd/wlgame/naming.go`：`namingModel`、`namingHotspotAt`、`namingGridIndexAt` |
| 畫面 | 同檔 `drawNaming`；版面常數照 `re/54` §1 與 `sub_1928A`／`sub_19223` 的立即值 |
| 接點 | `lordCardCustom` → `openNaming`；`updateLauncher` 開著時輸入全歸它；「確定」留 `customAdvisor`，開局時 `World.SetCustomAdvisor` |
| 規則層 | `internal/state/advisor.go`：`SetCustomAdvisor`／`HasCustomAdvisor`／`AdvisorNameRaw`；`TalkNoticeVars` 的 `{4}` 後援 |
| 選字表 | `internal/assets/namechars` |
| 驗收旗標 | `-open-naming` |
| 差異 | 聲母列的跳點用表裡的標記字定位（原版是 `sub_190A3`，動作相同）|

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 單元測試 | `TestNaming*`（`cmd/wlgame`，五條：選字與空隙、重來／繼續、聲母跳頁與翻頁夾制、肖像跳過、空名不放行）；`TestAdvisorNameRoundTripsThroughBlock`、`TestCustomAdvisorNameFeedsTalkVars`（`internal/state`）|
| 截圖 | `tools/parity_shot.sh out.png -open-naming`：視窗、六格、聲母列、16×6 選字格、翻頁列 |
| 對原版 | **未做**——原版這一頁要從新遊戲流程點「自定」進去，`tap:x,y,5` 能點選單之後已經可拍（`playtest/54`），排進下一輪 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~肖像的位置~~ | **(208, 144) confirmed**：`sub_1915B` 傳給 `sub_107D2` 的 `bx = 2D1Ah` 是 VRAM 位元組位移，`y*80 + x/8` ⇒ 144×80＋26 ⇒ (208,144)。正對照是君主卡的 `sub_18EA0`：`2817h` ⇒ (184,128)、`34A7h` ⇒ (312,168)，與 [`27`](27-lord-select-window.md) §1 記的座標一致 |
| ~~「重來」「繼續」的方向~~ | **confirmed**：`funcs_19037` 逐格對過（§1）|
| 六格與本頁字的顏色 | 六格 15、游標底線 15／1、字 9：`sub_19223`／`sub_1928A` 的屬性值直讀，但沒逐像素對過 |
