// Copyright (c) 2018 Australian Rivers Institute.

package explorer

import (
	"math"
	"math/rand"
	"time"

	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/pkg/errors"
)

type BaseExplorer struct {
	name                   string
	scenarioId             string
	objectiveValue         float64
	changeInObjectiveValue float64
	changeIsDesirable      bool
	changeAccepted         bool
	acceptanceProbability  float64
	randomNumberGenerator  *rand.Rand
	logHandler             logging.Logger
}

func (explorer *BaseExplorer) Initialise() {
	explorer.logHandler.Debug(explorer.scenarioId + ": Initialising Solution Explorer")
	explorer.SetRandomNumberGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func (explorer *BaseExplorer) TearDown() {
	explorer.logHandler.Debug(explorer.scenarioId + ": Triggering tear-down of Solution Explorer")
}

func (explorer *BaseExplorer) RandomNumberGenerator() *rand.Rand {
	return explorer.randomNumberGenerator
}

func (explorer *BaseExplorer) SetRandomNumberGenerator(generator *rand.Rand) {
	explorer.randomNumberGenerator = generator
}

func (explorer *BaseExplorer) Name() string {
	return explorer.name
}

func (explorer *BaseExplorer) SetName(name string) {
	explorer.name = name
}

func (explorer *BaseExplorer) WithName(name string) *BaseExplorer {
	explorer.name = name
	return explorer
}

func (explorer *BaseExplorer) ScenarioId() string {
	return explorer.scenarioId
}

func (explorer *BaseExplorer) SetScenarioId(id string) {
	explorer.scenarioId = id
}

func (explorer *BaseExplorer) WithScenarioId(id string) *BaseExplorer {
	explorer.scenarioId = id
	return explorer
}

func (explorer *BaseExplorer) LogHandler() logging.Logger {
	return explorer.logHandler
}

func (explorer *BaseExplorer) SetLogHandler(logHandler logging.Logger) error {
	if logHandler == nil {
		return errors.New("invalid attempt to set log handler to nil value")
	}
	explorer.logHandler = logHandler
	return nil
}

func (explorer *BaseExplorer) TryRandomChange(temperature float64) {}

func (explorer *BaseExplorer) SetObjectiveValue(objectiveValue float64) {
	explorer.objectiveValue = objectiveValue
}

func (explorer *BaseExplorer) ObjectiveValue() float64 {
	return explorer.objectiveValue
}

func (explorer *BaseExplorer) ChangeInObjectiveValue() float64 {
	return explorer.changeInObjectiveValue
}

func (explorer *BaseExplorer) SetChangeInObjectiveValue(change float64) {
	explorer.changeInObjectiveValue = change
}

func (explorer *BaseExplorer) AcceptanceProbability() float64 {
	return explorer.acceptanceProbability
}

func (explorer *BaseExplorer) SetAcceptanceProbability(probability float64) {
	explorer.acceptanceProbability = probability
}

func (explorer *BaseExplorer) DecideOnWhetherToAcceptChange(annealingTemperature float64, acceptFunction func(), rejectFunction func()) {
}

// newRandomValue returns the next random number in the range [0,1] from the supplied randomNumberGenerator.
// (which by default returns a random number in the range [0,1).
// See: http://mumble.net/~campbell/2014/04/28/uniform-random-float
func newRandomValue(randomNumberGenerator *rand.Rand) float64 {
	distributionRange := int64(math.Pow(2, 53))
	return float64(randomNumberGenerator.Int63n(distributionRange)) / float64(distributionRange-1)
}

func (explorer *BaseExplorer) ChangeIsDesirable() bool {
	if explorer.changeInObjectiveValue <= 0 {
		return true
	}
	return false
}

func (explorer *BaseExplorer) AcceptLastChange() {
	explorer.changeAccepted = true
}

func (explorer *BaseExplorer) RevertLastChange() {
	explorer.changeAccepted = false
}

func (explorer *BaseExplorer) ChangeAccepted() bool {
	return explorer.changeAccepted
}

func (explorer *BaseExplorer) Clone() Explorer {
	clone := *explorer
	return &clone
}
