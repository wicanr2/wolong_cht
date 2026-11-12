# 275：HD城池地圖重繪、工具列、頭像、戰術AI輔助與武將編輯

**狀態：DRAFT（整體）。** 頭像、顯示偏好、地圖縮放、色盤0地景／橫幅／五組城池景觀與精簡自動收起工具列已接桌面，限定驗證範圍見下方小節；其餘功能仍待完成。

HD美術方向已由使用者確認。DRAFT保留給其他色盤／場景的HD素材、局部美術修整、AI與編輯器。

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
- 2026-11-11使用者明確指定「工具列要單排」。擴充後所有入口維持單排，
  原版／HD及縮放常駐；分組雙排只保留歷史比較，不再待選。
- 2026-11-11使用者選「預設原版，記住之後的選擇」，HD素材缺漏時顯示原圖。
- 2026-11-11使用者選戰略地圖50%～200%、預設100%、每次25%，支援滾輪，介面固定。
- 頭像也要HD。
- HD版要參考rich2支援拉遠／拉近。地圖視野與介面大小分離，工具列提供減號、
  倍率重設、加號，支援滑鼠滾輪；目前比較頁的圖片縮放不代表正式地圖鏡頭完成。
- 另加工具列AI選項，作用在戰術畫面，以統帥加成改善小兵到城門及架梯兵效率。
  使用者已選擇同時加成找路效率、單兵走速與架梯動作時間。
  我方與敵方都套用，各依自己的統帥加成。
  使用者已選最高2倍：統帥1為1倍、15為2倍，中間線性增加，僅開啟AI輔助時套用。
- 另加武將能力編輯，包含隱藏的委任、陸戰、海戰、攻城欄位。

初始A為九項入口。加入AI輔助及武將編輯後為十一項，縮放分為三個實體按鈕，共十三鈕。
使用者已選擴充比較稿的單排方向。此選擇不代表編輯器、AI預設及正式功能接線均已定案。
正式遊戲現為 [spec266](266-mouse-only-desktop-toolbar.md) 的七項加頭像顯示切換與三個縮放鈕，共十一鈕單排；
AI與編輯器仍待接線。
本輪原版捲軸修正另見 [spec274](274-original-list-scrollbar-drag.md)。

## 已查證的工程入口

### 第五段：已同意城池景觀接入情報窗（CONFORMED，限定五組景觀抽驗）

2026-11-12接續使用者「ok 繼續囉」，沿已同意五城十圖中的五張景觀，接入正常據點情報窗。
不新增畫風、圖像、規則、視窗配置或存檔欄位。原版模式及缺圖回原圖的裁定沿用。

- 已證實的原資料契約沿re50與spec23：KYOGRF實際15頁，96×96，City.KindHigh選頁，與City.Kind獨立。
  本機既有manifest.json的scene樣本對應許昌6、滎陽14、舞陽2、虎牢關5、汝南0。
  每頁可供其他沿原版資料共用該景觀的據點使用，不以武將／城池ID或城池類型自行選圖。
- 素材為五張已同意的1254×1254原生PNG。許昌依hd-generated/manifest.json的user_approval，
  其餘四城依batch-2.json的user_approval；不重新生成或修改圖片。
  [接線前快照](../../workplace/scratch/ui-windows-20261104/hd-location-baseline.json)保存前段38檔來源包、
  現有manifest／originals SHA、五頁映射與選稿SHA。
- `-hd-locations`指向既有faithful目錄，只讀manifest.json的kind=scene樣本及originals.json，
  沿hd.file、hd_record.file／width／height／sha256／native_unmodified、city.scene_page取選稿。
  比對原映射與實際KYOGRF項目數，0為合法頁碼，重複頁、越界、缺欄、來源不一致或壞圖則整組停用。
  PNG由既有hd-generated子目錄讀取，驗證目錄與檔案不逃逸；沿HD共用檔案／尺寸／SHA限制。
  不新增資產散布格式，不讀未選attempts或地圖城池樣本。
- 只在HD模式與來源色盤0啟用；其他季節與液晶、未選頁及執行中讀圖失敗沿原Location。
  按需解碼，GPU快取最多兩張；切換顯示與讀檔不改遊戲資料。沿既有hd_display偏好，首次原版。
- 直接畫到既有原生合成面，顯示位置仍為(16,288)、96×96，情報窗外框、數值與文字保持原位置。
  先繪製的窗框在景觀後方，後繪製的文字／游標／視窗依原順序覆蓋。
  原頭像64×64繪製改用同一個可指定尺寸的原生矩形入口，保留原本結果與遮擋。
- 驗收包含五頁來源與0／缺頁／越界／重複／壞圖，原生細節與後畫遮擋、96與64尺寸，
  正常新遊戲經小地圖定位後逐一點五城，原版／HD及液晶往返、缺素材與兩種存檔未被顯示操作重寫。
  只作已接景觀與抽驗範圍的完成聲明，不外推其餘十頁、所有HD視窗或原生平台驗收。

本段先列READY再實作。入口為[景觀選稿載入](../../internal/assets/library/hdlocations.go)、
[GPU景觀／快取](../../cmd/wlgame/hdlocation.go)、[原生矩形共用繪製](../../cmd/wlgame/hdportrait.go)與
[正常據點情報窗](../../cmd/wlgame/cityinfo.go)。
[資料測試](../../internal/assets/library/hdlocations_test.go)涵蓋實際KYOGRF15頁、五張1254×1254選稿、
頁0／越界／缺頁／原映射不符／重複／逃逸／壞圖；
[GPU測試](../../cmd/wlgame/hdlocation_test.go)驗證96與64尺寸、原生單像素細節、後畫遮擋、兩頁快取及執行中缺圖回退。
獨立程序以`WOLONG_HD_LOCATION_GPU=1 go test -count=1 ./cmd/wlgame -run '^TestHDLocationGPUNativeSizeOcclusionAndCache$'`
在Xvfb容器執行。既有頭像GPU測試另跑通過，受影響冷測與vet通過。
記錄分為[素材](../../workplace/scratch/ui-windows-20261104/hd-location-asset-tests.log)、
[介面](../../workplace/scratch/ui-windows-20261104/hd-location-ui-tests.log)及
[獨立GPU](../../workplace/scratch/ui-windows-20261104/hd-location-gpu.log)。

[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/hd-location-gui.json)由
[流程工具](../../workplace/scratch/ui-windows-20261104/hd-location-gui.py)從首頁走新遊戲、用原版小地圖定位後逐一點城。
145個功能與同步檢查點包含重複收展同步，不當作145個獨立功能。
五城頁碼6／14／2／5／0、氾水關共用5、宛未選10回原圖、快取重載、液晶往返、正常讀檔及缺素材重開通過。
SAVE.DAT與SAVE-slot1.wlsave的SHA未被景觀／顯示操作改寫。沒有座標注入、direct-entry或改原版規則。

[13張實際畫面頁](../../workplace/scratch/ui-windows-20261104/hd-location-game-review.html)由
[建頁工具](../../workplace/scratch/ui-windows-20261104/build-hd-location-review.py)重生，附192×192截圖裁區的兩倍檢視，
裁區不增加圖片細節；[頁面工具](../../workplace/scratch/ui-windows-20261104/qa-hd-location-review.cjs)執行
[四項離線抽驗](../../workplace/scratch/ui-windows-20261104/hd-location-review-qa.json)，切換、搬移及窄視窗通過，外部請求與JS錯誤皆0。
[接線收據](../../workplace/scratch/ui-windows-20261104/hd-location-integration.json)與
[43檔來源覆蓋包](../../workplace/scratch/ui-windows-20261104/hd-location-source-r1.tar.gz)
保存程式、五張既有選稿、binary、GUI與測試SHA；五張PNG及既有美術頁未改。
[收尾核對](../../workplace/scratch/ui-windows-20261104/hd-location-final-checks.json)另驗來源覆蓋包、素材、離線頁面、文件入口及擁有權。
CONFORMED只指本段已同意的五頁與上述抽驗，未宣稱其餘十頁或全部HD介面完成。

### 第四段：HD橫幅與精簡自動收起工具列（CONFORMED，限定視窗模式抽驗）

使用者指出頂部仍非HD，確認範圍為「龍紋、標題與日期的橫幅」，另要求原版可縮小地圖、
工具列小一點並可隱藏。使用者已選「靠上緣自動顯示，離開自動收起」。
沿[28／6高度操作稿](../../workplace/scratch/ui-windows-20261104/toolbar-compact-hide-prototype.html)及
[選擇紀錄](../../workplace/scratch/ui-windows-20261104/toolbar-compact-hide-choice.json)，維持單排與畫布外置。

- 原版縮放已存在，本輪以HD關閉、HD素材仍可用的正常新遊戲抽驗50%／75%／100%／200%、滾輪及小地圖。
  [13項收據](../../workplace/scratch/ui-windows-20261104/original-map-zoom-gui.json)與
  [流程工具](../../workplace/scratch/ui-windows-20261104/original-map-zoom-gui.py)保存結果；小地圖沿左鍵開、右鍵關。
- 橫幅來源為Library.Banner的ICONGRF原圖640×32，八組色盤與來源SHA由
  [擷取工具](../../workplace/scratch/ui-windows-20261104/banner-source-probe.go)及
  [來源收據](../../workplace/scratch/ui-windows-20261104/banner-source-probe.json)保存。
  內建imagegen以色盤0原圖重繪，提示詞與原生來源沿
  [生成紀錄](../../workplace/scratch/ui-windows-20261104/banner-generation-r1.json)及
  [完整提示詞](../../workplace/scratch/ui-windows-20261104/banner-prompt-r1.txt)。
  原生2172×724 PNG保留bytes，顯示只取選稿記錄的內容矩形；白邊不拉進遊戲。
- `-hd-banner`指定含banner-selection.json的目錄。欄位沿file／width／height／sha256／native_unmodified，
  另以content_rect與palette_bank指定原生內容及色盤。沿HD素材共用路徑／尺寸／SHA驗證；錯誤回原圖。
  同一HD模式控制地圖、頭像及橫幅。只對已提供的色盤0啟用，液晶／其他季節沿原圖。
  原圖標題字形保留；文字造型細節仍屬重繪詮釋，不宣稱像素parity或新增字形校訂。
- 背景直接畫到高解析合成面，日期數字仍沿原字模與右緣496／528／560即時繪入，不能烘焙日期。
  原版四個視窗開關與第五個保留熱區仍沿x336起、每格32、y0–31；不因重繪搬移熱區。
- 工具列展開高度由40改28，按鈕上下各留2，文字垂直置中；收起高度6，頂緣把手喚出。
  初始收起，指標離開工具列後自動收起。有任一按鍵按住、工具列擁有手勢、原版選單／拖曳擁有輸入時不切換。
  必須消費放開幀，不能把工具列手勢交給背景。沒有新增快捷鍵與存檔欄位。
- 遊戲畫布仍640×400，工具列偏移只屬桌面層。一般視窗隨收展調整客戶區高度以保留遊戲顯示倍率；
  全螢幕／最大化沿固定視窗邊界重新適配。捕捉游標重映射並重建原始位移基準，避免收展被當作地圖捲動。
- 驗收：橫幅索引／壞圖／內容矩形、原生細節與後畫日期／視窗遮擋、原版回退；
  工具列單排四語、收展／完整手勢／捕捉座標與正常新遊戲、兩模式縮放及小地圖開關。
  素材與HTML留本機，不部署，不更改原版操作與世界規則。

本段先列READY再實作。入口為[橫幅載入](../../internal/assets/library/hdbanner.go)、
[原生繪製](../../cmd/wlgame/hdbanner.go)與[桌面收展](../../cmd/wlgame/desktopautohide.go)。
[選稿／壞圖測試](../../internal/assets/library/hdbanner_test.go)、
[GPU內容裁切／日期覆蓋](../../cmd/wlgame/hdbanner_test.go)及
[完整手勢／18次視窗適配座標檢查](../../cmd/wlgame/desktopautohide_test.go)通過。
獨立GPU測試用`WOLONG_HD_BANNER_GPU=1 go test -count=1 ./cmd/wlgame -run '^TestHDBannerGPUCropNativeDetailAndDateOrder$'`
於Xvfb容器執行。冷測與GPU記錄見
[素材](../../workplace/scratch/ui-windows-20261104/hd-banner-asset-tests.log)、
[介面](../../workplace/scratch/ui-windows-20261104/hd-banner-ui-tests.log)及
[GPU](../../workplace/scratch/ui-windows-20261104/hd-banner-gpu.log)，vet通過。

[正常滑鼠收據](../../workplace/scratch/ui-windows-20261104/banner-toolbar-gui.json)由
[流程工具](../../workplace/scratch/ui-windows-20261104/banner-toolbar-gui.py)從首頁進新遊戲，
通過96個功能與同步檢查點，含重複收展同步，不宣稱96個獨立功能。
實際視窗1280×812／856，遊戲均維持兩倍顯示；展開／收起不移動鏡頭，不把按住／放開交給背景。
原版模式縮放／小地圖、HD原生橫幅／頭像、模式／液晶往返、正常存讀檔與缺素材重開通過。
SAVE.DAT與原生槽檔的SHA未被顯示或工具列重寫。全螢幕／最大化僅核對座標轉換，未代證原生平台GUI。

[12張實際畫面頁](../../workplace/scratch/ui-windows-20261104/banner-toolbar-game-review.html)由
[建頁工具](../../workplace/scratch/ui-windows-20261104/build-banner-toolbar-review.py)重生，
[四項斷網抽驗](../../workplace/scratch/ui-windows-20261104/banner-toolbar-review-qa.json)由
[頁面工具](../../workplace/scratch/ui-windows-20261104/qa-banner-toolbar-review.cjs)執行。
頁內連結[原生PNG](../../workplace/hd-redraw-local-20261104/faithful/banner-hd-r1.png)與提示詞，
選稿紀錄為[banner-selection.json](../../workplace/hd-redraw-local-20261104/faithful/banner-selection.json)。
[接線收據](../../workplace/scratch/ui-windows-20261104/banner-toolbar-integration.json)及
[38檔來源覆蓋包](../../workplace/scratch/ui-windows-20261104/banner-toolbar-source-r1.tar.gz)
固定程式、工具、素材、binary與截圖SHA。重繪依既定原圖畫風，未宣稱本張美術已個別接受。
[收尾核對](../../workplace/scratch/ui-windows-20261104/banner-toolbar-final-checks.json)另驗38檔來源覆蓋包、
原生橫幅／既有57張地景、收據、文件入口及擁有權。

### 第三段正式接線：HD大地圖（CONFORMED，限定抽驗範圍）

使用者在第二段完成後回覆「ok go」，接續已展示的HD大地圖接線。本段沿既定畫風、
工具列模式、固定介面與50%～200%鏡頭，不改構圖、規則、路線、命中或存檔。
使用既有35區、22接縫的本機選稿，不重繪PNG，不以選稿接線代稱逐區美術接受。

- 輸入為faithful/world的regions.json及其record指向的選稿JSON／原生PNG。
  索引SHA-256為`27f4b3d291e31cd6f608062db06108ee96e76836bcc736aec2cdc47317b89748`。
  各項沿source_world_pixel_rect、core_world_pixel_rect、file、width、height、sha256及feather_world_pixels。
  世界範圍沿原資料6144×4096；core必須位於source內，區域不得重疊，接縫依清單順序覆蓋。
  原圖來源、資料SHA與parser沿本頁世界試繪清單；extract.go明確以色盤0重建原圖。
- 已證實的工程事實：原版MMAP只存無所屬中心，歸屬、首都、軍團與災害來自即時世界狀態。
  沿world/citymark.go與spec53／74共用標記分派。HD圖上另疊原版即時標記：
  歸屬／角旗採原圖塊差異的最小包圍框加一像素邊界；首都與軍團用原MCH透明遮罩。
  首都中心有軍團時仍抑制首都疊圖。雲層、災害及格游標在地景後繪製，不烘焙進背景。
  這是HD呈現差異，原模式仍走既有整幅RenderMarked，不改其像素輸出。
- 使用`-hd-world`指向現有目錄。載入驗證整組索引、唯一ID、世界矩形、SHA、尺寸與原生標記；
  不讀attempts。JSON上限4MiB，PNG上限16MiB、4096邊及4,194,304像素；沿既有HD路徑防逃逸契約。
  索引不合法整組停用；執行中讀圖失敗只讓該區沿原圖。接縫失敗沿下層，不蓋黑塊。
- 直接把可見原生圖片畫到高解析合成面，依世界矩形裁切；避免先壓到640×400再放大。
  接縫透明度沿本機預覽的8個世界像素邊界，來源PNG bytes不變。
  圖像按需解碼，GPU快取最多16張與128MiB；固定UI及頭像繼續保留原繪圖順序。
- 同一原版／HD按鈕控制地圖與頭像，首次原版、沿用hd_display偏好，單排幾何不變。
  無HD素材時停用切換且保留偏好。色盤0以外沿原版地圖，液晶沿原版地圖與頭像；
  不替尚未提供的季節／液晶HD變體自行調色。戰術畫面本段不換圖。
- 驗收包含索引／矩形／壞圖與fallback、實際57選稿、原版標記及首都抑制、
  原生GPU細節與遮擋、三倍率城池點選、全圖邊界／接縫、正常模式切換、讀檔與缺圖；
  實際截圖留本機。HD位置／筆觸依既有美術限制，不稱原版圖像parity或全平台完成。

[接線前來源索引](../../workplace/scratch/ui-windows-20261104/hd-world-baseline.json)
固定前段26份來源覆蓋包、二進位收據與素材索引，後續收據不得改寫前段的受測程式身分。

本段先列READY契約再接線。入口為[選稿載入](../../internal/assets/library/hdworld.go)、
[原生世界合成](../../cmd/wlgame/hdworld.go)及[共用原版標記](../../internal/assets/world/citymark.go)。
[索引／57張實際選稿測試](../../internal/assets/library/hdworld_test.go)、
[原版合成一致性測試](../../internal/assets/world/hdmarkers_test.go)與
[幾何／GPU測試](../../cmd/wlgame/hdworld_test.go)通過；GPU測試以
`WOLONG_HD_WORLD_GPU=1 go test -count=1 ./cmd/wlgame -run '^TestHDWorldGPUKeepsDetailAndFallback$'`
在獨立Xvfb容器程序執行，另設WOLONG_HD_ORIG為合法原版資料目錄。
測試保留原生單像素條紋，後畫UI覆蓋正確，單區圖片在載入後消失時顯示原地圖；
所有非0季節色盤與液晶拒用此組地景。這些是受控fixture，不代替正常玩家驗收。

[54項正常滑鼠抽驗](../../workplace/scratch/ui-windows-20261104/hd-world-gui.json)由
[流程工具](../../workplace/scratch/ui-windows-20261104/hd-world-gui.py)從首頁走新遊戲，
涵蓋原版／HD、液晶往返、三倍率點同城、小地圖定位、全圖西北／東南邊界、正常存讀檔及缺素材重開。
trace同步實際Draw的鏡頭與倍率；沒有注入狀態。SAVE.DAT與SAVE-slot1.wlsave兩檔的SHA在切換／縮放前後未變。
最大抽樣GPU快取75,491,140 bytes，低於128MiB；這是抽樣值，不稱所有玩家視野的效能上限。
正常GUI未抽到移動中的軍團，原版歸屬／首都／軍團疊圖以12種狀態在原地形上合成後逐byte核對，
不外推完整行軍流程或所有HD場景驗收。

[本機實際遊戲頁](../../workplace/scratch/ui-windows-20261104/hd-world-game-review.html)內嵌12張截圖，由
[建頁工具](../../workplace/scratch/ui-windows-20261104/build-hd-world-review.py)重生，
[四項斷網抽驗](../../workplace/scratch/ui-windows-20261104/hd-world-review-qa.json)沿
[頁面工具](../../workplace/scratch/ui-windows-20261104/qa-hd-world-review.cjs)執行。
[接線收據](../../workplace/scratch/ui-windows-20261104/hd-world-integration.json)保存32份來源、57張原生選稿、
binary、測試及截圖SHA；[來源覆蓋包](../../workplace/scratch/ui-windows-20261104/hd-world-source-r1.tar.gz)
配合收據的base_git_head可重建。冷測記錄分為[資料／標記](../../workplace/scratch/ui-windows-20261104/hd-world-core-tests.log)、
[介面](../../workplace/scratch/ui-windows-20261104/hd-world-ui-tests.log)與
[獨立GPU](../../workplace/scratch/ui-windows-20261104/hd-world-gpu.log)，vet通過。
原生尺寸有限，超過素材1:1的放大仍為插值；未生成更高倍率圖片。
CONFORMED限於本段色盤0選稿與上述抽驗，其他季節、液晶HD、據點景觀HD、戰術HD及局部美術仍待處理。
[收尾核對](../../workplace/scratch/ui-windows-20261104/hd-world-final-checks.json)驗證來源覆蓋包、57張原生PNG、
現行收據、文件入口及UID/GID；前段素材頁未改。

### 第二段正式接線：戰略地圖縮放（CONFORMED，限定抽驗範圍）

2026-11-11使用者選「50%～200%（建議）：兼顧視野與辨識」。沿已展示方案，
預設100%、每次25%，工具列減號／倍率重設／加號常駐單排，滾輪以游標為焦點。
只縮放戰略地圖、城池、軍團、災害與格游標；橫幅、視窗、文字、頭像與工具列尺寸固定。
現階段使用原版地圖，HD地圖素材接線另作一段。此為使用者授權的remake呈現差異。

- 已證實工程輸入：原世界384×256格、每格16像素，地圖視窗640×368，頂端32像素橫幅。
  資料入口為Library.RenderWorldMarked、mapPickTileAt、scrollMapPixels及drawMinimapViewBox。
  原版100%座標契約沿spec149，分別取整，不把camSub加入命中；本段保留該分支。
- 50%～200%採共同世界像素原點，畫面點反算為floor((原點＋畫面點／倍率)／16)。
  地圖來源依視野擴張，城池標記、軍團與災害合成後一同縮放，固定UI最後覆蓋。
  城池命中仍只接受資料登記的一格，不把整座圖形當熱區。
- 原點夾在世界邊界內；非100%原點取最近整數像素，滾動保留小數累積。
  回100%將原點吸附最近16像素格，焦點最多偏移8個世界像素，以保留原版繪製與命中。
  小地圖視野框沿原圖邊角調整寬高，點擊定位依當前可見視野置中。
- 工具列加入三個實體鈕後共十一鈕，保持單排與四語標籤。戰術、標題、選單與模態視窗
  停用縮放。滾輪只在未被UI遮擋的地圖區作用；有按鍵手勢時不消費滾輪。
- 倍率僅為執行中呈現狀態，不新增遊戲存檔欄位或偏好；新遊戲／讀檔回100%。
  此段不改規則、AI、計時、路徑與美術PNG。
- 驗收：100%既有命中／捲動測試；50%／75%／125%／200%的反算、邊界、焦點、
  小數累積及小地圖框；四語單排尺寸；正常滑鼠新遊戲、工具列與滾輪、同城不同倍率點選、
  模態遮擋、讀檔復原及HD頭像固定呈現。只聲明抽驗範圍，不稱原版新縮放parity。

[倍率比較原型](../../workplace/scratch/ui-windows-20261104/world-zoom-range-prototype.html)及
[原資料鏡頭探針](../../workplace/scratch/ui-windows-20261104/world-zoom-camera-probe.json)保存78項座標檢查及來源SHA。
[取樣診斷](../../workplace/scratch/ui-windows-20261104/world-zoom-pixel-phase.json)記錄瀏覽器與Go在縮小時的像素邊界取樣差：
兩側各自符合來源取樣40／40，跨實作僅29／40相同，不當成精確渲染通過。
[頭像接線程式快照索引](../../workplace/scratch/ui-windows-20261104/world-zoom-baseline.json)連同
[十八檔來源覆蓋包](../../workplace/scratch/ui-windows-20261104/hd-head-source-before-zoom.tar.gz)
固定縮放前的已驗收程式，避免後續修改使頭像收據失去可重生輸入。

本段先列READY契約再實作，已通過[座標與工具列冷測](../../workplace/scratch/ui-windows-20261104/world-zoom-tests.log)、
[獨立GPU抽驗](../../workplace/scratch/ui-windows-20261104/world-zoom-gpu.log)及vet。
192座城五種倍率共960次反算通過；GPU以原始全世界圖另合成劇本雲層，核對五種倍率35個像素與固定橫幅邊界。
[31項正常滑鼠抽驗](../../workplace/scratch/ui-windows-20261104/world-zoom-gui.json)由
[正常流程工具](../../workplace/scratch/ui-windows-20261104/world-zoom-gui.py)沿首頁新遊戲執行，
包括50%／100%／200%點選同一登記城、模態遮擋、滾輪焦點、小地圖定位及讀檔回100%。
只讀trace提供同步，不注入遊戲狀態。100%與其他倍率間的格吸附依上述契約保留。

[正式遊戲截圖頁](../../workplace/scratch/ui-windows-20261104/world-zoom-game-review.html)由
[建頁工具](../../workplace/scratch/ui-windows-20261104/build-world-zoom-review.py)重生，
[四項斷網頁面抽驗](../../workplace/scratch/ui-windows-20261104/world-zoom-review-qa.json)由
[頁面抽驗工具](../../workplace/scratch/ui-windows-20261104/qa-world-zoom-review.cjs)核對。
[接線收據](../../workplace/scratch/ui-windows-20261104/world-zoom-integration.json)保存26份來源、binary、截圖、測試及素材頁SHA，
[來源覆蓋包](../../workplace/scratch/ui-windows-20261104/world-zoom-source-r1.tar.gz)配合收據的base_git_head可重建本段。
正式入口為[世界鏡頭](../../cmd/wlgame/worldzoom.go)、[座標／GPU測試](../../cmd/wlgame/worldzoom_test.go)、
[工具列](../../cmd/wlgame/desktop.go)與[四語縮放文字](../../translations/desktop-zoom.json)。
CONFORMED僅指此原版地圖呈現與抽驗，不代表HD大地圖、全部戰場／視窗或原生平台驗收。
[收尾核對](../../workplace/scratch/ui-windows-20261104/world-zoom-final-checks.json)另驗來源覆蓋包、素材頁未變、
本機入口與擁有權。GUI存檔雜湊抽查的是SAVE.DAT；原生執行狀態的復原由正常讀檔抽驗涵蓋。

### 第一段正式接線：HD頭像與顯示偏好（CONFORMED，限定抽驗範圍）

2026-11-11使用者選「預設原版，記住之後的選擇」。本段僅授權頭像接線：
原七項工具列先加入顯示模式，維持單排；縮放、AI及編輯器仍待各自接線，不加入無作用按鈕。
此為分段實作，尚非十一項入口全部完成。首次原版；玩家切換後寫桌面偏好，原生與remake遊戲存檔不變。
缺圖使用原圖，無可用HD資料時不啟用切換，也不覆寫玩家原先保存的選擇。

已證實的工程事實：原版畫布640×400；全部八處Portrait呼叫均畫64×64頭像。
僅在64×64目標換圖會先失去原生細節。Ebiten 2.9.9的`FinalScreenDrawer`與
`DefaultDrawFinalScreen`契約取自本機鎖定模組的`run.go`及`gameforui.go`，不新增硬體時序推論。
頭像頁碼沿實際KAOGRF項目Count，不以版本字串或武將ID推導。

- 輸入沿本機既有`selection.json`的records與逐張JSON之portrait_page、file、width、height、sha256；
  使用`-hd-portraits`指定現有目錄，不另定正式散布包格式。只讀取選稿，不讀attempts。
- 載入先驗證整組索引、唯一頁碼、實際圖庫界限、檔名限制、PNG尺寸及SHA；失敗整組停用。
  圖像解碼與GPU快取按需要載入，最多保留八張GPU頭像；未選頁及其後讀取失敗回原圖。
  現有選稿為彩色PNG，無液晶變體；液晶模式使用原圖且停用HD切換，保留玩家HD偏好，切回彩色自動沿用。
- JSON限1MiB，PNG限16MiB，最大邊4096、像素總數4,194,304；路徑不得逃出指定目錄。
- 桌面HD繪圖以保序合成處理：遇到HD頭像時，先將此前的原版畫布按整數／顯示倍率貼到高解析合成面，
  再貼原生頭像；其後的視窗、文字及游標依原繪製順序覆蓋。不得在最後無條件把頭像蓋到所有視窗上。
  合成面倍率上限8以限制記憶體；此為呈現限制，不改原版邏輯座標。原模式沿DefaultDrawFinalScreen。
- 對拍`-shot`及錄影仍走原版game與640×400畫布；HD不成為原版oracle。
- 四語工具列文字進translations；桌面偏好version 1新增可省略的hd_display布林，缺欄位為false。
  原偏好檔仍能讀取，只有玩家操作才寫檔。AI與遊戲資料完全不讀此偏好。
- 驗收包含壞索引、頁碼0／越界、重複頁、路徑逃逸、檔案損壞、原生尺寸、保序遮擋、舊偏好、
  正常滑鼠新遊戲君主卡／世界HUD及切換後重開。HD頭像不作美術逐張接受或整款HD完成聲明。

本段按READY契約實作後，素材索引、原生保序合成、偏好、新遊戲君主卡與勢力情報窗通過驗證。
CONFORMED僅限此範圍。命名、編成、軍團與戰術對話已接共用入口，但未逐一走完正常GUI；
不外推這些場景的完整驗收。當時尚未接世界縮放與HD全圖，現況依第二／第三段；AI與編輯器維持DRAFT。

入口為[素材載入](../../internal/assets/library/hdportraits.go)、[GPU合成](../../cmd/wlgame/hdportrait.go)、
[共用頭像](../../cmd/wlgame/talkscene.go)、[桌面工具列](../../cmd/wlgame/desktop.go)與[偏好](../../cmd/wlgame/preferences.go)。
切換文字在[四語桌面顯示詞表](../../translations/desktop-display.json)。
[素材測試](../../internal/assets/library/hdportraits_test.go)及[合成／偏好測試](../../cmd/wlgame/hdportrait_test.go)包含實際150頁圖庫、
146張選稿、頁碼0／越界、重複頁、檔名逃逸、SHA與尺寸、舊偏好。GPU像素測試須以
`WOLONG_HD_GPU_PROBE=1 go test -count=1 ./cmd/wlgame -run '^TestHDCompositeKeepsNativeDetailAndLaterOcclusion$'`
在獨立容器程序與Xvfb下執行，不能在RunGame之前讀GPU像素；其合成fixture不代替正常玩家驗證。

[實際遊戲對照頁](../../workplace/scratch/ui-windows-20261104/hd-head-game-review.html)由
[建頁工具](../../workplace/scratch/ui-windows-20261104/build-hd-head-review.py)重生，七張截圖來自
[正常滑鼠抽驗](../../workplace/scratch/ui-windows-20261104/hd-head-gui.json)及其
[重播工具](../../workplace/scratch/ui-windows-20261104/hd-head-gui.py)，共29項檢查，包含液晶回原圖及切回彩色恢復HD偏好。
[接線收據](../../workplace/scratch/ui-windows-20261104/hd-head-integration.json)保存程式與測試雜湊、來源、二進位與限制；
[本機頁面抽驗](../../workplace/scratch/ui-windows-20261104/hd-head-review-qa.json)由
[頁面抽驗工具](../../workplace/scratch/ui-windows-20261104/qa-hd-head-review.cjs)核對單檔搬移離線開啟。
此頁是remake正常操作截圖，不是原版同狀態對拍。啟動時以`-hd-portraits`指向本機含selection.json的選稿目錄。

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

1. 已確認：工具列外置且單排、原版／HD與縮放常駐、頭像HD、戰術AI選項、武將能力編輯。
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
頭像、地圖縮放與色盤0地景已按上述分段接線；未更改存檔格式、戰鬥數值或公開素材政策。

## 比較稿與權利

本機比較稿：`workplace/scratch/ui-windows-20261104/hd-toolbar-prototype.html`，
原始A為 `hd-toolbar-A.svg`，B為 `hd-toolbar-B.svg`；生成程式同目錄。
採用A的配置選擇記錄於本頁。原版衍生畫面與自備字型形成的草稿留在本機，
原圖及雜湊由本輪manifest保存，不公開散布。

速度曲線概念稿：`workplace/scratch/ui-windows-20261104/ai-speed-prototype.html`，
比較最高1.5／2／3倍，統帥1為1倍、15為上限的線性提案。
這只是無障礙路線及架梯進度的概念動畫，未接遊戲、不作實戰效能或原版證據。
使用者已選2倍上限及上述線性曲線；其餘兩條不採用。這是新增設計，不稱為原版證據。

## 擴充工具列操作比較稿

既定九項入口加入AI輔助與武將編輯後，新增[本機單排／雙排比較稿](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1.html)。
原版／HD及縮放仍常駐，共11個功能入口，縮放拆為減號、倍率重設及加號，因此有13個實體按鈕。
分組雙排把設定、存讀檔、返回、下一步與離開獨立一列；緊湊單排保留高度，較窄視窗需左右捲動。
2026-11-11使用者回覆「工具列要單排」，已採單排，原建議雙排不採用，仍保留歷史比較。
[選擇紀錄](../../workplace/scratch/ui-windows-20261104/hd-toolbar-single-choice.json)保存原話、範圍及選定畫面雜湊；
[選擇前快照](../../workplace/scratch/ui-windows-20261104/hd-toolbar-before-single-choice.json)保存十一份檔案，
[單排參考畫面](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1-before-single-choice-single.png)保留不覆寫。
現行頁面開啟及重開皆預設單排，窄視窗仍維持一列並可左右捲動；不重開顯示模式常駐的決定。
單排確認稿另加零尺寸畫布防護，不在尺寸為零時改寫鏡頭座標。
[受控前後對照](../../workplace/scratch/ui-windows-20261104/hd-toolbar-single-size-control.json)確認舊稿座標失效、新稿保留有效值；
十項抽驗另核對每張截圖前的地圖像素、單排預設、窄視窗與重開。

本稿使用既有中原原圖／HD選稿及曹操原圖／HD頭像，沿其原始映射，不新增美術。
原版／HD切換會更換實際素材，支援工具列縮放、滾輪及拖曳，依顯示像素密度限制原生細節。
武將編輯取既有曹操資料，武術9、統帥13、政治13、委任類型2、攻城6、陸戰4、海戰0。
可用滑鼠增減、選委任類型、取消草稿、套用本頁示範及還原；只有姓名欄接受文字輸入。
姓名與數值均只留在頁面記憶體，重開即還原，沒有寫入正式或原生存檔。
AI開關僅預覽狀態及既定雙方1至2倍曲線，沒有加速戰鬥；其餘入口明示未接線。
本稿初始HD／AI關閉只供展示；正式顯示已選首次原版，AI設定邊界沿本頁決策樹。這份操作稿未接入遊戲或更動原版視窗。

[來源與建置清單](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1.json)、
[建置入口](../../workplace/scratch/ui-windows-20261104/build-hd-toolbar-expanded-r1.py)、
[頁面模板](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1.html.in)與
[互動程式](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1.js)保存在既有本機工作目錄。
[十項離線抽驗](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1-qa.json)由
[抽驗入口](../../workplace/scratch/ui-windows-20261104/qa-hd-toolbar-expanded-r1.cjs)執行，
涵蓋1280／960寬及640寬兩倍像素密度、實際分列、素材切換、縮放、編輯取消／套用／還原、
AI預覽、其餘入口、停用快捷鍵與單檔搬移重開。JS錯誤及外部請求皆0。
目視曾發現雙排被flex-basis擠在同列，舊抽驗只查切換狀態，無法證明分列；
[修正紀錄](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1-row-fix.json)保留修正前頁面、工具、截圖及收據。
修正後新增按鈕群座標與工具列高度檢查，並重新目視。
[完整性收據](../../workplace/scratch/ui-windows-20261104/hd-toolbar-expanded-r1-verification.json)由
[驗證入口](../../workplace/scratch/ui-windows-20261104/verify-hd-toolbar-expanded-r1.py)核對來源、重建、快照及既有三個美術頁未變。
工具列單排方向已確認。這些仍是操作稿的驗證，不代表正式工具列、編輯器、戰術加成或遊戲流程完成。

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

[81張歷史本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-eightyone-r1.html)預設顏良，
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
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-eightyone-r1.json)、
[第八批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-8.json)、
[69張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/sixtynine-r1-snapshot.json)、
[七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-eightyone-r1.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-eightyone-r1.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-8-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-8.json`準備。
本批提示詞與合併入口為 `portraits/prepare-batch8.py`與 `integrate-batch8.py`，快照入口為 `snapshot-sixtynine-r1.py`。
原生複製沿 `faithful/tooling/save-portraits.py originals-batch-8.json`，三份複製清單入口在第八批紀錄。
頁面、抽驗與核對沿用 `build.py`、`qa-check.cjs`與 `verify.py`，建頁日期改由選稿清單取得。
沿用既有Go、Python與Playwright映像，以1000:1000及有界斷網容器執行，原資料唯讀，只寫既有頭像工作目錄。

### 頭像擴充至93張與眼部／頭盔細修

使用者「繼續囉」接續已確認的原圖參考HD方向。先保存81張頁、程式、metadata、收據、
抽樣畫面與81張選稿，共109個檔案雜湊，再取下列12張64×64原圖。選稿皆1254×1254原生PNG。

| 新增武將 | 來源劇本索引 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|---|
| 凌統 | 0 | 104 | 80 | 第2稿 |
| 文聘 | 0 | 93 | 0 | 第2稿 |
| 廖化 | 0 | 102 | 58 | 第1稿 |
| 周倉 | 0 | 115 | 142 | 第1稿 |
| 蔣欽 | 0 | 58 | 71 | 第1稿 |
| 陳武 | 0 | 84 | 4 | 第2稿 |
| 朱桓 | 1 | 51 | 75 | 第1稿 |
| 曹休 | 1 | 61 | 38 | 第1稿 |
| 曹純 | 0 | 69 | 101 | 第2稿 |
| 夏侯霸 | 1 | 30 | 29 | 第1稿 |
| 孟達 | 0 | 98 | 30 | 第1稿 |
| 張任 | 0 | 79 | 95 | 第2稿 |

三張來自劇本索引1，其餘九張來自索引0，名字、武將ID與頁碼連同劇本核對。
各沿自己的黑底、臉部構圖、神情、冠盔與主要配色重畫。周倉與夏侯霸保留傾斜頭部，
文聘與朱桓保留側轉，孟達保留原有藍色帽前細節。

- 文聘首稿有重複眼睛與孤立面頰色點，第2稿修正眼部結構，將面頰色點融回膚色明暗；原有藍色側片保留。
- 陳武首稿多出大金色符號、額前菱形與棋盤衣紋，第2稿改回連續灰色盔面、少量金色亮點與連續衣料明暗。
- 凌統首稿張口過大且下巴多出細綠線，第2稿縮小張口、移除細線，保留連接頭盔的原有綠布。
- 曹純首稿將帽前淡色細節畫成鳥形，第2稿改為小型抽象亮面。原圖下唇座標35,40確有RGB(211,0,0)單一紅點，選稿保留為小紅色筆觸，意義維持未知。
- 張任首稿帽前金色碎片過散，第2稿融成連續淡金亮面，保留原有淡色／金色區；不宣稱已解出符號意義。

五張選第2稿，其餘七張選第1稿。12張首稿及五張修正稿全部保存，本批17次、累計117次。
93張選稿逐byte複製，未後製，最近鄰放大圖只作參考。細部髮鬚、材質與帽前線條仍屬重繪詮釋。
原圖、17次原生輸出與兩張總覽已目視，正常頁面抽看文聘並排、曹純窄視窗細節與捲動後來源視窗。

[93張歷史本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-ninetythree-r1.html)預設凌統，保留原版／HD／並排／拉桿、
同步縮放／拖曳／重設、HD 1:1及固定工具列，來源視窗擴充至93列。沿用分開內嵌PNG與暫存圖片網址，
不改正式遊戲資料格式、規則或引擎。

303,867,988 bytes單檔依原有2GB／1 CPU、30秒導覽及90秒外層完成七項斷網抽驗。
含186張內嵌PNG bytes／尺寸、來源劇本／武將／頁碼、93位切換、凌統預設、六種新增畫面、
四種模式、滑鼠拉桿／縮放／拖曳／原生上限／重設、93列來源表開頭／捲動／關閉／重開、
搬移單檔離線重開及390×844兩倍像素密度。三次載入約9.7、7.2、5.4秒，
各186個暫存網址全釋放、來源區塊全移除，JS錯誤0、外部請求0。
完整性核對93張原圖、117次原生bytes、四劇本名單、原始資源／解析器SHA、
趙雲／鄧艾／王平既有像素證據、曹純紅點、81張及更早快照與重建一致。
HTML SHA-256為 `b3caa5a2ae8f7c10c54844890c5f99f52d197877fd53999a458636832c0ac075`。
既有81張選稿、十張已同意城池PNG、城池頁及全圖頁未變。93張仍是150頁中的小樣，未宣稱逐張接受或正式接線完成。

來源：[第九批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-9-plan.json)、
[原圖與劇本映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-9.json)、
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-ninetythree-r1.json)、
[第九批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-9.json)、
[81張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/eightyone-r1-snapshot.json)、
[七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-ninetythree-r1.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-ninetythree-r1.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-9-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-9.json`準備。
本批提示詞與合併入口為 `portraits/prepare-batch9.py`與 `integrate-batch9.py`，快照入口為 `snapshot-eightyone-r1.py`。
原生複製沿 `faithful/tooling/save-portraits.py originals-batch-9.json`，兩份複製清單入口在第九批紀錄。
頁面、抽驗與核對沿用 `build.py`、`qa-check.cjs`與 `verify.py`。
沿既有Go、Python與Playwright映像，以1000:1000及有界斷網容器執行，原資料唯讀，只寫既有頭像工作目錄。

### 頭像擴充至105張與帽前細節／衣領修正

使用者「ok 繼續囉」接續原圖參考HD方向。先保存93張頁、程式、metadata、收據、抽樣畫面與93張選稿，
共121個檔案雜湊，再取下列12張64×64原圖。選稿皆1254×1254原生PNG。

| 新增武將 | 來源劇本索引 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|---|
| 曹叡 | 1 | 3 | 81 | 第1稿 |
| 曹彰 | 1 | 4 | 123 | 第1稿 |
| 關索 | 1 | 33 | 72 | 第2稿 |
| 王雙 | 1 | 22 | 35 | 第2稿 |
| 紀靈 | 0 | 46 | 27 | 第1稿 |
| 凌操 | 0 | 103 | 59 | 第2稿 |
| 孫桓 | 0 | 17 | 106 | 第2稿 |
| 孫韶 | 0 | 20 | 70 | 第2稿 |
| 劉封 | 0 | 33 | 141 | 第2稿 |
| 蔡瑁 | 0 | 54 | 11 | 第2稿 |
| 張紘 | 0 | 76 | 32 | 第1稿 |
| 孫乾 | 0 | 71 | 98 | 第2稿 |

四張來自劇本索引1，其餘八張來自索引0，名字、武將ID與頁碼連同劇本核對。
曹彰使用劇本1的原字；劇本0、武將14的曹漳也使用123頁，不據此改名或將一頁算成兩張。
各沿自己的黑底、頭部構圖、神情、冠盔與原色重畫。曹叡保留斜向側轉，
張紘保留寬綠布帽，曹彰與劉封保留原有藍色帽前細節。

- 關索第2稿收斂首稿獨立橢圓／菱形牌飾，改以連續金屬亮面處理，保留原有淡金區。
- 王雙第2稿移除新增灰色圓點群及凸起白色鉚釘，保留小範圍灰白反光、綠布與原有藍色細節。
- 凌操第2稿縮小過大的紅色紋樣，保留原有紅色；孫桓第2稿縮小金色紋樣，細部形狀及符號意義仍屬詮釋。
- 孫韶第2稿縮小新增金色捲紋，回到小型棕金色細節；孫乾第2稿將平面金色方塊融入布面明暗。
- 劉封第2稿將衣領棋盤紋改回連續暖灰／金色明暗，原有藍色帽前區與淡色盔面保留。
- 蔡瑁原圖帽前兩側有14個RGB(81,146,65)綠色像素，包含38,0與53,3。第2稿縮小新增捲紋，保留綠色及中央紅色；色塊意義維持未知。

八張選第2稿，其餘四張選第1稿。12張首稿與八張修正稿全部保存，本批20次、累計137次。
105張選稿逐byte複製，未後製，最近鄰放大圖只作參考。原圖、20次輸出及兩張總覽已目視，
正常頁面另抽看王雙並排、孫乾窄視窗細節及捲動後來源視窗；不宣稱美術逐像素相同或逐張接受。

[105張歷史本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-onehundredfive-r1.html)預設曹叡，保留四種比較模式、同步縮放／拖曳／重設、
HD 1:1及固定工具列，來源視窗擴充至105列。342,877,215 bytes單檔含210張內嵌PNG。
初次完整抽驗於開頁時出現「Page crashed」，相同HTML的獨立導覽約5.2秒通過。
驗證程式改讓PNG檢查子程序先結束，再啟動瀏覽器；該子程序RSS約1.09GB，啟動瀏覽器前的父程序約79MB。
同一2GB／1 CPU、30秒導覽及90秒外層通過七項完整斷網抽驗，未放寬原有斷言。

抽驗包含210張PNG bytes／尺寸、來源劇本／武將／頁碼、105位切換、曹叡預設、六種新增畫面、
四種模式及滑鼠操作、105列來源視窗開頭／捲動／關閉／重開、搬移單檔重開與390×844兩倍像素密度。
三次載入約5.1、6.5、15.2秒，各210個暫存網址全釋放，來源區塊全移除，JS錯誤0、外部請求0。

完整性核對的內嵌重建也曾以137結束。建頁程式原先同時保留全部base64區塊、合併字串與完整HTML，
現改為分段寫入暫存檔，再原子替換。原有1GB／1 CPU、60秒完整性核對重跑通過。
HTML SHA-256仍為 `0eb9fa5eddea1e9e871df3fe853fbc4f8a839bc96056d8b628f23db60959f690`。
兩次失敗均保留原工具與紀錄，未取得原失敗的OOM計數，不把OOM機制寫成已證實原因。
修正限於驗證工具，頁面程式、HTML bytes、原生素材與正式引擎未改。

完整性核對105張原圖、137次原生bytes、四劇本名單、原始資源／解析器SHA、既有趙雲／鄧艾／王平／曹純像素證據、
蔡瑁綠色樣本、共用123頁的原字、93張及更早快照與重建一致。既有93張選稿、十張已同意城池PNG、城池頁及全圖頁未變。
105張仍是150頁中的小樣，其餘頭像、全圖局部美術及正式接線仍待完成。

來源：[第十批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-10-plan.json)、[原圖與劇本映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-10.json)、
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set-onehundredfive-r1.json)、[第十批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-10.json)、
[93張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/ninetythree-r1-snapshot.json)、[七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-onehundredfive-r1.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-onehundredfive-r1.json)、[驗證工具檢視](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch10-browser-review-onehundredfive-r1.json)、
[失敗頁與工具快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch10-loading-before-snapshot.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-10-plan.json`，最近鄰參考由同目錄 `reference-portraits.py originals-batch-10.json`準備。
本批提示詞與合併入口為 `portraits/prepare-batch10.py`與 `integrate-batch10.py`，快照入口為 `snapshot-ninetythree-r1.py`。
原生複製沿 `faithful/tooling/save-portraits.py originals-batch-10.json`，兩份複製清單在第十批紀錄。
沿 `build.py`、`qa-check.cjs`及其 `qa-assets.cjs`子程序、`verify.py`重跑；最後由 `record-batch10-browser-review.py`記錄檢視。
獨立導覽診斷為 `batch10-navigation-probe.cjs`，不代替完整抽驗。全部沿既有映像、有界斷網、1000:1000容器，原資料唯讀。

### 頭像擴充至117張與側臉／色塊細修

使用者「繼續囉」接續原圖參考HD方向。先保存105張頁、程式、metadata、收據、抽樣畫面、105張選稿與前批失敗證據，
共145個檔案雜湊，再由既有解析器取下列12張64×64原圖。十二張皆來自劇本索引0，選稿皆1254×1254原生PNG。

| 新增武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---|---|---|
| 袁尚 | 1 | 90 | 第1稿 |
| 袁譚 | 4 | 52 | 第2稿 |
| 袁胤 | 0 | 5 | 第2稿 |
| 王朗 | 5 | 28 | 第1稿 |
| 李傕 | 27 | 66 | 第2稿 |
| 郭汜 | 6 | 67 | 第2稿 |
| 韓玄 | 7 | 7 | 第1稿 |
| 金旋 | 9 | 2 | 第2稿 |
| 嚴白虎 | 10 | 117 | 第1稿 |
| 公孫越 | 11 | 77 | 第1稿 |
| 黃祖 | 51 | 53 | 第2稿 |
| 劉繇 | 34 | 87 | 第1稿 |

名字、劇本、武將ID及頁碼一起核對。袁胤的ID確為0，單檔重開也抽驗其正常選取，不把0當成缺少人物。
各沿原圖保留黑底、頭部構圖、神情、髭鬚與原色；黃祖保留禿額，李傕保留灰白上唇細節。

- 袁譚第2稿收斂首稿過寬、過正面的臉部比例，較靠近原圖斜向側臉，幾何仍屬重繪詮釋。
- 袁胤第2稿減少新增的完整白色矩形框，保留綠色區及淡色亮邊。原圖綠色位於32–37、2–4，共18個RGB(130,162,97)像素。
- 李傕第2稿將孤立的大片白色口部亮條融回髭線／嘴唇明暗，保留原有灰白細節，不宣稱已判定每個淡色像素的材質。
- 郭汜第2稿縮小過大的平面藍色菱形。原圖帽前13個樣本含四個RGB(48,65,211)與九個RGB(0,32,97)，藍色保留，符號意義維持未知。
- 金旋第2稿移除凸起灰色圓釘，減少側盔規則綠金橫紋，保留灰白亮面、綠色底調與後方綠布。
- 黃祖第2稿移除首稿肩部過大的黑色重複菱格，恢復連續金色／淡色明暗及原有邊緣細節，禿額與面容保留。

六張選第2稿，其餘六張選第1稿。12張首稿與六張修正稿均保存，本批18次、累計155次。
原生PNG逐byte複製且未後製；最近鄰放大圖及兩張總覽只供參考／檢視。
原圖、18次輸出與兩張總覽已目視，正常頁面另抽看袁譚／李傕並排、黃祖窄視窗細節與捲動後來源表。

[現行117張本機對照](../../workplace/hd-redraw-local-20261104/faithful/portraits/index-117-r1.html)預設袁尚，保留四種模式、固定工具列、
同步縮放／拖曳／重設及HD 1:1，來源視窗擴充至117列。單檔380,383,085 bytes，含234張PNG。
前兩次完整抽驗都碰到90秒外層逾時，分別停在搬移後獨立重開及窄視窗載入；原HTML、工具、錯誤與第二次階段紀錄均保留。

驗證工具改逐張串流讀取內嵌PNG，保留每張bytes／尺寸及來源映射核對，子程序結束後才開瀏覽器。
各檢視依序關閉再重開；117次選圖由固定100毫秒等待改成等選項／標題實際更新，原有標題斷言保留。
同一2GB／1 CPU、30秒導覽、90秒外層與七項完整檢查通過，沒有減少117位切換。
階段紀錄約52秒到達第七項，三次載入約6.8、11.1、14.1秒，各234個暫存網址全釋放，來源區塊全移除。
JS錯誤0、外部請求0，包含四種模式、滑鼠拉桿／縮放／拖曳／上限／重設、117列來源視窗開頭／捲動／關閉／重開、
單檔搬移後選取袁胤，以及390×844兩倍像素密度與黃祖細節。
這些是本次工具環境的觀察，沒有隔離各項修正的個別效果，也不宣稱原失敗由OOM造成。

HTML SHA-256仍為 `3ad3a843bf5d368454b1c03e8218f001288d9edc7217284af7af44a57f01390b`，與初次失敗頁相同。
app、原生圖及正式引擎未因驗證工具調整而改動。完整性核對以原有1GB／1 CPU、60秒重跑通過，
含117張原圖、155次原生bytes、四劇本名單、來源／解析器SHA、既有像素證據、
袁胤綠色及郭汜藍色樣本、ID0、105張及更早快照與重建一致。
既有105張選稿、十張已同意城池PNG、城池頁及全圖頁未變。117張仍為150頁中的小樣，未宣稱逐張接受或正式接線完成。

來源：[第十一批計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-11-plan.json)、[原圖與映射](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-11.json)、
[提示詞與選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set.json)、[本批檢視與未選稿](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-11.json)、
[105張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/onehundredfive-r1-snapshot.json)、[七項抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-117-r1.json)、
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-117-r1.json)、[工具檢視](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch11-browser-review.json)、
[第一次逾時](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch11-browser-before-snapshot.json)、[第二次逾時與階段紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch11-browser-attempt2.json)。
原圖入口為 `faithful/tooling/extract-portraits.go batch-11-plan.json`，參考圖由同目錄 `reference-portraits.py originals-batch-11.json`準備。
提示詞與合併入口為 `portraits/prepare-batch11.py`、`integrate-batch11.py`，快照入口為 `snapshot-onehundredfive-r1.py`。
原生複製沿 `faithful/tooling/save-portraits.py originals-batch-11.json`，兩份複製清單在本批紀錄。
沿 `build.py`、`qa-check.cjs`及 `qa-assets.cjs`、`verify.py`重跑，最後執行 `record-batch11-browser-review.py`。
本機現行失敗輸出在雜湊核對封存後清除，歷史檔保留。工具皆沿既有映像、有界斷網及1000:1000容器，原資料唯讀。

### 頭像擴充至120張

新增劉勳、孫瑜與劉琦，沿原版64×64圖像，以內建imagegen重繪1254×1254原生PNG。
三張均來自劇本索引0，名字、武將ID與圖庫頁碼由原始劇本及圖庫解碼核對。

| 武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---:|---:|---|
| 劉勳 | 8 | 96 | r2 |
| 孫瑜 | 21 | 60 | r1 |
| 劉琦 | 28 | 78 | r1 |

劉勳第1稿將帽前原有綠色區塊畫成凸起寶石，第2稿改成平面綠色與淡色反光，保留綠色區塊。
孫瑜保留朝左側臉、外露耳部、深色鬍鬚與原有紅色區域；劉琦保留無鬚長臉、金褐頭飾、灰白衣領與藍色肩部。
四次原生輸出全部保留且未後製。細部材質及色塊含義仍屬重繪詮釋／未知，未作逐張接受聲明。

[本批紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-12.json)、[來源計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-12-plan.json)、
[原圖與參考圖](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-12.json)、[117張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/117-r1-snapshot.json)、
[完整提示詞](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-12-prompts.json)、[劉勳細修提示詞](../../workplace/hd-redraw-local-20261104/faithful/portraits/portrait-96-generation-r2.json)、
[原生輸出路徑](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-12-native-outputs.json)、[合併入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/integrate-batch12.py)、[文件登錄入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/record-batch12.py)及
[檢視用縮圖](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-12-contact-sheet.png)保存在既有本機目錄。縮圖只作檢視，原生PNG未改。

目前120張選稿、159次原生輸出、12份來源包。七項[離線抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-120-r1.json)通過，
包含240張內嵌PNG、全部120張切換、三張新增畫面、原版／HD／拉桿、滑鼠縮放及拖曳、
來源視窗、單檔搬移重開、窄視窗與兩倍像素密度。
沿原來2GB／1 CPU、30秒導覽與90秒外層執行，app.js及build.py未變。
[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-120-r1.json)由[本批驗證入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/verify-batch12.py)重生，
核對原生圖與工具輸出雜湊、全部來源及放大參考圖、117張快照398個檔案、城池／全圖頁與HTML重建一致。
現行verify.py亦使用同一內容。此為DRAFT小樣，尚未接入遊戲，未宣稱整個圖庫或HD版本完成。

### 頭像擴充至126張

新增六張原版參考重繪，皆來自劇本索引0。由原版KAOGRF及SINARIO資料核對名字、武將ID與圖庫頁碼，
以內建imagegen逐張重繪1254×1254原生PNG。原圖64×64，放大參考圖僅作模型輸入，未當作HD成品。

| 武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---:|---:|---|
| 劉度 | 30 | 113 | r1 |
| 夏侯尚 | 38 | 85 | r2 |
| 牛金 | 48 | 57 | r2 |
| 吳蘭 | 53 | 93 | r2 |
| 張勳 | 74 | 34 | r1 |
| 張松 | 77 | 19 | r2 |

- 劉度：保留斜向頭部、灰眉灰鬚、藍色帽頂與金褐頭飾，以及灰白／金色衣領。細部頭飾結構、藍色區域邊界與材質屬重繪詮釋，未宣稱原版逐像素幾何一致。
- 夏侯尚：保留長臉、黑鬚、灰綠高帽、紅色後布與綠色肩部。第2稿移除完整金色尖拱，改為少量原色筆觸，並縮小誇張口部開口；保留原有灰白細節。
- 牛金：保留抬頭長臉、黑鬚、藍色頭布、原有紅色區域與綠衣。第2稿將過亮白色嘴內條紋收斂為灰白口部陰影，保留原有淡色區域；嘴部材質不作原版已證實語意。
- 吳蘭：保留寬臉、灰綠帽、褐金側片、藍領與金色肩部。第2稿收斂帽緣新增翹角，減少新增鬍鬚量與亮白牙列；帽上細小金色曲線仍屬來源色塊的重繪詮釋。
- 張勳：保留較窄長臉、細上唇鬚與下巴長鬚、金褐頭飾、原有小紅色區域與綠色後布。肩部點陣轉成連續灰金明暗，未新增規律格紋；細部線條屬重繪詮釋。
- 張松：保留年輕長臉、明暗與灰綠高帽、褐金帽帶及灰白衣領。第2稿移除新增金色拱形圖案，保留少量細小金灰色表面筆觸；細部材質意義仍未知。

[本批紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-13.json)保存細修與未選稿，
[來源計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-13-plan.json)、[原圖與參考圖](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-13.json)、
[120張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/120-r1-snapshot.json)、[初稿完整提示詞](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-13-prompts.json)、
[全部提示詞索引](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set.json)、[原生輸出路徑](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-13-native-outputs.json)、
[合併入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/integrate-batch13.py)、[文件登錄入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/record-batch13.py)、
[檢視縮圖一](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-13-contact-sheet-1.png)與[檢視縮圖二](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-13-contact-sheet-2.png)皆存於既有本機目錄。
縮圖只供檢視，原生PNG未改。細修的完整提示詞檔由本批原生輸出路徑與全部提示詞索引連入。

目前126張選稿、169次原生輸出、13份來源包。[七項離線抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-126-r1.json)通過，
含252張內嵌PNG、全部126張切換、六張新增畫面、四種檢視、滑鼠縮放及拖曳、來源視窗、
單檔搬移重開、窄視窗及兩倍像素密度。沿原2GB／1 CPU、30秒導覽及90秒外層執行。
app.js及build.py未變。[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-126-r1.json)由[驗證入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/verify-batch13.py)重生，
核對全部原生輸出與來源、放大參考圖、262項120張快照、城池／全圖頁及HTML重建一致。
現行verify.py與本批入口相同。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未宣稱完整圖庫或HD完成。

### 頭像擴充至132張

新增六張原版參考重繪，皆來自劇本索引0。由原版KAOGRF及SINARIO資料核對名字、武將ID與圖庫頁碼，
以內建imagegen逐張重繪1254×1254原生PNG。原圖64×64，放大參考圖僅作模型輸入，未當作HD成品。

| 武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---:|---:|---|
| 雷銅 | 99 | 136 | r1 |
| 李恢 | 100 | 63 | r1 |
| 冷苞 | 106 | 86 | r1 |
| 伍習 | 108 | 129 | r2 |
| 李暹 | 109 | 143 | r1 |
| 張衛 | 110 | 18 | r2 |

- 雷銅：保留寬圓臉、較張開眼神、短髭與小下巴鬚、金褐頭飾、原有綠色額前區域及灰綠後布、紅色衣領。未新增長鬚或凸起寶石；五官與細部筆觸屬重繪詮釋。
- 李恢：保留較老的長臉、灰眉灰髭、灰綠高帽與原有小金色區域、灰白衣領及金邊。灰髭轉為自然髮絲明暗，未新增長白鬚；帽前細部結構屬重繪詮釋。
- 冷苞：保留斜向長臉與髭線、金褐頭飾、原有藍色帽前區域及淡色外緣、藍色後布、紅領與金色肩部。藍色與淡色輪廓原圖已有，其材質及意義未知，未當作已證實寶石。
- 伍習：保留灰綠帽、原有不對稱細金線、綠色後布、金色肩部與紅灰衣領。第2稿將新增的濃髭及下巴鬍鬚修薄，恢復較短上唇筆觸與下巴陰影；帽前原有金色線條保留。
- 李暹：保留寬臉、深色髭鬚、褐色高帽與原有綠色區域及淡色反光、綠色肩部與灰藍衣緣。嘴部灰白區域保留為收斂明暗，綠色區域未新增完整亮框；材質意義未知。
- 張衛：保留灰綠頭飾、原有灰白表面細節、小紅黃色區域、紅色後布及灰金衣領。第2稿縮小新增圓章式圖案，保留小色點，並收斂嘴部亮白牙列；色塊含義與表面材質仍未知。

[本批紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-14.json)保存細修與未選稿，
[來源計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-14-plan.json)、[原圖與參考圖](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-14.json)、
[126張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/126-r1-snapshot.json)、[初稿完整提示詞](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-14-prompts.json)、
[全部提示詞索引](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set.json)、[原生輸出路徑](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-14-native-outputs.json)、
[合併入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/integrate-batch14.py)、[文件登錄入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/record-batch14.py)、
[檢視縮圖一](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-14-contact-sheet-1.png)與[檢視縮圖二](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-14-contact-sheet-2.png)皆存於既有本機目錄。
縮圖只供檢視，原生PNG未改。細修的完整提示詞檔由本批原生輸出路徑與全部提示詞索引連入。

目前132張選稿、177次原生輸出、14份來源包。[七項離線抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-132-r1.json)通過，
含264張內嵌PNG、全部132張切換、六張新增畫面、四種檢視、滑鼠縮放及拖曳、來源視窗、
單檔搬移重開、窄視窗及兩倍像素密度。沿原2GB／1 CPU、30秒導覽及90秒外層執行。
app.js及build.py未變。[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-132-r1.json)由[驗證入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/verify-batch14.py)重生，
核對全部原生輸出與來源、放大參考圖、277項126張快照、城池／全圖頁及HTML重建一致。
現行verify.py與本批入口相同。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未宣稱完整圖庫或HD完成。

### 頭像擴充至138張

新增六張原版參考重繪，皆來自劇本索引0。由原版KAOGRF及SINARIO資料核對名字、武將ID與圖庫頁碼，
以內建imagegen逐張重繪1254×1254原生PNG。原圖64×64，放大參考圖僅作模型輸入，未當作HD成品。

| 武將 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---:|---:|---|
| 嚴圃 | 111 | 133 | r2 |
| 鞏志 | 112 | 45 | r1 |
| 笮融 | 113 | 131 | r3 |
| 張英 | 114 | 140 | r1 |
| 嚴輿 | 116 | 74 | r3 |
| 周昕 | 117 | 137 | r2 |

- 嚴圃：保留斜向寬臉、灰綠帽、原有綠色區域、灰白衣領與金褐肩部。第2稿修薄新增濃髭與下巴鬚，將帽上過亮平行綠帶融回不規則灰綠明暗；細部材質含義未知。
- 鞏志：保留較長窄臉、斜向眼神、金褐頭飾與原有小紅色區域、灰白／灰綠衣領及肩部。帽頂裁切與細部結構、髭線屬重繪詮釋，未宣稱原版逐像素幾何一致。
- 笮融：保留長臉、灰綠帽與金邊、灰綠衣領。第2稿移除新增折線圖案、修薄髭線，但過度移除原有灰白色塊；第3稿補回小灰白區域，未沿用第1稿新增的可讀圖案。原圖六個像素樣本另存本批紀錄，材質及含義未知。
- 張英：保留朝右側臉、外露耳部、深色頭髮與髭鬚、灰綠衣領及金褐肩部，未新增頭盔或帽子。原圖側臉及鬍鬚輪廓沿用，細部髮絲屬重繪詮釋。
- 嚴輿：保留原圖側臉、深色髭鬚、灰綠／褐色頭飾及原有淡色區域、灰綠衣領和金色肩部。並排檢視確認原圖本就大張口，第2稿收得太小而未選；第3稿回到較大張口，將過長上方亮條修為較短灰白區域，保留下方淡色嘴緣。細部材質及精確五官幾何仍屬重繪詮釋。
- 周昕：保留寬圓臉、微笑、灰綠帽與金邊、紅色肩部及灰白衣領。第2稿移除新增亮白直條，修薄髭線並縮小過大的微笑開口，保留嘴部原有灰白區域；細部材質意義未知。

圖庫頁133同時供劇本索引0的111嚴圃與劇本索引3的117閻圃使用。
沿[四劇本原始名單](../../workplace/hd-redraw-local-20261104/faithful/portraits/source-roster.json)及原始SINARIO雜湊核對，保留兩個來源名稱，
本批選用劇本0的嚴圃，不自行訂正譯名，也不把同頁重複計入選稿數。

笮融原圖帽前確有灰白色塊：六個原圖樣本在本批紀錄，含RGB243／243／243與162／178／178。
第2稿將這些色塊過度移除，第3稿恢復小灰白區域；第1稿新增的折線圖案未沿用。
灰白色塊的實際材質與含義仍未知，未將它宣稱為原版徽章。

嚴輿第2稿曾將嘴部收得太小。並排檢視確認原圖本就大張口，第3稿回到第1稿的張口幅度，只收斂過長亮條。
[第一份候選與勘誤](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch15-first-candidate.json)保存24項檔案雜湊及已通過技術抽驗的頁面，
[最終選稿入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/finalize-batch15.py)及[新增原生稿紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-15-final-native.json)保存修正。
先前技術抽驗不作美術接受聲明；變更選稿後重新執行完整七項抽驗與完整性檢查。

[本批紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-15.json)保存細修與未選稿，
[來源計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-15-plan.json)、[原圖與參考圖](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-15.json)、
[132張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/132-r1-snapshot.json)、[初稿完整提示詞](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-15-prompts.json)、
[全部提示詞索引](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set.json)、[原生輸出路徑](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-15-native-outputs.json)、
[合併入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/integrate-batch15.py)、[文件登錄入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/record-batch15.py)、
[檢視縮圖一](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-15-contact-sheet-1.png)與[檢視縮圖二](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-15-contact-sheet-2.png)皆存於既有本機目錄。
縮圖只供檢視，原生PNG未改。細修的完整提示詞檔由本批原生輸出路徑與全部提示詞索引連入。

目前138張選稿、189次原生輸出、15份來源包。[七項離線抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result-138-r1.json)通過，
含276張內嵌PNG、全部138張切換、六張新增畫面、四種檢視、滑鼠縮放及拖曳、來源視窗、
單檔搬移重開、窄視窗及兩倍像素密度。沿原2GB／1 CPU、30秒導覽及90秒外層執行。
app.js及build.py未變。[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification-138-r1.json)由[驗證入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/verify-batch15.py)重生，
核對全部原生輸出與來源、放大參考圖、289項132張快照、城池／全圖頁及HTML重建一致。
現行verify.py與本批入口相同。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未宣稱完整圖庫或HD完成。

### 頭像擴充至146張

新增八張原版參考重繪，前六張來自劇本索引0，趙廣與趙統為索引1。由原版KAOGRF及SINARIO資料核對名字、武將ID與圖庫頁碼，
以內建imagegen逐張重繪1254×1254原生PNG。原圖64×64，放大參考圖僅作模型輸入，未當作HD成品。

| 武將 | 來源劇本 | 武將ID | 圖庫頁碼 | 選稿 |
|---|---:|---:|---:|---|
| 李別 | 0 | 118 | 138 | r1 |
| 刑道榮 | 0 | 119 | 132 | r2 |
| 劉賢 | 0 | 120 | 61 | r2 |
| 武安國 | 0 | 121 | 130 | r2 |
| 劉磐 | 0 | 122 | 83 | r2 |
| 楊昂 | 0 | 123 | 139 | r1 |
| 趙廣 | 1 | 69 | 56 | r2 |
| 趙統 | 1 | 74 | 134 | r2 |

- 李別：保留斜向長臉、原有髭線與小下巴陰影、褐色帽頂和灰綠帽帶、帽前灰白區域及灰金衣領。左右明暗與裁切沿原圖，細部髭線與材質仍屬重繪詮釋。
- 刑道榮：保留金褐頭飾與原有淡色區域、灰綠後布、濃黑髭鬚及灰金衣領。第2稿將新增亮白牙列收成較短灰白嘴部區域，保留原圖的小紅色色點；紅色含義未知。
- 劉賢：保留朝左側臉、外露耳部、金褐頭飾、綠色後布及金綠衣領。第2稿移除新增長下巴鬚、修薄髭線，恢復較乾淨的下巴與原有明暗；精確五官幾何仍屬重繪詮釋。
- 武安國：保留長臉、濃黑髭鬚、灰綠／褐色頭飾及原有小紅色帽前區域、灰金衣領。第2稿收斂嘴部過長紅線，保留原圖小紅色區域；原圖兩個嘴部紅色像素另存，不將其解釋成傷口。
- 劉磐：保留斜向抬頭、寬臉、灰綠頭飾與金邊、紅色後布及藍色衣領。第2稿移除新增濃髭與下巴鬚，保留原圖短上唇筆觸與下巴陰影。
- 楊昂：保留長臉、原有髭線及下巴陰影、灰綠高帽、小金色帽前區域、左側灰白區域及灰金衣領。帽前色塊與左側灰白區域材質／含義未知，未宣稱原版徽章或已證實髮色。
- 趙廣：保留年輕長臉、原有眼神、灰綠／褐色頭飾及灰白表面區域、小金色區域、紅色後布與衣領。第2稿移除新增長髭尾、下巴鬚與濃重鬚影，保留較乾淨的臉部。
- 趙統：保留朝左側臉、長鼻、灰綠頭飾與原有淡色區域、灰綠後布、紅領與灰金肩部。第2稿修薄新增濃髭及下巴鬚，收斂過大的嘴部開口與亮條，保留原有灰白色細節。

四劇本名單有149個原始名稱，共使用146個圖庫頁，本批後這146頁都有重繪小樣。
另由[完整名單核對工具](../../workplace/hd-redraw-local-20261104/faithful/portraits/audit-roster-all.go)重讀[四劇本全部508筆非空姓名列](../../workplace/hd-redraw-local-20261104/faithful/portraits/source-roster-all.json)，
在姓名去重前收集圖庫頁，確認仍為146頁，同名跨劇本換圖的情況為0；沒有用去重名單掩蓋換圖。
圖庫實際150頁，145、147、148、149四頁未出現在現有名單，用途仍待確認；未將它們虛構成新武將。
以上只表示名單所用頁已有小樣，不代表逐張接受、完整圖庫或正式引擎完成。

前兩次完整抽驗在第三個窄視窗載入時碰到90秒外層逾時，桌面與單檔搬移重開已通過。
[首次紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch16-browser-first-attempt.json)與[第二次紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch16-browser-second-attempt.json)保留工具、頁面及失敗檔案雜湊。
靜態核對改為逐個Buffer區塊組行，並從數字重新建立要留在集合中的短ID，避免保留整個正規表示式擷取字串。
512MB探針原先在256MB JavaScript heap限制失敗；改存短ID後，同一檔案、同一限制通過。
[探針](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch16-assets-probe.json)、[先前探針失敗](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch16-assets-probe-failure.json)、
[工具檢視](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch16-browser-review.json)及[登錄入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/record-batch16-browser-review.py)保存前後工具與限制。
HTML、app.js、原生PNG與主QA工具均未改；仍以相同2GB／1 CPU、30秒導覽及90秒外層通過全部七項、三個離線視窗。
不將所有載入時間差異或兩次逾時歸因於單一因素，也不將驗證逾時當作遊戲缺陷。

武安國原圖嘴部在(29,46)／(30,46)、刑道榮在(25,40)／(24,41)各有兩個RGB211／0／0紅色像素。
樣本與原圖雜湊保存在本批紀錄；細修保留原有小色點，色塊含義未知，未自行解釋為傷口或其他新增標記。

[本批紀錄](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-16.json)保存細修與未選稿，
[來源計畫](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-16-plan.json)、[原圖與參考圖](../../workplace/hd-redraw-local-20261104/faithful/portraits/originals-batch-16.json)、
[138張快照](../../workplace/hd-redraw-local-20261104/faithful/portraits/138-r1-snapshot.json)、[初稿完整提示詞](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-16-prompts.json)、
[全部提示詞索引](../../workplace/hd-redraw-local-20261104/faithful/portraits/generation-set.json)、[原生輸出路徑](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-16-native-outputs.json)、
[合併入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/integrate-batch16.py)、[文件登錄入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/record-batch16.py)、
[檢視縮圖一](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-16-contact-sheet-1.png)、[檢視縮圖二](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-16-contact-sheet-2.png)與[檢視縮圖三](../../workplace/hd-redraw-local-20261104/faithful/portraits/batch-16-contact-sheet-3.png)皆存於既有本機目錄。
縮圖只供檢視，原生PNG未改。細修的完整提示詞檔由本批原生輸出路徑與全部提示詞索引連入。

目前146張選稿、203次原生輸出、16份來源包。[七項離線抽驗](../../workplace/hd-redraw-local-20261104/faithful/portraits/qa-result.json)通過，
含292張內嵌PNG、全部146張切換、八張新增畫面、四種檢視、滑鼠縮放及拖曳、來源視窗、
單檔搬移重開、窄視窗及兩倍像素密度。沿原2GB／1 CPU、30秒導覽及90秒外層執行。
app.js及build.py未變。[完整性收據](../../workplace/hd-redraw-local-20261104/faithful/portraits/verification.json)由[驗證入口](../../workplace/hd-redraw-local-20261104/faithful/portraits/verify-batch16.py)重生，
核對全部原生輸出與來源、放大參考圖、301項138張快照、城池／全圖頁及HTML重建一致。
現行verify.py與本批入口相同。仍為DRAFT小樣，未取得逐張接受或接入正式引擎，未宣稱完整圖庫或HD完成。

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
| 新入口分組及窄視窗安排 | 已選單排與28／6自動收起；依[spec89](89-siege-damage-report-toggle.md)加入戰報後為十二個實體鈕 | 不阻塞已確認操作；AI與編輯器入口待各自接線 |
| HD正式素材產製與接線 | 五城十張已同意，其中五組景觀已接情報窗；35區／22接縫、146頭像及色盤0橫幅已接，未新增逐張接受聲明 | 已接範圍限定抽驗通過；其餘十頁景觀、其他色盤／場景、局部美術及其餘正常視窗待完成 |
| HD世界鏡頭與固定介面 | 正式地圖已接50%～200%，色盤0可用HD全圖；視野、焦點、邊界與命中抽驗通過，介面固定 | 不阻塞已接範圍；其他季節與液晶HD仍沿原圖 |
| 編輯套用時機、取消及存檔衍生欄位 | 待設計／證據審查 | 是，先做單人垂直鏈草稿 |
