// Copyright (c) 2019 Australian Rivers Institute.

package variable

type UnitOfMeasure string

func (uom UnitOfMeasure) String() string {
	return string(uom)
}

type UnitOfMeasureContainer interface {
	UnitOfMeasure() UnitOfMeasure
}

type ContainedUnitOfMeasure struct {
	unitOfMeasure UnitOfMeasure
}

func (c *ContainedUnitOfMeasure) UnitOfMeasure() UnitOfMeasure {
	return c.unitOfMeasure
}

func (c *ContainedUnitOfMeasure) SetUnitOfMeasure(unitOfMeasure UnitOfMeasure) {
	c.unitOfMeasure = unitOfMeasure
}
