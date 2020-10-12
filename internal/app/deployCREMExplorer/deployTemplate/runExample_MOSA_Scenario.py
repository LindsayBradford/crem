# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import os


def main():
    configFilePrefix = 'Example_MOSA_Scenario'

    configFile = f'{configFilePrefix}.toml'
    outputPath = f'output/LOG_{configFilePrefix}.txt'

    commandArray = ['CREMExplorer.exe', '--ScenarioFile', configFile, '2>&1', '>', outputPath]
    command = ' '.join(commandArray)

    print (f'\nRunning "{command}".\n\n')

    os.system(command)

    print (f'\n\nPress <ENTER> to close window. Above log been written to "{outputPath}".\n')


if __name__ == '__main__':
    main()