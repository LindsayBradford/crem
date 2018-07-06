// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// shared package offers up definitions needed by more specialised annealing logging (sub-)packages that do not need
// to know about each-other, but do need a common frame of reference.
package shared

// NameValuePair is a struct allowing some Name as text to be associated with a matching Value.
type NameValuePair struct {
	Name  string
	Value interface{}
}

// LogAttributes is an array of name-value pairs that we want to log.
type LogAttributes []NameValuePair
