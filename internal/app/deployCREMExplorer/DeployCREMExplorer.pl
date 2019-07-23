#!/usr/local/bin/perl 

# Author: Lindsay Bradford, (c) Australian Rivers Institute, 2019.

# CREMEngine Deployment script.

use Archive::Zip;
use File::Basename;
use File::Copy;

# Removes whitespace from the head and tail of a string.
sub stringTrim() {
  my $string = $_[0];
  $string =~ s/^\s+//;
  $string =~ s/\s+$//;
  return $string;
}

sub initialise() {
  $sourceDir = "..\\..\\..\\cmd\\cremexplorer\\";
  $dataDir= "$sourceDir\\testdata\\";
  $targetTemplateDir = ".\\deployTemplate";
  $targetDir = ".\\deploy";
  $type = "Release";
  $binDir = "${sourceDir}";
  $baseExecutableName = "CREMExplorer";
  $executableName = "${binDir}\\${baseExecutableName}.exe";
}

sub getExecutableVerson($) {
  my $executable = $_[0];
  my $versionText = &stringTrim(`$executable --Version`);
  my @versionFields = split / /, $versionText;
  return $versionFields[1];
}

# Zips up an archive of the executable and its matching data.
 
sub zipDeploymentFiles() {

  my $versionNumber = getExecutableVerson($executableName);
  
  my $targetExecutableName = "${targetTemplateDir}\\${baseExecutableName}.exe";
  print "Copying $executableName to $targetExecutableName\n";
  copy($executableName, $targetExecutableName) or die "Copy failed: $!";

  my $changeLog = "${sourceDir}\\config\\ChangeLog.md";	
  my $targetChangeLogName = "${targetTemplateDir}\\ChangeLog.md";
  print "Copying $changeLog to $targetChangeLogName\n";
  copy($changeLog, $targetChangeLogName) or die "Copy failed: $!";
  
  my @directoriesToStore = (
     "${targetTemplateDir}",
  );

  my $zipFileName = "${targetDir}\\CREMExplorer_${versionNumber}.zip";

  my $zipHandle = Archive::Zip->new();   # new instance

  foreach $directory (@directoriesToStore) {
    print "Adding directory (${directory}) to archive (${zipFileName}).\n";
    $zipHandle->addTree($directory);  
  }

  if ($zipHandle->writeToFileNamed($zipFileName) != AZ_OK) {
    print "Error in archive creation!";
  } else {
    print "Archive created successfully!";
  }
}

####### Application below ########

initialise();
zipDeploymentFiles();
