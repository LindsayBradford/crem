package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/LindsayBradford/crem/cmd/cremconfigurator/mvp"
)

var view *FyneView

func New() mvp.View {
	view = new(FyneView).build()
	return view
}

type FyneView struct {
	app        fyne.App
	window     fyne.Window
	messageBar *widget.Label

	observers []mvp.ViewObserver
}

func (v *FyneView) build() *FyneView {
	v.createThemedApp()
	v.createWindow()
	return v
}

func (v *FyneView) Show() {
	v.window.CenterOnScreen()
	v.window.ShowAndRun()
}

func (v *FyneView) SetMessage(message string) {
	v.messageBar.SetText(message)
}

func (v *FyneView) Id() string {
	return "fyne"
}

func (v *FyneView) AddObserver(o mvp.ViewObserver) {
	if v.observers == nil {
		v.observers = *new([]mvp.ViewObserver)
	}
	v.observers = append(v.observers, o)
}

func (v *FyneView) raiseEvent(e mvp.ViewEvent) {
	e.View = v
	for _, o := range v.observers {
		o.EventRaised(e)
	}
}

func (v *FyneView) createThemedApp() {
	v.app = app.New()
	v.app.Settings().SetTheme(&Crem{})
}
