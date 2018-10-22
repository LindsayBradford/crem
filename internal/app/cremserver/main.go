// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/commandline"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components"
)

func main() {
	args := commandline.ParseArguments()
	if shouldRunScenario(args) {
		components.RunScenarioFromConfigFile(args.ScenarioFile)
	} else {
		components.RunServerFromConfigFile(args.ServerConfigFile)
	}
}
func shouldRunScenario(args *commandline.Arguments) bool {
	return args.ScenarioFile != ""
}
