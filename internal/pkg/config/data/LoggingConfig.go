// Copyright (c) 2019 Australian Rivers Institute.

package data

type LoggingConfig struct {
	Type                 LoggerType
	Formatter            FormatterType
	LogLevelDestinations map[string]string
}

type LoggerType struct {
	Value string
}

var (
	UnspecifiedLoggerType = LoggerType{""}
	NativeLibrary         = LoggerType{"NativeLibrary"}
	BareBones             = LoggerType{"BareBones"}
)

func (lt *LoggerType) UnmarshalText(text []byte) error {
	context := UnmarshalContext{
		ConfigKey: "Type",
		ValidValues: []string{
			NativeLibrary.Value, BareBones.Value,
		},
		TextToValidate: string(text),
		AssignmentFunction: func() {
			lt.Value = string(text)
		},
	}

	return ProcessUnmarshalContext(context)
}

type FormatterType struct {
	Value string
}

var (
	UnspecifiedFormatterType = FormatterType{""}
	RawMessage               = FormatterType{"RawMessage"}
	Json                     = FormatterType{"JSON"}
	NameValuePair            = FormatterType{"NameValuePair"}
)

func (ft *FormatterType) UnmarshalText(text []byte) error {
	context := UnmarshalContext{
		ConfigKey: "Formatter",
		ValidValues: []string{
			RawMessage.Value, Json.Value, NameValuePair.Value,
		},
		TextToValidate: string(text),
		AssignmentFunction: func() {
			ft.Value = string(text)
		},
	}

	return ProcessUnmarshalContext(context)
}
