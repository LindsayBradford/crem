// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	returnCode := m.Run()
	tearDown()

	os.Exit(returnCode)
}

var excelHandlerUnderTest *Handler

func setup() {
	excelHandlerUnderTest = new(Handler).Initialise()
}

func tearDown() {
	excelHandlerUnderTest.Destroy()
}
