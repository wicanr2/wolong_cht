# 08 — `wlgame` 正常編成／行軍路徑

**狀態：真實 `SINARIO.DAT` 下的編成、行軍、城兵攻城，以及敵方 AI 軍團遭遇選單
都已由鍵盤正常操作驗收；使用者提供的 `SAVE.DAT` 另有獨立回放證據，但第 1 槽
時鐘欄位為零，不能當正常開局／正常戰略 AI 的時序 oracle。**

- 日期：2026-08-09
- 執行環境：Docker／Xvfb，`demonwinter-go:latest`，Go cache 使用專案既有
  `wl-gomod`／`wl-gobuild` volume
- 原始輸入：`workplace/orig/dosv/SINARIO.DAT`、`workplace/orig/dosv/SAVE.DAT`、
  `workplace/eten/`，唯讀掛載
- 啟動：`cmd/wlgame`，正常 AI 回放使用 `-seed 17` 固定亂數；沒有使用
  `-open-battle`、`-open-battle-choice`、`-open-form` 或 `-open-corps`
- 初始狀態：劇本 1、玩家勢力 0（曹操）、196 年 4 月 1 日、`speed=0`

## 1. 正常操作序列

先按 `A` 開啟軍隊編成，按 `Enter` 兩次完成武將一覽表的兩段式選取。
開局預備兵只有騎馬 400、弓兵 600、步兵 1,000；畫面預設六槽全騎馬，
直接按 `Enter` 會拒絕「需要 6,000」——這是資源檢查，不是測試捷徑。

依現有資源完成編成的操作是：

```text
3                    # 第 1 槽改成一槽步兵
Down + Space × 5     # 第 2–6 槽清空
Enter                # 編成
```

畫面事件顯示「曹操 編成完畢」。接著按 `M`，在軍團一覽表按 `Enter`
兩次選軍團，再在目的地一覽表按 `Enter` 兩次選第一個距離排序項目
「虎牢關」。最後按 `=` 把速度由 0 調成 1，讓世界時鐘在行軍期間推進。
若要驗證攻城，在目的地一覽表從第一列往下 10 列選「汝南」（袁術據點），
再把速度逐步調高，等待道路上的軍團抵達。

## 2. 結果與證據

| 階段 | 證據 | 結果 |
|---|---|---|
| 編成完成 | [`wlgame-normal-formed.png`](../images/wlgame-normal-formed.png) | 真實預備兵限制下顯示「曹操 編成完畢」 |
| 目的地選取 | [`wlgame-normal-destination.png`](../images/wlgame-normal-destination.png) | 目的地一覽表按距離排序，第一項為虎牢關 |
| 行軍中 | [`wlgame-normal-march.png`](../images/wlgame-normal-march.png) | 日期由 196/4/1 推進至 196/4/2，事件顯示「曹操 向 虎牢關 行軍」 |
| 正常攻城 | [`wlgame-normal-garrison.png`](../images/wlgame-normal-garrison.png) | 無 `-open-*` 旗標攻打汝南；畫面事件顯示「曹操 對 城兵　攻下 汝南」 |
| 正常敵方 AI 遭遇 | [`wlgame-ai-normal-encounter.png`](../images/wlgame-ai-normal-encounter.png) | 固定種子 17；無 `-open-*` 旗標；事件佇列接入後於 196/6/28 顯示「呂布 對 曹操／攻城／戰鬥指揮／委任」 |
| 存檔遭遇回放 | [`wlgame-save-replay-choice.png`](../images/wlgame-save-replay-choice.png) | 將使用者提供的 DOS/V `SAVE.DAT` 第 1 槽複製成可寫 overlay；無 `-open-*` 旗標進入「張飛 對 許褚／攻城」的「戰鬥指揮／委任」選單 |

上述 PNG 都由 Docker/Xvfb 產生，檔案擁有者是目前 UID/GID `1000:1000`；
原始 `SINARIO.DAT` 與 `SAVE.DAT` 均未寫入。正常劇本切片證明編成／行軍以及
「城裡只有城兵時的自動攻城」玩家輸入接縫；依原版規則，這一類空城攻城不進
戰鬥指揮畫面。存檔回放則補足無除錯旗標的遭遇選單畫面接縫，但不是正常開局時序證據。

## 3. 正常敵方 AI 遭遇路徑

這次回放在真實劇本開局後只使用正常鍵盤輸入：

```text
A                         # 軍隊編成
Enter × 2                 # 選曹操
3；Down + Space × 5；Enter # 一槽步兵
M；Enter × 2              # 選軍團
Down × 22；Enter × 2      # 距離排序第 22 列：濮陽（據點 56）
= × 64                   # 正常調速，讓月結與敵方 AI 推進
```

`-seed 17` 只替驗收固定亂數來源，沒有跳過月結、宣戰、編成、行軍或遭遇判定。
真實劇本下，敵方勢力 13（呂布）在月結評估後編成軍團，沿現有 MMAP 道路向
濮陽行軍；當它與曹操軍團在濮陽發生攻城接觸時，`World` 正常建立
`EncounterChoice`，`wlgame` 顯示「戰鬥指揮／委任」。畫面截圖時日期為
196 年 6 月 28 日，沒有使用任何 `-open-*` 驗收旗標，也沒有傳送或直接改寫軍團座標。

### 存檔回放的採信界線

Docker 回放使用 `-save-file` 指向工作樹外的複本，未使用任何 `-open-*` 旗標。
`state.LoadScenario` 接受該 overlay，記錄的起始狀態是「劇本 1、0 年 0 月 0 日」，
600 幀截圖時已顯示「0 年 1 月 1 日」以及「張飛 對 許褚」。因此可證實：

- 儲存檔 overlay 能載入含軍團的結構化狀態；
- 正常遊戲迴圈能從該狀態進入遭遇選單，並顯示戰鬥指揮／委任選項；
- 不能據此宣稱該 `SAVE.DAT` 是可供正常開局或時鐘 parity 使用的原版 oracle；
  還需要時鐘欄位非零、來源可追溯的原版存檔或由正常戰略 AI 重新產生的狀態。

## 4. 尚未完成

- 原始劇本開局本身沒有敵方軍團；必須讓正常月結 AI 先宣戰、編成與行軍，
  本文件的固定種子回放已驗證這條路徑。
- 目前找到的 `SAVE.DAT` 雖能回放遭遇選單，時鐘欄位仍是零，不能證明敵方軍團
  是由正常戰略 AI 在可對拍的時間線上產生。
- 尚未完成的是原版有效時鐘／原版存檔與 remake 同狀態的逐畫面對拍；事件 1／8／13
  的 queue handler 已接入，但事件 2–7、9–12、多軍團請求與完整原版行軍狀態機仍未達
  parity。這些不再阻塞本文件已驗證的正常 AI 遭遇接點本身。
