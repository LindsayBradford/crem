/*
 * Copyright (c) 2018 Australian Rivers Institure. Author: Lindsay Bradford
 */

package reflect

import (
	"reflect"
	"runtime"
	"strings"
)

func DeriveMethodName() string {
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	callingMethod := strings.Split(details.Name(), ".")
	justMethodName := callingMethod[len(callingMethod)-1]
	return justMethodName
}

func CallMethodReturningFloat64(any interface{}, methodName string) float64 {
	return reflect.ValueOf(any).MethodByName(methodName).Call([]reflect.Value{})[0].Float()
}

func CallMethodReturningUint(any interface{}, methodName string) uint {
	return uint(reflect.ValueOf(any).MethodByName(methodName).Call([]reflect.Value{})[0].Uint())
}
