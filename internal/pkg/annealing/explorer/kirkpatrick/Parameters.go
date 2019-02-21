// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/pkg/errors"
)

const (
	DecisionVariableName  = "DecisionVariable"
	OptimisationDirection = "OptimisationDirection"
)

type optimisationDirection int

const (
	Invalid    optimisationDirection = iota
	Minimising optimisationDirection = iota
	Maximising optimisationDirection = iota
)

func (od optimisationDirection) String() string {
	switch od {
	case Minimising:
		return "Minimising"
	case Maximising:
		return "Maximising"
	default:
		return "Minimising"
	}
}

type Parameters struct {
	parameters.Parameters
}

func (kp *Parameters) Initialise() *Parameters {
	kp.Parameters.Initialise()
	kp.buildMetaData()
	kp.CreateDefaults()
	return kp
}

func (kp *Parameters) buildMetaData() {
	kp.AddMetaData(
		parameters.MetaData{
			Key:          DecisionVariableName,
			Validator:    kp.Parameters.IsString,
			DefaultValue: variable.ObjectiveValue,
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          OptimisationDirection,
			Validator:    kp.isOptimisationDirection,
			DefaultValue: Minimising.String(),
		},
	)
}

func (kp *Parameters) isOptimisationDirection(key string, value interface{}) bool {
	valueAsString, typeIsOk := value.(string)
	if !typeIsOk {
		kp.Parameters.AddValidationErrorMessage("Parameter [" + key + "] must be a string value")
		return false
	}

	_, parsingError := parseOptimisationDirection(valueAsString)
	return parsingError == nil
}

func parseOptimisationDirection(value string) (optimisationDirection, error) {
	directions := []optimisationDirection{Minimising, Maximising}

	for _, direction := range directions {
		if direction.String() == value {
			return direction, nil
		}
	}

	return Invalid, errors.New("not an optimisationDirection")
}
