# 89 — 戰後的損害報告改成可關的選項

**狀態：CONFORMED（2026-08-26 實作並實跑驗過）。**

- 日期：2026-08-26
- 出處：**原版戰術結束後沒有損害報告**（使用者裁定 2026-08-26）。
  remake 的結果畫面多了一行「攻城損害　N」，那是 remake 自己加的診斷資訊。
- 推論等級：不涉及原版事實——原版沒有這個畫面，所以這一整頁都是 remake 差異。

## 1. 現況

`drawBattleResult` 在攻城戰時多畫一行：

```go
if p.Mode == combat.Siege {
    g.td.Draw(screen, fmt.Sprintf("攻城損害　%d", b.CityDamage(p.CityWall)), …)
}
```

那一行的用處是**驗收**：`CityDamage` 是規則層算出來的，畫在螢幕上才看得出
它有沒有接對。對玩家而言它是多餘的，而且原版沒有。

## 2. 改法

| 項目 | 決定 |
|---|---|
| 開關位置 | 系統選單第 **8** 列「損害報告」 |
| 預設 | **關** |
| 值 | 「 開 」／「 關 」（與第 7 列「主君編成」的 `可`／`不可` 同一種兩值列）|
| 關掉時 | 那一行不畫，其餘結果照舊 |
| 命令列 | `-siege-damage`（驗收用，與 `-lord-corps` 同一種）|

⚠ **新列加在最後，不是插在中間。** 與 [`76`](76-lord-not-in-formation.md) 同一條理由：
插進去會把原版那六列往下推，於是**原版六列沒有一列還在原座標**，
[`../playtest/39`](../playtest/39-system-window-parity.md) 對那六列的比對就不再成立。
加在後面的話那六列一個 px 都不動，代價只是視窗再高一個列距（216 → 240）。

## 3. 滑鼠點一下就關掉報告

結果畫面現在只認 Enter／Space。**任何一個滑鼠鍵按下去也要關掉**——
玩家剛打完一場仗，手還在滑鼠上。

⚠ 這裡要用 `IsMouseButtonJustPressed` 而不是 `IsMouseButtonPressed`：
戰鬥是用滑鼠下令的，**最後一個指令的那一下如果還按著**，
用「持續按著」判定會讓結果畫面在出現的同一幀就被關掉，
玩家根本看不到它。

## 4. 改動

| 檔 | 內容 |
|---|---|
| `cmd/wlgame/strategyhud.go` | `sysRows` 7 → 8、標籤表加「損害報告」、`sysRowDamageReport` 常數、`dispatchSystemRow` 的 toggle、值格 |
| `cmd/wlgame/battle.go` | 那一行改成有開關才畫；結果畫面加滑鼠關閉 |
| `cmd/wlgame/main.go` | `-siege-damage` 旗標、`game.damageReport` 欄位 |

## 5. 驗證（2026-08-26）

| 項目 | 結果 |
|---|---|
| `TestDamageReportDefaultsOffAndTogglesFromSystemRow` | 預設關；第 8 列左右鍵都 toggle；值是「 開 」／「 關 」|
| `TestSystemMenuKeepsOriginalSixRowsInPlace` | 原版那六列的索引與值格 y **一格都沒動**；remake 的兩列在最後 |
| 實跑 | 系統選單八列畫得出來，最後一列是「損害報告　關」（`docs/images/system-menu-8rows.png`）|

⚠ 值格只有 48 px（六個半形字）。「 開 」／「 關 」是三個半形寬的全形字
加前後空白，與旁邊那幾列同寬——`sysValueLine` 會再補一次置中。

## 6. 未解

| 缺口 | 下手點 |
|---|---|
| 設定不進存檔 | 與語言同一條（`86` §7）：原版存檔沒有這一欄，要記住偏好得另存 remake 自己的設定檔 |
| 結果畫面本身原版沒有 | 原版打完直接回戰略畫面。要不要整頁拿掉是另一個裁定，本規格只讓多出來的那一行可關 |


## 7. 2026-09-08 戰後結果頁選項增補（CONFORMED）

使用者明確採用：結果頁變成選項、預設關閉，開啟時可設定幾秒後自動關閉，
避免每場都要確認。此決定取代 §6「整頁是否拿掉待裁定」，保留先前證據與歷史。
原版正常退卻直接返回，證據是 `workplace/battle-return-20260908/original/timeline.txt`
與 oracle.log；dosgolem fd746405cf0ef9e535121e215b935b0bf59cd1de，dosv KI.EXE
SHA-256 fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868。
本次不改規則、不新增反組譯，結果頁本身明列為重製便利性選項。

- 系統選單末尾加第 9 列「戰後結果」，值依序為關、3 秒、5 秒、10 秒、15 秒、30 秒。
- 預設關閉：戰鬥完成後自動結算並返回，不繪製結果頁、不等待任何確認。
- 開啟時從首次繪製結果頁開始，以單調實際時間計時，不受戰術速度或繪圖更新率影響。
- 到期自動結算；玩家仍可用新按下的滑鼠鍵、Enter 或 Space 提早關閉。
- 第一次顯示前不吃上一個戰術指令的確認輸入；每場重新開始計時。
- 設定沿用現行系統選單的工作階段生命週期，不改原版存檔或新增偏好檔案。
- 既有損害報告只控制結果頁中的額外一行；不會自行開啟結果頁。
- 只改 AppImage／Windows／macOS 桌面 UI，Android 不變。
- 驗收包含預設無確認、啟用後自動關閉、提前關閉、連續兩場重設計時，以及九列不超出 640×400。

正式 AppImage 的預設關閉與三秒模式皆由正常遭遇驗收，見 [playtest/112](../playtest/112-desktop-battle-result-options.md)。
