// Copyright (c) 2020 Australian Rivers Institute.

package cooling

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
)

type TemperatureCoolant interface {
	rand.Container
	parameters.Container

	SetTemperature(temperature float64) error
	Temperature() float64

	DecideIfAcceptable(variableChanges []float64) bool

	CoolingFactor() float64

	SetAcceptanceProbability(acceptanceProbability float64)
	AcceptanceProbability() float64

	CoolDown()
}
