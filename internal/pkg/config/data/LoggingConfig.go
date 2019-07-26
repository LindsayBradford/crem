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
	context := UnmarshalContext{
		ConfigKey: "[[Loggers]].Type",
		ValidValues: []string{
			NativeLibrary.value, BareBones.value,
		},
		TextToValidate: string(text),
		AssignmentFunction: func() {
			lt.value = string(text)
		},
	}

	return ProcessUnmarshalContext(context)
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
	context := UnmarshalContext{
		ConfigKey: "[[Loggers]].Formatter",
		ValidValues: []string{
			RawMessage.value, Json.value, NameValuePair.value,
		},
		TextToValidate: string(text),
		AssignmentFunction: func() {
			ft.value = string(text)
		},
	}

	return ProcessUnmarshalContext(context)
}
