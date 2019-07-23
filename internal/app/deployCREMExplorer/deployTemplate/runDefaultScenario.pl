#!/usr/bin/perl
my $outputPath = 'output/log.txt';
open (my $file, '>', $outputPath) or die "Could not open file: $!";

my $output = `CREMExplorer.exe --ScenarioFile DefaultScenario.toml`; 
die "$!" if $?;

print $file $output;

print $output; 
print "\n\nPress <ENTER> to close window. Above log been written to \"$outputPath\".\n"; <>;