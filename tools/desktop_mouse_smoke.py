#!/usr/bin/env python3
"""在已啟動 Xvfb 的 Docker 內，以真實滑鼠驗證桌面工具列與正常行軍。

python3 tools/desktop_mouse_smoke.py BINARY ORIG FONT OUTPUT
OUTPUT 必須存在。需要 xdotool 與 ImageMagick；不修改原版資料。
"""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import time


def run(binary, orig, font, output, mode="desktop"):
    assert mode in {"desktop", "advise", "march", "battle", "attack", "capture", "defend", "field", "water", "windows"}
    water = mode == "water"
    field_mode = mode in {"field", "water"}
    attacking = mode in {"attack", "capture", "defend", "field", "water"}
    attack_faction, attack_node, attack_xy = (10, 122, (304, 158)) if water else (12, 100, (216, 137))
    enemy_army = 72 if water else 3
    field_link = 0x1220 if water else 0x10B0
    field_target_xy = (320, 152) if water else (221, 144)
    output = Path(output)
    assert output.is_dir()
    trace = output / "mouse-toolbar-trace.jsonl"
    trace.unlink(missing_ok=True)
    checks = []
    log = (output / "mouse-toolbar-gui.log").open("w")
    with tempfile.TemporaryDirectory(prefix="wolong-mouse-") as temp:
        save = Path(temp) / "SAVE.DAT"
        args = [binary, "-orig", orig, "-font", font, "-seed", "1", "-save-file", str(save)]
        proc = subprocess.Popen(args, env=dict(os.environ, WOLONG_DESKTOP_TRACE=str(trace)), stdout=log, stderr=log)

        def command(*args):
            return subprocess.check_output(["xdotool", *map(str, args)], text=True, timeout=10).strip()

        def state():
            for _ in range(40):
                try:
                    return json.loads(trace.read_text().splitlines()[-1])
                except (OSError, IndexError, json.JSONDecodeError):
                    if proc.poll() is not None:
                        raise RuntimeError("game exited; inspect mouse-toolbar-gui.log")
                    time.sleep(.05)
            raise RuntimeError("missing desktop trace")

        def wait(predicate, label, seconds=8):
            end = time.monotonic() + seconds
            while time.monotonic() < end:
                s = state()
                if predicate(s):
                    checks.append(label)
                    print("PASS", label, flush=True)
                    return s
                time.sleep(.05)
            last_state = state()
            print("FOCUS", window, command("getwindowfocus"), flush=True)
            if label in {"field opening completes normally", "army moves while cursor rests on toolbar"}:
                subprocess.run(["kill", "-QUIT", str(proc.pid)], check=False)
                time.sleep(.3)
            raise AssertionError((label, last_state))

        def move(x, y, tolerance=1):
            for _ in range(140):
                s = state()
                dx, dy = x-s["x"], y-s["y"]
                if abs(dx) <= tolerance and abs(dy) <= tolerance:
                    return
                if s["captured"]:
                    command("mousemove_relative", "--", max(-24, min(24, dx))*2, max(-24, min(24, dy))*2)
                else:
                    command("mousemove", "--window", window, x*2, (y+40)*2)
                # 等待這次位移被遊戲消費，避免負載較高時重送舊座標的修正量。
                end = time.monotonic()+5
                while True:
                    moved = state()
                    if moved["pointerRaw"] != s["pointerRaw"] or (moved["x"], moved["y"]) != (s["x"], s["y"]) or moved["captured"] != s["captured"]:
                        break
                    assert time.monotonic() < end, "mouse move was not observed"
                    time.sleep(.01)
            raise AssertionError(("pointer target", x, y, state()))

        def click(x, y, button=1, *, exact=False):
            move(x, y, 0 if exact else 1)
            observed = {1: 0, 3: 1, 2: 2}[button]  # X11：左中右；trace：左右中。
            command("mousedown", button)
            end = time.monotonic() + 5
            while not state()["buttons"][observed]:
                assert time.monotonic() < end, "mouse down was not observed"
                time.sleep(.02)
            command("mouseup", button)
            end = time.monotonic() + 5
            while state()["buttons"][observed]:
                assert time.monotonic() < end, "mouse up was not observed"
                time.sleep(.02)
            time.sleep(.05)

        def toolbar(index):
            click(50+90*index, -20)

        def band_move(dy, row, y, label):
            assert state()["choiceActive"]
            raw = state()["pointerRaw"]
            command("mousemove_relative", "--", 0, dy*2)
            return wait(lambda s: s["pointerRaw"] != raw and s["choiceActive"] and s["choiceRow"] == row and s["choiceY"] == y, label)

        def band_click(button=1):
            s = state()
            assert s["choiceActive"]
            click(s["x"], s["y"], button)

        def shot(name):
            if mode == "windows":
                time.sleep(.25)  # trace在Update寫入；等穩定Draw，避免擷取上一幀。
            subprocess.check_call(["import", "-window", window, str(output/(name+".png"))], timeout=10)

        try:
            wait(lambda s: s["phase"] == 0, "normal title")
            window = command("search", "--sync", "--onlyvisible", "--name", "臥龍傳").splitlines()[-1]
            command("windowfocus", window)
            time.sleep(.2)
            move(300, 300)
            shot("mouse-toolbar-title")
            command("key", "F9", "F10", "Escape", "Return", "space", "Right", "g")
            time.sleep(.3)
            assert state()["phase"] == 0 and not state()["quit"]
            checks.append("gameplay keys disabled at title")
            if mode == "windows":
                for _ in range(3):
                    toolbar(5)
                wait(lambda s: s["phase"] == 3 and s["factionCount"] > 10, "normal launcher faction list")
                shot("windows-launcher-initial")

                def drag_preview(x, y, top, label, committed, *, launcher=False):
                    move(x, y, 0)
                    key = "factionTop" if launcher else "listTop"
                    return wait(lambda s: s["listDragging"] and s["listPreviewTop"] == top and s[key] == committed,
                                label)

                move(144, 144, 0)
                command("mousedown", 1)
                wait(lambda s: s["listDragging"] and s["buttons"][0], "launcher drag starts inside track")
                drag_preview(144, 200, 6 if state()["factionCount"] == 22 else state()["factionCount"]//2-5,
                             "launcher previews while committed top stays zero", 0, launcher=True)
                wait(lambda s: s["cursorMode"] == 1, "launcher native cursor hidden during drag")
                shot("windows-launcher-preview")
                drag_preview(300, 272, state()["factionCount"]-10, "launcher drag owns pointer outside track", 0, launcher=True)
                command("mouseup", 1)
                wait(lambda s: not s["listDragging"] and s["factionTop"] == s["factionCount"]-10 and not s["factionSelected"],
                     "launcher release commits without selecting")
                wait(lambda s: s["cursorMode"] == 0, "launcher native cursor restored after release")
                move(144, 144, 0)
                command("mousedown", 1)
                wait(lambda s: s["listDragging"] and s["buttons"][0], "launcher second drag starts inside track")
                drag_preview(140, -20, 0, "launcher drag crosses toolbar", state()["factionCount"]-10, launcher=True)
                shot("windows-launcher-over-toolbar")
                command("mouseup", 1)
                wait(lambda s: not s["listDragging"] and s["factionTop"] == 0 and s["phase"] == 3 and s["language"] == "zh-hant",
                     "release over toolbar commits without toolbar action")
                click(200, 128)
                wait(lambda s: s["factionSelected"], "launcher first click highlights")
                click(144, 200)
                wait(lambda s: s["phase"] == 4, "selected launcher confirms original row over scrollbar")
                toolbar(5)
                wait(lambda s: s["phase"] == -1 and s["pointerReady"], "normal Cao Cao new game")
                click(352, 15)  # resident command strip
                click(384, 47)
                wait(lambda s: s["list"] and len(s["listRows"]) == 22 and s["listTop"] == 0,
                     "normal in-game faction list has original 22 rows")
                move(430, 280, 0)
                shot("windows-list-initial")
                command("key", "Escape", "Return", "Down", "Page_Down", "F9", "g")
                time.sleep(.2)
                assert state()["list"] and state()["listTop"] == 0 and state()["listPhase"] == 0
                checks.append("gameplay keys disabled in list")
                move(32, 128, 0)
                command("mousedown", 1)
                wait(lambda s: s["listDragging"] and s["buttons"][0], "in-game drag starts inside track")
                for y, top, name in [(128, 0, "start"), (152, 0, "quarter"), (184, 6, "middle"), (240, 12, "bottom")]:
                    drag_preview(32, y, top, "original drag sample: "+name, 0)
                    shot("windows-list-held-"+name)
                drag_preview(300, 240, 12, "in-game drag owns pointer outside track", 0)
                command("mousedown", 3)
                wait(lambda s: s["buttons"][1] and s["listDragging"], "right button cannot cancel active drag")
                command("mouseup", 1)
                wait(lambda s: not s["buttons"][0] and s["buttons"][1] and s["listDragging"] and s["listTop"] == 0,
                     "right hold keeps preview after left release")
                command("mouseup", 3)
                wait(lambda s: not s["listDragging"] and s["listTop"] == 12 and s["listPhase"] == 0,
                     "all buttons released commits original top12")
                move(430, 280, 0)
                shot("windows-list-released")
                click(430, 280, 3)
                wait(lambda s: not s["list"], "right click closes list after drag")
                toolbar(1)
                wait(lambda s: s["system"], "formation settings by toolbar")
                if state()["lordCorps"]:
                    click(376, 304)
                toolbar(4)
                click(208, 47)
                wait(lambda s: s["list"], "normal formation candidates")
                click(150, 112)
                click(150, 112)
                wait(lambda s: s["form"], "normal formation child window")
                shot("windows-formation")
                click(430, 280, 3)
                wait(lambda s: not s["form"] and s["list"], "right cancel restores parent list")
                for _ in range(3):
                    if state()["hostReady"]: break
                    toolbar(4)
                wait(lambda s: s["hostReady"], "formation cancellation returns to map")
                toolbar(2)
                click(300, 152)
                wait(lambda s: not s["save"] and save.exists(), "save after window sampling")
                toolbar(3)
                click(300, 152)
                wait(lambda s: s["hostReady"], "reload after window sampling")
                result = {"passed": True, "checks": checks, "mode": mode, "args": args,
                          "binary_sha256": hashlib.sha256(Path(binary).read_bytes()).hexdigest(),
                          "window_geometry": command("getwindowgeometry", "--shell", window), "final_state": state()}
                (output/"mouse-toolbar-gui.json").write_text(json.dumps(result, ensure_ascii=False, indent=2)+"\n")
                print(json.dumps(result, ensure_ascii=False), flush=True)
                return
            for language in ([] if field_mode else ["zh-hans", "ja", "en", "zh-hant"]):
                toolbar(0)
                wait(lambda s: s.get("language") == language and s["phase"] == 0, "language button: "+language)
                shot("mouse-toolbar-"+language)
            for _ in range(8):
                if water and state()["phase"] == 3:
                    click(200, 144)
                toolbar(5)
                if state()["phase"] == -1:
                    break
            wait(lambda s: s["phase"] == -1 and s["world"] and s["pointerReady"], "normal new game")
            if water:
                assert state()["player"] == 1, "normal water fixture must choose Sun Ce"
            if not field_mode:
                move(320, -40)
                original_cam = state()["cam"]
                command("mousemove_relative", "--", 0, -192)
                wait(lambda s: s["cam"][1] < original_cam[1], "top edge pans map through toolbar")
                stopped_cam = state()["cam"]
                time.sleep(.4)
                assert state()["cam"] == stopped_cam
                checks.append("map stops when mouse stops")
                toolbar(1)
                wait(lambda s: s["system"], "settings by pointer")
                command("key", "Escape", "Return", "F9", "F10", "Right", "plus", "g")
                time.sleep(.2)
                assert state()["system"] and not state()["quit"]
                checks.append("gameplay keys disabled in settings")
                toolbar(4)
                wait(lambda s: not s["system"], "toolbar back closes settings")
                toolbar(2)
                wait(lambda s: s["save"], "save panel by pointer")
                toolbar(0)
                assert state()["save"] and state()["language"] == "zh-hant"
                checks.append("disabled language button cannot escape save owner")
                toolbar(4)
                wait(lambda s: not s["save"], "toolbar back closes save")
                toolbar(2)
                click(300, 152)
                wait(lambda s: not s["save"] and save.exists(), "save first slot by pointer")
                toolbar(3)
                wait(lambda s: s["save"], "load panel by pointer")
                click(300, 152)
                wait(lambda s: not s["save"], "load first slot by pointer")
                if mode == "advise":
                    def advise_menu(row):
                        if state()["advise"] == 0:
                            click(48, 47)
                        wait(lambda s: s["advise"] == 1 and s["choiceActive"], "advise command menu")
                        for n in range(row):
                            band_move(32, n+1, 100, "advise band row "+str(n+1))
                        band_click()

                    click(352, 15)
                    capital_before = state()["capital"]
                    advise_menu(3)
                    wait(lambda s: s["advise"] == 5 and s["list"], "relocation city list")
                    click(150, 128)
                    wait(lambda s: s["list"] and s["listCursor"] == 1 and s["listPhase"] == 1, "relocation first click selects noninitial row")
                    shot("advise-capital-selected")
                    move(150, 144)
                    assert state()["listCursor"] == 1 and state()["listPhase"] == 1
                    checks.append("moving pointer keeps originally highlighted city")
                    click(150, 144, 3)
                    wait(lambda s: s["list"] and s["listPhase"] == 0, "right click first clears city highlight")
                    click(150, 144, 3)
                    wait(lambda s: s["advise"] == 0 and s["hostReady"] and not s["list"], "second right click leaves relocation and returns to map")
                    assert state()["capital"] == capital_before
                    checks.append("cancelled city selection leaves capital unchanged")
                    advise_menu(3)
                    wait(lambda s: s["advise"] == 5 and s["list"], "reopen relocation city list")
                    expected_city = state()["listRows"][1]
                    click(150, 128)
                    click(150, 128)
                    wait(lambda s: not s["list"] and s["advise"] == 6 and s["adviseNode"] == expected_city, "second city click submits selected city")
                    shot("advise-relocation-verdict")
                    for _ in range(8):
                        s = state()
                        if s["advise"] == 0 and not s["messages"]:
                            break
                        click(300, 300, 3)
                    wait(lambda s: s["hostReady"], "relocation verdict returns to map with right clicks")
                    advise_menu(0)
                    wait(lambda s: s["advise"] == 3 and s["list"], "hostility target list")
                    top = state()["listTop"]
                    click(32, 255)
                    wait(lambda s: s["listTop"] == top+1, "advise target scrollbar moves one row")
                    click(140, 96)
                    wait(lambda s: s["advise"] == 3 and s["list"] and s["listPhase"] == 0, "advise target header sorts without confirming")
                    s = state()
                    expected_target = s["listRows"][s["listTop"]+1]
                    click(150, 128)
                    wait(lambda s: s["list"] and s["listPhase"] == 1, "first target click highlights after scrolling and sorting")
                    shot("advise-target-selected")
                    click(150, 128)
                    wait(lambda s: not s["list"] and s["advise"] == 4 and s["adviseTarget"] == expected_target, "second target click starts correct proposal")
                    shot("advise-proposal-opening")
                    before_queue = state()["adviseQueue"]
                    assert before_queue > 0
                    click(300, 300, 3)
                    wait(lambda s: s["adviseQueue"] < before_queue, "right click advances proposal dialogue")
                    for _ in range(8):
                        s = state()
                        if s["advise"] == 0:
                            break
                        if s["choiceActive"]:
                            band_click(3)
                        else:
                            click(300, 300, 3)
                    wait(lambda s: s["hostReady"], "proposal dialogue returns to map")
                    capital_after = state()["capital"]
                    toolbar(2)
                    click(300, 152)
                    wait(lambda s: not s["save"] and save.exists(), "save after advise")
                    toolbar(3)
                    click(300, 152)
                    wait(lambda s: s["hostReady"] and s["capital"] == capital_after, "reload preserves capital after advise")
                    result = {"passed": True, "checks": checks, "mode": mode, "args": args,
                              "binary_sha256": hashlib.sha256(Path(binary).read_bytes()).hexdigest(),
                              "city": expected_city, "target": expected_target,
                              "capital_before": capital_before, "capital_after": capital_after,
                              "final_state": state()}
                    (output/"mouse-toolbar-gui.json").write_text(json.dumps(result, ensure_ascii=False, indent=2)+"\n")
                    print(json.dumps(result, ensure_ascii=False), flush=True)
                    return
            if mode in {"march", "battle"} or attacking:
                toolbar(1)
                wait(lambda s: s["system"], "march setup opens settings")
                if state()["lordCorps"]:
                    click(376, 304)
                wait(lambda s: not s["lordCorps"], "original formation candidates by settings")
                wanted_speed = 0 if mode == "battle" else 2
                for _ in range(5):
                    if state()["speed"] == wanted_speed:
                        break
                    click(376, 232)
                wait(lambda s: s["speed"] == wanted_speed, "strategy speed by pointer: "+str(wanted_speed))
                if mode == "battle" or attacking:
                    for _ in range(5):
                        if state()["tacticalSpeed"] == 0:
                            break
                        click(376, 256)
                    wait(lambda s: s["tacticalSpeed"] == 0, "highest tactical speed by pointer")
                toolbar(4)
                wait(lambda s: not s["system"], "march setup returns to map")
            # 原版熱區：開指令列、編成、第一列兩段選取、確認編成。
            click(352, 15)
            before_corps = {c[0] for c in state().get("corps", [])}
            click(208, 47)
            wait(lambda s: s["list"], "formation candidates by pointer")
            click(150, 112)
            click(150, 112)
            wait(lambda s: s["form"], "formation editor by pointer")
            toolbar(5)
            wait(lambda s: bool(s.get("corps")), "confirm formation by toolbar")
            formed = [c[0] for c in state()["corps"] if c[0] not in before_corps]
            assert len(formed) == 1, ("formation did not create one army", formed)
            for _ in range(8):
                if not state()["messages"]:
                    break
                toolbar(5)
            for _ in range(4):
                if not state()["canBack"]: break
                toolbar(4)
            wait(lambda s: s["hostReady"], "formation returns to map")
            shot("mouse-toolbar-formed")
            if attacking:
                # 第一章真實資料：袁術12／汝南100。
                # 對張繡的開局提案會先被君主拒絕，不能當作攻城前提。
                click(48, 47)
                wait(lambda s: s["advise"] == 1 and s["choiceActive"], "attacker opens hostility proposal")
                band_click()
                wait(lambda s: s["advise"] == 3 and s["list"], "attacker chooses hostile faction by list")
                target_row = state()["listRows"].index(attack_faction)
                while target_row-state()["listTop"] > 8:
                    top = state()["listTop"]
                    click(32, 255)
                    wait(lambda s: s["listTop"] == top+1, "scroll to hostile faction")
                row_y = 112+16*(target_row-state()["listTop"])
                click(150, row_y)
                click(150, row_y)
                wait(lambda s: s["advise"] == 4 and s["adviseTarget"] == attack_faction, "hostility proposal targets faction"+str(attack_faction))
                for _ in range(8):
                    s = state()
                    if s["choiceActive"] or s["advise"] == 0:
                        break
                    click(300, 300, 3)
                if state()["choiceActive"]:
                    wait(lambda s: s["advise"] == 4, "lord requests a reason for hostility")
                    band_move(32, 1, 100, "select we are stronger reason")
                    band_click()
                for _ in range(10):
                    if state()["advise"] == 0 and not state()["messages"]:
                        break
                    click(300, 300, 3)
                wait(lambda s: s["hostReady"] and attack_faction in s["wars"], "normal advice declares war on faction"+str(attack_faction))
                shot("attack-war-declared")
                if mode == "attack":
                    toolbar(1)
                    for _ in range(5):
                        if state()["speed"] == 0:
                            break
                        click(376, 232)
                    toolbar(4)
                    move(500, -20)
                    end = time.monotonic()+120
                    while True:
                        s = state()
                        enemy = next((c for c in s["corps"] if c[0] == 3), None)
                        route = next((c for c in s["corpsRoutes"] if c[0] == 3), None)
                        if enemy and enemy[1:4] == [*attack_xy, 4] and route[1:3] == [attack_node, 0]:
                            break
                        assert tuple(s["date"][:3]) < (196, 4, 20), ("no natural defender at Runan before April20", s)
                        assert time.monotonic() < end, ("waiting for natural defender", s)
                        if s["messages"]:
                            toolbar(5)
                        else:
                            time.sleep(.03)
                    checks.append("natural enemy corps3 guards Runan100 before player march")
                    shot("attack-natural-defender")
                    toolbar(1)
                    for _ in range(5):
                        if state()["speed"] == 2:
                            break
                        click(376, 232)
                    toolbar(4)
            if mode == "battle":
                move(500, -20)
                end = time.monotonic()+120
                while not state().get("battle"):
                    s = state()
                    assert tuple(s["date"][:3]) < (196, 7, 1), ("no natural battle before July", s)
                    assert time.monotonic() < end, ("normal battle fixture time limit", s)
                    if s["messages"]:
                        toolbar(5)
                    elif s["choiceActive"]:
                        band_click()
                    else:
                        time.sleep(.03)
                entry = state()
                checks.append("natural encounter enters tactical battle")
                shot("battle-normal-entry")
                wait(lambda s: s.get("battle") and not s["battle"]["opening"] and not s["battle"]["talk"] and not s["battle"]["done"], "battle opening finishes without debug hooks", 45)
                # 原版熱區0x17完整涵蓋底列第一格，含左上角(0,368)。
                click(0, 368, exact=True)
                wait(lambda s: s.get("battle", {}).get("selected") == 4, "bottom squad left top boundary selects left wing")
                shot("battle-boundary-selected")
                click(79, 399, exact=True)
                wait(lambda s: s.get("battle", {}).get("selected") == 0, "bottom squad right bottom boundary clears left wing")
                click(400, 368, exact=True)
                wait(lambda s: s.get("battle", {}).get("selected") == 8, "last bottom slot boundary selects right wing")
                click(479, 399, exact=True)
                wait(lambda s: s.get("battle", {}).get("selected") == 0, "last pixel of bottom row clears right wing")
                click(480, 399, exact=True)
                assert state()["battle"]["selected"] == 0
                checks.append("sidebar pillar does not select a squad")
                click(496, 344, exact=True)
                wait(lambda s: s.get("battle", {}).get("standing") == 4, "guard command left top boundary issues order")
                click(543, 311, exact=True)
                wait(lambda s: s.get("battle", {}).get("standing") == 1, "attack command right bottom boundary issues order")
                shot("battle-boundary-command")
                click(544, 311, exact=True)
                assert state()["battle"]["standing"] == 1
                checks.append("gap beside command does not change order")
                click(552, 312, exact=True)
                wait(lambda s: s.get("battle", {}).get("line") == 28 and s["battle"]["standing"] == 1, "adjacent formation line remains independent")
                click(496, 360, exact=True)
                wait(lambda s: s.get("battle", {}).get("leaderOrder", [])[1:] == [5], "retreat command boundary issues irreversible order")
                wait(lambda s: not s.get("battle"), "normal retreat resolves back to strategy", 90)
                for _ in range(12):
                    if not state()["messages"]:
                        break
                    toolbar(5)
                wait(lambda s: s["hostReady"], "post battle returns to normal map")
                toolbar(2)
                wait(lambda s: s["save"], "save panel after normal battle")
                after = state()
                click(300, 152)
                wait(lambda s: not s["save"], "save after normal battle")
                toolbar(3)
                click(300, 152)
                wait(lambda s: s["hostReady"] and s["capital"] == after["capital"] and s["cityOwners"] == after["cityOwners"], "reload preserves post battle ownership and capital")
                shot("battle-post-load")
                result = {"passed": True, "checks": checks, "mode": mode, "args": args,
                          "binary_sha256": hashlib.sha256(Path(binary).read_bytes()).hexdigest(),
                          "battle_entry": entry, "final_state": state()}
                (output/"mouse-toolbar-gui.json").write_text(json.dumps(result, ensure_ascii=False, indent=2)+"\n")
                return
            if field_mode:
                click(240, 47)
                wait(lambda s: s["choiceActive"], "field first march opens corps selector")
                band_move(32, 1, 100, "field first march selects march command")
                band_click()
            else:
                # 工具列按下後拖到遊戲內放開，不得開啟背景據點或命令。
                move(230, -20)
                command("mousedown", 1)
                time.sleep(.1)
                move(320, 224)
                command("mouseup", 1)
                time.sleep(.2)
                assert state()["hostReady"]
                checks.append("toolbar drag release isolated")
                click(240, 47)
                wait(lambda s: s["choiceActive"] and s["choiceRow"] == 0 and s["choiceY"] == 100, "corps selector starts at first row and y100")
                band_move(31, 0, 131, "y131 keeps first row")
                shot("original-band-first")
                band_move(1, 1, 100, "y132 moves down once and resets y100")
                shot("original-band-second")
                time.sleep(.4)
                assert state()["choiceRow"] == 1 and state()["choiceY"] == 100
                checks.append("stationary pointer does not repeat selection")
                band_move(200, 1, 100, "last row clamps without wrapping")
                band_move(-1, 0, 131, "y99 moves up once and resets y131")
                band_move(-200, 0, 131, "first row clamps without wrapping")
                band_move(1, 1, 100, "down selects march command")
                band_click(3)
                wait(lambda s: not s["choiceActive"] and s["hostReady"] and abs(s["x"]-240) <= 1 and abs(s["y"]-47) <= 1, "right click cancels and restores pointer")
                click(240, 47)
                wait(lambda s: s["choiceActive"] and s["choiceRow"] == 0, "reopened selector resets first row")
                band_move(32, 1, 100, "reopened selector moves down")
                band_click()
            wait(lambda s: s["list"], "march list by pointer")
            click(150, 112)
            click(150, 112)
            wait(lambda s: s["mapPick"], "march target picker by pointer")
            click(416, 15)
            if attacking:
                click(440+attack_xy[0]//2, 40+attack_xy[1]//2)
            else:
                click(528, 102)
            s = state()
            target_x, target_y = attack_xy if attacking else (181, 127)
            tx, ty = (target_x-s["cam"][0])*16+8, (target_y-s["cam"][1])*16+40
            assert 0 <= tx < 432 and 64 <= ty < 400, ("target outside map", tx, ty, s)
            click(tx, ty)
            wait(lambda s: s["marchMode"] and s["choiceActive"] and s["choiceRow"] == 0, "march mode uses original band selector")
            band_click()
            for _ in range(4):
                s = state()
                if not s["canBack"]:
                    break
                toolbar(4)
            wait(lambda s: s["hostReady"], "march command returns to map")
            if water:
                current_army = next(c for c in state()["corps"] if c[0] == formed[0])
                click(440+current_army[1]//2, 40+current_army[2]//2)
            else:
                click(544, 98)
            move(500, -20)
            start = next(c for c in state()["corps"] if c[0] == formed[0])
            shot("mouse-toolbar-march-start")
            moved = wait(lambda s: any(c[0] == start[0] and c[1:3] != start[1:3] for c in s.get("corps", [])), "army moves while cursor rests on toolbar", 25)
            army = next(c for c in moved["corps"] if c[0] == formed[0])
            assert moved["cam"] == state()["cam"]
            assert 0 <= army[1]-moved["cam"][0] < 27 and 0 <= army[2]-moved["cam"][1] < 23, "moving army is outside visible map"
            time.sleep(.2)
            shot("mouse-toolbar-marching")
            march_checks = {}
            if mode == "march" or attacking:
                # 行軍可能剛觸發訊息；先以正常滑鼠推進，再開存檔。
                # 沿用原本8秒等待預算，不假設截圖後仍是空閒地圖。
                end = time.monotonic()+8
                while not state()["save"]:
                    current = state()
                    if current["messages"]:
                        toolbar(5)
                    elif current["hostReady"]:
                        toolbar(2)
                    else:
                        time.sleep(.02)
                    assert time.monotonic() < end, ("save panel during normal march", state())
                wait(lambda s: s["save"], "save panel during normal march")
                saved = next(c for c in state()["corps"] if c[0] == formed[0])
                saved_route = next(c for c in state()["corpsRoutes"] if c[0] == formed[0])
                assert saved_route[2] != 0, "saved army is not on a road"
                click(300, 152)
                wait(lambda s: not s["save"], "save marching army")
                toolbar(3)
                click(300, 152)
                loaded = wait(lambda s: s["hostReady"] and any(c[0] == saved[0] and c[1:3] == saved[1:3] and c[5] == saved[5] for c in s["corps"]) and saved_route in s["corpsRoutes"], "reload preserves marching position target and road metadata")
                march_checks["saved"] = saved
                march_checks["saved_route"] = saved_route
                march_checks["loaded_routes"] = loaded["corpsRoutes"]
            if attacking:
                toolbar(1)
                wait(lambda s: s["system"], "attack march opens speed settings")
                for _ in range(5):
                    if state()["speed"] == (4 if field_mode else 0):
                        break
                    click(376, 232)
                wait(lambda s: s["speed"] == (4 if field_mode else 0), "field approaches capture at lowest speed" if field_mode else "attacker continues at highest strategy speed")
                toolbar(4)
                move(500, -20)
                end = time.monotonic()+120
                while True:
                    s = state()
                    if mode in {"capture", "defend", "field", "water"}:
                        c = next(c for c in s["corps"] if c[0] == formed[0])
                        route = next(c for c in s["corpsRoutes"] if c[0] == formed[0])
                        if s["cityOwners"][attack_node] == s["player"] and c[1:4] == [*attack_xy, 4] and route[1:3] == [attack_node, 0]:
                            break
                        assert not s.get("battle"), ("unexpected tactical opponent in garrison fixture", s)
                    else:
                        if s.get("battle"):
                            break
                        assert s["cityOwners"][attack_node] != s["player"], ("garrison-only capture, no tactical opponent", s)
                    assert time.monotonic() < end, ("normal attacker did not enter battle", s)
                    if s["messages"]:
                        toolbar(5)
                    elif s["choiceActive"]:
                        band_click()
                    else:
                        time.sleep(.03)
                entry = state()
                if mode in {"capture", "defend", "field", "water"}:
                    checks.append("normal player captures garrison and arrives at city")
                    shot("attack-garrison-arrival")
                    if field_mode:
                        # Normal road order: Runan100 to Xinxi105, or Niuzhu122 to Jianye116.
                        # Wait for the enemy to commit to the road before approaching it.
                        move(500, -20)
                        end = time.monotonic()+120
                        while True:
                            current = state()
                            route = next((r for r in current["corpsRoutes"] if r[0] == enemy_army), None)
                            enemy = next((c for c in current["corps"] if c[0] == enemy_army), None)
                            if route and enemy and route[2] == field_link and enemy[5] == attack_node and (not water or enemy[1] <= 313):
                                break
                            assert not current.get("battle"), ("enemy arrived before road order", current)
                            assert tuple(current["date"][:3]) < (196, 4, 20)
                            assert time.monotonic() < end, ("enemy did not enter counterattack road", current)
                            if current["messages"]:
                                toolbar(5)
                            else:
                                time.sleep(.03)
                        checks.append("enemy commits to actual counterattack road")
                        toolbar(1)
                        for _ in range(5):
                            if state()["speed"] == 2:
                                break
                            click(376, 232)
                        toolbar(4)
                        click(240, 47)
                        wait(lambda s: s["choiceActive"], "field march opens corps selector")
                        band_move(32, 1, 100, "field march selects march command")
                        band_click()
                        wait(lambda s: s["list"], "field march opens corps list")
                        click(150, 112)
                        click(150, 112)
                        wait(lambda s: s["mapPick"], "field march chooses formed army")
                        target_x, target_y = field_target_xy
                        click(416, 15)
                        click(440+target_x//2, 40+target_y//2)
                        s = state()
                        tx, ty = (target_x-s["cam"][0])*16+8, (target_y-s["cam"][1])*16+40
                        assert 0 <= tx < 432 and 64 <= ty < 400
                        click(tx, ty)
                        wait(lambda s: s["marchMode"] and s["choiceActive"], "field march targets the enemy road endpoint")
                        band_click()
                        for _ in range(4):
                            if not state()["canBack"]:
                                break
                            toolbar(4)
                        wait(lambda s: s["hostReady"], "field march returns to strategy")
                        toolbar(1)
                        for _ in range(5):
                            if state()["speed"] == 0:
                                break
                            click(376, 232)
                        toolbar(4)
                        move(500, -20)
                        end = time.monotonic()+120
                        while not state().get("battle"):
                            current = state()
                            assert tuple(current["date"][:3]) < (196, 4, 20), ("no natural field collision before April20", current)
                            assert time.monotonic() < end, ("waiting for normal field collision", current)
                            if current["messages"]:
                                toolbar(5)
                            elif current["choiceActive"]:
                                band_click()
                            else:
                                time.sleep(.03)
                        entry = state()
                        b = entry["battle"]
                        assert b["mode"] == 1 and {b["attacker"], b["defender"]} == {enemy_army, formed[0]}, ("unexpected field battle", entry)
                        assert [b["attacker"], b["defender"]][b["side"]] == formed[0]
                        checks.append("normal road encounter enters player field battle")
                        shot("field-normal-entry")
                        wait(lambda s: s.get("battle") and not s["battle"]["opening"] and not s["battle"]["talk"] and not s["battle"]["done"], "field opening completes normally", 45)
                        click(496, 344, exact=True)
                        wait(lambda s: s.get("battle", {}).get("standing") == 4, "field guard command by mouse")
                        shot("field-guard")
                        click(496, 360, exact=True)
                        wait(lambda s: s.get("battle", {}).get("leaderOrder", [None, None])[b["side"]] == 5, "field retreat command by mouse")
                        wait(lambda s: not s.get("battle"), "field retreat resolves normally", 90)
                    if mode == "defend":
                        end = time.monotonic()+120
                        while not state().get("battle"):
                            current = state()
                            assert tuple(current["date"][:3]) < (196, 4, 20), ("no natural counterattack before April20", current)
                            assert time.monotonic() < end, ("waiting for normal counterattack", current)
                            if current["messages"]:
                                toolbar(5)
                            elif current["choiceActive"]:
                                band_click()
                            else:
                                time.sleep(.03)
                        entry = state()
                        assert entry["battle"]["attacker"] == 3 and entry["battle"]["defender"] == formed[0] and entry["battle"]["side"] == 1, ("unexpected counterattack", entry)
                        checks.append("natural counterattack enters battle with player defending")
                        shot("defender-normal-entry")
                        wait(lambda s: s.get("battle") and not s["battle"]["opening"] and not s["battle"]["talk"] and not s["battle"]["done"], "defender opening completes normally", 45)
                        click(496, 344, exact=True)
                        wait(lambda s: s.get("battle", {}).get("standing") == 4, "defender guard command by mouse")
                        shot("defender-guard")
                        click(496, 360, exact=True)
                        wait(lambda s: s.get("battle", {}).get("leaderOrder", [])[1:] == [5], "defender retreat command by mouse")
                        wait(lambda s: not s.get("battle"), "defender retreat resolves normally", 90)
                else:
                    assert entry["battle"]["attacker"] == formed[0] and entry["battle"]["side"] == 0, ("unexpected battle", entry)
                    checks.append("normal player march enters battle as attacker")
                    shot("attack-normal-entry")
                    wait(lambda s: s.get("battle") and not s["battle"]["opening"] and not s["battle"]["talk"] and not s["battle"]["done"], "attacker opening completes normally", 45)
                    click(496, 280, exact=True)
                    wait(lambda s: s.get("battle", {}).get("standing") == 2, "attacker issues charge by mouse")
                    shot("attack-charge")
                    wait(lambda s: not s.get("battle"), "normal attacker battle resolves", 120)
                for _ in range(20):
                    if not state()["messages"]:
                        break
                    toolbar(5)
                wait(lambda s: s["hostReady"], "attacker returns to strategy")
                toolbar(2)
                wait(lambda s: s["save"], "save after attacker battle")
                after = state()
                click(300, 152)
                wait(lambda s: not s["save"], "write attacker result")
                toolbar(3)
                click(300, 152)
                wait(lambda s: s["hostReady"] and s["capital"] == after["capital"] and s["cityOwners"] == after["cityOwners"], "reload preserves attacker result")
                shot("attack-post-load")
                result = {"passed": True, "checks": checks, "mode": mode, "args": args,
                          "binary_sha256": hashlib.sha256(Path(binary).read_bytes()).hexdigest(),
                          "march_checks": march_checks, "battle_entry": entry, "final_state": state()}
                (output/"mouse-toolbar-gui.json").write_text(json.dumps(result, ensure_ascii=False, indent=2)+"\n")
                return
            if mode == "march":
                # 原版和平折返應回舞陽87；只關閉自然通知，不注入時間或行軍狀態。
                end = time.monotonic()+90
                while True:
                    s = state()
                    c = next(c for c in s["corps"] if c[0] == formed[0])
                    if c[1:4] == [191, 119, 4] and c[5] == 87:
                        break
                    assert time.monotonic() < end, ("peaceful march did not return to Wuyang", s)
                    if s["messages"]:
                        toolbar(5)
                    else:
                        time.sleep(.05)
                checks.append("peaceful border turns back and arrives at Wuyang87")
                shot("march-turnback-arrived")
                march_checks["arrived"] = s
                toolbar(2)
                click(300, 152)
                wait(lambda s: not s["save"], "save after peaceful arrival")
                toolbar(3)
                click(300, 152)
                wait(lambda s: s["hostReady"] and any(c[0] == formed[0] and c[1:4] == [191, 119, 4] and c[5] == 87 for c in s["corps"]), "reload preserves peaceful arrival")
            toolbar(6)
            wait(lambda s: s["quit"], "quit button opens confirmation")
            toolbar(4)
            wait(lambda s: not s["quit"], "back cancels quit")
            result = {"passed": True, "checks": checks, "args": args, "binary_sha256": hashlib.sha256(Path(binary).read_bytes()).hexdigest(), "march_start": start, "march_end": moved["corps"], "final_state": moved}
            if mode == "march":
                result.update(mode=mode, march_checks=march_checks, final_state=state())
            (output/"mouse-toolbar-gui.json").write_text(json.dumps(result, ensure_ascii=False, indent=2)+"\n")
            print(json.dumps(result, ensure_ascii=False), flush=True)
        except Exception:
            if "window" in locals():
                shot("mouse-toolbar-failure")
            (output/"mouse-toolbar-gui.json").write_text(json.dumps({"passed": False, "checks": checks, "state": state()}, ensure_ascii=False, indent=2)+"\n")
            raise
        finally:
            proc.terminate()
            try:
                proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                proc.kill()
                proc.wait()
            log.close()


if __name__ == "__main__":
    run(*sys.argv[1:])
