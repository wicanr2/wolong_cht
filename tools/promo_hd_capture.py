"""在獨立 Xvfb 錄製目前桌面版，透過滑鼠操作 HD 和戰場縮放。"""
import argparse
import json
import os
from pathlib import Path
import re
import subprocess
import time


def run(*args):
    return subprocess.check_output(["xdotool", *map(str, args)], text=True).strip()


def capture(binary, out, name, packs, extra, duration, actions):
    trace = out / f"{name}.jsonl"
    trace.unlink(missing_ok=True)
    args = [str(binary), "-orig", "workplace/orig/dosv", "-font", "workplace/eten",
            "-seed", "1", "-direct", "-scenario", "0", "-player", "0",
            "-save-file", str(out / f"{name}.save"), *packs, *extra]
    with (out / f"{name}.log").open("w") as log:
        game = subprocess.Popen(args, env=dict(os.environ, WOLONG_DESKTOP_TRACE=str(trace)),
                                stdout=log, stderr=log)
        recording = None
        try:
            deadline = time.monotonic() + 45
            while time.monotonic() < deadline:
                try:
                    window = run("search", "--onlyvisible", "--pid", game.pid).splitlines()[-1]
                    break
                except (subprocess.CalledProcessError, IndexError):
                    assert game.poll() is None, name
                    time.sleep(.15)
            else:
                raise RuntimeError(f"找不到遊戲視窗：{name}")
            run("windowfocus", window)

            def state():
                try:
                    with trace.open("rb") as stream:
                        stream.seek(max(0, trace.stat().st_size - 32768))
                        rows = stream.read().splitlines()
                    for row in reversed(rows):
                        try:
                            return json.loads(row)
                        except json.JSONDecodeError:
                            continue
                except FileNotFoundError:
                    pass
                return {}

            def wait(predicate, label, seconds=15):
                until = time.monotonic() + seconds
                while time.monotonic() < until:
                    s = state()
                    if s and predicate(s):
                        return s
                    if game.poll() is not None:
                        break
                    time.sleep(.06)
                raise RuntimeError(f"{name}: {label}")

            def move(x, y):
                for _ in range(70):
                    s = state()
                    dy, dx = max(y, 1-s["barHeight"]) - s["y"], x-s["x"]
                    if abs(dx) <= 1 and abs(dy) <= 1:
                        return
                    if s["captured"]:
                        run("mousemove_relative", "--", max(-24, min(24, dx))*2, max(-24, min(24, dy))*2)
                    else:
                        run("mousemove", "--window", window, x*2, (y+s["barHeight"])*2)
                    wait(lambda t: t["pointerRaw"] != s["pointerRaw"] or
                         (t["x"], t["y"], t["captured"]) != (s["x"], s["y"], s["captured"]), "滑鼠位移", 20)
                raise RuntimeError(f"{name}: 滑鼠定位失敗")

            def click():
                run("mousedown", 1)
                wait(lambda t: t["buttons"][0], "左鍵按下")
                run("mouseup", 1)
                wait(lambda t: not t["buttons"][0], "左鍵放開")

            def toolbar(key):
                s = state()
                if not s["barVisible"]:
                    move(320, 2-s["barHeight"])
                    s = wait(lambda t: t["barVisible"] and t["pointerReady"], "工具列展開")
                x0, y0, x1, y1 = s["toolbar"][key]
                move((x0+x1)//2, (y0+y1)//2-s["barHeight"])
                click()

            wait(lambda s: s.get("pointerReady") and s.get("world") and s.get("hdBattleAvailable")
                 and (name == "map" or s.get("battle") is not None), "遊戲就緒", 45)
            if name != "map":
                wait(lambda s: s.get("battle") and not s["battle"]["opening"] and not s["battle"]["talk"],
                     "戰場開場結束", 50)
            if name == "siege":
                toolbar("display")
                wait(lambda s: s["hdBattleDrawn"], "HD 戰場繪製完成", 30)
                move(310, 240)
            geometry = subprocess.check_output(["xwininfo", "-id", window], text=True)
            dims = [int(re.search(pattern, geometry)[1]) for pattern in
                    (r"Absolute upper-left X:\s+(\d+)", r"Absolute upper-left Y:\s+(\d+)",
                     r"Width:\s+(\d+)", r"Height:\s+(\d+)")]
            x, y, w, h = dims
            clip = out / f"{name}.mp4"
            recording = subprocess.Popen(["ffmpeg", "-y", "-nostdin", "-loglevel", "error",
                "-f", "x11grab", "-draw_mouse", "0", "-framerate", "30", "-video_size", f"{w}x{h}",
                "-i", f"{os.environ['DISPLAY']}+{x},{y}", "-t", str(duration),
                "-an", "-c:v", "libx264", "-threads", "2", "-preset", "ultrafast", "-pix_fmt", "yuv420p", str(clip)])
            start = time.monotonic()
            evidence = []
            for seconds, key, expectation in actions:
                time.sleep(max(0, start+seconds-time.monotonic()))
                if key == "charge":
                    move(496, 280)
                    click()
                else:
                    toolbar(key)
                s = wait(expectation, key)
                evidence.append({"seconds": round(time.monotonic()-start, 3), "action": key,
                                 **{k: s.get(k) for k in ("hdDisplay", "hdBattleDrawn", "battleZoomPercent", "worldZoomPercent")}})
                move(310, 240)
            assert recording.wait(timeout=duration+20) == 0, name
            if name == "siege":
                assert max(action["seconds"] for action in evidence) < duration-2, evidence
            (out / f"{name}-capture.json").write_text(json.dumps({"args":args,"geometry":dims,"actions":evidence}, indent=2)+"\n")
            print(name, duration, evidence, flush=True)
        finally:
            for process in (recording, game):
                if process is not None and process.poll() is None:
                    process.terminate()
                    process.wait(timeout=8)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("binary", type=Path)
    parser.add_argument("output", type=Path)
    parser.add_argument("battles", type=Path)
    parser.add_argument("--scene", choices=("map", "siege", "field"))
    args = parser.parse_args()
    args.output.mkdir(exist_ok=True)
    art = Path("workplace/hd-redraw-local-20261104/faithful")
    assets = ["-hd-battle", str(args.battles), "-hd-world", str(art/"world"),
              "-hd-portraits", str(art/"portraits"), "-hd-locations", str(art), "-hd-banner", str(art)]
    if args.scene in (None, "map"):
        capture(args.binary, args.output, "map", assets, ["-speed", "0", "-cam", "206,114"], 11, [
            (1.6, "display", lambda s: s["hdDisplay"]),
        ])
    if args.scene in (None, "siege"):
        capture(args.binary, args.output, "siege", assets,
                ["-open-siege", "-siege-node", "76", "-battle-steps", "290", "-battle-cam", "36,14", "-tactical-speed", "0"],
                17, [(1.3, "charge", lambda s: s["battle"]["leaderOrder"][0] == 2),
                     (6, "zoomOut", lambda s: s["battleZoomPercent"] == 75 and s["hdBattleDrawn"])])
    if args.scene in (None, "field"):
        capture(args.binary, args.output, "field", assets,
                ["-open-battle", "-battle-steps", "105", "-battle-cam", "20,15", "-tactical-speed", "0"],
                11, [(1.6, "display", lambda s: s["hdBattleDrawn"])])


if __name__ == "__main__":
    main()
