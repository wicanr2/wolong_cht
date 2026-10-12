# 134：C 殘部群還原

**狀態：CONFORMED。10 函式＋1 raw 尾跳，89 案 O0／O2 全同，10 突變全拒。**

- 日期：2026-10-12。
- 規格：[spec/254](../spec/254-c-residual-faction-season.md)。
- 輸入：DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 固定IDA probe SHA-256：`63e12863b199542902260e5f6b917c00e0f268b923ec0149a7f321d3f07a8a4a`；DB SHA-256：`c36861d16e82ddc1fbee8a36e62b49441919fac37a3e5f106b2648b582b58dcd`。
- 工具與位址：IDA Pro 9.4（`ida-pro-9.4-idapython:py312-v1`），database linear段基址10000。
- 公開收據：[c-residual-verification.json](c-residual-verification.json)（`wolong-c-residual-verification-v1`，status `semantic-conformed`）；oracle golem `a9714ebdab2ad6b529f81225472680f2b11f2842`。

## 範圍

census 按位址差集的最後 10 支（247 指令 619 bytes），另加 raw 尾跳 `0x1AB48`（1 指令 `EB 32`，code 區無函式歸屬）。`nullsub_6／7` 與 `1F3CF／1F426／1F43B` 已在帳上是 code_blocks（名子差集會誤報，位址差集為準）。全部直接 callee 已在還原帳上，不擴張閉包；`15AFC` 內 2 個 INT61 滑鼠 far-call 走既有機制。

## 已證實的定位

- `sub_145F8`：`bx=al*64`，`ah==[bx+19h]` 寫 `FF`；再 `bx=ah*64`，`al==[bx+19h]` 寫 `FF`。兩勢力互指侵攻欄清除（[re/69](69-t2-cross-reference.md)）。
- `sub_15AD1`：`15CA4→19528→18853(4)→15AFC→18853(FFFF)`；`cs:198A6&2` 才 `15E1E`＋`ds=cs`；尾 `1950F→15C58`。入口：`off_159D2[0x17]`（`0x15A00` 字 `0x5AD1`）與 `sub_11A5E` 直呼。
- `sub_15AFC`：熱區 `1F`（`20,0A,0D08`）與 `20`（`208h,D8h,160Eh`，x 208h–278h × y 216–392）；標題 `106F5`（`5C04`「勢力一覽」）；22 列色 `90／9A（==10CFF）／93（==198A7）`；右鍵離場段 `15BEB`；左鍵 `1E453` 解碼，`dx-=D8`（`cx>=240h` 加 `B0`），`[si]<80h` 或 `==10CFF` 重掃，否則舊列還原＋`198A7` 更新＋新列選中＋`15DBB→15CC6`。
- `sub_15C14`：`si=al*64`；`bl>=11` 右欄（`dx=244h`）否則左欄（`dx=20Ch`）；`y=(bx<<4)+D8`；`[si]<80h` 走 `106F5` 佔位字（`9001／5C0D`）否則 `188B0`。
- `sub_19377`：`CF3!=1` 或 `CF0>10h` 早返；月份只認 3／6／9／12（`si／di` 對＋`10A65`＋曲號 2／3／4／5）。`sub_193D7`：`ah==1→102C2`（停）、`ah==2→10241`（播）。
- `sub_1AB7C`：`[si+19h]=min(40)`；`[si+4]` `=36→1ABFF`、`>36→1AC55`、否則 `1ABD2`。`0x1C0C4` 表字 `0xC300` 是 `sub_1C300` 取址入口（前導 `adc／and` 照抄，CF 敏感）。
- `sub_1E567／1E6FF`：8 byte 記錄＋`es:[di]=cs:1E47E／1E47F`，除計數器位址外逐位元組相同。

## 已證實的輸入與邊界

- 熱區座標全是十六進：`DX=208h`=520（矩形 x 520–632）、`CX>=240h`=576 分左右半。十進誤讀（208／240）點到矩形外，全部重掃。
- `15B51` 以 `cs:0D52` 蓋 DS：fixture 設 `0x8000`，否則列存在性讀低記憶體，選中靠運氣。
- `15AD1` prologue 以 entry-DS 寫遊戲狀態，會蓋掉 `0x8000` 記錄區；entry 案 entry-DS 用獨立段 `0x8300`（`15AFC` 經 `0D52` 切回 `0x8000`）。
- `15E1E→15F5D→1062F` 印記錄 `[si+20h]` 起 24 位元金額；`1062F` 首除是 16 位元除 10，值 ≥655360 即 `#DE`（實測 `dx:ax=0x13FCE5`）。真機值小不觸發；矩陣的 `15AD1` 案只收右鍵離開，`a6=2` 只留劇本 1 一案。
- 滑鼠等待以 `0x5B7B`（jump 目標）為檢查點：產生器加 `HOOK_ENTRIES`（同 startup 的 `PAUSE_ENTRIES` 機制），snapshot 對非 `functions` 位址早返，軌跡不差一筆。
- prePress 在迴圈前就被清掉（startup transfer 群同例）；picker 案一律 `pollFeed`。

## 驗證與實作界線

- 89案／opt：9群（diplomacy 5、draw-queue 4、map-node 4、picker-entry 4、picker-loop 5、picker-row 12、season 34、sidebar 8、tactical-cap 13）；stage 全 return；獨立稽核578。11入口全見。
- 10突變全以exit 1拒絕。m7 重設計過：只改比較立即值（`cmp 40→41`）是語意空操作（存值仍寫 40，各路皆同），改突變存值（`mov 40→41`）才拒絕。
- C不執行原版CPU指令，不借原版執行後狀態；兩側各自raw初始化。
- 殘部群的C補證；134-c-residual-restoration.md；254-c-residual-faction-season.md。

## 未解範圍

| 項目 | 狀態 |
|---|---|
| `15AD1` 選列後接尾段、`15E1E` 金額列印 | `15CA4` 按劇本建記錄的輸入未完整建模；部分（劇本，列）格金額超界。選列由 `15AFC` 直入案覆蓋。 |
| `funcs_1A827` 表內 `1AB7C` 槽位 | 取址存在已證實（`0x1A82E` 起連續 `AB9C／AB7C／…`），精確槽位未定；不擋 C 語意。 |
| 完整讀檔長鏈、正常玩家長程及原 C 機器碼 | 本輪局部契約不代證。 |
| INT61聲音驅動的真實硬體行為 | 以dosgolem聲音模型為界，不外推實機。 |
| 需IRQ0的等待路徑 | 只測就緒／邊界，不在矩陣內觸發tick。 |
