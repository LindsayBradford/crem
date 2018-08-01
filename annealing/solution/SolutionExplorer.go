// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package solution

import (
	"errors"
	. "github.com/LindsayBradford/crm/logging/handlers"
	"math"
	"math/rand"
	"time"
)

type SolutionExplorer interface {
	Name() string
	SetName(name string)

	Initialise()
	TryRandomChange(temperature float64)

	ObjectiveValue() float64
	SetObjectiveValue(objectiveValue float64)

	ChangeInObjectiveValue() float64
	SetChangeInObjectiveValue(change float64)

	AcceptanceProbability() float64
	SetAcceptanceProbability(probability float64)

	ChangeIsDesirable() bool
	AcceptLastChange()
	ChangeAccepted() bool
	RevertLastChange()

	SetRandomNumberGenerator(*rand.Rand)
	RandomNumberGenerator() *rand.Rand

	SetLogHandler(logger LogHandler) error
	LogHandler() LogHandler

	TearDown()
}

type BaseSolutionExplorer struct {
	name                   string
	objectiveValue         float64
	changeInObjectiveValue float64
	changeIsDesirable      bool
	changeAccepted         bool
	acceptanceProbability  float64
	randomNumberGenerator  *rand.Rand
	logHandler             LogHandler
}

func (explorer *BaseSolutionExplorer) Initialise() {
	explorer.logHandler.Debug("Initialising Solution Explorer")
	explorer.SetRandomNumberGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func (explorer *BaseSolutionExplorer) TearDown() {
	explorer.logHandler.Debug("Triggering tear-down of Solution Explorer")
}

func (explorer *BaseSolutionExplorer) RandomNumberGenerator() *rand.Rand {
	return explorer.randomNumberGenerator
}

func (explorer *BaseSolutionExplorer) SetRandomNumberGenerator(generator *rand.Rand) {
	explorer.randomNumberGenerator = generator
}

func (explorer *BaseSolutionExplorer) Name() string {
	return explorer.name
}

func (explorer *BaseSolutionExplorer) SetName(name string) {
	explorer.name = name
}

func (explorer *BaseSolutionExplorer) WithName(name string) *BaseSolutionExplorer {
	explorer.name = name
	return explorer
}

func (explorer *BaseSolutionExplorer) LogHandler() LogHandler {
	return explorer.logHandler
}

func (explorer *BaseSolutionExplorer) SetLogHandler(logHandler LogHandler) error {
	if logHandler == nil {
		return errors.New("invalid attempt to set log handler to nil value")
	}
	explorer.logHandler = logHandler
	return nil
}

func (explorer *BaseSolutionExplorer) TryRandomChange(temperature float64) {
	explorer.makeRandomChange()
	DecideOnWhetherToAcceptChange(explorer, temperature)
}

func (explorer *BaseSolutionExplorer) SetObjectiveValue(objectiveValue float64) {
	explorer.objectiveValue = objectiveValue
}

func (explorer *BaseSolutionExplorer) ObjectiveValue() float64 {
	return explorer.objectiveValue
}

func (explorer *BaseSolutionExplorer) ChangeInObjectiveValue() float64 {
	return explorer.changeInObjectiveValue
}

func (explorer *BaseSolutionExplorer) SetChangeInObjectiveValue(change float64) {
	explorer.changeInObjectiveValue = change
}

func (explorer *BaseSolutionExplorer) AcceptanceProbability() float64 {
	return explorer.acceptanceProbability
}

func (explorer *BaseSolutionExplorer) SetAcceptanceProbability(probability float64) {
	explorer.acceptanceProbability = probability
}

func (explorer *BaseSolutionExplorer) makeRandomChange() {}

func DecideOnWhetherToAcceptChange(explorer SolutionExplorer, annealingTemperature float64) {
	if explorer.ChangeIsDesirable() {
		explorer.SetAcceptanceProbability(1)
		explorer.AcceptLastChange()
	} else {
		probabilityToAcceptBadChange := math.Exp(-explorer.ChangeInObjectiveValue() / annealingTemperature)
		explorer.SetAcceptanceProbability(probabilityToAcceptBadChange)

		randomValue := newRandomValue(explorer.RandomNumberGenerator())
		if probabilityToAcceptBadChange > randomValue {
			explorer.AcceptLastChange()
		} else {
			explorer.RevertLastChange()
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

func (explorer *BaseSolutionExplorer) ChangeIsDesirable() bool {
	if explorer.changeInObjectiveValue <= 0 {
		return true
	}
	return false
}

func (explorer *BaseSolutionExplorer) AcceptLastChange() {
	explorer.changeAccepted = true
}

func (explorer *BaseSolutionExplorer) RevertLastChange() {
	explorer.changeAccepted = false
}

func (explorer *BaseSolutionExplorer) ChangeAccepted() bool {
	return explorer.changeAccepted
}
