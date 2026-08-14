# 21 — DOSBox-X AI Bridge：第一次動態取樣

**狀態：管線打通，跑得起來。第一輪取樣裁決了一條斷言
（據點 `+0x00` 低 4 位 ＝ 敵方鄰居遮罩，192/192），
另外兩條（`+0x18`／`+0x14`）**沒有取到證據**——開局沒有軍團，
兩欄全 0，相符不構成證據。**

- 日期：2026-08-14
- 範圍：**只驗 PC-98 日文原版**。松崗 DOS/V 版仍被防拷擋著（[`01`](01-dosbox-dosv.md)）
- 工具：`tools/dosboxx_bridge.sh`、`tools/dosboxx_probe.py`
- 上游：`DOSBox-X-MCP-Debugger`，submodule 釘在 `5fcf624b`（`pmanyeh/dosbox-x`）

## 1. 這套東西是什麼

```
Claude Code ──MCP/stdio──> Python server ──TCP 127.0.0.1:9876──> debug_ai.cpp
                                                                （DOSBox-X 內）
```

`debug_ai.cpp` 掛在 DOSBox-X 原生除錯器上，**不是第二顆模擬器**：斷點、單步、
暫存器與記憶體都走原生機制。MCP 那一層只是把工具呼叫轉成
newline-delimited JSON。

`tools/dosboxx_probe.py` **直接講底下那個 TCP 協定**，不經過 MCP。
兩條路背後是同一支 `debug_ai.cpp`，差別只在誰發請求——
MCP 版要重啟 session 才載入，腳本版可重跑，適合放進驗收流程。

## 2. 三個坑

| 症狀 | 實際原因 |
|---|---|
| 編譯到 `dos_programs.cpp` 掛掉 | 缺 `GL/glu.h`，補 `libglu1-mesa-dev` |
| build 明明成功卻被 Dockerfile 擋下 | **我自己的驗收條件是錯的**：`grep "Debugger" build.log` 那個字樣根本不在 log 裡。改成看 `config.h` 的 `C_DEBUG`／`C_HEAVY_DEBUG` |
| 每個 `memory.read` 都回 `DEBUGGER_NOT_STOPPED` | **除錯器的 UI 是 ncurses，開在控制終端上**，而 DOSBox-X 是看自己的 stdout 判斷有沒有終端。容器要 `-t`，而且 dosbox-x 的輸出**不能重導向**——接了 `\| tee` 或 `> log` 就不是 TTY 了 |

第三個最值得記：log 裡有一行
`Debugger is not available unless you start DOSBox-X from a terminal`，
但 bridge 本身照樣接受連線、`debug.status` 照樣回 `ok:true`。
**「連得上」不等於「用得了」**——健康檢查要問的是能不能讀到記憶體，
不是 port 開了沒。

## 3. 怎麼找到據點表：不猜 DS

據點表是從 `SINARIO.DAT` 的劇本區塊整段載進來的，所以拿本機解析出來的
前 64 bytes 當簽章，在 guest 的 640 KB 裡搜：

```
據點表在 linear 0317D0（317D:0000），掃了 208,000 bytes
```

這比「假設 DS ＝ 某個值」可靠，也**避開了跨版本外推**——
`CLAUDE.md` §7 第 9 條禁止拿 DOS/V 的程式碼位址推 PC-98，
但資料的版面是同一份原始碼編出來的，可以先靜態驗：

| 不變量 | 松崗 DOS/V | PC-98 |
|---|---|---|
| X 落在 4–370 | 192/192 | 192/192 |
| Y 落在 9–248 | 192/192 | 192/192 |
| 生產力 ≤ 上限 | 192/192 | 192/192 |
| 城兵 ≤ 上限 | 192/192 | 192/192 |
| 類型 0–4 | 192/192 | 192/192 |
| 所屬 ≤ 24 | 192/192 | 192/192 |

兩版的 `SINARIO.DAT` 大小同為 88,832、內容不同（城市名一個 Big5 一個
Shift-JIS），但**版面是同一套**。

## 4. 取樣結果

### 4.1 ⭐ 據點 `+0x00` 低 4 位：裁決了

[`../re/44`](../re/44-threat-and-reinforcement-ai.md) §5 把這四位從
「哪幾個方向有鄰接」改成「哪幾個鄰接槽屬於別的勢力」。
guest 記憶體裡的 192 筆全部符合新讀法：

```
+0x00 低 4 位與鄰居實況不符：0 / 192
+0x1B 與遮罩位元數不符：      0 / 192
```

**先做正對照**，否則「符合」可能只是兩個假說剛好給同一個答案：

| 假說 | 命中 |
|---|---|
| A「哪幾個槽屬於別的勢力」 | **192/192** |
| B「哪幾個槽有鄰居」 | 12/192 |

**兩個假說在 180 筆上給不同答案**，所以這是可分辨的檢定，不是同義反覆。

`+0x1B` 的分佈也有內容（不是常數）：0 有 77 筆、1 有 70、2 有 36、3 有 8、4 有 1。

### 4.2 ⚠ `+0x18` 與 `+0x14`：沒有取到證據

```
佔用軍團數（+0x18）：{0: 192}
周邊威脅量（+0x14）：{0: 192}
求援冷卻（+0x17）： {0: 192}
```

取樣點在 `NEW GAME` 對話框，**開局一支軍團都沒有**，所以三欄全 0 是必然的。
`docs/re/44` 說 `+0x18` ＝ 停在該格的軍團數——0 支軍團對 0，
**這個「相符」不構成證據**，它與「欄位語意完全是別的東西」相容。

要驗這兩條得讓遊戲跑到有軍團在移動的狀態：開局 → 編成 → 行軍幾個 tick，
再取樣。那是下一輪。

## 5. 怎麼重跑

```sh
tools/dosboxx_bridge.sh start      # headless DOSBox-X ＋ PC-98 原版，bridge 在 9876
tools/dosboxx_bridge.sh health     # 正對照：回 {"ok":true} 才算活著
# 除錯器要先停住，memory.read 才服務得到
printf '{"id":1,"method":"execution.pause"}\n' | nc 127.0.0.1 9876
python3 tools/dosboxx_probe.py locate
python3 tools/dosboxx_probe.py sample 317D0
tools/dosboxx_bridge.sh stop
```

`locate` 的位址**不要寫死**：每次開機的載入位置可能不同，重跑一次比較便宜。

## 6. 未解

| 項目 | 現況 |
|---|---|
| `+0x18`／`+0x14` 的動態證據 | 要先把遊戲開到有軍團的狀態（§4.2）|
| 據點換手之後遮罩會不會跟著變 | `sub_1890A` 的行為，靜態讀得出來，動態沒驗——要打下一座城才看得到 |
| MCP 那一層 | `.mcp.json` 已就位，但**要重啟 session 才載入**，本輪的取樣走的是 TCP 腳本 |
| 松崗 DOS/V 側 | 防拷未解，這套工具在那一版上仍然用不了 |
| 上游授權 | `DOSBox-X-MCP-Debugger` 的原創碼**尚未選定授權條款**（README 明講是刻意留白）。本專案只在本機使用，未再散布 |
