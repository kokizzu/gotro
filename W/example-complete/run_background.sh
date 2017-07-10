#!/usr/bin/env bash

mkdir -p `pwd`/logs
ofile=`pwd`/logs/bg_`date +%F_%H%M%S`.log
echo Logging into: $ofile
unbuffer time ./example-cron | tee $ofile

