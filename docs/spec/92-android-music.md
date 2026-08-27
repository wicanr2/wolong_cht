# 92 — Android 也要有原版的音樂

**狀態：CONFORMED（2026-08-27 實作並實跑驗過）。**

- 日期：2026-08-27
- 出處：使用者裁定「原版的曲子 ogg 請放入完整版裡面」。
  桌面四個包（Linux tar／AppImage／Windows／macOS）**本來就各有 32 個 ogg**，
  缺的只有 Android。
- 推論等級：曲目與場景對應是既有的 confirmed 結論
  （[`../re/58`](../re/58-bgm-scene-mapping.md)）；本規格只處理「怎麼放進 APK
  以及怎麼播」。

## 1. 現況盤點

| 包 | ogg | 說明 |
|---|---:|---|
| `linux-amd64.tar.gz` | 32 | `audio/` 目錄，`bundled_trees()` 收的 |
| `linux-amd64.AppImage` | 32 | `usr/share/wolong-remake/audio/` |
| `windows-amd64.tar.gz` | 32 | 同 Linux tar |
| `macos-universal.tar.gz` | 32 | 同上 |
| **`android-debug.apk`** | **0** | ⛔ 缺 |
| `linux-arm64-tools.tar.gz` | 0 | 這是工具包不是遊戲包，不用 |

APK 內嵌的是 `assets/gamedata/orig/*`（74 檔，含 `BGM.DAT`／`SOUND.DAT`）
與 `assets/gamedata/eten/*`。**原始資料在裡面，但沒有算好的 ogg**——
而手機端也**完全沒有接音訊**：`internal/ui/phone` 與 `mobile/wolong`
一行 `sound` 都沒有。

⭐ 所以「把 ogg 放進去」不能只放檔案：**放了不播等於在包裡塞 19 MB 死重量**。

## 2. 三段要做

### 2.1 打包：APK 多一個 assets 子目錄

`tools/android_build.sh` 在 `WOLONG_BUNDLE_DATA=1` 時多複製一份
`workplace/audio/*.ogg` → `assets/gamedata/audio/`。

⚠ **只收 `.ogg`**：`workplace/audio` 裡 ogg 旁邊躺著合成中間產物 wav，
整包 239 MB 而 ogg 只有 19 MB（`release_all_fs.py` 早就踩過同一條）。

Java 側 `ImportActivity` 的 `unpackBundled()` 目前只解 `orig` 與 `eten` 兩個
子目錄，加上 `audio`。**那個陣列就是唯一的清單**——漏了它，assets 進得去、
解不出來，而畫面上什麼都看不出來。

### 2.2 選曲規則抽成一份

桌面的 `musicTrack()`／`battleMusic()` 綁在 `cmd/wlgame` 的 `game` 上。
手機要同樣的行為，**不能抄第二份**（`CLAUDE.md` §7 第 6 條：
反組譯筆記重複只是浪費，程式碼重複會產生行為差異）。

新增 `internal/rules/bgm`，只吃純值：

```go
type Scene struct {
    Launcher, Ending, GameOver bool
    Message  bool          // 事件訊息或進言對話開著
    Battle   *Battle       // nil ＝ 不在戰術畫面
    Month    int           // 1–12；0 ＝ 沒有世界
}
type Battle struct { Field int; PlayerAttacks bool }
func Track(s Scene) string
```

規則照 [`../re/58`](../re/58-bgm-scene-mapping.md) 原樣搬，**不改行為**：
啟動殼層曲 0、結局 `endbgm-0`、勝負已定 `overbgm-0`、
戰術依戰場編號 7／8／9／10、事件與對話曲 6、其餘照月份挑四季曲 2–5。

### 2.3 手機播它

`mobile/wolong` 開 `sound.Open(DataRoot()/audio)`，每幀用
`bgm.Track(...)` 更新曲目——與桌面 `updateMusic()` 同一個形狀。
`internal/ui/phone.Session` 要多幾個查詢方法把場景狀態交出來
（戰術中／訊息開著／月份），**規則不進 phone**。

手機的系統面板第一頁加一列「音效」，點一下開關（桌面的系統選單早就有那一列）。

## 3. ⚠ 這是原版衍生物

ogg 是 `tools/bgm2ogg.sh` 從使用者自備的 `BGM.DAT` 算出來的
（[`29`](29-audio.md)）。**它與原版資料同一條界線**：只進「完整版」
（`WOLONG_BUNDLE_DATA=1`，⛔ 不可外流），可散布版一個都不帶。
`tools/denylist.py` 掃的是 git 追蹤的檔案，`workplace/` 早就 gitignore，
這一條不變。

## 4. 改動

| 檔 | 內容 |
|---|---|
| `tools/android_build.sh` | 內嵌時多複製 `audio/*.ogg` |
| `android/…/ImportActivity.java` | 解壓的子目錄陣列加 `audio` |
| `internal/rules/bgm`（新）| `Track(Scene)` |
| `cmd/wlgame/main.go` | `musicTrack()` 改成組 `bgm.Scene` 再呼叫 |
| `mobile/wolong/game.go` | 開 `sound.Bank`、每幀更新曲目 |
| `internal/ui/phone` | 場景查詢方法；系統面板加「音效」列 |

## 5. 驗證（2026-08-27）

| 項目 | 結果 |
|---|---|
| `TestTrackCoversEveryScene` | 16 條場景逐條對 `re/58` 的表（含「戰術蓋過訊息」與「結局蓋過勝負」兩個排序）|
| `TestSeasonMusicMatchesOriginalTable` 等三支 | 規則搬到 `internal/rules/bgm` 之後照樣過——**那三支是拿 `KI.EXE` 當 oracle 的**，跟著規則一起搬 |
| `TestMusicSceneFollowsSessionState` | Session 的狀態正確翻進 `bgm.Scene`；面板開著時放曲 6 |
| `TestSoundRowSeparatesMissingFromOff` | 「未接入」與「關」分得開；沒有音庫時點它不會變成「關」 |
| APK | **32 個 ogg**（14 首音樂 ＋ 18 個音效），內嵌檔數 74 → **106** |
| 模擬器 | `getFilesDir()/audio` 解出 **32 個 ogg**（`android_smoke.sh` 多一道檢查）；同一輪的指紋 `dc4773497b12206e`／`563d356df606873a`／`645434f9e3bd7991` 與桌面同幀逐字相同 ⇒ 選曲規則搬進 `internal/rules/bgm` 沒有動到規則層 |

⚠ **模擬器聽不到聲音不代表沒播**：`android_smoke.sh` 起的是 `-no-audio`
的模擬器。判準是**檔案解出來了**，不是「有沒有聲音」。

## 6. 未解

| 缺口 | 下手點 |
|---|---|
| 手機沒有實機聽過 | 沒有裝置；模擬器是 `-no-audio` 起的 |
| 桌面「沒有音效裝置就掛」那一條 | [`../release/08`](../release/08-full-20260826.md) §5 的老問題，手機端要確認 Ebiten 在 Android 上不會踩同一個 |
