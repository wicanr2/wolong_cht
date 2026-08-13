package main

import "testing"

func TestIdleClockGateRequiresStablePointerAndNoCommand(t *testing.T) {
	var gate idleClockGate
	cases := []struct {
		name        string
		x, y        int
		inputActive bool
		want        bool
	}{
		{name: "first observation is not idle", x: 12, y: 34, want: false},
		{name: "stable pointer starts idle clock", x: 12, y: 34, want: true},
		{name: "pointer movement pauses clock", x: 13, y: 34, want: false},
		{name: "command pauses clock even when pointer is stable", x: 13, y: 34, inputActive: true, want: false},
		{name: "stable pointer resumes on following idle frame", x: 13, y: 34, want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := gate.Allows(tc.x, tc.y, tc.inputActive); got != tc.want {
				t.Fatalf("Allows(%d, %d, %t) = %t, want %t", tc.x, tc.y, tc.inputActive, got, tc.want)
			}
		})
	}
}
