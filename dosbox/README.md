# DOSBox-X 設定

兩個版本各一份。**設定不是自己調的**，出自 joncampbell123 的 `msdostest`
測試庫（567 款老遊戲的實測 `dosbox.conf`），細節見
`~/.claude/knowledge-base/retro/dosbox-game-configs.md`。

| 檔案 | 版本 | 出處 |
|---|---|---|
| `pc98.conf` | PC-98 日文原版 | msdostest `garyouden-sangoku-seiha-no-kei-hokusho-neo-kobe-pc98-ia`，抓取日 2026-08-07。該目錄的 `__PASS__` 標記為 `20220526-155014`（DOSBox-X master） |
| `dosv.conf` | 松崗 DOS/V 繁中版 | msdostest 沒有收錄，自行撰寫，**未實測** |

## 為什麼不用打包版附的設定

`卧龙传.zip` 裡的 `.jsdos/dosbox.conf` 是 `machine=svga_s3`、`cycles=auto`——
那是打包者為了讓玩家能玩選的，不代表原版需求，而且
**`cycles=auto` 會讓即時制遊戲每次跑到不同的時間點**，截圖對不起來。
本專案一律用固定 `cycles`。

## 兩個與 msdostest 原設定的差異

1. **原設定掛的是 `.hdi` 硬碟映像，我們手上是五片 `.fdi`。**
   `pc98.conf` 改成 `imgmount 0 … -t floppy` ＋ `boot -l a`。
   磁片 D 有 `HDINST.EXE`，之後可以裝成 HD 映像省掉換片。
2. `machine=pc98` **只有 DOSBox-X 支援**，原版 DOSBox 與 Staging 都不行。

## 開發商

msdostest 的目錄名 `…-hokusho-…` 佐證了開發商是 **Hokusho（ホクショー）**，
與維基百科寫的「NEO･GETEN（ホクショー）」一致。這是第三個獨立來源。
