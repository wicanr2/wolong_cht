# 18 — 松崗 DOS/V 密碼頁輸入驗證

**狀態：已證實，在受控 DOSBox-X 重播中按「確定」即可越過密碼頁；密碼頁不再是 DOS/V 原版行為驗證的阻擋。**

- 日期：2026-08-12
- 範圍：只驗證松崗 DOS/V 原版的啟動輸入；密碼頁不納入 remake，也沒有修改、補丁或散布原版檔案。
- 推論等級：下列畫面轉場是**已證實**；原始程式為何接受空白仍是**未知**。

## 結論

使用 `wolong-dosboxx:latest`（DOSBox-X 2025.02.01 SDL2）以唯讀原版檔案啟動後，
密碼頁的四格保持空白時直接點「確定」，10 秒後進入原版開場敘事。獨立的新副本中輸入
`0000` 與 `1234` 後再點「確定」也都進入開場；`1234` 與空白控制組的最終 640×480
截圖 SHA-256 完全相同：

```text
10cd0e199e7bd944a4664f3a3e4debd94a3986df1c36ca7eac4ffaf208ebbd34
```

因此使用者所述「隨便輸入數字會過」在此環境成立；更精確地說，**數字不是必要條件，
空白確認也會放行**。四個方框沒有繪出鍵入數字，故本次不宣稱數字被原版儲存或驗證；
只能確認畫面上的原始「確定」流程不阻擋後續開場。

## 受控實驗

每一組都是新的一次性容器，將 `/orig/dosv` 唯讀掛載並複製到容器內 `/tmp/game`；
沒有寫入 `workplace/orig/dosv`。設定為：

```text
machine=vgaonly
core=normal
cputype=486
cycles=fixed 20000
mouse_emulation=integration
int33=true
int33 max x=640
int33 max y=480
```

DOSBox-X 開場會變更視窗尺寸，所以等待密碼頁後重新取得 X11 client geometry，再以
640×480 client 相對座標操作。四格中心為 `(240,300)`、`(326,300)`、`(411,300)`、
`(496,300)`；「確定」為 `(497,240)`。滑鼠使用 `mouse_emulation=integration` 的
INT 33 路徑，鍵盤以 XTEST 送入，避免 `--window` 的 XSendEvent 被 DOSBox-X 忽略。

| 組別 | 四格操作 | 確定後 10 秒 | 最終截圖 SHA-256 |
|---|---|---|---|
| 空白控制組 | 不輸入 | 原版開場敘事 | `10cd0e199e7bd944a4664f3a3e4debd94a3986df1c36ca7eac4ffaf208ebbd34` |
| `0000` | 逐格點選並送 `0` | 原版開場敘事 | `cba3762abb82dd14919acaa3da22849760a2621b04ba488d09b21df44b861055` |
| `1234` | 逐格點選並送 `1`、`2`、`3`、`4` | 原版開場敘事 | `10cd0e199e7bd944a4664f3a3e4debd94a3986df1c36ca7eac4ffaf208ebbd34` |

`0000` 的截圖在同一段文字漸現的不同時間點，所以雜湊不同；畫面狀態同樣已離開密碼頁。
空白與 `1234` 是同一時點，故雜湊相同。

## 輸入與可追溯性

```text
KI.EXE     fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868
YNFONT.EXE 35b8c12be3b73983bbca6a2b856897e494230d99182e7380fb57fddddb633761
PASS.MAP   6197da6a23d88bb427a3a44dabdc09b2eef330ac3a2abdb7e9c73cf1b3735a8a
PASS.SCH   9e60b14dff27d97bf8ca47de2c2c48aa4627da3be61adc8ed4d8e42036944dad
START.BAT  afd1c26c66104d322c22fbda0b1fc35abc4d6da83f9b147ec88586954fdf311f
```

暫存截圖已於驗證後清理，未納入儲存庫或發行包。若要解釋 `PASS.*`／`YNFONT.EXE`
的真正比較語意，須另做 IDA Pro 靜態研究；它不是讓 remake 可玩或讓 DOS/V oracle 前進的前置條件。
