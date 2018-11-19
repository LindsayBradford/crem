#!/usr/bin/perl

my $rootRepositoryDirectory = "../../../";

system("cloc-1.78.exe", $rootRepositoryDirectory, "--exclude-dir=vendor");  # https://github.com/AlDanial/cloc
system("gource", "--load-config", "gource.config");                         # https://github.com/acaudwell/Gource