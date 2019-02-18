// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

package observer

import "github.com/LindsayBradford/crem/pkg/name"

type Observer interface {
	ObserveEvent(event Event)
}

type Event struct {
	EventType   EventType
	EventSource name.Identifiable
	Note        string
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
