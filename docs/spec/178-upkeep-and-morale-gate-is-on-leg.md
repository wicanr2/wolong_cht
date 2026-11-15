# 178 — 軍費與士氣的判準是「有沒有走在路段上」

**狀態：CONFORMED。** `sub_12600` 的分岔是 `cmp word ptr [si+0Eh], 800h`——
**軍團有沒有走在兩個據點之間的那條路上**，不是它站在哪一種節點上。
remake 拿 `army.KindOf(Corps.Node)` 去問，而 `Node` 在行軍中留著出發那一站，
於是行軍中的軍團收便宜的軍費、**而且每小時回 10 點士氣**。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_12600`（IDA 線性 00012600–00012649）、
  `sub_1562B`（0001562B）
- 推論等級：confirmed（整支 74 B 逐條讀出）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §38
- 相關：[`50`](50-corps-upkeep-charges-funds.md)（軍費走哪一條扣款鏈）、
  [`177`](177-replan-on-leg-not-from-city.md)（同一個 `+0x0E` 判準的另一處）

## 1. 原版做什麼

```asm
sub_12600:
        cmp     cs:byte_10CF3, 1     ; 每天只有「一時」收
        jnz     retn
        mov     ax, [si+4]           ; 兵力
        cmp     word ptr [si+0Eh], 800h
        jb      short loc_12624
        ; ── 走在路段上 ──
        shr ax,1 / mov dx,ax / shr ax,1 / add ax,dx   ; 兵 × 3/4
        mov dh, [si+1] / call sub_1562B               ; 當場扣資金
        retn                                          ; ★ 士氣不回
loc_12624:
        ; ── 沒走在路段上（據點或野外節點都算）──
        mov cl,5 / shr ax,cl / inc ax                 ; 兵 ÷ 32 ＋ 1
        mov dh, [si+1] / call sub_1562B
        mov bh,[si+1] / bx = 勢力 × 64
        mov al, [bx+1Dh]             ; 勢力的士氣基準（+0x1D，開局 200）
        add byte ptr [si+6], 0Ah     ; ★ 士氣 +10
        cmp [si+6], al / jb → / mov [si+6], al        ; 夾到基準
```

⇒ 兩件事共用同一道閘：**軍費差 24 倍、士氣回不回**，
都看 `+0x0E` 有沒有 ≥ `800h`。

野外節點（`600h`–`7FFh`）走的是**便宜**那一邊——「野外」不是這條規則的分野。

## 2. remake 錯在哪

```go
inField := army.KindOf(c.Node) == army.FieldNode   // ⛔
```

`Corps.Node` 是 remake 自己的欄位，行軍中留著**出發那一站**（docs/spec/172）。
從據點出發的軍團，這個式子永遠回 false ⇒ 收便宜的軍費、每小時回 10 點士氣。

同局面對拍量到的：戰後退卻中的軍團 19，原版士氣停在 82，
remake 在 20 拍內回到 92。

`Upkeep`／`Recover` 的參數也跟著改名成 `onLeg`——原本叫 `inField`，
而註解把「道路上」歸到便宜那一邊，與原版相反。

## 3. 改了哪幾支

| 檔案 | 改什麼 |
|---|---|
| `internal/rules/combat/combat.go` | `Upkeep`／`Recover` 的參數 `inField` → `onLeg`，註解訂正 |
| `internal/state/corps.go` | 呼叫端改成 `onLeg := c.LinkAddr != 0` |

## 4. 驗證

`tools/parity_ck.sh` 的拍 2,140／2,200：軍團 19 的 `+0x06`（士氣）
從差 10 回到 0。

## 5. 未解

- `byte_10CF3`（「一時」的判準）remake 用 `hour == upkeepHour` 代替，
  兩者等價與否沒有逐拍證據。

## 2026-11-16：解體返回後仍執行本次軍費尾段

**狀態：CONFORMED，限解體尾段與欄位保留。** 原版正常重播的 196/5/17 一時，軍團56在 `125D2` 已是旗標0，
兵力總欄仍為303；接著執行 `12600` 扣10，士氣82變92。Go 當時在退兵時清總兵力，
又以 `Alive` 擋住尾段，因而少扣10。這項差異會進入 AI 每小時讀取的勢力資金。

來源為相同 DOS/V `KI.EXE`，SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；位址為 IDA Pro9.4 database linear。
`125CF → 125D2 → 125D5` 的 caller 不在移動返回後重驗旗標；`12600` 只檢查一時及路段欄。
`14732` 清單兵槽的兵力，`1466F` 清存在旗標，兩者都不清 `+4` 總兵力或 `+0x23` 階段。

- 已證實的正常原版觀測：[唯讀軌跡](../../workplace/scratch/ai-long-20261116/native-observations.jsonl)。
  新增觀測未改原版資料，4/16、5/1、5/16、6/1的區塊、RNG與佔用圖均重現前輪雜湊，見
  [重播核對](../../workplace/scratch/ai-long-20261116/native-replay-result.json)。
- 局部驗證固定四劇本、兩個主將槽、一時前後、兵力、兵源上限與資金，另含重編成，共480案。
  原版18函式的原始與重定位後 bytes 均核對，整個22208-byte區塊、佔用圖及258-byte RNG作比較，
  來源見[原版收據](../../workplace/scratch/ai-long-20261116/original-proof.json)。舊Go有464案不同。
- 待審查實作：退兵只清單兵槽，不清總欄；解體保留其餘原欄位並寫回存檔；
  本次已進入的軍團巡迴仍執行尾段。它不代表之後持續向不存在的軍團收費。
  重新編成的欄位沿 [spec173](173-corps-flag-bit0-and-sprite-fields.md) 同日補驗。

不變更時間基準、存檔格式、玩家操作或強化AI。長程六表與正常玩家路徑仍須驗證，
不以局部矩陣代證整體AI完成。

證據審查：候選480案全等，原本464案失敗的範圍均已辨明。
實際軍團巡迴原有兩條提前返回，因此只修尾段 helper 尚不足；須將本次尾段掛在共同返回點，
有戰術待續時仍由既有 `ResumeTick` 接回，不能提前執行或重複扣費。
`TestDisbandVisitAndReformationPreserveOriginalRecord` 覆蓋實際巡迴、保存、載入與重新編成，
已通過候選；這條鏈與原版正常觀測共同支持進入正式實作。

正式480案全等，前輪960案求援及391案分派也全等；全專案冷測、vet與建置通過。
六個錯版均被拒絕，其中「保留提前返回」的版本可通過480案局部函式比對，
卻被實際巡迴測試拒絕。這證明不能只查 helper 而漏掉 caller。
正常曹操滑鼠62項／9圖通過編成、停手行軍、自然野戰、退卻與存讀檔。
同一固定4/16來源的正式六表在5/1、5/16、6/1全等；6/16起仍有差異，
全區塊的其他敗走記錄也未據此閉合。
[整合收據](../../workplace/scratch/ai-long-20261116/integration.json)保存原版、正式版本、反例與正常路徑；
[正式長程輸出](../../workplace/scratch/ai-long-20261116/formal-long.log)保留六日期結果，整體AI目標仍進行中。

## 2026-10-10：軍費截斷與士氣回繞勘誤

已證實的原碼契約，來源為 DOS/V `KI.EXE`，SHA-256
`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
以下位址為 IDA Pro 9.4 database linear；固定 probe 與資料庫身分見
[re/132](../re/132-c-army-update-restoration.md)。

- `0x12613–0x12619` 的精確軍費為 `(men >> 1) + (men >> 2)`。
  §1 的「兵 × 3/4」只概述比例，兩次右移各自截斷，不能合併成
  `(men * 3) >> 2`。例如 `men=3` 時，原版軍費為1，合併式為2。
- `0x1263D–0x12646` 先對 `[si+6]` 做 byte 加10，再以無符號比較夾到勢力基準。
  加法先以256回繞，例如246先變成0；相鄰高byte不變。

此勘誤保留先前實測的歷史範圍，不代證新 C 驗收或 Go parity。
`CS:CF3` 與 Go 自然時鐘的排程等價仍未驗證。

## 軍團輪轉與行軍C補證

`sub_12600`（IDA 線性 `0x12600`）與 `sub_1562B`（`0x1562B`）已納入17函式的軍團切片，通過固定raw局部矩陣。
最終狀態與收據見 [re/132](../re/132-c-army-update-restoration.md) 與 [spec/252](252-c-army-update.md)。
較早實測與推論等級維持原範圍；`CS:CF3` 與 Go 自然時鐘的排程等價、完整主排程、正常玩家長程及原C機器碼匹配仍未代證。
