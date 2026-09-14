#!/usr/bin/env python3
"""驗證 `docs/re/43` 每一列的分流分類與 GitHub Issue 對應。

`docs/re/43-open-questions.md` 是生成的證據索引，不是工作清單。
本工具把每一列固定成可回查的指紋，並要求它有下列其中一種分流：

* `actionable`：目前玩家路徑、規則或交付 gate 的工作，連到既有 Issue；
* `merge-target`：可合併進既有 Issue 的重複／支援證據，連到目標 Issue；
* `evidence-only`：保留原版證據邊界，但不形成目前 remake 工作；
* `history`：已被現況取代的歷史敘述。

預設規則刻意保守：沒有明確命中現行 Issue 的列不會自動變成工作，
而會標成 `evidence-only`，並由報告列出供 #28 稽核。來源列的內容、
檔名或小節一旦改變，指紋就會改變，`verify` 會失敗而不是靜默沿用舊分類。
"""

from __future__ import annotations

import argparse
import collections
import datetime
import hashlib
import importlib.util
import json
import os
import re
import sys
from typing import Any


MANIFEST_REL = "docs/re/43-open-question-triage.json"
SCHEMA = "wolong-re-open-question-triage/1"
VALID_CLASSES = {"actionable", "merge-target", "evidence-only", "history"}


# 順序是有意義的：具體的玩家路徑先於較寬的「政略／戰鬥」關鍵字。
# 每條規則只把證據列連回已存在的 Issue，不會建立新的遠端工作。
RULES: tuple[tuple[int, re.Pattern[str], str], ...] = (
    (
        21,
        re.compile(r"Android|android|安卓|SAF|簽署|觸控|GPU|DPI|手機.*實機|實機.*手機|平板|docs/mobile/"),
        "Android 實機、觸控、GPU、DPI、SAF 或簽署驗收，併入 Issue #21。",
    ),
    (
        20,
        re.compile(r"公開包|公開 source package|Eten|外部.*字型|clean.*包|dist-public|-orig(?:\b|[`\\s])"),
        "公開包安裝、-orig 啟動或外部字型依賴，併入 Issue #20。",
    ),
    (
        19,
        re.compile(r"無音訊裝置|沒有音效裝置|音訊裝置|音效裝置|正常音訊|無音效|音訊.*崩潰"),
        "桌面音訊裝置缺失或正常播放驗收，併入 Issue #19。",
    ),
    (
        2,
        re.compile(r"Windows|macOS|原生.*實機|目標平台|跨平台.*驗收|原生 GUI"),
        "Windows／macOS 或目標平台原生驗收，併入 Issue #2。",
    ),
    (
        25,
        re.compile(r"存檔|SAVE\.DAT|SINARIO|sinario-save|序列化|round.?trip|savefile|存／讀檔|讀檔|spec/25"),
        "存檔、載入、序列化或事件狀態保存，併入 Issue #25。",
    ),
    (
        24,
        re.compile(r"BATTLE\.SCH|投射物|projectile|CH=0x20|普通箭|箭.*動畫"),
        "BATTLE.SCH 投射物、動畫幀或命中時序，併入 Issue #24。",
    ),
    (
        23,
        re.compile(r"事件\s*11\s*[／/]\s*12|11／12|災害|火災|暴動|物件.*timer|物件.*動畫|disaster"),
        "事件 11／12 災害物件、timer 或動畫，併入 Issue #23。",
    ),
    (
        26,
        re.compile(r"協力|協同進攻|合攻|Friendship.*gate|0xA3|④b|cooperation"),
        "協力／協同進攻的 ④b gate 或 Friendship 門檻，併入 Issue #26。",
    ),
    (
        27,
        re.compile(r"(?:sub_135AB|0x24).*(?:中立|neutral|門檻|越界)|(?:中立|neutral|門檻|越界).*(?:sub_135AB|0x24)"),
        "sub_135AB 的 0x24 中立／越界門檻，併入 Issue #27。",
    ),
    (
        29,
        re.compile(r"spec/132|spec/164|spec/169|spec/170|spec/171|spec/180|spec/183|spec/184|spec/194|sub_14325|sub_1699E|sub_14155|formAICorpsTo|AI.*分派|救援.*次數|AI.*行軍"),
        "AI 行軍分派、救援候選、佔用圖或君主出陣，併入 Issue #29。",
    ),
    (
        30,
        re.compile(r"playtest/68.*(被擋|到站|行軍)|playtest/69.*(玩家.*軍團|玩家攻|指令有沒有|行軍|到站)|re/85|spec/85-march|行軍途中被擋|玩家軍團.*不會到|玩家攻.*觸發"),
        "玩家軍團行軍阻塞、到站與戰鬥入口，併入 Issue #30。",
    ),
    (
        31,
        re.compile(r"spec/115|spec/121|spec/196|海戰|0xCA|類型 9|地形.*戰場|野戰.*地形|戰場.*選擇"),
        "野戰地形、海戰適性與戰場選擇分支，併入 Issue #31。",
    ),
    (
        32,
        re.compile(r"spec/85-latin|spec/87-latin|playtest/47-latin|半形語系|半形.*欄|戰場標題.*地名"),
        "半形語系固定欄位、戰場標題與多語系排版，併入 Issue #32。",
    ),
    (
        10,
        re.compile(r"196/7/16|90.?日|90天|召見|5/13|5 月 13|5/20|5/31|6/1|月結全面分歧|spec/161|playtest/116|playtest/118|playtest/119|spec/160|spec/163|spec/178|spec/181|spec/185|playtest/87|promo/dosv-realmachine|入佇列的隨機空格"),
        "召見後長流程與 90 日政略對拍，併入 Issue #10。",
    ),
    (
        11,
        re.compile(r"DI residual|Stage 1|1440F|143AF|residual|意圖.*據點|留守分支|spec/195"),
        "Stage 1 DI residual／意圖據點模型，併入 Issue #11。",
    ),
    (
        12,
        re.compile(r"parity_screens|畫面對拍.?gate|對拍 gate|總目標.*32|待補.*12|screen.?parity|畫面.*群組|docs/playtest/.*(對拍|畫面|截圖|parity)|docs/spec/.*(對拍|畫面|鏡頭)|playtest/70|playtest/81|playtest/90|re/88|spec/107|spec/125|spec/126|spec/147|spec/85"),
        "畫面對拍 gate、群組覆蓋或 32 組門檻，併入 Issue #12。",
    ),
    (
        13,
        re.compile(r"tick.?4|第.?4.?tick|兩個兵|2 個兵|士兵.*漂移|0/2/0|0/5/2|spec/200"),
        "戰術第 4 tick 後兩名士兵的中繼目標／狀態差異，併入 Issue #13。",
    ),
    (
        16,
        re.compile(r"移動排程|命令批次|移動批次|HP.?128|中間目標|movement.?schedul|tactical movement|體力.*128|spec/159|spec/133|playtest/114|同格多兵|spec/36|spec/63|spec/95|playtest/63"),
        "戰術移動批次、中間目標與 HP 更新時序，併入 Issue #16。",
    ),
    (
        14,
        re.compile(r"開場對白|對白框.*時間|對白框.*收|battle.*talk|field.*第.?52|talk.*timing|playtest/59|spec/112"),
        "戰場開場對白框生命週期與輸入阻擋，併入 Issue #14。",
    ),
    (
        17,
        re.compile(r"單挑|duel|玩家.*攻方|玩家.*守方|角色.*開場|野戰單挑|spec/117|playtest/117"),
        "單挑角色方向與開場可操作 tick，併入 Issue #17。",
    ),
    (
        18,
        re.compile(r"戰後.*結算|戰後.*離場|世界退出|世界離場|post.?battle.*outcome|敗北.*世界|死亡.*戰後|戰術完整結算|完整狀態對拍|戰後兵力|戰後士氣|playtest/115|playtest/80|spec/128|spec/129|spec/141|spec/47"),
        "戰後結果、世界離場與多種結束方式，併入 Issue #18。",
    ),
    (
        15,
        re.compile(r"委任戰力|委任.*混合|delegated|長流程|12 個窄|power mix|spec/158|playtest/45"),
        "委任戰力混合公式與長流程回歸，併入 Issue #15。",
    ),
    (
        22,
        re.compile(r"事件\s*[2-9]|政略事件|君主.*採納|TALK.*事件|事件.*完整流程|事件.*訊息|playtest/34|spec/44|spec/45|spec/49|spec/123|spec/142"),
        "政略事件 2–9 的輸入、採納、訊息與 UI，併入 Issue #22。",
    ),
    (
        8,
        re.compile(r"音色|諧波|頻譜|harmonic|spec/122|spec/29|playtest/25|playtest/26"),
        "原版音色與諧波量化，併入 Issue #8。",
    ),
    (
        7,
        re.compile(r"ICONGRF|segment.?1|段.?1.*圖塊|0x3000.*背板"),
        "ICONGRF segment 1 UI 圖塊語意，併入 Issue #7。",
    ),
    (
        4,
        re.compile(r"\+0x21|0x21.*軍團|軍團.*0x21"),
        "軍團記錄 +0x21 欄位，併入 Issue #4。",
    ),
    (
        5,
        re.compile(r"對峙|standoff|96 拍|對砍|standoffBlocks|96.*動畫"),
        "對峙 96 拍動畫與音效，併入 Issue #5。",
    ),
    (
        3,
        re.compile(r"等距|tie.?break|retreatEndpoint|loc_1491B|兩端到首都|廣度優先.*退"),
        "退卻等距端點與 BFS tie-break，併入 Issue #3。",
    ),
    (
        6,
        re.compile(r"退卻.*BFS|BFS|自我修改碼|inline.*立即值|149B[68]|交叉解碼.*bytes"),
        "退卻 BFS 自我修改碼與展開順序，併入 Issue #6。",
    ),
    (
        3,
        re.compile(r"等距|tie.?break|端點|retreatEndpoint|loc_1491B|兩端.*退|廣度優先.*退"),
        "退卻等距端點與 BFS tie-break，併入 Issue #3。",
    ),
    (
        1,
        re.compile(r"並排.*(畫面|兩版)|兩版本.*畫面|pc98golem.*wolong|同狀態.*(兩版|PC-98.*DOS/V)|docs/reference/02|playtest/41"),
        "DOS/V／PC-98 同狀態並排畫面，併入 Issue #1。",
    ),
)


HISTORY = re.compile(r"歷史快照|已被取代|過時工作|舊 worklist|obsolete|superseded", re.I)
EVIDENCE_EXPLICIT = (
    (
        re.compile(r"不擋|不是缺口|已標記的 remake 差異|記為等效差異|只影響|結果相同"),
        "來源明示這是已接受的 remake 差異、等效差異或不擋交付的限制，保留作證據邊界。",
    ),
    (
        re.compile(r"remake 不照抄|remake 照位置實作|remake 先鋪|兩版都沒有讀取端|沒有證據|未定位|未找到"),
        "列出的未知只涉及原版證據或刻意不複製的未定義行為，目前沒有可連結的 remake 玩家工作。",
    ),
)

# Issue body 指定的主要證據位置。命中規則但不在主要位置的列，是同一工作
# 的支援證據，標成 `merge-target`，避免把 805 列誤看成 805 個獨立工作。
CANONICAL_FILES: dict[int, tuple[str, ...]] = {
    1: ("docs/reference/02",),
    2: ("docs/release/", "docs/playtest/109", "docs/playtest/111", "docs/playtest/112", "docs/spec/155"),
    3: ("docs/spec/46",),
    4: ("docs/re/34", "docs/spec/175"),
    5: ("docs/spec/175", "docs/spec/186"),
    6: ("docs/re/11", "docs/re/20", "docs/spec/46"),
    7: ("docs/re/03",),
    8: ("docs/playtest/25", "docs/playtest/26", "docs/re/57", "docs/re/58"),
    10: ("docs/spec/161", "docs/playtest/118", "docs/playtest/119"),
    11: ("docs/spec/195",),
    12: ("docs/playtest/121",),
    13: ("docs/spec/200",),
    14: ("docs/spec/60",),
    15: ("docs/spec/158", "docs/playtest/115"),
    16: ("docs/spec/159", "docs/spec/133"),
    17: ("docs/playtest/117", "docs/spec/117", "docs/re/74"),
    18: ("docs/re/09", "docs/playtest/115", "docs/spec/129"),
    19: ("docs/release/10", "docs/spec/75"),
    20: ("docs/release/13",),
    21: ("docs/mobile/",),
    22: ("docs/mechanics/70-ai", "docs/re/12", "docs/playtest/10"),
    23: ("docs/formats/05", "docs/mechanics/70-ai", "docs/re/07"),
    24: ("docs/mechanics/30", "docs/re/09"),
    25: ("docs/formats/08",),
    26: ("docs/mechanics/70-ai", "docs/playtest/95", "docs/playtest/102"),
    27: ("docs/re/",),
    29: ("docs/spec/132", "docs/spec/164", "docs/spec/169", "docs/spec/170", "docs/spec/171", "docs/spec/180", "docs/spec/183", "docs/spec/184", "docs/spec/194"),
    30: ("docs/playtest/68", "docs/playtest/69", "docs/re/85", "docs/spec/85-march"),
    31: ("docs/spec/115", "docs/spec/121", "docs/spec/196"),
    32: ("docs/spec/85-latin", "docs/spec/87-latin", "docs/playtest/47-latin"),
}


def load_open_questions(repo: str) -> list[dict[str, Any]]:
    """用正式產生器的 collect 規則取得目前列，避免第二套 parser。"""
    path = os.path.join(repo, "tools", "re_open_questions.py")
    spec = importlib.util.spec_from_file_location("re_open_questions", path)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"無法載入 {path}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)

    files: list[str] = []
    for root in module.DOC_ROOTS:
        root_path = os.path.join(repo, root)
        for dirpath, dirnames, filenames in os.walk(root_path):
            dirnames[:] = [d for d in dirnames if d not in ("images", ".git")]
            for filename in sorted(filenames):
                if not filename.endswith(".md"):
                    continue
                path = os.path.join(dirpath, filename)
                rel = os.path.relpath(path, repo).replace(os.sep, "/")
                if rel in module.SKIP or filename == module.INDEX_NAME:
                    continue
                files.append(path)
    files.sort()

    rows: list[dict[str, Any]] = []
    occurrences: collections.Counter[tuple[str, str, str, str]] = collections.Counter()
    for path in files:
        if module.is_superseded(path):
            continue
        rel = os.path.relpath(path, repo).replace(os.sep, "/")
        for item, status, section in module.collect(path, rel)[0]:
            if module.PLATFORM.search(item + " " + status):
                continue
            identity = (rel, item, status, section)
            occurrence = occurrences[identity]
            occurrences[identity] += 1
            row = {
                "file": rel,
                "item": item,
                "status": status,
                "section": section,
                "domain": module.domain_of(rel),
                "verdict": module.verdict_of(item + " " + status),
                "occurrence": occurrence,
            }
            row["id"] = row_id(row)
            rows.append(row)
    return rows


def row_id(row: dict[str, Any]) -> str:
    payload = "\0".join(
        str(row.get(key, ""))
        for key in ("file", "section", "item", "status", "occurrence")
    )
    return "oq-" + hashlib.sha256(payload.encode("utf-8")).hexdigest()[:20]


def classify(row: dict[str, Any]) -> dict[str, Any]:
    haystack = " ".join(
        str(row.get(key, ""))
        for key in ("file", "section", "item", "status")
    )
    for issue, pattern, reason in RULES:
        if pattern.search(haystack):
            return {
                "classification": (
                    "actionable"
                    if any(path in row["file"] for path in CANONICAL_FILES.get(issue, ()))
                    else "merge-target"
                ),
                "github_issue": issue,
                "reason": reason,
            }
    if HISTORY.search(haystack):
        return {
            "classification": "history",
            "github_issue": None,
            "reason": "來源明確標示為歷史／過時資料，保留作證據但不重開工作。",
        }
    for pattern, reason in EVIDENCE_EXPLICIT:
        if pattern.search(haystack):
            return {
                "classification": "evidence-only",
                "github_issue": None,
                "reason": reason,
            }
    return {
        "classification": "evidence-only",
        "github_issue": None,
        "reason": "目前列只保留原版證據邊界，未命中現行玩家路徑、規則或交付 Issue；若影響範圍改變，先建立新 Issue 再調整分類。",
    }


def build_manifest(repo: str) -> dict[str, Any]:
    rows = load_open_questions(repo)
    entries = []
    for row in rows:
        decision = classify(row)
        entries.append({**row, **decision})
    return {
        "schema": SCHEMA,
        "generated_by": "tools/re_triage.py",
        "generated_on": datetime.date.today().isoformat(),
        "source": "tools/re_open_questions.py",
        "source_note": "每列以來源檔、章節、內容、狀態與 occurrence 固定指紋；內容變動必須重新稽核。",
        "classification_policy": {
            "actionable": "目前玩家路徑、規則或交付 gate；必須連到既有 GitHub Issue。",
            "merge-target": "可合併進既有 Issue 的重複／支援證據；必須連到目標 Issue。",
            "evidence-only": "保留原版證據邊界，不形成目前 remake 工作。",
            "history": "已被現況取代的歷史敘述，不重新開工。",
        },
        "rows": entries,
    }


def load_manifest(repo: str) -> dict[str, Any]:
    path = os.path.join(repo, MANIFEST_REL)
    with open(path, encoding="utf-8") as fh:
        return json.load(fh)


def verify(repo: str, manifest: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    if manifest.get("schema") != SCHEMA:
        errors.append(f"schema 不符：{manifest.get('schema')!r}")
    current = {row["id"]: row for row in load_open_questions(repo)}
    entries = manifest.get("rows")
    if not isinstance(entries, list):
        return errors + ["rows 不是陣列"]
    seen: set[str] = set()
    for entry in entries:
        row_id_value = entry.get("id")
        if row_id_value in seen:
            errors.append(f"重複 row id：{row_id_value}")
        seen.add(row_id_value)
        if row_id_value not in current:
            errors.append(f"manifest 有來源不存在或已變動的 row：{row_id_value}")
            continue
        source = current[row_id_value]
        for key in ("file", "section", "item", "status", "occurrence"):
            if entry.get(key) != source.get(key):
                errors.append(f"{row_id_value} 的 {key} 與目前來源不同")
        classification = entry.get("classification")
        if classification not in VALID_CLASSES:
            errors.append(f"{row_id_value} 分類無效：{classification!r}")
        issue = entry.get("github_issue")
        if classification in {"actionable", "merge-target"}:
            if not isinstance(issue, int) or issue < 1:
                errors.append(f"{row_id_value} 是工作分類但沒有 github_issue")
            elif issue == 28:
                errors.append(f"{row_id_value} 不得把分流治理 Issue #28 當成玩家工作目標")
        elif issue is not None:
            errors.append(f"{row_id_value} 非工作分類卻有 github_issue：{issue}")
    missing = sorted(set(current) - seen)
    if missing:
        errors.append(f"有 {len(missing)} 列沒有分流：{', '.join(missing[:8])}")
    if len(entries) != len(current):
        errors.append(f"manifest 列數 {len(entries)} != 目前來源 {len(current)}")
    return errors


def write_manifest(repo: str, manifest: dict[str, Any]) -> None:
    path = os.path.join(repo, MANIFEST_REL)
    with open(path, "w", encoding="utf-8") as fh:
        json.dump(manifest, fh, ensure_ascii=False, indent=2)
        fh.write("\n")


def report(manifest: dict[str, Any]) -> None:
    rows = manifest["rows"]
    counts = collections.Counter(row["classification"] for row in rows)
    issues = collections.Counter(
        row["github_issue"]
        for row in rows
        if row["github_issue"] is not None
    )
    print(f"docs/re/43 分流：{len(rows)} 列")
    print("分類：" + ", ".join(f"{name}={counts.get(name, 0)}" for name in sorted(VALID_CLASSES)))
    print("Issue 對應：" + ", ".join(f"#{n}={count}" for n, count in sorted(issues.items())))
    unlinked = [row for row in rows if row["classification"] == "evidence-only"]
    print(f"未連 Issue 的 evidence-only：{len(unlinked)} 列")
    for row in unlinked[:30]:
        print(f"  {row['id']} {row['file']} :: {row['item']}")
    if len(unlinked) > 30:
        print(f"  …其餘 {len(unlinked) - 30} 列省略")


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=("init", "verify", "report"))
    args = parser.parse_args(argv)
    repo = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    if args.command == "init":
        manifest = build_manifest(repo)
        write_manifest(repo, manifest)
        report(manifest)
        return 0
    try:
        manifest = load_manifest(repo)
    except (OSError, json.JSONDecodeError) as exc:
        print(f"讀取 {MANIFEST_REL} 失敗：{exc}", file=sys.stderr)
        return 1
    errors = verify(repo, manifest)
    if errors:
        print("分流驗證失敗：", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    if args.command == "report":
        report(manifest)
    else:
        print(f"分流驗證通過：{len(manifest['rows'])} 列")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
