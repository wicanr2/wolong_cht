#!/usr/bin/env python3
"""以固定舊組語覆蓋核對九個尋路與潰散函式及 raw 搜尋入口，獨立重組全部指令。"""

import argparse
import hashlib
import importlib.util
import json
import struct
from pathlib import Path


IMAGE_ID = "sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e"
EVIDENCE = "docs/re/129-c-route-restoration.md"
ROOTS = (0x1291A, 0x12977, 0x129C3, 0x12BA8, 0x147BB, 0x1487B, 0x14A0F, 0x19656, 0x196CF)
HISTORICAL = (
    ("docs/re/rectangle-handler-code.json", 370, 898),
    ("docs/re/c-list-code.json", 83, 215),
    ("docs/re/c-catalog-code.json", 170, 361),
)


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load(path):
    return json.loads(path.read_text(encoding="utf-8"))


def normalize_instructions(probe, original):
    """Restore file instruction bytes without changing IDA locator or operands."""
    fields = struct.unpack_from("<14H", original)
    header = fields[4] * 16
    table = fields[12]
    mz_relocations = [
        header + off + seg * 16
        for off, seg in (struct.unpack_from("<HH", original, table + i * 4) for i in range(fields[3]))
    ]
    if header != probe["mz_header_size"] or mz_relocations != probe["relocation_file_offsets"]:
        raise ValueError("MZ relocation table differs from the original file")
    if probe["ida_load_paragraph"] != 0x1000:
        raise ValueError("unexpected IDA load paragraph")
    rows, records = [], []
    for target in probe["targets"]:
        for chunk in target["chunks"]:
            start = chunk["start"] - 0x10000 + header
            raw = bytes.fromhex(chunk["file_bytes"])
            loaded = bytes.fromhex(chunk["bytes"])
            if len(raw) != chunk["end"] - chunk["start"] or original[start:start + len(raw)] != raw:
                raise ValueError("chunk file_bytes differs from the original")
            reconstructed = bytearray(raw)
            seen = []
            for relocation in chunk["loader_relocations"]:
                at = relocation["file_offset"]
                offset = at - start
                if at not in mz_relocations or offset < 0 or offset + 2 > len(raw):
                    raise ValueError("chunk relocation is outside the original MZ table")
                old = struct.unpack_from("<H", raw, offset)[0]
                new = (old + probe["ida_load_paragraph"]) & 0xFFFF
                if old != relocation["original_word"] or new != relocation["ida_word"]:
                    raise ValueError("chunk relocation word identity differs")
                struct.pack_into("<H", reconstructed, offset, new)
                seen.append(at)
            expected = [at for at in mz_relocations if start <= at < start + len(raw)]
            if sorted(seen) != sorted(expected) or bytes(reconstructed) != loaded:
                raise ValueError("IDA chunk differs beyond the declared MZ relocations")
            for insn in chunk["instructions"]:
                at = insn["ida_linear"] - 0x10000 + header
                offset = at - start
                size = len(bytes.fromhex(insn["bytes"]))
                if offset < 0 or offset + size > len(raw) or loaded[offset:offset + size].hex() != insn["bytes"]:
                    raise ValueError("IDA instruction bytes disagree with its loaded chunk")
                relocations = [
                    r for r in chunk["loader_relocations"] if at <= r["file_offset"] < at + size
                ]
                row = {**insn, "bytes": raw[offset:offset + size].hex(),
                       "ida_loaded_bytes": insn["bytes"], "loader_relocations": relocations}
                rows.append(row)
                if relocations:
                    if size != 5 or raw[offset] != 0x9A or len(relocations) != 1:
                        raise ValueError("unexpected relocated instruction form")
                    records.append({
                        "ida_linear": insn["ida_linear"], "file_start": at, "file_end": at + size,
                        "original_function_name": target["name"], "original_assembly": insn["assembly"],
                        "original_operands": insn["operands"], "ida_loaded_bytes": insn["bytes"],
                        "file_bytes": row["bytes"], "loader_relocations": relocations,
                        "normalization": "verify forward MZ loading, then restore file_bytes in a copied instruction row",
                        "locator_level": "proven",
                    })
    expected_offsets = []
    if sorted(r["loader_relocations"][0]["file_offset"] for r in records) != expected_offsets:
        raise ValueError("expected exactly zero recovery MZ relocations")
    return rows, records


def generate(repo, out):
    if not out.is_dir():
        raise ValueError("output must be an existing explicitly mounted directory")
    translator = repo / "tools/assembly_rebuild.py"
    spec = importlib.util.spec_from_file_location("assembly", translator)
    asm = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(asm)
    research = repo / "workplace/matching-decompilation/c-route"
    probe_path = research / "ida/ida-probe.json"
    database_path = research / "ida/input.exe.i64"
    full_probe = load(probe_path)
    if set(full_probe["recovery_targets"]) != set(ROOTS):
        raise ValueError("recovery target identity differs")
    nav = set(full_probe["navigation_only_targets"])
    if nav & set(ROOTS) or {t["ida_linear"] for t in full_probe["targets"]} != set(ROOTS) | nav:
        raise ValueError("navigation/recovery partition differs")
    probe = {**full_probe, "targets": sorted((t for t in full_probe["targets"] if t["ida_linear"] in ROOTS), key=lambda t: t["ida_linear"]),
             "decoded_blocks": []}
    original = (repo / "workplace/orig/dosv/KI.EXE").read_bytes()
    input_sha = hashlib.sha256(original).hexdigest()
    if input_sha != asm.EXPECTED or probe["input_sha256"] != input_sha:
        raise ValueError("original and IDA probe identities differ")
    if tuple(t["ida_linear"] for t in probe["targets"]) != ROOTS or probe["decoded_blocks"]:
        raise ValueError("route probe must contain exactly the nine fixed named functions")

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

    instructions, relocations = normalize_instructions(probe, original)
    raw_blocks = full_probe["decoded_blocks"]
    if len(raw_blocks) != 1:
        raise ValueError("expected exactly one raw search entry")
    block = raw_blocks[0]
    if (block["ida_linear"], block["end_ida_linear"], block["original_name_before_analysis"]) != (0x1491B, 0x14A0F, "loc_1491B"):
        raise ValueError("raw search locator differs")
    cursor = block["ida_linear"]
    raw_rows = []
    for insn in block["instructions"]:
        data = bytes.fromhex(insn["bytes"])
        at = insn["ida_linear"] - 0x10000 + probe["mz_header_size"]
        if insn["ida_linear"] != cursor or original[at:at + len(data)] != data:
            raise ValueError("raw search instructions are not contiguous original bytes")
        if any(at <= r < at + len(data) for r in probe["relocation_file_offsets"]):
            raise ValueError("unexpected raw-search MZ relocation")
        raw_rows.append({**insn, "ida_loaded_bytes": insn["bytes"], "loader_relocations": []})
        cursor += len(data)
    raw_bytes = b"".join(bytes.fromhex(row["bytes"]) for row in raw_rows)
    if cursor != 0x14A0F or len(raw_rows) != 107 or len(raw_bytes) != 244 or raw_rows[-1]["mnemonic"] != "retn":
        raise ValueError("raw search boundary/count differs")
    if hashlib.sha256(raw_bytes).hexdigest() != block["file_sha256"]:
        raise ValueError("raw search identity differs")
    instructions.extend(raw_rows)
    unique = {}
    for insn in instructions:
        at = insn["ida_linear"] - 0x10000 + probe["mz_header_size"]
        raw = bytes.fromhex(insn["bytes"])
        end = at + len(raw)
        if insn.get("file_offset", at) != at or not raw or original[at:end] != raw:
            raise ValueError(f"IDA instruction disagrees with original file at {at:#x}")
        if at in unique:
            raise ValueError(f"duplicate route instruction at {at:#x}")
        unique[at] = {
            **insn, "file_offset": at,
            "relocation_file_offsets": [r for r in probe["relocation_file_offsets"] if at <= r < end],
        }
    rows = [insn for _, insn in sorted(unique.items())]
    if len(rows) != 543 or sum(len(bytes.fromhex(r["bytes"])) for r in rows) != 1260:
        raise ValueError("route probe instruction totals differ")
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
    trials_path = out / "route-assembly-candidates"
    trials_path.mkdir(exist_ok=True)
    matched, unmatched, trials, symbols = asm.compile_candidates(rows, trials_path)
    if unmatched:
        (out / "route-unmatched-instructions.json").write_text(
            json.dumps(unmatched, indent=2) + "\n", encoding="utf-8")
        raise ValueError(f"{len(unmatched)} unmatched instructions; code-byte fallback forbidden")
    source = ["# Independent route instruction proof; gaps are zero placeholders.",
              ".intel_syntax noprefix", ".code16", '.section .image,"ax",@progbits']
    for insn in rows:
        source.extend([f'.org {insn["file_offset"]:#x}', matched[insn["file_offset"]]])
    source_path = out / "route-proof.S"
    source_path.write_text("\n".join(source) + "\n", encoding="utf-8")
    linker_path = out / "route-proof.ld"
    linker_path.write_text("\n".join(f"{name} = {value:#x};" for name, value in symbols.items()) + "\n",
                           encoding="utf-8")
    asm.run(["as", "--32", "-o", str(out / "route-proof.o"), str(source_path)])
    asm.run(["ld", "-m", "elf_i386", "-T", str(linker_path), "--section-start", ".image=0", "--entry", "0",
             "-o", str(out / "route-proof.elf"), str(out / "route-proof.o")])
    asm.run(["objcopy", "-O", "binary", "--only-section=.image", str(out / "route-proof.elf"),
             str(out / "route-proof.bin")])
    binary = (out / "route-proof.bin").read_bytes()
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
            "kind": "instruction",
            "original_name": "loc_1491B" if insn["ida_linear"] == 0x1491B else "",
            "function_name": next((t["name"] for t in probe["targets"] for c in t["chunks"]
                                   if c["start"] <= insn["ida_linear"] < c["end"]), None),
            "original_assembly": insn["assembly"], "gas": matched[at], "bytes": insn["bytes"],
            "operands": insn["operands"], "locator_level": "proven", "scope_sources": [EVIDENCE],
        }
        if at in new_offsets:
            entries.append(entry)
        ranges = [r for r in manifest["code_ranges"] if r["file_start"] <= at and end <= r["file_end"]]
        if not ranges and at not in new_offsets:
            raise ValueError(f"previously covered route instruction is outside the published manifest at {at:#x}")
        coverage_receipts.append({
            "ida_linear": insn["ida_linear"], "file_start": at, "file_end": end, "bytes": insn["bytes"],
            "fixed_prior_coverage": at not in new_offsets, "published_code_range": ranges[0] if ranges else None,
            "independent_native_assembly_match": True,
            "original_assembly": insn["assembly"], "gas": matched[at],
            "ida_loaded_bytes": insn["ida_loaded_bytes"],
            "loader_relocations": insn["loader_relocations"],
        })
    new_symbols = {s: v for s, v in symbols.items() if any(s in e["gas"] for e in entries)}
    verification = {
        "status": "byte-exact", "instructions": len(entries),
        "bytes": sum(e["file_end"] - e["file_start"] for e in entries),
        "input_instruction_rows": len(instructions), "unique_input_instructions": len(rows),
        "covered_input_instructions": len(rows) - len(entries),
        "coverage_source": "fixed original 24376 instruction rows plus rectangle370/list83/catalog170",
        "coverage_inputs": coverage_inputs, "instruction_byte_fallbacks": 0,
        "fixed_prior_instruction_count": 24999, "fixed_prior_instruction_bytes": 56866,
        "target_filter": "nine recovery functions plus raw loc_1491B through original RET; navigation functions excluded",
        "published_manifest": str(manifest_path.relative_to(repo)), "published_manifest_sha256": sha(manifest_path),
        "published_code_source": manifest["source"], "published_code_source_sha256": sha(code_source),
        "independent_reassembly": {
            "status": "byte-exact", "instructions": len(rows), "bytes": 1260,
            "source_sha256": sha(source_path), "linker_sha256": sha(linker_path),
            "binary_sha256": sha(out / "route-proof.bin"), "constant_symbols": symbols, "trials": trials,
        },
        "published_input_instructions": sum(row["published_code_range"] is not None for row in coverage_receipts),
        "pending_publication_instructions": sum(row["published_code_range"] is None for row in coverage_receipts),
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
        "loader_relocation_normalization": relocations,
        "named_instruction_count": 436, "named_instruction_bytes": 1016,
        "raw_instruction_count": 107, "raw_instruction_bytes": 244,
        "original_code_blocks": [{
            "ida_linear": block["ida_linear"], "end_ida_linear": block["end_ida_linear"],
            "original_name_before_analysis": block["original_name_before_analysis"],
            "boundary_authority": block["boundary_authority"],
            "source_entry": "raw_entry_1491B", "file_sha256": block["file_sha256"],
        }],
        "recovery_targets": list(ROOTS),
        "navigation_only_targets": full_probe["navigation_only_targets"],
        "original_functions": [{"name": t["name"], "ida_linear": t["ida_linear"],
                                **({"direct_callees": t["direct_callees"]} if "direct_callees" in t else {}),
                                "file_sha256": t["file_sha256"]} for t in probe["targets"]],
    }
    source_only = {
        "schema": result["schema"], "input_sha256": input_sha,
        "database_sha256": result["database_sha256"], "probe_sha256": result["probe_sha256"],
        "evidence": EVIDENCE, "constant_symbols": new_symbols, "instructions": entries,
    }
    source_bytes = (json.dumps(source_only, ensure_ascii=False, indent=2) + "\n").encode("utf-8")
    source_relative = "docs/re/route-handler-code.json"
    published_source = repo / source_relative
    if published_source.exists() and published_source.read_bytes() != source_bytes:
        raise ValueError("immutable route source differs from independently matched instructions")
    (out / "route-handler-code.json").write_bytes(source_bytes)
    result["instruction_source"] = {
        "path": source_relative, "sha256": hashlib.sha256(source_bytes).hexdigest(),
    }
    if source_only["instructions"] != result["instructions"] or source_only["constant_symbols"] != result["constant_symbols"]:
        raise ValueError("route source and coverage receipt instruction lists differ")
    (out / "route-code-supplement.json").write_text(
        json.dumps(result, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f'{len(entries)} previously uncovered instructions / {verification["bytes"]} bytes; '
          f'{len(rows)} instructions / 1260 bytes independently assembled exact')


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.output)
