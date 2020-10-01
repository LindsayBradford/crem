# (c) 2020, Australian Rivers Institute
# Author: Lindsay Bradford
import sys
import os
import io
import re

def main():
    logFilePath = deriveLogFilePath()
    qualityFilePath = deriveQualityFilePath(logFilePath)
    
    with open(logFilePath) as logFile, open(qualityFilePath, mode="w") as qualityFile:
        deriveQualityMetrics(logFile, qualityFile)

def deriveLogFilePath():
    def verifyArgumentLength():
        if len(sys.argv) != 2:
           print("Log file not specified on command-line. Exiting...")
           sys.exit()

    def deriveArgumentFilePath():
        logFilePath = sys.argv[1]

        if not os.path.isfile(logFilePath):
           print(f"Log file path {logFilePath} does not exist. Exiting...")
           sys.exit()
       
        return logFilePath

    verifyArgumentLength()
    return deriveArgumentFilePath()

def deriveQualityFilePath(logFilePath):
    qualityFilePath = logFilePath.replace('LOG_', 'QUALITY_')
    return qualityFilePath

def deriveQualityMetrics(logFile, qualityFile):
    writeHeadings(qualityFile)
    writeEntries(logFile, qualityFile)

def writeHeadings(qualityFile):
    headings = ['Iteration',  'LastReturnedToBase', 'Solutions', 'Temperature']
    
    stringIO = io.StringIO()
    print(*headings,  sep=', ', file=stringIO)
    
    qualityFile.write(stringIO.getvalue())

def writeEntries(logFile, qualityFile):
    def writeLineEntry(line, qualityFile):
                                                 
        if '[FinishedIteration]' in line:

            # Any combination of digits and ',' between 'Iteration [' and the first (reluctant) '/' found.
            iteration = extractValueViaPattern(line, 'Iteration \[([\d,]*?)/')
            # Any combination of digits, '.', and ',' between 'Temperature [' and the first (reluctant) ']' found.
            temperature = extractValueViaPattern(line, 'Temperature \[([\d,\.]*?)\]')
            # Any combination of digits and ',' between 'ArchiveSize [' and the first (reluctant) ']' found.
            archiveSize = extractValueViaPattern(line, 'ArchiveSize \[([\d,]*?)\]')
            # Any combination of digits and ',' between 'LastReturnedToBase [' and the first (reluctant) ']' found.
            lastReturnedToBase = extractValueViaPattern(line, 'LastReturnedToBase \[([\d,]*?)\]')

            qualityLine = f'{iteration}, {lastReturnedToBase}, {archiveSize}, {temperature}\n'
            qualityFile.write(qualityLine)
            
    for line in logFile:
        lastReturnedToBase = writeLineEntry(line, qualityFile)

def extractValueViaPattern(line, pattern):
      match = re.search(pattern, line)
      rawValue = match.group(1)
      value = rawValue.replace(',','')  # remove the commas to ensure CSV output is valid.
      return value

if __name__ == '__main__':
    main()
