# 128：正常行軍在道路相遇的野戰碰撞

**狀態：兩軍96個等待檢查點、野戰入場比對及50項正常滑鼠檢查通過。**

- 日期：2026-11-03。
- 原版：松崗DOS/V `KI.EXE`，SHA-256 `fffeba985231cda4d636e93d10f598470b1f691d00275e4aa38e285893d43868`。
- 執行器：dosgolem `f79eab6db8c64a524a06535e53e5d32b0bbcf7c8`，唯讀掛載。
- remake：`bd09fd40e344efee9d01e286a70f2a2303890caf` 加工作樹；沿用 [spec272](../spec/272-world-tick-battle-resume.md) 的程式，本輪未改遊戲Go來源。

- 程式：`workplace/scratch/world-resume-20261103/wlgame`，SHA-256 `af84322202bdd3317741791f59f5e35b572a12f6397b788c80cb91d64fcbf507`。
- 工具鏈：`san1-go:1.26.7-ebiten2.9.9-r1`。
- 既有規則：[spec175](../spec/175-encounter-standoff-countdown.md)、[spec186](../spec/186-standoff-corps-before-city.md)；快照映射沿用 [playtest126](126-player-march-collision-snapshot.md)。

## 正常原版場景

從首頁第一章曹操開局、正常編成37、向袁術12進言敵對並選「我國較有利」。
37正常攻下汝南100並走進城中心後，再由軍團選單與地圖下令往新息105 `(221,144)`。
原版軍團3已由新息往汝南行進，雙方在道路上相遇，沒有執行siege、跳入常式、搬軍團或注入勝負。

首次完整阻擋拍為196/4/9 2時子刻0：3停 `(221,141)`、37停 `(221,140)`。
兩軍均使用連結 `0x10B0`，目標分別100及105。其後95次完整更新加起點，共96個等待檢查點。
再一次更新時由3觸發野戰，待戰時鐘仍停在4/9 12時子刻5。

| 原版停點 | 指令數與狀態 |
|---|---|
| 首次完整阻擋拍 | 391,818,113；4/9 2時子刻0 |
| `0x14A7B` 野戰入口 | 399,364,343；SI=`2300`、DI=`2B80`，3攻37守 |
| `0x11B5A` 戰術入口 | 399,471,569；SI=`2300`、BX=`2B80`，戰場198、旗標`0x80`，不轉180度 |
| 下一戰術節拍 | 403,494,228；既有AfterTicks由0推到1，640×400戰術圖已目視 |

位址為IDA Pro9.4 database linear，base `0x10000`。
runtime軍團表從 `DS:2240` 起，軍團3為 `2300`、37為 `2B80`，與SAVE位移分開記錄。

## 同起點比較

gate快照保存完整22,208-byte區塊、258-byte raw RNG、98,304格佔用圖與事件節流。
raw RNG SHA-256為 `cea8741734034d452b89eee401ef8780b929940ba91cdc0055b9a5575aadb2a4`。
Go在執行前還原這份亂數狀態，接入實際道路、地形、AI及戰術，不重播種或重擲。
未解區塊 `0x52A0..0x52BF` 仍沿用原始SAVE模板，不稱完整執行期記憶體。

| 範圍 | 結果 |
|---|---|
| 等待96點，分別指定3及37 | 全域、兩軍必要欄位、玩家資金、完整raw RNG、全圖佔用低7位及事件節流相同 |
| 完整64-byte軍團記錄 | 3全部相同；37有86個byte差異，均僅為繪圖bit4；動畫`+21`無差異 |
| 完整世界區塊 | 10點逐byte相同；總差異為MCH物件區2,560個byte及上述繪圖bit4的86個byte，其餘0 |
| 野戰入場 | 對原版11B5A比較，子刻5、批次游標0與角色相同；排除表頭`3B..7F`、MCH及繪圖bit4後語意差異0 |
| 入場佔用表示 | 原版暫扣來攻3所在 `(221,141)` 一格，原版0／Go1；其餘全圖相同，沿用spec272的表示界線 |
| 比較器反例 | 只改Go第1拍raw RNG的一個bit，兩軍比較均於第1拍失敗，實際比較器退出1 |
| 重播穩定性 | 成功原版與前次已取得的510份狀態資料逐檔相同；比較輸入的480份檔案亦與成功原版逐檔相同 |

原版14A7B入口尚未選場及清雙方對峙，不能直接拿Go Pending與它比較。
本片用11B5A入口比較世界狀態。Go Pending已做戰術Spawn，該點RNG不與原版初始化前混比；
等待期仍比較完整raw RNG。入場raw區塊有41個byte差異，不宣稱整塊精確一致。

## 地形取樣限制

本場守方37的實際位置是 `(221,140)`；目前Provider仍從傳入的據點105中心 `(221,144)` 取樣。
兩處的十字五格地形類型恰好全為0，玩家37朝向3，均由配對 `(0,0)` 得到 `192+6=198`、不翻轉。
原版與Go本場編號相同，不能據此宣稱取樣座標正確。

| 取樣位置 | 上／左／右／中心／下的原始圖塊 |
|---|---|
| 守方實際位置 `(221,140)` | `C9 00 00 C9 C9` |
| Provider據點位置 `(221,144)` | `D8 F2 DB CD F5` |

這些圖塊經既有Terrain表都回0。`field-selection.json` 保存正式世界重播到待戰後的Provider編號、翻轉與兩處取樣。
取樣座標、同據點不同遭遇的快取與水域重複抽骰接線仍需獨立驗證，沿用
[Issue #31](https://github.com/wicanr2/wolong_cht/issues/31) 的範圍，不在本片猜補或擴大通過聲明。
地形五格是守方位置的十字；舊散文把它寫成下方一片，漏算了 `sub dx,18h` 上移一列，
原始指令與座標解釋見 [re05](../re/05-battle-selection.md) §2及既有battlefield規則註記。

## 正常滑鼠驗收

當時的 `tools/desktop_mouse_smoke.py` 已加入field路徑，與實際受測的可丟棄來源逐byte相同，
SHA-256為 `4e84a4dfe6b9cd05eb649f17e55fd37d87153029bc06912275f84945a20faf6b`。
固定seed1從首頁開局，經宣戰、編成、攻取汝南，待3走上真實道路後再下令105。
GUI的seed1與原版gate的raw RNG不同，正常操作接線不冒稱同一場戰鬥的精確狀態。

50項檢查通過，包含正常開局、設定、編成、進言宣戰、地圖出兵、行軍中原生存讀檔、道路相遇、
3攻37守且mode1、滑鼠守陣／退卻、返回戰略及戰後存讀檔。
GUI入口為4/9 17時、戰後讀檔為4/9 18時；汝南仍歸0、新息仍歸12、首都仍為許昌82。
`field-guard.png` 與戰後讀檔圖已目視，沒有以不同亂數、節拍及鏡頭的兩圖宣稱像素一致。
本場聚焦野戰所需路徑，語系、全選單邊界等共用抽測沿用既有收據。

遊戲359份來源與spec272封存逐檔相同，既有39套件冷測及67／68項攻城GUI未重跑，未冒稱為本輪新測試。

## 收據與失敗分類

本機根目錄為 `workplace/scratch/field-collision-20261103/`，原始素材及衍生快照不加入Git。
`original/` 為成功原版，`comparison-input/` 只替480份原始資料加上比較器所需檔名，mapping逐檔保存來源。
`remake-final/`、`final-comparison.json`、兩軍比較報告、`whole-block-summary.json` 及field-selection保存有效結果，
`gui-final/` 保存50項正常滑鼠收據。
入口為 `replay-original.sh`、`replay.sh field 空輸出目錄`、sampler及compare；來源與雜湊由verification及manifest索引。

- `original-attempt1` 只有觀測器截短後的未用符號編譯錯誤，未執行有效遊戲。
- `original-attempt2` 已取得道路野戰，右鍵按住期間跨過11B5A後，僅靠IP停止會漏停；
  改用既有OnCall旗標及AfterTicks契約後，同一路線乾淨重跑成功，未延長預算。
- `gui-attempt1` 在最高速截圖期間錯過離城窗口，實際成為汝南守城。
- `gui-attempt2` 我方靠近時，敵軍3尚未通過新息，之後改為留守105，實際成為37攻3守的攻城。
  trace能證實目標／Ordered改105；spec194的中繼受威脅留守分支是強推論，trace未含當時Threatened欄位。
  修正重播條件為敵軍3實際走上 `0x10B0` 道路後才下令，未更換seed、重擲、注入軍團或改遊戲規則。
- `gui-attempt3` 已進入正常野戰，但在frame148停住；`gui-attempt4` 同路線以SIGQUIT取到render thread
  停於OpenGL DrawElements，更新執行緒等待graphicscommand.Flush，未落在尋路或單挑。
  原堆疊及失敗圖保留，不把先前靜態尋路候選當成已確認的產品缺陷。
- `gui/` 的softpipe完整抽測在25秒內未觀測到行軍，`gui-softpipe-focused/` 因相同繪圖速度限制中止。
  沒有延長原本的操作或400秒容器預算。

沿用同一image，GLX探針確認Mesa22.3.6的llvmpipe LLVM15.0.6及softpipe均可啟動。
最終GUI以 `GALLIUM_DRIVER=llvmpipe LP_NUM_THREADS=0` 關閉工作執行緒後通過同一路徑，遊戲程式未改。
參數契約見 [Mesa官方環境變數文件](https://docs.mesa3d.org/envvars.html)，查閱日2026-11-03。
這是本機驗收環境修正，不外推其他平台或所有繪圖驅動的可用性。

## 未解範圍

完整戰鬥結果、不同地形／水域、野戰玩家攻方的原版矩陣、MCH動畫及完整戰役不由本片代證。
原版仍報未建模 `int61 AH0A/AL03`，不作音訊或硬體時序聲明。Issue #30與#31維持OPEN。

2026-11-04 後續：[spec273](../spec/273-encounter-battlefield-context.md) 已修正接戰座標、參戰朝向、
單次選場與讀檔接線。新程式重跑本片96點及50項正常滑鼠野戰通過，
收據另存 `workplace/scratch/field-selection-20261104/`，不覆寫本片原始紀錄。
原版原始快照沿用本片的固定輸入，沒有重生或替換；新局部terrain矩陣另行標示，仍不代證正常碼頭野戰。
