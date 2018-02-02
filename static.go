package main

// C  = 0
// Db = 1
// D  = 2
// Eb = 3
// E  = 4
// F  = 5
// Gb = 6
// G  = 7
// Ab = 8
// A  = 9
// Bb = 10
// B  = 11

// NoteNumber is a Midi Note Number
type NoteNumber = uint8

const lowestNote = 0
const highestNote = 127

// C  Pitch Class
const C = 0

// Cs represents the C Sharp Pitch Class
const Cs = 1

// Df represents the D Flat Pitch Class
const Df = 1

// Db represents the Flat Pitch Class
const Db = 1

// D Pitch Class
const D = 2

// Ds represents the D Sharp Pitch Class
const Ds = 3

// Ef represents E Flat Pitch Class
const Ef = 3

// Eb represents the E Flat Pitch Class
const Eb = 3

// E  Pitch Class
const E = 4

// F  Pitch Class
const F = 5

// Fs represents the F Sharp Pitch Class
const Fs = 6

// Gf represents the G Flat Pitch Class
const Gf = 6

// Gb represents the G Flat Pitch Class
const Gb = 6

// G represetns the G Pitch Class
const G = 7

// Gs Represents the G Sharp Pitch Class
const Gs = 8

// Af represents the A Flat Pitch Class
const Af = 8

// Ab represents the A Flat Pitch Class
const Ab = 8

// A  Pitch Class
const A = 9

// As represents the A  Sharp Pitch Class
const As = 10

// Bf represents the B Flat Pitch Class
const Bf = 10

// Bb represents the Flat Pitch Class
const Bb = 10

// B  Pitch Class
const B = 11

var pitchesFlats = [...]string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb"}
var pitchesSharps = [...]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#"}
var pitchesMap = map[string]int{
	"C":  C,
	"C#": Cs, "Db": Db,
	"D":  D,
	"D#": Ds, "Eb": Eb,
	"E":  E,
	"F":  F,
	"F#": Fs, "Gb": Gb,
	"G":  G,
	"G#": Gs, "Ab": Ab,
	"A":  A,
	"A#": As, "Bb": Bb,
	"B": B}

// AllNotes returns a slice with every note in notes
func AllNotes(notes ...NoteNumber) []NoteNumber {
	result := make([]NoteNumber, 0, 127)
	var i NoteNumber
	for i = lowestNote; i <= highestNote; i++ {
		for _, n := range notes {
			if i%12 == n {
				result = append(result, i)
				break
			}
		}
	}
	return result
}
