# AGENTS.md —《臥龍傳》remake 代理作業契約

本檔給 Codex、Claude 與其他自動化代理，只寫**代理專屬的**接手順序、外部知識引用與
執行環境邊界。

**專案目標、oracle 優先序、逆向規則、機制文件化要求、里程碑、目錄結構與硬規則
一律以 [`CLAUDE.md`](CLAUDE.md) 為準**，本檔不重複，也不得與它衝突。
面向使用者與新增文件預設繁體中文；程式識別字、命令、API、工具、產品名稱與檔名保留原文。

## 1. 接手順序

每次新 session、對話壓縮或工作交接時依序做：

1. 讀最新使用者需求。
2. 讀 `CONTEXT.md` 的現況與 §7.0 worklist。**不要把歷史段落當成目前待辦。**
3. 讀 `CLAUDE.md` 的目標、硬規則與證據契約。
4. 先看 [`docs/INDEX.md`](docs/INDEX.md) 的斷言總表，再讀任務直接相關的
   `docs/re/`、`docs/formats/`、`docs/mechanics/`、`docs/playtest/`。
   要碰反組譯時另外兩張表也要查：[`docs/re/21`](docs/re/21-function-census.md)
   的覆蓋地圖（這支函式有人讀過嗎）與
   [`docs/re/24`](docs/re/24-unread-function-catalogue.md)
   的未讀目錄（這支大概在做什麼、值不值得現在讀）。
   要了解某個子系統整體，入口是 [`docs/re/00-index.md`](docs/re/00-index.md)。
5. 執行 `git status --short`。**既有改動屬於使用者或前一輪工作，不得 reset、覆蓋或丟棄。**
6. 讀 [`WORKLIST.md`](WORKLIST.md) 只當快速入口；具體狀態仍以 `CONTEXT.md`、
   `docs/INDEX.md`、目前程式與可重現測試為準。

「動手之前查表」是 `[HARD]`，時機是**動手之前**不是下結論之前——
「還沒解」與「我不記得解過」在動手那一刻長得一模一樣。

**新增文件與規則時，不准寫進不存在的檔案、目錄、工具或 IDA 符號。**
`tools/phantom_scan.py` 會擋（已接進 `tools/check.sh`），但寫之前先 `ls` 更便宜。

## 2. 外部知識引用

- 命中復古遊戲、逆向、原版對拍、推廣片、音訊、跨平台或 Android 工作時，
  先讀 `~/.codex/knowledge-base/knowledge-router.md` 的「復古遊戲」路由，
  再只讀一份命中的任務入口與它直接指定的必要 reference；**不可預讀整個知識庫**。
- 專案 `AGENTS.md`、`CLAUDE.md`、`CONTEXT.md`、`docs/INDEX.md` 與可重現原版證據優先。
  外部知識只提供方法，**不能推翻本專案已定案的 DOS/V 座標、資產、流程或權利邊界**。
- `~/.claude` 一律唯讀。需要 Claude Code 的復古遊戲技巧時，使用
  `~/.codex/knowledge-base/sources/claude/` 的受控快照，清單在
  `~/.codex/knowledge-base/SYNC-MANIFEST.md`（在 kb 根目錄，不在 `sources/claude/` 裡）；
  不得從本專案或 `~/.codex` 回寫、同步覆蓋或污染 `~/.claude`。
- 推廣片工作固定命中
  `~/.codex/knowledge-base/sources/claude/retro-cht/game-promo-video-ffmpeg.md`；
  原版 AdLib／素材來源再加讀 `~/.claude/rulebook/93-promo-video-original-assets.md`；
  戰術畫面則加讀 remake 驗證 reference。
  這是階層式引用，不代表要把上述知識全部轉成技能。

## 3. 版本代號

`dosv` ＝ 1995 松崗 DOS/V 繁中版；`pc98` ＝ 1994 PC-9801 日文原版。
**一律用這兩個代號**，不要把「中文版」當成平台名稱。
`workplace/orig/` 的目錄分層就是為了讓路徑本身標明素材版本。

## 4. Docker-only 與資源衛生

- 分析、搜尋大量資料、轉檔、建置、測試、抓圖、執行程式、IDA、Xvfb、DOSBox、SDL、
  音訊與 GUI 自動化**全部在 Docker 內執行**。主機控制面只做必要的 `docker`、`git`、
  工作樹檢查與檔案編輯。**不在主機直接執行 Python、Go、專案程式、測試或 GUI**，
  文件索引與資產掃描也一樣；包裝器若違反這點，先修工具鏈。
- 一次性工作用 `docker run --rm`，帶相稱的 `--memory`、`--cpus`、`--pids-limit`、
  `--network none`（明確需要網路才開放），並以目前 UID/GID 執行可寫容器。
- **優先沿用既有包裝器**：`tools/go.sh`、`tools/py.sh`、`tools/ida.sh`、
  `tools/dosbox.sh`、`tools/dosboxx.sh`、`tools/shot.sh`、`wolong-dosboxx:latest`。
  **不可因一次失敗另造重複工具鏈。**
- 原始素材唯讀掛載；輸出只掛到工作樹的明確目錄。寫既有檔案前檢查 UID/GID，
  寫後抽查擁有權。每批工作後只清理**本輪建立且已無用途**的容器。
  **禁止** `docker system prune`、`docker image prune`、`docker volume prune`、
  `docker builder prune`、`docker rmi`、`docker container prune`。

## 5. 提交前的閘

`tools/check.sh` 是單一入口，它會跑 `go vet`、`go test`、`index.py generate`、
`phantom_scan.py`、`denylist.py` 與 `talkdat_selftest.py`。
另外自己做 `git diff --check` 與 dirty-tree 檢查。

**快取會把「跑不起來」偽裝成「通過」**：`go test` 對沒重編的套件直接回 `(cached)`，
所以需要顯示器的 Ebiten 測試在無 X 的環境下看起來是綠的。`tools/go.sh` 已內建 Xvfb，
但驗收環境的完整性仍要靠 `tools/go.sh clean -testcache` 後冷跑一次。

**測試綠不等於原版 parity。** 完成宣告還要有原版／reference 對照與正常玩家路徑證據。

- **除非使用者明確要求，不自行 commit 或 push。**
- 若有其他代理並行寫檔，**禁止 `git add -A` 或 `git add .`**，提交時只列明確檔案。

## 6. 派工邊界（主迴圈自己要寫進 prompt）

- 只能清理自己建立的 container；**禁止**任何 docker image／volume／system prune 或 rmi。
- 明列不准改的目錄：其他 repo、`~/.claude/`、`workplace/orig/`。
- 明列不准做的收尾動作：commit、push、重編、清理。**沒寫的等於允許**——
  agent 會為了「做得完整」而自行加碼。
- 收到「我順便做了 X」時，**X 一律當事故處理**：先查影響範圍再談結果。

## 7. 持續記憶與交接文件

| 檔案 | 角色 |
|---|---|
| `CONTEXT.md` | **專案狀態的單一真相來源** |
| `docs/INDEX.md` | 由 `tools/index.py generate` 產生的文件與斷言索引，**不手改** |
| `WORKLIST.md` | 交接、剩餘工作、命令閘與容器狀態的唯一入口；不另建 `HANDOFF.md` |
| `RESEARCH-LOG.md`、`REMAKE-PLAN.md` | 逆向研究的證據台帳；架構、法務邊界與刻意的 remake 差異 |

深層證據只放 `docs/` 對應文件，上表這幾份不複製證據。
