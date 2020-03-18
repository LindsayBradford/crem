#!/usr/bin/perl
my $configFilePrefix = 'Example_SOSA_Scenario';
my $configFile = "$configFilePrefix.toml";

my $outputPath = "output/LOG_$configFilePrefix.txt";
open (my $file, '>', $outputPath) or die "Could not open file: $configFile!";

my $output = `CREMExplorer.exe --ScenarioFile $configFile`; 
die "$!" if $?;

print $file $output;

print $output; 
print "\n\nPress <ENTER> to close window. Above log been written to \"$outputPath\".\n"; <>;