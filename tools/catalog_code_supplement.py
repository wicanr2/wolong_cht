#!/usr/bin/env python3
"""以本輪之前的固定組語來源，匹配清單家族新增的原始指令。"""

import argparse
import hashlib
import importlib.util
import json
from pathlib import Path


IMAGE_ID = "sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e"
EVIDENCE = "docs/re/119-c-list-families-restoration.md"


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def generate(repo, out):
    if not out.is_dir():
        raise ValueError("output must be an existing explicitly mounted directory")
    translator = repo / "tools/assembly_rebuild.py"
    spec = importlib.util.spec_from_file_location("assembly", translator)
    asm = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(asm)
    research = repo / "workplace/matching-decompilation/c-catalog"
    probe_path = research / "ida/ida-probe.json"
    database_path = research / "ida/input.exe.i64"
    probe = json.loads(probe_path.read_text(encoding="utf-8"))
    original = (repo / "workplace/orig/dosv/KI.EXE").read_bytes()
    input_sha = hashlib.sha256(original).hexdigest()
    if input_sha != asm.EXPECTED or probe["input_sha256"] != input_sha:
        raise ValueError("original and IDA probe identities differ")
    baseline_path = repo / "workplace/matching-decompilation/assembly/build/source-map.json"
    historical_paths = [repo / "docs/re/rectangle-handler-code.json", repo / "docs/re/c-list-code.json"]
    baseline = json.loads(baseline_path.read_text(encoding="utf-8"))
    covered_rows = [r for r in baseline if r["kind"] == "instruction"]
    coverage_inputs = {
        str(baseline_path.relative_to(repo)): {
            "sha256": sha(baseline_path), "instruction_count": len(covered_rows),
        },
    }
    for historical_path, expected_count in zip(historical_paths, [370, 83]):
        historical = json.loads(historical_path.read_text(encoding="utf-8"))
        if historical["input_sha256"] != input_sha or len(historical["instructions"]) != expected_count:
            raise ValueError("historical supplement identity or instruction count differs")
        for entry in historical["instructions"]:
            if original[entry["file_start"]:entry["file_end"]].hex() != entry["bytes"]:
                raise ValueError("historical supplement bytes differ from original")
        coverage_inputs[str(historical_path.relative_to(repo))] = {
            "sha256": sha(historical_path), "instruction_count": expected_count,
        }
        covered_rows.extend(historical["instructions"])
    coverage = []
    for start, stop in sorted((r["file_start"], r["file_end"]) for r in covered_rows):
        if coverage and start <= coverage[-1][1]:
            coverage[-1] = (coverage[-1][0], max(coverage[-1][1], stop))
        else:
            coverage.append((start, stop))
    instructions = [i for t in probe["targets"] for c in t["chunks"] for i in c["instructions"]]
    instructions += [i for b in probe["decoded_blocks"] for i in b["instructions"]]
    unique = {}
    for insn in instructions:
        at = insn["ida_linear"] - 0x10000 + probe["mz_header_size"]
        raw = bytes.fromhex(insn["bytes"])
        end = at + len(raw)
        if insn.get("file_offset", at) != at or not raw or original[at:end] != raw:
            raise ValueError(f"IDA instruction disagrees with original file at {at:#x}")
        if at in unique and unique[at]["bytes"] != insn["bytes"]:
            raise ValueError(f"conflicting instruction decodings at {at:#x}")
        unique[at] = {
            **insn, "file_offset": at,
            "relocation_file_offsets": [r for r in probe["relocation_file_offsets"] if at <= r < end],
        }
    rows = []
    for at, insn in sorted(unique.items()):
        end = at + len(bytes.fromhex(insn["bytes"]))
        if any(start <= at and end <= stop for start, stop in coverage):
            continue
        if any(max(start, at) < min(stop, end) for start, stop in coverage):
            raise ValueError(f"instruction partially overlaps historical coverage at {at:#x}")
        if rows and rows[-1]["file_offset"] + len(bytes.fromhex(rows[-1]["bytes"])) > at:
            raise ValueError(f"new instructions overlap at {at:#x}")
        rows.append(insn)
    trials_path = out / "assembly-candidates"
    trials_path.mkdir(exist_ok=True)
    matched, unmatched, trials, symbols = asm.compile_candidates(rows, trials_path)
    if unmatched:
        (out / "unmatched-instructions.json").write_text(json.dumps(unmatched, indent=2) + "\n", encoding="utf-8")
        raise ValueError(f"{len(unmatched)} unmatched instructions; code-byte fallback forbidden")
    source = ["# Local instruction patch only; gaps are not an EXE reconstruction.",
              ".intel_syntax noprefix", ".code16", '.section .image,"ax",@progbits']
    entries = []
    for insn in rows:
        at = insn["file_offset"]
        end = at + len(bytes.fromhex(insn["bytes"]))
        expression = matched[at]
        entries.append({
            "ida_linear": insn["ida_linear"], "file_start": at, "file_end": end,
            "kind": "instruction", "original_name": "", "function_name": None,
            "original_assembly": insn["assembly"], "gas": expression, "bytes": insn["bytes"],
            "original_ida_data_line": insn.get("original_ida_data_line"),
            "ida_code_classified": insn.get("originally_code", True),
            "operands": insn["operands"], "locator_level": "proven", "scope_sources": [EVIDENCE],
        })
        source.extend([f".org {at:#x}", expression])
    source_path = out / "catalog-patch.S"
    source_path.write_text("\n".join(source) + "\n", encoding="utf-8")
    linker_path = out / "catalog-patch.ld"
    linker_path.write_text("\n".join(f"{name} = {value:#x};" for name, value in symbols.items()) + "\n", encoding="utf-8")
    asm.run(["as", "--32", "-o", str(out / "catalog-patch.o"), str(source_path)])
    asm.run(["ld", "-m", "elf_i386", "-T", str(linker_path), "--section-start", ".image=0", "--entry", "0",
             "-o", str(out / "catalog-patch.elf"), str(out / "catalog-patch.o")])
    asm.run(["objcopy", "-O", "binary", "--only-section=.image", str(out / "catalog-patch.elf"), str(out / "catalog-patch.bin")])
    binary = (out / "catalog-patch.bin").read_bytes()
    for entry in entries:
        at, end = entry["file_start"], entry["file_end"]
        if binary[at:end] != original[at:end]:
            raise ValueError(f"linked encoding differs at {at:#x}")
    count = len(entries)
    byte_count = sum(e["file_end"] - e["file_start"] for e in entries)
    verification = {
        "status": "byte-exact", "instructions": count, "bytes": byte_count,
        "input_instruction_rows": len(instructions), "unique_input_instructions": len(unique),
        "covered_input_instructions": len(unique) - count,
        "coverage_source": "fixed original source-map instruction rows plus prior rectangle/list supplements",
        "coverage_inputs": coverage_inputs, "source_sha256": sha(source_path),
        "linker_sha256": sha(linker_path), "binary_sha256": sha(out / "catalog-patch.bin"),
        "image_id": IMAGE_ID, "assembler": asm.run(["as", "--version"]).splitlines()[0],
        "linker": asm.run(["ld", "--version"]).splitlines()[0],
        "generator_sha256": sha(Path(__file__)), "translator_sha256": sha(translator),
        "instruction_byte_fallbacks": 0, "constant_symbols": symbols, "trials": trials,
    }
    result = {
        "schema": "wolong-matching-code-supplement-v1", "input_sha256": input_sha,
        "tool": probe["tool"], "tool_version": probe["tool_version"],
        "database_sha256": sha(database_path), "probe_sha256": sha(probe_path),
        "evidence": EVIDENCE, "constant_symbols": symbols,
        "instructions": entries, "assembler_verification": verification,
    }
    (out / "catalog-code-supplement.json").write_text(json.dumps(result, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"{count} previously uncovered instructions / {byte_count} bytes: native GNU assembly exact")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.output)
