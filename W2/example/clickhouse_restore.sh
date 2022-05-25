#!/bin/bash
set -x
SERVER_USER="ubuntu"
SERVER_HOST="10.173.44.10"
SSH_PORT=22
SERVER_DIR=${SERVER_USER}@${SERVER_HOST}:/home/${SERVER_USER}
rsync -L -h -t -P -r -e "ssh ${SSH_PARAM} -p ${SSH_PORT}" backup ${SERVER_DIR}\
&& ssh -p $SSH_PORT $SERVER_USER@$SERVER_HOST 'bash -s' < clickhouse_restore_run_on_server.sh
