#!/usr/bin/env bash

packagename=$(pwd | grep -Eo "\b(\w+)\W*$")

instanceName="${packagename}-Clickhouse"
echo 
echo "building ${instanceName}"
sudo multipass launch --name "${instanceName}"

sudo multipass exec ${instanceName} -- bash -s < clickhouse_setup_local.sh