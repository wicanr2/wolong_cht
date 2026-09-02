# 58 — 接上兵的戰力之後重跑戰術對拍：野戰回到 0.05%，攻城沒有退步

**狀態：兩組都重量完，沒有回歸。** 野戰九區裡**七區逐像素相同**、
`field` 95 px（＝原版錄影裡的滑鼠游標）、`sb-minimap` 32 px；
攻城三個取樣點 `field` 0.85%／0.84%／2.00%，前兩點與 2026-08-27 同值。
這一輪挖出兩個**驗收路徑**的缺陷，都與規則層無關：
`-shot-frames` 推不動戰場（[`../spec/91`](../spec/91-tactical-parity.md) §6.1）、
驗收捷徑推完才武裝開場喊話（[`../spec/117`](../spec/117-fixture-arms-duel-before-stepping.md)）。

- 日期：2026-09-02
- 為什麼跑：[`../spec/115`](../spec/115-soldier-power.md) 把兵的戰力從
  軍團士氣（100–200）接回原版的 `((統率 + 適性) × 3 + 兵種係數) ÷ 4`（約 18）。
  傷害少一個數量級、命中率從 100% 降到約 59%，**戰鬥的節奏整個變了**——
  先前的對拍數字必須重驗
- 原版側：沿用受控擷取（`workplace/promo-live/parity-battle4/`、
  `parity-field14/b0.png`、`probe-march/e10.png`），
  640×480 由 `tools/parity_crop.py` 切成 640×400，**沒有重跑 DOSBox-X**
- remake 側：`tools/parity_shot.sh` 用現在的原始碼重截
- 比對：`tools/py.sh tools/parity_diff.py … --regions tactical`

## 1. 結果

### 1.1 野戰（`SAVE-FIELD.DAT`，夏侯惇 39 對呂布 35，開戰第 52 拍）

```
tools/parity_shot.sh out.png -direct -scenario 0 -player 0 -seed 1 \
  -save-file workplace/parity/SAVE-FIELD.DAT -load-slot 0 \
  -open-battle -siege-corps 39,35 -battle-steps 52 -shot-frames 1
```

| 區 | 這一輪 | 2026-08-25（[`43`](43-field-battle-parity.md)）| 判定 |
|---|---:|---:|---|
| `field` | **95 px（0.05%）** | 95 px（0.05%）| NEAR ＝ 原版錄影的滑鼠游標，消不掉 |
| `sb-minimap` | **32 px（0.16%）** | 0.20% | NEAR ＝ 部隊點的時刻差 |
| 其餘七區 | **0 px** | 0 px | **PASS** |

⭐ **一個像素都沒有退步。** 開場對白（呂布的挑戰、框的版面、代入的名字）
與地形層逐像素一致。

### 1.2 攻城（`SAVE-E.DAT`，張遼 81 攻／夏侯惇 39 守，據點 82）

```
tools/parity_shot.sh out.png -direct -scenario 0 -player 0 -seed 7 \
  -save-file workplace/promo-live/parity-battle4/SAVE-E.DAT -load-slot 0 \
  -open-siege -siege-node 82 -siege-corps 81,39 -battle-steps S -shot-frames 1
```

| 區 | 第 160 拍 | 第 240 拍 | 第 300 拍 | 2026-08-27 |
|---|---:|---:|---:|---|
| `field` | 0.85% | 0.84% | **2.00%** | 0.86%／0.81%／0.84% |
| `sb-minimap` | 1.64% | 1.64% | 1.21% | 1.64% |
| `sb-enemy` | 12 px | 12 px | 22 px | 14 px |
| `bottom` | 2 px | 2 px | 2 px | 2 px |
| `sb-title`／`sb-self`／三個純美術區 | **0 px** | **0 px** | **0 px** | 0 px |

前兩個取樣點與 2026-08-27 同值（差在小數第二位），
**第 300 拍從 0.84% 升到 2.00%**。把差異攤成 32×32 的方塊之後
（`workplace/parity/difflocate.py`），3,800 px 平鋪在 x 192–420、y 64–160
那一片，**沒有任何一塊佔到十分之一**——那是城門一帶的兵，不是某個
畫壞的元件。成因是局面不等價：戰力接對之後兵挨十一下才倒（先前一擊必殺），
第 300 拍城門下站著的人比原版那一格多。

⭐ **`sb-enemy` 反而從 14 px 降到 12 px。** 那一格畫的是對方大將的體力條，
而大將體力的算式（`max(70, (武力 × 4 + 50) × 士氣 ÷ 100)`）也是這一輪
接上的（[`../spec/115`](../spec/115-soldier-power.md) §1）——先前那一條畫的是軍團士氣。

## 2. ⛔ 第一個坑：`-shot-frames` 推不動戰場

第一輪照 [`../spec/91`](../spec/91-tactical-parity.md) §6 當時寫的
`-shot-frames 160`／`240`／`300` 量，得到 `field` **4.37%／4.37%／4.35%**——
數字變差、三個取樣點又幾乎一樣，看起來就像規則層改壞了。

**不動的地板要去裁圖，不要再掃步數**（[`49`](49-parity-retest-20260827.md) §3 的教訓）。
裁出來一看，7,724 px 裡有 6,808 px 集中在 x 224–479、y 0–31，
那就是**門強度視窗整塊**：remake 那一張根本沒畫它。

規則層是好的——直接跑 `Battle.Step()` 量，第 148 幀城壁開始掉、
條從那時起一直亮著：

| 戰場幀 | 條 | 最低城壁耐久 | 側 0（攻方）|
|---:|---|---:|---|
| 100 | false | 1,660 | X 56..60 |
| 140 | false | 1,660 | X 35..57 |
| **160** | **true** | 1,632 | X 28..57 |
| 300 | true | 1,334 | X 28..43 |

差別在**旋鈕**：`-battle-steps` 直接呼叫 `Battle.Step()`，
`-shot-frames` 推的是畫面，戰場靠節流器跟著走，戰術速度 2 之下
**每 6.59 個畫面才一步**（`speed.Steps`：`2913 / (2 × 16 × 600)`）。
`-shot-frames 300` 只推到戰場第 45 幀，門強度條當然還沒亮。
規格已改（[`../spec/91`](../spec/91-tactical-parity.md) §6.1）。

## 3. ⭐ 第二個坑：驗收捷徑推完才武裝開場喊話

野戰這一組先卡在另一件事：自然流程（`-shot-frames 400`）**停在遭遇訊息上**
不會進戰場——那是 2026-08-29 起的預期行為
（[`../spec/105`](../spec/105-encounter-goes-straight-to-battle.md) §4），
截圖模式沒有人按掉那則訊息。實測第 400、500、700、900、1,100 幀
**五張 PNG 逐位元組相同**，`51e64bf4…`。

改走 `-open-battle -siege-corps 39,35 -battle-steps 52` 之後八區立刻對上，
`field` 卻是 11.24%，差的整塊是呂布的挑戰對白框——remake 那一張沒有框。

⭐ **判準是拿一個「應該會改變結果」的旋鈕轉一圈。** 挑戰側由氣勢比較決定，
平手時看亂數尾（[`../spec/80`](../spec/80-duel-opening.md) §6），
所以換 seed 應該會換人挑戰。實測 seed 0–15 **十六個值一模一樣**：

```
seed 0   field 11.24%     …（中間十四個同值）…     seed 15  field 11.24%
```

一個吃亂數的機制對亂數完全免疫，只有一種解釋：**它沒有被啟動。**
成因與修法在 [`../spec/117`](../spec/117-fixture-arms-duel-before-stepping.md)——
`stageEncounter` 推完 N 拍才呼叫 `startBattleTalk`，而 `SetDuelInput`
在那支裡面，單挑的開場只有 50 tick。改成先武裝再逐拍推之後就是 §1.1 的 0.05%。

> **教訓寫成規則**：**驗收捷徑與正常迴圈的差別，會以「被驗收的東西不見了」
> 的形式出現。** 這是 [`49`](49-parity-retest-20260827.md) §2（`-shot` 把音效
> 關掉，於是那一格印「未接入」）的同一族。捷徑省掉的不是無關的雜項，
> 而是被驗收畫面的一部分。

## 4. 未解

| 項目 | 現況 | 下手點 |
|---|---|---|
| 攻城 `field` 的 0.84% 地板 | 局面不等價：原版擷取是 5月20日的張遼軍攻許昌，存檔是 5月10日（[`51`](51-siege-deadlock.md) §2）| 要對到 0 px 得有「存檔與影格出自同一次擷取」的攻城素材，同 [`52`](52-siege-timeseries-parity.md) 那一組 |
| 第 300 拍的 2.00% | 兵的密度不同（§1.2）| 取樣點要用局面條件挑，不是寫死步數（[`../spec/91`](../spec/91-tactical-parity.md) §6）|
| 野戰走自然流程 | 遭遇訊息擋住截圖（§3）| 要一個「訊息自動按掉」的驗收旗標；現在靠 `-open-battle -siege-corps` 繞過 |
| 原版側沒有重跑 | 用的是 08-16／17／24 的擷取 | 要重跑得先建 `wolong-dosboxx`（`docker/dosboxx/Dockerfile`）|
