# 184 — 求援調兵的 `want` 恆為 1

**狀態：CONFORMED。** `sub_14057` 在「這一格的軍團數 > 1」那一支傳給
`sub_14155` 的 `dl`（want）**永遠是 1**——它是最開頭那個「亂數 & 3」
在挑目標的迴圈裡被 `dec` 到 0 之後，再被 `mov al, 1` 設回來的。

remake 拿 `Requested(威脅量, 佔用數)` 當它——那是**求援的數量**，
走的是另一支。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_14057`（IDA 線性 00014057–000140B2）、
  `sub_14155`（00014155）
- 推論等級：confirmed（整支逐條讀出，而且與原版 `-watch` log 的暫存器對上）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §43
- 相關：[`180`](180-relief-dispatch-want-counts-candidates.md)（同一支的 `dec dl`）、
  [`../re/40`](../re/40-garrison-relief-request.md)

## 1. 原版做什麼

```asm
sub_14057:
        cmp  byte ptr [bp+0], 0FEh / jnb retn   ; 目標緩衝為空 → 結束
        call sub_1ECE0 / and al, 3              ; ★ 亂數 & 3
loc_14062:
        mov  di, bp                             ; ★ 掃到空槽就從頭再來
loc_14064:
        cmp  byte ptr ss:[di], 0FEh / jnb loc_14062
        dec  al / jz loc_14073                  ; ★ 只有 al 減到 0 才跳出
        add  di, 4 / jmp loc_14064
loc_14073:                                      ; ⇒ 到這裡時 al 一定是 0
        cmp  byte ptr [si+858h], 1 / ja loc_14099
        ; ── 佔用數 ≤ 1：求援 ──
        mov  al, [si+854h] / add al, 2 / sub al, [si+858h]
        jbe  retn                               ; 威脅量 ＋ 2 − 佔用數 ≤ 0 ⇒ 不求援
        （玩家的據點不走這條）
        call sub_140C9 / call sub_140B3
        retn
loc_14099:                                      ; ── 佔用數 > 1：調兵 ──
        and  al, al / jnz loc_1409F
        mov  al, 1                              ; ★ al 是 0 ⇒ 設成 1
loc_1409F:
        mov  cl, ss:[di] / mov ch, 0            ; 目標據點
        mov  dl, al                             ; ★ want ＝ 1
        mov  dh, [si+858h]                      ; 佔用數
        call sub_14155                          ; 進去先 `sub dh, dl` ⇒ skip
```

⇒ **`want` 恆為 1，`skip` ＝ 佔用數 − 1。**

⭐ 這個結論有第二個來源：原版 `-watch` log 在拍 3,415 的
`sub_14155` 第一次取亂數那一行是 `DX=0201`——`sub dh, dl` 之後
dh ＝ 2、dl ＝ 1，而那一格的佔用數是 3。

## 2. remake 錯在哪

`relieve` 兩支共用一個 `want := threat.Requested(c.Threat, c.Occupancy)`。
同局面拍 3,415 的據點 129：威脅量 3、佔用數 3 ⇒ `Requested` 給 **2**，
於是多掃一支候選、多取一次亂數，還調走了原版不會動的那一支。

`Requested` 本身沒錯——它是「佔用數 ≤ 1」那一支的求援數量
（`al = [si+854h] + 2 − [si+858h]`，`jbe` 就不求援）。錯在把它用到另一支。

## 3. 驗證

`tools/parity_ck.sh` **二十二個檢查點（拍 200 … 3,420）的據點、勢力、
軍團、全域四張表全部 0 個 byte**。

逐拍取數不一致 113 → **78 / 5,480（1.4%）**，第一個分歧從拍 3,415
推到 **3,566**。

## 4. 未解

- 拍 3,566：原版取了 `sub_14057` 的挑目標亂數而 remake 沒有
  ⇒ remake 的 `Targets` 是空的。侵攻目標或鄰居歸屬還沒對過。
- 挑目標的位置在 `al` 初值為 0 時不等價：原版 `dec al` → `0FFh`，
  要繞 255 步（掃到空槽會 `mov di, bp` 從頭再來）才回到 0；
  remake 是 `亂數 & 3 % len(Targets)`。目標只有一個時兩者相同，
  多個時會分岔。
