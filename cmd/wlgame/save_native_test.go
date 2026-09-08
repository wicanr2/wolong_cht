package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wicanr2/wolong_cht/internal/assets/library"
	"github.com/wicanr2/wolong_cht/internal/savepath"
	"github.com/wicanr2/wolong_cht/internal/state"
)

const nativeScenario = "../../workplace/orig/dosv/SINARIO.DAT"

func TestLauncherLoadRestoresNativeRuntime(t *testing.T) {
	if _, err := os.Stat(nativeScenario); err != nil {
		t.Skip("缺少原版資料")
	}
	w, err := state.LoadScenario(nativeScenario, 0)
	if err != nil {
		t.Fatal(err)
	}
	w.Player = 0
	snap := w.TakeSnapshot()
	snap.CityCursor = 91
	snap.Routes = map[int][][2]int{2: {{5, 6}, {7, 8}}}
	if err := w.Restore(snap); err != nil {
		t.Fatal(err)
	}
	save := filepath.Join(t.TempDir(), "SAVE.DAT")
	writer := &game{world: w, sourceFile: nativeScenario, saveBase: nativeScenario, saveFile: save}
	if err := writer.writeSave(1); err != nil {
		t.Fatal(err)
	}
	lib, err := library.Load(filepath.Dir(nativeScenario))
	if err != nil {
		t.Fatal(err)
	}
	reader := &game{lib: lib, origDir: filepath.Dir(nativeScenario), sourceFile: nativeScenario, saveFile: save,
		launcher: &launcherModel{phase: launcherLoad}, saveUI: saveUIState{active: true, action: saveRead}}
	if err := reader.applyLauncherResult(launcherResult{kind: launcherStartLoad, slot: 1}); err != nil {
		t.Fatal(err)
	}
	if reader.launcher != nil || reader.saveUI.active || reader.world == nil {
		t.Fatalf("讀檔未完成：%#v", reader.launcher)
	}
	got := reader.world.TakeSnapshot()
	if got.CityCursor != 91 || len(got.Routes[2]) != 2 {
		t.Fatalf("首頁讀檔遺失原生狀態：cursor=%d routes=%v", got.CityCursor, got.Routes)
	}
	nativePath, err := savepath.NativePath(save, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nativePath, []byte("broken native save"), 0644); err != nil {
		t.Fatal(err)
	}
	previous := reader.world
	reader.launcher = &launcherModel{phase: launcherLoad}
	reader.saveUI = saveUIState{active: true, action: saveRead}
	if err := reader.applyLauncherResult(launcherResult{kind: launcherStartLoad, slot: 1}); err != nil {
		t.Fatal(err)
	}
	if reader.launcher == nil || reader.launcher.notice == "" || !reader.saveUI.active || reader.world != previous {
		t.Fatal("原生檔損壞時應保留原世界及選槽錯誤，不可降級讀取")
	}
}

// 存檔要一次寫兩份：原版格式（拿去 DOSBox 的）與原生檔（遊戲讀的）。
// 讀檔優先原生檔，而且 routes／cityCursor 那些原版裝不下的狀態要活著。
func TestSaveWritesBothFormatsAndReadsBackNative(t *testing.T) {
	if _, err := os.Stat(nativeScenario); err != nil {
		t.Skip("找不到原版 SINARIO.DAT，跳過")
	}
	w, err := state.LoadScenario(nativeScenario, 0)
	if err != nil {
		t.Fatal(err)
	}
	snap := w.TakeSnapshot()
	snap.CityCursor = 91
	snap.Routes = map[int][][2]int{2: {{5, 6}, {7, 8}}}
	if err := w.Restore(snap); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	g := &game{
		world:      w,
		sourceFile: nativeScenario,
		saveBase:   nativeScenario,
		saveFile:   filepath.Join(dir, "SAVE.DAT"),
	}
	g.saveUI = saveUIState{active: true, action: saveWrite}
	g.clickSaveSlot(1)
	if g.saveUI.active {
		t.Fatalf("桌面點第二槽未直接完成存檔：%s", g.lastEvent)
	}

	// ① 原版格式那一份在，而且大小對得上（四個 22,208 B 區塊）。
	info, err := os.Stat(g.saveFile)
	if err != nil {
		t.Fatalf("原版格式存檔沒寫出來：%v", err)
	}
	if info.Size() != 22208*4 {
		t.Errorf("原版格式存檔大小 %d，want %d", info.Size(), 22208*4)
	}

	// ② 原生檔在，而且是另一個檔案。
	nativePath, err := savepath.NativePath(g.saveFile, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(nativePath); err != nil {
		t.Fatalf("原生存檔沒寫出來：%v", err)
	}

	// ③ 讀回來走原生檔，執行期狀態還在。
	// 先把記憶體裡的值改掉，否則「讀回來還是 91」可能只是沒動過。
	dirty := g.world.TakeSnapshot()
	dirty.CityCursor = 5
	dirty.Routes = nil
	if err := g.world.Restore(dirty); err != nil {
		t.Fatal(err)
	}
	if err := g.readSave(1); err != nil {
		t.Fatal(err)
	}
	got := g.world.TakeSnapshot()
	if got.CityCursor != 91 {
		t.Errorf("據點游標 = %d，want 91——原版格式裝不下它，代表讀的不是原生檔",
			got.CityCursor)
	}
	if len(got.Routes[2]) != 2 {
		t.Errorf("行軍路線 = %v，want 兩段", got.Routes[2])
	}

	// ④ 原生檔壞掉要炸，不能靜靜退回原版格式——那會讓玩家以為讀到了
	//    最新進度，實際上少了執行期狀態。
	if err := os.WriteFile(nativePath, []byte("{ 壞掉的 json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := g.readSave(1); err == nil {
		t.Error("壞掉的原生存檔被靜靜跳過了")
	}

	// ⑤ 原生檔不存在時要退回原版格式。
	if err := os.Remove(nativePath); err != nil {
		t.Fatal(err)
	}
	if err := g.readSave(1); err != nil {
		t.Fatalf("沒有原生檔時應該退回原版格式：%v", err)
	}
}
