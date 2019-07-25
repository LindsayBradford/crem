package scenario

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

type Scenario interface {
	SetAnnealer(annealer annealing.Annealer)
	Run()
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

	annealer.AddObserver(s.observer)
	s.runner.SetAnnealer(annealer)

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

func (s *BaseScenario) Run() {
	assert.That(s.annealer != nil)
	s.runner.Run()
}

var NullScenario Scenario = new(nullScenario)

type nullScenario struct{}

func (s *nullScenario) SetAnnealer(annealer annealing.Annealer) {}
func (s *nullScenario) Run()                                    {}
