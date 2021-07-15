package config

import "fmt"

const ShortApplicationName = "CREMEngine"
const LongApplicationName = "Catchment Resilience Exploration Modelling Engine "

const Version = "0.5"

func NameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", ShortApplicationName, Version)
}
