# macOS 發行交叉編譯工具鏈。
#
# `u5cht/osxcross:latest` 提供 osxcross 的 SDK 與 clang，但內含 Go 1.24，
# 建不了本專案（`go.mod` 要求 Go 1.25）。所以只換掉 GOROOT：
# 不掛載主機 Go，也不動 osxcross SDK。
#
# 建置：
#   docker build -t wolong-osxcross-go:20260827 \
#     -f docker/osxcross-go/Dockerfile docker/osxcross-go
#
# ⭐ **兩層 FROM 都要是公開映像。** 這一份曾經寫 `FROM wolong-go:20260809`
# ——指向另一顆只存在於本機的映像，那顆被清掉之後這一份就建不起來，
# 而發行流程要到跑 macOS 那一段才會發現（`docs/release/09` §4）。
# Go 版本改成釘死的公開 tag，鏈上每一環都能從零取得。
FROM golang:1.25.12-bookworm AS go125

FROM u5cht/osxcross:latest

# 基底 image 也有舊 GOROOT；COPY 不會移除其中在新版已不存在的互斥 runtime
# source，所以必須先精確移除這一個工具鏈目錄，再整份放入 Go 1.25 的 GOROOT。
RUN rm -rf /usr/local/go

COPY --from=go125 /usr/local/go /usr/local/go

ENV PATH="/usr/local/go/bin:${PATH}" \
    GOTOOLCHAIN=local
