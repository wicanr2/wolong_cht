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
    """在 guest 的常規記憶體裡搜簽章，回傳 (segment, offset)。

    一次讀 4 KB，640 KB 共 160 次請求。**找不到要能跟「沒跑到」分開**，
    所以最後印掃了幾個 byte。
    """
    needle = signature[:probe]
    scanned = 0
    for base in range(0, 0xA0000, 0x1000):
        seg, off = base >> 4, 0
        try:
            chunk = bridge.read(seg, off, 0x1000 + probe)
        except Exception as e:
            if "DEBUGGER_NOT_STOPPED" in str(e):
                raise RuntimeError(
                    "除錯器沒有停住——先送 execution.pause。"
                    "（沒有這一句，掃描會安靜地全部落空，看起來像「表不在記憶體裡」）")
            continue
        scanned += len(chunk)
        idx = chunk.find(needle)
        if idx >= 0:
            linear = base + idx
            return linear, scanned
    return None, scanned


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
    return 0


if __name__ == "__main__":
    sys.exit(main())

# `sample` 的欄位分佈快照。**零值要能跟「前提沒滿足」分開**：
# 開局沒有任何軍團，+0x18 與 +0x14 一定全 0，這時「相符」不構成證據。
def histogram(cities, off, label):
    from collections import Counter
    h = Counter(c[off] for c in cities)
    print("%s（+0x%02X）：%s" % (label, off, dict(sorted(h.items()))))
