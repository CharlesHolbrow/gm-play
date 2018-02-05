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

	afMajor := m.MajorTriad(m.Ab).AllOctaves().Over(m.C3).AllSubgroups(8)[:6]
	// cMinor := m.MinorTriad(m.C).AllOctaves().Over(m.C2).Under(m.C6)
	final := m.Append(afMajor...)
	fmt.Println(len(final))

	s := &Sequence{}
	s.Add(2, "Five")
	s.Add(0, "One")
	s.Add(1./3., "Two")
	s.Add(1./3., "Three")
	s.Add(1./3., "Four")

	fmt.Println(s.Sorted(time.Second))

	var v uint8 = 124
	for _, note := range final {
		fmt.Println(v)
		out.WriteShort(gm.Note{On: true, Note: note, Vel: v}.Midi())
		time.Sleep(time.Millisecond * 40)
		out.WriteShort(gm.Note{Note: note, Vel: 0}.Midi())
		v = (v - 2)
		if v == 0 {
			break
		}
	}

}
