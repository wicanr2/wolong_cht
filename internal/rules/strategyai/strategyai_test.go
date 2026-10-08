package strategyai

import (
	"testing"

	"github.com/wicanr2/wolong_cht/internal/rules/diplomacy"
)

func testFaction() Faction {
	return Faction{
		Alive: true, Cities: 4, Funds: 50000,
		Reserves: [3]int{400, 600, 1000}, Aggression: 14, InvasionTarget: NoTarget,
	}
}

func TestPowerUsesOriginalCityAndFundsGates(t *testing.T) {
	f := testFaction()
	if got := Power(f); got != 500 {
		t.Fatalf("Power = %d, want 500", got)
	}
	f.Reserves = [3]int{65500, 65500, 65500}
	if got := Power(f); got != 2000 {
		t.Fatalf("Power cap = %d, want 2000", got)
	}
	f.Funds = 19 << 8
	if got := Power(f); got != 0 {
		t.Fatalf("low funds Power = %d, want 0", got)
	}
	// 原始 24 位負數的高兩 byte 是 unsigned 0xFFFF；不能被 Go 算術右移
	// 誤判成小於 0x13。
	f.Funds = -1
	if got := Power(f); got != 2000 {
		t.Fatalf("negative 24-bit funds Power = %d, want 2000", got)
	}
}

func TestShouldDeclareWarKeepsStrictComparisons(t *testing.T) {
	self := testFaction()
	target := Faction{Alive: true, Cities: 2, Funds: 50000, Reserves: [3]int{300, 300, 300}}
	c := Candidate{Faction: 7, Friendship: diplomacy.Friendship(0xA9)}
	if !ShouldDeclareWar(self, target, c) {
		t.Fatal("eligible candidate was rejected")
	}
	self.Funds = FundLimit(self.Cities) << 8
	if ShouldDeclareWar(self, target, c) {
		t.Fatal("funds equal to the limit must be rejected")
	}
	self = testFaction()
	c.Friendship = diplomacy.Friendship(FriendshipLimit(self.Aggression))
	if !ShouldDeclareWar(self, target, c) {
		t.Fatal("friendship equal to the limit must be accepted")
	}
	c.Friendship = diplomacy.Friendship(FriendshipLimit(self.Aggression) + 1)
	if ShouldDeclareWar(self, target, c) {
		t.Fatal("friendship above the limit must be rejected")
	}
}

func TestSortCandidates(t *testing.T) {
	got := SortCandidates([]Candidate{
		{Faction: 4, Friendship: diplomacy.Friendship(0xB0)},
		{Faction: 2, Friendship: diplomacy.Friendship(0xA9)},
		{Faction: 1, Friendship: diplomacy.Friendship(0xA9)},
	})
	if got[0].Faction != 2 || got[1].Faction != 1 || got[2].Faction != 4 {
		t.Fatalf("sorted candidates = %+v", got)
	}
}

func TestSortCandidateSwapsPreserveOriginalTailOrder(t *testing.T) {
	in := []Candidate{{Faction: 1, Friendship: 0x60}, {Faction: 2, Friendship: 0x60}, {Faction: 3, Friendship: 0x50}}
	got := SortCandidates(in)
	for i, want := range []int{3, 2, 1} {
		if got[i].Faction != want {
			t.Fatalf("candidate %d is %d, want %d", i, got[i].Faction, want)
		}
	}
	if in[0].Faction != 1 {
		t.Fatal("input mutated")
	}
}

func TestSortCandidateRawByteBoundaries(t *testing.T) {
	got := SortCandidates([]Candidate{{Faction: 3, Friendship: 0xFF}, {Faction: 2, Friendship: 0x80}, {Faction: 1, Friendship: 0x7F}, {Faction: 4, Friendship: 0xFF}})
	for i, want := range []int{1, 2, 3, 4} {
		if got[i].Faction != want {
			t.Fatalf("raw candidate %d: got %d want %d", i, got[i].Faction, want)
		}
	}
}

func TestPlayerMonthlyDeclarationUsesSameRule(t *testing.T) {
	self, target := testFaction(), testFaction()
	target.Reserves = [3]int{40, 40, 40}
	candidate := Candidate{Faction: 1, Friendship: 0x80}
	before := ShouldDeclareWar(self, target, candidate)
	self.Player = true
	if !before || ShouldDeclareWar(self, target, candidate) != before {
		t.Fatal("extra player gate in monthly declaration")
	}
}

func TestAtWarValue(t *testing.T) {
	got := AtWarValue(diplomacy.Friendship(0xA9), diplomacy.Friendship(0xB6))
	if got != diplomacy.Friendship(20) || !got.AtWar() {
		t.Fatalf("AtWarValue = %#x, want war value 20", got)
	}
}
