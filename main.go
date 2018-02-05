package main

import (
	"log"
	"time"

	"github.com/CharlesHolbrow/gm"
	"github.com/CharlesHolbrow/m"
	"github.com/rakyll/portmidi"
)

// AddNotePattern adds notes (with the given pattern) to a sequence
func (s *Sequence) AddNotePattern(notes m.NoteGroup, pattern *m.Pattern, midiCh uint8) {
	for i, note := range notes {
		pEvent := pattern.Get(i)
		start := pEvent.StartPosition
		end := start + pEvent.Duration
		s.Add(start, gm.Note{On: true, Note: note, Vel: uint8(pEvent.Value), Ch: midiCh})
		s.Add(end, gm.Note{Note: note, Ch: midiCh})
	}
}

func main() {
	if err := portmidi.Initialize(); err != nil {
		panic("Error initializing portmidi: " + err.Error())
	}

	out, err := portmidi.NewOutputStream(2, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Melodies
	cMinorSubgroups := m.MinorTriad(m.G).AllOctaves().Over(m.G4).AllSubgroups(8)
	cMinor := m.Append(cMinorSubgroups[:4]...)
	final := cMinor.Append(cMinor.Transpose(6), cMinor.Transpose(7), cMinor.Transpose(-5))

	roots := m.Group(m.G3, m.E4, m.C4, m.C4) // Bassline
	roots = roots.Interleave(roots)

	// Patterns
	pBass := m.NewPattern(2)
	pBass.Push(1.5, 100).Advance(1.5).Push(0.5, 100)

	pSevenths := m.NewPatternSubdivisions(16, 2)
	pSevenths.RampValue(127, 2)

	// Create Sequence
	s := NewSequence()

	s.AddNotePattern(roots, pBass, 1)
	s.AddNotePattern(final, pSevenths, 0)

	for event := range s.Play(time.Second) {
		switch e := event.(type) {
		case gm.Note:
			out.WriteShort(e.Midi())
		}
	}

	// cMinor = cMinor.Append(cMinor.Reverse())
	// fMinor := cMinor.Transpose(5)
	// bFlatMinor := fMinor.Transpose(-7)

	// final := cMinor.Append(fMinor, bFlatMinor).Repeat(1)
	// fmt.Println(len(final))

	// s := NewSequence()

	// for i, number := range final {
	// 	s.Add(float64(i*20), gm.Note{On: true, Note: number, Vel: 120})
	// 	s.Add(float64(i*20+20), gm.Note{Note: number})
	// }

	// for event := range s.Play(time.Millisecond) {
	// 	switch e := event.(type) {
	// 	case gm.Note:
	// 		out.WriteShort(e.Midi())
	// 	}
	// }
	time.Sleep(10 * time.Millisecond)
}
