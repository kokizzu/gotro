#!/bin/bash
set -x

SUDOPASS=""

CHB_BIN="$(which clickhouse-backup)"
if [[ -z "${CHB_BIN}" ]] ; then
	wget -c https://github.com/AlexAkulov/clickhouse-backup/releases/download/v1.2.2/clickhouse-backup.tar.gz
	tar -xf clickhouse-backup.tar.gz
	cd clickhouse-backup/
	echo $SUDOPASS | sudo mv clickhouse-backup /usr/local/bin
  sudo mkdir -p /etc/clickhouse-backup
	sudo mv config.yml /etc/clickhouse-backup
	sudo chown ${USER} /etc/clickhouse-backup
	clickhouse-backup -v
	cd ..
	rmdir clickhouse-backup
	rm clickhouse-backup.tar.gz
fi

sudo usermod -aG clickhouse ${USER}

sudo rm -rf /var/lib/clickhouse/backup/
sudo mkdir -p /var/lib/clickhouse/backup/
sudo chown clickhouse:clickhouse /var/lib/clickhouse/backup
echo $SUDOPASS | sudo -S clickhouse-backup create
echo $SUDOPASS | sudo -S tar -czf ${HOME}/ch_backup.tgz -C /var/lib/clickhouse/backup/ .

