#!/usr/bin/env python3
"""做一份**受控存檔**給對拍用：從既有的 SAVE.DAT 出發，只改指定欄位。

    tools/py.sh tools/parity_save.py <來源 SAVE.DAT> <輸出目錄> [選項]
    tools/py.sh tools/parity_save.py --selftest

⭐ **為什麼要這個。** 即時制的對拍最怕「兩邊各自擲骰」：天災、雲的漂移、
AI 的判斷都吃亂數，取樣點一拉長兩邊就分開，而分開之後看到的殘差
**不是規則分歧，是亂數分歧**。改狀態比等亂數可靠——把會動的東西
先關掉，剩下的差異才有意義（`CLAUDE.md` §3.1、`docs/spec/90`）。

目前支援的改法：

  --no-clouds   把大地圖上那 16 朵會飄的雲停用（物件記錄 +0x00 的 bit 7）。
                原版 `sub_12459`／`sub_12533` 都以 `cmp [si], 80h / jb` 跳過
                停用的槽，所以清掉就等於整個關掉——**不移動也不繪製**
                （`docs/spec/146`）。

  --diplomat F:G  把武將 G 派駐到勢力 F 當外交官：勢力記錄 `+0x2A` ＝ G、
                武將記錄 `+0x17` ＝ 3。**兩個欄位都要寫**——前者是外交前置閘
                `sub_165EF` 讀的那一個，後者是一覽表「身分」欄與候選過濾
                讀的那一個（`docs/spec/143`／`docs/spec/150`）。
                停戰與請求協助**沒有外交官就走不到第二步**，所以那兩條
                狀態列（#7）要靠這個改法才拍得到。

  --player F     把玩家所仕的勢力改成 F（區塊 `+0x0D` ＝ F × 0x40、`+0x0F` ＝ F）。
                長時間對拍要「玩家什麼都不做」時，用它把觀察平台固定在
                指定勢力上（劉備＝勢力 2，`docs/spec/161`）。

⚠ **輸出目錄要自己備妥其餘原版檔案**（dosgolem 的 gamedir 需要整套）。
這一支只寫 `SAVE.DAT`，不碰來源，也不碰 `workplace/orig/`。
"""

import argparse
import os
import shutil
import sys

BLOCK = 22208
BLOCKS = 4
FILE_SIZE = BLOCK * BLOCKS

# 32 筆物件記錄在區塊裡的位移（段內 0x2040 ＋ 0x80），每筆 16 B。
# 後 16 筆是常駐的雲（docs/spec/146 §2）。
MAP_OBJECT_BASE = 0x20C0
MAP_OBJECT_SIZE = 16
CLOUD_FIRST, CLOUD_COUNT = 16, 16

# 勢力表 22 筆 × 64 B、武將表 128 槽 × 32 B（docs/formats/08 §1）。
FACTION_BASE, FACTION_SIZE, FACTION_COUNT = 0x0080, 64, 22
FACTION_DIPLOMAT = 0x2A          # 我方派駐在這個勢力的外交官（0xFF ＝ 無）
# 玩家所仕的勢力：+0x0D 是勢力表位址（編號 × 0x40）、+0x0F 是編號本身
# （docs/formats/08 §1.2.1，`sub_11AC3` 同時寫這兩個）。
PLAYER_PTR, PLAYER_ID = 0x0D, 0x0F
GENERAL_BASE, GENERAL_SIZE, GENERAL_COUNT = 0x42C0, 32, 127
GENERAL_DUTY = 0x17              # 職務值 0–4；3 ＝ 外交官（docs/spec/143）
DUTY_DIPLOMAT = 3
NO_ONE = 0xFF


def disable_clouds(data: bytearray, block: int) -> int:
    """清掉雲的存在旗標，回傳實際關掉幾朵。"""
    base = block * BLOCK + MAP_OBJECT_BASE
    n = 0
    for slot in range(CLOUD_FIRST, CLOUD_FIRST + CLOUD_COUNT):
        off = base + slot * MAP_OBJECT_SIZE
        if data[off] & 0x80:
            data[off] &= ~0x80 & 0xFF
            n += 1
    return n


def assign_diplomat(data: bytearray, block: int, faction: int, general: int) -> str:
    """派 general 到 faction 當外交官，回傳一行說明。

    ⚠ **兩個欄位都要寫。** 只寫勢力 `+0x2A` 的話前置閘會過，但那個武將的
    「身分」欄還是「－－－」，而且他仍然出現在任命候選裡——**同一件事在
    兩張表各存一份**（`docs/spec/143` §2），漏掉一半會做出一份原版自己
    走不到的狀態，對拍就變成在比一個不存在的局面。
    """
    if not 0 <= faction < FACTION_COUNT:
        raise ValueError(f"勢力編號要在 0–{FACTION_COUNT - 1}，收到 {faction}")
    if not 0 <= general < GENERAL_COUNT:
        raise ValueError(f"武將編號要在 0–{GENERAL_COUNT - 1}，收到 {general}")
    base = block * BLOCK
    data[base + FACTION_BASE + faction * FACTION_SIZE + FACTION_DIPLOMAT] = general
    data[base + GENERAL_BASE + general * GENERAL_SIZE + GENERAL_DUTY] = DUTY_DIPLOMAT
    return f"武將 {general} 派駐勢力 {faction} 當外交官"


def set_player(data: bytearray, block: int, faction: int) -> str:
    """把玩家所仕的勢力改成 faction，回傳一行說明。

    ⚠ **兩個欄位都要寫**（`docs/formats/08` §1.2.1，兩者都 confirmed）：
    `+0x0D` 是勢力表**位址**（＝ 勢力編號 × 0x40），`+0x0F` 是勢力**編號**。
    原版 `sub_11AC3` 選定勢力時同時寫這兩個，只改一個會做出一份
    原版自己走不到的狀態——而對拍最怕的就是比一個不存在的局面。
    """
    if not 0 <= faction < FACTION_COUNT:
        raise ValueError(f"勢力編號要在 0–{FACTION_COUNT - 1}，收到 {faction}")
    base = block * BLOCK
    ptr = faction * FACTION_SIZE
    data[base + PLAYER_PTR] = ptr & 0xFF
    data[base + PLAYER_PTR + 1] = (ptr >> 8) & 0xFF
    data[base + PLAYER_ID] = faction
    return f"玩家勢力 ＝ {faction}（+0x0D ＝ {ptr:#06x}、+0x0F ＝ {faction}）"


def selftest() -> int:
    """正對照 ＋ 負對照。少了負對照，一支「什麼都沒改」的工具也會通過。"""
    ok = True

    def check(label, cond):
        nonlocal ok
        print(f"  {'✓' if cond else '✗'} {label}")
        ok = ok and cond

    data = bytearray(FILE_SIZE)
    base = MAP_OBJECT_BASE
    for slot in range(32):
        data[base + slot * MAP_OBJECT_SIZE] = 0x80
    n = disable_clouds(data, 0)
    check("關掉 16 朵雲", n == 16)
    check("前 16 槽（火災／暴動）沒被動到",
          all(data[base + s * MAP_OBJECT_SIZE] == 0x80 for s in range(16)))
    check("後 16 槽的 bit 7 清掉了",
          all(data[base + s * MAP_OBJECT_SIZE] == 0x00 for s in range(16, 32)))
    check("其餘 byte 一個都沒動",
          sum(1 for i, b in enumerate(data)
              if b != 0 and i % MAP_OBJECT_SIZE == 0) == 16)
    # 負對照：本來就沒有雲的區塊要回 0，不能假裝有做事。
    empty = bytearray(FILE_SIZE)
    check("空區塊回 0 朵（負對照）", disable_clouds(empty, 0) == 0)
    # 負對照：只改指定的區塊。
    two = bytearray(FILE_SIZE)
    for slot in range(32):
        two[BLOCK + base + slot * MAP_OBJECT_SIZE] = 0x80
    check("只動指定的區塊（負對照）", disable_clouds(two, 0) == 0)

    # --diplomat：兩個欄位都要寫，而且只寫那兩個。
    d = bytearray(FILE_SIZE)
    for i in range(FILE_SIZE):
        d[i] = 0xFF if i % 7 == 0 else 0x11
    before = bytes(d)
    assign_diplomat(d, 0, 5, 9)
    fa = FACTION_BASE + 5 * FACTION_SIZE + FACTION_DIPLOMAT
    ge = GENERAL_BASE + 9 * GENERAL_SIZE + GENERAL_DUTY
    check("勢力 +0x2A 寫成武將編號", d[fa] == 9)
    check("武將 +0x17 寫成職務 3", d[ge] == DUTY_DIPLOMAT)
    check("其餘 byte 一個都沒動",
          sum(1 for i in range(FILE_SIZE) if d[i] != before[i]) <= 2)
    # 負對照：第 2 槽不受影響。
    d2 = bytearray(FILE_SIZE)
    assign_diplomat(d2, 1, 5, 9)
    check("只動指定的區塊（負對照）", d2[fa] == 0 and d2[BLOCK + fa] == 9)
    for bad in ((22, 0), (0, 127), (-1, 0)):
        try:
            assign_diplomat(bytearray(FILE_SIZE), 0, *bad)
            check(f"擋下超出範圍的 {bad}", False)
        except ValueError:
            check(f"擋下超出範圍的 {bad}", True)
    # --player：兩個欄位都要寫（位址與編號），而且只寫那兩個。
    p = bytearray(FILE_SIZE)
    for i in range(FILE_SIZE):
        p[i] = 0xFF if i % 5 == 0 else 0x22
    before_p = bytes(p)
    set_player(p, 0, 2)                      # 劉備＝勢力 2，2 × 0x40 ＝ 0x80
    check("+0x0D 低位元組 ＝ 0x80", p[PLAYER_PTR] == 0x80)
    check("+0x0E 高位元組 ＝ 0x00", p[PLAYER_PTR + 1] == 0x00)
    check("+0x0F ＝ 勢力編號 2", p[PLAYER_ID] == 2)
    check("其餘 byte 一個都沒動",
          sum(1 for i in range(FILE_SIZE) if p[i] != before_p[i]) <= 3)
    # 高位元組真的會用到：勢力 21 × 0x40 ＝ 0x540。
    p21 = bytearray(FILE_SIZE)
    set_player(p21, 0, 21)
    check("勢力 21 的位址跨到高位元組（0x540）",
          p21[PLAYER_PTR] == 0x40 and p21[PLAYER_PTR + 1] == 0x05)
    # 負對照：只動指定的區塊。
    p2 = bytearray(FILE_SIZE)
    set_player(p2, 1, 2)
    check("只動指定的區塊（負對照）",
          p2[PLAYER_ID] == 0 and p2[BLOCK + PLAYER_ID] == 2)
    for bad in (22, -1, 127):
        try:
            set_player(bytearray(FILE_SIZE), 0, bad)
            check(f"擋下超出範圍的勢力 {bad}", False)
        except ValueError:
            check(f"擋下超出範圍的勢力 {bad}", True)

    return 0 if ok else 1


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("source", nargs="?")
    ap.add_argument("outdir", nargs="?")
    ap.add_argument("--slot", type=int, default=0, help="改第幾個區塊（0–3）")
    ap.add_argument("--no-clouds", action="store_true")
    ap.add_argument("--diplomat", metavar="勢力:武將",
                    help="派一個外交官（勢力記錄 +0x2A ＋ 武將記錄 +0x17）")
    ap.add_argument("--player", type=int, metavar="勢力",
                    help="改玩家所仕的勢力（區塊 +0x0D ＋ +0x0F）")
    ap.add_argument("--selftest", action="store_true")
    ns = ap.parse_args()

    if ns.selftest:
        return selftest()
    if not ns.source or not ns.outdir:
        ap.error("要給來源 SAVE.DAT 與輸出目錄")
    if not 0 <= ns.slot < BLOCKS:
        ap.error("--slot 要在 0–3")

    data = bytearray(open(ns.source, "rb").read())
    if len(data) != FILE_SIZE:
        ap.error(f"{ns.source} 是 {len(data)} B，預期 {FILE_SIZE}")

    changes = []
    if ns.no_clouds:
        changes.append(f"停用 {disable_clouds(data, ns.slot)} 朵雲")
    if ns.diplomat:
        try:
            faction, general = (int(v) for v in ns.diplomat.split(":", 1))
        except ValueError:
            ap.error("--diplomat 要寫成 `勢力:武將`，兩個都是十進位整數")
        try:
            changes.append(assign_diplomat(data, ns.slot, faction, general))
        except ValueError as err:
            ap.error(str(err))
    if ns.player is not None:
        try:
            changes.append(set_player(data, ns.slot, ns.player))
        except ValueError as err:
            ap.error(str(err))
    if not changes:
        ap.error("沒有指定任何改法")

    os.makedirs(ns.outdir, exist_ok=True)
    out = os.path.join(ns.outdir, "SAVE.DAT")
    if os.path.abspath(out) == os.path.abspath(ns.source):
        ap.error("輸出會蓋掉來源；換一個目錄")
    with open(out, "wb") as f:
        f.write(data)
    print(f"{out}（第 {ns.slot + 1} 槽）：" + "、".join(changes))
    return 0


if __name__ == "__main__":
    sys.exit(main())
