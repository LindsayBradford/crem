// Copyright (c) 2019 Australian Rivers Institute.

package formatters

import "github.com/LindsayBradford/crem/pkg/logging"

// nullFormatMessage is the message supplied if the NullFormatter is left as the Formatter of a LogHandler.
const nullFormatMessage = "No formatter specified. Using the NullFormatter."

// NullFormatter implements a 'null object' placeholder formatter that is supplied by default if one is not specified.
// It returns a static message as per nullFormatMessage as a reminder that a proper formatter must be supplied
// for the observer to do anything meaningful.
type NullFormatter struct{}

func (nf *NullFormatter) Format(attributes logging.Attributes) string {
	return nullFormatMessage
}
