// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"fmt"
	"time"
)

func FormattedTimestamp() string {
	return fmt.Sprintf("%v", time.Now().Format(time.RFC3339Nano))
}
