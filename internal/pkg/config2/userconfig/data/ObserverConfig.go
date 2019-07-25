// Copyright (c) 2019 Australian Rivers Institute.

package data

import "github.com/LindsayBradford/crem/internal/pkg/parameters"

type ObserverConfig struct {
	ObserverType ScenarioObserverType
	Parameters   parameters.Map

	LoggingConfig
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
	context := unmarshalContext{
		configKey: "[Scenario.Observer].Type",
		validValues: []string{
			AttributeObserver.value, MessageObserver.value, InvariantObserver.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			ot.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}
