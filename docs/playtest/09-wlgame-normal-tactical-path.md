# 09 — `wlgame` 正常遭遇到戰術戰場

**狀態：正常開局、正常 AI 遭遇、戰鬥指揮選單、戰術戰場、戰後結果報告與 GUI 戰後回戰略已驗收；原版同狀態逐像素對拍仍未完成。**

- 日期：2026-08-09
- 執行環境：Docker／Xvfb，專案既有 `docker/go/Dockerfile` 建成的 `wolong-go:20260809`
- 原始輸入：`workplace/orig/dosv/SINARIO.DAT`、`workplace/eten/`，唯讀掛載
- 亂數：`-seed 17`，只固定驗收亂數來源
- 啟動：`cmd/wlgame`，未使用 `-open-battle`、`-open-battle-choice`、`-open-form` 或其他 `-open-*` 旗標
- 截圖：使用 `wlgame -shot` 內建 PNG 輸出；新檔案由 UID/GID `1000:1000` 產生

## 1. 正常輸入序列

先沿用 [`08-wlgame-normal-strategy-path.md`](08-wlgame-normal-strategy-path.md) 的真實劇本流程：

```text
A；Enter × 2；3；Down + Space × 5；Enter
M；Enter × 2；Down × 22；Enter × 2
= × 64
```

這條路徑在 196 年 6 月 28 日進入「呂布 對 曹操／攻城／戰鬥指揮／委任」。
再按一次 `Enter` 選第一列「戰鬥指揮」，進入攻城戰術畫面；在戰術畫面按 `2`
下達原版編號 2 的「攻擊」命令。戰術時間持續推進，沒有暫停或傳送。

## 2. 證據

| 階段 | 證據 | 可確認內容 |
|---|---|---|
| 正常 AI 遭遇 | [`wlgame-ai-afterpatch.png`](../images/wlgame-ai-afterpatch.png) | 196/6/28；呂布對曹操；戰鬥指揮／委任選單；沒有 `-open-*` |
| 戰術戰場 | [`wlgame-ai-battle-afterpatch.png`](../images/wlgame-ai-battle-afterpatch.png) | 真實攻城地形、等角城壁、雙方兵力橫幅、城壁耐久與 1–6 戰術命令列 |
| 下達攻擊後 | [`wlgame-ai-battle-attack-afterpatch.png`](../images/wlgame-ai-battle-attack-afterpatch.png) | 196/6/28；固定種子 17、正常鍵盤、戰略／戰術速度 16；輸入 `2` 號攻擊命令後仍在戰鬥中，攻方 112 兵、守方 100 兵、城壁耐久 1830 |
| 戰後結果報告 | [`wlgame-ai-battle-result.png`](../images/wlgame-ai-battle-result.png) | 事件佇列接入後重拍；守方勝；攻方 5590→0、守方 1000→1000；攻城損害 0；按 Enter 才回戰略畫面 |
| 狀態層戰後結算 | `TestNormalScenarioTacticalBattleTerminates` | 同一真實 `SINARIO.DAT`／固定種子 17／真實 `BATTLE.*`；第 549 幀守方勝，攻方 0、守方 100；`ResolvePending` 後清除 pending 遭遇 |
| GUI 戰後回戰略 | [`wlgame-ai-postbattle.png`](../images/wlgame-ai-postbattle.png) | 同一正常鍵盤流程在戰術完成後按 `Enter`；回到戰略地圖，顯示曹操情報與持續行軍事件；無 `-open-*` |

戰術核心本輪另以 `internal/rules/tactical` 測試固定了：`sub_1B618` 的亂數／戰力
近戰公式、突擊 `+0xC8`、大將碰撞旗標、隊長離場後七名隊員退卻、待機兵清除，
以及 `sub_1B00D` 抵達中繼點後才消費繞路點。這些規則不是由截圖推測，而是由
`KI.EXE` 的 IDA 線性位址證據回填，詳見 [`docs/re/11-tactical-battle.md`](../re/11-tactical-battle.md)。
`sub_1AD7F` 的 `CH=0x20` 特殊效果另由
`TestSpecialProjectileUsesCH20AndFallsVertically` 與
`TestClimbingInfantryCanUseSpecialProjectile` 驗證；這兩個測試是規則層證據，
不是這三張中途 GUI PNG 的畫面推測。

## 3. 結果報告與狀態回寫

`CorpsEvent` 現在沿事件流保留戰前／戰後戰略兵力點數與攻城損害；GUI 結果視窗
使用同一場 `tactical.Battle.Result()`，不從截圖或戰略層重新猜測。畫面上的 5590、
1000 是實際人數，狀態測試的 0、100 是每 10 人一點的戰略層數值；兩者換算一致。
按 Enter 才呼叫 `ResolvePending`，之後事件列也會保留勝負、兵力變化與據點損害。

## 4. 尚未完成

- GUI 戰後結果與狀態回寫已留證；仍需取得原版同狀態戰術畫面，完成正常與原版的逐畫面、
  戰後兵力與城壁損傷對照。
- `sub_1AD7F` 的 `CH=0x20` 高度／近距離效果已有 `shootSpecial` 的規則測試；仍需
  確認 `+0x1E` 原始欄位的完整來源、與一般 `sub_1AD2D` 箭路徑的動畫／逐狀態對拍。
