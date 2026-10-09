//go:build matching_engagement

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wicanr2/dosgolem/internal/dos"
	"github.com/wicanr2/dosgolem/internal/machine"
)

const engagementOriginalScratch = "/tmp/engagement-original-files"
const engagementNativeScratch = "/tmp/engagement-native-files"

type engagementSaveState struct {
	label                 []byte
	expected              []byte
	labelCalls, saveCalls int
	checks                int
}

var engagementSaveInput, engagementScenarioInput []byte
var engagementSaveStates = map[*machine.Machine]*engagementSaveState{}
var engagementSaveReceipts []map[string]any

func engagementResetSaveFiles(original, native *dos.DOS) {
	if len(engagementSaveInput) != 88832 || len(engagementScenarioInput) != 88832 {
		panic("engagement save inputs must be four original 56C0-byte slots")
	}
	for _, pair := range []struct {
		d    *dos.DOS
		path string
	}{{original, engagementOriginalScratch}, {native, engagementNativeScratch}} {
		if pair.d.Scratch != pair.path {
			panic("engagement save requires separate original/native Scratch")
		}
		if err := os.MkdirAll(pair.path, 0700); err != nil {
			panic(err)
		}
		if err := os.WriteFile(filepath.Join(pair.path, "SAVE.DAT"), engagementSaveInput, 0600); err != nil {
			panic(err)
		}
	}
	engagementSaveStates = map[*machine.Machine]*engagementSaveState{}
}

// Each side observes its own source memory before the original writer. These
// observations form the file oracle; they never initialize the other machine.
func engagementObserveSave(device *machine.Machine, t uint16, regs registers) {
	if currentInput.group != "save" || currentInput.engagement.phase != "archive-save" {
		return
	}
	state := engagementSaveStates[device]
	if state == nil {
		state = &engagementSaveState{}
		engagementSaveStates[device] = state
	}
	mem, code := device.Mem, int(regs.CS)*16
	world := int(strategyWord(mem, code+0xD52)) * 16
	if t == 0x8D4D {
		scenario := int(mem[code+0xD01])
		if scenario > 3 {
			panic("save label scenario is outside original SINARIO.DAT")
		}
		label := append([]byte(nil), engagementScenarioInput[scenario*0x56C0+0x40:scenario*0x56C0+0x60]...)
		copy(label[6:12], []byte{0xB6, 0xD5, 0xA4, 0x4F, 0xA1, 0x47})
		faction := int(strategyWord(mem, code+0xCFD))
		for _, pair := range []struct{ field, offset int }{{1, 12}, {2, 26}} {
			general := int(mem[world+faction+pair.field])
			if general >= 128 {
				panic("save label has invalid original general index")
			}
			name := world + 0x4242 + general*32
			copy(label[pair.offset:pair.offset+6], mem[name:name+6])
		}
		copy(label[18:26], []byte{0xA1, 0x40, 0xAD, 0x78, 0xAE, 0x76, 0xA1, 0x47})
		state.label, state.labelCalls = label, state.labelCalls+1
		return
	}
	if t != 0x8CFF {
		return
	}
	slot := int(byte(regs.AX))
	if slot != currentInput.profile || slot < 0 || slot > 3 || state.labelCalls != 1 {
		panic("save caller must select one legal slot and build one original label")
	}
	labelBase := int(strategyWord(mem, code+0x987C)) * 16
	if !bytes.Equal(mem[labelBase:labelBase+32], state.label) {
		panic("original save label differs from independent SINARIO/name model")
	}
	expected := append([]byte(nil), engagementSaveInput...)
	base := slot * 0x56C0
	copy(expected[base:base+0x3B], mem[code+0xCF0:code+0xD2B])
	copy(expected[base+0x40:base+0x60], state.label)
	copy(expected[base+0x80:base+0x52C0], mem[world:world+0x5240])
	queue := int(strategyWord(mem, code+0xD56)) * 16
	copy(expected[base+0x52C0:base+0x56C0], mem[queue:queue+0x400])
	state.expected, state.saveCalls = expected, state.saveCalls+1
	state.checks += 5
}

func engagementAuditSavedFile(device *machine.Machine, d *dos.DOS, c resumeCase) int {
	state := engagementSaveStates[device]
	if state == nil || state.labelCalls != 1 || state.saveCalls != 1 || len(state.expected) != 88832 {
		panic("save case did not execute exactly one real label/writer pair")
	}
	actual, err := os.ReadFile(filepath.Join(d.Scratch, "SAVE.DAT"))
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(actual, state.expected) {
		first := 0
		for first < len(actual) && first < len(state.expected) && actual[first] == state.expected[first] {
			first++
		}
		panic(fmt.Sprintf("independent save slot%d length%d first difference %X", c.profile, len(actual), first))
	}
	return state.checks + 1
}

func engagementCompareSaveFiles(original, native *dos.DOS, c resumeCase) {
	a, err := os.ReadFile(filepath.Join(original.Scratch, "SAVE.DAT"))
	if err != nil {
		panic(err)
	}
	b, err := os.ReadFile(filepath.Join(native.Scratch, "SAVE.DAT"))
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(a, b) {
		panic("independently written original/native SAVE.DAT differs")
	}
	engagementSaveReceipts = append(engagementSaveReceipts, map[string]any{
		"scenario": c.scenario, "slot": c.profile, "bytes": len(a),
		"original_sha256":                fmt.Sprintf("%x", sha256.Sum256(a)),
		"c_sha256":                       fmt.Sprintf("%x", sha256.Sum256(b)),
		"independent_expected_match":     true,
		"unchanged_gaps_and_other_slots": true,
		"separate_scratch":               true,
	})
}
