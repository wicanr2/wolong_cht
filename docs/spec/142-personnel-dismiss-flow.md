# 142 — 人事四條出口是一個迴圈：不過濾、選完回清單、官員自己說一句

**狀態：CONFORMED。** 原版的人事四條（內政官／外交官各任命與解任）
**是同一個迴圈**：不篩候選、選完不關清單、成功或失敗都回去繼續選，
**右鍵才離開**；而且成功時那位官員會自己說一句（八格一組，由 `+0x1E` 選）。
remake 先前四條都是「選一次就收掉」，解任還先過濾候選。

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
| 內政官／外交官任命 | 同檔 `pickCityForGovernor`／`pickFactionForDiplomat`：那裡已經有人就跳 TALK #52／#53；選武將那一步掛狀態列 #9 |
| **迴圈** | 四條的清單 callback 一律回 `false`（＝不關清單）；任命成功之後再呼叫自己一次，等於原版的 `jmp` 回迴圈開頭（連狀態列都重設）|
| **回迴圈的時機** | ⭐ **等那一句被按掉之後才回**（§3.1）：原版的 `sub_18810` 是**擋住的**，訊息還在畫面上時流程停在原地，所以那一刻看到的是**武將一覽 ＋ 狀態列 #9**，不是回去以後的據點一覽 ＋ #11 |
| 官員說的一句 | 同檔 `officialSays(base, who)`：`resolveBattleTalkIndex(base, 武將.TalkVariant)` ＋ 那位武將的肖像。四個組 `19Ch`／`19Dh`／`1A2h`／`1A3h` |
| 狀態列 | 四條出口各自 `setStatusTalk`（[`140`](140-status-message-box.md)）|
| 差異 | 無 |

### 3.1 `afterTalk`：把原版的「擋住」搬成續行

原版 `sub_18810` 畫完訊息就**等**，玩家按掉才 `jmp` 回迴圈開頭。
remake 的訊息是佇列式的（`enqueueTalk` 立刻回傳），照抄呼叫順序會讓
**訊息還掛著、底下的畫面已經是下一步**——逐像素比就是整張清單換掉
（[`../playtest/92`](../playtest/92-personnel-assign-parity.md)：19,927 px）。

所以 `messageDialog` 多一個 `then func()`，`afterTalk(fn)` 掛在剛入列的那一則上，
`updateMessageOnly` 把那一則收掉之後才呼叫。**這一條對每一個
「先說一句、再換畫面」的流程都成立**，不只人事。

## 4. 驗證

| 方式 | 證據 |
|---|---|
| 對原版 ✅ | [`../playtest/84`](../playtest/84-personnel-corps-list-parity.md)：內政官解任那一張的清單**十列全在**，與原版逐列相同 |
| 單元測試 | `TestDismissListsEveryTarget`（`cmd/wlgame`）：沒有任何人派駐時清單仍然開、列數 ＝ 全部 |
| 單元測試 | `TestDismissEmptySlotReportsNobody`：選到沒人的那一個回 TALK #54／#55，不寫壞資料 |
| 單元測試 | `TestPersonnelFlowsLoopBackToTheList`：四條的 callback 一律不關清單 |
| 單元測試 | `TestOfficialLineUsesVariantGroup`：四個組展開後的索引（457／461／465／502／510）|
| 突變測試 | 把 `return false` 改回 `true`、把過濾加回去，各要有測試變紅 |

## 5. 未解

| 項目 | 現況 |
|---|---|
| ~~武將記錄 `+0x17` 被收成 bool~~ | **已修**，另開的那一份是 [`143`](143-general-duty-field.md)：`General.Duty` 存 0–4，任命寫 2／3、解任連經費一起清 |
| 任命的「已經有人」訊息參數 | 原版 `push ax`（`ah = 0FFh`、`al` ＝ 武將編號）＋ `push bx`（**據點記錄位址**，直接位址式）。remake 直接代名字字串，**沒有走 formatter 的位址式**（[`../re/79`](../re/79-talk-marker-handlers.md) §2）|
