// Copyright (c) 2018 Australian Rivers Institute.

// threading supplies support in terms of forcing functions to run on the OS thread via a bound main goroutine
// (some behaviour, like OLE function invocations, demand this).
package threading

import "runtime"

// MainThreadFunction defines a function that should be run on the main goroutine
type MainThreadFunction func()

type MainThreadFunctionWrapper func(mainThreadFunction MainThreadFunction)

// MainThreadChannel defines a type of channel that should only be run on a goroutine locked to the  OS thread.
type MainThreadChannel chan MainThreadFunction

var (
	mainThreadChannel = make(MainThreadChannel)
)

// GetMainThreadChannel returns an internally initialised  MainThreadChannel and locks the main goroutine to the OS thread.
// Behaviour is not guaranteed if this is not called from the main go routine.
func GetMainThreadChannel() MainThreadChannel {
	runtime.LockOSThread()
	return mainThreadChannel
}

// RunHandler enters into a loop of receiving and executing functions for the channel. For expected behaviour, it should
// be the only function running on the main goroutine.
func (mtc MainThreadChannel) RunHandler() {
	for function := range mtc {
		function()
	}
}

// Call allows the invocation of a MainThreadFunction on the main goroutine via this channel.
func (mtc MainThreadChannel) Call(mainThreadFunction MainThreadFunction) {
	done := make(chan bool, 1)
	mtc <- func() {
		mainThreadFunction()
		done <- true
	}
	<-done
}

// Close will close communication with this channel, and unlock the main goroutine from the OS thread.
func (mtc MainThreadChannel) Close() {
	close(mtc)
	runtime.UnlockOSThread()
}
