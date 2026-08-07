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
| **⭐ 說明書 p.3–p.14 讀完** | 拿到說服系統全套判定（國力＝據點數、疲弊＝資金<0）、時間暫停規則、四劇本起始日、六位編制。見 `docs/mechanics/10`、`15`、`70` |
| **⭐⭐ `15-realtime.md` 標 READY，M5 的閘解除** | 反組譯出整條時間鏈（`docs/re/06`）：五層單位（子刻9→時24→日→月→年，**一日 216 tick**）、每月天數表（**二月固定 28**）、速度節流、月結 `sub_15358`（**勢力表 22×64 B**、「次月末生效」＝一次 4 word 複製）、季節是 3/6/9/12 月前 16 天的**漸變** |
| **⭐⭐ 好戰等級找到了：勢力記錄 `+0x28`** | 0–15，呂布 15、曹操 14、劉備 4、劉表 1、劉禪 0。**同一君主跨劇本是常數**，且對上說明書點名的兩個例子。連帶解出勢力表 22 欄位、據點表 `+1`/`+26`、武將表 `+28`。`docs/formats/08` §1.5 |
| **⭐ 三種災害的觸發公式全解 ＋ 找到一個原版 bug** | 火災比防災值、暴動比上昇值，**免疫門檻是 64 不是說明書暗示的 0**；暴風雨的「靠海」判定只比 X 的低位元組，**最東邊 49 個據點反而被判成內陸**。`docs/re/07` §17–18 |
| **⭐⭐ 月結的經濟公式全解，`internal/rules/economy` 已實作** | 三條說明書沒寫的公式：**收入依與首都的切比雪夫距離衰減**（÷2/3/4）、**募兵配比依 Y 座標分三區**（北騎59%／中步75%／南弓50%）、**赤字每兵種各扣 \|資金\|/16**。資金是 ±655,000 的有號 24 位。`docs/re/07` |
| **⭐ 規則層開工：`internal/rules/clock`** | 五層時鐘的 Go 實作 ＋ 10 個測試（全過）。測試直接拿反組譯出的常數當期望值：216 tick/日、二月 28 天、四劇本起始日、年封頂 999、季節 16 步漸變、進位階梯不變式 |
| **⭐ 劇本／存檔區塊三段式佈局定案** | 59 B 時鐘＋全域 / 21,056 B 勢力據點武將 / 1,024 B，合計正好 22,208。`docs/formats/08` §0 |
| **⛔ oracle 的滑鼠輸入無效（根因已定位）** | `docs/playtest/04`。pixel diff 證明**滑鼠事件從來沒進入 guest**（兩種天差地遠的設定，畫面 0 像素差異）；先前看到的「進展」都是遊戲自己在跑。**但鍵盤有效**（DOS/V 版靠 `type:START` 啟動過）→ 下一輪走 DOSBox-X 的 mapper 把鍵盤綁成滑鼠動作。仍擋住四件事 |
| **存檔槽 4 個（confirmed）** | `LOAD DATA` 畫面實測，全部 `0年 0月 0日`。日文說明書的說法從「說明書」升到 confirmed。**時間單位是年／月／日** |
| **oracle 可重現性通過** | 同一串操作跑兩次，截圖 **byte-for-byte 相同、0 個不同像素**。這是即時制專案的關鍵閘門 |
| **大地圖畫面拿到** | 推進到 NEW GAME 對話框。**畫面最上方的橫幅正是 `ICONGRF` 段 0**（640×32），位置完全對上 |
| **`MMAP.MDL` 全解** | 256 塊 16×16 地形圖塊，餘 0，與實機地形吻合（`docs/formats/05`）|
| **地圖尺寸定案** | **384 × 256 格**（`sub_1E4CE` 的 `cmp cx, 180h` ＋ 列跨距 0x18 段；98,304 ÷ 384 ＝ 256）|
| **`MMAP.MAP` 全解** | **READY**。RLE（用連續兩個相同 byte 當 run 觸發，無逃脫字元），80,716 → 98,308 B，取前 98,304 ＝ 384×256。畫出來是連貫的中國地圖 |
| **Go 世界地圖層** | `internal/assets/{rle,world}` ＋ 測試（尺寸、圖塊種類≥100、最常見圖塊不超過 1/3、兩版解出長度相同）|
| **劇本／存檔結構全解** | `docs/formats/08`。檔案 ＝ **4 個劇本 × 22,208 B**（4×0x56C0 ＝ 88,832 完全等於檔長）。每劇本 **192 據點 ＋ 127 武將**，記錄都是 32 B、名稱在 `+2` |
| **武將三圍 confirmed** | `+17` 武力、`+18` 統率、`+19` 政治，值域 1–15。**用歷史排名驗**：武力榜馬超/呂布/張飛 15、統率榜司馬懿/諸葛亮 15、政治榜諸葛亮/荀彧 15/14。三榜各自合理且不互相矛盾 |
| **據點記錄欄位（部分）** | 名稱是**定長 6 B、不足補全形空格 `A1 40`**（confirmed）；`+8`/`+10` 是地圖 x/y 座標（強證據：192/192 落在 384×256 內、偏移掃描前十名全是 `dy=0`）；`+1` 與 `+26` 是同一個值的兩份，疑似擁有勢力 |
| **⚠ 武將是 127 人不是 146** | 維基寫 146，一手資料每劇本 127（第 127 筆全零）。沒查證前不要引用 146 |
| **引擎畫得出遊戲畫面** | `cmd/wlview -world`。384×256 大地圖可捲動、頂端接上原版橫幅、四季切換有效。**與 PC-98 實機截圖視覺一致**：城池、河流、道路、橋、樹林、農田、山地全部正確 |
| **四季機制實地驗證** | 春→冬只換調色盤色號 14，**21 萬像素改變**而樹林/河流/道路/城池保持原色。`docs/formats/02` §4 記錄的機制第一次看到實際效果 |
| **補驗紀錄開張** | `docs/playtest/03`。專記「驗證」而非「發現」，**含失敗的驗證嘗試**。起因是兩輪內推翻自己兩條結論，產出速度超過驗證速度 |
| **地形類型已補驗** | 依類型為整張地圖上色 ＋ 四鄰接量化。橋（類型 1、2）**100% 鄰水零例外**（背景值 22%）→ confirmed；類型 5 的「河岸」判讀被 34.4% 打掉，改判林地。跨水構造（1、2、8、9）的值域全落在「可連接」範圍內，讓「行軍路網」從假說升為強證據 |
| **地形類型對映表** | `docs/mechanics/30-combat.md` §2。`sub_14C4C` 是 14 筆範圍查表（seg000 `0x982F`），值域已定案；地形名稱由圖塊外觀判讀，標強證據。**只有類型 1–7 產生野戰戰場**，平原（0）與 8、9 走別的分支 |
| **政略↔戰術接縫（M3 起步）** | `docs/re/05`。攻城戰用據點編號當戰場編號（**據點記錄 32 B，基址 0x840**）；**野戰依軍團所在格與下方四格的大地圖地形即時產生**（編號 ＝ 0xCE ＋ 類型 1–7）。軍團記錄 `+0x1A` 是格偏移 `y×384+x` |
| **戰場方向與旋轉** | `docs/formats/07` §2.2–2.3。戰場依進攻方向轉 180 度（`byte_1AB4F` 值域 0–0x3F，翻轉取 `3Fh−值`），且**旋轉時圖塊值要依三段規則換成鏡射版**。少做這一步，一半的戰鬥地形方向會錯 |
| **戰場格數** | 每個戰場 4,096 B ＝ 64 B 表頭 ＋ **3,968 格（64×62）** ＋ 64 B 尾段。範圍由旋轉常式的掃描邊界確定 |
| **`BATTLE.*` 分段結構** | `docs/formats/07`。`BATTLE.MAP` ＝ 512 B 索引 ＋ **214 個戰場** × 4,096 B；`BATTLE.MDL` ＝ 4,096 B 表頭 ＋ 3 組 × 63,488 B（合計 194,560 ＝ 檔長）。緩衝區配置四個數字同時對上 |
| **確認 `BATTLE.*` 不壓縮** | 只有 `MMAP.MAP` 走 RLE 那一支載入器。特地查過才寫——「都是地圖資料應該都壓縮」正是最容易犯的跨檔案外推 |
| **道路連接關係表**（訂正） | `docs/formats/05` §3。初版讀成「autotiling 換圖塊」是錯的——`0x80`–`0x83` 是**方向碼**不是圖塊編號。實際是掃四個方向建 768 筆 8-byte 記錄，並把道路格蓋成節點流水號。**強證據指向行軍用的路網，未定案** |
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
| `docs/re/04-mmap-entry-points.md` | 大地圖的入口點與記憶體佈局 |
| `docs/re/05-battle-selection.md` | **政略↔戰術的接縫**：戰場怎麼被選出來 |
| `docs/mechanics/00-index.md` | 機制索引與推論等級定義 |
| `docs/mechanics/30-combat.md` | 戰場：進場規則、地形類型對映、戰場形狀 |
| `docs/mechanics/70-ai.md` | **電腦 AI 的判斷邏輯**（蒐集中） |
| `docs/reference/01-jp-manual.md` | 日文原版說明書判讀紀錄（逐頁累加） |
| `docs/reference/02-jp-cht-diff.md` | 日中對照：15 則譯文缺陷 |
| `docs/reference/03-baked-japanese.md` | 燒進美術裡的日文（松崗沒重繪的部分） |
| `dosbox/` | 兩版的 DOSBox-X 設定 ＋ 出處 |
| `docs/playtest/01-dosbox-dosv.md` | DOS/V 版首次實跑：字型結案、防拷發現 |
| `docs/playtest/02-dosboxx-pc98.md` | PC-98 實跑：oracle 建立、合併抽檔陷阱 |
| `docs/playtest/03-verification-log.md` | **補驗紀錄**：撐得住的、撐不住的、失敗的驗證嘗試 |
| `docs/playtest/04-mouse-automation-blocked.md` | **⛔ 受阻**：動不了遊戲內游標，五種組合的紀錄與下一輪候選 |
| `docs/formats/04-map-sch-container.md` | `.MAP`/`.SCH` 容器格式 |
| `docs/formats/05-mmap-worldmap.md` | 大地圖：圖塊、384×256、自動連接 |
| `docs/formats/06-mmap-rle.md` | **READY** — `MMAP.MAP` 的 RLE 壓縮 |
| `docs/formats/07-battle.md` | 戰場：分段結構已解，像素格式未解 |
| `docs/formats/08-sinario-save.md` | **劇本／存檔**：4 劇本 × (192 據點 ＋ 127 武將)，武將三圍已定案 |
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

| 曾經寫過 | 實際 | 為什麼會錯 |
|---|---|---|
| `SINARIO.DAT` `+0x00` 是「恆為 1，未解」、`+0x03` 是「日」 | **`+0x00` 是日，`+0x03` 是「時」** | 兩欄在四個劇本裡都是 `1`（都從 1 日開始，且「時」初值就是 1），**光看資料分不出來**。要靠反組譯才分得開 |

> 每推翻一條就記一列，寫清楚誰推翻誰、憑什麼。**空著不是好事，是還沒開始驗。**

| 日期 | 原斷言 | 推翻依據 |
|---|---|---|
| 2026-08-07 | 「音源未解，候選 AdLib／SB／PC speaker」 | 日文說明書第 3.5 節明寫「FM音源とSSG音源のバランス」→ PC-98 側是 OPN。**但只推翻 PC-98 那一半，DOS/V 側仍未解** |
| 2026-08-07 | 「字型來源三個候選都沒驗」 | PC-98 版只用 843 B 的 `YNFONT.COM` ＋ 1,216 B 的 `FONTGRF.DAT` → 靠字型 ROM；DOS/V 的 60,888 B `YNFONT.EXE` 是在補這一塊 |
| 2026-08-07 | 「`TALK.DAT` 的訊息中間夾**控制位元組**做變數插入」 | 全解之後確認：訊息內除了 `0x00` 沒有任何控制位元組，變數是 ASCII `\1`–`\7`。**當初拿「Big5 掃出一堆短片段」反推機制，那是間接證據**——結論碰巧對，推導是錯的 |
| 2026-08-07 | 「`dosv` #225 有錯字」（`這哪<f9>堿O什麼好策略？`） | `F9D8` 是 Big5 罕用字區的「裏」，原文完全正確。**錯的是解碼器用了 `big5` 而不是 `cp950`**。差一點把工具的缺陷寫成原版的缺陷 |
| 2026-08-07 | 「`TALK.DAT` 偏移表 2,048 ÷ 4 ＝ 512 筆」 | 實際是 uint16、1024 筆。假說當時就標了「待驗」，驗完是錯的 |
| 2026-08-07 | 「地形類型 5（`0xB1`–`0xB3`）是河岸／淺灘」 | 量化推翻：只有 **34.4%** 鄰接水域，僅略高於 22% 的背景值；放大看是**綠色樹叢**。**原判讀是在 16×16 縮圖上做的——縮圖不足以判讀圖塊，要放大** |
| 2026-08-07 | 「大地圖有 autotiling：道路／河流依鄰格換成連接圖塊 `0x80`／`0x81`」 | 讀完 `sub_1E68C` 推翻：`0x80`–`0x83` 是**方向碼**，那段寫的是 8-byte 記錄不是圖塊值。**只看 `mov al, 80h` 就往「換圖塊」推，是拿指令片段當語意的典型錯誤。** 而且初版只讀到左右兩個方向，上下（`±180h`）漏了 |

---

## 7. Worklist（狀態的單一真相來源）

### 7.0 下一步（2026-08-08 基準）

**M5 的閘已解除**：`docs/mechanics/15-realtime.md` 標 READY
（整條時間鏈在機器碼裡讀出來了，`docs/re/06`），
規則層第一塊 `internal/rules/clock` 已實作 ＋ 10 個測試全過。

**月結已經整條解完。** 尾段八支只剩 `sub_12BD9` 沒讀
（對 22 個勢力各配一塊 0x30 緩衝區，疑似 AI 局勢評估或縮小地圖）。

**下一項：把三個規則層套件接成一個可跑的世界迴圈。**

理由：`clock`／`economy`／`general` 三塊都有測試護著，但**還沒有任何東西
把它們接起來跑**。做一個 headless 的世界模擬（載入 `SINARIO.DAT` →
跑 N 年 → 印出各勢力的軌跡），可以：

1. 用長期行為驗證公式（例如高稅十年後是否真的崩潰）
2. 建立 M5 需要的狀態容器（勢力／據點／武將三張表）
3. 之後接戰鬥與 AI 時有現成的驅動器

**不需要 oracle，也不需要畫面。**

> **oracle 卡住的時候，先問一次「這件事真的需要 oracle 嗎」。**
> `15-realtime.md` 被記成「擋在 oracle 後面」擋了好幾輪，
> 實際上它要的五項答案全在說明書與機器碼裡。
> 同一個教訓的另一面見 §7.0 舊版那條「機器碼有沒有直接寫」。

### 7.1 M0 剩餘

- [x] ~~即時制的可重現性~~ — `core=normal` ＋ 固定 `cycles`，已解
- [x] ~~np2 環境~~ — 不需要，DOSBox-X `machine=pc98` 兩版通吃
- [x] ~~字型來源結案~~ — 不用自備，`YNFONT.EXE` 自帶
- [x] ~~防拷檢查~~ — DOS/V 有、PC-98 沒有；oracle 改跑 PC-98
- [x] ~~弄一個 DOSBox-X~~ — `docker/dosboxx/`，PC-98 版已跑起來
- [x] ~~實測即時制的可重現性~~ — **通過**：同一串操作跑兩次，0 個不同像素
- [x] ~~推進到大地圖~~ — 已到 NEW GAME 對話框，地形畫面拿到了
- [x] ~~日文說明書 38 頁判讀~~ — **有實質機制的都讀完了**
      （剩 p.6 啟動操作、p.36–38 疑似附錄）
- [ ] **⛔ oracle 的滑鼠輸入無效**（`docs/playtest/04`）——
      根因已定位（事件從沒進 guest，但鍵盤有效）。
      **不再擋 `15-realtime`**；仍擋據點座標驗證、地形命名、戰場圖塊

### 7.2 M1 起手順序（依投報排）

> **前三項走通的路徑**（`docs/re/03` §5）：grep 檔名偏移的立即值 → 載入器讀出
> 「每筆多少 byte、怎麼定位」→ 繪製呼叫端的 `ax` 拆成寬高 → 渲染，餘數要是 0。
> **`MMAP.*` 這一組不適用**——資料先展開再畫，見 `docs/re/04`。

1. ~~`.BRG` 調色盤~~ ✅ **完成**，`docs/formats/02` READY
2. ~~`*GRF.DAT` 圖庫~~ ✅ **完成**（`ICONGRF` 段 1／段 3 待補）
2.5 ~~`MMAP.MDL` ＋ `MMAP.MAP`~~ ✅ **完成**（`MMAP.MCH` 待解）
3. `MMAP.*` 大地圖三件（兩版全同）
4. `BATTLE.*` 戰場（兩版全同）
5. ~~`SINARIO.DAT` ／ `SAVE.DAT` 的整體結構~~ ✅ **完成**
   （三段式佈局、時鐘 8 B、勢力表 22 × 64 B、據點表、武將表）

### 7.3 規則層（M5，閘已解除）

1. ~~`internal/rules/clock`~~ ✅ **完成**（五層時鐘，10 個測試）
2. ~~`internal/rules/economy`~~ ✅ **完成**（月結五步 ＋ 13 個測試）。
   `sub_15358` 尾段的批次呼叫還沒讀（生產力／上昇值／防災值／災害在那裡）
3. `internal/rules/diplomacy` — 條件式已知（`docs/mechanics/50` §8），
   **是布林判定不是擲骰**，數值門檻待反組譯
4. `internal/rules/persuasion` — 好戰等級已解（勢力記錄 `+0x28`），
   缺 `f(好戰等級)` 的修正量與信賴度增減幅度

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
