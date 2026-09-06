# 145 — 指令列最後兩格：武將是迴圈並自陳兵種，勢力是跳到首都開情報卡

**狀態：CONFORMED。** 指令列第 6 格「武將」與第 7 格「勢力」是八格裡最後兩個
沒對過的。原版兩格各有一則狀態列提示，而且**選完之後都還有下一步**：
武將那格會讓被點到的人自陳擅長的戰場（由三個適性挑一組台詞），
然後**跳回清單繼續選**；勢力那格把鏡頭移到該勢力的首都並開那個據點的情報卡。
remake 兩格都只寫了一行事件訊息就把清單關掉。

- 日期：2026-09-06
- 出處：`KI.EXE`（SHA-256 `fffeba98…d43868`）的 `sub_16366`（`00016366`，
  指令列 #6）與 `sub_163BF`（`000163BF`，#7）；分派表 `funcs_161FE`
- 推論等級：**confirmed**（兩支各 90／64 B，逐行讀完）
- 驗證：[`../playtest/86`](../playtest/86-general-faction-cells-parity.md)（兩張 `map` 0 px ＋ 台詞的原版錨點）、
  [`../playtest/91`](../playtest/91-general-boast-parity.md)（自陳那一張逐像素）
- 相關：[`140`](140-status-message-box.md)（狀態列）、
  [`142`](142-personnel-dismiss-flow.md)（同一種「選完回清單」的迴圈）、
  [`23`](23-city-info-window.md)（`sub_17E1F` 情報卡）、
  [`../re/26`](../re/26-list-window-engine.md) §4.1（`sub_175FA`／`sub_178A7`）

## 1. `sub_16366`：武將

```asm
sub_16366:
    call sub_12078
    mov  cx, 18h / call sub_18853        ; 狀態列 #24「確認指示之武將的能力。」
    call sub_175FA                       ; 選武將（本勢力一覽）
    pushf / call sub_120D6 / popf
    jnb  loc_16380
    mov  cx, 0FFFFh / call sub_18853     ; ★ 右鍵才走這裡：清狀態列 → retn
    retn
loc_16380:
    mov  ds, cs:word_10D52
    mov  cx, 1A8h
    cmp  byte ptr [bx+1Dh], 0FFh         ; 是俘虜嗎（+0x1D ≠ 0xFF）
    jnz  loc_163B3                       ; ★ 是 → 組 0x1A8
    mov  al,[bx+0Eh] / mov ah,[bx+0Fh] / mov cl,[bx+10h]
    cmp  al, ah / jnb .1 / mov al, ah    ; al = max(攻城, 野戰)
.1: cmp  al, cl / jnb .2 / mov al, cl    ; al = max(al, 水戰)
.2: mov  cx, 1A9h
    cmp  al, [bx+0Eh] / jz loc_163B3     ; 最大值就是攻城 → 0x1A9
    inc  cx
    cmp  al, [bx+0Fh] / jz loc_163B3     ; 是野戰 → 0x1AA
    inc  cx                              ; 否則水戰 → 0x1AB
loc_163B3:
    mov  ah, [bx+1Eh] / mov al, [bx+1]   ; 變體 ＋ 肖像
    call sub_18810
    jmp  short sub_16366                 ; ★ 回開頭（連狀態列都重設一次）
```

四組台詞（`0x196 + (組−0x196)×8`，[`../re/25`](../re/25-message-variants-and-personnel.md) §1）：

| 組 | 展開 | 內容 | 條件 |
|---|---|---|---|
| `0x1A8` | #550–557 | 俘虜的推託（「我不會投降的．．我一定要逃走。」）| `+0x1D ≠ 0xFF` |
| `0x1A9` | #558–565 | 城塞戰（「我最擅長城塞戰了。」）| 攻城適性 `+0x0E` 最高 |
| `0x1AA` | #566–573 | 野戰（「野戰的話，就由我出陣吧！」）| 野戰適性 `+0x0F` 最高 |
| `0x1AB` | #574–581 | 海戰（「我擅長於海戰！」）| 水戰適性 `+0x10` 最高 |

⭐ **平手歸前面那一個**：判斷是「最大值等不等於攻城」→「等不等於野戰」→
否則水戰，所以三個一樣高時說城塞戰。

⚠ 比較用的是**整個 byte**，remake 存的 `Aptitude` 已經 `>>4`。
四個劇本的三個適性欄位低半位元組全是 0（`../formats/08` §3），
所以兩種比法等價；**低半位元組哪天解出有值，這裡要回頭改**。

### 1.1 自陳的同時**不寫事件列**

remake 的大地圖底下有一條事件列（原版沒有，是 remake 自己加的）。
武將那格原本在 callback 裡寫一句「選擇了 ○○○」——那是台詞還沒接上以前的
暫時產物，接上之後就變成**同一件事講兩遍**，而且事件列剛好蓋在
地圖上，逐像素對拍時是這一張唯一的差異
（[`../playtest/91`](../playtest/91-general-boast-parity.md)：4,015 px 全在那條列上）。

**判準**：事件列補的是**玩家沒在看的時候發生的事**（時間推進帶來的：
月結、行軍、宣戰、災害）。**玩家自己剛下的指令不寫**——那些原版已經有
自己的回報（官員台詞、訊息框），再登記一次就是同一件事講兩遍，
而且那條列蓋在地圖上，逐像素對拍時往往是最大的一塊差異。

同一條判準套到人事四條（[`142`](142-personnel-dismiss-flow.md)）：
任命／解任成功的四句都拿掉了，那位官員自己說的那一句就是回報
（[`../playtest/92`](../playtest/92-personnel-assign-parity.md)：5,808 px）。
留著的只有「沒有據點」這種**原版走不到的狀態**——那是 remake 自己的防呆。

## 2. `sub_163BF`：勢力

```asm
sub_163BF:
    call sub_12078
    mov  cx, 19h / call sub_18853        ; 狀態列 #25「將游標移動至指示之勢力的首都據點。」
    call sub_178A7                       ; 選勢力（交友度一覽）
    pushf / call sub_120D6 / popf
    jb   loc_163FE                       ; 右鍵 → 直接清狀態列 retn
    mov  ds, cs:word_10D52
    mov  bh, [bx+3] / xor bl,bl / shr bx,1 ×3   ; bx = 該勢力首都 × 32（據點記錄）
    mov  si, bx
    mov  dx, [bx+848h] / mov bx, [bx+84Ah]      ; 據點的 X / Y
    mov  ax, 14h / mov cx, 0Ch
    call sub_12151                       ; 鏡頭移到 (X−20, Y−12)
    call sub_11F7F                       ; 更新游標
    call sub_11D46                       ; ★ 用新鏡頭重畫大地圖
    call sub_17E1F                       ; 那個據點的情報卡
loc_163FE:
    mov  cx, 0FFFFh / call sub_18853
    retn                                 ; ★ 不迴圈，一次就結束
```

⭐ **這一格會重畫大地圖**，所以彈出選單／清單的殘影在這裡是會被擦掉的
——與 [`126`](126-command-popup-menus.md) §1.2「選完不擦」不衝突：
留不留框是「有沒有重畫地圖」的副產物，不是選項屬性。

## 3. remake 現況

| | 原版 | remake |
|---|---|---|
| 武將 狀態列 | #24 | **沒有** |
| 武將 選完 | 自陳兵種 ＋ **回清單** | 寫一行 `lastEvent` 就關掉清單 |
| 勢力 狀態列 | #25 | **沒有** |
| 勢力 選完 | 鏡頭移到首都 ＋ 情報卡 | 寫一行 `lastEvent` 就關掉清單 |

## 4. remake 實作

| 檔案 | 改了什麼 |
|---|---|
| `cmd/wlgame/main.go` | `openGeneralList` 掛狀態列 #24；`pick` 改成 `generalSaysAptitude()` ＋ `return false`（回清單）|
| `cmd/wlgame/finance.go` | `openFactionList` 掛狀態列 #25；`pick` 改成 `focusCity(該勢力首都)` ＋ `return true` |
| `cmd/wlgame/personnel.go` | `officialSays` 已經是「組 ＋ 變體 ＋ 肖像」的通用形，直接重用 |

## 5. 驗證

| 怎麼驗 | 結果 |
|---|---|
| `TestGeneralAptitudeTalkBranches`：四條分支 ＋ 三種平手樣本 | 通過 |
| `TestGeneralAptitudeTalkMatchesOriginalCapture`：原版點夏侯淵說 #573，remake 展開到同一則 | 通過 |
| `TestFactionCellFocusesCapital`：鏡頭 ＝ 首都 −(20,12)、情報卡開著、狀態列還在 | 通過 |
| dosgolem 對拍兩格剛開的畫面 | **兩張 `map`／`minimap`／`faction`／`banner` 全 0 px**，只剩游標 88 px（[`../playtest/86`](../playtest/86-general-faction-cells-parity.md)）|
| 勢力選完之後的畫面 | 鏡頭、情報卡、指令列反白都對上；`map` 的殘差全部有名字（同上 §4）|

## 6. 未解

| 項目 | 現況 |
|---|---|
| `sub_175FA`／`sub_178A7` 的清單是不是只列本勢力 | 武將那張只有本勢力、勢力那張列全部活著的，兩張都拍過了（[`../playtest/86`](../playtest/86-general-faction-cells-parity.md)）|
| 適性欄位的低半位元組 | 四個劇本全是 0，語意未解 |
