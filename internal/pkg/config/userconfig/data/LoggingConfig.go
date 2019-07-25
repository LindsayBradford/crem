// Copyright (c) 2019 Australian Rivers Institute.

package data

type LoggingConfig struct {
	LoggingType          LoggerType
	LoggingFormatter     FormatterType
	LogLevelDestinations map[string]string
}

type LoggerType struct {
	value string
}

var (
	UnspecifiedLoggerType = LoggerType{""}
	NativeLibrary         = LoggerType{"NativeLibrary"}
	BareBones             = LoggerType{"BareBones"}
)

func (lt *LoggerType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "[[Loggers]].Type",
		validValues: []string{
			NativeLibrary.value, BareBones.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			lt.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}

type FormatterType struct {
	value string
}

var (
	UnspecifiedFormatterType = FormatterType{""}
	RawMessage               = FormatterType{"RawMessage"}
	Json                     = FormatterType{"JSON"}
	NameValuePair            = FormatterType{"NameValuePair"}
)

func (ft *FormatterType) UnmarshalText(text []byte) error {
	context := unmarshalContext{
		configKey: "[[Loggers]].Formatter",
		validValues: []string{
			RawMessage.value, Json.value, NameValuePair.value,
		},
		textToValidate: string(text),
		assignmentFunction: func() {
			ft.value = string(text)
		},
	}

	return processUnmarshalContext(context)
}
