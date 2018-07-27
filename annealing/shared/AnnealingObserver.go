// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

type AnnealingObserver interface {
	ObserveAnnealingEvent(event AnnealingEvent)
}

type AnnealingEvent struct {
	EventType AnnealingEventType
	Annealer  Annealer
	Note      string
}

type AnnealingEventType int

const (
	InvalidEvent AnnealingEventType = iota
	StartedAnnealing
	StartedIteration
	FinishedIteration
	FinishedAnnealing
	Note
)

func (eventType AnnealingEventType) String() string {
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
