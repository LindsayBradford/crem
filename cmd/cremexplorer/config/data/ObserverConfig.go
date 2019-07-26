// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
)

type ObserverConfig struct {
	Type       ScenarioObserverType
	Parameters parameters.Map

	data.LoggingConfig
}

type ScenarioObserverType struct {
	value string
}

var (
	UnspecifiedAnnealingObserverType = ScenarioObserverType{""}
	AttributeObserver                = ScenarioObserverType{"AttributeObserver"}
	MessageObserver                  = ScenarioObserverType{"MessageObserver"}
	InvariantObserver                = ScenarioObserverType{"InvariantObserver"}
)

func (ot *ScenarioObserverType) UnmarshalText(text []byte) error {
	context := data.UnmarshalContext{
		ConfigKey: "[Scenario.Observer].Type",
		ValidValues: []string{
			AttributeObserver.value, MessageObserver.value, InvariantObserver.value,
		},
		TextToValidate: string(text),
		AssignmentFunction: func() {
			ot.value = string(text)
		},
	}

	return data.ProcessUnmarshalContext(context)
}
