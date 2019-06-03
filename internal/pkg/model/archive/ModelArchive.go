// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"github.com/LindsayBradford/crem/pkg/archive"
	"github.com/LindsayBradford/crem/pkg/dominance"
)

func New() *ModelArchive {
	archive := new(ModelArchive)
	return archive
}

type ModelArchive struct {
	variableVector dominance.Float64Vector
	actionArchive  archive.BooleanArchive
}

func (ma *ModelArchive) WithVariables(vector dominance.Float64Vector) *ModelArchive {
	ma.variableVector = vector
	return ma
}

func (ma *ModelArchive) WithActions(actionArchive archive.BooleanArchive) *ModelArchive {
	ma.actionArchive = actionArchive
	return ma
}

func (ma *ModelArchive) Variables() *dominance.Float64Vector {
	return &ma.variableVector
}

func (ma *ModelArchive) Actions() *archive.BooleanArchive {
	return &ma.actionArchive
}
