#!/usr/bin/env bash
set -x 

CHB_BIN="$(which clickhouse-backup)"
if [[ -z "${CHB_BIN}" ]] ; then
	wget -c https://github.com/AlexAkulov/clickhouse-backup/releases/download/v1.2.2/clickhouse-backup.tar.gz
	tar -xf clickhouse-backup.tar.gz
	cd clickhouse-backup/
	sudo mv clickhouse-backup /usr/local/bin
	sudo mkdir -p /etc/clickhouse-backup
	sudo mv config.yml /etc/clickhouse-backup
	sudo chown ${USER} /etc/clickhouse-backup
	clickhouse-backup -v
	cd ..
	rmdir clickhouse-backup
	rm clickhouse-backup.tar.gz
fi

sudo usermod -aG clickhouse ${USER}

LAST_BACKUP=`ls -w 1 backup/ch_backup_*.tgz | tail -n 1`
echo ${LAST_BACKUP}

sudo rm -rf /var/lib/clickhouse/backup/
sudo mkdir -p /var/lib/clickhouse/backup/
sudo chown clickhouse:clickhouse /var/lib/clickhouse/backup
sudo tar xvfz ${LAST_BACKUP} -C /var/lib/clickhouse/backup/

sudo clickhouse-backup restore `sudo clickhouse-backup list | cut -d ' ' -f 1`
