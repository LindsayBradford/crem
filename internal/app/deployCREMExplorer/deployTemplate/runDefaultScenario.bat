@ECHO OFF

TITLE Running DefaultScenario of CREMExplorer

SET OUTPUT_PATH="output\log.txt"
CREMExplorer.exe --ScenarioFile DefaultScenario.toml > %OUTPUT_PATH%

TYPE %OUTPUT_PATH%

ECHO.
ECHO Above log has been written to %OUTPUT_PATH%.     
PAUSE