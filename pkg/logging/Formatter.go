// Copyright (c) 2018 Australian Rivers Institute.

package logging

import "github.com/LindsayBradford/crem/pkg/attributes"

// Formatter describes an interface for the formatters of Attributes into some observer-ready string.
// Instances of Logger are expected to delegate any formatting of the supplied attributes to a Formatter.
type Formatter interface {
	// Format converts the supplied attributes into a representative 'observer ready' string.
	Format(attributes attributes.Attributes) string
}
