# 275：HD城池地圖重繪、工具列、頭像、戰術AI輔助與武將編輯

**狀態：DRAFT。** 日期：2026-11-05。保存已確認需求與待決邊界，尚未進入正式實作。

HD美術方向已由使用者確認。DRAFT保留給正式素材接線、世界鏡頭、AI與編輯器等未完成部分。

- 日期：2026-11-04

## 使用者已確認

- 最新裁定重新定義HD：以原版為依據，先重繪城池與地圖。城池位置、等級與道路仍沿原版資料，
  先做本機畫風草稿，不再將超解析放大作為本輪HD美術方向。本機HTML不部署雲端的要求沿用；
  使用者已釐清只限制網頁部署，明確允許把素材上傳讓雲端手繪。必要原圖／參考圖可交內建圖片工具，
  成果及HTML留本機，不更新Sites或公開原版素材。此範圍不再重問。
- 使用者已否決A青綠彩繪與B古紙淡彩，認為兩者失去原版特色。這兩案不再待選。
  原版位置與路網一致不足以代表美術忠實；本輪回到實際城池圖與據點景觀圖，
  以原圖輪廓、視角、比例、色塊與既有線條作修繪約束，不再以通用城樓及地形筆觸替代。
- 許昌本機微量修繪也被指出仍糊。使用者進一步明確要求：基於原版手繪、不自創，
  提升HD清晰度與原版特色，強化色彩，品質如rich2 HD。沿用原圖造型不代表固定柔化放大像素；
  高解析線條與既有材質需要重繪，技術保護區檢查不能證明達到此要求。
- 使用者已對原生重繪回覆「ok 就是這樣」「這就是我要的」「hd 方向確定」。
  許昌地圖r1／景觀r2確立後續HD美術基準：原版造型、視角與配置，清楚線條／材質，強化原色。
  這項畫風不再待選，正式引擎接線與其他素材仍須依原圖逐項製作和驗證。
- 新增滎陽、舞陽、虎牢關、汝南各兩張後，使用者回覆「ok 同意」，五城十張皆已同意。
  原圖與SHA保持不變；使用者接著表示「期待大地圖hd」，後續優先推進大地圖試繪及預覽。
- 優先完成原版操作及視窗，以抽驗收斂，之後HD，再之後手機版。
- 全程滑鼠操作，只有姓名等文字欄位保留輸入；原版上下推動反白沿用既有裁定。
- 工具列在遊戲畫布外常駐，參考 `~/cht/rich2`。
- 使用者選A：原版／HD與縮放放在工具列。這次選的是入口配置，
  比較稿已明列點陣背景與字型為佔位，沒有選定HD美術。
- 頭像也要HD。
- HD版要參考rich2支援拉遠／拉近。地圖視野與介面大小分離，工具列提供減號、
  倍率重設、加號，支援滑鼠滾輪；目前比較頁的圖片縮放不代表正式地圖鏡頭完成。
- 另加工具列AI選項，作用在戰術畫面，以統帥加成改善小兵到城門及架梯兵效率。
  使用者已選擇同時加成找路效率、單兵走速與架梯動作時間。
  我方與敵方都套用，各依自己的統帥加成。
  使用者已選最高2倍：統帥1為1倍、15為2倍，中間線性增加，僅開啟AI輔助時套用。
- 另加武將能力編輯，包含隱藏的委任、陸戰、海戰、攻城欄位。

初始A為九鈕。新增AI輔助後為十項；若武將編輯也置於工具列則為十一項。
新增入口的間距與分組仍需比較稿，不能把九鈕草稿直接壓縮後稱為定版。
正式遊戲目前仍是 [spec266](266-mouse-only-desktop-toolbar.md) 的七鈕，
本輪原版捲軸修正另見 [spec274](274-original-list-scrollbar-drag.md)。

## 已查證的工程入口

| 項目 | 目前事實與來源 |
|---|---|
| 工具列參考 | `/home/anr2/cht/rich2/internal/ui/themetoolbar/toolbar.go`、`/home/anr2/cht/rich2/internal/frontend/theme_toolbar.go`、`toolbar_controls.go`；常駐，與遊戲畫布分離。rich2保留原版鍵盤相容，不複製到本作 |
| 城池原稿 | MMAP地圖城池與KYOGRF據點景觀是不同素材。`City.KindHigh` 選KYOGRF頁碼，不能只按 `City.Kind` 套用通用建築；來源沿 [re50](../re/50-city-info-window.md)。許昌82對應第6頁、96×96，實際圖庫15頁 |
| 拉遠／拉近參考 | `/home/anr2/cht/rich2/internal/frontend/native_zoom.go`：1–5倍、每次0.2、預設3.4、160ms過渡；`toolbar_controls.go` 接工具列與地圖內滾輪。世界鏡頭跟隨焦點，介面不縮放，zoom只屬呈現且不寫存檔。這些數字先作prototype參考，非本作正式倍率已定版 |
| 頭像 | `internal/state/state.go` 的 `General.Portrait` 是KAOGRF頁碼，150張；武將ID與頁碼不同，例如武將16曹操使用第50頁。既有 `docs/playtest/07` 是入口；HD替換須沿頁碼映射，保留缺圖fallback |
| 找路 | `internal/rules/tactical/path.go` 已有雙平面成本搜尋；`pathqueue.go` 每拍最多兩筆，走不動才排隊，見 [spec120](120-pathfind-request-queue.md)、[spec134](134-pathfind-wave-order.md) |
| 架梯 | [戰鬥機制](../mechanics/30-combat.md) §3.6區分固定架梯位置、就位與等待時間；改善排隊／就位和縮短動作時間是不同規則，不能默默合併 |
| 原始武將資料 | [格式08](../formats/08-sinario-save.md) §3、`internal/state.General`；松崗DOS/V KI.EXE身分沿spec274，原始資料由使用者自備 |

## 編輯欄位草案

| 使用者名稱 | typed欄位 | 原始武將記錄偏移 | 已知有效範圍／意義 |
|---|---|---|---|
| 武術 | Martial | +11h | 1–15 |
| 統帥／原版統率 | Command | +12h | 1–15 |
| 政治 | Politics | +13h | 1–15 |
| 委任 | Tactic | +16h | 0–7，八種戰術腳本類型，不把編號大小宣稱為強弱 |
| 攻城 | Aptitude[0] | +0Eh高半位元組 | 0–10 |
| 陸戰／原版野戰 | Aptitude[1] | +0Fh高半位元組 | 0–10 |
| 海戰／原版水戰 | Aptitude[2] | +10h高半位元組 | 0–10 |

以上偏移及型別沿既有已證實資料，不新增一個假的「委任能力百分比」。
適性寫回須保留不屬於編輯範圍的原始位元；若評價等衍生值必須更新，另列依賴及驗收。
正常存檔與原生存檔兩條路都須驗證，不能只改畫面上的數字。

## 決策樹與下一步

1. 已確認：工具列外置、原版／HD與縮放常駐、頭像HD、戰術AI選項、武將能力編輯。
2. 已選：AI加速包含找路、單兵走速與架梯動作時間，會改攻城節奏；排除只改善排隊／就位。
   雙方各依自己的統帥加成，排除僅我方受益。
3. 加成曲線已選 `倍率 = 1 + (統帥 - 1) / 14`，有效統帥1–15，最高2倍。
   下一步確認開關時機與原有移動／架梯計時接線，固定seed量測到門及架梯時間，
   保留相同初始狀態的原模式正對照。不能以總戰鬥時間變短代替找到原因。
4. HD美術先處理許昌一城，分開使用原版MMAP城池圖與KYOGRF第6頁據點景觀。
   使用者已否決青綠彩繪／古紙淡彩，停止擴展通用SVG造型；原圖直接作修繪輸入。
   使用者已確認地圖r1／景觀r2的HD方向，後續依此基準擴展其他城池與區域，逐份參考其原圖。
   已確認的畫風不重開選擇；資料及輪廓檢查不代替後續素材的品質檢查。
   頭像HD需求保留，排在這批地圖與城池之後。原版美術對拍仍用原模式。
5. 編輯器先做一名武將的滑鼠增減／類型選擇、預覽與取消草稿；確認套用時機與保存邊界，
   再走資料→typed欄位→規則→兩條存檔→讀回→正常UI驗收。

建議原版規則與可選AI輔助分開驗證；開關預設及切換時機尚未定案。
沒有把本頁新功能接進生產路徑，也未更改存檔格式、戰鬥數值或公開素材政策。

## 比較稿與權利

本機比較稿：`workplace/scratch/ui-windows-20261104/hd-toolbar-prototype.html`，
原始A為 `hd-toolbar-A.svg`，B為 `hd-toolbar-B.svg`；生成程式同目錄。
採用A的配置選擇記錄於本頁。原版衍生畫面與自備字型形成的草稿留在本機，
原圖及雜湊由本輪manifest保存，不公開散布。

速度曲線概念稿：`workplace/scratch/ui-windows-20261104/ai-speed-prototype.html`，
比較最高1.5／2／3倍，統帥1為1倍、15為上限的線性提案。
這只是無障礙路線及架梯進度的概念動畫，未接遊戲、不作實戰效能或原版證據。
使用者已選2倍上限及上述線性曲線；其餘兩條不採用。這是新增設計，不稱為原版證據。

## 原稿參考的原生HD重繪（方向已確認）

rich2的真實方法已核對於 `/home/anr2/cht/rich2/docs/hd-concepts/terrain-art-generation.json` 與spec086§36：
內建imagegen直接參考原版圖磚重繪，保留輪廓、數量、方向、邊緣交界與原投影，
另提高材質可讀性及原色色彩。本作依使用者這次的獨立明示授權採用同類工具流程，未沿用rich2的授權。

現行 [五座城池原版與HD重繪](../../workplace/hd-redraw-local-20261104/faithful/index.html)
是單檔離線HTML，嵌入十張原圖與十張1254×1254原生HD。採內建imagegen，未切CLI或外部模型服務。
圖像檢視用於查原稿及生成稿，HTML建置與瀏覽器驗證在斷網Docker完成。

- 許昌基準景觀原稿為KYOGRF第6頁96×96，地圖原稿為周圍7×7格112×112；兩者分別提供最近鄰8倍圖，
  只讓工具讀清原形，不作插值HD成品。rich2 atlas只作清晰度、色彩與材質完成度參考。
- 景觀第1稿保留正面大殿、屋頂、金色中央飾物與三組石階，但把原版抖色畫成棋盤格，
  並新增金色小牌框。第2稿針對這兩項修正，保留清楚屋瓦與石階，已選為景觀HD基準。
- 地圖第1稿依原稿方城、四門、建築群、中央金色建築與道路方向重繪，換成清楚屋瓦、石材和草木。
  繪製幾何不逐像素一致，牆體厚度與材質細節有重繪詮釋；使用者已確認此HD方向。
- 兩張選稿直接複製工具原生PNG bytes，未放大、濾鏡銳化或合成回舊柔化保護區；
  景觀SHA `e019af55ca00e54808a90652ed81d95fc42ec43cb7481db6cdc573c6f596591e`，
  地圖SHA `aae7498646ce4872dea4a16fb99990d51a436b143328d635d9246468037d7131`。
- 頁面提供五城切換及各城的地圖／景觀切換，保留原版／HD／並排／拉桿、工具列與滾輪縮放、拖曳同步及1:1細節。
  一般視野不再把低解析修繪圖放大；1:1按1254 CSS像素檢視，不冒稱實體螢幕像素相同。
  十一項離線抽驗通過，JS錯誤0、外部請求0，含二十張嵌入圖bytes、五城來源及頁碼、模式、縮放、1:1、搬移單檔及手機寬度。
  原稿及實際比較頁已目視檢查，這些檢查不取代使用者美術驗收。

來源、完整提示詞與選稿入口：
[重繪請求及提示詞清單](../../workplace/hd-redraw-local-20261104/faithful/hd-request/request.json)、
[明示授權](../../workplace/hd-redraw-local-20261104/faithful/hd-request/authorization.json)、
[輸出清單](../../workplace/hd-redraw-local-20261104/faithful/hd-generated/manifest.json)、
[瀏覽器抽驗](../../workplace/hd-redraw-local-20261104/faithful/qa/native-result.json)、
[來源及收尾收據](../../workplace/hd-redraw-local-20261104/faithful/qa/native-verification.json)。
重建比較頁使用 `faithful/build.py`；瀏覽器抽驗使用 `qa/native-check.cjs`，均在上述本機工作區與既有Docker工具執行。
前版HTML與重生程式另存 `index-local-retouch.html`、`build-local-retouch.py`、
`template-local-retouch.html`、`app-local-retouch.js`及 `manifest-local-retouch.json`，前稿PNG與技術收據未覆寫。
使用者確認原話、範圍與兩張選稿SHA在輸出清單的 `user_approval`。這兩張為已確認的HD美術基準。
已確認的許昌單城頁另存 `index-hd-baseline.html`，兩張基準圖的原始PNG bytes與SHA未變。

使用者確認方向後回覆「ok go」，第二批依同一基準新增滎陽、舞陽、虎牢關、汝南的地圖及景觀，共八張選稿：

| 城池 | 據點ID／世界座標 | KYOGRF頁碼 | 原圖特色及重繪檢查 |
|---|---|---|---|
| 滎陽 | 63／192,102 | 14 | 地圖北側水岸與封閉北牆、三向城門；景觀長城牆、偏左門洞、右上樓閣及前方水道 |
| 舞陽 | 87／191,119 | 2 | 小城的北門及向南彎路；景觀斜向城牆、門洞、階梯與紅黃旗幟 |
| 虎牢關 | 72／200,106 | 5 | 地圖狹長關牆及橫向道路；景觀右側大面積弧形城牆、左側岩壁與前方木橋 |
| 汝南 | 100／216,137 | 0 | 大城配置及右上農田，未借用許昌右上水域；景觀中央門樓、左塔與紅金主殿 |

每張都以自己的96×96或112×112原圖作唯一內容依據，許昌已確認圖只作材質、線條及色彩完成度參考。
舞陽地圖另參考滎陽已繪小城的共同建築筆觸，仍以舞陽原圖控制北門及周圍地形。
虎牢關地圖第1稿的下方關牆未延伸至畫面邊緣，第2稿作單一修正，選第2稿；原稿未覆寫。
八張新增選稿與九次生成記錄在 [第二批清單](../../workplace/hd-redraw-local-20261104/faithful/hd-generated/batch-2.json)，
原圖及最近鄰參考SHA在 [第二批輸入](../../workplace/hd-redraw-local-20261104/faithful/hd-request/batch-2-inputs.json)，
所有完整提示詞沿前述重繪請求清單進入。原生輸出未後製放大或貼回柔化像素。
使用者對此批回覆「ok 同意」。第二批清單的 `user_approval` 保存原話、八張圖的檔名及SHA，
連同許昌兩張，五城十張已成為確認的HD美術基準。
規格整體仍為DRAFT，不作原版像素parity、正式HD完成或人類畫師作品聲明；全圖、其餘城池與HD頭像尚未完成。

## 大地圖HD試繪與全圖導航

使用者表示「期待大地圖hd」，先做
[本機大地圖預覽](../../workplace/hd-redraw-local-20261104/faithful/world/index.html)。
此頁與城池比較頁分開保存，均為內嵌圖片及資料的離線HTML。

- 全圖使用既有world／palette解析器重建的6144×4096原版圖，對應384×256格、每格16像素。
  中原試繪取世界格半開範圍180,96至238,145，原圖928×784；裁切與原版重建全圖逐像素相同。
  13個據點及27段道路的資料沿既有原始資料解析輸出，不從模型圖片反推座標或改遊戲規則。
- 內建imagegen以此原圖為唯一內容依據，許昌大城、舞陽小城及虎牢關的已確認HD圖只供材質、線條與色彩參考。
  原生輸出1364×1153，PNG SHA `d0ea55e67af754381907486defd14c3c3021cccf68279097b333bce46b2fb511`，
  工具原檔與專案副本bytes相同，未作後製放大。畫面保留河道、林地丘陵、農田及城池群的原版配置，
  繪製細節及邊線不是逐像素幾何證據。
- 各原生試繪圖映射回原世界矩形，原版模式保留完整原圖，縮圖以金框標出各HD區域。全圖面積覆蓋已完成，局部美術仍待修整。
  滑鼠滾輪依游標位置縮放，拖曳平移，工具列切原版／HD、全圖／已繪範圍／本區及原生細節，縮圖可移動視野。
  據點定位與可關閉的城名依原資料，沒有畫進生成圖。放大上限依原生圖片與顯示比例計算，避免超過原生細節。
- 現行三十五區版十項斷網瀏覽器抽驗通過，JS錯誤0、外部請求0，涵蓋內嵌bytes、座標、顯示切換、游標縮放、拖曳、
  共同原生上限、全圖／已繪範圍／本區、區域選擇、縮圖、據點定位、名稱避讓／開關、來源說明、單檔搬移及390×844。
  這是本機美術與鏡頭原型，未接入Ebiten正式路徑，也不代證正式地圖鏡頭或可玩流程。

來源入口：[輸入與原座標](../../workplace/hd-redraw-local-20261104/faithful/world/inputs.json)、
[生成提示詞](../../workplace/hd-redraw-local-20261104/faithful/world/generation-r1.json)、
[圖檔清單](../../workplace/hd-redraw-local-20261104/faithful/world/manifest.json)、
[瀏覽器抽驗](../../workplace/hd-redraw-local-20261104/faithful/world/qa-result.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/world/verification.json)。
重建使用 `world/build.py`，全圖來源唯讀掛為 `/source`，world輸出目錄可寫掛為 `/local/world`；
同目錄 `app.js` 與 `template.html` 內嵌為單檔頁面，`qa-check.cjs`及 `verify.py` 分別驗證瀏覽器與來源完整性。
工具沿用既有rs2-hdtools及Playwright容器，有界、1000:1000、network none。圖片處理依既有雲端授權，
成果與HTML保存在本機，沒有Sites部署或公開散布原版素材。

### 區域擴展與交界補繪

使用者「go」後，以中原為中心逐批擴展，左右最後各五區已補齊，全圖由三十五區覆蓋。各區顯示範圍（core）採世界格半開範圍，新圖額外帶四格共同邊緣作生成參考；到原圖外緣即停止裁取，不補造外側地形。
這些是本機原型的圖像映射資料，沒有更動正式地圖、路徑、命中區或存檔格式。

| 區域 | 顯示範圍世界格 | 原圖參考尺寸 | 原生選稿尺寸 | 據點數 |
|---|---|---|---|---|
| 中原 | 180,96至238,145 | 928×784 | 1364×1153 | 13 |
| 北側 | 180,47至238,96 | 928×848 | 1312×1199，選第2稿 | 11 |
| 東側 | 238,96至296,145 | 992×784 | 1411×1115 | 8 |
| 南側 | 180,145至238,194 | 928×848 | 1312×1199 | 8 |
| 西側 | 122,96至180,145 | 992×784 | 1411×1115 | 7 |
| 西北 | 122,47至180,96 | 992×848 | 1356×1159 | 3 |
| 東北 | 238,47至296,96 | 992×848 | 1357×1159，選第2稿 | 14 |
| 西南 | 122,145至180,194 | 992×848 | 1356×1159 | 12 |
| 東南 | 238,145至296,194 | 992×848 | 1357×1159，選第2稿 | 9 |
| 西北外緣 | 64,47至122,96 | 992×848 | 1356×1159 | 4 |
| 西側外緣 | 64,96至122,145 | 992×912 | 1309×1202，選第2稿 | 10 |
| 西南外緣 | 64,145至122,194 | 992×848 | 1356×1159 | 3 |
| 東北外緣 | 296,47至354,96 | 992×848 | 1357×1159，選第2稿 | 5 |
| 東側外緣 | 296,96至354,145 | 992×912 | 1308×1202 | 3 |
| 東南外緣 | 296,145至354,194 | 992×848 | 1356×1159 | 8 |
| 北側西段 | 64,0至122,47 | 992×816 | 1383×1137 | 0 |
| 北側中西段 | 122,0至180,47 | 1056×816 | 1426×1103 | 3 |
| 北側中段 | 180,0至238,47 | 1056×816 | 1426×1103 | 4 |
| 北側東段 | 238,0至296,47 | 1056×816 | 1426×1103 | 5 |
| 北側海岸 | 296,0至354,47 | 992×816 | 1383×1137，選第2稿 | 1 |
| 南側西段 | 64,194至122,256 | 992×1056 | 1235×1274 | 3 |
| 南側中西段 | 122,194至180,256 | 1056×1056 | 1254×1254 | 5 |
| 南側中段 | 180,194至238,256 | 1056×1056 | 1254×1254，選第2稿 | 10 |
| 南側東段 | 238,194至296,256 | 1056×1056 | 1254×1254 | 6 |
| 南側海岸 | 296,194至354,256 | 992×1056 | 1216×1294 | 3 |
| 西側邊緣北段 | 0,0至64,47 | 1088×816 | 1448×1086 | 0 |
| 西側邊緣中北段 | 0,47至64,96 | 1088×912 | 1370×1148 | 4 |
| 西側邊緣中段 | 0,96至64,145 | 1088×912 | 1371×1148，選第2稿 | 11 |
| 西側邊緣中南段 | 0,145至64,194 | 1088×912 | 1370×1148 | 11 |
| 西側邊緣南段 | 0,194至64,256 | 1088×1056 | 1274×1235 | 5 |
| 東側邊緣北段 | 354,0至384,47 | 544×816 | 1024×1536，選第2稿 | 0 |
| 東側邊緣中北段 | 354,47至384,96 | 544×912 | 968×1625 | 0 |
| 東側邊緣中段 | 354,96至384,145 | 544×912 | 968×1625，選第3稿 | 0 |
| 東側邊緣中南段 | 354,145至384,194 | 544×912 | 968×1625 | 1 |
| 東側邊緣南段 | 354,194至384,256 | 544×1056 | 900×1747 | 2 |

三十五個顯示矩形不重疊，完整覆蓋6144×4096原圖，合計全部192個原資料據點及25,165,824個原圖像素。
這是試繪的100%面積覆蓋；局部美術驗收、頭像與正式遊戲接線未完成，不以此宣稱完整HD完成。
地圖與資料的裁切／座標驗證，不代表生成圖裡的河岸、道路或符號逐像素一致。
北側第1稿借用了風格參考裡的方格農田，原圖此區沒有這些物件；第2稿單獨移除新增農田。
中原原生PNG保持SHA不變，原單片HTML、程式與收據另存帶 `central-r1` 的歷史檔；兩版均可重建。

交界先以正常滑鼠在原生細節倍率巡查。東側出現河岸台階，西側出現同一小城的重複片段，
因此從原資料取兩條384×784的原圖走廊，並用瀏覽器將當前HD拼接渲染為參考擷取圖，
再交內建imagegen補繪。參考擷取圖不是新美術成品。
原生東側補繪878×1792、西側補繪1013×1552，均直接保存工具PNG，未後製改圖。
檢視器將它們映射回原世界走廊的中間192像素範圍，外緣用8個原圖像素的透明度接合，
沒有模糊濾鏡，也不將柔化原版像素覆蓋回HD。這是顯示層合成，不改PNG或遊戲資料。
主要河岸台階與重影已改善；仍有局部筆觸、配色及邊線差異，不稱最終無縫或原版精確幾何。

目前可切本區、已繪範圍及原版全圖；原生細節上限依三十五區與補繪圖的共同像素密度計算，
原本單片的「1:1」按鈕改標為「原生細節」，避免將不同原生密度宣稱為同一嚴格1:1。
城名以實際字寬計算矩形，避開已繪名稱及畫布四邊；資料本身保留全部192個據點。

來源與選稿入口：[三十五區及完整提示詞清單](../../workplace/hd-redraw-local-20261104/faithful/world/regions.json)、
[新增原圖參考](../../workplace/hd-redraw-local-20261104/faithful/world/expansion-inputs.json)、
[交界參考與原座標](../../workplace/hd-redraw-local-20261104/faithful/world/seam-inputs.json)、
[交界檢視與限制](../../workplace/hd-redraw-local-20261104/faithful/world/seam-review.json)。
四角直接以各自原圖為內容依據，相鄰兩張HD只供筆觸、原色與邊界參考。西北原圖沒有農田；東北與東南第1稿仍誤加金色田格，
第2稿針對此項移除。既有五區原生PNG未改，五區HTML／manifest／清單／收據保存帶 `five-r1` 的歷史副本。
四角初稿合成後，以正常滑鼠在四個交會處抽驗，發現橫向河岸斷差與道路偏移。再取四條1024×384原版走廊，
以當前拼接的2倍離屏擷取作編修目標，交內建imagegen補繪四張2048×768原生圖；中央1024×192範圍用於顯示，
四邊各8個原圖像素作透明度接合。外側參考區不當成新美術成品。東北補繪第1稿誤加向右的L形路段，
第2稿移除錯路，恢復大城南門直接縱向接到下方小城。既有兩張縱向交界保留，新橫向補繪只在自己的顯示矩形覆蓋。
四個交會處的修正前後畫面已目視，主要斷差與錯路改善；局部筆觸、配色及邊線差異仍在，不宣稱最終無縫。

新增來源與選稿：[四角原圖及原座標](../../workplace/hd-redraw-local-20261104/faithful/world/diagonal-inputs.json)、
[四角完整提示詞及選稿](../../workplace/hd-redraw-local-20261104/faithful/world/diagonal-generation-set.json)、
[橫向交界來源](../../workplace/hd-redraw-local-20261104/faithful/world/diagonal-seam-inputs.json)、
[檢視與限制](../../workplace/hd-redraw-local-20261104/faithful/world/diagonal-review.json)。
重生入口沿 `build.py`／`app.js`／`template.html`；`integrate-diagonals.py`保存本次一次性五區→九區清單轉換，
`prepare-diagonal-seams.cjs`保存補繪參考擷取流程，`save-neighbours.py`與 `save-diagonal-seams.py`保存原生輸出複製及metadata格式。
同目錄 `verification.json` 核對三十五張原圖裁切、全部192個據點錨點、完整來源Path、四十六次區域生成及二十五次交界生成的原生bytes，
共五十七張選稿。現行頁與中原歷史頁重建同SHA，二十五區及更早快照SHA未變；五城十張已同意美術及城池比較頁未改。
此處仍為DRAFT原型，正式世界資產接線尚未完成。

### 西側三區與接續抽驗

2026-11-05沿既有HD方向加入長安、漢中所在西側三區，新增17個原資料據點。各區只用自己的原圖決定內容，
東鄰HD只作筆觸、色彩與邊界參考。西側外緣第1稿多畫上方河岸田格，第2稿移除原圖沒有的田格；原生初稿均保留。
以正常滑鼠在三處新縱向接縫與兩處橫向接縫抽驗，再補繪中段／南段兩條縱向及上下兩條橫向走廊，
原生尺寸分別878×1792及2048×768。顯示方式沿既有core裁切與8原圖像素透明度接合，沒有模糊或PNG後製。
主要河岸斷差、支流與山林錯位已改善；北段縱向接縫的局部筆觸與色差仍列可選修整，不作最終無縫聲明。
當時放大上限由所有選稿的原生像素密度共同決定。十二區單檔HTML內嵌原圖、二十二張選稿與來源資料，可離線開啟。

來源入口：[西側原圖與座標](../../workplace/hd-redraw-local-20261104/faithful/world/west-column-inputs.json)、
[完整提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/world/west-column-generation-set.json)、
[接縫來源](../../workplace/hd-redraw-local-20261104/faithful/world/west-column-seam-inputs.json)、
[目視檢查與限制](../../workplace/hd-redraw-local-20261104/faithful/world/west-column-review.json)。
`prepare-west-column.py`保存九區快照與原圖裁取，`integrate-west-column.py`保存這次一次性清單擴展；
`capture-seam-references.cjs`依 `west-column-seam-plan.json`保存接縫參考擷取流程，原生圖沿既有save腳本複製保存。
九區歷史入口為 `index-nine-r1.html`，快照雜湊在 `nine-r1-snapshot.json`；現行重建與抽驗仍沿 `build.py`、`qa-check.cjs`及 `verify.py`。

### 東側三區、接縫與離線載入

沿使用者正向回饋及繼續指示，加入臨淄、北海、廣陵、建業、會稽所在東側三區，共新增16個原資料據點。
保留原版海岸、河口及湖泊輪廓，平原區不因風格參考而新增山林；東北第1稿誤加田格，由第2稿移除。
兩條縱向與兩條橫向補繪改善東北大城碎片、南側河岸與兩條海岸／河口接縫；尺寸與合成方式沿既有原生輸出及8像素透明度接合。
中段縱向接縫的局部筆觸與色差列可選修整，不作最終無縫或逐像素幾何聲明。

來源：[東側原圖與座標](../../workplace/hd-redraw-local-20261104/faithful/world/east-column-inputs.json)、
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/world/east-column-generation-set.json)、
[接縫來源](../../workplace/hd-redraw-local-20261104/faithful/world/east-column-seam-inputs.json)、
[目視檢查與限制](../../workplace/hd-redraw-local-20261104/faithful/world/east-column-review.json)。
`prepare-east-column.py`與 `integrate-east-column.py`保存這次裁圖、快照與一次性清單擴展，十二區歷史頁為 `index-twelve-r1.html`。

單檔初稿115.8MB將所有PNG base64放在同一JSON，斷網容器初次就緒約20–25秒，同頁第二次開啟超過既有30秒門檻。
`build-payload-probe.py`製作可撤回原型：PNG原始bytes分別存於不執行的ASCII script區塊，JSON只存圖片ID；
`app.js`讀取後清空文字區塊，仍等整組素材全部載入才顯示。沒有放大、轉碼或更改PNG，也沒有改正式遊戲資產／存檔格式。
相同2 CPU／3GB容器、相同30秒門檻，115.8MB原型初次就緒約16–17秒、同頁重開約7.5秒；這些是局部比較，非所有機器的速度承諾。
最終129.7MB十五區頁已通過完整十項抽驗及搬至不同目錄後同頁重新開啟，所有內嵌PNG的SHA與尺寸仍和原生檔相同。
來源與重播：[載入方法](../../workplace/hd-redraw-local-20261104/faithful/world/payload-loading.json)、
[原方式量測](../../workplace/hd-redraw-local-20261104/faithful/world/qa-load-probe.json)、
[分塊原型量測](../../workplace/hd-redraw-local-20261104/faithful/world/qa-load-probe-chunked.json)；重播程式 `qa-load-probe.cjs`。

接縫參考擷取改由 `capture-seam-references.cjs`先核對本機PNG的SHA，再透過容器內暫時loopback同源讀取。
Docker仍為network none，沒有主機port、外部介面或公開網站，工作完成即關閉伺服器。這只產生補繪參考圖，
正式預覽仍直接以單一HTML離線開啟，瀏覽器抽驗沒有外部請求。完整世界及正式引擎接線仍待完成。

### 北側五區、原圖上緣與載入抽驗

新增晉陽、中山、真定、南皮與黃等13個原版據點。西段沒有據點，仍保留兩處水塘、河道與細長城牆，沒有增建城堡；
中段與東段不添加原圖沒有的農田。北側海岸第1稿誤加金色田格，第2稿移除，保留南側真正的綠色農田。
原版資料中的「北京」沿用原名，不憑通行地名改寫。城名及座標來自原資料，不從生成圖讀取。

九處新交界以正常滑鼠抽樣。中段出現重複城池及分離屋頂，另有東段大城及海岸斷差，
取兩條1056×384及一條992×384原圖走廊補繪。原生選稿2079×756／2079×756／2016×780，直接保存工具PNG，
顯示沿既有中央192像素走廊與四邊8像素透明度接合。主要城池碎片與橫向斷差已改善；其他筆觸、配色及縱向邊線仍有差異。
原圖上緣的鏡頭沿既有世界邊界限制，抽驗改查目標矩形完整可見，不要求鏡頭固定在目標矩形中心。
零據點的北側西段可正常切換，清空前一區據點按鈕；原版／HD、縮放、拖曳與縮圖仍沿用既有操作。

151.1MB初稿在完整操作後搬移重開超過30秒。獨立重開測試三次皆通過，沒有OOM事件，尚未證實唯一原因。
可撤回原型改用原始PNG bytes建立暫存圖片網址，先清除base64文字，載入成功或失敗即釋放網址，仍等待整組圖片就緒。
相同2 CPU／3GB與30秒門檻，151.1MB原型三次開啟約6.3–7.3秒；不把此數字當成最終162.3MB頁或所有機器的速度承諾。
最終162,275,252 bytes頁面通過完整十項抽驗、單檔搬移後同頁重開與390×844；JS錯誤0、外部請求0。

來源：[北側原圖及座標](../../workplace/hd-redraw-local-20261104/faithful/world/north-band-inputs.json)、
[完整提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/world/north-band-generation-set.json)、
[接縫原圖](../../workplace/hd-redraw-local-20261104/faithful/world/north-band-seam-inputs.json)、
[前後檢視與限制](../../workplace/hd-redraw-local-20261104/faithful/world/north-band-review.json)、
[載入比較](../../workplace/hd-redraw-local-20261104/faithful/world/north-loading.json)。
`prepare-north-band.py`保存十五區快照與裁圖，`integrate-north-band.py`與 `integrate-north-seams.py`保存一次性擴展；
歷史頁為 `index-fifteen-r1.html`，快照雜湊入口為 `fifteen-r1-snapshot.json`。目前頁面由既有 `build.py`重生。
所有圖片沿已確認方向製作，沒有把本次「go」記為新圖逐張已獲使用者接受；其餘世界、頭像及正式引擎接線仍待完成。

### 南側五區、原圖下緣與接縫抽驗

南側五區新增長沙、武陵、桂陽、柴桑及建安等27個原版據點，各區顯示範圍延伸至世界格y=256。
西段原圖沒有農田，保留密集的低矮山林、溪流與三座小城；其他區域依各自原圖保留大湖、島岸、河道與真正的綠色田格。
中段第1稿把右上河邊的橘黃色抖色畫成金色田格，第2稿針對這一項修正；原生初稿及修稿均保存。
五區仍映射到原世界矩形，繪製邊線與原生輸出比例不作逐像素幾何證據。放大上限沿既有程式取所有選稿的最低真實像素密度，沒有放大PNG補數字。

九處新交界及兩處既有交界以正常滑鼠抽樣。西段的城鎮、窄河與山林在列邊界出現裁切斷差，中段湖岸與小城也有碎片，
分別取992×384及1056×384原圖走廊，依當前組版2倍參考補繪2014×780及2079×756原生圖。
顯示沿既有中央192像素走廊與四邊8原圖像素透明度接合，PNG未後製，主要新接縫斷差改善；局部配色、筆觸及縱向邊線仍列可選修整。

最終195,611,745 bytes單檔通過完整十項抽驗、操作後搬移至不同目錄同頁重開及390×844；JS錯誤0、外部請求0。
保留原有世界邊界限制，逐區檢查原圖上下緣的目標矩形完整可見。已繪範圍的視野可能比全圖更寬，此時世界邊界會固定中心；
游標縮放檢查先用工具列選回中原，避免把合理的世界邊界限制判為焦點錯誤，沒有更改縮放規則。
瀏覽器收據固定起始HTML雜湊，結束時確認頁面未改，且選稿數與清單相同，避免以末端新檔案的SHA取代受測檔。

來源：[南側原圖及座標](../../workplace/hd-redraw-local-20261104/faithful/world/south-band-inputs.json)、
[完整提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/world/south-band-generation-set.json)、
[接縫原圖](../../workplace/hd-redraw-local-20261104/faithful/world/south-band-seam-inputs.json)、
[前後檢視與限制](../../workplace/hd-redraw-local-20261104/faithful/world/south-band-review.json)。
`prepare-south-band.py`保存二十區快照與裁圖，`integrate-south-band.py`及 `integrate-south-seams.py`保存一次性擴展；
歷史頁為 `index-twenty-r1.html`，快照雜湊入口 `twenty-r1-snapshot.json`。目前頁面仍由 `build.py`重生，
抽驗與完整性核對仍用 `qa-check.cjs`及 `verify.py`，均在既有斷網容器執行。
這批圖沿已確認畫風產製，沒有新增逐張接受的裁定。左右外側帶、局部邊線、頭像及正式引擎接線仍待完成。

### 左右邊緣十區與全圖覆蓋

左右各五區新增34個原版據點，包括西涼、成都、吳、章安及溫州。西北零據點區保留三處原水域及細長城牆，
西側關隘保留窄關牆，不改成大型城堡；東側三個零據點區保留空海面與原岸線。各圖只依自己的原圖決定內容。
西側中段及東北海岸第1稿誤加金色田格，各選第2稿移除；東側中段水紋選第3稿減少粗條帶，局部水紋與色差仍在。
新增區域共十四次原生生成，所有PNG及未選稿都保存。既有二十五區PNG及已同意城池圖未改。

正常滑鼠抽驗十四處新接縫及兩處既有接縫，另用縮圖抵達四個世界角落。西側中南段的宕渠出現重複城鎮片段，
另有兩處列邊界切到屋頂及河岸。依原圖補繪一張878×1792縱向、兩張2141×734橫向原生圖，沿既有中央192像素走廊及四邊8像素透明度接合。
主要重複城鎮與岸線斷差已改善。中北／中段交界仍有孤立屋頂；該橫向補繪第2稿產生額外小物件，未選用，保留第1稿並列待修。
局部配色、水紋、筆觸與邊線仍有差異，不宣稱最終無縫或逐像素幾何。

全圖巡查第一次在最西跨最東時碰到舊測試32段拖曳上限；左側七張擷取已完成，第一張右側擷取未到場。
測試改以距離、顯示倍率及每段滑鼠距離計算有界步數，並先將目標鏡頭換算到原世界邊界內，沒有改遊戲或預覽鏡頭規則。
相同資源與門檻下乾淨重跑通過，最終再依同一流程驗證新增補繪。沒有以直接注入鏡頭代替滑鼠。

最終250,657,792 bytes單檔通過十項斷網抽驗，包括三十五區、零據點區、四邊／四角、原版／HD、滑鼠縮放／平移、
來源與城名、搬移後同頁重開及390×844。JS錯誤0、外部請求0，保留原有30秒載入門檻。
全圖覆蓋後，頁面改顯示「全圖HD試繪」及局部待修說明，避免繼續顯示其餘區域仍用原圖。
完整性收據明列 `full_original_world_painted_coverage=true`、`full_world_HD_complete=false`及 `engine_integration=false`。

來源：[左右外側原圖及座標](../../workplace/hd-redraw-local-20261104/faithful/world/outer-bands-inputs.json)、
[完整提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/world/outer-bands-generation-set.json)、
[接縫原圖](../../workplace/hd-redraw-local-20261104/faithful/world/outer-bands-seam-inputs.json)、
[修正、未選稿及限制](../../workplace/hd-redraw-local-20261104/faithful/world/outer-bands-review.json)。
`prepare-outer-bands.py`保存二十五區快照及原圖，`integrate-outer-bands.py`與 `integrate-outer-seams.py`保存一次性擴展；
歷史入口為 `index-twentyfive-r1.html`，快照入口 `twentyfive-r1-snapshot.json`。目前頁面仍以 `build.py`重生，
`qa-check.cjs`及 `verify.py`分別保存瀏覽器與來源檢查。全圖仍為本機DRAFT原型，新圖未宣稱已獲逐張美術驗收。

### 全圖屋頂細修與原版頭像三張小樣

全圖覆蓋後，西側中北／中段交界的孤立藍屋頂仍在。先以洋紅框標明原生補繪中的錯誤物件，
再由內建imagegen依完整補繪、框選參考與原圖生成第3稿。孤立屋頂及短底座已移除，沒有用新橋或小屋替代，
左下完整城鎮、既有河道、田格及道路保留。PNG為2141×734原生輸出，未手動後製。
三十五區原生區域圖未改，修正前HTML、metadata、程式、收據與局部畫面保存 `thirtyfive-r1` 快照。
目前仍為57張選稿加原圖，原生生成累計71次。250,663,972 bytes單檔依原資源及門檻完成十項斷網抽驗，
JS錯誤0、外部請求0，局部配色、水紋與邊線仍留待修整。

頭像沿原版 `KAOGRF.DAT` 及 `General.Portrait` 頁碼製作，沒有把武將ID當成圖庫索引。
來源使用既有解碼原圖，核對圖庫、色盤、劇本及三張64×64 PNG的雜湊，沒有重畫成通用三國人物。

| 武將 | 武將ID | 圖庫頁碼 | 原圖 | 選稿原生尺寸 | 選稿 |
|---|---|---|---|---|---|
| 曹操 | 16 | 50 | 64×64 | 1254×1254 | 第1稿 |
| 關羽 | 40 | 114 | 64×64 | 1254×1254 | 第1稿 |
| 諸葛亮 | 60 | 125 | 64×64 | 1254×1254 | 第2稿 |

各以自己的原圖保留頭部構圖、臉型、神情、服飾色塊及黑底，重畫清楚髮鬚、衣料與輪廓。
諸葛亮首稿轉得太正且多出金色衣肩，第2稿恢復較明確左轉角度與深藍衣肩；四次原生輸出皆保存且bytes未改。
64×64原圖以最近鄰放大只作參考，不當成HD成果。細部髮絲與衣料屬重繪詮釋，臉部定位不逐像素相同，
三張小樣尚未獲使用者逐張接受，不代表150頁頭像全數完成。

[首批頭像對照（歷史）](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-three-r1.html)可選三位武將，
切原版／HD／並排／拉桿，同步縮放及滑鼠拖曳。工具列常駐且不隨圖片縮放，HD 1:1依螢幕像素密度限制原生上限。
七項斷網抽驗包括六張內嵌PNG bytes、武將及頁碼映射、三位切換、四種模式、滑鼠拉桿、滾輪／工具列縮放、
拖曳／重設、來源說明、單檔搬移重開、390×844與兩倍像素密度；JS錯誤0、外部請求0。
本機頁8,975,872 bytes，重建同SHA。這是本機美術及圖片檢視原型，沒有接入正式頭像繪製、遊戲或存檔。

來源：[屋頂修正與前後畫面](../../workplace/hd-redraw-local-20261104/faithful/world/detail-polish-review.json)、
[全圖修正前快照](../../workplace/hd-redraw-local-20261104/faithful/world/thirtyfive-r1-snapshot.json)、
[頭像原圖及映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals.json)、
[首批提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-three-r1.json)、
[首批目視檢視及限制](../../workplace/hd-redraw-local-20261104/faithful/portraits/review-three-r1.json)、
[首批七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-three-r1.json)、
[首批來源及原生bytes核對](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-three-r1.json)。
頭像頁由同目錄 `build.py`重生，`qa-check.cjs`及 `verify.py`保存瀏覽器及完整性檢查。
來源準備與複製腳本在 `faithful/tooling/prepare-portraits.py`、`portrait-jobs.py`及 `save-portraits.py`，
只處理參考圖與原生檔複製，最終美術修訂交內建imagegen。

### 頭像擴充至九張及趙雲原圖勘誤

使用者「ok go」延續原圖參考重繪。新增六張頭像，各依原劇本中的 `General.Portrait` 取圖，
沿用 `internal/assets/library`、`internal/state` 及原色盤，沒有以武將ID代替圖庫頁碼。
來源仍是64×64原圖，全部選稿1254×1254，最近鄰參考放大不當成HD成果。

| 新增武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|
| 劉備 | 32 | 94 | 第1稿 |
| 孫權 | 18 | 128 | 第1稿 |
| 司馬懿 | 55 | 9 | 第2稿 |
| 張飛 | 80 | 118 | 第1稿 |
| 趙雲 | 73 | 6 | 第1稿 |
| 周瑜 | 57 | 110 | 第1稿 |

各圖保留原版人物的配色、頭部朝向及服飾輪廓。孫權沿原圖的棕色鬍鬚，周瑜保留細鬍鬚，
張飛保留張口露齒的神情。司馬懿首稿臉部較居中，第2稿恢復更明確的右轉及右緣裁切。
八次新增原生輸出均保留，連同首批四次共十二次；九張選稿的PNG bytes未後製。

趙雲的初次人工判讀有誤：將藍色額飾與深色上唇線當成多出的細節，要求第2稿移除。
重新檢視原圖的1024×1024最近鄰參考，並直接讀64×64原圖RGB後確認，零起算座標 `(20,7)` 等六個額前像素是
`RGB(48,65,211)`；上唇 `(21,36)` 等位置有黑色／深綠筆畫。這些顏色與位置屬已證實，
每條深色線究竟是鬍鬚還是陰影仍未知。依可見原線形與顏色保留第1稿，第2稿不選用。
錯誤提示詞及未選稿保留作歷史，後續以 [原圖細節勘誤](../../workplace/hd-redraw-local-20261104/faithful/portraits/portrait-6-source-detail-review.json)為準，
不能再由歷史提示詞宣稱原圖沒有藍色或必須無鬍鬚。

[九張對照（歷史）](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-nine-r1.html)預設開新增劉備，
武將選單、標題、數量及來源表同步，可切並排／拉桿／原版／HD，保留既有滑鼠縮放、拖曳與HD 1:1。
選稿、嘗試稿與來源由同目錄 `selection.json`供 `build.py`重生，這只是本機原型資料，未改正式資產或存檔格式。
最終29,763,303 bytes單檔通過七項斷網抽驗，包括十八張內嵌PNG bytes、九位切換、新增六張畫面、四種模式、
滑鼠拉桿／縮放／拖曳、來源表、搬移重開，以及390×844兩倍像素密度下的九列來源視窗開關與原生上限。
JS錯誤0、外部請求0。十二次原生bytes、九張原圖、原始資源雜湊及勘誤像素核對通過，重建同SHA。
首批三張選稿、十張已同意城池圖及全圖頁均未改；三張舊頁與誤選趙雲第2稿的九張候選各有獨立快照。

來源：[新增原圖及映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-2.json)、
[九張版本提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-nine-r1.json)、
[第二批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-2.json)、
[首批快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/three-r1-snapshot.json)、
[已被取代候選](../../workplace/hd-redraw-local-20261104/faithful/portraits/nine-first-candidate-snapshot.json)、
[九張版本七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-nine-r1.json)、
[九張版本完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-nine-r1.json)。
原圖解碼腳本為 `faithful/tooling/extract-portrait-batch2.go`，使用既有 `hr-go-ebiten:1.26.7-2.9.9-r1` 容器，
唯讀專案與原資料，只寫既有頭像工作目錄。九張仍是150頁中的小樣，細部雕飾及髮鬚是重繪詮釋，
未宣稱逐像素相同、逐張接受、完整頭像或正式引擎完成。

### 頭像擴充至21張及來源視窗修正

使用者「go」延續原圖參考HD方向。本批新增12張，原圖64×64，選稿皆1254×1254原生PNG。
先由 `SINARIO.DAT` 四個出貨劇本與既有解析器取149個不重複名字，確認實際武將ID、圖庫頁碼及來源劇本。
這個名字集合與 `KAOGRF.DAT` 的150頁分開計數，不作一人一圖或頭像完成數。
預選董卓、孫堅不在該集合，改取實際存在的荀彧、陸遜；未創造人物或拿其他頁假代。
本批12筆均可在劇本索引0找到。來源解析器、計畫、原始資料及解碼PNG的雜湊皆保存。

| 新增武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|
| 呂布 | 35 | 102 | 第1稿 |
| 荀彧 | 62 | 121 | 第1稿 |
| 袁紹 | 2 | 105 | 第1稿 |
| 袁術 | 3 | 68 | 第1稿 |
| 陸遜 | 101 | 144 | 第1稿 |
| 孫策 | 19 | 127 | 第1稿 |
| 馬超 | 25 | 22 | 第1稿 |
| 黃忠 | 52 | 119 | 第1稿 |
| 夏侯惇 | 39 | 49 | 第1稿 |
| 夏侯淵 | 37 | 62 | 第1稿 |
| 張遼 | 81 | 44 | 第1稿 |
| 龐統 | 95 | 104 | 第2稿 |

各依自己的原圖重畫細線與材質，保留原色色塊、朝向與冠盔。袁紹及黃忠保留灰白鬍鬚，
夏侯惇沿原圖保留雙眼，不加眼罩；馬超保留金色冠盔及紅色額飾，不換成其他作品常見的銀盔。
龐統第1稿將原圖的單一圓形帽飾畫成三瓣形，第2稿只修回圓／橢圓金框，首稿也保存。
本批13次原生輸出，連同前批共25次，PNG bytes與圖片工具輸出相同，未手動後製。
既有九張選稿、已同意城池圖及大地圖頁未改，九張頁與收據保存 `nine-r1` 快照。

[當時21張本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-twentyone-r1.html)預設呂布，
保留原版／HD／並排／拉桿、滑鼠縮放／同步拖曳及依像素密度限制的HD 1:1。
第一次21張抽驗雖通過，目視來源視窗截圖發現開啟時看不到標題。正常滑鼠探針確認，
原生對話框聚焦長表格後的第一個連結，造成 `scrollTop=522`、標題y=-478，焦點為「提示詞與生成紀錄」。
將關閉鈕移到上方標題列，標題／關閉鈕隨捲動保留，開啟時重設到頂部；標題背景遮住上緣捲過的文字碎片。
來源視窗修正前頁、收據與截圖另存，沒有把原本七項技術檢查當成視覺已通過。

最終69,191,415 bytes單檔依同一2GB／1 CPU、30秒載入及90秒外層門檻完成七項斷網抽驗。
包括42張內嵌PNG bytes／尺寸、21位切換、新批預設、六種新增畫面抽樣、四種模式、滑鼠拉桿／縮放／拖曳／重設，
21列來源表、搬移單檔重開及390×844兩倍像素密度；另外檢查來源視窗開頭、滑鼠捲動後的標題／關閉位置，以及關閉再開回頂。
JS錯誤0、外部請求0。完整性檢查核對21張原圖、25次原生bytes、武將及頁碼、四劇本名單、原始資源雜湊、
既有趙雲像素勘誤及歷史快照，並確認重建同SHA。這仍是150頁中的21張美術小樣，細部筆畫是重繪詮釋，
未宣稱逐像素幾何、逐張接受、完整頭像或正式遊戲接線。

來源：[四劇本名單](../../workplace/hd-redraw-local-20261104/faithful/portraits/source-roster.json)、
[第三批原圖及映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-3.json)、
[當時提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-twentyone-r1.json)、
[第三批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-3.json)、
[九張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/nine-r1-snapshot.json)、
[來源視窗探針及前後畫面](../../workplace/hd-redraw-local-20261104/faithful/portraits/source-window-review.json)、
[當時七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-twentyone-r1.json)、
[當時完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-twentyone-r1.json)。
可重生入口仍為同目錄 `build.py`、`selection.json`、`qa-check.cjs`及 `verify.py`。
`faithful/tooling/extract-portraits.go`依 `batch-3-plan.json`取圖，`--roster`輸出原始名單；
沿用 `hr-go-ebiten:1.26.7-2.9.9-r1`、`rs2-hdtools:r1`與既有Playwright容器，只有本機工作目錄可寫。

### 頭像擴充至33張及單檔載入調整

使用者「繼續囉」延續原圖參考HD方向。本批再加12張，皆由既有解析器依原始名字、武將ID與圖庫頁碼取圖，
可在劇本索引0找到。原圖64×64，選稿皆1254×1254原生PNG；既有21張選稿與PNG雜湊不變。

| 新增武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|
| 賈詡 | 36 | 36 | 第2稿 |
| 荀攸 | 63 | 69 | 第2稿 |
| 徐庶 | 65 | 48 | 第2稿 |
| 法正 | 94 | 17 | 第1稿 |
| 陳宮 | 82 | 26 | 第2稿 |
| 魯肅 | 107 | 31 | 第2稿 |
| 呂蒙 | 105 | 40 | 第1稿 |
| 甘寧 | 42 | 39 | 第1稿 |
| 太史慈 | 72 | 20 | 第2稿 |
| 許褚 | 45 | 24 | 第1稿 |
| 典韋 | 88 | 120 | 第2稿 |
| 祝融 | 125 | 23 | 第1稿 |

逐張沿原圖保留黑底、朝向、頭部構圖與主要色塊。法正保留原圖淡灰額前方塊；祝融保留棕髮、褐色高冠、金紅頭帶與原有膚色。
賈詡、荀攸、陳宮與魯肅首稿將點陣陰影畫成大棋盤格，第2稿只將該處還原為連續衣料／髮絲；
徐庶第2稿移除新增黃色帽塊，太史慈與典韋第2稿移除新增紅色額塊。原圖本來就有的帽飾及額前色塊保留。
本批12張首稿及七張修正稿全數保存，共19次新增輸出；連同前批共44次，原生PNG逐byte複製，未手動後製。
細部雕飾、髮絲與筆觸是原圖的美術詮釋，不稱原版逐像素幾何。

[當時33張本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-thirtythree-r1.html)預設賈詡，
保留原版／HD／並排／拉桿、同步縮放／拖曳／重設、HD 1:1與33列來源表。
21張頁、metadata、程式及收據先保存不可覆寫的快照；來源視窗21列畫面也另存，歷史來源紀錄改連對應收據與截圖。

最初108,426,835 bytes的33張頁在原有30秒導覽門檻逾時，失敗頁、程式與錯誤訊息保存。
沿用大地圖的短期圖片網址方法，將未改動PNG分開內嵌，解碼後移除來源區塊，載入成功或失敗時釋放暫存網址。
這是本機檢視原型的載入調整，不改正式資料格式或引擎。大量資料網址及保留記憶體是原因假說，
沒有OOM事件或足以證明唯一原因的受控比較，不能寫成已確認的記憶體耗盡。

最終108,432,389 bytes單檔依相同2GB／1 CPU、30秒導覽及90秒外層門檻乾淨重跑，七項斷網抽驗通過。
三次正常／搬移後／窄視窗開啟約6.8、2.5、2.8秒；各66個暫存網址全部釋放，來源區塊全部移除。
抽驗包含66張內嵌PNG bytes與尺寸、33位切換及賈詡預設、六種新增畫面、四種模式、滑鼠拉桿／縮放／拖曳／重設，
33列來源表開頭／捲動／關閉／重開、搬移單檔離線重開與390×844兩倍像素密度，JS錯誤0、外部請求0。
已目視兩張新增總覽、正常頁面的賈詡並排、祝融細節及捲動後來源視窗。
完整性收據核對33張原圖、44次原生bytes、四劇本名單與頁碼、原始資源及解析器雜湊、趙雲既有像素勘誤、歷史快照與重建一致。
既有21張選稿、十張已同意城池PNG、城池頁與全圖頁未變。33張仍是150頁中的美術小樣，未宣稱逐張接受或正式接線完成。

來源：[第四批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-4-plan.json)、
[原圖及映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-4.json)、
[當時提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-thirtythree-r1.json)、
[第四批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-4.json)、
[21張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/twentyone-r1-snapshot.json)、
[載入調整前快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch4-loading-before-snapshot.json)、
[載入紀錄與限制](../../workplace/hd-redraw-local-20261104/faithful/portraits/loading-review.json)、
[當時七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-thirtythree-r1.json)、
[當時完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-thirtythree-r1.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-4-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-4.json`準備。
本機頁由 `portraits/build.py`重生，`qa-check.cjs`與 `verify.py`核對；本批產製入口為 `prepare-batch4.py`、`integrate-batch4.py`及 `finish-batch4-loading.py`。
沿用既有Go、Python與Playwright映像，以1000:1000及有界斷網容器執行，原資料唯讀，只寫既有頭像工作目錄。

### 頭像擴充至45張與來源劇本核對

使用者「ok go go」延續原圖參考HD方向。先保存33張頁、metadata、程式、收據、抽樣畫面及原生選稿雜湊，
再由既有解析器依實際名字取下列12張原圖。原圖64×64，選稿皆1254×1254原生PNG，本批皆選第1稿。

| 新增武將 | 來源劇本索引 | 武將ID | 圖庫頁碼 |
|---|---|---|---|
| 徐晃 | 0 | 64 | 108 |
| 張郃 | 0 | 75 | 76 |
| 曹仁 | 0 | 70 | 10 |
| 魏延 | 0 | 47 | 51 |
| 龐德 | 0 | 96 | 65 |
| 周泰 | 0 | 56 | 82 |
| 黃蓋 | 0 | 50 | 116 |
| 程昱 | 0 | 85 | 135 |
| 張昭 | 0 | 78 | 43 |
| 姜維 | 1 | 37 | 112 |
| 鄧艾 | 1 | 82 | 124 |
| 孟獲 | 1 | 107 | 146 |

武將ID連同來源劇本核對，不以不同劇本中的相同ID代替人物身份。原圖及資料雜湊、解析器雜湊與版本均保留。
沿各自原圖保留黑底、臉部構圖、冠盔輪廓及主要配色。黃蓋保留灰白眉鬚，曹仁保留單一紅色圓飾，
周泰及姜維保留藍色帽飾，張昭及鄧艾保留綠色帽塊，孟獲沿自己的原盔及紅色布帶繪製。
鄧艾上唇的小紅塊先回查原圖：64×64母本座標38,41確為RGB(211,0,0)，全圖只有這一個相同色值像素。
首稿有保留，沒有將它當成多出的裝飾移除。色點的敘事或材質意義未知，檢視紀錄只確認可見色值。
本批新增12次原生輸出，連同前批共56次；45張選稿逐byte複製且未後製，細部紋樣、髮絲及筆觸仍屬重繪詮釋。

[當時45張本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-fortyfive-r1.html)預設徐晃，
保留原版／HD／並排／拉桿、同步縮放／拖曳／重設及HD 1:1，來源視窗擴充至45列。
沿用33張時的分開內嵌PNG及短期圖片網址方法，沒有再改載入方式、原生PNG或正式遊戲資料格式。
載入調整的歷史紀錄改連33張收據，沒有將新的45張收據冒作當時的結果。

150,522,248 bytes單檔依原有2GB／1 CPU、30秒導覽與90秒外層完成七項斷網抽驗。
包含90張內嵌PNG bytes與尺寸、來源劇本／武將／頁碼、45位切換、徐晃預設、六種新增畫面，
四種模式、滑鼠拉桿／縮放／同步拖曳／原生上限／重設，45列來源表開頭／捲動／關閉／重開，
搬移單檔離線重開及390×844兩倍像素密度。三次載入約5.8、6.5、4.1秒，各90個暫存網址全部釋放、來源區塊全部移除，JS錯誤0、外部請求0。
原圖、原生輸出及兩張新增總覽已目視，正常頁面另抽看黃蓋、鄧艾並排及捲動後的來源視窗。
完整性檢查核對45張原圖、56次原生bytes、四劇本名單與來源劇本、原始資源及解析器SHA、趙雲既有像素勘誤、鄧艾紅點樣本、歷史快照及重建一致。
既有33張選稿、十張已同意城池PNG、城池頁及全圖頁未變。45張仍是150頁中的美術小樣，未宣稱逐張接受或正式接線完成。

來源：[第五批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-5-plan.json)、
[原圖與劇本映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-5.json)、
[當時提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-fortyfive-r1.json)、
[第五批檢視與限制](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-5.json)、
[鄧艾原圖色點核對](../../workplace/hd-redraw-local-20261104/faithful/portraits/portrait-124-source-detail-review.json)、
[33張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/thirtythree-r1-snapshot.json)、
[當時七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-fortyfive-r1.json)、
[當時完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-fortyfive-r1.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-5-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-5.json`準備。
提示詞與合併入口為 `portraits/prepare-batch5.py`及 `integrate-batch5.py`；頁面、抽驗與核對沿用 `build.py`、`qa-check.cjs`及 `verify.py`。
沿用既有Go、Python與Playwright映像，以1000:1000及有界斷網容器執行，原資料唯讀，只寫既有頭像工作目錄。

### 頭像擴充至57張與關興亮面修正

使用者「繼續囉」延續既定原圖參考HD方向。先保存45張頁、程式、metadata、收據、抽樣畫面及45張原生選稿，
共73個檔案雜湊，再由既有解析器取下列12張64×64原圖。選稿皆1254×1254原生PNG。

| 新增武將 | 來源劇本索引 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|---|
| 曹丕 | 0 | 124 | 126 | 第1稿 |
| 曹植 | 0 | 15 | 115 | 第1稿 |
| 劉禪 | 1 | 18 | 99 | 第1稿 |
| 關平 | 0 | 43 | 47 | 第1稿 |
| 關興 | 1 | 32 | 55 | 第2稿 |
| 張苞 | 0 | 126 | 41 | 第1稿 |
| 王平 | 1 | 23 | 103 | 第1稿 |
| 郝昭 | 1 | 25 | 42 | 第1稿 |
| 郭淮 | 1 | 26 | 64 | 第1稿 |
| 司馬師 | 1 | 47 | 25 | 第1稿 |
| 司馬昭 | 1 | 48 | 15 | 第1稿 |
| 諸葛恪 | 1 | 53 | 111 | 第1稿 |

四張來自劇本索引0、八張來自索引1，名字、武將ID與頁碼連同來源劇本核對。
各沿自己的黑底、臉部構圖、神情、冠盔及主要配色重畫。曹丕保留俯視角度，劉禪保留藍／金頭冠及圓臉，
張苞保留張口與深棕鬍鬚，司馬師保留傾斜構圖，諸葛恪保留綠色帽塊。
關興第1稿將原圖頭盔的淡色區畫成平面月牙，第2稿保留淡色區位置，改以金屬亮面筆觸處理，臉部及其餘構圖維持。
這是淡色區的美術詮釋，不宣稱已解出原版物件意義。兩稿都保存，其餘十一張選第1稿。
王平原圖額前的斜紅色圖形也保留；區域40,8起的5×5範圍內有十個RGB(211,0,0)像素，樣本保存於本批紀錄，符號意義未知。
本批12張首稿與一張修正稿共13次輸出，連同前批共69次。57張選稿逐byte複製，沒有後製或將放大圖當成HD原生稿。

[當時57張本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-fiftyseven-r1.html)預設曹丕，
保留原版／HD／並排／拉桿、同步縮放／拖曳／重設、HD 1:1及固定工具列，來源視窗擴充至57列。
沿用既有分開內嵌PNG與暫存圖片網址方式，不改正式遊戲資料格式、規則或引擎。

189,359,703 bytes單檔依原有2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗。
包含114張內嵌PNG bytes與尺寸、來源劇本／武將／頁碼、57位切換、曹丕預設、六種新增畫面，
四種模式、滑鼠拉桿／縮放／同步拖曳／原生上限／重設，57列來源表開頭／捲動／關閉／重開，
搬移單檔離線重開與390×844兩倍像素密度。本次三次載入皆約3.8秒，各114個暫存網址全部釋放、來源區塊全部移除，JS錯誤0、外部請求0。
原圖、原生輸出及兩張總覽已目視，正常頁面另抽看劉禪、關興第2稿並排、張苞窄視窗細節及捲動後來源視窗。
完整性核對57張原圖、69次原生bytes、四劇本名單、原始資源及解析器SHA、既有趙雲／鄧艾像素證據、王平原圖紅色樣本、歷史快照與重建一致。
既有45張選稿、十張已同意城池PNG、城池頁及全圖頁未變。57張仍是150頁中的小樣，細部筆觸屬重繪詮釋，未宣稱逐張接受或正式接線完成。

來源：[第六批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-6-plan.json)、
[原圖與劇本映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-6.json)、
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-fiftyseven-r1.json)、
[第六批檢視、原圖樣本與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-6.json)、
[45張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/fortyfive-r1-snapshot.json)、
[七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-fiftyseven-r1.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-fiftyseven-r1.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-6-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-6.json`準備。
提示詞與合併入口為 `portraits/prepare-batch6.py`及 `integrate-batch6.py`；頁面、抽驗與核對沿用 `build.py`、`qa-check.cjs`及 `verify.py`。
沿用既有Go、Python與Playwright映像，以1000:1000及有界斷網容器執行，原資料唯讀，只寫既有頭像工作目錄。

### 頭像擴充至69張與帽飾細修

使用者「ok go」延續既定原圖參考HD方向。先保存57張頁、程式、metadata、收據、抽樣畫面及57張原生選稿，
共85個檔案雜湊，再由既有解析器取下列12張64×64原圖。選稿皆1254×1254原生PNG。

| 新增武將 | 來源劇本索引 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|---|
| 公孫瓚 | 0 | 12 | 79 | 第1稿 |
| 劉表 | 0 | 31 | 89 | 第1稿 |
| 劉璋 | 0 | 29 | 91 | 第1稿 |
| 馬騰 | 0 | 26 | 107 | 第1稿 |
| 韓遂 | 0 | 41 | 122 | 第1稿 |
| 孔融 | 0 | 13 | 21 | 第1稿 |
| 田豐 | 0 | 89 | 14 | 第1稿 |
| 審配 | 0 | 61 | 92 | 第2稿 |
| 馬良 | 0 | 91 | 88 | 第1稿 |
| 馬謖 | 0 | 90 | 109 | 第1稿 |
| 諸葛瑾 | 0 | 59 | 97 | 第2稿 |
| 陳登 | 0 | 83 | 16 | 第1稿 |

十二張均來自劇本索引0，名字、武將ID與頁碼連同來源劇本核對。
各沿自己的黑底、臉部構圖、神情、冠帽及主要配色重畫。公孫瓚保留側轉與綠色帽塊，田豐保留側面，
韓遂保留張口與紅色帽塊，馬良保留淡灰眉線，馬謖保留灰帽淡色線。
審配首稿將灰帽明暗畫成獨立白色方塊，第2稿回到連續灰綠布面，保留淡色上緣及金邊。
諸葛瑾首稿的淡金帽飾過大，第2稿縮小輪廓，保留原圖額前淡色／金色細節。
兩張修正都是美術詮釋，不宣稱帽飾材質或意義已證實。兩者選第2稿，其餘十張選第1稿，未選稿皆保存。
本批12張首稿與兩張修正稿共14次輸出，連同前批共83次。69張選稿逐byte複製，沒有後製或將放大圖當成HD原生稿。

[當時69張本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-sixtynine-r1.html)預設公孫瓚，
保留原版／HD／並排／拉桿、同步縮放／拖曳／重設、HD 1:1及固定工具列，來源視窗擴充至69列。
沿用既有分開內嵌PNG與暫存圖片網址方式，不改正式遊戲資料格式、規則或引擎。

226,002,096 bytes單檔依原有2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗。
包含138張內嵌PNG bytes與尺寸、來源劇本／武將／頁碼、69位切換、公孫瓚預設、六種新增畫面，
四種模式、滑鼠拉桿／縮放／同步拖曳／原生上限／重設，69列來源表開頭／捲動／關閉／重開，
搬移單檔離線重開與390×844兩倍像素密度。本次三次載入約4.5、4.5、4.8秒，
各138個暫存網址全部釋放、來源區塊全部移除，JS錯誤0、外部請求0。
原圖、原生輸出及兩張總覽已目視，正常頁面另抽看審配與諸葛瑾第2稿並排，以及捲動後窄視窗來源表。
完整性核對69張原圖、83次原生bytes、四劇本名單、原始資源及解析器SHA、既有趙雲／鄧艾像素證據、王平原圖紅色樣本、歷史快照與重建一致。
HTML SHA-256為 `466efc943c6444f9d968070dda7b97fdeb4f9af61611dd73dc0b97881f40368f`。
既有57張選稿、十張已同意城池PNG、城池頁及全圖頁未變。69張仍是150頁中的小樣，細部筆觸屬重繪詮釋，未宣稱逐張接受或正式接線完成。

來源：[第七批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-7-plan.json)、
[原圖與劇本映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-7.json)、
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-sixtynine-r1.json)、
[第七批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-7.json)、
[57張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/fiftyseven-r1-snapshot.json)、
[七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-sixtynine-r1.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-sixtynine-r1.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-7-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-7.json`準備。
提示詞與合併入口為 `portraits/prepare-batch7.py`及 `integrate-batch7.py`；頁面、抽驗與核對沿用 `build.py`、`qa-check.cjs`及 `verify.py`。
沿用既有Go、Python與Playwright映像，以1000:1000及有界斷網容器執行，原資料唯讀，只寫既有頭像工作目錄。

### 頭像擴充至81張與原色色塊細修

使用者「ok go」延續既定原圖參考HD方向。先保存69張頁、程式、metadata、收據、抽樣畫面及69張原生選稿，
共97個檔案雜湊，再由既有解析器取下列12張64×64原圖。選稿皆1254×1254原生PNG。

| 新增武將 | 來源劇本索引 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|---|
| 顏良 | 0 | 44 | 8 | 第1稿 |
| 文醜 | 0 | 92 | 12 | 第1稿 |
| 張繡 | 0 | 22 | 54 | 第1稿 |
| 張魯 | 0 | 23 | 46 | 第1稿 |
| 馬岱 | 0 | 24 | 84 | 第1稿 |
| 嚴顏 | 0 | 49 | 1 | 第2稿 |
| 程普 | 0 | 86 | 37 | 第2稿 |
| 徐盛 | 0 | 66 | 33 | 第1稿 |
| 丁奉 | 0 | 87 | 13 | 第1稿 |
| 曹洪 | 0 | 67 | 100 | 第2稿 |
| 曹真 | 0 | 68 | 73 | 第1稿 |
| 滿寵 | 0 | 97 | 3 | 第3稿 |

十二張均來自劇本索引0，名字、武將ID與頁碼連同來源劇本核對。
各沿自己的黑底、臉部構圖、神情、冠盔及主要配色重畫。張魯保留傾斜笑臉，嚴顏保留灰白眉鬚，
文醜保留側轉，程普保留斜向頭部，丁奉保留紅色額前細節，曹真保留側轉笑容。

- 嚴顏原圖確有小紅點，第1稿畫成兩個過大的平面方塊；第2稿縮小並增加表面明暗，保留原有紅色。
- 程普原圖確有綠色區，第1稿擴大為明確捲紋；第2稿縮小。細部形狀仍屬重繪詮釋，不宣稱已解出符號。
- 曹洪原圖確有藍色塊，第1稿過大且平整；第2稿縮小，加入細筆觸及明暗，原有藍色與鄰接淡色亮面保留。
- 滿寵原圖下唇有暖色亮點，前兩稿畫成獨立橙色橫條；第3稿融回嘴唇及皮膚明暗，保留暖光與原有深色髭線。

前三張選第2稿，滿寵選第3稿，其餘八張選第1稿。色點的材質與符號意義仍未知。
本批12張首稿與五張修正稿共17次輸出，連同前批共100次。81張選稿逐byte複製，未後製，也未把放大參考圖當成原生HD稿。
原圖、17次原生輸出及兩張總覽已目視，正常頁面另抽看四張修正稿並排、嚴顏窄視窗細節與捲動後來源視窗。

[現行81張本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index.html)預設顏良，
保留原版／HD／並排／拉桿、同步縮放／拖曳／重設、HD 1:1及固定工具列，來源視窗擴充至81列。
沿用分開內嵌PNG與暫存圖片網址方式，不改正式遊戲資料格式、規則或引擎。

265,009,434 bytes單檔依原有2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗。
含162張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、81位切換、顏良預設、六種新增畫面、
四種模式、滑鼠拉桿／縮放／拖曳／原生上限／重設、81列來源表開頭／捲動／關閉／重開、
搬移單檔離線重開及390×844兩倍像素密度。三次載入約3.3、5.3、4.1秒，
各162個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。
完整性核對81張原圖、100次原生bytes、四劇本名單、原始資源／解析器SHA、
趙雲／鄧艾／王平既有像素證據、69張及更早快照與重建一致。
HTML SHA-256為 `00bda6411edeab10a7ff60aa633d5cffb780acb0bd4fba52fd6da9b95bd1b162`。
既有69張選稿、十張已同意城池PNG、城池頁及全圖頁未變。81張仍是150頁中的小樣，未宣稱逐張接受或正式接線完成。

來源：[第八批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-8-plan.json)、
[原圖與劇本映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-8.json)、
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set.json)、
[第八批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-8.json)、
[69張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/sixtynine-r1-snapshot.json)、
[七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-8-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-8.json`準備。
本批提示詞與合併入口為 `portraits/prepare-batch8.py`與 `integrate-batch8.py`，快照入口為 `snapshot-sixtynine-r1.py`。
原生複製沿 `faithful/tooling/save-portraits.py originals-batch-8.json`，三份複製清單入口在第八批紀錄。
頁面、抽驗與核對沿用 `build.py`、`qa-check.cjs`與 `verify.py`，建頁日期改由選稿清單取得。
沿用既有Go、Python與Playwright映像，以1000:1000及有界斷網容器執行，原資料唯讀，只寫既有頭像工作目錄。

## 許昌本機原圖修繪（清晰度未通過）

使用者否決前兩款SVG重繪，指出失去原版特色。本輪停止擴展通用城池與新畫風，
先做 [許昌原圖修繪對照（歷史）](../../workplace/hd-redraw-local-20261104/faithful/index-local-retouch.html)。
這是本機模型輔助的原圖細節修繪試作，不宣稱人類畫師完成手繪，也不把技術檢查當作美術驗收。

- 直接讀原版 `KYOGRF.DAT`、`MMAP.MAP`／`MMAP.MDL`、`GAMEPAL.BRG`及劇本。
  以既有library／state解碼，KYOGRF實際15頁；許昌82的 `KindHigh=6`。
  前批只保留位置與路線，建築由通用SVG依據點類型合成，未遵循原版建築圖像；這是使用者否決後的修正重點。
- 地圖城池使用世界 `(206,114)` 周圍7×7格的112×112原圖；據點景觀直接取KYOGRF第6頁96×96原圖。
  分別建立原圖四倍參照與保護區；模型以384×384處理，再合成回448×448／384×384。
  兩份原圖不互相替代，也不按大城／中城類型猜建築。
- 原圖深色像素、強邊緣與圖框設保護區，直接保留原圖四倍參照放大版本。
  其餘部分由本機Stable Diffusion 1.5原圖條件修繪加入有界表面細節，再沿原圖的廣域色塊合成。
  固定seed `20261104`、strength 0.25、16個排程步，實際各4個低噪聲步；沒有重擲後挑圖。
  原圖、保護區、原始模型輸出、合成小樣與參數全部留存。
- 低記憶體工具使用既有CPU image與鎖版Diffusers0.35.2純Python套件層，UNet線性層動態int8，
  卷積float32，文字編碼後釋放CLIP。這是本機工具的近似計算方式，不作原版規則或像素parity的聲明。
- 公開模型只在未掛原圖的下載容器取得，固定模型commit及檔案SHA；本機推論、建頁與QA均 `--network none`。
  模型與套件來源在同目錄 `models/sources.json`／`diffusers-wheel.json`，工具版本在 `tooling/runtime.json`。
  沒有呼叫imagegen、Sites或雲端模型視覺工具，也沒有上傳原圖或截圖。
- 比較頁只顯示許昌兩份來源，可切原版／修繪／並排／拉桿及查看輪廓保護區。
  rich2式滾輪、減號／倍率重設／加號、拖曳平移與雙圖同步保留。全圖及其他城池尚未擴展。
- 9項瀏覽器抽驗通過，JS錯誤0、外部請求0。五個原始輸入SHA未變，兩張小樣共198,092個保護區像素
  與原圖四倍參照完全相同，單通道改動上限29。此處的參照是Lanczos四倍圖，不能把它冒稱原版原始像素對拍。
  使用者後續指出仍糊，未達HD。保護區約51%／65%的像素直接保留柔化四倍參照，
  模型又只以384×384處理；這限制了原輪廓重描與細節重建，不能繼續拿像素相同作為HD品質的成功證據。

歷史重生順序：`faithful/extract.go` → `prepare_inputs.py` → 本機 `paint.py` → `build-local-retouch.py`。
來源／數值檢查在 `faithful/qa/verify.py`，瀏覽器抽驗在 `qa/check.cjs`；收據入口為
`faithful/qa/verification.json`、`qa/pixels.json`及 `qa/result.json`。此分支仍為DRAFT，未進正式美術或遊戲路徑。

## 城池與地圖SVG重繪（已被使用者否決）

使用者於本機超解析比較後指定「重新定義HD」，以原版為依據，先重新手繪城池與地圖。
當時入口是 [本機城池與地圖重繪比較](../../workplace/hd-redraw-local-20261104/index.html)。
這批採SVG線條繪製手繪風草稿，不宣稱人類畫師逐筆手繪；也沒有呼叫雲端影像生成服務。

- 地圖底稿取松崗DOS/V版 `MMAP.MAP`／`MMAP.MDL`／`GAMEPAL.BRG`，沿既有
  `internal/assets/world`、`palette` 解碼；原世界384×256格、256圖塊。
  據點沿 `SINARIO.DAT` 第一劇本的 `state.City`，Big5姓名由專案解碼器處理。
  城池資料與254條道路不新增玩法或推測規則。
- 第一批區域為世界格 `(180,96)` 至 `(238,145)` 的半開區間，原尺寸928×784。
  涵蓋許昌、虎牢關、舞陽、汝南等13個據點，27條道路與區域相交。
  新地圖將據點錨點設在原座標的格中心，路線使用原版逐格路徑，沒有自己補一條連線。
- 新繪屋頂、瓦片、城牆、垛口、城門、山石、樹木及水紋均為本機SVG幾何。
  水域外形由原圖顏色描出，陸地筆觸用原色分佈作視覺詮釋，尚非已證實的地形分類；
  城池新造型與顯示尺寸也是草案，不冒稱原版建築形制或pixel parity。
- 四張近看圖：許昌82／大城0、滎陽63／中城1、舞陽87／小城2、虎牢關72／關隘3。
  左側保留原版9×9格資源合成，右側放大新城池造型供確認細節，不以此頁證明城牆可走性、
  正式命中區、戰術城門位置或戰鬥效果。
- A青綠彩繪與B古紙淡彩各五張，共10張新SVG；五張原版PNG與新SVG都內嵌單檔HTML。
  可切原版／重繪／並排、畫風、據點名、縮放、拖曳及同步檢視；畫風切換保留焦點，工具列固定。
  SVG可直接保存，QA另在本機輸出3倍PNG，放大仍可檢查新繪線條。
- 核對入口為同目錄 `qa/geometry.json`、`qa/result.json`、`qa/verification.json`；
  12項瀏覽器抽驗通過，JS錯誤0、外部請求0。資料核對只證明來源、原圖裁切、據點及路徑一致，
  圖像輸出只證明能渲染，不代表使用者接受畫風。使用者後續已明確否決A、B，
  認為失去原版特色；兩者未接正式遊戲，也不再作待選或定版素材。

本機重生順序：`extract.go` 用既有Go parser產生原圖／資料，`art.py`繪製SVG與比較資料，
`build.py`嵌入單檔頁面，`qa/verify.py`核對來源與幾何，`qa/check.cjs`抽驗瀏覽器並輸出PNG。
所有工作都在有界、1000:1000、`--network none` 容器；原始輸入唯讀，只有指定輸出可寫。
沒有原圖、重繪或本機截圖上傳OpenAI，沒有更新先前雲端網站。

## 前版本機超解析比較（歷史）

使用者於先前網站交付後要求：原版與HD都以原版圖像為基礎，檔案留在本機。
前版入口是 [本機原版／HD比較HTML](../../workplace/hd-original-local-20261104/index.html)，
可直接雙擊開啟，不需要伺服器。來源、28張圖片及程式均嵌入同一HTML；
`manifest.json` 保存輸入、輸出、原版來源與雜湊，`build.py`可在離線容器重生。

- 四個場景是真實原版dosgolem截圖：戰略地圖、汝南守城、道路野戰與勢力視窗。
  來源分別沿playtest127／128及spec274；沒有用remake畫面冒充原版。
- 曹操、關羽、諸葛亮直接取自 `KAOGRF.DAT` 第50／114／125頁。
  映射依劇本一武將16／40／60的 `General.Portrait`，實際圖庫150張，原尺寸64×64。
  初始能力也讀自原版劇本，沒有沿用先前新繪網頁的示範數字。
- 使用者回饋前版差異不足，希望基於原版更好看。新增「原畫HD」，預設以曹操頭像並排開啟。
  由本機Real-ESRGAN重建4倍色階與輪廓，再作輕度色彩、明暗及銳化；不使用換臉模型。
  頭像使用general-x4v3，場景使用animevideov3，均為官方v0.2.5.0權重。
  原圖維持原始PNG bytes，並留輪廓／柔和放大兩種前版效果及 `index-v1.html` 供追溯。
- 戰術側欄及底列、戰略頂列與提示框、勢力清單區域從原圖以整數倍合成，
  提示框中的頭像保留強化。區域及例外在manifest，座標依既有battlelayout／listlayout／statusbox；
  這是prototype合成選擇，沒有新增原版行為斷言。地圖內文字及場景細節仍待使用者確認。
- 可切原圖、HD、並排、拉桿及三種強化。參考rich2加入減號／倍率／加號、160ms過渡、
  滾輪拉遠拉近、拖曳平移及雙圖同步；點倍率回到3.4的預設檢視，範圍1–5、步進0.2。
  此處倍率是圖片檢視尺度，3.4對應適合視窗；正式HD應由世界鏡頭縮放地圖，介面保持固定。
  截圖沒有畫面外地圖，不能用此頁代證世界鏡頭或正式HD完成。
  另保留AI與武將編輯選項示意。
  固定截圖不偽造即時部隊動畫，編輯不改圖片、不寫正式存檔。
- 超解析及建頁均在 `--network none` 容器內；公開模型下載獨立容器，僅掛模型目錄，
  未掛任何原版圖像。來源固定在官方tag／release，URL、SHA-256及BSD授權保存於
  本機 `models/sources.json`／`models/LICENSE`，重生入口為 `models/download.py` 與 `superresolve.py`。
  模型架構沿 [官方v0.3.0 SRVGG](https://github.com/xinntao/Real-ESRGAN/blob/v0.3.0/realesrgan/archs/srvgg_arch.py)，
  權重及設定沿 [官方推論入口](https://github.com/xinntao/Real-ESRGAN/blob/v0.3.0/inference_realesrgan.py)。
  沒有呼叫imagegen或Sites，原圖與瀏覽器截圖
  也沒有交給 `view_image` 或其他雲端影像工具。HTML的CSP禁止外部連線及資源。
- 本機收據入口 `workplace/hd-original-local-20261104/qa/result.json`，
  包含file網址開啟、28張嵌入圖雜湊、畫面／頭像、拉桿、縮放、編輯與手機瀏覽抽驗。
  `qa/beauty-verification.json`保存原圖雜湊、指定介面區域精確保留、與前版像素差異，
  `beauty/superresolution.json`保存離線推論及各圖模型；差異數字不作好看或原版parity的證據。
  畫面截圖僅存本機，尚待使用者親自確認強化視覺。

## 先前新繪概念網站（歷史，素材與交付方式已被取代）

2026-11-04 使用者要求先做網頁確認概念。私人預覽：
[臥龍傳HD概念預覽](https://wolong-hd-concept-anr2.cy-wang.chatgpt.site)。
來源位於 `workplace/hd-concept-web-20261104/`，網站身分保存於該目錄
`.openai/hosting.json`。此版只保留歷史定位；依最新裁定不再同步、更新或使用它交付原版素材。

- 常駐工具列有原版／HD外觀、縮放、AI輔助、武將編輯及原有操作入口。
  原版風格是新素材的程式示意，不是原版畫面的對拍收據。
- 曹操、關羽、諸葛亮使用內建imagegen新製頭像，提示詞及雜湊在網站來源的
  `asset-provenance.json`。古典寫實畫風先供比較，尚未定版。
- 戰略與攻城地圖是原創SVG。可點選城池、試編成兩段確認與右鍵回父窗，
  文字選單保留上下推動反白。
- 武將編輯可調基本能力、三種適性與委任腳本，取消不套用，確認後改示範資料；
  四槽保存只寫瀏覽器localStorage。沒有讀寫原版或正式原生存檔。
- 攻城動畫演示雙方統帥1–2倍、部隊移動與架梯進度。這是無引擎的概念模擬，
  不代表已實作正式AI尋路、戰鬥規則、真正架梯耗時或戰局公平性。
  概念頁允許即時切換；正式遊戲的開關時機仍未定案。
- 可在頁面記錄意見並下載文字，不會自動傳送給其他人。

Sites私人網站只上傳新製素材與網站程式，未放入原版EXE、遊戲圖片、字型或資料檔。
桌面及390×844瀏覽器已抽驗；這不代證正式手機版完成。
本機 `.qa/qa-result.json`、畫面、部署收據與來源commit保存此份概念的定位，
不把此版新繪畫面作為目前HD原圖強化的實作基準。

## 未解

| 項目 | 等級 | 是否阻塞正式實作 |
|---|---|---|
| AI切換時機及移動／架梯接線 | 待設計／證據審查；曲線已選雙方1–2倍 | 是，先查現有計時與做prototype |
| 新入口分組及窄視窗安排 | 未決 | 是，不能直接縮小已選A |
| HD正式素材產製與接線 | 五城十張已獲同意；三十五區試繪覆蓋原圖全部面積與192據點，孤立屋頂已修；81張頭像小樣未宣稱逐張接受 | 畫風不阻塞；局部配色／水紋／邊線、其餘頭像、資產映射與正式接線仍待完成 |
| HD世界鏡頭與固定介面 | 本機原型可依世界座標切換全圖／已繪範圍／本區、滑鼠縮放與平移，工具列固定 | 是，正式Ebiten視野、焦點、邊界及命中座標仍待接線與驗證 |
| 編輯套用時機、取消及存檔衍生欄位 | 待設計／證據審查 | 是，先做單人垂直鏈草稿 |
