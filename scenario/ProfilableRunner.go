// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"github.com/LindsayBradford/crem/excel"
	"github.com/LindsayBradford/crem/logging/handlers"
	"github.com/LindsayBradford/crem/scenario/profiling"
)

type ProfilableRunner struct {
	base        CallableRunner
	profilePath string
}

func (runner *ProfilableRunner) ThatProfiles(base CallableRunner) *ProfilableRunner {
	runner.base = base
	return runner
}

func (runner *ProfilableRunner) ToFile(filePath string) *ProfilableRunner {
	runner.profilePath = filePath
	return runner
}

func (runner *ProfilableRunner) LogHandler() handlers.LogHandler {
	return runner.base.LogHandler()
}

func (runner *ProfilableRunner) Run() error {
	runner.LogHandler().Info("About to collect cpu profiling data to file [" + runner.profilePath + "]")
	defer runner.LogHandler().Info("Collection of cpu profiling data to file [" + runner.profilePath + "] complete.")

	return profiling.CpuProfileOfFunctionToFile(runner.base.Run, runner.profilePath)
}

type SpreadsheetSafeScenarioRunner struct {
	base CallableRunner
}

func (runner *SpreadsheetSafeScenarioRunner) ThatLocks(base CallableRunner) *SpreadsheetSafeScenarioRunner {
	runner.base = base
	return runner
}

func (runner *SpreadsheetSafeScenarioRunner) LogHandler() handlers.LogHandler {
	return runner.base.LogHandler()
}

func (runner *SpreadsheetSafeScenarioRunner) Run() error {
	runner.LogHandler().Debug("Making scenario runner spreadsheet interaction safe")

	if err := excel.EnableSpreadsheetSafeties(); err != nil {
		return err
	}

	defer excel.DisableSpreadsheetSafeties()
	defer runner.LogHandler().Debug("Released scenario runner spreadsheet interaction safeties")

	return runner.base.Run()
}
