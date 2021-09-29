package presenter

import (
	"fmt"
	"github.com/LindsayBradford/crem/cmd/cremconfigurator/mvp"
)

func New() mvp.Presenter {
	return new(implementation).build()
}

type implementation struct{}

func (v *implementation) build() *implementation {
	return v
}

func (v *implementation) EventRaised(ve mvp.ViewEvent) {
	fmt.Printf("Event [%d] raised\n", ve)
}
