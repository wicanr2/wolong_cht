"""Compose the approved local 47-second Wolong remake trailer."""
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess

STAGE = Path("/out/promo-story")
CANDIDATE = STAGE / "candidate"
OLD = Path("/repo/dist-all/v.1.1.0-20261122/promo/wolong-remake-promo-v.1.1.0-20261122.mp4")
CAPTURE = Path("/repo/workplace/desktop-package-v.1.1.0-20261122/capture")
REVIEW = Path("/review")
IMAGE_ROOT = Path("/repo/docs/images")
VERSION = os.environ["WOLONG_RELEASE_VERSION"]
assert re.fullmatch(r"v\.[0-9]+\.[0-9]+\.[0-9]+-[0-9]{8}", VERSION)

BG = "0x121c25"
FPS = 30
VIDEO = ["-an", "-c:v", "libx264", "-threads", "2", "-preset", "fast",
         "-crf", "18", "-pix_fmt", "yuv420p"]


def ffmpeg(*args):
    subprocess.run(["ffmpeg", "-y", "-nostdin", "-hide_banner", "-loglevel", "error",
                    "-filter_threads", "2", "-filter_complex_threads", "2", *map(str, args)], check=True)


def scaled(w=1152, h=608, top=56):
    return (f"fps={FPS},scale={w}:{h}:force_original_aspect_ratio=decrease:"
            f"force_divisible_by=2:flags=lanczos,pad=1280:720:(ow-iw)/2:{top}+({h}-ih)/2:"
            f"color={BG},setsar=1")


def clip(name, path, duration, start=0, still=False, fullscreen=False):
    args = ["-loop", "1", "-framerate", str(FPS)] if still else ["-ss", str(start)]
    args += ["-i", path, "-t", str(duration)]
    args += ["-vf", "fps=30,setsar=1" if fullscreen else scaled()]
    ffmpeg(*args, *VIDEO, f"{name}.mp4")


def split(name, sources, panel, duration):
    w, h, top = panel
    pad_filter = (f"scale={w}:{h}:force_original_aspect_ratio=decrease:force_divisible_by=2:"
                  f"flags=lanczos,pad={w}:{h}:(ow-iw)/2:(oh-ih)/2:color={BG},setsar=1")
    a = f"[0:v]{pad_filter}[left];[1:v]{pad_filter}[right]"
    a += (f";[left][right]hstack,pad=1280:720:(ow-iw)/2:{top}:color={BG},"
          "fps=30,tpad=stop_mode=clone:stop_duration=0.25[v]")
    ffmpeg(*sources[0], *sources[1], "-filter_complex", a,
           "-map", "[v]", "-t", duration, *VIDEO, f"{name}.mp4")


def timestamp(t):
    dec = round(t * 100)
    return f"{dec//360000}:{dec//6000%60:02}:{dec//100%60:02}.{dec%100:02}"


ASS_HEADER = """[Script Info]
ScriptType: v4.00+
PlayResX: 1280
PlayResY: 720
WrapStyle: 2
ScaledBorderAndShadow: yes
[V4+ Styles]
Format: Name,Fontname,Fontsize,PrimaryColour,SecondaryColour,OutlineColour,BackColour,Bold,Italic,Underline,StrikeOut,ScaleX,ScaleY,Spacing,Angle,BorderStyle,Outline,Shadow,Alignment,MarginL,MarginR,MarginV,Encoding
Style: Title,Noto Serif CJK TC,34,&H00ADE7F8,&H00FFFFFF,&H00101721,&H00101721,-1,0,0,0,100,100,1,0,1,2,0,8,20,20,15,1
Style: Text,Noto Sans CJK TC,27,&H00EEF3FA,&H00FFFFFF,&H20121720,&H40121720,0,0,0,0,100,100,0,0,1,2,1,2,20,20,18,1
Style: Status,Noto Sans CJK TC,25,&H00FFFFFF,&H00FFFFFF,&H00101721,&H00101721,-1,0,0,0,100,100,0,0,1,2,0,8,20,20,20,1
[Events]
Format: Layer,Start,End,Style,Name,MarginL,MarginR,MarginV,Effect,Text
"""


def make_titles(rows, replay):
    events = []

    def cue(a, b, style, text, x=None, y=None):
        pos = rf"{{\pos({x},{y})}}" if x is not None else ""
        events.append(f"Dialogue: 0,{timestamp(a)},{timestamp(b)},{style},,0,0,0,,{pos}{text}")

    cue(3.15, 6.2, "Text", rows[0]["caption"])
    cue(6.2, 11.7, "Text", "你提議出兵，君主可能不答應。")
    cue(12, 15, "Title", "一款即時三國戰略，值得留下來", 640, 12)
    cue(12, 15, "Text", "松崗 DOS/V 原版畫面")
    cue(15, 22, "Title", "用 Go 重寫《臥龍傳》", 640, 12)
    cue(15, 19, "Text", rows[1]["caption"])
    cue(19, 22, "Text", "保留即時戰略與軍師進言")
    cue(22, 24, "Title", "沿原版造型重繪", 640, 12)
    cue(22, 24, "Text", "武將頭像・城池景觀")
    cue(24, 30, "Title", "HD 戰場", 640, 12)
    cue(24, 30, "Text", rows[2]["caption"])
    cue(30, 41, "Title", "椎陽攻城・同一梯位・受控測試", 640, 14)
    cue(30, 41, "Status", "戰術 AI 關閉", 320, 82)
    cue(30, 41, "Status", "戰術 AI 開啟・統帥 15", 960, 82)
    cue(30, 41, "Text", rows[3]["caption"])
    for rating, x, completed in ((0, 320, 81), (15, 960, 41)):
        for tick in range(0, completed, 4):
            sample = replay[rating][491 + tick]
            progress = (80 - sample["durability"]) * 100 // 80
            cue(30 + tick / 8, 30 + min(tick+4, completed) / 8, "Status",
                f"架梯 {progress:>2}%　｜　第 {tick:>2} 拍", x, 604)
        cue(30 + completed/8, 41, "Status",
            rf"{{\c&H9CE9F5&}}{completed} 拍完成", x, 600)
    cue(41, 44, "Title", "編輯武將能力", 640, 11)
    cue(41, 44, "Text", rows[4]["caption"])
    cue(44, 47, "Title", "臥龍傳 Remake", 640, 11)
    cue(44, 47, "Text", "Linux・Windows・macOS・Android")
    Path("titles.ass").write_text(ASS_HEADER+"\n".join(events)+"\n")


def length(path):
    result = subprocess.check_output(["ffprobe", "-v", "error", "-of", "json",
                                      "-show_format", path], text=True)
    return float(json.loads(result)["format"]["duration"])


def audio_mix(rows, voices):
    filters = []
    timing = []
    for index, (row, voice) in enumerate(zip(rows, voices), 2):
        d = length(voice)
        available = row["start"] + row["duration"] - row["at"] - .15
        speed = max(1., d / available)
        assert speed <= 2.0, (voice, d, available)
        delay = round(row["at"] * 1000)
        filters.append(f"[{index}:a]aresample=48000,atempo={speed:.5f},"
                       f"afade=t=in:d=0.02,volume=1.6,adelay={delay}:all=1[a{index}]")
        timing.append({"input_seconds": d, "start": row["at"], "speed": speed})
    windows = [f"between(t,{row['at']-.14},{row['at'] + t['input_seconds']/t['speed'] + .16})"
               for row, t in zip(rows, timing)]
    music = f"[1:a]aresample=48000,volume='if({'+'.join(windows)},0.25,0.8)':eval=frame[score]"
    mix = "[score]"+"".join(f"[a{i}]" for i in range(2, len(voices)+2))
    mix += (f"amix=inputs={len(voices)+1}:normalize=0:duration=first,"
            "alimiter=limit=.80:level=false:latency=1,atrim=duration=47,apad=whole_dur=47[a]")
    return ";".join([music, *filters, mix]), timing


def main():
    os.chdir(STAGE)
    rows = json.loads(Path("story.json").read_text())
    replay = {0: {}, 15: {}}
    for line in Path("replay.jsonl").read_text().splitlines():
        row = json.loads(line)
        if row["y"] == 26:
            replay[row["rating"]][row["frame"]] = row
    for rating, completed in ((0, 81), (15, 41)):
        run = replay[rating]
        assert run[491]["durability"] == 80
        assert run[492]["durability"] == (79 if rating == 0 else 78)
        assert not run[490+completed]["broken"]
        assert run[491+completed]["broken"] and run[492+completed]["builder_z"] > 0
    make_titles(rows, replay)

    original = REVIEW / "battle-hd-victory-native-native-245.png"
    portraits = IMAGE_ROOT / "readme-hd-person-50.png"
    city = next(iter(sorted(IMAGE_ROOT.glob("readme-hd-city-*.png"))))
    editor = REVIEW / "command-tools-joined/formed-commander-edited.png"
    android = REVIEW / "android-compact-install/android-battle-selected.png"
    if not Path("08.mp4").exists():
        clip("01", OLD, 12, fullscreen=True)
        clip("02", original, 3, still=True)
        clip("03", CAPTURE / "map.mp4", 7, start=1)
        split("04", [["-loop", "1", "-i", portraits], ["-loop", "1", "-i", city]], (610, 510, 85), 2)
        clip("05", CAPTURE / "siege.mp4", 6, start=9)
        clip("07", editor, 3, still=True)
        clip("08", android, 3, still=True)
        Path("cuts.txt").write_text("".join(f"file '{i:02d}.mp4'\n" for i in range(1, 9)))
    if not Path("06.mp4").exists():
        split("06", [["-framerate", "8", "-i", "frames/00-%03d.png"],
                     ["-framerate", "8", "-i", "frames/15-%03d.png"]], (616, 472, 119), 11)

    use_piper = os.environ.get("WOLONG_PROMO_VOICE") == "piper"
    voice_dir = STAGE / ("voice-piper" if use_piper else "voice")
    narrator = "Piper 1.4.1 Huayan medium，本機合成" if use_piper else "Qwen3-TTS 1.7B CustomVoice Uncle_Fu，本機合成"
    voices = [voice_dir / f"{i}.wav" for i in range(len(rows))]
    film = CANDIDATE / f"wolong-remake-promo-{VERSION}.mp4"
    if not all(voice.is_file() for voice in voices):
        print("Visual segments ready. Rendering a local subtitle preview while narration finishes.", flush=True)
        ffmpeg("-f", "concat", "-safe", "0", "-i", "cuts.txt", "-i", OLD, "-vf", "ass=titles.ass",
               "-map", "0:v", "-map", "1:a", "-t", "47", "-c:v", "libx264", "-threads", "2",
               "-preset", "fast", "-crf", "19", "-c:a", "aac", "-movflags", "+faststart", "preview.mp4")
        return
    audio, timing = audio_mix(rows, voices)
    args = ["-f", "concat", "-safe", "0", "-i", "cuts.txt", "-i", OLD]
    for voice in voices:
        args += ["-i", voice]
    ffmpeg(*args, "-filter_complex", audio, "-vf", "ass=titles.ass", "-map", "0:v", "-map", "[a]",
           "-t", "47", "-c:v", "libx264", "-threads", "2", "-preset", "fast", "-crf", "18",
           "-c:a", "aac", "-b:a", "192k", "-movflags", "+faststart", film)
    with (CANDIDATE / "ffprobe.json").open("w") as out:
        subprocess.run(["ffprobe", "-v", "error", "-show_format", "-show_streams", "-of", "json", film],
                       stdout=out, check=True)
    probe = json.loads((CANDIDATE / "ffprobe.json").read_text())
    video, audio_stream = probe["streams"]
    assert video["nb_frames"] == "1410" and video["r_frame_rate"] == "30/1"
    assert abs(float(video["duration"]) - float(audio_stream["duration"])) < .04
    with (CANDIDATE / "media-check.log").open("w") as log:
        subprocess.run(["ffmpeg", "-nostdin", "-hide_banner", "-nostats", "-xerror", "-i", film,
                        "-vf", "blackdetect=d=0.3:pix_th=0.1,freezedetect=n=-60dB:d=1.5",
                        "-af", "silencedetect=n=-45dB:d=0.5,volumedetect", "-f", "null", "-"],
                       stderr=log, check=True)
    ffmpeg("-i", film, "-vf", "fps=1/4,scale=320:180,tile=4x3",
           "-frames:v", "1", CANDIDATE / "contact-sheet.png")
    ffmpeg("-i", film, "-t", "12", "-vf", "fps=1/3:start_time=0,scale=480:270,tile=2x2",
           "-frames:v", "1", CANDIDATE / "opening-contact-sheet.png")
    for timestamp_seconds in (4, 13, 18, 23, 28, 33, 36, 40.5, 42, 45):
        ffmpeg("-ss", timestamp_seconds, "-i", film, "-frames:v", "1",
               CANDIDATE / f"sample-{timestamp_seconds}.png")
    sources = [OLD, original, portraits, city, CAPTURE/"map.mp4", CAPTURE/"siege.mp4",
               editor, android, *voices, Path("story.json").resolve(), Path("replay.jsonl").resolve(),
               Path("replay.go").resolve(), Path("/repo/tools/promo_story_render.py"),
               REVIEW/"ladder-timing-20261122/ladder-comparison.json",
               Path("/repo/workplace/dosbox/audio/d7open  _000.wav")]
    if use_piper:
        sources += [STAGE / "piper" / name for name in
                    ("MODEL_CARD", "model-provenance.json", "zh_CN-huayan-medium.onnx.json")]
        sources.append(voice_dir / "input.txt")
    with (CANDIDATE / "source-SHA256SUMS.txt").open("w") as output:
        for path in sources:
            with path.open("rb") as stream:
                output.write(f"{hashlib.file_digest(stream, 'sha256').hexdigest()}  {path}\n")
    (CANDIDATE / "story.json").write_text(json.dumps({"cuts":rows,"voice_timing":timing,
        "narration":narrator,
        "ladder":{"frame_origin":491,"ticks_per_second":8,"target":[39,26],"off":81,"on_15":41},
        "scope":"本機保存"}, ensure_ascii=False, indent=2)+"\n")
    with film.open("rb") as stream:
        digest = hashlib.file_digest(stream, "sha256").hexdigest()
    (CANDIDATE / "source.json").write_text(json.dumps({
        "capture_version": "v.1.1.0-20261122 與現行遊戲引擎重播",
        "opening_seconds": 12, "sha256": digest, "scope": "本機保存",
    }, ensure_ascii=False, indent=2)+"\n")


if __name__ == "__main__":
    main()
