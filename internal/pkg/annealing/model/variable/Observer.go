// Copyright (c) 2019 Australian Rivers Institute.

package variable

type Observer interface {
	ObserveDecisionVariable(variable DecisionVariable)
}

type Observable interface {
	Subscribe(observers ...Observer)
}
