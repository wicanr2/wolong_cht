#!/usr/bin/env python3
"""比較dosgolem原始快照與Go行軍重播，不修改輸入。

以 tools/py.sh 執行。原版資料與操作時間軸入口見
docs/playtest/126-player-march-collision-snapshot.md。
"""
import argparse
import hashlib
import json
from pathlib import Path


def read(path, size):
    data = path.read_bytes()
    if len(data) != size:
        raise ValueError(f"{path}: {len(data)} bytes，預期{size}")
    return data


def diff(a, b, offset=0):
    return [{"offset": offset+i, "original": x, "remake": y}
            for i, (x, y) in enumerate(zip(a, b)) if x != y]


def compare(original, remake, initial, ticks, army):
    rows = []
    for tick in range(ticks+1):
        left = original / (initial if tick == 0 else f"{'' if initial == 'start' else initial+'-'}tick-{tick:03d}")
        right = remake / f"tick-{tick:03d}"
        a = read(left.with_name(left.name+".save.dat"), 4*22208)[:22208]
        b = read(right.with_name(right.name+".save.dat"), 4*22208)[:22208]
        ar = read(left.with_name(left.name+".rng.bin"), 258)
        br = read(right.with_name(right.name+".rng.bin"), 258)
        ao = read(left.with_name(left.name+".occupancy.bin"), 384*256)
        bo = read(right.with_name(right.name+".occupancy.bin"), 384*256)
        am = json.loads(left.with_name(left.name+".json").read_text())
        bm = json.loads(right.with_name(right.name+".json").read_text())
        base = 0x22C0+army*64
        army_a, army_b = bytearray(a[base:base+64]), bytearray(b[base:base+64])
        # 沿用corps_diff.py：bit4為繪圖旗標。+21是未建模的對峙動畫欄位。
        # 兩者的原始差異仍列於whole_block_differences，不稱完整記錄逐byte一致。
        army_a[0] &= 0xEF
        army_b[0] &= 0xEF
        army_a[0x21] = army_b[0x21] = 0
        player = a[0x0F]
        if player >= 22 or b[0x0F] != player:
            raise ValueError("玩家勢力不合法或兩側不同")
        funds = 0x80+player*64+0x20
        occupancy_diff = [(i, x&0x7F, y) for i, (x, y) in enumerate(zip(ao, bo)) if (x&0x7F) != y]
        primary = {
            "globals_0_to_3a": diff(a[:59], b[:59]),
            "player_army": diff(army_a, army_b, base),
            "player_funds": diff(a[funds:funds+3], b[funds:funds+3], funds),
            "rng": diff(ar, br),
            "occupancy_count": len(occupancy_diff),
            "occupancy_first": [{"x": i%384, "y": i//384, "original": x, "remake": y}
                                for i, x, y in occupancy_diff[:8]],
            "event_delay": [] if am["event_delay"] == bm["event_delay"] else [am["event_delay"], bm["event_delay"]],
        }
        whole = diff(a, b)
        # 起點另要求完整區塊一致，不能把中途範圍排除當成起點已對齊。
        passed = not any(primary.values()) and (tick != 0 or not whole)
        rows.append({"tick": tick, "passed": passed, "primary": primary,
                     "whole_block_difference_count": len(whole), "whole_block_differences": whole,
                     "original_block_sha256": hashlib.sha256(a).hexdigest(),
                     "remake_block_sha256": hashlib.sha256(b).hexdigest()})
    return {"passed": all(row["passed"] for row in rows), "army": army, "ticks": ticks,
            "checkpoints": len(rows), "occupancy_cells_per_checkpoint": 384*256,
            "rng_bytes_per_checkpoint": 258, "rows": rows,
            "scope": "globals0..3A, player army except renderer bit4 and animation +21, player funds, raw RNG, full occupancy low7, event delay",
            "whole_block_exact": all(not row["whole_block_difference_count"] for row in rows),
            "limits": "Does not assert other armies, object animation, full battle or full campaign parity."}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--original", type=Path, required=True)
    parser.add_argument("--remake", type=Path, required=True)
    parser.add_argument("--initial", choices=["start", "gate"], required=True)
    parser.add_argument("--ticks", type=int, required=True)
    parser.add_argument("--army", type=int, default=37)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    if not 1 <= args.ticks <= 10000 or not 0 <= args.army < 127:
        parser.error("ticks或army超出範圍")
    result = compare(args.original, args.remake, args.initial, args.ticks, args.army)
    args.output.write_text(json.dumps(result, ensure_ascii=False, indent=2)+"\n")
    bad = [row["tick"] for row in result["rows"] if not row["passed"]]
    print(json.dumps({"passed": result["passed"], "checkpoints": result["checkpoints"],
                      "first_failed_tick": bad[0] if bad else None,
                      "whole_block_exact": result["whole_block_exact"]}, ensure_ascii=False))
    return 0 if result["passed"] else 1


if __name__ == "__main__":
    raise SystemExit(main())
