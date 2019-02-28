// Copyright (c) 2019 Australian Rivers Institute.

package action

var _ Observer = new(spyObserver)

type spyObserver struct {
	lastObserved        ManagementAction
	observationsCounted uint
}

func (os *spyObserver) ObserveAction(action ManagementAction) {
	os.lastObserved = action
	os.observationsCounted++
}

func (os *spyObserver) ObserveInitialisationAction(action ManagementAction) {
	os.lastObserved = action
	os.observationsCounted++
}

func (os *spyObserver) LastObserved() ManagementAction {
	return os.lastObserved
}

func (os *spyObserver) ObservationsCounted() uint {
	return os.observationsCounted
}

func (os *spyObserver) Reset() {
	os.lastObserved = nil
	os.observationsCounted = 0
}
