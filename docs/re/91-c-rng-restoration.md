# 91：第一個 C 函式：sub_1ECE0 的語意還原與對照

**狀態：局部語意對照通過。可讀 C 取數函式、呼叫介面適配層與 Go 資料規則均有收據。**

- 日期：2026-10-08
- 範圍：松崗 DOS/V `KI.EXE` 的 `sub_1ECE0`；不外推播種或整個玩家流程
- EXE SHA-256：`fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`
- 原始指令：IDA 9.4 線性 `0x1ECE0`–`0x1ECFC`，檔案 `0xEEE0`–`0xEEFC`，13 條／28 bytes
- 指令 SHA-256：`f56e5a1d08956bc9c68476dd5245b85a7e0169d52757b5e476c32b28e1911a9d`
- 契約：[`spec/201`](../spec/201-c-rng-function.md)；資料語意：[`10`](10-rng.md)
- 組語基準：[`90`](90-assembly-reconstruction.md)

## 1. 還原的 C

[`rng.c`](../../tools/c_recovery/rng.c) 的資料函式為：

```c
uint8_t sub_1ECE0(KiRngState *state) {
    uint8_t value = (uint8_t)(state->table[state->s] + state->c);
    state->c = (uint8_t)(state->c + 0x89);
    state->s = value;
    return value;
}
```

正式研究來源另有四個受編譯參數控制的負對照，預設 `KI_MUTATION=0`。
欄位沿用已證實的原始定位：`byte_1ECFC` 為 c、`byte_1ECFD` 為 s、`0x1ECFE` 為表。
這個 C 介面是本次建立的 typed representation，不代表原作者的名稱或來源寫法。

適配層 `sub_1ECE0_abi` 處理原始 SS:SP 推入、POP 與近返回的記憶體副作用，
並保留 AX 高位元組與運算旗標。輸入及副作用契約在 spec/201，沒有改動正式 Go 規則。

## 2. 原版 oracle 與控制條件

dosgolem 直接載入原始 EXE，再以 `IDAIn` 保留原始 CS，執行原始 13 條指令。
`IDA()` 的正規化段只適合記憶體讀取；拿它作 CS 會改變這支函式的絕對定址。

| 條件 | 固定方式 |
|---|---|
| 表 | `T[i] = (mul × i + add) mod 256`；四組 `(1,0)`、`(255,255)`、`(197,13)`、`(137,39)` |
| c／s | 每張表枚舉全部 65,536 種組合；執行前設定 |
| 暫存器 | 由案例索引與固定常數產生，實際入口值納入比較 |
| 旗標 | IF／TF 為 0；DF 與六個運算旗標依案例索引變化 |
| 一般堆疊 | SS=`0x4000`、入口 SP=`0x8000` |
| 邊界 SP | `0`、`2`、`0xFFFC`、`0xFFFE` |
| 重疊 SP | SS=原始 CS，SP=`0xED00`、`0xED04`；堆疊不覆蓋函式指令 |
| 停止點 | 近返回後的受控返回位址，每次必須剛好執行 13 條指令 |

初始狀態由公開的 `CallNear` 與記憶體 API 擺入。初始化旗標／SS／SP 的五條測試指令
放在獨立 RAM 區，原始 28 bytes 保持相同。此方法是局部常式注入，不是正常玩家路徑。
表是明示受控輸入，沒有測試原版的 BIOS 時鐘播種，也沒有重擲或挑選 seed。

## 3. 結果

GCC `-O0` 與 `-O2` 各自獨立建置；兩側的輸入與枚舉順序相同。

| 比較 | 每個最佳化版本的結果 |
|---|---|
| 原版與 C：四張表 | 262,144 組相同 |
| 原版與 C：堆疊邊界 | 1,024 組相同 |
| 原版與 C：堆疊重疊 | 512 組相同 |
| 合計 | 263,680 組，14 個暫存器／段／FLAGS 欄位及指定記憶體相同 |
| Go `rng.Next` | 非重疊的 263,168 組，回傳值與 258-byte 狀態相同 |
| 全記憶體核對 | 258 次完整 1 MB 相同 |

原版寫入監看限制在 c、s 與堆疊六個 bytes。兩個最佳化版本的原版／C 摘要皆為
`7de79c400d7fba4d87eb6c34d44f13e38b0f3c93876eb5a63142a0dbed20a74f`。
原版／Go 的資料狀態摘要皆為
`2635412d88976d173c826ab8624b4e342fd25dd39be47013213eb7bc80ed1a70`。

## 4. 負對照與容易漏掉的副作用

| 故意改錯 | 首次拒絕 |
|---|---|
| 計數器加 `0x88` | 第 1 組 |
| AX 高位元組歸零 | 第 1 組 |
| 旗標取自回傳值的加法 | 第 1 組 |
| BX 直接還原本機變數的舊值 | 第 1,281 組，第一個堆疊重疊案例 |

BX 的 PUSH／POP 保存的是記憶體，堆疊若與 c／s 重疊，取數寫回會改掉稍後 POP 的值。
因此適配層必須真的讀回堆疊。此條是受控 alias 邊界，沒有把它當成玩家正常會遇到的狀態。

## 5. 工具身分與來源

| 工具 | 本次身分 |
|---|---|
| dosgolem | `a9714ebdab2ad6b529f81225472680f2b11f2842`；讀取的 oracle／internal 原始碼雜湊在執行前後相同 |
| dosgolem 工作樹 | `/home/anr2/cht/dosgolem/cmd/probe/main.go` 有既有修改；本次沒有修改它，且該 command 不在此工具的 import 路徑 |
| Go | `go1.26.7 linux/amd64` |
| GCC | Debian `12.2.0-14+deb12u1`，原生 C11 |
| 編譯參數 | `-O0`／`-O2 -std=c11 -D_GNU_SOURCE`；負對照另加 `KI_MUTATION=1..4` |

C 原始碼 SHA-256：`1e6500c692d652c07cd31b99afa5693275bdab91d0f2d2c9ce9c8ef1188cb0fd`。
C header SHA-256：`71d82996afa8df2387e9899f8a78872ad35d617ed1276b31aaf12c68dd66c0f0`。
兩個最佳化版本的結果 JSON 雜湊相同：
`aca7a14bbeb93593802bd64fe922fb05211b1bd475de83fb0ef6fdb479be6e93`。

## 6. 未解範圍

| 項目 | 邊界 |
|---|---|
| C 機器碼匹配 | 本次沒有達成或聲稱；此為原生 C 的語意還原 |
| 輸入完整性 | 四張受控置換表，不是全部可能的表；其他 CS、奇數 SP、IF／TF 開啟不在本次範圍 |
| 原版播種與取數時機 | 沒有新增驗證；局部取數通過不證明整段亂數流同步 |
| GUI 與玩家流程 | 沒有本輪正常玩家路徑或時鐘對拍收據 |

組語還原台帳保存建立時的導航快照。C 的目前完成狀態查
[`c-recovery-status.json`](c-recovery-status.json)，避免把舊的 `not-started` 快照當成現況。
現行玩家工作仍查 [GitHub Issues](https://github.com/wicanr2/wolong_cht/issues)。

## 7. 重跑與產物

```sh
tools/c_recovery_rng.sh
```

對照器在 [`c_recovery_rng.go`](../../tools/c_recovery_rng.go)，收據驗證器在
[`c_recovery_verify.py`](../../tools/c_recovery_verify.py)。原始 EXE、dosgolem 與專案來源唯讀掛載；
Go workspace 與 runtime 只存在容器內，沒有在 go.mod 加不可取得的版本。

本機產物根目錄為 `workplace/matching-decompilation/c-rng/`；`results/O0.json`、
`results/O2.json`、四份 `mutant-*.json` 保存局部收據，`verification.json` 保存跨來源身分驗證。
完整原始碼／工具版本與 SHA-256 清單均隨本機產物保存。
