//go:build matching_engagement

package main

// engagementExtraCases requires separate, freshly seeded DOS Scratch layers
// before any save case executes. The runner owns that file-state preparation
// and independently checks the four 18CFF write ranges on each side.
func engagementExtraCases(bases []resumeCase) []resumeCase {
	if len(bases) != 4 {
		panic("engagement extra cases require four scenario bases")
	}
	world := func(scenario int, target uint32, group string) resumeCase {
		c := engagementClone(bases[scenario])
		c.group, c.target = group, target
		c.engagement.warmup = false
		c.engagement.worldWarmup = true
		c.engagement.ds, c.engagement.es = 0xFFFF, 0xFFFF
		if c.engagement.fullEntrySP == 0 {
			c.engagement.fullEntrySP = c.entrySP
		}
		c.engagement.postPatches = nil
		c.engagement.phase = "world-extra"
		c.events = []listEvent{{100, 120, 1}}
		c.mapEvents, c.menuChoices = nil, nil
		return c
	}
	var cases []resumeCase
	for scenario := range bases {
		for slot := 0; slot < 4; slot++ {
			c := world(scenario, 0x18B5D, "save")
			c.profile = slot
			c.engagement.phase = "archive-save"
			// 18B5D supplies SAVE.DAT and AL=2. The real selector uses
			// DX=256, BX=144+48*slot, CX=020F; it accepts one click.
			// 18D4D builds the label, then 18CFF writes the selected slot.
			// The second event cancels the caller's next selector loop.
			c.events = []listEvent{{264, uint16(152 + 48*slot), 0}, {100, 120, 1}}
			cases = append(cases, c)
		}
		if len(bases[scenario].outcome.armies) < 2 {
			panic("engagement world UI lacks the original foreign army owner")
		}
		foreign := bases[scenario].outcome.armies[1].owner
		for _, hover := range []byte{foreign, 0xFF} {
			for _, target := range []uint32{0x15A3A, 0x15DBB} {
				c := world(scenario, target, "world-ui")
				// 15DBB reads this code byte as a faction index, with FF
				// taking its original blank-label branch. 15A3A calls it,
				// the 192-city 15CC6 scan, and the 19528/1950F map images.
				c.engagement.postPatches = []engagementPatch{{
					segment: 0xFFFF, offset: 0x98A7, bytes: []byte{hover},
				}}
				cases = append(cases, c)
			}
		}
		// Original world-panel title and hotspot entry; no synthetic frame
		// close or restore is called in isolation.
		cases = append(cases, world(scenario, 0x1614A, "world-ui"))
	}
	return cases
}
