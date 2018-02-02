package main

import (
	"fmt"
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

	notes := AllNotes(C, Eb, G)
	fmt.Println(notes)
	for _, note := range notes[4:] {
		out.WriteShort(gm.Note{On: true, Note: note, Vel: 127}.Midi())
		time.Sleep(time.Millisecond * 40)
		out.WriteShort(gm.Note{Note: note, Vel: 127}.Midi())
	}
}
