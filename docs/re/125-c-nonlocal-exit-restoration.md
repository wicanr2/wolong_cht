# 125：信賴度退出、原堆疊返回與調色盤淡出的 C

**狀態：CONFORMED。三函式，O0／O2各2,520完整裝置與八個錯版通過。**

- 日期：2026-10-09
- 範圍：松崗 DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；下列位址為 IDA database linear，段基址 10000。
- Database SHA-256：`5166d55917adb9e528eabf7e4325b002455b802e348209daa9f1c6dc488e5740`。
- Probe SHA-256：`9ae600545d4fa8bdf99dcb96ba2987f44bc94b29e7aebf0b5ecc98ed0892c743`。
- 探針：[ida_main_probe.py](../../tools/ida_main_probe.py)；前輪界線：[re/124](124-c-strategy-command-restoration.md)；規格：[spec/245](../spec/245-c-nonlocal-exit.md)。

## 已證實原始契約

| 入口 | 契約 |
|---|---|
| `sub_11CB1`，24 bytes | `0x11CB1` 從 CS:9903 恢復 SS，`0x11CB6` 從 CS:9901 恢復 SP；保存 AX、真正呼叫滑鼠隱藏與淡出，再由 `0x11CC8` 的 RET 消費外層返回位址。 |
| `sub_10A1C`，73 bytes | 從 CS:1964 讀 16 色工作色盤，強度 16→0 共 17 級，每級色號 0→15，共 272 次 DAC writer；保存 AX／BX／CX／DX／SI。 |
| `sub_1EBDC`，82 bytes | BL 為強度；原 BH／DL／DH 三個 byte 只取低 nibble，先乘強度並加 128 捨入，再以第二／第三／第一 byte 的順序寫 3C9。 |

三函式共 89 指令／179 bytes；`0x11CBF` 的滑鼠 farcall 含一筆 MZ 重定位。
DAC 原始輸出公式是 `byte(4*floor(((component&15)*level+8)/16))`，保留 byte 回繞。
3DA 的等待讀取 bit3 兩組低／高轉換；本輪只驗固定模擬器的 I/O 契約，不推實機時間。

## 真正外層 frame 與原 caller

原 `0x10067` 的 CALL 進入 `sub_11BE0`，返回位址為 `0x1006A`。
`0x11BEC`／`0x11BF1` 保存當時 SS／SP；子 caller 之後使用較低或不同 segment 的堆疊。
`sub_11CB1` 的 RET 必須取回這個外層位址，不能繼續 `sub_13DC9` 的側欄更新與 POP。

- 原版準備階段執行 CALL 到 `0x11BF6`，包含 189F0／1533D；1533D 會畫日期，不能只複製 RAM。
- 前置輸入先提供已選玩家 0、數字字模 bank 與 CS:D44=5000；D850 是另一個指標。189F0 會按勢力重寫 decoded map 的城市圖塊，155A6 會重算武將評分，因此先驗磁碟來源 bytes，再記錄前置程式產生的 runtime 狀態。
- 兩側以 dosgolem 的深層 Snapshot／Restore 取得相同 RAM、VGA planes、latches、DAC 與埠狀態，各自持有獨立裝置。
- 測試段原版用 CPU，C 用原生編譯函式。C 的原 RET 完成後，C runner 以同一個仍有效的 `setjmp/longjmp` 脈絡離開舊 C 呼叫鏈；共用 header 與已驗證函式不改。
- 在 dosgolem 執行期 CS:IP=`0110:006A` 停止，尚未執行外層 RET；不代證完整主迴圈或正常玩家路徑。
- 原 `GAMEPAL.BRG` 共 384 bytes，SHA-256 `1f0119c75ea5cd333bd3ac75ef92030f93924f011728edc9b1f727c483263708`。工作色盤取 `[48:96]`；原 EXE 的 CS:1964 字串不是合法遊戲色盤。

`sub_13DC9` 借位時強制信賴度 0 與 CX=19E，通知後退出。
CX 原為 FFFF 且恰好扣到零時可正常返回；CX 非 FFFF 且通知後為零時也退出。
理由與進言場景須同時驗正常返回與非區域退出，不以飽和減法模型假設所有 caller 都返回。

## 右鍵表定位勘誤

已證實，`sub_159B7` 內 `0x159C0` 的原 bytes `FF 97 06 5A` 是 DS:[BX+5A06] 的間接 CALL。
`funcs_159C0` 的符號後綴不能當成表位址；真正右表為 IDA linear `0x15A06`，26 個 word。
`0x159CE` 的 MOV ES,AX 後落入 `0x159D0` 的 RET。
此處只校正定位，未把右鍵動作或整個主迴圈計入三函式 C 閉包。

## 驗證結果

- O0／O2 各 2,520 組：DAC helper 960、淡出 16、直接退出 96、信賴度 160、進言場景 248、理由 960、理由迴圈 80。
- 原版先經獨立模型；完整 RAM／四 plane／768-byte DAC／暫存器／FLAGS／SS frame／IN／OUT／裝置與非區域返回 metadata 一致，八個實編譯錯版均由狀態比較拒絕。
- 十二份原版前置快照涵蓋四劇本、三組外層 SS／SP；另驗不同的 nested SS。固定玩家 0 的地圖改動為 180／450／538／584 bytes，四組 SHA 與獨立 byte 模型相同。
- 正常進言的六個 locals 會在 ADD SP,6 後被 CALL19321／PUSH AX 覆寫；驗證其信賴度、CF 與決策軌跡。非區域進言的理由 frame 仍存留，則逐欄核對；不猜已退休 stack bytes 的語意。
- [版控收據](c-main-verification.json)、[原檔指令與 MZ](c-main-code.json)、[Go 冷測](c-main-go-verification.json)保存來源與範圍；重跑 `bash tools/c_recovery_main.sh`。
- 89 指令／179 bytes 全已有組語覆蓋，一筆 MZ 反向重定位與獨立重組通過，整檔 67,099-byte EXE 仍與原版一致。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整 11BE0 主迴圈與其他退出分支 | 本輪只驗原始 frame 準備及信賴度 caller 的退出閉包 |
| 原 TSR、實機淡出時間與原 C 機器碼 | 固定平台 I/O 與原生 C 比較不代證 |

## 初始化前綴的後續原生 C 驗證

初始化前段已由原生 C 驗證。`sub_189F0` 的原始 `0x189F0` 城市／軍團處理及評分日期，現由兩側各自的 raw 初態產生，有限前綴停在11BF6而無RET；來源見 [re/126](126-c-world-bootstrap-restoration.md)。本頁原Snapshot實驗仍保留當時範圍，完整主迴圈未因此完成。
