// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package commandline

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/LindsayBradford/crem/cmd/cremengine/config"
	"github.com/pkg/errors"
)

//go:embed LICENCE.md
var licence string

// ParseArguments processes the command-line arguments supplied
// to the utility and returns a populated Arguments struct containing
// relevant argument values for use later in the utility.

func ParseArguments() *Arguments {
	args := new(Arguments)

	args.define()
	args.process()

	return args
}

type Arguments struct {
	Version             bool
	Licence             bool
	EngineConfigFile    string
	ScenarioFile        string
	SolutionFile        string
	SolutionSummaryFile string
}

// THe define sets up the relevant command-line
// arguments that the utility will accept via the 'flags' package.

func (args *Arguments) define() {

	flag.StringVar(
		&args.EngineConfigFile,
		"EngineConfigFile",
		"",
		"file dictating engine run-time behaviour",
	)

	flag.BoolVar(
		&args.Version,
		"Version",
		false,
		"Prints the version number of this utility and exits.",
	)

	flag.BoolVar(
		&args.Licence,
		"Licence",
		false,
		"Prints the copyright licence under which this software is released.",
	)

	flag.StringVar(
		&args.ScenarioFile,
		"ScenarioFile",
		"",
		"file dictating scenario run-time behaviour",
	)

	flag.StringVar(
		&args.SolutionFile,
		"SolutionFile",
		"",
		"file dictating scenario run-time management action state",
	)

	flag.StringVar(
		&args.SolutionSummaryFile,
		"SolutionSummaryFile",
		"",
		"file detailing a set of solutions to a particular run of scenario",
	)

	flag.Usage = usageMessage

	flag.Parse()
}

// The process method does some simple "utility-stopping" processing
// once the command-line arguments have been parsed into args.
// It catches invalid show-stopping settings, and basic usage message display.

func (args *Arguments) process() {

	if flag.NFlag() == 0 {
		flag.Usage()
	}

	if args.Version == true {
		fmt.Println(
			GetVersionString(),
		)
		Exit(0)
	}

	if args.Licence == true {
		fmt.Println(
			GetLicenceString(),
		)
		Exit(0)
	}

	if args.EngineConfigFile != "" {
		validateFilePath(args.EngineConfigFile)
	}

	if args.ScenarioFile != "" {
		validateFilePath(args.ScenarioFile)
	}

	if args.SolutionFile != "" {
		validateFilePath(args.SolutionFile)
	}

	if args.SolutionSummaryFile != "" {
		validateFilePath(args.SolutionSummaryFile)
	}

	validateArgInterdependence(args)
}

func validateFilePath(filePath string) {
	pathInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		exitError := errors.Errorf("file specified [%s] does not exist", filePath)
		Exit(exitError)
	}
	if pathInfo.Mode().IsDir() {
		exitError := errors.Errorf("file specified [%s] is a directory, not a file", filePath)
		Exit(exitError)
	}
}

func validateArgInterdependence(args *Arguments) {
	if args.SolutionFile != "" && args.ScenarioFile == "" {
		exitError := errors.Errorf("solution file [%s] needs a scenario file (--ScenarioFile flag) specified", args.SolutionFile)
		Exit(exitError)
	}

	if args.SolutionSummaryFile != "" && args.ScenarioFile == "" {
		exitError := errors.Errorf("solution summary file [%s] needs a scenario file (--ScenarioFile flag) specified", args.SolutionSummaryFile)
		Exit(exitError)
	}
}

func Exit(exitValue interface{}) {
	var exitCode int
	switch exitValue.(type) {
	case error:
		exitingError, _ := exitValue.(error)
		fmt.Fprintf(os.Stderr, "Critical Error; forcing application exit: %v\n", exitingError)
		exitCode = 1
	case int:
		exitValueAsInt, _ := exitValue.(int)
		exitCode = exitValueAsInt
	case nil:
		exitCode = 0
	default:
		fmt.Fprintf(os.Stderr, "Critical Error; forcing application exit for unknown error type %v\n", exitValue)
		exitCode = 1
	}
	os.Exit(exitCode)
}

// usageMessage is the function we supply to the flags package to upon
// a request for how to use the utility from the command-line

func usageMessage() {
	fmt.Printf("Usage of %s\n", GetVersionString())
	fmt.Println("  --Help                             Prints this help message.")
	fmt.Println("  --Version                          Prints the version number of this utility.")
	fmt.Println("  --Licence                          Prints the copyright licence of this utility.")
	fmt.Println("  --EngineConfigFile  <FilePath>     File describing the engine run-time behaviour.")
	fmt.Println()
	fmt.Println("  --ScenarioFile  <FilePath>         File describing a scenario to run and its  run-time behaviour.")
	fmt.Println("  --SolutionFile  <FilePath>         File describing a scenario run-time management action state.")
	fmt.Println("  --SolutionSummaryFile  <FilePath>  File describing a set of solutions to a particular run of supplied scenario.")

	Exit(0)
}

// Returns a formatted string, identifying the utility, and it's
// version number as defined in the utility's configuration.

func GetVersionString() string {
	return fmt.Sprintf("%s v%s (%s)", justExecutableName(), config.Version, runtime.Version())
}

func GetLicenceString() string {
	return fmt.Sprintf("%s", licence)
}

func justExecutableName() string {
	appName := filepath.Base(os.Args[0])
	return appName
}
