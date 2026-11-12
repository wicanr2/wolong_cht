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

## 2026-10-08：月結呼叫鏈、六個 C 函式與赤字捨位訂正

- 前一輪時鐘成果 `228f30d` 已推到 origin/main，分類為 progress。開工確認工作樹乾淨，再按逆向、IDA 及文件職責路由載入當前入口。
- 新 IDA 9.4 probe 匯出月結主流程、資金加減、距離除數、預備兵上限與赤字函式，共 306 原始 bytes。DB、原始 bytes、operand、xref 與原始名稱保留在本機 c-economy/ida。
- spec/205 READY 後實作六個 C 函式。月結的 C 資金／赤字／RNG 實際相接，原版側保留相同真實指令；未還原的據點結算與尾端 callee 用明示 MOV／RET fixture。
- 小量 probe 通過 192 組。固定 RNG 的原版赤字 probe 揭露 Go 捨位錯誤：原版先讀負資金高位再 NEG，Go 曾先取絕對值再右移；合法欠款非 256 倍數時每兵種少扣 16。
- 依 spec/206 READY 修正 Go，將舊 -16000 的 992 預期訂正為原版 1008，新增六個邊界與非零固定 RNG 測試。原始證據、機制現況及 CONTEXT.md 訂正台帳同步。
- O0/O2 各 1,775,568 組原版/C 與 1,326,175 組 Go 對照相同，各版 867 次完整 1 MB 記憶體抽查相同。高低位資金、carry、距離、赤字與月結順序六個突變全部被拒絕。
- 經濟及 internal/state 冷測通過。狀態層最初缺少 x/text v0.40.0，既有 image 的 GOPROXY=off 使首次開網路仍無法載入；明示官方 proxy、按 go.mod／go.sum 下載後，以同一 image、同一測試範圍重跑通過，分類為環境問題。
- C 初次編譯的 hex literal 與 + 貼連成 preprocessing token，修正空白；Go 對照器的 Bytes 長度由 uint32 改成正式 API 的 int。沒有放寬判準、挑樣本或忽略差異。
- 六個新函式加入現行 C 台帳，共九個；語意索引附 proven 等級及 re/94 出處。組語與先前 C 收據保留，沒有完整 C matching 或完整玩家月結的完成聲明。
- 本輪依授權完成後 commit、push，完整 Goal 保持 active。
- 收尾文件與工具檢查 21／25 通過，833 列分流、嚴格研究索引、新工具語法及所有正對照通過。四項既有失敗仍為歷史圖片／研究輸入缺檔與教訓資料不同步／防線缺失；新文件沒有新增引用錯誤。
- DB 雜湊與探針路徑原先寫在同一行，過期斷言工具誤綁為腳本雜湊；改為分列工具、DB 路徑與 DB 雜湊後，原工具重跑通過。沒有放寬掃描器。
- 經濟及狀態層 Go vet 通過，九個 C 函式的現行來源雜湊核對通過。研究產物 UID/GID 為 1000:1000，原版 EXE 雜湊相同，既有 root-owned 路徑仍為 27 筆。

## 2026-10-08：五個據點結算 C 函式與原始整數寬度

- 前輪 `4e50f80` 已推送，分類為 progress。核對乾淨工作樹、現行 C 台帳與逆向／IDA／文件職責路由，接續據點收入與募兵。
- 新 IDA 9.4 probe 保留五個原始函式的 420 bytes、operand、chunk、xref 及 DB 身分。spec/207 READY 後實作 C，沿用未改動的經濟算術／堆疊介面。
- 據點結算接入既有 C 月結，替換前輪 MOV 收入 fixture。距離、收入、募兵、玩家／AI、資金、預備兵、赤字與 RNG 在 C 實際相接；九個世界更新尾端及重畫仍是明示 RET fixture。
- 固定原版反例查出玩家 low-word carry 與募兵 16-bit 累計兩個 Go 差異。依 spec/208 READY 修正；經濟、internal/state 冷測及 Go vet 通過。
- 密集募兵測試最初用概略比率推期望，忽略分次移位的捨位。直接讀原版後改為 51584／5952／7808，四組稅率／累計反例的原版/C/Go 相同。舊差異 probe 保留，沒有挑樣本。
- 首次單獨執行探針缺少 /orig 掛載，補上原版唯讀掛載後，用同一 binary、同一案例重跑；分類為驗證環境問題。
- O2 的 1,779,300 組原版/C、1,385,882 組 Go 與 435 次完整記憶體抽查相同。整批 600 秒上限在 O0 結束前送達；已確認工作終止，保留 O2，只補跑 O0 與負對照，沒有減少案例或變更判準。
- 既有六函式收據已用新 Go 來源完整重生，1,775,568 組原版/C、1,326,175 組 Go 與六個負對照仍通過；先前收據保留在 c-economy/pre-settlement-proof。
- O0 補跑完成，與已取得的 O2 收據逐位元組相同；五組負對照皆被拒絕，來源、DB、原版 EXE／SINARIO.DAT 與收據核對通過。
- 現行 C 台帳已加入五個函式，共十四個；語意索引保留原始名稱、位址與 proven／re/95 出處，沒有改寫先前 DB 或證據。
- 文件與工具檢查 21／25 通過，842 列分流、嚴格研究索引、過期斷言、新工具語法、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗仍保留，新文件沒有新增引用錯誤。
- 本輪按既有授權 commit、push；完整 matching decompilation Goal 保持 active，十四個 C 函式與經濟接線不代表完整世界月結、玩家垂直鏈或 C 機器碼匹配完成。
- 原版 EXE／SINARIO.DAT 身分、十四個 C 函式來源雜湊與分級語意均核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出並移除。

## 2026-10-08：十一個月結世界更新 C 函式與真實 writer

- 前輪 `259df4f` 已推送，分類為 progress。開工核對工作樹、現行台帳與逆向／IDA／文件職責路由，選擇七個月結尾端及其 event writer、關係取址與暴風雨標記依賴。
- IDA 9.4 probe 保留十一個原始函式，共 1068 bytes，原始名稱、operand、chunk、xref 及 DB 身分不改寫。spec/209 READY 後實作 C，沿用原有經濟、據點及 RNG 介面。
- 原版與 C 均執行真實事件 writer 與 RNG，並比較 queue code 的低 byte 判空、指定／隨機起點、滿槽、向後搜尋及 CF。月底的官員、外交、災害與赤字通知不再用事件清單或 RET 替代。
- O0/O2 各 142,864 組原版/C 相同，各版 1117 次完整 1 MB 記憶體抽查相同；四劇本接線與七個刻意錯誤版本均通過原定拒絕閘門。
- 固定原版生產力 probe 揭露 Go 普通整數加法及乘積符號的寬度差異：65000／上昇值100／稅率30 的原版為 12114。依 spec/210 READY 修正後，三個邊界冷測與 560 個原版/C/Go 向量相同，最終 RNG 相同。
- 經濟、internal/state 冷測與 Go vet 通過。新矩陣最初缺少八個群組的外層括號，依語法定位補回後在同一工具鏈重跑，沒有改案例範圍或判準。
- 十一函式加入現行 C 台帳，共二十五個；每筆保留原始位址、C source hash、routine hash、proven 及 re/96 出處。
- 先前經濟及據點收據已保留在各工作根的 pre-world-proof，現行收據按新 Go 檔案身分重生。原始 EXE／DAT 一直唯讀。
- 剩餘 sub_1585F、sub_12BD9 與 UI／音效仍為明示 fixture，完整玩家月結及 C 機器碼尚未驗證。Goal 保持 active，本輪依授權完成後 commit、push。
- 收尾文件與工具檢查 21／25 通過，852 列分流、嚴格研究索引、過期斷言、新工具語法、校訂、資產 deny-list 與所有正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件沒有新增引用錯誤。
- 二十五個 C 函式來源身分、分級語意、原版 EXE／DAT 雜湊與輸出擁有權均核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出並移除。

## 2026-10-08：完整 C 月結規則與政治／俘虜依賴

- 前輪 `5402deb` 已推送，分類為 progress。中斷前只讀取依賴，恢復後確認工作樹乾淨、沒有存活容器，再沿逆向／IDA／文件職責契約續作。
- 二十一個原始函式的 1508 bytes、operand、chunk、xref 與固定 IDA DB 匯出，形成政治／俘虜依賴閉包。spec/211 READY 後實作 C。
- 九個月結尾端全部以真實 C 接線，含政治、俘虜、queue 月度初始化、scratch 邊境去重排序與 writer。僅通知、音效、重畫仍是明示 RET fixture。
- O0/O2 各 11,632 組原版/C 相同，各版 91 次完整記憶體核對相同，含四個原版劇本 × 三個玩家選擇。新政治／俘虜原版證據沒有外推完整 Go 月結 parity。
- 八個 writer、relation war bit、queue 搬移、任官、力量、timer、排序及協力門檻突變均被拒絕。協力門檻曾通過原先 80 個政治樣本，查明沒有滿足深層前提；新增八個相鄰門檻後，正常版本通過、突變拒絕。
- 收據驗證器舊計數 11,624 在新增案例後拒絕新收據，按實際新增群組與 source identity 更新為 11,632，沒有改來源身分或降低比較判準。
- 現行 C 台帳加入二十一函式，共四十六個，每筆保留原始位址、routine／source hash 與 proven／re/97 出處。正式 Go 原始碼本輪沒有變更。
- 本輪依授權完成後 commit、push；完整 matching decompilation Goal 保持 active，月結規則接線不代表整檔 C 或正常玩家 UI 完成。
- 收尾文件與工具檢查 21／25 通過，860 列分流、嚴格研究索引、過期斷言、新工具語法、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件沒有新增引用錯誤。
- 四十六個 C 函式 source identity 與分級出處核對，原版 EXE／DAT 雜湊相同。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：完整月結原版／C／Go 36 向量比較

- 前輪 `e2fe24f` 已推送，分類為 progress。核對現況、逆向／IDA／spec／文件職責路由，將已還原 C 呼叫鏈用來直接審查 remake。
- 從正式 tick 純抽取 153 行原有月結規則至 private monthlyRules，正常 tick 的前置據點／軍團／時鐘與最後每時更新次序保留。MatchingMonthly 僅在 matching build tag 可見，沒有 production shortcut。
- 同四個原版劇本、玩家 0／7／21、raw c=0／77／255 與 s=uint8(c×37)，原版/C 每向量全暫存器及 1 MB 相同。Go 在同一 rules 邊界比較全部 22,208 bytes 與 258-byte RNG，所有載入基準零差異，沒有 mask 或剔除 queue。
- 最初 0／36 完全相同，general +0x1F 共 3429 bytes 差異。按 READY spec/213 接回 typed score、月結 byte 計算與保存，9／36 通過。
- 原始交換式排序修正依 READY spec/214，12／36 通過；中立 FF18 producer 依 spec/215，20／36 通過；storm globals 與玩家宣戰 gate 依 spec/216、217，最終 36／36 完整區塊與最終 RNG 相同。
- 各 before-score／after-score／after-order／after-neutral checkpoint 的 audit、vectors 與 source hashes 保留。本機原版區塊與 RNG binary 不進 Git。
- Storm globals 使用 typed 更新欄位，Bytes 套到原版位置，RawBlock 原始錨點保持不變。新增原版格式 round-trip、停用／終止槽、byte wrap、同值尾端、FF 邊界、中立非 producer 與玩家規則冷測。
- strategyai／state 正常與 matching tag 冷測、Go vet 通過；驗證器 globals／勢力／據點／武將／queue／RNG 六種單 byte 突變全部拒絕。
- 首次排序實作以 byte 暫存 Raw() 的 int，型別檢查拒絕，改用正式 API 型別後同一向量重跑。沒有放寬原始比較或挑 seed。
- spec/212–217 更新 CONFORMED，範圍限這 36 固定月結 rules 向量，不能當作任意後期世界、正常 GUI／音效／長程、完整 C 或 C 機器碼完成。Goal 保持 active，本輪依授權 commit、push。
- 收尾文件與工具檢查 21／25 通過，885 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件與研究入口沒有新增引用錯誤。
- 36 向量 source identity、原版 EXE／DAT 雜湊、46 個未修改 C 函式來源與產物擁有權核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：八個每時更新與 dispatcher C 核心

- 前輪 `8251289` 已推送，分類為 progress。開工核對乾淨工作樹與現況，命中逆向／IDA／文件職責路由，沿時鐘的每時下游續作。
- 原始八函式共 445 bytes，IDA 9.4 probe 保留原始名稱、位址、operand、chunk、xref 與固定 DB。spec/218 READY 後以 C 接線財政、維持費、外交與原始 event table。
- event 10／13 和 trust helper 為真實 C，其他 11 event handler 與 UI／音效／退出／重畫仍為明示 RET fixture，沒有額外 Alive gate 或更換表格。
- O0/O2 各 21,871 組原版/C 相同，各版 171 次全 1 MB 記憶體核對相同，含全部 cadence byte、13 碼及空事件、四劇本每時入口、外交與支出／兵池邊界。
- 初次 trace 的暫存器與資料已一致，根因是 event 10 位址重複登錄、sub_1310A 沒登錄。按位址去重並補正式依賴後，用同樣矩陣、原版與 source 重跑，沒有放寬 trace 判準。
- cadence 突變被原始 comparator 拒絕，但失敗 JSON 沒保存真正不同的 globals／cadence 欄位；補齊 schema 並重生收據，validator 才能機器核對差異。
- 支出上界、carry、cadence、cursor stride、trust borrow、次序與外交 byte 經費七個突變全拒絕。正式 Go 本輪沒有修改，完整 Go 每時對拍與其餘 handler 尚待後續。
- 現行 C 台帳加入八函式，共五十四個；每筆 source／routine hash、原始位址與 proven／re/99 出處保留。Goal 保持 active，本輪依授權 commit、push。
- 收尾文件與工具檢查 21／25 通過，893 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓檢查失敗保留，新文件沒有新增引用錯誤。
- 五十四個 C 函式來源身分、分級語意、原版 EXE／DAT 雜湊與產物擁有權均核對。輸出 UID/GID 為 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：十三碼事件 handler 與政治／災害 C 閉包

- 前輪 `9f5a843` 已推送，分類為 progress。開工核對乾淨工作樹、現行狀態、逆向／IDA／文件職責路由與 GitHub Issue #27。
- 三十個原始函式共 1919 bytes，IDA 9.4 probe 保留名稱、linear address、operand、完整 chunks、xref 與 DB 身分。spec/219 READY 後實作 C 原始名稱入口。
- 十一 handler、event 8 落入 `sub_133FD` 的原始尾端及十八個政治／災害依賴均真實接線。沒有為落下尾端插入 CALL／RET；modal、金額、戰鬥準備、UI 仍是明示 fixture。
- O0/O2 各 43,476 組原版/C 相同，各版 340 次完整 1 MB 核對相同；三十函式都有實際入口次數，包含四劇本原始 dispatcher 與受控 queue／四種回應。
- 門檻、關係 SHR、raw 0x24、災害容量／timer、軍團分支、原始 count 比較與協力上界八個突變全拒絕。補齊 C 原始名入口後，以最終來源完整重生收據。
- 原版與 C 的狀態摘要分別讀各自記憶體，仍相同。正式 Go 本輪沒有改動，中立 byte 行為沒有升為玩家語意，Issue #27 保持 OPEN。
- 更正 spec/218 的每時矩陣為四劇本獨立入口，原收據保留。較早 spec/218、re/99、spec/64 的範圍 backlink 由新 verifier 核對。
- 現行 C 台帳加入三十函式，共八十四個；source／routine hash、位址、實際入口及 proven／re/100 出處保留。完整 Goal 保持 active，本輪依授權 commit、push。
- 收尾檢查的 DB 雜湊少了明確檔案路徑，stale_scan 誤配到上一個 Markdown 連結；補上 DB 路徑後通過。兩個自我測試需要 workplace 暫存案例，改掛有界 tmpfs 後以同一工具鏈重跑通過，首次環境失敗 log 保留。
- 收尾腳本用 basename 去比完整路徑，漏寫 docs/re/43；改為明確重生，並正查新兩份文件都已收錄。文件與工具檢查 21／25 通過、901 列分流與嚴格索引通過，四項既有缺檔／教訓失敗保留，本輪沒有新增引用錯誤。
- 八十四個 C 函式來源身分、三十筆分級出處、原版 EXE／DAT、實際入口與 backlink 均核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：十七個外交／金額視窗 C 函式

- 前輪 `b5ad756` 已推送，分類為 progress。開工核對乾淨工作樹、現行狀態、逆向／IDA／視窗／文件職責路由與 GitHub Issue #22。
- 十七函式共 1280 bytes。第一次直接比 IDA 與 file chunks，在二十個 FAR segment word 失敗；確認全都來自 MZ loader 加 paragraph。保留舊 probe／DB，按原始 116 筆 MZ 表逐 word 記錄兩種 bytes 與 runtime paragraph，沒有遮罩差異。
- IDA 9.4 probe 審查通過、spec/220 READY 後實作 C 原始名入口。真正的 FAR／RETF、原始 SS frame、PUSHF／POPF、CF 輪詢／取消重試、RNG／trust、金額／扣款、code word 修改及恢復已接線。
- 受控裝置 FAR leaf 會 clobber FLAGS，讓 PUSHF／POPF 成為可拒絕的 gate。選單忙碌與數值取消分開，原始 caller 迴圈真實執行，底層繪圖／輸入／音效仍為明示 primitive fixture。
- O0/O2 各 53,322 組原版/C 相同，每版 417 次完整 1 MB 核對相同；十七函式皆有實際入口、252-byte callee 快照、控制欄位與中途 code word 證據，十個突變均拒絕。
- 收尾抽查查出訊息向量寫到 world bank，實際讀取端是 SS:BP。補 typed message、改寫真正的 SS frame 後，以完整矩陣與全部負對照重生，舊來源與收據保留在 c-modal 的 pre-message 檔案。
- 較早 spec/219、re/100 範圍 backlink 由 verifier 核對，另十一份同位址規格逐份確認不受影響。正式 Go 未改，原始未初始化 stack 槽不補值、raw 大數值不冒稱玩家輸入，Issue #22 保持 OPEN。
- 現行 C 台帳加入十七函式，共 101 個。完整 Goal 保持 active，本輪依授權完成後 commit、push；本輪不宣稱完整視窗／正常玩家／Go 三方或 C 機器碼完成。
- 收尾文件與工具檢查 21／25 通過，909 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓失敗保留，新文件沒有新增引用錯誤。
- 101 個 C 函式來源身分、十七筆分級出處、最終 driver／收據、原版 EXE／DAT、重定位與擁有權核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：數值輸入器與財政 caller 的真實 C

- 前輪 `7bdafd3` 已推送，分類為 progress。開工核對乾淨工作樹、現行狀態、逆向／IDA／spec／文件職責路由及 GitHub Issue #22。
- 十七函式共 598 bytes，IDA 9.4 probe 保留原始名稱、完整 chunks、operand、xref、file／IDA bytes、三筆重定位、CS:7D01 間接 table 與 CS:7D93 十八格。spec/221 READY 後實作 C。
- 真實六鍵、主迴圈、LAHF／SAHF、word MUL／carry／cap、DIV model、CLD／LODSB、裝置保護／popup wrapper 與財政四 caller 接線。圖形／文字／輸入／popup leaf 為明示 primitive。
- 新 driver 少一個 scenario 外層括號，補齊後編譯器又指明 goto 跨 plans 宣告；回查路由，將宣告移到所有跳躍前，用同一工具鏈重跑。FAR fixture 跳距多一 byte，誤入未覆寫的 code；按指令邊界修正，沒有把它當成 CPU 或產品缺陷。
- 完整 O0/O2 各原版/C 788,931 組、Go 有效數值 787,839 組相同，每版 6164 次完整 1 MB 核對相同，十七函式都有實際入口。Go 使用既有正式 API，原版 CF 與 Go action-valid 不混稱。
- 第一版跳表錯版忽略完成鍵，讀到第 513 筆非宣告事件才返回；加入明示尾端取消並完整重生，最終錯版由第二筆事件取消，verifier 正查 index=2／cancel=1。舊來源／收據留在 c-numeric 的 pre-terminal-cancel 檔案，正常矩陣與判準不變。
- 十個 carry／乘百／除十／最大／FLAGS／glyph step／table／SAHF／取消寫回／徵兵單位突變全拒絕；六個原始 handler 都有真正 table 主路徑向量。正式 Go 原始碼未改，圖形／自然玩家／完整 Go transaction 與 C 機器碼仍未驗證。
- spec/220、re/101 與 spec/78 的同版範圍 backlink 由 verifier 核對。現行 C 台帳累計 118 個，Go 欄只記每函式直接 scalar 比較；完整 Goal 保持 active，本輪依授權完成後 commit、push。
- 收尾文件與工具檢查 21／25 通過，917 列分流、嚴格研究索引、過期斷言、校訂、資產 deny-list 與全部正對照通過。四項既有缺檔／教訓失敗保留，本輪沒有新增引用錯誤。
- 118 個 C 函式來源身分、十七筆分級出處、最終原版/C/Go 收據、原版 EXE／DAT、重定位／終止輸入與擁有權核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍為 27 筆，本輪容器已退出移除。

## 2026-10-08：熱區 map／query 與數值 pixel 接線

- 前輪`99d914b`已推送，分類為progress。核對乾淨工作樹、現況、逆向／IDA／GUI／文件職責路由與Issue #22。
- 十函式491 bytes，IDA9.4保留原始names／chunks／operand／xref／file bytes，無重定位。Writer／reader與舊re/22、re/47證據確認0x1E3D7為熱區，不是glyph，spec/222 READY後實作C。
- 原始map init／register／clear／query、window wrapper、tile flags與border caller真實接線。輸入改事先固定CX／DX pixels／cancel，從十八個raw key table定位，不直接覆寫query回傳AL。
- 生成器找舊模板位置失敗，核對後修正；新C的hex literal尾E與加號連為preprocessing number，補空白再跑，均為工具問題。DI=FFFFh抽樣找出word第二byte應同段wrap，C init改byte write，不移除邊界。
- Init修正後結果仍舊，回查build發現外部C include未觸發cgo cache。加入source manifest SHA-256至CGO flags，另匯出每個binary的build info並驗證，強制來源變更真正重編。
- O0/O2各261,367組原版/C、80組正式Go數值相同，各版2042次完整1 MB相同。全640×400 query、raw word／base wrap、原始0尺寸、重疊CF、flags／border、四劇本及財政／數值pixel鏈均比完整ABI。
- 八個init／high-byte／collision／pitch／tile flags／border／Y offset突變全拒絕，query丟high Y在第163,841組被拒絕，不以隨機pattern相同掩蓋。正式Go原始碼沒改，VGA primitive和正常玩家仍未驗證。
- 勘誤現行glyph導覽與FFF8h word公式，保留舊群組名、operand與歷史收據，取消spec/78的一列假glyph未讀前提；其餘分流保留，backlink由verifier機器核對。
- C台帳累計128函式，完整Goal保持active；本輪依授權完成後commit、push，不把map／scalar結果稱完整UI或C機器碼匹配。
- 收尾文件與工具21／25通過、924列分流與嚴格索引通過；四項既有缺檔／教訓失敗保留。最終研究改寫後四列source fingerprint按實際來源重審，檢查容器cwd補為/repo後同一矩陣重跑，沒有放寬分類。
- 128個C來源、十筆分級證據、最終ABI／map／pixel／compiled source digest與原版EXE／DAT身分核對。輸出UID/GID 1000:1000，既有root-owned路徑仍27筆，本輪容器已全部退出移除。

## 2026-10-08：八函式VGA四plane與保存／恢復

- 前輪`a8605b5`已推送，分類為progress。核對乾淨工作樹、現行狀態、逆向／IDA／GUI／平台規格／文件職責路由與Issue #22。
- 八個原始chunk467 bytes，IDA9.4無relocation；平台probe確認MZ loader與兩個獨立plane／latch狀態。標準VGA引用IBM規格，遊戲RE只追原始port參數、source連續性、dummy read與loop。
- C算法保留原始register／FLAGS／DF／SS、MOVSB／REP、single/double byte列、GC序列與保存major layout；bus callback到第二台獨立成熟VGA。不讀原版結果回填C、不用平面RAM假證明VRAM。
- 新Go工具goto跨configs宣告，移到所有跳躍前以同一矩陣重跑；smoke數量8截斷第三次roundtrip，改為9保留完整三步，沒有改保存／恢復算法。
- O0/O2各9,942例全部RAM／四plane／GC／seq／latch／index／port序列與原始register相同，310次640×400 indexed pixels相同，18條save/draw/restore全plane回到保存前。
- 四個原始2KB頭像×兩destination，原版／C本機debug-palette PNG逐byte相同；來源檔count／大小／SHA核對。自動審查拒絕PNG base64外送檢視，理由可能披露未授權美術；未執行或繞過，改以本機技術收據完成驗證。
- 原版port log的Step保留定位，C不宣稱相同CPU instruction count／wall-clock。所有CGO build flags含source manifest digest，由實際binary build info核對，研究檔用matching_vga tag隔離後完整重生。
- 八個dummy latch／列寬／pitch／OR／plane／read map／wrapper X／DF突變全拒絕，原版資料與衍生圖沒有進Git。正常玩家／整體UI／完整Go圖形與C機器碼仍未驗證，Issue #22保持OPEN。
- re/103、spec/222與re/03加同版後續範圍backlink。C台帳累計136，完整Goal保持active，本輪依持續授權完成後commit、push。

## 2026-10-08：還原程式碼加入 GitHub 紀錄

- 使用者要求把驗證通過的 matching decompilation 程式碼放進現有 GitHub repo。C 來源已逐輪版控，本輪另加入完整組語指令來源與 linker script。
- 沿用 re/90 與既有研究目錄，新增來源索引、公開驗證摘要與 Docker 重跑入口。所有指令保留原始 IDA 定位，沒有加入資料宣告，原版資料區只從本機自備 EXE 匯入。
- 版控來源冷組譯 24,376 條／55,392 bytes 全部匹配；匯入 11,707 個非指令 bytes 後，完整 67,099-byte EXE 相同。指令、linker 常數與私有輸入身分三個負對照均拒絕。
- README、CONTEXT、re/90、re/104與研究索引同輪連結來源、台帳與驗證紀錄；原版 EXE、DB、圖庫、RAM、plane 與 PNG 留在 ignored 本機研究目錄。
- 只掛載版控指令來源、索引與自備 EXE 的第二個容器同樣冷重建通過，不依賴本機 DB 或含資料的組語。公開摘要與實際收據逐 byte 相同。
- 文件與工具檢查 21／25 通過，932 列分流與嚴格索引通過。四項既有失敗仍是歷史圖片／研究產物缺檔、教訓文件不同步與一條教訓文字缺失，沒有新增來源引用問題。
- Docker stdin 檢查首次未開 interactive 而沒有執行，補 `-i` 後同一檢查矩陣有完整 probe。一般 tools 套件本來沒有非研究 Go 檔，改用 `go list -e` 核對新兩檔確實被排除；格式檢查通過。
- 最終核對全部 136 個 C 來源、八筆 VGA 分級、組語自動合併的語意／原始運算元、公開與本機收據。原版與來源雜湊不變，輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，本輪容器全部退出移除。

## 2026-10-09：五個位元對齊與戰術按鈕 C 函式

- 前輪 `47e16c9` 已推送，分類為 progress。核對乾淨 main、現況、逆向／IDA／平台／文件職責路由與 Issue #22，沒有重開已還原的外框函式。
- IDA9.4 十八函式原始 chunks／xref／operand 與三筆既有 MZ 重定位核對，新五函式共 429 bytes。原始位置／hit 表另核對，契約 READY 後才實作 C。
- 新 driver include 首次用了不存在的 world.c，核對現有名稱後改 world_update.c，以同矩陣重跑。C 首次用連續實體 word 讀取，SI=FFFFh 反例只有 plane 不同；回查 CPU.read16 後修正同段回繞，保留反例與原邊界。
- FAR fixture 首次漏 CS 段覆寫，原版與 C 的裝置 restore 分支分歧；補 fixture 契約後同矩陣重生，沒有改遊戲 caller。原始失敗收據與編譯來源摘要留在本機。
- O0/O2 各 4,954 組全部 RAM／四 plane／GC／seq／latch／port／十四暫存器／FLAGS／入口快照相同，683 次內容區像素、36 次 hit 消費通過。六按鈕來源連續、單一重畫與外框直接接真實 C／VGA。
- 十一個首／尾遮罩、dummy read、來源步幅、word 交換、平面重設、JL／OF、六圖重設、位置表、DF 與外框 pitch 錯版全拒絕。ROW 遮罩不再與忽略資料的 mode 1 關聯，不放寬判準。
- 五個 C 來源與分級證據加入台帳，累計 141 個。re/18、re/103、re/104 與 spec/223 回填新範圍入口，原始圖片／收據保留。正式 Go 規則、UI 與存檔未改，正常玩家與 C 機器碼仍未驗證，完整 Goal 保持 active。
- 更新版控組語的五函式分級註記後，24,376 條／55,392 bytes 與完整 67,099-byte EXE 冷重建仍相同，三個負對照拒絕；來源與公開收據身分核對。
- 收尾文件／工具 21／25 通過，940 列分流、嚴格索引、來源引用與研究 tag 隔離通過。四項既有缺檔／教訓問題保留，沒有新增引用錯誤。
- 全部 141 個 C 來源、五筆分級、原版／資產與編譯摘要核對。輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，本輪容器已全部退出移除；依授權提交並推送來源與證據。

## 2026-10-09：十六矩形／選取／計量 C 函式與 Go 校訂

- 前輪 `6139f9e` 已推送，分類為 progress。核對現況、乾淨 main、逆向／IDA／文件／平台路由與 Issue #22；直接沿用 renderer／VGA adapter。
- 新 16 函式 1,169 bytes、36 原始定位、四筆 MZ 核對；spec/225–226 READY 後才實作 C 與正式 Go 修正。兩次右移與取 byte 原樣保留；體力 3=1、1024 兵力回繞與 9999 體力=74 成為相鄰反例。
- 舊素材矩陣給單一命令／外框第 1 段資料；依 sub_100DF 配置與圖庫實際 47,776 bytes，補第 3 段九命令圖與 0x6C0 框 bank。舊 raw receipts 保留，新的實際素材 54／18 組通過。
- 本機原 go.sh 預設 SDK 缺失，沿用現成 hr-go-ebiten image；僅在容器下載 go.mod／go.sum 鎖版依賴，Xvfb 有界 trap。正式 sidebar／selection／layout 冷測通過。
- O0/O2 各 134,958 全 RAM／四 plane／I/O／ABI 相同，131,072 原版／C／Go 長度與 2,441 次內容區相同；十二錯版全拒絕。正式 Go 函式與上限常數逐字擷取，摘要進實際 binary flags。
- IDA 資料項 mnemonic 空白造成窄 decoder 首次無法停止，改用 decoded insn 的 canonical mnemonic；原資料行與分類保留。Opcode 02／03 八條／20 bytes 合併到完整版控組語，新 24,384 條／55,412 bytes 冷組譯、整檔 67,099 bytes 相同，三個負對照拒絕。
- C 台帳累計 157，完整 Goal 保持 active。新的來源與原始位址／分級／收據進 GitHub；正式計量只依原版校訂，規則／存檔格式未改，自然戰術與 C 機器碼仍未完成。
- 收尾文件／工具 21／25 通過、950 列分流與嚴格索引通過。原始雜湊曾被掃描器誤綁到同段工具連結，分開 DB／manifest 檔案段落後過期斷言檢查通過；四項既有缺檔／教訓問題保留。
- 正式 wlgame 全套 `go test -count=1` 在 SDK／Xvfb 通過。最終 157 個 C 來源、16 分級、8 追加組語、公開／本機收據、compiled source 與擁有權核對；輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，容器全部退出移除，依授權提交並推送。

## 2026-10-09：顯示清單九 C 函式與十七原始入口

- 前輪 `c7dc8df` 已推送，分類為 progress。核對現況、乾淨 main、逆向／IDA／文件與 Issue #22，沿用獨立成熟 VGA，還原原始登記／DS table 與九個 opcode。
- 臨時 DB 解碼先保存原名稱／資料行／bytes，遇 ret 後 word 資料與切在前綴中間的邊界後回查寫入／分支；保留 1EE62 與 1F45B 資料，完整 code 不以資料硬解。
- 固定 IDA 832 指令產生 native C 具體運算與 label，不讀 opcode／不呼叫 guest CPU。明示 MUL／DIV operand 修正，原始像素 EFEF 與舊矩形 helper 給兩側對称擷取；有限字串 fixture DS 改讀宣告 bank。
- 十個原始場景與九 opcode、替代 DS／交換 table、線段／底紋／雙色／文字迴圈 O0/O2 各 382 全 RAM／四 plane／I/O／ABI 相同。X／Y 相同使錯源突變不能拒絕，改不對称輸入後同矩陣重生，十個錯版全拒絕。
- 新增九個原始 named C 函式與 17 原始 code 入口，函式界線導覽與無函式 data 項分開記錄，glyph raster 還是明示 fixture。正式 Go 未改，完整 Goal 保持 active。
- 317 條／746 bytes 未分類 code 經實際組譯匹配，來源同輪進 GitHub。合併為 24,693 條／56,138 bytes，完整 67,099-byte EXE 與 SHA 相同，三個錯版拒絕，資料槽維持原始 bytes。
- 生成 C 在乾淨容器從同一 IDA 證據逐 byte 重生相同。完整 C 來源摘要、166 named 函式／17 code 入口、317 補充指令與公開／本機收據核對；只保留 glyph raster 的明示邊界。
- 收尾文件／工具 21／25 通過、958 列分流與嚴格索引通過，四項既有缺檔／教訓問題保留。舊 verifier 的補充範圍改為原兩個 handler 子集，修正縮排後全部 AST／語法通過，未改其原矩陣。
- 研究 driver 以 matching_display 隔離，輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆，容器已全部退出移除；依授權提交並推送原始來源與收據。

## 2026-10-09：原始 C glyph raster 與三字名稱

- 前輪 `4d31520`已推送，分類progress。核對現況、逆向／IDA／文件／字型平台路由與Issue #22，目標為移除 glyph raster no-op。
- 原始C四named函式282bytes與兩入口157bytes由IDA9.4核對，三字caller最初遺漏尾端epilogue後補完整邊界。spec/228 READY後產生C，所有原位址、far與self operand保留。
- 兩側各DOS／Machine／Font cache；C adapter只呼叫原平台IntHook，沒有CPU.Step。原始F720取得向量、32-byteSSbuffer、字庫讀取與真正VGAraster皆執行。未用import移除後同輸入乾淨重跑。
- O0/O2各252組全RAM／四plane／latch／port／ABI相同，每組取內容區像素。兩側740全形／150半形、缺字0，八種對齊、上界、透明／背景、三字live字色與十場景覆蓋，九錯版全部拒絕。
- 服務stub早先導覽名錯寫成sub_10410，改真正0080:0410／0414再重生所有最新收據，KI位址與平台位址不混用。固定生成C在乾淨容器逐byte再生一致。
- 組語補充保留先前317指令，再追加21／59bytes全部實際組譯相同；總24714／56197與完整67099-byteEXE匹配，資料槽、字庫與原版不進Git。
- C台帳170函式，19個原始code入口；glyph fixture已取消，正式Go未改，自然UI、原TSR硬體時序與C機器碼仍未完成。完整Goal保持active。
- 最終公開摘要在所有錯版程序結束後重新綁定最新收據，四C來源／兩入口與真實字庫／compiled flags 逐項核對，生成C在乾淨容器重生相同。
- 文件／工具21／25通過、965列分流與嚴格索引通過，四項既有缺檔／教訓問題保留。舊display verifier保留原317指令子集，新增glyph指令不破壞歷史檢查條件。
- 170個C來源／19code入口、24714組語與整檔、原版／字庫、source摘要及擁有權核對；輸出UID/GID1000:1000，既有root-owned路徑仍27筆，本輪容器已全部退出移除，依授權提交與推送來源證據。

## 2026-10-09：數字 C 還原與 GitHub 紀錄

- 前輪 `85f06de` 已推送，核對乾淨工作樹、復古逆向／IDA／文件與 backlink 路由，讀取 Issue #22 現況。沿用每輪 commit／push 授權，完整 Goal 保持 active。
- spec/229 READY 後還原 `sub_1062F`、`sub_1069A`、`sub_106DE`、`sub_11E17` 與原始無名稱 `0x10984` 入口。固定 IDA 143 指令生成 native C，字庫取實際 ICONGRF 第三段加 `0x840`。
- O0/O2 各 4,881 組完整 RAM／四 plane／GC／seq／latch／port／ABI 與內容區像素相同。八個錯版全部拒絕；產生器在乾淨容器重生相同，source digest 綁定實際編譯 flags。
- 新收據限數字 raster、日期與 SS 參數入口，較早 editor 收據仍保留原明示邊界。第一個 DIV 限制、原作者工具鏈與自然玩家 UI 不猜補，正式 Go 本輪沒有改動。
- 分級語意、四函式／一 code 入口與公開驗證摘要進版控，台帳目前 174 函式／20 code 入口。組語分級註解重生後，24,714 指令與完整原 EXE 相同，三個組語錯版拒絕。
- 文件與工具 21／25 通過，971 列分流與嚴格問題索引通過。新增規格份數已回填 README；四項既有缺圖／研究產物／教訓同步問題維持基線。Go package 檢查補上容器 HOME 與 cache 後同命令通過，研究 driver 正常被預設 build 排除。
- 公開與本機收據、174 函式／20 入口來源、整檔組語、原版字庫、source digest 與 UID/GID 逐項核對。既有 root-owned 路徑仍 27 筆，沒有新 root 輸出；本輪容器已退出移除，依授權提交並推送來源及驗證紀錄。

## 2026-10-09：TALK、參數與原始肖像載入

- 前輪 `06fff77` 已推送，分類 progress，核對乾淨工作樹、逆向／IDA／backlink 路由與 Issue #22。依既有 commit／push 授權繼續完整 Goal，不把本輪 renderer 完成當成全專案完成。
- 八個 named C 函式與六個原始 handler 共 743 bytes／376 指令，spec/230 READY 後接到真實字庫、數字、背景與肖像。保留原 near table、word offset、CF、LAHF／SAHF、原始四格替換與 DOS open／seek／read／close。
- IDA 探針原先只按 basename 找出處，上一輪 spec/229 使匯出失敗。改完整 docs 唯讀掛載與完整路徑 resolver，同名錯目錄與越界拒絕通過，再從原 EXE 建一次性 DB。
- driver 少了入口計數後補回，status 樣本的 AX personality 與 speaker 分離，避免越界索引。初期 corpus 格／像素單位混用讓 glyph 被裁切，字庫計數暴露假完整；改格座標 0，停止舊工作並以同 wrapper 全部重跑。
- O0/O2 各 1,277 組全 RAM／四 plane／ABI／port／DOS API／像素相同，1,022 槽各一次、150 肖像、15,287 全形／178 半形、缺字 0。十個錯版與裁切工作量護欄拒絕通過，生成 C 在乾淨容器重生相同。
- 正式 Go 保持既有實作，INT50 媒體 UI、自然事件／producer 與 C 機器碼沒有擴大完成聲明。分級語意與公開來源台帳目前 182 函式／26 code 入口；整檔 24,714 組語與 67,099-byte 原 EXE 仍匹配。
- 較早三份訊息／快取文件回填可驗 backlink；直接位址 consumer 的「五支」改為三支，數字格式移出已解缺口，其餘 producer／屬性限制保留。976 列分流與嚴格問題索引通過，文件與工具 21／25 通過，四項既有缺檔／教訓問題維持基線。
- 最新公開收據在驗完六個唯讀素材與 resolver 拒絕對照後重生，182 函式／26 入口的來源 binding、生成 C、編譯 flags、所有八個新入口與完整 RAM／plane 收據相符。研究 driver 預設 build 排除，所有輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆。
- 舊裁切工作已停止且刪除，最後全矩陣與每批工具容器皆已退出移除；依授權把還原來源、原始定位及驗證紀錄提交並推送 GitHub，完整 Goal 保持 active。

## 2026-10-09：訊息等待與完整 mouse segment

- 前輪 `52c2b09` 已推送，分類 progress，核對原始現況與逆向／IDA 路由後接訊息顯示／等待／擦除。19 named 函式與兩 raw 入口共 957 bytes／418 指令，spec/231 READY 後還原。
- mouse 遠呼叫與 16 項原 table、cursor 保存／恢復／兩層繪圖、CLI／STI／FLAGS 及 callback 為真實 C。INT33 用獨立平台 API，C 不跑原 guest，timer 以明示 counter checkpoint 輸入，不外推硬體 cadence。
- 原 `EB 0F` 跳過左鍵但問右鍵及 counter，NOP 模式加入左鍵，兩者均有十個 counter 間隔 timeout。spec/45、re/42、re/25 與 mechanics/15 回填分級勘誤，正式 Go 沒有猜改。
- Driver 原先錯認 CPU 有 InHook；查實 Bus 接口後用只記錄 IN 的委派 bus，其他記憶體／OUT 契約不變。生成 C 的 NULL 尾端與 raw callback／patch 保留原邊界和實際 CS。
- O0/O2 各 210 全 RAM／四 plane／IN／OUT／mouse 狀態／ABI 相同；三個背景還原、十九 named／兩 raw 入口與十二錯版通過，生成 C 乾淨重生。C 台帳 201 函式／28 raw 入口。
- 組語新增21指令／71bytes與三個far重定位，比對 file bytes 與 IDA loader bytes。匯出器舊補充數量約束修正後重生，24,735指令／56,268codebytes與整檔相同，三組語錯版拒絕；舊glyph verifier只驗原338指令子集，保留歷史矩陣。
- 收尾讀取 Issue #22 現況，仍為 OPEN，未以局部 C 換成正常玩家完成聲明。982 列分流、嚴格問題索引、生成 C／source flags 與所有來源 binding 通過；文件與工具 21／25，四項既有缺檔／教訓問題維持基線。
- 六個唯讀資料／字庫身分重驗，201 函式／28 入口與公開／本機收據相同。舊 glyph verifier 在隔離暫存輸出通過，不覆寫其歷史收據；本輪所有輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆。
- 全矩陣與每批 Docker 容器皆已退出移除；依授權提交並推送來源、原始定位與驗證紀錄，完整 Goal 保持 active。

## 2026-10-09：selector 與原進言 GUI C 接線

- 前輪 `95b951a` 已推送，分類 progress，核對原始狀態、逆向／IDA 路由和 Issue #22 後接 popup 選列、捲動與 cursor 保存還原；spec/232 READY 後還原 14 named 函式和兩 raw 入口，474 指令 native C。
- 原 `0x9441`／`0x9443`／`0x9465`／`0x9467` live operand與 `0x36D` row callback保留。map／pixel兩分支與原 `sub_193E9`／`sub_13B7E`、speaker／advisor caller接入真正字形／mouse／等待C；原 nullsub_5 照原image，沒有猜補。
- Prototype未初始化CS:987C畫面保存段，原版保存覆蓋服務區導致超budget；查原sub_19796後設定7200獨立段，用同wrapper和輸入重新跑。沒有調高steps或以fixture代替原分支。
- O0/O2各138完整RAM／四plane／IN／OUT／Mouse／ABI／像素相同，三個XOR雙切還原、2153全形／6半形、缺字0；十二錯版全拒絕，生成C在乾淨容器重生相同。原版正常玩家與scroll helper動態改寫仍有獨立界線，正式Go不改。
- 原re/84錯置file／IDA位址與Y回置值已按完整原branch和固定回播勘誤，保留原斷言形成背景。新增10組語／19bytes，先前369以前的補充原記錄保留，24,745指令／56,287codebytes與完整EXE匹配。
- 全錯版終止後獨立驗證最新完整來源與 flags；987 列分流與嚴格索引通過，re/84 已閉合缺口加明示無缺口標記。文件／工具21／25通過，四項既有缺檔／教訓問題維持基線；gofmt與預設build排除通過。
- 舊glyph與input verifier在隔離暫存輸出通過，原始歷史收據不覆寫；215函式／30raw入口binding、六個原始唯讀素材、完整組語／C收據與擁有權核對。所有輸出UID/GID1000:1000，既有root-owned路徑仍27筆。
- 最後全矩陣與每批工具容器皆已退出移除，依授權把原始定位、C來源與驗證紀錄提交並推送GitHub；完整Goal保持active，未把局部consumer取代正常玩家驗收。

## 2026-10-09：完整資源讀檔、BGM cue 與分段配置 C

- 前輪 3c5ca536 已推送，核對現況、逆向／IDA 路由與 GitHub Issue #22。原 scroll helper 的直接 xref／逐段候選沒有定位寫入端；以三個已知 writer 正對照封存查詢範圍，未據此排除間接改寫。
- spec/233 READY 後還原 sub_187AF、sub_1E378、sub_1F4A2、sub_10241、sub_102C2 與 sub_100DF。六函式 369 bytes／165 原指令；使用獨立 DOS／INT61 平台服務，C 不執行 guest CPU。
- 初版大檔案測試覆蓋測試堆疊，按原讀檔區間移開 SS；沒有更改 read count／stride 或放寬原控制流。第 6 個錯版在 0x6000 bank 經三次左移回繞而漏過，改用 0x6100 測試段後完整重跑。
- O0／O2 各 93 完整 RAM／VGA／ABI／API／sound 狀態相同，21 個來源 buffer／尾端守衛與十二錯版通過。生成 C 乾淨重生相同，完整來源 digest 與實際編譯 flags 綁定，原平台來源前後一致。
- 六函式分級索引、公開收據與原資料／IDA DB 身分入版控；舊 MMAP、BGM 與 allocation 證據加回鏈。C 台帳 221 函式／30 code 入口；正式 Go 沒有更動，INT50／失敗退出、正常玩家、原 TSR 音訊與 C 機器碼保留界線。
- 組語只重生分級註解，24,745 指令／56,287 code bytes 與完整 67,099-byte EXE 仍相同，三個組語錯版拒絕。嚴格索引／分流與來源 binding 通過；文件／工具 21／25，四項既有歷史資產／教訓缺口維持基線。
- 索引檢查曾多報一列 spec/233 把已採用 API 契約寫成未解。改列真正未驗的 TSR 音色／硬體時序後重跑，新增文件沒有剩餘索引錯誤。
- 所有本輪輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆。依授權把來源、定位與驗證收據 commit／push；完整 Goal 保持 active。

## 2026-10-09：原始世界顯示格與完整圖塊合成 C

- 前輪 d713d43 已推送，分類 progress；重新讀取現況、逆向／IDA／文件職責路由與 Issue #22。由場景退出鏈定位世界顯示格；完整場景入口仍留待上游 producer／重畫接線。
- spec/234 READY 後還原十函式 710 bytes／328 原指令，保留初始化、40×23 記錄、鏡頭取樣、四張推入上限、髒格／快取、128-byte 底圖與 160-byte 覆蓋圖合成及四 plane。
- 第一個 IDA 匯出對 far relocation 的檔案／loader bytes 直接比較失敗，按既有 modal 契約分列重定位。兩個錯版位址按固定 IDA 指令修正，原始控制流未更改。
- 一次自動權限審查因逾時未完成；依回覆只重試一次，採明確 tools／研究輸出可寫、原始資料唯讀掛載後成功。未退回主機執行工作負載。
- DF=1 的直接底圖拷貝反向覆蓋自己的原指令，導致 guest 控制流不同；保留失敗實驗與輸入／source 身分。按原 renderer CLD 的合法前提驗 DF=0，其他安全 DF 分支與位移邊界保留。
- 舊 200 萬指令預算不足整幅重畫；依 920 格 × 最大 6 次合成／搬移／控制的靜態上界設定 400 萬，實測最高 2,969,875。trace 分成兩個既有 8192 snapshot 區段，涵蓋最長的 8,281 次入口，原迴圈不截短。
- O0／O2 各 2,622 完整 RAM／plane／ABI／FLAGS／I/O 相同，256 直接圖塊與 512 遮罩合成獨立比真實來源；遮罩底圖使用真實 MDL 避免零背景漏過錯版。12 個 O0 錯版全拒絕，C 乾淨重生相同。
- 台帳更新後固定所有十個 probe 根，從新的唯讀原始輸入重建 IDA，避免已還原 helper 被 frontier 篩除；新 DB 下再次驗證來源／原函式 bytes 與 C 重生相同，沒有重跑無變更的 native 矩陣。
- C 台帳 231 函式／30 code 入口與來源 binding 通過；組語 24,745 指令／56,287 code bytes 與 67,099-byte EXE 仍相同，三個錯版拒絕。文件／工具 21／25、同四個歷史缺檔／教訓問題，phantom 列表未新增。
- 所有輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆。依授權 commit／push 來源與證據；正常場景／上游 producer／原作者 C 機器碼未由本閉包代證，完整 Goal 保持 active。

## 2026-10-09：軍團／物件 producer 與矩陣 C

- 前輪 32ed4f0 已推送，分類 progress；重讀現況、逆向／IDA／文件職責路由與 Issue #22，接世界顯示格上游的軍團／物件。
- 固定 IDA 證據定位六函式 461 bytes、raw handler 181 bytes、267 指令。矩陣 CMP 被標成資料，第一次 create_insn 遭既有 data item 阻擋；只在一次性 DB 清該項後解碼，保留原資料行和原名稱。
- spec/235 READY 後接 native C 兩次軍團巡覽、單格／矩陣 producer、物件舊相位更新與原小地圖 VGA 點；live capacity 字元由入口 AH 改寫，軍團 3／物件 5，沒有改成單格上限 4。
- 64 冒煙案例後加入獨立 MCH pattern 高層裁切／透明／保護／容量格子表比對，修正 corpus 的相機設定後重跑。曾誤輪詢已完成的 smoke session，依終止狀態不重啟該 session，直接執行新的完整矩陣。
- O0／O2 各 14,622 完整 RAM／plane／ABI／FLAGS／I/O 相同，55 個真實 pattern 共 12,320 次 source 矩陣核對；四劇本原 producer 接完整 C renderer，12 個 O0 錯版全拒絕，C 乾淨重生相同。
- 新 decoded CMP 80 FB 05 的三 bytes 以 GNU 組語匹配，總 24,746 指令／56,290 code bytes 與 67,099-byte 原 EXE 相同，三個組語錯版拒絕。歷史 supplement rows 保留，selector verifier 明確排除此新增指令，隔離輸出重驗通過，不改寫其原收據。
- C 台帳 237 函式／31 raw 入口、所有來源 binding／公共與本機收據通過。文件／工具 21／25、同四個歷史缺檔／教訓問題，phantom 列表無新增；1002 列分流與嚴格索引通過。
- 所有輸出 UID/GID 1000:1000，既有 root-owned 路徑仍 27 筆；本輪容器完成後移除。依授權 commit／push 原始定位、來源與證據；type 3 自然來源、正常場景與 C 機器碼不由局部閉包代證，完整 Goal 保持 active。

## 2026-10-09：原場景主入口與鏡頭／小地圖退出 C

- 前輪 b8ea729 已推送，分類 progress；重讀現況、逆向／IDA 路由與 Issue #22。查原主入口與依賴，固定九個新函式及既有 sub_11D46，578 bytes／255 指令。
- spec/236 READY 後接 tail／far、鏡頭 dirty、小地圖四 plane 背景及 carry 位元對齊框、三句 TALK／等待、停播曲、IVENT→MMAP 恢復與完整 world producer／renderer，沿用原唯一 C sub_11D46。
- 第一版資產準備把 TALK 壓到 MDL、BGM／hotspot 落入解壓 MAP。移至獨立段並加入每次 actual source-bank bytes 核對。第二個分支的 original budget failure 以 trace／RAM 查證：無原 code bytes 變動，世界 DS 誤設程式 DS，性格及 TALK index 讀錯；依原呼叫契約修正後同 budget 重跑。
- O0／O2 各 171 完整 RAM／plane／ABI／IN／OUT／API／Mouse／sound 狀態相同；16 主場景三句／資源／還原完整通過，12 個 O0 錯版全拒絕。驗證器曾 guessed font >1000，改由原文第一行 432 glyph 保守下限及原三句 48 框入口驗收；每版 816 全形 glyph、缺字 0。
- 新 root 固定兩支小地圖 helper，台帳更新後再從唯讀原檔建立一次性 IDA；C 重生、來源／原函式 bytes／新 DB 身分一致，不重跑未變更的 native 矩陣。
- C 台帳 246 函式／31 raw 與來源 binding／公私收據通過；組語 24,746 指令／56,290 code bytes 與 67,099-byte EXE 相同，三個錯版拒絕。文件／工具 21／25、同四個歷史缺檔／教訓問題，無新增 phantom；1006 列分流與嚴格索引通過。
- 所有產物 UID/GID 1000:1000，既有 root-owned 路徑仍 27；容器完成後移除。依授權 commit／push 原始定位、C 來源與證據；自然事件／長程玩家／原 C compiler machine code 仍是獨立界線，完整 Goal 保持 active。

## 2026-10-09：君主出陣、自動編成與側欄 C

- 由原上游 caller 固定 14 函式、1,056 bytes／501 指令。spec/237 READY 後還原完整君主出陣、六格選兵、退兵／補兵、佔用圖與四個側欄 bit。遷都清單保留後續獨立閉包。
- IDA 匯出與生成器曾錯誤併行，輸入暫時不存在；待匯出成功後依序生成，未更改原始控制流。佔用圖使用 B000 bank，避免 VGA／背景保存重疊。
- O0／O2 各 451 完整 RAM／plane／ABI／FLAGS／IN／OUT／API 相同；64 個完整出陣場景接原三句對話、音樂、游標與世界重畫。十個編譯錯版全拒絕，C 乾淨重生一致。第一版錯版只改比較門檻，案例下仍輸出同值；改為實際上限常數 100→99，並完整重跑。
- 資金高位符號分支跳過不足判斷，舊機制表公式只適用 BH<80h；以原指令及實際主入口補回例外，不改 Go 規則。原 cmp bx,12Ch 保留。
- C 台帳 260 函式／31 raw 入口。來源、固定輸入、IDA DB、工具版本及公開收據依授權提交 GitHub；正常玩家長程與原 C compiler machine code 仍未完成。
- 嚴格索引與 1011 列分流通過；文件／工具 21／25，同四項歷史資產／教訓缺口，phantom 清單無新增。新文件已補齊規定狀態／日期欄位後乾淨重跑。
- 完整組語 24,746 指令／56,290 code bytes 與 67,099-byte 原 EXE 相同，三個錯版拒絕。全部 C／raw 來源與公開／本機收據 binding 通過；新產物 UID/GID 1000:1000，既有 root-owned 路徑仍 27，無 .md 目錄。
- 本輪一次性容器均已移除，沒有殘留執行中或停止的專案容器；正式原版資產與 EXE 留在本機。

## 2026-10-09：據點清單、排序與完整遷都 C

- 前輪 60e2fb5 已推送。由原遷都 caller 追至 18 函式與 5 raw，2,002 bytes／835 指令；spec/238 READY 後接完整原清單、五個排序欄與首欄重建、四 plane 滑塊、特殊返回與遷都 writer／側欄。
- 第一個差異是四段 trace 只清總計數，留下上例 snapshot；逐段清除後重跑。欄位夾具超出有效 0–5，按六欄描述子修正。遷都接受原版進入錯誤 opcode，是漏掉 sub_1030F 初始化遠指標；依原參數補齊，同預算通過。一次修正命令因字串跳脫失敗，改用獨立腳本。
- 取消選列後的 JMP 重入 sub_18412 漏一個觀測點；保留無新增 stack frame 的重入 snapshot。O0／O2 496 組通過後，補六表頭實際點擊，停止自己的舊矩陣後統一重跑最終 520 組。
- O0／O2 各 520 全 RAM／plane／ABI／FLAGS／IN／OUT／API 相同；120 原交換排序、36 builder、48 捲軸公式另核對原版。24 表頭、16 完整遷都包含四次原首都重選及八次 writer，十二錯版拒絕；C 乾淨重生一致。
- 勘誤：re/26 於 2026-09-02 宣稱 word_183D3「沒有讀取端、寫了不用」。新原 bytes 與 live patch 表明它是 0x183D2 MOV SI 的立即值，0x1821E 寫入後由 CPU 取指消費；舊 grep 地址漏掉指令立即值。已修現行說明，保留此歷史與新證據入口。
- 七個 patch 點帶出的 83 個原基準未覆蓋指令 215 bytes 以固定 GNU binutils 匹配，另存 c-list-code.json；歷史 rectangle supplement 370 行保持原樣。完整組語 24,829 指令／56,505 code bytes，原私有資料匯入降至 10,594 bytes。
- C 台帳 278 函式／36 raw；正式 Go 未更動，其他清單家族、自然玩家長程與 C compiler machine code 不由局部閉包代證。來源、固定輸入、原始定位與收據依授權 commit／push，完整 Goal 保持 active。
- 字型觀測定位 299 次缺字都在直接捲動 helper，字碼是指令 bytes；1851A／18546／184DD 的夾具 DS 誤設 CS。按原 callback 世界 DS 修正後，同 520 組與十二錯版重跑，兩側缺字 0，字型服務未改。
- 逐條以原 24,376 指令與歷史 370 行覆蓋核對完整 835 指令，查得 83 指令／215 bytes 尚未進組語基準，包含七個 patch 點的後續區段。固定轉換器以原編碼及 load hint 全數匹配。台帳更新後明列 18 個固定根，保留 sub_1748F renderer 匯出，再驗 C 重生與收據。
- 最終 1014 列分流與嚴格索引通過；文件／工具 21／25，同四項歷史資產／教訓缺口，phantom 無新增。三條不可變位址回鏈護欄通過，全部 835 條本輪指令已在完整組語範圍內。
- 最終 24,829 指令／56,505 code bytes 與 67,099-byte 原 EXE 相同，三個組語錯版拒絕；278 C／36 raw 來源與公開／本機收據 binding 通過。產物 UID/GID 1000:1000，既有 root-owned 仍 27，無 .md 目錄。
- 本輪一次性容器均已移除，沒有殘留執行中或停止的專案容器；原版資產、EXE 與研究 binary 只留本機。

## 2026-10-09：四家族清單與回呼 C

- 前輪 5fd5145 已推送。讀現況／路由／AGENTS 尾段與 Issue #22；依專案規則派唯讀代理核對契約、盤點 Go 工具鏈，再分工單一資料參照檔與獨立 ASM 補充，均未越界改檔或 Git 寫入。
- spec/239 READY 後還原 11 函式／12 raw、1,870 bytes／902 指令，共用既有 list 引擎。保留 126／128／22 範圍、X/Y patch、DS:CFF、資格／職務／哨兵、byte 士氣、原三位數字與表頭。
- 獨立 cache 參照起初誤把 XOR DL 當 loop 前序；原版參照抓出差異，實查 LOOP 回到 17A9A，改成每筆清0。原 C 生成控制流未修改；保留交錯邊界向量並以錯版驗跨項殘留。
- O0／O2 各 804 全 RAM／plane／FLAGS／SS frame／IN／OUT／API 相同，168 builder／240 排序／16 cache 獨立核對、156 表頭、28 成功選取與十三錯版通過。字型缺字0，C 乾淨重生一致。
- 170 指令／361 bytes 以固定 binutils 全數匹配，包含 imm_at_17284=FFFF linker 符號。補充取原 source-map、370 舊行與83 list行為固定覆蓋來源，兩次重跑同SHA；未改舊補充。
- 預設 Go image／module volume 已不存在，沿用既有 eob Go1.26.7／Ebiten SDK。go.sh 加專用快取、原素材唯讀、有界 timeout／Xvfb trap；精確 go.mod/go.sum 依賴先由官方 proxy 取得，再離線冷跑 vet／test，39 套件通過。實際 Docker inspect 證明原素材RW=false、network none、UID1000與資源上限；不碰共享 cache volume。
- 舊 re/26／27、mechanics/10、spec/38 的排序、士氣寬度、勢力數字位數與cache入口已勘誤；四格幾何保留，正式 Go 引擎無改。開局 renderer 舊未解列有新證據回鏈。
- C 台帳 289 函式／48 raw；完整組語 24,999 指令／56,866 code bytes，私有匯入10,233 bytes。正常任命／外交／新遊戲長程、非法空清單排序與原 C 機器碼保持證據界線。來源與收據依授權 commit／push，完整 Goal 保持 active。
- 正式 check.sh 完整執行，Go stages 通過，文件階段停在既有兩張playtest/19 PNG與277 phantom；逐項25檢查為21通過／4已知失敗，無新增引用。re/27原未解全解後補缺口無標記，嚴格索引與1018列分流通過。
- 最後逐條證明本輪902指令均在編譯組語範圍；完整67,099-byte EXE相同、三組ASM錯版拒絕。289 C／48 raw、公私收據與五回鏈、39套件冷Go測試及實際唯讀掛載binding全通過；產物1000:1000，既有root-owned27、無.md目錄。
- 代理已完成且沒有後續檔案寫入；本輪一次性容器均已移除。只提交還原來源、定位、工具與收據，原始資產／EXE／研究binary／module cache留本機。

## 2026-10-09：玩家軍團編成 C

- 前輪 ad020ad 已推送。沿用復古逆向／IDA 路由與現有工具；依 AGENTS §10 分派唯讀契約審查、單一獨立資料模型、指令覆蓋及正式 Go 冷測，均按限定檔案交付。
- spec/240 READY 後還原七函式／247 指令／561 bytes，選將、六槽圖示／數字、兵種切換、確定／取消接回既有唯一 C 規則、清單、TALK、滑鼠及 VGA。
- O0／O2 各 200 完整狀態相同，十個編譯錯版拒絕；獨立分兵模型先驗原版，字型缺字 0。保留僅初始化 type、不清 men、取消不重算總兵力與原熱區邊界。
- 初次 smoke 的資料模型未涵蓋 switch/offpanel 群組，屬 harness 分支缺項；補齊後同一工具乾淨重跑72組與完整200組通過，原 C 邏輯未改。
- 全部247指令已存在前輪組語，新增0；固定binutils另獨立重組561bytes相同，外部語意索引附加原始定位，整檔基準重新冷組譯。
- 正式Go冷vet與39套件通過；C台帳296函式／48raw，正常玩家長程與原C機器碼仍未驗。來源與收據依使用者授權commit／push，完整Goal保持active。
- 提交前25項檢查21通過，4個既有失敗不變：歷史playtest19兩張缺圖、277筆phantom引用、lessons render過期、舊guard文字缺失。完整check已跑並停於既有phantom，個別後續檢查另跑；新失敗0，正規化phantom增刪均0。檢查收據為workplace/matching-decompilation/checks/formation-checks.json，SHA-256 `2b6481cbea1e4cb99ae7549f6e4759abd57ee8c55f576af1aadbfeba2400c17a`。
- root-owned維持27個既有路徑，沒有.md目錄或本輪root-owned輸出；本輪容器已收尾，沒有新增image。

## 2026-10-09：原人事選單與任免 C

- 前輪34e3b1f已推送。路由與現況／Issue22查核後，按AGENTS§10分工唯讀契約、單一資料模型、獨立指令覆蓋與Go冷測，未越界修改。
- spec/241 READY後還原七函式／185指令／459bytes，保留原四項live table、兩層任命、拒絕重試、解任XCHG、原POP BX/SI差異與唯一共用清單／UI。
- 初編譯的pn_call與既有politics.c同名，改為psn前綴後同命令smoke112組通過。矩陣編輯曾多插入非選單重複案例，確認執行handle後停止自己的容器，修正後同容器同命令乾淨重跑，未把兩項工具問題記成產品缺陷。
- O0／O2各2,332完整裝置相同，十二錯版拒絕；2,064 helper覆蓋全部合法官員索引與FF，120任命／96解任／52選單，含16次live table修改。獨立欄位與CF/BX模型先驗原版；缺字0，來源乾淨重生一致。
- 原版任命保留非零經費、解任清零已實跑正對照；正常自然撥款玩家長程保持證據界線。更正re/22舊未讀說法、spec/148前序XOR解讀，spec/143收斂至自然撥款長程限制，mechanics/60同步。
- 185指令全已有覆蓋，binutils獨立重組459bytes完全相同；組語註記重生後整檔67,099bytes仍匹配，三反例拒絕。303 C／48raw，完整Goal保持active，來源與收據依授權提交推送。
- 錯版包裝器曾沿用編成group而產生0案例，被exit guard拒絕；第7錯版曾使用不存在的byte setter，編譯失敗未當成錯版拒絕。修正後加入未宣告函式編譯護欄，13種正常／錯版語法審查全通過，並以最終來源重跑O0/O2各2,332組與十二個實編譯錯版。
- 正式Go冷vet／39套件通過，無cached；完整check已跑並停於既有缺圖／phantom，25個單項另跑21通過、4舊失敗不變，新失敗0。phantom正規化277無增刪，lessons舊日誌相同；三條resolution backlink護欄通過。檢查收據workplace/matching-decompilation/checks/personnel-checks.json，SHA-256 `af6fdd61d5e19247ffa302b2499699e4d0f5fc227c037b0180b797c9c9bbc3e0`。
- 檔案擁有權維持27個既有root-owned路徑，無新root-owned輸出或.md目錄；本輪容器清理完成，未建立image。

## 2026-10-09：據點資訊與軍團面板 C

- 前輪74e50b6已推送。查現況／路由／Issue22與原RE後，依AGENTS§10分工契約、單檔獨立模型、指令覆蓋、13版語法預檢及Go冷測。
- spec242 READY後還原10函式310指令711bytes，據點完整局部視窗、景觀讀取、軍團三名／數值／六槽／退出與自勢力恢復，共用原C服務。軍團行軍指令上游未代證。
- 初probe因範本替換把具名roots誤列raw，修正為10個明確roots與空RAW_ROOTS後由IDA乾淨重匯出。smoke136通過；完整矩陣兩種最佳化700通過後，唯讀審查指出範本註記殘留，停止已確認自己的容器，更正說明並以最終來源同命令重跑。
- O0/O2各700完整裝置、十二實編譯錯版、獨立參數／KYOGRF buffer／世界保存／ABI通過；字型缺字0，C乾淨重生一致。
- 回填re32/50、spec23/24：SI是索引×32，首都比較左移，軍團數字0F、byte士氣、type0下溢D200、bit02恢復自勢力情報，原版據點入口不只地圖。正式Go未改。
- 310指令全已有覆蓋；固定binutils獨立重組711bytes相同。組語語意索引重生後整檔67,099bytes仍匹配，三反例拒絕。313 C／48raw，完整Goal保持active，依授權提交推送來源與收據。
- 正式Go冷vet／39套件通過，無cached；完整check已跑並停於既有index兩缺圖／phantom，後續25單項另驗21pass/4舊fail，新failure0。phantom初次275與當前來源矛盾，只重跑第3項；docs43前後SHA穩定時為277並與基線逐列相同，不宣稱移除舊問題，原因未證實。13/14隔離自測的審查deadline逾時後依工具指示重試成功，未改權限範圍。
- 六條resolution backlink通過。檢查收據workplace/matching-decompilation/checks/details-checks.json，SHA-256 `a58a64fcfc845a0524b737b69e118ce60267e2116c3e6643353df0e4d9011db5`；完整組語、313 C/48raw來源雜湊與公開/私有收據一致。擁有權檢查維持27個既有root-owned，無.md目錄；本輪容器清理完成，未建立image。

## 2026-10-09：軍團行軍選點與狀態分派 C

- 前輪b4c2cb4已推送。查現況／路由／Issue30，依AGENTS§10分工原契約、單檔獨立模型、MZ重組／17版語法預檢及Go冷測。
- spec243 READY後還原25函式613指令1477bytes，補原熱區表159D2[22]的15AB6與兩個鏡頭callee；五個farcall用runtime operand，原檔段1000／IDA載入2000分開核對。
- 初smoke缺既有RNG ABI路由，接sub_1ECE0_abi後同命令352組通過；1,308全矩陣與16錯版通過後，為明示熱區優先／重畫新增12個必要循環案例，最終來源356抽樣與O0/O2各1,320、16實編譯錯版通過。
- 原版先經狀態模型；controller捕原版獨有分派前snapshot，固定raw RNG/table/C/S，不重擲。全部12stage/9handler、隊列、DI residual、補兵／解散、兩鍵mask、模式取消重選、小地圖及日期停止均驗。
- 回填spec39/149、re47/85與mechanics20：取消仍分派/OR2、解散只清byte00、移動仍回mask、搜尋完整性前提、日期停不等於動畫RAM不變。正式Go未改。
- 613指令已有覆蓋，五筆MZ反向重定位與獨立重組相同；組語註記重生後67,099bytes整檔仍匹配。338 C/48raw，完整Goal保持active，依授權提交推送。
- 正式Go冷vet／39套件通過，無cached；25單項先完成21pass/4舊fail、新0，phantom前後來源SHA穩定、277零增刪，13/14隔離自測通過。之後才跑完整check，停於既有index兩缺圖／phantom，避免並行生成docs43。四條回鏈、5筆MZ、raw RNG每版16次／clock0均通過。檢查收據workplace/matching-decompilation/checks/march-checks.json，SHA-256 `8506b746b4ac5c1a727359cac3f37777b66647afdde2b7c799aecb1d21badc79`。
- 338 C/48raw全部來源雜湊與原組語整檔收據一致；擁有權基線27、無.md目錄，本輪輸出UID1000、容器清理完成，未建立image。

## 2026-10-09：政略指令與進言理由 C

- 抽樣先修共用政治／世界 helper 路由，並依 OR10h 後四次 SHR 訂正獨立模型的理由位元順序。正式 Go 未改。
- 原版借位扣信賴度會經11CB1恢復主循環保存的SS/SP；孤立場景缺外層堆疊，不能當作繪字破壞或正常RET。矩陣保留無借位恰好扣到零，非區域退出另列界線。
- 依 spec244 READY 還原29函式934指令2184bytes，保留原定位、byte溢位、DS／CF與尾跳。
- O0／O2各3,348組、十六實編譯錯版拒絕。原版先驗獨立模型，再比完整裝置狀態；INT61回覆明示為受控條件。
- 十五個錯版由狀態比較拒絕；第九個錯版另由原生DIV保護中止，以固定退出碼2、精確assertion／SIGABRT與編譯來源核對。兩種拒絕分開記錄。
- 25單項為21通過、4既有失敗：index兩個歷史PNG、phantom277筆、lessons render／verify。沒有新增失敗；phantom正規化零增刪，lessons失敗日誌與前輪逐字相同。
- 輸出UID/GID1000；擁有權衛生掃描維持27個既有root-owned路徑，沒有新增或異常md目錄。原生矩陣與單項檢查容器已清理。
- 六筆MZ反向重定位與獨立組譯通過；整檔67,099bytes仍exact，367 C／48raw；正式Go未改。
- 依使用者授權提交與推送。原TSR音效、自然長程、原C機器碼與完整Goal仍待完成。
- 最終正式Go vet／39套件冷測通過，cached0；完整tools/check.sh停於既有index兩缺圖與phantom277，後續檢查由先行25單項覆蓋。Go與C最終收據已重綁；容器已清理，依授權提交推送。

## 2026-10-09：信賴度非區域退出與原調色盤淡出 C

- spec245 READY後還原3函式89指令179bytes；KiMachine16與共用hooks header不變，原SS/SP與RET完成後由活著的C runner轉移，舊caller不續行。
- prefix前缺D44時城市資料寫入程式，city2把10676的D0改成CF IRET，造成原版跳02C0；補完整地圖指標、數字字模與已選玩家。rawmap先驗，prefix城市tile／評分視為原版runtime資料。
- 正常scene locals在釋放後被19321覆寫，修正獨立模型只在仍存活的frame驗attempt/budget；其餘用信賴度／CF／原決策軌跡。
- O0/O2各2520完整RAM/planes/DAC/registers/FLAGS/stack/IN/OUT/API、12原frame快照、四獨立maphash與八錯版通過。370C/48raw，正式Go未改，完整Goal保持active。
- 右鍵表原DS5A06與159D0落下已勘誤，撤回三條由錯誤地址產生的未解；此導航證據不代證右鍵actions或完整mainloop C。
- 本輪Go冷vet／39套件通過，cached0；最後完整check的39套件可重用cache，兩者收據分列。25單項21pass／4既有fail／新增0，完整check停既有index兩缺圖／phantom277；來源掃描前後穩定，13/14隔離自測通過。
- C與Go最終收據重綁，source／89指令／1筆MZ及6筆回鏈一致。依授權提交推送；原版資產與私有RAM／圖片不入Git。
- 本輪輸出UID/GID1000；root-owned基線維持27、無異常.md目錄，Docker容器已清理。

## 2026-10-09：城市與軍團原生初始化 C

- spec246 READY後還原五函式140指令317bytes，另列11BE0七指令22bytes無RET前段，不註冊完整mainloop。
- 兩側從各自raw RAM/VGA開始；原生prefix不再借原版執行後Snapshot。192城、127army槽、active80、descriptor、occupancy byte累加／回繞及score/date原DS契約逐項驗。
- O0/O2各5008組、5416分階段全狀態，10實編譯錯版全拒絕；D44/D850分離與C000 occupancy反例、slot126/127、poison欄位及重複bulk均涵蓋。
- 五完整函式與一有限prefix分開記為375C/49raw。147指令339bytes全有舊組語覆蓋，獨立重組exact、無新MZ，原EXE67099仍exact。正式Go未改，完整Goal保持active。
- 本輪最終來源與收據重綁完成：Go vet／39套件冷測通過，cached0；最後完整check的39套件重用cache，兩次證據分列。25單項21通過／4既有失敗／新增0；phantom277正規化零增刪，lessons失敗與基線相同。完整check停於既有index兩缺圖及phantom，其後階段由25單項補足。
- 最終native驗證確認147指令／339bytes、十錯版、375C／49raw與兩筆回鏈一致；公開C收據已綁定最終Go proof。新文件已列入專案索引，原版資產、RAM與圖片不入Git。
- 本輪輸出UID/GID1000；root-owned維持27個既有路徑，沒有異常.md目錄。本輪容器已清理；依使用者授權commit與push，完整mainloop、正常玩家長程及原C機器碼仍未完成。

## 2026-10-09：世界點選、右鍵取消與淡入 C

- spec247 READY後還原八函式159指令396bytes；原DS表26索引與159B7落入nullsub_1只RET一次。
- O0/O2各2148完整裝置狀態與十二實編譯錯版通過；128組DS≠CS、live表修改、完整低byte座標、城市圖塊與軍團重選由獨立模型先驗原版。
- C來源與較早右鍵證據回鏈已登錄，383C/49raw；159指令已有組語覆蓋，完整EXE67099仍exact。更新排程11CD0、完整mainloop、自然玩家長程及原C機器碼仍待完成。
- 世界互動本輪最終收據已重綁：正式Go vet與39套件冷測通過、cached0；最後完整check的39套件重用cache，兩者分列。25單項21通過／4既有失敗／新增0，phantom277正規化零增刪；完整check停在既有index兩缺圖及phantom，後續由25單項補足。
- 發布程序初次誤用組語收據檔名，依實際code-record-verification.json修正後完成；原生對拍不需重跑。最終159指令／396bytes、十二錯版、383C／49raw及兩筆回鏈核對通過。
- 本輪輸出UID/GID1000；root-owned維持27個既有路徑、沒有異常.md目錄，Docker容器已清理。依授權commit與push；原版資產與私有RAM／圖片不入Git，完整Goal保持active。

## 2026-10-10：據點輪轉、增援與動畫更新 C

- spec248 READY後還原18函式562指令1340bytes；較大主排程普查180函式及7筆未封閉邊界另存私有census，不當已還原。
- O0/O2各8144完整狀態與十六錯版；固定raw258-byte RNG、據點192游標、原距離欄位、內政官／災害與32物件邊界均保留原式。
- 原增援成功／部分成功、武將選取及提示訊息接既有真實C；401C/49raw，完整軍團／交戰／戰術／主排程及C機器碼仍未完成。
- 原第14錯版將生產力借位夾零，因倍率取生產力高byte而在合法輸入域等價。舊8144收據、來源與binary保存在c-tick/equivalent-control-review；核對全部65536個P的最大D界線後，改用錯讀倍率lowbyte，正常演算法不變。新來源hash的O0/O2與十六錯版全部重跑通過。
- 機制補註另設原生C小節，避免生成索引把已證實內容當未解；triage保留1040筆，沒有移除舊列。雲的速度／回繞舊斷言已在兩份舊文件及CONTEXT勘誤，現有Go未改。
- 據點輪轉本輪最終收據已重綁：Go vet／39套件冷測通過、cached0；最後完整check的39套件重用cache，兩次證據分列。25單項21通過／4既有失敗／新增0，phantom277正規化零增刪；完整check停既有index兩缺圖及phantom，後續由25單項補足。
- 最終562指令／1340bytes、十六state-mismatch、401C／49raw及四筆回鏈通過。輸出UID/GID1000，root-owned維持27個既有路徑、沒有異常.md目錄，Docker容器已清理；依授權commit與push，完整Goal保持active。

## 2026-10-10：自我修改尋路與軍團潰散 C

- spec249 READY後還原九函式436指令1016bytes與完整raw搜尋107指令244bytes；三個CS即時operand與原queue/visited保留。
- O0/O2各4588完整狀態與十六state-mismatch；使用可重生synthetic圖表，原版全道路與正常玩家長程不在此證據範圍。
- 新增79指令179bytes；25078/57045/67099整檔exact與3組譯錯版拒絕，source-only supplement與coverage分離且維持whole-file SHA。410C/50raw。
- 舊成本未逐條讀列oq-2ef4b20c0b00482c0266已由新來源取代，原merge-target/Issue3對照與舊內容見Git4324137；Go等距端點與遠端Issue狀態不變。
- 最終來源覆蓋核對明分464條舊覆蓋、79條新增及543條已發布；最初驗證器把兩個scope等同，依實際收據契約修正後通過，原生結果未重跑或替換。
- 尋路本輪最終收據已重綁：正式Go vet／39套件冷測通過、cached0；最後完整check的39套件重用cache，兩次收據分列。25單項21通過／4既有失敗／新增0；phantom277正規化零增刪，完整check停既有index兩缺圖及phantom，後續由25單項補足。
- 最終410C／50raw、4588完整狀態、十六錯版、三筆回鏈與25078/57045/67099整檔重建核對通過；原版全roadtable及Go等距／正常玩家仍未代證。
- 本輪輸出UID/GID1000，root-owned維持27個既有路徑、沒有異常.md目錄，Docker容器已清理；依授權commit與push，完整Goal保持active。

## 2026-10-10：自動戰鬥、退卻與據點易主 C

- spec250 READY後還原21函式796指令1721bytes，保留比例／六槽／原RNG／原易主順序與SS/flags。
- O0/O2各8364完整狀態與二十錯版，含144個實際outer-frame退出；C不得續行舊caller。原完整交戰／戰術／軍團與正常玩家仍待。
- 訂正196CF的文字拷貝方向：DS minimap→ES VRAM；195C9/1963F才是反向保存。原bytes與既有C不變，舊文字由Git a8161c8可回查。
- 431C/50raw；796指令已有組語覆蓋，完整EXE仍exact，三筆較早證據回鏈已登錄。

- 獨立模型首差2289確認軍團+9的乘5，原17012–1701A兩次SHL再加原值；小地圖MOVSW在FFFF的高byte須分段回繞。只修模型，案例／seed／原C保持不變，316組煙霧及完整矩陣均通過。
- 最終native收據已綁431C／50raw、三筆回鏈與Go proof；39套件冷測cached0。25單項21通過／4既有失敗[2,3,7,8]／新增0，phantom277無增刪；完整check的39套件cached，停點保留既有文件缺口。
- 首輪專案檢查新增1筆外部cpu.go相對引用，改成固定版本GitHub連結後用相同命令重跑通過基線；沒有放寬檢查。
- 當輪輸出UID/GID1000，root-owned維持27個舊路徑，無異常.md目錄，驗證容器已清理；依授權commit與push，完整Goal保持active。

## 2026-10-10：完整交戰入口 C 研究

- 從已發布431C／50raw接續野戰與攻城入口；先以[IDA探針](tools/ida_engagement_probe.py)核對尚未還原的直接呼叫閉包。原版bytes與輸入保持唯讀，完整軍團與主排程仍待。
- 原始直接與九張固定表閉包為234named／12raw／7817有效指令／18563碼bytes；C原型完整生成，原始guest stack與七個live patch保留，正常及20錯版語法檢查通過，尚待完整動態矩陣。
- 12條舊IDA錯邊界由31條新／替換指令取代，source-map及舊補充不覆寫；25097／57101／67099整檔組語exact、三個錯版拒絕。匯出器整合時遇到來源欄位及變數遮蔽問題，依實際schema修正後用相同命令重跑通過，未改產品行為或放寬guard。
- 原版道路建圖bounded已追到malloc讀buffer1154與手動arena1200重疊，壓縮檔10AC4偏移與RAM逐byte一致；以真AH4A保留到9B00後四劇本完整真退卻原版/C通過，城市座標假說已由primary否定。
- 暖機分成戰術A156前與完整世界還原兩種，各側各自執行及比較，不cross-copy Snapshot。設定選單需world前置，保留全部六鍵案例後修正前提。
- 目前擴充抽樣304例／604階段通過，194／234named與10／12raw已進入；完整444例與剩餘存檔／分支正繼續核對，未標全C或Goal完成。
- 後續580例／1156階段原版與C全部相同，234個named均進入、12個raw有11個進入；存檔抽樣12筆與地圖還原兩側各222次通過。舊收據及來源清單另存smoke-580-pre-report-fix.json與smoke-580-c-source.sha256，不拿它代替最終完整矩陣。
- 我在smoke執行中修改了外層wrapper，原生PASS後shell因讀取到修改後行段而報status未定義；此為編排錯誤，沒有原版/C差異。後續執行前凍結wrapper，O2／O0各用獨立有界輪次。
- 報告曾無條件沿用actual_nonlocal_exit=true，但實際outer退出為0。改由outer計數衍生，戰術19FDC→11B76的SS/SP還原另以同案例真11B5A入口觀測計數；warmup暫停與真正返回另外揭露。
- 垂直尋路初態誤寫D2FE，改為原1BE10實讀的D2FC後四劇本／八階段相同，raw1BFBF進入四次。凍結來源後啟動完整O2；本次先提交整檔組語exact及READY研究工作區的checkpoint，C台帳不升級，O0／O2最終矩陣、二十錯版與專案檢查仍待後續收據。
- checkpoint `4ac06c1` 已推送。其後O2在654例／1302個已接受階段時中止，未產生完整O2.json。host session退出143，Docker事件記錄容器執行1772秒後遭signal9、退出137並移除；沒有OOM事件，尚未達原2400秒期限。外層程序提前終止的上游原因未知，不能把進度檔當完整通過收據。
- O0與前14個錯版已預編譯，尚未執行完整驗證。改用容器內有界期限及持久退出紀錄的長工作入口，避免host等待程序結束時刪除尚未完成的原生矩陣；原C、Go、seed、矩陣與50M原版指令上限維持相同。
- 新入口的背景grid四例／八階段通過，job退出碼0且容器自行移除。原樣supervisor尾段的兩個獨立探針分別得到正常0與期限124，收據為`workplace/matching-decompilation/c-engagement/checks/supervisor-probes-1791584303468.json`，SHA-256 `e39b66fce4b83cd4b13d1af839302237c415d82b73d84d52f74c994a1628b31b`；探針只替換payload，不代證遊戲矩陣。
- 六項語意索引更新後，補充收據仍記舊SEM來源雜湊，嚴格source-only因此拒絕。重新執行原pinned工具的7817指令獨立組譯，再匯出、綁定及冷重建後source-only通過；未直接改歷史hash，舊收據另存私有code-proof，163份編譯來源仍為同一digest。
- 第一批錯版1–9由狀態比較拒絕；第10個在28例／56階段仍通過，證實原負例未觸發BE33差異。其三筆旗標同樣誤寫D2FE。四個matching_engagement Go驗證檔集中改由各側真初始化的CS:D2FC解析FFFD，runner共用既有patch入口；C演算法、錯版定義、seed及案例數不變，正式Go引擎不變。舊十組收據與bc37來源清單另存before-d2fc-fix前綴，新的二十錯版與兩級完整矩陣均須重跑。
- 新來源清單SHA-256 `05ce526895dcda6d13ae9b2ccd842e6922410c1601ac2a8fa282de4f6a5fd6e8`的正常vertical-route四例／八階段通過。二十個O0實編譯錯版全由狀態差異拒絕，逐ELF/buildinfo、macro id、source digest、退出1與原版／C前置digest核對通過，無panic、unsupported或SIG；收據`workplace/matching-decompilation/c-engagement/checks/engagement-controls-05ce526895dcda6d-999116.json`，SHA-256 `191cb676bcffb648741f96c31935636f7195cf65fa86cc1dc7cce29fdc5e7db8`。開始以持久背景工作重跑O2，O0與C台帳升級仍待完整收據。
- 同來源O2完整748例／1492階段通過，234個named及12個raw皆有實際入口，16組完整SAVE與兩側各256次MMAP還原通過。256次真戰術frame還原與outer退出0分列；1492個未outer轉移階段含508個warmup暫停，不能全稱真RET。收據`workplace/matching-decompilation/c-engagement/results/O2.json`，SHA-256 `3e3f7fdc1c50b839795d20e4122ce1e351a53532c80fd8ddf04ada17408d4eb1`。開始O0背景工作`wolong-c-engagement-normal-o0-1114652-29032`，C狀態仍READY。
- 下一軍團切片可沿用c-tick/census的原始probe `937f8738225ae0f9e5efc133245316dd075d845ffded97a22f5f211b4804f2a7`與DB `4f01452cc2df29a9306aa4b444d1d4125e3f27f7730fcfdadbefd9f9da3a4667`，16個導航函式加直接依賴1562B共434指令／1068bytes已逐chunk核對KI原bytes。此為後續RE入口，未宣稱新C完成；仍須審查間接分派，保留126FF→12708與12804→12808兩條落下邊。
- 同來源O0完整通過，收據與O2逐byte相同。最終驗證核234函式／12raw／7817指令、二十狀態差異、四筆回鏈及665／62台帳，公開收據SHA-256 `d575c695a8bfc1e3854f82ebea9797ec4f04f76a359dd30dc3a8ac1e9bed1fff`；Go proof `0c270b534c6469f410c9781bf88bd074d6a9ff921eaacc46310c6d5d50ebcc9d`已重綁。2599筆全台帳來源／介面／fixture／產生檔綁定逐筆吻合。
- 25項檢查21通過、4既有失敗[2,3,7,8]，phantom277正規化零增刪；完整check的Go vet及39套件通過並使用cache，停於既有文件問題。原39套件冷測cache0另存，沒有把暖測當冷測。檢查收據SHA-256 `6e7d5f53ec3870bffbaee204e65fcb7ec1be876ca826eace6dfd9add7fe3ab8c`。初次新增兩筆歷史hash缺路徑，核真v3檔後補路徑並同命令重跑，未改checker或基準。
- 分流舊1047筆分類完整保留，新增四筆局部C不代證的邊界後為1051；strict與索引已重生。README現況數字同步生成索引，沒有用未解列數當完成率。原版資產、RAM、ELF及完整EXE仍留本機。
- 本輪輸出UID/GID1000，root-owned維持27個既有路徑，無異常.md目錄；原生與檢查容器均已清理。依授權提交推送本輪CONFORMED來源與證據，完整軍團／主排程、正常玩家長程及原C機器碼仍未完成，整體Goal保持active。

## 2026-10-10：軍團輪轉 C 研究

- 從已推送7969082的665C／62raw接續軍團更新。沿用固定c-tick/census原始IDA證據，17函式434指令／1068bytes、兩個內部落下邊、十個已CONFORMED外部callee均核對，沒有新間接CALL／JMP或MZ重定位。
- 新建並索引re132與spec252 DRAFT。保留16格一批、尾批4200、軍費兩次右移各自截斷、士氣byte回繞及原順序碰撞，不把Go的127常駐軍團模型套回原始C。
- 尚待初態／獨立模型審查與原生C驗證；C台帳不提升，正式Go引擎不改。
- 窄契約審查後spec252升READY，17函式靜態C及十二錯版已產生；固定來源、兩份獨立組語切片、乾淨重生與十三種編譯變體通過，尚不等於動態驗收。
- 新Go/C回呼最初放在含完整C定義的preamble旁，造成cgo重複符號；移用既有平台回呼後，同命令重跑。其後補接原先未涵蓋的真`sub_1563B`扣款callee及唯一RET適配，軍費／士氣窄矩陣32例全狀態通過。
- 第一批完整直接入口矩陣停在停戰helper的原版first模型。fixture誤把道路點4004放進BX，原`0x142AD/0x142B9`要求道路link800。修正caller ABI及模型顯式輸入，不改原C、seed或停戰門檻。未完成收據不當作通過。
- 補正spec178及機制40的軍費截斷、路段判準與士氣byte回繞；原碼契約與正式Go排程／新C驗收分列。
- READY提交前25項檢查21通過、4既有失敗[2,3,7,8]，phantom277正規化零增刪。完整`tools/check.sh`的Go vet與39套件通過，39套件使用cache；其後停在既有文件缺口，完整log SHA-256 `f1a31f020f6e240bc3a7a05d4ea265faf5aeca4898e3abaefcb0720f323d7019`，已保存在`workplace/matching-decompilation/c-army/checks/ready-checkpoint-army-full-check.log`。
- 分流保留原1051筆內容及順序，只追加四筆新研究限制，主表1055／平台6／總數1061。索引已重生，README規格238份為232CONFORMED／5READY／1DRAFT。
- 軍費32例、舊來源manifest／ELF／buildinfo與停戰首差紀錄封存於本機`before-five-links`前綴；它們只作歷史範圍，不代替後續工廠收據。
- 五個根更新／移動／到站接點補入後凍結研究來源，O2有界工作`wolong-c-army-normal-o2-1320903-2151`已啟動；本次提交保存READY原型與來源證據，不把進行中矩陣或單一最佳化當作完整C驗收。O0、十二錯版及真戰鬥／outer接線仍待，C台帳維持665／62，整體Goal保持active。
- `1320903-2151`退出2，停在停戰獨立模型。根因為模型把端點絕對偏移6／8做XOR2，8變成10；原`0x142E7/0x142EA`是DI0／2先XOR2再加6。只修模型順序後，工作`wolong-c-army-normal-o2-1328798-32638`退出0，560例／588階段全狀態一致，全部17原始入口均進入，14108個獨立檢查通過。原C、門檻、seed與工廠未為此改動。
- 560例收據`workplace/matching-decompilation/c-army/results/direct-560-pass-O2.json`的SHA-256為`30e0c0fe68a1bcc66484198acc9ccc99adacbf664fcd4bc0be02b29a65b09858`，完整167份來源文字、ELF、buildinfo亦已封存；索引`workplace/matching-decompilation/c-army/results/direct-560-pass-archive.json`的SHA-256為`8b4a9ff8d8ee3c32032f6d7534ad2b16b9a8f8794ffc907b4b139a1a9e096bbd`，來源`workplace/matching-decompilation/c-army/results/direct-560-pass-sources.tar.gz`的SHA-256為`ae1bfbdf69b6beaa0c9faf93af50d11b2bc1fa600e9155c130a5765205297446`。此成功範圍尚不含真戰鬥或outer退出，不能代替最終驗收。
- 原生入口新增獨立verify模式及每job編譯來源文字封存。暫存Go副本先與manifest核對，再保存來源tar與SHA，最終驗證器逐檔核對封存內容；避免後續改檔讓舊收據只剩雜湊而無法重現。
- 野戰／攻城16例32階段通過，含16次前置及8次內層真戰術frame還原；四劇本真末城退出4例8階段通過。額外的`125A3→12662→12708→12880→14ADE→11CB1`根退出1例2階段通過，沒有續收軍費、倒數或寫回cursor。零城兵必勝分支由原碼證實，沿用固定seed3，沒有重擲挑結果。
- 最終凍結來源manifest為`workplace/matching-decompilation/c-army/results/c-source.sha256`，SHA-256 `ddd23c55db576cd820c73549fbe72e9122600b10b35b74730562dd6788d186d4`。完整工作`wolong-c-army-full-1398126-31183`的O0／O2各581例630階段均通過，兩份`workplace/matching-decompilation/c-army/results/O0.json`與`workplace/matching-decompilation/c-army/results/O2.json`逐byte相同，SHA-256均為`d2fd454ef87cbef9d71b0d5f66b5f7a9be869fed50fc839f77356e3979d7b851`；十二錯版尚在同一工作序列核對，未因此提前提升台帳。
- 完整工作最終退出0，十二錯版全部由實際狀態差異拒絕。嚴格驗證器核對14個ELF、167份封存來源、七回鏈及完整矩陣通過；`docs/re/c-army-verification.json`與本機收據逐byte相同，SHA-256 `6a609c92de9d420ed7036eb619504b76798e17c3a2304397133f4d7d07137624`。spec252升CONFORMED，台帳682C／62raw的2667筆來源綁定全部吻合。
- 最終25項檢查仍為21通過／4既有失敗[2,3,7,8]，phantom277正規化零增刪；`workplace/matching-decompilation/c-army/checks/army-checks.json`的SHA-256為`a2a02f501ca2d6ad6c6b91b568810b5675e66c88be81757e55c562367a06bcf4`。完整`tools/check.sh`的Go vet與39套件通過、cached39，後續停在既有文件缺口；`workplace/matching-decompilation/c-army/checks/army-full-check.log`的SHA-256為`78d7b56e6333fed80e98e069153ee33d1278204ae1cf00d41c45f5ad68770350`。
- 只移除已驗收的兩筆C矩陣未解列，其他1053筆分類與順序保留，另有6筆平台層；規格238份為233CONFORMED／4READY／1DRAFT。台帳提升後重跑原pinned GNU工具，434指令／1068bytes仍exact、來源ed6d與組語收據e399皆不變，沒有新增指令或改全域組語。
- 本輪輸出皆1000:1000，沒有新root-owned檔案或偽.md目錄，原生與檢查容器已清理。依授權commit與push本輪CONFORMED來源及證據；完整主排程、讀檔、新遊戲、正常玩家長程與原C機器碼仍待，整體Goal保持active。

## 2026-10-12：開機、主排程與命名UI C 研究

- 42函式＋3 raw片段、1041指令2520 bytes，O0／O2各217例217階段逐byte同收據，十三錯版全由實際狀態差異拒絕。
- 九案開機導覽：解出18DC8二選框幾何、18B7C／18E5A／18FC9槽位選單塗漆座標、18412處理器選單右鍵取消鏈；0x1005B直接進入因fixture與真實呼叫點差約300位元組退役，改劇本多樣性。
- transfer群以pollFeed把0x222B raw等待做成檢查點餵食（右鍵直通INT33 AX=5讀後清零，兩則訊息需兩次右鍵）。
- mainloop的a5=1／d2d≥A0h與transfer的d2d<FFh因子需IRQ0改測就緒邊界；12216補丁路徑無呼叫端。
- 工具鏈：-O0補function-sections＋gc-sections修未定義引用；稽核流分叉改走mismatch exit 1（保留非確定性panic護欄）；verify fallthrough放寬已知入口；報表補new_routine_sha256鍵。
- 三筆突變重設計為值變（m1離場碼、m2訊息位址、m12視訊模式）；C台帳724／65。
- spec253 CONFORMED；公開收據docs/re/c-startup-verification.json。

## 2026-10-12：殘部群 C 研究（census 收工）

- 位址差集最後 10 支＋raw 尾跳 1AB48（247＋1 指令 621 bytes），O0／O2各89案逐byte同收據，十錯版全由實際狀態差異拒絕。
- 名子差集誤報：nullsub×2 與 VGA×3 已在帳上是 code_blocks；位址差集為準。
- picker 熱區座標全十六進（208h／240h）；15B51 以 cs:0D52 蓋 DS，entry-DS 另給 0x8300 段；prePress 不足，pollFeed 一律開。
- 15E1E→15F5D→1062F 金額列印 16 位元首除，值 ≥655360 即 #DE；15AD1 案只收右鍵離開，選列由 15AFC 直入案覆蓋。
- m7 只改比較是語意空操作，改突變存值才拒絕；C台帳734／66，census 739／739 全覆蓋。
- spec254 CONFORMED；公開收據docs/re/c-residual-verification.json。

## 2026-11-02：工作區清理與全滑鼠桌面工具列

- 使用者指定順序為完成 remake、HD、手機版，並指定參考 `~/cht/rich2` 工具列，明選全程滑鼠操作，只保留文字輸入。現行工作回到正常玩家路徑，既有 C 研究作查證來源。
- 基線 `bd09fd40e344efee9d01e286a70f2a2303890caf`，開工工作樹乾淨。本輪未 commit、push、改 tag 或製作發行包；現行 GitHub Issues 已透過主機 gh 核對，沒有寫入遠端。
- workplace 清除 44 個 `matching-decompilation/*/gocache`，54,299 檔、4,843,557,786 bytes。27,713 個其餘檔案逐檔內容前後一致，彙總 SHA-256 `ced8b2aa4d934c72f8bab513a91eb3e781ba8550830875241014f6e15b9c0aa6`。原版、IDA DB、封存來源、收據、存檔及截圖保留。
- 個別修復 24 個已確認的舊 dosbox 截圖／log 與 manual 衍生圖擁有權，逐檔內容雜湊不變；刪除 `manual/workplace/parity` 與 `manual/workplace` 兩個空目錄。workplace 已無 root-owned 檔案或偽 `.md` 目錄，沒有遞迴 chown 儲存庫。
- spec266 經使用者決定建立 READY，實作獨立 40px 工具列，原版遊戲畫布仍為 640×400。語言、設定、存檔、讀檔、返回、確認、離開皆以實際滑鼠分派；快捷鍵與方向鍵停用。新的 UI 文字與滑鼠提示放在四語 `translations/desktop-controls.json`。
- 返回／確認送到目前輸入擁有者，補齊存檔、表格、編成、選點等取消入口；工具列擁有按下到放開的整段手勢。模態期間主機功能停用；單純停留在工具列仍依 idle gate 恢復世界，頂邊持續外推可向上捲動。
- 沿用既有 `san1-go:1.26.7-ebiten2.9.9-r1`，沒有建立重複 image。初次失敗來自快取唯讀缺依賴、image 預設非 root 令新 volume 初始化失敗、以及預設 GOPROXY=off。只修正 `wl-gobuild` 根目錄擁有權，透過 go.mod／go.sum 下載鎖定依賴後恢復離線驗證。
- `tools/go.sh test -count=1 ./cmd/... ./internal/... ./mobile/... ./translations/...` 的全部測試套件通過；同範圍 `vet` 通過。新增手勢隔離、模態邊界、返回一層、語系文字寬度與座標測試。桌面最後增補後另冷跑 cmd/wlgame 通過。
- 正常 GUI 重播器 `tools/desktop_mouse_smoke.py` 只觀測實際狀態，用 X11 按下／放開握手同步，不注入世界、時間或亂數。固定 seed 1、原版唯讀、存檔在容器可寫暫存，走正常首頁開始。
- 失敗收據保留為 scratch 的 attempt1／attempt2／attempt4／attempt5 前綴。首次編成的固定 90ms 輸入改為實際按下／放開握手後通過；未據此宣稱原版或產品規則有缺陷。另補正合法的母本 `uiLang=nil` 觀測值、清除舊 trace 並等待新視窗可見，避免腳本把前輪狀態當成就緒。
- 最終 GUI 30 項通過，收據 `workplace/scratch/mouse-toolbar-gui.json`，實際程式 SHA-256 `fa362568d528734a10205a474e8958e1c3556a402b35f9596ae78f6b311b0cc8`。包含四語切換、頂邊外推與停手、禁用遊戲鍵盤、設定／存讀檔、正常編成／行軍、拖出放開隔離、離開確認與取消。
- 軍團 16 由正常編成產生，目標宛 90，固定鏡頭 `(188,104)`；位置 `(206,114)` → `(204,114)`，朝向 4 → 0。`mouse-toolbar-march-start.png` 與 `mouse-toolbar-marching.png` 目視可見軍團離開城池向西移動。證明 remake 會移動，不代證原版動畫、時間、到站、完整戰鬥或完整戰役。
- 重新以 dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8` 跑原版 NEW GAME，生成 `workplace/shots/mouse-original-boot.png`／`.log`。實際唯讀工作副本在 `/home/anr2/cht/dosgolem`，預設 `dosgolem-wolong` 不存在；本輪未改 golem。這份原版新收據僅驗啟動，不冒充行軍同狀態比較。
- re84 原版 y 100–131 列帶選取與現行 `talkChoiceClick` 直接點文字列的差異，已由正常輸入重現。早期「行為等價」只限規則出口，不代表滑鼠手感。使用者尚未回覆保留直接點選或還原列帶的選擇，正式選單模式本輪沒有另作裁定。完整 remake、HD、手機與原生平台完成度維持未完成。
- spec266 升 CONFORMED 僅限工具列與上述正常抽驗；現況入口與 README 操作說明同步，失敗與歷史證據保留。Docker 工作皆有資源限制、外層逾時與自動清理；Xvfb／遊戲有界回收，最終容器與擁有權狀態另行確認。
- 最終文件索引 568 份／135 條斷言檢查通過。補上 spec 分區索引的 266 入口；playtest/19 兩個既有歷史圖片缺檔改為保留原位置與缺檔註記，不補造圖片證據。
- 本機 `mouse-toolbar-source.tar.gz` 與 `mouse-toolbar-source.json` 保存受測來源及 SHA-256；來源封存 SHA-256 `77052f84931db42724425cb8b07ceab2dce49dc1816a7acbcd2dc04ec35d3bdf`。GUI 後只有 desktop.go 的 gofmt，受測二進位保持原雜湊。
- 最終核對：GUI 30 項、受測二進位雜湊、新文件入口、檔案 UID/GID、原版 KI.EXE 雜湊均吻合；workplace root-owned 與偽 .md 目錄皆 0。專案與本輪 Docker 容器均未殘留。

## 2026-11-03：原版上下推動反白

- 使用者明選原版上下推動反白，承接 spec266 全滑鼠外部工具列，不再等待相同操作取捨。依復古遊戲、GUI 試玩、共同決策與文件職責路由載入既有入口，spec267 先 READY 再實作。
- 檢查 grill-me 上游：入口最後提交 `fcf0071560d32913c9d4f820e0d7ca467c881619`，2026-08-15；grilling 核心最後提交 `7739eaed7f8f3915437a8def198a6d0d7206dbc7`，2026-10-29，內容為子代理使用全新上下文。以 GitHub API 查證，收據 `workplace/scratch/grilling-upstream-check.json`。本機是繁中獨立改編，沒有正式版號；未覆寫技能或上游來源。
- 沿用 re84／re112、spec232 與固定 IDA 9.4 DB 的已證實原始 selector，不新增反組譯切片或修改 DB。新 `choiceband.go` 接回文字指令、進言理由、行軍模式、重疊選單、外交、撥款與原生結束選單，保留各自規則出口。
- 文字選單保存並隱藏游標，初始 `(0,100)`、反白首列；y 超過131下移一列回100，低於100上移一列回131；首尾夾住。左鍵在框外也確認，優先於位置處理；右鍵取消後還原。表格、存讀檔、編成、YES／NO 沿用各自滑鼠處理器。文字選單獨占輸入時停用工具列，離開後恢復。
- `tools/go.sh test -count=1 ./cmd/... ./internal/... ./mobile/... ./translations/...` 全部通過，同範圍 vet 通過。選列器測試涵蓋兩列／四列／五列、邊界與端點、停手、確認優先、游標保存恢復、子視窗返回與工具列擁有權。
- 正常 GUI 首次失敗是端點選列結果與前一個狀態相同，腳本未等候新位移就送下一步；加入唯讀 `pointerRaw` 觀測並等待實際位移。第二次失敗是腳本把 X11 右鍵3誤當觀測陣列第三格，trace 已顯示取消與還原成功；改為明確對映左1→0、右3→1、中2→2。兩次收據保留於 `workplace/shots/selector-20261103/gui-attempt1/`、`gui-attempt2/`，沒有放寬產品判準。
- 修正腳本後以同一 Docker 工具鏈乾淨重跑，41 項通過。固定 remake seed 1、正常首頁至開新遊戲、四語、設定、存讀檔、編成、軍團上下反白、框外確認、右鍵還原、行軍模式、移動與離開取消。沒有 direct-entry、強制行軍、世界注入或自動關訊息旗標。
- 受測程式 `workplace/scratch/wlgame-original-band`，SHA-256 `9ccc5a0f9338076f81918e5fc65bc4531ce417038c8dcad69b7c40608c1e1cd4`。軍團16下令前往宛90，固定鏡頭 `(188,104)`，位置 `(206,114)` → `(204,114)`，朝向4→0；兩張真窗截圖已目視確認軍團由城池往西移動。只證明 remake 正常路徑移動，不代證完整原版動畫與時間。
- 原版透過 dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8` 正常第一章／曹操開局重跑，無直接注入。文字選單送 y131、132、399、99、0、132；每次之後推進300,000道指令。日期維持196/4/1 1時，滑鼠回置與既有原始 selector 一致，框外確認後進入行軍清單並還原 `(240,47)`。未涉及亂數判定的雙側比較，不宣稱兩側世界或亂數序列相同。
- 真窗截圖以 point 縮放還原原生640×400，軍團選單 `(192,56,112,56)` 第一列／第二列各6,272 pixels與原版全同；整幅背景與軍團狀態不同，不宣稱整幅同狀態 parity。原版／remake 選單及移動截圖均已目視。
- 最終收據、原版完整時間軸、區域比較與來源封存集中在 `workplace/shots/selector-20261103/`，入口為 spec267；前一輪工具列收據與二進位保留。spec267 升 CONFORMED，現況表、README、spec266、re84追加勘誤及政略機制同步。完整戰役、原生平台、HD及手機完成度維持未完成。
- 封存357份受測來源並逐檔回讀核對；`source.tar.gz` SHA-256 `25a2627d8837807305b24f1142947f6db685594e219034274e90934e6c325e41`，明細在 `source-manifest.json`。文件索引569份／135條斷言檢查通過，git diff --check通過。原版KI.EXE雜湊不變，workplace root-owned與偽.md目錄皆0，本輪及專案標記的Docker容器均未殘留。沒有commit、push、發行或修改技能。

## 2026-11-03：接回進言表格的滑鼠操作

- 接續使用者確認的全滑鼠操作與原版上下推動反白。本輪核對主機 gh 登入及28個現行 open Issues，只讀取遠端，未修改 Issue 狀態。依 remake、GUI 還原、規格閘門、文件職責與驗證路由載入現有入口，沒有重新訪談已確認決定。
- 發現 `updateAdvise` 的遷都、協助勢力、目標勢力選取階段在共用 `updateListUI` 前返回，仍只讀取已停用的鍵盤或工具列確認。修正前版本加入唯讀狀態觀測，正常首頁開局後選遷都，點 `(150,128)` 仍停在 Cursor0／Browsing，失敗收據保存在 `workplace/shots/advise-mouse-20261103/before/`。
- 重生原版正常第一章／曹操開局的遷都收據。第一次先後點不同列，第二下已直接確認原先反白的洛陽，後面帶 clear／cancel 名稱的圖其實是對話，沒有當成取消證據。已保留至 `oracle-attempt1/`，依原始 C 及 re26／re118 的 `sub_18463` 回查，確認第二段只等待左右鍵、不重做座標命中。
- 原版重新分支重跑：列表點洛陽反白、右鍵解除、再次右鍵返回地圖；還原同一次正常流程抵達的完整機器狀態，再點同列兩次進判決，右鍵從君主首句推進至軍師第二句。每步使用指令數與控制流，未注入世界或直接進函式。原版日期維持196/4/1 1時；完整時間軸與 log 在本輪收據根目錄。
- spec268 READY 後，進言表格接回既有滑鼠分派，選取結果由回呼接入原規則，避免重複 Confirm 或把排序後列號當勢力 ID。保留外交官與使者閘門；閘門失敗保留選列，協助方確認後的新目標清單從未反白狀態開始。
- 共用表格反白後左鍵於任何遊戲區位置確認原列，右鍵先清反白；反白等待時隱藏遊戲游標，子視窗保留自己的游標。遷都列表第二次右鍵依原版返回地圖。進言說服對話補上右鍵推進目前一句，不同時取消或選中背景。遊戲規則、存檔結構與亂數均未更改。
- 新增五個回歸測試，涵蓋原列確認、遷都取消及資料 ID、排序後勢力 ID、外交官閘門／子表格、右鍵只推進一句。正式範圍 `tools/go.sh test -count=1 ./cmd/... ./internal/... ./mobile/... ./translations/...` 及同範圍 vet 全部通過。
- `tools/desktop_mouse_smoke.py` 新增可選 `advise` 模式，呼叫方式為既有四參數後加 `advise`；default 模式仍驗編成與行軍。兩條流程都使用真實 X11 按下／放開與實際狀態握手，不注入遊戲狀態，固定 remake seed1，原版唯讀，存檔位於容器暫存目錄。
- 進言 GUI45項通過：正常首頁、四語／設定／存讀檔、遷都非首列反白、取消後首都不變、確認洛陽58並演完拒絕判決、外交清單捲動／排序後選劉璋3、右鍵推進及退出、再次存讀檔保持首都82。同一二進位的編成／行軍GUI41項通過，軍團16目標宛90，固定鏡頭 `(188,104)`，位置 `(206,114)` → `(204,114)`、朝向4→0。
- 受測二進位 `workplace/scratch/wlgame-advise-mouse` SHA-256 `ea861393c13eaae4c2a91699abc4562c7b0361cf5041ed8740405d86967fb049`。原版與 remake 的代表清單、判決、提案及軍團移動圖已目視。兩側背景、部分城市數值與操作提示未對齊，本輪不宣稱完整進言全畫面 parity 或整體 remake 完成。
- 受測357份來源封存為 `workplace/shots/advise-mouse-20261103/source.tar.gz`，SHA-256 `3c690f54c520d2a60dfedef44094768682ba684a4c7eb493c38df7637c4aedec`；逐檔來源、二進位及全部收據雜湊記在 `source-manifest.json`。本輪收據及 `regression/` 不覆蓋前輪 selector 或工具列收據。
- spec268升CONFORMED，現況表、spec分區索引、spec38輸入勘誤、政略機制及README操作說明同步。完整事件、完整戰役、原生平台、HD與手機驗收維持原有範圍；本輪沒有commit、push、發行或技能修改。
- 收尾核對：文件索引570份／135條斷言通過，來源封存與工作樹逐檔一致、兩條GUI收據綁定同一二進位。原版KI.EXE雜湊不變，抽查修改檔UID/GID皆1000:1000，workplace root-owned與偽.md目錄皆0；本輪Docker容器已清理。

## 2026-11-03：和平折返同拍移動與原生行軍存檔

- 接續使用者已確認的全滑鼠與原版上下推動反白，先完成 remake。開工沿用唯一現況表，核對主機 gh 登入及 Issue #30，只讀遠端。本輪沒有 commit、push、發行、改 Issue 或修改技能。
- 依 remake、規格閘門、GUI試玩與文件職責路由載入現有入口。Issue #30 引用的舊「攻宛沒到站」已由 re86／spec132 裁決為和平折返，沒有重新開啟為行軍故障；本片只修查證後的同拍與存檔差異。
- dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8` 從正常首頁、第一章／曹操、編成及攻宛操作抵達和平邊界。原版 `KI.EXE` SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；固定 IDA9.4 linear 匯出 `c-army/ida-source.json` SHA-256 `ed6d851a394023ab58809aa504391bd20e7c649aa721ac57b8ef99ee925c26b7`，99道相關指令與原始EXE bytes核對相同，未修改資料庫。
- 原版196/4/3 5時，軍團37於 IDA `0x142EE` 的 `(190,119)`、目標90、步進4，108道指令後於 `0x126FE` 已回舞陽 `(191,119)`、目標／意圖87、步進`0xFC`。路徑點仍`0x45CC`，目標座標仍`(181,127)`。以 `DS=2754`、`SI=2B80` 的明確段讀取64-byte記錄，來源及原始欄位保存在 `workplace/shots/turnback-20261103/original-turnback.json`。
- 原版 `attempt1/` 保留多送一次右鍵造成的腳本失敗；`attempt2/` 已取到移動，但未指定段的peek讀到程式碼。現行收據使用 `peek:2754:2B80:64` 乾淨重跑，不引用錯誤記錄。取樣依指令與遊戲時間，不用牆上秒數；此局部比較不涉及亂數條件，不宣稱雙側完整亂數流一致。
- spec269先READY。原版取樣的兩個回歸先重現舊程式多停一拍、步進量未反向及目標座標提前改寫。正式道路圖的折返改到 `tickOneCorps` 的重算／碰撞前，交由 `replanOnLeg` 消費旗標並於同拍繼續走。無素材的直線降級保留，主君編成預設仍依 spec76 的既有選擇。
- 完整GUI第一次在行軍中存讀檔後仍往宛走，收據在 `gui-attempt1/`。確認原生 `Snapshot.Routes` 有格子即跳過恢復，出發節點留成目標90，路徑標記缺失，跨舞陽後連結位址一直為0。spec269補上快取恢復契約，再先以原生Snapshot回歸重現。
- `Restore` 清除舊標記；`restoreMarchRoutes` 對既有兩種成本規劃的候選逐格核對已存路線，才恢復出發節點及標記。無法吻合時由衍生的 `routeUnresolved` 停止移動並計數，不修改存檔格式或默默改路。初版用標記長度直接擋行軍，誤擋既有合法反向段；反向存檔回歸指出後，改成只標記實際恢復失敗的軍團，狀態層全測通過。
- 正式範圍 `tools/go.sh test -count=1 ./cmd/... ./internal/... ./mobile/... ./translations/...` 冷測39個有測試套件通過，同範圍vet通過。後加的原生存檔跨城取樣另跑目標測試通過；共出發前、路段中、跨城後5個取樣逐格與逐拍比對，另有2個損壞快取反例及原版同拍／反向存檔回歸。
- GUI第二次停在設定工具列，trace顯示相對位移尚未消費便被重送，實際點擊落到遊戲頂邊；記錄保存在 `gui-attempt2/`。腳本改為每次位移都等候唯讀觀測更新後，用同一二進位、同一條件乾淨重跑，沒有修改遊戲來配合腳本。
- 最終GUI51項通過，固定remake seed1，正常首頁、四語、設定、原版上下選列、編成、行軍、途中存讀檔、和平折返、到站存讀檔與離開取消。透過既有設定關閉主君編成並選最高戰略速度；不注入世界、時間、軍團或改用direct-entry。軍團37從 `(206,114)` 出發，途中存於 `(197,118)`，載入後節點82、連結`0x0F30`、路徑點`0x436C`與目標座標相同；4/3 8時觀測到停舞陽87，步進`0xFC`、路徑點`0x45CC`。後續4/3 10時正常重新指派到站座標發生於再次存檔面板之前，並非讀檔覆寫。
- 受測程式 `workplace/scratch/wlgame-turnback` SHA-256 `8c23c1f0bc0a625a6b9de8bccdbecbaea274dd17384b2e5d38d1e09c04c40049`。原版邊界圖及remake移動／到站真窗圖已目視；兩側世界、鏡頭與操作耗時不同，原版108道指令未必重畫，不宣稱全畫面或wall-clock parity。
- 沿用 `san1-go:1.26.7-ebiten2.9.9-r1`，原版唯讀，GUI存檔在容器暫存。357份受測來源封存在 `workplace/shots/turnback-20261103/source.tar.gz`，SHA-256 `f353dcbc4ed259e9a9a5f6801970a00ca40583bc6e0734d86f1e8710c4a0fcc9`；逐檔回讀一致，來源／52份收據／二進位／原版雜湊見 `source-manifest.json`。
- spec269升CONFORMED，現況表、spec132與172、軍事機制及規格索引同步。只完成和平折返及上述存檔抽驗；特殊舊快取與現行成本路線不吻合時會停止並記錄。攻方／守方、碰撞與完整戰役仍由Issue #30管理，不宣稱完整remake完成。
- 整個工作根的首次擁有權稽核發現一個舊 `tools/__pycache__/dosboxx_probe.cpython-312.pyc` 屬root。只修正該單一檔案為1000:1000，前後SHA-256同為 `1804e8309e1dc4fc7db21ba2e5c171fceea4544b1a5fa5ef9cd54f82f69dc449`，未遞迴更改其他檔案；再跑完整find檢查為空。
- 收尾核對：文件索引571份／135條斷言及 `git diff --check` 通過；來源與收據逐檔不變，原版KI.EXE雜湊不變。修改檔抽查UID/GID皆1000:1000，整個工作根root-owned及偽.md目錄皆0；專案標記與本輪具名Docker容器均無殘留。

## 2026-11-03：戰術命令與選隊的點擊邊界

- 接續使用者確認的全滑鼠與原版上下推動反白，保留上輪和平折返修正。依remake、規格閘門、GUI還原、驗證交接與文件職責路由載入現有入口；核對主機gh登入、28個現行open Issues及Issue #30，本地worklist verify作輔助，沒有寫入遠端。
- 依AGENTS的簡易平行工作規則，委派一名唯讀代理檢查戰術滑鼠與原版取樣入口。邊界明列不得改檔、其他repo、原版、技能、commit、push、重編或清理；只使用既有Docker。代理未修改任何檔案。
- 發現 `battleHitRect` 四邊各內縮2px，讓側欄六命令及底列六選隊格的邊線點不到。re60已證實側欄每列48×16、底列每格80×32，playtest69保存歷史20個熱區紀錄。先建立spec270 READY與索引，再改用原版矩形的全畫布回歸；舊程式首先於 `(496,280)` 失敗。
- 增加唯讀桌面觀測：玩家、交戰、城市歸屬及戰術幀、側別、選取位元、陣線、常令與大將命令，用於正常GUI同步。修正前程式正常首頁編成守軍37，196/6/6 5時由軍團81攻許昌82，點 `(0,368)` 仍無選取，收據保存在 `workplace/shots/battle-hit-20261103/before/`。
- 正式實作只將 `battleHitRect` 改用完整正尺寸格子；移除不存在的2px內縮，修正右欄舊2×3註解為單欄六列。不改繪圖、命令碼、選隊排列或戰鬥規則。回歸依原版常數檢查258,084個畫布內外座標，各檢查側欄與底列，另保留相鄰陣形與換位表測試。
- GUI `gui-attempt1/` 使用最低戰術速度，45秒內只走到第11幀，未越過開場；改由同一設定入口選最高速再跑。`gui-attempt2/` 在相鄰x479／480測試被腳本1px容差誤點，改為邊界操作要求精確座標，沒有修改產品來配合。
- `gui-attempt3/` 的退卻檢查誤用 `Standing`。起初懷疑輸入太晚，但trace顯示第444幀已點擊、第560幀才完成；規則層明定退卻不改常令。改觀測大將實際 `Cmd/Next` 後確認5/5，而常令仍1，沒有為測試更動規則。
- 最終正常戰術GUI44項通過，固定remake seed1、正常首頁與編成守軍，沒有open-battle、世界／時間／軍團注入或強制勝負。涵蓋第一格與末格四個角的選取／解除、柱與空隙不命中、守陣左上角、攻擊右下角、鄰接陣線、退卻、返回戰略及存讀檔。退卻後許昌歸勢力13、首都移至洛陽58，存讀檔保留城市歸屬與首都；這只是remake正常結果，不代證原版結算。
- 同一二進位另跑行軍GUI51項通過，涵蓋上下推動反白、編成、行軍中存讀檔、和平折返及到站存讀檔。`regression/attempt1/` 的最高戰略速度令存檔時已在中繼城，連結位址為0，不符合行軍中取樣前提；改透過既有設定選普通速度2，保留非零連結位址的原判準，沒有放寬斷言或改seed。
- 正式範圍 `tools/go.sh test -count=1 ./cmd/... ./internal/... ./mobile/... ./translations/...` 的39個有測試套件通過，同範圍vet通過。最後增加大將命令的唯讀觀測後，另冷跑及vet `./cmd/wlgame` 通過。遊戲與原版素材均在容器內執行，原版唯讀、GUI存檔在容器暫存目錄。
- 原版 `KI.EXE` SHA-256仍為 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`，dosgolem工作副本 `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8` 唯讀。本輪嘗試重生原版熱區，但未取得成功新收據；不將失敗腳本或remake圖登錄為原版對拍。
- `oracle-attempt1/`、`oracle-attempt2/` 的日期條件取樣分別用40M／200M預算未達成；隨後回查路由與既有腳本。有限profile證實主循環在執行，60M指令後日期為4/4 12時、再1M至13時；120M取樣到4/8 4時。不能由先前的until失敗宣稱時鐘停住、訊息擋住或oracle整體不能跑。
- `oracle-attempt3/` 只有軍團19、37、72，舊腳本預期的35未出現，`siege:37,76`未達戰場。`oracle-attempt4/` 改用實際駐122的72，仍未達戰場入口，原因未知。未固定原版亂數作玩法比較，沒有以這些探索宣稱雙側seed或正常戰役相同。
- `oracle-probe/` 是由golem CLI複製的可丟棄介面建構探針，只寫本輪輸出目錄，未改golem或遊戲EXE。直接呼叫原始建構到預期註冊完成點亦未達成，不算熱區成功收據；保留程式與輸出供回查，停止猜測前置及增加等待預算。
- 受測程式 `workplace/scratch/wlgame-battle-hit` SHA-256 `bcd7904992f486eb64fe4513b0a02c78a3d0594a8a1ca303ab53dc647a2e64c0`，修正前觀測版本 `905e46a6629ced4e996ac8a4f2f6df174db9a3e24aa67038f244680efd10d438`。正常選隊、下令與戰後載入真窗圖已目視；不宣稱整幅像素、原版時間、攻方入口或完整戰役parity。
- 沿用 `san1-go:1.26.7-ebiten2.9.9-r1`，沒有新建工具image。358份受測來源及spec270封存為 `workplace/shots/battle-hit-20261103/source.tar.gz`，SHA-256 `9c0d8952340da989d6e1a94d7b42b6daf868393c796312c29b161607a93ca3c0`；逐檔回讀一致。112份收據、二進位與原版雜湊由 `source-manifest.json` 索引，`verification.json` 分開標示remake通過與原版新取樣未知。
- spec270保持READY，沒有升CONFORMED；現況表、規格索引、spec37、戰鬥機制及README操作說明同步。README舊的「另外加鍵盤」敘述改為已確認的全滑鼠外殼。原版新熱區、玩家攻方、逐拍碰撞與完整攻守仍由Issue #30管理，不另建待辦鏡像。
- 本輪沒有commit、push、發行或技能修改。整個工作根的root-owned檔案及偽.md目錄檢查為空，Docker程序均有逾時、資源上限與自動回收；最終索引、來源雜湊與容器檢查結果接下列收尾紀錄。
- 最終核對：文件索引572份／135條斷言通過，以已知spec269作正對照確認新spec270有入口；358份來源及112份收據逐檔雜湊一致，兩條GUI綁定同一受測程式，原版雜湊不變。修改檔抽查1000:1000，整個工作根root-owned及偽.md目錄皆0，`git diff --check`通過；本輪具名與專案標記Docker容器均無殘留。

## 2026-11-03：補齊原版戰術點擊熱區收據

- 接續spec270的原版取樣缺口。依復古remake、規格閘門、驗證交接及文件職責路由載入既有入口，沿用目前狀態表與Issue #30範圍；主機gh登入及現行Issue已核對，未寫入遠端。未重開已完成的選單、行軍與存檔修正。
- 依專案適用的平行檢查規則，唯讀代理核對原版攻城呼叫前置。發現任意CS傳入 `IDAIn` 有段回繞風險，擷取副本先以絕對位址等待 `sub_1E453`，再驗證CS與世界表。本次實際CS為0110、換算實體位址均0F553，沒有觀察到段錯置，不能把風險寫成這次失敗的原因。
- 固定IDA來源為 `workplace/matching-decompilation/c-engagement/ida-closed-v5/ida-probe.json`，SHA-256 `32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9`；相關入口指令的檔案偏移與bytes對照原版KI.EXE。位址採IDA Pro9.4 database linear、base0x10000；沒有修改資料庫、原始符號或其他專案。
- `resume/attempt1/` 停在開戰訊息。既有re133已證實 `sub_1222B` raw等待只讀右鍵；送出右鍵後，原版實際進入11B5A、19946及1C7A9。`resume/attempt2/` 在右鍵處理已越過入口後再次等待同一入口，誤報逾時；移除重複等待後成功。playtest72追加日期化勘誤，保留歷史觀測，撤回「所有開戰訊息都可等候自過」的外推。
- 最終原版取樣：正常第一章曹操開局及編成37，196/4/8 4時輔助 `siege:37,122`、守軍72，節拍6→56；298,798,791道指令讀得20個固定戰術熱區，另有開場對話碼27。28個邊界與鄰接座標查詢一致，640×400原版戰場截圖已目視。這是固定介面的輔助取樣，沒有代證正常玩家攻方行軍、亂數、時間、畫面或結算。
- 新收據在 `workplace/shots/battle-hit-20261103/resume/`。4000-byte熱區表SHA-256 `17b08e0147010a49e94f046564533872777aaa2f2c1423026c88164fc7fdc4bf`，原版戰場PNG SHA-256 `da18b0a9bcd16434287a7f7b1b84f09cffcce57a98e579efb625fd89e40740e3`。保存修改後的golem CLI副本、時間軸及有界Docker重播器，原golem工作副本仍乾淨且唯讀。
- `resume/replay.sh` 在獨立空目錄乾淨重跑exit0；熱區表、三張PNG與完整指令紀錄逐檔相同。重播使用既有 `san1-go:1.26.7-ebiten2.9.9-r1`、2CPU／2GiB／256程序／240秒限制、network none及使用者1000:1000；原版與golem來源唯讀，沒有新建image。
- 新增 `TestBattleHitTestsAgainstOriginalHotmap`，由 `WOLONG_TACTICAL_HOTMAP` 指定私人原版表；`tools/go.sh` 明文傳入此環境變數。測試逐點比較256,000個座標的命中與分派，側欄4,608／底列15,360像素均相同；缺檔、長度錯誤或空控制區不會通過。未指定資料時明確略過，這次正向與反例均實際提供原版表。
- 正向新熱區測試及既有完整矩形測試通過；用Go overlay僅還原舊2px內縮的反例，於 `(496,280)` 原版碼09、remake未命中而失敗，非顯示器或環境錯誤。另以實際原版表冷跑完整桌面套件通過，桌面vet及wrapper語法檢查通過。紀錄為 `hitmap-test.log`、`negative-test.log`、`desktop-test.log`、`go-vet.log`。
- 本次不改遊戲執行邏輯。上輪封存的341份Go來源逐檔相同；同一受測二進位SHA-256仍為 `bcd7904992f486eb64fe4513b0a02c78a3d0594a8a1ca303ab53dc647a2e64c0`，上輪44項正常戰術與51項行軍GUI仍適用。新增的測試、wrapper及來源封存另存於resume，不覆寫上輪READY封存與失敗收據。
- spec270升CONFORMED，限側欄命令與底列選隊的點擊邊界。現況表、規格索引、spec37與README同步；完整攻守、碰撞、完整戰役、HD及手機仍未完成。沒有commit、push、發行或技能修改。新來源與逐檔收據由 `resume/source-manifest.json` 索引，最終文件、擁有權與容器狀態接續本段收尾紀錄。
- 最終封存365份來源、30份新收據，`resume/source.tar.gz` SHA-256 `c2dd3be718fb1020ce6ff1c498547321bc487414b3ae90947bb4d99fb8d85cf9`；封存逐檔回讀一致，舊112份收據及READY來源封存不變。首次收據匯出漏傳Docker標準輸入，確認沒有產物後以 `-i` 重跑並取得明確的檔數與雜湊輸出，屬工具命令問題。
- 文件索引572份／135條斷言及 `git diff --check` 通過；以既有spec269作正對照，spec270索引與CONFORMED狀態一致。整個工作根root-owned與偽.md目錄皆0，原版及受測二進位雜湊不變，專案標記與本輪具名容器無殘留。

## 2026-11-03：正常滑鼠攻方行軍與攻城抽驗

- 核對目前狀態表、主機gh登入及現行Issue #30；Issue要求原版固定存檔、亂數、碰撞與正常攻守路徑，遠端仍保持OPEN，未寫入。載入remake、正常玩家路徑、驗證交接與文件職責入口；沿用既有Go／Ebiten image及上輪受測程式。
- 依AGENTS平行檢查規則委派唯讀代理查原始第一章、進言、遭遇與既有helper。SINARIO.DAT SHA-256 `21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`，確認宛90屬張繡21，開局敵對提案會先被君主拒絕；改選可由正常進言說服的袁術12／汝南100，原始鄰接82→87→100。代理沒有改檔、編譯、執行GUI或存取網路，未改其他repo。
- 擴充 `tools/desktop_mouse_smoke.py` 的capture／attack模式。正常首頁第一章曹操、固定remake seed1，透過設定關閉主君編成、正常編成37；進言宣戰選「我國較有利」，軍團上下反白、清單及地圖選點下令。途中存讀檔保留 `(202,114)`、節點82、連結0x0F30、路徑點0x4348與目標100 `(216,137)`。沒有direct-entry、時間／世界／軍團注入或強制勝負。
- `workplace/shots/attacker-march-20261103/attempt1/` 已正常宣戰並攻陷汝南，腳本錯把城兵守城當成必開戰術畫面。依既有corps.go及spec175確認：沒有駐城軍團就走城兵自動判定。保留失敗來源與trace，再明確拆成城兵及駐城軍團兩條驗收，不修改遊戲規則。
- capture分支63項通過。4/8 6時汝南歸屬由12改0，4/8 8時軍團37到 `(216,137)`、節點100、連結0；戰後再存讀檔保留歸屬與首都82。收據與實際使用的重播器版本在capture/。
- `battle-attempt1/` 嘗試呂布13／76時，實際進言已直接同意，腳本仍等待理由選單而失敗。改成依實際回覆處理。既有軌跡能確認袁術軍團3由壽春106經弋陽112、新息105往汝南，故軍團分支回到袁術場景，先正常等待3駐守，再下令37出兵；沒有猜補駐軍或改seed。
- `battle-attempt2/` 的attack分支67項通過。4/9 19時觀測軍團3在 `(216,137)`／100駐守；4/17 19時由37主攻3、汝南戰場、玩家側0。完整trace從frame0開始，重播器入口檢查點frame4；等開場及對話結束後點 `(496,280)` 突擊，常令與大將命令2／2，第900幀正常結束。汝南歸玩家，戰後存讀檔保留歸屬與首都82；玩家軍團仍在 `(216,134)`，沒有把取得城池說成已到城中心。
- 城兵攻陷訊息、實際汝南戰術畫面與戰後返回已目視。入口 `attack-normal-entry.png` 在PendingBattle首個檢查點擷取，仍為戰略過渡，沒有當作戰術畫面證據；戰術圖採 `attack-charge.png`。63與67項含共用操作，沒有當成130個獨立情境或全遊戲完成率。
- 本輪沒有改動遊戲Go來源，也沒有新發現需修改的規則。受測二進位仍為 `bcd7904992f486eb64fe4513b0a02c78a3d0594a8a1ca303ab53dc647a2e64c0`。新增playtest125並掛入CONTEXT及全專案文件索引；Issue #30的原版同狀態碰撞／亂數閘門仍未取得新收據，完整戰役、HD及手機仍未完成。
- GUI皆在network none、2CPU、2GiB、256程序、400秒限制及使用者1000:1000的Docker中執行。原版與程式工作樹唯讀，只有各案例輸出目錄可寫，遊戲存檔留在容器暫存。每個案例保存原始trace、JSON、截圖及使用的腳本來源；本機replay.sh保存同一啟動方式並通過語法檢查。未commit、push、發行、修改技能或原版。
- 封存362份來源及73份收據，入口 `workplace/shots/attacker-march-20261103/source-manifest.json`；`source.tar.gz` SHA-256 `5072bac8956e5a7a0eca624b189d414a102c271e5682d6b1dce3d232b0936c14`，逐檔回讀一致。上輪341份Go來源、356份遊戲來源／語系／module輸入及112份既有收據不變，原版KI.EXE與SINARIO.DAT雜湊吻合。
- 收尾文件索引573份／135條斷言通過，以playtest124作正對照確認125有入口；Python語法與 `git diff --check` 通過。修改檔抽查UID/GID皆1000:1000，整個工作根root-owned及偽.md目錄皆0，專案Docker容器無殘留；本機封存與驗收JSON維持gitignore。

## 2026-11-03：原版玩家城門碰撞的固定快照比較

- 接續Issue #30。核對CONTEXT、主機gh登入及Issue完成條件，遠端仍OPEN，沒有寫入。載入remake、dosgolem能力、文件職責及驗證交接路由；依專案規則委派唯讀代理檢查快照／佔用圖／AI上下文，未修改其他repo或原版，未commit、push、發行或改技能。
- 歷史 `workplace/parity/`、`workplace/dosgolem/` 已不存在，沒有拿舊路徑冒充成功輸入。原始SAVE第一槽玩家−1、年份0，只作診斷；重新從原版正常第一章曹操開局、編成37、向袁術12提出敵對並選「我國較有利」，逐句右鍵回覆後從地圖下令攻汝南100。
- 沿用dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`、KI.EXE SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；工具image `san1-go:1.26.7-ebiten2.9.9-r1`／`sha256:aac706e67af842cea8a6ad6011a1828161722381e1d7e2f21398aa31f28c07a6`。來源與原版唯讀，輸出在 `workplace/shots/collision-same-state-20261103/`。
- 快照沿用orig_snapshot已解映射，實際段tables2754、queue2514、occupancy46B3；另存258-byte raw RNG、byte_131AD事件節流及完整384×256佔用圖。未解區塊52A0..52BF保留來源模板，沒有假稱完整原版記憶體。Go先還原EventDelay，再接道路、地形、政略AI與戰術；不覆寫玩家、時鐘或游標，不用skip追齊。
- `march-probe/` 少關宣戰通知；`march-probe2/` 已回地圖後又多送右鍵。均保留截圖、時間軸與錯誤。依實際訊息與熱區修正輸入，`march-probe3/` 正常下令成功，軍團37目標100，沒有direct-entry或CallNear進行軍／戰鬥。
- 最初對拍器用Tick漏掉MCH更新；改TickMap後，發現Restore在SetRoads之後清掉道路標記，首拍軍費變19而原版450。調整初始化順序後差異消失。這是對拍器造成的假差異，沒有改正式遊戲；失敗版本與修正後來源分別保存於remake-probe／remake-map／remake-wired。
- 出發起點為196/4/1 1時子刻2，軍團37 `(206,114)`；首拍到 `(204,114)`、連結0F30、軍費450。依已解映射重建的完整區塊、raw RNG與佔用低7位在起點全同；後16拍的必要欄位也一致，17個檢查點通過 `departure-comparison.json`。
- 首次城門碰撞在 `(216,134)`、連結0FA0、路徑點46AC，目標汝南100 `(216,137)`。`gate-oracle/` 誤要求完整tick後倒數仍12；既有125A3呼叫順序及新1264A入口顯示，12880設12後同訪問即減到11。修正觀測前提後，gate-oracle2取到起點及80拍，沒有改遊戲規則。spec175現行「先倒數」文字改正，歷史原因保留本段。
- 全世界80拍的81個檢查點中，玩家37、資金、時鐘、事件節流及98,304格佔用值皆同；raw RNG從第10拍分岔，`gate-comparison.json` 保持失敗。首次世界記錄差異為AI軍團74的timer3／Stage1對Go timer1／Stage2。
- 新143AF原版入口實測DS2754、DI0D40、[DI+18]讀值0；Go只讀觀測gate-context-probe確認起點、第9及10拍stage1DI均−1，依spec195既有未知值近似回3，因而出擊並多取一次亂數。這是已記錄的殘留寄存器模型界線；不追整個DI生命週期、不猜填DI，也不修改正式規則讓全世界亂數綠燈。spec195追加新具體收據及限制。
- 為隔離背景AI，原版在正常執行到125BD／SI0940及125D8／SI2B80時讀取玩家37訪問前後狀態，沒有跳入函式。12次訪問各自從實際原版快照及raw RNG起跑，Go透過本機overlay僅呼叫正式tickOneCorps與tickStandoff；不改位置、timer、交戰或勝負，overlay沒有落成正式internal檔案。
- `visit-comparison.json` 12次皆通過必要欄位、完整raw RNG與佔用圖；9次訪問後區塊逐byte相同，另3次僅軍團繪圖bit4不同。這是12個局部訪問，沒有換稱為連續96拍整世界Go重播；正常GUI接線沿用playtest125的63／67項收據。
- 原版從4/7 17時子刻8的首次碰撞，經96個遊戲tick，在4/8 5時子刻5自然到14ADE（335,040,416道指令、SI2B80、DI14C0）。只驗觸發入口，未比較其後軍損或結算。正常下令、首次碰撞、自然攻城入口640×400截圖已目視。原版報一筆int61 AH0A／AL03未建模服務，本片不驗音效或硬體時序。
- 新比較器 `tools/march_snapshot_compare.py` 明確區分全域前59 bytes、玩家軍團、資金、事件節流、raw RNG與全圖佔用，並保留完整區塊差異。出發正對照通過；錯誤初始化版本於第1拍被拒絕；單改raw RNG一個bit的反例也於第1拍被拒絕。兩個反例均exit1，沒有靠缺檔、顯示器或環境錯誤產生假紅燈。
- 新playtest126、CONTEXT、spec175及195連結已更新。正式遊戲Go來源仍未改動，既有受測程式bcd79049與正常GUI收據維持有效。Issue #30仍需其餘碰撞／守方矩陣與全世界上下文，不宣稱完整remake、戰役、HD或手機完成。
- 所有工作在有界、network none、使用者1000:1000的既有Docker工具鏈內完成；本機replay.sh保存原版三種取樣時間軸及同一資源／掛載方式並通過語法檢查。來源封存、逐檔收據、擁有權、索引與容器收尾接續下列紀錄。
- 收尾：來源387份、收據1,467份逐檔雜湊皆通過，來源封包387個成員也逐一吻合。`source.tar.gz` SHA-256為`e4a50d83bbf99247e58c7980b023981510dbe13c623b938d6a0789e7c7468e8a`；來源與收據清單在同目錄`source-manifest.json`，結果摘要在`verification.json`，皆保留本機且受gitignore排除。文件索引檢查通過574份文件，`git diff --check`通過，dosgolem工作樹乾淨。全工作根未發現root-owned檔案或誤建的`.md`目錄，封包擁有者1000:1000；本輪一次性容器已退出並刪除。

## 2026-11-03：修正抵達分派前漏接的AI上下文

- 接續Issue #30；主機gh登入正常，遠端Issue仍OPEN，沒有改遠端狀態。依remake與dosgolem路由載入既有證據，先讀spec195及playtest126；使用既有Go工具鏈，沒有新建image或改其他儲存庫。
- 新正常原版重播保存於 `workplace/scratch/ai-context-20261103/`。從原版首頁、曹操、編成37、向袁術宣戰及正常地圖行軍到城門，才逐指令只讀追蹤DI；不跳入常式、不覆寫世界、不重擲。gate及前16拍共68份狀態資料與前輪收據逐檔一致。
- 實際第10拍，軍團74在143AF前的最後三個DI寫入為128F4／128F6／128F8：7D00→0350→06A0→0D40。既有固定IDA9.4匯出的12662及128F4共78條指令與KI原始bytes全同。來源是本次抵達前的node×32，沒有加840；106號據點因而讀據點40的佔用值0，留守且不取出擊延遲亂數。spec195先前只追兩個writer不足以解釋此路徑，已追加勘誤並保留舊失敗收據。
- 先寫spec271 DRAFT及可丟棄overlay。獨立唯讀審查發現途中留守的14300會額外改BX，初版把兩分支共用node×32不正確；限縮為目標相同、自軍據點分支後才升READY。初版與輸出保留。最初新測試有Factions[32]的fixture越界，修正工具資料後重跑，沒有把panic算成產品缺陷或錯版拒絕。
- 正式corps.go只在目標相同、自軍據點分派前寫node×32；途中留守不改。aimarch.go對據點表前的DI以目前世界區塊映射讀byte，涵蓋半筆勢力與表間區域；未知近似保留。原生存檔格式未變，Snapshot缺DI與其他上游仍列限制。state.go與aimarch.go的過期註解已更新。
- 新行為測試16例，涵蓋據點106／42／43／64、讀值2／3及未知／錯誤舊上下文，比較Stage、Timer和完整raw RNG。正式80拍81個檢查點全通過，完整世界區塊的2,144筆差異均僅在MCH物件區21C0..22BF，區域外零byte差異。沒有把它宣稱整塊或完整戰役一致。
- 實際移除本次DI寫入的錯版有8個行為案例失敗，完整80拍也在第10拍因狀態與RNG差異失敗；不是缺檔、顯示器或panic造成。正式冷測範圍為cmd/internal/mobile/translations，39套件通過、12個無測試，不掃入scratch探針。
- 正式程式 `wlgame` SHA-256 `e1d987e089e47724998d7b29feff7d8a776b3def5e5933c43a9a91fd6d660e44`。正常攻方67項GUI通過，包括自然守軍、行軍中原生存讀檔、玩家主攻、滑鼠突擊、戰後存讀檔。attack-charge及attack-post-load兩圖已目視。GUI與原版快照不是同一場戰鬥，不代證完整戰鬥或像素一致。
- 受測state.go在完成測試及建置後只訂正註解；`build-state.go`與`build-overlay.json`保留當時精確輸入，已確認去除整行註解後程式相同。其餘正式受測邏輯未變。spec271已CONFORMED，CONTEXT與AI機制文件回填，Issue #30其餘碰撞／守方／完整戰役仍開啟。
- 原版素材與dosgolem保持唯讀，所有工作在network none、1000:1000、有界資源與逾時的Docker內執行。來源封存、文件檢查與容器清理結果接續下列收尾紀錄。
- 收尾：來源409份、收據1,452份與外部輸入410份已建立雜湊清單，409個來源封包成員逐一核對通過。`source.tar.gz` SHA-256為`8b4be5aa31903f5530a00e117cbe643f2bc273352ab85fbd5b4ed2f045c74837`，與受測程式、verification及重播入口均只留本機且受gitignore排除。文件索引575份、135條斷言，索引檢查與`git diff --check`通過；全工作根未見root-owned檔案或誤建`.md`目錄，封包與程式擁有者1000:1000。dosgolem工作樹乾淨，本輪臥龍傳容器全部退出刪除，未清理其他專案容器。本輪未commit或push。

## 2026-11-03：正常玩家守方的城門碰撞

- 接續Issue #30。載入remake、dosgolem及文件路由；主機gh登入與遠端Issue讀取正常，Issue仍OPEN，未寫遠端。本地worklist verify僅作輔助，沒有把舊C研究或六月份守方快照當成本輪成功輸入。
- 唯讀稽核找到既有capture63的較短場景：37正常攻下汝南後，3仍按AI目標往汝南前進。由正常首頁、第一章曹操、編成37、向袁術12宣戰並正常下令，原版重生此場景；沒有執行siege、跳入函式、覆寫世界或改勝負。
- 本機收據根目錄 `workplace/scratch/defender-collision-20261103/`。原版4/8 6時子刻5由37守汝南中心；4/9 18時子刻0，3停在 `(219,137)`、Link10B0／Path4D4C、倒數11。軍團3的runtime位移2300與save位移2380分開記錄。
- 正常原版14ADE於381,275,913道、4/10 5時子刻5自然進入，SI2300／DI14C0；11B5A於381,381,065道，SI2300／BX2B80，確認3攻37守汝南100，這支沒有走14A7B野戰。既有AfterTicks把戰術節拍0推到1，385,693,354道後擷取640×400原版戰術圖並目視。
- 同一原版gate及凍結258-byte raw RNG起跑，等待94拍共95點；來攻3與守軍37分別比較，必要全域、玩家資金、事件節流、完整98,304格佔用低7位與raw RNG均同。另直接比兩軍完整64-byte記錄均零差異，整區塊2,528筆差異全限於MCH物件區21C0..22BF，區外零差異。
- 終版前81點共405份狀態資料及兩個自然入口metadata，與凍結attempt2逐檔相同。單改Go第1拍軍團3的X座標一個bit，對3的比較器即退出1。完整gate拍至攻城入口相隔95遊戲tick；96從碰撞前訪問點算起，沒有把不同起點混用。
- 新增滑鼠工具defend模式。正常首​​頁攻取汝南後留守、自然反攻、玩家側1、守陣及退卻、返回戰略與戰後存讀檔共68項通過。退卻後汝南歸12、首都仍82，原生讀檔保留。正式工具與實跑可丟棄來源逐byte一致，SHA-256為67f48c97c81f78638ad701d1439c2e229920126575bf5d51e04984df04e818ae。原版、Go守陣及戰後讀檔圖均已目視，不作同鏡頭、同節拍或完整戰鬥像素聲明。
- 遊戲程式仍為e1d987e0；本輪未改cmd/internal/mobile/translations的受測來源，已逐檔對前輪manifest核對，因此不重跑或冒稱新39套件冷測。GUI seed1與原版gate的raw RNG分開記錄。
- 新觀測也確認入戰排程缺口：Go第96次更新建立3攻37守的Pending後，仍完成軍團尾段／剩餘批次、MCH及時鐘，停在子刻6、軍團游標400；原版在函式中途入戰，子刻5、游標0。ResolvePending沒有續跑紀錄。不能只擋Clock；軍費、倒數、輪轉游標、剩餘軍團及存檔閘門都需一起驗證，未把未達READY的修法放入正式遊戲。詳見playtest127的未解及現況表，Issue #30保留開啟。
- setup-attempt是分派字串非唯一及o.Step不存在的工具初始化錯誤；original-attempt1監看offset錯誤，漏停自然入口；original-attempt2已取真實入口，卻等待戰略idle而失敗。回到既有對拍路由、改用已存在的戰術AfterTicks契約後乾淨重跑成功，未加等待預算。來源與紀錄保留，未當成產品缺陷。另一次唯讀來源檢查的自動權限審查逾時，依回覆重試一次成功，未擴權或改用主機分析。
- playtest127已連入CONTEXT，spec175與軍事機制回填。所有工作在既有network none、1000:1000、有界Docker內進行；來源封存、索引、擁有權與容器收尾另記如下。
- 收尾：來源428份、收據2,555份、外部輸入73份，全部3,056項雜湊與來源封包逐檔核對通過。`source.tar.gz` SHA-256為`3de631ac42ce81b646a324ee6573597a68435eb53e5d1b5f131ea589809f9bad`。來源、圖像、raw快照與manifest皆留本機且受gitignore排除。文件索引576份、135條斷言，索引檢查及`git diff --check`通過；全工作根無root-owned檔案或誤建`.md`目錄，封包擁有者1000:1000。dosgolem工作樹乾淨，專案容器已全數退出刪除。本輪未commit或push。

## 2026-11-03：戰術中斷世界更新與戰後續跑

- 接續Issue #30。先核對CONTEXT、原版守方收據及遠端Issue，仍OPEN且未寫入。載入remake規格閘門、dosgolem、文件職責與驗證交接路由；沿用既有Go／Ebiten image，未改原版、dosgolem或其他專案。
- 本輪使用playtest127固定的正常玩家守方擷取，沒有重新執行原版。既有IDA9.4匯出的130條指令逐條與KI.EXE原始bytes核對，確認125CF移動呼叫尚未返回便可入戰，125D2軍費／士氣、125D5倒數與125F6批次游標均在戰鬥之後。不是只延後Clock即可修正。
- 建立spec272 DRAFT與可丟棄overlay，從相同gate與258-byte raw RNG起跑。95個等待點仍通過，正式待戰改停子刻5、批次游標0；header3B..7F、MCH與軍團繪圖bit4以外無語意差異。原版入戰暫扣 `(219,137)` 一格佔用，Go保留邏輯佔用；Go Pending已做戰術Spawn而原版11B5A尚未初始化，故入戰RNG不混比，等待RNG仍完整比較。
- 唯讀代理審查續跑、事件傳遞及存檔入口，要求保留手機既有paused優先順序，已訂正。新增事件只送一次及ResolvePending實際勝利退出案例後通過，規格升READY才放入正式來源。最初測試的rng型別名稱錯誤是編譯用的工具錯誤，修正後重跑，不算錯版被拒絕。
- 正式世界層保存當次批首、剩餘配額、當前尾段與最初Tick／TickMap種類。建立Pending立即停；戰後按當前兵力處理軍費／士氣與倒數，再續完原批次，最後更新MCH、時鐘與每時／月結。同批第二場仍會暫停，勝負已定則丟棄續跑。桌面沿用共用Event報告，手機保留通知與暫停順序。
- World.SaveInto、原生Encode、桌面writeSave及手機SaveSlot在任何檔案寫入前拒絕未完更新；續跑期間也擋地圖指令與工具列存檔。續跑狀態只存在記憶體，沒有新增原生欄位或更改存檔版本。
- 七個新增測試涵蓋戰後兵力軍費、移動timer不重扣、原16格配額、第二場、MCH一次、事件一次、勝負退出與原版／原生存檔保護。隔離四套件整合通過；手機專用過濾命令未匹配測試，另跑完整手機套件通過，未把零測試當成暫停驗證。
- 正式冷測命令為 `go test -count=1 ./cmd/... ./internal/... ./mobile/... ./translations/...`，39套件通過、12個無測試；同範圍go vet通過。新程式SHA-256為 `af84322202bdd3317741791f59f5e35b572a12f6397b788c80cb91d64fcbf507`，保存於 `workplace/scratch/world-resume-20261103/wlgame`。
- 正式同起點重播再次通過95點及入戰時鐘／游標。真正換回舊排程的錯版，等待95點仍通過，入戰以子刻6、游標400兩項差異退出1；沒有依靠panic、缺檔或顯示器故障產生假紅燈。入戰完整raw區塊仍有37個排除項byte差異，未宣稱全區塊精確一致。
- 新程式正常攻方67項、守方68項滑鼠驗收通過。包含首頁第一章曹操、向袁術提案宣戰、編成37、正常出兵、行軍存讀檔、37攻3的突擊、3自然反攻37的守陣及退卻、戰後返回與存讀檔。兩側GUI固定seed1，未使用direct-entry、座標／軍團注入或強制勝負。守陣、突擊與守方戰後讀檔圖已目視。
- spec272升CONFORMED，更新唯一目前狀態表、即時制機制與spec索引。playtest127保留舊結果並連到後續修正。完整戰鬥結果、野戰矩陣、MCH動畫、完整戰役、HD與手機發行仍未代證；Issue #30仍開啟。
- 本輪所有工作使用network none、1000:1000與有界Docker資源，原始資料維持唯讀。原始失敗／修正來源、比較器與重播入口保存在本機同一案例目錄，未加入Git、未commit、push或發行。封存與清理結果接續下列紀錄。
- 收尾：封存480份來源、1,266份收據及599份外部輸入；其中518份原版收據逐檔對前輪manifest核對。全清單與480個來源封包成員雜湊均通過，`source.tar.gz` SHA-256為 `7b8a58041e2323d1ab555bb4b88b1b8e7ead3efa48c1f607c04cca86dc48b969`。入口為同目錄verification與source-manifest，產物維持gitignore。文件索引577份／135條斷言及`git diff --check`通過，以spec271作正對照確認272有入口。初次索引因歷史收據的狀態欄改名而拒絕，保留既有欄名後以同一命令成功，未改索引判準。
- 全工作根未見root-owned檔案或誤建`.md`目錄，新來源、程式與封包抽查皆1000:1000。dosgolem工作樹乾淨，本輪臥龍傳容器已全部退出刪除，未清理其他專案資源。下一個最小驗收是Issue #30的野戰碰撞：先從正常原版行軍取得兩軍相遇的固定快照，再驗角色與觸發前狀態，不沿用本片城門收據。

## 2026-11-03：正常道路野戰碰撞驗收

- 接續Issue #30，主機gh登入正常，遠端#30及#31均OPEN且未寫入。載入remake、dosgolem、文件職責及驗證交接路由，沿用既有image `san1-go:1.26.7-ebiten2.9.9-r1`，沒有新建image、改原版或其他專案，未commit、push或發行。
- 唯讀代理檢查既有野戰證據與接線，舊SAVE-FIELD是驗收捷徑，不當成正常道路相遇的新證據。主迴圈以playtest127的正常攻取汝南場景為起點，37到城中心後再從軍團選單及地圖下令往新息105 `(221,144)`，與已來攻的原版3在道路相遇。沒有siege、跳入常式、搬軍團、改交戰或強制勝負。
- 原版391,818,113道指令到首次完整阻擋拍，4/9 2時子刻0，3在 `(221,141)`、37在 `(221,140)`。其後95次完整更新共96點，399,364,343道自然走14A7B、399,471,569道走11B5A，3攻37守；403,494,228道後戰術節拍0→1並取640×400圖，已目視。原版raw RNG gate SHA-256為 `cea8741734034d452b89eee401ef8780b929940ba91cdc0055b9a5575aadb2a4`。
- 從原版gate與258-byte raw RNG起跑，正式Go兩軍96點的必要欄位、資金、時鐘、完整亂數、98,304格佔用低7位及事件節流均同。完整軍團3記錄皆同，37僅86筆繪圖bit4不同；全區塊另有MCH的2,560筆byte差異，其餘0。起點及另外9點整塊相同。入戰比較用11B5A，排除表頭3B..7F、MCH與bit4後語意差異0，子刻5、軍團游標0相同；佔用僅 `(221,141)` 原版暫扣0／Go邏輯1。Go Spawn後RNG不與原版初始化前混比，完整raw區塊入戰仍有41筆差異。
- 單改Go等待第1拍raw RNG一個bit，兩軍比較均在第1拍失敗，實際比較器退出1。成功原版與前次已取得的510份狀態資料相同，僅改比較檔名的480份輸入也逐檔核對，沒有挑選不同seed或重擲。
- 確認本場戰場198、不轉180度，但Provider仍用攻方Node105中心 `(221,144)`，原版用守方實際 `(221,140)`。十字圖塊分別D8/F2/DB/CD/F5及C9/00/00/C9/C9，經既有Terrain表都為0，因此本場碰巧同結果。取樣座標、快取及水域抽骰仍限於Issue #31，未擴修正式接線。spec196與戰鬥機制的舊下方五格散文依既有re05及規則程式補正為十字，未新增逆向完成聲明。
- 工具原版attempt1是截短觀測器留下未用符號的編譯錯；attempt2已取真野戰，但按住右鍵跨過11B5A後以IP停止會漏停。沿用上一輪OnCall旗標與AfterTicks契約後，同一路線重生成功，沒有加預算。失敗來源及有效擷取均留在 `workplace/scratch/field-collision-20261103/`。
- GUI seed1第一次在最高速截圖時錯過離城窗口而守汝南；第二次我方先靠近新息，使3在105改為留守，實際是攻城。trace只證實目標／Ordered／Stage變化，對應spec194受威脅中繼城留守是強推論，未因腳本場景不符改遊戲。重查GUI路由後，以3實際進入0x10B0道路才下令，正常取得3攻37守的mode1野戰。
- GUI第三次在frame148停住。第四次同輸入擷取SIGQUIT，實際render thread在OpenGL DrawElements→DrawTriangles內syscall，另一goroutine等待graphicscommand.Flush；未在尋路、單挑或音訊。未把先前靜態尋路候選寫成產品缺陷。SIGQUIT後失效視窗抓圖留有界容器，診斷完成後只刪除本輪該容器；驗收器抓圖補10秒逾時。
- 回到平台規格優先路由，參考Mesa官方envvars契約。既有image的GLX小探針確認預設llvmpipe LLVM15.0.6／OpenGL4.5、softpipe／OpenGL3.3，均為Mesa22.3.6；不重建image，以GALLIUM_DRIVER=softpipe重跑同一程式與玩家輸入。結果及收尾接續記錄。
- 新playtest128已掛入CONTEXT及軍事／戰鬥機制；正式359份cmd/internal/mobile/translations來源逐檔與spec272封存相同，39套件冷測與既有攻方67／守方68仍屬上一輪收據，本輪未重跑或冒稱新增。新GUI與完整來源封存待本片收尾。
- softpipe完整抽測在既有25秒行軍等待內未觀測到移動，聚焦版本亦因繪圖速度限制中止，沒有加預算。依Mesa官方envvars契約改為 `GALLIUM_DRIVER=llvmpipe LP_NUM_THREADS=0`，同一image與程式的正常野戰50項通過。範圍保留正常首頁、設定、宣戰、編成、出兵、行軍存讀檔、敵軍走上10B0後才下令、mode1的3攻37守、開場、守陣、退卻與戰後存讀檔；省去本場前已另驗的語系與全選單邊界抽測，沒有略過玩家動作。
- 最終GUI於4/9 17時入戰、4/9 18時戰後讀檔，汝南0／新息12／首都82保留。守陣與戰後讀檔圖已目視；同一畫布的地形、側欄與按鈕正常顯示，未作同狀態像素聲明。遊戲沒有改動，不將本機Mesa參數選擇寫進產品規則或宣稱其他平台均已驗收。
- 正式 `tools/desktop_mouse_smoke.py` 與實際受測來源逐byte相同，SHA-256 `4e84a4dfe6b9cd05eb649f17e55fd37d87153029bc06912275f84945a20faf6b`。可重生腳本以修改前工具為固定輸入，不因正式工具升級而重複加分支。來源封存、索引與清理結果接續下列紀錄。
- 收尾：516份來源、2,481份收據、602份外部輸入及516個封包成員逐檔核對通過。`source.tar.gz` SHA-256為 `0847f00c95c47794e4f88150be8bc3fa10c56373ecfd2eee10160db566a1eb4c`，入口為同目錄verification與source-manifest，原版資料、圖像與封包維持本機gitignore。文件索引578份／135條斷言通過，以playtest127作正對照確認128有入口；Python語法、wrapper語法及`git diff --check`通過。
- 全工作根未見root-owned檔案或誤建`.md`目錄，正式工具、新文件及封包抽查1000:1000。dosgolem工作樹乾淨、HEAD仍f79eab6；本輪臥龍傳容器全部退出刪除，未清理其他專案資源。Issue #30與#31未關閉。後續查證入口為Issue #31、本片的兩處十字取樣及Provider／rollBattlefield接線，先建立非平原差異的同狀態證據再改正式路徑。

## 2026-11-04：接戰位置與單次選場接線

- 使用者繼續既有remake修正。先讀CONTEXT、現行Issue #31、復古路由的spec閘門及文件職責，沿用逆向remake技能；不擴入HD或手機設計。Issue #31在主機gh核對為OPEN，沒有遠端寫入。
- 固定IDA匯出SHA-256 `32ebc39f…296a5d9` 的14A7B／14B63／14BDD／14C1A／14C4C逐條與KI.EXE bytes相同。14A8E選場在14A9E挑守軍之前；座標屬最初接觸守軍，朝向取參戰玩家且守方優先。新spec273經獨立唯讀證據審查後READY，再改正式程式。審查修正了14B7D／14B80位址及朝向敘述。
- 類型9舊「鄰格CA」解釋不成立：14C33 LDS切換DS，14C41仍讀切換後DS的SI+1，不是軍團勢力。既有固定C與獨立model已處理此資料流；本輪不把錯誤XOR公式寫入Go，維持已標示的不翻轉近似。兩側非玩家的入口AL也未深挖。spec121追加可回查勘誤，Issue #31仍OPEN。
- fieldAt改為先選場再挑代表守軍，以顯式參數送到fight／beginTactical。Pending新增執行期編號及翻轉，Node仍供據點與戰果使用。Provider移除世界參考、按據點快取與抽骰，僅按明示編號載資產；Field／Category／Script、桌面／手機圖像、標題、喊話、單挑及音樂共用結果。正常入口掛地形，桌面讀檔重新掛地形，手機讀檔重新掛地形與戰術。存檔格式未改。
- 四個受影響套件冷測通過；正式範圍`./cmd/... ./internal/... ./mobile/... ./translations/...`的40套件冷測與11個無測試套件通過，vet通過。移除Provider最後旋轉快取後另重跑該套件。錯座標overlay實際選出198而非193，多抽一次overlay在建場前觀測到2次而非1次，兩者均以斷言退出1，不是編譯或環境失敗。
- 沿用spec251固定O2驗證二進位，只跑terrain：dosgolem原版與固定C的60案／120階段完整RAM／plane通過，固定258-byte合成RNG在執行前設好，沒有重擲。不把這份函式收據冒稱Go正常海戰。所有原始輸入唯讀掛載。
- 新Go取樣器沿用playtest128正常原版gate及raw RNG。兩軍96點等待完整RNG／佔用圖通過；入場編號198、不翻轉、子刻5與游標0相同，排除既有表頭／MCH／renderer後語意差異0。原版一格暫扣及Spawn前後RNG階段仍分開，未放寬比較器。
- 新程式SHA-256 `eb13ad8822e3f99eb07caebeb8b6c2c58c597e819f7628dd08a4fa20c59766f2`。同一Docker image與LP_NUM_THREADS=0下，從首頁正常滑鼠取得野戰50項、攻城67項、守城68項通過，含行軍與戰後存讀檔。三場截圖已目視；不同亂數、時鐘與鏡頭不作像素對拍。手機只驗共用讀檔／戰場接線，未宣稱裝置實機完成。
- 守城第一次停於行軍中開存檔，當時自然訊息使工具列正確停用。驗收器改為在原8秒預算內正常滑鼠推進訊息後開存檔，同一seed1、程式與400秒總上限重跑通過。失敗`gui-defend/`保留，成功在`gui-defend-retry/`；未因腳本時序修改產品規則。
- 初次測試的兩個未用匯入／變數屬編譯接線問題，修正後同命令乾淨重跑。整批go fmt曾誤觸30份原先乾淨來源；逐檔證明內容等於HEAD經同一gofmt結果後，只還原這批自有格式變動，保留所有原先髒檔。.git為唯讀，未用git restore改索引。沒有commit、push或Issue狀態寫入。
- 正常非平原／碼頭海戰、類型9翻轉、AI雙方入口AL、MCH動畫及完整戰役仍是證據限制。spec273只就已列選場接線與抽驗CONFORMED。來源封存、索引及Docker清理接續收尾紀錄。
- 收尾：索引579份／135條斷言通過，spec272作正對照確認spec273有入口。第一次索引只缺規定的日期欄位，補齊後同命令重跑通過；失敗log保留。1,862份來源與封包成員逐檔核對，489份收據、564份外部輸入保存雜湊。`source.tar.gz` SHA-256為 `f669484faf446227dc32470edd9379252d24f9048387beac4c23563a0261c41a`；入口 `workplace/scratch/field-selection-20261104/verification.json`，封包維持gitignore。全工作根未見root-owned檔案或`.md`目錄；本輪Docker容器全部退出並移除。工作樹保留原先未提交修改，沒有commit或push。

## 2026-11-04：正常CA水域野戰與海戰適性

- 使用者繼續既有remake。核對CONTEXT、復古重製路由、規格閘門與文件職責；主機gh認證成功，Issue #31仍OPEN，沒有遠端寫入。沿用既有Go／Ebiten image與dosgolem；所有分析、原版與GUI在有界、無網路、1000:1000容器內進行，原始資料唯讀。
- 從原始地圖254條路掃到35條含CA／C0–C3，第一章曹操領內無此道路。選孫策1／首都129濡須口，原版首列編成18孫權；對手劉繇10／首都116建業。沒有搬軍團或用StageBattle／OpenSiege取得正常原版收據。
- 唯讀委派查證發現playtest119的舊AI碼頭樣本是4/16受控存檔後自然演進，引用的原始存檔／RNG目前缺失；playtest48橋上圖屬許昌攻城，均不當成本輪正常水戰起點。沿現存道路122牛渚至116建業建立新玩家路線。
- 原版正常首頁→第一章孫策→編成→進言宣戰劉繇→出兵牛渚→空城自動攻城→太史慈72自然反攻→孫權向建業出迎。4/3 21時自然進14A7B，18攻72守，實際位置304155／304154皆CA，戰場212、不翻轉。原版兩次相同輸入490份狀態資料逐byte相同；未挑seed或重擲。
- 原版14A7B步數275,647,559，14B63為275,647,569，11B5A為275,753,980。從選場前258-byte raw RNG取一次Next就得到212及11B5A的完整狀態。Go從同一等待gate raw RNG初始化，兩軍96點完整RNG／佔用圖等比較通過；入戰子刻7及游標同原版，排除既有表頭／MCH／renderer後語意差異0，原版暫扣一格另列。Go選場後Spawn前完整RNG也逐byte相同。
- 原版第一戰術節拍讀unit段1D30E與summary段1D30A，類別2、將領18／72正對照。普通兵+18h戰力與大將戰力／HP均與Go相同：18六隊13,13,9,9,7,7、大將26／138；72六隊14,14,9,9,7,7、大將56／196。與誤讀陸戰時的18大將24、72大將54等明確不同。初始化能力比對不擴張成原版tick1與Go frame0的全部兵位置一致。
- 相同gate實跑錯座標Go overlay，96點等待仍通過，入場卻選198；選場、Spawn前RNG及海戰戰力均不符。同一比較器退出1，非編譯失敗或panic。
- 遊戲與語系365份來源逐檔與spec273封存相同，程式SHA仍eb13ad88…59766f2；不重跑或冒稱本輪新增40套件冷測。正式desktop_mouse_smoke新增water模式，與實際受測來源逐byte相同，SHA-256 `38b4ffe95149022d068e78137d58169635e732c7a0ff71dd22354c0b00f97715`。46項正常滑鼠水戰通過，含行軍及戰後存讀檔。GUI seed1為72攻18守、4/3 22時，兩軍均在CA，與原版固定快照攻守／時點不同，沒有像素對拍聲明。
- 診斷保留：第一份開局沿用曹操選列，依playtest122改144／152；孫策直接接受宣戰，移除曹操理由選單的腳本假設；觀測器原先把空城自動攻城當成停止點，改為繼續；道路偵測把+0C當連結，依spec172校正+0E，前次實際守城不改名當水戰。GUI第一次沿用曹操鏡頭而使孫權在畫面外，改為正常縮圖跟隨軍團，同程式／seed／400秒上限重跑通過。未因工具失敗改產品規則或增加等待預算。
- 新playtest129掛入CONTEXT、spec115／121／273、戰鬥機制及生成索引。海戰適性與「未實跑碼頭」的舊缺口已回填；類型9翻轉、AI雙方AL、正常自動野戰結果及完整海戰結算仍待驗，Issue #31維持OPEN。沒有commit或push。封存與清理結果接續收尾紀錄。
- 收尾：獨立唯讀審查從3,072-byte單位區及64-byte摘要重算雙側能力，與原版／Go收據相同，並核對GUI不同攻守、frame0／tick1及單場界線；非water工具參數維持原值。正式工具可由固定舊來源重生，語法及實跑來源身分通過。索引580份／135條斷言通過，以128作正對照確認129有入口；`git diff --check`通過。1,891份來源含私有重播器與封包成員逐檔驗證，5,479份收據、451份必要外部輸入保存雜湊。`source.tar.gz` SHA-256為 `9d92d84fba13d56daee6359b5cc293ccaa520a55f684da95afacf0c239f0e01f`，入口 `workplace/scratch/water-field-20261104/verification.json`。封包與原版衍生物維持本機gitignore。全工作根無root-owned檔案或`.md`目錄，本輪Docker容器全部退出移除；dosgolem仍乾淨f79eab6，未清理其他專案資源。

## 2026-11-04：原版清單拖曳與HD規劃入口

- 使用者要求抽驗、優先操作手感與視窗，接近後規劃HD。讀CONTEXT、GitHub Issues與本輪路由，載入GUI還原、規格閘門及文件職責，沿用reverse-engineer-retro-game-remake技能。主機gh認證成功，沒有Issue寫入。水域維持既有playtest129完成範圍，不擴矩陣。
- 既有唯讀代理核對rich2目前工具列，採常駐外部畫布參考；rich2的原版鍵盤相容不帶入本作，全滑鼠與上下推動反白沿用既有決定。未改rich2、dosgolem或原始資料。
- 固定c-list IDA9.4匯出四函式、c-catalog的loc1796C共五組bytes均與KI.EXE相同。正常曹操開局開22列勢力表，在18503回呼返回後取樣：Y128/152/184/240預覽0/0/6/12，滑塊偏移0/0/35/70，正式byte181BC到放開才寫12，時鐘不變。
- 初稿誤把正式top未提交解讀成畫面不換列。原版實際截圖及1796C回呼確認列與滑塊同步預覽，立即訂正DRAFT；獨立審查後READY才寫生產實作。原始錯誤探針及未消費位移的診斷取樣均保留，正式收據在original-stable。
- spec274新增共用listscroll.go，按住跨槽／工具列仍由清單捕捉，任一滑鼠鍵維持拖曳，放開提交一次。新遊戲勢力表也接同一換算；渲染預覽列表、隱藏游標、換表失效，沒有增加快捷鍵或改世界規則與存檔格式。
- 正式40套件冷測、11套件無測試、vet通過。獨立審查抓到工具列箭頭漏接，並補啟動器原生游標隱藏／恢復；受影響cmd/wlgame再次冷測、vet、重建通過。最後程式SHA-256為c6cff2ffd948cb8a88971aec64dcb1d0ba4c19b7688f45c02dcbc971c7bff907。
- 正常滑鼠windows模式34項通過：首頁、兩種清單拖曳、放開在工具列不觸發按鈕、右鍵維持捕捉、反白確認原列、編成子視窗取消回父表／地圖、存讀檔。五個取樣的2560-pixel捲軸均相同；清單預覽6、12及放開後400×192視窗各76800 pixels相同。預覽0只有曹操武將數12/11的48 pixels差異，原因未解。鏡頭與完整世界狀態不同，不宣稱整張畫面或整款完成。
- 錯span overlay的Y152得到top3而非0，測試退出1；單次點擊overlay正常到遊戲清單後無法取得拖曳擁有權，GUI斷言退出1。兩者非編譯失敗或panic，正式來源未被覆寫。
- 首次GUI的按下尚未被接收就移入工具列，屬驗證腳本事件合併，改等按下及捕捉已觀測後再移動。同程式同seed重跑，沒有調遊戲規則。截圖另等穩定Draw，避免把上一幀當成現在。SVG工具預設缺rsvg-convert，內建MSVG也不可用；改由現有Go圖形函式畫草稿預覽，未安裝新工具或image。
- 依grill-me／grilling建立可丟棄工具列A/B HTML、SVG及Go預覽。使用者選A，原版/HD及縮放常駐，並新增頭像HD、戰術AI工具列選項與武將編輯。新增入口不能直接壓縮九鈕稿；正式介面仍七鈕。
- 使用者確認AI同時加成找路效率、單兵走速與架梯動作，我方敵方各依自己的統帥。已展示無障礙概念動畫比較1.5/2/3倍，使用者選最高2倍，統帥1為1倍、15為2倍、線性增加，只在AI輔助開啟時生效。這是新增設計，非原版規則。
- 武將編輯已查現有General型別及原始偏移：武術/統帥/政治1–15；委任Tactic是0–7腳本類型；攻城/陸戰/海戰適性0–10。需求及下一閘門在spec275 DRAFT，尚未改正式AI數值、頭像素材或存檔。HD美術、開關時機與編輯保存邊界仍待草稿和確認。
- spec274就本片CONFORMED；spec275新增並掛入CONTEXT及索引。原版衍生物、字型草稿與封包維持本機gitignore，沒有commit、push或發行。來源封存、索引及容器清理接續收尾紀錄。
- 收尾：索引582份／135條斷言通過，以spec273正對照確認274／275均有入口；git diff --check通過。1886份來源及封包成員逐檔核對，143份收據與931份外部輸入保存雜湊。source.tar.gz SHA-256為 `6de27a5a1bae7f1b6787a7ca71f6da975d7b23ac9178778e916ebf8b93d02f9a`；入口 `workplace/scratch/ui-windows-20261104/verification.json`。全工作根無root-owned檔案或`.md`目錄，最後寫入檔仍為1000:1000。本輪介面與HD草稿容器皆退出移除，未建立image、未清理其他專案資源，沒有commit或push。

## 2026-11-04：可操作HD概念網站

- 使用者要求先做網頁確認HD概念。核對CONTEXT、復古GUI路由、文件職責，載入Sites及imagegen技能；沿用既有共同決策與spec275 DRAFT。網站只作可丟棄概念，不把新美術或AI規則接進正式遊戲。
- 在workplace/hd-concept-web-20261104建立獨立靜態網站與Git來源。Sites僅建立一次，精確ID已保存於.openai/hosting.json；保留owner-private。helper經Docker隱藏stdin接收短效憑證，不寫入檔案或命令列。原版遊戲儲存庫沒有commit或push；網站來源依Sites技能保存到自己的私人儲存庫。
- 內建imagegen新製曹操／關羽／諸葛亮三聯頭像，另以原生SVG重畫中原與攻城示意地圖。委派只處理maps.js及瀏覽器QA。網站沒有原版EXE、圖像、字型或資料檔；原版風格模式使用新頭像的低解析呈現，不宣稱原版畫面對拍。提示詞、素材SHA與來源在asset-provenance.json。
- 工具列常駐於畫布外，整合原版/HD、縮放、AI、武將編輯及原有入口。正常點城池、兩點確認原武將、編成子窗、右鍵返回及上下推動反白皆可操作。編輯支援基本能力1–15、三種適性0–10、委任腳本0–7，取消不套用；四槽只存瀏覽器localStorage。意見可本機保存及下載，不自動傳送。
- 攻城頁演示雙方統帥1–2倍、部隊移動及架梯；屬無引擎概念，不能作正式AI／戰局結果證據。即時切換只供比較，正式模式切換時機仍待確認。畫風、地圖及示範數值沒有升格為已選定的HD設計。
- 沿用既有Playwright1.56.1映像與鎖版module，沒有新建工具image。初輪QA三個關閉鈕strict locator歧義屬腳本問題，改選唯一的X後8/8通過，來源不變。目視另發現0%時靜態梯子已出現，修為SVG clip隨進度增長並重設。修正版9/9通過、頁面與網路錯誤0、來源四檔雜湊前後相同，含四槽reload、手機390×844橫捲與編輯，梯子0→部分→100%高度73→重設0。
- Sites首版正常保存後，修正版用同一網站另存新版本，保留原版本，未移動或覆寫已存版本。最終來源commit f2e3197157b4a2edd1695d900ae2fb1cb363b8bf，私人部署appgdep_6aead54e218c8191a2c8d742ffee95ff已由原生工具回報succeeded。網址 https://wolong-hd-concept-anr2.cy-wang.chatgpt.site 。未提供open_in_codex工具，以可點擊網址交付。
- spec275與CONTEXT掛入網站來源、範圍、素材及驗證入口；狀態仍DRAFT。索引582份／135條斷言通過，工作根無root-owned檔案或.md目錄，網站輸出1000:1000。所有預覽server、瀏覽器與工作容器已收尾，沒有長駐本機網站程序。封包與完整收據見workplace/hd-concept-web-20261104/.qa/verification.json。

## 2026-11-04：原版圖像的本機HTML強化比較

- 使用者指定概念圖改用原版，HD也從原版強化，只在local建立HTML，不上拋OpenAI。此決定取代前一版新繪素材與Sites交付；本輪沒有任何Sites、imagegen或網路呼叫，也沒有將原圖或瀏覽器截圖交給view_image。先前網站未更新、未同步；所有新原版衍生物留本機。
- 讀CONTEXT、spec275、復古GUI路由、文件職責及Sites local-only分支。建立與雲端來源分離的workplace/hd-original-local-20261104，沒有.openai或.git。沿用既有rs2-hdtools:r1，Pillow11.3.0/numpy2.3.3；確認無離線神經放大模型後，使用本機鄰域取色4倍輪廓平滑，以及可選Lanczos／輕度銳化，不自稱AI補畫。所有容器network none、資源有界、1000:1000，原始檔唯讀。
- 唯讀查證選定playtest128的march-start與field-tactical、playtest127的defender-tactical、spec274的parked勢力表。四張均原版640×400，輸入SHA與既有收據相同，不用remake編成圖冒充原版。
- 以既有library/state解碼原版KAOGRF與劇本，實際150頁；曹操16→50、關羽40→114、諸葛亮60→125，均64×64。最初兩次將General.Name原始Big5誤當UTF8比對，回查text.Decode契約後以專案CP950解碼器修正；未改原始資料。能力值也由劇本一讀取，不沿用前一版示範數字。
- build.py保持七張原圖PNG bytes不變，產生兩種4倍強化，共21張圖。輪廓版本逐圖確認取色為原圖顏色子集。所有圖像、CSS、JS嵌入單一index.html，約19.6MB，無外部URL／script／stylesheet／字型，CSP禁止外部連線。可切原版、HD、並排、拉桿、兩種強化、2×/4×檢視與同步捲動，另有原圖頭像的武將編輯／AI選項示意。固定畫面不造假即時部隊動畫，沒有正式存檔寫入。
- 沿用本機Playwright1.56.1在完全斷網容器，以file:///local/index.html開啟，11項抽驗通過、JS錯誤0、外部請求0。驗證21張嵌入圖片雜湊、全部七樣本、拉桿、縮放、編輯取消／夾限／套用、來源、筆記、390×844及HTML搬到另一目錄仍可獨立開啟。第一輪腳本點擊位於視窗下方的拉桿，改為先scrollIntoView後同環境重跑；首輪失敗保留，產品與圖像未因此修改。截圖只存本機，未送任何模型目視。
- README、CONTEXT與spec275的現行入口改為本機HTML，舊網站留歷史，規格仍DRAFT。原版強化程度待使用者親自確認，不聲稱正式HD完成，也未改遊戲規則或存檔。索引582份／135條斷言通過，來源與完整收據在qa/verification.json，沒有commit、push或發行。
- 收尾：index.html SHA-256 `a92db68ee4d10705e4ad6e733f1b7b4735409c93ce350aa82cf2e7b63741bda0`，19627119 bytes；工作根無root-owned檔案或.md目錄，輸出1000:1000。沒有新建image，所有本輪處理與驗證容器均使用--rm收尾，沒有長駐server。

## 2026-11-04：原畫HD細節重建及rich2式縮放原型

- 使用者認為前版差異不足，要求基於原版更好看，另指定HD版參考rich2支援拉遠拉近。沿用本機HTML、不將原圖上傳OpenAI的裁定。核對CONTEXT、復古GUI路由、文件職責及grilling可撤回prototype分支，不把視覺樣本接進正式規則。
- 沿用既有san1-voice:qwen3-tts-0.1.1-cpu-r1的PyTorch2.9.1+cpu作本機推論，未建立重複image。公開模型下載只掛models目錄，沒有原圖；官方Real-ESRGAN v0.3.0架構、v0.2.5.0 general-x4v3／animevideov3權重及BSD授權均存本機並記SHA。官方primary來源已核對，沒有使用imagegen、Sites或雲端影像處理。
- superresolve.py在network none容器以CPU執行，原圖及模型唯讀，只有beauty輸出目錄可寫。頭像用general，場景用anime模型，另作輕度色彩、明暗與銳化，不使用換臉或擴散模型。七張原始PNG與構圖沿用，與前版原圖雜湊核對；index-v1.html保留前版SHA。模型原始輸出及統計在beauty/superresolution.json。
- build.py將戰術側欄／底列、戰略頂列／提示框、勢力清單合成回原版整數倍像素，只在提示框頭像保留強化。區域及例外在manifest，這是prototype選擇，不作原版行為新斷言。指定區域4,947,968個輸出像素精確相同；場景內文字及視覺效果待使用者確認，像素差異不等於美術驗收。
- 同一個index.html新增「原畫HD」並預設曹操並排，保留前版輪廓／柔和放大，共28張內嵌圖。單檔約51.5MB，可搬移離線開啟；CSP禁止外部連線與資源。原圖及瀏覽器截圖未經view_image或其他雲端模型目視，所有圖像處理留本機。
- rich2以唯讀掛載核對native_zoom.go、toolbar_controls.go、themetoolbar/toolbar.go。比較頁加入減號／倍率重設／加號，1–5、步進0.2、預設3.4、160ms過渡，滾輪及拖曳平移，雙圖同步。比較頁3.4代表適合視窗的預設圖片檢視；正式HD需世界鏡頭獨立縮放、介面固定，相關焦點、邊界、命中座標及倍率仍在spec275 DRAFT，未以固定截圖宣稱引擎鏡頭完成。
- 完全斷網Playwright以file網址抽驗13項通過，JS錯誤0、外部請求0；驗證28圖、七樣本、三種強化、拉桿、縮放上下限／重設／滾輪、拖曳與雙圖同步、切換強化保留焦點、編輯取消／套用、AI、來源、筆記、390×844及搬移後獨立開啟。圖像核對另見qa/beauty-verification.json，來源及環境入口見qa/verification.json。
- 補焦點保留後，viewer的data-sample狀態與測試的樣本按鈕選擇器撞名，手機最後一步得到兩個匹配。屬驗證選擇器／狀態標記的命名衝突，改為data-viewer-sample後同容器、同命令乾淨重跑13項通過，沒有調美術或遊戲規則。失敗診斷保留於qa/zoom-selector-attempt.txt。
- CONTEXT及spec275更新需求與目前樣本，README的穩定本機入口不變。沒有修改正式遊戲、存檔或原版資料，沒有commit、push或發行。Docker與擁有權收尾結果由同目錄qa/verification.json保存。
- 收尾：HTML 51,487,406 bytes，SHA-256 `875276f842e85b2a5642d0f28aeafc1567a96c7716bd371a3a3d34116d6819ab`；索引582份／135條斷言與git diff --check通過。全工作根無root-owned檔案或.md目錄，本機樣本皆1000:1000。本輪wolong-hd／wolong-local-hd／wolong-rich2容器已全部退出移除，沒有新image或長駐server。

## 2026-11-04：HD重新定義為原版依據的城池與地圖重繪

- 使用者要求重新定義HD，基於原版先重新手繪城池與地圖。原圖不傳OpenAI、只建本機HTML的要求沿用。此決定取代上一輪超解析放大的美術方向，HD頭像及rich2式拉遠拉近需求保留。核對CONTEXT、復古地圖／GUI路由、文件職責、grill-me／grilling及既有資料解碼技能，不新增正式規則或逆向切片。
- 新批次放在workplace/hd-redraw-local-20261104，與前版及先前雲端網站分開。以既有world／palette／state parser讀取原版MMAP.MAP、MMAP.MDL、GAMEPAL.BRG及SINARIO.DAT，來源唯讀、輸出1000:1000。確認世界384×256格、256圖塊、192據點及254條道路，原始圖像只在本機重建，沒有呼叫imagegen、Sites、網路或view_image。
- 第一批取世界格半開區間(180,96)至(238,145)，928×784，涵蓋許昌、虎牢關、舞陽、汝南等13個據點、27條相交道路。據點錨點及路線點列沿原資料；保留原圖PNG，另以SVG新線條繪製屋頂、瓦片、城牆、垛口、城門、山石、樹木及水紋。水域由原色外形描出，山林採色彩分佈的視覺詮釋，未當成已證實的地形分類。新城池外觀與尺寸屬設計草稿，不宣稱原版建築形制或戰術城門幾何一致。
- 提供A青綠彩繪與B古紙淡彩，每款一張區域圖、許昌大城82、滎陽中城63、舞陽小城87、虎牢關72四張近看圖，共10張新SVG。近看圖刻意放大新城池造型方便審視，原版對照仍為9×9格素材合成，不冒充原版執行截圖。這批是SVG線條繪製的手繪風草稿，不宣稱人類畫師逐筆手繪。
- HTML內嵌五張原版PNG及十張SVG，可切原版／重繪／並排、畫風與據點名。rich2式減號／倍率重設／加號、滾輪縮放、拖曳平移及雙圖同步沿用；畫風切換保留焦點、工具列固定。可在本機保存SVG。正式世界鏡頭與完整世界美術仍待設計和使用者確認，沒有以網頁縮放代證正式HD引擎。
- 幾何核對通過：四個原始輸入SHA未變，原圖裁切逐像素相同；13個新圖據點錨點與27條逐格路徑等於原資料；新SVG沒有內嵌舊圖像。這只證明來源及幾何，不能證明好看、地形解讀或原版parity。
- 完全斷網Playwright以file網址抽驗12項通過，JS錯誤0、外部請求0。驗證15個嵌入資產、五樣本兩畫風、據點名、原版／重繪／並排、縮放／重設／上下限、拖曳及雙圖同步、切風格保留焦點、來源、SVG保存、HTML搬移及390×844。十張新SVG另可本機輸出3倍PNG，非空白且色層存在；截圖與PNG全留本機，未送模型目視。抽驗沒有失敗或產品數值調參。
- 以來源唯讀、只掛明確輸出重生art.py及build.py，HTML、manifest與全部15資產的前後SHA完全相同，qa/rebuild-proof.json保存收據。沿用既有Go、pod-intro-img:r2、rs2-hdtools:r1及Playwright映像，沒有新增image或下載套件。
- CONTEXT、README、spec275及規格索引改為新入口，前版超解析及雲端網站保留歷史。A／B選擇已向使用者展示，尚不把初始顯示A視為定版。所有新美術維持本機，沒有改正式遊戲或存檔，也沒有commit、push或發行。來源、資產及驗證入口為同目錄qa/verification.json，收尾狀態接續記錄。
- 收尾：HTML 3,604,161 bytes，SHA-256 `0f225663e363625e2913114603377dbb52a0c67d668fc3a26a0aadac7ed0b91b`。索引582份／135條斷言與git diff --check通過；全工作根無root-owned檔案或.md目錄，本批輸出1000:1000。所有wolong-redraw容器均已退出移除，沒有新image、長駐server或雲端同步。畫風仍待使用者回覆，未擴大成正式世界地圖資產。

## 2026-11-04：否決通用SVG造型，回到許昌原圖修繪

- 使用者明確否決A青綠彩繪、B古紙淡彩，指出失去原版特色；後續「繼續」沿用原圖修繪與本機限制，未解讀為接受前兩案。decision.json保存否決及理由，CONTEXT與spec275撤除A／B待選狀態。前批只對齊位置與路線，建築用通用SVG依Kind合成，這不足以保存原版美術；不再擴展該模板。
- 核對原版素材／GUI還原路由與grilling原型閘門，補讀re50及library。許昌82的KindHigh=6才是KYOGRF據點景觀頁碼；實際15頁、每張96×96。它與MMAP上的地圖城池是兩份素材，不能只按大城類型猜造型。沿用原解析器、CP950解碼與原版調色盤，輸出五座城池的實際來源，這輪只修繪許昌。
- faithful/extract.go讀取KYOGRF、MMAP、調色盤及原劇本；原始資料唯讀。許昌地圖取原座標206,114周邊7×7格112×112；景觀取第6頁96×96。prepare_inputs.py建立四倍參照448／384、深色像素及強邊緣保護區，沒有通用城樓或新山林幾何。來源、原圖與映射在faithful/originals.json。
- 為滿足素材不傳OpenAI的限制，沿用既有CPU PyTorch image，單獨下載公開SD1.5模型固定commit 451f4fe16113bff5a5d2269ed5ad43b0592e9a14及Diffusers0.35.2鎖版wheel。下載容器只掛models，未掛原圖；安裝在同一image的容器內完成純Python套件層，網路關閉，未掛主機Python或venv，也未建立重複image。官方模型與Diffusers API來源已查證，下載檔逐一驗SHA。
- 最初metadata腳本名inspect.py遮蔽Python標準inspect，改名public_model_probe.py後成功；不屬於產品缺陷。全float32首輪在第一張輸出前exit137，未取得可確證的OOM事件，不能將推測寫成確診。改為同一image／同一命令、UNet線性層原地動態int8、卷積float32，文字編碼後釋放CLIP；兩張各51.63／66.26秒完成。精度變更是本機影像工具近似，不改遊戲規則。
- 本機paint.py直接使用兩份原圖，各以384×384輸入，seed固定20261104，strength0.25、16個排程步中各4個低噪聲步；沒有重擲後挑圖。原始模型輸出保留，合成小樣限制表面細節改動並沿原圖廣域色塊；保護區直接取四倍參照。模型、工具、原圖均唯讀，只有painted輸出可寫，容器network none，未呼叫imagegen、Sites或view_image。
- 新入口faithful/index.html只展示許昌地圖城池與據點景觀兩份原圖／修繪小樣。內嵌原圖、參照、修繪及輪廓共8張PNG，可切原版、修繪、並排、拉桿與輪廓保護區；rich2式拉遠拉近、倍率重設、拖曳及雙圖同步保留。這是模型輔助原圖修繪，沒有宣稱人類畫師完成手繪或原版特色已經合格。
- 9項離線瀏覽器抽驗通過，JS錯誤0、外部請求0，含兩份來源、輪廓、各比較模式、拉桿、縮放／拖曳、來源說明、HTML搬移及390×844。五個原始輸入SHA未變，兩張共198,092個保護區像素與四倍Lanczos參照完全相同，單通道改動上限29；此參照不是原版原始像素parity。截圖與模型輸出全留本機，沒有交給模型視覺工具，是否保有原版特色仍由使用者目視確認。
- README、CONTEXT及spec275／索引改指向原圖修繪對照，前兩款SVG保留否決紀錄。正式美術、完整地圖、其他城池與頭像沒有因這兩張小樣宣稱完成；遊戲規則及存檔未改。來源、模型／工具SHA與收據在faithful/qa/verification.json，沒有commit、push或發行，清理結果接續記錄。
- 收尾：HTML 1,340,785 bytes，SHA-256 `11fa4d26dc2c7fa81143fa76a5362192e9a3bb50648debb2ce2e22a59159d20a`。索引582份／135條斷言及git diff --check通過，faithful全目錄1000:1000，工作根無root-owned檔案或.md目錄。全部wolong-faithful及wolong-original-style容器已退出移除，沒有新增image、長駐server或雲端同步。模型與鎖版工具層留在本機供重生，原版特色的視覺驗收仍待使用者確認。

## 2026-11-04：清晰HD與色彩要求，對照rich2原圖重繪流程

- 使用者指出許昌小樣仍糊，並明確要求「基於原版手繪，不要自創，更高清、更有特色、強化色彩，如rich2 HD」。將原先小樣標記為清晰度未通過，9項技術抽驗及保護區相同不再被誤讀為HD美術通過；不重開A／B造型選擇。
- 依原版呈現／縮放路由載入GUI還原及文件職責，唯讀核對rich2的CONTEXT、spec086與terrain-art-generation.json。rich2實際使用內建imagegen對原圖參考重繪，保留輪廓、物件數量、原投影、邊緣交界，另重畫清楚材質及原色深淺；其原圖傳送有自己的明示授權，不能替本專案授權。讀取並套用imagegen技能的原圖編修與內建工具規則，未呼叫圖片工具。
- 模糊根因從程式與實際瀏覽器尺寸核對：原方案約51%／65%像素固定為Lanczos柔邊四倍參照，模型僅以384處理，無法重新畫清楚被鎖住的主要線條；網頁又把448／384原生圖顯示為512，最高約753。這是資產／檢視管線限制，不是新增原版遊戲缺陷。
- faithful/hd-request備妥許昌原版第6頁96×96裁圖來源、768×768最近鄰技術參考圖及完整重繪約束。最近鄰圖只為圖片工具讀原形，沒有新增內容，不稱為新HD成品。request.json要求以原圖唯一決定物件、構圖、視角與形狀，重描原有線條及材質、強化原色，原生1024以上；rich2只供清晰度／手繪完成度／色彩參考，不借用其建築或物件。未實際產生新候選。
- 就先前「不用上拋到OpenAI」提出單一釐清：是否只限HTML不部署，還是原圖也不能交OpenAI圖片工具。過去按後者處理，但不能把agent解讀寫成使用者已明確認定的範圍。問題已明說參考裁圖會交OpenAI處理、生成後保存本機且不發布網站；答覆前不以imagegen或view_image傳送原圖，不用沉默或預選項當作授權。
- 獨立工作僅在斷網Docker準備參考圖、量測原生／網頁顯示尺寸與修正CONTEXT、spec275及索引的目前狀態。沒有圖片工具、雲端影像讀取、Sites、外部發佈或CLI fallback，也未改正式遊戲、原始素材或前版HTML。重繪動作與提示詞已具體可審，待使用者釐清圖片處理範圍後才開啟依賴分支。
- 收尾：尺寸量測448／384→512，最高752.9375，與程式推算相同；結果在hd-request/blur-review.json。索引582份／135條斷言及git diff --check通過，來源及準備資料核SHA存hd-request/verification.json。輸出1000:1000，工作根無root-owned檔案或.md目錄，本輪wolong-hd-sharp／rich2／reference／diagnosis容器均已退出移除，沒有新image或server。圖片處理範圍仍待使用者答覆，沒有默認允許傳圖。

## 2026-11-04：授權雲端圖片重繪，許昌原生HD與本機比較頁

- 使用者已明確釐清「不用部署網頁在OpenAI」「直接在本機端」，並允許「把素材上傳讓雲端手繪」。這取代先前對原圖也不得傳送的解讀。授權原話存於faithful/hd-request/authorization.json；原圖／參考圖可交內建imagegen及影像檢視，成果與HTML留本機，不更新Sites或公開原版素材。沒有再次要求相同授權。
- 沿用已查證的復古GUI／原版呈現路由與imagegen技能，採內建工具模式，不用CLI/API fallback或本機SD1.5續修。先目視許昌KYOGRF第6頁96×96與MMAP城池112×112的最近鄰8倍參考，rich2 atlas僅供清楚筆觸、材質與原色色彩完成度參考。各份完整提示詞與參考SHA由hd-request/request.json及generation-inputs.json索引。
- 景觀第1稿原生1254×1254，屋瓦與石階清楚，但將原版抖色畫成棋盤格，並多出金色牌框。保留原始輸出，針對這兩項編修第2稿；背景改連續深藍、斜坡改赭色材質，保留正面大殿、單層屋頂、金色中央飾物、紅黃藍形狀與三組台階。地圖第1稿同為1254×1254，依方城、四門、建築群及道路方向重繪。繪製幾何不逐像素一致，牆厚、比例與新增材質仍有詮釋，不稱原版美術parity。
- 三張生成PNG從工具的generated_images目錄逐byte複製至faithful/hd-generated，保留工具原檔，沒有後製放大、銳化或覆蓋舊柔化像素。現選景觀r2 SHA e019af55ca00e54808a90652ed81d95fc42ec43cb7481db6cdc573c6f596591e、地圖r1 SHA aae7498646ce4872dea4a16fb99990d51a436b143328d635d9246468037d7131。選稿與人工檢視限制見hd-generated/manifest.json及visual-review.json。
- faithful/index.html改為兩份原圖與兩張原生HD的單檔比較，7,291,268 bytes，SHA c532fb90ff34d541a105726e0ca304804079cb35f49a5b127ef232896072c355。保留原版／HD／並排／拉桿、rich2式縮放、拖曳同步，加入1:1原生細節；移除舊保護區相同的誤導性現行說明。舊HTML、資料與重生程式另存index-local-retouch.html等版本，舊PNG及收據未覆寫。
- 既有Playwright容器斷網抽驗9項通過，JS錯誤0、外部請求0，含四張內嵌PNG原始bytes、兩份來源、比較模式、滑鼠縮放／平移同步、1:1尺寸、來源／授權說明、移動HTML及390×844。實際原圖、新稿及瀏覽器比較畫面已經授權的影像工具目視；技術抽驗不代替使用者的美術確認。結果在qa/native-result.json。
- qa/native-record.py核對五個原版輸入SHA未變、兩張最近鄰參考及rich2參考SHA、三張輸出與工具原始PNG相同；重建HTML／manifest前後SHA相同，前版HTML／manifest SHA亦保留。收據為qa/native-verification.json。所有分析、保存、建頁及瀏覽器工作使用既有有界、1000:1000、network none的Docker容器，沒有新image或模型下載。
- CONTEXT、README及spec275／規格索引改為原生HD候選入口，狀態仍DRAFT。沒有改正式遊戲、規則、存檔或原版資料，也沒有擴展全圖、其他城池或頭像；沒有commit、push、網站部署或發行。全工作根無root-owned檔案或.md目錄，faithful皆1000:1000；本輪工作容器均--rm，主機docker ps確認無wolong容器。使用者下一步可直接開本機HTML確認風格與細部忠實度。

## 2026-11-04：使用者確認HD美術方向

- 使用者對許昌原生HD樣本明確回覆「ok 就是這樣」「這就是我要的」「hd 方向確定」。依文件職責路由與grilling決策後記錄，將地圖r1／景觀r2設為後續HD美術基準。保留原版造型、構圖、視角與配置，重新描繪清楚材質並強化原色。此畫風不再等待決定，不重開前批SVG與柔化小樣選擇。
- 確認原話、範圍、檔名、尺寸與SHA鎖定於faithful/hd-generated/manifest.json的user_approval。兩張PNG原始bytes未改；生成當時的記錄保留，現行manifest與visual-review／request標示方向已確認。CONTEXT、README、spec275及規格索引同步；spec275整體仍DRAFT，正式素材接線、世界鏡頭、AI與編輯器未因此聲稱完成。
- 離線頁面撤掉等待畫風確認的文字，改為已確認HD樣本。重建HTML 7,292,127 bytes，SHA 4a1843e9433bd4b0bc96c1bce6d71d77f8f735a2889cb57efad1eef5cc9eaaa2。既有瀏覽器與來源收據隨新頁面重跑；驗證入口仍為qa/native-result.json及qa/native-verification.json，不新增重複工作文件。
- 沿用有界、1000:1000、network none、--rm的既有Docker工具。沒有重新生成圖片、變更遊戲或原版資料、部署網站、commit、push或發行；素材仍留本機。收尾由同一收據記錄來源完整性、擁有權及容器狀態。

## 2026-11-04：依已確認HD方向擴展四城，共五城十張

- 使用者「ok go」承接已確認的許昌原圖參考HD方向。本輪命中復古重製與文件職責路由，讀取remake／imagegen技能及驗證reference，沿用雲端圖片處理授權與本機HTML交付。沒有重開畫風、工具列、AI曲線或全滑鼠決定，也沒有委派代理。
- 直接使用originals.json中滎陽63、舞陽87、虎牢關72、汝南100的地圖112×112及景觀96×96；景觀頁碼分別14、2、5、0。各自最近鄰8倍參考先目視，許昌已確認HD圖只提供材質、線條及色彩參考。新輸入與SHA在hd-request/batch-2-inputs.json。滎陽北牆封閉及三向城門、舞陽北門、虎牢關狹長關牆、汝南右上農田均按自己的原圖處理。
- 內建imagegen分別產生八張1254×1254原生圖，完整提示詞及有序參考路徑在hd-request/request.json的prompt_set。舞陽地圖另讀滎陽小城稿，保持原圖共同建築的筆觸一致，未借用其河岸或封閉北牆。虎牢關地圖第1稿下方關牆提早收尾，單獨編修第2稿延伸到下緣，選第2稿；兩稿皆保留，共九次新增生成。沒有以反覆重擲挑選偽裝一次生成。
- 八張選稿與九次原始PNG均逐byte複製到faithful/hd-generated，工具原檔保留。每份記錄保存原圖、工具原始路徑、SHA、尺寸與提示詞入口；選稿／修正範圍在hd-generated/batch-2.json。繪製幾何及細節屬原圖參考的美術重繪，不稱原版像素parity或正式地圖幾何。
- 既有build.py改讀原始五城映射，index.html整合十張原圖與十張HD。增加城池與地圖／景觀兩層切換，來源彈窗依當前城池、世界座標、KYOGRF頁碼與圖SHA顯示；保留原版／HD／並排／拉桿、rich2式縮放、同步拖曳與1:1。許昌兩張定稿PNG未變，先前已確認單城頁與manifest另存index-hd-baseline.html及manifest-hd-baseline.json，SHA仍4a1843e9433bd4b0bc96c1bce6d71d77f8f735a2889cb57efad1eef5cc9eaaa2。
- 完全斷網Playwright十一項抽驗通過，JS錯誤0、外部請求0。二十張內嵌PNG bytes及尺寸、五城十組切換、頁碼與來源、各比較模式、縮放及平移同步、1:1、單檔搬移、390×844皆已抽驗；實際虎牢關桌面與手機畫面另目視。技術檢查不代替原版精確幾何或正式引擎驗收。
- qa/native-record.py核對五個原版輸入、十張原圖、八張新參考及十二張累積原生生成圖的SHA；所有本機圖與工具原始PNG相同，許昌已確認兩張圖SHA未變。HTML／manifest乾淨重建前後SHA一致。現行HTML 41,626,539 bytes，SHA 0ab6f124649466d74b3da12a9ef88babc0f18f775fdf1f69fa00d2ca5cfa7234。現行收據仍qa/native-result.json／native-verification.json，單城舊收據另存帶hd-baseline的歷史檔。
- CONTEXT、spec275與規格索引回填五城十張及其來源。規格仍DRAFT；全圖、其餘城池、HD頭像、世界鏡頭與正式素材接線尚未完成。沒有改正式遊戲、規則、存檔或原始資料，沒有部署網站、公開素材、commit、push或發行。只沿用既有rs2-hdtools與Playwright image，沒有下載模型或新建image。斷網工作容器有界且1000:1000，工作根無root-owned檔案或.md目錄，輸出擁有權皆正確；主機Docker核對無wolong容器殘留。

## 2026-11-04：同意五城十張，先行中原大地圖HD

- 使用者「ok 同意」接受第二批四城八張，與許昌兩張合為五城十張已確認基準。依文件職責路由與grilling決策後記錄，將原話、範圍與八張PNG的檔名／SHA寫入hd-generated/batch-2.json的user_approval，主manifest與CONTEXT、spec275同步。所有已確認圖的原始PNG bytes未變，不重開畫風選擇。
- 使用者接著表示「期待大地圖hd」，本輪轉向原版大地圖的可撤回美術與鏡頭預覽。沿復古重製路由使用既有remake／imagegen技能，來源取hd-redraw-local的world-original.png與original-data.json，不增加RE切片或改正式引擎。全圖6144×4096，區域世界格180,96至238,145，原圖928×784，含13據點、27段道路；裁切與原版重建全圖逐像素相同。
- 內建imagegen以區域原圖為唯一內容依據，許昌、舞陽、虎牢關已確認的地圖稿只供清楚材質、線條與原色色彩參考。首張原生1364×1153，SHA d0ea55e67af754381907486defd14c3c3021cccf68279097b333bce46b2fb511，保存為faithful/world/region-hd-r1.png，與工具原檔bytes一致，沒有放大或貼回柔化原圖。提示詞、原生路徑與來源入口在world/generation-r1.json、inputs.json與manifest.json；繪製幾何不稱原版pixel parity。
- world/index.html嵌入原版全圖、區域HD與資料。試繪放回原世界矩形，工具列提供原版／HD、全圖／中原、減號／倍率重設／加號、1:1；支援游標位置縮放、拖曳、點縮圖移動與五城定位。據點名由原資料另繪且可關閉。中原檢視只顯示試繪矩形；全圖導航時範圍外仍為原版，縮圖以金框標示HD範圍。原生尺寸決定放大上限，沒有將全圖原版區域冒稱已完成HD。
- 建頁首輪因容器掛載與host相對路徑不同，誤找/world-original.png。改用明確唯讀/source/world-original.png後同容器建置成功，屬工具掛載問題。游標縮放抽驗首次期望值含小數螢幕座標，但Chromium wheel事件將clientY取整；qa-wheel-diagnosis.json記錄523.6875要求對應523事件，造成0.126原圖像素的期望差。測試改送整數螢幕位置，保留0.1驗證門檻，產品縮放程式未因此調參。
- 來源核對發現初始區域metadata沿用了歷史預覽折線點，現已改為既有original-data.json解析器的完整Path陣列，27段逐項相同。此修正只補來源資料，原圖、生成圖及據點錨點未變；沒有以這些資料聲稱原版執行期找路對拍完成。
- 五城比較頁更新同意紀錄後十一項抽驗再次通過。大地圖八項抽驗通過，兩者JS錯誤0、外部請求0；含內嵌bytes、原圖來源、原座標、原版／HD切換、游標縮放、拖曳、1:1上限、全圖／中原、縮圖與據點定位、城名開關、單檔開啟及390×844。桌面及手機畫面已目視，收據分別在qa/native-result.json與world/qa-result.json。
- world/verify.py核對原圖裁切、13據點錨點、27段解析器Path、原生PNG bytes、八張新同意圖SHA及乾淨重建。大地圖HTML 12,242,411 bytes，SHA 9bd2174dff1d403e784a88427f51fd761197a2554de921aa0954f844c3798f78；五城頁41,626,637 bytes，SHA 4f6becd322eed35a9046636520e7331f045fc57fa3c98835ad59a96e9e03aa87。兩頁重建均同SHA。完整收據為world/verification.json與qa/native-verification.json。
- README、CONTEXT、spec275與規格索引加入大地圖入口，正式規格仍DRAFT。完整世界HD、其餘素材及Ebiten鏡頭／素材接線尚未完成；沒有改遊戲、規則、存檔或原始資料，沒有部署、公開素材、commit、push或發行。工具沿用既有rs2-hdtools與Playwright，容器有界、1000:1000、network none且--rm；來源唯讀，工作根無root-owned檔案或.md目錄，收尾確認無專案容器殘留。

## 2026-11-04：大地圖向四方擴展，五區47據點及交界補繪

- 使用者「go」承接已展示的大地圖方向。命中復古重製與文件職責路由，讀取既有remake／imagegen技能與文件規則，沿用原圖交雲端重繪及本機HTML的授權，沒有重問畫風或圖片工具範圍。工作仍為DRAFT原型，未改正式引擎或啟用未審核規則。
- 從原版6144×4096重建圖取中原北、東、南、西四區，各附四格共同邊緣作生成參考。顯示範圍北180,47至238,96，東238,96至296,145，南180,145至238,194，西122,96至180,145；原圖均與全圖裁切相同。加原中原區共47個來源據點、3,637,760原圖像素，約14.455%全圖；原生圖只供美術顯示，不作正式道路或碰撞資料。
- 內建imagegen各繪一區，原生北／南1312×1199、東／西1411×1115。北側第1稿借入原圖沒有的農田格，第2稿針對此項移除，選第2稿。原生PNG直接複製工具輸出，來源及SHA在world/expansion-inputs.json、各區r1/r2.json；完整提示詞由regions.json索引。四個方向皆逐張看過原圖與新稿，中原原生PNG的SHA d0ea55e67af754381907486defd14c3c3021cccf68279097b333bce46b2fb511未變。
- world/build.py及app.js改為五區清單，顯示各自原世界矩形，依顯示範圍裁掉共同邊緣，避免重複覆繪區塊。增加區域選擇、本區／五區／全圖，保留游標縮放、平移、縮圖與據點定位。原生上限按所有區域及補繪圖共同像素密度計算，按鈕改標「原生細節」。城名以實際字寬做邊界與互相避讓，資料保留全部47個據點。
- 先以正常滑鼠平移到四處交界，目視發現東側河岸台階及西側同一小城重複片段。從原圖取兩條384×784走廊，並以瀏覽器離屏Canvas擷取當前拼接作參考；這是預覽渲染擷取，不作新美術成品。兩條交界再由內建imagegen補繪，原生東878×1792、西1013×1552，PNG均未後製。檢視器在中間192原圖像素範圍使用補繪圖，外緣8原圖像素作透明度接合，無模糊濾鏡。主要缺陷已改善；仍有局部筆觸、配色及邊線差異，不宣稱最終無縫或原版幾何精確。
- 十項斷網Playwright抽驗通過，JS錯誤0、外部請求0。核對全圖、五區及兩條交界嵌入bytes／尺寸、五區切換及世界座標、原版／HD、游標縮放、平移、共同原生上限、全圖／五區／本區、縮圖、據點定位、城名避讓／開關、搬至不同目錄仍可單檔開啟，以及390×844。四個交界均經滑鼠到場抓圖，東西交界修正後與手機畫面再目視；結果與限制在world/qa-result.json、seam-review.json。
- world/verify.py核對五張原圖裁切、47個據點錨點及各區完整來源Path、六次區域生成及兩次交界生成的原生bytes。七張選稿不作原版像素parity。五城已確認圖及中原PNG均未改；中原舊頁另存index-central-r1.html，連同舊程式可重建為原SHA 9bd2174dff1d403e784a88427f51fd761197a2554de921aa0954f844c3798f78。現行五區HTML 39,300,402 bytes，SHA 55a4e0ba90d524d84427b625308fe4511a4909a9f5c1ffa9635fda16aa7d8853，乾淨重建相同，收據為world/verification.json。
- CONTEXT、README、spec275及索引更新五區入口。完整世界HD、對角區域、局部交界品質及正式Ebiten接線仍未完成；沒有修改遊戲、規則、存檔、原始資料，也沒有commit、push、公開素材或網站部署。沿用既有rs2-hdtools與Playwright映像，沒有新建image或下載模型；容器有界、1000:1000、network none且--rm，來源唯讀。所有輸出擁有權正確，工作根無root-owned檔案或.md目錄，主機核對無wolong容器殘留。

## 2026-11-04：補齊大地圖四角，九區85據點及橫向交界補繪

- 使用者「ok go」接續既有大地圖方向。本輪命中復古重製與文件職責路由，使用remake／imagegen技能及抽驗reference；沿用原圖可交雲端圖片工具、成果與HTML留本機的授權，未重開畫風、工具列或操作決定，也未委派代理。正式引擎仍在spec275 DRAFT，沒有新增逆向切片或玩法規則。
- 從原版6144×4096重建圖裁取四角，各在兩個相鄰方向帶四格參考邊緣，原圖均992×848。顯示範圍西北122,47至180,96、東北238,47至296,96、西南122,145至180,194、東南238,145至296,194；各含3／14／12／9個原版據點。與既有五區合為九個不重疊矩形，完整覆蓋2784×2352矩形、85據點、6,547,968原圖像素，約26.0193%全圖。原座標、完整來源Path與裁圖SHA在world/diagonal-inputs.json。
- 內建imagegen以每區原圖作唯一內容依據、相鄰兩區HD作材質與色彩參考，先生成四張1356×1159。東北、東南第1稿誤加金色田格，針對原圖沒有田格的位置修出第2稿1357×1159；西北原圖沒有農田，選稿維持無田格。共六次區域生成，全部原始PNG保存本機且與工具bytes一致。完整提示詞及選稿由world/diagonal-generation-set.json索引，不將手繪邊線當作原版精確幾何。
- 五區HTML、manifest、清單與收據另存帶five-r1的快照，SHA入口為world/five-r1-snapshot.json；五區HTML仍為55a4e0ba90d524d84427b625308fe4511a4909a9f5c1ffa9635fda16aa7d8853。九區頁增加四角選擇，已繪範圍從十字形補為矩形；標題、據點數、面積與來源說明改從清單計算，避免下一批仍顯示五區文字。保留全滑鼠縮放、平移、縮圖、據點定位與共同原生細節上限。
- 正常滑鼠平移抽驗四個交會處，發現橫向河岸斷差及道路偏移。取四條1024×384原版走廊，以瀏覽器離屏Canvas擷取目前拼接作2倍編修參考，並用內建imagegen補繪四張2048×768原生條帶。只在中央1024×192顯示矩形合成，四邊各8原圖像素作透明度接合，沒有模糊濾鏡或修改PNG。東北補繪第1稿誤加向右的L形道路，追加第2稿移除錯路並恢復大城向南直通小城；兩稿均保存。共五次新增交界生成，沿用舊兩條縱向補繪，共九張區域與六張交界選稿。
- 修正前後四個交會畫面已目視，主要河岸斷差與道路錯接改善。局部筆觸、配色及邊線仍有差異，含部分縱向交界；不宣稱最終無縫、原版像素parity、完整世界HD或新大地圖已獲使用者逐張接受。來源、選稿、修正原因與畫面入口在world/diagonal-review.json，既有seam-review.json保留前批記錄。
- 最終十項斷網Playwright抽驗通過，JS錯誤0、外部請求0，涵蓋全圖及十五張選稿內嵌bytes／尺寸、九區切換、原版／HD、世界座標、游標縮放、拖曳、原生上限、已繪範圍／本區／全圖、縮圖、據點定位、城名避讓／開關、來源說明、單檔搬移及390×844。四個中央交界及四個四區交會處以正常滑鼠到場保存畫面；沒有以測試注入鏡頭代替滑鼠操作。
- world/verify.py核對九張原圖裁切、85據點及完整來源道路，累積十二次區域生成與七次交界生成共十九張原生輸出bytes。五城十張已同意PNG及城池比較頁SHA未改，中原原生PNG未改；中原歷史HTML與九區現行HTML重建同SHA，五區歷史快照未變。九區HTML 73,471,787 bytes，SHA 3ce9aab1833bd49f5cb4696862b7c39e7bdff30b9cb2eac20530453914e43c5b；收據為world/qa-result.json及verification.json。
- CONTEXT唯一狀態表、README穩定入口、spec275與索引更新九區範圍；後續仍需其餘世界素材、局部邊線修整、HD頭像及正式Ebiten素材／鏡頭接線。HEAD仍bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行或網站部署。沿用既有rs2-hdtools與Playwright映像，沒有新建映像或下載模型；Docker工作有界、1000:1000、network none、來源唯讀且--rm，所有輸出擁有權正確，工作根無root-owned檔案或.md目錄，收尾檢查專案容器清理狀態。

## 2026-11-05：向西擴展三區，十二區102據點

- 使用者「go」接續原圖參考HD大地圖。本輪命中復古重製、文件職責與README摘要路由，載入remake／imagegen技能及抽驗reference；沿用圖片可交雲端重繪、成果與HTML留本機的授權，未重開畫風或操作決定，也未委派代理。只擴展可撤回的本機DRAFT預覽，沒有改正式引擎、規則或存檔。
- 原版全圖唯讀裁出西北外緣64,47至122,96、西側外緣64,96至122,145、西南外緣64,145至122,194，皆為世界格半開範圍。參考圖包含右側及上下相鄰邊緣，原尺寸992×848／992×912／992×848；新增4／10／3個據點，包括長安及漢中。十二區不重疊且填滿3712×2352矩形，共102據點、8,730,624個原圖像素，約全圖34.6924%。裁圖、原座標、完整來源Path與SHA在world/west-column-inputs.json。
- 內建imagegen逐區依自己的原圖重繪，東鄰已繪HD僅供筆觸、原色及邊界參考。西北／西南選稿1356×1159；西側第1稿多畫上方河岸田格，局部第2稿移除，選1309×1202。四次區域原生輸出全部保存，PNG與工具原始bytes相同，提示詞與選稿由world/west-column-generation-set.json索引。既有九區PNG未改，九區HTML、manifest、清單、程式與收據保存為nine-r1快照，雜湊入口world/nine-r1-snapshot.json。
- 先用正常滑鼠巡查三處縱向及兩處橫向新交界，發現中段／南段河岸斷差、支流與山林錯位。從原版取兩條384×784縱向及兩條1024×384橫向走廊，以當前組版的2倍離屏擷取作編修參考；再由內建imagegen生成兩張878×1792及兩張2048×768原生補繪。顯示沿既有core裁切與四邊8原圖像素透明度接合，沒有模糊或PNG後製。來源在west-column-seam-inputs.json，計畫在west-column-seam-plan.json，capture-seam-references.cjs保存參考擷取流程。
- 修正前後四處畫面已目視，主要河岸斷差及山林錯位改善。北段縱向交界沒有本次抽樣的重大連接缺口，局部筆觸及色差保留為可選修整；其他局部邊線差異仍在，不宣稱最終無縫或原版逐像素幾何。選稿、修正與前後畫面入口在world/west-column-review.json。新圖沿既有畫風產製，不將「go」記成新大地圖已獲逐張美術驗收。
- 最終十項斷網Playwright抽驗通過，JS錯誤0、外部請求0。涵蓋全圖及十二區／十條交界的內嵌bytes與尺寸、十二區切換、原版／HD、世界座標、游標縮放、拖曳、共同原生上限、已繪範圍／本區／全圖、縮圖、據點定位、城名避讓與開關、來源說明、單檔搬移及390×844。新五處交界均以正常滑鼠到場抓圖，四處修正後及長安／漢中區域畫面再次目視；沒有注入鏡頭代替操作。
- world/verify.py核對十二張原圖裁切、102據點與原始包圍範圍內據點集合、完整來源Path，累積十六次區域生成及十一次交界生成，共二十七張原生輸出，二十二張選稿。五城十張已同意PNG及城池比較頁SHA未改，中原PNG未改；中原歷史頁與現行頁重建同SHA，九區及五區歷史快照未變。現行HTML 103,388,895 bytes，SHA 8dcb1ffd82a1dc1dcb9fec611307b64ee661e1fb957c62cdf58174ac2fccb0cd；收據仍world/qa-result.json及verification.json。
- CONTEXT唯一目前狀態表、README穩定入口、spec275與索引更新十二區範圍。其餘世界美術、局部邊線修整、HD頭像及正式Ebiten資產／鏡頭接線尚未完成，spec275維持DRAFT。HEAD核對仍bd09fd40e344efee9d01e286a70f2a2303890caf，保留原有工作樹，沒有commit、push、發行或網站部署。Python工具容器未附rg，改用容器內Python作少量檔案搜尋，未當成產品缺陷。沿用既有rs2-hdtools與Playwright，沒有新映像或模型下載；Docker工作有界、1000:1000、network none、來源唯讀且--rm。擁有權與工作根root-owned／.md目錄檢查通過，收尾另核對容器清理與文件差異。

## 2026-11-05：向東增加三區，十五區118據點及單檔載入修正

- 使用者回覆「ok 繼續囉 很棒」，沿既有原圖參考HD方向續做。命中復古重製、文件職責及README摘要路由，沿用remake／imagegen技能與抽驗reference。這是正向回饋與進度授權，沒有擴張成每張新地圖已獲逐張驗收；未重開已確認畫風，也未委派代理。
- 原版全圖唯讀裁出東北外緣296,47至354,96、東側外緣296,96至354,145、東南外緣296,145至354,194，各有相鄰邊緣參考。原尺寸992×848／992×912／992×848，新增5／3／8據點，包括臨淄、北海、廣陵、建業、會稽；十五區填滿4640×2352矩形，共118據點、10,913,280原圖像素，約全圖43.3655%。裁圖、原座標、完整來源Path及SHA在world/east-column-inputs.json，十二區歷史頁及程式／metadata保存twelve-r1快照。
- 內建imagegen依各區原圖決定內容，西鄰HD只供筆觸、材質及原色參考，保留原有海岸、河口與湖泊，平原不新增山林。東北第1稿多畫金色田格，局部第2稿移除，選1357×1159；東側1308×1202、東南1356×1159。四次區域原生輸出均保存，PNG與工具原始bytes相同，完整提示詞與選稿由world/east-column-generation-set.json索引。
- 首次十五區瀏覽器抽驗已走完區域及滑鼠接縫取樣，在搬至/tmp後同頁重開逾時。qa-load-probe.cjs在相同2 CPU／3GB容器重現：115.8MB單一JSON內嵌頁初次就緒約20–25秒，同頁重開超過30秒，cgroup無OOM。建立build-payload-probe.py原型，將PNG bytes分別放在不執行的ASCII script區塊，元資料只存ID；載入後清空原文字區塊，整組圖片仍全部就緒後才顯示。相同門檻與資源，原型初次約16–17秒、重開約7.5秒。原型與前版HTML均保留，量測在qa-load-probe.json、qa-load-probe-chunked.json，方法與SHA在payload-loading.json；不把原型時間冒稱最終129.7MB頁面或所有機器的效能。
- 分塊方式已接回build.py／template.html／app.js，PNG未放大、轉碼或改bytes，正常預覽仍為單檔、CSP禁止外部連線。接縫參考擷取初次批次傳送大量base64時renderer關閉，原因未確證；改讀file URL後遇到canvas origin限制。最後改為先逐檔核對SHA，再由同一Node程序在斷網容器loopback提供同源PNG，避免大量除錯通道複製。沒有主機port或外部介面，finally關閉伺服器；這只產生補繪參考，未部署預覽網站，兩次工具失敗不當作遊戲缺陷。
- 正常滑鼠巡查東側三處縱向與兩處橫向接縫，再按原圖走廊補繪東北大城／南側兩條縱向及上下兩條橫向，原生分別878×1792與2048×768。四次交界輸出均保存，沿既有顯示矩形與四邊8原圖像素透明度接合，沒有模糊或PNG後製。前後畫面目視確認主要城池碎片、河岸與海岸斷差改善；中段縱向筆觸、色差及其他局部邊線仍列可選修整，不稱最終無縫。來源在east-column-seam-inputs.json，選稿／限制／前後畫面在east-column-review.json。
- 最終十五區十項斷網Playwright抽驗通過，JS錯誤0、外部請求0，保留原有30秒載入門檻。涵蓋十五區與十四張交界選稿、原圖內嵌bytes／尺寸、圖片載入後文字區塊釋放、區域切換、原版／HD、世界座標、游標縮放、拖曳、共同原生上限、已繪範圍／本區／全圖、縮圖、據點定位、城名避讓與開關、來源說明、單檔搬移後同頁重開及390×844。四處接縫修正後、建業／會稽所在區域與手機尺寸畫面已目視；沒有注入鏡頭代替滑鼠到場。
- world/verify.py核對十五張原圖裁切、118據點及包圍範圍原資料集合、完整來源Path，累積二十次區域與十五次交界生成共三十五張原生輸出，二十九張選稿，HTML共內嵌三十張PNG。五城十張已同意PNG、城池比較頁、中原PNG及十二區／九區／五區歷史快照未改；中原歷史與現行頁重建同SHA，載入比較證據也核對通過。現行HTML 129,723,954 bytes，SHA 5e2cd38af2343325825cc25a666a08fe0bf1973bbdb62a6583deb46a8d5d5b3e，收據仍qa-result.json及verification.json。
- CONTEXT唯一目前狀態表、README穩定入口、spec275與索引更新十五區範圍及載入證據。其餘世界、局部邊線、HD頭像與正式Ebiten素材／鏡頭接線尚未完成，spec275仍DRAFT；沒有更動正式規則、資料或存檔格式。HEAD仍bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、公開素材、發行或網站部署。沿用既有工具映像，沒有新增映像或模型下載；Docker工作有界、1000:1000、network none、來源唯讀且--rm，輸出擁有權正確，收尾核對工作根root-owned／.md目錄與容器清理。

## 2026-11-05：北側五區接到原圖上緣，二十區131據點

- 使用者「go」接續已確認的原圖參考HD方向。本輪命中復古重製、文件職責及README摘要路由，沿用remake／imagegen技能與抽驗reference。圖片工具授權沿既有明示決定，不重開畫風、不部署網站，也沒有委派代理；正式引擎仍在spec275 DRAFT。
- 原版全圖裁取北側五區，世界格半開範圍依序為64,0至122,47、122,0至180,47、180,0至238,47、238,0至296,47、296,0至354,47。各區向南及相鄰列帶四格參考，原尺寸992×816或1056×816，新增0／3／4／5／1個據點，包括晉陽、中山、真定、南皮與黃；原名「北京」沿資料保留。二十區不重疊且填滿4640×3104矩形，共131據點、14,402,560原圖像素，約全圖57.2306%。來源與完整Path在world/north-band-inputs.json。
- 西段原圖沒有據點，保留兩處水塘、河道、山林及細長城牆，不增建城堡；中段與東段不添農田。內建imagegen逐區參考自己的原圖，相鄰HD只供筆觸、材質與原色參考。西段／海岸選稿1383×1137，中西段／中段／東段1426×1103；海岸第1稿多畫金色田格，第2稿移除並保留真正南側綠色田格。六次區域輸出及完整提示詞保存本機，入口north-band-generation-set.json；既有十五區PNG未改，歷史頁、程式、metadata與收據保存fifteen-r1快照。
- 首次151,122,438 bytes頁面通過正常操作與交界擷取，但搬移至/tmp後同頁重開超過既有30秒門檻。build-north-loading-probe.py保存初稿及失敗，再建立原始PNG bytes轉暫存圖片網址的可撤回原型；讀取後立即清除base64文字區塊，成功／失敗即釋放網址，整組素材仍全部載入才顯示。相同2 CPU／3GB與30秒門檻，原型三次就緒約6.3–7.3秒；原方式的獨立重開測試三次也通過，沒有OOM事件，因此不宣稱已證實唯一原因。長data URL與圖片快取保留只列假說。方法、SHA、失敗及兩側量測在north-loading.json，未以加長逾時或提高資源通過。
- 以正常滑鼠巡查九處新交界及兩處舊交界。中段出現重複城池與分離屋頂，東段大城及海岸也有橫向斷差；取兩條1056×384及一條992×384原版走廊，以2倍離屏組版作編修參考，交內建imagegen補繪2079×756／2079×756／2016×780。原生PNG不後製，顯示沿中央192像素走廊與四邊8原圖像素透明度接合。前後畫面目視確認主要碎片與斷差改善，其他局部配色、筆觸及縱向交界仍留可選修整，不稱最終無縫。來源、計畫及檢視在north-band-seam-inputs.json、north-band-seam-plan.json、north-band-review.json。
- 最終十項斷網Playwright抽驗通過，JS錯誤0、外部請求0。涵蓋原版與三十七張選稿的內嵌bytes／尺寸、二十區逐一切換、世界上緣完整可見、零據點區域清空按鈕、原版／HD、世界座標、游標縮放、拖曳、共同原生上限、已繪範圍／本區／全圖、縮圖、據點定位、城名避讓／開關、來源說明、完整操作後搬移單檔同頁重開及390×844。北側區域、三處修正後交界與窄視窗已目視；正常操作沒有以鏡頭注入代替。原有世界邊界限制維持，fit檢查改驗整個目標矩形可見，容許y=0處鏡頭中心被邊界限制。
- world/verify.py核對二十張原圖裁切、131個據點及原資料包圍範圍集合、完整來源Path、二十六次區域與十八次交界生成，共四十四次原生輸出bytes；選稿三十七張，HTML共內嵌三十八張PNG。五城十張已同意PNG、城池比較頁、中原PNG及十五區／十二區／九區／五區快照未改；中原歷史頁與現行頁重建同SHA。現行HTML 162,275,252 bytes，SHA a7b6b48bb3c4e452f952ed3869ce0ea0222349bc25bc0dada405f253039434fc，收據為qa-result.json與verification.json。
- CONTEXT唯一狀態表、README、spec275及索引更新二十區入口，沒有把「go」記成新增圖片逐張已獲驗收。其餘世界、頭像與正式Ebiten素材／鏡頭接線仍待完成；未更動玩法、存檔或原始資料，沒有commit、push、發行或網站部署。沿用既有工具映像，容器有界、1000:1000、network none、來源唯讀且--rm；所有輸出擁有權正確，工作根root-owned／.md目錄檢查通過，收尾核對專案容器清理與git diff --check。

## 2026-11-05：南側五區接到原圖下緣，二十五區158據點

- 使用者「ok 繼續啦」延續已確認的原圖參考HD方向。命中復古重製、文件職責及README摘要路由，載入remake／imagegen技能及抽驗reference。先讀CONTEXT及主機gh的遠端open Issues，這批仍在spec275本機DRAFT範圍，未修改Issue、正式引擎或存檔。沒有重開畫風決定、委派代理或部署網站。
- 南側五區世界格半開範圍依序為64,194至122,256、122,194至180,256、180,194至238,256、238,194至296,256、296,194至354,256。原圖向北及相鄰列帶四格參考，尺寸992×1056或1056×1056，新增3／5／10／6／3個據點，包括長沙、武陵、桂陽、柴桑與建安。二十五區填滿4640×4096矩形，共158據點、19,005,440原圖像素，約全圖75.5208%；原圖左右外側帶尚未重繪。來源、原座標及完整Path在world/south-band-inputs.json。
- 內建imagegen逐區依自己的原圖重繪，北鄰HD只作材質、筆觸與原色參考。西段沒有農田，保留三座小城、密集低矮山林與溪流；其他區保留原有湖泊、島岸與綠色田格。西段原生1235×1274，中西段／中段／東段1254×1254，海岸1216×1294。中段第1稿把右上河邊橘黃色抖色畫成金色田格，第2稿針對此項修正；六次區域輸出全部保存，選稿及完整提示詞在south-band-generation-set.json。既有二十區PNG未改，二十區HTML、metadata、程式及收據保存twenty-r1快照；歸檔程式是歷史證據，未宣稱可直接重建二十區舊頁。
- 原生輸出複製第一次缺少腳本既有的/output掛載別名，尚未寫入素材就停止。補正同一工作目錄掛載後，在原容器工具鏈重新執行，屬工具呼叫問題，沒有改遊戲或重建映像。
- 初次二十五區187,239,501 bytes頁面十項抽驗通過。正常滑鼠巡查九處新交界及兩處舊交界，西段城鎮／河岸／山林與中段湖岸／小城有明顯裁切斷差。取992×384與1056×384原圖走廊，依當前組版2倍參考補繪2014×780及2079×756。沿中央192像素顯示走廊及四邊8原圖像素透明度接合，PNG原生bytes未後製。前後目視確認主要斷差改善；局部配色、筆觸及縱向交界仍留可選修整，不宣稱最終無縫。入口south-band-seam-inputs.json、south-band-seam-plan.json及south-band-review.json。
- 最終頁十項抽驗通過後，將瀏覽器收據補強為固定起始HTML雜湊，流程結束確認頁面未改，並在verify.py比對受測選稿數與清單；完整流程依相同門檻再次乾淨通過。這避免建頁與驗證若交錯時，把流程末端另一份檔案SHA誤當受測檔。沒有以延長逾時或增加容器資源通過，也沒有把工具時序疑慮當成產品缺陷。
- 最終斷網Playwright十項抽驗包含原圖與四十四張選稿的內嵌bytes／尺寸、二十五區切換、上下世界邊界及零據點區、原版／HD、游標縮放、拖曳、共同原生上限、本區／已繪範圍／全圖、縮圖、據點定位、城名避讓／開關、來源說明、正常滑鼠交界擷取、操作後搬移單檔同頁重開及390×844。JS錯誤0、外部請求0。總覽可能比原世界更寬，游標焦點檢查先透過工具列選中原再量測，保留既有世界邊界限制；沒有為測試更改鏡頭規則。五區、兩段修正後交界及窄視窗均已目視。
- world/verify.py核對二十五張原圖裁切、158據點與原資料包圍範圍集合、完整來源Path、三十二次區域及二十次接縫生成，共五十二次原生輸出bytes。選稿四十四張，HTML共內嵌四十五張PNG。十張已同意城池圖、城池比較頁、中原PNG及二十區／十五區／十二區／九區／五區快照未變；中原歷史頁與現行頁重建同SHA。最終HTML 195,611,745 bytes，SHA 46fafd0c3f46905d6db92f0ba77ecf3079e9a79672c457ffa2ed2b5099700c54，收據為qa-result.json及verification.json。
- CONTEXT唯一現況表、README、spec275與索引更新二十五區範圍。剩餘左右外側帶、局部邊線、頭像及正式Ebiten素材／鏡頭接線仍待完成；新圖不記成逐張已獲使用者驗收。HEAD仍bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行或網站部署。沿用既有rs2-hdtools與Playwright映像，未下載模型或新增映像。Docker工作有界、1000:1000、network none、來源唯讀且--rm；輸出擁有權及工作根root-owned／.md目錄核對通過，收尾檢查容器清理與git diff --check。

## 2026-11-05：補左右十區，三十五區試繪覆蓋全部原圖

- 使用者「ok go」延續原圖參考HD方向。先讀CONTEXT、復古重製路由、remake／imagegen技能及抽驗reference，並依主機gh規則核對遠端open Issues。沿已確認的畫風補左右兩側，不開新玩法、架構或公開素材決定，也沒有委派代理。正式接線仍在spec275 DRAFT。
- 西側五區世界格x=0至64，東側x=354至384；y依序為0至47、47至96、96至145、145至194、194至256，皆為半開區間。來源帶相鄰四格參考且不越過原圖邊緣，西側原尺寸1088×816／912／1056，東側544×816／912／1056。西側新增0／4／11／11／5據點，東側0／0／0／1／2，共34據點，包括西涼、成都、吳、章安及溫州。三十五區不重疊且完整覆蓋6144×4096、25,165,824原圖像素及全部192個原資料據點。裁圖、原座標及完整Path在world/outer-bands-inputs.json。
- 內建imagegen逐區依自己的原圖重繪，相鄰HD只作材質、筆觸與原色參考。西北零據點區保留三處水域與細長城牆，關隘沿原圖保留窄關牆；東側三個零據點區保留空海面及原岸線。西側選稿1448×1086、1370×1148、1371×1148、1370×1148、1274×1235；東側1024×1536、三張968×1625及900×1747。西側中段與東北海岸第1稿誤加金色田格，各選第2稿修正；東側中段選第3稿減少粗水紋，局部水紋色差仍留待修整。共十四次區域生成，全部原生PNG與提示詞保存，入口outer-bands-generation-set.json。
- 既有二十五區PNG未改。二十五區HTML、程式、metadata及收據保存twentyfive-r1快照，包含原QA及verify腳本；歷史程式是證據副本，未宣稱可直接重建舊範圍。全圖組版後改顯示「全圖HD試繪」與局部待修說明，原版模式仍可看完整原圖；沒有更動原世界座標、命中規則或正式引擎。
- 首次238,437,142 bytes全圖頁完成區域及左側交界擷取，在由最西跨最東的正常滑鼠巡查碰到舊32段拖曳上限。左側七張擷取存在，第一張右側擷取未到場，失敗保存在qa-outer-first-failure.txt。測試改依距離、倍率及每段滑鼠距離計算有界步數，並將要求鏡頭換算到原世界合法邊界；未改產品鏡頭規則，也未直接注入鏡頭。相同2 CPU／3GB、30秒載入及180秒外層門檻乾淨重跑通過。
- 正常滑鼠抽驗十四處新交界與兩處既有交界，另以縮圖到四個世界角落。西側中南段宕渠有重複城鎮片段，兩處橫向邊界切到屋頂與河岸；依原圖走廊補繪878×1792縱向及兩張2141×734橫向圖，沿中央192像素走廊、四邊8原圖像素透明度接合，沒有模糊或PNG後製。主要重複城鎮及岸線斷差已改善。中北／中段交界的孤立屋頂細修第2稿產生額外小物件，未選用；保留較忠實第1稿並明列殘留屋頂。四次交界生成都保存，入口outer-bands-seam-inputs.json及outer-bands-review.json。
- 最終十項斷網Playwright抽驗通過，JS錯誤0、外部請求0。涵蓋原圖及五十七張選稿的內嵌bytes／尺寸、三十五區切換、零據點區、四邊／四角、原版／HD、游標縮放、拖曳、共同原生上限、本區／已繪範圍／全圖、縮圖、據點定位、城名避讓／開關、來源說明、十六處交界擷取、操作後搬移單檔同頁重開及390×844。收據固定起始HTML雜湊，流程中頁面未改，選稿數與清單相同；保留原門檻與資源。全圖、十區、三段修正後西側交界、東側海面與窄視窗已目視。
- world/verify.py核對三十五張原圖裁切、全部192據點及全圖原資料集合、完整來源Path、四十六次區域及二十四次交界生成，共七十次原生bytes。選稿五十七張，HTML共內嵌五十八張PNG。十張已同意城池圖、城池比較頁、中原PNG及二十五區／二十區／十五區／十二區／九區／五區快照未變，中原歷史頁與現行頁重建同SHA。最終HTML 250,657,792 bytes，SHA 506c1474cbb397885184e57e2d9c1f89a4c7fdf63a4ea34b7e00d87976189aab；收據為qa-result.json及verification.json。
- CONTEXT唯一現況表、README、spec275與索引更新為全圖面積覆蓋完成。這不代表完整HD或引擎接線完成，也不記成新圖逐張已獲使用者驗收；孤立屋頂、局部配色／水紋／邊線、頭像與正式資產映射仍待完成。HEAD仍bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行、改Issue或部署網站。沿用既有工具映像，Docker有界、1000:1000、network none、來源唯讀且--rm，原生圖及輸出擁有權核對通過，工作根無root-owned檔案或.md目錄；收尾核對容器清理、來源索引及git diff --check。

## 2026-11-05：修正全圖孤立屋頂，新增三張原版參考HD頭像

- 使用者「ok 繼續啦～」延續原圖參考HD方向。載入復古重製路由、remake及imagegen技能、文件職責與README規則，先讀目前CONTEXT及現行完整性收據。主機gh auth成功；首次Issue查詢誤用目錄名dragon，依git remote修正為wicanr2/wolong_cht後成功核對遠端open Issues。未修改Issue、正式引擎、規則或存檔，沒有委派代理或重開已定畫風。
- 全圖原圖面積覆蓋已完成，不重做三十五張區域圖。針對西側中北／中段交界的孤立屋頂，先建立洋紅框定位參考，再將完整補繪、框選參考及原圖交內建imagegen。第3稿移除孤立藍屋頂與短底座，沒有換成額外橋或小屋；完整城鎮、河道、道路及真正田格保留。2141×734原生PNG bytes未後製，SHA b747adf5cb4f21aededc36e40ea6812c48cb898a9edbaa8f95ef69be4fa1ce7a。舊第1稿、第2稿及生成紀錄均保留，修正前全圖、metadata、程式、QA、完整性收據及局部畫面保存thirtyfive-r1快照。入口world/detail-polish-review.json。
- 頭像沿既有KAOGRF解碼原圖取曹操16／50、關羽40／114、諸葛亮60／125，核對圖庫150頁、色盤、劇本及三張64×64原圖雜湊。最近鄰16倍圖只作雲端參考。三張HD皆1254×1254原生重繪，沿各自原圖保留黑底、臉型、神情、頭部構圖及原色色系；諸葛亮首稿太正面且多金色衣肩，第2稿恢復較明確左轉角度與深藍衣肩。四次原生PNG、原圖、提示詞與選稿均保存於既有faithful工作區下portraits，入口originals.json、generation-set.json、review.json。這是三頁小樣，未宣稱150頁完成或已獲使用者美術驗收。
- 建立portraits/index.html本機單檔，可選三位武將，切並排／拉桿／原版／HD，同步縮放、拖曳及重設；工具列固定，HD 1:1依螢幕像素密度限制原生上限。原圖放大保留像素邊緣，HD重繪沒有經放大或後製。建頁及抽驗程式在同目錄build.py、qa-check.cjs、verify.py，準備及原生複製腳本沿用faithful/tooling。新增入口已掛入README、CONTEXT及spec275。
- 首次原生複製少傳docker的-i，兩支讀stdin腳本在寫入前取得空字串而停止。補回stdin傳遞後用相同容器乾淨重跑成功，歸類工具呼叫問題，沒有修改圖片或遊戲，也未更換映像。
- 頭像七項斷網瀏覽器抽驗通過，包括武將／頁碼及六張內嵌PNG bytes／尺寸、三位正常切換、四種顯示模式、滑鼠拉桿、滾輪焦點／工具列縮放、原生上限、同步拖曳／重設、來源說明、單檔搬移重開及390×844兩倍像素密度。JS錯誤0、外部請求0。三張1254×1254及原圖並排、拉桿、細節與窄視窗已擷取；曹操、諸葛亮並排及窄視窗已目視。HTML 8,975,872 bytes，SHA 4c0125d7426b8328fcc2847ecc190b04a66f43a1f46da614816f526f75f5d242，重建相同。完整性收據核對四次原生bytes、三張原圖、原始資源雜湊及十張已同意城池PNG／城池頁未變。
- 大地圖依原30秒載入、180秒外層及2 CPU／3GB完成原十項斷網抽驗，JS錯誤0、外部請求0，修正後屋頂及宕渠交界已目視。最終HTML 250,663,972 bytes，SHA 4c6e4e49beaca2cc45863d9a5a5af6e2dd93bd5f4222dba553a017d84e42efde。完整性核對三十五張原圖、全部192據點、完整Path、71次原生輸出、57張選稿、歷史快照及重建一致。保留full_world_HD_complete=false與engine_integration=false，局部配色／水紋／邊線、其餘頭像與正式接線仍待完成。
- CONTEXT唯一現況表、README、spec275及索引更新目前入口與限制。未新增正式交付目錄，沒有commit、push、發行或網站部署。沿用rs2-hdtools:r1與既有Playwright映像，所有工作有界、1000:1000、network none、來源唯讀且--rm；原生輸出與檔案擁有權正確，工作根無root-owned檔案或.md目錄。收尾核對文件索引、git diff --check及專案容器清理。

## 2026-11-05：新增六張原版頭像，九張本機對照及趙雲判讀勘誤

- 使用者「ok go」延續原圖參考HD方向。先讀CONTEXT、復古重製路由、remake／imagegen技能、抽驗reference及文件職責，更新首頁前核對README規則。主機gh auth成功並讀取wicanr2/wolong_cht的遠端open Issues；未寫入Issue。這批只擴充既有本機美術原型，沒有變更正式規則、存檔或引擎，沒有委派代理，也未把普通繼續當成新畫風的逐張接受。
- 沿用既有faithful/portraits及faithful/tooling，沒有建立新根目錄。先將三張頁、metadata、程式、QA、完整性收據及畫面保存three-r1快照，並鎖定三張既有選稿SHA。以既有hr-go-ebiten:1.26.7-2.9.9-r1容器執行extract-portrait-batch2.go，唯讀專案及原資料，只有頭像工作目錄可寫，既有wl-gomod唯讀，建置快取在容器/tmp。核對原圖實際150頁，新增劉備32／94、孫權18／128、司馬懿55／9、張飛80／118、趙雲73／6、周瑜57／110，數字依序為武將ID／頁碼。來源、圖庫／色盤／劇本及64×64 PNG雜湊在originals-batch-2.json；最近鄰16倍圖只作參考。
- 內建imagegen依六張自己的原圖生成六張1254×1254首稿。孫權保留棕色鬍鬚，周瑜保留細鬍鬚，張飛保留張口神情，劉備保留黃金頭帶與淺色衣肩。司馬懿首稿鼻尖較往內，第2稿恢復右轉及右邊裁切，選第2稿。這批共八次原生輸出，連同首批四次共十二次，原生PNG從工具輸出逐byte複製且全部保存。入口generation-set.json、batch-2.json及review.json。
- 趙雲初次人工判讀錯誤，將藍色額飾與深色上唇線當成新增細節，要求第2稿移除。首個九張候選在瀏覽器七項抽驗通過，但技術通過不能證明美術忠實。重新放大1024參考並直接讀64×64 RGB後，確認額前六個像素為RGB(48,65,211)，上唇也有黑色／深綠筆畫。保留可見原色及線形，改回第1稿；鬍鬚／陰影語意維持未知，不由人物常見形象補定。第2稿、含錯誤斷言的提示詞及首個九張候選保留，勘誤及可重播像素樣本在portrait-6-source-detail-review.json；歷史提示詞不能再當成原圖沒有藍色或必須無鬍鬚的證據。候選快照為nine-first-candidate-snapshot.json，現行入口已換正確選稿。
- 本機頁擴充為九位，預設顯示新增劉備，數量／選單／標題與來源表同步。selection.json保存選稿、嘗試稿與來源，build.py由此重生，仍是本機原型資料而非正式遊戲格式。保留並排／拉桿／原版／HD、滑鼠縮放／同步拖曳／重設及依像素密度限制的HD 1:1，沒有加熱鍵。來源表可在390×844窄視窗開關九列清單。
- 修正趙雲選稿後，使用相同2GB／1 CPU／90秒外層及30秒載入門檻，乾淨重跑七項斷網瀏覽器抽驗通過，JS錯誤0、外部請求0。包含十八張內嵌PNG bytes／尺寸、九位切換、新批預設、六張新增畫面、四種模式、滑鼠拉桿／滾輪／工具列／拖曳／原生上限／重設、九列來源表、單檔搬移後重開及390×844兩倍像素密度。六張原圖及原生輸出、新批並排總覽、司馬懿與趙雲正式頁面及窄視窗已目視。原生圖未後製，總覽只是檢視用縮圖。
- verify.py核對九張原圖、圖庫／色盤／劇本、武將及頁碼映射、十二次原生bytes、趙雲原像素、三張與候選快照、選稿清單及重建一致。三張既有選稿、十張已同意城池PNG／城池頁及大地圖頁未變。最終頭像HTML 29,763,303 bytes，SHA 187339ede1d807aaf2bf6ada89fd223e84d9a51ccaa6df4e8af450fa838b3917；收據qa-result.json與verification.json。全圖本輪未改，不重跑既有十項；正式引擎未改，不以Go測試代替美術檢視。
- README、CONTEXT唯一狀態表、spec275與索引更新九張入口，首批三張技術敘述改連到歷史頁與對應收據。仍只有150頁中的九張小樣，細部雕飾／髮鬚屬重繪詮釋，未宣稱逐像素相同、全部頭像完成、逐張接受或正式接線完成。HEAD仍bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行、網站部署或上傳網站。沿用既有三顆工具映像，沒有新增映像；Docker工作有界、1000:1000、network none、來源唯讀且--rm。輸出擁有權與工作根root-owned／.md目錄核對通過，收尾檢查索引、git diff --check及專案容器清理。

## 2026-11-05：新增12張原版參考頭像，21張對照與來源視窗修正

- 使用者「go」接續已確認的HD方向。先核對CONTEXT、原圖勘誤及現行選稿，命中復古重製路由，載入remake／imagegen技能、抽驗reference與文件職責，更新首頁前讀README規則。主機gh auth及遠端open Issues核對成功，未寫入Issue、委派代理或重開畫風；正式引擎及spec275仍維持既有DRAFT邊界。
- 既有九張頁、metadata、程式、QA、完整性收據及九張選稿雜湊保存nine-r1快照。新增extract-portraits.go沿既有library／state解析器取圖，先驗名字與頁碼再寫PNG。最初預選董卓在劇本0查無，擴到四劇本仍查無；一個相依參考圖工作也因輸入不存在停止，沒有產生錯圖。依資料格式入口核對後，先匯出四劇本source-roster.json，共149個不重複名字，確認董卓、孫堅不在該集合，改選實際存在的荀彧、陸遜。沒有反覆調名字猜配圖，也沒有把查無名字當成引擎缺陷。所有相依步驟在原圖匯出成功後才續跑。
- 新增呂布35／102、荀彧62／121、袁紹2／105、袁術3／68、陸遜101／144、孫策19／127、馬超25／22、黃忠52／119、夏侯惇39／49、夏侯淵37／62、張遼81／44、龐統95／104，數字依序為武將ID／頁碼，皆可在劇本索引0找到。原圖64×64，圖庫150頁，與149個不重複名字分開計數。originals-batch-3.json保存原始資源、解碼PNG、解析腳本SHA及go1.26.7；batch-3-plan.json保存預選、替換與名單入口。沿用hr-go-ebiten:1.26.7-2.9.9-r1，專案及wl-gomod唯讀，建置快取在容器/tmp，只寫既有頭像工作目錄。
- 內建imagegen依每張原圖及最近鄰參考重繪12張1254×1254首稿，沒有加入其他作品的慣用造型。袁紹、黃忠保留灰白鬍鬚，夏侯惇保留原圖雙眼，馬超保留金盔／紅色額飾，呂布未加長雉羽。龐統首稿將圓形帽飾改成三瓣形，第2稿修回單一圓／橢圓金框，選第2稿；其餘選第1稿。13次新增原生輸出全數保存，連同前批共25次。原生PNG bytes逐一與工具輸出核對，未後製。來源、提示詞、選稿及兩張比較總覽在generation-set.json及batch-3.json，總覽只是檢視用縮圖。
- 對照頁擴充為21位，預設呂布，保留並排／拉桿／原版／HD、滑鼠縮放／同步拖曳／重設與依像素密度限制的HD 1:1。最初69,191,074 bytes版本七項抽驗通過，但目視窄視窗來源表時發現標題不見。正常滑鼠探針確認原生dialog聚焦長表格後第一個連結，scrollTop=522、標題y=-478，activeTag=A、文字為「提示詞與生成紀錄」。將關閉鈕移到上方標題列，標題／關閉鈕隨捲動保留，每次開啟回頂，並遮住標題上緣捲過的文字碎片。修正前頁、收據與截圖保存source-focus-before快照，原因與前後畫面在source-window-review.json。
- 原圖、原生輸出與兩張總覽均已目視，另由正常操作擷取六種新增頁面抽樣，檢視袁紹、夏侯惇、龐統與來源視窗。修正後使用相同2GB／1 CPU、90秒外層與30秒載入門檻乾淨重跑七項斷網抽驗通過，JS錯誤0、外部請求0。包含42張內嵌PNG bytes／尺寸、21位切換、新批預設、四種模式、滑鼠拉桿／縮放／拖曳／重設／原生上限、21列來源表、單檔搬移重開及390×844兩倍像素密度。來源視窗另驗開頭可見標題／關閉、滑鼠捲動後兩者仍可見、正常關閉及再次開啟回頂；上緣背景調整後也依同門檻重新驗證，沒有沿用另一份HTML的測試SHA。
- verify.py核對21張解碼原圖、25次原生bytes、所有嘗試稿集合、四劇本名單及頁碼、原始資源與解析腳本雜湊、趙雲既有像素勘誤、九張／三張／候選／來源視窗修正前快照及重建一致。既有九張選稿、十張已同意城池PNG／城池頁及大地圖頁未變。最終HTML 69,191,415 bytes，SHA 9c2352b3459795e0cd2e95fbaa72265e73d4f74c93bb1555e78dcf2fc6ab87f7；收據qa-result.json及verification.json。正式引擎與大地圖未改，不擴跑其既有測試。
- README、CONTEXT唯一狀態表、spec275與索引更新21張入口，九張技術敘述改連對應歷史頁及收據。21張仍是150頁中的小樣，細部線條屬重繪詮釋，未宣稱逐像素相同、逐張接受、完整頭像或正式接線。HEAD仍bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行、網站部署或原版素材公開散布。沿用三顆既有工具映像，沒有新映像或新交付目錄；Docker工作有界、1000:1000、network none、來源唯讀且--rm，輸出擁有權與工作根root-owned／.md目錄核對通過，收尾檢查索引、git diff --check及專案容器清理。

## 2026-11-05：新增12張原版參考頭像，33張離線對照與載入調整

- 使用者「繼續囉」接續既定原圖參考HD方向。先讀CONTEXT、復古重製路由、remake／imagegen技能、抽驗reference及文件職責，更新首頁前核對README規則；主機gh auth與遠端open Issues讀取成功，未寫入Issue或委派代理。現行21張頁、程式、metadata、QA與選稿雜湊先保存twentyone-r1快照，既有工作樹與正式引擎不動。
- 沿用extract-portraits.go及既有Go容器，以batch-4-plan.json從實際原始名單取賈詡36／36、荀攸63／69、徐庶65／48、法正94／17、陳宮82／26、魯肅107／31、呂蒙105／40、甘寧42／39、太史慈72／20、許褚45／24、典韋88／120、祝融125／23，數字為武將ID／頁碼，皆出自劇本索引0。originals-batch-4.json保存64×64原圖、來源及解析器雜湊；最近鄰1024參考與四張一組總覽只作檢視。沿用唯讀原資料與wl-gomod，建置快取留容器/tmp，只寫既有頭像工作目錄。
- 內建imagegen依12張自己的原圖重畫。法正保留淡灰額前方塊，祝融保留原冠形、棕髮及金紅頭帶。賈詡、荀攸、陳宮、魯肅首稿將點陣陰影畫成棋盤格，第2稿修回連續衣料／髮絲；徐庶移除新增黃色帽塊，太史慈與典韋移除新增紅額塊。七張選第2稿，五張選第1稿，19次新增原生輸出全留，累計44次。33張選稿皆1254×1254，PNG bytes與工具輸出相同且未後製。入口generation-set.json、batch-4.json及review.json；兩張新增比較總覽為檢視縮圖。
- 33張初版108,426,835 bytes在原30秒page.goto門檻逾時，未產生成功收據。保存失敗HTML、程式、錯誤及21列來源視窗截圖，再沿world/north-loading.json的既有方式分開內嵌PNG區塊、解碼成短期blob網址、移除已讀來源並釋放網址。大量資料網址與保留記憶體只列原因假說，無OOM或唯一成因證據。loading-review.json與batch4-loading-before-snapshot.json保存前後範圍及限制；不改圖片、正式資料格式、資源上限或逾時。
- 同一Playwright映像、2GB／1 CPU、30秒導覽與90秒外層乾淨重跑七項斷網抽驗通過。包含66張內嵌PNG bytes／尺寸、33位切換、賈詡預設、六種新增畫面、四種模式、滑鼠拉桿／縮放／同步拖曳／重設／原生上限、33列來源表開頭／捲動／關閉／重開、搬移單檔離線重開及390×844兩倍像素密度。三次開啟約6.8、2.5、2.8秒，各66個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。原圖、輸出、兩張總覽、賈詡正常並排、祝融細節及捲動後來源視窗已目視。
- verify.py核對33張原圖、44次原生bytes、四劇本名單、原始資源及解析器SHA、趙雲既有像素勘誤、21張及更早快照、失敗頁保存與重建一致。既有21張選稿、十張已同意城池PNG／城池頁及全圖頁未變。最終HTML 108,432,389 bytes，SHA 1cc877c6388eb800cfd734b8bcee5d37d43f912466020495c32bf811d3b78b7a，收據qa-result.json與verification.json。更新README、CONTEXT唯一狀態表、spec275及索引，21張技術敘述改連歷史頁與收據；仍為150頁中的33張小樣，未宣稱逐張接受、完整頭像或正式接線。沒有commit、push、發行、網站部署、新映像或新交付根目錄。Docker工作均有界、1000:1000、斷網、來源唯讀且--rm；擁有權、工作根root-owned／.md目錄、文件入口、git diff --check與專案容器清理列入收尾核對。

## 2026-11-05：新增12張原版參考頭像，45張對照及來源劇本核對

- 使用者「ok go go」接續既定HD方向。先核對CONTEXT、現行選稿、復古重製路由、remake／imagegen技能與文件職責，更新首頁前讀README規則。主機gh auth及遠端open Issues讀取成功，未寫入Issue或委派代理。先保存33張頁、程式、metadata、收據、抽樣畫面與33張選稿，共60個檔案雜湊，入口thirtythree-r1-snapshot.json。
- batch-5-plan.json依四劇本實際名單選徐晃、張郃、曹仁、魏延、龐德、周泰、黃蓋、程昱、張昭、姜維、鄧艾與孟獲。沿用extract-portraits.go、Go工具映像與唯讀wl-gomod，原資料唯讀，只寫既有頭像目錄。前九張來自劇本索引0，後三張來自索引1，武將ID連同來源劇本核對；原圖、映射、解析器及輸入雜湊在originals-batch-5.json。參考圖步驟一度未等Go匯出結束就啟動，因輸入尚不存在而在寫檔前停止；等同一匯出程序exit 0後，以相同容器與命令乾淨重跑成功，歸類相依步驟啟動順序問題。
- 內建imagegen依各自64×64原圖與最近鄰參考重畫12張1254×1254首稿，本批皆選第1稿，原生PNG未後製，累計56次輸出及45張選稿。黃蓋保留灰白眉鬚，曹仁保留紅色圓飾，周泰及姜維保留藍色帽飾，張昭及鄧艾保留綠色帽塊。鄧艾上唇紅點回查母本，確認座標38,41為RGB(211,0,0)，全圖只有一個該色像素，首稿保留；色點意義未知，不推作傷痕或裝飾。證據在portrait-124-source-detail-review.json，提示詞及目視範圍在generation-set.json、batch-5.json；兩張比較總覽只作檢視縮圖。
- 45張頁預設徐晃，沿用暫存圖片網址、固定工具列、並排／拉桿／原版／HD、同步縮放／拖曳與HD 1:1。來源視窗擴充至45列，載入調整的歷史紀錄改連33張收據。相同2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗，含90張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、45位切換、六種新增畫面、四種模式、滑鼠拉桿／縮放／拖曳／重設／原生上限、45列來源視窗開頭／捲動／關閉／重開、單檔搬移重開及390×844兩倍像素密度。三次載入約5.8、6.5、4.1秒，各90個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。原圖、輸出、兩張總覽、黃蓋及鄧艾正常並排、姜維窄視窗細節與捲動後來源視窗已目視。
- verify.py核對45張原圖、56次原生bytes、四劇本名單及來源劇本、原始資源與解析器SHA、趙雲既有像素勘誤、鄧艾紅點、33張及更早快照、載入歷史收據與重建一致。既有33張選稿、十張已同意城池PNG／城池頁及全圖頁未變。最終HTML 150,522,248 bytes，SHA c1be8ca53c0e2549f456a134fe158b3b6aac64f55dc7b31226e0ddb20ac5581f，收據qa-result.json及verification.json。正式引擎與全圖未改，不擴跑其測試；45張仍是150頁中的小樣，細部線條屬重繪詮釋，未宣稱逐張接受或正式接線完成。
- 更新README、CONTEXT唯一狀態表、spec275與索引，33張技術敘述改連歷史頁及對應收據。HEAD沿用bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行、網站部署、新映像或新交付根目錄。Docker工作均有界、1000:1000、斷網、來源唯讀且--rm；檔案擁有權、工作根root-owned／.md目錄、文件入口、git diff --check與專案容器清理列入收尾核對。

## 2026-11-05：新增12張原版參考頭像，57張對照與關興亮面修正

- 使用者「繼續囉」接續既定原圖參考HD方向。先核對CONTEXT、現行45張選稿、復古重製路由、remake／imagegen技能及文件職責，更新首頁前讀README規則。主機gh auth及遠端open Issues讀取成功，未寫入Issue或委派代理。先保存45張頁、程式、metadata、收據、抽樣畫面及45張選稿，共73個檔案雜湊，入口fortyfive-r1-snapshot.json。
- batch-6-plan.json選曹丕、曹植、劉禪、關平、關興、張苞、王平、郝昭、郭淮、司馬師、司馬昭與諸葛恪。沿用extract-portraits.go及既有Go映像，以唯讀原資料、唯讀wl-gomod及容器/tmp建置快取取12張64×64原圖，只寫既有頭像目錄。曹丕、曹植、關平與張苞來自劇本索引0，其餘八張來自索引1；來源劇本、名字、武將ID及頁碼一起核對，映射與雜湊保存originals-batch-6.json。確認匯出程序exit 0後才產生最近鄰參考及提示詞。
- 內建imagegen依每張自己的原圖重畫12張1254×1254首稿。曹丕保留俯視，劉禪保留藍／金冠與圓臉，張苞保留張口及深棕鬍鬚，諸葛恪保留綠色帽塊。關興首稿將頭盔淡色區畫成平面月牙，第2稿保留淡色區，局部改以金屬亮面筆觸處理；選第2稿，其餘十一張選第1稿。兩稿都保存，本批13次輸出，累計69次及57張選稿，原生PNG逐byte複製且未後製。王平額前斜紅色圖形回查原圖，區域40,8起5×5內有十個RGB(211,0,0)像素，樣本保存batch-6.json，符號意義未知；未因不熟悉該圖形而刪掉原色。提示詞與選稿入口generation-set.json，兩張新增總覽僅作檢視縮圖。
- 本機頁擴充57位、預設曹丕，沿用既有載入方式、固定工具列、原版／HD／並排／拉桿、同步縮放／拖曳／重設及HD 1:1。相同2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗，含114張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、57位切換、六種新增畫面、四種模式、滑鼠拉桿／縮放／拖曳／重設／原生上限、57列來源視窗開頭／捲動／關閉／重開、單檔搬移後重開及390×844兩倍像素密度。本次三次載入皆約3.8秒，各114個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。原圖、原生輸出、兩張總覽、劉禪及關興第2稿正常並排、張苞窄視窗細節與捲動後來源視窗已目視。
- verify.py核對57張原圖、69次原生bytes、四劇本名單與來源劇本、原始資源及解析器SHA、趙雲／鄧艾既有像素證據、王平紅色樣本、45張及更早快照、載入歷史收據與重建一致。既有45張選稿、十張已同意城池PNG／城池頁及全圖頁未變。最終HTML 189,359,703 bytes，SHA db0636009b6b8a6b5548a17f0067d1fdd7a7e6d025d29f907c095951ee23ac5e，收據qa-result.json及verification.json。正式引擎與全圖未改，不擴跑其測試；57張仍是150頁中的小樣，細部筆觸屬重繪詮釋，未宣稱逐張接受或正式接線完成。
- 更新README、CONTEXT唯一狀態表、spec275及索引，45張技術敘述改連歷史頁與對應收據。HEAD沿用bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行、網站部署、新映像或新交付根目錄。Docker工作均有界、1000:1000、斷網、來源唯讀且--rm；檔案擁有權、工作根root-owned／.md目錄、文件入口、git diff --check及專案容器清理列入收尾核對。

## 2026-11-05：新增12張原版參考頭像，69張對照與帽飾細修

- 使用者「ok go」接續既定原圖參考HD方向。先核對CONTEXT、現行57張選稿、復古重製路由、remake／imagegen技能及文件職責，更新首頁前讀README規則。主機gh auth及遠端open Issues讀取成功，未寫入Issue或委派代理。先保存57張頁、程式、metadata、收據、抽樣畫面及57張選稿，共85個檔案雜湊，入口fiftyseven-r1-snapshot.json。
- batch-7-plan.json選公孫瓚、劉表、劉璋、馬騰、韓遂、孔融、田豐、審配、馬良、馬謖、諸葛瑾與陳登。沿用extract-portraits.go及既有Go映像，以唯讀原資料、唯讀wl-gomod及容器/tmp建置快取取12張64×64原圖，只寫既有頭像目錄。十二張均來自劇本索引0，名字、武將ID及頁碼一起核對，映射與雜湊保存originals-batch-7.json。確認匯出程序exit 0後才產生最近鄰參考及提示詞。
- 內建imagegen依每張自己的原圖重畫12張1254×1254首稿。審配首稿將灰帽明暗畫成獨立白色方塊，第2稿恢復連續灰綠布面；諸葛瑾首稿淡金帽飾過大，第2稿縮小輪廓，保留原圖淡色／金色細節。兩者選第2稿，其餘十張選第1稿。全部版本保存，本批14次輸出，累計83次及69張選稿，原生PNG逐byte複製且未後製。提示詞與選稿入口generation-set.json，兩張新增總覽僅作檢視縮圖；細部筆觸與帽飾材質仍屬重繪詮釋。
- 本機頁擴充69位、預設公孫瓚，沿用既有載入方式、固定工具列、原版／HD／並排／拉桿、同步縮放／拖曳／重設及HD 1:1。相同2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗，含138張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、69位切換、六種新增畫面、四種模式、滑鼠拉桿／縮放／拖曳／重設／原生上限、69列來源視窗開頭／捲動／關閉／重開、單檔搬移後重開及390×844兩倍像素密度。三次載入約4.5、4.5、4.8秒，各138個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。原圖、原生輸出、兩張總覽、審配與諸葛瑾第2稿正常並排及捲動後窄視窗來源表已目視。
- verify.py核對69張原圖、83次原生bytes、四劇本名單與來源劇本、原始資源及解析器SHA、趙雲／鄧艾既有像素證據、王平紅色樣本、57張及更早快照、載入歷史收據與重建一致。既有57張選稿、十張已同意城池PNG／城池頁及全圖頁未變。最終HTML 226,002,096 bytes，SHA 466efc943c6444f9d968070dda7b97fdeb4f9af61611dd73dc0b97881f40368f，收據qa-result.json及verification.json。正式引擎與全圖未改，不擴跑其測試；69張仍是150頁中的小樣，未宣稱逐張接受或正式接線完成。
- 更新README、CONTEXT唯一狀態表、spec275及索引，57張技術敘述改連歷史頁與對應收據。索引582份文件及135條斷言檢查通過，87個頭像連結存在，檔案擁有權皆1000:1000，工作根root-owned檔案與.md目錄皆0，git diff --check通過。HEAD沿用bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行、網站部署、新映像或新交付根目錄。Docker工作均有界、1000:1000、斷網、來源唯讀且--rm；最終docker ps -a --filter name=wolong無執行中或停止容器。

## 2026-11-06：新增12張原版參考頭像，81張對照與原色色塊細修

- 使用者「ok go」接續既定原圖參考HD方向。核對CONTEXT、69張選稿、復古重製路由、remake／imagegen技能、抽驗交接及README／文件職責。主機gh auth及遠端open Issues讀取成功，未寫入Issue或委派代理。先保存69張頁、程式、metadata、收據、抽樣畫面及69張選稿，共97個檔案雜湊，入口sixtynine-r1-snapshot.json。
- batch-8-plan.json選顏良、文醜、張繡、張魯、馬岱、嚴顏、程普、徐盛、丁奉、曹洪、曹真與滿寵。沿用extract-portraits.go及既有Go映像，以唯讀原資料、唯讀wl-gomod及容器/tmp快取取12張64×64原圖，只寫既有頭像目錄。十二張均來自劇本索引0，名字、武將ID及頁碼一起核對，來源與SHA保存originals-batch-8.json。確認匯出程序exit 0後才產生最近鄰參考及提示詞。
- 內建imagegen依每張自己的原圖重畫12張1254×1254首稿。嚴顏紅色點、程普綠色區、曹洪藍色塊與滿寵下唇暖色亮點皆回查原圖。嚴顏／程普／曹洪選第2稿，收斂過大的比例與平面感；滿寵第2稿仍有橙色橫條，第3稿融回嘴唇及皮膚明暗，其餘八張選第1稿。全部版本保存，本批17次輸出、累計100次及81張選稿，原生PNG逐byte複製且未後製。細部筆觸、材質及符號意義仍屬詮釋／未知。提示詞與選稿入口generation-set.json，兩張總覽僅作檢視縮圖。
- 本機頁擴充81位、預設顏良，沿用既有載入方式、固定工具列、四種比較模式與同步縮放／拖曳／重設及HD 1:1。相同2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗，含162張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、81位切換、六種新增畫面、滑鼠操作、81列來源視窗開頭／捲動／關閉／重開、單檔搬移後重開及390×844兩倍像素密度。三次載入約3.3、5.3、4.1秒，各162個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。原圖、原生輸出、兩張總覽、四張修正稿並排、嚴顏窄視窗細節及捲動後來源表已目視。
- verify.py核對81張原圖、100次原生bytes、四劇本名單及來源劇本、原始資源／解析器SHA、趙雲／鄧艾／王平既有像素證據、69張及更早快照、載入歷史收據與重建一致。既有69張選稿、十張已同意城池PNG／城池頁及全圖頁未變。HTML 265,009,434 bytes，SHA 00bda6411edeab10a7ff60aa633d5cffb780acb0bd4fba52fd6da9b95bd1b162，收據qa-result.json及verification.json。正式引擎與全圖未改，不擴跑其測試；81張仍是150頁中的小樣，未宣稱逐張接受或正式接線完成。
- 更新README、CONTEXT唯一狀態表、spec275及索引，69張技術敘述改連歷史頁與收據。HEAD沿用bd09fd40e344efee9d01e286a70f2a2303890caf，保留既有工作樹，沒有commit、push、發行、網站部署、新映像或新交付根目錄。索引582份文件及135條斷言檢查通過，97個頭像連結存在，檔案擁有權皆1000:1000，工作根root-owned檔案與.md目錄皆0，git diff --check通過。Docker工作均有界、1000:1000、斷網、來源唯讀且--rm；最終docker ps -a --filter name=wolong無執行中或停止容器。

## 2026-11-06：本週進度提交與推送

- 使用者明確要求「本週進度先這樣，先push」。本輪只整理提交、驗證與交接，不續做HD素材或新增玩法。分支main、遠端origin為wicanr2/wolong_cht；fetch後起始HEAD與origin/main同為bd09fd40e344efee9d01e286a70f2a2303890caf。GitHub主機auth及Issues讀取成功，未寫Issue。
- 提交範圍為本週累積的全滑鼠工具列、原版上下推動反白、清單拖曳、行軍／讀檔恢復、入戰中斷與續跑、戰場上下文、測試及spec266至275等文件。HD現況保留五城十張、三十五區全圖試繪與81張頭像小樣的分級，不宣稱整款remake或HD已完成。
- 沿既有san1-go:1.26.7-ebiten2.9.9-r1及tools/go.sh，正式package patterns的go test -count=1與go vet通過。tools/check.sh亦實跑，Go段及文件索引通過；整體仍停於歷史幽靈引用，不宣稱完整檢查全綠。其餘文件／工具逐項抽驗的正對照、worklist、資產deny-list及TALK校訂工具均通過。
- 本週文件增加後，分流索引缺34列、舊3列被新證據取代。保留1105列既有分類，逐項核對34列；另補四份記錄的限制標題後加入5列，分流與嚴格生成均驗證1144列。未重開已完成分支、改原版語意或寫遠端Issue。首次嚴格生成失敗曾清空輸出，先恢復未修改的暫存區版本，再由正式工具重生。另明列rich2外部來源完整路徑，將資料庫／程式SHA與文件連結分段，兩份舊腳本SHA標明當時版本，所有原始雜湊值保留。
- README現況同步261份規格、255 CONFORMED／4 READY／2 DRAFT；未解主表1144列與6筆平台層分開計。過期斷言掃描已通過。幽靈引用204筆主要是既有本機歷史路徑；lessons render --check與silent-fixture-drop防線檢查在原HEAD也失敗，保留既有問題，不以改checker或放寬規則通過。
- 原版EXE／資料、HD原圖及重繪PNG、離線HTML、workplace與dist-all均不入提交。deny-list掃描1868個追蹤檔，逐內容比對120個原版雜湊，未命中原版資產。提交前會再次核對暫存區、檔案型態、git diff --check、擁有權及容器清理；不建Release、移動tag、部署網站或force push。
- 本機收據為[提交檢查與遠端核對](workplace/weekly-push-20261106-checks.json)，保存檢查輸出、原HEAD比較與收尾的git_delivery。獨立檢查彙整器最初因既有Python的tarfile不支援filter參數中止，改用同容器的tar讀取可信git archive後完成；這是驗證腳本環境差異，不是產品失敗。依本輪授權建立「保存本週進度：原版滑鼠操作、行軍存讀檔與HD規劃」提交並正常推送origin/main，實際新SHA及遠端一致結果記於收據及回覆。

## 2026-11-11：新增12張原版參考頭像，93張對照與眼部／頭盔細修

- 使用者「繼續囉」接續既定原圖參考HD方向。核對CONTEXT、81張選稿、復古重製路由、remake／imagegen技能、抽驗交接及README／文件職責。主機gh auth及遠端open Issues讀取成功，未寫入Issue或委派代理。起始main與origin/main同為14c7700818071e3710d921e6223daaf1858cb348，工作樹乾淨。先保存81張頁、程式、metadata、收據、抽樣畫面及81張選稿，共109個檔案雜湊，入口eightyone-r1-snapshot.json。
- batch-9-plan.json選凌統、文聘、廖化、周倉、蔣欽、陳武、朱桓、曹休、曹純、夏侯霸、孟達與張任。沿用extract-portraits.go與既有Go映像，以唯讀原資料、唯讀wl-gomod及容器/tmp快取取12張64×64原圖，只寫既有頭像目錄。朱桓／曹休／夏侯霸來自劇本索引1，其餘九張來自索引0，名字、武將ID與頁碼連同劇本核對，保存originals-batch-9.json。確認匯出exit 0後才產生最近鄰參考及提示詞。
- 內建imagegen依各自原圖重畫12張1254×1254首稿。文聘第2稿修正重複眼睛與面頰孤立色點，陳武第2稿移除新增大符號／額前菱形／棋盤衣紋，凌統第2稿收斂過度張口並移除下巴細綠線，曹純第2稿將鳥形細節改成小型抽象帽前亮面，張任第2稿將金色碎片融成連續淡金亮面。五張選第2稿，其餘七張選第1稿；原有藍色、綠布及淡金色區保留。曹純原圖下唇35,40有RGB(211,0,0)單一紅點，選稿保留，樣本記於batch-9.json，意義未知。全部17次輸出逐byte複製且未後製，累計117次、93張選稿，提示詞入口generation-set.json。兩張新增總覽僅作檢視縮圖。
- 本機頁擴充93位、預設凌統，沿用既有載入方式、固定工具列、四種比較模式與同步縮放／拖曳／重設及HD 1:1。相同2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗，含186張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、93位切換、六種新增畫面、滑鼠操作、93列來源視窗開頭／捲動／關閉／重開、單檔搬移後重開及390×844兩倍像素密度。三次載入約9.7、7.2、5.4秒，各186個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。原圖、17次原生輸出、兩張總覽、文聘並排、曹純窄視窗細節與捲動後來源表已目視。
- verify.py核對93張原圖、117次原生bytes、四劇本名單及來源劇本、原始資源／解析器SHA、趙雲／鄧艾／王平既有像素證據、曹純紅點、81張及更早快照、載入歷史收據與重建一致。既有81張選稿、十張已同意城池PNG／城池頁及全圖頁未變。HTML 303,867,988 bytes，SHA b3caa5a2ae8f7c10c54844890c5f99f52d197877fd53999a458636832c0ac075，收據qa-result.json及verification.json。正式引擎與全圖未改，不擴跑其測試；93張仍是150頁中的小樣，未宣稱逐張接受或正式接線完成。
- 更新README、CONTEXT唯一狀態表、spec275及索引，81張技術敘述改連歷史頁與對應收據。本輪不commit、push、發行、部署網站、建新映像或交付根目錄。Docker工作均有界、1000:1000、斷網、來源唯讀且--rm；文件索引582份／135條斷言與嚴格未解分流1144列均通過，只刷新HD現況一列，保留1143列既有分類。107個頭像文件連結與109個81張快照雜湊核對通過，擁有權皆1000:1000，工作根root-owned檔案與.md目錄皆0，git diff --check通過；最終docker ps -a --filter name=wolong無執行中或停止容器。

## 2026-11-11：新增12張原版參考頭像，105張對照與驗證工具分段處理

- 使用者「ok 繼續囉」接續原圖參考HD方向。核對CONTEXT、93張現況、復古重製路由、remake／imagegen技能與文件職責，沿既有映像及工作目錄。主機gh auth與遠端open Issues讀取成功，未寫Issue或委派代理。HEAD維持14c7700818071e3710d921e6223daaf1858cb348，保留前輪八份未提交文件變更。先保存93張頁、程式、metadata、收據、抽樣畫面與93張選稿，共121個檔案雜湊，入口ninetythree-r1-snapshot.json。
- batch-10-plan.json選曹叡、曹彰、關索、王雙、紀靈、凌操、孫桓、孫韶、劉封、蔡瑁、張紘與孫乾。沿原extract-portraits.go與Go映像，以唯讀原資料、唯讀wl-gomod及容器/tmp快取匯出12張64×64原圖；確認程序exit 0後才產生最近鄰參考。前四張來自劇本索引1，其餘八張索引0，名字／劇本／武將ID／頁碼與來源SHA保存originals-batch-10.json。曹彰為劇本1武將4、頁123，劇本0武將14的曹漳也用同頁，保留兩版原字，頁碼只計一次。
- 內建imagegen依每張自己的原圖重畫12張1254×1254首稿。關索／王雙／凌操／孫桓／孫韶／劉封／蔡瑁／孫乾選第2稿，收斂帽前牌飾或過大的紅／金／綠細節，移除新增灰色圓點／白色鉚釘及棋盤衣紋，其餘四張選第1稿。孫桓金色紋樣的細部形狀與符號意義仍屬詮釋。蔡瑁原圖帽前14個RGB(81,146,65)綠色像素已逐點保存，第2稿縮小捲紋而保留綠色。全部20次輸出逐byte複製且未後製，累計137次、105張選稿；原生圖、未選稿、提示詞與兩張檢視總覽留本機，入口batch-10.json及generation-set.json。
- 本機頁105位、預設曹叡，沿既有四種模式、固定工具列與同步縮放／拖曳／重設及HD 1:1。初次完整QA開頁出現Page crashed，先保存HTML／app／QA／錯誤於batch10-loading-before-snapshot.json。相同HTML單獨導覽約5.2秒通過；驗證程式改讓qa-assets.cjs的PNG讀檔子程序結束後才啟動瀏覽器，子程序RSS約1.09GB，瀏覽器啟動前父程序約79MB。相同2GB／1 CPU、30秒導覽及90秒外層完成原七項完整抽驗，沒有放寬斷言。原失敗未取得OOM計數，原因不作已證實OOM聲明。
- 七項抽驗含210張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、105位切換、六種新增畫面、四種模式與滑鼠控制、105列來源視窗開頭／捲動／關閉／重開、單檔搬移及390×844兩倍像素密度。三次載入約5.1、6.5、15.2秒，各210個暫存網址全釋放，來源區塊全移除，JS錯誤0、外部請求0。原圖、20次原生輸出、兩張總覽、王雙正常並排、孫乾窄視窗細節與捲動後來源表已目視。
- 完整性核對在內嵌重建時以137終止，保留當時build.py／verify.py後，將建頁的全部base64列表／合併字串／完整HTML同時保留，改為對齊區塊分段寫入暫存檔再原子替換。相同1GB／1 CPU、60秒的verify.py完整重跑exit 0，核對105張原圖、137次原生bytes、四劇本名單、來源與解析器SHA、既有色點及蔡瑁綠色樣本、123頁原字、93張及更早快照與重建一致。HTML保持342,877,215 bytes與SHA 0eb9fa5eddea1e9e871df3fe853fbc4f8a839bc96056d8b628f23db60959f690，與已通過瀏覽器抽驗及初版逐byte相同。既有93張選稿、十張已同意城池PNG／城池頁及全圖頁未變。工具前後、診斷及收據在batch10-browser-review.json；正式引擎與原生素材未改，不擴跑其測試。
- 更新README、CONTEXT唯一狀態表、spec275及索引，93張技術敘述改連歷史頁與收據。105張仍為150頁中的小樣，未宣稱逐張接受或正式接線完成。本輪不commit、push、發行、部署網站、建立新映像或交付根目錄。Docker均有界、1000:1000、斷網、來源唯讀且--rm；文件索引582份／135條斷言與嚴格未解分流1144列通過，只刷新HD現況一列，保留1143列既有分類。121個頭像文件連結、121個93張快照雜湊與六個失敗頁／工具檔案核對通過，擁有權皆1000:1000，工作根root-owned檔案與.md目錄皆0，git diff --check通過。原生圖與HTML維持gitignore；最終docker ps -a --filter name=wolong無執行中或停止容器。

## 2026-11-11：新增12張原版參考頭像，117張對照與逐張串流核對

- 使用者「繼續囉」接續既定原圖參考HD方向。先核對CONTEXT、105張現況、復古重製路由、remake／imagegen技能、文件職責與驗證交接。主機gh auth及遠端open Issues讀取成功，未寫Issue或委派代理。沿用HEAD 14c7700818071e3710d921e6223daaf1858cb348與前輪八份未提交文件變更。先保存105張頁、程式、metadata、收據、抽樣畫面、105張選稿及前批失敗證據，共145個檔案雜湊，入口onehundredfive-r1-snapshot.json。
- batch-11-plan.json選袁尚、袁譚、袁胤、王朗、李傕、郭汜、韓玄、金旋、嚴白虎、公孫越、黃祖與劉繇。沿既有extract-portraits.go、Go映像、唯讀原資料與wl-gomod及容器/tmp快取匯出12張64×64原圖，確認exit 0後才產生最近鄰參考。十二張皆來自劇本索引0，名字／劇本／武將ID／頁碼及來源SHA保存originals-batch-11.json；袁胤的ID確為0。
- 內建imagegen依各自原圖重畫12張1254×1254首稿。袁譚／袁胤／李傕／郭汜／金旋／黃祖選第2稿：收斂側臉比例與完整白框、將口部孤立亮條融回髭線／唇部明暗、縮小藍色塊、移除凸起灰釘與規則側盔橫紋，以及移除肩部過大的黑色菱格。其餘六張選第1稿。黃祖禿額、原有髭線及色塊保留，細部材質與幾何仍屬詮釋。袁胤帽前18個RGB(130,162,97)綠色像素與郭汜13個深淺藍樣本保存batch-11.json，意義未知。全部18次原生PNG逐byte複製且未後製，累計155次、117張選稿，原圖／提示詞／未選稿及兩張檢視總覽留本機。
- 117張單檔沿原有分段建頁、四種模式、固定工具列、同步縮放／拖曳／重設與HD 1:1，預設袁尚。第一次完整QA碰90秒外層逾時，停在搬移後重開；保留HTML、QA與錯誤於batch11-browser-before-snapshot.json。將PNG檢查改為逐張串流讀取，各桌面／獨立／窄視窗依序關閉再重開；第二次仍於窄視窗載入時逾時。第二次階段顯示18.5秒完成靜態核對、39.9秒初次ready、57.0秒完成全部117位切換、68.7秒完成獨立重開；該版工具／錯誤／階段紀錄保存batch11-browser-attempt2.json。
- 兩次同類失敗後回查驗證路由，並查app.js確認選取會更新標題並同步重畫。將每次切換固定等100毫秒，改成等選項／標題實際更新，原標題斷言及全部117位切換保留。相同2GB／1 CPU、30秒導覽、90秒外層完成原七項完整斷網抽驗；沒有減少檢查項目。階段約51.9秒到達第七項，三次載入約6.8、11.1、14.1秒，各234個暫存網址全釋放，來源區塊全移除，JS錯誤0、外部請求0。未隔離每個修正的效果，也不宣稱原失敗的OOM機制已證實。
- 七項包含234張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、117位切換、新批預設與六種畫面、四模式與滑鼠控制、117列來源視窗開頭／捲動／關閉／重開、單檔搬移重開後袁胤選取及390×844兩倍像素密度。原圖、18次原生輸出與兩張總覽已目視；正常頁面另看袁譚／李傕並排、黃祖窄視窗細節及捲動後來源表。
- verify.py依原有1GB／1 CPU、60秒完整重跑exit 0，核對117張原圖、155次原生bytes、四劇本名單、原始資源／解析器SHA、既有像素證據、袁胤綠色與郭汜藍色、ID0、105張及更早快照、前兩次失敗頁／工具與重建一致。105張階段的QA與初版改以歷史快照核對，沒有拿新117張收據代稱105張。HTML保持380,383,085 bytes，SHA 3ad3a843bf5d368454b1c03e8218f001288d9edc7217284af7af44a57f01390b，與初次失敗頁逐byte相同。既有105張選稿、十張已同意城池PNG／城池頁及全圖頁未變。工具檢視入口batch11-browser-review.json，通用失敗輸出在雜湊核對封存後清除；歷史檔保留。
- 更新README、CONTEXT唯一狀態表、spec275及索引，105張技術敘述改連歷史頁與收據。117張仍為150頁中的小樣，未宣稱逐張接受或正式接線完成；正式引擎與原生素材未因驗證工具調整而改動，不擴跑其測試。本輪不commit、push、發行、部署網站、建立新映像或交付根目錄。Docker均有界、1000:1000、斷網、來源唯讀且--rm；文件索引582份／135條斷言及嚴格未解分流1144列通過，只刷新HD現況一列，保留1143列既有分類。135個頭像文件連結、145個105張快照雜湊與八個逾時封存檔案核對通過，擁有權皆1000:1000，工作根root-owned檔案與.md目錄皆0，git diff --check通過。原生圖與HTML仍由gitignore排除；最終docker ps -a --filter name=wolong無執行中或停止容器。

## 2026-11-11：原版參考HD頭像累計120張

- 新增劉勳、孫瑜、劉琦，來源劇本索引0，武將ID／頁碼為8／96、21／60、28／78。內建imagegen四次原生輸出保存在本機，累計159次；劉勳第2稿修平新增寶石，保留原有綠色區塊。未選稿保留，PNG未後製。
- [第十二批紀錄](workplace/hd-redraw-local-20261104/faithful/portraits/batch-12.json)保存完整提示詞、來源與選稿。[117張快照](workplace/hd-redraw-local-20261104/faithful/portraits/117-r1-snapshot.json)保留398項檔案雜湊，既有117張、城池及全圖頁未變。
- 七項離線抽驗通過，120張切換、240張內嵌PNG、三張新增畫面、縮放／拖曳、單檔搬移與窄視窗來源操作均通過；JS錯誤及外部請求皆0。原2GB／1 CPU與時間限制沿用，app.js和build.py未變。
- 完整性核對159次原生圖、120張原圖及來源映射、參考圖、快照與HTML重建一致。仍為DRAFT小樣，未接入引擎或取得逐張接受；未提交或推送原版素材、圖片與HTML。
- 文件索引582份／135條斷言及嚴格未解分流1144列通過，保留1143列既有分類。規格109個頭像連結存在，工作根root-owned檔案與.md目錄皆0，頭像擁有權皆1000:1000。git diff --check通過，HTML與原生圖仍由gitignore排除，專案容器已清理。

## 2026-11-11：原版參考HD頭像累計126張

- 新增劉度、夏侯尚、牛金、吳蘭、張勳、張松，來源劇本索引0，名字、武將ID與頁碼逐張核對。內建imagegen原生輸出累計169次，PNG未後製，未選稿保留。選稿與細修依[第十三批紀錄](workplace/hd-redraw-local-20261104/faithful/portraits/batch-13.json)，完整提示詞與原生圖皆保存在本機。
- [120張快照](workplace/hd-redraw-local-20261104/faithful/portraits/120-r1-snapshot.json)保存262項檔案雜湊。既有120張選稿、城池與全圖頁未變。
- 七項離線抽驗通過，核對252張內嵌PNG及126張切換、六張新增畫面、滑鼠縮放／拖曳、來源視窗、單檔搬移與窄視窗操作；JS錯誤及外部請求皆0。資源、時間限制及app.js／build.py沿用前批。
- 完整性核對全部169次原生圖、126張原圖、來源映射、參考圖與重建一致。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未部署網站、提交或推送素材。
- 文件索引582份／135條斷言與嚴格未解分流1144列通過，保留1143列既有分類。規格123個頭像連結存在，工作根root-owned檔案與.md目錄皆0，頭像擁有權皆1000:1000。git diff --check通過，HTML與原生圖仍由gitignore排除，專案容器已清理。

## 2026-11-11：原版參考HD頭像累計132張

- 新增雷銅、李恢、冷苞、伍習、李暹、張衛，來源劇本索引0，名字、武將ID與頁碼逐張核對。內建imagegen原生輸出累計177次，PNG未後製，未選稿保留。選稿與細修依[第十四批紀錄](workplace/hd-redraw-local-20261104/faithful/portraits/batch-14.json)，完整提示詞與原生圖皆保存在本機。
- [126張快照](workplace/hd-redraw-local-20261104/faithful/portraits/126-r1-snapshot.json)保存277項檔案雜湊。既有126張選稿、城池與全圖頁未變。
- 七項離線抽驗通過，核對252張內嵌PNG及132張切換、六張新增畫面、滑鼠縮放／拖曳、來源視窗、單檔搬移與窄視窗操作；JS錯誤及外部請求皆0。資源、時間限制及app.js／build.py沿用前批。
- 完整性核對全部177次原生圖、132張原圖、來源映射、參考圖與重建一致。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未部署網站、提交或推送素材。
- 文件索引582份／135條斷言與嚴格未解分流1144列通過，保留1143列既有分類。規格137個頭像連結存在，工作根root-owned檔案與.md目錄皆0，頭像擁有權皆1000:1000。git diff --check通過，HTML與原生圖仍由gitignore排除，專案容器已清理。

## 2026-11-11：原版參考HD頭像累計138張

- 新增嚴圃、鞏志、笮融、張英、嚴輿、周昕，來源劇本索引0，名字、武將ID與頁碼逐張核對。內建imagegen原生輸出累計189次，PNG未後製，未選稿保留。選稿與細修依[第十五批紀錄](workplace/hd-redraw-local-20261104/faithful/portraits/batch-15.json)，完整提示詞與原生圖皆保存在本機。
- [132張快照](workplace/hd-redraw-local-20261104/faithful/portraits/132-r1-snapshot.json)保存289項檔案雜湊。既有132張選稿、城池與全圖頁未變。
- 七項離線抽驗通過，核對276張內嵌PNG及138張切換、六張新增畫面、滑鼠縮放／拖曳、來源視窗、單檔搬移與窄視窗操作；JS錯誤及外部請求皆0。資源、時間限制及app.js／build.py沿用前批。
- 完整性核對全部189次原生圖、138張原圖、來源映射、參考圖與重建一致。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未部署網站、提交或推送素材。
- 並排檢視後修正嚴輿張口判讀：原圖本就大張口，第2稿收得太小。第3稿保留張口，只收斂嘴部亮條；第一次候選的24項檔案雜湊與通過的技術收據另存，最終選稿重新通過七項抽驗與完整性核對。
- 文件索引582份／135條斷言與嚴格未解分流1144列通過，保留1143列既有分類。規格155個頭像連結存在，工作根root-owned檔案與.md目錄皆0，頭像擁有權皆1000:1000。git diff --check通過，HTML與原生圖仍由gitignore排除，專案容器已清理。
- 圖庫頁133保留劇本0的111嚴圃與劇本3的117閻圃兩個原始名稱，只計一張選稿。另依封存的[132張QA收據](workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-132-r1.json)修正spec275第十四批文內的內嵌圖數：原寫252，實際為264；圖片與收據未改。此次文件數字直接由收據取值。

## 2026-11-11：原版參考HD頭像累計146張

- 新增李別、刑道榮、劉賢、武安國、劉磐、楊昂、趙廣、趙統，前六張來源劇本索引0，趙廣與趙統為索引1，名字、武將ID與頁碼逐張核對。內建imagegen原生輸出累計203次，PNG未後製，未選稿保留。選稿與細修依[第十六批紀錄](workplace/hd-redraw-local-20261104/faithful/portraits/batch-16.json)，完整提示詞與原生圖皆保存在本機。
- [138張快照](workplace/hd-redraw-local-20261104/faithful/portraits/138-r1-snapshot.json)保存301項檔案雜湊。既有138張選稿、城池與全圖頁未變。
- 七項離線抽驗通過，核對292張內嵌PNG及146張切換、八張新增畫面、滑鼠縮放／拖曳、來源視窗、單檔搬移與窄視窗操作；JS錯誤及外部請求皆0。資源、時間限制及app.js／build.py沿用前批。
- 完整性核對全部203次原生圖、146張原圖、來源映射、參考圖與重建一致。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未部署網站、提交或推送素材。
- 四劇本149個原始名稱共用146頁，這146頁均有小樣；圖庫145、147、148、149四頁未出現在現有名單，用途待確認。另依[132張封存選稿表](workplace/hd-redraw-local-20261104/faithful/portraits/selection-132-r1.json)與[138張封存選稿表](workplace/hd-redraw-local-20261104/faithful/portraits/selection-138-r1.json)回填前兩批文內來源包數14／15，原收據與圖片未改；本批數字直接取收據。
- 前兩次完整抽驗碰到90秒外層逾時。靜態核對工具改用有界Buffer組行，並重新建立短ID，避免集合保留正規表示式擷取的大字串。相同512MB容器與256MB JavaScript heap探針通過，完成時RSS為84,631,552 bytes；完整七項抽驗在原90秒限制通過。HTML、app.js、主QA與原生PNG未改，失敗收據及舊工具保留於[本批工具檢視](workplace/hd-redraw-local-20261104/faithful/portraits/batch16-browser-review.json)，不將全部時間差歸因於單一因素。
- 收尾核對通過：索引582份文件／135條斷言，triage共1144列，保留其餘1143列分類；spec275的236個本機連結可解析。既有138張快照與城池／全圖頁未變，檔案擁有權符合1000:1000，工作根root-owned檔案及誤建.md目錄皆0。git diff --check通過，專案Docker容器已清空。

## 2026-11-11：擴充工具列本機操作比較稿

- 依既定HD方向，補[單排／雙排比較稿](workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1.html)，原版／HD及縮放維持常駐，加入AI與武將編輯，共11個功能入口、13個實體按鈕。依grilling先提供可撤回操作稿，建議分組雙排，尚未記為使用者選擇。版面選擇不包含編輯器正式套用時機、保存邊界或AI正式預設。
- 使用既有中原原圖與HD、曹操原圖與HD選稿，沒有新增或改寫美術。編輯取既有實際資料，支援姓名文字輸入、滑鼠增減、委任類型、取消、套用示範及還原；僅保留頁面記憶體。AI開關僅改預覽狀態，戰鬥及存檔均未接線，正式遊戲未改。
- 首次十項QA為綠，目視卻發現雙排仍在同列，原因為session的flex-basis覆蓋width。修正為100%基準，補群組座標及工具列高度驗證，重跑十項離線抽驗通過。修正前六份檔案雜湊保存在[版面修正紀錄](workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1-row-fix.json)，舊綠收據不作分列證據。
- [目前抽驗](workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1-qa.json)涵蓋兩版面、1280／960及640寬兩倍像素密度、真實素材切換、滾輪／拖曳、編輯取消與還原、AI預覽、入口及快捷鍵停用、單檔搬移。七張截圖已保存，目視雙排、單排、編輯器與窄視窗；JS錯誤0、外部請求0。
- [完整性收據](workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1-verification.json)確認8份輸入、原始曹操資料、修正前快照、HTML重建一致，以及城池／大地圖／146張頭像頁雜湊未變。工具入口及來源已連入spec275；沒有新增Markdown檔案、提交、推送或部署網站。
- 收尾：索引582份文件／135條斷言通過，spec275的246個本機連結可解析；triage仍1144列，strict未解文件與現況一致。兩個工作根root-owned及誤建.md目錄皆0，工具列工作檔UID/GID符合1000:1000。git diff --check通過，專案相關Docker容器已清空，比較稿與截圖仍受gitignore保護。

## 2026-11-11：工具列採單排

- 使用者明確指定「工具列要單排」。採用單排，原版／HD、縮放、AI及武將編輯維持同列；雙排僅保留歷史比較。[選擇紀錄](workplace/scratch/ui-windows-20261104/hd-toolbar-single-choice.json)保存原話與範圍，[選擇前快照](workplace/scratch/ui-windows-20261104/hd-toolbar-before-single-choice.json)保存十一份檔案與選定畫面，不覆寫原參考截圖。此選擇不擴張為編輯器與AI正式行為均已確認。
- 本機確認稿已改為預設單排，重開仍為單排，窄視窗不換列。姓名輸入、編輯草稿、HD切換、滾輪／拖曳及AI狀態預覽沿用既有範圍，尚未接入正式遊戲。
- 前兩次抽驗出現黑圖，[首次失敗](workplace/scratch/ui-windows-20261104/hd-toolbar-single-first-attempt.json)與[第二次失敗](workplace/scratch/ui-windows-20261104/hd-toolbar-single-second-attempt.json)均保留。第二次像素探針記錄地圖已載入、畫布context未失效，但鏡頭座標為NaN；獨立[診斷探針](workplace/scratch/ui-windows-20261104/hd-toolbar-single-probe.json)未重現，沒有將其當作修正通過。
- 新增畫布零尺寸防護。[受控前後對照](workplace/scratch/ui-windows-20261104/hd-toolbar-single-size-control.json)以同一原圖隱藏畫布並呼叫paint，再還原，證明舊稿座標失效、新稿保留有限值。前兩次自然失敗沒有零尺寸軌跡，未將所有失敗原因宣稱為已證實。十項完整離線抽驗重跑通過，另核對每張截圖前地圖非空、預設單排、640寬不換列及重開；JS錯誤與外部請求皆0。
- 收尾核對通過：HTML重建一致，城池／大地圖／146張頭像頁未變，選定截圖及十一份快照雜湊一致。索引582份文件／135條斷言、spec275的250個本機連結及1144列triage通過。失敗收據封存後移除過期入口檔，工作目錄UID/GID符合1000:1000，root-owned及誤建.md目錄皆0。git diff --check通過，專案Docker容器已清空，未提交、推送或部署。

## 2026-11-11：HD頭像與顯示偏好首段接線

- 使用者選「預設原版，記住之後的選擇」。在spec275先寫READY小節，完成素材、合成、偏好與限定正常路徑驗證後，該小節標為限定範圍CONFORMED；整體仍DRAFT。正式工具列由七項加頭像顯示切換，共八鈕單排，縮放、AI與編輯器其餘入口繼續待接線。
- 新增唯讀HD選稿載入：沿selection.json及逐張metadata，不另建散布包格式；驗證實際150頁圖庫與146張選稿，包含頁碼0、重複／越界、檔名逃逸、SHA及PNG尺寸。每次解碼再核對原生bytes，GPU頭像快取最多八張，缺頁回原圖。
- 八處Portrait呼叫收斂到drawPortrait。保序合成在每張原生頭像前貼入先前原版畫布，後畫的視窗、文字與游標可正常蓋住頭像；最後由Ebiten的FinalScreenDrawer繪出，避免先縮成64×64。原模式沿DefaultDrawFinalScreen，對拍與錄影仍走原版640×400；顯示偏好hd_display不寫入遊戲存檔。
- 受影響冷測與vet通過。首次GPU測試在RunGame之前讀像素而失敗，改為獨立程序、實際Draw內核對後通過，舊log存hd-head-tests-before-runtime-probe.log。這是測試生命週期修正，非產品缺陷。
- [正常滑鼠抽驗](workplace/scratch/ui-windows-20261104/hd-head-gui.json)共29項通過，從首頁走新遊戲君主卡、原版／HD切換、離開確認遮擋、勢力情報窗、正常存檔與重開；另驗未指定素材時保留偏好且回原圖。首輪錯把首頁當作已捕捉游標；次輪錯把君主卡的上一幀計數當成世界HUD證據。已保留兩份失敗與trace，補Draw階段識別，改用原版橫幅滑鼠開關真正打開情報窗後乾淨重跑，未注入狀態或使用direct-entry。
- [實際遊戲對照頁](workplace/scratch/ui-windows-20261104/hd-head-game-review.html)內嵌七張實際截圖及兩張原尺寸裁圖，原版頭像裁圖10色、HD裁圖7394色，只作原生細節技術核對，不作美術接受或原版parity聲明。[接線收據](workplace/scratch/ui-windows-20261104/hd-head-integration.json)保存來源、測試、binary SHA；現行binary為wlgame-hd-head-r1，SHA-256 788f189159fd996d4b8caa1c22cec4bc8b5cec4d4d114fd63476375a4f63c104。
- 正常GUI只抽君主卡與勢力情報窗。命名、編成、軍團及戰術對話的共用入口不代證完整正常流程；地圖HD、世界縮放、AI加成、編輯器及原生平台驗收仍未完成。素材與畫面留本機，未提交、推送、發行或部署。

- 最終補查液晶設定：無HD液晶變體時沿用原圖並停用切換，保留HD偏好，切回彩色恢復。新29項GUI包含此往返，受影響冷測、獨立GPU核對與vet再次通過。舊偏好缺少hd_display時明確回false，語系套用失敗前不提前改HD偏好。
- 最終首次GUI重跑曾因素材驗證尚未完成、內層state只等2秒而提前失敗；改由既有12秒wait統一等待啟動trace，沒有放寬外層限制或改遊戲。失敗與啟動log存hd-head-gui-before-startup-sync-fix.json及hd-head-game-before-startup-sync-fix.log。實際君主卡、確認視窗遮擋與勢力情報畫面已目視，七張遊戲截圖、29項GUI及四項離線對照頁抽驗通過。
- 收尾核對通過：索引582份文件／135條斷言、spec275的265個本機連結、1144列未解分流及strict文件一致。接線收據核對18份來源、最終binary、七張截圖、測試與頁面雜湊；城池／世界／146張頭像畫廊未變。工作檔UID/GID符合1000:1000，root-owned及誤建.md目錄皆0；git diff --check通過，專案Docker容器已清空。引擎程式、測試與翻譯詞表留工作樹，素材／binary／存檔／HTML仍受gitignore保護，未提交或推送。

## 2026-11-11：戰略地圖50%～200%縮放

- 使用者選「50%～200%（建議）：兼顧視野與辨識」。沿已展示的預設100%、每次25%、工具列及滾輪、介面固定契約，先在spec275列READY，再接正式桌面。頭像首段十八檔來源另存[快照索引](workplace/scratch/ui-windows-20261104/world-zoom-baseline.json)，舊binary與收據保留。
- 正式工具列加入三個縮放鈕，共十一個實體鈕單排。地圖、城池、軍團、雲層／災害與格游標共用鏡頭；橫幅、視窗與HD頭像固定。100%沿既有整格分支，回100%吸附最近整格，非100%累積小數捲動，小地圖框隨視野調整。倍率不進遊戲存檔，新遊戲與讀檔回100%。
- [原資料探針](workplace/scratch/ui-windows-20261104/world-zoom-camera-probe.json)保存五種倍率與78次可見城池反算，僅為可丟棄原型。初版以Go與Canvas逐像素比對失敗；[取樣診斷](workplace/scratch/ui-windows-20261104/world-zoom-pixel-phase.json)確認縮小時整數邊界選擇不同，各自40／40符合來源，跨實作29／40相同，沒有把原型測試改稱精確渲染通過。
- 正式受影響冷測通過，192城五倍率共960次座標檢查。首次測試誤要求100%無操作清掉殘值，已改回原版契約。GPU首次比較漏掉劇本雲層，75%一處黑像素來自原有cloud slot24；以原全圖獨立合成MCH雲層後，五倍率35處取樣及橫幅邊界通過。兩項均為測試基準問題，未據此改遊戲規則。vet及本機建置通過。
- [31項正常滑鼠抽驗](workplace/scratch/ui-windows-20261104/world-zoom-gui.json)由首頁走新遊戲，通過50%／100%／200%點選同一座城82、城池與設定視窗停用縮放、游標焦點、小地圖定位、存檔未被縮放重寫及讀檔回100%，HD頭像與偏好保留。第一輪容器未設工作目錄而找不到translations/corrections.json，補-w /repo後相同容器乾淨重跑通過，分類為驗收環境問題。
- [實際遊戲截圖頁](workplace/scratch/ui-windows-20261104/world-zoom-game-review.html)內嵌七張正常操作畫面，目視50%、200%與小地圖。四項斷網瀏覽器抽驗通過，包含單檔搬移、畫面切換及窄視窗，外部請求與JS錯誤皆0。[接線收據](workplace/scratch/ui-windows-20261104/world-zoom-integration.json)固定26份來源、binary、截圖與測試SHA；來源覆蓋包配合base_git_head可重建。
- spec275第二段僅就原版地圖呈現及上述抽驗標CONFORMED。地圖HD、AI、編輯器及原生平台驗收仍待完成，既有城池／世界／146張頭像頁SHA未變。未提交、推送、發行或部署。
- [收尾核對](workplace/scratch/ui-windows-20261104/world-zoom-final-checks.json)七項通過。索引仍582份文件／135條斷言，1144列分流與strict生成結果一致，git diff --check通過。來源UID/GID1000:1000，工作區root-owned與誤建.md目錄皆0，本輪Go、GUI、頁面與資料容器均已刪除，其他專案容器未動。GUI檔案雜湊抽查範圍為SAVE.DAT，原生執行狀態另經正常讀檔驗證，不擴大存檔聲明。

## 2026-11-11：HD大地圖接入正式桌面

- 使用者接續「ok go」，沿上一段已展示的HD大地圖接線。載入復古重製技能、規格閘門與既有spec275，先列第三段READY契約再實作。[接線前索引](workplace/scratch/ui-windows-20261104/hd-world-baseline.json)固定縮放版26檔來源覆蓋包與收據，前段binary未覆寫。主機gh驗證成功後唯讀核對open Issues，未寫入遠端。
- `-hd-world`沿本機既有regions.json與逐項選稿JSON載入35區、22接縫；驗證矩形、ID、原生標記、PNG SHA與尺寸，按需解碼並限制GPU快取16張／128MiB。與頭像共用HD檔案驗證入口，沒有新增散布格式。原版／HD按鈕現在同時控制頭像與地景，首次原版及hd_display偏好沿用，單排十一鈕不變。
- 地景依世界座標直接畫到原生合成面，接縫沿既有8世界像素透明度。即時標記沿原版markerTables共用分派；角旗／歸屬採差異框加一像素邊界，首都與軍團用MCH遮罩，災害／格游標後畫，固定UI最後覆蓋。原版100%分支與遊戲規則不改，季節／液晶無對應HD素材時沿原圖。
- [資料／標記冷測](workplace/scratch/ui-windows-20261104/hd-world-core-tests.log)、[介面冷測](workplace/scratch/ui-windows-20261104/hd-world-ui-tests.log)、[獨立GPU抽驗](workplace/scratch/ui-windows-20261104/hd-world-gpu.log)及vet通過。57張實際選稿驗證通過；12種歸屬／首都／軍團狀態將透明標記疊回原地形後，逐byte等於既有RenderMarked。GPU保留原生條紋細節，後畫UI與單區缺圖回退正確，另驗非0季節色盤與液晶不套錯地景。
- [54項正常滑鼠抽驗](workplace/scratch/ui-windows-20261104/hd-world-gui.json)由首頁進新遊戲，涵蓋三倍率點許昌、原版／HD、液晶往返、模態遮擋、小地圖、全圖西北／東南邊界、正常存讀檔及缺圖重開。SAVE.DAT與SAVE-slot1.wlsave兩檔SHA在顯示／縮放前後未變。trace同步實際Draw的鏡頭與倍率；沒有狀態注入。GUI未抽移動中軍團，證明範圍限於標記合成測試與已走正常路徑。
- 首次GUI工具的字串替換誤改helper分割字串而產生SyntaxError，遊戲尚未啟動；修正並先compile語法檢查，再用同一容器／命令乾淨重跑通過，分類為驗收腳本問題。沒有據此改產品輸入或遊戲規則。
- [本機實際遊戲頁](workplace/scratch/ui-windows-20261104/hd-world-game-review.html)內嵌12張截圖，目視100%、50%、200%與東南邊界。四項斷網頁面抽驗通過，包含搬移單檔、12張切換與窄視窗，JS錯誤及外部請求皆0。[接線收據](workplace/scratch/ui-windows-20261104/hd-world-integration.json)固定32份來源、57選稿、binary與截圖SHA；來源覆蓋包配合base_git_head可重建。抽樣GPU快取最高75,491,140 bytes。
- spec275第三段僅就色盤0選稿與上述範圍標CONFORMED。原生圖像有限，超過素材1:1的放大仍為插值。其他季節／液晶HD、據點景觀／戰術HD、局部美術、AI與編輯器仍待完成；不宣稱逐區美術接受或原版圖像parity。圖像與頁面留本機，未提交、推送、發行或部署。
- [收尾核對](workplace/scratch/ui-windows-20261104/hd-world-final-checks.json)七項通過：32份來源與覆蓋包、57張原生PNG、現行binary／收據／頁面、文件入口與擁有權一致。索引仍582份文件／135條斷言，1144列分流與strict生成結果一致，git diff --check通過；工作區root-owned與誤建.md目錄皆0，來源UID/GID1000:1000，本輪Docker容器已清空。

## 2026-11-11：HD龍紋橫幅、原版縮放抽驗與精簡自動工具列

- 使用者指出最上方尚未HD，確認為龍紋／標題／日期橫幅；另要求原版可縮小地圖、工具列小一點並能隱藏。沿grilling做[28／6高度操作稿](workplace/scratch/ui-windows-20261104/toolbar-compact-hide-prototype.html)，使用者選靠上緣顯示、離開收起。新裁定取代工具列始終展開，維持單排與畫布外置。
- [原版模式13項抽驗](workplace/scratch/ui-windows-20261104/original-map-zoom-gui.json)證明既有50%～200%、滾輪及小地圖可用，無須改縮放規則。初輪誤以左鍵關小地圖，核對既有左開／右關後只修測試，失敗trace與腳本另存before-close-fix。工具列操作稿首輪hover驗收的目標在喚出後消失，改用滑鼠座標觸發並檢查狀態後四項斷網抽驗通過；皆為驗收預期／同步問題。
- 使用[擷取工具](workplace/scratch/ui-windows-20261104/banner-source-probe.go)取得八組ICONGRF原版橫幅與來源SHA。內建imagegen依色盤0原圖重繪，原生2172×724 PNG保存在[原生選稿](workplace/hd-redraw-local-20261104/faithful/banner-hd-r1.png)，[完整提示詞](workplace/scratch/ui-windows-20261104/banner-prompt-r1.txt)與[生成紀錄](workplace/scratch/ui-windows-20261104/banner-generation-r1.json)可回查。PNG bytes未改，顯示內容矩形為(0,309,2172,106)，排除白色外緣；未把本張美術記為已個別接受。
- spec275第四段先READY再接線。`-hd-banner`載入既有HD共用驗證契約下的選稿，直接畫原生背景，日期沿原字模即時繪入，原版四個視窗開關及保留熱區未移動。液晶／未提供的季節沿原圖，模式仍用hd_display，沒有新增遊戲存檔欄位。
- 正式工具列40縮到28高，初始收起6高，靠上緣喚出、離開收起。按住、放開幀、原版選單或拖曳擁有輸入時不切換。一般視窗同步調整客戶區高度，保持遊戲兩倍顯示；捕捉游標重映射並重建raw基準，防止收展引發假捲動。共用滑鼠工具改讀動態高度並保留舊binary的40高fallback。
- [素材冷測](workplace/scratch/ui-windows-20261104/hd-banner-asset-tests.log)、[介面冷測](workplace/scratch/ui-windows-20261104/hd-banner-ui-tests.log)、[獨立GPU](workplace/scratch/ui-windows-20261104/hd-banner-gpu.log)與vet通過。涵蓋內容矩形、壞圖、原生條紋、即時日期後畫、手勢完整性及18次固定視窗邊界座標轉換。全螢幕／最大化只核對座標數學，未作原生GUI驗收。
- [正常滑鼠流程](workplace/scratch/ui-windows-20261104/banner-toolbar-gui.json)從首頁走新遊戲，通過96個功能及同步檢查點，含重複收展同步，不稱96個獨立功能。實際1280×812／856，原版縮放／小地圖、HD橫幅／頭像、拖出放開、模式／液晶往返、存讀檔與缺圖重開通過，兩種存檔SHA未被顯示操作重寫。目視HD橫幅與展開後較小工具列。
- [本機實際遊戲頁](workplace/scratch/ui-windows-20261104/banner-toolbar-game-review.html)內嵌12張截圖並附原生PNG與提示詞，四項斷網頁面抽驗通過。[接線收據](workplace/scratch/ui-windows-20261104/banner-toolbar-integration.json)固定38份來源、binary、工具、素材與截圖SHA；前段32檔覆蓋包與受測程式保留。其他HD場景／色盤、美術修整、AI與編輯器仍待完成，未提交、推送、發行或部署。
- [收尾核對](workplace/scratch/ui-windows-20261104/banner-toolbar-final-checks.json)七項通過；原生橫幅、既有57張地景及三個美術頁SHA未變，來源與38檔覆蓋包一致。索引仍582份文件／135條斷言，1144列分流與strict生成結果一致；來源UID/GID1000:1000，工作區root-owned與誤建.md目錄皆0，git diff --check通過，本輪Docker容器已清空。

## 2026-11-12：五組已同意景觀接入情報窗

- 接續使用者「ok 繼續囉」，命中復古遊戲路由，載入規格閘門入口、復古重製技能與共同決策規則。沿已同意的五張景觀先列spec275第五段READY，再實作；沒有新增畫風或產品決定。主機gh驗證成功後唯讀核對open Issues，未寫入遠端。
- `-hd-locations`讀既有manifest.json與originals.json，依實際KYOGRF15頁及City.KindHigh對應許昌6、滎陽14、舞陽2、虎牢關5、汝南0。驗證選稿路徑、尺寸、SHA、原生標記及原映射，兩頁GPU快取按需載入。景觀仍在原情報窗(16,288)、96×96，與64×64頭像共用原生矩形合成，文字與後畫視窗保留順序。
- [素材冷測](workplace/scratch/ui-windows-20261104/hd-location-asset-tests.log)、[受影響介面冷測](workplace/scratch/ui-windows-20261104/hd-location-ui-tests.log)、[獨立GPU](workplace/scratch/ui-windows-20261104/hd-location-gpu.log)、既有頭像GPU回歸及vet通過。涵蓋實際五張1254×1254選稿、頁0／越界／重複／原映射不符／逃逸／壞圖、原生細節、96與64尺寸、後畫遮擋、快取淘汰及執行中缺圖回退。
- [正常滑鼠抽驗](workplace/scratch/ui-windows-20261104/hd-location-gui.json)由首頁走新遊戲、小地圖定位、逐一點五城；145個功能與同步檢查點包含重複收展同步，不當作獨立功能數。氾水關共用頁5、宛未選頁10回原圖、原版／HD及液晶往返、快取重載、正常讀檔及缺素材重開通過。SAVE.DAT與SAVE-slot1.wlsave的SHA未被顯示操作改寫。
- [本機實際畫面頁](workplace/scratch/ui-windows-20261104/hd-location-game-review.html)內嵌13張截圖，附截圖裁區方便比較，未增加圖片細節。[四項離線抽驗](workplace/scratch/ui-windows-20261104/hd-location-review-qa.json)通過，搬移單檔與窄視窗正常，JS錯誤及外部請求皆0。[接線收據](workplace/scratch/ui-windows-20261104/hd-location-integration.json)固定43份來源、覆蓋包、binary、測試與截圖SHA，前段來源與受測程式保留。
- 五張原生PNG、既有素材索引及美術頁未修改。CONFORMED僅限五頁與上述抽驗，其餘十頁、其他色盤／場景、AI及編輯器仍待完成，不宣稱原版parity或全部HD完成。未提交、推送、發行或部署。
- [收尾核對](workplace/scratch/ui-windows-20261104/hd-location-final-checks.json)七項通過，43檔來源覆蓋包與本機受測程式一致。索引582份文件／135條斷言，1144列分流與strict生成結果一致；來源UID/GID1000:1000，工作區root-owned與誤建.md目錄皆0，git diff --check通過。本輪Docker容器已清空，其他專案容器未動。

## 2026-11-12：移除額外「右鍵關閉」提示

- 接續後使用者指出「右鍵關閉」訊息視窗多餘，優先處理本項。命中復古遊戲及GUI還原路由，載入既有規格閘門與GUI還原入口，未開始新美術生成。主機gh驗證成功後唯讀核對open Issues，未寫遠端。
- 據點及軍團情報窗的額外框在程式內明標remake差異；進言結果亦使用同一字串。先在spec266列READY，移除三處額外提示、無用座標與四語字串，保留原版TALK #4、右鍵／返回分派及對白推進。沒有改遊戲規則或存檔欄位。
- [受影響冷測、vet與本機建置](workplace/scratch/ui-windows-20261104/no-close-hint-tests.log)通過。Go對唯讀模組快取寫入stat cache的提示未造成失敗，未修改原版素材或借用主機工具鏈。
- GUI工具初輪漏掉新軍團與許昌同格時的據點／軍團二選一；次輪把君主編成軍團，依既有規則阻止進言；第三輪誤用單像素位移選取反白，原契約需越過100–131的範圍。分別查證mapChoice、openAdvise與choiceBandMove，修正測試流程與相對位移，不改產品程式。失敗工具／trace／log／收據以before-city-army-fix、before-lord-corps-fix及before-band-sync-fix保存於同一工作目錄。
- 以相同容器／命令乾淨重跑後，[正常滑鼠抽驗](workplace/scratch/ui-windows-20261104/no-close-hint-gui.json)69個功能及同步檢查點通過，包含重複工具列同步，不宣稱69個獨立功能。正常首頁、新遊戲、原版／HD據點窗、軍團編成與情報、同格選單、進言結果的右鍵關閉及正常存讀檔通過；兩種存檔未被顯示操作重寫。說服分支結束提示只經程式審查與既有測試，不擴大正常GUI聲明。
- [本機實際畫面](workplace/scratch/ui-windows-20261104/no-close-hint-game-review.html)附四張修正後截圖與前版比較；[接線收據](workplace/scratch/ui-windows-20261104/no-close-hint-integration.json)保存45檔來源、覆蓋包、binary與截圖SHA，原五組景觀、素材索引及三個美術頁未修改。現行binary為wlgame-no-close-hint-r1，SHA-256 ab9cc9422bad48e6bcf2f0d36350d36d896708fff34b555c8d1074681114ac1b。未提交、推送、發行或部署。
- [三項離線頁面抽驗](workplace/scratch/ui-windows-20261104/no-close-hint-review-qa.json)通過，五張實際截圖完整、窄視窗無橫向溢出，JS錯誤與外部請求皆0。接線收據的final_checks核對45檔來源覆蓋包、受測程式、文件入口及擁有權通過。索引582份文件／135條斷言、1144列分流與strict生成結果一致；來源UID/GID1000:1000，工作區root-owned與誤建.md目錄皆0，git diff --check通過。本輪容器皆已結束刪除，其他專案容器未動。

## 2026-11-12：編成與存讀檔成功不再重複通知

- 使用者指出「夏侯淵編成完畢」「已儲存第一槽」是原版沒有的訊息，要求考慮必要性與玩家體驗。命中GUI還原與規格閘門路由；核對spec88的額外六秒事件框、spec109的原版編成對白及存讀檔完成關窗流程。先列spec266 READY，再移除編成／存檔／讀檔成功通知，並在成功後清除舊通知，避免前次錯誤殘留。
- 編成仍顯示原版主將肖像對白並建立部隊；存讀檔成功直接關窗返回。失敗處理、原版事件與其他訊息不改，沒有為此新增提示音、工具列動畫、存檔欄位或遊戲規則。
- [既有編成、存檔、桌面及事件列冷測、vet與建置](workplace/scratch/ui-windows-20261104/quiet-success-tests.log)通過。[正常滑鼠抽驗](workplace/scratch/ui-windows-20261104/quiet-success-gui.json)從首頁走新遊戲、編成、原版主將對白、存檔、讀檔及HD切換，共50個功能與同步檢查點，包含重複工具列同步，不當作獨立功能數。兩種存檔存在、可正常讀回，顯示操作未重寫存檔。失敗分支由既有測試與程式審查保留，不宣稱全部儲存裝置錯誤GUI完成。
- [本機比較頁](workplace/scratch/ui-windows-20261104/quiet-success-game-review.html)包含五張本輪截圖與兩張前版畫面；已目視編成原版對白及存檔後地圖。[接線收據](workplace/scratch/ui-windows-20261104/quiet-success-integration.json)固定45檔來源覆蓋包、binary、流程工具、測試與截圖SHA；前版程式與收據保留。未提交、推送、發行或部署。
- [離線頁面抽驗](workplace/scratch/ui-windows-20261104/quiet-success-review-qa.json)通過七圖載入、窄視窗及離線檢查，JS錯誤與外部請求皆0。收尾核對曾抓到spec109引用檔名錯誤，修正為既有109-formation-leader-line.md後相同核對通過。來源覆蓋包、程式、截圖與既有美術雜湊一致；索引582份文件／135條斷言、1144列分流及strict輸出一致。UID/GID1000:1000、工作區root-owned與誤建.md目錄皆0，git diff --check通過，本輪容器已清空，其他專案容器未動。

## 2026-11-12：通知清理、訊息右鍵與曹操對呂布比較

- 接續移除行軍、戰術下令及已有TALK的外交／撥款摘要，保留失敗與重要事件。使用者指定第一劇本曹操對呂布，抽驗操作、動畫及停手時間推進；沿復古GUI與規格閘門，先READY再實作，不改規則或動畫速度。
- 預設dosgolem-wolong目錄為空，改用既有dosgolem來源f79eab6db8c64a524a06535e53e5d32b0bbcf7c8唯讀建工具，沒有另建image。初始腳本誤點表頭及選單前提，改查熱區逐步操作；boot失敗紀錄保留，分類為腳本問題。原版素材唯讀，所有工作負載在有界Docker內執行。
- 原版首次sub_1ECE0取數前設定NewFixed(1)的258-byte狀態，Go用-seed 1，事前固定、不重擲挑結果。正常新遊戲、第一劇本曹操、宣戰呂布13、編成夏侯淵37、向徐州79下令，沒有直接叫戰鬥或修改日期、軍團、勝負。不宣稱自然RTC播種或完整骰序已對齊。
- [首次失敗收據](workplace/scratch/ui-windows-20261104/cao-lubu-remake-failure-before-message-right-fix.json)重現4/1 9時、(210,114)訊息不接受右鍵。原版同一點是陳留請求援軍，右鍵可繼續。補READY後接回右鍵逐頁、末頁才續行，保留idleGate.Pause。[冷測、vet及建置](workplace/scratch/ui-windows-20261104/cao-lubu-tests-r2.log)通過，[負對照](workplace/scratch/ui-windows-20261104/message-right-negative.log)確實拒絕僅左鍵版本。
- 修正後首跑在存檔前遇到新訊息，工具未先處理輸入擁有者；補正常右鍵同步後同程式、同種子通過，存檔閘門未改，before-save-owner-sync-fix紀錄保留。[正常GUI](workplace/scratch/ui-windows-20261104/cao-lubu-remake-gui.json)60個功能／同步檢查點、14圖通過，包含宣戰、編成、下令、停手、途中存讀檔及戰鬥退場，不當作60個獨立功能。
- [可見行軍補充](workplace/scratch/ui-windows-20261104/cao-lubu-visible-gui.json)由正常小地圖移到我軍。初版誤要求第一格移動時小時欄立即變更，改為等到座標及時欄皆前進，不改子刻或速度。兩側同於4/6 6時、夏侯淵(241,108)攻呂布35(243,108)、椎陽76入戰。原版134步及33圖SHA完全重現，加入唯讀1D318寫入累計後仍相同。
- 戰術同拍尚未通過：原版932拍到sub_19F58結算入口，Go最新最後活動幀440，首次通過的另一取樣為583，兩份皆保留。下令時點、完整入戰RNG與單兵狀態未鎖定，不調速度掩蓋或單由數字判定規則錯誤。原版有梯子，Go本案未取得同事件畫面；投射物、受擊與倒地仍沿Issue #13／#16／#24。原版僅取到寫回入口，七種int61音效呼叫尚未實作，完整戰後返回與音訊不作本輪完成聲明。
- [比較矩陣](workplace/scratch/ui-windows-20261104/cao-lubu-comparison.json)、[本機頁](workplace/scratch/ui-windows-20261104/cao-lubu-game-review.html)及[接線收據](workplace/scratch/ui-windows-20261104/cao-lubu-integration.json)固定50檔來源、兩側程式、種子、工具及畫面身分。使用者要求完成後commit與push；原版素材、HD圖、存檔及HTML留本機，不作Release或整款完成聲明。
- [收尾核對](workplace/scratch/ui-windows-20261104/cao-lubu-final-checks.json)通過，12張原版、14張remake及2張視野內行軍圖的離線切換正常，JS錯誤與外部請求皆0。50檔來源及兩側程式／覆蓋包雜湊一致，既有素材未改；582份文件／135條斷言與1144列分流一致。來源UID/GID1000:1000，工作區root-owned與誤建.md目錄皆0，git diff --check通過。本輪容器皆已結束刪除，其他專案容器未動。依授權準備提交並推送，實際提交及遠端核對結果保存於[本機推送收據](workplace/scratch/ui-windows-20261104/repo-update-20261112.json)。

### 2026-11-12：訂正戰術開場移動倒數

- 使用者續行曹操對呂布的操作與動畫比較。先讀CONTEXT、現行Issues #13／#16／#24及復古戰術對拍路由，沿既有工作區；使用reverse-engineer-retro-game-remake與use-ida-pro-9-4技能，不新增美術、AI規則或速度調整。
- dosgolem來源f79eab6db8c64a524a06535e53e5d32b0bbcf7c8乾淨，原版固定seed1的正常開局重播不改部隊、日期或戰鬥結果。新增唯讀19C45、19FA0、1ADC8觀測，保存完整96槽及258-byte RNG；兩次共同230快照含指令數均相同。正式Go來源未接觀測器，早期觀測版只使用可丟棄Go build overlay。
- 原自然入戰的兵種、血量、戰力及X相同，但94槽Y不同，RNG置換表相同而狀態byte不同。另見原版96槽倒數0、Go均1。以原版擺位前的RNG隔離Spawn後，其他8項初始化欄位全等，只剩96個倒數不同。[修正前](workplace/scratch/ui-windows-20261104/cao-spawn-before.json)保留，不搜尋種子或略過亂數湊結果。
- IDA Pro9.4一次性DB的[指令／bytes／xref匯出](workplace/scratch/ui-windows-20261104/ida-probe.json)證實19ADF先清零兵士，19C95只設存在旗標。19CC0–19CC7將ES移到兵士段+0xC0個段落，19D20的0x180屬+0xC00城牆物件。舊spec200的設定端推論因此推翻；保留歷史並先READY，再將Spawn.Stun改0。資料庫SHA為03d00e84f396d3fc828a6d0f2407a0e4515d735b334e7a73c37c1ab5ef938321，KI.EXE SHA沿spec200。
- [修正後初始化](workplace/scratch/ui-windows-20261104/cao-spawn-after.json)96槽9項全等，[首次移動](workplace/scratch/ui-windows-20261104/cao-first-movement-after.json)另驗X／Y、血量、餘力、倒數及步行相位全部相同。這是局部同輸入規則驗證，不宣稱原版主迴圈入口之後的全部RNG、腳本或結算已同步。
- [四套件冷測、vet及建置](workplace/scratch/ui-windows-20261104/cao-delay-tests.log)通過，含戰術、世界、桌面及手機。新增TestSpawnHasNoArtificialMovementDelay驗96槽倒數與首批餘力，[舊版設1的負對照](workplace/scratch/ui-windows-20261104/cao-delay-negative.log)確實被拒絕。正式binary wlgame-cao-delay-r1 SHA為98d871656c47cacc66e869304a33e5d6a16810d897bcb476fc84afb3c14bef15。
- 工具環境問題分開保留：login shell移除Go PATH、探針誤用欄位名、GUI未指定工作目錄、IDA版本API名稱不符，均修正隔離腳本後重跑，未當作產品缺陷。正常GUI再遇兩個取樣錯誤：等待喊話消失才下令，以及首格移動便要求小時改變。依spec155及15-realtime原有契約修正，before-talk-gate-fix與before-idle-hour-gate-fix失敗紀錄保留，不改戰術速度、子刻或遊戲程式來遷就腳本。
- 正式程式[正常滑鼠收據](workplace/scratch/ui-windows-20261104/cao-delay-remake-gui.json)通過61個功能／同步檢查點與14圖，從首頁曹操、宣戰呂布、編夏侯淵、停手行軍、途中存讀檔，自然入戰直到戰後返回。spec200僅將本次初始化／首批移動範圍標CONFORMED；機制文件與CONTEXT追加訂正。相鄰物件的完整分類改為不阻塞證據，不再掛到兩兵漂移Issue #13。
- [接線收據](workplace/scratch/ui-windows-20261104/cao-delay-integration.json)保存來源覆蓋包、正式binary、原版工具、快照、測試與SHA；[本機對照頁](workplace/scratch/ui-windows-20261104/cao-delay-game-review.html)五圖離線載入，[390寬頁面抽驗](workplace/scratch/ui-windows-20261104/cao-delay-review-qa.json)通過，JS錯誤及外部請求0。截圖屬不同自然入戰狀態的流程檢視，沒有原版／remake同拍畫面聲明。
- 完整入戰RNG與命令相位、投射物、架梯及後段結算仍沿Issues #13／#16／#24；未關閉遠端Issue，未作Release。原版素材、HTML、快照、存檔及二進位繼續留本機。
- [收尾核對](workplace/scratch/ui-windows-20261104/cao-delay-final-checks.json)通過；582份文件／135條斷言、1144列分流與嚴格生成輸出一致。來源覆蓋包、正式binary、固定seed、觀測工具、快照及瀏覽器收據雜湊相符，新證據連結存在，工作區root-owned與.md目錄皆0。容器已結束刪除，未動其他專案資源。沿使用者先前提交與推送指示更新repo，實際提交與遠端核對寫入[本機推送收據](workplace/scratch/ui-windows-20261104/repo-delay-update-20261112.json)。

### 2026-11-12：鎖敵Y方向與隊長就位修正

- 使用者續行戰術操作比較。從CONTEXT、GitHub Issues #13／#16與戰術對拍路由接手，沿reverse-engineer-retro-game-remake及IDA技能；目前base為6fb853e。原版輸入沿用上輪正常曹操對呂布、seed1的完整記錄，檔案SHA未變，沒有新跑一輪原版的聲明。
- 局部Go探針以同一擺位前RNG初始化96槽，再以同一主迴圈入口RNG重播，對齊命令後、移動前的原版1ADC8邊界，不搜尋拍點偏移。原錯版第一拍有6個Target不同；第5拍起另見就位隊長的餘力下降。正式遊戲沒有接入探針或移植戰鬥快照。
- 新IDA輸出放在既有工作區內的movement-20261112子目錄，因包裝器固定產生ida-probe.json與input.exe.i64，必須避免覆寫上一輪已雜湊的證據。`tools/ida.sh probe dosv tools/ida_tactical_targeting.py`沿用既有非root image，輸入KI.EXE SHA為fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868，IDA Pro9.4、739函式、UID/GID1000。資料庫SHA為64dde76e1f5a52645760f5596a3ad51292ec5c0cb4b5771b226f28ddfb09f66c，未改名或修改正式資料庫。
- [原指令與bytes](workplace/scratch/ui-windows-20261104/movement-20261112/ida-probe.json)顯示1A889有借位時會先經1A88D加Y差，再於1A8A9加一次。原版候選Y較大時距離為`|dx|+2|dy|`，另一方向為`|dx|+|dy|`；保留平面懲罰與等距先掃。舊re/11摘要省略第一道加法，已訂正，同時保留CONTEXT勘誤。守陣接敵與攻擊碰撞的其他距離判斷未合併或更改。
- 隊長表1A7E7索引7指向1AA2C，隊員表1A82D索引7指向1A82C的retn。隊長就位後仍重算陣形目標，抵達時補餘力；Go原先將兩者都當空操作。依spec159補READY並先跑[可丟棄候選](workplace/scratch/ui-windows-20261104/cao-movement-candidate-comparison.json)，確認前3拍96槽吻合後才修改正式規則。
- [正式規則40拍比較](workplace/scratch/ui-windows-20261104/cao-movement-after-comparison.json)核對27個欄位，前三拍每拍96槽全等，40個完整RNG狀態相同。第4拍仍有槽29的StepX原版2／Go5；其後更多走位差異保留，未猜補中繼點、改尋路預算或調速度。
- [冷測、vet及建置](workplace/scratch/ui-windows-20261104/cao-move-tests.log)通過戰術、世界、桌面及手機四套件。新增Y方向／等距選擇、六隊隊長跨雙側的就位餘力與陣形線移動、隊員命令7空分派測試；[舊錯版負對照](workplace/scratch/ui-windows-20261104/cao-move-negative.log)分別被鎖敵與就位測試拒絕。桌面既有opt-in trace新增隊長位置、目標及餘力，沒有玩家通知或存檔格式變更。
- 正式binary wlgame-cao-move-r1的SHA為c3f8a7ff3e797bd05bb0c5978c3d10a4d574a3fa71716d83f7b95e647b26b68c。[正常滑鼠](workplace/scratch/ui-windows-20261104/cao-move-remake-gui.json)通過66個功能／同步檢查點、15圖：首頁曹操、宣戰、編夏侯淵、下令徐州、停手行軍與途中存讀檔，自然遇呂布後改中央線。就位隊長目標由(3,32)成(26,32)，指標固定(300,180)時X由3至4、5、6；正常改回、突擊、240拍後點退卻並返回戰略。本片驗主動退卻分支，不聲稱完整突擊勝負與原版相同。
- [本機對照頁](workplace/scratch/ui-windows-20261104/cao-move-game-review.html)四張正式畫面與[斷網抽驗](workplace/scratch/ui-windows-20261104/cao-move-review-qa.json)通過，390寬無溢出，JS錯誤及外部請求皆0。[接線收據](workplace/scratch/ui-windows-20261104/cao-move-integration.json)固定5檔來源覆蓋包、binary、原版快照、探針、測試與頁面雜湊。
- spec159只有本片兩項標CONFORMED，整體排程仍READY，Issues #13／#16／#24保留。下一個最小查證點是第4拍槽29的尋路輸入及中繼點輸出，完整投射物、架梯及結算未代證。原版、HTML、截圖、快照及binary留本機，不作Release。
- [收尾核對](workplace/scratch/ui-windows-20261104/cao-move-final-checks.json)通過：5檔來源、覆蓋包、binary、原版快照及頁面雜湊一致，1144列分流與嚴格輸出未變，582份文件／135條斷言索引通過，新證據連結存在。工作區root-owned與.md目錄皆0，來源UID/GID1000:1000。本輪有界容器皆結束刪除，其他專案資源未動。沿先前指示提交與推送，實際Git提交及遠端核對保存於[本機推送收據](workplace/scratch/ui-windows-20261104/repo-move-update-20261112.json)。

### 2026-11-12：戰術尋路的反向擴散與成本回溯

- 沿CONTEXT、現行Issues #13／#16及復古戰術對拍路由接手，base為69e7f05。使用既有逆向與IDA技能，聚焦第4拍槽29的中繼點，不改速度、AI加成、存檔或美術。
- dosgolem乾淨來源f79eab6db8c64a524a06535e53e5d32b0bbcf7c8以正常曹操開局重播，首次RND前固定seed1；新增唯讀1BD46、1AF44、1AF4C觀測，17個共同快照含指令數、96筆原記錄與完整RNG皆未變。12筆請求有11筆成功、1筆失敗，不重擲或選樣。
- IDA Pro9.4經既有非root image及tools/ida.sh probe匯出[原始bytes](workplace/scratch/ui-windows-20261104/path-20261112/ida-probe.json)，KI.EXE與DB身分見spec134。擴充既有tools/ida_range.py的/output模式，不另造專用工具、不改符號。最初文字輸出留在暫存/work，已改正重跑並只刪自己產生的檔案；原版觀測器生成時的leading-plus編譯錯誤及原稿保留，屬工具問題。
- 靜態資料流證實從GoalXY往目前XY擴散；停止格一出隊就返回，回溯預算為波數減2。成本／軸順序決定轉角，不讀前驅。第一筆(1,44)往(5,43)，原版為(2,44)、(2,43)、(5,43)，舊Go為(5,44)、(5,43)。候選驗證後才修改正式明確平面路徑；goalPlane=-1的通用helper保留原相容語意，正式computePath明示平面，兩者共用擴散。
- [正式40拍比較](workplace/scratch/ui-windows-20261104/cao-path-after-comparison.json)每拍96槽27欄位與完整RNG全等；[12筆尋路](workplace/scratch/ui-windows-20261104/cao-path-evidence.json)的11筆完整點列及1筆失敗均吻合。檢查工具最初錯把所有請求視為成功，已依原版失敗分支修正，不丟棄樣本。全低平面表仍有586個導航旗標與26個成本值不同，語意未逐項映射，保留evidence-only限制，不另開全表重寫或宣稱全場一致。
- [四套件冷測／vet／建置](workplace/scratch/ui-windows-20261104/cao-path-tests.log)通過；新增精確加權轉角與下降／門格／不可攀登測試，既有上升獨立Z點亦通過。[舊錯版負對照](workplace/scratch/ui-windows-20261104/cao-path-negative.log)被精確三轉角測試拒絕。正式binary wlgame-cao-path-r1 SHA為901d4e58879dab7831d16a890ef9959e2701666492bee38143d22bb7f7a0e870。
- [正常滑鼠流程](workplace/scratch/ui-windows-20261104/cao-path-remake-gui.json)66個功能／同步檢查點、15圖通過：曹操宣戰、編夏侯淵、下令徐州、停手行軍與途中存讀檔，自然入戰、改線、停手移動、突擊及主動退卻返回。沒有直接開戰、注入部隊或強制勝負；自然GUI不當作同狀態整場對拍。
- [來源收據](workplace/scratch/ui-windows-20261104/cao-path-integration.json)封存三檔來源、binary、工具、種子及原始輸入雜湊；[本機頁](workplace/scratch/ui-windows-20261104/cao-path-game-review.html)四張正式畫面通過390寬[離線抽驗](workplace/scratch/ui-windows-20261104/cao-path-review-qa.json)，JS錯誤與外部請求皆0。spec134限本片CONFORMED，機制文件及目前狀態已同步；投射物、架梯、40拍後命令與完整結算仍沿Issues #13／#16／#24，未關閉Issue或作Release。
- [收尾核對](workplace/scratch/ui-windows-20261104/cao-path-final-checks.json)通過：三檔來源覆蓋包、32份證據、正式binary、五份原版輸入及頁面雜湊相符，1145列分流與嚴格生成輸出一致，582份文件／135條斷言索引通過。工作區root-owned及誤建.md目錄皆0，來源UID/GID1000:1000。所有本輪有界容器已結束刪除，未動其他專案資源。依使用者指示commit並push，實際提交與遠端核對見[本機推送收據](workplace/scratch/ui-windows-20261104/repo-path-update-20261112.json)。

### 2026-11-12：依玩家攻守選擇戰術AI腳本

- 從CONTEXT及現行Issues #16／#24、復古戰術對拍路由接手，base為d70097d。使用逆向重製與IDA技能，先延伸既有40拍收據，沒有重開已完成尋路或全導航表重寫。技能安裝連結位於另一個唯讀掛載，補掛IDA專案後讀取權威技能；沒有環境受阻聲明。
- 延伸到50拍仍吻合。探針按舊正式來源載入呂布戰法0、能力類別0的段0後，第69拍先有RNG差異、第70拍起命令不同；原版等待與取數指令本身未證明有錯。新增唯讀1A426／1A57A／1A587／1A591／1A5CA觀測，沿正常開局重播至116拍，117個共同快照的指令數、96槽、RNG及全域資料全部重現。
- [原版腳本觀測](workplace/scratch/ui-windows-20261104/cao-phase-script.jsonl)的256bytes與BATTLE.DAT段1完全相同，SHA為43af82e9e62ab72fe852d52029cada98f071e26e49bc1bfb877f69af9f271dbe。以實際段1隔離後前73拍吻合。第74拍有13槽差異，包含Target、地形及移動；證據檢查器最初只假定Target，讀完全部差異後已修正，保留完整欄位，不丟棄樣本。
- IDA Pro9.4經既有非root image及tools/ida.sh probe，匯出[腳本主分派](workplace/scratch/ui-windows-20261104/script-20261112/ida-probe.json)及[載入函式](workplace/scratch/ui-windows-20261104/script-load-20261112/ida-probe.json)。1CBEB先取0／1／2再加1，只對玩家守城的攻城類改槽0；1CC0C取對方戰法乘4加槽。原始bytes、KI.EXE與DB身分在spec59，未改原符號或正式DB。舊摘要省略加1與攻守閘，已同步訂正機制及格式。
- spec59先READY再實作。battle.ScriptSlot唯一處理AI主攻攻城0／防守攻城1／陸戰2／海戰3；TacticalSetup.Script新增玩家守方參數，由參戰軍團決定，桌面與手機共用來源。能力Category、地形、RNG、腳本等待及存檔不變，戰術測試fixture同步用正確腳本。
- [正式73拍比較](workplace/scratch/ui-windows-20261104/cao-phase-after-comparison.json)每拍96槽27欄位及完整RNG全等，腳本PC／wait／cond亦全等；[獨立證據](workplace/scratch/ui-windows-20261104/cao-phase-evidence.json)保存原版共同快照與腳本SHA。局部探針使用原版擺位／入口RNG及單兵初值，沒有成為正式快照載入功能，不代替自然GUI。
- [六套件冷測／vet／建置](workplace/scratch/ui-windows-20261104/cao-phase-tests.log)通過資料、戰場來源、戰術、世界、桌面及手機。新增四槽、戰法0／7、分類邊界與玩家攻守回呼測試；[舊選法負對照](workplace/scratch/ui-windows-20261104/cao-phase-negative.log)被拒絕。Go探針最初從/output執行碰到internal套件路徑限制，改用模組內同一檔案後重跑，屬驗證工具問題。
- 正式binary wlgame-cao-phase-r1 SHA為aabb3732f3b4a9d4bed42bff4e56372190cf86b10e4006a6cacf32d4077fec1b。[正常滑鼠](workplace/scratch/ui-windows-20261104/cao-phase-remake-gui.json)66個功能／同步檢查點、15圖通過，含途中存讀檔、自然入戰、改線、停手移動、突擊及正常主動退卻返回。首次因戰場在兩次讀取間退出而KeyError，原工具／log／trace／圖保留於before-exit-race-fix覆蓋包。第二次在舊30秒上限內只跑至開場30拍，trace持續前進，保留before-opening-timeout-fix覆蓋包；只將三個驗收等待上限延至120秒，拍點、返回條件與binary不變。修正取樣後同程式重跑通過，沒有以環境逾時改遊戲速度。
- [來源收據](workplace/scratch/ui-windows-20261104/cao-phase-integration.json)封存七檔受測來源、正式binary、原版工具與輸入；[本機頁](workplace/scratch/ui-windows-20261104/cao-phase-game-review.html)四張正式流程畫面供檢視。spec59本片CONFORMED，整體排程仍READY；第74拍鎖敵／地形／移動及投射物、架梯、完整結算仍沿Issues #13／#16／#24。其他角色只有回呼／資料抽驗，不作全部正常路徑完成聲明，原版素材、HTML及binary留本機。
- [收尾核對](workplace/scratch/ui-windows-20261104/cao-phase-final-checks.json)通過：七檔來源覆蓋包、48份證據、正式程式、五份原版輸入與頁面雜湊一致，1145列分流及嚴格生成輸出相符，582份文件／135條斷言索引通過。四圖離線頁在390寬無溢出、無外部請求及JS錯誤。來源UID/GID1000:1000，工作區root-owned及誤建.md目錄皆0。所有本輪容器已結束移除，未動其他專案資源。依使用者既有指示commit並push，實際提交及遠端核對保存於[本機推送收據](workplace/scratch/ui-windows-20261104/repo-phase-update-20261112.json)。

### 2026-11-12：地形原版列號與斜坡高度修正

- 從CONTEXT、現行Issue #16與復古對拍路由接手，base為64e53e6；使用逆向重製與IDA技能，補讀舊規格回填與文件職責入口。只追本案第74拍首次差異，不重開全導航表重寫。對照第73拍命令後／移動前全等，先查該拍移動完成點。
- dosgolem乾淨來源f79eab6db8c64a524a06535e53e5d32b0bbcf7c8沿正常曹操開局、固定seed1重播至116拍，新增唯讀19FA0及第73拍1B240觀測。[117共同快照](workplace/scratch/ui-windows-20261104/cao-ground-original.jsonl)的指令數、96槽、RNG及全域資料與前輪全等；[原始地圖／模型及96筆移動點](workplace/scratch/ui-windows-20261104/cao-ground-oracle.jsonl)證實第76場4096bytes與BATTLE.MAP相同、2048-byte定義與BATTLE.MDL組0相同，沒有另換地形。
- 根因為士兵一直保存含表頭的原版Y，規則Field卻吃Tiles的62列資料。原版(55,36)低平面0、Go誤讀下一列1；原版(55,26)低平面1、Go誤讀下一列無地面。完整緩衝候選消除走位差異後，另剩平面與StepZ接縫。
- IDA Pro9.4由既有非root image及tools/ida.sh probe取得[原指令／bytes](workplace/scratch/ui-windows-20261104/ground-20261112/ida-probe.json)。1B1D2依格號及+1E選平面；1B228–1B23E同時增減當前Z與StepZ，普通水平升降不切換+1E。原版第73拍槽61／62／63的Z1、StepZ1、PlaneHigh0直接分開了「Z大於0即高平面」的錯解。KI.EXE與DB雜湊、位址空間均在spec36。
- spec36先READY再實作。Field改64列，X／Y有效1–62；正式Provider讀完整MinimapRows，翻轉只動中間62列。水平移動同步StepZ差值，平面旗標改依OnWall。地形繪圖仍62列，ScreenPos與單兵／投射物／倒地改用原版座標；場景旗仍按地形列轉換，破壞重繪讀Field Y＋1，鏡頭／游標公式保留。
- [正式84拍比較](workplace/scratch/ui-windows-20261104/cao-ground-after-comparison.json)每拍96槽27欄位、完整RNG及腳本PC／wait／cond相同。低平面4096個原始導航byte全等，舊586格殘差由列號修正消除；高平面仍有455格差異，不猜成因或展開全表重寫。[獨立收據](workplace/scratch/ui-windows-20261104/cao-ground-evidence.json)保留第85拍槽73／79僅StepZ不同的新閘門及spec57勘誤backlink。
- [七套件冷測／vet／建置](workplace/scratch/ui-windows-20261104/cao-ground-tests.log)通過資料、Provider、戰術、世界、繪圖、桌面及手機。[四項舊錯版負對照](workplace/scratch/ui-windows-20261104/cao-ground-negative.log)分別拒絕錯列號、舊Y上限、斜坡高度／平面與物件投影。首次冷測發現舊邊界test寫死MaxY=Height−1，以及結構fixture從移除表頭處多讀兩列；依原版64列契約修正後同工具鏈、同命令冷跑通過，失敗log保留。破壞重繪測試原素材路徑少上一層而一直跳過，這次修正並以真實素材執行通過。
- 正式binary wlgame-cao-ground-r1 SHA為2b02bc377bdf1052aca262640449f46686f0b30aa859f2f16c3aa140d80ab93b。[正常滑鼠](workplace/scratch/ui-windows-20261104/cao-ground-remake-gui.json)66個功能／同步檢查點、15圖通過，含途中存讀檔、自然入戰、改線、停手移動、突擊及正常主動退卻返回。沒有直接開戰、改部隊或強制勝負；自然GUI不是同狀態整場對拍。
- [來源收據](workplace/scratch/ui-windows-20261104/cao-ground-integration.json)封存12檔受測來源與53份證據，[本機頁](workplace/scratch/ui-windows-20261104/cao-ground-game-review.html)四張正式畫面供檢視。spec36本片CONFORMED，spec57的規則座標舊說、re63、格式、機制及現況同步勘誤；整體排程仍READY。第85拍StepZ、高平面、其餘成本、投射物、架梯及完整結果仍沿既有Issues邊界。未作Release，原版、HTML、快照及binary留本機。
- [收尾核對](workplace/scratch/ui-windows-20261104/cao-ground-final-checks.json)通過：12檔來源、53份證據、正式程式及五份原版輸入雜湊一致，1145列分流與嚴格生成輸出相符，582份文件／135條斷言索引通過。四圖離線頁在390寬無溢出、無外部請求及JS錯誤。來源UID/GID1000:1000，工作區root-owned及誤建.md目錄皆0。所有本輪容器已結束移除，未動其他專案資源。依使用者既有指示commit並push，實際提交及遠端核對保存於[本機推送收據](workplace/scratch/ui-windows-20261104/repo-ground-update-20261112.json)。

### 2026-11-12：坡地換位與重畫旗標修正

- 從CONTEXT及現行Issue #16接手，base為9f751aa。復古對拍路由命中逆向重製技能、戰術驗證、規格閘門與文件職責；換位定位沿IDA技能，候選暴露後續更新缺口時回到RE，收尾再讀試玩與回填路由。唯一現況表沿CONTEXT，未重開已完成的地形列號與全導航表工作。
- [原版工具](workplace/scratch/ui-windows-20261104/cao-swap-original.go)沿正常曹操開局、固定seed1及前輪輸入重播。dosgolem來源f79eab6db8c64a524a06535e53e5d32b0bbcf7c8，[117共同快照](workplace/scratch/ui-windows-20261104/cao-swap-original.jsonl)的指令數、96槽、RNG及全域資料與前輪全同。新增1B732／1B798唯讀觀測取得467對換位，17對低平面坡地、3對不同StepZ、4對不同地形旗標；沒有高平面換位樣本。
- IDA Pro9.4使用既有py312-v1及tools/ida.sh probe，原EXE唯讀、一次性DB、非root輸出。[換位原指令](workplace/scratch/ui-windows-20261104/swap-20261112/ida-probe.json)證實雙方平面0且任一Z非0才交換+12h；[重畫原指令](workplace/scratch/ui-windows-20261104/swap-redraw-20261112/ida-probe.json)證實bit1按每兵重畫的時機依目前圖塊重算。輸入及兩份DB雜湊、原位址與bytes均在spec62。
- StepZ單項候選先消除第85拍差異，下一個旗標差異在第87拍。單純移除旗標交換會提早於第84拍不同，因此spec62回到DRAFT，查清1B2D3–1B300後再READY。正式swapWith按原條件交換StepZ，保留各自HighTerrain；drawUnitCollision在活兵重畫點依目前格重算。沒有更動換位資格、碰撞格表、AI加成、存檔格式或遊戲速度。
- [正式115拍比較](workplace/scratch/ui-windows-20261104/cao-swap-after-comparison.json)每拍96槽27欄位、完整RNG與腳本PC／wait／cond全等。[獨立收據](workplace/scratch/ui-windows-20261104/cao-swap-evidence.json)核對全部467對與第84拍槽73／79的高度交換；第85拍兩槽差異已消除。原版計數器116尚未進入本輪入口收據，故不把候選的116拍或更早長程資料當作新正式上限。
- [七套件冷測／vet／建置](workplace/scratch/ui-windows-20261104/cao-swap-tests.log)通過資料、Provider、戰術、世界、繪圖、桌面及手機。[四錯版負對照](workplace/scratch/ui-windows-20261104/cao-swap-negative.json)分別拒絕不交換StepZ、無條件交換、立即交換地形旗標及漏掉重畫更新。最初隔離probe以/output路徑執行而碰到Go internal限制，改用同檔案的/repo工作樹路徑後重跑通過，未改正式程式。
- GUI前兩次分別在退卻第289拍及等待命令取樣第139拍碰到120秒觀測上限。第139拍同時有喊話，一度誤判它阻擋時間；讀battletalk.go與battle.go後證實喊話自動隨節拍消失、不阻擋下令，撤回該判斷及試作觀測器。[第一份失敗包](workplace/scratch/ui-windows-20261104/cao-swap-before-retreat-timeout-fix.tar.gz)、[第二份失敗包](workplace/scratch/ui-windows-20261104/cao-swap-before-native-talk-fix.tar.gz)及[撤回觀測器](workplace/scratch/ui-windows-20261104/cao-swap-withdrawn-native-talk-observer.tar.gz)保留，勘誤見[觀測器訂正](workplace/scratch/ui-windows-20261104/cao-swap-observer-correction.json)。最後恢復基準滑鼠輸入，只將三段等待及退卻上限改300秒，CPU配額2改4，原image、binary、遊戲速度及拍點條件不變；乾淨重跑通過。
- 正式程式wlgame-cao-swap-r1 SHA-256為d242565a3129d7ca7d839cc7d0c4752a315d3909f81867655f3672600e9e5f70。[正常滑鼠](workplace/scratch/ui-windows-20261104/cao-swap-remake-gui.json)66個功能／同步檢查點、15圖通過，含曹操對呂布開局、編夏侯淵、停手行軍、途中存讀檔、自然入戰、改線、突擊與正常主動退卻返回。這是remake正常流程，沒有直接入戰或強制勝負，不等於整場原版同狀態對拍。
- [來源收據](workplace/scratch/ui-windows-20261104/cao-swap-integration.json)、[本機頁](workplace/scratch/ui-windows-20261104/cao-swap-game-review.html)保存本輪正式程式與證據。spec62本片CONFORMED，spec36、re11、機制、正常試玩與現況已回填；整體排程仍READY，Issues #13／#16／#24保持開啟。高平面導航、其餘成本、投射物、架梯及完整結果不由本片推定完成；其他角色與手機GUI未擴張驗收。素材、快照、binary及HTML只留本機，未作Release。
- 文件嚴格檢查要求spec62明列本片無新增未解項，回填護欄另抓到規格只列範圍而缺少精確1B783鍵；補齊範圍說明及原指令後重跑通過，未改正式程式。[收尾核對](workplace/scratch/ui-windows-20261104/cao-swap-final-checks.json)驗證三檔來源、98份證據、五份原版輸入、115個邊界及467次原版換位，1145列分流與嚴格生成輸出一致，索引582份文件／135條斷言通過。四圖離線頁在390寬無溢出，無外部請求及JS錯誤。來源UID/GID1000:1000，工作區root-owned及誤建.md目錄皆0；本輪容器均已結束移除，未動其他專案資源。依使用者既有指示commit並push，提交及遠端核對保存於[本機推送收據](workplace/scratch/ui-windows-20261104/repo-swap-update-20261112.json)。

## 2026-11-12：戰後三秒短提示與可見軍團補驗

- 使用者指出長條戰果與原版不同，另指出許昌截圖看不到夏侯淵與呂布。沿復古remake驗收、規格閘門與文件職責路由查證，使用grill-me／grilling的本機三案比較；使用者明確選「B：三秒短提示，點開明細」。spec89先列READY再實作，原版結果頁的既有選項保持獨立。
- 接通戰術返回與快速事件的共用UI快照。提示由首次繪製計三秒，不暫停世界；點提示或工具列戰報可查最近一場，明細沿關閉鈕、返回、確認與右鍵操作，閱讀才暫停。新戰鬥取代上一場，讀檔／新遊戲清除舊戰報，不改存檔格式或戰術規則。移除長條戰果與原始據點損害值；原版對白及其他事件路徑保留。工具列十二鈕單排，仍28高展開／6高收起，四語文字存JSON。
- 軍團圖示核對沿用前輪未修改的wlgame-cao-swap-r1。前輪行軍及戰後截圖的視野沒有涵蓋兩軍，不能由狀態欄位代證可見。正常軍團情報定位後補拍許昌編成、停手行軍與椎陽兩勢力同框；dosgolem從正常輸入重生原版。原版三個／remake四個局部圖塊的形狀及顏色均0差，限曹操靜止／朝X增加及呂布靜止，不稱完整110張或整場parity。原版兩次錯定位與移走游標但未執行指令的擷取留存；補正輸入座標、執行300,000道指令消費滑鼠後核對，未改比較容差。入口為[playtest50](docs/playtest/50-corps-on-map.md)、[像素結果](workplace/scratch/ui-windows-20261104/corps-map-pixels.json)。
- 正常戰鬥仍是夏侯淵攻、呂布守椎陽。另從正常曹操開局編夏侯淵後留守許昌，觀察至196/4/20 2時，呂布留在椎陽，未發生交戰。此指定情境未驗收，不以強制入戰或袁術反攻汝南的舊案例替代；[有界觀察](workplace/scratch/ui-windows-20261104/xuchang-defence-gui.json)保存結果。
- 受影響桌面／排字冷測與vet通過。第一輪測試把語系提示誤當按鈕標籤，既有單排測試亦未識別新增action；修正測試列舉後乾淨重跑。首輪正常GUI106點通過，但HD截圖只有底框，該圖不作文字驗收；獨立GPU逐像素合成通過，未重現持續空白，未據此改HD合成。補實際截圖六個內文區域的文字像素檢查，並修正短提示內容區的8px欄寬不足後，以wlgame-battle-notice-r2重跑109個功能／同步檢查點、20圖通過，包含重複同步，不作獨立功能數。正常兩次自然遭遇驗到三秒自行收起、停手日期前進、工具列重開、原版／HD明細、四種關閉、讀檔失效及點提示無地圖穿透；原版與原生存檔兩者SHA未被介面重寫。
- 第一回GUI因容器工作目錄錯誤讀不到translations/corrections.json，分類為驗證環境錯誤；修正工作目錄後用同一image重跑。先前單元／首輪GUI與不完整HD截圖保存在battle-notice-before-text-capture-gate.tar.gz。新增十圖[本機實際畫面頁](workplace/scratch/ui-windows-20261104/battle-notice-game-review.html)，四項離線抽驗通過390／1280寬、十圖載入、無JS錯誤／外部請求；[來源收據](workplace/scratch/ui-windows-20261104/battle-notice-integration.json)含十檔來源覆蓋包、binary、工具與原版身分。HTML、原圖、原始狀態與binary皆留本機，不作Release。驗收範圍沿[spec89](docs/spec/89-siege-damage-report-toggle.md)，前115拍以外、高平面、投射物、架梯與完整戰況維持既有Issues限制。

- 收尾檢查通過：十檔受測來源及覆蓋包、100份證據、109點GUI／20圖、七項軍團局部像素與原版／seed身分一致；來源UID/GID1000:1000，root-owned與誤建.md目錄皆0。嚴格未解索引1145列、三筆分流回填、582份文件／135條斷言、過期斷言與資產deny-list通過。第一次索引在本機頁尚未產生時拒絕連結；頁面完成後乾淨重跑，未將此驗證順序問題列為產品缺陷。[收尾核對](workplace/scratch/ui-windows-20261104/battle-notice-final-checks.json)保存雜湊與範圍。本輪容器均已結束移除，沒有臥龍傳殘留容器；未重建或清理其他專案image。依使用者既有指示提交並推送，提交與遠端核對保存於[本機推送收據](workplace/scratch/ui-windows-20261104/repo-battle-notice-update-20261112.json)。

## 2026-11-13：戰報預設安靜與省略右鍵教學

- 使用者指定「以滑鼠右鍵恢復」提示不必要，另指定戰報改成選項、預設不打開，避免多處交戰頻繁打斷。命中復古remake與規格閘門路由，沿用逆向重製技能及既有spec。查明右鍵教學為原版TALK #4，勢力選擇及軍團情報各有一個入口；依新裁定省略顯示，保留原文、索引、右鍵返回及其他任務提示，明列remake差異。
- spec89／266先列READY再實作。系統末尾加入「戰報提示」第十列，預設關閉，原九列座標不動，視窗底部仍在640×400內。關閉時只保留最近一場供工具列手動看；啟用後只提示下一場，停用或重開不補播舊戰報。沿既有獨立偏好新增相容布林欄位，舊偏好缺欄位仍為關；不改遊戲存檔或戰鬥規則。
- 受影響桌面冷測與vet通過。r1正常曹操對呂布流程118個功能／同步檢查點、19圖通過，驗到自然戰鬥返回四秒無提示、日期繼續前進、手動戰報、切換不補播、開關兩種重啟記憶，兩種遊戲存檔SHA不變。r2另補軍團情報同一句教學的省略，最終程式再跑44個功能／同步檢查點、三圖，正常勢力選擇、編成後軍團情報及右鍵返回均目視通過。數字包含重複工具列同步，不作獨立功能數。短驗收腳本首次漏開原版指令列，分類為驗證腳本問題；補正常展開操作後，以同一image、binary乾淨重跑通過，原log保留。r2未再重跑整場戰鬥；啟用後新結果提示由本輪模型回歸及前片顯示收據支持，不冒稱新的正常戰鬥GUI驗收。
- [本輪來源收據](workplace/scratch/ui-windows-20261104/quiet-notice-integration-20261113.json)核對13檔來源、覆蓋包、兩個binary及22圖；來源UID/GID1000:1000，root-owned與誤建.md目錄皆0。截圖、binary與原版素材只留本機。當輪行為與限制見[spec89](docs/spec/89-siege-damage-report-toggle.md)及[spec266](docs/spec/266-mouse-only-desktop-toolbar.md)。
- 收尾文件檢查通過：582份文件／135條斷言、嚴格未解索引與分流1145列、過期斷言及原版資產deny-list均通過。受測13檔來源與兩個binary仍一致。本輪容器已結束移除，沒有臥龍傳殘留容器；未重建工具image或清理其他專案資源。依既有指示提交推送，提交及遠端核對見[本機推送收據](workplace/scratch/ui-windows-20261104/repo-quiet-notice-update-20261113.json)。

## 2026-11-13：登城首拍中繼點與高層地面

- 沿復古remake、規格閘門、IDA9.4、文件職責與RE回填路由接續既有Issue #16。維持已確認的滑鼠操作、預設關閉戰報及原版規則，不開新設計訪談。原版正常曹操開局、固定seed1的輸入只延長最後戰術節拍，取得182個快照；原先117筆指令數、單位、RNG與全域資料前綴完全相同。
- IDA Pro9.4重新讀`1A7B7–1A839`、`1AB39–1AB7C`、`1BBA6–1BC39`及`1BCA6–1BD07`。新登城命令同時寫StepXY／Z，原Go只更新Goal；無低層地面時原版仍掃高層，原Go提前回傳。spec36先列READY再修改，正式實作限這兩處；同命令後續拍保留繞路點。未改RNG、時間基準、存檔格式或AI加成。
- 原指令另證實`1BC18`只清DH導航旗標，DL高度仍保留，推翻舊re63「整格清0」摘要。已勘誤原始文件與程式註解，未擴張為全部正規化重寫。低層4096原始byte仍全等，高層差異455降為60。勘誤與未受影響引用見[回填表](workplace/scratch/ui-windows-20261104/cao-scale-backlinks.json)。
- [正式逐兵](workplace/scratch/ui-windows-20261104/cao-scale-after-comparison.json)前138拍96槽27欄位及完整RNG吻合。[腳本比對](workplace/scratch/ui-windows-20261104/cao-scale-script-comparison.json)發現第132拍cond原版26、Go7，故含腳本的完整狀態聲明限前131拍；第139拍守陣命令仍不同。原版入口受控規則重播與正常GUI是兩種證據，不宣稱整場同狀態或全部動畫完成。
- [七套件冷測／vet／建置](workplace/scratch/ui-windows-20261104/cao-scale-tests.log)通過；新測試涵蓋獨立高層、門格條件、城壁替代圖塊及雙側隊長／隊員登城首拍與後續路徑。兩個舊錯版分別被回歸測試拒絕。最後只再修正正規化註解，[Go詞元核對](workplace/scratch/ui-windows-20261104/cao-scale-comment-check.json)確認與受測來源的執行詞元相同，原受測來源及binary保留。
- [正常滑鼠](workplace/scratch/ui-windows-20261104/cao-scale-remake-gui.json)66個功能／同步檢查點、15圖通過，含曹操開局宣戰呂布、編夏侯淵、停手行軍、途中存讀檔、自然入戰、改線後移動、突擊及主動退卻返回；計數含重複工具列同步。已目視戰術畫面有城牆與單兵，沒有強制入戰或勝負。仍是夏侯淵攻、呂布守椎陽，未代證許昌留守情境或手機GUI。
- IDA預設locked-v1標籤不存在，檢查官方manifest後沿用已存在且相容的py312-v1，不另建image。早期探針缺overlay、Go internal路徑及日誌命名錯誤均屬驗證腳本／環境，修正後同容器工具鏈乾淨重跑，未列產品缺陷。早期共用nav輸出曾被候選覆寫，已排除該檔，以命名固定的cao-scale-after-nav.bin為正式收據。
- [來源收據](workplace/scratch/ui-windows-20261104/cao-scale-integration.json)保存四檔來源、77份證據、binary與五份原版輸入身分。spec36本片CONFORMED，整體spec159仍READY；腳本、後續守陣、高層導航、架梯、投射物及完整結果沿既有Issues保留。原版素材、IDA資料庫、binary與截圖皆只留本機，不作Release。
- 收尾通過索引582份文件／135條斷言、嚴格未解與分流1145列、過期斷言及原版資產deny-list。全庫引用掃描仍有204筆歷史workplace缺件，與HEAD用同一檢查器核對相同，本輪新增0；歷史census.tsv缺席的符號層未執行，不宣稱全庫掃描全綠。來源覆蓋包、77份證據及原版輸入雜湊一致，UID/GID1000:1000，工作區root-owned與誤建.md目錄0；本輪容器已結束移除，未重建image或動其他專案資源。[收尾核對](workplace/scratch/ui-windows-20261104/cao-scale-final-checks.json)及[引用差量](workplace/scratch/ui-windows-20261104/cao-scale-doc-phantom-delta.json)保存結果。依既有指示commit並push，提交與遠端核對另存本機repo-scale-update-20261113.json。
