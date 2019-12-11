// Copyright (c) 2019 Australian Rivers Institute.

package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/pkg/errors"
)

// RoundFloat implements a simple (and thus, probably flawed) floating-point rounding function, rounding the
// supplied value the 'precision' number of decimal places.  If the value supplied is too large to accurately
// convert, the function will panic.
func RoundFloat(value float64, precision int) float64 {
	shift := math.Pow10(precision)

	if math.Abs(value) > math.MaxFloat64/shift {
		panic(errors.New("Attempt to round floating point number too big for precision required."))
	}

	shiftedValue := math.Round(value*shift) / shift
	return shiftedValue
}

// https://stackoverflow.com/a/55769252/772278
func DerivePrecision(value float64) int {
	decimalPlaces := fmt.Sprintf("%f", value-math.Floor(value))  // produces 0.xxxx0000
	decimalPlaces = strings.Replace(decimalPlaces, "0.", "", -1) // remove 0.
	decimalPlaces = strings.TrimRight(decimalPlaces, "0")        // remove trailing 0s
	return len(decimalPlaces)
}
