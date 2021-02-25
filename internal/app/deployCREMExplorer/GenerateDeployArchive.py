# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import subprocess
import shutil

def main():
    config = initialise()
    zipDeploymentFile(config)


def initialise():
    sourceDir = '../../../cmd/cremexplorer/'
    binDir = sourceDir
    dataDir = f'{sourceDir}/testdata/'

    baseExecutableName = 'CREMExplorer'
    executableName = f'{binDir}/{baseExecutableName}.exe'

    return {
        'sourceDir': sourceDir,
        'dataDir': dataDir,
        'targetTemplateDir': './deployTemplate',
        'targetDir': './deploy',
        'type': 'Release',
        'baseExecutableName': baseExecutableName,
        'executableName': executableName
    }


def zipDeploymentFile(config):
    versionNumber = getExecutableVersion(config)

    targetTemplateDir = config['targetTemplateDir']
    baseExecutableName = config['baseExecutableName']
    targetExecutableName = f'{targetTemplateDir}/{baseExecutableName}.exe'

    executableName = config['executableName']
    print (f'Copying {executableName} to {targetExecutableName}\n')
    shutil.copy(executableName, targetExecutableName)

    sourceDir = config['sourceDir']
    changeLog = f'{sourceDir}/config/ChangeLog.md'
    targetChangeLogName = f'{targetTemplateDir}/ChangeLog.md'
    print (f'Copying {changeLog} to {targetChangeLogName}\n')
    shutil.copy(changeLog, targetChangeLogName)
   
    targetDir = config['targetDir']
    zipFileName = f'{targetDir}/CREMExplorer_{versionNumber}'

    print (f'Adding directory ({targetTemplateDir}) to archive ({zipFileName}.zip).\n')
    shutil.make_archive(zipFileName, 'zip', targetTemplateDir)        

 
def getExecutableVersion(config):
    commandArray = [config['executableName'], '--Version']
    output = subprocess.run(commandArray, capture_output=True, text=True)
    version = output.stdout.split()[1]
    return version

def logCommand(commandArray):
    command = ' '.join(commandArray)
    print (f'\nRunning "{command}".\n\n')

if __name__ == '__main__':
    main()
