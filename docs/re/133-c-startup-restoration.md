# 133：C 開機、主排程與命名 UI 還原

**狀態：CONFORMED。42 函式＋3 raw 片段，217 案 O0／O2 全同，13 突變全拒。**

- 日期：2026-10-11／12。
- 規格：[spec/253](../spec/253-c-startup-main-loop.md)。
- 輸入：DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 固定IDA probe SHA-256：`3d3761a281fa2ef1b50d1e8088e29053335ec7317473448e5624681e87cd8ea7`；DB SHA-256：`3c02d203809b36b1996e5cfbf636b1301f8c6daf061ac9955b0f33bcae93ab75`。
- 工具與位址：IDA Pro 9.4（`ida-pro-9.4-idapython:py312-v1`），database linear段基址10000。
- 公開收據：[c-startup-verification.json](c-startup-verification.json)（`wolong-c-startup-verification-v1`，status `semantic-conformed`）；oracle golem `a9714ebdab2ad6b529f81225472680f2b11f2842`。

## 範圍

42個原始命名函式、1021條指令，另加3個raw片段（離場尾段`0x10043` 11指令、共用RETN位元組`0x1006A`、INT61 tick回呼`0x10356` 8指令）。閉包由14個開機根經未還原直接呼叫展開，28個額外函式全在命名UI群與視訊群內；跳表`cs:8F97h` 8筆全在閉包內。

## 已證實的定位

開機選單鏈（`0x10000→1005B→11A6E→11AC3`，`11AC3`是新遊戲／讀檔開關）：

- `sub_18DC8`是二選框：點擊框為`[DX+28h, +80h)×[BX+28h, +48h)`，`y_rel>>3`為0–1回`AX=0`（上選項）、3以上回`AX=1`（下選項），列2（`y_rel` 16–23）是死區重取；右鍵（CF）直接重取。`11AC3`傳`DX=D0h／BX=80h`，框原點為（248,168）。
- `sub_18B7C`是槽位選單：`BP=AX(entry)`；碼`<0x20`重取；`BP=2`跳過槽檢查；否則槽內字`≠D0A1h`才離開，`CF=1`（右鍵）帶carry返回。`18BD1`用`1E3D7`漆四列槽碼`0x20–0x23`（`BX=90h／C0h／F0h／120h`，`DX=100h`，`CX=20Fh`）。
- `sub_18E5A`是同形小選單：`18F37`漆`0x20`（`BX=110h`）與`0x21`（`BX=F8h`），`DX=148h／CX=206h`（x格41–46，x 328–375；早先誤讀`BX`高位元組，`(400,252)`落空，實測修正為`(350,256)`）。
- `sub_18FC9`是確認選單：碼`==0x20`置`[BX+2]=7Fh`後返回（`19195h`漆`DX=1D0h／BX=C0h／CX=208h`，確認點`(495,200)`）；`>0x20`走`cs:funcs_19037`調整處理器後重取。
- `sub_18412`是處理器選單（劇本2設定流）：碼`<0x3D`重取，`≥0x3D`分派`cs:funcs_18450`（5筆）或`1857F`後重取；唯一離開是右鍵（CF，往回取消）或處理器前進。劇本2的`(160,140)`（碼`0x3E`）經處理器1前進。
- `sub_1E3D7`塗漆ABI：`AL`=碼、`DX`=x像素、`BL`=y像素、`CX=(CH`列數`,CL`欄數`)`，列步幅80格。`sub_1E453`查80×50格圖（`x/8+(y/8)*80`），純滑鼠位置，無鍵盤分支。

九案boot導覽：`0x10000`四案（劇本列0／1／2＋`pspCmd=12`讀檔取消後走上選項）、`0x11AC3`三案（`BX=0／64／128`，框內`(310,175)`起手）、`0x11A6E`二案。`0x1005B`直接進入已退役：fixture堆疊／arena／seed-fill與真實`0x10040`呼叫點差約300位元組（含PRNG流），改以劇本多樣性取代，`1005B`本體經`0x10000`覆蓋。

## 已證實的輸入與邊界

- `sub_1222B` raw等待（`sub_18810`訊息等待，不經`0x21E7／0x84DD`檢查點）：免補丁直入第二半，只查右鍵（`AX=5／BX=1`經`0x1110`驅動直通INT33 `AX=5`，讀後清零）。`11CD0`的`d2a=1`路徑連發兩則訊息（`CX=4Bh／197h`），兩次進入需兩次右鍵；harness以`pollFeed`把`0x222B`入口做成transfer群專用檢查點餵食。`sub_12216`的補丁路徑（`EB 0F→9090`）在矩陣內無呼叫端。
- `sub_11D8E`是tick等待：`[0CFAh]==0`返回，否則自旋等`D2C≠0`再等`D2D≥[0CFAh]`，通過後雙雙歸零。矩陣無IRQ0，故`mainloop`的`a5=1／d2d≥A0h`（先歸零`D2C`再進`11D8E`）改用`d2d∈{2,3,9Eh,9Fh}`；`transfer`只測`d2d=FFh`就緒路。`a5=1／d2d<A0h`早返、`a5=2／3`遞減返回皆可測。
- `sub_11CD0`的`0x11D43（call 11CB1）`落下`0x11D46`：已知函式入口的跨程序邊，非截斷；verify的fallthrough檢查已放寬為`rows ∪ old_entries`（與call／jmp目標檢查一致）。

## 驗證與實作界線

- 217案／opt：13群（boot 9、bootchain 27、dispatcher 3、entry 3、exit 6、load 14、mainloop 22、naming-integ 10、naming-standalone 91、prompt 10、tick 8、transfer 4、video 10）；stage（exit 7、pause 10、return 196、transfer 4）；獨立稽核818。
- 13突變全以exit 1拒絕。重設計三筆（值變不改控制流，否則C遊蕩觸界panic）：m1 `0x1002D`離場碼`AL 0→1`、m2 `0x10024`訊息位址`DX+1`、m12 `0x10051`視訊模式`3→2`。
- 工具鏈修正：`-O0`補`-ffunction-sections／-fdata-sections`＋`--gc-sections`（未用static函式在-O0保留造成未定義引用）；稽核流分叉改記入mismatch走exit 1（保留「有分叉但無狀態差」的非確定性panic護欄）；報表補`new_routine_sha256`鍵。
- C不執行原版CPU指令，不借原版執行後狀態；兩側各自raw初始化。

## 未解範圍

| 項目 | 狀態 |
|---|---|
| 完整讀檔長鏈、正常玩家長程 | 本輪局部契約不代證。 |
| 未註冊熱區碼的自然可達性 | 只有位元組證據，無執行收據。 |
| INT61聲音驅動的真實硬體行為 | 以dosgolem聲音模型為界，不外推實機。 |
| 需IRQ0的等待路徑 | `11D8E`的tick分支與`1222B`的`D2D`飽和分支只測就緒／邊界，不在矩陣內觸發tick。 |
| `sub_12216`補丁路徑 | 矩陣內無呼叫端，未覆蓋。 |
