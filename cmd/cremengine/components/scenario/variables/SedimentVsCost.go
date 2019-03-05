// Copyright (c) 2019 Australian Rivers Institute.

package variables

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
)

const SedimentVsCostName = "SedimentVsCost"

var _ variable.InductiveDecisionVariable = new(SedimentVsCost)

type SedimentVsCost struct {
	variable.CompositeInductiveDecisionVariable
}

func (sc *SedimentVsCost) Initialise() *SedimentVsCost {
	sc.SetName(SedimentVsCostName)
	return sc
}
