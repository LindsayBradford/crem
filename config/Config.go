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
	Annealer                 AnnealingConfig
	loggerConfigList         []LoggerConfig
	annealingObserversConfig []AnnealingObserverConfig
}

type AnnealingConfig struct {
	Type                string
	StartingTemperature float64
	CoolingFactor       float64
	MaxIterations       uint64
	EventNotifier       string
}

type LoggerConfig struct {
}

type AnnealingObserverConfig struct {
}

func RetrieveConfig(configFilePath string) *CRMConfig {
	var conf CRMConfig
	if _, err := toml.DecodeFile(configFilePath, &conf); err != nil {
		panic(errors.New("failed to load config at[" + configFilePath + "]"))
	}
	return &conf
}
