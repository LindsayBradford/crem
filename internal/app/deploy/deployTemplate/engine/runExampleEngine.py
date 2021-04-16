# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import os


def main():
    configFilePrefix = 'ExampleEngineConfig'
    configFile = f'{configFilePrefix}.toml'

    scenarioConfigFilePrefix = 'Example_SOSA_Scenario'
    scenarioConfigFile = f'{scenarioConfigFilePrefix}.toml'

    outputPath = f'output/LOG_{configFilePrefix}.txt'

    commandArray = ['CREMEngine.exe', '--EngineConfigFile', configFile, '--ScenarioFile', scenarioConfigFile, '2>&1', '>', outputPath]
    command = ' '.join(commandArray)

    print (f'\nRunning "{command}".\n\n')

    os.system(command)

    print (f'\n\nPress <ENTER> to close window. Above log been written to "{outputPath}".\n')


if __name__ == '__main__':
    main()
