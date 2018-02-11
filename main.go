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

	progression := []m.NoteGroup{
		// m.Sus4Triad(m.F4),h\
		m.MajorChord(m.G4),
		m.MinorChord(m.A4),
		m.MajorChord(m.F4),
		m.Sus4Chord(m.A4),
		// Sus4s
		m.Sus4Chord(m.A4).AllOctaves().Over(m.G3 + 1).Under(m.G5),
		m.Sus4Chord(m.A4).AllOctaves().Over(m.A3 + 1).Under(m.G5),
		m.Sus4Chord(m.A4).AllOctaves().Over(m.F3 + 1).Under(m.G5),
		m.Sus4Chord(m.A4).AllOctaves().Over(m.G3 + 1).Under(m.G5),
		// final
		m.DominantSeventh(m.A4),
	}
	bassProg := []*m.Sequence{
		// 4
		Bass(m.G3),
		Bass(m.A3),
		Bass(m.F3),
		Bass(m.G3),
		// 4
		Bass(m.G3),
		Bass(m.A3),
		Bass(m.F3),
		Bass(m.G3),
		// 4
		Bass(m.D4),
	}

	main := m.NewSequence()

	// Bass sequence
	minor := progression[1]
	root := minor[0]
	sus4 := m.Sus4Triad(root)
	d7 := m.DominantSeventh(root)

	bass := Bass(root - 12)
	oneKick := Kick(1, 1)
	twoKick := Kick(2, 1)

	// Begin Copying Sequences to Main
	// Progression
	for _, chord := range progression[:4] {
		main.CopyFrom(ChordProgression([]m.NoteGroup{chord}, 1))
		main.CopyFrom(twoKick)
		main.AdvanceCursor(1)
		main.CopyFrom(twoKick)
		main.AdvanceCursor(1)
	}

	//
	main.Add(0, "Increasing Groups of Four A")
	for i := 0; i < 4; i++ {
		subdivisions := i + 1
		main.Add(0, fmt.Sprintf("Groups of four: %v", subdivisions))
		chips := ChipArp(m.Sus4Triad(root), subdivisions, 4)
		bass := bassProg[i]
		main.CopyFrom(bass, oneKick, chips)
		main.AdvanceCursor(1)
	}

	// Measure
	main.CopyFrom(bass)
	main.CopyFrom(ChipArp(sus4, 6, 5))
	main.CopyFrom(ChipArp(sus4.Transpose(-24), 6, 5))
	main.AdvanceCursor(1)

	// Measure
	main.CopyFrom(Bass(d7[0] - 12 + 4))
	main.CopyFrom(ChipArp(d7, 4, 4))
	main.AdvanceCursor(1)

	// Measure - Stabs
	main.CopyFrom(oneKick)
	main.CopyFrom(ChordProgression(chords[0:3], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "One")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:4], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Two")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:4], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Three")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords, 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Increasing Groups of Four B")
	for i := 0; i < 4; i++ {
		subdivisions := i + i
		c := ChipArp(m.Sus4Triad(root), subdivisions, 4)
		b := bassProg[i]
		main.CopyFrom(c, b).AdvanceCursor(1)
	}

	main.Add(0, "Increasing groups of four Interpolated")
	for i := 1; i < 2; i++ {
		c := IntArp(m.Sus4Triad(root), i, 7).CopyFrom(bass)
		main.CopyFrom(c).AdvanceCursor(1)
	}

	// Measure 7
	main.CopyFrom(oneKick)
	main.CopyFrom(ChordProgression(chords[0:3], 0.75))
	main.AdvanceCursor(0.75)

	// Measure 8
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:4], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Stage 5")
	main.CopyFrom(ThinArp(m.MinorTriad(root).AllOctaves().Under(root+36).Over(root), 24, 1.))
	main.CopyFrom(bass)
	main.AdvanceCursor(1)

	for i, chord := range progression {
		// bass := Bass(chord[0] - 12)
		bass := bassProg[i]
		main.CopyFrom(ChordProgression([]m.NoteGroup{chord}, 1))
		main.CopyFrom(bass, twoKick)
		main.AdvanceCursor(1)
		main.CopyFrom(bass, twoKick)
		main.AdvanceCursor(1)
	}
	// Progression
	for _, chord := range progression[:4] {
		main.CopyFrom(ChordProgression([]m.NoteGroup{chord}, 1))
		main.CopyFrom(ChipArp(chord, 3, 4))
		main.CopyFrom(twoKick)
		main.AdvanceCursor(1)

		main.CopyFrom(twoKick)
		main.AdvanceCursor(1)
	}

	//
	main.Add(0, "Increasing Groups of Four A")
	for i := 0; i < 4; i++ {
		subdivisions := i + 1
		// Ingredients
		chord := m.Sus4Triad(root)
		chips := ChipArp(chord, subdivisions, 4)
		pad := ChordProgression([]m.NoteGroup{chord}, 1)
		bass := bassProg[i]

		//
		main.CopyFrom(bass, twoKick, chips, pad)
		main.Add(0, fmt.Sprintf("Groups of four: %v", subdivisions))
		main.AdvanceCursor(1)
	}

	// Measure
	main.CopyFrom(bass, oneKick)
	main.CopyFrom(ChipArp(sus4, 6, 5))
	main.CopyFrom(ChordProgression([]m.NoteGroup{sus4}, 0.5))
	main.CopyFrom(ChipArp(sus4.Transpose(-24), 6, 5))
	main.AdvanceCursor(1)

	// Measure
	main.CopyFrom(oneKick)
	main.CopyFrom(Bass(d7[0] - 12 + 4))
	main.CopyFrom(ChordProgression([]m.NoteGroup{d7}, 0.5))
	main.CopyFrom(ChipArp(d7, 4, 4))
	main.AdvanceCursor(1)

	// Measure - Stabs
	main.Add(0, "Zero")
	main.CopyFrom(oneKick)
	main.CopyFrom(ChordProgression(chords[0:3], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "One")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:4], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Two")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:4], 0.75))

	main.AdvanceCursor(0.75)

	main.Add(0, "Three")
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords, 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Increasing Groups of Four")
	for i := 0; i < 4; i++ {
		subdivisions := i + i
		chord := m.Sus4Triad(root)
		c := ChipArp(chord, subdivisions, 4)
		p := ChordProgression([]m.NoteGroup{chord}, 1./.8*0.75)
		b := bassProg[i]
		main.CopyFrom(c, b, p, oneKick).AdvanceCursor(1)
	}

	main.Add(0, "Increasing groups of four Interpolated B")
	for i := 1; i < 2; i++ {
		c := IntArp(m.Sus4Triad(root), i, 7).CopyFrom(bass)
		main.CopyFrom(c).AdvanceCursor(1)
	}

	// Measure 7
	main.CopyFrom(oneKick)
	main.CopyFrom(ChordProgression(chords[0:3], 0.75))
	main.AdvanceCursor(0.75)

	// Measure 8
	main.CopyFrom(bass, twoKick)
	main.CopyFrom(ChordProgression(chords[1:4], 0.75))
	main.AdvanceCursor(0.75)

	main.Add(0, "Stage 5")
	r := ThinArp(m.MinorTriad(root).AllOctaves().Under(root+36).Over(root), 24, 1.)
	main.CopyFrom(r, bass)
	main.AdvanceCursor(1)

	for i, chord := range progression {
		bass := bassProg[i]
		main.CopyFrom(ChordProgression([]m.NoteGroup{chord}, 1))
		c := ThinArp(chord.AllOctaves().Over(chord[0]).Under(chord[0]+24), 24, 0.2)
		main.CopyFrom(c)
		main.CopyFrom(bass, twoKick)
		main.AdvanceCursor(1)
		main.CopyFrom(bass, twoKick)
		main.CopyFrom(ThinArp(chord.AllOctaves().Over(chord[0]-12).Under(chord[0]+24), 24, .3))
		main.AdvanceCursor(1)
	}

	////////////////////////////////////////////////////
	// repeat the whole thing with different chord group
	////////////////////////////////////////////////////

	chords = []m.NoteGroup{
		// m.Sus4Triad(m.F4),h\
		m.MajorChord(m.G4),
		m.MinorChord(m.A4),
		m.MajorChord(m.F4),
		m.MajorChord(m.C4),
	}

	progression = []m.NoteGroup{
		// m.Sus4Triad(m.F4),h\
		m.MajorChord(m.G4),
		m.MinorChord(m.A4),
		m.MajorChord(m.F4),
		m.MajorChord(m.C5),
		// Sus4s
		m.Sus4Chord(m.A4).AllOctaves().Over(m.G3 + 1).Under(m.G5),
		m.Sus4Chord(m.A4).AllOctaves().Over(m.A3 + 1).Under(m.G5),
		m.Sus4Chord(m.A4).AllOctaves().Over(m.F3 + 1).Under(m.G5),
		m.MajorChord(m.A4).AllOctaves().Over(m.G3 + 1).Under(m.G5),
		// final
		m.MajorChord(m.C5),
	}
	bassProg = []*m.Sequence{
		// 4
		Bass(m.G3),
		Bass(m.A3),
		Bass(m.F3),
		Bass(m.C3),
		// 4
		Bass(m.G3),
		Bass(m.A3),
		Bass(m.F3),
		Bass(m.C3),
		// 4
		Bass(m.C3),
	}

	final := m.NewSequence()

	// Bass sequence
	minor = progression[1]
	root = minor[0]
	sus4 = m.Sus4Triad(root)
	d7 = m.DominantSeventh(root)

	bass = Bass(root - 12)
	oneKick = Kick(1, 1)
	twoKick = Kick(2, 1)

	// Begin Copying Sequences to final

	// Progression
	for _, chord := range progression[:4] {
		final.CopyFrom(ChordProgression([]m.NoteGroup{chord}, 1))
		final.CopyFrom(ChipArp(chord, 3, 4))
		final.CopyFrom(twoKick)
		final.AdvanceCursor(1)

		final.CopyFrom(twoKick)
		final.AdvanceCursor(1)
	}

	//
	final.Add(0, "Increasing Groups of Four A")
	for i := 0; i < 4; i++ {
		subdivisions := i + 1
		// Ingredients
		chord := m.Sus4Triad(root)
		chips := ChipArp(chord, subdivisions, 4)
		pad := ChordProgression([]m.NoteGroup{chord}, 1)
		bass := bassProg[i]

		//
		final.CopyFrom(bass, twoKick, chips, pad)
		final.Add(0, fmt.Sprintf("Groups of four: %v", subdivisions))
		final.AdvanceCursor(1)
	}

	// Measure
	final.CopyFrom(bass, oneKick)
	final.CopyFrom(ChipArp(sus4, 6, 5))
	final.CopyFrom(ChordProgression([]m.NoteGroup{sus4}, 0.5))
	final.CopyFrom(ChipArp(sus4.Transpose(-24), 6, 5))
	final.AdvanceCursor(1)

	// Measure
	final.CopyFrom(oneKick)
	final.CopyFrom(Bass(d7[0] - 12 + 4))
	final.CopyFrom(ChordProgression([]m.NoteGroup{d7}, 0.5))
	final.CopyFrom(ChipArp(d7, 4, 4))
	final.AdvanceCursor(1)

	// Measure - Stabs
	final.Add(0, "Zero")
	final.CopyFrom(oneKick)
	final.CopyFrom(ChordProgression(chords[0:3], 0.75))
	final.AdvanceCursor(0.75)

	final.Add(0, "One")
	final.CopyFrom(bass, twoKick)
	final.CopyFrom(ChordProgression(chords[1:4], 0.75))
	final.AdvanceCursor(0.75)

	final.Add(0, "Two")
	final.CopyFrom(bass, twoKick)
	final.CopyFrom(ChordProgression(chords[1:4], 0.75))

	final.AdvanceCursor(0.75)

	final.Add(0, "Three")
	final.CopyFrom(bass, twoKick)
	final.CopyFrom(ChordProgression(chords, 0.75))
	final.AdvanceCursor(0.75)

	final.Add(0, "Increasing Groups of Four")
	for i := 0; i < 4; i++ {
		subdivisions := i + i
		chord := m.Sus4Triad(root)
		c := ChipArp(chord, subdivisions, 4)
		p := ChordProgression([]m.NoteGroup{chord}, 1./.8*0.75)
		b := bassProg[i]
		final.CopyFrom(c, b, p, oneKick).AdvanceCursor(1)
	}

	final.Add(0, "Increasing groups of four Interpolated B")
	for i := 1; i < 2; i++ {
		c := IntArp(m.Sus4Triad(root), i, 7).CopyFrom(bass)
		final.CopyFrom(c).AdvanceCursor(1)
	}

	// Measure 7
	final.CopyFrom(oneKick)
	final.CopyFrom(ChordProgression(chords[0:3], 0.75))
	final.AdvanceCursor(0.75)

	// Measure 8
	final.CopyFrom(bass, twoKick)
	final.CopyFrom(ChordProgression(chords[1:4], 0.75))
	final.AdvanceCursor(0.75)

	final.Add(0, "Stage 5")
	r = ThinArp(m.MinorTriad(root).AllOctaves().Under(root+36).Over(root), 24, 1.)
	final.CopyFrom(r, bass)
	final.AdvanceCursor(1)

	for i, chord := range progression {
		bass := bassProg[i]
		final.CopyFrom(ChordProgression([]m.NoteGroup{chord}, 1))
		c := ThinArp(chord.AllOctaves().Over(chord[0]).Under(chord[0]+24), 24, 0.2)
		final.CopyFrom(c)
		final.CopyFrom(bass, twoKick)
		final.AdvanceCursor(1)
		final.CopyFrom(bass, twoKick)
		final.CopyFrom(ThinArp(chord.AllOctaves().Over(chord[0]-24).Under(chord[0]+24), 24, .3))
		final.AdvanceCursor(1)
	}

	// Copy Paste from Beginning!

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
	for event := range final.Play(time.Millisecond * 1200) {
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
