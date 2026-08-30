# 「經典再現」推廣片

**狀態：已產出研究／推廣用 60 秒比較片；不把代表幀比較宣稱為同狀態逐像素 parity。**

- 日期：2026-08-11
- 影片：[`dist/promo/wolong-remake-classic-revival.mp4`](../../dist/promo/wolong-remake-classic-revival.mp4)
- 規格：60 秒、1280×720、H.264／AAC、44.1 kHz 立體聲
- SHA-256：`ca474a01c874739f9e4047e1bc5912ad72b4bdd73fb9eb6c2a64c58278c8d695`
- 重現腳本：[`tools/promo_classic_revival.sh`](../../tools/promo_classic_revival.sh)
- 代表畫面：[`classic-revival-frame.png`](classic-revival-frame.png)

## 敘事與來源

影片主軸是「經典再現」：先放使用者提供的 YouTube 遊玩錄影代表幀，再與目前
DOS/V remake 的可重播畫面並排，依序展示自然策略 HUD、事件 2–5 TALK 選擇、戰術、
投射物、戰果與戰後回到自然畫面。

- 原版來源：使用者提供的 YouTube 錄影
  [`af6xqcicXoI`](https://www.youtube.com/watch?v=af6xqcicXoI)，只取已保存的代表幀。
- 原版欄位在片中標為「YouTube 代表幀」；沒有把原版影片、執行檔、`BGM.DAT`、
  `SOUND.DAT` 或其他原版資產封裝進影片／發行包。
- Remake 欄位來自已驗收的實際截圖，固定 `seed=17`；這是流程與視覺方向展示，
  不是同日期、同輸入、同戰局的 parity 證明。
- 配樂沿用本專案原創 `tools/promo_score.py` 合成聲，不取樣原版音樂。

## 可重現的重播基準

依 `~/.claude/knowledge-base/retro/dosbox-game-configs.md` 的 deterministic DOSBox
原則，原版參考條件以 DOSBox-X／DOSBox 設定中的固定節拍標記保存為：

```text
machine=pc98（若測試 PC-98 oracle）
core=normal
cputype=486
cycles=20000
```

這些設定的用途是讓後續重新拍攝時時序可重現，不代表目前 DOS/V 密碼保護已被繞過，
也不代表 YouTube 代表幀就是上述條件下的無損擷取。影片結尾同步標出 `remake seed=17`，
讓觀眾知道兩側各自的來源與限制。

## 影像驗收界線

`classic-revival-natural-side-by-side.png` 與 `classic-revival-natural-difference.png`
保留自然畫面並排／差異證據。差異圖可用來檢查外框、HUD 幾何、色彩層級與畫面方向；
由於原版代表幀與 remake 的日期、鏡頭、輸入狀態及壓縮流程不同，不將 raw 像素數值解讀
為原版同狀態 parity。
