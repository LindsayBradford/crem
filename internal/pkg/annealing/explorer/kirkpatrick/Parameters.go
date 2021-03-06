// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"fmt"

	"github.com/LindsayBradford/crem/internal/pkg/parameters"

	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
	"github.com/pkg/errors"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.
		Initialise("Kirkpatrick Explorer Parameter Validation").
		Enforcing(ParameterSpecifications())
	return p
}

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

func ParameterSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          DecisionVariableName,
			Validator:    IsString,
			DefaultValue: ObjectiveValue,
		},
	).Add(
		Specification{
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
		return NewInvalidSpecificationError("Parameter [" + key + "] must be a string value")
	}
	if _, parsingError := parseOptimisationDirection(valueAsString); parsingError == nil {
		return NewValidSpecificationError(key, value)
	} else {
		return NewInvalidSpecificationError(parsingError.Error())
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
