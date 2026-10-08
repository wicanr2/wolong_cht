# 102：C 數值輸入器與財政 caller

**狀態：十七函式 O0/O2 各 788,931 組原版/C、787,839 組 Go 數值相同，十個負對照拒絕；限原始 ABI 與明示 primitive。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- IDA Pro 9.4 database linear；檔案偏移另列
- DB：`workplace/matching-decompilation/c-numeric/ida/input.exe.i64`
- DB SHA-256：`8a4d51f1057f192aba3f3e24369737121d9155cf28ed2329b3805e67a5e4f59e`
- 原始 probe：[`ida_numeric_probe.py`](../../tools/ida_numeric_probe.py)
- 契約：[`spec/221`](../spec/221-c-numeric-editor.md)

## 1. 原始定位與分級語意

十七函式共 598 bytes。原始名稱、chunks、operand、xref、file bytes 與 IDA bytes 保留，
語意只附加分級。以下 raw 控制流為已證實，限本輪 corpus 與 primitive contract。

| 原始函式 | IDA 起點 | 檔案起點 | bytes | MZ 重定位 | 已證實的 raw 語意 |
|---|---|---|---:|---:|---|
| `sub_17C6E` | `0x17C6E` | `0x7E6E` | 147 | 0 | 原始數值主迴圈、raw key table、完成／取消、SS frame 與 LAHF／SAHF |
| `sub_17D0D` | `0x17D0D` | `0x7F0D` | 58 | 0 | 原始數值外框、word 座標運算與圖形 callee 參數 |
| `sub_17D47` | `0x17D47` | `0x7F47` | 24 | 0 | 原始數值外框恢復參數與保存 |
| `sub_17D5F` | `0x17D5F` | `0x7F5F` | 52 | 0 | 十八格 raw byte、CLD／LODSB、三列六欄位置與熱區登記 callee |
| `sub_17DA5` | `0x17DA5` | `0x7FA5` | 30 | 0 | word 乘十加 digit、carry／飽和、cap 及 CF |
| `sub_17DC3` | `0x17DC3` | `0x7FC3` | 26 | 0 | word 乘百、高 word／飽和、cap 及 CF |
| `sub_17DDD` | `0x17DDD` | `0x7FDD` | 13 | 0 | 原始 word 除十、DIV FLAGS 模型與 CF |
| `sub_17DEA` | `0x17DEA` | `0x7FEA` | 2 | 0 | 原始完成動作只設 CF |
| `sub_17DEC` | `0x17DEC` | `0x7FEC` | 5 | 0 | 原始最大動作讀 SS cap 與 CF |
| `sub_17DF1` | `0x17DF1` | `0x7FF1` | 4 | 0 | 原始清零動作與 CF |
| `sub_101B4` | `0x101B4` | `0x3B4` | 33 | 1 | 原始 FAR 狀態／pointer 保存及 PUSHF／POPF／register |
| `sub_101DB` | `0x101DB` | `0x3DB` | 51 | 2 | 原始有號狀態三分支、FAR 恢復及 FLAGS／register |
| `sub_193E9` | `0x193E9` | `0x95E9` | 32 | 0 | 原始 DL／DH operand byte 交換、popup wrapper 與恢復 |
| `sub_167CD` | `0x167CD` | `0x69CD` | 25 | 0 | 原始財政上限 100、取消 gate 與 CS byte 寫回 |
| `sub_167E6` | `0x167E6` | `0x69E6` | 32 | 0 | 原始財政騎兵上限 10000、取消／除十與 CS word 寫回 |
| `sub_16806` | `0x16806` | `0x6A06` | 32 | 0 | 原始財政弓兵上限 10000、取消／除十與 CS word 寫回 |
| `sub_16826` | `0x16826` | `0x6A26` | 32 | 0 | 原始財政步兵上限 10000、取消／除十與 CS word 寫回 |

原始 MZ 表有 116 筆，本輪二個裝置 helper 含其中三筆。file segment word 0x1000，
IDA loader 加 0x1000 為 0x2000；runtime load paragraph 0x0110，遠 callee 為 0x1110:0000。
兩種位址基準分開記錄、逐 word 核對，不遮罩差異。

## 2. 原始 table、frame 與數值

間接 table 位於 IDA 0x17D01（CS:7D01）；六個 word 為 7DA5／7DC3／7DDD／7DF1／7DEC／7DEA。
十八格資料位於 IDA 0x17D93（CS:7D93），raw bytes：
`59 5A 5B 5D 5E 5E 56 57 58 52 5F 5F 53 54 55 5C 60 60`。
C 真正讀 table 與 LODSB，不用預列 handler 結果或熱區 layout 取代原始資料。

乘十加數字用 DX high-word 與 ADD carry 判斷超出 word，先飽和 FFFF 再按 cap 鉗制；
乘百同樣保留 high-word。退位以 word DIV 除十，未定義 FLAGS 沿用 dosgolem 的 high-half
比較模型，不宣稱實機的未定義旗標。數值主迴圈保留原始 LAHF／SAHF、SS frame、CLD、
DS／ES／register 保存、完成和右鍵取消。完成 helper 設 STC，主迴圈最終回 CLC；
0x5E 清零鍵不等同真正的取消。

四個財政 caller 由 C 真正接到數值輸入器，取消不寫 CS globals。上限 100／10000、
後三列的除十與寫回位置照原始 bytes；完整 UI 與 Go transaction 不由 scalar core 代證。

## 3. 原版／C／Go 矩陣

| 群組 | O0/O2 各原版/C | Go 有效數值 |
|---|---:|---:|
| digit-full | 655,360 | 655,360 |
| hundred-full | 65,536 | 65,536 |
| delete-full | 65,536 | 65,536 |
| action-boundary | 1,620 | 975 |
| glyph | 75 | 0 |
| device | 36 | 0 |
| editor | 432 | 432 |
| finance | 192 | 0 |
| popup | 48 | 0 |
| scenario-editor | 96 | 0 |
| 合計 | **788,931** | **787,839** |

- digit-full：全部 65536 個 SI × 十 digit，cap=FFFF。
- hundred-full／delete-full：各全部 u16 SI。
- action-boundary：九 cap × 十二 SI × 六動作，digit 含十值。
- glyph：三 helper × 五個 X × 五個 Y，驗十八格／保存與原始 word 邊界。
- device：二 helper × 六個有號狀態 × 三種 pointer，驗 FAR 及保存。
- editor：九 cap × 十六輸入計畫 × 三個裝置狀態。
- finance：四 caller × 十六輸入計畫 × 三裝置狀態。
- popup：四回應 × 四 CF 忙碌 × 三裝置狀態。
- scenario-editor：四劇本 × 三玩家 × 二 modal × 四真實數值輸入計畫。

每版 6164 次完整 1 MB 抽查相同；每案例比十四 register／FLAGS、world、queue、globals、
RNG、stack、cadence 與觀測欄位。每 callee 入口記 276 bytes，包含原始 ABI、64-byte UI 與
24-byte 輸入／裝置保存／popup operand／cap。十七函式都有實際 `entries_seen`。

Go 使用正式 `state.EditAmountValue`，有效條件在執行前固定為 0≤current≤cap≤65535。
只比六動作／主迴圈的結果值；Go action-valid bool 不當作原版 CF。raw SI 超過 cap 的
C 案例仍完整驗 ABI，不把不可比的 Go 正規化當作產品差異。本輪沒有正式 Go 變更。

O0/O2 JSON 相同，SHA-256：`7df4ef6af49a087461c427abad12a5e922b86bf1161f292d6f56e36d3103c712`。
兩側分別讀各自記憶體，狀態摘要同為 `e6a2f13f423396eaf9f9fc1c70617b7f959994238c153cd804ea0aec546c3f6b`。

| 刻意改錯的 C | 首次拒絕 |
|---|---:|
| digit carry 丟棄 | digit-full 第 65537 組 |
| 00 改乘十 | hundred-full 第 2 組 |
| 退位改除百 | delete-full 第 11 組 |
| 最大鍵讀零 | action-boundary 第 194 組 |
| 裝置 helper 不恢復 FLAGS | device 第 1 組 |
| glyph 水平步長改八 | glyph 第 51 組 |
| 跳表索引偏一 | editor 第 1 組 |
| 主迴圈不做 SAHF | editor 第 1 組 |
| 取消仍寫財政 globals | finance 第 4 組 |
| 徵兵值不除十 | finance 第 49 組 |

## 4. 明示輸入與 primitive

原版 RAM 與 C RAM 同時寫入事先固定的 raw key／cancel pairs。輸入查詢 IDA 0x121E7 每次
消費一筆、保存 last key；0x1E453 讀取該 raw key。真正的數值主迴圈處理它們。
popup leaf IDA 0x19409 提供固定回應／CF 忙碌，數值取消後的下一回應固定為 2。
裝置 FAR leaf 返回固定有號狀態／pointer／座標並 clobber FLAGS，原始 helper 必須自行保存。

每份輸入計畫後追加明示尾端取消，讓錯版忽略完成鍵時仍只讀已宣告事件。
第一次跳表負對照讀到第 513 筆，誤入原始 code bytes 才返回；舊收據／driver 保留於
`pre-terminal-cancel` 檔案。最終負對照由第二筆取消返回，verifier 機器核對 index=2／cancel=1，
不能用偶然 bytes 作拒絕的終止條件。

圖形／文字／其他裝置／世界 primitive 位址（IDA linear）沿用 re/101 的 RET 邊界，
另含 0x19796、0x18853、0x1062F、0x1895D、0x1E3D7。0x193E9、0x101B4／0x101DB、
0x17C6E 與四財政 caller 均有真正 C；底層畫面與自然玩家 input 沒有冒稱完成。

## 5. 回填與重跑

較早 [spec/220](../spec/220-c-modal-control.md)、[re/101](101-c-modal-restoration.md)
與 [spec/78](../spec/78-amount-input-editor.md) 已加同版位址範圍 backlink，verifier 自動核對。
0x17D5F caller 的十八格／位置／LODSB／參數閉合；0x1E3D7 的熱區 writer 新證據見 re/103，實際 VGA blit 仍未驗證。
現行玩家政略視窗 Issue #22 保持 OPEN。

```sh
WOLONG_IDA_PY_IMAGE=ida-pro-9.4-idapython:py312-v1 tools/ida.sh probe dosv tools/ida_numeric_probe.py workplace/matching-decompilation/c-numeric/ida KI.EXE
tools/c_recovery_numeric.sh
```

Go 1.26.7、GCC 12.2.0、IDA 9.4；非 root、無網路、限資源 Docker，原版／專案／dosgolem／
既有 x/text module cache 唯讀。dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`，import 來源前後 SHA-256 相同。
RNG 表固定 12:34:56、raw c/s 執行前寫兩側，沒有重擲或挑種子；IF/TF=0、分離 SS／資料段，
shift／logic／DIV 未定義 FLAGS 沿用固定 dosgolem 模型，不外推實機逐週期／自然輸入時間。

| 來源 | SHA-256 |
|---|---|
| numeric.c | `b68f2a4e4149c7cf5fff8aba2da0a6a437202c0626519ed7c6b7ccd73df00e5b` |
| numeric.h | `c0abfde294d876d1b7c49119f134bff56b0cdc9e76ac4ae5a23abe2cd6cea0f0` |
| numeric_fixture.h | `4780a1ff675bc17d77553035c0b2a2eeae790c607d1d613878a59d300267a2d5` |
| Go 對拍工具 | `dfed24a8f90943da352e3f5b229b883e8b6ac34b2654a83354030e3dd2c41aaf` |

本機收據在 `workplace/matching-decompilation/c-numeric/`，驗證入口是
[`c_recovery_numeric_verify.py`](../../tools/c_recovery_numeric_verify.py)，現行 C 分級台帳
[`c-recovery-status.json`](c-recovery-status.json) 共 118 個函式。Go 比較來源逐檔雜湊在 verification.json。

## 6. 未解範圍

| 項目 | 邊界 |
|---|---|
| 真實圖形／實際 input polling 與 popup leaf | primitive fixture 不取代 UI／像素／正常玩家證據 |
| 完整 Go 財政／視窗垂直鏈 | 局部數值核心不能外推完整玩家流程 |
| 硬體 wall-clock／自然輸入時間 | 未驗證 |
| C 機器碼與原作者工具鏈 | 未驗證，整檔 binary match 仍是組語基準 |

## 熱區 callee 勘誤（2026-10-08）

`sub_1E3D7` 是熱區 writer，由同一 segment／base 的 `sub_1E453` 回讀。先前 glyph
導覽名稱不成立。歷史群組名 `glyph` 與舊收據保留，只量 caller／RET primitive；
真正 map／pixel query 接線見 [`re/103`](103-c-hotspot-restoration.md)，不能把舊 ABI 通過當作 glyph 畫面證據。
