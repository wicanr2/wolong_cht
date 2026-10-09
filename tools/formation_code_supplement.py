#!/usr/bin/env python3
"""以固定前輪組語覆蓋核對七個玩家編成函式，並獨立重組全部指令。"""

import argparse
import hashlib
import importlib.util
import json
from pathlib import Path


IMAGE_ID = "sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e"
EVIDENCE = "docs/re/120-c-player-formation-restoration.md"
ROOTS = (0x16C5E, 0x16C92, 0x16D56, 0x16D6F, 0x16DA8, 0x16DFD, 0x16E80)
HISTORICAL = (
    ("docs/re/rectangle-handler-code.json", 370, 898),
    ("docs/re/c-list-code.json", 83, 215),
    ("docs/re/c-catalog-code.json", 170, 361),
)


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load(path):
    return json.loads(path.read_text(encoding="utf-8"))


def generate(repo, out):
    if not out.is_dir():
        raise ValueError("output must be an existing explicitly mounted directory")
    translator = repo / "tools/assembly_rebuild.py"
    spec = importlib.util.spec_from_file_location("assembly", translator)
    asm = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(asm)
    research = repo / "workplace/matching-decompilation/c-formation"
    probe_path = research / "ida/ida-probe.json"
    database_path = research / "ida/input.exe.i64"
    probe = load(probe_path)
    original = (repo / "workplace/orig/dosv/KI.EXE").read_bytes()
    input_sha = hashlib.sha256(original).hexdigest()
    if input_sha != asm.EXPECTED or probe["input_sha256"] != input_sha:
        raise ValueError("original and IDA probe identities differ")
    if tuple(t["ida_linear"] for t in probe["targets"]) != ROOTS or probe["decoded_blocks"]:
        raise ValueError("formation probe must contain exactly the seven fixed named functions")

    # This fixed union decides additions. The current manifest is only an
    # independent coverage check; it cannot cancel a newly decoded instruction.
    baseline_path = repo / "workplace/matching-decompilation/assembly/build/source-map.json"
    covered = [r for r in load(baseline_path) if r["kind"] == "instruction"]
    if len(covered) != 24376 or sum(r["file_end"] - r["file_start"] for r in covered) != 55392:
        raise ValueError("fixed original source-map baseline differs")
    coverage_inputs = {
        str(baseline_path.relative_to(repo)): {
            "sha256": sha(baseline_path), "instruction_count": 24376, "instruction_bytes": 55392,
        },
    }
    for relative, count, byte_count in HISTORICAL:
        path = repo / relative
        historical = load(path)
        entries = historical["instructions"]
        if historical["input_sha256"] != input_sha or len(entries) != count:
            raise ValueError("historical supplement identity or instruction count differs")
        if sum(e["file_end"] - e["file_start"] for e in entries) != byte_count:
            raise ValueError("historical supplement byte count differs")
        for entry in entries:
            if original[entry["file_start"]:entry["file_end"]].hex() != entry["bytes"]:
                raise ValueError("historical supplement bytes differ from original")
        coverage_inputs[relative] = {
            "sha256": sha(path), "instruction_count": count, "instruction_bytes": byte_count,
        }
        covered.extend(entries)
    fixed_ranges = []
    for start, stop in sorted((r["file_start"], r["file_end"]) for r in covered):
        if fixed_ranges and start <= fixed_ranges[-1][1]:
            fixed_ranges[-1] = (fixed_ranges[-1][0], max(fixed_ranges[-1][1], stop))
        else:
            fixed_ranges.append((start, stop))

    instructions = [i for t in probe["targets"] for c in t["chunks"] for i in c["instructions"]]
    unique = {}
    for insn in instructions:
        at = insn["ida_linear"] - 0x10000 + probe["mz_header_size"]
        raw = bytes.fromhex(insn["bytes"])
        end = at + len(raw)
        if insn.get("file_offset", at) != at or not raw or original[at:end] != raw:
            raise ValueError(f"IDA instruction disagrees with original file at {at:#x}")
        if at in unique:
            raise ValueError(f"duplicate formation instruction at {at:#x}")
        unique[at] = {
            **insn, "file_offset": at,
            "relocation_file_offsets": [r for r in probe["relocation_file_offsets"] if at <= r < end],
        }
    rows = [insn for _, insn in sorted(unique.items())]
    if len(rows) != 247 or sum(len(bytes.fromhex(r["bytes"])) for r in rows) != 561:
        raise ValueError("formation probe instruction totals differ")
    new = []
    for insn in rows:
        at = insn["file_offset"]
        end = at + len(bytes.fromhex(insn["bytes"]))
        if any(start <= at and end <= stop for start, stop in fixed_ranges):
            continue
        if any(max(start, at) < min(stop, end) for start, stop in fixed_ranges):
            raise ValueError(f"instruction partially overlaps historical coverage at {at:#x}")
        new.append(insn)

    manifest_path = repo / "docs/re/assembly-code-record.json"
    manifest = load(manifest_path)
    if manifest["input_sha256"] != input_sha:
        raise ValueError("published assembly manifest input differs")
    code_source = repo / manifest["source"]
    if sha(code_source) != manifest["source_sha256"]:
        raise ValueError("published assembly source hash differs")
    for code_range in manifest["code_ranges"]:
        raw = original[code_range["file_start"]:code_range["file_end"]]
        if hashlib.sha256(raw).hexdigest() != code_range["sha256"]:
            raise ValueError("published assembly code-range bytes differ")

    # Compile every input instruction, including the already covered ones.
    # This gives a nonempty native-assembly proof when there are no additions.
    trials_path = out / "formation-assembly-candidates"
    trials_path.mkdir(exist_ok=True)
    matched, unmatched, trials, symbols = asm.compile_candidates(rows, trials_path)
    if unmatched:
        (out / "formation-unmatched-instructions.json").write_text(
            json.dumps(unmatched, indent=2) + "\n", encoding="utf-8")
        raise ValueError(f"{len(unmatched)} unmatched instructions; code-byte fallback forbidden")
    source = ["# Independent formation instruction proof; gaps are zero placeholders.",
              ".intel_syntax noprefix", ".code16", '.section .image,"ax",@progbits']
    for insn in rows:
        source.extend([f'.org {insn["file_offset"]:#x}', matched[insn["file_offset"]]])
    source_path = out / "formation-proof.S"
    source_path.write_text("\n".join(source) + "\n", encoding="utf-8")
    linker_path = out / "formation-proof.ld"
    linker_path.write_text("\n".join(f"{name} = {value:#x};" for name, value in symbols.items()) + "\n",
                           encoding="utf-8")
    asm.run(["as", "--32", "-o", str(out / "formation-proof.o"), str(source_path)])
    asm.run(["ld", "-m", "elf_i386", "-T", str(linker_path), "--section-start", ".image=0", "--entry", "0",
             "-o", str(out / "formation-proof.elf"), str(out / "formation-proof.o")])
    asm.run(["objcopy", "-O", "binary", "--only-section=.image", str(out / "formation-proof.elf"),
             str(out / "formation-proof.bin")])
    binary = (out / "formation-proof.bin").read_bytes()
    entries = []
    coverage_receipts = []
    new_offsets = {r["file_offset"] for r in new}
    for insn in rows:
        at = insn["file_offset"]
        end = at + len(bytes.fromhex(insn["bytes"]))
        if binary[at:end] != original[at:end]:
            raise ValueError(f"linked encoding differs at {at:#x}")
        entry = {
            "ida_linear": insn["ida_linear"], "file_start": at, "file_end": end,
            "kind": "instruction", "original_name": "", "function_name": None,
            "original_assembly": insn["assembly"], "gas": matched[at], "bytes": insn["bytes"],
            "operands": insn["operands"], "locator_level": "proven", "scope_sources": [EVIDENCE],
        }
        if at in new_offsets:
            entries.append(entry)
        ranges = [r for r in manifest["code_ranges"] if r["file_start"] <= at and end <= r["file_end"]]
        if not ranges:
            raise ValueError(f"formation instruction is outside the published manifest at {at:#x}")
        coverage_receipts.append({
            "ida_linear": insn["ida_linear"], "file_start": at, "file_end": end, "bytes": insn["bytes"],
            "fixed_prior_coverage": at not in new_offsets, "published_code_range": ranges[0],
            "independent_native_assembly_match": True,
        })
    new_symbols = {s: v for s, v in symbols.items() if any(s in e["gas"] for e in entries)}
    verification = {
        "status": "byte-exact", "instructions": len(entries),
        "bytes": sum(e["file_end"] - e["file_start"] for e in entries),
        "input_instruction_rows": len(instructions), "unique_input_instructions": len(rows),
        "covered_input_instructions": len(rows) - len(entries),
        "coverage_source": "fixed original 24376 instruction rows plus rectangle370/list83/catalog170",
        "coverage_inputs": coverage_inputs, "instruction_byte_fallbacks": 0,
        "published_manifest": str(manifest_path.relative_to(repo)), "published_manifest_sha256": sha(manifest_path),
        "published_code_source": manifest["source"], "published_code_source_sha256": sha(code_source),
        "independent_reassembly": {
            "status": "byte-exact", "instructions": len(rows), "bytes": 561,
            "source_sha256": sha(source_path), "linker_sha256": sha(linker_path),
            "binary_sha256": sha(out / "formation-proof.bin"), "constant_symbols": symbols, "trials": trials,
        },
        "image_id": IMAGE_ID, "assembler": asm.run(["as", "--version"]).splitlines()[0],
        "linker": asm.run(["ld", "--version"]).splitlines()[0],
        "generator_sha256": sha(Path(__file__)), "translator_sha256": sha(translator),
        "constant_symbols": new_symbols,
    }
    result = {
        "schema": "wolong-matching-code-supplement-v1", "input_sha256": input_sha,
        "tool": probe["tool"], "tool_version": probe["tool_version"],
        "database_sha256": sha(database_path), "probe_sha256": sha(probe_path),
        "evidence": EVIDENCE, "constant_symbols": new_symbols, "instructions": entries,
        "instruction_coverage": coverage_receipts, "assembler_verification": verification,
    }
    (out / "formation-code-supplement.json").write_text(
        json.dumps(result, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f'{len(entries)} previously uncovered instructions / {verification["bytes"]} bytes; '
          f'{len(rows)} instructions / 561 bytes independently assembled exact')


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.output)
