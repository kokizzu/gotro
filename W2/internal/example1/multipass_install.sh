#!/usr/bin/env bash

SNAP_BIN=$(which snap)
if [[ -z "${SNAP_BIN}" ]] ; then
	echo "need to install snap manually. visit https://snapcraft.io/docs/installing-snapd";
    exit;
fi

sudo snap install multipass
echo "if multipass not detected try restart computer first."
echo "if cant enter to an instance use: 'snap install multipass --edge' instead "