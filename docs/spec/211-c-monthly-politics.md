# 211：C 月結政治、俘虜與佇列初始化

**狀態：CONFORMED。二十一函式 O0/O2 各 11,632 組原版/C 相同，完整 C 月結規則接線通過。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 位址空間：IDA Pro 9.4 database linear；檔案偏移另列
- 出處：[`re/07`](../re/07-monthly-settlement.md)、[`re/15`](../re/15-event10-producer.md)、[`re/96`](../re/96-c-monthly-world-update-restoration.md)

## 1. 原始依賴與審查契約

| 原始函式 | 角色 |
|---|---|
| `sub_1585F`／`sub_15899`／`sub_15940` | 127 武將、月度 timer、入仕、俘虜逃脫與歸降 |
| `sub_12AD2`／`sub_15990` | 原始勢力武將數 byte 加減、玩家通知參數與堆疊 |
| `sub_1301C`／`sub_12FB1` | 256-slot raw writer、勢力 code.high 適配及既有 writer 接線 |
| `sub_12BD9` | 佇列前 64 格丟棄、後 192 格前移／後 64 格清零、cursor=0、cadence=7、scratch 清 FFFF、兩輪勢力更新 |
| `sub_12C52`／`sub_12CDF` | 真實地圖邊境清單、重複排除、24 個 scratch word 區塊與關係排序 |
| `sub_12D3A`／`sub_12D58` | RNG 遷都事件、玩家／AI 關係更新及政治順序／分支 |
| `sub_12DB8`／`sub_12DF3`／`sub_130F0`／`sub_1310A` | 原始雙向關係取址、低 7-bit 減值、war bit 保存及玩家／AI 漂移 |
| `sub_13091` | 三兵種預備兵右移、word 和、據點數及資金高位的原始力量指標 |
| `sub_12E33`／`sub_12E89`／`sub_12EFB`／`sub_12F71` | 協力、停戰、宣戰與中立政治分支、byte／word 邊界及真實事件 writer |

保留原始函式、offset、CS／DS／ES／SS、operand 與 FLAGS。合法勢力 0–21、首都與地圖鄰接索引、
武將 0–126、至少一個存在勢力、分離段及 IF/TF=0。
中立 24、FF sentinel 與排序 bit 按原始格式處理，不作合法勢力外推。
不為無限尋找或非法表格猜補終止行為。

## 2. 驗證閘門

O0/O2 比全部暫存器、FLAGS、世界、queue、scratch、globals、RNG 與堆疊。
每個 callee 入口快照保留表格、SS stack 與 source identity。
涵蓋 timer、存在旗標、任官 sentinel、逃脫／歸降門檻、writer 滿槽與 byte 判空、
關係 war bit／相鄰門檻、scratch 去重與排序、政治順序與四個原版劇本完整月結規則接線。
刻意錯的 timer、任官計數、queue 搬移／writer、關係 bit、力量、排序或政治門檻版本須被拒絕。

八個 timer、任官、queue 搬移／writer、關係 bit、力量、排序、協力門檻突變全部被拒絕。
原版/C 矩陣、四劇本接線與 source identity 見 [`re/97`](../re/97-c-monthly-politics-restoration.md)。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 通知、音效、重畫與正常玩家流程 | 未驗證，fixture 只驗原始參數與堆疊，沒有 UI 完成聲明 |
| 非法索引、無勢力、段重疊與資料錯誤 | 未驗證，不擴張正常記錄契約 |
| 完整 Go 月結規則同狀態比較 | 未驗證，本輪新政治依賴以原版/C 為主 |
| C 機器碼匹配 | 未驗證，原版工具鏈仍未知 |
