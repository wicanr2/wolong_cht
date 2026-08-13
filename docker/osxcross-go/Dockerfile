# macOS 發行交叉編譯工具鏈。
#
# `u5cht/osxcross:latest` 提供既有的 osxcross SDK／clang，但內含 Go 1.24，
# 無法建置本專案（go.mod 要求 Go 1.25）。此 image 是其明確替代 revision：
# 不掛載主機 Go，也不改變 osxcross SDK，只從已鎖在本機的 Go 1.25 base 取 toolchain。
# 建置：
#   docker build --network=none -t wolong-osxcross-go:20260811-event10-r4 \
#     -f docker/osxcross-go/Dockerfile docker/osxcross-go
# 使用專案既有的 `wolong-go:20260809` Go 1.25.12。它的 GOROOT 已依其 baseline
# GOEXPERIMENT 篩過互斥 runtime source，直接用通用 Go image 交叉重建 darwin 標準庫
# 會得到重複宣告，不能作為可重現 release toolchain。
FROM wolong-go:20260809 AS go125

FROM u5cht/osxcross:latest

# 基底 image 也有舊 GOROOT；COPY 不會移除其中在新版已不存在的互斥 runtime source，
# 所以必須先精確移除這一個工具鏈目錄，再放入 project image 的完整 Go 1.25 GOROOT。
RUN rm -rf /usr/local/go

COPY --from=go125 /usr/local/go /usr/local/go

ENV PATH="/usr/local/go/bin:${PATH}" \
    GOTOOLCHAIN=local
