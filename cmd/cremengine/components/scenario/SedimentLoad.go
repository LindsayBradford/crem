// Copyright (c) 2019 Australian Rivers Institute.

package scenario

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
)

const SedimentLoadVariableName = "SedimentLoad"

type SedimentLoad struct {
	model.VolatileDecisionVariable
	bankSedimentContribution BankSedimentContribution
}

func (sl *SedimentLoad) Initialise(planningUnitTable *tables.CsvTable, parameters Parameters) *SedimentLoad {
	sl.SetName(SedimentLoadVariableName)
	sl.bankSedimentContribution.Initialise(planningUnitTable, parameters)
	sl.SetValue(sl.deriveInitialSedimentLoad())
	return sl
}

func (sl *SedimentLoad) deriveInitialSedimentLoad() float64 {
	return sl.bankSedimentContribution.OriginalSedimentContribution() +
		sl.gullySedimentContribution() +
		sl.hillSlopeSedimentContribution()
}

func (sl *SedimentLoad) gullySedimentContribution() float64 {
	return 0 // TODO: implement
}

func (sl *SedimentLoad) hillSlopeSedimentContribution() float64 {
	return 0 // implement
}
