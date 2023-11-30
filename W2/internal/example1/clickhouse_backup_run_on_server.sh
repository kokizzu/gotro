#!/bin/bash
set -x

SUDO_PASS=""

CHB_BIN="$(which clickhouse-backup)"
if [[ -z "${CHB_BIN}" ]] ; then
	wget -c https://github.com/AlexAkulov/clickhouse-backup/releases/download/v1.3.2/clickhouse-backup.tar.gz
	tar -xf clickhouse-backup.tar.gz
	cd clickhouse-backup/
	echo $SUDO_PASS | sudo mv clickhouse-backup /usr/local/bin
  sudo mkdir -p /etc/clickhouse-backup
  sed -e '/^ *skip_tables:/b ins' -e b -e ':ins' -e 'a\'$'\n''  - information_schema.\*\n  - INFORMATION_SCHEMA.\*' -e ': done' -e 'n;b done' config.yml
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
echo $SUDO_PASS | sudo -S clickhouse-backup create
echo $SUDO_PASS | sudo -S tar -czf ${HOME}/ch_backup.tgz -C /var/lib/clickhouse/backup/ .

