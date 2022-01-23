#!/usr/bin/env bash
set -x
CONNSTR=myusername:mysecretpassword@10.204.28.21:3301
SERVERUSER="ubuntu"
SERVERHOST="10.204.28.21"
SSHPORT=22

BACKDATE=$(date '+%Y%m%d_%H%M%S')
BACKDIR=./backup/tt_${BACKDATE}

echo 'box.backup.start()' | tarantoolctl connect $CONNSTR #> ./backup/tt.log
#cat ./backup/tt.log | grep '/var/lib/tarantool/' | cut -d '/' -f 5 > ./backup/tt_snapshot.txt
mkdir -p ${BACKDIR}
rsync -avP -e "ssh -p ${SSHPORT}" $SERVERUSER@$SERVERHOST:/home/$SERVERUSER/tarantool-data ${BACKDIR} # --files-from=./backup/tt_snapshot.txt
echo 'box.backup.stop()' | tarantoolctl connect $CONNSTR 

tar cvfz ./backup/tt_backup_${BACKDATE}.tgz -C ${BACKDIR} .
rm -rf ${BACKDIR}
