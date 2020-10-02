# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford

# http://cloc.sourceforge.net/

import subprocess


def main():
    showCountOfLinesOfCode()
    visualiseRepositoryViaGource()


def showCountOfLinesOfCode():
    rootRepositoryDirectory = "../../../"
    clocArray = ['cloc-1.88.exe', rootRepositoryDirectory]

    clocOutput = subprocess.run(clocArray, capture_output=True, text=True)
    print(clocOutput.stdout)


def visualiseRepositoryViaGource():
    gourceArray = ['gource', '--load-config', 'gource.config']
    subprocess.run(gourceArray, shell=True)


if __name__ == '__main__':
    main()