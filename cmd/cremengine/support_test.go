// +build windows
// Copyright (c) 2019 Australian Rivers Institute.

package main

import (
	"os"
	"testing"

	"github.com/onsi/gomega/gexec"
)

var cremExceutablePath string

const withFailure = 1

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
		os.Exit(withFailure)
	}
}

func tearDown() {
	gexec.CleanupBuildArtifacts()
}
