#!/usr/bin/env python3
"""保存已匹配的指令來源；從自備原版匯入非指令區，驗證完整 EXE。"""
from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

EXPECTED = "fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868"
SIZE, INSTRUCTIONS, CODE_BYTES = 67099, 24714, 56197
SOURCE = "tools/c_recovery/KI.code.S"
LINKER = "tools/c_recovery/KI.code.ld"
MANIFEST = "docs/re/assembly-code-record.json"
SUPPLEMENT = "docs/re/rectangle-handler-code.json"


def digest(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def require(ok: bool, message: str) -> None:
    if not ok:
        raise ValueError(message)


def write_json(path: Path, value: dict) -> None:
    path.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def original_bytes(repo: Path) -> bytes:
    raw = (repo / "workplace/orig/dosv/KI.EXE").read_bytes()
    require(len(raw) == SIZE and digest(raw) == EXPECTED, "original identity differs")
    return raw


def export(repo: Path) -> None:
    """只取已驗證 source-map 的 instruction 行，不複製任何資料宣告。"""
    raw = original_bytes(repo)
    local = repo / "workplace/matching-decompilation/assembly/build"
    report = json.loads((local / "report.json").read_text(encoding="utf-8"))
    source = (local / "KI.reconstructed.S").read_bytes()
    linker = (local / "KI.ld").read_bytes()
    require(report["whole_file_exact"] and report["standalone_source_exact"], "baseline not exact")
    require(report["input_sha256"] == EXPECTED, "baseline input differs")
    require(report["source_sha256"] == digest(source), "baseline source differs")
    require(report["linker_source_sha256"] == digest(linker), "baseline linker differs")
    mapping = json.loads((local / "source-map.json").read_text(encoding="utf-8"))
    semantic_path = repo / "docs/re/matching-semantic-index.json"
    semantics = json.loads(semantic_path.read_text(encoding="utf-8"))
    require(semantics["input_sha256"] == EXPECTED, "semantic index input differs")
    lines = source.decode("utf-8").splitlines()
    public = [
        "# DOS/V KI.EXE: matched instruction source, GNU binutils 2.40.",
        "# File offsets include the 512-byte MZ header. IDA linear base: 0x10000.",
        "# Gaps are zero placeholders; import private data with matching_code_record.sh.",
        "# Original labels are navigation. Graded meanings: docs/re/matching-semantic-index.json.",
        ".intel_syntax noprefix", ".code16", '.section .image,"ax",@progbits',
    ]
    instructions = []
    cursor = 0
    for row in mapping:
        start, end = row["file_start"], row["file_end"]
        require(start == cursor and start < end <= SIZE, "baseline map hole or overlap")
        cursor = end
        if row["kind"] != "instruction":
            continue
        instructions.append({**row, "gas": lines[row["assembly_line"] - 1]})
    require(cursor == SIZE and len(instructions) == 24376, "baseline coverage differs")
    require(sum(r["file_end"] - r["file_start"] for r in instructions) == 55392, "baseline code bytes differ")
    supplement_raw = (repo / SUPPLEMENT).read_bytes()
    supplement = json.loads(supplement_raw)
    require(supplement["schema"] == "wolong-matching-code-supplement-v1" and supplement["input_sha256"] == EXPECTED, "supplement identity differs")
    extra = supplement["instructions"]
    require(len(extra) == 338 and sum(r["file_end"] - r["file_start"] for r in extra) == 805, "supplement coverage differs")
    for row in extra:
        start, end = row["file_start"], row["file_end"]
        require(raw[start:end] == bytes.fromhex(row["bytes"]), "supplement bytes differ")
        require(any(r["kind"] == "data-or-unknown" and r["file_start"] <= start < end <= r["file_end"] for r in mapping), "supplement overlaps baseline code")
    instructions.extend(extra)
    instructions.sort(key=lambda row: row["file_start"])
    ranges = []
    count, cursor = 0, 0
    for row in instructions:
        start, end = row["file_start"], row["file_end"]
        require(cursor <= start < end <= SIZE, "instruction overlap or overflow")
        cursor = end
        count += 1
        public.append(f".org {start:#x}")
        name = row["original_name"]
        if name and re.fullmatch(r"[A-Za-z_][A-Za-z_0-9]*", name):
            public.append(name + ":")
        original = row["original_assembly"].split(";", 1)[0].strip()
        function = row["function_name"]
        match = re.fullmatch(r"sub_([0-9A-Fa-f]+)", function or "")
        semantic = semantics["functions"].get(f"0x{int(match[1], 16):X}") if match else None
        warning = "" if semantic and semantic["level"] == "proven" else "⚠ "
        level = semantic["level"] if semantic else "unknown"
        meaning = semantic["meaning"] if semantic else "未附加語意"
        sources = ", ".join(semantic["sources"]) if semantic else "無"
        public.append(f'# IDA {row["ida_linear"]:#x}; file {start:#x}; locator=proven; {original}')
        public.append(f"# {warning}函式 {function or '無'}; [{level}] {meaning}; 出處: {sources}")
        if "opcode" in row:
            public.append(f'# [proven] opcode {row["opcode"]}; 原 IDA 資料行: {row["original_ida_data_line"]}; 出處: {supplement["evidence"]}')
        for i, op in enumerate(row["operands"]):
            extra = semantics["operands"].get(f'0x{row["ida_linear"]:X}:{i}')
            if extra:
                require(extra["original"] == op["original"], "graded operand locator differs")
                warning = "" if extra["level"] == "proven" else "⚠ "
                public.append(f'# {warning}operand {i} {op["original"]}; [{extra["level"]}] {extra["meaning"]}; 出處: {", ".join(extra["sources"])}')
        public.append(row["gas"])
        if ranges and ranges[-1][1] == start:
            ranges[-1][1] = end
        else:
            ranges.append([start, end])
    require(count == INSTRUCTIONS, "instruction coverage differs")
    require(sum(end - start for start, end in ranges) == CODE_BYTES, "instruction byte count differs")
    public.append(f".org {SIZE:#x}")
    out = ("\n".join(public) + "\n").encode("utf-8")
    (repo / SOURCE).write_bytes(out)
    (repo / LINKER).write_bytes(linker)
    write_json(repo / MANIFEST, {
        "schema": "wolong-matching-code-record-v1",
        "input": "DOS/V KI.EXE", "input_sha256": EXPECTED, "input_size": SIZE,
        "instruction_count": count, "instruction_bytes": CODE_BYTES,
        "baseline_instruction_count": 24376, "supplement_instruction_count": 338,
        "supplement": SUPPLEMENT, "supplement_sha256": digest(supplement_raw),
        "supplement_database_sha256": supplement["database_sha256"],
        "private_import_bytes": SIZE - CODE_BYTES,
        "address_space": "IDA database linear base 0x10000; file = linear - 0x10000 + 512",
        "ida_version": report["ida_version"], "ida_image_id": report["ida_image_id"],
        "database_sha256": report["database_sha256"],
        "inventory_sha256": report["inventory_sha256"],
        "semantic_index": "docs/re/matching-semantic-index.json",
        "semantic_index_sha256": digest(semantic_path.read_bytes()),
        "record_tool": "tools/matching_code_record.py",
        "record_tool_sha256": digest((repo / "tools/matching_code_record.py").read_bytes()),
        "baseline_source_sha256": digest(source), "baseline_linker_sha256": digest(linker),
        "source": SOURCE, "source_sha256": digest(out),
        "linker": LINKER, "linker_sha256": digest(linker),
        "build_image_id": report["build_image_id"], "assembler_version": report["assembler_version"],
        "evidence": "docs/re/90-assembly-reconstruction.md", "byte_match_level": "proven",
        "semantic_level": "unknown; consult graded semantic index per original address",
        "code_ranges": [{"file_start": a, "file_end": b, "sha256": digest(raw[a:b])} for a, b in ranges],
    })
    print(f"Recorded {count} instructions, {CODE_BYTES} code bytes, {len(ranges)} ranges; private bytes omitted")


def check_manifest(manifest: dict) -> None:
    require(manifest["schema"] == "wolong-matching-code-record-v1", "manifest schema differs")
    require(manifest["input_sha256"] == EXPECTED and manifest["input_size"] == SIZE, "manifest input differs")
    require(manifest["instruction_count"] == INSTRUCTIONS, "manifest instruction count differs")
    require(manifest["instruction_bytes"] == CODE_BYTES, "manifest code byte count differs")
    require(manifest["private_import_bytes"] == SIZE - CODE_BYTES, "private byte count differs")
    cursor, total = 0, 0
    for row in manifest["code_ranges"]:
        a, b = row["file_start"], row["file_end"]
        require(cursor <= a < b <= SIZE, "code range overlap or overflow")
        cursor = b
        total += b - a
    require(total == CODE_BYTES, "code range coverage differs")


def run(args: list[str]) -> str:
    result = subprocess.run(args, check=True, capture_output=True, text=True, timeout=45)
    return result.stdout


def compile_source(source: Path, linker: Path, output: Path, stem: str) -> bytes:
    obj, elf, binary = (output / (stem + suffix) for suffix in (".o", ".elf", ".code.bin"))
    run(["as", "--32", "-o", str(obj), str(source)])
    run(["ld", "-m", "elf_i386", "-T", str(linker), "--entry", "0", "-o", str(elf), str(obj)])
    run(["objcopy", "-O", "binary", "--only-section=.image", str(elf), str(binary)])
    return binary.read_bytes()


def assemble(repo: Path, output: Path, image_id: str) -> None:
    raw = original_bytes(repo)
    manifest_path = repo / MANIFEST
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    check_manifest(manifest)
    require(digest((repo / manifest["supplement"]).read_bytes()) == manifest["supplement_sha256"], "supplement source differs")
    require(digest(Path(__file__).read_bytes()) == manifest["record_tool_sha256"], "record tool source differs")
    require(image_id == manifest["build_image_id"], "build image differs from verified toolchain")
    source, linker = repo / SOURCE, repo / LINKER
    require(digest(source.read_bytes()) == manifest["source_sha256"], "recorded code source differs")
    require(digest(linker.read_bytes()) == manifest["linker_sha256"], "recorded linker source differs")
    text = source.read_text(encoding="utf-8")
    require(len(re.findall(r"^# IDA ", text, re.M)) == INSTRUCTIONS, "source instruction markers differ")
    require(not re.search(r"\.(?:byte|word|long|quad|incbin|ascii|asciz)\b", text), "data directive in code source")
    output.mkdir(parents=True, exist_ok=True)
    compiled = compile_source(source, linker, output, "KI")
    require(len(compiled) == SIZE, "assembled image size differs")
    reconstructed = bytearray(SIZE)
    cursor = 0
    for row in manifest["code_ranges"]:
        a, b = row["file_start"], row["file_end"]
        require(not any(compiled[cursor:a]), "assembled bytes outside recorded code ranges")
        require(digest(compiled[a:b]) == row["sha256"], "assembled code range hash differs")
        require(compiled[a:b] == raw[a:b], "assembled instruction bytes differ")
        reconstructed[cursor:a] = raw[cursor:a]
        reconstructed[a:b] = compiled[a:b]
        cursor = b
    require(not any(compiled[cursor:]), "assembled tail outside recorded code ranges")
    reconstructed[cursor:] = raw[cursor:]
    require(bytes(reconstructed) == raw and digest(reconstructed) == EXPECTED, "whole EXE differs")
    (output / "KI.EXE").write_bytes(reconstructed)
    controls = []
    for kind in ("instruction", "linker-constant"):
        mutant_source, mutant_linker = text, linker.read_text(encoding="utf-8")
        if kind == "instruction":
            mutant_source, n = re.subn(r"^clc$", "stc", mutant_source, count=1, flags=re.M)
        else:
            mutant_linker, n = re.subn(
                r"^(imm_at_[A-F0-9]+) = (0x[0-9a-f]+);$",
                lambda m: f"{m[1]} = {int(m[2], 16) ^ 1:#x};", mutant_linker, count=1, flags=re.M,
            )
        require(n == 1, "negative control target missing")
        ms, ml = output / (kind + ".S"), output / (kind + ".ld")
        ms.write_text(mutant_source, encoding="utf-8")
        ml.write_text(mutant_linker, encoding="utf-8")
        altered = compile_source(ms, ml, output, kind)
        changed = sum(altered[a:b] != raw[a:b] for a, b in
                      ((r["file_start"], r["file_end"]) for r in manifest["code_ranges"]))
        require(changed > 0, "negative control accepted")
        controls.append({"kind": kind, "rejected": True, "changed_code_ranges": changed})
    # 私有資料只接受固定身分的原版，不能以不同版本填補通過。
    altered_input = bytes([raw[0] ^ 1]) + raw[1:]
    require(digest(altered_input) != EXPECTED, "input identity control accepted")
    controls.append({"kind": "private-input-identity", "rejected": True})
    result = {
        "schema": "wolong-matching-code-build-v1", "whole_file_exact": True,
        "input_sha256": EXPECTED, "rebuilt_sha256": digest(reconstructed), "rebuilt_size": SIZE,
        "instruction_count": INSTRUCTIONS, "compiled_instruction_bytes": CODE_BYTES,
        "imported_instruction_bytes": 0, "imported_private_bytes": SIZE - CODE_BYTES,
        "all_noncode_placeholders_zero": True, "manifest_sha256": digest(manifest_path.read_bytes()),
        "source_sha256": digest(source.read_bytes()), "linker_sha256": digest(linker.read_bytes()),
        "record_tool_sha256": manifest["record_tool_sha256"],
        "build_image_id": image_id, "assembler_version": run(["as", "--version"]).splitlines()[0],
        "linker_version": run(["ld", "--version"]).splitlines()[0],
        "objcopy_version": run(["objcopy", "--version"]).splitlines()[0],
        "negative_controls": controls, "c_machine_code_match": False,
    }
    write_json(output / "code-record-verification.json", result)
    print(f"Recorded source: {INSTRUCTIONS} instructions / {CODE_BYTES} bytes exact; whole {SIZE}-byte EXE exact; 3 controls rejected")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("action", choices=("export", "build"))
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--output", type=Path)
    parser.add_argument("--image-id")
    args = parser.parse_args()
    if args.action == "export":
        export(args.repo)
    else:
        if args.output is None or args.image_id is None:
            parser.error("build requires --output and --image-id")
        assemble(args.repo, args.output, args.image_id)
