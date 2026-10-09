#!/usr/bin/env python3
"""獨立重組 engagement v5 指令，保留固定舊覆蓋與錯誤邊界的替換證據。"""
from __future__ import annotations

import argparse
from bisect import bisect_left
from collections import Counter
import hashlib
import importlib.util
import json
from pathlib import Path
import struct

IMAGE_ID = "sha256:474f41ef91c354dd4754b08ef9302e965271e32417d6fba1772aecca0a5f9e2e"
PROBE = "workplace/matching-decompilation/c-engagement/ida-closed-v5/ida-probe.json"
PROBE_SHA = "32ebc39f59e53f16d18f2f4d5bb16ae8f6971ef2dd4e733976a93dd8f296a5d9"
EVIDENCE = "docs/re/131-c-engagement-tactical-restoration.md"
BASELINE = "workplace/matching-decompilation/assembly/build/source-map.json"
BASELINE_ASM = "workplace/matching-decompilation/assembly/build/KI.reconstructed.S"
HISTORICAL = (
    ("docs/re/rectangle-handler-code.json", 370, 898),
    ("docs/re/c-list-code.json", 83, 215),
    ("docs/re/c-catalog-code.json", 170, 361),
    ("docs/re/route-handler-code.json", 79, 179),
)


def sha(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))


def write_json(path: Path, value) -> None:
    path.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def normalized_instructions(probe, original):
    """驗證 IDA 載入後 bytes，另建 file-encoding 副本；原運算元保持不動。"""
    fields = struct.unpack_from("<14H", original)
    header, table = fields[4] * 16, fields[12]
    relocations = [header + off + segment * 16 for off, segment in
                   (struct.unpack_from("<HH", original, table + i * 4) for i in range(fields[3]))]
    if header != probe["mz_header_size"] or relocations != probe["relocation_file_offsets"]:
        raise ValueError("original MZ header/relocation table differs from the candidate")
    if header != 512 or probe["ida_load_paragraph"] != 0x1000:
        raise ValueError("unexpected IDA address-space contract")
    loaded = bytearray(original)
    for at in relocations:
        struct.pack_into("<H", loaded, at, (struct.unpack_from("<H", original, at)[0] + 0x1000) & 0xffff)
    unique, normalization, inline = {}, [], []
    memberships, membership_bytes = Counter(), Counter()
    named_chunks = 0

    def add(insn, owner, kind, original_name):
        at = insn["ida_linear"] - 0x10000 + header
        data = bytes.fromhex(insn["bytes"])
        end = at + len(data)
        if insn.get("file_offset", at) != at or loaded[at:end] != data:
            raise ValueError(f"candidate loader bytes differ at IDA {insn['ida_linear']:#x}")
        rr = []
        for offset in relocations:
            if not at <= offset < end:
                continue
            if offset + 2 > end:
                raise ValueError("MZ relocation splits an instruction boundary")
            rr.append({"file_offset": offset, "original_word": struct.unpack_from("<H", original, offset)[0],
                       "ida_word": struct.unpack_from("<H", loaded, offset)[0]})
        row = {**insn, "file_offset": at, "bytes": original[at:end].hex(),
               "ida_loaded_bytes": insn["bytes"], "loader_relocations": rr,
               "relocation_file_offsets": [r["file_offset"] for r in rr],
               "scope_owners": [{"kind": kind, "ida_linear": owner, "original_name": original_name}]}
        if at in unique:
            previous = unique[at]
            if previous["bytes"] != row["bytes"] or previous["ida_loaded_bytes"] != row["ida_loaded_bytes"]:
                raise ValueError("candidate contains conflicting interpretations at one instruction start")
            previous["scope_owners"].extend(row["scope_owners"])
        else:
            unique[at] = row
        if rr:
            if len(data) != 5 or original[at] != 0x9a or len(rr) != 1:
                raise ValueError("unexpected recovery MZ relocation instruction form")
            normalization.append({"ida_linear": insn["ida_linear"], "file_start": at, "file_end": end,
                                  "original_assembly": insn["assembly"], "original_operands": insn["operands"],
                                  "ida_loaded_bytes": insn["bytes"], "file_bytes": row["bytes"],
                                  "loader_relocations": rr, "scope_owners": row["scope_owners"],
                                  "normalization": "Forward-load original words, verify the candidate, then restore file bytes only in a copied instruction row",
                                  "locator_level": "proven"})
        memberships[kind] += 1
        membership_bytes[kind] += len(data)

    recovery = set(probe["recovery_targets"])
    for target in probe["targets"]:
        if target["ida_linear"] not in recovery:
            continue
        for chunk in target["chunks"]:
            start = chunk["start"] - 0x10000 + header
            end = chunk["end"] - 0x10000 + header
            if original[start:end].hex() != chunk["file_bytes"] or loaded[start:end].hex() != chunk["bytes"]:
                raise ValueError("named chunk identity differs beyond the original MZ loader transform")
            expected = [offset for offset in relocations if start <= offset < end]
            if sorted(r["file_offset"] for r in chunk["loader_relocations"]) != sorted(expected):
                raise ValueError("named chunk does not declare all original MZ relocations")
            for declared in chunk["loader_relocations"]:
                offset = declared["file_offset"]
                if (declared["original_word"], declared["ida_word"]) != (
                    struct.unpack_from("<H", original, offset)[0], struct.unpack_from("<H", loaded, offset)[0]):
                    raise ValueError("declared MZ relocation words differ")
            named_chunks += end - start
            cursor = chunk["start"]
            for insn in chunk["instructions"]:
                if cursor < insn["ida_linear"]:
                    a, b = cursor - 0x10000 + header, insn["ida_linear"] - 0x10000 + header
                    inline.append({"ida_start": cursor, "ida_end": insn["ida_linear"],
                                   "file_start": a, "file_end": b, "bytes": original[a:b].hex(),
                                   "original_owner": target["name"], "instruction": False})
                add(insn, target["ida_linear"], "named", target["name"])
                cursor = insn["ida_linear"] + len(bytes.fromhex(insn["bytes"]))
            if cursor < chunk["end"]:
                a, b = cursor - 0x10000 + header, chunk["end"] - 0x10000 + header
                inline.append({"ida_start": cursor, "ida_end": chunk["end"], "file_start": a, "file_end": b,
                               "bytes": original[a:b].hex(), "original_owner": target["name"], "instruction": False})
    for block in probe["decoded_blocks"]:
        for insn in block["instructions"]:
            add(insn, block["ida_linear"], "raw", block.get("original_name_before_analysis", ""))
    rows = [row for _, row in sorted(unique.items())]
    if memberships != {"named": 7080, "raw": 737} or membership_bytes != {"named": 16696, "raw": 1867}:
        raise ValueError("candidate effective instruction counts differ")
    if len(rows) != 7817 or named_chunks != 16706 or sum(r["file_end"]-r["file_start"] for r in inline) != 10:
        raise ValueError("candidate unique rows or preserved noninstruction ranges differ")
    if len(normalization) != 24:
        raise ValueError("expected all 24 recovery far-call MZ relocations")
    if any(a["file_offset"] + len(bytes.fromhex(a["bytes"])) > b["file_offset"] for a, b in zip(rows, rows[1:])):
        raise ValueError("candidate instruction boundaries overlap")
    return rows, normalization, inline


def fixed_coverage(repo, original):
    """固定五份輸入；目前 assembly manifest 不參與新增／退役判定。"""
    baseline = load(repo / BASELINE)
    baseline_sha = sha(repo / BASELINE)
    lines = (repo / BASELINE_ASM).read_text(encoding="utf-8").splitlines()
    old = []
    inputs = {}
    for row in baseline:
        if row["kind"] == "instruction":
            old.append({**row, "bytes": original[row["file_start"]:row["file_end"]].hex(),
                        "gas": lines[row["assembly_line"]-1], "coverage_source": BASELINE,
                        "coverage_source_sha256": baseline_sha})
    if len(old) != 24376 or sum(r["file_end"]-r["file_start"] for r in old) != 55392:
        raise ValueError("original source-map baseline differs")
    inputs[BASELINE] = {"sha256": baseline_sha, "instructions": 24376, "bytes": 55392}
    for name, count, byte_count in HISTORICAL:
        doc = load(repo / name)
        if doc["input_sha256"] != hashlib.sha256(original).hexdigest() or len(doc["instructions"]) != count:
            raise ValueError("historical supplement identity differs")
        if sum(r["file_end"]-r["file_start"] for r in doc["instructions"]) != byte_count:
            raise ValueError("historical supplement bytes differ")
        inputs[name] = {"sha256": sha(repo / name), "instructions": count, "bytes": byte_count}
        for row in doc["instructions"]:
            if original[row["file_start"]:row["file_end"]].hex() != row["bytes"]:
                raise ValueError("historical instruction bytes differ from original")
            old.append({**row, "coverage_source": name, "coverage_source_sha256": inputs[name]["sha256"]})
    old.sort(key=lambda r: r["file_start"])
    if len(old) != 25078 or sum(r["file_end"]-r["file_start"] for r in old) != 57045:
        raise ValueError("fixed coverage union differs")
    if any(a["file_end"] > b["file_start"] for a, b in zip(old, old[1:])):
        raise ValueError("historical instruction rows overlap")
    return old, inputs


def generate(repo: Path, out: Path):
    if not out.is_dir() or out.stat().st_uid != 1000:
        raise ValueError("output must be an existing explicitly mounted user-owned directory")
    if not (repo / EVIDENCE).is_file() or sha(repo / PROBE) != PROBE_SHA:
        raise ValueError("evidence entry or fixed candidate identity differs")
    translator = repo / "tools/assembly_rebuild.py"
    spec = importlib.util.spec_from_file_location("assembly", translator)
    asm = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(asm)
    original = (repo / "workplace/orig/dosv/KI.EXE").read_bytes()
    full_probe = load(repo / PROBE)
    input_sha = hashlib.sha256(original).hexdigest()
    if input_sha != asm.EXPECTED or full_probe["input_sha256"] != input_sha:
        raise ValueError("original/candidate input identity differs")
    if len(full_probe["recovery_targets"]) != 234 or len(full_probe["decoded_blocks"]) != 12:
        raise ValueError("candidate entry partition differs")
    protected_paths = [BASELINE, BASELINE_ASM, PROBE, "docs/re/matching-semantic-index.json"]
    protected_paths += [name for name, _, _ in HISTORICAL]
    before = {name: sha(repo / name) for name in protected_paths}
    rows, normalization, inline = normalized_instructions(full_probe, original)
    old, coverage_inputs = fixed_coverage(repo, original)
    starts = [r["file_start"] for r in old]
    by_start = {r["file_start"]: r for r in old}
    classes, overlap_map, new_rows, retired = {}, {}, [], {}
    for insn in rows:
        a, b = insn["file_offset"], insn["file_offset"] + len(bytes.fromhex(insn["bytes"]))
        exact = by_start.get(a)
        i, overlaps = max(0, bisect_left(starts, a)-1), []
        while i < len(old) and old[i]["file_start"] < b:
            if old[i]["file_end"] > a:
                overlaps.append(old[i])
            i += 1
        if exact and exact["file_end"] == b:
            classes[a] = "same-boundary"
        elif overlaps:
            classes[a] = "retired-boundary"
            new_rows.append(insn)
            overlap_map[a] = overlaps
            for row in overlaps:
                retired[row["file_start"]] = row
        else:
            classes[a] = "new-uncovered"
            new_rows.append(insn)
    if len(new_rows) != 31 or sum(len(bytes.fromhex(r["bytes"])) for r in new_rows) != 84:
        raise ValueError("expected exactly 31 replacement/addition rows and 84 bytes")
    if len(retired) != 12 or sum(r["file_end"]-r["file_start"] for r in retired.values()) != 28:
        raise ValueError("expected exactly 12 retired historical rows and 28 bytes")
    proposed = [(r["file_start"], r["file_end"]) for r in old if r["file_start"] not in retired]
    proposed += [(r["file_offset"], r["file_offset"]+len(bytes.fromhex(r["bytes"]))) for r in new_rows]
    proposed.sort()
    if any(a[1] > b[0] for a, b in zip(proposed, proposed[1:])):
        raise ValueError("replacement leaves overlapping instruction intervals")

    trials_path = out / "engagement-assembly-candidates"
    trials_path.mkdir(exist_ok=True)
    matched, unmatched, trials, symbols = asm.compile_candidates(rows, trials_path)
    if unmatched:
        write_json(out / "engagement-unmatched-instructions.json", unmatched)
        raise ValueError(f"{len(unmatched)} instructions did not assemble exact; byte fallback is forbidden")
    text = ["# Independent engagement v5 instruction proof; noncode gaps are zero.",
            ".intel_syntax noprefix", ".code16", '.section .image,"ax",@progbits']
    for insn in rows:
        text += [f'.org {insn["file_offset"]:#x}', matched[insn["file_offset"]]]
    source = out / "engagement-proof.S"
    source.write_text("\n".join(text)+"\n", encoding="utf-8")
    linker = out / "engagement-proof.ld"
    linker.write_text("\n".join(f"{name} = {value:#x};" for name, value in symbols.items())+"\n", encoding="utf-8")
    asm.run(["as", "--32", "-o", str(out/"engagement-proof.o"), str(source)])
    asm.run(["ld", "-m", "elf_i386", "-T", str(linker), "--section-start", ".image=0", "--entry", "0",
             "-o", str(out/"engagement-proof.elf"), str(out/"engagement-proof.o")])
    asm.run(["objcopy", "-O", "binary", "--only-section=.image", str(out/"engagement-proof.elf"), str(out/"engagement-proof.bin")])
    binary = (out/"engagement-proof.bin").read_bytes()
    mask = bytearray(len(binary))
    entries, coverage = [], []
    for insn in rows:
        a, b = insn["file_offset"], insn["file_offset"]+len(bytes.fromhex(insn["bytes"]))
        if binary[a:b] != original[a:b]:
            raise ValueError(f"linked instruction mismatch at IDA {insn['ida_linear']:#x}")
        mask[a:b] = b"\1" * (b-a)
        owners = insn["scope_owners"]
        named = next((owner["original_name"] for owner in owners if owner["kind"] == "named"), None)
        entry = {"ida_linear": insn["ida_linear"], "file_start": a, "file_end": b, "kind": "instruction",
                 "original_name": next((o["original_name"] for o in owners if o["ida_linear"] == insn["ida_linear"]), ""),
                 "function_name": named, "original_assembly": insn["assembly"], "gas": matched[a], "bytes": insn["bytes"],
                 "operands": insn["operands"], "scope_owners": owners, "locator_level": "proven", "scope_sources": [EVIDENCE]}
        if classes[a] != "same-boundary":
            entries.append(entry)
        coverage.append({**entry, "coverage_class": classes[a], "fixed_prior_coverage": classes[a] == "same-boundary",
                         "independent_native_assembly_match": True, "ida_loaded_bytes": insn["ida_loaded_bytes"],
                         "loader_relocations": insn["loader_relocations"],
                         "retired_prior_rows": [{"ida_linear": r["ida_linear"], "file_start": r["file_start"], "file_end": r["file_end"],
                                                 "source": r["coverage_source"], "source_sha256": r["coverage_source_sha256"]}
                                                for r in overlap_map.get(a, [])]})
    if any(value and not keep for value, keep in zip(binary, mask)):
        raise ValueError("independent proof emitted nonzero bytes outside instruction intervals")
    entries_by_start = {r["file_start"]: r for r in entries}
    retired_entries = []
    for a, row in sorted(retired.items()):
        replacements = [entries_by_start[start] for start, prior in overlap_map.items() if any(r["file_start"] == a for r in prior)]
        retired_entries.append({**row, "retirement_level": "proven",
                                "retirement_note": "Boundary conflict against candidate v5; original bytes remain unchanged",
                                "retirement_reason": "The effective candidate instruction begins at a different boundary and overlaps this historical decoding; append-only emission would duplicate instruction coverage",
                                "retirement_sources": [EVIDENCE, PROBE], "replacement_rows": replacements})
    new_symbols = {name: value for name, value in symbols.items() if any(name in entry["gas"] for entry in entries)}
    after = {name: sha(repo / name) for name in protected_paths}
    if before != after:
        raise ValueError("immutable candidate or historical sources changed during assembly")
    result = {
        "schema": "wolong-matching-code-supplement-v1", "input_sha256": input_sha,
        "tool": full_probe["tool"], "tool_version": full_probe["tool_version"],
        "database_sha256": sha(repo/Path(PROBE).parent/"input.exe.i64"), "probe_sha256": sha(repo/PROBE),
        "evidence": EVIDENCE, "constant_symbols": new_symbols, "instructions": entries,
        "retired_instructions": retired_entries, "instruction_coverage": coverage,
        "preserved_inline_data": inline, "loader_relocation_normalization": normalization,
        "named_instruction_count": 7080, "named_instruction_bytes": 16696, "named_chunk_extent_bytes": 16706,
        "raw_instruction_count": 737, "raw_instruction_bytes": 1867,
        "recovery_targets": full_probe["recovery_targets"], "navigation_only_targets": full_probe["navigation_only_targets"],
        "original_functions": [{"name": t["name"], "ida_linear": t["ida_linear"], "file_sha256": t["file_sha256"]}
                               for t in full_probe["targets"] if t["ida_linear"] in full_probe["recovery_targets"]],
        "original_code_blocks": [{k: block[k] for k in ["ida_linear", "end_ida_linear", "original_name_before_analysis"] if k in block}
                                 for block in full_probe["decoded_blocks"]],
        "retirement_summary": {"instructions": 12, "bytes": 28, "replacement_or_addition_instructions": 31,
                               "replacement_or_addition_bytes": 84, "new_uncovered_instructions": 21, "new_uncovered_bytes": 49,
                               "boundary_replacement_instructions": 10, "boundary_replacement_bytes": 35,
                               "predicted_total_instructions": 25097, "predicted_total_instruction_bytes": 57101},
        "assembler_verification": {
            "status": "byte-exact", "instructions": 31, "bytes": 84,
            "input_instruction_rows": len(rows), "unique_input_instructions": len(rows), "covered_input_instructions": 7786,
            "coverage_source": "fixed original source-map plus rectangle370/list83/catalog170/route79; no current-manifest cancellation",
            "coverage_inputs": coverage_inputs, "fixed_prior_instruction_count": 25078, "fixed_prior_instruction_bytes": 57045,
            "instruction_byte_fallbacks": 0, "target_filter": "234 recovery named entries plus 12 raw CFG entries; navigation excluded",
            "published_source_binding": "pending root merge; current S/linker hashes intentionally not bound",
            "independent_reassembly": {"status": "byte-exact", "instructions": len(rows), "bytes": 18563,
                "source_sha256": sha(source), "linker_sha256": sha(linker), "binary_sha256": sha(out/"engagement-proof.bin"),
                "constant_symbols": symbols, "trials": trials, "noninstruction_output_bytes": "zero placeholders only"},
            "image_id": IMAGE_ID, "assembler": asm.run(["as", "--version"]).splitlines()[0],
            "linker": asm.run(["ld", "--version"]).splitlines()[0],
            "generator_sha256": sha(Path(__file__)), "translator_sha256": sha(translator),
            "sources_before": before, "sources_after": after,
        },
        "address_space": "IDA database linear base 0x10000; original-file offsets include 512-byte MZ header; loader-normalized bytes are separate fields",
        "notes": ["Retired rows remain complete and traceable; the historical source-map and four supplements are immutable.",
                  "This instruction proof does not rebuild or modify the complete published executable.",
                  "Ten inline noninstruction bytes remain outside instruction coverage and must not be emitted as instruction fallbacks."],
    }
    write_json(out/"engagement-code-supplement.json", result)
    print(json.dumps({"instructions": len(rows), "bytes": 18563, "new_rows": len(entries), "new_bytes": 84,
                      "retired_rows": len(retired_entries), "retired_bytes": 28, "mz_relocated_instructions": len(normalization),
                      "new_constant_symbols": new_symbols, "all_constant_symbols": symbols,
                      "output_sha256": sha(out/"engagement-code-supplement.json")}, ensure_ascii=False))


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    generate(args.repo, args.output)
