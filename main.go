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

	// Reason doesn't obey cc123. Send a note off message to every note.
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

	chords := []m.NoteGroup{
		// m.Sus4Triad(m.F4),h\
		m.MajorChord(m.G4),
		m.MinorChord(m.A4),
		m.MajorChord(m.F4),
		m.MajorChord(m.C4),
	}

	main := m.NewSequence()

	// Bass sequence
	major := chords[1]
	root := major[0]
	sus4 := m.Sus4Triad(root)
	d7 := m.DominantSeventh(root)

	bass := Bass(root - 12)
	oneKick := Kick(1, 1)
	twoKick := Kick(2, 1)

	// Measure 1 - kick
	main.CopyFrom(bass).AdvanceCursor(1)
	// Measure 2 - kick/bass
	main.CopyFrom(bass).CopyFrom(oneKick)
	main.AdvanceCursor(1)
	// Measure 3 - add synth
	seq := ChipArp(sus4, 2, 4)
	main.CopyFrom(bass, oneKick, seq)
	main.AdvanceCursor(1)
	// Measure 4
	main.CopyFrom(bass)
	main.CopyFrom(oneKick)
	main.AdvanceCursor(1)
	// Measure 5
	main.CopyFrom(bass)
	main.CopyFrom(ChipArp(sus4, 3, 4))
	main.AdvanceCursor(1)
	// Measure 6
	main.CopyFrom(bass)
	main.CopyFrom(ChipArp(sus4, 4, 4))
	main.AdvanceCursor(1)
	// Measure 5
	main.CopyFrom(bass)
	main.CopyFrom(ChipArp(sus4, 5, 4))
	main.CopyFrom(ChipArp(sus4.Transpose(-24), 6, 4))
	main.AdvanceCursor(1)
	// Measure 6
	main.CopyFrom(bass)
	main.CopyFrom(ChipArp(d7, 4, 4))
	main.AdvanceCursor(1)

	// Measure 7
	main.CopyFrom(oneKick)
	main.CopyFrom(ChordProgression(chords[:3], 0.75))
	main.AdvanceCursor(0.75)

	// Measure 8
	main.Add(0, "One")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Two")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Three")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords, 0.75))
	main.AdvanceCursor(0.75)

	// for i := 3; i <= 7; i++ {
	// 	main.Add(0, fmt.Sprintf("%d", i))
	// 	for _, chord := range chords {

	// 		root := chord[0] - 12
	// 		inv := m.Group(chord[2]-12, chord[0], chord[1])
	// 		bass := Bass(root)

	// 		c := main.CopyFrom(ChipArp(chord, i, i).CopyFrom(bass).CopyFrom(twoKick))
	// 		main.AdvanceCursor(1)
	// 		main.CopyFrom(bass).CopyFrom(twoKick)
	// 		main.AdvanceCursor(1)

	// 		c = IntArp(inv, i+1, i).CopyFrom(bass).CopyFrom(twoKick)
	// 		main.CopyFrom(c).AdvanceCursor(1)
	// 	}
	// }

	main.Add(0, "Increasing Groups of Four")
	for i := 1; i < 8; i++ {
		c := ChipArp(m.Sus4Triad(root), i, 4).CopyFrom(bass)
		main.CopyFrom(c).AdvanceCursor(1)
	}

	main.Add(0, "Increasing groups of four Interpolated")
	for i := 1; i < 2; i++ {
		c := IntArp(m.Sus4Triad(root), i, 7).CopyFrom(bass)
		main.CopyFrom(c).AdvanceCursor(1)
	}

	// Measure 7
	main.CopyFrom(oneKick)
	main.CopyFrom(ChordProgression(chords[:3], 0.75))
	main.AdvanceCursor(0.75)

	// Measure 8
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Stage 5")
	r := ThinArp(m.MinorTriad(root).AllOctaves().Under(root+36).Over(root), 24, 1.)
	main.CopyFrom(r, bass)
	main.AdvanceCursor(1)

	for event := range main.Play(time.Millisecond * 1400) {
		switch e := event.(type) {
		case gm.Note:
			// fmt.Println(e.Midi())
			out.WriteShort(e.Midi())
		case gm.CC:
			out.WriteShort(e.Midi())
		case string:
			fmt.Println(e)
		}

	}

	time.Sleep(10 * time.Millisecond)
}
