// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package commandline

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/LindsayBradford/crm/config"
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
	Version    bool
	CpuProfile string
	ConfigFile string
}

// THe define sets up the relevant command-line
// arguments that the utility will accept via the 'flags' package.

func (args *Arguments) define() {

	flag.StringVar(
		&args.CpuProfile,
		"CpuProfile",
		"",
		"write cpu profile to file",
	)

	flag.StringVar(
		&args.ConfigFile,
		"ConfigFile",
		"",
		"file dictating run-time behaviour",
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

	if args.ConfigFile == "" {
		usageMessage()
		os.Exit(0)
	}

	if args.Version == true {
		fmt.Println(
			GetVersionString(),
		)
		os.Exit(0)
	}

	if args.ConfigFile != "" {
		pathInfo, err := os.Stat(args.ConfigFile)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Could not find configuration file [%s]. Exiting.\n", args.ConfigFile)
			os.Exit(1)
		}
		if pathInfo.Mode().IsDir() {
			fmt.Fprintf(os.Stderr, "Error: Configuration file [%s] is a directory. Exiting.\n", args.ConfigFile)
			os.Exit(1)
		}
	}
}

// usageMessage is the function we supply to the flags package to upon
// a request for how to use the utility from the command-line

func usageMessage() {
	fmt.Printf("Help for %s\n", GetVersionString())
	fmt.Println("  --Help                        Prints this help message.")
	fmt.Println("  --Version                     Prints the version number of this utility.")
	fmt.Println("  --ConfigFile  <FilePath>      File that configures the applications run-time behaviour.")
	fmt.Println("  --CpuProfile  <FilePath>      Capture CPU profiling to file.")
	fmt.Println()
	fmt.Println("General usage takes the form:")
	fmt.Printf("  %s --ConfigFile <FilePath>\n", justExecutableName())
	os.Exit(0)
}

// Returns a formatted string, identifying the utility, and it's
// version number as defined in the utility's configuration.

func GetVersionString() string {
	return fmt.Sprintf("%s v%s (%s)", justExecutableName(), config.VERSION, runtime.Version())
}

func justExecutableName() string {
	appName := filepath.Base(os.Args[0])
	return appName
}
