// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package strings offers an extension of functionality to the default golang strings package.
package strings

import "strings"

// FluentBuilder offers a 'fluent' interface (see https://en.wikipedia.org/wiki/Fluent_interface) to the golang
// strings.Builder, allowing us to chain .Add() calls
type FluentBuilder struct {
	strings.Builder
}

// Add appends the contents of each supplied strings to its buffer, returning a reference to the FluentBuilder,
// allowing a chaining of a number of Add() calls
func (fb *FluentBuilder) Add(strings ...string) *FluentBuilder {
	fb.writeStrings(strings...)
	return fb
}

// AddIf appends the contents of each supplied strings to its buffer, returning a reference to the FluentBuilder,
// only when shouldAdd is true, allowing a chaining of a number of Add() calls
func (fb *FluentBuilder) AddIf(shouldAdd bool, strings ...string) *FluentBuilder {
	if !shouldAdd {
		return fb
	}
	fb.writeStrings(strings...)
	return fb
}

// writeStrings appends the contents of each supplied strings to its buffer via the Builder.WriteString() method
func (fb *FluentBuilder) writeStrings(strings ...string) (int, error) {
	var fullLength = 0
	for _, str := range strings {
		strLength, _ := fb.WriteString(str)
		fullLength += strLength
	}
	return fullLength, nil
}
