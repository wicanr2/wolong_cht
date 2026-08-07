# 04 — 受阻：PC-98 oracle 動不了遊戲內的滑鼠游標

**狀態：卡住。已試六種組合全部失敗，停手記錄。**

- 日期：2026-08-07
- 依據 `rulebook/41`（特例越補越多、「再差一點」卡很久 → 停下重想架構）
  與 `CLAUDE.md` §8 第 10 條（猜三次沒中就停手）

## 1. 現象

`tools/dosboxx.sh` 可以把 PC-98 原版開機、播完開場、走到 `NEW GAME` 對話框。
**點擊事件進得去**（會觸發選項），**但遊戲內的滑鼠游標從頭到尾停在同一個位置**
（約畫面座標 `(325, 200)`）。

結果就是：每次點擊都落在游標所在的那一格，也就是 `NEW GAME` 的 `NO`，
於是一路跳到 `LOAD DATA`，**進不了遊戲本體**。

## 2. 試過的五種組合（全部失敗）

| # | 做法 | 結果 |
|---|---|---|
| 1 | `xdotool mousemove --window $WID x y` | 座標被當成 root 座標。視窗實際在 `(320,312)`，所以**點擊落在視窗外**。這是真 bug，已修 |
| 2 | `autolock=true` ＋ root 絕對座標 | 點擊進得去，**游標不動** |
| 3 | `autolock=true` ＋ 「先灌 −2000,−2000 歸零再走相對位移」 | 游標不動 |
| 4 | `autolock=false` ＋ root 絕對座標 | 游標不動 |
| 5 | 鍵盤（`Up` / `Return` / `space`） | 對話框完全不反應 —— 這個介面是純滑鼠的 |
| 6 | **`mouse_emulation=never` ＋ `autolock=false`** | 游標**還是不動** |

### 第 6 項的來由（找到了設定，但沒解決）

`/usr/share/dosbox-x/dosbox-x.reference.full.conf` 裡有：

```
#   mouse_emulation: When is mouse emulated ?
#                      integration: when not locked
#                      locked:      when locked      ← 預設
#                      always:      every time
#                      never:       at no time
#                      If disabled, the mouse position in DOSBox-X is exactly
#                      where the host OS reports it.
```

**這本來是很好的假說**：預設 `locked` 配上我用的 `autolock=false`，
等於「從來不模擬」，正好解釋游標不動。改成 `never`（官方說明白寫
「位置就是 host 回報的位置」）照理應該直通。

**實測沒有用。** 游標仍停在原位。

> 記下這一項是因為**它的推理看起來完全正確**——
> 找到設定、讀懂說明、對上現象、改了、還是不動。
> 下一輪如果又想到「應該是 mouse_emulation」，這裡已經試過了。

第 5 項本身是有用的資訊：**遊戲的選單不吃鍵盤**。
`KI.EXE` 的 `ERROR: Mouse driver not install ?` 與 PC-98 版另附的
`YNMOUSE.COM` 都印證這一點。

## 3. 為什麼要記下來

這一輪在同一個問題上花了**六次執行、約 25 分鐘**（每次要重播 90 秒開場）。
再試更多變體不會變便宜，而且沒有新的假說 —— 符合停手的訊號。

**不寫下來的話，下一輪會從第 1 種再試一次。**

## 4. 下一輪的候選（依可能性排，都還沒試）

1. **PC-98 的滑鼠是匯流排介面（8255），不是 PS/2 或 INT 33h。**
   DOSBox-X 可能需要另外的設定才會模擬它。
   查 `dosbox-x` 的完整設定範本裡與 `mouse` 相關的項目 ——
   本輪查失敗（無 X 環境下 `dosbox-x -c config -wc` 會掛住），
   下一輪在有 Xvfb 的容器裡再查一次。
2. **DOSBox-X 的存檔狀態（save state）。** 手動玩一次到遊戲開始，存狀態，
   之後每次從那裡續 —— 順便省掉 90 秒開場的重播成本。
3. **改用松崗 DOS/V 版當 oracle。** 那一版是標準 PC 滑鼠（INT 33h），
   DOSBox 的支援成熟得多。代價是要先過防拷（`docs/playtest/01` §2）。
4. `sensitivity` 設定、mapper 檔、`-c` 啟動指令。
5. **用真的 X server 而不是 Xvfb。** 六次失敗橫跨兩種 `autolock`、
   絕對與相對定位、以及 `mouse_emulation=never` —— 共同點是**都在 Xvfb 下**。
   Xvfb 的指標處理（沒有真實輸入裝置）是目前最大的共同嫌疑。
   這一項升為第一候選。
6. `mouse_emulation=always`（唯一沒試過的值）。
7. **先確認畫面上那個紅色標記到底是不是滑鼠游標。** 它在不同畫面會換位置
   （`NEW GAME` 時在 `NO` 旁、`LOAD DATA` 時在第二格旁），
   但也可能是**選單的選取指示器**而不是游標 ——
   若是後者，整個「移動游標」的方向就是錯的。
   **這一項最便宜，應該先做**：在同一個畫面送幾次不同位置的點擊，
   看標記會不會跟著跳。

## 5. 這件事擋住了什麼

| 被擋住的 | 為什麼需要它 |
|---|---|
| 據點座標驗到 confirmed | 要在遊戲裡比對某據點的畫面位置與 `SINARIO.DAT` 的座標（`docs/formats/08` §4） |
| `docs/mechanics/15-realtime.md` | 要量「一個遊戲月 ＝ 多少真實時間」「戰略速度四檔各差多少」——
必須進到遊戲本體才量得到。**這份規格擋著整個規則層** |
| 地形類型的命名升 confirmed | 要看野戰戰場實際生成出來長什麼樣（`docs/playtest/03` §3） |
| 戰場圖塊的像素格式 | 同上，要有實機畫面當對照 |

**四件事同時卡在這一個點上**，所以下一輪它的優先度最高。

## 6. 已經拿到的（沒有白跑）

- **存檔槽是 4 個**，`LOAD DATA` 畫面實測看到，全部顯示 `0年 0月 0日`。
  日文說明書第 3.5 節說的「4カ所のどこにでもセーブが可能」
  從「說明書」等級升到 **confirmed**。
- **遊戲的時間單位是「年／月／日」**（存檔槽的顯示格式）。
  這是 `15-realtime.md` 的第一筆事實。
- `tools/dosboxx.sh` 的 timeline 多了 `move` 與 `rclick`，
  並修好了視窗座標換算的 bug（第 1 項）。
