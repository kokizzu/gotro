#!/usr/bin/env bash
set -x

LAST_BACKUP=`ls -w 1 backup/tt_backup_*.tgz | tail -n 1`
echo ${LAST_BACKUP}
docker-compose stop

RESTORE_DIR=./backup/tarantool-data
sudo rm -rf $RESTORE_DIR
mkdir -p $RESTORE_DIR
tar xvfz $LAST_BACKUP -C $RESTORE_DIR

# check with:
# tarantoolctl connect myusername:mysecretpassword@10.204.28.21:3301
