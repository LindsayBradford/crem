# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import subprocess
def main():
    configFilePrefix = 'Example_MOSA_Scenario'

    configFile = f'{configFilePrefix}.toml'
    outputPath = f'output/LOG_{configFilePrefix}.txt'

    commandArray = ['CREMExplorer.exe', '--ScenarioFile', configFile]
    command = ' '.join(commandArray)

    print (f'\nRunning "{command}", capturing output to "{outputPath}".\n\n')

    with open(outputPath, mode="w") as outputFile:
        output = subprocess.run(commandArray, capture_output=True, text=True)

        print(output.stdout, file=outputFile)
        print(output.stdout)

    print (f'\n\nPress <ENTER> to close window. Above log been written to "{outputPath}".\n')

if __name__ == '__main__':
    main()