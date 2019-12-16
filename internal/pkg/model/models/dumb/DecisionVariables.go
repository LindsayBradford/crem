// Copyright (c) 2019 Australian Rivers Institute.

package dumb

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/variableOld"
)

type DecisionVariables struct {
	variableOld.InductiveDecisionVariables
}

func (dv *DecisionVariables) Initialise() *DecisionVariables {
	dv.InductiveDecisionVariables = variableOld.NewInductiveDecisionVariables()
	dv.buildDecisionVariables()
	return dv
}

func (dv *DecisionVariables) buildDecisionVariables() {
	dv.NewForName("ObjectiveValue")
}
