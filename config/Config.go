// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package config contains configuration global to the Catchment Resilience Modelling tool.
package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/LindsayBradford/crm/strings"
	"github.com/pkg/errors"
)

// Version number of the Catchment Resilience Modelling tool
const VERSION = "0.1.1"

func Retrieve(configFilePath string) (*CRMConfig, error) {
	var conf CRMConfig
	metaData, decodeErr := toml.DecodeFile(configFilePath, &conf)
	if decodeErr != nil {
		return nil, errors.Wrap(decodeErr, "failed retrieving config from file")
	}
	if len(metaData.Undecoded()) > 0 {
		errorMsg := fmt.Sprintf("unrecognised configuration key(s) %q", metaData.Undecoded())
		return nil, errors.New(errorMsg)
	}
	conf.FilePath = configFilePath
	return &conf, nil
}

type CRMConfig struct {
	FilePath string

	ScenarioName               string
	RunNumber                  uint64
	MaximumConcurrentRunNumber uint64
	CpuProfilePath             string

	Annealer           AnnealingConfig
	Loggers            []LoggerConfig
	AnnealingObservers []AnnealingObserverConfig
	SolutionExplorers  []SolutionExplorerConfig
}

type EventNotifierType struct {
	value string
}

var (
	UnspecifiedEventNotifierType = EventNotifierType{""}
	Sequential                   = EventNotifierType{"Sequential"}
	Concurrent                   = EventNotifierType{"Concurrent"}
)

type unmarshalContext struct {
	configKey          string
	validValues        []string
	textToValidate     string
	assignmentFunction func()
}

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

type AnnealerType struct {
	value string
}

var (
	UnspecifiedAnnealerType = AnnealerType{""}
	OSThreadLocked          = AnnealerType{"OSThreadLocked"}
	ElapsedTimeTracking     = AnnealerType{"ElapsedTimeTracking"}
	Simple                  = AnnealerType{"Simple"}
)

func (at *AnnealerType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "Annealer.Type",
		validValues: []string{
			OSThreadLocked.value, ElapsedTimeTracking.value, Simple.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			at.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}

type AnnealingConfig struct {
	Type                AnnealerType
	StartingTemperature float64
	CoolingFactor       float64
	MaximumIterations   uint64
	EventNotifier       EventNotifierType
	SolutionExplorer    string
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

type LoggerConfig struct {
	Name                 string
	Type                 LoggerType
	Formatter            FormatterType
	LogLevelDestinations map[string]string
}

type AnnealingObserverConfig struct {
	Type                   AnnealingObserverType
	Logger                 string
	IterationFilter        IterationFilter
	NumberOfIterations     uint64
	PercentileOfIterations float64
	SecondsBetweenEvents   uint64
}

type AnnealingObserverType struct {
	value string
}

var (
	UnspecifiedAnnealingObserverType = AnnealingObserverType{""}
	AttributeObserver                = AnnealingObserverType{"AttributeObserver"}
	MessageObserver                  = AnnealingObserverType{"MessageObserver"}
)

func (ot *AnnealingObserverType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "[[AnnealingObservers]].Type",
		validValues: []string{
			AttributeObserver.value, MessageObserver.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			ot.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}

type IterationFilter struct {
	value string
}

var (
	UnspecifiedIterationFilter          = IterationFilter{""}
	EveryNumberOfIterations             = IterationFilter{"EveryNumberOfIterations"}
	EveryElapsedSeconds                 = IterationFilter{"EveryElapsedSeconds"}
	EveryPercentileOfFinishedIterations = IterationFilter{"EveryPercentileOfFinishedIterations"}
)

func (filter *IterationFilter) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "[[AnnealingObservers]].IterationFilter",
		validValues: []string{
			EveryNumberOfIterations.value,
			EveryElapsedSeconds.value,
			EveryPercentileOfFinishedIterations.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			filter.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}

type SolutionExplorerConfig struct {
	Type      string
	Name      string
	Penalty   float64
	InputFile string
}

func processUnmarshalContext(context unmarshalContext) error {
	if valueIsInList(context.textToValidate, context.validValues...) {
		context.assignmentFunction()
		return nil
	}
	return generateErrorFromContext(context)
}

func valueIsInList(value string, list ...string) bool {
	for _, listEntry := range list {
		if value == listEntry {
			return true
		}
	}
	return false
}

func generateErrorFromContext(context unmarshalContext) error {
	const errorTemplate = "invalid value \"%v\" specified for key \"%s\"; should be one of: %s"
	return fmt.Errorf(errorTemplate, context.textToValidate, context.configKey, listToString(context.validValues...))
}

func listToString(list ...string) string {
	builder := strings.FluentBuilder{}
	needsComma := false
	for _, entry := range list {
		if needsComma {
			builder.Add(", ")
		}

		builder.Add("\"", entry, "\"")
		needsComma = true
	}
	return builder.String()
}
