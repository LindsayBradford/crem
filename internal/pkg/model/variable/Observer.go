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
	c.observers = observers
}
