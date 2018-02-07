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
	// BUG(charles) output stream should be locked with a mutex
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

	root := m.NoteNumber(m.B3)
	subdivisions := 4 // Divide each beat into n parts
	repetitions := 8  // repeat

	// Chip Melody
	cMinorSubgroups := m.MinorTriad(root).AllOctaves().Over(root + 12).AllSubgroups(subdivisions)
	chipMelody := m.Append(cMinorSubgroups[:repetitions]...)

	// Chip Pattern
	chipPattern := m.NewSequence()
	chipPattern.AddSubdivisions(subdivisions*repetitions, 1, .8)
	chipPattern.Cursor = 1
	chipPattern.RampSustainVelocity(100, 0)

	// Bass Melody
	bassMelody := m.Group(root, root) // Bassline

	// Bass Pattern
	bassPattern := m.NewSequence()
	bassPattern.AddSustain(0, 0.75, 100)
	bassPattern.AddSustain(0.75, .25, 50)
	bassPattern.Cursor = 1

	s := m.NewSequence()
	s.AddRhythmicMelody(chipPattern, chipMelody, 0)
	s.AddRhythmicMelody(bassPattern, bassMelody, 1)
	s.Cursor = 1
	s.CopyFrom(s)
	s.Cursor = 2
	s.CopyFrom(s)
	s.Cursor = 3
	s.CopyFrom(s)
	s.Cursor = 4

	// pSevenths := m.NewPatternSubdivisions(16, 2)
	// pSevenths.RampValue(127, 2)

	for event := range s.Play(time.Millisecond * 2000) {
		switch e := event.(type) {
		case gm.Note:
			// fmt.Println(e.Midi())
			out.WriteShort(e.Midi())
		}
	}

	time.Sleep(10 * time.Millisecond)
}
