// +build windows

// Copyright (c) 2020 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/bootstrap"
	"github.com/LindsayBradford/crem/cmd/cremengine/commandline"
)

func main() {
	args := commandline.ParseArguments()
	bootstrap.RunMainThreadBoundEngineFromConfigFile(args.EngineConfigFile)
}
