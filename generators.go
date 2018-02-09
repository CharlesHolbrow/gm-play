package main

import (
	"github.com/CharlesHolbrow/m"
)

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

func ChordProgression(chords []m.NoteGroup, length float64) *m.Sequence {
	pattern := m.NewSequence()
	pattern.AddSubdivisions(len(chords), length, 0.8)
	pattern.Cursor = length
	s := m.NewSequence()
	s.AddRhythmicChords(pattern, chords, 2)
	return s
}

// ThinArp creates a random arpeggiation from notes. Then notes are randomly
// removed. `chance` is the liklyhood that a note will be kept.
func ThinArp(notes m.NoteGroup, subdivisions int, chance float64) *m.Sequence {
	repetitions := subdivisions/len(notes) + 1
	notes = notes.Repeat(repetitions).Permute()

	pattern := m.NewSequence()
	pattern.AddSubdivisions(subdivisions, 1, 1.1)
	pattern.Cursor = 1
	// pattern.RampSustainVelocity(110, 10)
	pattern.RandomRemove(chance)

	notes = notes[:pattern.Len()]

	s := m.NewSequence().AddRhythmicMelody(pattern, notes, 1)
	return s
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

// inserv n equally spaced kicks over length
func Kick(n int, length float64) *m.Sequence {
	pattern := m.NewSequence()
	pattern.AddSubdivisions(n, length, 0.5)
	pattern.Cursor = length
	melody := m.Group(36).Repeat(n)
	return m.NewSequence().AddRhythmicMelody(pattern, melody, 15)
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
