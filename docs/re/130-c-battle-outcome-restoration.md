# 130：自動戰鬥、退卻與據點易主的 C

**狀態：CONFORMED。21函式，O0／O2各8,364組完整狀態與二十錯版通過。**

- 日期：2026-10-10
- 輸入：松崗DOS/V KI.EXE；SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 工具：IDA Pro 9.4；位址均為IDA database linear，段基址10000。
- Database SHA-256：`6efd571f0d44cf6d52e3d1eeecdc41079a082cfbf15a798a8fb1b4dae656d282`。
- Probe SHA-256：`b20f688289b06ca5384f3ac810cd4b9b93cc5ddb82d5be4d12556701a8988ce9`。
- 探針：[ida_outcome_probe.py](../../tools/ida_outcome_probe.py)；規格：[spec/250](../spec/250-c-battle-outcome.md)。

## 已證實原始閉包

| 原始定位 | 範圍 |
|---|---|
| `sub_15130`，0x15130 | 雙方戰力、各加8、32-bit比例與100上限；回傳AL及兩個壞滅bits，保留同值分支及16-byte SS區域。 |
| `sub_15285`／`sub_152D7` | 六槽兵種係數、士氣右移、word乘法、武將raw能力及XLAT，保留原截斷與byte回繞。 |
| `sub_151B3` | 城市與雙方六槽傷亡、最低大將槽、總兵力與兩種士氣換算；RNG順序及byte DIV分母照原式。 |
| `sub_1474A`，0x1474A | 重算、士氣／大將槽歸零、己方據點停留、退一站、300門檻與8／10階段。 |
| `sub_14CF3`，0x14CF3 | owner交換、內政官卸任、兩方據點數、遷都／滅亡、舊軍團轉向、城市標記及小地圖。 |
| `sub_14DF0`／`sub_14FCE` | 遷都失敗時舊勢力失效；玩家勢力由0x14FE5呼叫既有`sub_11CB1`恢復外層frame退出，原caller不續行。 |
| 其餘13個helper | 內政官／外交官、武將／軍團、鄰接遮罩及小地圖依賴；原始邊界及bytes在固定探針。 |

共21個命名函式、796指令／1,721bytes，直接呼叫邊界已封閉。
`sub_14A7B`、`sub_14ADE`與完整軍團更新`sub_125A3`只作導航；本輪不代證完整戰術戰鬥或正常玩家流程。
既有規則見[戰鬥](../mechanics/30-combat.md)、[戰後退卻](../spec/46-post-battle-retreat.md)及[易主轉向](../spec/47-city-fall-corps-redirect.md)。

## 小地圖拷貝方向勘誤

已證實，原`sub_19656`將DS設為minimap資料、ES設為A0C8，`sub_196CF`（0x196CF）的MOVSB由DS資料還原到ES的VRAM；前面的ES讀取只供VGA latch。
本輪`sub_195C9`（0x195C9）反向把DS設為A0C8、ES設為minimap，`sub_1963F`（0x1963F）的MOVSW保存VRAM至資料。
上一輪語意索引對`sub_196CF`的文字方向寫反，原始bytes與已驗證C並未改動；原註記可由Git基準a8161c8回查。

## 驗證入口

[Go收據](c-outcome-go-verification.json)保存本輪39套件冷測及隔離自測，原生C完整狀態與二十錯版驗證已完成。
[指令覆蓋](c-outcome-code.json)保存21函式796指令／1,721bytes的獨立重組及既有來源對映。

## 獨立模型勘誤

原版首差位於世界區塊`2289`，即軍團`2280+9`。模型給3，原版給5。
既有`sub_16FD2`的IDA線性位址`0x17012–0x1701A`先保存原值，再左移兩次並相加，已證實是byte乘5。
初版獨立模型誤寫乘3；修正模型，保持原始C、案例與seed。

`sub_195C9`在座標0,0產生目的位移`FFCF`，第三個word落在`FFFF`。
`sub_1963F`的IDA線性位址`0x19648`是MOVSW；dosgolem的[固定版本CPU來源](https://github.com/wicanr2/dosgolem/blob/a9714ebdab2ad6b529f81225472680f2b11f2842/internal/cpu/cpu.go#L250)之`read16`／`write16`各自對兩個byte的位移做16-bit回繞。
獨立模型須在分段位移相加後再轉線性位址，不能直接以線性位址加1。

## 驗證結果

- O0／O2各8,364組完整RAM／四plane／DAC／暫存器／FLAGS／SS／IN／OUT／API一致；群組為battle 5,376、retreat 384、capture 1,376、fall 96、diplomat 32、history 8、army 320、neighbors 512、minimap 256、alert 4。
- 原版先經獨立資料／算式／固定raw RNG／平台輸入模型，再與各自raw初態的C比較；不從原版已執行Snapshot初始化C。
- 每版144組真正恢復outer SS／SP及RET位置，另有8,220組正常返回。玩家滅亡時舊caller不續行；淡出沿既有來源核對272個DAC setter。
- 二十個實編譯錯版均由狀態差異拒絕；原六槽、士氣、XLAT、比例、退卻、owner／官員／鄰接與滅亡順序保留。
- [C來源](../../tools/c_recovery/outcome.c)、[獨立模型](../../tools/c_recovery_outcome_data.go)及[收據](c-outcome-verification.json)可回查；重跑`bash tools/c_recovery_outcome.sh`。
- 796指令／1,721bytes獨立重組一致，沒有新增組語指令；整檔25,078指令仍重建相同67,099-byte EXE。
- 較早戰後退卻、易主轉向與小地圖文字均回鏈至本頁；Go與正常玩家完成條件不變。

## C交戰／戰術補證

`sub_14A7B`的原始`0x14A7B`已納入交戰／戰術C閉包。O0／O2各748組完整狀態、全部原始入口及二十個實編譯錯版已核對，見[re/131](131-c-engagement-tactical-restoration.md)與[spec/251](../spec/251-c-engagement-tactical.md)。較早證據保留其原範圍；完整軍團／主排程、正常玩家長程與C機器碼匹配仍未完成。

## 未解範圍

| 項目 | 限制 |
|---|---|
| 完整交戰caller、戰術／軍團更新、主排程與正常玩家長程 | 本輪閉包不代證 |
| 非法DIV／記錄、原硬體時間與原C機器碼 | 不由局部比較推廣 |
