package main

import "github.com/CharlesHolbrow/m"

func arpRiser(chord m.NoteGroup, subdivisions, repetitions int) (melody m.NoteGroup) {
	// Chip Melody
	cMinorSubgroups := chord.AllOctaves().Over(chord[0]).AllSubgroups(subdivisions)
	chipMelody := m.Append(cMinorSubgroups[:repetitions]...)
	return chipMelody
}

func arpInterleaver(chord m.NoteGroup, subdivisions, repetitions int) (melody m.NoteGroup) {
	cMinorSubgroups := chord.AllOctaves().Over(chord[0]).AllSubgroups(subdivisions)
	chipMelody := m.Append(cMinorSubgroups[:repetitions]...)
	totalNotes := subdivisions * repetitions
	if totalNotes%2 == 1 {
		chipMelody = chipMelody[:len(chipMelody)-1]
	}
	firstHalf := chipMelody[len(chipMelody)/2:]
	secondHalf := chipMelody[:len(chipMelody)/2]
	return firstHalf.Interleave(secondHalf)
}

func arpRandom(chord m.NoteGroup, subdivisions, repetitions int) (melody m.NoteGroup) {
	return arpRiser(chord, subdivisions, repetitions).Permute()
}

func rampDown(melody m.NoteGroup) (pattern *m.Sequence) {
	pattern = m.NewSequence()
	pattern.AddSubdivisions(len(melody), 1, .8)
	pattern.Cursor = 1
	pattern.RampSustainVelocity(100, 0)
	return
}

func ChipArp(chord m.NoteGroup, subdivisions, repetitions int) *m.Sequence {
	melody := arpRiser(chord, subdivisions, repetitions)
	pattern := rampDown(melody)

	s := m.NewSequence()
	s.AddRhythmicMelody(pattern, melody, 1)
	return s
}

func IntArp(chord m.NoteGroup, subdivisions, repetitions int) *m.Sequence {
	melody := arpInterleaver(chord, subdivisions, repetitions)
	pattern := rampDown(melody)
	s := m.NewSequence()
	s.AddRhythmicMelody(pattern, melody, 1)
	return s
}

func RandArp(chord m.NoteGroup, subdivisions, repetitions int) *m.Sequence {
	melody := arpRandom(chord, subdivisions, repetitions)
	pattern := rampDown(melody)
	s := m.NewSequence()
	s.AddRhythmicMelody(pattern, melody, 1)
	return s
}

func Bass(root m.NoteNumber) *m.Sequence {
	// Bass Melody
	bassMelody := m.Group(root, root) // Bassline
	// Bass Pattern
	bassPattern := m.NewSequence()
	bassPattern.AddSustain(0, 0.75, 100)
	bassPattern.AddSustain(0.75, .25, 50)
	bassPattern.Cursor = 1

	s := m.NewSequence()
	s.AddRhythmicMelody(bassPattern, bassMelody, 0)
	return s
}
