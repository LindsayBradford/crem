// Copyright (c) 2019 Australian Rivers Institute.

package averaged

import (
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/cooling"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
)

var _ cooling.TemperatureCoolant = NewCoolant()

func NewCoolant() *Coolant {
	newCoolant := new(Coolant)
	newCoolant.Initialise()
	newCoolant.parameters.Initialise()
	return newCoolant
}

type Coolant struct {
	rand.RandContainer
	parameters Parameters

	acceptanceProbability float64
	temperature           float64
	coolingFactor         float64
}

func (c *Coolant) Initialise() *Coolant {
	c.parameters.Initialise()
	return c
}

func (c *Coolant) WithParameters(params parameters.Map) *Coolant {
	c.SetParameters(params)
	return c
}

func (c *Coolant) SetParameters(params parameters.Map) error {
	c.parameters.AssignOnlyEnforcedUserValues(params)

	c.temperature = c.parameters.GetFloat64(StartingTemperature)
	c.coolingFactor = c.parameters.GetFloat64(CoolingFactor)

	return nil
}

func (c *Coolant) ParameterErrors() error {
	return c.parameters.ValidationErrors()
}

func (c *Coolant) DecideIfAcceptable(variableChanges []float64) bool {
	c.calculateAcceptanceProbability(variableChanges)
	randomValue := c.RandomNumberGenerator().Float64Unitary()
	return c.acceptanceProbability > randomValue
}

func (c *Coolant) calculateAcceptanceProbability(variableChanges []float64) {
	numberOfChanges := len(variableChanges)

	probabilities := make([]float64, numberOfChanges)
	for index, individualChange := range variableChanges {
		absoluteChangeInObjectiveValue := math.Abs(individualChange)
		probabilities[index] = math.Exp(-absoluteChangeInObjectiveValue / c.temperature)
	}

	finalProbability := float64(1)
	for _, probability := range probabilities {
		finalProbability = finalProbability + probability
	}

	numberOfChangesAsFloat := float64(numberOfChanges)
	finalProbability = finalProbability / numberOfChangesAsFloat

	c.acceptanceProbability = finalProbability
}

func (c *Coolant) Temperature() float64 {
	return c.temperature
}

func (c *Coolant) SetTemperature(temperature float64) error {
	c.temperature = temperature
	return nil // TODO: Why do I need an error return?
}

func (c *Coolant) CoolingFactor() float64 {
	return c.coolingFactor
}

func (c *Coolant) SetAcceptanceProbability(acceptanceProbability float64) {
	c.acceptanceProbability = acceptanceProbability
}

func (c *Coolant) AcceptanceProbability() float64 {
	return c.acceptanceProbability
}

func (c *Coolant) CoolDown() {
	c.temperature *= c.coolingFactor
}
