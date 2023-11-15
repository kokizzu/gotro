#!/usr/bin/env bash

set -x

SUDO_PASS=""

CH_USER="userC"
CH_PASS="passC"

CHB_BIN="$(which clickhouse-backup)"
if [[ -z "${CHB_BIN}" ]] ; then
  mkdir -p clickhouse-backup
  cd clickhouse-backup
  wget -c https://github.com/Altinity/clickhouse-backup/releases/download/v2.4.2/clickhouse-backup-linux-amd64.tar.gz
  tar -xf clickhouse-backup-linux-amd64.tar.gz --strip-components=3

  echo $SUDO_PASS | sudo -S mv clickhouse-backup /usr/local/bin
  echo $SUDO_PASS | sudo -S mkdir -p /etc/clickhouse-backup

  clickhouse-backup default-config > config.yml
  sed -e '/^ *skip_tables:/b ins' -e b -e ':ins' -e 'a\'$'\n''  - information_schema.\*\n  - INFORMATION_SCHEMA.\*' -e ': done' -e 'n;b done' config.yml
  sed -i "0,/username/ s/username:.*/username: \"${CH_USER}\"/g" config.yml
  sed -i "0,/password/ s/password:.*/password: \"${CH_PASS}\"/g" config.yml
  echo $SUDO_PASS | sudo -S mv config.yml /etc/clickhouse-backup
  echo $SUDO_PASS | sudo -S chown ${USER} /etc/clickhouse-backup
  clickhouse-backup -v

  cd ..
  rm clickhouse-backup-linux-amd64.tar.gz
  rmdir clickhouse-backup
fi

LAST_BACKUP=`ls -w 1 backup/ch_backup_*.tgz | tail -n 1`

sudo rm -rf /var/lib/clickhouse/backup/
sudo mkdir -p /var/lib/clickhouse/backup/
sudo chown ${user}:clickhouse /var/lib/clickhouse/backup
sudo tar xvfz ${LAST_BACKUP} -C /var/lib/clickhouse/backup/

sudo clickhouse-backup restore `ls /var/lib/clickhouse/backup`
