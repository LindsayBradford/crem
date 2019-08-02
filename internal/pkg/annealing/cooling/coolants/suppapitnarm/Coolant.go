// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package suppapitnarm

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
)

type Coolant struct {
	rand.RandContainer
	parameters Parameters

	AcceptanceProbability float64
	Temperature           float64
	CoolingFactor         float64
}

func (c *Coolant) Initialise() *Coolant {
	c.parameters.Initialise()
	return c
}

func (c *Coolant) WithParameters(params parameters.Map) *Coolant {
	c.parameters.AssignOnlyEnforcedUserValues(params)

	c.Temperature = c.parameters.GetFloat64(StartingTemperature)
	c.CoolingFactor = c.parameters.GetFloat64(CoolingFactor)

	return c
}

func (c *Coolant) ParameterErrors() error {
	return c.parameters.ValidationErrors()
}

func (c *Coolant) DecideIfAcceptable(variableChanges []float64) bool {
	c.calculateAcceptanceProbability(variableChanges)
	randomValue := c.RandomNumberGenerator().Float64Unitary()
	return c.AcceptanceProbability > randomValue
}

func (c *Coolant) calculateAcceptanceProbability(variableChanges []float64) {
	probabilities := make([]float64, len(variableChanges))
	for index, individualChange := range variableChanges {
		absoluteChangeInObjectiveValue := math.Abs(individualChange)
		probabilities[index] = math.Exp(-absoluteChangeInObjectiveValue / c.Temperature)
	}

	finalProbability := float64(1)
	for _, probability := range probabilities {
		finalProbability = finalProbability * probability
	}
	c.AcceptanceProbability = finalProbability
}

func (c *Coolant) CoolDown() {
	c.Temperature *= c.CoolingFactor
}
