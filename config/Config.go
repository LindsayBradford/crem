// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package config contains configuration global to the Catchment Resilience Modelling tool.
package config

import (
	"errors"

	"github.com/BurntSushi/toml"
)

// Version number of the Catchment Resilience Modelling tool
const VERSION = "0.1.1"

// Compile-time debug flag.
const DEBUG = false

type CRMConfig struct {
	Title              string
	FilePath           string
	Annealer           AnnealingConfig
	Loggers            []LoggerConfig
	AnnealingObservers []AnnealingObserverConfig
	SolutionExplorers  []SolutionExplorerConfig
}

type EventNotifierType string

const (
	Unspecified EventNotifierType = ""
	Synchronous EventNotifierType = "Synchronous"
	Concurrent  EventNotifierType = "Concurrent"
)

type AnnealingConfig struct {
	Type                string
	StartingTemperature float64
	CoolingFactor       float64
	MaximumIterations   uint64
	EventNotifier       EventNotifierType
	SolutionExplorer    string
}

type LoggerConfig struct {
	Name                 string
	Type                 string
	Formatter            string
	LogLevelDestinations map[string]string
	Default              bool
}

type AnnealingObserverConfig struct {
	Type            string
	Logger          string
	IterationFilter string
	FilterRate      int64
}

type SolutionExplorerConfig struct {
	Type      string
	Name      string
	Penalty   float64
	InputFile string
}

func Retrieve(configFilePath string) *CRMConfig {
	var conf CRMConfig
	if _, err := toml.DecodeFile(configFilePath, &conf); err != nil {
		panic(errors.New("failed to retrieve config [" + configFilePath + "]"))
	}
	conf.FilePath = configFilePath
	return &conf
}
