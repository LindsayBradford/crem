# Catchment Resilience Exploration Modeller

[![Build Status](https://travis-ci.com/LindsayBradford/crem.svg?token=Xt8jEnqxCbgTcvvxNK8e&branch=master)](https://travis-ci.com/LindsayBradford/crem)
[![Go Report Card](https://goreportcard.com/badge/github.com/LindsayBradford/crem)](https://goreportcard.com/report/github.com/LindsayBradford/crem)
[![GoDoc](https://godoc.org/github.com/LindsayBradford/crem?status.svg)](https://godoc.org/github.com/LindsayBradford/crem)

#### Getting Started

CREM makes use of a number of 3rd-party libraries that are not included in this source repository. The [govendor](https://github.com/kardianos/govendor) utility governs library dependencies via  [vendor.json](https://github.com/LindsayBradford/crem/blob/master/vendor/vendor.json). 
Once you've git-cloned this repository, run:

```
> go get github.com/kardianos/govendor
> cd <new CREM repository folder>
> govendor sync
```

to download compatible versions of the libraries CREM depends on as [vendor libraries](https://golang.org/cmd/go/#hdr-Vendor_Directories). 

From there a `go build` from  within [cmd/cremexplorer](https://github.com/LindsayBradford/crem/blob/master/cmd/cremexplorer) should produce 
a `cremexplorer.exe` executable.  Then run your new executable  from the command-line, specifying a scenario config file like this:

```> cremexplorer.exe --ScenarioFile <someScenarioFile>```

You'll find a [simple test scenario configuration](https://github.com/LindsayBradford/crem/blob/master/cmd/cremexplorer/testdata/TestCREMExplorer-Kirkpatrick-WhiteBox.toml) here. Further detail on configuring a scenario can be found in the [wiki](https://github.com/LindsayBradford/crem/wiki/Configuration#scenario-configuration). 

#### Overview:

The heart of the Catchment Resilience Exploration Modeller is a highly configurable [simulated annealer]( https://en.wikipedia.org/wiki/Simulated_annealing) that allows for the exploration of stakeholder objectives around river catchment resilience. The following objectives have been implemented:
 - Sediment Produced
 - Implementation Cost

The explorer is configured via [TOML](https://github.com/toml-lang/toml) files, based on a '[convention over configuration](https://en.wikipedia.org/wiki/Convention_over_configuration)' approach.

#### General Usage Notes:

- This software was constructed and tested on a 64-bit Windows 10 platform using [GoLang](https://golang.org/) 1.12.7.

- Continuous integration via travis-ci is also employed.

#### Contact Information:

- This software is produced on behalf of [Griffith University](http://www.griffith.edu.au/) within the [Australian Rivers Institute](http://www.griffith.edu.au/environment-planning-architecture/australian-rivers-institute), and originally authored by [Dr Lindsay Bradford](https://github.com/LindsayBradford).

- E-Mail: [ari@griffith.edu.au](mailto:ari@griffith.edu.au), or [l.bradford@griffith.edu.au](mailto:l.bradford@griffith.edu.au)
- Voice: +61 7 3735 7402, or +61 7 3735 6598

#### Copyright:

The Catchment Resilience Exploration Modeller (CREM) software is licensed under a BSD 3-clause "New" or "Revised" licence,
detailed in [LICENCE.md](LICENCE.md).

#### Dependencies:

CREM makes use of the following libraries:

- [Gomega](https://github.com/onsi/gomega)  for a Fluent-API based approach to test assertions
- [go-ole](https://github.com/go-ole/go-ole) for I/O via Excel files
- [BurntSushi/toml](https://github.com/BurntSushi/toml) for TOML config file support
- [pkg/errors](https://github.com/pkg/errors) For error wrapping
- [nu7hatch/gouuid](https://github.com//nu7hatch/gouuid) for UUID generation
