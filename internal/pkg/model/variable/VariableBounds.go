// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"fmt"
	strings2 "strings"

	"github.com/LindsayBradford/crem/pkg/strings"
)

type Bounded interface {
	WithinBounds(value float64) bool
	BoundErrorAsText(value float64) string
}

var converter = strings.NewConverter().Localised().WithFloatingPointPrecision(6).PaddingZeros()

var _ Bounded = new(VariableBounds)

type VariableBounds struct {
	hasMinimum bool
	minimum    float64

	hasMaximum bool
	maximum    float64
}

func (vb *VariableBounds) SetMinimum(minimum float64) {
	vb.hasMinimum = true
	vb.minimum = minimum
}

func (vb *VariableBounds) SetMaximum(maximum float64) {
	vb.hasMaximum = true
	vb.maximum = maximum
}

func (vb *VariableBounds) WithinBounds(value float64) bool {
	if vb.hasMinimum && value < vb.minimum {
		return false
	}

	if vb.hasMaximum && value > vb.maximum {
		return false
	}

	return true
}

func (vb *VariableBounds) BoundErrorAsText(value float64) string {

	boundMessages := make([]string, 0)

	if vb.hasMinimum && value < vb.minimum {
		lowerBoundAsString := converter.Convert(vb.minimum)
		boundMessages = append(boundMessages, fmt.Sprintf("< lower bound %s", lowerBoundAsString))
	}

	if vb.hasMaximum && value > vb.maximum {
		upperBoundAsString := converter.Convert(vb.maximum)
		boundMessages = append(boundMessages, fmt.Sprintf("> upper bound %s", upperBoundAsString))
	}

	if len(boundMessages) > 0 {
		boundMessagesAsString := strings2.Join(boundMessages, ", ")
		return fmt.Sprintf("%s %s", converter.Convert(value), boundMessagesAsString)
	}
	return ""
}
