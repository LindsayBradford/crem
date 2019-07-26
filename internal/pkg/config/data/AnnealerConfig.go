// Copyright (c) 2019 Australian Rivers Institute.

package data

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
)

type AnnealerConfig struct {
	Type          AnnealerType
	EventNotifier EventNotifierType
	Parameters    parameters.Map
}

type AnnealerType struct {
	Value string
}

func (at AnnealerType) String() string {
	return string(at.Value)
}

var (
	UnspecifiedAnnealerType = AnnealerType{""}
	Kirkpatrick             = AnnealerType{"Kirkpatrick"}
	Suppapitnarm            = AnnealerType{"Suppapitnarm"}
)

func (at *AnnealerType) UnmarshalText(text []byte) error {
	context := UnmarshalContext{
		ConfigKey: "Annealer.Type",
		ValidValues: []string{
			Kirkpatrick.Value, Suppapitnarm.Value,
		},
		TextToValidate: string(text),
		AssignmentFunction: func() {
			at.Value = string(text)
		},
	}

	return ProcessUnmarshalContext(context)
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
	context := UnmarshalContext{
		ConfigKey: "Annealer.EventNotifier",
		ValidValues: []string{
			Sequential.value, Concurrent.value,
		},
		TextToValidate: string(text),
		AssignmentFunction: func() {
			ent.value = string(text)
		},
	}

	return ProcessUnmarshalContext(context)
}
