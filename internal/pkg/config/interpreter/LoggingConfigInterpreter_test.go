// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/config/data"
	"github.com/LindsayBradford/crem/pkg/logging/formatters"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

// const equalTo = "=="

func TestConfigInterpreter_NewLoggingConfigInterpreter_DefaultLoggerNoErrors(t *testing.T) {
	// given
	g := NewGomegaWithT(t)
	expectedLogHandlerName := "DefaultLogHandler"

	// when
	interpreterUnderTest := NewLoggingConfigInterpreter()

	// then
	actualLogHandler := interpreterUnderTest.LogHandler()
	g.Expect(actualLogHandler.Name()).To(Equal(expectedLogHandlerName))
	g.Expect(actualLogHandler.SupportsLogLevel("Annealer")).To(BeTrue())
	g.Expect(actualLogHandler.SupportsLogLevel("Model")).To(BeTrue())

	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_ValidNativeDefaultingLoggingConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.LoggingConfig{
		LoggingType:      data.NativeLibrary,
		LoggingFormatter: data.Json,
	}

	// when
	interpreterUnderTest := NewLoggingConfigInterpreter().Interpret(&configUnderTest)

	// then
	actualLogger := interpreterUnderTest.LogHandler()
	actualFormatter := actualLogger.Formatter()

	g.Expect(actualLogger).To(BeAssignableToTypeOf(&loggers.NativeLibraryLogger{}))
	g.Expect(actualLogger.SupportsLogLevel("Annealer")).To(BeTrue())
	g.Expect(actualLogger.SupportsLogLevel("Model")).To(BeTrue())

	g.Expect(actualFormatter).To(BeAssignableToTypeOf(&formatters.JsonFormatter{}))

	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_ValidBareBonesDefaultingLoggingConfig_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.LoggingConfig{
		LoggingType:      data.BareBones,
		LoggingFormatter: data.NameValuePair,
	}

	// when
	interpreterUnderTest := NewLoggingConfigInterpreter().Interpret(&configUnderTest)

	// then
	actualLogger := interpreterUnderTest.LogHandler()
	actualFormatter := actualLogger.Formatter()

	g.Expect(actualLogger).To(BeAssignableToTypeOf(&loggers.BareBonesLogger{}))
	g.Expect(actualLogger.SupportsLogLevel("Annealer")).To(BeTrue())
	g.Expect(actualLogger.SupportsLogLevel("Model")).To(BeTrue())

	g.Expect(actualFormatter).To(BeAssignableToTypeOf(&formatters.NameValuePairFormatter{}))

	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_InvalidLoggingTypeConfig_Panics(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.LoggingConfig{
		LoggingType:      data.LoggerType{Value: "NotRecognised"},
		LoggingFormatter: data.NameValuePair,
	}

	// when
	interpreterUnderTest := NewLoggingConfigInterpreter()

	interpretFunction := func() {
		interpreterUnderTest.Interpret(&configUnderTest)
	}

	// then
	g.Expect(interpretFunction).To(Panic())
}

func TestConfigInterpreter_InvalidLoggingFormatterConfig_Panics(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.LoggingConfig{
		LoggingType:      data.BareBones,
		LoggingFormatter: data.FormatterType{"Unrecognised"},
	}

	interpreterUnderTest := NewLoggingConfigInterpreter()

	// when
	interpretFunction := func() {
		interpreterUnderTest.Interpret(&configUnderTest)
	}

	// then
	g.Expect(interpretFunction).To(Panic())
}

func TestConfigInterpreter_ValidLoggingConfigWithLogLevelDestinations_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.LoggingConfig{
		LoggingType:      data.BareBones,
		LoggingFormatter: data.RawMessage,
		LogLevelDestinations: map[string]string{
			"Debugging":   "StandardError",
			"Information": "StandardOutput",
			"Warnings":    "Discarded",
			"Errors":      "Discarded",
		},
	}

	// when
	interpreterUnderTest := NewLoggingConfigInterpreter().Interpret(&configUnderTest)

	// then
	actualLogger := interpreterUnderTest.LogHandler()
	actualFormatter := actualLogger.Formatter()

	g.Expect(actualLogger).To(BeAssignableToTypeOf(&loggers.BareBonesLogger{}))
	g.Expect(actualFormatter).To(BeAssignableToTypeOf(&formatters.RawMessageFormatter{}))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_InvalidLogLevelDestinations_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.LoggingConfig{
		LoggingType:      data.BareBones,
		LoggingFormatter: data.RawMessage,
		LogLevelDestinations: map[string]string{
			"Debugging": "NotAValidDestination",
		},
	}

	// when
	interpreterUnderTest := NewLoggingConfigInterpreter().Interpret(&configUnderTest)

	// then
	actualLogger := interpreterUnderTest.LogHandler()
	actualFormatter := actualLogger.Formatter()

	g.Expect(actualLogger).To(BeAssignableToTypeOf(&loggers.BareBonesLogger{}))
	g.Expect(actualFormatter).To(BeAssignableToTypeOf(&formatters.RawMessageFormatter{}))

	g.Expect(interpreterUnderTest.Errors()).To(Not(BeNil()))
	t.Log(interpreterUnderTest.Errors())
}
