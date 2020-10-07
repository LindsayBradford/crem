# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

import subprocess

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


def getExecutableVersion(config):
    commandArray = [config['executableName'], '--Version']
    output = subprocess.run(commandArray, capture_output=True, text=True)




if __name__ == '__main__':
    main()