// +build windows

// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	"github.com/LindsayBradford/crem/internal/pkg/commandline"
	"github.com/LindsayBradford/crem/pkg/logging"
)

var (
	defaultLogHandler logging.Logger
)

func main() {
	args := commandline.ParseArguments()

	if shouldRunScenario(args) {
		scenario.RunExcelCompatibleScenarioFromConfigFile(args.ScenarioFile)
	} else {
		components.RunServerFromConfigFile(args.ServerConfigFile)
	}
}

func shouldRunScenario(args *commandline.Arguments) bool {
	return args.ScenarioFile != ""
}
