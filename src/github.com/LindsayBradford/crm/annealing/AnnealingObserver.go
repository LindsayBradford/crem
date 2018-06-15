// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

type AnnealingObserver interface {
	ObserveAnnealingEvent(event AnnealingEvent)
}

type AnnealingEvent struct {
	eventType AnnealingEventType
	annealer  Annealer
	note      string
}

type AnnealingEventType int

const (
	INVALID_EVENT AnnealingEventType = iota
	STARTED_ANNEALING
	STARTED_ITERATION
	FINISHED_ANNEALING
	NOTE
)

func (eventType AnnealingEventType) String() string {
	labels := [...]string{
		"INVALID_EVENT",
		"STARTED_ANNEALING",
		"STARTED_ITERATION",
		"FINISHED_ANNEALING",
		"NOTE"}

	if eventType < STARTED_ANNEALING || eventType > NOTE {
		return labels[INVALID_EVENT]
	}

	return labels[eventType]
}
