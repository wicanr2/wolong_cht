# 180 — 求援調兵的 `want` 數的是候選不是選中

**狀態：CONFORMED。** `sub_14155` 的 `dec dl` 在 `loc_1418A`——位在
「收」與「不收」兩條路**匯流之後**。所以 `want` 是「掃過幾支候選」，
不是「收到幾支」；一支 `+0x23 >= 8` 的守軍會吃掉一個額度，
額度歸零就整支結束。

- 日期：2026-11-16；初版為 2026-09-10
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

## 6. 歷史差異與抽驗範圍

- 2026-09-10 記錄的差異：原版**沒有比勢力**，Go 額外檢查勢力，所以不同勢力同格時未分派。
  2026-11-16 已依下節原版矩陣修正；仍不以受控同格場景宣稱自然遊玩已重現此極端局面。

### 2026-11-16：同格不同勢力的分派補驗

**狀態：CONFORMED。** 本節只修 `sub_14155` 的候選條件，不改求援次數、編成額度或 UI。
沿既有 `c-tick` 的 IDA Pro 9.4 匯出重新核對 63 bytes，並以 dosgolem 原始指令
檢查同格不同勢力、非委任、Stage 8、死亡及行軍中軍團的分派與 RNG。
兩側先固定相同 258-byte 亂數狀態，再比較完整軍團記錄與呼叫後亂數狀態。
原始輸入為松崗 DOS/V `KI.EXE`，SHA-256 同 [re/44](../re/44-threat-and-reinforcement-ai.md)。
已證實：原版沒有比較軍團勢力。四劇本、八種軍團配置、三種跳過額度、
兩種需求與兩份固定完整亂數狀態，共 384 案；舊 Go 有 240 案軍團記錄不符，
RNG 全同。僅移除額外勢力限制的可丟棄候選實作，完整軍團記錄與 258-byte RNG 全等。
正常玩家的極端同格局面未重現，不把受控局部測試稱為正常玩家路徑。
本機[證據審查](../../workplace/scratch/ai-dispatch-20261115/relief-ready.json)保存輸入與工具雜湊、
矩陣、測試入口修正與限制。據此進入正式實作，保留其餘委任與 Stage 條件。

正式修改限 `internal/state/strategy.go` 的 `dispatchGarrison`。存檔資料與型別不變，
`TestReliefDispatchDoesNotFilterFaction` 保留外軍委任／非委任／Stage 8 三案，
並確認被跳過的候選仍耗用額度。局部原版驗證覆蓋原始軍團 bytes → `LoadBlock` →
分派 → `Bytes`；正常 GUI 與原版受控矩陣分別記錄，不宣稱同狀態 UI 對拍。

全專案冷測、受影響 vet 與建置通過。正常曹操滑鼠 62 項／9 圖通過編成、停手行軍、
途中存讀、自然野戰、退卻與戰後存讀。這是既有玩家路徑回歸，沒有注入外軍同格狀態。
本機[整合收據](../../workplace/scratch/ai-dispatch-20261115/relief-integration.json)保存來源、
binary、固定亂數、原版與正式矩陣、正常 GUI 及限制。

### 2026-11-16：比較範圍勘誤與尾端軍團槽

前片所稱「完整軍團記錄」需限縮：384 案實際只比較前 120 個槽，
全部有變動的記錄都在前八槽，因此原結論對這些候選仍成立，不能稱已驗尾端。
本片重跑原 384 案，再將唯一候選分別放在槽 120–126，共 391 案；
比較 128 個原始槽的全部 8,192 bytes 及完整 258-byte RNG，正式 Go 全等。
第 128 槽是原始暫存槽，Go 正式軍團仍為 127 支，不改資料格式。
原收據保留，追加[本片結果](../../workplace/scratch/ai-relief-20261116/dispatch-current-result.json)與
[整合收據](../../workplace/scratch/ai-relief-20261116/integration.json)，不重寫舊測試輸出。
