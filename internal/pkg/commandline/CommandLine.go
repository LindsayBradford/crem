// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package commandline

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/pkg/errors"
)

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
	Version          bool
	ScenarioFile     string
	ServerConfigFile string
}

// THe define sets up the relevant command-line
// arguments that the utility will accept via the 'flags' package.

func (args *Arguments) define() {

	flag.StringVar(
		&args.ScenarioFile,
		"ScenarioFile",
		"",
		"file dictating scenario run-time behaviour",
	)

	flag.StringVar(
		&args.ServerConfigFile,
		"ServerConfigFile",
		"",
		"file dictating HTTP server runtime behaviour",
	)

	flag.BoolVar(
		&args.Version,
		"Version",
		false,
		"Prints the version number of this utility and exits.",
	)

	flag.Usage = usageMessage

	flag.Parse()
}

// The process method does some simple "utility-stopping" processing
// once the command-line arguments have been parsed into args.
// It catches invalid show-stopping settings, and basic usage message display.

func (args *Arguments) process() {

	if args.Version == true {
		fmt.Println(
			GetVersionString(),
		)
		Exit(0)
	}

	if args.ScenarioFile != "" {
		validateFilePath(args.ScenarioFile)
	}

	if args.ServerConfigFile != "" {
		validateFilePath(args.ServerConfigFile)
	}
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

func Exit(exitValue interface{}) {
	var exitCode int
	switch exitValue.(type) {
	case error:
		exitingError, _ := exitValue.(error)
		fmt.Fprintf(os.Stderr, "Critical Error: %v. Exiting.\n", exitingError)
		exitCode = 1
	case int:
		exitValueAsInt, _ := exitValue.(int)
		exitCode = exitValueAsInt
	default:
		exitCode = 0
	}
	os.Exit(exitCode)
}

// usageMessage is the function we supply to the flags package to upon
// a request for how to use the utility from the command-line

func usageMessage() {
	fmt.Printf("Help for %s\n", GetVersionString())
	fmt.Println("  --Help                         Prints this help message.")
	fmt.Println("  --Version                      Prints the version number of this utility.")
	fmt.Println("  --ScenarioFile  <FilePath>     File describing a scenario to run and its  run-time behaviour.")
	fmt.Println("  --ServerConfigFile <FilePath>  File describing how the application is to run as a web server.")
	fmt.Println()
	fmt.Println("Web-server usage takes the form:")
	fmt.Printf("  %s [--ServerConfigFile <FilePath>]\n", justExecutableName())
	fmt.Println()
	fmt.Println("  If no server config fle is specified, configuration will first attempt to load from the relative")
	fmt.Println("  path \"./config/server.toml\". If no such path exists, a web-server will start with default values")
	fmt.Println("  for all entries in a server config file. ")
	fmt.Println()
	fmt.Println("Running a single scenario takes the form:")
	fmt.Printf("  %s --ScenarioFile <FilePath>\n", justExecutableName())
	fmt.Println()
	fmt.Println("If a scenario file is specified, the application runs the scenario instead of a web server.")

	Exit(0)
}

// Returns a formatted string, identifying the utility, and it's
// version number as defined in the utility's configuration.

func GetVersionString() string {
	return fmt.Sprintf("%s %s (%s)", justExecutableName(), config.Version, runtime.Version())
}

func justExecutableName() string {
	appName := filepath.Base(os.Args[0])
	return appName
}
