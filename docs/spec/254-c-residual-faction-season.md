# 254：C 殘部群（勢力一覽／季節／戰術／地圖／外交／側欄）

**狀態：CONFORMED。10 支 247 指令＋1 raw 尾跳，89 案 O0／O2 全同，10 突變全拒；census 739／739 全覆蓋。**

- remake 實作：`tools/c_recovery/residual.c／.h／_fixture.h／residual_generated.inc`，harness `tools/c_recovery_residual.go／_data.go／_platform.go／.sh／_container.sh／_verify.py／_generate.py／_mutants.tsv`。
- 驗證：`WOLONG_RESIDUAL_MODE=normal-o2／normal-o0` 各 89 案 pass=true 且 JSON 全同；`mutants` 10／10 exit 1；`verify` PASS（`10 named／1 raw／89 cases／10 state mismatches`）；公開收據 [c-residual-verification.json](../re/c-residual-verification.json)。

- 日期：2026-10-12。
- 探針：`workplace/matching-decompilation/c-residual/ida/ida-probe.json`，SHA-256 `63e12863b199542902260e5f6b917c00e0f268b923ec0149a7f321d3f07a8a4a`；DB SHA-256：`c36861d16e82ddc1fbee8a36e62b49441919fac37a3e5f106b2648b582b58dcd`。
- 輸入：DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具與位址：IDA Pro 9.4（`ida-pro-9.4-idapython:py312-v1`，image `sha256:4ac62de83339c215bab10e455cee3a22d9c6efed9fd0d8ed0f068327b83a06ab`），database linear段基址10000。
- 差集口徑：按 linear address 對帳（`code_blocks` 用位址鍵，名子差集會誤報）。`nullsub_6／7` 與 `1F3CF／1F426／1F43B` 已在帳上是 code_blocks，本輪只做剩下的 10 支。

## 原始契約

- 已證實，10 支 247 指令，所有直接 callee 都已在還原帳上（探針 `unrestored_direct_callees` 全空），不擴張閉包：`sub_15AFC` 內 2 個 INT61 滑鼠 far-call（`0x15B81`／`0x15B8F`，`AX=5`），走既有 far-call 機制。
- 已證實，`sub_145F8`（17 指令 37B，無 callee）：`bx=al*64`，`ah==[bx+19h]` 則寫 `FF`；再 `bx=ah*64`，`al==[bx+19h]` 則寫 `FF`（[`re/69`](../re/69-t2-cross-reference.md) §2、[`mechanics/70`](../mechanics/70-ai.md) §外交清理：兩勢力互指侵攻欄的清除）。
- 已證實，`sub_15AD1`（15 指令）：`15CA4→19528→18853(cx=4)→15AFC→18853(cx=FFFF)`；`cs:198A6&2` 非零才 `call 15E1E` 並 `ds=cs`；尾 `1950F→15C58`。`off_159D2[0x17]`（`0x15A00` 字值 `0x5AD1`）與 `sub_11A5E` 直接呼叫是兩個入口（[`re/31`](../re/31-faction-picker-screen.md)）。
- 已證實，`sub_15AFC`（100 指令 264B）：勢力一覽互動迴圈。開場註冊熱區（`1895D` al=`1F` 矩形 `20,0A,0D08`；`1E3D7` al=`20` 矩形 `208,D8,160E`）、兩次 `1F1A3` 畫框（ah=7／9）、`106F5` 印標題（si=`5C04`「勢力一覽」，ax=`7001`，pos `220,C8`）；`al=0..21` 逐列 `call 15C14`，顏色 `90`／`9A`（`al==cs:10CFF`）／`93`（`al==cs:198A7`）；滑鼠等待先 `AX=5,BX=1`（右鍵就走離場段 `15BEB`：`101B4→1895D` 還原→`101DB→120D6`），再 `AX=5,BX=0` 左鍵命中才往下；`1E453` 解熱區碼，`!=20h` 重掃；列號 `dx-=D8`（`cx>=240h` 則 `dx+=B0`；熱區矩形 x 208h–278h，左半 x<240h 取列 0–10、右半取列 11–19），`ax=dx>>4`，`si=ax<<6`，`[si]<80h` 或 `al==cs:10CFF` 重掃；否則 `101B4`、舊列 `15C14(90)`、寫 `cs:198A7`、新列 `15C14(93)`、`15DBB→15CC6→101DB` 後重掃（[`re/31`](../re/31-faction-picker-screen.md)）。
- 已證實，`sub_15C14`（31 指令）：畫一列。`si=al*64`；`bl>=11` 則 `dx=244` 並 `bx-=11` 否則 `dx=20C`；`y=(bx<<4)+D8`；`[si]<80h` 走 `106F5` 印佔位字（ax=`9001`，si=`5C0D`，前後 `ds=cs`／還原），否則 `call 188B0`（[`re/31`](../re/31-faction-picker-screen.md) §1.3）。
- 已證實，`sub_19377`（39 指令）＋`sub_193D7`（8 指令）：季節漸變。`ds:CF3!=1` 或 `ah=ds:CF0>10h` 直接返；`al=ds:CF4` 月份只認 3／6／9／12（si／di 對 `1934→18A4`、`18A4→18D4`、`18D4→1904`、`1904→1934`，`call 10A65`，al=2／3／4／5 傳曲號）；`193D7` 看 `ah`：1 走 `102C2`（停），2 走 `10241`（播，曲號=al），其餘直接返（[`re/06`](../re/06-game-clock.md) §6）。
- 已證實，`sub_1AB7C`（13 指令）：`[si+19h]=min([si+19h],40)`；`bl=[si+4]`：`=36` 走 `1ABFF`，`>36` 走 `1AC55`，否則走 `1ABD2`（[`re/11`](../re/11-tactical-battle.md)：攻擊上限 40，兵種 `×18` 分派）。另有 raw 片段 `0x1AB48`（1 指令 2 bytes：`EB 32` 尾跳進 `1AB7C`，code 區無函式歸屬），兩側都實際進入。
- 已證實，`sub_1C300`（6 指令 13B）：`adc bl,dl`／`and [di+2],dh`／`xor ah,ah`／`al=5`／`call 1C6BF`／`retn`；後 4 條即 `sub_1C30D` 本體，前 2 條是本入口自帶的前導（CF 敏感，C 照抄語意）。`0x1C0C4` 表字值 `0xC300` 是取址入口，另有 `sub_1C031` 直接呼叫（[`re/60`](../re/60-tactical-sidebar.md)）。
- 已證實，`sub_1E567`／`sub_1E6FF`（各 9 指令）：寫 8 byte 記錄 `[si..si+7]=di,bp,cx,ah,al`，`es=dx`，`al=cs:1E47E／1E47F`，`es:[di]=al`；兩支除計數器位址外逐位元組相同（[`re/08`](../re/08-hourly-update.md)、[`re/69`](../re/69-t2-cross-reference.md) §2.8、[`formats/05`](../formats/05-mmap-worldmap.md)）。
- 假說，`sub_1AB7C` 的 `0x1A82D` 取址在 `funcs_1A827` 表區內（`0x1A82E` 起連續 `AB9C／AB7C／AB39…` 字值），精確槽位未定；不影響本函式 C 語意（入口暫存器契約相同），維持假說不升級。

## 實作與驗收

- 改 `tools/c_recovery/residual.*`（新群）＋ `tools/c_recovery_residual_*` harness（沿 `startup` 群骨架：generate／go／data／verify／sh／container.sh／mutants.tsv）。`1C300` 前導的 `adc` 吃 CF，harness 固定 CF 初值兩側同給。
- 固定矩陣涵蓋：`145F8` 雙向互指／單向／皆無／同勢力／端點五態；`15C14` 左右欄（`bl<11／>=11`）× 空槽／實槽 × 三色；`19377` 月份 3／6／9／12／其他 × `CF3!=1` × `CF0>10` × `ah=1／2／其他`，另 `193D7` 直入五 `ah` × 二曲號；`1AB7C` 上限夾取（`>40／=40／<40`）× 兵種三路（`<36／=36／>36`）× `1AB48` 尾跳進入；`1C300` CF=0／1 × 進位／遮罩邊界；`1E567／1E6FF` 記錄寫入＋計數器回寫 × 四組暫存器；`15AFC` 直入：列掃全列顏色、`cx>=240h` 右欄位移、`[si]<80` 重掃、`10CFF` 重掃、右鍵離場段、左右欄選列切換（舊列還原＋新列選中）；`15AD1` 只走右鍵離開（`98A6` 位元 0／1 兩路、`ds=cs` 還原、尾呼叫全過），選列後接尾段不在矩陣內（見未解範圍）。
- `15AFC` 是滑鼠等待迴圈，以有界事件序列驗收（左鍵命中列／右鍵離場／無效熱區重掃各至少一案）；`int61` 軌跡比對沿 startup 群機制。
- 兩側各自 raw 初始化，固定原始 RNG 與等價輸入；C 不執行原版 CPU 指令，不借原版執行後狀態。
- 以 O0／O2 比較完整 RAM、VGA、DAC、暫存器、FLAGS、SS／SP、IN／OUT、API 及 callee 軌跡；正常返回與轉移分列。每個原始入口都需實際進入；突變只做值變不改控制流，全部拒絕才算過。
- 初態模型需先獨立檢查原版，再比較 C；未確認的玩家語意維持原推論等級。
- 正式 Go 引擎、存檔格式及玩家路徑不在本輪變更範圍；原版素材與私有 dump 不入 Git。

## 未解範圍

| 項目 | 限制 |
|---|---|
| `15AD1` 選列後接尾段、`15E1E` 金額列印 | `15E1E→15F5D→1062F` 印記錄 `[si+20h]` 起 24 位元金額；`1062F` 首除是 16 位元，值 ≥655360 即 `#DE`。prologue（`15CA4`）按劇本建記錄，劇本 2＋ 的部分（劇本，列）格金額超界，真機值小不觸發；矩陣只收右鍵離開案，選列由 `15AFC` 直入案覆蓋。 |
| 完整讀檔長鏈、正常玩家長程及原 C 機器碼 | 本輪局部契約不代證。 |
| `funcs_1A827` 表內 `1AB7C` 槽位 | 取址存在已證實，精確槽位未定；不擋 C 語意。 |
| INT61 聲音驅動的真實硬體行為 | 以 dosgolem 聲音模型為界，不外推實機。 |
| 需 IRQ0 的等待路徑 | 只測就緒／邊界，不在矩陣內觸發 tick。 |
