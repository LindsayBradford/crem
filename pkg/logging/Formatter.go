// Copyright (c) 2018 Australian Rivers Institute.

package logging

// Formatter describes an interface for the formatters of Attributes into some observer-ready string.
// Instances of LogHandler are expected to delegate any formatters of the supplied attributes to a Formatter.
type Formatter interface {
	// Format converts the supplied attributes into a representative 'observer ready' string.
	Format(attributes Attributes) string
}
