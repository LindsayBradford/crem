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
	Title                    string
	FilePath                 string
	Annealer                 AnnealingConfig
	loggerConfigList         []LoggerConfig
	annealingObserversConfig []AnnealingObserverConfig
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
}

type LoggerConfig struct {
}

type AnnealingObserverConfig struct {
}

func Retrieve(configFilePath string) *CRMConfig {
	var conf CRMConfig
	if _, err := toml.DecodeFile(configFilePath, &conf); err != nil {
		panic(errors.New("failed to retrieve config [" + configFilePath + "]"))
	}
	conf.FilePath = configFilePath
	return &conf
}
