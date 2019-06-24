// Copyright (c) 2019 Australian Rivers Institute.

package strings

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/LindsayBradford/crem/pkg/math"
	"golang.org/x/text/message"
)

const (
	englishLocalisation = "en"

	integerFormat = "%d"

	defaultPrecision = 6

	escapedQuote = "\""
)

var (
	localised *message.Printer
)

func init() {
	localised = message.NewPrinter(message.MatchLanguage(englishLocalisation))
}

func NewConverter() *Converter {
	return new(Converter).
		WithFloatingPointPrecision(defaultPrecision).
		NotPaddingZeros()
}

type Converter struct {
	precision   int
	floatFormat string

	quoting bool
}

func (c *Converter) WithFloatingPointPrecision(precision int) *Converter {
	c.precision = precision
	return c
}

func (c *Converter) PaddingZeros() *Converter {
	finalFormat := fmt.Sprintf("%%.%df", c.precision)
	c.floatFormat = finalFormat
	return c
}

func (c *Converter) NotPaddingZeros() *Converter {
	c.floatFormat = "%g"
	return c
}

func (c *Converter) QuotingStrings() *Converter {
	c.quoting = true
	return c
}

func (c *Converter) Convert(value interface{}) string {
	switch value.(type) {
	case string, fmt.Stringer:
		return c.convertString(value)
	default:
		return c.convertNonString(value)
	}
}

func (c *Converter) convertString(value interface{}) string {
	if c.quoting {
		return escapedQuote + c.convertRawString(value) + escapedQuote
	}
	return c.convertRawString(value)
}

func (c *Converter) convertRawString(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case fmt.Stringer:
		valueAsStringer := value.(fmt.Stringer)
		return valueAsStringer.String()
	}
	panic(errors.New("could not convert value to string"))
}

func (c *Converter) convertNonString(value interface{}) string {
	switch value.(type) {
	case bool:
		return strconv.FormatBool(value.(bool))
	case int:
		return localised.Sprintf(integerFormat, value.(int))
	case uint64:
		return localised.Sprintf(integerFormat, value.(uint64))
	case float64:
		roundedValue := math.RoundFloat(value.(float64), c.precision)
		return localised.Sprintf(c.floatFormat, roundedValue)
	}
	panic(errors.New("could not convert value to string"))
}
