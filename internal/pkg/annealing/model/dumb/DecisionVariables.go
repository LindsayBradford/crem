// Copyright (c) 2019 Australian Rivers Institute.

package dumb

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
)

type DecisionVariables struct {
	variable.InductiveDecisionVariables
}

func (dv *DecisionVariables) Initialise() *DecisionVariables {
	dv.InductiveDecisionVariables = variable.NewInductiveDecisionVariables()
	dv.buildDecisionVariables()
	return dv
}

func (dv *DecisionVariables) buildDecisionVariables() {
	dv.NewForName(variable.ObjectiveValue)
}
