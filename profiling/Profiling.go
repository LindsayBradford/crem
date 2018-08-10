// Copyright (c) 2018 Australian Rivers Institute.

package profiling

import (
	"github.com/pkg/errors"
	"os"
	"runtime/pprof"
)

type ProfileableFunction func() error

type OptionalProfilingFunctionPair struct {
	UnProfiledFunction ProfileableFunction
	ProfiledFunction   ProfileableFunction
}

// CpuProfileOfFunctionToFile establishes profiling based on what is passed as the cpuProfile
// parameter. If an empty string, it's assumed profiling is not needed.  A non-empty
// string is assumed to contain a path to a file in which profiling data is to be collated.

func CpuProfileOfFunctionToFile(functionToProfile ProfileableFunction, cpuProfilePath string) error {
	if cpuProfilePath == "" {
		return errors.New("empty cpu profile path supplied")
	}

	fileHandle, createErr := os.Create(cpuProfilePath)
	if createErr != nil {
		return errors.Wrap(createErr, "creation of cpu profiling file failed")
	}

	pprof.StartCPUProfile(fileHandle)
	functionErr := functionToProfile()
	defer pprof.StopCPUProfile()

	if functionErr != nil {
		return errors.Wrap(functionErr, "cpu profiling function error")
	}

	return nil
}
