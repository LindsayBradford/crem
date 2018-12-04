// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package explorer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/dumb"
)

func NewDumbExplorer() *DumbExplorer {
	explorer := new(DumbExplorer).WithModel(dumb.New())
	explorer.parameters.Initialise()
	return explorer
}

type DumbExplorer struct {
	KirkpatrickExplorer
}

func (nde *DumbExplorer) WithName(name string) *DumbExplorer {
	nde.KirkpatrickExplorer.WithName(name)
	return nde
}

func (nde *DumbExplorer) WithParameters(params parameters.Map) *DumbExplorer {
	nde.KirkpatrickExplorer.WitParameters(params)
	return nde
}

func (nde *DumbExplorer) ParameterErrors() error {
	return nde.parameters.ValidationErrors()
}

func (nde *DumbExplorer) WithModel(model model.Model) *DumbExplorer {
	nde.KirkpatrickExplorer.WithModel(model)
	return nde
}
