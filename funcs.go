package main

// Group creates a new NoteGroup with the supplied notes
func Group(notes ...NoteNumber) NoteGroup {
	return notes
}

// MinorTriad creates a minor triad based on root
func MinorTriad(root NoteNumber) NoteGroup {
	third := (root + 3) % 12
	fifth := (root + 7) % 12
	return Group(root, third, fifth)
}

// MajorTriad creates a minor triad based on root, and repeats that
// triad in every octave.
func MajorTriad(root NoteNumber) NoteGroup {
	third := (root + 4) % 12
	fifth := (root + 7) % 12
	return Group(root, third, fifth)
}

// Dedupe iterates over a NoteGroup, and removes any notes that are not occuring
// for the first time.
func (notes NoteGroup) Dedupe() (result NoteGroup) {
	noteMap := make(map[NoteNumber]int, len(notes))
	result = make(NoteGroup, 0, len(notes))

	for _, note := range notes {
		noteMap[note]++
		if noteMap[note] == 1 {
			result = append(result, note)
		}
	}
	return result
}

// AllOctaves creates a new NoteGroup, with all the original notes repeated in
// every octave.
func (notes NoteGroup) AllOctaves() NoteGroup {
	result := make(NoteGroup, 0, 128)
	var i NoteNumber
	for i = lowestNote; i <= highestNote; i++ {
		for _, n := range notes.Dedupe() {
			if i%12 == n%12 {
				result = append(result, i)
				break
			}
		}
	}
	return result
}

// Append NoteGroups into one larger group
func (notes NoteGroup) Append(appendages ...NoteGroup) (result NoteGroup) {
	result = notes
	for _, group := range appendages {
		result = append(result, group...)
	}
	return result
}

// Interleave multiple groups together. This chooses the shortest group, and
// creates a new group with all of the others interleaved.
//
// In this example, there are 3 groups, and the shortest group has a length of
// 2, so the result is 6 units long:
//
// ([1,1]).Interleave([2,2], [5,6,7,8]) == [1,2,5,1,2,6]
func (notes NoteGroup) Interleave(others ...NoteGroup) (result NoteGroup) {
	// find the shortest group
	shortest := notes
	for _, group := range others {
		if len(group) < len(shortest) {
			shortest = group
		}
	}
	totalGroups := len(others) + 1
	groupSize := len(shortest)
	resultSize := groupSize * totalGroups
	result = make(NoteGroup, 0, resultSize)

	for i := 0; i < groupSize; i++ {
		result = append(result, notes[i])
		for _, group := range others {
			result = append(result, group[i])
		}
	}

	return result
}
