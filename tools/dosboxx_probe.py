#!/usr/bin/env python3
"""在跑著的 DOSBox-X 裡定位據點表，取樣驗證靜態解出來的欄位。

    tools/dosboxx_bridge.sh start
    python3 tools/dosboxx_probe.py locate          # 找據點表在 guest 記憶體哪裡
    python3 tools/dosboxx_probe.py sample <seg:off>

⚠ 這一支**直接講 bridge 的 TCP 協定**（newline-delimited JSON，
`127.0.0.1:9876`），不經過 MCP 那一層。兩者背後是同一個 `debug_ai.cpp`，
差別只在誰發請求：MCP 版由 agent 在對話中呼叫（要重啟 session 才載入），
這一支是可重跑的腳本，適合放進驗收流程。

## 怎麼找到表

**不用猜 DS。** 據點表是從 `SINARIO.DAT` 的劇本區塊整段載進來的，
所以拿本機解析出來的前幾筆記錄當簽章，在 guest 記憶體裡搜就好——
搜得到就是它，搜不到就是還沒載入（或劇本選的不是這一個）。
這比「假設 DS ＝ 某個值」可靠，也不必跨版本外推位址
（`CLAUDE.md` §7 第 9 條：PC-98 與 DOS/V 的程式碼位址不能互推，
但**資料的版面是同一份原始碼編出來的**）。

## 取樣要驗什麼

`docs/re/44` 這一輪解出來的三條，都是讀記憶體就能驗、不需要斷點的：

    +0x18  停在該據點那一格的軍團數（每 tick 從佔用圖抄回來）
    +0x14  鄰接敵方據點的 +0x18 總和
    +0x00  低 4 位 ＝ 哪幾個鄰接槽屬於別的勢力（**這一條推翻了舊結論**）

第三條最需要動態證據：靜態只看得到寫入端的形狀，
而「據點換手之後這四位真的會變」只有跑起來才看得到。
"""
import json
import socket
import sys

HOST, PORT = "127.0.0.1", 9876
CITY_SIZE = 32
NUM_CITIES = 192


class Bridge:
    def __init__(self):
        self.sock = socket.create_connection((HOST, PORT), timeout=10)
        self.buf = b""
        self.id = 0

    def request(self, method, params=None):
        self.id += 1
        payload = {"id": self.id, "method": method}
        if params:
            payload["params"] = params
        self.sock.sendall(json.dumps(payload).encode() + b"\n")
        while b"\n" not in self.buf:
            chunk = self.sock.recv(65536)
            if not chunk:
                raise RuntimeError("bridge 關掉了連線")
            self.buf += chunk
        line, self.buf = self.buf.split(b"\n", 1)
        reply = json.loads(line)
        if "error" in reply and reply["error"]:
            raise RuntimeError("bridge 回錯誤：%s" % reply["error"])
        return reply.get("result", reply)

    def read(self, seg, off, length):
        r = self.request("memory.read",
                         {"address": "%04X:%04X" % (seg, off), "length": length})
        data = r.get("bytes") if r.get("bytes") is not None else r.get("data")
        if isinstance(data, str):
            return bytes.fromhex(data.replace(" ", ""))
        # bridge 回的是十六進位字串的陣列（["00","1F",…]），不是整數陣列。
        if data and isinstance(data[0], str):
            return bytes(int(x, 16) for x in data)
        return bytes(data)


def scenario_cities(path, block=0, block_size=0x5400, city_base=0x08C0):
    """從本機的 SINARIO.DAT 取劇本的據點表，當搜尋簽章。"""
    raw = open(path, "rb").read()
    b = raw[block * block_size:(block + 1) * block_size]
    return b[city_base:city_base + NUM_CITIES * CITY_SIZE]


def locate(bridge, signature, probe=64):
    """在 guest 的常規記憶體裡找據點表，回傳 (linear, 掃過的 bytes)。

    ⚠ **不能拿劇本檔的原始 bytes 當簽章。** 遊戲一跑起來，上昇值、防災值、
    城兵、所屬、+0x14／+0x18 全都在變，整段比對必然落空——第一版就是這樣
    在遊戲跑了三分鐘之後回「沒找到」，而表明明就在那裡。

    改用**執行期不會變的欄位**：X 與 Y（`+0x08`／`+0x0A`）。
    連續 16 筆的座標全中才算，等於 32 個 u16 同時對上，誤中機率可以忽略。
    """
    want = []
    for i in range(16):
        r = signature[i * CITY_SIZE:(i + 1) * CITY_SIZE]
        want.append((r[8] | (r[9] << 8), r[10] | (r[11] << 8)))

    mem = bytearray()
    for base in range(0, 0xA0000, 0x1000):
        try:
            mem += bridge.read(base >> 4, 0, 0x1000)
        except Exception as e:
            if "DEBUGGER_NOT_STOPPED" in str(e):
                raise RuntimeError(
                    "除錯器沒有停住——先送 execution.pause。"
                    "（沒有這一句，掃描會安靜地全部落空，"
                    "看起來像「表不在記憶體裡」）")
            mem += b"\x00" * 0x1000

    need = 16 * CITY_SIZE
    for off in range(0, len(mem) - need, 2):
        hit = True
        for k, (wx, wy) in enumerate(want):
            b0 = off + k * CITY_SIZE
            if (mem[b0 + 8] | (mem[b0 + 9] << 8)) != wx or \
               (mem[b0 + 10] | (mem[b0 + 11] << 8)) != wy:
                hit = False
                break
        if hit:
            return off, len(mem)
    return None, len(mem)


def sample(bridge, linear):
    """把 192 筆據點記錄讀回來，檢查三條不變量。"""
    seg, off = linear >> 4, linear & 0x0F
    data = b""
    remaining = NUM_CITIES * CITY_SIZE
    cur = linear
    while remaining > 0:
        n = min(0x1000, remaining)
        data += bridge.read(cur >> 4, cur & 0x0F, n)
        cur += n
        remaining -= n
    cities = [data[i * CITY_SIZE:(i + 1) * CITY_SIZE] for i in range(NUM_CITIES)]

    bad_mask = bad_threat = 0
    for i, c in enumerate(cities):
        owner = c[0x01]
        mask = c[0x00] & 0x0F
        want_mask = 0
        threat = 0
        for k in range(4):
            n = c[0x1C + k]
            if n >= NUM_CITIES:
                continue
            if cities[n][0x01] != owner:
                want_mask |= 1 << k
                threat += cities[n][0x18]
        if mask != want_mask:
            bad_mask += 1
            if bad_mask <= 5:
                print("  據點 %3d：+0x00 低 4 位 ＝ %04b，由鄰居算出來是 %04b"
                      % (i, mask, want_mask))
        if c[0x1B] != bin(mask).count("1"):
            print("  據點 %3d：+0x1B ＝ %d，而遮罩有 %d 位"
                  % (i, c[0x1B], bin(mask).count("1")))
    print("+0x00 低 4 位與鄰居實況不符：%d / %d" % (bad_mask, NUM_CITIES))
    return cities


def main():
    cmd = sys.argv[1] if len(sys.argv) > 1 else "locate"
    b = Bridge()
    print("debug.status:", b.request("debug.status"))
    if cmd == "locate":
        sig = scenario_cities(sys.argv[2] if len(sys.argv) > 2
                              else "workplace/orig/pc98/SINARIO.DAT")
        linear, scanned = locate(b, sig)
        if linear is None:
            print("沒找到（掃了 %d bytes）——據點表可能還沒載入，"
                  "或這個劇本不是區塊 0" % scanned)
            return 1
        print("據點表在 linear %06X（%04X:%04X），掃了 %d bytes"
              % (linear, linear >> 4, linear & 0x0F, scanned))
    elif cmd == "sample":
        sample(b, int(sys.argv[2], 16))
    elif cmd == "verify":
        return 0 if verify(b, int(sys.argv[2], 16)) else 2
    return 0



# `sample` 的欄位分佈快照。**零值要能跟「前提沒滿足」分開**：
# 開局沒有任何軍團，+0x18 與 +0x14 一定全 0，這時「相符」不構成證據。
def histogram(cities, off, label):
    from collections import Counter
    h = Counter(c[off] for c in cities)
    print("%s（+0x%02X）：%s" % (label, off, dict(sorted(h.items()))))


CORPS_BASE_IN_SEG = 0x2240   # 軍團表：127 × 64 B
CITY_BASE_IN_SEG = 0x0840    # 據點表：192 × 32 B
CORPS_SIZE, NUM_CORPS = 64, 127
FRIEND_BASE_IN_SEG = 0x0600  # 交友度表：22 列 × 24 B
FRIEND_STRIDE, NEUTRAL, PEACE_BIT = 24, 0x18, 0x80


def read_block(bridge, linear, length):
    out = b""
    cur = linear
    while length > 0:
        n = min(0x1000, length)
        out += bridge.read(cur >> 4, cur & 0x0F, n)
        cur += n
        length -= n
    return out


def verify(bridge, city_linear):
    """拿軍團表對據點的 +0x18／+0x14，補上 §4.2 缺的那兩條證據。

    段基址由據點表反推：`ds` 基底 ＝ 據點表 linear − 0x840。
    軍團表在同一段的 0x2240（`docs/re/44`）。**這是資料段的版面，
    不是程式碼位址**，所以不算跨版本外推——而且下面會自己驗：
    軍團記錄要嘛全零、要嘛欄位落在合理範圍，兩者都不成立就表示推錯了。
    """
    seg_base = city_linear - CITY_BASE_IN_SEG
    cities = [read_block(bridge, city_linear, NUM_CITIES * CITY_SIZE)
              [i * CITY_SIZE:(i + 1) * CITY_SIZE] for i in range(NUM_CITIES)]
    craw = read_block(bridge, seg_base + CORPS_BASE_IN_SEG, NUM_CORPS * CORPS_SIZE)
    corps = [craw[i * CORPS_SIZE:(i + 1) * CORPS_SIZE] for i in range(NUM_CORPS)]

    friend = read_block(bridge, seg_base + FRIEND_BASE_IN_SEG, 22 * FRIEND_STRIDE)
    alive = [c for c in corps if c[0x00] >= 0x80]
    print("活著的軍團：%d 支" % len(alive))
    if not alive:
        print("⚠ 一支軍團都沒有——這一輪驗不了 +0x18／+0x14。"
              "**全 0 相符不構成證據**，先讓遊戲跑到有軍團在動。")
        return False
    for c in alive[:5]:
        print("   勢力 %2d  座標 (%3d,%3d)  節點 %3d  兵力 %5d"
              % (c[0x01], c[0x10] | (c[0x11] << 8), c[0x12] | (c[0x13] << 8),
                 (c[0x0E] | (c[0x0F] << 8)) // 8, c[0x04] | (c[0x05] << 8)))

    occ = {}
    for c in alive:
        occ[(c[0x10] | (c[0x11] << 8), c[0x12] | (c[0x13] << 8))] = \
            occ.get((c[0x10] | (c[0x11] << 8), c[0x12] | (c[0x13] << 8)), 0) + 1

    bad18 = bad14 = 0
    nonzero18 = nonzero14 = 0
    for i, c in enumerate(cities):
        x = c[0x08] | (c[0x09] << 8)
        y = c[0x0A] | (c[0x0B] << 8)
        want18 = occ.get((x, y), 0)
        if c[0x18] != want18:
            bad18 += 1
            if bad18 <= 5:
                print("  據點 %3d：+0x18 ＝ %d，實際站著 %d 支"
                      % (i, c[0x18], want18))
        nonzero18 += c[0x18] != 0
        owner = c[0x01]
        # 中立據點根本不做威脅判斷——`sub_13EFD` 的
        # `cmp byte [si+841h], 18h / jz` 直接跳過 `sub_13F74`，
        # 所以它們的 +0x14 永遠是 0。驗證器不跳過的話，
        # 這幾筆會變成穩定的假不符（實測就是最後剩下的那四筆）。
        if owner == NEUTRAL:
            continue
        # ⚠ 威脅量不是「所有非我方鄰居」的和。`sub_13FA9` 有兩道濾網，
        # 第一版的驗證器沒照抄，於是穩定出現「算出來 1、實際 0」的假不符：
        #   中立（0x18）不算威脅——原版是 `jz` 跳過累加那一段
        #   交友度 bit 7 ＝ 和平 → 不算威脅
        want14 = 0
        for n in c[0x1C:0x20]:
            if n >= NUM_CITIES:
                continue
            no = cities[n][0x01]
            if no == owner or no == NEUTRAL:
                continue
            if owner < 22 and no < 22 and friend[owner * FRIEND_STRIDE + no] & PEACE_BIT:
                continue
            want14 += cities[n][0x18]
        if c[0x14] != min(want14, 255):
            bad14 += 1
            if bad14 <= 5:
                print("  據點 %3d：+0x14 ＝ %d，由鄰居算出來是 %d"
                      % (i, c[0x14], want14))
        nonzero14 += c[0x14] != 0
    print("+0x18 不符：%d / %d（非零的有 %d 筆）" % (bad18, NUM_CITIES, nonzero18))
    print("+0x14 不符：%d / %d（非零的有 %d 筆）" % (bad14, NUM_CITIES, nonzero14))
    return bad18 == 0 and bad14 == 0


if __name__ == "__main__":
    sys.exit(main())
