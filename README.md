# 臥龍傳 remake

把 NEO･GETEN《臥龍傳－三國制霸之計》（1994 日文 PC-98 版 ／ 1995 松崗 DOS/V 繁中版）
完整逆向，在 Go / Ebiten 上重寫引擎，保存松崗繁中版原文並以日文原版對照校訂。

定位是文化資產保存。**不散布原版執行檔、資料檔、美術或音樂**——
公開產出只有引擎程式碼與校訂紀錄，玩家自備合法原版。

| | |
|---|---|
| 原名 | 臥竜伝 三国制覇の計 |
| 開發／原發行 | NEO･GETEN（ホクショー） |
| 日文版 | 1994，PC-9801 |
| 繁中版 | 松崗，1995，DOS/V |
| 類型 | **即時制**戰略。玩家扮演軍師，不是君主 |
| 規模 | 武將 146 人、城池 192 個、據點近 300 處 |
| 劇本 | 呂布歸天、赤壁之戰、蜀地偏安、劉禪繼位 |

## 現在做到哪裡

**逆向與格式解析階段。遊戲本體還沒開始寫。**
狀態的單一真相來源是 [`CONTEXT.md`](CONTEXT.md)。

### 已解出的格式

| 格式 | 狀態 | 文件 |
|---|---|---|
| `TALK.DAT` 訊息表 | READY，兩版 byte-for-byte round-trip | [`docs/formats/01`](docs/formats/01-talk-dat.md) |
| `.BRG` 調色盤 | READY | [`docs/formats/02`](docs/formats/02-brg-palette.md) |
| `*GRF.DAT` 圖庫 | READY（`ICONGRF` 剩兩段） | [`docs/formats/03`](docs/formats/03-grf-images.md) |
| `MMAP.MDL` 地形圖塊 | READY，256 塊 16×16 | [`docs/formats/05`](docs/formats/05-mmap-worldmap.md) |
| `MMAP.MAP` 世界地圖 | READY，RLE → 384×256 格 | [`docs/formats/06`](docs/formats/06-mmap-rle.md) |
| `.MAP`/`.SCH` 容器 | 索引層 READY | [`docs/formats/04`](docs/formats/04-map-sch-container.md) |
| `BATTLE.*` 戰場 | 分段結構已解，像素格式未解 | [`docs/formats/07`](docs/formats/07-battle.md) |

### 引擎已經畫得出遊戲畫面

`cmd/wlview` 的大地圖模式：384×256 格的世界地圖，可捲動，頂端是原版的
標題橫幅（`ICONGRF` 段 0）。

![大地圖（春）](docs/images/wlview-world-spring.png)

**按 `4` 換到冬天** —— 只換調色盤的色號 14 這**一個顏色**，
21 萬個像素改變，而樹林、河流、道路、城池全部保持原色。
這是原版 1994 年在 16 色機器上的做法，remake 照做（`docs/formats/02` §4）。

![大地圖（冬）](docs/images/wlview-world-winter.png)

![素材檢視器](docs/images/wlview-kyogrf.png)

### 已有的程式

```
internal/assets/palette   .BRG 解碼（純函式，不認識 Ebiten）
internal/assets/gfx       *GRF.DAT 4bpp planar 解碼
internal/assets/text      TALK.DAT 解析與寫回
internal/assets/rle       原版的 RLE 解壓（MMAP.MAP 用）
internal/assets/world     大地圖：384×256 格 + 256 塊地形圖塊
internal/assets/library   把素材目錄載成可檢視的項目（不 import Ebiten）
cmd/wlshot                解素材成 PNG，無頭環境可跑
cmd/wlview                Ebiten 互動檢視器（素材模式 / 大地圖模式，Tab 切換）
```

分層是刻意的：**Ebiten 在 init 期就要求顯示器**，
所以解碼層與載入層一律不 import 它，否則無頭環境連截圖工具都跑不起來。

## 怎麼跑

原版素材不隨本專案散布。自備之後放進 `workplace/orig/dosv/`（松崗版）
或 `workplace/orig/pc98/`（PC-98 版）。PC-98 的磁片映像可以用
`tools/fdi_extract.py` 抽出來。

**建置一律走 docker**，不污染主機環境：

```sh
tools/go.sh test ./...                       # 測試
tools/go.sh run ./cmd/wlshot -list           # 列出素材
tools/go.sh run ./cmd/wlshot -asset 0 -sheet 15 -out kao.png
tools/go.sh run ./cmd/wlview                 # 互動檢視器
tools/go.sh run ./cmd/wlview -world          # 直接開大地圖
tools/shot.sh out.png KEYS=Right,Down        # headless 截圖驗收
```

Python 工具（只用標準函式庫，不裝套件）：

```sh
python3 tools/fdi_extract.py <image.fdi> <輸出目錄>
python3 tools/talkdat.py verify workplace/orig/dosv/TALK.DAT cp950
python3 tools/brg.py swatch workplace/orig/dosv/GAMEPAL.BRG out.png
python3 tools/grf.py sheet workplace/orig/dosv/KAOGRF.DAT 64 64 \
        workplace/orig/dosv/GAMEPAL.BRG 0 out.png 15
```

Python 與 Go 兩套解碼器是刻意重複的——**兩邊的輸出逐像素比對過，完全相同**。
獨立實作互相驗證，比單一實作加測試更有說服力。

反組譯用 IDA Pro 9.4（`tools/ida.sh`），DOSBox-X 設定在 [`dosbox/`](dosbox/)。

## 這個專案的兩條硬性原則

1. **完整性 > 投報。** 不得以「成本高、效益低」為由跳過任何素材、任何格式。
2. **SDD：spec 齊了才實作。** 反組譯 → 收攏成規格 → 才動手寫程式。
   只有標 READY 的規格可以動手。

細節見 [`CLAUDE.md`](CLAUDE.md)。

## 兩版對照的副產品

PC-98 日文原版與松崗繁中版是同一份程式的兩次編譯，
**23 個資料檔 byte-for-byte 完全相同**。這讓兩件事變得很便宜：

- **日中對照**：`TALK.DAT` 兩版索引一一對應，
  已掃出 [15 則譯文缺陷](docs/reference/02-jp-cht-diff.md)
  （漏變數、漏名詞、`#192`–`#195` 錯位、`#257`/`#258` 對調）。
- **哪些美術被重繪過**：diff 直接告訴你。
  順帶找到主畫面標題橫幅上寫的是日文「臥竜伝」——
  [松崗版根本沒重繪](docs/reference/03-baked-japanese.md)。
