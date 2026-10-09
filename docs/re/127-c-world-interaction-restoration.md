# 127：世界點選、右鍵分派與淡入的 C

**狀態：CONFORMED。八函式，O0／O2各2,148組完整狀態與十二錯版通過。**

- 日期：2026-10-09
- 範圍：松崗DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；位址均為IDA database linear，段基址10000。
- Database SHA-256：`f8bc52847a56fa72da6fc0d140d07608f96967861208de94c27530358d94b85e`。
- Probe SHA-256：`328f9f6bd299452d717e92828b11e92de1c17ef49a66b4bca4c55246f9c0e3bf`。
- 探針：[ida_loop_probe.py](../../tools/ida_loop_probe.py)；規格：[spec/247](../spec/247-c-world-interaction.md)。

## 已證實原始契約

| 原始入口 | 範圍與行為 |
|---|---|
| `sub_109D0`，0x109D0 | 40指令／76bytes；淡入強度0..16，每級讀16色並等待兩次垂直回掃。 |
| `sub_11E46`，0x11E46 | 74指令／200bytes；地圖點選城市與同格軍團，八byte SS區域，面板與取消迴圈。 |
| `sub_11F0E`，0x11F0E | 13指令／34bytes；低byte座標限制與既有選單呼叫。 |
| `sub_159B7`，0x159B7 | 10指令／25bytes；AL經原DS表間接call，重畫後DS／ES=CS，落入0x159D0。 |
| `nullsub_1`，0x159D0 | 1指令／1byte；原RET，保留原名稱及獨立邊界。 |
| `sub_15AA2`，0x15AA2 | 7指令／20bytes；清98A6 bit2，原顯示參數AL=0、DX=27、BX=0、CX=0A0D。 |
| `sub_15E4C`，0x15E4C | 7指令／20bytes；清98A6 bit1，原顯示參數AL=0、DX=27、BX=10、CX=0D0D。 |
| `sub_161B6`，0x161B6 | 7指令／20bytes；清98A6 bit0，原顯示參數AL=0、DX=0、BX=0、CX=021B。 |

原右鍵表位於IDA 0x15A06，26個word只指向上述三個清除動作與`nullsub_1`。
表經DS讀取；函式邊界不包含最後的RET，不能把25bytes視為完整返回路徑。
八函式共159指令／396bytes；主迴圈0x11BE0仍為導航，更新排程0x11CD0及其依賴尚未還原。

## 驗證紀錄

[Go檢查收據](c-interaction-go-verification.json)保存本輪39套件冷測、工具版本與隔離自測。
原生結果與指令來源見[C收據](c-interaction-verification.json)及[指令覆蓋](c-interaction-code.json)。
- O0／O2各2,148組完整RAM／四plane／DAC／暫存器／FLAGS／SS／IN／OUT／API一致。fade12、clear108、right488、popup1,144、point396。
- 原版先驗獨立資料模型，再與各自raw初態的C比較。128組DS≠CS反例可分開錯誤CS讀表／旗標，26個右鍵索引及72組原表修改均涵蓋。
- popup每個低byte與高byte污染、城市CB..D3／CA／D4邊界、城市存在flags未檢查、兩項取消／選擇、多軍團重選與槽125／126邊界均已驗。
- 淡入每例核對272次原setter參數、1,088次DAC writes與完整DAC；十二實編譯錯版均以狀態差異拒絕。
- 點選的軍團面板分支使用外國軍團，完整執行顯示、關閉與重選；本輪不據此聲稱正常自軍長程或自然玩家流程。
- [C來源](../../tools/c_recovery/interaction.c)、[獨立模型](../../tools/c_recovery_interaction_data.go)、[收據](c-interaction-verification.json)與[指令覆蓋](c-interaction-code.json)可回查；重跑`bash tools/c_recovery_interaction.sh`。
- 159指令／396bytes均有既有組語覆蓋，獨立重組逐byte一致，無新增MZ重定位；完整67,099-byte EXE仍相同。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整主迴圈、更新排程與正常玩家長程 | 本輪閉包不代證 |
| 原硬體時間與原C機器碼 | 固定平台契約及局部比較不代證 |
