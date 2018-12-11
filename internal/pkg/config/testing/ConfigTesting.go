// Copyright (c) 2018 Australian Rivers Institute.

package testing

import (
	"testing"

	. "github.com/onsi/gomega"
)

type Context struct {
	Name           string
	T              *testing.T
	ConfigFilePath string
	Runner         ScenarioFileRunningFunction
}

type ScenarioFileRunningFunction func(scenarioPath string)

func (tc *Context) VerifyScenarioRunViaConfigFileDoesNotPanic() {
	if testing.Short() {
		tc.T.Skip("skipping " + tc.Name + " in short mode")
	}
	g := NewGomegaWithT(tc.T)

	scenarioRun := func() {
		tc.Runner(tc.ConfigFilePath)
	}

	g.Expect(scenarioRun).To(Not(Panic()), tc.Name+" should not panic")
}
