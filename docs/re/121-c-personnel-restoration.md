# 121：原人事選單與內政官、外交官任免的 C

**狀態：CONFORMED。七函式，O0／O2 各 2,332 完整裝置與十二錯版通過。**

- 日期：2026-10-09
- 範圍：松崗 DOS/V `KI.EXE`；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；位址均為 IDA database linear address，段基址 `0x10000`。
- IDA database SHA-256：`cfc39dbd3238559d30c92c900a08d44bb98c3fe27b17f6123958418c899a1485`。
- Probe SHA-256：`643ee2cee9da253cb371533d96db87783a78df26f7d6d2704d9bb8c8d5200252`。
- 既有證據：[re/25](25-message-variants-and-personnel.md)、[re/22](22-strategy-command-tree.md)。
- 規格：[spec/241](../spec/241-c-personnel.md)。

## 已證實契約

| 原始入口 | 原始定位與效果 |
|---|---|
| `sub_16265` | 四項原選單、TALK #78、位置 `DX=0403h`；取消返回，否則 AL×2 查 CS:6280。 |
| `sub_16A9B` | 玩家據點清單，再選無職非君主武將；寫武將+17=2與據點+19=索引。已有人顯示 #52，成功組19C。 |
| `sub_16B08` | 玩家據點清單，解任內政官；空缺顯示 #54，成功組1A2。 |
| `sub_16B71` | 存活他勢力清單，再選 eligible 武將；寫武將+17=3與勢力+2A=索引。已有人顯示 #53，成功組19D。 |
| `sub_16BE3` | 存活他勢力清單，解任外交官；空缺顯示 #55，成功組1A3。 |
| `sub_16B4F`、`sub_16C2A` | 無條件 XCHG 官員欄為 FF；原FF回CF=1，否則清原官員+1A經費、+17職務並回CF=0。 |

七函式共 185 指令／459 bytes，與固定原始 bytes 一致。
四支完整 caller 只保存入口 DS，退出清提示、重畫大地圖並還原游標。
第二層選將取消回第一層；成功與拒絕也回第一層，直到第一層取消。

## 原版邊界

- 任命不清武將經費；解任才清經費與職務。
- `0x16B2E` 空缺訊息後 POP BX，`0x16C09` 是 POP SI，保留原差異。
- `0x16279` 查原始可變 table，不新增 AL<4 檢查；helper 不檢查官員索引、職務種類或所屬勢力。
- 任命第二層共用 `sub_17663`，保留其既有排序狀態與資格條件。
- 受控輸入用真實選單、清單、TALK、mouse與VGA，正式 Go 引擎不修改。

## 驗證結果

- O0／O2 各 2,332 組：2,064 helper、120 任命、96 解任、52 真實人事選單案例。選單含16次原CS跳表受控修改。
- 原版先經獨立欄位與CF／BX模型 15268 項核對，再與C比完整1MiB RAM、四plane、暫存器／FLAGS、SS frame、呼叫、IN／OUT及DOS／font／mouse／sound。
- 官員索引0–127與FF、經費0／1／255；任命保留非零經費，`0x16B67`／`0x16C42` 解任清零。這是受控初值正對照，不代證自然撥款玩家路徑。
- 清單排序初值固定0，依原掃描順序選首筆；共用引擎的其他排序欄仍依re/119的獨立驗證範圍判讀。
- 十二個實際編譯錯版全拒絕，兩側缺字0，C來源乾淨重生一致。
- [版控收據](c-personnel-verification.json)、[指令覆蓋](c-personnel-code.json)與[正式Go冷測](c-personnel-go-verification.json)保存來源及驗證雜湊。重跑 `bash tools/c_recovery_personnel.sh`。
- 全185指令／459bytes已在既有組語基準，另以固定binutils獨立重組完全一致，本輪新增組語0條。
- 共用 `sub_17663` 的 `0x1768A` MOV CL,CS:98AA 已由 [re/119](119-c-list-families-restoration.md) 證實覆寫前序XOR；本輪沿用既有唯一C，不把它解釋為重設清單游標。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然玩家長程、非法分派與原作者 C 機器碼 | 不由局部 C 語意代證 |
