# AGENTS.md —《臥龍傳》remake 交接與作業規範

本檔是本儲存庫給 Codex、Claude 與其他自動化代理的快速恢復契約。它補充
[`CLAUDE.md`](CLAUDE.md) 與 [`CONTEXT.md`](CONTEXT.md)，不取代反組譯證據、格式規格或
`docs/INDEX.md`。面向使用者與新增文件預設使用繁體中文；程式識別字、命令、API、工具、
產品名稱與檔名保留原文。

## 接手順序

每次新 session、對話壓縮或工作交接時依序做：

1. 讀最新使用者需求。
2. 讀 `CONTEXT.md` 的現況與 §7.0 worklist；不要把歷史段落當成目前待辦。
3. 讀 `CLAUDE.md` 的目標、硬規則與證據契約。
4. 先看 [`docs/INDEX.md`](docs/INDEX.md) 的斷言總表，再讀任務直接相關的
   `docs/re/`、`docs/formats/`、`docs/mechanics/`、`docs/playtest/`。
5. 執行 `git status --short`。既有改動屬於使用者或前一輪工作，不得 reset、覆蓋或丟棄。
6. 讀 [`MEMORY.md`](MEMORY.md) 與 [`WORKLIST.md`](WORKLIST.md)，只把它們當快速入口；
   具體狀態仍以 `CONTEXT.md`、`docs/INDEX.md`、目前程式與可重現測試為準。

## 外部知識引用

- 本專案命中復古遊戲、逆向、原版對拍、推廣片、音訊、跨平台或 Android 工作時，
  先讀 `~/.codex/knowledge-base/knowledge-router.md` 的「復古遊戲」路由，再只讀
  一份命中的任務入口與它直接指定的必要 reference；不可預讀整個知識庫。
- 專案 `AGENTS.md`、`CONTEXT.md`、`docs/INDEX.md` 與可重現原版證據優先。外部知識
  只提供方法，不能推翻本專案已定案的 DOS/V 座標、資產、流程或權利邊界。
- `~/.claude` 一律唯讀。需要 Claude Code 的復古遊戲技巧時，使用
  `~/.codex/knowledge-base/sources/claude/` 的受控快照及其 `SYNC-MANIFEST.md`；
  不得從本專案或 `~/.codex` 回寫、同步覆蓋或污染 `~/.claude`。
- 推廣片工作固定命中 `game-promo-video-ffmpeg.md`；原版 AdLib／素材來源再加讀
  `rulebook/93-promo-video-original-assets.md`；戰術畫面則加讀 remake 驗證 reference。
  這是階層式引用，不代表要把上述知識全部轉成技能。

## 專案目標與邊界

- 以 NEO･GETEN《臥龍傳－三國制霸之計》的 `dosv` 與 `pc98` 原版為行為 oracle，
  在 Go／Ebiten 上乾淨重寫，保存松崗繁中母本並以日文版逐句校訂。
- `dosv`＝1995 松崗 DOS/V 繁中版；`pc98`＝1994 PC-9801 日文原版。
  使用這兩個代號，不把「中文版」當成模糊平台名稱。
- 核心規則對齊原版；高解析、縮放、現代輸入、固定時間基準或其他外殼變更都要明列為
  remake 差異，不能默默改玩法。
- 原版執行檔、資料、圖像、音樂、倚天字型與原始存檔不進版控、不進發行包。
  `workplace/orig/` 與原始研究輸入唯讀；存檔實驗使用副本或明確輸出目錄。

## 目前接手基線（2026-08-10）

- M0–M2 已完成；M3–M6 大致完成。格式解碼、五層即時時鐘、月結、內政、外交、說服、
  戰術戰鬥、沿原版道路行軍、信賴度存檔改寫與 Ebiten 呈現已有實作及歷史驗收證據。
- 主要未完成是 M7 校訂後的硬換行／行寬／formatter／畫面 parity，以及 M8 的目標平台實機／發行驗收。
- M7 的 `translations/corrections.json` 目前有 60 筆；60 筆已有
  `tools/talkdat.py correct` 與 selftest 證據，#0–#1021 已完成第一輪逐句讀取。`#751` 採保留既有中文、在「我」後補 `{1}` 的
  最小修正；插入位置由同一戰場台詞池的 `{1}` 用法提供強證據，但不宣稱完整語氣／排版已逐句等價。`#321` 已由
  格式化器（formatter）資料流定案，`#192`–`#195` 已由 `sub_13BA9`／`sub_13C99` 證實為三變體直取槽位，
  `#257/#258` 已由內容對照定案；校訂清單與裁定紀錄見 `docs/reference/02` §12。
- M8 要在目標環境實際建置、跑正常玩家路徑、做發行 deny-list、檢查存檔與原始素材隔離。
- 獨立未解項包括 `ICONGRF` 段 1、視窗內龍紋、`MMAP.MCH`／`BATTLE.MCH` 語意、
  DOS/V 音源、協力第 ④b 道閘，以及 `sub_135AB` 等 `CONTEXT.md` §7.0
  列出的項目。未知不等於可以自行補玩法。

## 證據與逆向契約

- Oracle 優先序：固定狀態原版實測 > IDA Pro 9.4 反組譯 > 說明書／當年資料 > 社群資料。
  社群資料與 2017 MOD 只能指路；MOD 數值不能當 1995 原版數值。
- 反組譯先用 `ida-pro-9.4-ver2` 與 `.i64` 的函式邊界、交叉參照、資料流；Ghidra 只交叉驗證。
  不以攤平 `.asm` 取代 IDA 關係圖，不以推測性改名覆蓋原函式名、全域名、位址或運算元。
- 每項語意都標 `confirmed`、強證據、假說或未知，並附輸入檔、SHA-256（若已記錄）、
  工具版本、位址空間與出處。直接 xref 抓不到取址後的間接讀寫；零寫入不能直接解讀成無寫入。
- 推翻舊結論時保留舊證據索引，將新證據與理由追加到 `CONTEXT.md` 的「已被推翻的斷言」，
  不重寫歷史。不因欄位名稱看起來合理就升級推論等級。
- 每解出一條機制，同步更新 `docs/re/` 與 `docs/mechanics/`；規格只有 `READY` 才能作為
  實作閘。存檔從原始 bytes 改寫已解欄位，未知 bytes 必須保持不變並做 round-trip。

## Docker-only 與資源衛生

- 分析、搜尋大量資料、轉檔、建置、測試、抓圖、執行程式、IDA、Xvfb、DOSBox、SDL、音訊與
  GUI 自動化全部在 Docker 內執行。主機控制面只做必要的 `docker`、`git`、工作樹檢查與檔案編輯。
- 一次性工作使用 `docker run --rm`，帶相稱的 `--memory`、`--cpus`、`--pids-limit`、
  `--network none`（明確需要網路才開放），並以目前 UID/GID 執行可寫容器。
- 優先沿用既有映像與包裝器：`wolong-dosboxx:latest`、`tools/go.sh`、`tools/ida.sh`、
  `tools/dosbox.sh`、`tools/dosboxx.sh`、`tools/shot.sh`。不可因一次失敗另造重複工具鏈。
- 原始素材唯讀掛載；輸出只掛到工作樹的明確目錄或 `/tmp`。寫既有檔案前檢查 UID/GID，
  寫後抽查擁有權。每批工作後檢查專案相關執行中／停止容器，只清理本輪建立且已無用途的容器。
  禁止 `docker system prune`、`docker image prune`、`docker volume prune`、`docker rmi`。
- 不在主機直接執行 Python、Go、專案程式、測試、DOSBox 或 GUI。即使是文件索引／資產掃描，
  也放到隔離容器；包裝器若違反這點，先修工具鏈或在容器內直接執行。

## 程式、文本與驗收

- 分層維持 `internal/ui` → `internal/rules`／`internal/state` → `internal/assets`；排版邏輯
  不依賴 Ebiten。規則只保留一份實作，新增前先查既有程式與文件。
- 所有玩家可見文字走語系資料；繁中以原版為母本，先對 `pc98`，不得憑語感修成另一套譯文。
- 原版 `0xFF` 等哨兵不可被 Go 零值悄悄取代；在唯一入口正規化並測試。即時制截圖必須記錄
  狀態、座標、種子、時間步、模式與已知差異。
- 信賴度不是勢力 `+0x1D`：它是每個存檔區塊 `+0x10`，對應原版 `cs:0D00h`／IDA
  `byte_10D00`；原始位址與證據等級必須和 remake 的 `World.Trust` 分開記錄。
- 提交前至少跑 `tools/go.sh vet ./...`、`tools/go.sh test ./...`、`tools/index.py generate`
  與 `tools/denylist.py --selftest`，再做資產掃描、`git diff --check`、dirty-tree 檢查。
  測試綠不等於原版 parity；完成宣告還要有原版／reference 對照與正常玩家路徑證據。
- 除非使用者明確要求，不自行 commit 或 push。若有其他代理並行寫檔，禁止 `git add -A`／`.`，
  提交時只列明確檔案。

## 持續記憶與交接文件

- `CONTEXT.md`：專案狀態唯一真相來源。
- `docs/INDEX.md`：由 `tools/index.py generate` 產生的文件與斷言索引，不手改。
- `MEMORY.md`：代理快速恢復摘要與不可重犯的陷阱。
- `WORKLIST.md`：唯一的交接、剩餘工作、命令閘與容器狀態入口；不再建立 `HANDOFF.md`。
- `RESEARCH-LOG.md`、`REMAKE-PLAN.md`、`VERIFICATION-MATRIX.md`：逆向研究、架構／切片計畫、
  驗證矩陣台帳。深層證據仍只放 `docs/` 對應文件。
