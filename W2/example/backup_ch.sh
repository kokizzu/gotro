#!/bin/bash
set -x
SUDOPASS="andri"

now=`date '+%Y%m%d_%H%M%S'`;
isExist="$(which clickhouse-backup)"

if [[ -z "${isExist}" ]] ; then
	echo $SUDOPASS | sudo -S wget -c https://github.com/AlexAkulov/clickhouse-backup/releases/download/v0.5.2/clickhouse-backup.tar.gz
	echo $SUDOPASS | sudo -S tar -xf clickhouse-backup.tar.gz
	cd clickhouse-backup/
	sudo cp clickhouse-backup /usr/local/bin
	clickhouse-backup -v
fi

echo $SUDOPASS | sudo -S rm -rf /var/lib/clickhouse/backup/*
echo $SUDOPASS | sudo -S clickhouse-backup clean
echo $SUDOPASS | sudo -S clickhouse-backup create
echo $SUDOPASS | sudo tar -czf /var/lib/clickhouse/backup/ch_backup_$now.tgz -C /var/lib/clickhouse/backup/ .
