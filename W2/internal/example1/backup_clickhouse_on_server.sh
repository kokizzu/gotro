#!/usr/bin/env bash

set -x

SUDO_PASS=""

CH_USER="userC"
CH_PASS="passC"
CH_BACKUP_FILE=${HOME}/ch_backup.tgz

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

  rm clickhouse-backup-linux-amd64.tar.gz
  cd ..
  rmdir clickhouse-backup
fi

echo $SUDO_PASS | sudo -S groupadd clickhouse
echo $SUDO_PASS | sudo -S usermod -aG clickhouse ${USER}

echo $SUDO_PASS | sudo -S rm -rf /var/lib/clickhouse/backup/
echo $SUDO_PASS | sudo -S mkdir -p /var/lib/clickhouse/backup/
echo $SUDO_PASS | sudo -S chown ${USER}:clickhouse /var/lib/clickhouse/backup
echo $SUDO_PASS | sudo -S clickhouse-backup create
echo $SUDO_PASS | sudo -S tar -czf $CH_BACKUP_FILE -C /var/lib/clickhouse/backup/ .
