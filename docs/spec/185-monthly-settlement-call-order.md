# 185 — 月結內部的呼叫順序：政略在災害之前

**狀態：CONFORMED。** `sub_15358` 的呼叫序列裡，政略評估（`sub_12BD9`）排在
撥款請求與災害之前；remake 把災害的骰子提到最前、把政略推到最後，
於是同一拍內的亂數**次數相同而值全部錯位**。

逐拍取數比對只比「這一拍取了幾次」，兩邊都是 647 次——
**形狀完全正常，而序列早就分岔了**。

- 日期：2026-09-10
- 出處：`KI.EXE`（松崗 DOS/V，`workplace/ida/dosv/KI.EXE.i64`，
  SHA-256 `6deb8e9c…e5dd62`）`sub_15358`（IDA 線性 00015358–000153C5）、
  `sub_12286`（00012286–000122DA）、`sub_122DB`（000122DB–0001237D）、
  `sub_12FBF`（00012FBF–0001300D）、`sub_1301C`（0001301C–0001304D）
- 推論等級：confirmed（四支逐條讀出，並與原版 `-watch` 攔到的
  `sub_1ECE0` 返回位址序列逐筆對上）
- 實測：[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §45
- 相關：[`176`](176-event-queue-cursor-not-restored.md)（同一條佇列的游標）、
  [`182`](182-monthly-settlement-before-hourly.md)（月結與每小時的先後）、
  [`../re/07`](../re/07-monthly-settlement.md) §21

## 1. 原版做什麼

```asm
sub_15358:
        …22 個勢力的迴圈（sub_1563B／sub_153C6／sub_15609／sub_15828）
        call sub_15695
        call sub_1585F          ; 在野武將／俘虜
        call sub_155A6
        call sub_12BD9          ; ★ 佇列壓縮 ＋ 政略（遷都／合作／停戰／宣戰）
        call sub_15715          ; 內政官撥款請求
        call sub_1578F          ; 外交官撥款請求
        call sub_122DB          ; ★ 暴風雨
        call sub_12286          ; ★ 逐據點火災／暴動
        call sub_157FE
        …來月設定生效（4 個 word）
        mov  al, 0Eh / call sub_15E80
```

### 1.1 `sub_12286`：骰子與事件寫入是交錯的

```asm
        mov si, 840h / mov cx, 0C0h            ; 192 座，runtime 據點表
loc_12297:
        call sub_1ECE0 / cmp al, 18h / jnb loc_122B4     ; 火災閘一
        call sub_1ECE0 / and al, 3Fh
        cmp al, [si+11h] / jb loc_122B4                  ; 火災閘二（防災值）
        mov dx, si / mov ax, 10Ch / mov bl, 0FFh
        call sub_12FBF / jmp loc_122CF                   ; ★ 火災中了就不骰暴動
loc_122B4:
        call sub_1ECE0 / cmp al, 18h / jnb loc_122CF     ; 暴動閘一
        call sub_1ECE0 / and al, 3Fh
        cmp al, [si+10h] / jb loc_122CF                  ; 暴動閘二（上昇值存值）
        mov dx, si / mov ax, 20Ch / mov bl, 0FFh
        call sub_12FBF
loc_122CF:
        add si, 20h / loop loc_12297
```

**每一座據點骰完就當場寫佇列**，不是先骰完 192 座再統一寫。
`sub_12FBF` 在 `bl == 0FFh` 時自己會取一次亂數當搜尋起點，
所以那一次取數夾在兩座據點的骰子之間。

### 1.2 `sub_122DB`：暴風雨的槽位提示**不是** `0FFh`

```asm
        call sub_1ECE0 / test al, 1 / jz retn            ; 50%
        call sub_1ECE0 / cmp al, 0C0h / jnb retn         ; 據點編號 < 192
        bx = al >> 3 ← 中心據點
        cmp byte ptr [bx+848h], 0C0h / jnb loc_12332     ; 靠海就跳過加骰
        call sub_1ECE0 / test al, 1 / jz retn            ; 內陸再 50%
loc_12332:
        push bx
        call sub_1ECE0 / and al, 7 / add al, 8
        shl al, 1 / shl al, 1 / mov bl, al               ; ★ bl ＝ (亂數&7 ＋ 8) × 4
        mov al, 0Bh / xor ah, ah / xor dx, dx
        call sub_12FBF
        pop bx
        jb  loc_1237A                                    ; ★ 寫不進去就不設範圍
        …word_10D22／10D24／10D26／10D28 ＝ 暴風雨矩形
```

兩件事跟 remake 現況不同：

1. **`bl` 是算出來的槽位提示**（32、36、…、60），走 `sub_12FBF` 的
   `bl × 4 ＋ 游標`那一支，**函式內部不會再取亂數**。
   remake 傳 `0xFF`，於是那一次亂數改由 `queueEvent` 內部取——
   次數相同、值的去向完全不同。
2. **入佇列失敗（CF＝1）就不設暴風雨範圍。** remake 先設 `stormArea`
   再入佇列，而且不看回傳值。

### 1.3 `sub_12FBF` 與 `sub_1301C` 的分工

| | `sub_12FBF` | `sub_1301C` |
|---|---|---|
| `bl == 0FFh` | 起點 ＝ `(亂數 & 7Ch) ＋ 游標` | 沒有這一支（`bl` 一律當數值用）|
| 否則 | 起點 ＝ `bl × 4 ＋ 游標` | 同左 |
| 起點上限 | `>= 100h` 直接失敗（STC）| 不檢查 |
| 搜尋範圍 | `< 100h`（前 64 格）| `< 400h`（全部 256 格）|
| 回傳 | CF ＝ 0 成功／1 失敗 | 不回報 |

remake 的 `queueEvent`／`queueFullEvent` 形狀已經對，這一份不動它們。

## 2. remake 現況與要改的地方

`World.settleMonth`（`internal/state/state.go`）目前的順序是

```
經濟 → RollCityDisaster ×192 → RollStorm → 在野武將 → 壓縮佇列
     → 撥款請求 → 事件 11 → 事件 12 ×N → runStrategicAI → 事件 10
```

要改成原版的

```
經濟 → 在野武將 → 壓縮佇列 ＋ runStrategicAI → 撥款請求
     → RollStorm ＋ 事件 11 → 逐據點（骰子 ＋ 事件 12 交錯）→ 事件 10
```

| 改動 | 對應原版 |
|---|---|
| `runStrategicAI` 移到撥款請求**之前** | `sub_12BD9` 在 `sub_15715` 之前 |
| `RollStorm` 移到撥款請求**之後**，與事件 11 綁在一起 | `sub_122DB` |
| 事件 11 的 `slotHint` 改成 `((亂數&7)+8)*4`，並改成先入佇列成功才設 `stormArea` | `loc_12332`–`loc_1237A` |
| `RollCityDisaster` 移到最後，**逐座骰完就當場入佇列** | `sub_12286` 的迴圈 |

## 3. 怎麼驗

1. `tools/parity_pace_diff.py`：月結那一拍（拍 2,967）兩邊的**來源序列**
   要能逐筆對上（經濟 192 → 政略 → 災害），不是只有總數 647 相同。
2. 原版側 `WOLONG_DOSGOLEM_WATCH=12FBF` 攔到的 8 筆事件
   （遷都 8／9／15／17／18、勢力 0 宣戰 13、勢力 2 合作、勢力 2 天災清除）
   要與 remake 的 `OnEvent` 軌跡同序同值。
3. `tools/parity_ck.sh` 的**佇列**那一欄要歸零——
   在 `tools/orig_snapshot.py` 補上 `peek:2514` 之後它才是活狀態
   （[`../playtest/119`](../playtest/119-rng-pace-comparison.md) §44）。

## 4. 未解

- `sub_15695`／`sub_155A6`／`sub_157FE` 在 remake 的對應位置還沒逐條核對，
  這一份只動已證實會取亂數的四支。
- **拍 4,911：同一場攻城戰的勝負判定相反。** 戰前兩邊逐 byte 相同、
  逐拍取數的前 17 筆也逐筆對上，而原版的軍團 1 掉 116 兵（敗方）、
  remake 只掉 20 兵（勝方）並攻下據點 2。落點在 `sub_15130` 的
  `cmp cx, dx`——`sub_15285`（基礎戰力，守城方多加一次城兵數）
  或 `sub_152D7`（將領修正）。
- **拍 5,200** 仍然雪崩，要等上面那一項解掉再看。
