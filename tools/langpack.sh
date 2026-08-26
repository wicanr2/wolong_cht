#!/usr/bin/env bash
# 在帶網路的一次性容器裡跑語系包產生器（需要 pip 裝 OpenCC／pypinyin）。
#
#   tools/langpack.sh                      # 簡體語系包（tools/langpack.py）
#   tools/langpack.sh tools/namepack.py en # 英文人名羅馬化
#
# 與 tools/py.sh 分開的理由：py.sh 刻意 --network none 且不裝第三方套件；
# 語系包生成是離線一次性工作，需要網路拉 opencc，所以獨立一支、
# 一樣 --rm 用完即丟（docs/spec/84 §2）。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

docker run --rm --log-opt max-size=10m --log-opt max-file=3 \
    --memory 768m --cpus 1 --pids-limit 96 \
    -v "$REPO_ROOT":/work -w /work \
    python:3.12-slim \
    sh -c "pip install -q opencc-python-reimplemented pypinyin && python ${*:-tools/langpack.py zh-hans}"
