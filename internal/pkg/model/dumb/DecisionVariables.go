// Copyright (c) 2019 Australian Rivers Institute.

package dumb

import "github.com/LindsayBradford/crem/internal/pkg/model"

type DecisionVariables struct {
	model.VolatileDecisionVariables
}

func (dv *DecisionVariables) Initialise() *DecisionVariables {
	dv.VolatileDecisionVariables = model.NewVolatileDecisionVariables()
	dv.buildDecisionVariables()
	return dv
}

func (dv *DecisionVariables) buildDecisionVariables() {
	dv.NewForName(model.ObjectiveValue)
}
