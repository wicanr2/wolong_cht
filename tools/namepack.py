#!/usr/bin/env python3
"""從兩版 SINARIO.DAT 抽人名／地名，產生日文與英文名表（docs/spec/84）。

    tools/py.sh tools/namepack.py ja        # 日文名表（含外字反推）
    tools/langpack.sh names-en              # 英文羅馬化（需要 pypinyin）
    tools/py.sh tools/namepack.py selftest

## 日文名怎麼來的

PC-98 日文原版的 `SINARIO.DAT` 與松崗繁中版**同結構、只有語言不同**，
武將記錄 `+2`（6 byte）與據點記錄 `+2`（6 byte）逐槽對應同一個人／同一座城。
所以日文名不需要翻譯，直接讀原版。

⚠ 日文名裡有 **PC-98 外字**（Shift-JIS 未定義的 `0xEC40`–`0xEC61`）：
`汜`、`瓚`、`繡`、`傕` 這些字不在 JIS X 0208 裡，原版把它們放進外字區
自帶字形。外字碼本身不帶語意，**靠兩版逐字對齊反推**：同一個位置，
繁中版是什麼字，那個外字碼就是什麼字。反推出來的表寫進輸出檔，
呈現時交給字型鏈（JIS 字型沒有 → 退到倚天 Big5 字型，見 spec/84）。
"""

import argparse
import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
DOSV = ROOT / "workplace" / "orig" / "dosv" / "SINARIO.DAT"
PC98 = ROOT / "workplace" / "orig" / "pc98" / "SINARIO.DAT"
OUT_JA = ROOT / "translations" / "names-ja.json"
OUT_EN = ROOT / "translations" / "names-en.json"

BLOCK = 22208          # 一個劇本區塊
SCENARIOS = 4
GENERAL_BASE = 0x42C0  # 武將表（127 × 32B），名稱在 +2、6 byte
GENERAL_N = 127
CITY_BASE = 0x8C0      # 據點表（192 × 32B），名稱在 +2、6 byte
CITY_N = 192
NAME_OFF, NAME_LEN = 2, 6


def records(path, base, count):
    """逐劇本、逐槽回傳名稱欄的原始 bytes。"""
    data = path.read_bytes()
    for s in range(SCENARIOS):
        blk = data[s * BLOCK:(s + 1) * BLOCK]
        for i in range(count):
            off = base + i * 32 + NAME_OFF
            yield s, i, blk[off:off + NAME_LEN]


def big5_chars(raw):
    """繁中名 → 字元 list（去掉全形空白補位）。"""
    out = []
    for i in range(0, len(raw), 2):
        pair = raw[i:i + 2]
        if pair in (b"\xa1\x40", b"\x00\x00", b"  "):
            continue
        try:
            out.append(pair.decode("cp950"))
        except UnicodeDecodeError:
            out.append("?" + pair.hex())
    return out


def sjis_tokens(raw):
    """日文名 → token list。已定義的字回字元，外字回 ('gaiji', code)。"""
    out = []
    for i in range(0, len(raw), 2):
        pair = raw[i:i + 2]
        if pair in (b"\x81\x40", b"\x00\x00", b"  "):
            continue
        try:
            out.append(pair.decode("cp932"))
        except UnicodeDecodeError:
            out.append(("gaiji", int.from_bytes(pair, "big")))
    return out


def build_ja():
    gaiji = {}       # 外字碼 → 繁中字（逐字對齊反推）
    conflicts = []
    names = {}       # 繁中名 → 日文名
    stats = {"general": 0, "city": 0, "skipped": 0}

    for kind, base, count in (("general", GENERAL_BASE, GENERAL_N),
                              ("city", CITY_BASE, CITY_N)):
        zh_all = {(s, i): raw for s, i, raw in records(DOSV, base, count)}
        for s, i, ja_raw in records(PC98, base, count):
            zh = big5_chars(zh_all[(s, i)])
            ja = sjis_tokens(ja_raw)
            if not zh or not ja:
                continue
            if len(zh) != len(ja):
                # 字數不一致就不能逐字對齊——記下來，名字仍以日文側可解的
                # 部分為準（外字保持未解，呈現時是缺字框）。
                stats["skipped"] += 1
                continue
            for zc, jt in zip(zh, ja):
                if not isinstance(jt, tuple):
                    continue
                code = jt[1]
                if code in gaiji and gaiji[code] != zc:
                    conflicts.append((hex(code), gaiji[code], zc))
                gaiji.setdefault(code, zc)
            stats[kind] += 1

    # 第二輪：外字表建好之後才組名字（同一個外字可能第一次出現在
    # 字數對不上的那一筆）。
    for kind, base, count in (("general", GENERAL_BASE, GENERAL_N),
                              ("city", CITY_BASE, CITY_N)):
        zh_all = {(s, i): raw for s, i, raw in records(DOSV, base, count)}
        for s, i, ja_raw in records(PC98, base, count):
            zh = "".join(big5_chars(zh_all[(s, i)]))
            ja = "".join(
                t if not isinstance(t, tuple) else gaiji.get(t[1], "〓")
                for t in sjis_tokens(ja_raw))
            if not zh or not ja:
                continue
            if zh in names and names[zh] != ja:
                conflicts.append((zh, names[zh], ja))
            names.setdefault(zh, ja)

    if conflicts:
        print(f"⚠ {len(conflicts)} 筆對齊衝突（前 5）：{conflicts[:5]}", file=sys.stderr)
    payload = {
        "note": "PC-98 日文原版的人名與地名，逐槽對齊松崗繁中版取得；"
                "外字碼由兩版逐字對齊反推（tools/namepack.py）",
        "gaiji": {hex(k): v for k, v in sorted(gaiji.items())},
        "names": dict(sorted(names.items())),
    }
    OUT_JA.write_text(json.dumps(payload, ensure_ascii=False, indent=1) + "\n",
                      encoding="utf-8")
    diff = sum(1 for k, v in names.items() if k != v)
    print(f"names-ja.json：{len(names)} 個名字（{diff} 個與繁中不同形）、"
          f"外字 {len(gaiji)} 碼、字數對不上 {stats['skipped']} 筆")


# ── 英文羅馬化 ──────────────────────────────────────────────────────
#
# 使用者裁定（2026-08-26）：人名用三國人物名稱的拼音，大寫、以連字號分隔
# 姓與名（`CAO-CAO`）。複姓要當成一個整體（`ZHUGE-LIANG`，不是
# `ZHU-GE-LIANG`），所以複姓表不能省。

COMPOUND_SURNAMES = [
    "諸葛", "司馬", "夏侯", "皇甫", "公孫", "太史", "上官", "歐陽", "尉遲",
    "淳于", "令狐", "長孫", "宇文", "慕容", "東方", "西門", "南宮", "軒轅",
    "第五", "呼廚", "沮授",
]

# 姓氏的固定讀法：pypinyin 給的是常用讀音，不是姓氏讀音。
# **只套用在姓**——同一個字當地名時讀法不同（樂進 Yuè 進、樂成 Lè 城）。
SURNAME_OVERRIDES = {
    "單": "SHAN",   # 單福／單經
    "解": "XIE",
    "區": "OU",
    "樂": "YUE",    # 樂進（地名的「樂」讀 Lè，所以只套在姓）
    "紀": "JI",
    "種": "CHONG",  # 種輯
    "繆": "MIAO",
    "召": "SHAO",
    "重": "CHONG",
    "闞": "KAN",
    "覃": "QIN",
    "查": "ZHA",
    "任": "REN",
    "燕": "YAN",
    "牟": "MOU",
    "朴": "PIAO",
    "冒": "MO",
    "貟": "YUN",
}

# 名（不含姓）的破音字覆寫。
GIVEN_OVERRIDES = {
    "繇": "YAO",    # 劉繇
}

# 地名的固定讀法：這些是歷史地名的專讀，pypinyin 給的是今讀或常用讀音。
CITY_OVERRIDES = {
    "會稽": "KUAIJI",    # kuài jī，不是 huì jī
    "信都": "XINDU",     # 地名的「都」讀 dū
    "武都": "WUDU",
    "瑯邪": "LANGYA",    # 邪 讀 yá
    "長阪": "CHANGBAN",  # 長坂坡
    "沓中": "TAZHONG",   # 沓 讀 tà
    "陜": "SHAN",        # 陝縣
}

# 整個名字的覆寫：通行英文寫法與上面規則不同時以這裡為準。
NAME_OVERRIDES = {
    "曹操": "CAO-CAO",
    "劉禪": "LIU-SHAN",   # 禪 讀 shàn
    "孫乾": "SUN-QIAN",   # 乾 讀 qián
    "劉備": "LIU-BEI",
    "孫策": "SUN-CE",
    "孫權": "SUN-QUAN",
    "諸葛亮": "ZHUGE-LIANG",
    "司馬懿": "SIMA-YI",
    "呂布": "LU-BU",
    "貂蟬": "DIAO-CHAN",
    "董卓": "DONG-ZHUO",
    "袁紹": "YUAN-SHAO",
    "關羽": "GUAN-YU",
    "張飛": "ZHANG-FEI",
    "趙雲": "ZHAO-YUN",
    "馬超": "MA-CHAO",
    "黃忠": "HUANG-ZHONG",
    "周瑜": "ZHOU-YU",
    "魯肅": "LU-SU",
    "呂蒙": "LU-MENG",
    "陸遜": "LU-XUN",
    "許褚": "XU-CHU",
    "典韋": "DIAN-WEI",
    "郭嘉": "GUO-JIA",
    "荀彧": "XUN-YU",
    "賈詡": "JIA-XU",
    "貂嬋": "DIAO-CHAN",
}


def split_surname(name):
    """回傳 (姓, 名)。單姓一律取第一個字。"""
    for s in COMPOUND_SURNAMES:
        if name.startswith(s) and len(name) > len(s):
            return s, name[len(s):]
    if len(name) == 1:
        return name, ""
    return name[0], name[1:]


def build_en():
    """人名與地名分開處理。

    人名：姓與名以連字號分隔（`CAO-CAO`、`ZHUGE-LIANG`）。
    地名：**不分隔**——`成都` 是 `CHENGDU` 不是 `CHENG-DU`，而且要用
    詞組讀音（`長安` 是 CHANGAN 不是 ZHANGAN、`成都` 的「都」讀 dū）。
    逐字轉會兩個都錯，所以地名整串丟給 pypinyin 讓它吃詞庫。
    """
    from opencc import OpenCC
    from pypinyin import lazy_pinyin, Style

    # ⚠ **pypinyin 的詞組字典是簡體鍵**：繁體字串查不到詞組，會逐字退回
    # 常用讀音，於是「長安」變成 ZHANGAN、「會稽」變成 HUIJI。
    # 先轉簡體再查拼音，key 仍保持繁體。
    t2s = OpenCC("t2s")

    def phrase(text, overrides=None, as_list=False):
        """整串轉拼音（先轉簡體才吃得到 pypinyin 詞庫），overrides 逐字優先。"""
        simp = t2s.convert(text)
        if len(simp) != len(text):  # 轉換改變字數就退回逐字，寧可保守
            simp = text
        py = lazy_pinyin(simp, style=Style.NORMAL, errors=lambda x: [""] * len(x))
        if len(py) != len(text):
            py = lazy_pinyin(text, style=Style.NORMAL, errors=lambda x: [""] * len(x))
        out = []
        for ch, syl in zip(text, py):
            if overrides and ch in overrides:
                out.append(overrides[ch])
                continue
            out.append(syl.upper())
        return out if as_list else "".join(out)

    people, cities = set(), set()
    for kind, base, count, bucket in (("general", GENERAL_BASE, GENERAL_N, people),
                                      ("city", CITY_BASE, CITY_N, cities)):
        for _, _, raw in records(DOSV, base, count):
            name = "".join(big5_chars(raw))
            if name:
                bucket.add(name)

    def apostrophe(syllables):
        """隔音符號：後一個音節以 a／o／e 起頭時要加 `'`，
        否則 `長安` 會變成讀不出來的 CHANGAN、`西安` 變成 XIAN（單音節）。"""
        out = ""
        for i, syl in enumerate(syllables):
            if i and syl[:1] in "AOE":
                out += "'"
            out += syl
        return out

    out, unresolved = {}, []
    for name in sorted(cities):
        # 地名優先：同名時城名的讀法對玩家比較常見（清單上看得到）。
        if name in CITY_OVERRIDES:
            out[name] = CITY_OVERRIDES[name]
            continue
        v = apostrophe(phrase(name, as_list=True))
        if v:
            out[name] = v
        else:
            unresolved.append(name)
    for name in sorted(people):
        if name in NAME_OVERRIDES:
            out[name] = NAME_OVERRIDES[name]
            continue
        surname, given = split_surname(name)
        a = phrase(surname, SURNAME_OVERRIDES)
        b = phrase(given, GIVEN_OVERRIDES) if given else ""
        if not a or (given and not b):
            unresolved.append(name)
            continue
        out[name] = f"{a}-{b}" if b else a
    if unresolved:
        print(f"⚠ {len(unresolved)} 個名字轉不出拼音：{unresolved}", file=sys.stderr)
    payload = {
        "note": "三國人物與地名的漢語拼音，大寫（使用者裁定 2026-08-26）。"
                "人名的姓與名以連字號分隔、複姓當成一個整體；"
                "地名不分隔並使用詞組讀音。覆寫表在 tools/namepack.py",
        "names": dict(sorted(out.items())),
    }
    OUT_EN.write_text(json.dumps(payload, ensure_ascii=False, indent=1) + "\n",
                      encoding="utf-8")
    print(f"names-en.json：{len(out)} 個名字（人 {len(people)}／地 {len(cities)}）")


def selftest():
    bad = 0
    ja = json.loads(OUT_JA.read_text(encoding="utf-8"))
    names_ja = ja["names"]
    if len(names_ja) < 300:
        sys.exit(f"日文名表只有 {len(names_ja)} 筆，太少")
    # 外字必須全部反推出來——留著 〓（U+3013）代表沒對齊到。
    left = [k for k, v in names_ja.items() if "〓" in v]
    if left:
        print(f"⚠ {len(left)} 個日文名still有未解外字：{left[:5]}", file=sys.stderr)
        bad += 1
    for zh, expect in (("曹操", "曹操"), ("孫權", "孫権"), ("郭汜", "郭汜")):
        if names_ja.get(zh) != expect:
            print(f"⚠ {zh} 的日文名 = {names_ja.get(zh)!r}，預期 {expect!r}",
                  file=sys.stderr)
            bad += 1
    if OUT_EN.exists():
        en = json.loads(OUT_EN.read_text(encoding="utf-8"))["names"]
        for zh, expect in (("曹操", "CAO-CAO"), ("劉備", "LIU-BEI"),
                           ("諸葛亮", "ZHUGE-LIANG"), ("許昌", "XUCHANG"),
                           ("成都", "CHENGDU"), ("長安", "CHANG'AN"),
                           ("會稽", "KUAIJI")):
            if en.get(zh) != expect:
                print(f"⚠ {zh} 的英文名 = {en.get(zh)!r}，預期 {expect!r}",
                      file=sys.stderr)
                bad += 1
        empty = [k for k, v in en.items() if not v or v.startswith("-") or v.endswith("-")]
        if empty:
            print(f"⚠ {len(empty)} 個英文名不完整：{empty[:5]}", file=sys.stderr)
            bad += 1
        print(f"selftest：日文 {len(names_ja)} 名、英文 {len(en)} 名、問題 {bad}")
    else:
        print(f"selftest：日文 {len(names_ja)} 名（英文名表尚未產生）、問題 {bad}")
    return 1 if bad else 0


if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("cmd", choices=["ja", "en", "selftest"])
    args = ap.parse_args()
    sys.exit({"ja": build_ja, "en": build_en, "selftest": selftest}[args.cmd]())
