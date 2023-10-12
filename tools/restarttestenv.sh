#!/bin/bash
LocalDir=$(cd "$(dirname $0)";pwd)
echo $LocalDir
$LocalDir/stoptestenv.sh
$LocalDir/starttestenv.sh /home/rds/document/objectStorage/pkg