package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// Reaper doesn't obey cc123. Send a note off message on every note.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)

		// Block until a signal is received.
		<-c
		fmt.Println("\nSending note-off event on each midi note")
		for ch := uint8(0); ch <= 15; ch++ {
			for n := uint8(0); n <= 127; n++ {
				out.WriteShort(gm.Note{Note: n, Ch: ch}.Midi())
			}
		}
		time.Sleep(10 * time.Millisecond)
		portmidi.Terminate()
		os.Exit(0)
	}()

	// Melodies
	// cMinorSubgroups := m.MinorTriad(m.G).AllOctaves().Over(m.G4).AllSubgroups(8)
	// cMinor := m.Append(cMinorSubgroups[:4]...)
	// final := cMinor.Append(cMinor.Transpose(6), cMinor.Transpose(7), cMinor.Transpose(-5))

	// Melody
	notes := m.Group(m.G3, m.E4, m.C4, m.C4).Repeat(40) // Bassline
	notes = notes.Interleave(notes)

	// Patterns
	rhythm := m.NewSequence()
	rhythm.AddSustain(0, 1.5, 100)
	rhythm.AddSustain(1.5, 0.5, 50)
	rhythm.Cursor = 2

	s := m.NewSequence()
	s.AddRhythmicMelody(rhythm, notes, 1)

	// pSevenths := m.NewPatternSubdivisions(16, 2)
	// pSevenths.RampValue(127, 2)

	for event := range s.Play(time.Second) {
		switch e := event.(type) {
		case gm.Note:
			fmt.Println(e.Midi())
			out.WriteShort(e.Midi())
		}
	}

	time.Sleep(10 * time.Millisecond)
}
