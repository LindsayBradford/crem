# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import subprocess
import shutil

def main():
    config = deriveConfiguration()
    generateDeployArchive(config)

def deriveConfiguration():
    explorerSourceDir = '../../../cmd/cremexplorer/'
    baseExecutableName = 'CREMExplorer'
    targetTemplateDir = './deployTemplate'
    return {
        'explorerSourceDir': explorerSourceDir,

        'targetTemplateDir': targetTemplateDir,
        'targetDir': './deploy',

        'sourceChangeLog': f'{explorerSourceDir}/config/ChangeLog.md',
        'targetChangeLog': f'{targetTemplateDir}/ChangeLog.md',

        'baseExecutableName': baseExecutableName,
        'sourceExecutable': f'{explorerSourceDir}/{baseExecutableName}.exe',
        'targetExecutable': f'{targetTemplateDir}/{baseExecutableName}.exe'
    }

def generateDeployArchive(config):
    updateExplorerDeployTemplate(config)
    generateArchiveFromTemplate(config)

def updateExplorerDeployTemplate(config):
    updateTemplate(config['sourceExecutable'], config['targetExecutable'])
    updateTemplate(config['sourceChangeLog'], config['targetChangeLog'])

def updateTemplate(sourceFile, targetFile):
    print (f'Copying {sourceFile} to {targetFile}\n')
    shutil.copy(sourceFile, targetFile)

def generateArchiveFromTemplate(config):
    versionNumber = getExecutableVersion(config)
    zipFileName = f'{config["targetDir"]}/CREMExplorer_{versionNumber}'

    print (f'Adding directory ({config["targetTemplateDir"]}) to archive ({zipFileName}.zip).\n')
    shutil.make_archive(zipFileName, 'zip', config["targetTemplateDir"])        

def getExecutableVersion(config):
    commandArray = [config['sourceExecutable'], '--Version']
    output = subprocess.run(commandArray, capture_output=True, text=True)
    version = output.stdout.split()[1]
    return version

def logCommand(commandArray):
    command = ' '.join(commandArray)
    print (f'\nRunning "{command}".\n\n')

if __name__ == '__main__':
    main()
