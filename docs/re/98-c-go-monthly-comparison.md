# 98：完整月結規則的原版／C／Go 三方比較

**狀態：36 固定向量的完整原版區塊及最終 RNG 相同，限月結規則直接入口與明示 UI fixture。**

- 日期：2026-10-08
- 松崗 DOS/V KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- SINARIO.DAT SHA-256：`21acf8a8c4d406b4deb3a184ec0a95f3670d3e0bfff02df63d5d218f46f0754c`
- 原版入口：IDA Pro 9.4 database linear `sub_15358`，C 呼叫鏈與 DB 身分見 [`97`](97-c-monthly-politics-restoration.md)
- Go 入口：原有 `tick` 月結區段純抽出為 `monthlyRules`；研究適配只存在於 `matching` build tag
- 比較工具：[`monthly_compare.go`](../../tools/monthly_compare.go)、[`monthly_compare_verify.py`](../../tools/monthly_compare_verify.py)
- 契約：[`spec/212`](../spec/212-c-go-monthly-comparison.md)，修正契約為 spec/213–217

## 1. 相同輸入與比較範圍

四個原版劇本，每個取玩家 0／7／21、raw RNG c=0／77／255，s=uint8(c×37)，共 36 向量。
置換表固定由 12:34:56 產生，原版/C 及 Go 在執行前收到同一份 258 bytes；沒有重擲或篩選 seed。
兩側讀同一 22,208-byte 區塊，設定同一當月／下月財政參數與佇列。
Go 啟用正式 StrategicAI，關閉已有可關閉的近似 event 10；其他正式規則保留。

原版側保留完整月結規則的真實指令，僅通知、音效、重畫為 RET fixture。
C 每向量與原版比較全部暫存器及完整 1 MB 記憶體，再重建同一原版區塊。
Go 在同一月結規則邊界執行，沒有據點 tick、軍團 tick 或每時更新混入。
載入／Bytes 前置基準 36 向量均為零差異。

比較完整 22,208 bytes 及完整 258-byte RNG，沒有 comparison mask，也不省略 queue、globals、corps、
未知欄位或載入空隙。runtime cursor／delay 另記錄兩側設定與結果。
這個規則直接入口不替代正常玩家 UI／存檔／長程流程。

## 2. 原版證據定位的五個 Go 缺口

| 差異 | 原始證據／修正 |
|---|---|
| 所有存在武將 +0x1F 評分沒有寫回 | `sub_155A6`；新增 typed byte、載入／保存及月結更新，見 [`213`](../spec/213-general-monthly-score.md) |
| 同值候選用勢力編號重排 | `sub_12C52`；改原始交換式選擇排序，尾端相對順序也一致，見 [`214`](../spec/214-political-candidate-order.md) |
| 中立邊境 FF18 宣戰事件漏接 | `sub_12F71`；合法前提及有號資金高位、真實 writer，見 [`215`](../spec/215-neutral-monthly-declaration.md) |
| StormArea 已算出但四個 globals 沒保存 | `sub_122DB` 的 D22／D24／D26／D28；同步原版 word 位置，見 [`216`](../spec/216-monthly-storm-globals.md) |
| 玩家月結宣戰被額外 Player gate 丟棄 | `sub_12EFB` 沒有該 gate；去除兩層額外條件，見 [`217`](../spec/217-player-monthly-declaration.md) |

推論等級已證實，限上述原始位址、受控向量與局部規則。沒有把 UI 名稱或 C 變數名當成證據。
初次 36 向量全部有差異，score +0x1F 共 3429 byte 差異；這些 bytes 經修正全部歸零。

| 研究 checkpoint | 完整 raw block 與 RNG 相同 |
|---|---:|
| before-score，原始月結 | 0／36 |
| after-score | 9／36 |
| after-order | 12／36 |
| after-neutral | 20／36 |
| 評分／排序／中立／globals／玩家 producer 完成後 | **36／36** |

各 checkpoint 原始 audit、vectors 與 Go source hashes 保留在本機根目錄的同名子目錄。
不是只消除差異數字，最終 validator 逐向量讀完整原版與 Go 區塊及 RNG，重新算摘要並驗 equality。

## 3. 工具與拒絕閘門

```sh
tools/monthly_compare.sh
```

工具在非 root、無網路、限資源 Docker 執行，原版、C／Go／dosgolem 及鎖定模組唯讀。
Go 1.26.7、GCC 12.2.0，dosgolem revision `a9714ebdab2ad6b529f81225472680f2b11f2842`。
import 來源在前後雜湊相同，C／Go 完整 source hashes、每個 vector 的輸入／輸出 hash 記在 verification.json。
沒有向 Git 加入原版資料、向量區塊或 RNG binary。

驗證器 selftest 分別改 globals、勢力、據點、武將、queue 或 RNG，各種單 byte 差異都拒絕。
最終 audit SHA-256：`8c0d0f549be777db05b62a91d95cd98162cbd74accffcadb71de83367707f2d8`。
原版/C 若完整機器狀態不同，執行器立即停止，不產生可被當成通過的 Go 收據。
原有正常 state 冷測、matching tag 冷測、strategyai 邊界／排序與 Go vet 通過。

## 4. 未解範圍

| 項目 | 邊界 |
|---|---|
| 原版任意後期世界、俘虜／政治所有條件 | 未驗證，36 向量是完整矩陣的明示抽樣，不外推全部世界 |
| 正常玩家 UI／音效／長程與完整存檔垂直鏈 | 未驗證，rules 直接入口不取代正常玩家驗收 |
| C 全執行檔與機器碼匹配 | 未完成，C 仍為還原語意；整檔 byte match 仍由組語基準提供 |
