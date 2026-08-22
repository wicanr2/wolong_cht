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
import shutil
import sys
import tarfile
from pathlib import Path


REPO = Path.cwd().resolve()
DIST = Path(
    os.environ.get("WOLONG_DIST_ROOT", str(REPO / "dist-all"))
).resolve()
WORK = DIST / ".work"
# 版本字串。由 `WOLONG_RELEASE_VERSION` 決定，預設是 20260812 那一次的三平台交付。
# **所有產物的檔名都從這裡長出來**——先前是逐處硬寫 `20260812`，
# 於是「換一次版本」等於改十幾個字串，而漏改一處會產出名字對不上內容的檔案。
RELEASE_VERSION = os.environ.get("WOLONG_RELEASE_VERSION", "wolong-remake-20260812")
# RELEASE_STAMP 是版本字串尾端的日期，檔名用它。
RELEASE_STAMP = RELEASE_VERSION.rsplit("-", 1)[-1]
PROMO_FILES = (
    "wolong-remake-trailer.mp4",
    "wolong-remake-classic-revival.mp4",
    "wolong-remake-yt-comparison.mp4",
    "wolong-remake-dosv-live-comparison.mp4",
    "wolong-remake-android.mp4",
)


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


def write_template(template: str, dest: Path) -> None:
    source = REPO / "packaging" / "release" / template
    if not source.is_file():
        raise SystemExit(f"缺少說明模板：{source}")
    checked(dest)
    dest.parent.mkdir(parents=True, exist_ok=True)
    dest.write_text(
        source.read_text(encoding="utf-8").replace("@RELEASE_VERSION@", RELEASE_VERSION),
        encoding="utf-8",
    )


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


def package_stage(
    name: str,
    files: list[tuple[str, str, bool]],
    directories: list[tuple[str, str]] | None = None,
) -> None:
    stage = WORK / "stage" / name
    if stage.exists():
        shutil.rmtree(stage)
    stage.mkdir(parents=True)
    write_template("README-RELEASE.md", stage / "README-RELEASE.md")
    for src, relative, executable in files:
        copy_file(WORK / "raw" / src, stage / relative, executable)
    for src, relative in directories or []:
        copy_tree(WORK / "raw" / src, stage / relative)
    write_hashes(stage)
    package_dir(stage, DIST / "packages" / f"{name}.tar.gz")


def prepare() -> None:
    if DIST.exists() and any(DIST.iterdir()):
        raise SystemExit(
            f"{DIST} 已存在且非空；為避免覆蓋，請先明確移除該目錄後再建置。"
        )
    DIST.mkdir(parents=True, exist_ok=True)
    WORK.mkdir(parents=True, exist_ok=True)
    for rel in ("packages", "promo", "verification", "experimental/android"):
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
    target = DIST / "experimental" / "android"
    copy_file(apk, target / name)
    for stale in target.glob("wolong-remake-android-debug-*.apk"):
        if stale.name != name:
            checked(stale).unlink()
    write_template("ANDROID-EXPERIMENTAL.md", target / "README.md")
    return name


def verify_manifest_paths(manifest: dict) -> None:
    """manifest 列到的每一個路徑都必須真的存在。

    ⚠ 少了這一關，「交付目錄缺一個檔」會安靜地通過——manifest 與 SHA256SUMS
    都照樣產出，直到有人照著清單去抓才發現。沉默的成功比失敗難發現
    （CLAUDE.md §7 第 21 條）。
    """
    paths: list[str] = []
    for key in ("desktop_full_packages", "promo_videos"):
        paths.extend(manifest[key])
    for key in ("linux_appimage", "linux_arm64_tools", "android_experimental"):
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
    write_template("README-RELEASE.md", root / "usr" / "share" / "doc" / "wolong-remake" / "README-RELEASE.md")


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
        "android_experimental": f"experimental/android/{android_name}",
        "original_assets_included": False,
        "complete_original_talk_table_included": False,
        "native_gui_smoke": {"linux_amd64": True, "windows_amd64": False, "macos": False},
    }
    verify_manifest_paths(manifest)
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
    # ⭐ APK 要重新同步再回填 manifest：Android 是另一條管線，
    # 重建之後 `refresh` 是唯一會跑到的一步。
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    manifest["android_experimental"] = f"experimental/android/{sync_android()}"
    verify_manifest_paths(manifest)
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
