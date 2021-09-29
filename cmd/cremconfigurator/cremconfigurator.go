//go:build windows
// +build windows

// Copyright (c) 2021 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremconfigurator/view"
)

func main() {
	window := view.BuildWindow()
	window.ShowAndRun()
}
