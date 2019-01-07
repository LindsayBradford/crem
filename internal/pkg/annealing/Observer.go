// Copyright (c) 2018 Australian Rivers Institute.

package annealing

type Observer interface {
	ObserveAnnealingEvent(event Event)
}

type Event struct {
	EventType EventType
	Annealer  Observable
	Note      string
}

type EventType int

const (
	InvalidEvent EventType = iota
	StartedAnnealing
	StartedIteration
	FinishedIteration
	FinishedAnnealing
	Note
)

func (eventType EventType) String() string {
	labels := [...]string{
		"InvalidEvent",
		"StartedAnnealing",
		"StartedIteration",
		"FinishedIteration",
		"FinishedAnnealing",
		"Note",
	}

	if eventType < StartedAnnealing || eventType > Note {
		return labels[InvalidEvent]
	}

	return labels[eventType]
}
