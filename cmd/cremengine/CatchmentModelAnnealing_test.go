// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var cremExceutablePath string

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup(m *testing.M) {
	var err error
	cremExceutablePath, err = gexec.Build("github.com/LindsayBradford/crem/cmd/cremengine")

	if err != nil {
		os.Exit(1)
	}

}

func tearDown() {
	gexec.CleanupBuildArtifacts()
}

const timeoutSeconds = 5

const withSuccess = 0
const withFailure = 1

func TestCremEngine_ScenarioOneRun_ExitWithSuccess(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	commandLineArguments := buildScenarioArguments("testdata/CatchmentConfig-OneRun.toml")

	// when
	session, err := startCremExecutableWith(commandLineArguments)

	// then
	g.Expect(err).ToNot(HaveOccurred())
	g.Eventually(session, timeoutSeconds).Should(gexec.Exit(withSuccess))
}

func TestCremEngine_ScenarioBadInputs_ExitWithError(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	commandLineArguments := buildScenarioArguments("testdata/CatchmentConfig-BadInputs.toml")

	// when
	session, err := startCremExecutableWith(commandLineArguments)

	// then
	g.Expect(err).ToNot(HaveOccurred())
	g.Eventually(session, timeoutSeconds).Should(gexec.Exit(withFailure))
}

func startCremExecutableWith(commandLineArguments []string) (*gexec.Session, error) {
	command := exec.Command(cremExceutablePath, commandLineArguments...)
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	return session, err
}

func buildScenarioArguments(scenarioFilePath string) []string {
	return []string{"--ScenarioFile", scenarioFilePath}
}
