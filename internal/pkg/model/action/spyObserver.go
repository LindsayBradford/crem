// Copyright (c) 2019 Australian Rivers Institute.

package action

var _ Observer = new(spyObserver)

// spyObserver implements the management action Reporting interface, offering "test spy" functionality, allowing
// test frameworks to "spy" on the observation of management action state changes.
// See: https://martinfowler.com/bliki/TestDouble.html
type spyObserver struct {
	lastObserved        ManagementAction
	observationsCounted uint
}

func (os *spyObserver) ObserveAction(action ManagementAction) {
	os.lastObserved = action
	os.observationsCounted++
}

func (os *spyObserver) ObserveActionInitialising(action ManagementAction) {
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
