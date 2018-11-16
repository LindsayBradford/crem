// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"os"
	"testing"

	"github.com/go-ole/go-ole"
)

func TestMain(m *testing.M) {
	setup()
	returnCode := m.Run()
	tearDown()

	os.Exit(returnCode)
}

var excelHandlerUnderTest *Handler

func setup() {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		panic(err)
	}
	excelHandlerUnderTest = new(Handler).Initialise()
}

func tearDown() {
	excelHandlerUnderTest.Destroy()
	ole.CoUninitialize()
}
