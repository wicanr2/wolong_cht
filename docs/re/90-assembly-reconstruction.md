# 90：組語基準：完整重建松崗版 KI.EXE

**狀態：整檔組語重建通過。指令來源、MZ 封裝與資料宣告可獨立產生相同 EXE。**

- 日期：2026-10-08
- 使用者決定：先使用組語建立基準，後續由組語還原 C 函式
- 範圍：松崗 DOS/V `KI.EXE`，不外推 PC-98 或其他執行檔
- 原始檔與重建檔 SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 一次性 IDA DB SHA-256：`e32713237d223145834531bffc8979146ed4bacaa98f73ec4bbf692b747c2c4f`
- 工具：IDA Pro 9.4、GNU assembler／ld／objcopy 2.40
- 位址空間：IDA database linear address；檔案偏移與原始段值另列，不混用
- Go 基線：`e8aeb99`，本輪沒有修改遊戲規則

本檔保存組語基準建立時的收據與台帳快照。後續 C 還原現況見
[`91`](91-c-rng-restoration.md) 與 [`c-recovery-status.json`](c-recovery-status.json)。

## 1. 已完成的基準

| 驗證項目 | 結果 | 等級與範圍 |
|---|---|---|
| 完整 EXE | 67,099 bytes，逐位元組與 SHA-256 相同 | 已證實，整檔重建 |
| IDA 辨識的指令 | 24,376 條，55,392 bytes，全部用助憶碼組譯 | 已證實，IDA 指令 inventory |
| 未組譯指令的原始 bytes fallback | 0；任一指令不匹配即停止 | 已證實，生成工具的檢查 |
| MZ 標頭 | 512 bytes，欄位與 padding 保留 | 已證實，檔案格式與 bytes |
| MZ 重定位表 | 116 筆完整保留，113 條指令還原載入前的段值 | 已證實，原始檔與 IDA loader 對照 |
| 其餘映像 bytes | 11,195 bytes，由資料宣告保存 | 已證實，bytes；用途仍依既有分級證據 |
| 來源對映 | 從檔案起點到檔尾沒有缺口或重疊 | 已證實，67,099 bytes 全覆蓋 |
| C 函式定位台帳 | 739 筆原始 IDA 名稱、chunks、來源行與語意等級 | 導航資訊，不能當作 C 還原完成率 |

所有指令都先以 IDA 結構化運算元轉成組語，再以 GNU assembler 的實際輸出逐條比對。
資料宣告只使用於標頭、padding 與非指令區段。指令不能靠 `.byte` 或 `.incbin` 繞過匹配。

原始檔從未改寫。生成的整份來源含原版程式與資料，只留在被 Git 忽略的本機研究目錄，
沒有加入公開引擎包或發行產物。

## 2. 編碼與載入位址

| 問題 | 處理方式 |
|---|---|
| 同一條指令有 load／store 兩種方向編碼 | 逐條試組譯，選擇 bytes 相同的形式 |
| 原版保留 16-bit 立即數，assembler 會縮短 | 在 linker script 明示 66 個立即數常數，組譯時保留 16-bit 重定位，再由 ld 解決 |
| IDA 把 far call 的段值加上載入基址 | 用 MZ 重定位表定位欄位，還原檔案裡的原始段值 |
| 程式碼與資料混排 | 保留原始檔案位置與 `.org`，來源對映逐段檢查 |
| IDA 的 item 尾 byte | 分類時回到 item head，避免把指令尾算成未知資料 |

位寬與編碼方向的表示依
[GNU assembler 2.40 文件](https://sourceware.org/binutils/docs-2.40/as/i386_002dMnemonics.html)，
查閱日期 2026-10-08。具體匹配由組譯與整檔 bytes 證明，不從文件推定通過。

## 3. 獨立組譯與負對照

只掛載 `KI.reconstructed.S` 與 `KI.ld` 的新容器，經 `as → ld → objcopy` 冷重建，
輸出的 SHA-256 仍為原始檔的完整雜湊。此驗證保存於 `assembly/source-only/KI.EXE`。

| 改動 | 驗證結果 |
|---|---|
| 把一個 `clc` 改成 `stc` | 整檔雜湊不同，被拒絕 |
| 改一個 linker 立即數常數 | 整檔雜湊不同，被拒絕 |
| 改第一筆 MZ 重定位位移 | 整檔雜湊不同，被拒絕 |

比較沒有遮掉位址、立即數、資料或重定位欄位。

## 4. 組語到 C 的入口與限制

函式入口先查 `function-ledger.json`，再以原始 IDA 名稱、線性位址、chunks 與
`assembly_line` 開啟來源。`source-map.json` 保留原始反組譯與運算元，並附分級語意索引。
沒有分級語意的項目帶警示，不以 assembler 的新常數名稱替代原始定位。

| 台帳資訊 | 含義 |
|---|---|
| `instruction_bytes_exact` | 該 IDA 函式 chunks 中的已辨識指令有相同 bytes |
| `semantic.level` | 既有證據的語意等級；本次匹配不自動升級 |
| `boundary_authority` | 函式邊界來自 IDA，沒有證明原作者的來源切分 |
| `c_status` | 本次全部為 `not-started` |

C 還原依既有閘門進行：原版證據、DRAFT 介面與副作用契約、證據審查、READY，
再實作與驗證。輸入、回傳值、暫存器保存、旗標、記憶體寫入及 caller／consumer
都要有可回查證據。C 的機器碼匹配與行為等價分開記錄，不由組語匹配代為證明。
原版動態行為仍使用專案的 dosgolem 入口，見 [`spec/131`](../spec/131-dosgolem-oracle.md)。

## 5. 未解與證據限制

| 項目 | 邊界 |
|---|---|
| 原作者語言、編譯器與組譯器 | 仍未知；本次使用 GNU assembler 能重建，不代表原作者使用它 |
| IDA 程式碼／資料分類 | 保留分析工具的定位；指令能匹配不證明每個被標為程式碼的區段曾經執行 |
| 函式語意與間接控制流 | 本次沒有逐函式新增玩法語意或宣稱全部解讀 |
| C 還原與 C matching | 尚未開始；台帳提供入口，沒有 C 通過聲明 |
| Go remake 行為 | 本次沒有新的同狀態或正常玩家路徑對拍 |

這條研究線不增加 remake 的既有發行閘門。當前玩家工作仍查
[GitHub Issues](https://github.com/wicanr2/wolong_cht/issues)；函式台帳是證據索引。

## 6. 重跑與本機產物

```sh
tools/assembly_rebuild.sh
```

包裝器經 [`tools/ida.sh`](../../tools/ida.sh) 建一次性 DB，再由
[`ida_assembly_inventory.py`](../../tools/ida_assembly_inventory.py) 匯出全部 segments。
[`assembly_rebuild.py`](../../tools/assembly_rebuild.py) 組譯、核對 bytes、建立來源對映與
函式台帳，並執行來源冷重建與負對照。全部 workload 都在非 root、限資源、無網路的 Docker
容器內，原始輸入唯讀掛載。

研究根目錄：`workplace/matching-decompilation/assembly/`。

| 產物 | 用途 |
|---|---|
| `dosv/ida-probe.json` | 固定輸入身分、指令、重定位、原始運算元與分級語意 |
| `dosv/input.exe.i64` | 本次一次性 IDA DB |
| `build/KI.reconstructed.S`、`build/KI.ld` | 可獨立組譯的完整來源 |
| `build/KI.EXE` | 與原版相同的本機研究執行檔 |
| `build/source-map.json` | 檔案 bytes 到來源行與原始定位的對映 |
| `build/function-ledger.json` | C 函式還原的原始定位與證據入口 |
| `build/report.json` | 整檔與逐指令匹配、工具身分、雜湊與負對照 |

| 本次身分 | SHA-256 |
|---|---|
| 組語來源 | `47571181edc1a1fe7a2554f0e94fc1c9bcf87cc2e42e06c7fcc093c871fcae62` |
| Linker script | `a301287e77bf1a79a717c206dc0754e2c0ccb3a75b3eae38815bfff0a1995ac7` |
| 收據 | `6cc91dc43c529b8281bceae7b36036a4a744d0c90cd0e5ea7500e3c03a3abf73` |

IDA image ID：`sha256:4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab`。
組譯 image ID：`sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e`。
恢復入口與先行試點在 [`89`](89-matching-decompilation-pilot.md)，本輪過程在
[`WORKLOG.md`](../../WORKLOG.md)。

## 7. 版控程式碼與冷重建紀錄

2026-10-08，依使用者要求，完整指令來源另存至 GitHub repo。
本節追加版控產物與驗證，前述含資料的完整來源與歷史收據仍保留在本機。

| 版控產物 | 用途 |
|---|---|
| [`KI.code.S`](../../tools/c_recovery/KI.code.S) | 全部 24,714 條已匹配指令，保留原始 IDA 名稱、線性位址、檔案偏移與運算元，自動合併分級語意、出處與未知警示 |
| [`KI.code.ld`](../../tools/c_recovery/KI.code.ld) | 66 個編碼所需的立即數常數與區段配置 |
| [`assembly-code-record.json`](assembly-code-record.json) | 輸入、DB、來源與工具雜湊，75 個指令範圍與逐範圍雜湊 |
| [`assembly-code-verification.json`](assembly-code-verification.json) | 冷組譯、完整 EXE 比較與三個負對照的收據 |
| [`matching_code_record.py`](../../tools/matching_code_record.py) | 從已驗證基準匯出指令來源，以及本機組譯與資料匯入 |
| [`matching_code_record.sh`](../../tools/matching_code_record.sh) | 使用固定 GNU binutils 2.40 image 的 Docker 重跑入口 |

組語來源不含 `.byte`、`.word`、`.incbin` 或字串資料宣告。非指令範圍使用 `.org`
零值佔位。組譯後先比較全部 56,197 個指令 bytes 與逐範圍雜湊，再只從使用者自備的
固定 SHA-256 原版匯入 10,902 bytes。匯入範圍包含 512-byte MZ 標頭與其餘非指令區，
匯入的指令 bytes 為零。完整 67,099-byte EXE 與原版逐 byte 及 SHA-256 相同，等級為已證實。

把 `clc` 改為 `stc`、改一個 linker 常數，都使實際組譯的指令範圍不同而被拒絕。
改動私有輸入的標頭則被完整輸入雜湊拒絕。原版 EXE、資料區與重建 EXE 均留在本機。
來源與 linker SHA-256 分別為 `5b7e8e04b7ad80ba0d5176da5b33d150e12b16b32a3c083ff31e8b88b27884ba`
與 `a301287e77bf1a79a717c206dc0754e2c0ccb3a75b3eae38815bfff0a1995ac7`。

```sh
tools/matching_code_record.sh
```

重建輸出在 `workplace/matching-decompilation/assembly/code-record/`。
首次重生指令紀錄時，以 `matching_code_record.py export --repo /repo` 在容器中讀取
既有 `assembly/build/` 基準，只有 `tools/c_recovery/` 與 `docs/re/` 可寫；原版與基準唯讀。
一般驗證只需要版控來源、索引與自備原版，不依賴本機 IDA DB 或含資料的組語來源。

指令 bytes 已匹配不會提升函式語意等級。語意仍查
[`matching-semantic-index.json`](matching-semantic-index.json)，C 現況查
[`c-recovery-status.json`](c-recovery-status.json)。C 來源同樣進版控，局部行為等價與 C 機器碼匹配分開記錄。

2026-10-09 補充：原始分派表的二十 bytes 雖為 IDA 資料項，直接解碼已確認八條指令。
[rectangle-handler-code.json](rectangle-handler-code.json) 與 [re/106](106-c-rectangle-bars-restoration.md)
保存原資料分類、原始位置與工具身分；目前版控來源合併後為 24,714 條。§1–6 保留初建基準快照。

顯示清單後續 code／data 審查見 [re/107](107-c-display-interpreter-restoration.md)，追加其餘
handler、直線、底紋與雙色框的固定 decoder／組譯證據。原始初建與八指令研究快照保留。

原始 glyph 的 21 條／59 bytes 補充見 [re/108](108-c-glyph-raster-restoration.md)，
包含原始兩 patched far call 與完整三字 code；資料字庫仍只作本機輸入。

尋路自我修改區段的79條新增來源見[re/129](129-c-route-restoration.md)及[route-handler-code.json](route-handler-code.json)。現行版控組語為25,078指令／57,045bytes，整檔仍67,099bytes；來源與驗證入口沿用本頁，歷史checkpoint數字不改寫。

2026-10-10，交戰與戰術來源審查以31條新／替換指令取代12條舊錯邊界，最新版控來源為25,097條指令／57,101指令bytes，完整67,099-byte EXE仍與原版相同。原始與退休指令的定位及雜湊保存在[補充收據](c-engagement-code.json)，研究入口為[re/131](131-c-engagement-tactical-restoration.md)。三個組譯負對照均被拒絕；此結果不代表C原生矩陣或C機器碼匹配已完成。
