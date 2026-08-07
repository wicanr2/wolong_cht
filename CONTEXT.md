# CONTEXT — 專案脈絡與文件索引

> **這份是全專案的單一入口。** 對話被壓縮、或換一個新 session 接手時，先讀這份，
> 就能重建完整全局，再依索引跳到需要的文件。
>
> 硬規則、目標、工作紀律在 [`CLAUDE.md`](./CLAUDE.md)。這份只管**狀態**。
>
> 最後更新：2026-08-07

---

## 1. 這個專案在做什麼

把 NEO･GETEN《臥龍傳－三國制霸之計》（1994 日文 PC-98 版 / 1995 松崗 DOS/V 繁中版）
完整逆向，在 Go / Ebiten 上重寫引擎，保存松崗繁中版原文並以日文原版對照校訂。

**這一款與前兩個 remake 專案最大的不同：它是即時制。**
時間模型是第一等問題，`docs/mechanics/15-realtime.md` READY 之前不寫規則層。

---

## 2. 一句話現況

**M0 進行中，`TALK.DAT` 已全解（M2 的核心已經拿下）。**
兩版素材入庫、兩版 `KI.EXE` 進 IDA、逐檔比對完成；
`TALK.DAT` 兩版都 byte-for-byte round-trip 通過，日中對照表已產出並找到 15 則譯文缺陷。
`.BRG` 調色盤與**四個圖庫全部解完**（`ICONGRF` 剩兩段）。
**Go 解碼層與素材檢視器已經跑起來**（`internal/assets/*` ＋ `cmd/wlshot`／`cmd/wlview`），
Go 版與 Python 版的輸出逐像素完全相同。
**下一步：`MMAP.*` 大地圖三件**——入口已找到（`docs/re/04`），
但這一組不像圖庫那樣「載入器直接寫死尺寸」，資料在記憶體裡會先被展開，
要先解 `sub_1E48A`。

---

## 3. 現況一覽

### 已完成

| 項目 | 狀態 |
|---|---|
| 松崗 DOS/V 版入庫 | 69 檔 → `workplace/orig/dosv/`（唯讀） |
| **PC-98 日文原版入庫** | 5 片 `.fdi` → 抽出 69 檔 → `workplace/orig/pc98/`（唯讀） |
| FDI 抽檔工具 | `tools/fdi_extract.py`（FAT12 / 1024 B sector，只用標準函式庫） |
| **兩版逐檔比對** | **23 檔 byte-for-byte 相同**、3 檔大小同內容異、其餘不同。見 `docs/re/01` |
| 兩版 `KI.EXE` 反組譯 | `workplace/ida/{dosv,pc98}/KI.EXE.i64`，732 ／ 725 函式 |
| IDA 包裝腳本 | `tools/ida.sh {batch\|raw} {dosv\|pc98} …` |
| 日文說明書初判 | 38 頁掃描，目次與第 3.5 節已讀。見 `docs/reference/01` |
| 網路查證 | 發行資訊、四劇本、規模。見 `CLAUDE.md` §2 |
| **`TALK.DAT` 全解** | **READY**。`docs/formats/01`。兩版 byte-for-byte round-trip 通過，1,022 則、零逃脫位元組 |
| **日中對照第一批** | `docs/reference/02`。15 則譯文缺陷（漏變數 6、漏名詞 3、變數換號 1、錯位 4、對調 2） |
| 訊息抽取工具 | `tools/talkdat.py`（dump／export／build／verify／diff） |
| **`.BRG` 調色盤全解** | **READY**。`docs/formats/02` ＋ `docs/re/02`。通道順序 B,R,G、4 bit、每組 16 色、亮度 0–16，兩版機器碼互證 |
| **顯示模式定案** | 兩版都是 **16 色 4 平面 planar**（DOS/V VGA GC+Seq、PC-98 GRCG）→ `*GRF.DAT` 是 4 bpp |
| **四季調色盤** | `GAMEPAL` banks 0–3 只差色號 14：灰綠→鮮綠→橙褐→雪白 |
| **四個圖庫全解** | `docs/formats/03`。`KAOGRF` 150×64×64、`KYOGRF` 15×96×96（據點景觀）、`IVENTGRF` 3×288×176（劇情過場）**全部餘 0 B 並渲染驗收**；`ICONGRF` 是四段組合檔，段 0（640×32 標題橫幅）與段 2（192×128 縮小地圖）已解 |
| **中文化缺口第一項** | `docs/reference/03`。主畫面標題橫幅寫的是日文「臥竜伝」，`ICONGRF` 兩版相同 → **松崗版沒重繪** |
| **DOSBox-X 設定** | `dosbox/`。出自 msdostest（實測 `__PASS__`），`core=normal` ＋ 固定 `cycles=20000` → **即時制的可重現性解決了** |
| **DOS/V 版實跑** | `tools/dosbox.sh` ＋ `docs/playtest/01`。**字型懸案結案**（不用自備，`YNFONT.EXE` 自帶）、**防拷確認**、淡入行為與 `docs/re/02` 對上 |
| **PC-98 版實跑（oracle 建立）** | `tools/dosboxx.sh` ＋ `docker/dosboxx/` ＋ `docs/playtest/02`。**沒有防拷，開場完整播放**。掛目錄跑（原版本來就裝硬碟），五片分掛會報「ファイルが異常です」 |
| **逐片抽檔修正** | `workplace/orig/pc98_discs/{A..E}` 是權威來源；合併目錄 A 片優先。踩到「同名檔被後片靜靜蓋掉」（`YNSHELL.COM` 2699 vs 2675 B），工具已加比對警告 |
| **`.MAP`/`.SCH` 容器格式** | `docs/formats/04`。索引 (offset,length) ＋ 壓縮片段，六個檔的末端全部等於檔長。**只有松崗自加的檔用這格式**，原版的 `MMAP`/`BATTLE` 是裸資料 |
| **PC-98 有無符號表** | **沒有**（使用者提問）。`KI.EXE` 兩版都無 CodeView／Borland debug 資訊、無原始檔名殘留。PC-98 的價值在「兩版 diff」與「沒有防拷」，不在符號 |
| **Go 解碼層** | `internal/assets/{palette,gfx,text,library}`。三份 READY 規格全部接上，測試綠（含 Go 版的 byte-for-byte round-trip）|
| **素材檢視器** | `cmd/wlview`（Ebiten，實跑截圖驗過）＋ `cmd/wlshot`（無頭出 PNG）|
| **交叉驗證** | Go 版與 Python 版的圖庫輸出**逐像素完全相同**（KAOGRF 1.84 MB、KYOGRF 414 KB）|
| **Go 建置環境** | `tools/go.sh`（docker）＋ `tools/shot.sh`（Xvfb 截圖）＋ `docker/go/Dockerfile` |
| **AI 機制文件開張** | `docs/mechanics/70-ai.md`（＋ `00-index.md`）。內容目前全來自日文說明書第 10、11 章，每條都標了等級 |
| GitHub repo | https://github.com/wicanr2/wolong_cht（private） |

### 進行中／受阻

| 項目 | 狀態 |
|---|---|

| **DOS/V 防拷** | **確認有**（`docs/playtest/01`）：查說明書第 NN 頁的四字密碼，擋住 DOS/V 的 oracle。**繞路：PC-98 版沒有這一關**，規則類的 oracle 改跑 PC-98 |
| 日文說明書全文判讀 | 38 頁只讀了 3 頁。純掃描無文字層，要 OCR 或逐頁看 |

---

## 4. 文件索引

| 文件 | 內容 |
|---|---|
| `CLAUDE.md` | 目標、原則、硬規則、里程碑、工作紀律 |
| `docs/re/01-first-recon.md` | 首輪偵查：檔案清單、執行結構、兩版比對 |
| `docs/formats/01-talk-dat.md` | **READY** — `TALK.DAT` 訊息表格式 |
| `docs/formats/02-brg-palette.md` | **READY** — `.BRG` 調色盤格式 |
| `docs/formats/03-grf-images.md` | **READY**（僅 `KAOGRF`）— 圖庫格式 |
| `docs/re/02-palette-routine.md` | 兩版調色盤常式的反組譯 |
| `docs/re/03-image-blitter.md` | 圖庫載入器與 VGA 四平面繪製常式 |
| `docs/re/04-mmap-entry-points.md` | 大地圖的入口點與記憶體佈局（進行中） |
| `docs/mechanics/00-index.md` | 機制索引與推論等級定義 |
| `docs/mechanics/70-ai.md` | **電腦 AI 的判斷邏輯**（蒐集中） |
| `docs/reference/01-jp-manual.md` | 日文原版說明書判讀紀錄（逐頁累加） |
| `docs/reference/02-jp-cht-diff.md` | 日中對照：15 則譯文缺陷 |
| `docs/reference/03-baked-japanese.md` | 燒進美術裡的日文（松崗沒重繪的部分） |
| `dosbox/` | 兩版的 DOSBox-X 設定 ＋ 出處 |
| `docs/playtest/01-dosbox-dosv.md` | DOS/V 版首次實跑：字型結案、防拷發現 |
| `docs/playtest/02-dosboxx-pc98.md` | PC-98 實跑：oracle 建立、合併抽檔陷阱 |
| `docs/formats/04-map-sch-container.md` | `.MAP`/`.SCH` 容器格式 |
| `README.md` | 公開的專案入口 |
| `translations/extract/talk-{dosv,pc98}.json` | 抽出的 1,022 則訊息（兩版） |
| `tools/fdi_extract.py` | PC-98 FDI 磁片抽檔 |
| `tools/ida.sh` | IDA Pro 9.4 headless 包裝 |

---

## 5. 術語表

| 詞 | 意義 |
|---|---|
| **dosv** | 松崗 1995 DOS/V 繁體中文版。專案裡一律用這個代號，不寫「中文版」 |
| **pc98** | NEO･GETEN 1994 PC-9801 日文原版 |
| `YN*` | 開發方的工具鏈前綴（`YNSHELL`／`YNFONT`／`YNSOUND`／`YNMOUSE`），兩版都有 |
| `KI.EXE` | 遊戲本體。兩版同名 |
| 據點 | 遊戲內的城池／關口／戰場，共約 300 處。**不寫「城市」** |
| 軍團 | 玩家指揮的部隊單位，編成含主將／前鋒／左翼／右翼／左備／右備 |
| 戰略／戰術 | 原版用語。戰略＝大地圖政略，戰術＝戰場。**不改寫成「大地圖／戰鬥」** |

---

## 6. 已被推翻的斷言

> 每推翻一條就記一列，寫清楚誰推翻誰、憑什麼。**空著不是好事，是還沒開始驗。**

| 日期 | 原斷言 | 推翻依據 |
|---|---|---|
| 2026-08-07 | 「音源未解，候選 AdLib／SB／PC speaker」 | 日文說明書第 3.5 節明寫「FM音源とSSG音源のバランス」→ PC-98 側是 OPN。**但只推翻 PC-98 那一半，DOS/V 側仍未解** |
| 2026-08-07 | 「字型來源三個候選都沒驗」 | PC-98 版只用 843 B 的 `YNFONT.COM` ＋ 1,216 B 的 `FONTGRF.DAT` → 靠字型 ROM；DOS/V 的 60,888 B `YNFONT.EXE` 是在補這一塊 |
| 2026-08-07 | 「`TALK.DAT` 的訊息中間夾**控制位元組**做變數插入」 | 全解之後確認：訊息內除了 `0x00` 沒有任何控制位元組，變數是 ASCII `\1`–`\7`。**當初拿「Big5 掃出一堆短片段」反推機制，那是間接證據**——結論碰巧對，推導是錯的 |
| 2026-08-07 | 「`dosv` #225 有錯字」（`這哪<f9>堿O什麼好策略？`） | `F9D8` 是 Big5 罕用字區的「裏」，原文完全正確。**錯的是解碼器用了 `big5` 而不是 `cp950`**。差一點把工具的缺陷寫成原版的缺陷 |
| 2026-08-07 | 「`TALK.DAT` 偏移表 2,048 ÷ 4 ＝ 512 筆」 | 實際是 uint16、1024 筆。假說當時就標了「待驗」，驗完是錯的 |

---

## 7. Worklist（狀態的單一真相來源）

### 7.0 下一步（2026-08-07 基準）

**`TALK.DAT` 已完成**（`docs/formats/01` READY，兩版 round-trip 通過）。

**下一項：`.BRG` 調色盤 ＋ `*GRF.DAT` 圖庫。**

理由：`.BRG` 是全部素材裡最小的（48–576 B）、四個檔兩版全同，
是驗證整條「解碼 → 渲染 → 對照 DOSBox 截圖」管道最便宜的起點。
管道通了再上 `*GRF.DAT`（307 KB 的 `KAOGRF` 是最大的一個）。

> **原本以為要先建 DOSBox 才驗得了調色盤——結果不用。**
> 反組譯直接給出通道順序，比截圖比對更硬。
> 教訓：**撞到「需要視覺 oracle」的判斷時，先問一次「機器碼有沒有直接寫」。**
> 截圖仍然有用（驗證整張圖畫對），但不是解格式的前置條件。

### 7.1 M0 剩餘

- [x] ~~即時制的可重現性~~ — msdostest 的設定用 `core=normal` ＋ 固定 `cycles`，已解
- [ ] DOSBox-X docker 化並實跑兩版（設定已在 `dosbox/`）
- [x] ~~np2 環境~~ — 不需要，DOSBox-X `machine=pc98` 兩版通吃
- [x] ~~字型來源結案~~ — 不用自備，`YNFONT.EXE` 自帶
- [x] ~~防拷檢查~~ — DOS/V 有、PC-98 沒有；oracle 改跑 PC-98
- [x] ~~弄一個 DOSBox-X~~ — `docker/dosboxx/`，PC-98 版已跑起來
- [ ] **實測即時制的可重現性**：同一串操作跑兩次要得到同一張圖。
      這是 oracle 能不能用的關鍵，`cycles` 固定了但還沒驗
- [ ] 推進到標題選單與大地圖（開場動畫很長，要更長的自動化序列）
- [ ] 日文說明書 38 頁全判讀 → `docs/reference/01`

### 7.2 M1 起手順序（依投報排）

> **前三項走通的路徑**（`docs/re/03` §5）：grep 檔名偏移的立即值 → 載入器讀出
> 「每筆多少 byte、怎麼定位」→ 繪製呼叫端的 `ax` 拆成寬高 → 渲染，餘數要是 0。
> **`MMAP.*` 這一組不適用**——資料先展開再畫，見 `docs/re/04`。

1. ~~`.BRG` 調色盤~~ ✅ **完成**，`docs/formats/02` READY
2. ~~`*GRF.DAT` 圖庫~~ ✅ **完成**（`ICONGRF` 段 1／段 3 待補）
3. `MMAP.*` 大地圖三件（兩版全同）
4. `BATTLE.*` 戰場（兩版全同）
5. `SINARIO.DAT` ／ `SAVE.DAT`（兩版大小同內容異 → diff 定位字串欄位）

> **從 PC-98 側入手解格式**：640×400 固定規格，比 DOS/V 少一層猜測。
> 但結論要在 dosv 上覆驗（`CLAUDE.md` §8 第 9 條）。

---

## 8. 已建好的工具

| 工具 | 用法 |
|---|---|
| `tools/fdi_extract.py` | `python3 tools/fdi_extract.py <image.fdi> <輸出目錄>` |
| `tools/ida.sh` | `tools/ida.sh batch dosv KI.EXE` ／ `tools/ida.sh raw pc98 idat -A -S/work/tools/x.idc KI.EXE.i64` |
| `tools/talkdat.py` | `dump` ／ `export` ／ `build` ／ `verify` ／ `diff`。**驗收指令是 `verify`，要求 byte-for-byte** |
| `tools/brg.py` | `info` ／ `swatch`。純標準函式庫的 PNG 輸出 |
| `tools/grf.py` | `sheet` ／ `one`。**`sheet` 會印「餘 N byte」，不是 0 就代表尺寸猜錯** |
| `tools/ida_xref.idc` | 查 IDA 的 xref 圖。⚠ 立即值形式的位址參考 IDA 不建 xref，回零筆不等於沒人用（`docs/re/03` §0） |
| `tools/go.sh` | `tools/go.sh test ./...`。image 沿用 demonwinter-go，volume 是自己的 `wl-gomod`／`wl-gobuild` |
| `tools/shot.sh` | `tools/shot.sh out.png KEYS=Right,Down [參數]`。Xvfb + xdotool，驗呈現層 |
| `tools/dosbox.sh` | `tools/dosbox.sh dosv "wait:2;type:START;key:Return;wait:10;shot:x"`。DOS/V oracle（會撞防拷）|
| `tools/dosboxx.sh` | `tools/dosboxx.sh "wait:20;click:320,200;shot:x"`。**PC-98 oracle，無防拷**。timeline 支援 wait／key／type／click／shot |

**還沒建但 `CLAUDE.md` §5.1 要求的**：`tools/addr.py`（位址反查三張表）、
`tools/dump_func.py`（dump 時自動翻語意 ＋ grep docs）、`tools/ida_xref.idc`。
開始讀組語之前要先補上。
