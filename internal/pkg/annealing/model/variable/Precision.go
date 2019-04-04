// Copyright (c) 2019 Australian Rivers Institute.

package variable

type Precision int

type PrecisionContainer interface {
	Precision() Precision
}

type ContainedPrecision struct {
	precision Precision
}

func (c *ContainedPrecision) Precision() Precision {
	return c.precision
}

func (c *ContainedPrecision) SetPrecision(precision Precision) {
	c.precision = precision
}
