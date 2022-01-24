#!/bin/bash
set -x
SERVER_USER="ubuntu"
SERVER_HOST="10.204.28.21"
SSH_PORT=22
ssh -p $SSH_PORT $SERVER_USER@$SERVER_HOST 'bash -s' < clickhouse_backup_run_on_server.sh
rsync --remove-source-files -r -e "ssh -p ${SSH_PORT}" $SERVER_USER@$SERVER_HOST:/home/$SERVER_USER/ch_backup.tgz ./backup
mv ./backup/ch_backup.tgz ./backup/ch_backup_$(date '+%Y%m%d_%H%M%S').tgz

#wget -O 'https://github-releases.githubusercontent.com/150444746/e7087300-cec2-11eb-850c-5369b804e074?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20210830%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20210830T171747Z&X-Amz-Expires=300&X-Amz-Signature=73cb021c64453d37d392940f87536846360953c5c2956dc141d42b7f9a5d3620&X-Amz-SignedHeaders=host&actor_id=1061610&key_id=0&repo_id=150444746&response-content-disposition=attachment%3B%20filename%3Dclickhouse-backup_1.0.0_amd64.deb&response-content-type=application%2Foctet-stream'

