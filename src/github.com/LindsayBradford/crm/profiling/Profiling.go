// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package profiling

import (
	"os"
	"runtime/pprof"

	"github.com/LindsayBradford/crm/logging/handlers"
)

type NoParameterFunction func()

// ProfileIfRequired establishes profiling based on what is passed as the cpuProfile
// parameter. If an empty string, it's assumed profiling is not needed.  A non-empty
// string is asssumed to contain a path to a file in which profiling data is to be collated.

func ProfileIfRequired(cpuProfilePath string, humanLogHandler handlers.LogHandler, functionToProfile NoParameterFunction) {
	if cpuProfilePath != "" {
		f, err := os.Create(cpuProfilePath)
		if err != nil {
			humanLogHandler.ErrorWithError(err)
			os.Exit(1)
		}
		humanLogHandler.Info("About to profile cpu data to [" + cpuProfilePath + "]")

		pprof.StartCPUProfile(f)
		functionToProfile()
		defer pprof.StopCPUProfile()
	} else {
		functionToProfile()
	}
}