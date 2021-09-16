package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/LindsayBradford/crem/cmd/cremconfigurator/theme"
)

func main() {
	app := createApp()

	window := createWindow(app)
	window.ShowAndRun()
}

func createApp() fyne.App {
	a := app.New()
	a.Settings().SetTheme(&theme.Crem{})
	return a
}

func createWindow(app fyne.App) fyne.Window {
	w := app.NewWindow("CREM Configurator")
	messageLabel := widget.NewLabel("Generating CREM Explorer configuration...")

	scenarioItem := widget.NewAccordionItem("[Scencario]", buildScenarioContainer())
	annealerItem := widget.NewAccordionItem("[Annealer]", buildAnnealerContainer())
	modelItem := widget.NewAccordionItem("[Model]", buildModelContainer())

	scopeAccordian := widget.NewAccordion(scenarioItem, annealerItem, modelItem)
	scopeAccordian.MultiOpen = true

	generate := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Generate Configuration",
			func() {
				messageLabel.SetText("Configuration Generated")
			}),
		layout.NewSpacer(),
	)

	w.SetContent(
		container.NewVBox(
			scopeAccordian,
			layout.NewSpacer(),
			generate,
			layout.NewSpacer(),
			messageLabel,
		))

	return w
}

func buildScenarioContainer() *fyne.Container {
	nameLabel := widget.NewLabel("     [Scenario.Name]")
	nameEntry := widget.NewEntry()

	name := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry,
	)

	runNumberLabel := widget.NewLabel("     [Scenario.RunNumber]")
	runNumberEntry := widget.NewEntry()

	runNumber := container.New(layout.NewFormLayout(),
		runNumberLabel, runNumberEntry,
	)

	maxConcurrentRunNumberLabel := widget.NewLabel("     [Scenario.MaximumConcurrentRunNumber]")
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

	outputPathLabel := widget.NewLabel("     [Scenario.OutputPath]")
	outputPathEntry := widget.NewEntry()

	outputPath := container.New(layout.NewFormLayout(),
		outputPathLabel, outputPathEntry,
	)

	outputTypeLabel := widget.NewLabel("     [Scenario.OutputType]")
	var outputTypes = []string{"CSV", "Excel", "JSON"}
	outputTypeSelect := widget.NewSelect(outputTypes, func(selected string) {})

	outputType := container.New(layout.NewFormLayout(),
		outputTypeLabel, outputTypeSelect,
	)

	outputLevelLabel := widget.NewLabel("     [Scenario.OutputLevel]")
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
	)

	return l
}

func buildAnnealerContainer() *fyne.Container {
	typeLabel := widget.NewLabel("     [Annealer.Type] ")
	var annealerTypes = []string{"Simple", "Suppapitnarm", "AveragedSuppapitnarm"}
	typeEntry := widget.NewSelect(annealerTypes, func(selected string) {})

	f := container.New(layout.NewFormLayout(), typeLabel, typeEntry)

	return f
}

func buildModelContainer() *fyne.Container {
	typeLabel := widget.NewLabel("     [Model.Type]     ")
	var modelTypes = []string{"CatchmentModel"}
	typeEntry := widget.NewSelect(modelTypes, func(selected string) {})

	f := container.New(layout.NewFormLayout(), typeLabel, typeEntry)

	return f
}
