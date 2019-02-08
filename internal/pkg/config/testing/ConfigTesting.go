// Copyright (c) 2018 Australian Rivers Institute.

package testing

import (
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const WithSuccess = 0
const WithFailure = 1

type BinaryTestingContext struct {
	Name              string
	ExecutablePath    string
	T                 *testing.T
	ConfigFilePath    string
	TimeoutSeconds    uint
	ExpectedErrorCode int
}

func TestExecutableAgainstConfigFile(tc BinaryTestingContext) {
	if testing.Short() {
		tc.T.Skip("skipping " + tc.Name + " in short mode")
	}

	g := NewGomegaWithT(tc.T)

	// given
	commandLineArguments := buildScenarioArguments(tc.ConfigFilePath)

	// when
	session, err := startExecutableWithArguments(tc.ExecutablePath, commandLineArguments)

	// then
	g.Expect(err).ToNot(HaveOccurred())
	g.Eventually(session, tc.TimeoutSeconds).Should(gexec.Exit(tc.ExpectedErrorCode))
}

func startExecutableWithArguments(executablePath string, commandLineArguments []string) (*gexec.Session, error) {
	command := exec.Command(executablePath, commandLineArguments...)
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	return session, err
}

func buildScenarioArguments(scenarioFilePath string) []string {
	return []string{"--ScenarioFile", scenarioFilePath}
}

type ScenarioFileRunningFunction func(scenarioPath string)

type Context struct {
	Name           string
	T              *testing.T
	ConfigFilePath string
	Runner         ScenarioFileRunningFunction
}
