// Copyright (c) 2019 Australian Rivers Institute.

package kirkpatrick

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/rand"
)

type Coolant struct {
	rand.RandContainer

	AcceptanceProbability float64
	Temperature           float64
	CoolingFactor         float64
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
