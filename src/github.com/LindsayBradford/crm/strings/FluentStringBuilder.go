package strings

import . "strings"

// A Fluent wrapper to Go strings.Builder, allowing us to chain .Add() calls

type FluentBuilder struct {
	Builder
}

// WriteStrings appends the contents of each s to this buffer.
// via the Builder.WriteString() method
func (this *FluentBuilder) WriteStrings(strings ...string) (int, error) {
	var fullLength = 0
	for _, str := range strings {
		strLength, _ := this.WriteString(str)
		fullLength += strLength
	}
	return fullLength, nil
}

// A string building method, making use of WriteStrings() to allow
// related strings to be built per Add() call, and the Add() calls themselves
// to be chained

func (this *FluentBuilder) Add(str ...string) *FluentBuilder {
	this.WriteStrings(str...)
	return this
}
