# 工作歷程

較早工作紀錄保存在既有 [`WORKLIST.md`](WORKLIST.md) 與
[`RESEARCH-LOG.md`](RESEARCH-LOG.md)。目前狀態回查 [`CONTEXT.md`](CONTEXT.md)，
研究文件索引回查 [`docs/INDEX.md`](docs/INDEX.md)。

## 2026-10-08：比對式反編譯初步試點

- 使用者採用工具鏈辨識與三函式試點；Go 基線 `e8aeb99`，遊戲規則沒有修改。
- 核對既有 IDA archive 的 SHA-256，恢復 `ida-pro-9.4-idapython:py312-v1`，image ID 與私有 manifest 相符。沒有新建同功能映像。
- `tools/ida.sh` 新增 `probe` 入口。原始 EXE 唯讀、一次性 DB、非 root、限 CPU／記憶體／PID 與外層逾時；本次沒有使用舊 root 分支。
- 研究輸出使用既有 `workplace/` 下的 `matching-decompilation/`，沒有另建根層研究目錄。
- `tools/matching_decomp.sh` 重跑 IDA 與編譯試驗；三個函式的組語重組共 233 bytes 相同，沒有遮掉呼叫或重定位運算元。
- 每個函式的位元組突變、時鐘的錯誤呼叫目標與錯誤輸入雜湊均被拒絕。Borland 標記掃描有 `LOGO.EXE` 正對照。
- GCC 12.2 `-m16` 的十二組 C 試驗沒有匹配；原版編譯器、組譯器、連結器與語言維持未知，整檔重建與 Go 行為對照沒有完成聲明。
- 證據、輸入／輸出雜湊與重跑入口在 [`docs/re/89`](docs/re/89-matching-decompilation-pilot.md)。
- 新匯出自動合併 `docs/re/matching-semantic-index.json`，保留原始函式名、位址與運算元。未證實或未收錄的語意帶警示；索引本身綁定原始 EXE SHA-256。
- 新研究文件的三筆分流以 `evidence-only` 回填到既有台帳，原有 799 筆內容、分類與 Issue 對應保持不變。
- 完整重跑包含冷啟動 IDA、三函式組語匹配、十二組 C 試驗與正負對照。Shell 與 Python 語法檢查通過。
- 文件／資產門禁逐項執行。文件索引有兩張歷史圖片缺失，幽靈引用檢查有歷史研究檔缺失，教訓的 render／verify 有既有資料不同步與防線文字缺失。用 Docker 內的 `HEAD` 快照確認這四項同樣報錯，沒有以放寬門檻處理。
- `worklist` 與 `lessons` 自測的第一次執行被唯讀 `workplace/` 掛載擋住，分類為驗證環境問題。同一 Python image、同一命令改用容器 tmpfs 後通過。原始失敗與後續檢查 log 保存在 `workplace/matching-decompilation/checks/`。
- 分流驗證與嚴格研究索引生成通過，共 802 列。文件／資產檢查最終 19／23 通過，另外四項是上述已由 `HEAD` 重現的既有問題；沒有宣稱全專案門禁全綠。
- 研究輸出抽查均為 UID/GID `1000:1000`。本輪容器使用 `--rm` 並已退出；保留恢復的 IDA image 作目前研究工具鏈，其他專案容器與映像沒有清理。
- 尚未 commit 或 push。

## 2026-10-08：採用組語基準，完成完整 KI.EXE 重建

- 使用者決定先以組語建立基準，後續還原 C 函式。新增 `tools/assembly_rebuild.sh` 與組譯工具，沿用先行試點的 Docker images。
- IDA 9.4 匯出全 segments，共 24,376 條指令、55,392 bytes；113 條指令有 loader 重定位差異，按原始 MZ 表還原段值。
- 所有 IDA 指令以助憶碼組譯回相同 bytes。任一指令不匹配即停止，沒有程式碼 bytes fallback。MZ 標頭、116 筆重定位與其餘 11,195 bytes 的資料宣告一起重建。
- 整份 67,099 bytes 的輸出與原始 KI.EXE SHA-256 相同。只掛載組語與 linker script 的獨立容器冷重建也相同。指令、立即數常數與 MZ 重定位三組負對照均拒絕。
- 來源對映完整覆蓋檔案，函式台帳保存 739 筆 IDA 定位與分級語意；沒有宣稱 C 還原、原作者來源切分或 Go 行為對拍完成。
- 全量匯出初次讀取中文 JSON 時缺少 UTF-8，修正後以同一入口、同一 image 從乾淨輸入重跑。分類 item 尾 byte 時改查 item head，並以完整指令 bytes 總數作正對照。
- 現況與收據：[`docs/re/90`](docs/re/90-assembly-reconstruction.md)。先行試點及其 DB／收據保留，沒有覆寫歷史 checkpoint。
- 新工具 Shell／Python 語法、803 列分流、嚴格研究索引、過期斷言與資產 deny-list 均通過。文件索引與幽靈引用仍回報先行試點已由 HEAD 複現的歷史缺檔，本輪新文件與工具沒有新增引用錯誤。
- 原始 EXE、完整重建與來源獨立重建的雜湊一致；來源／DB／收據與文件交叉核對通過。研究輸出與本輪修改檔案擁有權抽查為 UID/GID `1000:1000`。
- 本輪容器已退出並自動移除。只清理由本輪 IDA wrapper 留下的個別 `.ida-*` 暫存目錄，其他研究輸入、既有 root-owned 檔案與其他專案資源保持原狀。沒有 commit 或 push。

## 2026-10-08：第一個 C 函式 sub_1ECE0

- 先以固定雜湊的 IDA 指令建立 spec/201，審查取數、寫回順序、AX、FLAGS、PUSH／POP／RET 及堆疊重疊契約，READY 後寫研究用 C。
- `tools/c_recovery/rng.c` 提供可讀的 typed state 函式與 16-bit 呼叫介面適配層。正式 Go 規則沒有修改，也沒有在原版資料上寫入。
- dosgolem 載入原始 EXE，保留真實 CS 的 IDAIn 入口，每次執行原始 13 條指令。四張固定表各枚舉 65,536 種 c/s，另測堆疊邊界與重疊。
- GCC O0/O2 各 263,680 組原版/C 對照與 263,168 組 Go 狀態對照相同，每版 258 次完整 1 MB 記憶體核對。計數器增量、AH、旗標來源與 BX 取回方式四種突變皆被拒絕；BX 突變在 alias 案例才暴露。
- spec/201 更新 CONFORMED，範圍限受控局部常式與 IF/TF=0。C 機器碼匹配、播種、呼叫時機、GUI 與正常玩家路徑沒有新的完成聲明。
- 對照器的 C11 全域旗標使 cgo 的 POSIX prototype 不可見，加入 GNU feature macro 後同一工具鏈重跑通過。gofmt 的單檔掛載無法建立相鄰暫存檔，改用容器 tmp 產生內容後回寫。
- 全量枚舉第一次在原版計時器中斷送達時超出 20 條局部預算。根因是入口 IF=1；按已定的局部隔離契約固定 IF/TF=0，保留原始 13 條判準，沒有放寬預算或忽略差異。
- 收據逐項驗證 source hashes、O0/O2 摘要與負對照，dosgolem 被 import 的 oracle/internal 原始碼在執行前後雜湊相同。其既有 cmd/probe/main.go 修改沒有被本次改寫。
- 現況在 docs/re/91 與 c-recovery-status.json；組語基準台帳保留建立時的快照，避免改寫先行收據。
- C 完成索引的來源雜湊、807 列分流、嚴格研究索引、過期斷言與資產 deny-list 通過；README 的規格數已包含新增的研究契約。文件索引與幽靈引用仍有已複現的歷史缺檔，新 C 文件與工具沒有新增引用錯誤。
- 收尾核對原版 EXE 雜湊、研究輸出擁有權與容器清理；本輪輸出均由 UID/GID 1000:1000 擁有，沒有 commit 或 push。

## 2026-10-08：持續 Goal 與每輪提交

- 使用者設定持續目標「完成臥龍傳 matching decompilation」，並明確授權每輪完成後 commit、push。
- 前一輪分類為 progress：整檔組語來源、第一個 C 函式、原版/C/Go 局部收據已形成可回查交付物。
- 本輪先核對 main 與 origin/main 沒有差異，再提交這批已驗證研究來源與文件；本機原版、DB、整檔組語、重建 EXE、含原版資料的研究產物保持 Git 忽略。
- 完整 Goal 維持 active。現有成果只證明組語基準及 sub_1ECE0 的局部 C 語意，不當成整個 C 還原與 remake 比較完成。
- 基準與第一個 C 函式已 commit 為 `1d21147`，push 成功，origin/main 指向同一提交。

## 2026-10-08：第二個 C 函式 sub_1EC82

- 先核對完整 IDA inventory 的 94 bytes 與 RTC 呼叫流程，建立 spec/202。原版讀 BIOS 回覆的 AL，不能默默假設為零；C 介面保留 RTC_AL。
- dosgolem 的現有 RTC 成功回覆不在支援範圍，以明示的 RAM fixture ISR 提供 CX、DX、AL，保留原始 INT／IRET 與原始 94 bytes。來源與固定回覆隨收據記錄。
- 全部 86,400 個合法時分秒、12 個非零 AL 案例與 8 個堆疊邊界，O0/O2 各 86,420 組完整原版/C 狀態相同；AL=0 的 86,408 組 Go 播種狀態相同。每版 85 次完整記憶體核對相同。
- 每案原版固定 3,369 條指令，沒有挑選時間或放寬停點。索引步長、小時倍率、XOR 來源三個 C 突變全部被拒絕。
- XOR 的 AF 本輪按 dosgolem 模型為零，明示硬體未定義；合法 BCD、IF/TF=0、分離堆疊與固定 RTC 是證據邊界，未外推實機時序或正常玩家流程。
- spec/202 更新 CONFORMED，C 狀態索引加入第二個函式；原版取數的來源與先前收據保持原樣。Goal 未達整個 C 還原完成。
- 完整收據、來源雜湊與重跑入口在 docs/re/92；所有 workload 仍在 Docker 內，原始輸入唯讀。
- 新工具語法、815 列分流、嚴格研究索引、過期斷言及資產 deny-list 通過。文件索引與幽靈引用的既有歷史缺檔保留，本輪新文件沒有新增引用錯誤。
- 本輪完成後依授權 commit、push；Goal 保持 active，兩個 C 函式不代表全 executable 的 C 還原與 remake 比較已完成。

## 2026-10-08：第三個 C 函式 sub_11D8E 與 Go 年份訂正

- 先核對固定 IDA inventory 與原始 137 bytes，建立 spec/203 的時鐘契約。原版年份比較是 1000，999 會再加到 1000；五組原版 probe 揭露舊 Go 的固定 999 上限差異。
- spec/204 審查 READY 後，修正 Go 換年，新增連續 999→1000→999 與五種入口年份的冷測。機制現況更新在 docs/mechanics/15，舊結論的訂正追加在 CONTEXT.md §6。
- 原版/C O0/O2 各 292,297 組相同：226,665 組合法日期、65,536 個年份入口、96 組 callee／等待案例。Go 排除 MOV callee 後 292,249 組相同，各版 286 次完整 1 MB 記憶體抽查相同。
- 四個 callee 的入口快照比較全部 14 個暫存器／段／FLAGS、日期與堆疊；等待比較輪詢次數、FLAGS 與清零欄位。原始受測 137 bytes 不變，RET／MOV 與 ready／count fixture 是明示的輸入控制。
- 日進位門檻、年份比較、呼叫順序與等待比較四種突變全部被拒絕。O0/O2 收據相同，來源、EXE、原始函式 bytes 與 dosgolem import 來源前後身分核對通過。
- 最初年份 probe 放在 /output，遭 Go internal import 邊界拒絕；移到工作樹的研究區並補正明確輸出掛載後，以同一原版與條件重跑。這是驗證環境問題。
- spec/203、spec/204 更新 CONFORMED，研究入口為 docs/re/93；現行 C 索引已有三個函式。證據限局部時鐘，完整 callee、硬體時序、1000 年 UI 與長期玩家流程未驗證，沒有 C 機器碼匹配完成聲明。
- 所有 workload 在非 root Docker 內執行，原版與外部 dosgolem 來源唯讀；本輪依既有授權 commit、push，完整 Goal 保持 active。
- 文件與工具檢查 21／25 通過。825 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list、新工具語法及所有正對照通過。四項失敗沿用已在 HEAD 重現的歷史缺檔、教訓資料不同步與既有防線缺失；新文件與工具沒有新增引用錯誤。
- 三個 C 狀態索引的來源身分全部通過。時鐘產物 UID/GID 為 1000:1000，原版 EXE 雜湊相同，專案既有 root-owned 路徑仍為 27 筆，沒有擴大修復範圍。
