# 143 — 武將記錄 `+0x17` 是職務值，不是布林

**狀態：CONFORMED。** 原版的 `+0x17` 存 0–4 五個職務值，武將一覽的「身分」欄
直接拿它查字串表；remake 收成 `General.Posted bool`，於是職務靠一覽表
**反推**、俘虜推不出來、存檔寫回把 2／3／4 一律壓成 1。
順帶一起修的是解任時原版會把 `+0x1A`（經費餘額）歸零，remake 沒做。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）
  - `sub_1770C`（`0001770C`）：一覽表逐列繪製，職務值 → 字串表 `cs:75A4h`
  - `sub_16F26`（`00016F26`）寫 1、`sub_16A9B`（`00016A9B`）寫 2、
    `sub_16B71`（`00016B71`）寫 3、`sub_129C3`（`000129C3`）寫 4
  - `sub_16B4F`（`00016B4F`）／`sub_16C2A`（`00016C2A`）清 0 **並清 `+0x1A`**
  - 字串表原始位元組取自 `workplace/orig/dosv/KI.EXE` 檔案位移 `0x77A4`
    （＝ `cs:75A4h` ＋ MZ 標頭 `0x200`）
- 推論等級：**confirmed**（寫入端逐支讀完、讀取端全庫掃過、字串表逐 byte 取出）
- 驗證：[`../playtest/85`](../playtest/85-general-duty-field.md)（dosgolem 逐步讀記錄）＋ §7 的四組單元測試
- 相關：[`38`](38-list-windows.md) §4（身分欄）、
  [`138`](138-state-table-parity.md) §4（壓縮本身）、
  [`142`](142-personnel-dismiss-flow.md)（四條出口的迴圈）、
  [`../re/26`](../re/26-list-window-engine.md) §9（身分名稱表）、
  [`../re/08`](../re/08-hourly-update.md) §4（`sub_16F26` 寫 1）

## 1. 職務值與字串表

`sub_1770C` 畫「身分」欄那三行：

```asm
mov  al, [si+17h]
and  al, al
jnz  short loc_1777F
test byte ptr [si], 40h        ; 職務為 0 才看君主位元
jz   short loc_1777F
mov  al, 5
loc_1777F:
mov  ah, al / shl al,1 / add al,ah / shl al,1   ; al ×= 6
add  ax, 75A4h / mov si, ax
mov  ax, 9001h / call loc_10701                 ; 畫 3 全形字
```

字串表在 `cs:75A4h`，每項 **6 byte ＝ 3 全形字**（Big5，取自檔案位移 `0x77A4`）：

| `+0x17` | bytes | 顯示 |
|---|---|---|
| 0 | `A1 D0 A1 D0 A1 D0` | `－－－` |
| 1 | `AD 78 B9 CE AA F8` | `軍團長` |
| 2 | `A4 BA AC 46 A9 78` | `內政官` |
| 3 | `A5 7E A5 E6 A9 78` | `外交官` |
| 4 | `AB 52 B8 B8 A1 40` | `俘虜　` |
| （5） | `A7 67 A5 44 A1 40` | `君主　` |

⭐ **5 不存在 `+0x17` 裡**，是「職務 ＝ 0 **而且** 記錄 `+0x00` 的 bit 6 立起」時
臨時算出來的。**順序是職務優先**：君主若也編了軍團，`+0x17` 是 1，
顯示「軍團長」而不是「君主」。這與 remake 現行的反推順序相反（§4）。

清單那一列整格是 `FFFF`（空列）時走另一條，畫 `cs:7579h` 的 `－－－`。

## 2. 誰寫

武將表基址 `4240h`、stride 32，所以 `[bx+4257h]` 與 `[bx+17h]` 是同一格。

| 值 | 位置 | 函式 | 時機 |
|---|---|---|---|
| 1 | `00016F26`+… | `sub_16F26` | 編成軍團，主將 → 軍團長 |
| 2 | `000169…`（`loc_16ACB`）| `sub_16A9B` | 任命內政官（緊接著寫據點 `+0x19`）|
| 3 | `loc_16BA6` | `sub_16B71` | 任命外交官（緊接著寫勢力 `+0x2A`）|
| 4 | `000129E5` | `sub_129C3` | 武將被俘 |
| 0 | `sub_16B4F` | 解任內政官 | **同時 `[bx+1Ah] = 0`** |
| 0 | `sub_16C2A` | 解任外交官 | **同時 `[bx+1Ah] = 0`** |
| 0 | `sub_12A7E` | 軍團記錄倒數歸零而消滅 | |
| 0 | `sub_14651` | 軍團解散 | |
| 0 | `sub_14D63` | 據點易主，派駐的內政官遣回 | 不清 `+0x1A` |
| 0 | `sub_15074` | 勢力滅亡，派駐的外交官遣回 | 不清 `+0x1A` |
| 0 | `sub_150B4` | 武將脫離（連軍團一起收）| |
| 0 | `sub_150D7` | 清 `+0x1D`（俘虜關係）時一併清職務 | |
| 0 | `sub_15940` | 俘虜歸順原勢力 | |

⚠ 排除三處長得像但不是這張表的：`sub_13EFD` 的 `[bx+17h]`（勢力表，stride 64）、
`sub_1440F` 的 `xchg al,[bx+17h]`（同樣是勢力表）、`loc_1B03D` 的
`[si+17h]`（戰術段的兵記錄）。

## 3. 誰讀

全庫只有八處讀，**七處是零／非零測試**：

| 函式 | 用途 |
|---|---|
| `sub_137F5` | AI 挑同勢力、無職、政治值最高的武將 |
| `sub_145C1` | AI 挑同勢力、無職、`+0x11` 最小的武將 |
| `sub_14FCE` | 勢力滅亡的逐人處理 |
| `sub_150B4` | 同上的子程序 |
| 線性 `0x176A0`（`sub_17663` 的建清單 callback）| **編成候選清單** （同勢力、存在、無職、不是自己）|
| `sub_16224` | 「進言」的擋下條件：**君主**（`sub_187FF` 由勢力 `+0x01` 取）有職務就跳 TALK #64「主公正在出征中，無法進言。」 |
| `sub_13771` | 每小時處理裡的一個分支 |
| `sub_1770C` | **唯一用到值本身的**——身分欄（§1）|

也就是說：**規則層只要「有沒有職務」，呈現層才要「是哪一種」。**

## 4. remake 現況與差距

`internal/state/state.go` 的 `General.Posted bool ＝ r[0x17] != 0`。四個差距：

1. **存檔寫回不對稱。** `r[0x17] = 1 if Posted else 0`，
   把原版存檔裡的 2／3／4 壓成 1。四個劇本的初始值全是 0
   （[`../re/26`](../re/26-list-window-engine.md) §9），所以 `SINARIO.DAT` 沒事，
   但**玩過的 `SAVE.DAT` round-trip 會壞**。
2. **任命不寫職務。** `cmd/wlgame/personnel.go` 的四條出口只動據點 `+0x19`
   ／勢力 `+0x2A`，沒碰武將那一格 —— 所以派出去的內政官
   **仍然出現在編成候選裡**（`corps.go` 的 `!gen.Posted` 濾不掉他）。
3. **身分欄反推，而且順序相反。** `cmd/wlgame/listfamily.go` 的 `generalRank`
   先查君主再查外交官／內政官／`Posted`；原版是**職務優先、君主墊底**（§1）。
   俘虜（4）反推不出來，永遠顯示 `－－－`。
4. **解任不清經費。** 原版 `sub_16B4F`／`sub_16C2A` 連 `+0x1A` 一起歸零
   （[`../re/25`](../re/25-message-variants-and-personnel.md) §3.1 已記載，未實作），
   remake 只清職務，官員換人後新任會帶著前任沒花完的錢。

## 5. remake 實作

| 檔案 | 改了什麼 |
|---|---|
| `internal/state/state.go` | `Posted bool` → `Duty int`（0–4）＋ 常數 `DutyNone`／`DutyCorpsLeader`／`DutyGovernor`／`DutyDiplomat`／`DutyCaptive`；`Posted()` 收成方法 `Duty != DutyNone`；載入 `int(r[0x17])`、寫回 `byte(g.Duty)` |
| `internal/state/personnel.go`（新）| `AssignGovernor`／`DismissGovernor`／`AssignDiplomat`／`DismissDiplomat` 四支，含 `+0x17` 與 `+0x1A` 的寫入；`cmd/wlgame/personnel.go` 只負責 UI |
| `internal/state/corps.go`、`corpsorder.go`、`events.go`、`freelance.go`、`outcome.go`、`event10_approx.go`、`strategy.go`、`invariant.go` | 賦值端換成對應的 `Duty` 常數；判斷端改 `g.Posted()` |
| `cmd/wlgame/listfamily.go` | `generalRank` 直接讀 `Duty`，只有 `Duty == DutyNone && Sovereign` 才回 5 |
| `cmd/wlgame/statetables.go` | `generalDuty` 從「反查兩張表」改成直接讀 `Duty` |
| `cmd/wlgame/corps.go`、`personnel.go`、`battle.go` | 跟著改 `Posted()` |

俘虜那一格走既有路徑（`outcome.go` 的 `disperseFaction`，對應 `sub_129C3`），
值從 `true` 換成 `DutyCaptive`。

## 6. 存檔 round-trip

`r[0x17] = byte(g.Duty)`，值域 0–4，與原版一致。
既有的 byte-for-byte round-trip 測試（`internal/state` 的 save 測試）
必須在**含已任命官員的存檔**上跑過才算數——目前的樣本是劇本初始局面，
`+0x17` 全 0，壓縮的錯誤在那份樣本上看不出來（§4 第 1 點）。
驗收要另外造一份「先任命再存檔」的樣本。

## 7. 驗證

| 怎麼驗 | 結果 |
|---|---|
| `TestGeneralDutyRoundTripsAllValues`：`+0x17` 灌成 0–4 載入再寫回 | 通過 |
| `TestAssignAndDismissGovernor`／`…Diplomat`：任命寫兩格、解任清三格（含 `Budget == 0`）| 通過 |
| `TestGeneralRankFollowsDutyFirst`：君主兼軍團長回 1 不是 5 | 通過 |
| `TestFormCandidatesExcludeAppointedOfficials`：任命後不在編成候選、解任後回來 | 通過 |
| dosgolem：原版走完「人事 → 內政官任命 → 選城 → 選將 → 解任」，逐步讀武將記錄 | 任命後 `+0x17 = 02`、解任後 `= 00`，**整筆只有那一格變**；武將一覽同時顯示曹操「君主」、夏侯淵「內政官」、夏侯惇「軍團長」（[`../playtest/85`](../playtest/85-general-duty-field.md)）|

## 8. 未解

| 項目 | 現況 |
|---|---|
| `sub_13771` 讀 `+0x17` 的那個分支 | 只知道是每小時處理裡的一支，判斷後 `sub_137F5` 挑人；分支語意未解 |
| `+0x17` 有沒有第六個值 | 字串表只有 6 項而第 6 項要靠 bit 6 算出來，所以存得下的上限是 4；沒有反證 |
| **解任把 `+0x1A` 歸零沒有實跑正對照** | 機器碼確定（`mov byte [bx+1Ah], 0`），但實跑那一輪官員的經費本來就是 0，等於沒比。要先撥款再解任才驗得到 |
| ~~身分欄的畫面對拍~~ | **拍了**（[`../playtest/91`](../playtest/91-general-boast-parity.md)）：武將一覽整張 0 px，含「君王」與「軍團長」兩格 |
