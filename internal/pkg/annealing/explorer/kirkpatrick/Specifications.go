// Copyright (c) 2019 Australian Rivers Institute.

package kirkpatrick

import (
	"fmt"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/pkg/errors"
)

const (
	DecisionVariableName  = "DecisionVariable"
	OptimisationDirection = "OptimisationDirection"
)

type optimisationDirection int

const (
	Invalid optimisationDirection = iota
	Minimising
	Maximising
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

func DefineSpecifications() *specification.Specifications {
	specs := specification.New()
	specs.Add(
		specification.Specification{
			Key:          DecisionVariableName,
			Validator:    specification.IsString,
			DefaultValue: variable.ObjectiveValue,
		},
	).Add(
		specification.Specification{
			Key:          OptimisationDirection,
			Validator:    isOptimisationDirection,
			DefaultValue: Minimising.String(),
		},
	)
	return specs
}

func isOptimisationDirection(key string, value interface{}) error {
	valueAsString, typeIsOk := value.(string)
	if !typeIsOk {
		return specification.NewInvalidSpecificationError("Parameter [" + key + "] must be a string value")
	}
	if _, parsingError := parseOptimisationDirection(valueAsString); parsingError == nil {
		return specification.NewValidSpecificationError(key, value)
	} else {
		return specification.NewInvalidSpecificationError(parsingError.Error())
	}
}

func parseOptimisationDirection(value string) (optimisationDirection, error) {
	directions := []optimisationDirection{Minimising, Maximising}

	for _, direction := range directions {
		if value == direction.String() {
			return direction, nil
		}
	}

	errorMsg := fmt.Sprintf("Parameter value [%s] is not a valid OptimisationDirection, should be one of %v", value, directions)
	return Invalid, errors.New(errorMsg)
}
