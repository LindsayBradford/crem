package scenario

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type Scenario interface {
	LogHandler() logging.Logger
	SetAnnealer(annealer annealing.Annealer)
	Run() error
}

func NewBaseScenario() *BaseScenario {
	scenario := new(BaseScenario)
	return scenario
}

var _ Scenario = new(BaseScenario)

type BaseScenario struct {
	annealer annealing.Annealer
	runner   CallableRunner
	observer observer.Observer
}

func (s *BaseScenario) SetAnnealer(annealer annealing.Annealer) {
	assert.That(s.observer != nil)

	s.runner.SetAnnealer(annealer)
	annealer.AddObserverAsFirst(s.observer)

	s.annealer = annealer
}

func (s *BaseScenario) WithRunner(runner CallableRunner) *BaseScenario {
	s.runner = runner
	return s
}

func (s *BaseScenario) WithObserver(observer observer.Observer) *BaseScenario {
	assert.That(s.runner != nil)
	s.observer = observer
	return s
}

func (s *BaseScenario) LogHandler() logging.Logger {
	assert.That(s.runner != nil)
	return s.runner.LogHandler()
}

func (s *BaseScenario) Run() error {
	assert.That(s.annealer != nil)
	return s.runner.Run()
}

var NullScenario Scenario = new(nullScenario)

type nullScenario struct{}

func (s *nullScenario) SetAnnealer(annealer annealing.Annealer) {}
func (s *nullScenario) LogHandler() logging.Logger              { return nil }
func (s *nullScenario) Run() error                              { return nil }
