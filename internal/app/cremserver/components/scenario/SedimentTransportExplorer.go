// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"github.com/LindsayBradford/crem/annealing/explorer"
)

const (
	_                                  = iota
	BankErosionFudgeFactor      string = "BankErosionFudgeFactor"
	WaterDensity                string = "WaterDensity"
	LocalAcceleration           string = "LocalAcceleration"
	GullyCompensationFactor     string = "GullyCompensationFactor"
	SedimentDensity             string = "SedimentDensity"
	SuspendedSedimentProportion string = "SuspendedSedimentProportion"
)

type SedimentTransportSolutionExplorer struct {
	explorer.SingleObjectiveAnnealableExplorer

	parameters *SedimentTransportParameters
}

func (stse *SedimentTransportSolutionExplorer) WithName(name string) *SedimentTransportSolutionExplorer {
	stse.SingleObjectiveAnnealableExplorer.WithName(name)
	return stse
}

func (stse *SedimentTransportSolutionExplorer) WithParameters(params map[string]interface{}) *SedimentTransportSolutionExplorer {
	stse.parameters = new(SedimentTransportParameters).Initialise()
	stse.parameters.Merge(params)
	return stse
}

func (stse *SedimentTransportSolutionExplorer) ParameterErrors() error {
	return stse.parameters.ValidationErrors()
}
