#!/usr/bin/env python3
"""產生 remake 推廣片使用的原創、可重現短配樂。

不讀取原版 BGM／SOUND.DAT，也不使用外部取樣。五音動機是 D4–C4–F4–E4–A4，
以簡單的三角波、方波與合成鼓聲表現「地圖巡行 → 戰術脈衝 → 收束」；輸出
固定 60 秒、44.1 kHz、立體聲 WAV，供 ffmpeg 混入推廣片。
"""

import argparse
import math
import struct
import wave


RATE = 44100
DURATION = 60
MOTIF = (62, 60, 65, 64, 69)


def hz(midi):
    return 440.0 * (2.0 ** ((midi - 69) / 12.0))


def envelope(t, length, attack=0.025, release=0.16):
    if t < 0 or t >= length:
        return 0.0
    if t < attack:
        return t / attack
    if t > length - release:
        return max(0.0, (length - t) / release)
    return 1.0


def triangle(phase):
    return 4.0 * abs((phase % 1.0) - 0.5) - 1.0


def add_tone(buf, start, length, midi, amp, shape="triangle"):
    frequency = hz(midi)
    first = max(0, int(start * RATE))
    last = min(len(buf), int((start + length) * RATE))
    for index in range(first, last):
        t = index / RATE - start
        phase = t * frequency
        if shape == "square":
            value = 1.0 if phase % 1.0 < 0.5 else -1.0
        else:
            value = triangle(phase)
        buf[index] += amp * envelope(t, length) * value


def add_kick(buf, start, amp=0.28):
    first = max(0, int(start * RATE))
    last = min(len(buf), int((start + 0.22) * RATE))
    for index in range(first, last):
        t = index / RATE - start
        phase = (72.0 - 38.0 * min(1.0, t / 0.22)) * t
        buf[index] += amp * math.exp(-18.0 * t) * math.sin(2.0 * math.pi * phase)


def add_hat(buf, start, amp=0.09):
    first = max(0, int(start * RATE))
    last = min(len(buf), int((start + 0.07) * RATE))
    for index in range(first, last):
        t = index / RATE - start
        noise = math.sin(index * 12.9898) * math.sin(index * 78.233)
        buf[index] += amp * math.exp(-48.0 * t) * noise


def compose():
    buf = [0.0] * (RATE * DURATION)

    # 0–15 秒：低密度地圖探索。保留足夠空白，讓畫面與 UI 有呼吸。
    for offset, note in enumerate(MOTIF):
        add_tone(buf, 1.0 + offset * 2.35, 1.15, note, 0.16)
    for start in (8.0, 10.4, 12.8):
        add_tone(buf, start, 1.0, 50, 0.07)

    # 15–40 秒：戰術畫面，動機壓縮成每兩秒一格，加入低音與鼓。
    for beat in range(12):
        start = 16.0 + beat * 2.0
        note = MOTIF[beat % len(MOTIF)]
        add_tone(buf, start, 0.72, note, 0.22, "square")
        add_tone(buf, start, 0.9, note - 24, 0.12)
        add_kick(buf, start)
        add_hat(buf, start + 0.5)
        add_hat(buf, start + 1.5)

    # 40–52 秒：高潮，五音動機與開放五度疊合；勝利不落俗套地保留未完感。
    for beat in range(12):
        start = 40.0 + beat
        note = MOTIF[beat % len(MOTIF)]
        add_tone(buf, start, 0.68, note + 12, 0.22)
        add_tone(buf, start, 0.85, note - 12, 0.11)
        add_kick(buf, start, 0.22)
        if beat % 2 == 1:
            add_hat(buf, start + 0.5, 0.11)

    # 52–60 秒：回到地圖材質，最後停在 A4，不做完整大調終止。
    for offset, note in enumerate(MOTIF):
        add_tone(buf, 52.0 + offset * 1.35, 0.9, note, 0.15)
    add_tone(buf, 58.5, 1.2, MOTIF[-1], 0.12)

    peak = max(1.0, max(abs(value) for value in buf))
    return [max(-0.92, min(0.92, value * 0.82 / peak)) for value in buf]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("out", help="輸出的 WAV 路徑")
    args = parser.parse_args()
    samples = compose()
    with wave.open(args.out, "wb") as stream:
        stream.setnchannels(2)
        stream.setsampwidth(2)
        stream.setframerate(RATE)
        frames = bytearray()
        for value in samples:
            packed = struct.pack("<h", int(value * 32767))
            frames += packed + packed
        stream.writeframes(frames)
    print(f"promo score: {args.out} ({DURATION}s, motif={'-'.join(map(str, MOTIF))})")


if __name__ == "__main__":
    main()
