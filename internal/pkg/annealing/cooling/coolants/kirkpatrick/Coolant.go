// Copyright (c) 2019 Australian Rivers Institute.

package kirkpatrick

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
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
	c.parameters.CreateEmpty().
		WithSpecifications(
			DefineSpecifications(),
		).AssigningDefaults()
	return c
}

func (c *Coolant) WithParameters(params parameters.Map) *Coolant {
	c.parameters.AssignUserValues(params)

	c.Temperature = c.parameters.GetFloat64(StartingTemperature)
	c.CoolingFactor = c.parameters.GetFloat64(CoolingFactor)

	return c
}

func (c *Coolant) ParameterErrors() error {
	return c.parameters.ValidationErrors()
}

func (c *Coolant) DecideIfAcceptable(objectiveFunctionChange float64) bool {
	c.calculateAcceptanceProbability(objectiveFunctionChange)
	randomValue := c.RandomNumberGenerator().Float64Unitary()
	return c.AcceptanceProbability > randomValue
}

func (c *Coolant) calculateAcceptanceProbability(objectiveFunctionChange float64) {
	absoluteChangeInObjectiveValue := math.Abs(objectiveFunctionChange)
	c.AcceptanceProbability = math.Exp(-absoluteChangeInObjectiveValue / c.Temperature)
}

func (c *Coolant) CoolDown() {
	c.Temperature *= c.CoolingFactor
}
