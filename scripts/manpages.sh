#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run . man | gzip -c >manpages/org-stats.1.gz
