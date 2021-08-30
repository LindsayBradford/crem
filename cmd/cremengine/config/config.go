package config

import "fmt"

const ShortApplicationName = "CREMEngine"
const LongApplicationName = "Catchment Resilience Exploration Modelling Engine "

const Version = "0.6"

func NameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", ShortApplicationName, Version)
}
