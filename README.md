# Catchment Resilience Modeller

[![Build Status](https://travis-ci.com/LindsayBradford/crm.svg?token=Xt8jEnqxCbgTcvvxNK8e&branch=master)](https://travis-ci.com/LindsayBradford/crm)
[![GoDoc](https://godoc.org/github.com/LindsayBradford/crm?status.svg)](https://godoc.org/github.com/LindsayBradford/crm)

## Release 0.1 - 04/06/2018

TBD: <Add blurb>

#### Overview:

The heart of the Catchment Resilience Modeller is a highly configurable [simulated annealer]( https://en.wikipedia.org/wiki/Simulated_annealing). 
A number of simpler annealing applications that exercise the code-base from a testing perspective can be found in the 
[/internal/app](https://github.com/LindsayBradford/crm/tree/master/internal/app) directory. The easiest place to start is with the [dumb annealer](https://github.com/LindsayBradford/crm/blob/master/internal/app/dumbannealer/main.go). 

Below is a very brief introduction to each testing annealer:
- DumbAnnealer: 
   - Attempts to anneal to an "objective value" of 0 a solution space, where each step has a 50/50 chance of  being better or worse by 1 "point".
   - Its primary purpose is to exercise temperature cooling to ensure that over time, the algorithm is less likely to choose a worse option.
   - The solution space is essentially stateless. 
- SimpleExcelAnnealer:
  - Its primary purpose is to exercise a relatively trivial annealing example with a solution-space that is retrieved 
    and stored to [an excel spreadsheet](https://github.com/LindsayBradford/crm/blob/master/internal/app/SimpleExcelAnnealer/testdata/SimpleExcelAnnealerTestFixture.xls). 
  - The algorithm itself started in an excel macro spreadsheet, was converted to VB.NET and has been cut across to 
    golang. This annealer is also used to profile and benchmark the core annealing algorithm across languages. 
    
The Annealers are configure via [TOML](https://github.com/toml-lang/toml) files that can be found in their respective /testdata directories. 
There is currently no documentation describing the structure of config files. However, a test config file 
covering a rich set of features [can be found here](https://github.com/LindsayBradford/crm/blob/master/config/testdata/DumbAnnealerRichValidConfig.toml). 
Configuration operates on a 'convention over configuration' basis, allowing [far smaller configurations](https://github.com/LindsayBradford/crm/blob/master/config/testdata/NullAnnealerMinimalValidConfig.toml).

#### General Usage Notes:

- This software was constructed and tested on a 64-bit Windows 10 platform using [GoLang](https://golang.org/) 1.10.3.

- Continuous integration via travis-ci is also employed.

#### Contact Information:

- This software is produced on behalf of [Griffith University](http://www.griffith.edu.au/) within the [Australian Rivers Institute](http://www.griffith.edu.au/environment-planning-architecture/australian-rivers-institute), and originally authored by [Dr Lindsay Bradford](https://github.com/LindsayBradford).

- E-Mail: [ari@griffith.edu.au](mailto:ari@griffith.edu.au), or [l.bradford@griffith.edu.au](mailto:l.bradford@griffith.edu.au)
- Voice: +61 7 3735 7402, or +61 7 3735 6598

#### Copyright:

The Catchment Resilience Modeller (CRM) software is licensed under a BSD 3-clause "New" or "Revised" licence,
detailed in [LICENCE.md](LICENCE.md).

#### Dependencies:

- [Gomega](https://github.com/onsi/gomega)  For a Fluent-API based approach to test assertions
- [go-ole](https://github.com/go-ole/go-ole) For I/O via Excel files
- [pkg/errors](https://github.com/pkg/errors) For error wrapping
