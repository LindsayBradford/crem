// Copyright (c) 2019 Australian Rivers Institute.

package data

import "github.com/LindsayBradford/crem/internal/pkg/parameters"

type ObserverConfig struct {
	LoggingType          LoggerType
	LoggingFormatter     FormatterType
	ObserverType         ScenarioObserverType
	LogLevelDestinations map[string]string
	Parameters           parameters.Map
}

type LoggerType struct {
	value string
}

var (
	UnspecifiedLoggerType = LoggerType{""}
	NativeLibrary         = LoggerType{"NativeLibrary"}
	BareBones             = LoggerType{"BareBones"}
)

func (lt *LoggerType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "[[Loggers]].Type",
		validValues: []string{
			NativeLibrary.value, BareBones.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			lt.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}

type FormatterType struct {
	value string
}

var (
	UnspecifiedFormatterType = FormatterType{""}
	RawMessage               = FormatterType{"RawMessage"}
	Json                     = FormatterType{"JSON"}
	NameValuePair            = FormatterType{"NameValuePair"}
)

func (ft *FormatterType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "[[Loggers]].Formatter",
		validValues: []string{
			RawMessage.value, Json.value, NameValuePair.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			ft.value = string(text)
		},
	}

	return processUnmarshalContext(context)
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
