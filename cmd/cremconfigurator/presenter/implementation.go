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
	switch ve.Type {
	case mvp.FocusGained:
		message := fmt.Sprintf("Event [%s] raised on field [%v] from view [%s]", ve.Type.String(), ve.Context, ve.View.Id())
		fmt.Println(message)
		//ve.View.SetMessage(message)
	default:
		message := fmt.Sprintf("Event [%s] raised from view [%s]", ve.Type.String(), ve.View.Id())
		fmt.Println(message)
		//ve.View.SetMessage(message)
	}
}
