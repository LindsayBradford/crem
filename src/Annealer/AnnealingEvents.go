// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

type AnnealingEvent int

const (
	INVALID_EVENT AnnealingEvent = iota
	STARTED_ANNEALING
	STARTED_ITERATION
	FINISHED_ANNEALING
	NOTE
)

func (event AnnealingEvent) String() string {
	labels := [...]string{
		"INVALID_EVENT",
		"STARTED_ANNEALING",
		"STARTED_ITERATION",
		"FINISHED_ANNEALING",
		"NOTE"}

	if event < STARTED_ANNEALING || event > NOTE {
		return labels[INVALID_EVENT]
	}

	return labels[event]
}
