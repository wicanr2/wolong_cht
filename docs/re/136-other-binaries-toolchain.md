# 136：其餘七支執行檔的工具鏈判定

**狀態：READY。2026-10-12。八支全判完：3 支 Borland C++ 1991，5 支手寫組語。**

- 日期：2026-10-12。
- 出處：[`reference/04`](../reference/04-first-survey.md) §2.2（三支 banner 既有結論）、
  [`re/135`](135-toolchain-fingerprint.md)（KI 側方法與探針）、本輪八支 IDA 探針。
- 方法：`tools/ida.sh probe` 一次性 DB（`ida-pro-9.4-idapython:py312-v1`，
  非 root，原檔唯讀）＋ `tools/ida_toolchain.py`；banner 用 `grep -a`；
  MZ 標頭只讀解析。收據在 `/tmp/tc-{logo,ynfont,install,str,d7open,d7end,d7over}`，
  不入版控。

## 1. 總表

| 執行檔 | 大小 | banner | IDA 函式／FLIRT lib | cc 欄位 | MZ linker 形 | 判定 |
|---|---|---|---|---|---|---|
| `KI.EXE` | 67,099 | 無 | 739／**0** | 129（預設） | `reloff=0x1E`、`max=1`、`cksum=764a` | 手寫 ASM（`re/135`） |
| `LOGO.EXE` | 57,418 | Borland 1991 | 103／39 | 2 Borland C++ | TLINK（`reloff=0x3E`、`fb506a72`） | **Borland C++ 1991** |
| `YNFONT.EXE` | 60,888 | Borland 1991 | 132／63 | 2 Borland C++ | TLINK（同上形） | **Borland C++ 1991** |
| `INSTALL.EXE` | 66,710 | Borland 1991 | 189／114 | 2 Borland C++ | TLINK（同上形） | **Borland C++ 1991** |
| `STR.EXE` | 916 | 無 | 1／**0** | 129（預設） | TLINK（`reloff=0x3E`、`fb506a72`、`max=65535`） | 手寫 ASM，TLINK 鏈接 |
| `D7OPEN.EXE` | 4,734 | 無 | 56／**0** | 129（預設） | 零 reloc、`min=max=17`、`sp=0100` | 手寫 ASM |
| `D7END.EXE` | 6,677 | 無 | 43／**0** | 129（預設） | 零 reloc、`min=max=17`、`sp=0100` | 手寫 ASM |
| `D7OVER.EXE` | 2,747 | 無 | 31／**0** | 129（預設） | 零 reloc、`min=max=17`、`sp=0100` | 手寫 ASM |

三支 Borland 的判定三重一致：banner 字串＋FLIRT runtime 命名
（`_main`、`_malloc`、`_exit`、`__checknull`…）＋`cc_id=2`。
**confirmed**。

## 2. 逐支證據（新判定的四支）

- `STR.EXE`：全檔只有 1 支函式（`start`，`0x1011E`，`ip:cs=011e:0000`）。
  入口 `mov ax,seg dseg；mov ds,ax；cli；mov ax,3515h；int 21h`
  （取 `15h` 向量），接著比 `es:0／es:2` 是否為 `4F44h／5653h`
 （常駐簽章檢查）。無 banner、無 runtime——ASM 寫的小工具，
  用 TLINK 鏈（故 MZ 是 TLINK 形）。**confirmed**（語言）；用途尚未研究，
  不屬本輪範圍。
- `D7OPEN／D7END／D7OVER`：三支同一家族——零 reloc、`min=max=17`、
  `ss:sp=*:0100`（COM 式堆疊）、`cksum` 非零、入口
  `mov ax,cs；mov ds,ax` 後一串直線 `call`（`sub_…` 全無命名）。
  FLIRT 零命中、無 banner。**confirmed**（手寫 ASM）。

## 3. 判讀邊界

- `cc_id=129／Visual C++` 在五支手寫檔上都是 loader 預設值（見 `re/135` §7），
  不代表 MSVC；FLIRT 零命中才是 signature 面的答案。
- `STR.EXE` 證明「TLINK 形 MZ ≠ Borland C++」：linker 與語言是兩個獨立問題，
  本表兩欄分開判。
- Ghidra 本環境無，八支皆未跑；三重證據（banner／FLIRT／MZ）已收斂，
  不預期翻案。
- 本判定只處理「什麼語言、什麼 linker」；各支的用途與行為仍以既有 `docs/re/`
  筆記為準（`D7*` 三支的行為筆記若缺，那是另一條工作，不在本輪）。
