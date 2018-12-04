// Copyright (c) 2018 Australian Rivers Institute.

package explorer

import (
	"errors"
	"math"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
)

const guaranteed = 1

type KirkpatrickExplorer struct {
	BaseExplorer
	parameters KirkpatrickExplorerParameters
	model      model.Model

	decisionVariable      model.DecisionVariable
	optimisationDirection optimisationDirection
}

func New() *KirkpatrickExplorer {
	newExplorer := new(KirkpatrickExplorer)
	newExplorer.parameters.Initialise()
	return newExplorer
}

func (ke *KirkpatrickExplorer) Initialise() {
	ke.BaseExplorer.Initialise()
}

func (ke *KirkpatrickExplorer) WithModel(model model.Model) *KirkpatrickExplorer {
	ke.model = model
	return ke
}

func (ke *KirkpatrickExplorer) WitParameters(params parameters.Map) *KirkpatrickExplorer {
	ke.parameters.Merge(params)

	ke.setOptimisationDirectionFromParams()
	ke.setDecisionVariableFromParams()

	return ke
}

func (ke *KirkpatrickExplorer) setOptimisationDirectionFromParams() {
	optimisationDirectionParam := ke.parameters.GetString(OptimisationDirection)
	ke.optimisationDirection, _ = parseOptimisationDirection(optimisationDirectionParam)
}

func (ke *KirkpatrickExplorer) setDecisionVariableFromParams() {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	decisionVariable, dvError := ke.model.DecisionVariable(decisionVariableName)
	if dvError == nil {
		ke.decisionVariable = decisionVariable
	} else {
		ke.decisionVariable = model.NullDecisionVariable
		panic("Decision variable [" + decisionVariableName + "] not recognised by model") // TODO: log instead
	}
}

func (ke *KirkpatrickExplorer) ParameterErrors() error {
	return ke.parameters.ValidationErrors()
}

func (ke *KirkpatrickExplorer) ObjectiveValue() float64 {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	if dv, dvError := ke.model.DecisionVariable(decisionVariableName); dvError == nil {
		return dv.Value()
	}
	return 0
}

func (ke *KirkpatrickExplorer) TryRandomChange(temperature float64) {
	ke.model.TryRandomChange()
	ke.acceptOrRevertChange(temperature)
}

func (ke *KirkpatrickExplorer) acceptOrRevertChange(annealingTemperature float64) {
	if ke.ChangeTriedIsDesirable() {
		ke.SetAcceptanceProbability(guaranteed)
		ke.AcceptLastChange()
	} else {
		absoluteChangeInObjectiveValue := math.Abs(ke.ChangeInObjectiveValue())
		probabilityToAcceptBadChange := math.Exp(-absoluteChangeInObjectiveValue / annealingTemperature)
		ke.SetAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue := newRandomValue(ke.RandomNumberGenerator())
		if probabilityToAcceptBadChange > randomValue {
			ke.AcceptLastChange()
		} else {
			ke.RevertLastChange()
		}
	}
}

func (ke *KirkpatrickExplorer) ChangeTriedIsDesirable() bool {
	switch ke.optimisationDirection {
	case Minimising:
		ke.changeIsDesirable = ke.changeInObjectiveValue() < 0
		return ke.changeIsDesirable
	case Maximising:
		ke.changeIsDesirable = ke.changeInObjectiveValue() > 0
		return ke.changeIsDesirable
	}
	return false
}

func (ke *KirkpatrickExplorer) ChangeIsDesirable() bool {
	return ke.changeIsDesirable
}

func (ke *KirkpatrickExplorer) changeInObjectiveValue() float64 {
	if change, changeError := ke.model.Change(ke.decisionVariable); changeError == nil {
		ke.SetChangeInObjectiveValue(change)
		return change
	}
	return 0
}

func (ke *KirkpatrickExplorer) AcceptLastChange() {
	ke.model.AcceptChange()
	ke.changeAccepted = true
}

func (ke *KirkpatrickExplorer) RevertLastChange() {
	ke.model.RevertChange()
	ke.changeAccepted = false
}

func (ke *KirkpatrickExplorer) SetAcceptanceProbability(probability float64) {
	ke.acceptanceProbability = math.Min(guaranteed, probability)
}

func (ke *KirkpatrickExplorer) Clone() Explorer {
	clone := *ke
	return &clone
}

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

type KirkpatrickExplorerParameters struct {
	parameters.Parameters
}

func (kp *KirkpatrickExplorerParameters) Initialise() *KirkpatrickExplorerParameters {
	kp.Parameters.Initialise()
	kp.buildMetaData()
	kp.CreateDefaults()
	return kp
}

func (kp *KirkpatrickExplorerParameters) buildMetaData() {
	kp.AddMetaData(
		parameters.MetaData{
			Key:          DecisionVariableName,
			Validator:    kp.Parameters.IsString,
			DefaultValue: model.ObjectiveValue,
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

func (kp *KirkpatrickExplorerParameters) isOptimisationDirection(key string, value interface{}) bool {
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
