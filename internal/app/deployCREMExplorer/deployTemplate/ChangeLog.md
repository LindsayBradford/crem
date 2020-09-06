# Change Log

## Version 0.10 (04 Sept 2020):
### New Features
* Altered approach to Hillslope calculation in Catchment Model based on pre-processed hillslope hot-spots. 
* Introduction of ParticulateNitrogen variable for tracking particulate nitrogen sources (as per sediment)
* Replacement of Implementation Cost as a parameter function to being defined per management action.
* Introduction of Opportunity Cost, defined per management action.
* Now deploys with Laidley_data_v1_7.xlsx (supporting the above changes) 
### Bug Fixes
* Instead of failing silently, we now report a log error and quit if an invalid objective variable has been specified.  

## Version 0.9 (07 April 2020):
### New Features
* Averaged Suppapitnarm Annealer is now a usable option.
* Deploy environment now has example configuration for both Annealer types.
* Logging now captures and reports which management actions are randomly activated at model initialisation.

## Version 0.8 (10 March 2020):
### Bug Fixes
* Fixed bug in how "SedimentProduction" reports changes to the annealer.
* Changed Kirkpatrick explorer to consider "no objective change" updates to be undesirable.
* Fixed bug in Suppapitnarm Explorer that allowed duplicate solutions to be archived into the solution _set_.
### New Features
* Reworking Kirkpatrick explorer and catchment CoreModel events for better logging.

## Version 0.7 (21 January 2020):
### New Features
* Reintroduced (only) upper bound limit parameters 
  MaximumImplementationCost and MaximumSedimentProduction to limit single 
  objective annealing runs on SedimentProduction and ImplementationCost 
  decision variables respectively.
* Removed decision variable "SedimentVsProduction" and parameters 'SedimentProductionDecisionWeight' & 'ImplementationCostDecisionWeight'. 

## Version 0.6 (15 November 2019):
### New Features
* Now deploys with Laidley_data_v1_5.xlxs (with fixes to hillslope RSLK calculations)
### Fixed
* Removing Upper and Lower bound limits (breaks model in a very hard-to-isolate way)

## Version 0.5 (16 September 2019):
### New Features
* Adding Upper and Lower bound limits to CatchmentModel decision variables SedimentProduction & ImplementationCost.
* Now deploys with Laidley_data_v1_4.xlxs (with fixes to gully volumes)

## Version 0.4 (06 September 2019):
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