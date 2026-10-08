# 214：政治候選的原始排序順序

**狀態：CONFORMED。原始交換排序、同值尾端、raw byte 邊界與完整 36 三方向量通過。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA 9.4 database linear `sub_12C52`，`0x12C8A`–`0x12CDB`
- bytes SHA-256：`7d7a2b71bcb161ef0e00db6ca6012ee769a704b8e7dd92cc16d89e1efacba8ba`
- 來源：[`re/97`](../re/97-c-monthly-politics-restoration.md)；三方月結 audit 的 `s1-p21-seed0`

## 1. 已證實契約

原版按據點編號、四個鄰接槽的遇見次序去重。勢力編號不保證與該次序相同。
每個位置從自己開始掃，DL=FF，只在 relation byte 嚴格更小時立即交換到該位置。
同值不作新交換，先前交換也可能改變剩餘同值候選的相對順序。

Go 移除勢力編號 tie-break，逐步照原版交換。不能用 stable sort 只保留第一個同值結果，
因為尾端政治分支也會讀完整候選順序。輸入 Candidate 的 Friendship 是建立時快照，不讀漂移後 live 值。

## 2. 驗證

三個候選 `[4:B0,2:A9,1:A9]` 原版結果為 `[2,1,4]`。
三個 `[1:60,2:60,3:50]` 原版交換後為 `[3,2,1]`，尾端同值也要一致。
追加相鄰 raw byte 與 FF 同值，冷測 strategyai／state；三方完整月結對照繼續保留 queue／RNG 差異。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 其餘政治 producer、queue 與 RNG | 未驗證，排序修正不能宣稱整段月結 parity |
| 正常玩家 UI 與長程流程 | 未驗證，本輪只修已證實候選順序 |
