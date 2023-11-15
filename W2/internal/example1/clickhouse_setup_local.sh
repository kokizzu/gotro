#!/bin/bash
sudo apt-get install -y apt-transport-https ca-certificates dirmngr
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv E0C56BD4

echo "deb https://repo.clickhouse.tech/deb/stable/ main/" | sudo tee \
    /etc/apt/sources.list.d/clickhouse.list
sudo apt-get -y update

sudo apt-get install -y clickhouse-server clickhouse-client

sudo clickhouse start

CHB_BIN="$(which clickhouse-backup)"
if [[ -z "${CHB_BIN}" ]] ; then
	wget -c https://github.com/AlexAkulov/clickhouse-backup/releases/download/v1.3.2/clickhouse-backup-linux-amd64.tar.gz
	tar -xf clickhouse-backup-linux-amd64.tar.gz
	cd build/linux/amd64
	sudo mv clickhouse-backup /usr/local/bin
	sudo mkdir -p /etc/clickhouse-backup
	sed -e '/^ *skip_tables:/b ins' -e b -e ':ins' -e 'a\'$'\n''  - information_schema.\*\n  - INFORMATION_SCHEMA.\*' -e ': done' -e 'n;b done' config.yml
	sudo mv config.yml /etc/clickhouse-backup
	sudo chown ${USER} /etc/clickhouse-backup
	clickhouse-backup -v
	cd ~
	rm -rf build
	rm clickhouse-backup-linux-amd64.tar.gz
fi