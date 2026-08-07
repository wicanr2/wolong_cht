# 10 — 亂數產生器：`sub_1ECE0` 與 `sub_1EC82`

**狀態：全解，已實作成 `internal/rules/rng`。**

- 日期：2026-08-08
- 輸入：`workplace/ida/dosv/KI.EXE.i64`

規則層每一支公式都在做同一件事：**拿一個 byte 去比門檻**。
傷亡量、災害、外交官成效、武將逃脫、將領修正的那個 25% 分支，
全部是 `and al, 7` 或 `and al, 3Fh` 之後跟一個 `cmp`。
門檻抄對了但分佈不同，長期行為還是會偏，而且偏得很難查——
所以這一支值得單獨解乾淨。

## 1. 取數：`sub_1ECE0`

```asm
mov  ax, cs / mov ds, ax
mov  bx, 0ECFEh              ; 256 byte 的置換表
mov  al, byte_1ECFD          ; 狀態 s
xlat                         ; al = T[s]
add  al, byte_1ECFC          ; + 計數器 c
add  byte_1ECFC, 89h         ; c += 0x89
mov  byte_1ECFD, al          ; s = al
retn                         ; 回傳 al（0–255）
```

```
s ← T[s] + c
c ← c + 0x89
回傳 s
```

狀態只有兩個 byte，加上一張 256 byte 的置換表。
`0x89` 是奇數，所以 `c` 自己會走完 256 個值才回頭。

## 2. 播種：`sub_1EC82`

由 `sub_1006B` 在啟動時呼叫一次。

```asm
xor  al, al / mov bx, 0ECFEh
.fill: mov [bx], al / inc bx / inc al / jnz .fill   ; T[i] = i

mov  ah, 2 / int 1Ah         ; BIOS 即時時鐘，CH/CL/DH 全是 BCD
mov  bl, dh                  ; bl = 秒
add  al, dh / add al, cl     ; al = 秒 + 分
shl  ch, 1 / shl ch, 1
add  al, ch                  ; al += 時 × 4
xor  ah, ah / mov bh, ah / mov dh, ah
push ax / push bx
mov  dl, bl / inc dx         ; 第二個索引 ＝ 秒 + 1

.shuffle:                    ; 256 次交換
  mov  al, [bx-1302h]        ; −0x1302 ≡ +0xECFE (mod 0x10000)
  xchg dx, bx
  xchg al, [bx-1302h]
  xchg dx, bx
  mov  [bx-1302h], al
  add  bl, 4Fh               ; 索引 i 步長 0x4F
  add  dl, 89h               ; 索引 j 步長 0x89
  dec  ah / jnz .shuffle

pop  bx / pop ax
mov  byte_1ECFC, al          ; c = 時×4 + 分 + 秒
xor  al, bl
mov  byte_1ECFD, al          ; s = c ⊕ 秒
```

表先填成 `T[i] = i`，再做 **256 次交換**——所以它在任何時候都是
0–255 的置換，不會有值消失。

⚠ **檔案裡那 256 byte 是全 0。** 直接從 `KI.EXE` 偏移 `0xEEFE`
讀出來會以為表是空的；它是啟動時才建起來的。
拿靜態內容當結論會得到一個「輸出永遠等於計數器」的假產生器。

### 種子只有時鐘，而且是 BCD

`int 1Ah` 回傳的時分秒都是 **BCD**——「35 分」是 `0x35` ＝ 53，不是 35。
個位 A–F 永遠不出現，所以**種子的分佈本來就不均勻**。
洗牌只吃「秒」，時與分只進計數器初值。

## 3. 品質

用參照模型跑出來的數字（`internal/rules/rng` 的測試釘住了同樣的性質）：

| 項目 | 結果 |
|---|---|
| 週期 | 約 **63,000**（狀態空間上限 65,536） |
| 輸出值域 | 0–255 全部出現 |
| 平均 | 127.44（理想 127.5） |
| 低 3 位分佈 | 8,178 – 8,208（理想 8,192） |
| 不同時刻的種子 | 前 8 個輸出有 690 種相異序列 |

對一個 1995 年的 16 位遊戲來說夠用。一場遊戲跑不到一個週期。

## 4. remake 的取捨

`internal/rules/rng` 照抄演算法，但多開一個入口：

- `New(時, 分, 秒)` —— 與原版一致，含 BCD 轉換
- `Now()` —— 用系統時間，等同原版開機時做的事
- `NewFixed(byte)` —— **原版沒有的入口**。長跑驗證需要同一個種子
  跑兩次結果一樣，所以把「秒」那個 byte 直接拉出來當參數，
  **不過 BCD**，讓 256 個種子都能用

在這之前 `cmd/wlsim` 用的是一個線性同餘產生器充數，
註解寫著「等 `sub_1ECE0` 解出來再換掉」——現在換掉了。
