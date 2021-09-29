package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/LindsayBradford/crem/cmd/cremconfigurator/mvp"
)

func New() mvp.View {
	return new(FyneView).build()
}

type FyneView struct {
	app    fyne.App
	window fyne.Window

	observers []mvp.ViewObserver
}

func (v *FyneView) build() *FyneView {
	v.createThemedApp()
	v.createWindow()
	return v
}

func (v *FyneView) Show() {
	v.window.ShowAndRun()
}

func (v *FyneView) AddObserver(o mvp.ViewObserver) {
	if v.observers == nil {
		v.observers = *new([]mvp.ViewObserver)
	}
	v.observers = append(v.observers, o)
}

func (v *FyneView) raiseEvent(e mvp.ViewEvent) {
	for _, o := range v.observers {
		o.EventRaised(e)
	}
}

func (v *FyneView) createThemedApp() {
	v.app = app.New()
	v.app.Settings().SetTheme(&Crem{})
}
