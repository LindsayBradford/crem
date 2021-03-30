# (c) 2021, Australian Rivers Institute
# Author: Lindsay Bradford

import shutil, os


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
        'versionNumber': 'v0.3',
        
        'baseArchiveName': baseArchiveName,
        'explorerSourceDir': explorerSourceDir,

        'targetTemplateDir': targetTemplateDir,
        'targetDir': './deploy',

        'sourceExplorerChangeLog': f'{explorerSourceDir}/config/ChangeLog.md',
        'targetExplorerChangeLog': f'{targetTemplateDir}/explorer/ChangeLog.md',

        'baseExplorerName': baseExplorerName,
        'explorerMain': f'{explorerSourceDir}/main.go',
        'sourceExplorer': f'{explorerSourceDir}/{baseExplorerName}.exe',
        'targetExplorer': f'{targetTemplateDir}/explorer/{baseExplorerName}.exe'
,
        'sourceEngineChangeLog': f'{engineSourceDir}/config/ChangeLog.md',
        'targetEngineChangeLog': f'{targetTemplateDir}/engine/ChangeLog.md',

        'baseEngineName': baseEngineName,
        'engineMain': f'{engineSourceDir}/main.go',
        'sourceEngine': f'{engineSourceDir}/{baseEngineName}.exe',
        'targetEngine': f'{targetTemplateDir}/engine/{baseEngineName}.exe'
    }


def generateDeployArchive(config):
    print(f'Generating deployable {config["baseArchiveName"]} archive {config["versionNumber"]}...\n')
    
    updateExplorerDeployTemplate(config)
    updateEngineDeployTemplate(config)
    generateArchiveFromTemplate(config)


def updateExplorerDeployTemplate(config):
    print(f'Updating {config["baseExplorerName"]} deploy template from source repository...')

    compileGoProgram(config['explorerMain'], config['sourceExplorer'])
    updateTemplate(config['sourceExplorer'], config['targetExplorer'])
    updateTemplate(config['sourceExplorerChangeLog'], config['targetExplorerChangeLog'])

    print('\n')


def updateEngineDeployTemplate(config):
    print(f'Updating {config["baseEngineName"]} deploy template from source repository...')

    compileGoProgram(config['engineMain'], config['sourceEngine'])
    updateTemplate(config['sourceEngine'], config['targetEngine'])
    updateTemplate(config['sourceEngineChangeLog'], config['targetEngineChangeLog'])

    print('\n')


def compileGoProgram(mainGoFile, targetExecutableName):
    commandArray = ['go', 'build', '-o', targetExecutableName, mainGoFile]
    runCommand(commandArray) 


def runCommand(commandArray):
    command = ' '.join(commandArray)
    print (f'  Running "{command}"...')
    os.system(command)


def updateTemplate(sourceFile, targetFile):
    print (f'  Copying ({sourceFile}) to ({targetFile})...')
    shutil.copy(sourceFile, targetFile)


def generateArchiveFromTemplate(config):
    zipFileName = f'{config["targetDir"]}/{config["baseArchiveName"]}_{config["versionNumber"]}'

    print (f'Adding directory ({config["targetTemplateDir"]}) to archive ({zipFileName}.zip)...')
    shutil.make_archive(zipFileName, 'zip', config["targetTemplateDir"])        
    print (f'Creation of archive ({zipFileName}.zip) complete.')


if __name__ == '__main__':
    main()
