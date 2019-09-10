// Copyright (c) 2019 Australian Rivers Institute.

package variable

type Bounded interface {
	WithinBounds(value float64) bool
}

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
