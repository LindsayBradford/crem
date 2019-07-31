// Copyright (c) 2018 Australian Rivers Institute.

package suppapitnarm

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"

	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

const DefaultExplorableDecisionVariables = "SedimentProduced,ImplementationCost"

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.
		Initialise("Suppapitnarm Explorer Parameter Validation").
		Enforcing(ParameterSpecifications())
	return p
}

const (
	ExplorableDecisionVariables = "ExplorableDecisionVariables"
)

type optimisationDirection int

const (
	Invalid optimisationDirection = iota
	Minimising
	Maximising
)

func ParameterSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          ExplorableDecisionVariables,
			Validator:    IsString,
			DefaultValue: DefaultExplorableDecisionVariables,
		},
	)
	return specs
}
