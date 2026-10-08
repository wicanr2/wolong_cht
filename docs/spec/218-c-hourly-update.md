# 218：C 每時更新、維持費與事件分派核心

**狀態：CONFORMED。八函式 O0/O2 各 21,871 組原版/C 相同，七個負對照皆被拒絕。**

- 日期：2026-10-08
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 位址空間：IDA Pro 9.4 database linear；檔案偏移另列
- 出處：[`re/08`](../re/08-hourly-update.md)、[`re/15`](../re/15-event10-producer.md)、[`re/98`](../re/98-c-go-monthly-comparison.md)

## 1. 原始契約

| 原始函式 | IDA 起點／bytes | 規則 |
|---|---|---|
| `sub_13E11` | 0x13E11／84 | event dispatch → 財政／LowFunds → 維持費 → 外交官 → 下一勢力 cursor → redraw |
| `sub_13E65` | 0x13E65／41 | 三兵種相加 carry 進 DL，五次 SHR／RCR 做 24-bit 除 32，真實支出 helper |
| `sub_13E8E` | 0x13E8E／111 | 原始外交官 byte sentinel、兩個 RNG 閘、官員經費、雙向 relation 更新 |
| `sub_15673` | 0x15673／34 | 24-bit 支出累加、只有上界鉗制與 AX／DX 保存 |
| `sub_131AE` | 0x131AE／68 | cadence byte 遞減、queue cursor <0x100、四 byte record、code.low 的原始間接表 |
| `sub_13496` | 0x13496／16 | event 10 的原始 formatter word、CX param 及 SS stack |
| `sub_13507` | 0x13507／19 | event 13 通知後 AL=50、CX=Param 的信賴度處理 |
| `sub_13DC9` | 0x13DC9／72 | 信賴度 byte SUB／borrow 歸零、原始通知／退出呼叫次序與保存暫存器 |

raw code.low 僅取 0–13，合法 cursor、段、勢力／官員索引，SS 分離、IF/TF=0。
維持費與外交官在不在場勢力的原版每時入口仍會被呼叫，原始 source 不加額外 Alive gate。
事件表的其他 11 個 handler 使用明示 RET fixture；不能因此聲稱政治／戰術 handler 完成。

## 2. 驗證閘門

O0/O2 比全暫存器／FLAGS、世界、globals、queue、RNG、cadence 與 stack。
每 callee 入口記 raw 參數、關係／官員記錄、SS stack。涵蓋 cadence／空事件／13 碼分派、
勢力 cursor、財政高位有號邊界、24-bit 維持費及 clamp、外交官 budget／relation／RNG、
event 10／13 的控制流及四劇本連續每時規則。
刻意錯的次序、cadence、cursor、carry、支出上限、外交 byte 閘與 trust borrow 需拒絕。

七個次序、cadence、cursor、carry、支出上限、外交 byte 經費及 trust borrow 突變皆被拒絕；
完整矩陣、callee 入口快照與 171 次全記憶體核對見 [`re/99`](../re/99-c-hourly-update-restoration.md)。

## 3. 未解範圍

| 項目 | 邊界 |
|---|---|
| 其他 event handler 與 UI／音效／退出 | 未驗證，fixture 只量本輪 caller 與原始參數 |
| 正常玩家日期流與完整原版硬體等待 | 未驗證，局部每時入口不替代玩家流程 |
| 完整 Go 每時規則 parity | 未驗證，另按原版／C 範圍建立三方 audit |
| C 機器碼匹配 | 未驗證，原版工具鏈仍未知 |
