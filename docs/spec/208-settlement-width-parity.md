# 208：玩家收入與募兵累計的原始整數寬度

**狀態：CONFORMED。原版 word carry 與募兵 wrap 已對齊，局部矩陣、四劇本接線及狀態層冷測通過。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 工具與位址：IDA Pro 9.4 database linear，`sub_1548F` 的 `0x154B8`–`0x154C1`、`sub_15547` 的 `0x1559A`–`0x155A3`
- 動態證據：`workplace/matching-decompilation/c-settlement/results/go-probe-tax.json`、`go-probe-recruit.json`

## 1. 已證實的差異

玩家稅率計算先處理收入高 16 位，再處理低 byte。最後只對 SS:[BP] 的 word 做 ADD，
沒有對 SS:[BP+2] 做 ADC。低 word 的 carry 因此不進高 byte。
gross=66303、tax=99 時，原版與 C 都得到 103，舊 Go 普通乘除得到 65639。

募兵對 SS:[BP+4]／[BP+6]／[BP+8] 各做 word ADD，累計每次以 16-bit 繞回，
之後才套玩家募兵上限與預備兵上限。192 個北方據點，各 production=65535、除數 2，
原版累計得到騎馬 51584、弓兵 5952、步兵 7808；舊 Go 無限整數累計再鉗制得到 65500、5952、65500。
原版與 C 的完整 ABI、記憶體及 callee 快照在同一輸入相同，差異發生在 Go。

## 2. 審查後的 Go 契約

玩家收入保留原始兩段乘除與 low-word ADD，稅率限 0–100、gross 是合法 24-bit 累計。
三兵種的累計在每個據點加完後遮為 16-bit，再走既有募兵上限與預備兵加法。
AI 收入除 2、AI 的錯步進與節制閘、赤字次序及 RNG 次數保持原版契約。
這是原版對齊，沒有在 remake 修掉原版整數溢位。

## 3. 驗證

單元測試涵蓋稅率 0／99／100、進位前後、192 據點累計與三兵種結果。
原版矩陣另驗所有稅率、收入相鄰值、每個 production、區域邊界、玩家／AI 分支與四劇本月結。
結果與來源身分見 [`re/95`](../re/95-c-city-settlement-restoration.md) 與本機 verification.json。

## 4. 未解範圍

| 項目 | 邊界 |
|---|---|
| 非法稅率、除法 fault 與異常記錄 | 未驗證，不用模數運算猜補原版例外路徑 |
| 正常玩家長程月結與完整世界更新 | 未驗證，保留尾端 callee 的明示 fixture 範圍 |
