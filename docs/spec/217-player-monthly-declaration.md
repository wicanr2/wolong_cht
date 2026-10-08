# 217：玩家勢力也執行月度宣戰 producer

**狀態：CONFORMED。額外 Player gate 已移除，同規則玩家／AI 冷測與 36 三方向量通過。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- IDA 9.4 database linear：`sub_12D58`→`sub_12EFB`，`0x12EFB`–`0x12F71`
- bytes SHA-256：`0e02dae298d25b7894fdbcbd566a4217378a8f17ab9f7008ad682a3981fdda4a`
- 直接反例：`s1-p7-seed255` 原版 Code=0701、Param=FF08，Go 省略

## 1. 已證實契約

玩家／AI 的 relation 漂移入口不同，停戰與協力另有玩家條件，宣戰 producer 自身沒有。
去除 Go 的額外 self.Player 與 state i==Player 返回，保留原版既有目標、資金、關係、力量閘與 raw queue。
這是玩家仍由原版君主政治處理的對齊，不增加直接玩家命令或改 GUI。

## 2. 驗證

相同規則資料只切換 Player 標記，ShouldDeclareWar 結果需一致；正式 state 冷測與 36 月結向量。
queue 與 RNG 不能因仍有其他差異而遮罩或捨棄。正常玩家長程／完整 UI 仍未驗證。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 玩家宣戰後正常 UI 與長程流程 | 未驗證，本輪驗局部月結 producer 與完整 raw 向量 |
