"""Render the local extended trailer using the existing captures and a timed voice plan."""
import hashlib
import json
import math
import os
from pathlib import Path
import re
import subprocess

import promo_story_render as short

WORK = Path("/out/promo-extended")
SOURCES = Path("/out/promo-story")
OUTPUT = WORK / "candidate"
DURATION = 102
FPS = 30
VERSION = short.VERSION
FILM = OUTPUT / f"wolong-remake-promo-{VERSION}-long.mp4"
MUSIC = Path("/repo/workplace/dosbox/audio/d7open  _000.wav")
ORIGINAL = short.REVIEW / "battle-hd-victory-native-native-245.png"
FORMED = short.REVIEW / "command-tools-joined/mouse-toolbar-formed.png"
EDITOR = short.REVIEW / "command-tools-joined/formed-commander-edited.png"
PHONE_MAP = Path("/delivery/smoke/android-startup.png")
PHONE_BATTLE = short.REVIEW / "android-compact-install/android-battle-selected.png"


def digest(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()


def media_length(path):
    return short.length(str(path))


def visual(name, path, seconds, start=0, still=False, full=False, crop=None):
    inputs = ["-loop", 1, "-framerate", FPS] if still else ["-ss", start]
    viewport = "fps=30,setsar=1" if full else short.scaled(1152, 552, 68)
    if crop:
        viewport = f"crop={crop},{viewport}"
    short.ffmpeg(*inputs, "-i", path, "-vf", viewport, "-t", seconds, *short.VIDEO, name)


def visuals():
    edits = [
        ("opening", SOURCES/"01.mp4", 12, 0, False, True, None),
        ("dosv", ORIGINAL, 4, 0, True, False, None),
        ("map-original", short.CAPTURE/"map.mp4", 7, 0, False, False, None),
        ("formation", FORMED, 4, 0, True, False, None),
        ("march", short.CAPTURE/"map.mp4", 4, 2, False, False, None),
        ("movement", short.CAPTURE/"field.mp4", 3, 1, False, False, None),
        ("gallery", SOURCES/"04.mp4", 5, 0, False, True, None),
        ("map-hd", short.CAPTURE/"map.mp4", 7, 4, False, False, None),
        ("siege", short.CAPTURE/"siege.mp4", 11, 5, False, False, None),
        ("editor", EDITOR, 5, 0, True, False, None),
        ("editor-detail", EDITOR, 6, 0, True, False, "592:530:520:177"),
        ("phone-map", PHONE_MAP, 4, 0, True, False, None),
        ("phone-battle", PHONE_BATTLE, 7, 0, True, False, None),
        ("closing", short.CAPTURE/"siege.mp4", 7, 8, False, False, None),
    ]
    for key, path, duration, start, still, full, crop in edits:
        target = Path(key+".mp4")
        if target.exists():
            continue
        if key == "gallery":
            short.ffmpeg("-i", path, "-vf", "tpad=stop_mode=clone:stop_duration=4",
                         "-t", duration, *short.VIDEO, target)
        else:
            visual(target, path, duration, start, still, full, crop)
    camera = "[0:v]scale=592:454,setsar=1[a];[1:v]scale=592:454,setsar=1[b];"
    camera += "[a][b]hstack,pad=1280:720:48:104:color=0x121c25,fps=30,"
    camera += "tpad=stop_mode=clone:stop_duration=2[v]"
    if not Path("ai.mp4").is_file():
        short.ffmpeg("-framerate", 6, "-i", SOURCES/"frames/00-%03d.png",
                     "-framerate", 6, "-i", SOURCES/"frames/15-%03d.png",
                     "-filter_complex", camera, "-map", "[v]", "-t", 16, *short.VIDEO, "ai.mp4")
    keys = [edit[0] for edit in edits]
    keys.insert(keys.index("editor"), "ai")
    expected = [edit[2] for edit in edits]
    expected.insert(keys.index("ai"), 16)
    for key, duration in zip(keys, expected, strict=True):
        assert math.isclose(media_length(key+".mp4"), duration, abs_tol=.005), (key, duration)
    Path("cuts.txt").write_text("".join(f"file '{key}.mp4'\n" for key in keys))
    return sorted({edit[1] for edit in edits})


def timed_narration(plan):
    clips = []
    for chapter, scene in enumerate(plan):
        offsets = []
        for index, line in enumerate(scene["lines"]):
            path = WORK/"voice"/f"{chapter:02d}-{index:02d}.wav"
            offsets.append((line["text"], path, media_length(path)))
        start = scene["start"]+scene["lead"]
        room = scene["end"] - start - .32 - .24*(len(offsets)-1)
        speed = max(1., sum(item[2] for item in offsets)/room)
        assert speed <= 1.12, (scene["title"], speed)
        for text, path, length in offsets:
            end = start + length/speed
            clips.append({"text":text,"wav":str(path),"start":round(start,3),"end":round(end,3),
                          "speed":round(speed,5),"sha256":digest(path)})
            start = end+.24
    assert all(b["start"]>=a["end"] for a,b in zip(clips,clips[1:]))
    return clips


def wrap(text, separator):
    if len(text)<=31:
        return text
    near = [index+1 for index,char in enumerate(text) if char in "，、。" and 14<=index<=30]
    cut = near[-1] if near else 30
    return text[:cut]+separator+text[cut:]


def captions(plan, clips):
    events = []
    srt = []
    def caption(a,b,style,text,x=None,y=None):
        pos = rf"{{\pos({x},{y})}}" if x is not None else ""
        events.append(f"Dialogue: 0,{short.timestamp(a)},{short.timestamp(b)},{style},,0,0,0,,{pos}{text}")
    def timecode(t):
        ts = round(t*1000)
        return f"{ts//3600000:02}:{ts//60000%60:02}:{ts//1000%60:02},{ts%1000:03}"
    for index, clip in enumerate(clips, 1):
        text = clip["text"]
        caption(clip["start"],clip["end"],"Text",wrap(text,r"\N"))
        srt.extend((str(index), f"{timecode(clip['start'])} --> {timecode(clip['end'])}", wrap(text,"\n"),""))
    for scene in plan[1:]:
        if scene["start"] != 57:
            caption(scene["start"],scene["end"],"Title",scene["title"],640,16)
    caption(12,16,"Status","松崗 DOS/V 原版擷取",640,608)
    caption(57,73,"Title","椎陽攻城・同一梯位・受控重播",640,13)
    caption(57,73,"Status","戰術 AI 關閉",344,68)
    caption(57,73,"Status","戰術 AI 開啟・統帥 15",936,68)
    replay = {0:{},15:{}}
    for line in (SOURCES/"replay.jsonl").read_text().splitlines():
        item = json.loads(line)
        if item["y"] == 26:
            replay[item["rating"]][item["frame"]] = item
    for rating,x,finish in ((0,344,81),(15,936,41)):
        track = replay[rating]
        assert not track[490+finish]["broken"] and track[491+finish]["broken"]
        assert track[492+finish]["builder_z"] > 0
        for tick in range(0,finish,3):
            progress=(80-track[491+tick]["durability"])*100//80
            caption(57+tick/6,57+min(tick+3,finish)/6,"Status",
                    f"架梯 {progress}%　｜　第 {tick} 拍",x,573)
        caption(57+finish/6,73,"Status",rf"{{\c&H9CE9F5&}}{finish} 拍完成",x,573)
    caption(84,95,"Status","Android 模擬器畫面",640,612)
    caption(95,102,"Status","Linux・Windows・macOS・Android",640,610)
    Path("titles.ass").write_text(short.ASS_HEADER+"\n".join(events)+"\n")
    (OUTPUT/f"wolong-promo-{VERSION}-long.zh-Hant.srt").write_text("\n".join(srt)+"\n")


def mix(clips):
    filters = []
    inputs = []
    for index, clip in enumerate(clips):
        inputs += ["-i", clip["wav"]]
        delay = round(clip["start"]*1000)
        filters.append(f"[{index}:a]aresample=48000,atempo={clip['speed']},"
                       "loudnorm=I=-18:TP=-2:LRA=7,aresample=48000,aformat=channel_layouts=stereo,"
                       f"adelay={delay}:all=1,apad,atrim=duration={DURATION}[v{index}]")
    marks = "".join(f"[v{i}]" for i in range(len(clips)))
    filters.append(f"{marks}amix=inputs={len(clips)}:normalize=0:duration=longest[voice]")
    short.ffmpeg(*inputs, "-filter_complex", ";".join(filters), "-map", "[voice]", "-c:a", "pcm_s16le", "narration.wav")
    music_inputs = ["-i", MUSIC] * 8
    score = [f"[{i}:a]atrim=start=0.65:end=14.65,asetpts=PTS-STARTPTS[m{i}]" for i in range(8)]
    score += ["[m0][m1]acrossfade=d=0.55[b1]"]
    score += [f"[b{i-1}][m{i}]acrossfade=d=0.55[b{i}]" for i in range(2,8)]
    score += ["[b7]asetpts=N/SR/TB,atrim=duration=102,aresample=48000,volume=17dB,"
              "alimiter=limit=0.72:level=false,afade=t=in:d=0.45,afade=t=out:st=99:d=3[music]"]
    short.ffmpeg(*music_inputs, "-filter_complex", ";".join(score), "-map", "[music]", "-c:a", "pcm_s16le", "score.wav")
    assert math.isclose(media_length("score.wav"), DURATION, abs_tol=.01)
    mix_filter = ("[2:a]asplit=2[side][voice];"
                  "[1:a][side]sidechaincompress=threshold=0.025:ratio=8:attack=30:release=450[bed];"
                  "[bed][voice]amix=inputs=2:duration=first:normalize=0,alimiter=limit=0.82:level=false:latency=1[audio]")
    short.ffmpeg("-f","concat","-safe",0,"-i","cuts.txt","-i","score.wav","-i","narration.wav",
                 "-filter_complex",mix_filter,"-vf","ass=titles.ass",
                 "-map","0:v","-map","[audio]","-t",DURATION,"-c:v","libx264","-threads",2,"-preset","fast",
                 "-crf",18,"-pix_fmt","yuv420p","-c:a","aac","-b:a","192k","-movflags","+faststart",FILM)


def main():
    os.chdir(WORK)
    plan=json.loads(Path("plan.json").read_text())
    clips=timed_narration(plan)
    sources=visuals()
    captions(plan,clips)
    mix(clips)
    probe=json.loads(subprocess.check_output(["ffprobe","-v","error","-show_format","-show_streams","-of","json",FILM],text=True))
    video,audio=probe["streams"]
    assert (video["width"],video["height"],video["nb_frames"],video["r_frame_rate"]) == (1280,720,str(DURATION*FPS),"30/1")
    assert abs(float(video["duration"])-float(audio["duration"]))<.04 and audio["channels"]==2
    cmd=["ffmpeg","-hide_banner","-nostdin","-nostats","-xerror","-i",str(FILM),
         "-vf","blackdetect=d=0.3:pix_th=0.1,freezedetect=n=-60dB:d=1.5",
         "-af","silencedetect=n=-45dB:d=0.5,volumedetect","-f","null","-"]
    log=subprocess.run(cmd,capture_output=True,text=True,check=True).stderr
    peak=float(re.search(r"max_volume: ([-\d.]+) dB",log)[1])
    mean=float(re.search(r"mean_volume: ([-\d.]+) dB",log)[1])
    quiet=[float(d) for d in re.findall(r"silence_duration: ([\d.]+)",log)]
    assert -32<mean<-9 and -10<peak<-.5 and "black_start:" not in log
    assert len(quiet)==log.count("silence_start:") and all(d<1 for d in quiet)
    (OUTPUT/"long-media-check.log").write_text(log)
    short.ffmpeg("-i",FILM,"-vf","fps=1/6,scale=320:180,tile=6x3","-frames:v",1,OUTPUT/"long-contact-sheet.png")
    checks=[Path("/repo/tools/promo_extended_render.py"),Path("/repo/tools/promo_story_render.py"),
            WORK/"plan.json",MUSIC,SOURCES/"replay.go",SOURCES/"replay.jsonl",
            SOURCES/"piper/model-provenance.json",*sources,*[Path(clip["wav"]) for clip in clips]]
    report={"duration_seconds":DURATION,"film_sha256":digest(FILM),"narration":"Piper Huayan medium, offline",
            "reference":"rich2 per-scene captions, voice normalization and sidechain ducking",
            "mean_db":mean,"peak_db":peak,"longest_quiet_seconds":max(quiet,default=0),
            "soundtrack":"原版 DOS/V 片頭曲擷取","ladder_ticks_per_second":6,
            "source_sha256":{str(p):digest(p) for p in checks},"narration_cues":clips,
            "scenes":plan,"private":True,"ffprobe":probe}
    (OUTPUT/"long-story.json").write_text(json.dumps(report,ensure_ascii=False,indent=2)+"\n")
    print(DURATION,"seconds",len(clips),"voice clips",peak,"dB",flush=True)


if __name__=="__main__":
    main()
