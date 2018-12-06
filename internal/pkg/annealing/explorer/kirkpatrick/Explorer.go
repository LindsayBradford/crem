// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/logging"
)

const guaranteed = 1

type Explorer struct {
	name       string
	scenarioId string
	model      model.Model

	parameters            Parameters
	optimisationDirection optimisationDirection

	randomNumberGenerator *rand.Rand
	acceptanceProbability float64
	changeIsDesirable     bool
	changeAccepted        bool
	objectiveValueChange  float64
	logHandler            logging.Logger
}

func New() *Explorer {
	newExplorer := new(Explorer)
	newExplorer.parameters.Initialise()
	newExplorer.SetModel(model.NullModel)
	return newExplorer
}

func (ke *Explorer) Initialise() {
	ke.logHandler.Debug(ke.scenarioId + ": Initialising Solution Explorer")
	ke.SetRandomNumberGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func (ke *Explorer) RandomNumberGenerator() *rand.Rand {
	return ke.randomNumberGenerator
}

func (ke *Explorer) SetRandomNumberGenerator(generator *rand.Rand) {
	ke.randomNumberGenerator = generator
}

func (ke *Explorer) Name() string {
	return ke.name
}

func (ke *Explorer) SetName(name string) {
	ke.name = name
}

func (ke *Explorer) WithName(name string) *Explorer {
	ke.SetName(name)
	return ke
}

func (ke *Explorer) Model() model.Model {
	return ke.model
}

func (ke *Explorer) SetModel(model model.Model) {
	ke.model = model
}

func (ke *Explorer) WithModel(model model.Model) *Explorer {
	ke.model = model
	return ke
}

func (ke *Explorer) ScenarioId() string {
	return ke.scenarioId
}

func (ke *Explorer) SetScenarioId(id string) {
	ke.scenarioId = id
}

func (ke *Explorer) WithScenarioId(id string) *Explorer {
	ke.scenarioId = id
	return ke
}

func (ke *Explorer) WithParameters(params parameters.Map) *Explorer {
	ke.parameters.Merge(params)

	ke.setOptimisationDirectionFromParams()
	ke.checkDecisionVariableFromParams()

	return ke
}

func (ke *Explorer) setOptimisationDirectionFromParams() {
	optimisationDirectionParam := ke.parameters.GetString(OptimisationDirection)
	ke.optimisationDirection, _ = parseOptimisationDirection(optimisationDirectionParam)
}

func (ke *Explorer) checkDecisionVariableFromParams() {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	_, dvError := ke.Model().DecisionVariable(decisionVariableName)
	if dvError != nil {
		ke.parameters.AddValidationErrorMessage("decision variable [" + decisionVariableName + "] not recognised by model")
	}
}

func (ke *Explorer) ParameterErrors() error {
	return ke.parameters.ValidationErrors()
}

func (ke *Explorer) ObjectiveValue() float64 {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	if dv, dvError := ke.Model().DecisionVariable(decisionVariableName); dvError == nil {
		return dv.Value()
	}
	return 0
}

func (ke *Explorer) TryRandomChange(temperature float64) {
	ke.Model().TryRandomChange()
	ke.acceptOrRevertChange(temperature)
}

func (ke *Explorer) acceptOrRevertChange(annealingTemperature float64) {
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

// newRandomValue returns the next random number in the range [0,1] from the supplied randomNumberGenerator.
// (which by default returns a random number in the range [0,1).
// See: http://mumble.net/~campbell/2014/04/28/uniform-random-float
func newRandomValue(randomNumberGenerator *rand.Rand) float64 {
	distributionRange := int64(math.Pow(2, 53))
	return float64(randomNumberGenerator.Int63n(distributionRange)) / float64(distributionRange-1)
}

func (ke *Explorer) ChangeTriedIsDesirable() bool {
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

func (ke *Explorer) ChangeIsDesirable() bool {
	return ke.changeIsDesirable
}

func (ke *Explorer) changeInObjectiveValue() float64 {
	decisionVariableName := ke.parameters.GetString(DecisionVariableName)
	if change, changeError := ke.Model().DecisionVariableChange(decisionVariableName); changeError == nil {
		ke.SetChangeInObjectiveValue(change)
		return change
	}
	return 0
}

func (ke *Explorer) ChangeInObjectiveValue() float64 {
	return ke.objectiveValueChange
}

func (ke *Explorer) SetChangeInObjectiveValue(change float64) {
	ke.objectiveValueChange = change
}

func (ke *Explorer) AcceptLastChange() {
	ke.Model().AcceptChange()
	ke.changeAccepted = true
}

func (ke *Explorer) RevertLastChange() {
	ke.Model().RevertChange()
	ke.changeAccepted = false
}

func (ke *Explorer) ChangeAccepted() bool {
	return ke.changeAccepted
}

func (ke *Explorer) AcceptanceProbability() float64 {
	return ke.acceptanceProbability
}

func (ke *Explorer) SetAcceptanceProbability(probability float64) {
	ke.acceptanceProbability = math.Min(guaranteed, probability)
}

func (ke *Explorer) LogHandler() logging.Logger {
	return ke.logHandler
}

func (ke *Explorer) SetLogHandler(logHandler logging.Logger) error {
	if logHandler == nil {
		return errors.New("invalid attempt to set log handler to nil value")
	}
	ke.logHandler = logHandler
	return nil
}

func (ke *Explorer) Clone() explorer.Explorer {
	clone := *ke
	modelClone := ke.Model().Clone()
	clone.SetModel(modelClone)
	return &clone
}

func (ke *Explorer) TearDown() {
	ke.logHandler.Debug(ke.scenarioId + ": Triggering tear-down of Solution Explorer")
}
