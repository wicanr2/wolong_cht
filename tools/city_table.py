#!/usr/bin/env python3
"""印出某個劇本的據點表：編號、名稱、所屬勢力、座標、四個鄰接據點。

    tools/py.sh tools/city_table.py [劇本 0-3] [--near 據點編號]

欄位出處 `docs/formats/08-sinario-save.md` §1.6（`+0x01` 勢力、`+0x02` 名稱、
`+0x08/0x0A` 座標、`+0x1C`–`+0x1F` 四個鄰接編號）。

為什麼要這一支：**原版的行軍目的地一覽是按編號排的**（`sub_1703C`），
一頁十列，要點哪一列得先知道那一列是誰。用猜的等於拿一次 3.5 分鐘的
實機擷取去換一個查表就有的答案。
"""
import sys

BLOCK = 0x56C0
CITY_OFF = 0x08C0
FACTION_OFF = 0x0080
PATH = "workplace/orig/dosv/SINARIO.DAT"
NO_OWNER = 0x18


def load(scenario):
    data = open(PATH, "rb").read()
    base = scenario * BLOCK
    cities = []
    for i in range(192):
        r = data[base + CITY_OFF + 32 * i: base + CITY_OFF + 32 * i + 32]
        name = r[2:8].split(b"\x00")[0].decode("cp950", "replace").strip()
        cities.append({
            "i": i, "name": name, "owner": r[1],
            "x": r[8] | (r[9] << 8), "y": r[10] | (r[11] << 8),
            "kind": r[0x16] & 0x0F,
            "adj": [v for v in r[0x1C:0x20] if v != 0xFF],
        })
    factions = []
    for f in range(22):
        r = data[base + FACTION_OFF + 64 * f: base + FACTION_OFF + 64 * f + 64]
        factions.append({"f": f, "capital": r[3],
                         "name": r[0x20:0x26].split(b"\x00")[0]
                         .decode("cp950", "replace").strip()})
    return cities, factions


def main():
    scenario = int(sys.argv[1]) if len(sys.argv) > 1 and sys.argv[1].isdigit() else 0
    cities, factions = load(scenario)
    near = None
    if "--near" in sys.argv:
        near = int(sys.argv[sys.argv.index("--near") + 1])
    print("劇本 %d：勢力首都" % scenario)
    for f in factions:
        cap = f["capital"]
        if cap < len(cities):
            print("  勢力 %2d %-8s 首都 %3d %s" % (f["f"], f["name"], cap, cities[cap]["name"]))
    print()
    if near is not None:
        c = cities[near]
        print("據點 %d %s（勢力 %d，(%d,%d)）的鄰接：" % (near, c["name"], c["owner"], c["x"], c["y"]))
        for j in c["adj"]:
            n = cities[j]
            print("  %3d %-6s 勢力 %2d  (%3d,%3d)  類型 %d" %
                  (j, n["name"], n["owner"], n["x"], n["y"], n["kind"]))
        print()
    print("前 20 筆（＝目的地一覽的第一頁與第二頁）")
    for c in cities[:20]:
        print("  %3d %-6s 勢力 %2d  (%3d,%3d)  類型 %d  鄰接 %s" %
              (c["i"], c["name"], c["owner"], c["x"], c["y"], c["kind"], c["adj"]))


if __name__ == "__main__":
    main()
