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

	messageLabel := widget.NewLabel("Generating CREM Explorer configuration...")
	scenarioItem := widget.NewAccordionItem("Scencario", buildScenarioContainer())

	annealerItem := widget.NewAccordionItem("Annealer", buildAnnealerContainer())
	modelItem := widget.NewAccordionItem("Model", buildModelContainer())

	scopeAccordian := widget.NewAccordion(scenarioItem, annealerItem, modelItem)
	scopeAccordian.MultiOpen = true

	generate := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Generate Configuration",
			func() {
				messageLabel.SetText("Configuration Generated")
				v.raiseEvent(mvp.GenerationRequested)
			}),
		layout.NewSpacer(),
	)

	v.window.SetContent(
		container.NewVBox(
			scopeAccordian,
			layout.NewSpacer(),
			generate,
			layout.NewSpacer(),
			messageLabel,
		))
}

func buildScenarioContainer() *fyne.Container {
	nameLabel := widget.NewLabel("     Name")
	nameLabel.Alignment = fyne.TextAlignTrailing

	nameEntry := widget.NewEntry()

	name := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry,
	)

	runNumberLabel := widget.NewLabel("     Run Number")
	nameLabel.Alignment = fyne.TextAlignTrailing
	runNumberEntry := widget.NewEntry()

	runNumber := container.New(layout.NewFormLayout(),
		runNumberLabel, runNumberEntry,
	)

	maxConcurrentRunNumberLabel := widget.NewLabel("     Maximum concurrent run number")
	nameLabel.Alignment = fyne.TextAlignTrailing

	maxConcurrentRunNumberEntry := widget.NewEntry()

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
	outputPathEntry := widget.NewEntry()

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
