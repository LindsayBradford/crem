# Change Log

## Version 0.9 (06 June 2022):
### New Features
* Compiled against v0.22 of CremExplorer.

## Version 0.8 (08 October 2021):
### New Features
* Compiled against v0.20 of CremExplorer, permitting usage of new catchment model variable "TotalNitrogenProduction"

## Version 0.7 (24 September 2021):
### New Features
* Removal of running engine api behaviour:
  * GET /api/v1/model/subcatchment/[0-9]*/applicableActions -- Returns the mgt actions that can be applied for the given
    model's subcatchment
  * GET /api/v1/model/actions -- Returns the active management actions of the current running model
  * POST /api/v1/model/actions -- Supplies a solution of active management actions to the current running model
* Addition of new running engine api behaviour:
  * GET /api/v1/model/actions/applicable -- Returns the mgt actions that can be applied for the given model, broken down
    by subcatchment
  * GET /api/v1/model/actions/active -- Returns the active management actions of the current running model
  * PUT /api/v1/model/actions/active -- Supplies a solution of active management actions to the current running model

## Version 0.6 (22 September 2021):
### New Features
* Addition of new running engine api behaviour:
  * GET /api/v1/model/subcatchment/[0-9]*/applicableActions -- Returns the mgt actions that can be applied for the given
    model's subcatchment

## Version 0.5 (27/08/2021):
### New Features
* Warnings are now posted to the log for attempting to activate management actions via CSV tables that aren't supported.
### Bugs Fixed
* Fixed bug in model PATCH parser, where bad encodings were not reporting errors in the message responses.
* Fixed bug where key attributes were missing on solution GET responses after a successful solution set upload

## Version 0.4 (15 July 2021):
### New Features
* Addition of new running engine api behaviour:
  * POST /api/v1/solutions                   -- Supplies a solution summary to the engine, which the engine then uses
  * GET  /api/v1/solutions                   -- Returns the current solution summary data loaded into the engine
  * GET  /api/v1/solutions/<solution-label>  -- Returns full model solution detail for a solution present in the solution summary.
  * PATCH /api/v1/model                      -- Allows attributes to be uploaded to the currently running model.
    * An engine-reserved attribute 'Encoding' can be supplied that alters the model's management action state as per solution summaries.
    * An engine-reserved attribute 'ParetoFrontMember' reports whether the model is a member of the solution summary pareto front.
    * An engine-reserved attribute 'ValidAgainstScenario' reports on whether the model is valid against its scenario. 
    * An engine-reserved attribute 'ValidationErrors' reports why the model is invalid valid against its scenario.
* Command-line now allows a solution summary file to be specified via the new --SolutionSummaryFile command-line argument.
### Bugs Fixed
* POST /api/v1/model/actions is now a PUT operation to better match RESTful best-practices.
* POST /api/v1/model/subcatchment/[0-9]* is now a PUT operation to better match RESTful best-practices.
* Solution attributes are now being correctly exported with JSON encodings of model solutions.
* Engine will now correctly report errors to the log if a badly formed solution csv file is specified on the command-line.

## Version 0.3 (03 June 2021):
### Bugs Fixed
* Fixed bug in Subcatchment parser, stopping WetlandsEstablishment action from having state set.
* Fixed bug in Subcatchment parser, where semantically invalid action state changes would respond with success.
* Fixed bug in actions parser, where strings in the posted CSV table would trigger a panic attack.

## Version 0.2 (02 March 2021):
### New Features
* Addition of new running engine api behaviour:
  * GET  /api/v1/scenario                   -- Returns the configuration for the currently running scenario
  * POST /api/v1/scenario                   -- Supplies a scenario to the engine, which the engine then uses
  * GET  /api/v1/model                      -- Returns the state of the running model derived from supplied scenario
  * GET  /api/v1/model/actions              -- Returns the active management actions of the current running model
  * POST /api/v1/model/actions              -- Supplies a solution of active management actions to the current running model
  * GET  /api/v1/model/subcatchment/[0-9]*  -- Returns the state of management actions for the given model's subcatchment
  * POST /api/v1/model/subcatchment/[0-9]*  -- Supplies updates to the state of mgt actions for the given model's subcatchment
* Command-line now allows an initial scenario to be specified via --ScenarioFile <FileName>.
* Command-line now allows an initial scenario's management action state to be specified via --SolutionFile <FileName>.

## Version 0.1 (18 July 2019):
### New Features
* Basic CREMEngine executable, supplying only admin /status and /shutdown behaviour.