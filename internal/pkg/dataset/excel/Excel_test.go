// Copyright (c) 2018 Australian Rivers Institute.

// +build windows
// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package excel

import (
	"os"
	"testing"

	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/go-ole/go-ole"
)

func TestMain(m *testing.M) {
	setup()
	returnCode := m.Run()
	tearDown()

	os.Exit(returnCode)
}

var (
	mainThreadChannel = make(threading.MainThreadChannel)
)

func setup() {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		panic(err)
	}

	go mainThreadFunctionHandler()
}

func tearDown() {
	ole.CoUninitialize()
	// threading.GetMainThreadChannel().Close()
}

func callOnMainThread(function threading.MainThreadFunction) {
	done := make(chan bool, 1)
	mainThreadChannel <- func() {
		function()
		done <- true
	}
	<-done
}

func mainThreadFunctionHandler() {
	for function := range mainThreadChannel {
		function()
	}
}
