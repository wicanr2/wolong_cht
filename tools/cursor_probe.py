#!/usr/bin/env python3
"""量出兩張 DOSBox-X 截圖之間變動的像素範圍。

    tools/py.sh tools/cursor_probe.py a.png b.png

兩張只差一個滑鼠游標時，變動範圍就是**guest 認為游標在哪**——
這是「主機視窗座標 → 遊戲座標」對映的直接證據，
不必靠「點了有沒有反應」去反推（那會把對映錯誤與熱區錯誤混在一起）。
"""
import sys

sys.path.insert(0, "tools")
from parity_diff import read_png


def main():
    if len(sys.argv) != 3:
        raise SystemExit("用法: cursor_probe.py a.png b.png")
    wa, ha, a = read_png(sys.argv[1])
    wb, hb, b = read_png(sys.argv[2])
    if (wa, ha) != (wb, hb):
        raise SystemExit("尺寸不同：%dx%d vs %dx%d" % (wa, ha, wb, hb))
    pts = [(x, y) for y in range(ha) for x in range(wa) if a[y][x] != b[y][x]]
    if not pts:
        print("兩張完全相同")
        return
    xs = [p[0] for p in pts]
    ys = [p[1] for p in pts]
    print("變動 %d 點  x %d..%d  y %d..%d" % (len(pts), min(xs), max(xs), min(ys), max(ys)))
    print("視窗座標左上角 (%d, %d)" % (min(xs), min(ys)))
    print("若遊戲畫面在 y 偏移 40：減法對映 → 遊戲 y=%d；" % (min(ys) - 40))
    print("整窗等比對映（×400/480）→ 遊戲 y=%.1f" % (min(ys) * 400 / 480))


if __name__ == "__main__":
    main()
