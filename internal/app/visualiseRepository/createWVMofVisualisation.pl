#!/usr/bin/perl

my $visualisationVideoName = "repositoryVisualisation";
my $ppmFile = "$visualisationVideoName.ppm";
my $wmvFile = "$visualisationVideoName.wmv";

system("ffmpeg","-y","-r","60","-f","image2pipe","-vcodec","ppm","-i","$ppmFile","-vcodec","wmv1","-r","60","-qscale","0","$wmvFile")  # https://www.ffmpeg.org/download.html#build-windows