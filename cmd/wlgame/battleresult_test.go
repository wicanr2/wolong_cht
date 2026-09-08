package main

import (
	"github.com/wicanr2/wolong_cht/internal/rules/tactical"
	"testing"
	"time"
)

func TestBattleResultDefaultAndVisibleTimeout(t *testing.T) {
	var r battleResultTimer
	b := &tactical.Battle{}
	now := time.Unix(100, 0)
	if !r.ready(b, 0, now, false) {
		t.Fatal("預設關閉仍等待確認")
	}
	if r.ready(b, 3, now, true) {
		t.Fatal("首次顯示前吃掉先前輸入")
	}
	r.show(b, now)
	if r.ready(b, 3, now.Add(3*time.Second-time.Nanosecond), false) {
		t.Fatal("提前到期")
	}
	if !r.ready(b, 3, now.Add(3*time.Second), false) {
		t.Fatal("到期仍等待確認")
	}
	if !r.ready(b, 3, now.Add(time.Second), true) {
		t.Fatal("無法提前關閉")
	}
	r.show(b, now.Add(time.Second))
	if !r.ready(b, 3, now.Add(3*time.Second), false) {
		t.Fatal("重繪重設倒數")
	}
	next := &tactical.Battle{}
	if r.ready(next, 3, now.Add(time.Hour), false) {
		t.Fatal("新戰鬥沿用舊期限")
	}
	r.show(next, now.Add(time.Hour))
	if r.ready(next, 3, now.Add(time.Hour+time.Second), false) {
		t.Fatal("新戰鬥提前關閉")
	}
}

func TestBattleResultOptionPreservesOriginalRowsAndFits(t *testing.T) {
	g := &game{}
	if g.battleResultSeconds != 0 {
		t.Fatal("結果頁預設不是關")
	}
	for _, want := range []int{3, 5, 10, 15, 30, 0} {
		g.dispatchSystemRow(sysRowBattleResult, true)
		if g.battleResultSeconds != want {
			t.Fatalf("設定=%d，預期%d", g.battleResultSeconds, want)
		}
	}
	if sysRowBattleResult != 8 || sysRowQuit != 5 || sysWinY+sysWinH > screenH {
		t.Fatal("新增選項推動原版列或超出畫面")
	}
}
