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

	root := m.NoteNumber(m.B4)

	main := m.NewSequence()

	// Bass sequence
	b := Bass(root - 12)
	bEnd := Bass(root - 12 + 4)

	// Sequence
	c1 := ChipArp(m.Sus4Triad(root), 7, 4).CopyFrom(b)
	c2 := IntArp(m.MajorTriad(root), 7, 4).CopyFrom(b)
	c3 := RandArp(m.MinorTriad(root), 7, 4).CopyFrom(b)
	c4 := RandArp(m.MajorTriad(root), 7, 4).CopyFrom(bEnd)

	main.CopyFrom(c1)
	main.Cursor = 1
	main.CopyFrom(c2)
	main.Cursor = 2
	main.CopyFrom(c3)
	main.Cursor = 3
	main.CopyFrom(c4)

	// pSevenths := m.NewPatternSubdivisions(16, 2)
	// pSevenths.RampValue(127, 2)

	for event := range main.Play(time.Millisecond * 2000) {
		switch e := event.(type) {
		case gm.Note:
			// fmt.Println(e.Midi())
			out.WriteShort(e.Midi())
		}
	}

	time.Sleep(10 * time.Millisecond)
}
