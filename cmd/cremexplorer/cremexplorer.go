//go:build windows
// +build windows

// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremexplorer/bootstrap"
	"github.com/LindsayBradford/crem/cmd/cremexplorer/commandline"
)

func main() {
	args := commandline.ParseArguments()
	bootstrap.RunExcelCompatibleScenarioFromConfigFile(args.ScenarioFile)
}
