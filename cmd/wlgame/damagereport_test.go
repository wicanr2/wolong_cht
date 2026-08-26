package main

import "testing"

// 損害報告預設關著，而且系統選單第 8 列點得動（docs/spec/89）。
//
// ⭐ **原版戰術結束後沒有損害報告**——那一行是 remake 的驗收資訊。
// 這一支釘的是「預設關」與「點得開」，不是那一行的內容。
func TestDamageReportDefaultsOffAndTogglesFromSystemRow(t *testing.T) {
	g := &game{}
	if g.damageReport {
		t.Error("預設應該是關的——原版沒有這個報告")
	}
	if got := damageReportValue(false); got != " 關 " {
		t.Errorf("關的值 ＝ %q", got)
	}
	if got := damageReportValue(true); got != " 開 " {
		t.Errorf("開的值 ＝ %q", got)
	}
	g.dispatchSystemRow(sysRowDamageReport, true)
	if !g.damageReport {
		t.Fatal("點第 8 列沒有打開")
	}
	g.dispatchSystemRow(sysRowDamageReport, false)
	if g.damageReport {
		t.Error("右鍵也該是 toggle")
	}
}

// 新列一律加在最後：插進去會把原版那六列往下推，
// docs/playtest/39 對那六列的比對就不再成立。
func TestSystemMenuKeepsOriginalSixRowsInPlace(t *testing.T) {
	if sysRowSave != 0 || sysRowVideo != 1 || sysRowSound != 2 ||
		sysRowStrategySpeed != 3 || sysRowTacticalSpeed != 4 || sysRowQuit != 5 {
		t.Fatal("原版那六列的索引被動到了")
	}
	if sysRowLordCorps != 6 || sysRowDamageReport != 7 {
		t.Errorf("remake 的兩列不在最後：lordCorps=%d damageReport=%d",
			sysRowLordCorps, sysRowDamageReport)
	}
	if sysRows != 8 {
		t.Errorf("列數 ＝ %d，預期 8", sysRows)
	}
	// 原版那六列的值格 y 座標不能因為加列而改變。
	for k := 0; k < 6; k++ {
		if got := sysRowRect(k).Min.Y; got != sysValueY+k*sysRowStep {
			t.Errorf("第 %d 列的 y ＝ %d，被推動了", k, got)
		}
	}
}
