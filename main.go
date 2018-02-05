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

	cMinor := m.Append(m.MinorTriad(m.C).AllOctaves().Over(m.C2).AllSubgroups(7)[:12]...)
	cMinor = cMinor.Append(cMinor.Reverse())
	fMinor := cMinor.Transpose(5)
	bFlatMinor := fMinor.Transpose(-7)
	// afMajor := m.Append(m.MajorTriad(m.Ab).AllOctaves().Over(m.C3).AllSubgroups(7)[:12]...)
	// final := afMajor.Interleave(cMinor)
	final := cMinor.Append(fMinor, bFlatMinor).Repeat(10)
	fmt.Println(len(final))

	s := &Sequence{}

	for i, number := range final {
		s.Add(float64(i)*20, gm.Note{On: true, Note: number, Vel: 120})
		s.Add(float64(i)*20+32, gm.Note{Note: number})
	}

	for event := range s.Play(time.Millisecond) {
		switch e := event.(type) {
		case gm.Note:
			out.WriteShort(e.Midi())
		}
	}
}
