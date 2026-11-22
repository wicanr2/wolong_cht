"""以封包從乾淨設定目錄開啟遊戲並截圖。"""
import argparse
from collections import deque
import json
import os
from pathlib import Path
import subprocess
import time


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("platform", choices=("linux", "windows"))
    parser.add_argument("launcher", type=Path)
    args = parser.parse_args()
    work = Path("/out/smoke")
    work.mkdir(exist_ok=True)
    out = Path("/delivery/smoke")
    trace = work / f"{args.platform}.jsonl"
    trace.unlink(missing_ok=True)
    env = dict(os.environ, XDG_DATA_HOME=str(work/"data"), XDG_CONFIG_HOME=str(work/"config"),
               WOLONG_DESKTOP_TRACE=str(trace))
    if args.platform == "windows":
        env["WOLONG_DESKTOP_TRACE"] = "Z:"+str(trace).replace("/", "\\")
        command = ["wine", "cmd", "/c", args.launcher.name]
        cwd = args.launcher.parent
    else:
        command = [str(args.launcher)]
        cwd = work
        env["APPIMAGE_EXTRACT_AND_RUN"] = "1"
    with (out / f"{args.platform}-launcher.log").open("w") as log:
        game = subprocess.Popen(command, cwd=cwd, env=env, stdout=log, stderr=log)
        try:
            deadline = time.monotonic()+100
            state = {}
            while time.monotonic() < deadline:
                if trace.is_file():
                    with trace.open(encoding="utf-8") as stream:
                        rows = deque(stream, maxlen=6)
                    for row in reversed(rows):
                        try:
                            state = json.loads(row)
                            break
                        except json.JSONDecodeError:
                            continue
                if state.get("phase") == 0 and state.get("hdBattleAvailable"):
                    assert not state["hdDisplay"] and not state["hdBattleDrawn"]
                    break
                if game.poll() is not None:
                    raise RuntimeError(f"launcher exited {game.returncode}")
                time.sleep(.25)
            else:
                raise RuntimeError("launcher did not report the opening screen")
            output = out / f"{args.platform}-startup.png"
            for _ in range(70):
                windows = subprocess.run(["xdotool", "search", "--onlyvisible", "--name", ".*"],
                                         capture_output=True, text=True).stdout.splitlines()
                visible = []
                for window in windows:
                    geom = subprocess.run(["xdotool", "getwindowgeometry", "--shell", window],
                                          capture_output=True, text=True).stdout.splitlines()
                    if "WIDTH=1280" in geom and any(row in geom for row in ("HEIGHT=800", "HEIGHT=812", "HEIGHT=856")):
                        visible.append(window)
                if visible:
                    break
                time.sleep(.2)
            else:
                raise RuntimeError("the desktop game window did not appear")
            subprocess.run(["xdotool", "windowfocus", visible[-1], "mousemove", "--window",
                            visible[-1], "640", "330"], check=True)
            for _ in range(60):
                subprocess.run(["import","-window","root",str(output)], check=True)
                colors = int(subprocess.check_output(["identify","-format","%k",str(output)]))
                if colors > 8:
                    break
                time.sleep(.2)
            else:
                raise RuntimeError("the startup window never drew visible game graphics")
            data = {key:state.get(key) for key in ("phase","hdBattleAvailable","hdDisplay","hdBattleDrawn")}
            data.update(startup=True, colors=colors, native_gui=args.platform=="linux",
                        environment="Linux Xvfb" if args.platform=="linux" else "Wine 9.0 / Xvfb",
                        launcher=args.launcher.name)
            (out / f"{args.platform}.json").write_text(json.dumps(data,indent=2)+"\n")
            print(data, flush=True)
        finally:
            if args.platform == "windows":
                subprocess.run(["wineserver","-k"], env=env, timeout=10, check=False)
            if game.poll() is None:
                game.terminate()
            try:
                game.wait(timeout=10)
            except subprocess.TimeoutExpired:
                game.kill()
                game.wait()


if __name__ == "__main__":
    main()
