// Copyright (c) 2018 Australian Rivers Institute.

package testing

import (
	"testing"

	. "github.com/onsi/gomega"
)

type TestingContext struct {
	Name           string
	T              *testing.T
	ConfigFilePath string
	Runner         ScenarioFileRunningFunction
}

type ScenarioFileRunningFunction func(scenarioPath string)

func (tc *TestingContext) VerifyScenarioConfigFilesDoesNotPanic() {
	if testing.Short() {
		tc.T.Skip("skipping " + tc.Name + " in short mode")
	}
	g := NewGomegaWithT(tc.T)

	scenarioRun := func() {
		tc.Runner(tc.ConfigFilePath)
	}

	g.Expect(scenarioRun).To(Not(Panic()), tc.Name+" should not panic")
}
