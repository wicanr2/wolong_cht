#!/usr/bin/env python3
"""把 DOSBox-X 截圖裡的 640×400 遊戲畫面切出來。

    tools/py.sh tools/parity_crop.py 進來.png 出去.png

## 為什麼需要這一支

`tools/dosv_live_capture.sh` 截的是 X11 視窗，而 DOSBox-X 的視窗是
**640×480**：遊戲的 640×400 放在 **y 偏移 40**，沒有縮放，上下各 40 px 黑邊。

兩件事因此要小心：

1. **對拍前不切就比不了**——`tools/parity_diff.py` 尺寸不同會直接報錯
   （那是刻意的，縮放後硬比會把「尺寸不對」這個線索洗掉）。
2. **這個偏移只管截圖，不管滑鼠。** INT 33 的對映吃的是**整個視窗**
   （黑邊也算），所以送點擊是 `視窗 y ＝ 遊戲 y × 1.2`，不是加 40。
   量法與證據在 `tools/dosv_capture.sh` 的說明。
   兩者搞混的症狀是「點了沒反應」——差十幾 px 正好落進按鈕之間的空隙。

偏移是**量出來的不是寫死的**：找出上下第一列與最後一列有內容的位置，
並檢查高度確實是 400。量不到就報錯，不要猜。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png, write_png

GAME_W, GAME_H = 640, 400


def find_offset(w, h, px):
    def lit(y):
        return sum(1 for x in range(w) if px[y][x] != (0, 0, 0))

    rows = [y for y in range(h) if lit(y) > 20]
    if not rows:
        raise SystemExit("整張圖都是黑的，切不出遊戲畫面")
    top, bottom = rows[0], rows[-1]
    if bottom - top + 1 != GAME_H:
        raise SystemExit("量到的內容高是 %d 不是 %d——先確認那真的是遊戲畫面"
                         % (bottom - top + 1, GAME_H))
    return top


def main():
    if len(sys.argv) != 3:
        raise SystemExit("用法: parity_crop.py 進來.png 出去.png")
    src, dst = sys.argv[1], sys.argv[2]
    w, h, px = read_png(src)
    if w != GAME_W:
        raise SystemExit("寬 %d 不是 %d，這支只處理 640 寬的 DOSBox-X 截圖" % (w, GAME_W))
    if h == GAME_H:
        top = 0
    else:
        top = find_offset(w, h, px)
    rows = [px[y] for y in range(top, top + GAME_H)]
    write_png(dst, GAME_W, GAME_H, rows)
    print("切出 %d×%d（y 偏移 %d）→ %s" % (GAME_W, GAME_H, top, dst))


if __name__ == "__main__":
    main()
