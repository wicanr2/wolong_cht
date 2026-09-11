#!/usr/bin/env python3
"""畫面對拍組的清單與判定（權威資料在 `tools/parity_screens.json`）。

    tools/py.sh tools/parity_screens.py list            # 給 parity_screens.sh 迴圈用
    tools/py.sh tools/parity_screens.py judge 組名 差分表.md
    tools/py.sh tools/parity_screens.py --selftest

## 為什麼要分成資料與執行

規則層的 `tools/parity_ck.sh` 從 2026-09-10 起每一輪都跑得動，是因為
「哪些檢查點、拿什麼比」是資料。畫面層先前沒有這一層——每一組的原版圖、
命令與取樣點散在二十幾份 `docs/playtest/`，有的只寫了目錄沒寫哪一張圖，
有的 remake 側漏了 `-save-file`。**寫不齊就重跑不出來**，而重跑不出來的
對拍在下一輪規則改動之後完全不會開口（`docs/spec/197` 就是這樣拖了三天）。

## rects：文件當初對的是視窗本體，不是整個分區

有幾組（說服場景、判決、液晶、遊戲結束）原版與 remake 只在**視窗內**對齊，
視窗外的大地圖是另一個局面。那幾組的 `rects` 記下當初真正比過的矩形；
`regions` 那一欄仍然照跑，**兩份數字都印**——把沒對齊的範圍藏起來，
下一輪就會有人以為整張都對過。

## allow 與 gap 是兩件事

`allow` 是**有正當理由**的殘差：原版錄影自己畫的滑鼠游標、M7 的校訂字、
刻意的 remake 差異。理由寫在第二欄。

`gap` 是**已知但還沒收掉的缺口**。它照樣擋回歸（超過上限就不過），
但在報表上分開印，總結行也分開數——把缺口混進 allow，閘就會在
「全部通過」的外表下把待辦藏起來。

## allow 不是及格線

`allow` 記的是**有正當理由**的殘差（原版錄影自己畫的滑鼠游標、
fixture 與原版停在不同日期），理由寫在第二欄。沒有理由就是 0。
判定只有兩種：符合預期（`=` 或 `<`，後者代表變好了要回填）與
**超出**（`>`，那是回歸）。
"""
import json
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
DATA = os.path.join(HERE, "parity_screens.json")

# `tools/parity_diff.py` 的輸出：| `區名` | n / total | x.xx% | d | 判定 |
ROW = re.compile(r"^\|\s*`([^`]+)`\s*\|\s*(\d+)\s*/\s*(\d+)\s*\|")


def load():
    with open(DATA, encoding="utf-8") as f:
        return json.load(f)["groups"]


def parse_diff(text):
    """把差分表讀成 {區名: (不同, 總數)}。

    ⚠ **空的結果要與「全部相同」分開**。`parity_diff.py` 尺寸不符時印的是
    一行 `✗ 尺寸不同…`，一列表格都沒有——那時回空 dict，呼叫端要當成失敗，
    不能當成「零個差異」。這兩者在 shell 的 exit code 上長得一樣。
    """
    out = {}
    for line in text.splitlines():
        m = ROW.match(line.strip())
        if m:
            out[m.group(1)] = (int(m.group(2)), int(m.group(3)))
    return out


def judge(name, path):
    groups = {g["name"]: g for g in load()}
    if name not in groups:
        print("✗ 沒有這一組：%s" % name)
        return 2
    g = groups[name]
    allow = g.get("allow", {})
    gaps = g.get("gap", {})
    with open(path, encoding="utf-8") as f:
        text = f.read()
    stats = parse_diff(text)
    if not stats:
        print("✗ %s：差分表是空的——比對本身失敗了" % name)
        print(text.strip()[:400])
        return 2
    bad, better, open_gaps = [], [], []
    for region, (n, total) in sorted(stats.items()):
        if region in gaps:
            cap = gaps[region][0]
            if n > cap:
                bad.append((region, n, total, cap))
            else:
                open_gaps.append((region, n, total, gaps[region][1]))
                if n < cap:
                    better.append((region, n, cap))
            continue
        cap = allow.get(region, [0])[0]
        if n > cap:
            bad.append((region, n, total, cap))
        elif n < cap:
            better.append((region, n, cap))
    for region, n, total, cap in bad:
        why = "（上限 %d）" % cap if cap else ""
        print("  ✗ %-14s %d / %d%s" % (region, n, total, why))
    for region, n, total, why in open_gaps:
        print("  ⚠ %-14s %d / %d 缺口：%s" % (region, n, total, why))
    for region, n, cap in better:
        print("  ⭐ %-14s %d（allow 記的是 %d，變好了 → 回填 JSON）" % (region, n, cap))
    if bad:
        print("✗ %s：%d 區超出" % (name, len(bad)))
        return 1
    kept = ", ".join("%s=%d" % (r, allow[r][0]) for r in sorted(allow) if r in stats)
    tag = "（已知殘差 %s）" % kept if kept else "（全區 0 px）"
    if open_gaps:
        print("◐ %s：%d 個已知缺口%s" % (name, len(open_gaps), tag if kept else ""))
        return 3
    print("✓ %s%s" % (name, tag))
    return 0


def selftest():
    """⚠ 這一支自己也要有正對照。

    只驗「合格的表判成過」不夠——一支永遠回 0 的壞判定器也會通過。
    所以第二組刻意超出上限，必須判成不過。
    """
    ok = "| `banner` | 0 / 20480 | 0.00% | 0 | PASS |\n| `map` | 95 / 145152 | 0.07% | 9 | NEAR |"
    assert parse_diff(ok) == {"banner": (0, 20480), "map": (95, 145152)}
    assert parse_diff("✗ 尺寸不同：原版 640x480、remake 640x400") == {}
    groups = load()
    assert groups, "清單是空的"
    for g in groups:
        for key in ("name", "note", "doc", "orig", "crop", "args", "regions"):
            assert key in g, "%s 缺 %s" % (g.get("name"), key)
        assert g["regions"] in ("strategy", "tactical"), g["name"]
        for field in ("allow", "gap"):
            for region, why in g.get(field, {}).items():
                assert len(why) == 2 and why[1], "%s 的 %s[%s] 沒寫理由" % (g["name"], field, region)
        both = set(g.get("allow", {})) & set(g.get("gap", {}))
        assert not both, "%s 的 %s 同時在 allow 與 gap" % (g["name"], both)
        for rname, rect in g.get("rects", {}).items():
            # 名稱進 shell 的空白分隔清單，含空白就會被切成兩段而靜默比錯一塊。
            assert " " not in rname and ":" not in rname, "%s 的矩形名 %r 不能含空白或冒號" % (g["name"], rname)
            assert len(rect.split(",")) == 4, "%s 的 %s 不是 x,y,w,h" % (g["name"], rname)
    print("✓ parity_screens selftest：%d 組、解析與理由欄都在" % len(groups))
    return 0


def main():
    if len(sys.argv) >= 2 and sys.argv[1] == "--selftest":
        return selftest()
    if len(sys.argv) >= 2 and sys.argv[1] == "list":
        for g in load():
            # ⚠ tab 是 IFS 的空白字元，連續兩個會被 `read` 合併成一個——
            #   沒有矩形的組若輸出空欄位，`note` 就會被讀進 `rects`，
            #   而 note 裡的 `/` 會讓後面那道 sed 炸掉。補一個 `-` 佔位。
            rects = " ".join("%s:%s" % (k, v) for k, v in g.get("rects", {}).items()) or "-"
            print("\t".join([g["name"], g["orig"], "1" if g["crop"] else "0",
                             g["regions"], " ".join(g["args"]), rects, g["note"]]))
        return 0
    if len(sys.argv) == 4 and sys.argv[1] == "judge":
        return judge(sys.argv[2], sys.argv[3])
    print(__doc__)
    return 2


if __name__ == "__main__":
    sys.exit(main())
