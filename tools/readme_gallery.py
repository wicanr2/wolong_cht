#!/usr/bin/env python3
"""在 Docker 內重生 README 展示圖；來源固定於 readme-gallery.json。

只排版、縮放展示圖，不修改原圖或生成新美術。需要 Pillow 與自備中文字型。
原版採最近鄰縮放；HD 採 Lanczos 縮小；兩側完整顯示相同取景範圍。
"""

import argparse
import hashlib
import json
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont, ImageOps


def sha256(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()


def render(entry, root, font_path):
    pair = len(entry["sources"]) == 2
    box = (480, 480) if pair else (960, 600)
    canvas = Image.new("RGB", (1024, box[1] + 152), "#101b1b")
    draw = ImageDraw.Draw(canvas)
    title = ImageFont.truetype(str(font_path), 30)
    label = ImageFont.truetype(str(font_path), 22)
    note = ImageFont.truetype(str(font_path), 19)
    draw.text((32, 18), entry["title"], font=title, fill="#f1dfba")
    draw.line((32, 61, 992, 61), fill="#65705d")
    for index, source in enumerate(entry["sources"]):
        path = root / source["path"]
        if sha256(path) != source["sha256"]:
            raise ValueError(f"來源雜湊不符：{path}")
        with Image.open(path) as image:
            image.load()
            if list(image.size) != source["size"]:
                raise ValueError(f"來源尺寸不符：{path}")
            resample = Image.Resampling.LANCZOS if source["mode"] == "hd" else Image.Resampling.NEAREST
            preview = ImageOps.contain(image.convert("RGB"), box, method=resample)
        x = 24 + index * 496 if pair else 32
        draw.text((x, 72), source["label"], font=label, fill="#d8dfd5")
        draw.rectangle((x, 108, x + box[0] - 1, 108 + box[1] - 1), fill="#080f10")
        canvas.paste(preview, (x + (box[0] - preview.width) // 2, 108 + (box[1] - preview.height) // 2))
    draw.text((32, canvas.height - 32), entry["note"], font=note, fill="#b8c4b6")
    return canvas


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--font", type=Path, required=True)
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parents[1])
    parser.add_argument("--output", type=Path)
    parser.add_argument("--receipt", type=Path, required=True)
    args = parser.parse_args()
    root = args.root.resolve()
    output = args.output or root / "docs/images"
    # 輸出目錄由呼叫端預先建立，避免誤寫到不存在的掛載。
    if not output.is_dir() or not args.receipt.parent.is_dir() or not args.font.is_file():
        parser.error("字型、輸出目錄或收據目錄不存在")
    catalog = root / "tools/readme-gallery.json"
    manifest = json.loads(catalog.read_text())
    rows = []
    for entry in manifest["images"]:
        target = output / entry["file"]
        image = render(entry, root, args.font)
        image.save(target, optimize=True)
        rows.append({"file": str(target.relative_to(root)) if target.is_relative_to(root) else str(target),
                     "sha256": sha256(target), "size": list(image.size), "bytes": target.stat().st_size})
    receipt = {"catalog_sha256": sha256(catalog), "script_sha256": sha256(Path(__file__)),
               "font_sha256": sha256(args.font), "pillow_version": Image.__version__,
               "scope": "12組美術對照、4張戰場現況展示；不證明操作或戰場HD已完成", "images": rows}
    args.receipt.write_text(json.dumps(receipt, ensure_ascii=False, indent=2) + "\n")
    print(f"完成 {len(rows)} 張展示圖，共 {sum(row['bytes'] for row in rows):,} bytes")


if __name__ == "__main__":
    main()
