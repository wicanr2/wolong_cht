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
