#!/usr/bin/env bash
# IDA Pro 9.4 headless 包裝。image 來源 /home/anr2/ida_94_official/dist
#
#   tools/ida.sh batch <版本> <執行檔>       產 .i64 + .asm
#   tools/ida.sh raw   <版本> <idat 參數…>   直接下 idat 指令（會改寫 .i64）
#   tools/ida.sh script <版本> <腳本.idc|.py>  對 .i64 的唯讀副本跑腳本
#
# <版本> = dosv | pc98，對應 workplace/ida/<版本>/
#
# ⭐ 腳本副檔名是 .py 就自動換到修好 IDAPython 的 image
# （ida-pro-9.4-idapython:py312-v1）。**基底 image 跑 IDAPython 是零輸出的
# 靜默失敗，而且 exit code 不可信**——所以這裡用副檔名決定 image，
# 不讓呼叫端記得要換。優先寫 IDAPython：有 idautils／ida_funcs／ida_xref，
# 不必跟 IDC 缺哪個內建函式搏鬥。
#
# script 模式存在的理由：idat 一開啟 .i64 就會改寫它的雜湊，
# 而筆記要靠雜湊標明「這個結論是在哪一份資料庫上驗的」。
# 對副本跑，原始 .i64 的身分才穩定。輸出落在 workplace/ida/<版本>/census/。
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE=ida-pro-9.4-ver2
IMAGE_PY=ida-pro-9.4-idapython:py312-v1
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
    # .py 走修好 IDAPython 的 image；.idc 走基底。副檔名決定，不靠紀律。
    IMG="$IMAGE"
    case "$IDC" in
      *.py)
        IMG="$IMAGE_PY"
        docker image inspect "$IMG" >/dev/null 2>&1 || {
          echo "[ida.sh] 找不到 $IMG。IDAPython 在基底 image 上是靜默失敗，" >&2
          echo "         不要退回基底重試——先建 image，Dockerfile 見" >&2
          echo "         ~/.claude/knowledge-base/retro/assets/ida-pro-9.4-idapython.Dockerfile" >&2
          exit 3; }
        ;;
    esac
    # 中斷的執行會在 SCRATCH 留下未打包的資料庫（.id0/.id1/.nam/.til）。
    # IDA 看到它們就不肯開同名的 .i64，訊息是
    # 「Failed to initialize IDA as library (error code 4)」——**看起來像授權或
    # image 壞掉，其實只是殘檔**。每次開跑前清掉，只清 SCRATCH 裡的副本，
    # 原始 .i64 在上一層目錄不受影響。
    rm -f "$SCRATCH/${DB%.i64}".id0 "$SCRATCH/${DB%.i64}".id1 \
          "$SCRATCH/${DB%.i64}".nam "$SCRATCH/${DB%.i64}".til
    echo "來源資料庫："; sha256sum "$WORK/$DB"
    echo "image：$IMG"
    cp -f "$WORK/$DB" "$SCRATCH/$DB"
    chmod u+w "$SCRATCH/$DB"
    docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
      --network none --memory 4g --pids-limit 256 \
      -v "$SCRATCH:/work" -v "$ROOT/tools:/tools:ro" -w /work \
      "$IMG" idat -A "-S/tools/$(basename "$IDC")" "$DB"
    fix_owner "$SCRATCH"
    echo "原始資料庫雜湊（應與上方相同）："; sha256sum "$WORK/$DB"
    ;;
  *) echo "用法: ida.sh {batch|raw|script} {dosv|pc98} …" >&2; exit 2 ;;
esac
