# 17 — DOSBox 原版／remake 可玩性專家驗證

**狀態：remake 正常策略路徑與存檔／讀檔通過；DOS/V 原版密碼頁已可進入開場，尚未展開完整自然長程驗證；PC-98 保留為歷史規則 oracle。**

- 日期：2026-08-11

本輪只把原版執行檔當行為 oracle，原始資料以唯讀掛載；所有 DOSBox、Xvfb、截圖與
remake 執行均在一次性 Docker 容器內完成。remake 使用目前工作樹建出的 Linux 執行檔、
`-seed 17`，正常策略路徑沒有使用 `-open-*` 旗標。

## 1. 結論

| 範圍 | 結果 | 證據界線 |
|---|---|---|
| 松崗 DOS/V 原版 | **PASS（啟動至開場）** | 2026-08-12 證實空白確認／`0000`／`1234` 均越過密碼頁；完整自然長程驗證尚未執行 |
| PC-98 原版 | **開機與自然 oracle 可用** | DOSBox-X 640×400 可進 `NEW GAME`；既有受控流程已到劇本選單與戰略地圖，見 `pc98-oracle-scenarios.png`／`pc98-oracle-in-game.png` |
| remake 正常策略 | **PASS** | `A → 編成 → 關閉通知 → M → 選軍團／目的地 → 提高速度`，在 196/6/28 進入「呂布 對 曹操／攻城／戰鬥指揮／委任」 |
| remake 存檔／讀檔 | **PASS** | 實際寫入 88,832 bytes 的 overlay，第 1 槽儲存後再讀回；原始 `SINARIO.DAT` 唯讀 |
| remake 戰術 GUI | **debug smoke PASS；正常路徑沿用既有證據** | 目前建置可開戰術畫面並接受 `2` 號攻擊輸入；本輪沒有把 debug fixture 升格成原版自然 parity，正常無旗標戰術證據仍見 playtest 09 |

## 2. DOSBox 原版觀察

### DOS/V 密碼保護

以 `dosbox-run:latest`、固定 `cycles=20000`、`machine=vgaonly` 啟動
`workplace/orig/dosv/START.BAT`，15 秒後截到複製保護頁。頁碼每輪可變，本輪為 15；
既有紀錄曾觀察到其他頁碼，因此頁碼不是 parity 條件。

- 截圖：[`original-dosv-password.png`](../images/original-dosv-password.png)
- SHA-256：`a9a972edbb4c896a914a84acfe65b4e55a8f93a6ea532b6feef7036d527ed5bf`
- 截圖解析度：640×480；DOSBox 視窗包含外圍區域，遊戲邏輯畫面仍是 640×400。

### 2026-08-12 勘誤：原版確認流程可進入開場

既有 `dosbox-run` 歷史 smoke 只到密碼頁；改以既有 DOSBox-X 的 INT 33 integration
輸入後，空白確認、`0000`、`1234` 均在新副本中進入原版開場。這不是二進位修改或
密碼猜測，且不會重製密碼頁；完整證據見
[`18-dosv-password-verification.md`](18-dosv-password-verification.md)。

因此 DOS/V 現在可用於密碼頁後的自然流程採樣；本輪尚未把它擴大為完整長程玩法或
同狀態逐像素 parity 驗收。

### PC-98 規則 oracle

DOSBox-X 使用 `machine=pc98`、`cputype=486`、`cycles=20000`，原始 PC-98 檔案先複製到
容器 `/tmp/game`，沒有修改 `workplace/orig/pc98`。本輪重現到 `NEW GAME` 入口：

- [`expert-original-00-newgame.png`](../images/expert-original-00-newgame.png)，640×400，
  SHA-256：`2ead3a9e27f8b2dafece4191248944f9844470f691119674bea857114746ec7f`
- 既有受控流程已通過劇本選單與遊戲內戰略地圖：
  [`pc98-oracle-scenarios.png`](../images/pc98-oracle-scenarios.png)、
  [`pc98-oracle-in-game.png`](../images/pc98-oracle-in-game.png)。後者 SHA-256：
  `318088b52e1a52bc205d7197826554dce69d42389b229c81c9faceb597bd6459`。

目前 headless DOSBox-X image 沒有 window manager，重播新局時的 X11 bus-mouse／焦點
輸入沒有穩定前進；這是 oracle 操作橋接限制，不是以此宣稱 PC-98 不可玩。既有
`docs/playtest/06-mouse-solved.md` 已記錄可用的閉迴路滑鼠與正常新局入口。

## 3. remake 正常玩家路徑

### 策略／時鐘／遭遇

目前建置以 `-seed 17 -speed 0` 開始，沒有 `-open-form`、`-open-corps`、
`-open-battle`、`-open-siege` 或其他 debug 旗標。可重播輸入為：

```text
A；Enter；Enter；3；Down + Space × 5；Enter
Enter（關閉「曹操 編成完畢」通知）
M；Enter；Enter；Down × 22；Enter；Enter
= × 64
```

結果：編成成功、目的地列表可選取、已下達的軍團在玩家不再下命令時持續行軍，
日期從 196/4/1 流逝，固定種子的一次穩定重播在 196/6/28 顯示「呂布 對 曹操／攻城／
戰鬥指揮／委任」。

代表幀：

- [`expert-remake-debug-01-formed.png`](../images/expert-remake-debug-01-formed.png)：
  編成完成；SHA-256 `8b2874a41757801dd62036be8c1536b4237ccdc67ff8052571864635967826f4`
- [`expert-remake-debug-06-encounter.png`](../images/expert-remake-debug-06-encounter.png)：
  196/6/28 遭遇選單；SHA-256 `c60e1e1a801e5fec741af7c3a5fb6edecf0d110f64809a51df96407e761ae130`

另一次重播在提高速度後未於同一短時間窗碰到遭遇便越過多個月；這表示實機驗收應
使用 bounded clock／到站條件，不應把固定秒數當成自然路徑 parity。該次不推翻已通過
的策略輸入接縫，也不把它升格為完整長程測試。

### 存檔／讀檔

在正常戰略畫面開系統視窗後執行：

```text
4 → S → 1 → Enter → L → 1 → Enter
```

remake 以 `-save-file /tmp/remake-save.dat` 啟動，實際產生 88,832 bytes overlay，
畫面先顯示「已儲存第 1 槽」，再顯示「已讀取第 1 槽；信賴度 255」。

- [`expert-remake-save-03-saved.png`](../images/expert-remake-save-03-saved.png)：儲存成功
- [`expert-remake-save-05-loaded.png`](../images/expert-remake-save-05-loaded.png)：讀取成功
- 保存區的原始 DOS/V 檔案以 `:ro` 掛載，沒有被覆寫。

## 4. 戰術畫面邊界

目前工作樹建置以 `-open-siege` 做了短 GUI smoke：戰場、城壁耐久、雙方兵力與 1–6
命令列都能繪製，送出 `2` 號攻擊按鍵後程式仍維持執行。

- [`expert-remake-debug-siege-start.png`](../images/expert-remake-debug-siege-start.png)
- [`expert-remake-debug-siege-attack.png`](../images/expert-remake-debug-siege-attack.png)

這兩張是 **debug fixture smoke**，不是自然流程證據。正常無 `-open-*` 的遭遇→戰術→
戰後結果→回戰略流程已有 [`docs/playtest/09-wlgame-normal-tactical-path.md`](09-wlgame-normal-tactical-path.md)
與固定種子 17 的狀態測試；本輪不把受控 tactical smoke 改寫成新的原版逐像素 parity。

## 5. DOSBox／remake 畫面比較

| 比較項目 | 原版 DOSBox／DOSBox-X | remake | 判定 |
|---|---|---|---|
| 邏輯畫面 | 原版 640×400 | 1280×800，為 2× 邏輯畫面輸出 | 尺寸轉換一致；不是同解析度 pixel diff |
| 地圖／道路／據點 | PC-98 oracle 的原版圖與既有資料 | DOS/V 自然 HUD 與右側情報欄 | 結構與資料接線可比；鏡頭／日期不同，不宣稱逐像素 |
| 常駐 HUD | 原版頂端 banner／日期／圖示 | banner、命令列、minimap、自勢力數值欄 | remake 依影片與原版幾何重建，保留明確現代化欄位 |
| 事件／戰鬥入口 | 原版玩法 oracle 受 DOS/V 密碼與本輪 PC-98 輸入橋接限制 | 正常策略路徑與戰術 GUI 已可走 | remake 可玩性通過；原版同狀態 pixel parity 保持開放 |

## 6. 尚未關閉的驗證

- DOS/V 密碼頁後的原版自然流程，沒有合法密碼不驗證。
- PC-98 原版在目前無 window manager 的容器內需要修正輸入橋接，才能重播完整
  新局→戰術流程；既有 oracle 截圖與 `docs/playtest/06` 保留。
- Windows／macOS 原生 GUI 與三平台正式包仍是獨立 release gate。
- 完整長程遊戲測試依使用者先前要求不執行。
