# 244：C 政略指令與進言理由

**狀態：CONFORMED。29 函式，O0／O2 各 3,348 完整裝置與十六錯版通過。**

- 日期：2026-10-09
- 出處：[re/124](../re/124-c-strategy-command-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約與閘門

- 還原 29 個原入口，既有清單、鏡頭、數字、TALK、場景與隊列保留唯一 C 實作。
- 八項指令列、君主任務閘門、軍團／據點／武將／勢力入口與財政進出保留原 ABI。
- 進言 AL 結果、DL 理由、信賴度、已嘗試位元與剩餘理由數逐項對照原版。
- 保留 byte 溢位、有號 24-bit 金額、CF 保存、DS 恢復、自身尾跳與六筆 farcall。
- 原版由獨立模型先驗，再比完整 RAM／四 plane／暫存器／FLAGS／SS frame／裝置與 API 軌跡。
- 同一矩陣以 O0／O2 編譯執行，實編譯錯版必須被拒絕；記錄零缺字及來源、工具、素材雜湊。
- INT61/AH0A 狀態為明示平台 fixture；原 TSR 與硬體時間不在範圍內。
- 正式 Go 不改；原版長程玩家流程與原 C 機器碼仍待另驗。
- 進言拒絕驗恰好扣到零且無借位的返回路徑；扣減借位後 `sub_11CB1` 的啟動堆疊恢復另列範圍，不偽裝為局部函式返回。

## 實作與驗證

[strategy.c](../../tools/c_recovery/strategy.c)與產生來源保留原入口及運算元。
[獨立模型](../../tools/c_recovery_strategy_data.go)先驗原版；[版控收據](../re/c-strategy-verification.json)、[原檔指令與 MZ](../re/c-strategy-code.json)、[Go 冷測](../re/c-strategy-go-verification.json)保存來源及結果。
正式 Go 玩家路徑未改，自然長程不由局部矩陣代證。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 自然長程、原 TSR 音效與原 C 機器碼 | 本輪不代證 |
| ~~信賴度借位後非區域跳出~~ | 非區域返回已由三函式閉包驗證，見 [spec/245](245-c-nonlocal-exit.md) |

## 非區域退出的後續驗證

非區域返回已由三函式閉包驗證；完整原 outer frame 與 RET、DAC、13DC9／13830／13BA9／13B5A 借位分支見 [spec/245](245-c-nonlocal-exit.md)。原 29 函式正常返回矩陣未冒充這條新增驗證。
