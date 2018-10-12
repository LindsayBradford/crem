// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"fmt"
	"time"

	"github.com/LindsayBradford/crm/config"
)

const ContentTypeHeaderKey = "Content-Type"

const TomlMimeType = "application/toml"
const JsonMimeType = "application/json"

func FormattedTimestamp() string {
	return fmt.Sprintf("%v", time.Now().Format(time.RFC3339Nano))
}

func NameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", config.LongApplicationName, config.Version)
}
