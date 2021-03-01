# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import subprocess
import shutil

def main():
    config = deriveConfiguration()
    generateDeployArchive(config)

def deriveConfiguration():
    explorerSourceDir = '../../../cmd/cremexplorer/'
    engineSourceDir = '../../../cmd/cremengine/'

    baseArchiveName = 'CREM'
    baseExplorerName = f'{baseArchiveName}Explorer'
    baseEngineName = f'{baseArchiveName}Engine'

    targetTemplateDir = './deployTemplate'
    return {
        'baseArchiveName': baseArchiveName,
        'explorerSourceDir': explorerSourceDir,

        'targetTemplateDir': targetTemplateDir,
        'targetDir': './deploy',

        'sourceExplorerChangeLog': f'{explorerSourceDir}/config/ChangeLog.md',
        'targetExplorerChangeLog': f'{targetTemplateDir}/explorer/ChangeLog.md',

        'baseExplorerName': baseExplorerName,
        'sourceExplorer': f'{explorerSourceDir}/{baseExplorerName}.exe',
        'targetExplorer': f'{targetTemplateDir}/explorer/{baseExplorerName}.exe'
,
        'sourceEngineChangeLog': f'{engineSourceDir}/config/ChangeLog.md',
        'targetEngineChangeLog': f'{targetTemplateDir}/engine/ChangeLog.md',

        'baseEngineName': baseEngineName,
        'sourceEngine': f'{engineSourceDir}/{baseEngineName}.exe',
        'targetEngine': f'{targetTemplateDir}/engine/{baseEngineName}.exe'
    }

def generateDeployArchive(config):
    updateExplorerDeployTemplate(config)
    updateEngineDeployTemplate(config)
    generateArchiveFromTemplate(config)

def updateExplorerDeployTemplate(config):
    updateTemplate(config['sourceExplorer'], config['targetExplorer'])
    updateTemplate(config['sourceExplorerChangeLog'], config['targetExplorerChangeLog'])

def updateEngineDeployTemplate(config):
    updateTemplate(config['sourceEngine'], config['targetEngine'])
    updateTemplate(config['sourceEngineChangeLog'], config['targetEngineChangeLog'])

def updateTemplate(sourceFile, targetFile):
    print (f'Copying {sourceFile} to {targetFile}\n')
    shutil.copy(sourceFile, targetFile)

def generateArchiveFromTemplate(config):
    versionNumber = getExecutableVersion(config)
    zipFileName = f'{config["targetDir"]}/{config["baseArchiveName"]}_{versionNumber}'

    print (f'Adding directory ({config["targetTemplateDir"]}) to archive ({zipFileName}.zip).\n')
    shutil.make_archive(zipFileName, 'zip', config["targetTemplateDir"])        

def getExecutableVersion(config):
    commandArray = [config['sourceExplorer'], '--Version']
    output = subprocess.run(commandArray, capture_output=True, text=True)
    version = output.stdout.split()[1]
    return version

def logCommand(commandArray):
    command = ' '.join(commandArray)
    print (f'\nRunning "{command}".\n\n')

if __name__ == '__main__':
    main()
