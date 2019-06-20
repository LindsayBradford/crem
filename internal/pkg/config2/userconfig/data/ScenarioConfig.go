// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package data

type ScenarioConfig struct {
	Name string

	RunNumber                  uint64
	MaximumConcurrentRunNumber uint64

	OutputPath string
	OutputType ScenarioOutputType

	CpuProfilePath string

	Observer ObserverConfig

	UserDetail map[string]interface{}
}

type ScenarioOutputType struct {
	value string
}

func (sot *ScenarioOutputType) String() string {
	return sot.value
}

var (
	CsvOutput   = ScenarioOutputType{"CSV"}
	JsonOutput  = ScenarioOutputType{"JSON"}
	ExcelOutput = ScenarioOutputType{"EXCEL"}
)

func (sot *ScenarioOutputType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "OutputType",
		validValues: []string{
			CsvOutput.value, JsonOutput.value, ExcelOutput.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			sot.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}
