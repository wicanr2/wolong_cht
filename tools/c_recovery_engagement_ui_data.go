//go:build matching_engagement

package main

// engagementUICases first completes the real 11B5A retreat and world restore.
// The settings cancellation tail calls 11C8D, whose map/display banks must
// already have passed through 11B76 -> 119CA. Tactical A156 state is not this
// caller state. The original 18B7C archive selector ABI uses entry DS to
// address the code globals; AL selects one of three title rows, and DX is
// the original filename offset. Cancellation completes the real open/read/
// draw/close path. It never accepts a row or enters the 18CFF write routine.
func engagementUICases(bases []resumeCase) []resumeCase {
	if len(bases) != 4 {
		panic("engagement UI requires four original scenario bases")
	}
	var cases []resumeCase
	ui := func(scenario int, target uint32) resumeCase {
		c := engagementClone(bases[scenario])
		c.group, c.target = "ui-archive", target
		c.engagement.warmup = false
		c.engagement.worldWarmup = true
		if c.engagement.fullEntrySP == 0 {
			c.engagement.fullEntrySP = c.entrySP
		}
		c.engagement.ds, c.engagement.es = 0xFFFF, 0xFFFF
		c.engagement.phase = "archive-cancel"
		c.engagement.postPatches = nil
		// 11B5A restores the outer caller's SI/DI. Retain the factory's
		// legal world-army locators; these UI routines preserve or replace
		// them internally instead of taking an invented tactical unit index.
		c.ax, c.bx, c.cx, c.dx = 0, 0, 0, 0x0DA6
		c.events = []listEvent{{100, 120, 1}}
		c.mapEvents, c.menuChoices = nil, nil
		return c
	}
	for scenario := range bases {
		// 18B5D sets AL=2 and SAVE.DAT itself. Its carry branch exits before
		// 18D4D/18CFF, so this is the complete read-only cancellation caller.
		cases = append(cases, ui(scenario, 0x18B5D))
		for _, mode := range []byte{0, 1, 2} {
			c := ui(scenario, 0x18B7C)
			c.ax = uint16(mode)
			cases = append(cases, c)
		}
		// 18C20 reads four 128-byte headers at 0x56C0-byte intervals. Both
		// original files contain four complete slots and are mounted read-only.
		for _, filename := range []uint16{0x0D9A, 0x0DA6} {
			c := ui(scenario, 0x18C20)
			c.dx = filename
			cases = append(cases, c)
		}
		// 160CC registers six rows: DX=352, BX=152+24*row and
		// CX=0206. The real 6056 table dispatches 6084,6097,60A1,
		// 5FF1,60A5,60B4. The four settings counts at 5FF4 are
		// 2,5,5,5, so zero is a legal starting value for each option.
		for row := 0; row < 6; row++ {
			c := ui(scenario, 0x15FAA)
			c.group, c.engagement.phase = "settings", "settings-cancel"
			c.engagement.postPatches = []engagementPatch{{
				segment: 0xFFFF, offset: 0x0CF8, bytes: []byte{0, 0, 0, 0},
			}}
			c.events = []listEvent{{368, uint16(160 + 24*row), 0}, {100, 120, 1}}
			// Row 0 opens the archive selector, canceled by the next
			// right event. Row 5 opens 193E9: explicitly cancel that
			// popup, so 160C8 never calls the unrelated main exit.
			c.menuChoices = []int{-1}
			cases = append(cases, c)
		}
	}
	return cases
}
