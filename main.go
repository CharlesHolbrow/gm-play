package main

import (
	"log"
	"time"

	"github.com/CharlesHolbrow/gm"
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

	cMajor := MajorTriad(Ab).AllOctaves()[4:]
	cMinor := MinorTriad(C).AllOctaves()[10:]

	for _, note := range cMinor.Interleave(cMajor) {
		out.WriteShort(gm.Note{On: true, Note: note, Vel: 127}.Midi())
		time.Sleep(time.Millisecond * 40)
		out.WriteShort(gm.Note{Note: note, Vel: 127}.Midi())
	}

}
