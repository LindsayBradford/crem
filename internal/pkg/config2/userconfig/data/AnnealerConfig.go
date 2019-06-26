// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package data

import "github.com/LindsayBradford/crem/internal/pkg/parameters"

type AnnealerConfig struct {
	Type          AnnealerType
	EventNotifier EventNotifierType
	Parameters    parameters.Map
}

type AnnealerType struct {
	value string
}

func (at AnnealerType) String() string {
	return string(at.value)
}

var (
	UnspecifiedAnnealerType = AnnealerType{""}
	Kirkpatrick             = AnnealerType{"Kirkpatrick"}
	Suppapitnarm            = AnnealerType{"Suppapitnarm"}
)

func (at *AnnealerType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "Annealer.Type",
		validValues: []string{
			Kirkpatrick.value, Suppapitnarm.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			at.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}

type EventNotifierType struct {
	value string
}

var (
	UnspecifiedEventNotifierType = EventNotifierType{""}
	Sequential                   = EventNotifierType{"Sequential"}
	Concurrent                   = EventNotifierType{"Concurrent"}
)

func (ent *EventNotifierType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "Annealer.EventNotifier",
		validValues: []string{
			Sequential.value, Concurrent.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			ent.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}
