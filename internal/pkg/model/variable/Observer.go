// Copyright (c) 2019 Australian Rivers Institute.

package variable

type Observer interface {
	ObserveDecisionVariable(variable DecisionVariable)
}

type Observable interface {
	Subscribe(observers ...Observer)
}

type ContainedDecisionVariableObservers struct {
	observers []Observer
}

func (c *ContainedDecisionVariableObservers) Observers() []Observer {
	return c.observers
}

func (c *ContainedDecisionVariableObservers) Subscribe(observers ...Observer) {
	if c.observers == nil {
		c.observers = make([]Observer, 0)
	}

	for _, newObserver := range observers {
		c.observers = append(c.observers, newObserver)
	}
}
