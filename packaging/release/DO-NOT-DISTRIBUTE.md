# ⛔ 這一批不可外流

版本：`@RELEASE_VERSION@`

這個 `dist-all` 是**內含遊戲檔案的完整版**：四個平台的包裡都有松崗 DOS/V 的
69 個原版檔與倚天點陣字。它是給自己在自己的機器上玩的，**不是交付物**。

## 不可以做的事

- 上傳到 GitHub Release、雲端硬碟、任何公開或半公開的位置
- 傳給別人，包括「只是給朋友試玩」
- 放進任何會被自動備份到公開位置的目錄

## 為什麼

這是文化資產保存專案。公開產出只有**引擎程式碼與譯文校訂紀錄**；
原版執行檔、資料檔、美術、音樂與倚天字型一律不散布，玩家自備。
把原版資產夾帶出去，會讓整個專案的定位站不住。

## 要一份可散布的怎麼辦

```sh
WOLONG_BUNDLE_DATA=0 tools/release_all.sh <YYYYMMDD>
```

出來的 `dist-all` 不含任何原版資產，`release-manifest.json` 的
`distributable` 會是 `true`，而且這個檔案不會存在。

## 怎麼分辨手上這一批是哪一種

| 判準 | 完整版（不可外流） | 可散布 |
|---|---|---|
| 這個檔案 | 存在 | 不存在 |
| `release-manifest.json` 的 `distributable` | `false` | `true` |
| `release-manifest.json` 的 `original_assets_included` | `true` | `false` |
| 包裡有沒有 `gamedata/` | 有 | 沒有 |

⚠ **看檔名分不出來。** 兩種批次的檔名完全一樣，只有大小差約 5 MB。
所以判準要看上面那四項，不要看名字。
