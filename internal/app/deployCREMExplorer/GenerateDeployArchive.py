# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import subprocess
import shutil

def main():
    config = deriveConfiguration()
    generateDeployArchive(config)

def deriveConfiguration():
    sourceDir = '../../../cmd/cremexplorer/'
    baseExecutableName = 'CREMExplorer'
    return {
        'sourceDir': sourceDir,

        'targetTemplateDir': './deployTemplate',
        'targetDir': './deploy',

        'baseExecutableName': baseExecutableName,
        'executableName': f'{sourceDir}/{baseExecutableName}.exe'
    }

def generateDeployArchive(config):
    updateExplorerDeployTemplate(config)
    generateArchiveFromTemplate(config)

def updateExplorerDeployTemplate(config):
    targetExecutableName = f'{config["targetTemplateDir"]}/{config["baseExecutableName"]}.exe'
    updateTemplate(config['executableName'], targetExecutableName)

    changeLog = f'{config["sourceDir"]}/config/ChangeLog.md'
    targetChangeLogName = f'{config["targetTemplateDir"]}/ChangeLog.md'
    updateTemplate(changeLog, targetChangeLogName)

def updateTemplate(sourceFile, targetFile):
    print (f'Copying {sourceFile} to {targetFile}\n')
    shutil.copy(sourceFile, targetFile)

def generateArchiveFromTemplate(config):
    versionNumber = getExecutableVersion(config)
    zipFileName = f'{config["targetDir"]}/CREMExplorer_{versionNumber}'

    print (f'Adding directory ({config["targetTemplateDir"]}) to archive ({zipFileName}.zip).\n')
    shutil.make_archive(zipFileName, 'zip', config["targetTemplateDir"])        

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
