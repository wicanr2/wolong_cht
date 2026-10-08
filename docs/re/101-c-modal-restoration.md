# 101：C 外交與金額視窗控制流

**狀態：十七函式 O0/O2 各 53,322 組原版/C 相同，十個負對照拒絕；限明示底層 primitive／輸入／裝置 fixture。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-modal/ida/input.exe.i64`
- DB SHA-256：`42181f593431ef87bb95cb5e66c6945fc7e92c86d0d8ef8f9e4df23209f338d2`
- 原始 chunks／operand／xref／重定位：[`ida_modal_probe.py`](../../tools/ida_modal_probe.py)
- 契約：[`spec/220`](../spec/220-c-modal-control.md)

## 1. 原始定位與分級語意

十七函式共 1280 bytes。每筆原始名稱、位址、operand、chunks、IDA bytes 與 file bytes
保留，語意只作分級註記。以下 raw 控制流為已證實，限本輪 primitive 與受控回應矩陣。

| 原始函式 | IDA 起點 | 檔案起點 | bytes | loader 重定位 | 已證實的 raw 語意 |
|---|---|---|---:|---:|---|
| `sub_12078` | `0x12078` | `0x2278` | 94 | 7 | 局部滑鼠範圍／座標的原始 FAR 呼叫與保存 |
| `sub_120D6` | `0x120D6` | `0x22D6` | 123 | 7 | 世界滑鼠範圍、word 座標與原始 FAR／保存 |
| `sub_138C7` | `0x138C7` | `0x3AC7` | 31 | 0 | 停戰 modal 包裝、主公取址、保存與遠呼叫接線 |
| `sub_138E6` | `0x138E6` | `0x3AE6` | 28 | 0 | 協力 modal 包裝、主公取址、保存與遠呼叫接線 |
| `sub_13902` | `0x13902` | `0x3B02` | 230 | 3 | 原始外交回應、金額、CF retry、RNG／trust 與 SS frame |
| `sub_139E8` | `0x139E8` | `0x3BE8` | 288 | 3 | 原始金額下界、回應、CF retry、raw +0x1A byte／扣款與 SS frame |
| `sub_13C3D` | `0x13C3D` | `0x3E3D` | 92 | 0 | 回應 formatter、byte 條件及信賴度 raw 參數 |
| `sub_13B7E` | `0x13B7E` | `0x3D7E` | 43 | 0 | 原始三選一的 CF 輪詢、參數與返回 |
| `sub_13D09` | `0x13D09` | `0x3F09` | 60 | 0 | 原始場景 block 參數、word MUL／保存與 renderer ABI |
| `sub_13D45` | `0x13D45` | `0x3F45` | 35 | 0 | 原始場景還原參數、保存與 renderer ABI |
| `sub_13C99` | `0x13C99` | `0x3E99` | 39 | 0 | 原始說話型減三一次、訊息與肖像參數 |
| `sub_13CDC` | `0x13CDC` | `0x3EDC` | 45 | 0 | 原始玩家軍師記錄與下框訊息／肖像參數 |
| `sub_19321` | `0x19321` | `0x9521` | 21 | 0 | 原始 byte 表格取址與恢復服務參數 |
| `sub_187FF` | `0x187FF` | `0x89FF` | 17 | 0 | 玩家主公 raw 記錄的原始取址 |
| `sub_13D68` | `0x13D68` | `0x3F68` | 41 | 0 | 原始場景框參數、記錄與 renderer ABI |
| `sub_11D46` | `0x11D46` | `0x1F46` | 72 | 0 | 原始 resume、byte flags gate、world／畫面 helper 呼叫次序與保存 |
| `sub_12216` | `0x12216` | `0x2416` | 21 | 0 | 原始 code word patch、callee 入口觀測及恢復 |

## 2. MZ 重定位與遠返回

原始 MZ 表共 116 筆重定位，本輪十七 chunk 包含二十筆。每一筆 file segment word
為 0x1000，IDA loader 加 paragraph 0x1000 後為 0x2000。probe 逐筆驗證原始表，
同列保存兩種 bytes 與 word，沒有遮罩或忽略差異。第一次未處理 relocation 的直接比較
在這二十個位置失敗；原始輸出與 DB 保留為 `ida/ida-probe-pre-relocation.json`、
`ida/input-pre-relocation.i64`，沒有把工具位址空間差異當作遊戲缺陷。

dosgolem 本輪 load paragraph 是 0x0110，實際遠 callee 是 0x1110:0000；IDA 對應位址
單獨記為 0x20000。runtime 指令 bytes 依原始 MZ 表加 0x0110 逐 word 核對，C 記憶體
採相同載入影像。FAR 入口記錄實際 CS:IP 及兩個 return words，最後執行 RETF，不能替換 near RET。

## 3. 接線與受控 leaf 契約

[`modal.c`](../../tools/c_recovery/modal.c) 接回原始名稱的 C 入口，以及既有事件、經濟、RNG／trust。
真正執行外交採納、回應 formatter、金額下界、扣款、說話型、SS frame、CF retry 及原始恢復次序。
PUSHF／POPF 隔離數值輸入回傳 CF 與後續 FAR clobber。SS 未初始化槽不補預設值，
兩側保留同一初始 RAM 與固定案例序列的原始 stack。0x12216 的中途 code word 修改
在 callee 入口追加觀測，原始 word 隨後恢復，不用最終 bytes 相同代替中途證據。

| 位址基準 | fixture 行為 |
|---|---|
| runtime 0x1110:0000，IDA 0x20000 | AX=3 時提供相同固定 CX／DX，其餘回傳不改 register；CMP AX,AX clobber FLAGS，RETF |
| IDA 0x19796 | 固定初始 AL，FLAG 不變 |
| IDA 0x193E9 | 固定選單 AX，先回 0–3 次 CF=1，再回 CF=0；由真正選單迴圈輪詢 |
| IDA 0x17C6E | 固定金額 AX，先回 0–3 次取消 CF=1，再回 CF=0；由真正 caller 回選單重試 |
| IDA 0x1E453 | 固定 AL，兩側相同，原始 resume branches 仍執行 |
| IDA 0x1222B | RET；caller 的暫時 code word 另記錄 |

RET primitive 位址（IDA linear）為 0x15E80、0x10CE7、0x10CDE、0x18810、0x102F5、0x11CB1、
0x145F8、0x14236、0x15E60、0x197C3、0x1075B、0x1E38C、0x189A4、0x187AF、0x10C14、0x1FA37、
0x102C2、0x10241、0x11F30、0x15C58、0x1D615、0x11CC9、0x1D66A。
這些圖形／文字／音效／裝置／其餘世界 callee 的完整實作不在本輪證據範圍。

數值編輯器實際 AX=30000 是上限，CF=1 是取消，已有 [spec/78](../spec/78-amount-input-editor.md)
的原版契約。本矩陣另注入 32768／65535 raw 回應來量 caller 的 word 邊界，不能當作合法玩家輸入。
選擇 CF 忙碌與數值 CF 取消的控制流分開，不以它們推導硬體時間或真實滑鼠時序。

## 4. 完整矩陣與負對照

| 群組 | O0/O2 各原版/C |
|---|---:|
| window：兩入口 × 四個座標軸各五個 word 邊界 | 1,250 相同 |
| helper：九 helper × 全部 byte variant／flags／selector | 2,304 相同 |
| selector：四初值 × 四回應 × 四 CF 輪詢 × 三訊息狀態 | 192 相同 |
| reply：四回應 × 五個 AH × 兩玩家 × 五信賴度 | 200 相同 |
| diplomacy：四模式 × 三 AH × 五原金額 × 四回應 × 七金額 × 三取消重試 × 八 raw RNG | 40,320 相同 |
| trust-rng：全部 raw RNG × 三信賴度 × 四金額 | 3,072 相同 |
| amount：七原金額 × 八新金額 × 四回應 × 三取消重試 × 三說話型 | 2,016 相同 |
| wrapper：兩包裝 × 八 raw RNG × 五原金額 × 四回應 × 七金額 | 2,240 相同 |
| scenario-modal：四劇本 × 三玩家 × 六 handler × 四回應 × 三金額 × 兩種 CF 回應 | 1,728 相同 |
| 合計 | **53,322 相同** |

所有案例比十四 register／segment／FLAGS、世界、queue、globals、RNG、cadence、stack
與控制欄位。每個原始 callee 入口另記 252 bytes：上述原始 ABI 與 64-byte UI／input／code patch。
每版 417 次完整 1 MB 抽查相同，十七函式皆有 `entries_seen` 實際入口。
四劇本用受控輸入直接從六個原始 handler 進入，不能當自然 producer 或正常玩家 GUI 路徑。

O0/O2 JSON 相同，SHA-256：`7b26942d563e4c0b5ddc490852fcd7024a906b9bea52a4f64009aa4c37cd3552`。
原版及 C 分別讀取各自記憶體，狀態摘要同為
`3a12c347e0868495ca15eb24b22307ae5b8a6f8e2f8a250b0bf06b849fa2f665`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| 局部滑鼠 X 上界少 1 | window 第 1 組 |
| RNG／trust 的 ≤ 改 < | trust-rng 第 145 組 |
| 非零金額下界改 499 | amount 第 289 組 |
| 回應 2 仍扣款 | amount 第 298 組 |
| 數值取消旗標不經 POPF 還原 | amount 第 301 組 |
| 選單 CF 忙碌時不繼續 | selector 第 4 組 |
| 中途 code word 改 9091h | helper 第 1 組 |
| 回應 3 的 trust 參數改 29 | reply 第 151 組 |
| 說話型多減一次 3 | helper 第 4 組 |
| 恢復滑鼠 X 位移改成減法 | window 第 631 組 |

## 5. 範圍回填與重跑

較早 [spec/219](../spec/219-c-event-handlers.md)、[re/100](100-c-event-handlers-restoration.md)
已加後續視窗 backlink，由本輪 verifier 機器核對。舊收據的 modal fixture 保留，
新 C 收據只提升 caller／控制流，不把底層繪圖與正常玩家 UI 提升為完成。
另逐份核對以下同版位址的舊規格，它們的 Go 行為與原版未解邊界不受本輪影響：

- [189](../spec/189-lord-may-ignore-the-strategists-answer.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [50](../spec/50-corps-upkeep-charges-funds.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [145](../spec/145-general-and-faction-cells.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [103](../spec/103-phone-diplomacy-amount-keypad.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [38](../spec/38-list-windows.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [42](../spec/42-event-scene-speakers.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [142](../spec/142-personnel-dismiss-flow.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [126](../spec/126-command-popup-menus.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [45](../spec/45-advise-scene-layout.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [78](../spec/78-amount-input-editor.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。
- [13](../spec/13-main-window-toggles.md)：不受影響，已有玩家／呈現／規則證據沿用；本輪沒有 production 變更。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_modal_probe.py workplace/matching-decompilation/c-modal/ida KI.EXE
tools/c_recovery_modal.sh
```

IDA 9.4、Go 1.26.7、GCC 12.2.0，原版／專案／dosgolem 唯讀，非 root、無網路、限資源 Docker。
IDA image `sha256:4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab`；
Go image `sha256:e8c859f5632dcfde7b32d2012b4351728f6437930887c2f6a91ea242459e5514`。
dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，import 來源前後雜湊相同。
RNG 表固定 12:34:56，每案例 raw c/s 執行前寫兩側，沒有重擲或挑 seed。
未定義 AF 沿用 dosgolem 既有模型，IF/TF=0，分離 SS／資料段；不外推實機逐週期一致。

| 來源 | SHA-256 |
|---|---|
| modal.c | `37105677c93746e880ab89c7940eb02a1b318bfe18e76615b12870f17294da8e` |
| modal.h | `db3f1ea75cfd9d3026382fc1e803364b67c594c9874103f2de98090c3249f106` |
| modal_fixture.h | `427b17352aecd5e4d84522b1620cd41464cc3d571f52e77b5b565c6e9e71d072` |
| Go 對拍工具 | `89a20cf5b7ecd455874e3347a8522911e802b8b49b37f4b0a6e95ac2a708b8a4` |

本機產物在 `workplace/matching-decompilation/c-modal/`，驗證入口是
[`c_recovery_modal_verify.py`](../../tools/c_recovery_modal_verify.py)，現行 C 台帳
[`c-recovery-status.json`](c-recovery-status.json) 共 101 個函式。正式 Go 本輪沒有修改。

## 6. 未解範圍

| 項目 | 邊界 |
|---|---|
| 完整底層繪圖／輸入／音效及自然玩家視窗 | primitive fixture 不取代正常玩家或像素證據，Issue #22 保持 OPEN |
| 完整 Go 視窗及每時三方 | 尚未驗證 |
| 硬體 wall-clock／自然輸入時序 | 尚未驗證，CF 控制流不外推實機時間 |
| C 機器碼與原作者工具鏈 | 尚未驗證，整檔 binary match 仍是組語基準 |
