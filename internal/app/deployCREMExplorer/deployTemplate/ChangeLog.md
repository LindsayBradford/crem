# Change Log

## Version 0.4 (_TBD_):
### New Features
### Fixed
* Fixed dangling Excel Ole handler on panic over saving to open files, causing Excel to generally misbehave. 
* Fixed bug in Hillslope calculation of Vegetation cover 

## Version 0.3 (29 July 2019):
### New Features
* Simplified approach to application configuration.
* Introduced parameters 'SedimentProductionDecisionWeight' & 'ImplementationCostDecisionWeight' to influence 
  Kirkpatrick Annealer decisions on minimising SedimentVsProduction decision variable. 

## Version 0.2 (23 July 2019):
### New Features
* Introduction of hillslope sediment and hillslope revegetation management action to CatchmentModel.
* Introduced this change log. 
### Fixed
* Model parameter 'GullySedimentReductionTarget' was being mis-read.  Instead of a reduction of sediment 
  produced from active gully repair action _by_ 0.8  of original gully sediment, the model was reducing sediment 
  produced _to_ 0.8 of original sediment (assuming parameter default). 

## Version 0.1 (18 July 2019):
### New Features
* Kirkpatrick Annealer, weighting the minimisation of SedimentProduction against ImplementationCost variables evenly.
* CatchmentModel, implementing sediment from river banks and gullies, along with matching management actions.