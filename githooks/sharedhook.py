from enum import Enum

class ExitStatus(Enum):
    SUCCESS = 0
    FAILURE = 1

def reportStatusAndExit(errorCount, checkPrefix):
    exitStatus = ExitStatus.SUCCESS

    if errorCount > 0:
        exitStatus = ExitStatus.FAILURE

    print (f'{checkPrefix} {exitStatus.name}\n')

    exit(exitStatus.value)
