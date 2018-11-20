#!/usr/bin/perl
# Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared;

use Exporter qw(import);
 
our @EXPORT_OK = qw(reportStatusAndExit);

use constant {
	SUCCESS => 0,
	FAILURE => 1,
};


sub reportStatusAndExit {
	my ($errorCount, $checkPrefix) = @_;

	my $statusString = "SUCCESS";
	my $exitStatus = SUCCESS;
	
	if ($errorCount > 0) {
		$statusString = "FAILURE";
		$exitStatus = FAILURE;
	} 

	print $checkPrefix, " $statusString\n";	
	
	exit $exitStatus;
}