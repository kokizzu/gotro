#!/usr/bin/env bash

set -x

# Remote server information
SERVER_USER="user"
SERVER_HOST="benalu.dev"
SSH_PORT=22

# Tarantool information
TT_USER="userT"
TT_PASS="passT"
TT_PORT=3301
TT_PORT_LOCAL=3302
TT_DATA_DIR=/home/$SERVER_USER/tmpdb/var_lib_tarantool
CONN_STR=$TT_USER:$TT_PASS@localhost:$TT_PORT_LOCAL

BACKUP_DATE=$(date '+%Y%m%d_%H%M%S')
BACKUP_DIR=./backup/tt_${BACKUP_DATE}

mkdir -p $BACKUP_DIR
ssh -N -f -L $TT_PORT_LOCAL:localhost:$TT_PORT $SERVER_USER@$SERVER_HOST
echo 'box.snapshot()' | tarantoolctl connect $CONN_STR
echo 'box.backup.start()' | tarantoolctl connect $CONN_STR > ./backup/tt_backup.log
cat ./backup/tt_backup.log | grep '/var/lib/tarantool/' | cut -d '/' -f 5 > ./backup/tt_snapshot.txt
rsync -avP -e "ssh -p ${SSH_PORT}" $SERVER_USER@$SERVER_HOST:$TT_DATA_DIR $BACKUP_DIR --files-from=./backup/tt_snapshot.txt
echo 'box.backup.stop()' | tarantoolctl connect $CONN_STR

tar czvf ./backup/tt_backup_${BACKUP_DATE}.tgz -C $BACKUP_DIR .
rm -rf $BACKUP_DIR
pgrep -n ssh | xargs kill
