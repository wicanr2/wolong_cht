# 180 — 求援調兵的 `want` 數的是候選不是選中

**狀態：CONFORMED。** `sub_14155` 的 `dec dl` 在 `loc_1418A`——位在
「收」與「不收」兩條路**匯流之後**。所以 `want` 是「掃過幾支候選」，
不是「收到幾支」；一支 `+0x23 >= 8` 的守軍會吃掉一個額度，
額度歸零就整支結束。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_14155`（IDA 線性 00014155–00014193，63 B）、
  `sub_14499`（00014499）、`sub_16FD2`（00016FD2）
- 推論等級：confirmed（整支逐條讀出）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §40
- 相關：[`../re/40`](../re/40-garrison-relief-request.md) §5、
  [`179`](179-recalc-on-battle-and-arrival-test.md)（`sub_16FD2` 的計時）

## 1. 原版做什麼

```asm
sub_14155:
        sub     dh, dl              ; dh ＝ 佔用數 − want ＝ 可以跳過幾支
        mov     di, 2240h           ; 軍團表
        mov     bx, si / shr bx,1 ×2 ; bx ＝ 據點編號 × 8
loc_14160:
        cmp     bx, [di+0Eh]        ; 這支軍團的現在節點就是這個據點？
        jnz     short loc_1418E
        cmp     byte ptr [di], 80h
        jb      short loc_1418E     ; 不活著 → 下一支
        and     dh, dh
        jz      short loc_14179     ; 跳過額度用完 → 直接檢查
        call    sub_1ECE0
        cmp     al, 40h
        jnb     short loc_14179     ; ≥ 40h → 不跳過
        dec     dh                  ; ★ 額度只在真的跳過時才減
        jmp     short loc_1418E
loc_14179:
        test    byte ptr [di], 4    ; 位元 2 ＝ 委任
        jz      short loc_1418A     ;   沒設 → 不收
        cmp     byte ptr [di+23h], 8
        jnb     short loc_1418A     ;   Stage ≥ 8 → 不收
        mov     [di+20h], cl        ; 收：意圖 ← 目標
        mov     [di+23h], ch        ;     Stage ← ch
loc_1418A:
        dec     dl                  ; ★ 收不收都減
        jz      short locret_14193  ; ★ 歸零就整支結束
loc_1418E:
        add     di, 40h / jmp loc_14160
```

⇒ **`dec dl` 在匯流點**。掃到一支「人在這個據點、活著、沒被跳過」的軍團，
不論收不收，額度就少一個。

## 2. 這一條的代價

拿「收到幾支」當額度的話，一支 `Stage >= 8` 的守軍不算數，迴圈會往下
找第二支，於是多取一次亂數（`skip` 還有額度時每一支候選都要抽），
而且**調走了原版不會動的那一支**。

同局面拍 2,455 的據點 129：`occ=2`、`threat=1`、`want=1`、`skip=1`。
兩支守軍是軍團 19（`Stage 8`）與軍團 20。原版掃到軍團 19 就結束。

## 3. 改了哪幾支

| 檔案 | 改什麼 |
|---|---|
| `internal/rules/threat/threat.go` | `Dispatch` 的 `want` 改成每掃過一支候選就減，收不收都一樣 |
| `internal/state/corpsorder.go` | `resupplyCorps`（`sub_14499`）補上 `recalcCorps`——它也是 `sub_16FD2` 的呼叫者，計時寫 1 讓 Stage 3 的分派提前一個週期跑（見 §4）|

## 4. 為什麼補兵也要重算

`sub_14499` 在 `sub_16FD2` 的四個呼叫者裡（[`../re/30`](../re/30-corps-formation-ui.md) §5）。
少了那一步，Stage 3 會多停 16 拍——**剛好夠讓求援的調兵把它掃走**
（`+0x23 < 8` 才調得動）。同局面的軍團 19：

| | Stage 9 | Stage 3 | Stage 8 |
|---|---|---|---|
| 原版 | 拍 2,415 | 拍 2,439 | **拍 2,447** |
| remake（修之前）| 拍 2,415 | 拍 2,439 | 沒走到——2,455 被調走 |

## 5. 驗證

`TestDispatchWantCountsCandidatesNotPicks`：一支 `Stage >= 8` 的候選吃掉
`want` 就結束（正對照：額度夠時掃得到第二支；亂數次數也釘住）。

`tools/parity_ck.sh` 十七個檢查點裡**只剩拍 2,430 的一個 byte**
（軍團 50 的朝向）。逐拍取數不一致 178 → **117 / 5,480（2.1%）**，
第一個分歧從拍 2,455 推到 **3,304**。

## 6. 未解

- 原版**沒有比勢力**（`sub_14155` 只看節點欄、活著、位元 2、Stage），
  所以別勢力的軍團站在同一格時原版會收它。remake 的 `Ready` 含勢力檢查
  ⇒ 不收，但一樣消耗額度。那個局面（對峙中）還沒在對拍裡出現過。
