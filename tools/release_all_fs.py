#!/usr/bin/env python3
"""在 Docker 中整理 dist-all 的可交付發行目錄。

此工具只允許寫入儲存庫根目錄下的 dist-all；建置器以各自專用 Docker image
把 binary 放進 .work/raw，之後再由這裡製作封裝、雜湊與根目錄 manifest。原版
資料、完整 TALK 表與建置中間檔不會進最終交付目錄。
"""
from __future__ import annotations

import datetime
import hashlib
import json
import os
import re
import shutil
import sys
import tarfile
from pathlib import Path


REPO = Path.cwd().resolve()
DIST = Path(
    os.environ.get("WOLONG_DIST_ROOT", str(REPO / "dist-all"))
).resolve()
WORK = DIST / ".work"
# 版本字串。**所有產物的檔名都從這裡長出來**——先前是逐處硬寫日期，
# 於是「換一次版本」等於改十幾個字串，而漏改一處會產出名字對不上內容的檔案。
def _release_version() -> str:
    """版本字串：環境變數優先，否則從 `packages/` 現有的產物反推。

    ⚠ **不要用寫死的日期當預設。** 先前這裡的預設是 `wolong-remake-20260812`，
    而 `tools/py.sh` 不轉送 `WOLONG_*`——所以文件寫的
    `tools/py.sh tools/release_all_fs.py refresh` 這一步，會把
    `finalise` 剛寫對的版本**改回 20260812**，於是
    `dist-all/README.md` 在 `20260827` 的批次裡宣稱自己是 20260812。
    整條流程 exit 0，只有打開那個檔的人看得到。
    """
    env = os.environ.get("WOLONG_RELEASE_VERSION")
    if env:
        return env
    found = sorted(DIST.glob("packages/wolong-remake-linux-amd64-*.tar.gz"))
    if found:
        stamp = found[-1].name.rsplit("-", 1)[-1].split(".")[0]
        return f"wolong-remake-{stamp}"
    raise SystemExit(
        "無法決定版本：請設 WOLONG_RELEASE_VERSION，"
        f"或先讓 {DIST}/packages/ 裡有 wolong-remake-linux-amd64-*.tar.gz"
    )


RELEASE_VERSION = _release_version()
# RELEASE_STAMP 是版本字串尾端的日期，檔名用它。
RELEASE_STAMP = RELEASE_VERSION.rsplit("-", 1)[-1]
# ⭐ **內含遊戲檔案的完整版**（docs/spec/72）。預設開啟——使用者裁定
# 2026-08-22：dist-all 就是四平台完整版。設 `WOLONG_BUNDLE_DATA=0`
# 回到「不含原版資產、可散布」的舊行為。
#
# ⚠ 開著的時候 dist-all **不可外流**：裡面有松崗 DOS/V 的 69 個原版檔
# 與倚天字型。`DO-NOT-DISTRIBUTE.md` 與 manifest 的 `distributable: false`
# 是給人與機器各一份的標記。
BUNDLE_DATA = os.environ.get("WOLONG_BUNDLE_DATA", "1") != "0"
ORIG_SRC = REPO / "workplace" / "orig" / "dosv"
FONT_SRC = REPO / "workplace" / "eten"
# 音檔是原版音樂的合成渲染，屬於原版衍生資產——**只有完整版收**
# （docs/spec/75 §2）。可散布批次一個都不收，與 gamedata 同一條界線。
AUDIO_SRC = REPO / "workplace" / "audio"
# 包內的目錄名，與 cmd/wlgame 的 bundledOrigDir／bundledFontDir 對應。
BUNDLED_ORIG = "gamedata"
BUNDLED_FONT = "fonts"
BUNDLED_AUDIO = "audio"
# ⭐ **發行只帶一支合成片**（使用者裁定 2026-08-30）。
# 主預告、原版實機對照與手機片接成 `wolong-remake-promo.mp4`，
# 配樂全片統一鋪原版曲子（`tools/promo_combined.sh`，docs/promo/combined.md）。
# 那三支與兩支研究用對照片仍留在 `dist/promo/` 當素材，製作紀錄與重錄命令
# 也都還在——**只是不進發行目錄**。
PROMO_FILES = ("wolong-remake-promo.mp4",)


def checked(path: Path) -> Path:
    resolved = path.resolve()
    if resolved != DIST and DIST not in resolved.parents:
        raise SystemExit(f"只允許操作 {DIST}：{resolved}")
    return resolved


def copy_file(src: Path, dst: Path, executable: bool = False) -> None:
    if not src.is_file():
        raise SystemExit(f"缺少封裝輸入：{src}")
    checked(dst)
    dst.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(src, dst)
    if executable:
        dst.chmod(dst.stat().st_mode | 0o111)


def copy_tree(src: Path, dst: Path) -> None:
    if not src.is_dir():
        raise SystemExit(f"缺少封裝輸入目錄：{src}")
    checked(dst)
    shutil.copytree(src, dst)


def template_vars() -> dict[str, str]:
    """模板的代換表。

    ⭐ **界線敘述由 `BUNDLE_DATA` 決定，不是手寫的。** 兩種批次的說明只差
    這幾段；分成兩份模板的話，改一份忘了改另一份，就會出現一份說
    「不含原版資產」的說明躺在含原版資產的包裡——而那正是最危險的錯誤方向。
    """
    if BUNDLE_DATA:
        return {
            "@RELEASE_VERSION@": RELEASE_VERSION,
            "@RELEASE_STAMP@": RELEASE_STAMP,
            "@ROOT_BOUNDARY@": (
                "⛔ **這一批內含原版資產，不可外流**（見 `DO-NOT-DISTRIBUTE.md`）。"
                "四個平台的完整包裡都有松崗 DOS/V 的 69 個原版檔與倚天點陣字，"
                "解開就能玩，不必自備資料。\n\n"
                "要一份可散布的：`WOLONG_BUNDLE_DATA=0 tools/release_all.sh <YYYYMMDD>`。"
            ),
            "@PKG_BOUNDARY@": (
                "⛔ **此封裝內含原版資產，不可外流。** `gamedata/` 是松崗 DOS/V 的"
                "原版資料，`fonts/` 是倚天點陣字；兩者都不得散布。"
            ),
            "@VERIFY_BOUNDARY@": (
                "⛔ **這一批是完整版，`packages/` 裡本來就有原版資產**："
                "`gamedata/` 是松崗 DOS/V 的 69 個原版檔、`fonts/` 是倚天點陣字、"
                "`audio/` 是由使用者自備 `BGM.DAT`／`SOUND.DAT` 算出來的 ogg。"
                "deny-list 掃的是**版控裡的檔案**，不是這些包——"
                "它擋的是「原版資產被 commit 進 repo」，不是「包裡有沒有」。"
            ),
            "@PKG_LAUNCH@": (
                "解開之後直接跑，**不必帶任何旗標**——`gamedata/` 與 `fonts/`"
                "就在執行檔旁邊：\n\n"
                "```text\n./wlgame\n```\n\n"
                "要換成別的資料目錄再明講 `-orig` 與 `-font`；"
                "明講的旗標一律優先。"
            ),
        }
    return {
        "@RELEASE_VERSION@": RELEASE_VERSION,
        "@RELEASE_STAMP@": RELEASE_STAMP,
        "@ROOT_BOUNDARY@": (
            "桌面包與 APK 都不含任何原版執行檔、資料、美術、音樂、字型或完整原版"
            "文字表；玩家必須自行提供合法松崗 DOS/V 資料與字型。"
        ),
        "@PKG_BOUNDARY@": (
            "此封裝只包含 remake 程式與公開的 `translations/corrections.json`。"
            "**不包含**松崗 DOS/V 的執行檔、資料檔、圖像、音樂、完整 TALK 文字表"
            "或倚天字型；請自行準備合法的松崗繁中版資料與中文字型。"
        ),
        "@VERIFY_BOUNDARY@": (
            "`packages/` 已經過封裝內容檢查與 deny-list 掃描；原版執行檔、資料檔、"
            "美術、音樂、字型與 `talk-dosv-corrected.json` 都不應存在。"
        ),
        "@PKG_LAUNCH@": (
            "```text\n./wlgame -orig /path/to/songgang-cht -font /path/to/eten-fonts\n```\n\n"
            "Windows：\n\n"
            "```text\nwlgame.exe -orig C:\\path\\to\\songgang-cht -font C:\\path\\to\\eten-fonts\n```"
        ),
    }


def write_template(template: str, dest: Path) -> None:
    source = REPO / "packaging" / "release" / template
    if not source.is_file():
        raise SystemExit(f"缺少說明模板：{source}")
    checked(dest)
    dest.parent.mkdir(parents=True, exist_ok=True)
    text = source.read_text(encoding="utf-8")
    for key, value in template_vars().items():
        text = text.replace(key, value)
    # ⚠ **沒有被替換掉的變數要當錯誤**。模板用了一個 template_vars() 沒定義的
    # 名字時，字面上的 `@FOO@` 會直接出現在交付檔裡，而流程一路 exit 0——
    # 這種失敗只有人打開那個檔才看得到，而那通常是在使用者手上。
    leftover = re.findall(r"@[A-Z_]{3,}@", text)
    if leftover:
        raise SystemExit(f"{template} 有沒被替換的變數：{sorted(set(leftover))}")
    dest.write_text(text, encoding="utf-8")


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as fp:
        for chunk in iter(lambda: fp.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def write_hashes(root: Path, name: str = "SHA256SUMS.txt", prefix: str = "") -> None:
    checked(root)
    own_checksum = (root / name).resolve()
    rows: list[str] = []
    for path in sorted(p for p in root.rglob("*") if p.is_file() and p.resolve() != own_checksum):
        relative = path.relative_to(root).as_posix()
        rows.append(f"{sha256(path)}  {prefix}{relative}")
    (root / name).write_text("\n".join(rows) + "\n", encoding="utf-8")


def package_dir(source: Path, target: Path) -> None:
    checked(source)
    checked(target)
    if not source.is_dir():
        raise SystemExit(f"缺少 staging 目錄：{source}")
    with tarfile.open(target, "w:gz", format=tarfile.PAX_FORMAT) as archive:
        for path in sorted(source.rglob("*")):
            info = archive.gettarinfo(str(path), arcname=str(path.relative_to(source.parent)))
            info.uid = info.gid = 0
            info.uname = info.gname = ""
            info.mtime = 0
            if info.isfile():
                with path.open("rb") as fp:
                    archive.addfile(info, fp)
            else:
                archive.addfile(info)


def copy_bundled_tree(src: Path, dst: Path) -> None:
    """複製原版資料並**正規化權限**。

    ⚠ `workplace/orig/` 整個是 `chmod a-w`（原版資產唯讀），而 `copytree`
    連目錄的 mode 一起複製——於是 staging 目錄自己也變成不可寫，
    收尾的 `shutil.rmtree(WORK)` 會在 `unlink` 時吃 `PermissionError`
    （unlink 要的是**父目錄**的寫入權，不是檔案的）。

    ⭐ 目錄放寬到 0755、檔案維持 0444：清得掉，而玩家拿到的資料仍然唯讀。
    """
    checked(dst)
    if not src.is_dir():
        raise SystemExit(f"缺少封裝輸入目錄：{src}")
    dst.mkdir(parents=True, exist_ok=True)
    # ⭐ **只收最上層、不收隱藏項**——與 `tools/android_build.sh` 的
    # `cp "$ORIG_DIR"/*` 同一個檔案集。`workplace/orig/dosv` 底下有一個
    # `.jsdos/`（js-dos 的 dosbox.conf），`copytree` 會把它一起帶走，
    # 於是桌面包 73 個檔、APK 69 個。**兩條管線挑不同的集合，
    # 之後任何「兩邊對不起來」的問題都會先被誤判成別的原因。**
    n = 0
    for entry in sorted(src.iterdir()):
        if entry.name.startswith(".") or not entry.is_file():
            continue
        # ⚠ `workplace/audio` 裡 ogg 旁邊還躺著合成中間產物 wav
        # （整包 239 MB，其中 ogg 只有 19 MB）。**只收 ogg**——
        # 照單全收會讓每個桌面包多出兩百多 MB 的中間檔。
        if src == AUDIO_SRC and entry.suffix.lower() != ".ogg":
            continue
        shutil.copy2(entry, dst / entry.name)
        (dst / entry.name).chmod(0o444)
        n += 1
    if n == 0:
        raise SystemExit(f"{src} 裡沒有可收的檔案")
    dst.chmod(0o755)


def bundled_trees() -> list[tuple[Path, str]]:
    """完整包要另外收進去的兩棵樹：原版資料與點陣字。

    ⚠ 只給**完整遊戲**的包用。`linux-arm64-tools` 只有 wlsim／wlshot，
    塞 4.4 MB 進去只是讓包變大（docs/spec/72 §2）。
    """
    if not BUNDLE_DATA:
        return []
    if not (ORIG_SRC / "SINARIO.DAT").is_file():
        raise SystemExit(f"要內嵌遊戲檔案卻找不到 {ORIG_SRC}/SINARIO.DAT")
    trees = [(ORIG_SRC, BUNDLED_ORIG)]
    if FONT_SRC.is_dir():
        trees.append((FONT_SRC, BUNDLED_FONT))
    if AUDIO_SRC.is_dir():
        trees.append((AUDIO_SRC, BUNDLED_AUDIO))
    return trees


def package_stage(
    name: str,
    files: list[tuple[str, str, bool]],
    directories: list[tuple[str, str]] | None = None,
    extra_trees: list[tuple[Path, str]] | None = None,
) -> None:
    stage = WORK / "stage" / name
    if stage.exists():
        shutil.rmtree(stage)
    stage.mkdir(parents=True)
    write_template("README-RELEASE.md", stage / "README-RELEASE.md")
    # ⚠ **授權條款要跟著包走。** `LICENSE` 第 3 條 (a) 要求散布時保留全文，
    # 而拿到 tar.gz 的人看不到儲存庫——包裡沒有這一份，收到的人就不知道
    # 自己被授權了什麼。
    copy_file(REPO / "LICENSE", stage / "LICENSE")
    for src, relative, executable in files:
        copy_file(WORK / "raw" / src, stage / relative, executable)
    for src, relative in directories or []:
        copy_tree(WORK / "raw" / src, stage / relative)
    for src, relative in extra_trees or []:
        copy_bundled_tree(src, stage / relative)
    write_hashes(stage)
    package_dir(stage, DIST / "packages" / f"{name}.tar.gz")


def prepare() -> None:
    if DIST.exists() and any(DIST.iterdir()):
        raise SystemExit(
            f"{DIST} 已存在且非空；為避免覆蓋，請先明確移除該目錄後再建置。"
        )
    DIST.mkdir(parents=True, exist_ok=True)
    WORK.mkdir(parents=True, exist_ok=True)
    for rel in ("packages", "promo", "verification"):
        (DIST / rel).mkdir(parents=True, exist_ok=True)


def rebuild() -> None:
    """只清除 dist-all 的已知發行輸出，供明確重建流程使用。

    不接受未知檔案或目錄，避免把使用者放進 dist-all 的資料誤刪；來源
    ``dist/``、``android/`` 與整個儲存庫都不在此清理範圍內。
    """
    allowed_dirs = {
        "packages",
        "promo",
        "verification",
        "experimental",
        ".work",
    }
    allowed_files = {
        "README.md",
        "SHA256SUMS.txt",
        "release-manifest.json",
        "DO-NOT-DISTRIBUTE.md",
        # ⚠ `finalise()` 會把 `LICENSE` 寫進發行根目錄，所以它也要在這裡列名。
        # 漏列的症狀是**下一次完整重建在第一步就停**（「含未列入發行輸出的
        # 項目，拒絕清理：LICENSE」）——而那一步之前什麼都還沒做，
        # 看起來像流程壞掉而不是清單少一行。授權定案那一輪只跑過 `refresh`，
        # 沒有跑過 `rebuild`，所以這個洞到 2026-08-30 要重打包時才浮出來。
        "LICENSE",
    }
    if DIST.exists():
        unknown = [
            path.name
            for path in DIST.iterdir()
            if (path.is_dir() and path.name not in allowed_dirs)
            or (path.is_file() and path.name not in allowed_files)
        ]
        if unknown:
            raise SystemExit(
                f"{DIST} 含未列入發行輸出的項目，拒絕清理：{', '.join(sorted(unknown))}"
            )
        for name in sorted(allowed_dirs):
            target = DIST / name
            if target.exists():
                shutil.rmtree(target)
        for name in sorted(allowed_files):
            target = DIST / name
            if target.exists():
                target.unlink()
    DIST.mkdir(parents=True, exist_ok=True)


def android_apk_name(apk: Path) -> str:
    """APK 的發行檔名。

    ⭐ Android 與另外三個平台**並列在 `packages/`**（docs/spec/72 §2）。
    先前它在 `experimental/android/`，界線是簽章與驗收——那是驗收狀態，
    寫在 README 就夠，用目錄分級表達會讓人以為它功能不完整。

    日期取**檔案本身的 mtime**，不是發行日期：APK 是另一條管線建的，
    重打發行時未必重建，掛上發行日期會讓檔名宣稱一個沒發生過的建置。
    """
    stamp = datetime.date.fromtimestamp(apk.stat().st_mtime).strftime("%Y%m%d")
    return f"wolong-remake-android-debug-{stamp}.apk"


def sync_android(apk: Path | None = None) -> str:
    """把最新的 debug APK 放進交付目錄，回傳它的發行檔名。

    ⭐ **這一步不能只留在 `stage()`。** APK 是另一條管線（`tools/android_build.sh`）
    建的，重建之後若只跑 `refresh`，交付目錄會留著上一批的檔案——或者更糟，
    一個都沒有，而 manifest 照樣寫著那個不存在的路徑（2026-08-22 踩過）。
    複製完順手清掉別批的 APK，避免同一個目錄裡並存兩個日期。
    """
    if apk is None:
        apk = REPO / "android" / "app" / "build" / "outputs" / "apk" / "debug" / "app-debug.apk"
    name = android_apk_name(apk)
    target = DIST / "packages"
    copy_file(apk, target / name)
    for stale in target.glob("wolong-remake-android-*.apk"):
        if stale.name != name:
            checked(stale).unlink()
    return name


def write_distribution_marker() -> None:
    """散布界線的標記。**兩個方向都要寫**，不能只在內嵌時放一個檔。

    ⚠ 只在內嵌時建立、不內嵌時不刪除的話，一份舊的
    `DO-NOT-DISTRIBUTE.md` 會留在一個其實乾淨的交付目錄裡，
    而下一個人會照著它把一批可散布的東西當成機密收起來——
    反過來也一樣危險。
    """
    marker = DIST / "DO-NOT-DISTRIBUTE.md"
    if BUNDLE_DATA:
        write_template("DO-NOT-DISTRIBUTE.md", marker)
    elif marker.is_file():
        checked(marker).unlink()


def verify_manifest_paths(manifest: dict) -> None:
    """manifest 列到的每一個路徑都必須真的存在。

    ⚠ 少了這一關，「交付目錄缺一個檔」會安靜地通過——manifest 與 SHA256SUMS
    都照樣產出，直到有人照著清單去抓才發現。沉默的成功比失敗難發現
    （CLAUDE.md §7 第 21 條）。
    """
    paths: list[str] = []
    for key in ("desktop_full_packages", "promo_videos"):
        paths.extend(manifest[key])
    for key in ("linux_appimage", "linux_arm64_tools", "android_full_package"):
        paths.append(manifest[key])
    missing = [rel for rel in paths if not (DIST / rel).is_file()]
    if missing:
        raise SystemExit("manifest 指到不存在的檔案：" + "、".join(missing))


def promo_source(name: str) -> Path:
    """找推廣片的來源。

    ⭐ **`dist-all/promo` 也算來源。** 影片是另外一條管線產的，重打發行時
    未必還留在 `dist/promo/`；而 `stage` 寫的是 staging 目錄、`promote` 之後
    才換掉 `dist-all`，所以這時候讀 `dist-all/promo` 是安全的。
    少了這條退路，整批重建會在編完三平台之後才因為缺影片而停掉。
    """
    for base in (REPO / "dist" / "promo", REPO / "dist-all" / "promo"):
        candidate = base / name
        if candidate.is_file():
            return candidate
    raise SystemExit(
        f"找不到推廣片 {name}：dist/promo 與 dist-all/promo 都沒有")


def stage() -> None:
    raw = WORK / "raw"
    correction = REPO / "translations" / "corrections.json"
    if not correction.is_file():
        raise SystemExit(f"缺少公開校訂覆蓋：{correction}")
    for platform in ("linux-amd64", "windows-amd64", "linux-arm64-tools"):
        target = raw / platform / "translations" / "corrections.json"
        copy_file(correction, target)
    for platform in ("darwin-amd64", "darwin-arm64"):
        target = raw / platform / "translations" / "corrections.json"
        copy_file(correction, target)

    package_stage(
        f"wolong-remake-linux-amd64-{RELEASE_STAMP}",
        [
            ("linux-amd64/wlgame", "wlgame", True),
            ("linux-amd64/wlview", "wlview", True),
            ("linux-amd64/wlsim", "wlsim", True),
            ("linux-amd64/wlshot", "wlshot", True),
        ],
        [("linux-amd64/translations", "translations")],
        bundled_trees(),
    )
    package_stage(
        f"wolong-remake-windows-amd64-{RELEASE_STAMP}",
        [
            ("windows-amd64/wlgame.exe", "wlgame.exe", True),
            ("windows-amd64/wlview.exe", "wlview.exe", True),
            ("windows-amd64/wlsim.exe", "wlsim.exe", True),
            ("windows-amd64/wlshot.exe", "wlshot.exe", True),
        ],
        [("windows-amd64/translations", "translations")],
        bundled_trees(),
    )
    package_stage(
        f"wolong-remake-macos-universal-{RELEASE_STAMP}",
        [
            ("darwin-amd64/wlgame", "darwin-amd64/wlgame", True),
            ("darwin-amd64/wlview", "darwin-amd64/wlview", True),
            ("darwin-amd64/wlsim", "darwin-amd64/wlsim", True),
            ("darwin-amd64/wlshot", "darwin-amd64/wlshot", True),
            ("darwin-arm64/wlgame", "darwin-arm64/wlgame", True),
            ("darwin-arm64/wlview", "darwin-arm64/wlview", True),
            ("darwin-arm64/wlsim", "darwin-arm64/wlsim", True),
            ("darwin-arm64/wlshot", "darwin-arm64/wlshot", True),
        ],
        [
            ("darwin-amd64/translations", "darwin-amd64/translations"),
            ("darwin-arm64/translations", "darwin-arm64/translations"),
        ],
        # ⚠ macOS 包是**兩個架構共用一份資料**：資料放包根，
        # 而執行檔在 darwin-<arch>/ 底下——resolveDataDir 的
        # `<exe 目錄>/../<名稱>` 那一條正好對上（docs/spec/72 §3）。
        bundled_trees(),
    )
    package_stage(
        f"wolong-remake-linux-arm64-tools-{RELEASE_STAMP}",
        [
            ("linux-arm64-tools/wlsim", "wlsim", True),
            ("linux-arm64-tools/wlshot", "wlshot", True),
        ],
        [("linux-arm64-tools/translations", "translations")],
    )

    for name in PROMO_FILES:
        copy_file(promo_source(name), DIST / "promo" / name)
    write_template("PROMO-README.md", DIST / "promo" / "README.md")

    sync_android()

    write_template("ROOT-README.md", DIST / "README.md")
    copy_file(REPO / "LICENSE", DIST / "LICENSE")


def appdir() -> None:
    root = WORK / "appdir"
    if root.exists():
        shutil.rmtree(root)
    (root / "usr" / "bin").mkdir(parents=True)
    copy_file(REPO / "packaging" / "appimage" / "AppRun", root / "AppRun", True)
    copy_file(REPO / "packaging" / "appimage" / "wolong-remake.desktop", root / "wolong-remake.desktop")
    copy_file(REPO / "docs" / "images" / "wlgame-dosv-natural-remake-skeleton.png", root / "wolong-remake.png")
    for name in ("wlgame", "wlview", "wlsim", "wlshot"):
        copy_file(WORK / "raw" / "linux-amd64" / name, root / "usr" / "bin" / name, True)
    copy_file(
        REPO / "translations" / "corrections.json",
        root / "usr" / "share" / "wolong-remake" / "translations" / "corrections.json",
    )
    # ⭐ 資料放 `usr/share/wolong-remake/`：執行檔在 `usr/bin/`，
    # resolveDataDir 的 `<exe 目錄>/../share/wolong-remake/<名稱>`
    # 那一條正好對上（docs/spec/72 §3）。**這兩處要一起改**。
    for src, name in bundled_trees():
        copy_bundled_tree(src, root / "usr" / "share" / "wolong-remake" / name)
    write_template("README-RELEASE.md", root / "usr" / "share" / "doc" / "wolong-remake" / "README-RELEASE.md")
    copy_file(REPO / "LICENSE", root / "usr" / "share" / "doc" / "wolong-remake" / "LICENSE")


def record_verification() -> None:
    raw = WORK / "raw"
    rows = [
        "# 建置與檔頭驗證",
        "",
        "Linux amd64 已由 Xvfb GUI smoke 驗收；Windows／macOS 以下僅為交叉建置檔頭驗證。",
        "",
        "## SHA-256",
        "",
    ]
    for rel in (
        "linux-amd64/wlgame",
        "windows-amd64/wlgame.exe",
        "darwin-amd64/wlgame",
        "darwin-arm64/wlgame",
    ):
        path = raw / rel
        rows.append(f"- `{rel}`：`{sha256(path)}`")
    rows.extend(
        [
            "",
            "## ABI 摘要",
            "",
            "- Linux amd64：ELF x86-64。",
            "- Windows amd64：PE32+ x86-64。",
            "- macOS Intel：Mach-O x86-64。",
            "- macOS Apple Silicon：Mach-O arm64。",
            "",
            "Windows／macOS 尚未在目標作業系統完成 GUI runtime smoke，不以檔頭取代實測。",
        ]
    )
    (DIST / "verification" / "BUILD-ABI.md").write_text("\n".join(rows) + "\n", encoding="utf-8")


def finalise() -> None:
    if not (DIST / "packages" / f"wolong-remake-linux-amd64-{RELEASE_STAMP}.AppImage").is_file():
        raise SystemExit("缺少已建立的 AppImage")
    android_name = sync_android()
    record_verification()
    write_template("VERIFICATION.md", DIST / "verification" / "README.md")
    manifest = {
        "release_version": RELEASE_VERSION,
        "desktop_full_packages": [
            f"packages/wolong-remake-linux-amd64-{RELEASE_STAMP}.tar.gz",
            f"packages/wolong-remake-windows-amd64-{RELEASE_STAMP}.tar.gz",
            f"packages/wolong-remake-macos-universal-{RELEASE_STAMP}.tar.gz",
        ],
        "linux_appimage": f"packages/wolong-remake-linux-amd64-{RELEASE_STAMP}.AppImage",
        "linux_arm64_tools": f"packages/wolong-remake-linux-arm64-tools-{RELEASE_STAMP}.tar.gz",
        "promo_videos": [f"promo/{name}" for name in PROMO_FILES],
        "android_full_package": f"packages/{android_name}",
        "original_assets_included": BUNDLE_DATA,
        "distributable": not BUNDLE_DATA,
        "complete_original_talk_table_included": False,
        "native_gui_smoke": {"linux_amd64": True, "windows_amd64": False, "macos": False},
    }
    verify_manifest_paths(manifest)
    write_distribution_marker()
    (DIST / "release-manifest.json").write_text(
        json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8"
    )
    if WORK.exists():
        shutil.rmtree(WORK)
    write_hashes(DIST, prefix="dist-all/")


def promote() -> None:
    """在 staging 完整驗證後，以可回復的目錄交換更新 dist-all。"""
    staging = (REPO / "dist-all.staging").resolve()
    live = (REPO / "dist-all").resolve()
    backup = (REPO / ".dist-all.previous").resolve()
    if DIST != staging:
        raise SystemExit("promote 必須以 WOLONG_DIST_ROOT=dist-all.staging 執行")
    if not (staging / "release-manifest.json").is_file():
        raise SystemExit("staging 缺少 release-manifest.json，拒絕替換舊 dist-all")
    if backup.exists():
        raise SystemExit(f"已有未清理的交換備份：{backup}")
    if live.exists():
        live.rename(backup)
    try:
        staging.rename(live)
    except Exception:
        if backup.exists() and not live.exists():
            backup.rename(live)
        raise
    if backup.exists():
        shutil.rmtree(backup)


def refresh() -> None:
    """清除中斷建置遺留物後重算最終交付雜湊。

    只可在已完成的發行根目錄上執行；不重建、不讀取原版資料，也不修改包檔。
    """
    manifest_path = DIST / "release-manifest.json"
    if not manifest_path.is_file():
        raise SystemExit("缺少完成的 release-manifest.json；不能 refresh 未完成的發行目錄")
    if WORK.exists():
        shutil.rmtree(WORK)
    write_template("ROOT-README.md", DIST / "README.md")
    write_template("PROMO-README.md", DIST / "promo" / "README.md")
    copy_file(REPO / "LICENSE", DIST / "LICENSE")
    # ⭐ APK 要重新同步再回填 manifest：Android 是另一條管線，
    # 重建之後 `refresh` 是唯一會跑到的一步。
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    manifest.pop("android_experimental", None)
    manifest["android_full_package"] = f"packages/{sync_android()}"
    # 推廣片清單由 PROMO_FILES 決定，不沿用舊 manifest——換片之後
    # 沿用舊清單會讓 verify_manifest_paths 指向一個已經不存在的檔名。
    manifest["promo_videos"] = [f"promo/{name}" for name in PROMO_FILES]
    manifest["original_assets_included"] = BUNDLE_DATA
    manifest["distributable"] = not BUNDLE_DATA
    verify_manifest_paths(manifest)
    write_distribution_marker()
    manifest_path.write_text(
        json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8"
    )
    write_template("VERIFICATION.md", DIST / "verification" / "README.md")
    write_hashes(DIST, prefix="dist-all/")


def main() -> int:
    if len(sys.argv) != 2:
        raise SystemExit("用法：release_all_fs.py prepare|stage|appdir|finalise|refresh")
    commands = {
        "rebuild": rebuild,
        "prepare": prepare,
        "stage": stage,
        "appdir": appdir,
        "finalise": finalise,
        "promote": promote,
        "refresh": refresh,
    }
    try:
        commands[sys.argv[1]]()
    except KeyError:
        raise SystemExit("用法：release_all_fs.py prepare|stage|appdir|finalise|refresh") from None
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
