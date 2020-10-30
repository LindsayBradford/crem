# Change Log

## Version 0.13.1 (30/10/2020):
### New Features
* Now deploys with Laidley_data_v1_7_1.xlsx (containing a small data fix)
* Annealing Observer now take precedence over generic logging, ensuring better ordering of logged events.
### Bug Fixes
* Fixed Excel resource handlers released in wrong order, triggering runtime error in SOSA solution saving. 

## Version 0.12 (07/10/2020):
### New Features
* Minor logging changes to allow for easier annealing quality analysis.
* Now deploys with a 'scripts' directory
    * 'MOSA_QualityExtractor.py' that extracts MOSA quality metrics from logs.
    * 'SOSA_QualityExtractor.py' that extracts SOSA quality metrics from logs.
* ALl existing scripts have been cut across to python, with duplication across scripting languages removed.
* Logging now reports application name and version in opening log entry. 
### Bug Fixes
* Fixed issue where management actions for a model were created in differing orders, breaking compression.

## Version 0.11 (28 Sept 2020):
### New Features
* Planning Units renamed to Subcatchments for CatchmentModel output.
* Example config files now ship with rich configuration detail, showcasing full config potential.
* A CSV-formatted summary file of variable values is now created per MOSA run for ease of analysis.   
### Bug Fixes
* Fixed a MOSA solution set decompression bug that saw multiple solutions with same variable values output. 
* Redundant CatchmentModel parameters around coarse-grained costing have now been removed.
* Parameter "ExplorableDecisionVariables" removed, as it does nothing.

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