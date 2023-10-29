#!/usr/bin/env bash
set -x
CONN_STR=myusername:mysecretpassword@10.204.28.21:3301
SERVER_USER="ubuntu"
SERVER_HOST="10.204.28.21"
SSHPORT=22

BACKUP_DATE=$(date '+%Y%m%d_%H%M%S')
BACKUP_DIR=./backup/tt_${BACKUP_DATE}

echo 'box.snapshot()' | tarantoolctl connect $CONN_STR
echo 'box.backup.start()' | tarantoolctl connect $CONN_STR #> ./backup/tt.log
#cat ./backup/tt.log | grep '/var/lib/tarantool/' | cut -d '/' -f 5 > ./backup/tt_snapshot.txt
mkdir -p $BACKUP_DIR
rsync -avP -e "ssh -p ${SSHPORT}" $SERVER_USER@$SERVER_HOST:/home/$SERVER_USER/tarantool-data $BACKUP_DIR # --files-from=./backup/tt_snapshot.txt
echo 'box.backup.stop()' | tarantoolctl connect $CONN_STR 

tar cvfz ./backup/tt_backup_${BACKUP_DATE}.tgz -C $BACKUP_DIR .
rm -rf $BACKUP_DIR
