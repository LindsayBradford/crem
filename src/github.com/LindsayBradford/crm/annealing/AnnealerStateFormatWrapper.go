/*
 * Copyright (c) 2018 Australian Rivers Institure. Author: Lindsay Bradford
 */

package annealing

import (
	"fmt"

	. "github.com/LindsayBradford/crm/reflect"
)

type AnnealerStateFormatWrapper struct {
	annealerToFormat Annealer

	methodFormats map[string]string
}

const DEFAULT_FLOAT64_FORMAT = "%f"
const DEFAULT_UINT_FORMAT = "%d"

func NewAnnealerStateFormatWrapper(annealer Annealer) *AnnealerStateFormatWrapper {
	wrapper := AnnealerStateFormatWrapper{
		annealerToFormat: annealer,
		methodFormats: map[string]string{
			"Temperature":      DEFAULT_FLOAT64_FORMAT,
			"CoolingFactor":    DEFAULT_FLOAT64_FORMAT,
			"MaxIterations":    DEFAULT_UINT_FORMAT,
			"CurrentIteration": DEFAULT_UINT_FORMAT,
		},
	}

	return &wrapper
}

func (this *AnnealerStateFormatWrapper) Temperature() string {
	thisMethodsName := DeriveMethodName()
	valueNeedingFormatting := CallMethodReturningFloat64(this.annealerToFormat, thisMethodsName)
	return this.applyFormatting(thisMethodsName, valueNeedingFormatting)
}

func (this *AnnealerStateFormatWrapper) CoolingFactor() string {
	thisMethodsName := DeriveMethodName()
	valueNeedingFormatting := CallMethodReturningFloat64(this.annealerToFormat, thisMethodsName)
	return this.applyFormatting(thisMethodsName, valueNeedingFormatting)
}

func (this *AnnealerStateFormatWrapper) MaxIterations() string {
	thisMethodsName := DeriveMethodName()
	valueNeedingFormatting := CallMethodReturningUint(this.annealerToFormat, thisMethodsName)
	return this.applyFormatting(thisMethodsName, valueNeedingFormatting)
}

func (this *AnnealerStateFormatWrapper) CurrentIteration() string {
	thisMethodsName := DeriveMethodName()
	valueNeedingFormatting := CallMethodReturningUint(this.annealerToFormat, thisMethodsName)
	return this.applyFormatting(thisMethodsName, valueNeedingFormatting)
}

func (this *AnnealerStateFormatWrapper) applyFormatting(formatKey string, valueToFormat interface{}) string {
	formatToApply := this.methodFormats[formatKey]
	return fmt.Sprintf(formatToApply, valueToFormat)
}