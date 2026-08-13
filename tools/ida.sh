#!/usr/bin/env bash
# IDA Pro 9.4 headless 包裝。image 來源 /home/anr2/ida_94_official/dist
#
#   tools/ida.sh batch <版本> <執行檔>       產 .i64 + .asm
#   tools/ida.sh raw   <版本> <idat 參數…>   直接下 idat 指令（會改寫 .i64）
#   tools/ida.sh script <版本> <腳本.idc>    對 .i64 的唯讀副本跑腳本
#
# <版本> = dosv | pc98，對應 workplace/ida/<版本>/
#
# script 模式存在的理由：idat 一開啟 .i64 就會改寫它的雜湊，
# 而筆記要靠雜湊標明「這個結論是在哪一份資料庫上驗的」。
# 對副本跑，原始 .i64 的身分才穩定。輸出落在 workplace/ida/<版本>/census/。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE=ida-pro-9.4-ver2
MODE="$1"; VER="$2"; shift 2
WORK="$ROOT/workplace/ida/$VER"
mkdir -p "$WORK"

# IDA 在這個 image 裡以非 root 執行會靜默失敗（無輸出、非 0 離開），
# 所以容器內仍是 root。代價是產出歸 root，之後在容器外改檔／清理都要 sudo。
# 對策是跑完再由容器自己 chown 回呼叫者——不需要主機端 sudo。
UIDGID="$(id -u):$(id -g)"
fix_owner() {
  docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --entrypoint chown \
    -v "$1:/out" "$IMAGE" -R "$UIDGID" /out
}

run() {
  docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --network none --memory 4g --pids-limit 256 \
    -v "$WORK:/work" -v "$ROOT/tools:/work/tools:ro" -w /work \
    "$IMAGE" "$@"
}

case "$MODE" in
  batch)
    BIN="$1"
    cp -f "$ROOT/workplace/orig/$VER/$BIN" "$WORK/$BIN"
    chmod u+w "$WORK/$BIN"
    sha256sum "$WORK/$BIN"
    run idat -A -B "$BIN"
    fix_owner "$WORK"
    ;;
  raw) run "$@"; fix_owner "$WORK" ;;
  script)
    IDC="$1"; DB="${2:-KI.EXE.i64}"
    SCRATCH="$WORK/census"
    mkdir -p "$SCRATCH"
    echo "來源資料庫："; sha256sum "$WORK/$DB"
    cp -f "$WORK/$DB" "$SCRATCH/$DB"
    chmod u+w "$SCRATCH/$DB"
    docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
      --network none --memory 4g --pids-limit 256 \
      -v "$SCRATCH:/work" -v "$ROOT/tools:/tools:ro" -w /work \
      "$IMAGE" idat -A "-S/tools/$(basename "$IDC")" "$DB"
    fix_owner "$SCRATCH"
    echo "原始資料庫雜湊（應與上方相同）："; sha256sum "$WORK/$DB"
    ;;
  *) echo "用法: ida.sh {batch|raw|script} {dosv|pc98} …" >&2; exit 2 ;;
esac
