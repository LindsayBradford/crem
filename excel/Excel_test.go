// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

var excelHandlerUnderTest *ExcelHandler

func setup() {
	excelHandlerUnderTest = InitialiseHandler()
}

func tearDown() {
	excelHandlerUnderTest.Destroy()
}
