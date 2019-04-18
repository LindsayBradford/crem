// Copyright (c) 2019 Australian Rivers Institute.

package strings

import (
	"errors"
	"fmt"
	"strconv"
)

func Convert(value interface{}) string {
	switch value.(type) {
	case bool:
		return strconv.FormatBool(value.(bool))
	case int:
		return strconv.Itoa(value.(int))
	case uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case float64:
		return strconv.FormatFloat(value.(float64), 'g', -1, 64)
	case string:
		return value.(string)
	case fmt.Stringer:
		valueAsStringer := value.(fmt.Stringer)
		return valueAsStringer.String()
	}
	panic(errors.New("could not convert value to string"))
}
