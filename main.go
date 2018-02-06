package main

import (
	"fmt"
	"log"
	"time"

	"github.com/CharlesHolbrow/gm"
	"github.com/CharlesHolbrow/m"
	"github.com/rakyll/portmidi"
)

func main() {
	if err := portmidi.Initialize(); err != nil {
		panic("Error initializing portmidi: " + err.Error())
	}

	out, err := portmidi.NewOutputStream(2, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Melodies
	// cMinorSubgroups := m.MinorTriad(m.G).AllOctaves().Over(m.G4).AllSubgroups(8)
	// cMinor := m.Append(cMinorSubgroups[:4]...)
	// final := cMinor.Append(cMinor.Transpose(6), cMinor.Transpose(7), cMinor.Transpose(-5))

	// Melody
	roots := m.Group(m.G3, m.E4, m.C4, m.C4).Repeat(40) // Bassline
	roots = roots.Interleave(roots)

	// Patterns
	pBass := NewSequence(2)
	pBass.Add(0, m.Note{Vel: 100, Length: 1.5})
	pBass.Add(1.5, m.Note{Vel: 100, Length: 0.5})

	melody := NewSequence(1)
	midiCh := uint8(1)
	for i, root := range roots {
		seqEvent := pBass.Get(i)
		if m, ok := seqEvent.Event.(m.Note); ok {
			melody.Add(seqEvent.position, gm.Note{On: true, Note: root, Ch: midiCh, Vel: 100})
			melody.Add(seqEvent.position+m.Length, gm.Note{Note: root, Ch: midiCh})
		}
	}

	// pSevenths := m.NewPatternSubdivisions(16, 2)
	// pSevenths.RampValue(127, 2)

	// Create Sequence

	for event := range melody.Play(time.Second) {
		switch e := event.(type) {
		case gm.Note:
			fmt.Println(e.Midi())
			out.WriteShort(e.Midi())
		}
	}
	time.Sleep(10 * time.Millisecond)
}
