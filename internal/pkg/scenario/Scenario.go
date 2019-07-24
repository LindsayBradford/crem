package scenario

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

type Scenario interface {
	SetAnnealer(annealer annealing.Annealer)
	SetRunner(runner CallableRunner)
	SetObserver(observer observer.Observer)

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

func (s *BaseScenario) SetRunner(runner CallableRunner) {
	s.runner = runner
}

func (s *BaseScenario) SetObserver(observer observer.Observer) {
	assert.That(s.runner != nil)
	s.observer = observer
}

func (s *BaseScenario) Run() {
	assert.That(s.annealer != nil)
	s.runner.Run()
}

var NullScenario Scenario = new(nullScenario)

type nullScenario struct{}

func (s *nullScenario) SetAnnealer(annealer annealing.Annealer) {}
func (s *nullScenario) SetRunner(runner CallableRunner)         {}
func (s *nullScenario) SetObserver(observer observer.Observer)  {}
func (s *nullScenario) Run()                                    {}
