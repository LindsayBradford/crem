// +build windows
// Copyright (c) 2019 Australian Rivers Institute.

package main

import (
	"testing"
)

func TestDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := BinaryTestingContext{
		Name:              "Single run of Dumb Annealer",
		T:                 t,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-OneRun.toml",
		ExpectedErrorCode: withSuccess,
	}

	testCremExecutableAgainstConfigFile(context)
}

func TestDumbAnnealerIntegrationThreeRunsSequentially(t *testing.T) {
	context := BinaryTestingContext{
		Name:              "Three sequential runs of Dumb Annealer",
		T:                 t,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-ThreeRunsSequentially.toml",
		ExpectedErrorCode: withSuccess,
	}

	testCremExecutableAgainstConfigFile(context)
}

func TestDumbAnnealerIntegrationThreeRunsConcurrently(t *testing.T) {
	context := BinaryTestingContext{
		Name:              "Three concurrent runs of Dumb Annealer",
		T:                 t,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-ThreeRunsConcurrently.toml",
		ExpectedErrorCode: withSuccess,
	}

	testCremExecutableAgainstConfigFile(context)
}

func TestKirkpatrickDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := BinaryTestingContext{
		Name:              "Single run of Kirkpatrick Dumb Annealer",
		T:                 t,
		ConfigFilePath:    "testdata/KirkpatrickDumbAnnealerTestConfig-OneRun.toml",
		ExpectedErrorCode: withSuccess,
	}

	testCremExecutableAgainstConfigFile(context)
}
