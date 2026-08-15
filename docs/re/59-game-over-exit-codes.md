# 59 — ⭐ 結局與敗北是靠離開碼交出去的

**狀態：三個離開碼與各自的觸發點 confirmed。
⭐ 結局的閘門是**存活勢力數 `== 1`**，而那個計數器就是劇本區塊的 `+0x3A`。**

這解掉 [`43`](43-open-questions.md) 掛的「勢力『滅亡』的精確判定」，
也給了 [`mechanics/80-victory.md`](../mechanics/80-victory.md) 第一個機器碼出處
——那一份先前整份都是「說明書」等級。

- 日期：2026-08-16
- 輸入：`workplace/ida/dosv/KI.EXE.i64`　SHA-256
  `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`；
  `YNVSHELL.COM`（979 B）、`SINARIO.DAT`（88,832 B）
- 工具：`tools/ida_dump.py`、`tools/ida_callsite_args.py`、`tools/ida_xref.idc`、
  `tools/ida_disp_users.py`
- 位址空間：IDA linear（段基址 `0x10000`）

## 1. 怎麼找到的：從「誰啟動 D7END」反追

`KI.EXE` 裡**沒有** `D7END`／`D7OVER` 的字串——所以不是它啟動的。
逐檔 grep 檔名字串，命中的是 979 byte 的 `YNVSHELL.COM`，
裡面一張 FCB 格式的檔名表：

```
02CB  YNFONT  .EXE    02D9  STR     .EXE    02E6  YNSOUND .COM
02F3  LOGO    .EXE    0300  D7OPEN  .EXE    030D  KI      .EXE
031A  D7END   .EXE    0327  D7OVER  .EXE
```

**那是串接器，所以勝負一定是靠 `KI.EXE` 的離開碼傳出去的。**
`start` 的結尾把 `AX` 跨過收尾程序保存起來：

```asm
  call sub_1005B          ; 主迴圈
  call sub_102C2          ; 停音樂
  push ax                 ; ⭐ 保存
  xor ax,ax / int 33h     ; 重設滑鼠
  mov ax,3 / int 10h      ; 回文字模式
  pop ax                  ; ⭐ 取回
  mov ah,4Ch / int 21h    ; AL ＝ 離開碼
```

⚠ **主迴圈 `sub_11BE0` 是無窮迴圈**（`jmp short loc_11C1F`），
出口靠 `cs:word_19901`／`word_19903` 存的 `sp`／`ss`——`sub_11CB1`
把它們還原再 `retn`，`AX` 用 push/pop 保住。
所以「誰決定離開碼」＝「誰呼叫 `sub_11CB1`」，一共四個呼叫點。

## 2. 三個離開碼

| `AL` | 呼叫點 | 去向 | 觸發 |
|---:|---|---|---|
| 0 | `0x160C8` | 回 shell | 正常離開 |
| 1 | `sub_13DC9` `0x13E04` | `D7OVER.EXE` | **信賴度歸零** |
| 1 | `sub_14FCE` `0x14FE5` | `D7OVER.EXE` | **玩家所仕的勢力滅亡** |
| **2** | `sub_11CD0` `0x11D43` | **`D7END.EXE`** | **存活勢力數 ＝ 1**（§3）|

離開碼 2 的那一段還先送兩則訊息：

```asm
loc_11D20:
  call sub_10CDE                     ; PC 喇叭
  cx = 4Bh  / al = 93h / sub_18810   ; TALK #75
  ds = [10D52] / call sub_187FF
  ah = [bx+425Eh] / al = [bx+4241h]  ; 武將記錄的兩個欄位
  cx = 197h / call sub_18810         ; TALK #407
  al = 2 / call sub_11CB1
```

## 3. ⭐ `byte_10D2A` ＝ 存活勢力數，而它是存檔欄位

閘門只有一行：

```asm
00011D0B  cmp cs:byte_10D2A, 1
00011D11  jz  → loc_11D20            ; ⭐ 剩一個勢力 → 結局
```

三件事一起才定案：

1. **全庫只有一個寫入點**，而且是 `dec`（`sub_14FCE` 的 `0x14FE8`）。
   `tools/ida_disp_users.py` 掃位移 `D2A` 也是 0 處——沒有間接寫入。
2. **靜態初值是 0**（`KI.EXE` 的資料區）。只減不加又從 0 起，
   代表它一定是**載入時被覆寫**的。
3. `0x0D2A` 落在劇本區塊前 59 byte 的**最後一格**
   （`cs:0CF0h + 0x3A`，[`../formats/08`](../formats/08-sinario-save.md) §1）。

拿 `SINARIO.DAT` 直接驗：

| 劇本 | `+0x3A` | 該劇本的勢力數 |
|---:|---:|---:|
| 0 | **22** | 22 |
| 1 | **11** | 11 |
| 2 | **6** | 6 |
| 3 | **4** | 4 |

**四個劇本全中。**

## 4. 滅亡 ＝ 據點數歸 0

`sub_14CF3`（據點易主，由攻城入口 `sub_14ADE` 呼叫）：

```asm
00014CF5  xchg bh, [si+1]            ; 換持有者，bh ＝ 舊持有者
00014CFC  cmp bh, 18h / jz → 跳過     ; ⭐ 0x18 ＝ 無主城，不算滅亡
00014D0A  dec byte [bx+23h]          ; ⭐ 舊持有者的據點數 −1
00014D0D  call sub_14DF0             ; 回 CF ＝ 1 表示滅亡
00014D1C  jnb → 跳過
00014D1E  call sub_14FCE
```

`sub_14FCE` 的順序有意義：

```asm
00014FD9  and byte [bx], 7Fh         ; 清勢力的存活位元
00014FDC  cmp bx, cs:word_10CFD      ; ⭐ 滅亡的是玩家所仕的勢力嗎
00014FE1  jnz → 0x14FE8
00014FE3    al = 1 / call sub_11CB1  ; 是 → 敗北，不再往下
00014FE8  dec cs:byte_10D2A          ; ⭐ 才減計數器
```

**先判玩家、再減計數器**——所以玩家自己滅亡時走敗北，
不會因為「剩一個」被誤判成結局。

## 5. 「佔領所有城池」在程式碼裡不是另一條規則

`sub_11CB1` 的四個呼叫點全部列在 §2，**離開碼 2 只有一個來源**。
沒有第二個結局觸發、沒有據點總數的檢查。

關係是單向的：所有據點歸一個勢力 ⇒ 其餘勢力據點數為 0 ⇒ 全部滅亡
⇒ 存活勢力數 ＝ 1。反過來不成立——留著無主城（持有者 `0x18`）時，
存活勢力數照樣可以是 1。

## 6. 未解

| 項目 | 現況 |
|---|---|
| 結局的兩則訊息 | TALK `#75` 與 `#407` 的內容還沒對出來 |
| `sub_14DF0` 的 CF | 「找不到替代據點」與「據點數 0」是不是同一件事，還沒逐行讀 |
| 無主城 `0x18` | 值 24 落在 22 個勢力之外，但劇本裡有沒有無主城沒查過 |
| `D7END.EXE` | 結局過場本身完全沒讀（`END_S*.DAT` 的用法未解）|
