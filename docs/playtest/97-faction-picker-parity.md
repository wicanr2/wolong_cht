# 97 — 22 勢力選擇視窗：四區 0 px，抓到兩個缺陷

**狀態：通過。** 縮小地圖圖例右半格開的那個視窗（原版熱區 `0x17` →
`sub_15AD1`）逐區比對，`banner`／`command`／`map`／`faction` **全 0 px**，
`minimap` 只剩原版自己畫的滑鼠游標（82 px）。抓到兩個缺陷。

- 日期：2026-09-06
- 規格：[`../spec/140`](../spec/140-status-message-box.md) §1.1（#4）、
  [`../re/31`](../re/31-faction-picker-screen.md)（版面與顏色規則）
- 原版側：`WOLONG_DOSGOLEM_GAMEDIR=dosgolem/root-noclouds tools/dosgolem.sh …
  "…;runto:11CD0;sclick:416,15;steps:400000;stap:600,175;steps:900000;shot:orig-picker"`
  （`416,15` ＝ 橫幅第三格開關 ＝ 縮小地圖；`600,175` ＝ 圖例右半格）
- remake 側：`-save-file workplace/dosgolem/root-noclouds/SAVE.DAT -load-slot 0
  -open-faction-picker -cam 0,0
  -fixture-when clock:196/4/17/6 -shot-when clock:196/4/17/6`
  ⚠ **受控存檔不能省**（[`../spec/147`](../spec/147-controlled-parity-save.md)）：原版側跑的是 `root-noclouds`，remake 只用 `-direct` 從劇本跑到同一時刻會有雲、也會有軌跡分歧。這一行本輪補上，補之前照著跑對不出文中的數字
  （`-open-faction-picker` 本輪新增）

## 1. 一路修下來的數字

| 改了什麼 | `minimap` | `faction` |
|---|---:|---:|
| 第一次拍 | 190 | 170 |
| 玩家那一列的字色（§2）| 190 | **0** |
| 視窗開著時不畫視野框（§3）| **82** | **0** |

## 2. 玩家那一列畫成紅的，原版是藍的

原版的著色是**兩次連續判斷**（[`../re/31`](../re/31-faction-picker-screen.md) §2）：

```asm
ah = 90h                              ; 一般
cmp al, cs:byte_10CFF / jnz +         ; 是玩家自己 → ah = 9Ah（前景 A ＝ 色 10 紅）
cmp al, cs:byte_198A7 / jnz +         ; 是目前選中 → ah = 93h（前景 3 ＝ 色 3 藍）
```

**後面那次覆蓋前面**。而開局時圖例第二格盯的就是玩家自己
（`byte_198A7` ＝ 玩家勢力），兩個條件同時成立——所以原版畫的是**藍的**。
remake 寫成 `switch`，把「是自己」放在前面，於是畫成紅的。

⭐ **這種 bug 兩邊都自洽**：兩個規則各自都對，錯的只有**誰先誰後**，
而且只在「兩個條件同時成立」時才顯現。逐像素對拍是唯一看得到的地方。

## 3. 視野框：原版蓋掉之後就不再回來

remake 在縮小地圖上畫視野框（原版 `sub_15C58` → `sub_196ED` 的那張點陣），
但**這個視窗開著時原版沒有那個框**。

原因在建構順序（[`../re/31`](../re/31-faction-picker-screen.md) §1.3）：
`sub_15A3A` 把縮小地圖整個重鋪一次——底圖 → 上方兩個勢力名 → 外框 →
192 個據點 → `sub_15C58`——**這份清單裡沒有視野框**，所以框被底圖蓋掉之後
不再畫回來。與清單、彈出選單「關掉不擦」是同一族現象的反面：
**這一次是重畫把東西擦掉了，而且沒有人補回來**。

## 4. 順帶：狀態列 #4 接上了

`sub_15AD1` 進 `sub_15AFC` 之前 `mov cx, 4 / call sub_18853`，
掛的是 #4「以滑鼠的右鍵回復。」——**27 個帶索引的狀態列呼叫點到此全部接上**
（[`../spec/140`](../spec/140-status-message-box.md) §3）。

## 5. 未解

| 項目 | 現況 |
|---|---|
| 原版擷取裡的滑鼠游標 | 82 px，停在剛點的圖例上 |
| 只重畫兩列 | 原版換選中時只重畫舊的與新的那兩列，remake 整張重畫。**視覺上沒有差異**，但要知道原版沒有全畫面刷新（[`../re/31`](../re/31-faction-picker-screen.md) §2.2）|
