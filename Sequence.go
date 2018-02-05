package main

import (
	"fmt"
	"sort"
	"time"
)

// An Event is just a interface{} type. You can add any time to a sequence
type Event interface{}

// SequenceEvent wraps an event, and adds metadata about the location.
type SequenceEvent struct {
	subPosition int
	position    float64       // dimensionless floating point value
	Time        time.Duration // Duration from seq start to event time
	Event       Event
}

// Sequence is an ordered collection of Events.
type Sequence struct {
	// content is important, because it stores the order that events were added
	content map[float64][]SequenceEvent

	// list is where the actual sorting happens
	list []SequenceEvent
}

// NewSequence creates and initializes a new Sequence
func NewSequence() *Sequence {
	return &Sequence{
		list:    make([]SequenceEvent, 0),
		content: make(map[float64][]SequenceEvent),
	}
}

// Add an event to the sequence. Position is a dimensionless point to place the
// event. The dimension can be set with the sequence.Sorted() function.
func (s *Sequence) Add(position float64, event Event) {
	if position < 0 {
		fmt.Printf("Bad event position: %f (%v)\n", position, event)
		panic("Cannot add event to with negative position")
	}

	events, ok := s.content[position]
	if !ok {
		events = make([]SequenceEvent, 0, 10)
		s.content[position] = events
	}

	timeEvent := SequenceEvent{
		Event:       event,
		position:    position,
		subPosition: len(events),
	}

	s.content[position] = append(events, timeEvent)
	s.list = append(s.list, timeEvent)
}

// Sorted creates a slice of TimeEvents. The .Time property of each event will
// be populated. To Add an event, you had to specify a dimensionless time
// position. Set that dimension now with the `unit` argument.
func (s *Sequence) Sorted(unit time.Duration) []SequenceEvent {
	sort.Sort(s)
	result := make([]SequenceEvent, len(s.list))
	for i, tEvent := range s.list {
		tEvent.Time = time.Duration(tEvent.position * float64(unit))
		result[i] = tEvent
	}
	return result
}

// Play back the sequence on the supplied channel. If out is nil, create a
// channel. returns the playback channel.
func (s *Sequence) Play(unit time.Duration) chan interface{} {
	start := time.Now()
	out := make(chan interface{})
	go func() {
		for _, tEvent := range s.Sorted(unit) {
			time.Sleep(time.Until(start.Add(tEvent.Time)))
			out <- tEvent.Event
		}
		close(out)
	}()
	return out
}

// sort.Interface methods
func (s *Sequence) Len() int      { return len(s.list) }
func (s *Sequence) Swap(i, j int) { s.list[i], s.list[j] = s.list[j], s.list[i] }
func (s *Sequence) Less(i, j int) bool {
	if s.list[i].position == s.list[j].position {
		return s.list[i].subPosition < s.list[j].subPosition
	}
	return s.list[i].position < s.list[j].position
}
