#!/usr/bin/env python3
"""用「立即值 → 呼叫」的軌跡替未讀函式判角色，並把 TALK 索引翻成中文。

    tools/py.sh tools/re_classify.py <calltrace.tsv> <census.tsv> [--tier T4]

為什麼可行：這支程式的 UI 與流程函式幾乎都走「設常數 → 呼叫共用常式」，
而那些常數多半是 TALK 訊息索引。TALK.DAT 已全解，所以查表就能拿到
**遊戲自己的文字**——比從呼叫端的參數順序反推可靠得多。

輸出是 markdown，直接進 docs/re/。每支函式給的是**證據**不是命名：
「呼叫 sub_18853(cx=22)，訊息是『將游標移動至軍團的現在位置。』」
是事實；「這是軍團位置確認」是解讀，要人來下。
"""
import collections
import json
import os
import re
import sys

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# 已定案語意的共用常式。每一項都要能指到寫下它的文件，
# 沒有出處的不准列進來——這張表是解讀的地基。
KNOWN = {
    "sub_18853": ("狀態列提示", "cx", "talk", "docs/re/22 §4"),
    # ⚠ 索引是 cx 不是 al。`sub_18810` 把 cx 原封傳給 sub_1075B，由後者
    # 決定索引與 ×8 展開（docs/re/25 §1、§2）；ax 是變數值（ah 選變體）。
    # 早先記成 al 讓「al=0x93」被翻成 #147，全庫產生十幾個假陽性——
    # 人事四支明明是 cx=0x0B/0x0C/0x0D/0x0E，卻全被標成「不能准許出兵進攻」。
    "sub_18810": ("顯示訊息", "cx", "talk", "docs/re/25 §2"),
    "sub_193E9": ("彈出選單", "cx", "talk", "docs/re/22 §4"),
    "sub_1E453": ("讀滑鼠熱區編號", None, None, "docs/re/22 §2"),
    "sub_121E7": ("等待按鍵（CF=1 為右鍵取消）", None, None, "docs/re/22 §2"),
    "sub_12078": ("游標狀態保存", None, None, "docs/re/22 §3.3"),
    "sub_120D6": ("游標狀態復原", None, None, "docs/re/22 §3.3"),
    "sub_1716D": ("選軍團（回傳 bx＝記錄位址）", None, None, "docs/re/22 §3.3"),
    "sub_12151": ("游標／鏡頭移到座標", None, None, "docs/re/22 §3.3"),
    "sub_1FA37": ("圖塊 blit（ax＝尺寸、si＝資源偏移）", "ax", None, "docs/re/18 §2"),
    "sub_1F888": ("位元對齊 blit", None, None, "docs/re/18 §2"),
    "sub_1E3D7": ("寫熱區圖", None, None, "docs/re/22 §2"),
    "sub_20000": ("系統服務分派（ax＝服務號）", "ax", None, "docs/re/16"),
    "sub_181C0": ("開視窗", None, None, None),
    "sub_1820E": ("一覽表選取", None, None, None),
    # 以下由 tools/ida_hot_helpers.idc 這一輪讀出／確認，見 docs/re/24 §2
    "sub_1895D": ("繪矩形／視窗（al＝樣式）", "al", None, "docs/re/24 §2"),
    "sub_189A4": ("繪訊息框", "al", None, "docs/re/22 §4"),
    "sub_15E80": ("重畫狀態畫面（al＝畫面號）", "al", None, "docs/re/06"),
    "sub_17C6E": ("數值編輯視窗（ax＝上限）", "ax", None, "docs/re/13"),
    "sub_17663": ("選武將", None, None, "docs/re/22 §3.0"),
    "sub_1ECE0": ("亂數", None, None, "docs/re/10"),
    "sub_1062F": ("印數字", None, None, "docs/re/06"),
    "sub_101B4": ("繪圖狀態保存", None, None, None),
    "sub_101DB": ("繪圖狀態復原", None, None, None),
}

# 這幾個暫存器的立即值才可能是參數；其餘忽略以免噪音。
ARG_REGS = ("ax", "al", "ah", "bx", "bl", "cx", "cl", "dx", "dl", "si", "di", "bp")


def load_talk():
    p = os.path.join(REPO, "workplace", "tmp", "talk.json")
    if not os.path.exists(p):
        return None
    d = json.load(open(p, encoding="utf-8"))
    msgs = d["messages"] if isinstance(d, dict) and "messages" in d else d
    out = []
    for m in msgs:
        s = "".join(x for x in m if x).replace("　", " ").strip()
        out.append(s)
    return out


def load_census(path):
    funcs = {}
    with open(path, encoding="utf-8", errors="replace") as fh:
        for line in fh:
            f = line.rstrip("\n").split("\t")
            if f[0] == "FUNC" and len(f) >= 12:
                funcs[int(f[1], 16)] = {
                    "start": int(f[1], 16), "bytes": int(f[3]),
                    "insns": int(f[4]), "name": f[5], "callers": int(f[6]),
                }
    return funcs


def load_trace(path):
    """回傳 start -> 依位址排序的事件串。"""
    ev = collections.defaultdict(list)
    with open(path, encoding="utf-8", errors="replace") as fh:
        for line in fh:
            if line.startswith("#"):
                continue
            f = line.rstrip("\n").split("\t")
            if f[0] == "FN" or len(f) < 3:
                continue
            ev[int(f[1], 16)].append(f)
    for k in ev:
        ev[k].sort(key=lambda f: int(f[2], 16))
    return ev


def pair_args(events):
    """把 IMM 配到後面最近的 CALL 上。

    16-bit 真實模式的呼叫慣例就是「先設暫存器再 call」，所以往前找最近的
    賦值是正確的方向。但**同一個暫存器可能被覆寫**，所以只保留最後一次。
    回傳 [(call_ea, target_name, {reg: value}), …]。
    """
    out, pending = [], {}
    for f in events:
        if f[0] == "IMM":
            reg = f[3]
            if reg in ARG_REGS:
                # IDC 以 %d 輸出，16-bit 立即值會是負的；還原成無號。
                pending[reg] = int(f[4]) & 0xFFFF
        elif f[0] == "CALL":
            out.append((int(f[2], 16), f[4], dict(pending)))
            pending = {}
        elif f[0] in ("ICALL", "INT", "PORT"):
            pending = {}
    return out


def talk_of(talk, idx):
    if talk is None or not (0 <= idx < len(talk)):
        return None
    s = talk[idx]
    return s if s else "（空訊息）"


def describe(calls, talk):
    """回傳這支函式的證據列（字串）。只寫看得到的，不做命名。"""
    lines, seen = [], set()
    for ea, target, args in calls:
        info = KNOWN.get(target)
        if not info:
            continue
        label, argreg, kind, _src = info
        piece = label
        if argreg and argreg in args:
            v = args[argreg]
            if kind == "talk":
                if v == 0xFFFF:
                    piece += "（清除）"
                    if piece not in seen:
                        seen.add(piece); lines.append(piece)
                    continue
                txt = talk_of(talk, v)
                if txt:
                    piece += f"　#{v}「{txt[:34]}」"
                else:
                    piece += f"　{argreg}={v}"
            else:
                piece += f"　{argreg}=0x{v:X}"
        if piece not in seen:
            seen.add(piece)
            lines.append(piece)
    return lines


def main():
    trace_p, census_p = sys.argv[1], sys.argv[2]
    want = "T4"
    if "--tier" in sys.argv:
        want = sys.argv[sys.argv.index("--tier") + 1]
    # 預設把目錄本身排除——它逐支登記了每個未讀函式，不排除就會把
    # 「已登記」誤算成「已讀懂」，T4 直接歸零（見 re_coverage.scan_mentions）。
    exclude = ["docs/re/24-unread-function-catalogue.md"]
    if "--include-catalogue" in sys.argv:
        exclude = []

    talk = load_talk()
    funcs = load_census(census_p)
    events = load_trace(trace_p)

    # 分級沿用 re_coverage.py 的定義，這裡重算一次以免兩支不同步。
    sys.path.insert(0, os.path.join(REPO, "tools"))
    from re_coverage import scan_mentions, tier, segment_of, SEGMENTS
    hits = scan_mentions(REPO, exclude=exclude)
    for ea, f in funcs.items():
        f["tier"] = tier(hits.get(ea, set()))
        f["seg"] = segment_of(ea)

    if talk is None:
        print("⚠ 找不到 workplace/tmp/talk.json，**TALK 訊息那一層沒有跑**。\n"
              "先跑 tools/talkdat.py export workplace/orig/dosv/TALK.DAT big5 "
              "workplace/tmp/talk.json\n")

    sel = [f for f in funcs.values() if f["tier"] == want]
    withev = 0
    for lo, hi, segname in SEGMENTS:
        grp = sorted((f for f in sel if lo <= f["start"] < hi),
                     key=lambda f: f["start"])
        if not grp:
            continue
        print(f"\n### `{lo:05X}`–`{hi:05X}` {segname}（{len(grp)} 支）\n")
        print("| 函式 | B | 呼叫者 | 由呼叫的共用常式與訊息看到的證據 |")
        print("|---|---:|---:|---|")
        for f in grp:
            ev = describe(pair_args(events.get(f["start"], [])), talk)
            if ev:
                withev += 1
            cell = "<br>".join(ev) if ev else "—"
            print(f"| `{f['name']}` | {f['bytes']} | {f['callers']} | {cell} |")
    print(f"\n{want} 共 {len(sel)} 支，其中 {withev} 支可由共用常式與 TALK 訊息"
          f"取得證據，{len(sel) - withev} 支要逐支讀組語。")


if __name__ == "__main__":
    main()
