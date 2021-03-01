# Change Log

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