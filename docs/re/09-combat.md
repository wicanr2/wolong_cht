# 09 — 戰鬥：觸發、自動判定、傷亡與武將的下場

**狀態：戰略層自動判定與政略↔戰術入口已解；戰術本體的核心規則已在
[`docs/re/11`](11-tactical-battle.md) 初步讀出並實作。完整結算、少數攻擊分支與同狀態對拍未完。**

- 日期：2026-08-09
- 輸入：`workplace/ida/dosv/KI.EXE.i64`
- 前置：`docs/re/05`（戰場選擇）、`docs/re/08`（每小時更新）

戰鬥其實是兩套系統。一套是**戰術畫面**——玩家的軍團碰上敵人時，
畫面切到戰場，用 `BATTLE.MAP`／`BATTLE.MDL`／`BATTLE.SCH` 那組資產。
另一套是**自動判定**：AI 對 AI 的戰鬥沒有畫面，直接算出勝負與傷亡。

這份文件解的是第二套，外加兩套共用的前後段（觸發條件、壞滅判定、
武將的下場、佔領）。第一套的入口找到了，內容留待日後。

## 1. 戰鬥怎麼被觸發

軍團每走一步會做兩件檢查，都在行軍程序 `sub_12708` 裡。

### 1.1 野戰：撞到人（`sub_12831`）

```asm
mov  bx, [si+0Eh]
mov  di, 2240h                    ; 軍團表
mov  cx, 7Fh                      ; 127 筆
.loop:
  cmp  byte ptr [di], 80h / jb .next   ; 該筆不存在
  cmp  ax, [di+12h] / jnz .next        ; Y 不同
  cmp  dx, [di+10h] / jz  .hit         ; X 相同 → 撞上
.next:
  add  di, 40h ; ★ stride 64
  loop .loop
```

撞上之後比勢力（`[si+1]` 對 `[di+1]`）：**同勢力就穿過去**，
不同勢力才進戰鬥。

### 1.2 攻城：走進別人的據點（`sub_12880`）

```asm
mov  di, [si+0Eh]
cmp  byte ptr [si+0Ah], 4         ; 行進方向
jnz  .a
mov  di, es:[di+8]                ; 連結表：另一條路
jmp  .b
.a: mov di, es:[di+6]
.b: shl di,1 / shl di,1           ; 節點 × 32
mov  al, [si+1]
cmp  [di+841h], al                ; ★ 據點的所屬勢力
jz   .pass                        ; 自家的城 → 直接進去
```

`ds:840h + 據點編號 × 32` 就是**據點表在記憶體裡的位置**
（劇本檔 `+0x8C0`，兩者差 0x80，與交友度表同一個位移關係）。

### 1.3 三岔路：誰來打這一場

兩支的分派邏輯一模一樣，看的是移動中軍團的 `[si+3]`：

| `[si+3]` | 動作 |
|---:|---|
| 0 | `mov byte ptr [si+3], 0Ch` —— 只設狀態，這一步不打 |
| 1 | `call sub_14A7B`／`sub_14ADE` —— **進戰鬥** |
| ≥2 | `mov al,3 / call sub_102F5` —— 丟訊息，不打 |

而且**不論走哪條，都會先 `or byte ptr [si], 20h`**：軍團被擋住了。

## 2. 玩家在不在場，決定是否出現戰鬥方式選單

`sub_14E5C`（野戰）與 `sub_14ED7`（攻城）都先問同一個問題：

```asm
mov  al, cs:byte_10CFF        ; 玩家的勢力編號
cmp  al, [si+1] / jz .playerAttacks
cmp  al, [di+1] / jz .playerDefends
mov  al, 1 / call sub_15130   ; ★ 都不是玩家 → 自動判定
```

⭐ **只有玩家的勢力捲進去，才會出現「戰鬥指揮／委任」選單。**
選「戰鬥指揮」才開戰術畫面；選「委任」則走同一套自動判定。玩家未捲入時
直接自動判定。
另外還有兩個例外會退回自動判定：

- **玩家那一方的軍團委任中**——兩條路各檢查各自那一方：
  玩家是攻方走 `loc_14E70` 的 `test byte ptr [si], 4`、
  玩家是守方走 `loc_14E8A` 的 `test byte ptr [di], 4`。
  ⭐ **判準是「玩家那一方」，與攻守無關**
- **攻城時城裡沒有軍團**（`cmp bx, 4200h`）—— 打空城不進戰術畫面，
  只用城兵自動判定，玩家收到一則訊息

進戰術畫面時走 `sub_11B5A`，戰場編號與旋轉旗標先寫進
`byte_10D34`／`byte_10D35`（`docs/re/05`），交戰雙方寫進
`word_10D2E`（我方）與 `word_10D30`（對方）。

remake 對應為 `internal/state.EncounterChoice`：選單掛起時 `World.Tick`
不推進時鐘；`ChooseBattleCommand` 建立 `Pending` 戰術戰鬥，
`ChooseBattleDelegate` 直接呼叫自動判定並回傳同型的 `CorpsEvent`。
這是執行期狀態，不是存檔欄位。

### 戰鬥前的訊息

| `cx` | 內容 |
|---:|---|
| `0x1A` | 「{2}受到{1}兵馬的攻擊，被攻陷了！！」 |
| `0x1B` | 「{1}的兵馬，向{2}進攻過來了！！」 |
| `0x1C` | 「{1}大人的兵馬，向{2}進攻了！！」 |
| `0x1D` | 「{1}大人的兵馬，遇上{1}的兵馬了！！」 |

`0x1B`／`0x1C` 的差別就是玩家是守方還是攻方，`0x1A` 是空城被拿下。

### ⭐ 野戰與攻城誰先問：據點佔好幾格地圖

`sub_12708`（每走一步）的順序是**先野戰後攻城**，但兩個分支看的是不同的圖：

```asm
cmp byte ptr [di], 0        ; ★ 佔用圖：目標格有沒有軍團
jz  .empty
  call sub_12831            ;   有 → 找同格的敵軍團 → sub_14A7B（野戰）
  jb  .empty                ;   找不到／同勢力 → 落到下面
  retn
.empty:
test byte ptr [si], 1
cmp al, 0CEh / jb .move
cmp al, 0DDh / ja .move     ; ★ 目標格的圖塊值是不是據點（0xCE–0xDD）
call sub_12880              ;   是 → sub_14ADE（攻城）
```

⭐ **關鍵是 `0xCE`–`0xDD` 是一整段圖塊值——一個據點佔好幾格地圖。**
守軍站在其中一格，攻方踏進的通常是**別的那幾格**，那幾格佔用圖是 0，
所以走攻城；接著 `sub_14ADE` 用**據點自己的座標**呼叫 `sub_14C72`
（`mov dx, [di+8]` / `mov ax, [di+0Ah]`，di 是據點記錄）把守軍找出來。

佔用圖是個計數器：`sub_12662` 出發前 `dec byte ptr [di]`、
到達後 `inc byte ptr [di]`（透過軍團記錄 `+0x1A`／`+0x1C` 存的遠指標）。
**待在城裡的軍團一樣佔著它那一格**——所以踏在守軍身上那一格確實會打野戰。

> ⚠ **本專案的據點是一個點**，攻方必然踏在守軍那一格上。照抄順序的話
> 攻城那條路永遠走不到，所以 `internal/state` 的 `resolveContact`
> **把據點放在野戰前面**。這是補上地圖模型的差異，不是規則不同。

## 3. ⭐ 自動判定：`sub_15130`

```asm
mov  ah, al                   ; 保存模式：1 ＝ 野戰、0 ＝ 攻城
and  al, al / jnz .1
mov  al, 3                    ; 攻城 → 攻方用 3 號地形列
.1:
call sub_15285 / mov al, ah / call sub_152D7   ; 攻方戰力 → dx
mov  cx, dx
xchg si, di
call sub_15285 / call sub_152D7                ; 守方戰力 → dx
xchg si, di
add  cx, 8 / add dx, 8        ; 各加 8，避免除以 0
cmp  cx, dx / pushf           ; 記住誰大
jnb  .2 / xchg dx, cx         ; cx ＝ 大的、dx ＝ 小的
.2:
mov  ax, cx / mov bx, dx / xor dx, dx / mov cx, 3
.3: shl ax,1 / rcl dx,1 / loop .3      ; 大的 × 8
div  bx                                ; ★ 比值 ＝ 大 × 8 ÷ 小
cmp  ax, 64h / jb .4 / mov ax, 64h     ; 上限 100
.4: popf / jnb .win_si
```

**勝負只看戰力大小，不擲骰。** 骰子只出現在傷亡量。

`sub_152D7` 出口會把 `ax` 還原，所以第二次呼叫 `sub_15285` 拿到的
`al` 還是原始模式（0 或 1）——這正是攻守用不同地形列的機制。

| | 攻方 | 守方 |
|---|---|---|
| **野戰** | 列 1 | 列 1 |
| **攻城** | 列 3 | 列 0 ＋ 城兵加成 |

### 3.1 基礎戰力：`sub_15285`

```asm
xor ah,ah / push ax
shl ax,1 / shl ax,1 / mov di, ax      ; di ＝ 地形列 × 4
add si, 28h                           ; ★ 軍團 +0x28 起是六個部隊槽
mov cx, 6
.loop:
  mov bl, [si+2] / dec bl / xor bh,bh      ; 兵種 − 1
  mov al, [si+1]                           ; 該槽兵力
  mul byte ptr cs:[bx+di+5120h]            ; ★ × 地形係數
  add dx, ax
  add si, 4                                ; ★ 每槽 4 byte
  loop .loop
pop ax / and al,al / jnz .skip
mov bx, cs:word_10D32
add dl, [bx+13h] / adc dh, 0          ; ★ 守城方加城兵數
.skip:
sub si, 40h
mov al, [si+6] / shr al,1 ×3          ; ★ 士氣 ÷ 8
xor ah,ah / mul dx / mov dx, ax
```

```
基礎戰力 ＝ (士氣 ÷ 8) × Σ(各槽兵力 × 地形係數[兵種])
```

六個槽 × 4 byte ＝ 24，`0x28 + 24 = 0x40`——**軍團記錄剛好 64 byte**。
這是 stride 64 的第三個獨立證據。

### 3.2 地形係數表（`cs:5120h`，檔案偏移 `0x5320`）

十六個 byte，四列 × 四欄（第四欄沒用到）：

```
02 03 03 00 | 03 02 01 00 | 01 03 02 00 | 02 01 02 00
```

| 列 | 場合 | 騎馬 | 弓兵 | 步兵 |
|---:|---|---:|---:|---:|
| 0 | **守城** | 2 | **3** | **3** |
| 1 | **野戰** | **3** | 2 | 1 |
| 2 | （沒有呼叫點會用到） | 1 | 3 | 2 |
| 3 | **攻城** | 2 | 1 | 2 |

**野戰騎馬 3、守城騎馬 2 而弓步各 3。** 這就是說明書
「騎馬隊のみの軍団は移動速度が速くなります」與
「騎馬のみの編成では城壁に登れない」在數值上的體現：
騎兵野外強、城牆邊弱。

列 2 是死資料——`sub_15130` 只會傳 0、1、3。它可能是戰術層在用的，
也可能是開發期留下的。**沒有證據之前不要替它命名。**

### 3.3 將領修正：`sub_152D7`

```asm
mov bh, [si+2] / xor bl,bl / shr bx,1 ×3   ; 武將編號 × 32
push ax
mov dx, [bx+4251h]                ; dl ＝ 武力(+0x11)、dh ＝ 統率(+0x12)
mov ch, dl / shl ch, 1            ; 預設 ch ＝ 武力 × 2
cmp dl, dh / jb .mix
  call sub_1ECE0 / and al,3 / jnz .done    ; 武力 ≥ 統率 → 75% 用 武力×2
  mov dh, dl
.mix:
  mov ch, dh / shr dh,1 / shr dh,1 / sub ch, dh   ; ch ＝ dh × 3/4
  add ch, [bx+4252h]                              ; ＋統率
.done:
pop ax
add bx, 424Eh / xlat              ; ★ 武將 +0x0E(攻城) 或 +0x0F(野戰)
mov cl, 4 / shr al, cl            ; 適性 ＝ 高半位元組，0–10
mov bl, 10h / sub bl, al / xor bh,bh
mov al, ch / mov ah, bh / shl ax, cl
xor dx,dx / div bx                ; bx ＝ 將領值 × 16 ÷ (16 − 適性)
mov bx, ax
pop dx / mul dx                   ; × 基礎戰力
… >> 10
```

```
將領值 ＝  武力 × 2                （武力 ≥ 統率，75%）
        ＝  武力 × 3/4 ＋ 統率      （武力 ≥ 統率，25%）
        ＝  統率 × 7/4              （武力 < 統率）

最終戰力 ＝ 基礎戰力 × 將領值 ÷ (64 × (16 − 適性))
```

兩件事值得記下來。

**第一，`+0x0E`／`+0x0F` 不是兵種適性，是場合適性。** 索引它們的是
「攻城 0／野戰 1」，不是兵種。原本 `docs/formats/08` 把
`+0x0E`–`+0x10` 記成「疑似兵種適性」——三個欄位裡至少有兩個
是**攻城適性**與**野戰適性**。`+0x10` 沒有呼叫點會取到（要 `al = 2`），
與地形係數表的列 2 是同一個缺口。

**第二，將領能力主導戰力。** 適性 0–10 讓分母在 6 到 16 之間，
將領值在 2 到 30 之間；兩者相乘，最強與最弱的將領差到約 **40 倍**。
兵力與士氣的影響遠小於這個。

## 4. ⭐ 傷亡與士氣：`sub_151B3`

輸入 `si` ＝ **勝方**、`di` ＝ **敗方**、`al` ＝ 戰力比值（8–100）。

### 4.1 城市損傷（只在攻城時）

```asm
mov dl, cs:byte_10D34
cmp dl, 0C0h / jnb .skip          ; ≥ 0xC0 是野戰的戰場編號 → 跳過
mov ah, 3Fh / sub ah, al / shr ah,1 / shr ah,1
mov bx, cs:word_10D32
sub [bx+13h], ah / jnb .1 / mov byte ptr [bx+13h], 0    ; 城兵
.1: sub [bx+10h], ah / …                                 ; 上昇值
.2: sub [bx+11h], ah / …                                 ; 防災值
```

攻城戰的戰場編號就是據點編號（0–191），野戰的是 `0xC0` 以上
（`sub_14BDD`／`sub_14C1A` 都 `add cl, 0C0h` 以上），所以這個檢查
把兩種戰鬥分得乾淨。

三個欄位一起扣同一個量，**不分勝敗**：城兵、上昇值、防災值。

### ⚠ 一個溢位

`ah = (0x3F − 比值) >> 2` 是 byte 運算。比值最小 8（勢均力敵）、
最大 100，而 `0x3F` ＝ 63：

| 比值 | 兵力對比 | 城市損傷 |
|---:|---|---:|
| 8 | 1 : 1 | **13** |
| 16 | 2 : 1 | 11 |
| 32 | 4 : 1 | 7 |
| 63 | 約 8 : 1 | 0 |
| **64** | 8 : 1 | **63** ← 減出負數，繞回去了 |
| 100 | ≥ 12 : 1 | 54 |

設計意圖顯然是「苦戰傷城多、輾壓傷城少」，而**懸殊到 8 倍以上時
損傷突然跳到最大**。這是原版的行為，remake 要照做還是修掉，
是另一個決定（見 `docs/mechanics/30-combat.md`）。

### 4.2 逐槽扣兵

```asm
add si, 28h / add di, 28h
xor bx,bx / mov dx, bx
mov ch, al / mov cl, 6
.loop:
  call sub_1ECE0 / and al,7 / add al,2      ; ★ 勝方：2–9
  mov ah, [si+1] / sub ah, al / ja .1
    xor ah,ah / cmp cl,6 / jnz .1 / mov ah,1   ; ★ 大將槽保底 1
  .1: add bl, ah / adc bh,0 / mov [si+1], ah

  call sub_1ECE0 / xor ah,ah / inc ch / div ch
  mov al, ah / add al, 8                    ; ★ 敗方：8 ＋ rand mod (比值+i)
  mov ah, [di+1] / sub ah, al / ja .2
    xor ah,ah / cmp cl,6 / jnz .2 / mov ah,1
  .2: add dl, ah / adc dh,0 / mov [di+1], ah
  add si,4 / add di,4 / dec cl / jnz .loop
```

| | 每槽損失 |
|---|---|
| **勝方** | `rand(0..7) + 2` ＝ 2–9，**與戰力比值無關** |
| **敗方** | `8 + rand mod (比值 + i)`，`i` 是槽序 1–6 |

`inc ch` 在迴圈內，所以除數逐槽遞增。實作要照抄，不要化簡成同一個除數。

**第一槽（大將）永遠留 1。** 自動判定打不死大將的部隊。

### 4.3 士氣

```asm
xchg bx, [si+4]        ; 勝方：新總兵力進 +0x04，bx ＝ 舊值
xchg dx, [di+4]
mov cx, dx
xor ax, ax
and bx,bx / jz .a
cmp byte ptr [si+6], 64h / jb .a          ; ★ 戰前士氣 < 100
mov al, [si+6] / mul word ptr [si+4] / div bx
.a: mov [si+6], al
xor ax, ax
and cx,cx / jz .b
cmp byte ptr [di+6], 64h / jb .b
mov al, 64h / xor ah,ah / mul word ptr [di+4] / div cx
.b: mov [di+6], al
```

```
勝方士氣 ＝ 戰前士氣 × 新兵力 ÷ 舊兵力
敗方士氣 ＝ 100      × 新兵力 ÷ 舊兵力     ← 不是按戰前值衰減
戰前士氣 < 100 的一方 → 士氣歸零 → 壞滅
```

⭐ **敗方的士氣被重設成「100 × 兵力比」，戰前有多高都一樣。**
所以打輸一場之後士氣必定低於 100，**下一場不論勝負都會壞滅**——
輸了就得撤回據點養士氣，這是原版的硬性節奏。

⚠ 說明書 5.5 說的是「100 を切った状態で**戦闘して負けると**軍団が壊滅」，
但機器碼對**勝方一樣清零**。士氣不足 100 的軍團打贏也會散掉。

## 5. 壞滅判定：`sub_1474A`

```asm
call sub_16FD2                       ; 重算總兵力／是否純騎馬
cmp byte ptr [si+6], 0   / jz .dead  ; ★ 士氣 0
cmp byte ptr [si+29h], 0 / jz .dead  ; ★ 大將槽兵力 0
and cl, cl / jz .stay                ; cl ＝ 0 → 勝方，原地不動
mov bx, [si+0Eh] / cmp bx, 600h / jnb .retreat
shl bx,1 / shl bx,1
mov al, [bx+841h] / cmp al, [si+1] / jz .stay   ; 站在自家城裡 → 不退
.retreat:
call sub_1487B / jb .dead            ; ★ 找不到退路 → 壞滅
… 設新目標 [si+14h]／[si+20h]，or byte ptr [si], 2
```

自動判定裡大將槽保底 1，所以實際上**壞滅的唯一入口是士氣 0**。
`[si+29h]` 那一條是給戰術層用的（那裡部隊會被打光）。

退路由 `sub_1487B` 找，開頭就是：

```asm
mov bh, [si+1] / … ; 勢力 × 64
mov al, [bx+3] / cmp al, 0FFh / jnz .ok
stc / retn                            ; ★ 沒有首都 → 無處可退
```

## 6. 武將的下場：`sub_1291A`

輸入 `al` ＝ 勝方勢力編號。

```asm
mov bh, [si+1] / … ; bx ＝ 敗方勢力 × 64
mov al, [bx+1]     ; 君主的武將編號
mov ah, [bx+3]     ; 首都據點編號
mov bx, si / sub bx, 2240h / shr bx,1 / add bx, 4240h   ; ★ 軍團 i ↔ 武將 i

cmp ah, 0FFh          / jz .captured    ; 沒有首都 → 必被擒
cmp al, [si+2]        / jz .escaped     ; ★ 君主親征絕不被擒
mov al, cs:byte_12919
cmp al, [si+1]        / jz .escaped     ; 同勢力
cmp al, 18h           / jz .escaped     ; 勝方是無主勢力(0x18)
call sub_1ECE0 / and al, 7Fh
mov ah, [bx+1Fh] / shr ah,1 / add ah, 28h
cmp al, ah / ja .captured               ; ★ rand(0..127) > 評價÷2+40 → 被擒
.escaped: call sub_12977
.captured: call sub_129C3
```

```
逃脫機率 ＝ (評價 ÷ 2 + 40 + 1) ÷ 128
```

評價 66（呂布、趙雲）→ 約 58%；評價 10（文官）→ 約 36%。
**能力越高越容易脫身。**

### 軍團編號 ＝ 武將編號

`(si − 0x2240) ÷ 2 + 0x4240` 把軍團記錄位址直接換成武將記錄位址：
軍團 64 byte、武將 32 byte，兩張表**同索引平行**。記憶體佈局也吻合：

```
ds:2240h  軍團表 127 × 64 ＝ 8,128
ds:4200h  城兵用的臨時軍團 1 × 64
ds:4240h  武將表 127 × 32 ＝ 4,064
```

中間那 64 byte 不是填充，是 §7 的城兵。

### 兩種結局

| | 函式 | 效果 | 訊息（敗方視角／勝方視角） |
|---|---|---|---|
| **逃脫** | `sub_12977` | 軍團 `[si]=8` 解散、武將保住 | `0x1F`「部隊遭殲滅！幸好以逃過敵軍之手」／`0x20`「沒能將{1}捉到手」 |
| **被擒** | `sub_129C3` | `[bx+1Ch]` ← 勝方、`[bx+1Dh]` ← 舊主 | `0x21`「很遺憾，似乎遭敵軍所擒了」／`0x22`「捉到{1}了」 |

被擒之後還有一條分支：

```asm
mov ah, al / xor al,al / shr ax,1 / shr ax,1   ; 舊主勢力 × 64
mov di, ax
cmp byte ptr [di], 80h / jnb .normal           ; 舊主還在
test byte ptr [bx], 10h / jnz .suicide         ; ★ 武將旗標 bit 4
.suicide:
mov byte ptr [bx], 0
mov word ptr [bx+1Ch], 0FFFFh                  ; 變成在野
… 訊息 0x43
```

訊息 `0x43` 是「{1}在即將被我軍擒拿之前，自刎而死了」。
**舊主已滅 ＋ 武將旗標 bit 4 → 自刎。** 那個 bit 就是「不事二主」。

### ⚠ 訂正：武將 `+0x1D` 不是「派駐狀態」

`docs/formats/08` 原本把它記成「派駐狀態，`0xFF` ＝ 未派駐」。
四處用法一致指向另一件事：

| 位置 | 用法 |
|---|---|
| `sub_129C3` | 被擒時寫入**舊主的勢力編號** |
| `sub_150D7` | `xchg al(0FFh), [di+1Dh]`，並通知該編號的勢力 → 訊息 `0x25`「被敵軍所擒的{1}大人回來了」 |
| `sub_1585F`（月結） | `[si+1Dh] != 0FFh` → `sub_15940`，成功則清成 `0xFF` ＋ 訊息 `0x42`「俘虜的{1}已歸降我軍了」 |
| `sub_16366`（顯示） | `[bx+1Dh] == 0FFh` 才走一般的能力值顯示 |

```
武將 +0x1D ＝ 捕虜狀態（記著舊主的勢力編號），0xFF ＝ 非捕虜
```

開局四個劇本全是 `0xFF`，與「派駐狀態」的觀察一致——
**光看初始值分不出來，是用法把它定下來的。**

## 7. 城兵：`sub_14F8A`

城裡沒有軍團時，攻城戰在 `ds:4200h` 現搭一支臨時軍團：

```asm
mov bx, 4200h
mov al, [di+1]        / mov [bx+1], al       ; 勢力 ＝ 據點的所屬
mov byte ptr [bx+2], 7Fh                     ; ★ 將領 ＝ 127 號（空記錄）
mov byte ptr [bx+6], 0FFh                    ; ★ 士氣 255
mov al, [di+13h] / xor ah,ah / mov [bx+4], ax  ; 兵力 ＝ 城兵數
mov cx, 6 / div cl / add bx, 28h
.loop:
  mov [bx+1], al / mov byte ptr [bx+2], 3    ; ★ 兵種 3 ＝ 步兵
  and ah,ah / jz .1 / dec ah / inc byte ptr [bx+1]   ; 餘數散給前幾槽
  .1: add bx,4 / loop .loop
```

**城兵是六隊步兵、士氣 255、將領是 127 號空記錄。**
`sub_14FC8` 在戰後把 `[bx]` 清 0，臨時記錄用完就丟。

兵種編號也在這裡定案：`3 ＝ 步兵`，與勢力記錄預備兵的
`+0x04 騎馬／+0x06 弓兵／+0x08 步兵` 順序一致，所以
`1 ＝ 騎馬、2 ＝ 弓兵、3 ＝ 步兵`。

⚠ 守城方在 `sub_15285` 裡**又加了一次** `[bx+13h]`（城兵數）。
城兵已經化成六個槽的兵力了，這一加是重複計算——
原版如此，還沒有證據說它是刻意的。

## 8. 佔領：`sub_14CF3`

```asm
mov bh, al / xchg bh, [si+1]      ; 據點的所屬勢力 ← 勝方
mov [si+1Ah], bh                  ; 記下原本的所屬
cmp bh, 18h / jz .neutral         ; ★ 原本是無主(0x18) → 沒有「奪取」
call sub_14D63
… dec byte ptr [bx+23h]           ; 舊主的據點數 −1
call sub_14DF0 … call sub_14FCE
.neutral:
mov bh, al / … / inc byte ptr [bx+23h]     ; 新主的據點數 +1
```

`0x18`（24）是**無主勢力**的編號——勢力只有 22 個（0–21），
`0x18` 是保留給空城的哨兵值。`sub_1291A` 裡「勝方是 `0x18` 就不抓人」
用的是同一個值。

## 9. 補給與士氣回復：`sub_12600`

不是戰鬥的一部分，但決定了戰後怎麼恢復，寫在這裡。

```asm
cmp cs:byte_10CF3, 1 / jz .run        ; ★ ds:0CF3h 是「時」
retn                                   ;   只有一時那個小時才收
.run:
mov ax, [si+4]                        ; 兵力
cmp word ptr [si+0Eh], 800h           ; 節點 × 8 ≥ 0x800 → 節點 ≥ 256 ＝ 野外
jb  .inTown
shr ax,1 / mov dx,ax / shr ax,1 / add ax,dx   ; ★ 野外：兵力 × 3/4
… 扣款
retn
.inTown:
mov cl,5 / shr ax,cl / inc ax                 ; ★ 城／路上：兵力 ÷ 32 ＋ 1
… 扣款
mov bh, [si+1] / … ; 勢力 × 64
mov al, [bx+1Dh]                              ; 士氣基準（＝200）
add byte ptr [si+6], 0Ah                      ; ★ 士氣 +10
cmp [si+6], al / jb .ret / mov [si+6], al     ; 上限 ＝ 基準
```

| 位置 | 每天軍費 | 士氣 |
|---|---|---|
| **據點或道路上** | 兵力 ÷ 32 ＋ 1 | **+10**，上限 200 |
| **野外** | **兵力 × 3/4** | 不回復 |

⚠ **不是每 tick，是每天。** 這一支掛在每 tick 的 `sub_125A3` 上，
但它開頭就檢查 `ds:0CF3h == 1`——那是**小時**（`sub_11D8E` 在
`0x17` ＝ 23 進位），所以實際上一天只收一次。
只看呼叫點會寫成「每 tick」，差 216 倍。

### 而且每 tick 只更新 16 支軍團

```asm
sub_125A3:
  mov si, cs:word_10D18 / mov cx, 10h    ; ★ 從游標開始，只跑 16 筆
  …
  add si, 40h / loop …
  cmp si, 1FC0h / jb .1 / xor si, si     ; 127 × 64 繞回
.1: mov cs:word_10D18, si
```

**軍團是輪流被更新的**，掃完 127 支要 8 個 tick。
這不是效能取巧而已——它讓同一個 tick 裡不同軍團的行動有先後，
先被處理的那一支會先移動、先撞上人。

差距是 24 倍。**野外駐留貴到不可能長期維持**，
而且士氣只在有路的地方回。這解釋了為什麼輸掉的軍團要往城裡退。

## 10. 戰術層的入口

```asm
sub_11B5A:
  push …                     ; 存全部暫存器
  call sub_12078
  mov ax, 2 / call sub_20000 ; 滑鼠：隱藏游標
  call sub_10A1C
  call sub_19946
  call sub_19FA0             ; ★ 戰術戰鬥主體，回傳 al ＝ 勝方
  push ax
  … 還原畫面 …
  pop ax
  mov si, cs:word_10D2E / mov di, cs:word_10D30
  cmp cs:byte_10D35, 80h / jb .1 / xor al,1 / xchg si,di
  .1: … 兩邊各呼叫 sub_1474A，結果進 ah 的 bit 0/1
```

出口與自動判定完全一致：**`al` ＝ 勝方、`ah` 的 bit 0/1 ＝ 哪一方壞滅**，
所以 `sub_14A7B`／`sub_14ADE` 之後的處理兩條路共用。

`sub_20000` 是滑鼠驅動的分派表（`int 33h`），不是 overlay 載入器——
戰術層在 `KI.EXE` 裡面，用 `BATTLE.MAP`／`BATTLE.MDL`／`BATTLE.SCH`
那組資產。

## 11. 目前仍待補完的戰術邊界

| 缺口 | 位置 |
|---|---|
| **戰術完整結算** | `TestNormalScenarioTacticalBattleTerminates` 已證實真實正常攻城的狀態層勝負／傷亡回寫，`wlgame-ai-postbattle.png` 證明正常 GUI 回戰略；GUI 戰後訊息、完整狀態對拍與少數分支仍未完 |
| `sub_1AD7F` 攻擊分支 | `shootSpecial` 已接入 `CH=0x20` 的相鄰格／垂直效果；`+0x1E` 的初始化／上移／下移／交換來源與 `sub_1AC55` 的 raw 比較已確認並接成 `PlaneHigh`，普通箭原版 SCH 單幀圖形已接回，完整投射物動畫／同狀態對拍仍待確認 |
| 原版／remake 同狀態對拍 | 需有效時序原版存檔或可重建的同狀態 oracle |
| `sub_14C72`：怎麼挑出對手軍團 | 野戰與攻城共用，回傳 `bx` |
| 地形係數表的列 2、武將適性 `+0x10` | 兩者都要 `al = 2` 才取得到，戰略層沒有呼叫點 |
| 武將旗標 `+0x00` 的 bit 4（自刎） | 值域已知有 7 種，只解出這一個位元 |
| 據點 `+0x10`／`+0x11` 被攻城扣減 | 欄位語意已知（上昇值／防災值），但「被打過的城成長變慢」還沒在數值上驗過 |
| `[si+3]` 的 0／1／≥2 是誰設的 | 決定哪一支軍團會進戰鬥畫面 |
| `sub_129C3` 的 `[bx+17h] = 4` | 武將 `+0x17` 未解 |
