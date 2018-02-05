package main

// Pattern stores a collection of On/Off events along with a start time and duration for each event
type Pattern struct {
	length float64
	cursor float64
	events []patternEvent
}

type patternEvent struct {
	startPosition float64
	duration      float64
}

// NewPattern creates and initializes a new Pattern
func NewPattern(length float64) *Pattern {
	return &Pattern{
		length: length,
		events: make([]patternEvent, 0),
	}
}

// Push an Event to the pattern, advancing the cursor
func (p *Pattern) Push(advance, duration float64) {
	p.cursor = (p.cursor + advance)
	p.events = append(p.events, patternEvent{
		startPosition: p.cursor,
		duration:      duration,
	})
}
