# 245：C 非區域退出與 DAC 淡出

**狀態：CONFORMED。O0／O2各2,520完整裝置與八個錯版通過。**

- 日期：2026-10-09
- 出處：[re/125](../re/125-c-nonlocal-exit-restoration.md)。
- KI.EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。

## 契約與閘門

- 還原 11CB1／10A1C／1EBDC，完整執行堆疊恢復、原 mouse farcall、17 級淡出、DAC byte 算式與原 RET。
- 以原 CALL10067 到 11BF6 準備合法外層 frame，深層複製 VGA／RAM 等狀態；舊 C helper 與 header 保持唯一實作。
- 原生 C runner 使用活著的局部逃逸脈絡，阻止舊 caller 在非區域 RET 後繼續 POP；不以 fake RET 或 guest CPU 取代 C。
- 同時驗原 13DC9／13830／13BA9／13B5A 的正常與非區域分支，含借位、恰好到零、CX=FFFF／非 FFFF、重複理由及撤回。
- 原版先經獨立模型，再比 O0／O2 完整 RAM、四 plane、DAC、暫存器／FLAGS、SS frame、IN／OUT 與裝置軌跡。
- DAC helper 驗低 nibble、byte 回繞、捨入與三通道順序；淡出驗 272 次實際 writer 與最終前 16 色黑色。
- 八個實編譯錯版必須被拒絕；來源、原檔 bytes、MZ 與獨立重組綁定，C 機器碼匹配仍標 false。
- 正式 Go 不改。原版初始化只作明示輸入準備，不據此宣稱 11BE0 已由原生 C 完整還原。

## 實作與收據

[main.c](../../tools/c_recovery/main.c)保留原始 SS／SP、滑鼠、淡出與 RET；原生 runner 只負責 C ABI 跨 frame 轉移。
[獨立模型](../../tools/c_recovery_main_data.go)先驗原版，完整 C 比較與來源綁定見[收據](../re/c-main-verification.json)、[原檔指令與 MZ](../re/c-main-code.json)、[Go 冷測](../re/c-main-go-verification.json)。
完整主迴圈只用作原版輸入準備，未計入三函式 C 完成範圍。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整主迴圈、其他退出與自然玩家路徑 | 本輪閉包不代證 |
| 原 TSR、實機時間及原 C 機器碼 | 固定模擬器平台與 C 行為比較不代證 |

## 初始化前段的後續驗證

初始化前段已由原生 C 驗證。五函式與無RET前段由獨立raw輸入產生完整狀態，再接本頁的退出閉包，見 [spec/246](246-c-world-bootstrap.md)。完整主迴圈仍在未解範圍。

## 右鍵分派的後續原生 C 驗證

右鍵分派已由原生 C 驗證。`sub_159B7` 的原始 `0x159B7` 已接原DS間接表與三個清除動作，保留落入`nullsub_1`的單一RET；全部26索引、live表與DS≠CS反例見 [spec/247](247-c-world-interaction.md)。完整主迴圈及其更新排程仍未完成。
