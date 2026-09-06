# 142 — 解任不先過濾：清單列全部，選到沒人的才回一則訊息

**狀態：CONFORMED。** 原版的「內政官解任」「外交官解任」**不篩掉沒派人的
那些**——清單照列全部，選到沒人的那一個才跳訊息，然後**回到清單繼續選**。
remake 先前先過濾候選，沒有人派駐時直接回一句事件列訊息，連清單都不開。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_16B08`（`00016B08`）
  ＋ `sub_16B4F`（`00016B4F`）、`sub_16BE3`（`00016BE3`）＋ `sub_16C2A`；
  四條出口的分派見 [`../re/25`](../re/25-message-variants-and-personnel.md) §3
- 推論等級：**confirmed**（逐行讀完 ＋ dosgolem 擷取，[`../playtest/84`](../playtest/84-personnel-corps-list-parity.md)）
- 相關：[`126`](126-command-popup-menus.md)（人事那張彈出選單）、
  [`140`](140-status-message-box.md)（狀態列提示）

## 1. 原版做什麼

```asm
sub_16B08（內政官解任）
loc_16B11:
    mov cx, 0Dh / call sub_18853        ; 狀態列 #13「要解任哪個據點的內政官？」
    call sub_17400                      ; 選據點——★ **不過濾**
    jb  離開                            ; CF=1 ＝ 右鍵
    mov si, bx / call sub_16B4F
    jnb 成功
    mov al, 93h / mov cx, 36h / call sub_18810   ; #54「是否有所差錯？{2}並未派遣任何人。」
    jmp loc_16B11                       ; ★ 回去再選
成功:
    mov cx, 1A2h / mov ah,[bx+1Eh] / mov al,[bx+1] / call sub_18810
    jmp loc_16B11                       ; ★ 也回去再選
離開:
    mov cx, 0FFFFh / call sub_18853 / call sub_11D46 / call sub_120D6

sub_16B4F（34 B）
    mov bh, 0FFh
    xchg bh, [si+19h]                   ; ★ 先寫 0FFh，bh ＝ 舊值
    cmp bh, 0FFh / jnz 有人 / stc / retn ; 舊值就是 0FFh ⇒ 那裡沒有內政官
```

外交官那一支（`sub_16BE3`）**逐行同形**，只差四個立即值：
狀態列 `0Eh`（#14）、選勢力用 `sub_17906`、失敗訊息 `37h`（#55
「{3}勢力仍未派遣任何人．．．」）、成功的變體組 `1A3h`。

⭐ **三個怪癖，照抄不修正**：

1. **不過濾**。清單是「全部的據點／勢力」，不是「有派人的那些」。
2. **先寫再檢查**。`xchg` 無條件把 `0FFh` 寫進去，才看舊值是不是 `0FFh`
   ——選到沒人的那一格等於寫了一次同樣的值，沒有副作用。
3. **選完回到清單**。成功或失敗都 `jmp` 回迴圈開頭（連狀態列都重設一次），
   **右鍵才離開**。任命那兩支（`sub_16A9B`／`sub_16B71`）也是同一個迴圈。

## 2. 演算法

```
解任（內政官／外交官各一份，只差欄位與四個訊息編號）：
  迴圈：
    掛狀態列提示
    開清單（玩家的全部據點／全部他勢力）；右鍵 → 清狀態列並離開
    舊值 ← 目標欄位；目標欄位 ← 0xFF
    舊值 == 0xFF → 跳「並未派遣任何人」的訊息
    否則         → 那位官員說一句（變體組 ＋ 他的肖像）
    回到迴圈開頭
```

## 3. remake 實作

| 項目 | 位置 |
|---|---|
| 內政官解任 | `cmd/wlgame/personnel.go` 的 `removeGovernor`：候選改成 `playerCities()` 全部；選到沒有內政官的跳 TALK #54、有人的照舊 |
| 外交官解任 | 同檔 `removeDiplomat`：候選改成「活著且不是自己」的全部勢力；沒人跳 TALK #55 |
| 狀態列 | 四條出口各自 `setStatusTalk`（[`140`](140-status-message-box.md)）|
| 差異 | **選完不回到清單**——remake 的清單 callback 回 `true` 就收掉。原版是迴圈，右鍵才離開。這一條列在 §5 |

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 對原版 ✅ | [`../playtest/84`](../playtest/84-personnel-corps-list-parity.md)：內政官解任那一張的清單**十列全在**，與原版逐列相同 |
| 單元測試 | `TestDismissListsEveryTarget`（`cmd/wlgame`）：沒有任何人派駐時清單仍然開、列數 ＝ 全部 |
| 單元測試 | `TestDismissEmptySlotReportsNobody`：選到沒人的那一個回 TALK #54／#55，不寫壞資料 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| 選完回到清單的迴圈 | 原版成功或失敗都回清單繼續選，**remake 選完就收掉清單**。要改得動 `listPick` 的回傳語意，影響四條出口以外的地方，這一輪沒動 |
| 成功時那位官員說的話 | 變體組 `1A2h`／`1A3h`（＝ TALK 418／419 那兩組八個）。remake 目前沒有這一則 |
