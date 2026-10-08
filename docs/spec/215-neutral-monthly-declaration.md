# 215：中立邊境的月度事件 1

**狀態：CONFORMED。FF18 raw producer、門檻、無邊境／既有目標與完整 36 三方向量通過。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA 9.4 database linear `sub_12F71`，`0x12F71`–`0x12FB1`，64 bytes
- bytes SHA-256：`100559ce2c31af8dc150d1b49387d86f1ba839f7e0f71c876b0ab0eb38fd28b1`
- 出處：[`re/97`](../re/97-c-monthly-politics-restoration.md)、三方 `s3-p0-seed0` raw queue

## 1. 已證實契約

sub_12EFB 不成立後才走中立分支。InvasionTarget <24 直接返回；有中立鄰接才繼續。
門檻 `min(Cities*16+96,1757)` 與資金高 16 位做有號比較，資金必須嚴格更大。
條件不成立就清目標 FF；已有中立目標 24 不重發，其餘合法無目標 FF 發事件 1、Param=FF18。
writer 與 RNG 次數依原版，不把中立 24 當成普通勢力或增加存活檢查。

Go 使用同一據點 Adjacency bit／Neighbours 記錄判中立邊境，在普通宣戰未發事件後處理。
正常分支保留已有普通侵攻目標的原始更新，不把中立補入外交候選或改候選排序。

## 2. 驗證

冷測有／無中立鄰接、邊界資金、原有中立目標與 raw queue／RNG 寫入。
36 向量中的原版 FF18 事件需由 Go 真實 writer 產生；所有剩餘 raw／RNG 差異繼續報告。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 其餘普通政治分支與 RNG 漂移 | 未驗證，中立補線不代表完整月結 parity |
| 正常玩家 UI、長程與存檔垂直鏈 | 未驗證，本輪驗局部原版 producer |
