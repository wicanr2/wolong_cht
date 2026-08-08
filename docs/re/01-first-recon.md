# 01 — 首輪偵查：兩版檔案清單與比對

**狀態：READY。兩版檔案清單、執行結構、逐檔比對都完成。**

- 日期：2026-08-07
- 輸入：`workplace/orig/dosv/`（松崗 DOS/V 繁中版，69 檔）、
  `workplace/orig/pc98/`（PC-98 日文原版，69 檔，由 `tools/fdi_extract.py` 從
  五片 `.fdi` 抽出）
- 推論等級：**檔案大小與 SHA-256 是 confirmed**；語意欄一律標等級

## 1. PC-98 磁片結構（confirmed）

五片 `.fdi`，各 1,265,664 B ＝ 4,096 B 檔頭 ＋ 1,261,568 B 裸映像
（1024 B/sector × 1232 sectors，2HD）。檔案系統是 **FAT12**，
BPB 完整，只有 OEM 欄位被換掉：

| 片 | OEM | 內容 |
|---|---|---|
| A | `FGDOS1.0` | 遊戲本體：`KI.EXE`、`MMAP.*`、`KAOGRF`／`ICONGRF`、`SINARIO`／`SAVE`、`TALK.DAT`、`YN*` |
| B | `FGDOS1.0` | 開場：`D7OPEN.EXE`、`OPEN_S1..S6.DAT`、`ROGO.*` |
| C | `YN  1.00` | 戰場與音樂：`BATTLE.MAP/MDL`、`BGM.DAT`、`SOUND.DAT`、`IVENTGRF`、`KYOGRF` |
| D | `YN  1.00` | 結局：`D7END.EXE`、`END_S1..S12.DAT`、`HDINST.EXE` |
| E | `NEC  5.0` | 原廠 MS-DOS 5.0 系統片 |

A／B 片的開機磁區內有：

```
* FGDOS Ver 1.00 Loading module version 1.2 *
Copyright(C)1993
```

**假說**：`FGDOS` 是開發方自製的載入器，**防拷檢查的頭號嫌疑犯**。
1993 這個年份與日文版 1994 發行對得上（開發期）。未驗。

## 2. 兩版逐檔比對（confirmed）

**23 個檔 byte-for-byte 完全相同**，3 個大小相同但內容不同。

| 檔案 | dosv | pc98 | 比對 |
|---|---:|---:|---|
| `BATTLE.DAT` | 8,192 | 8,192 | 大小同、內容異 |
| `BATTLE.MAP` | 877,056 | 877,056 | **完全相同** |
| `BATTLE.MDL` | 194,560 | 194,560 | **完全相同** |
| `BATTLE.SCH` | 115,200 | 115,200 | **完全相同** |
| `BGM.DAT` | 20,826 | 20,790 | 大小不同 |
| `D7END.EXE` | 6,677 | 3,889 | 大小不同 |
| `D7OPEN.EXE` | 4,734 | 4,836 | 大小不同 |
| `D7OVER.EXE` | 2,747 | 1,916 | 大小不同 |
| `ENDBGM.DAT` | 2,218 | 2,294 | 大小不同 |
| `ENDPAL.BRG` | 576 | 576 | **完全相同** |
| `END_S1.DAT` | 40,999 | 40,999 | **完全相同** |
| `END_S10.DAT` | 40,806 | 55,175 | 大小不同 |
| `END_S11.DAT` | 10,637 | 27,413 | 大小不同 |
| `END_S12.DAT` | 81,596 | 81,596 | **完全相同** |
| `END_S2.DAT` | 64,331 | 64,331 | **完全相同** |
| `END_S3.DAT` | 28,391 | 43,292 | 大小不同 |
| `END_S4.DAT` | 35,798 | 53,190 | 大小不同 |
| `END_S5.DAT` | 29,443 | 47,022 | 大小不同 |
| `END_S6.DAT` | 31,196 | 42,560 | 大小不同 |
| `END_S7.DAT` | 25,935 | 40,726 | 大小不同 |
| `END_S8.DAT` | 29,131 | 39,039 | 大小不同 |
| `END_S9.DAT` | 26,600 | 47,806 | 大小不同 |
| `GAMEOVER.DAT` | 16,963 | 26,641 | 大小不同 |
| `GAMEPAL.BRG` | 384 | 384 | **完全相同** |
| `ICONGRF.DAT` | 47,776 | 47,776 | **完全相同** |
| `IVENTGRF.DAT` | 76,032 | 76,032 | **完全相同** |
| `KAOGRF.DAT` | 307,200 | 307,200 | **完全相同** |
| `KI.EXE` | 67,099 | 65,823 | 大小不同 |
| `KYOGRF.DAT` | 69,120 | 69,120 | **完全相同** |
| `MMAP.MAP` | 80,716 | 80,716 | **完全相同** |
| `MMAP.MCH` | 43,058 | 43,058 | **完全相同** |
| `MMAP.MDL` | 32,768 | 32,768 | **完全相同** |
| `OPENBGM.DAT` | 3,424 | 2,990 | 大小不同 |
| `OPENPAL.BRG` | 288 | 288 | **完全相同** |
| `OPEN_S1.DAT` | 51,403 | 51,403 | **完全相同** |
| `OPEN_S2.DAT` | 319,264 | 319,264 | **完全相同** |
| `OPEN_S3.DAT` | 305,649 | 305,649 | **完全相同** |
| `OPEN_S4.DAT` | 212,513 | 212,513 | **完全相同** |
| `OPEN_S5.DAT` | 55,506 | 55,506 | **完全相同** |
| `OPEN_S6.DAT` | 71,150 | 71,150 | **完全相同** |
| `OVERBGM.DAT` | 684 | 888 | 大小不同 |
| `OVERPAL.BRG` | 48 | 48 | **完全相同** |
| `SAVE.DAT` | 88,832 | 88,832 | 大小同、內容異 |
| `SINARIO.DAT` | 88,832 | 88,832 | 大小同、內容異 |
| `SOUND.DAT` | 304 | 372 | 大小不同 |
| `START.BAT` | 26 | 14 | 大小不同 |
| `TALK.DAT` | 34,182 | 45,718 | 大小不同 |
| `YNSOUND.COM` | 3,463 | 4,553 | 大小不同 |

只在 dosv：`.jsdos`、`DISK1`、`DISK2`、`DISK3`、`DISK4`、`END_S13.DAT`、`END_S14.DAT`、`END_S15.DAT`、`INSTALL.EXE`、`INSTALL.MAP`、`INSTALL.SCH`、`LABEL.$$$`、`LOGO.EXE`、`MOUSE.MCH`、`MOUSE.SCH`、`PASS.MAP`、`PASS.SCH`、`PLAY.BAT`、`SHOW.O`、`STR.EXE`、`YNFONT.EXE`、`YNVSHELL.COM`

只在 pc98：`AUTOEXEC.BAT`、`CCP.EXE`、`COMMAND.COM`、`CONFIG.SYS`、`FGDOS.SYS`、`FONTGRF.DAT`、`HDINST.EXE`、`IO.SYS`、`MSDOS.SYS`、`ROGO.EXE`、`ROGO0.DAT`、`ROGO1.DAT`、`STARTUP.CMD`、`STARTUP.SYS`、`YNFONT.COM`、`YNMOUSE.COM`、`YNSHELL.COM`、`臥竜伝_A`、`臥竜伝_B`、`臥竜伝_C`、`臥竜伝_D`

## 3. 從比對讀出來的三件事

### 3.1 松崗版是移植不是重寫（強證據）

四個圖庫、四個調色盤、大地圖三件、戰場三件、整段開場動畫——
全部 byte-for-byte 相同。資料檔原封不動搬過去。

**推論**：格式解一次通吃兩版。從 PC-98 側入手更省事（640×400 固定規格）。

### 3.2 diff 就是「哪些素材燒了文字」的偵測器（強證據）

| 素材 | 比對 | 讀出來的事 |
|---|---|---|
| `OPEN_S1`–`S6` | 六張全相同 | **開場沒有燒字**，文字是疊繪上去的 |
| `END_S1`／`S2`／`S12` | 相同 | 這三張沒有文字 |
| `END_S3`–`END_S11` | 九張全不同 | **圖裡燒了日文，中文版重繪過** |

省掉逐張目視檢查。M1 解圖像格式時，這九張是驗證解碼器的最佳樣本
——同一張圖的兩種語言版本，結構相同、像素不同。

### 3.3 松崗版動過結局段（confirmed 有差異，語意未解）

PC-98 只到 `END_S12.DAT`，松崗版多了 `END_S13`／`S14`／`S15`。
**是加了新過場、還是把原本的長段拆開，未解。**

## 4. 只在單邊存在的檔（confirmed）

見上表。其中值得追的：

| 檔 | 只在 | 意義 |
|---|---|---|
| `FGDOS.SYS`／`STARTUP.SYS`／`CCP.EXE`／`STARTUP.CMD` | pc98 | 自製 DOS，防拷候選 |
| `FONTGRF.DAT`（1,216 B）／`YNFONT.COM`（843 B） | pc98 | 字型。對照 dosv 的 `YNFONT.EXE`（60,888 B）→ **PC-98 靠字型 ROM** |
| `YNMOUSE.COM` | pc98 | 滑鼠驅動。dosv 版把它併進 `KI.EXE`？未驗 |
| `ROGO.EXE`／`ROGO0.DAT`／`ROGO1.DAT` | pc98 | 開機 logo（「ロゴ」羅馬字）。dosv 換成 `LOGO.EXE` |
| `HDINST.EXE` | pc98 | 硬碟安裝 |
| `PASS.MAP`／`PASS.SCH` | dosv | **PC-98 沒有**。關隘資料，移植時新增或改名。未解 |
| `MOUSE.MCH`／`MOUSE.SCH` | dosv | 游標 |
| `SHOW.O`／`INSTALL.*`／`STR.EXE` | dosv | 松崗自寫的安裝與外殼 |

## 5. 反組譯（confirmed）

| 版本 | 檔 | 大小 | SHA-256 | 函式 |
|---|---|---:|---|---:|
| dosv | `KI.EXE` | 67,099 | `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868` | 732 |
| pc98 | `KI.EXE` | 65,823 | `061917f9f3f5c03e29397a9c636d546052128a99b8c8ce31ded0e84cf2a481e8` | 725 |

兩版載入位址都是 `10000h`–`2041Bh`，入口 `1000:0`。
IDA 猜的編譯器是 "Visual C++ (guessed)" —— **那是猜的，不採信**；
`LOGO.EXE`／`YNFONT.EXE`／`INSTALL.EXE` 有明確的 Borland C++ 1991 banner，
但 `KI.EXE` 沒有（`docs/re/02` 待補）。

函式數只差 7 → **同一份原始碼的兩次編譯**（強證據）。

指令：

```sh
tools/ida.sh batch dosv KI.EXE
tools/ida.sh batch pc98 KI.EXE
```

## 6. 下一步

`TALK.DAT` 的索引結構（`CONTEXT.md` §7.0）。
兩版同名同結構、只有語言不同（日 45,718 B ↔ 中 34,182 B），
解開就同時得到中文原文的可寫回抽取與日中對照表。
入口是 `0x000`–`0x7FF` 那 2,048 B，假說是偏移表。
