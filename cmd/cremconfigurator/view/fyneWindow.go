package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/LindsayBradford/crem/cmd/cremconfigurator/mvp"
)

func (v *FyneView) createWindow() {
	v.window = v.app.NewWindow("CREM Configurator")

	v.window.SetContent(
		container.NewVBox(
			v.createScopeAccordion(),
			layout.NewSpacer(),
			v.createBottomContainer(),
		))
}

func (v *FyneView) createBottomContainer() *fyne.Container {
	v.messageBar = widget.NewLabel("CREM Configurator: The CREM configuration tool")

	bottomContainer := container.NewVBox(
		v.createGenerateButtonContainer(),
		layout.NewSpacer(),
		v.messageBar,
	)

	return bottomContainer
}

func (v *FyneView) createGenerateButtonContainer() *fyne.Container {
	buttonContainer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Generate Configuration",
			func() {
				buttonEvent := mvp.ViewEvent{Type: mvp.GenerationRequested}
				v.raiseEvent(buttonEvent)
			}),
		layout.NewSpacer(),
	)
	return buttonContainer
}

func (v *FyneView) createScopeAccordion() *widget.Accordion {
	scenarioItem := widget.NewAccordionItem("Scencario", buildScenarioContainer())
	annealerItem := widget.NewAccordionItem("Annealer", buildAnnealerContainer())
	modelItem := widget.NewAccordionItem("Model", buildModelContainer())

	accordion := widget.NewAccordion(scenarioItem, annealerItem, modelItem)
	//scopeAccordion.MultiOpen = true
	return accordion
}

func buildScenarioContainer() *fyne.Container {
	nameLabel := widget.NewLabel("     Name")
	nameLabel.Alignment = fyne.TextAlignTrailing

	nameEntry := NewCremEntry("Scenario.Name")

	name := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry,
	)

	runNumberLabel := widget.NewLabel("     Run Number")
	nameLabel.Alignment = fyne.TextAlignTrailing

	runNumberEntry := NewCremEntry("Scenario.RunNumber")

	runNumber := container.New(layout.NewFormLayout(),
		runNumberLabel, runNumberEntry,
	)

	maxConcurrentRunNumberLabel := widget.NewLabel("     Maximum concurrent run number")
	nameLabel.Alignment = fyne.TextAlignTrailing

	maxConcurrentRunNumberEntry := NewCremEntry("Scenario.MaxConcurrentRunNumber")

	concurrentRunNumbers := container.New(layout.NewFormLayout(),
		maxConcurrentRunNumberLabel, maxConcurrentRunNumberEntry,
	)

	baseNumbers := container.NewHBox(
		layout.NewSpacer(),
		runNumber,
		layout.NewSpacer(),
		concurrentRunNumbers,
		layout.NewSpacer(),
	)

	base := container.NewVBox(
		name,
		baseNumbers,
	)

	outputPathLabel := widget.NewLabel("     Output Path")
	outputPathEntry := NewCremEntry("Scenario.OutputPath")

	outputPath := container.New(layout.NewFormLayout(),
		outputPathLabel, outputPathEntry,
	)

	outputTypeLabel := widget.NewLabel("     Output Type")
	var outputTypes = []string{"CSV", "Excel", "JSON"}
	outputTypeSelect := widget.NewSelect(outputTypes, func(selected string) {})

	outputType := container.New(layout.NewFormLayout(),
		outputTypeLabel, outputTypeSelect,
	)

	outputLevelLabel := widget.NewLabel("     OutputLevel")
	var outputLevels = []string{"Summary", "Detail"}
	outputLevelSelect := widget.NewSelect(outputLevels, func(selected string) {})

	outputLevel := container.New(layout.NewFormLayout(),
		outputLevelLabel, outputLevelSelect,
	)

	outputDetail := container.NewHBox(
		layout.NewSpacer(),
		outputType,
		layout.NewSpacer(),
		outputLevel,
		layout.NewSpacer(),
	)

	output := container.NewVBox(
		outputPath,
		outputDetail,
	)

	l := container.NewVBox(
		base,
		widget.NewSeparator(),
		output,
		widget.NewSeparator(),
	)

	return l
}

func buildAnnealerContainer() *fyne.Container {
	typeLabel := widget.NewLabel("     Type ")
	var annealerTypes = []string{"Simple", "Suppapitnarm", "AveragedSuppapitnarm"}
	typeEntry := widget.NewSelect(annealerTypes, func(selected string) {})

	c := container.NewVBox(
		container.New(layout.NewFormLayout(), typeLabel, typeEntry),
		widget.NewSeparator(),
	)

	return c
}

func buildModelContainer() *fyne.Container {
	typeLabel := widget.NewLabel("     Type     ")
	var modelTypes = []string{"CatchmentModel"}
	typeEntry := widget.NewSelect(modelTypes, func(selected string) {})

	c := container.NewVBox(
		container.New(layout.NewFormLayout(), typeLabel, typeEntry),
		widget.NewSeparator(),
		widget.NewSeparator(),
	)

	return c
}
