# Catchment Resilience Exploration Modeller (CREM)

[![Build Status](https://travis-ci.com/LindsayBradford/crem.svg?token=Xt8jEnqxCbgTcvvxNK8e&branch=master)](https://travis-ci.com/LindsayBradford/crem)
[![Go Report Card](https://goreportcard.com/badge/github.com/LindsayBradford/crem)](https://goreportcard.com/report/github.com/LindsayBradford/crem)
[![GoDoc](https://godoc.org/github.com/LindsayBradford/crem?status.svg)](https://godoc.org/github.com/LindsayBradford/crem)

## Overview:

This repository produces two key applications:

- The [CREMExplorer](https://github.com/LindsayBradford/crem/blob/master/cmd/cremexplorer)
- The [CREMEngine](https://github.com/LindsayBradford/crem/blob/master/cmd/cremengine)

Both applications are configured via [TOML](https://github.com/toml-lang/toml) files, based on
a '[convention over configuration](https://en.wikipedia.org/wiki/Convention_over_configuration)' approach.

### CREMExplorer:

The [CREMExplorer](https://github.com/LindsayBradford/crem/blob/master/cmd/cremexplorer) is a highly configurable
modelling platform that wraps a river catchment model in either a single or
multi-objective [simulated annealer]( https://en.wikipedia.org/wiki/Simulated_annealing) to explore stakeholder
objectives around river catchment resilience.

The [river catchment model](https://github.com/LindsayBradford/crem/blob/master/internal/pkg/model/models/catchment)
tracks the following stakeholder objectives:

- Sediment Produced
- Particulate Nitrogen Produced
- Dissolved Nitrogen Produced
- Management Action Implementation Cost
- Management Action Opportunity Cost

This river catchment model allows management actions to be applied in order to mitigate pollutants entering a river
system. The following management actions have been implemented:

- Riparian revegetation
    - predominately targeting sediment and particulate nitrogen
- Gully Restoration
    - predominately targeting sediment
- Hillslope revegetation
    - predominately targeting sediment and particulate nitrogen
- Wetland establishment
    - predominately targeting dissolved nitrogen

[Single-objective simulated annealing](https://github.com/LindsayBradford/crem/blob/master/internal/pkg/annealing/explorer/kirkpatrick)
is used to find optimised solutions to minimising/maximising a particular stakeholder objective, optionally limited by a
2nd objective. For instance, a scenario can be configured to answer a question like "Find a near-optimal minimised
sediment producted for a budget of $10M in implementation costs."

[Multi-objective simulated annealing](https://github.com/LindsayBradford/crem/blob/master/internal/pkg/annealing/explorer/suppapitnarm)
is used to explore trade-offs between all the supplied stakeholder objectives, optionally limited by one of those
objectives. For instance, a scenario can be configured to answer a question like
"Find me a set of tradeoffs between Sediment and Dissolved Nitrogen produced for a budget of $10M in implementation
costs"

### CREMEngine:

The [CREMEngine](https://github.com/LindsayBradford/crem/blob/master/cmd/cremengine) wraps the river catchment model
described above in a web-server interface, allowing the following:

- The river catchment model can be deployed independent of the annealing, and manipulated it in real-time for
  visualiation purposes.
- Ability to configure the model to run exactly as per scenario definitions supplied to the CREMExplorer.
- Ability to set the model's state to match solutions produced by the CREMExplorer, to showcase optimised solutions, or
  individual tradeoff solutions of interest.

## Getting Started:

CREM makes use of a number of 3rd-party libraries that are not included in this source repository. Go's
built-in [3rd-party module support](https://golang.org/ref/mod) is used to track and integrate needed 3rd-party
libraries. Once you've git-cloned this repository, run:

```
> cd <new CREM repository folder>
> go mod vendor
```

to download compatible versions of the libraries CREM depends on
as [vendor libraries](https://golang.org/cmd/go/#hdr-Vendor_Directories).

From there a `go build` from
within [cmd/cremexplorer](https://github.com/LindsayBradford/crem/blob/master/cmd/cremexplorer) should produce
a `cremexplorer.exe` executable. Then run your new executable from the command-line, specifying a scenario config file
like this:

```> cremexplorer.exe --ScenarioFile <someScenarioFile>```

You'll find
a [simple test scenario configuration](https://github.com/LindsayBradford/crem/blob/master/cmd/cremexplorer/testdata/TestCREMExplorer-Kirkpatrick-WhiteBox.toml)
here. Further detail on configuring a scenario can be found in
the [wiki](https://github.com/LindsayBradford/crem/wiki/Configuration#scenario-configuration).

## General Usage Notes:

- This software was constructed and tested on a 64-bit Windows 10 platform using [GoLang](https://golang.org/) 1.16.

- [Continuous integration](https://travis-ci.com/LindsayBradford/crem) via travis-ci is also employed.

## Contact Information:

- This software is produced on behalf of [Griffith University](http://www.griffith.edu.au/) within the [Australian Rivers Institute](http://www.griffith.edu.au/environment-planning-architecture/australian-rivers-institute), and originally authored by [Dr Lindsay Bradford](https://github.com/LindsayBradford).

- E-Mail: [ari@griffith.edu.au](mailto:ari@griffith.edu.au), or [l.bradford@griffith.edu.au](mailto:l.bradford@griffith.edu.au)
- Voice: +61 7 3735 7402, or +61 7 3735 6598

## Copyright:

The Catchment Resilience Exploration Modeller (CREM) software is licensed under a BSD 3-clause "New" or "Revised" licence,
detailed in [LICENCE.md](LICENCE.md).

## Dependencies:

The following 3rd-party libraries are required for this code-base:

- [Gomega](https://github.com/onsi/gomega)  for a Fluent-API based approach to test assertions
- [go-ole](https://github.com/go-ole/go-ole) for I/O via Excel files
- [BurntSushi/toml](https://github.com/BurntSushi/toml) for TOML config file support
- [pkg/errors](https://github.com/pkg/errors) For error wrapping
- [nu7hatch/gouuid](https://github.com//nu7hatch/gouuid) for UUID generation
