#!/usr/bin/env bash
set -x

LAST_BACKUP=`ls -w 1 backup/tt_backup_*.tgz | tail -n 1`
echo ${LAST_BACKUP}
docker-compose stop

RESTOREDIR=./backup/tarantool-data
rm -rf $RESTOREDIR
mkdir -p $RESTOREDIR
tar xvfz $LAST_BACKUP -C $RESTOREDIR

