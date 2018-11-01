// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/errors"
)

type SedimentTransportSolutionExplorer struct {
	explorer.SingleObjectiveAnnealableExplorer

	parameters  map[string]interface{}
	paramErrors *errors.CompositeError
}

func (stse *SedimentTransportSolutionExplorer) WithName(name string) *SedimentTransportSolutionExplorer {
	stse.SingleObjectiveAnnealableExplorer.WithName(name)
	return stse
}

func (stse *SedimentTransportSolutionExplorer) WithParameters(params map[string]interface{}) *SedimentTransportSolutionExplorer {
	stse.parameters = params
	stse.validateParameters()
	return stse
}

func (stse *SedimentTransportSolutionExplorer) validateParameters() {
	stse.paramErrors = errors.NewComposite("SedimentTransportSolutionExplorer parameters")

	for key, value := range stse.parameters {
		stse.validateParameter(key, value)
	}
}

func (stse *SedimentTransportSolutionExplorer) validateParameter(key string, value interface{}) {
	if key != "Penalty" {
		stse.paramErrors.AddMessage("Key is not [Penalty]")
	}
}

func (stse *SedimentTransportSolutionExplorer) ParameterErrors() error {
	if stse.paramErrors.Size() > 0 {
		return stse.paramErrors
	}
	return nil
}
