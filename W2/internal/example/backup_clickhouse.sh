#!/usr/bin/env bash

set -x

# Remote server information
SERVER_USER="user"
SERVER_HOST="benalu.dev"
SSH_PORT=22

BACKUP_DATE=$(date '+%Y%m%d_%H%M%S')
BACKUP_DIR=./backup

mkdir -p $BACKUP_DIR

CH_BACKUP_FILE=/home/$SERVER_USER/ch_backup.tgz

# Run clickhouse backup script on remote server
ssh -p $SSH_PORT $SERVER_USER@$SERVER_HOST 'bash -s' < backup_clickhouse_on_server.sh

# Copy backup file from remote server to local machine
rsync --remove-source-files -r -e "ssh -p ${SSH_PORT}" $SERVER_USER@$SERVER_HOST:$CH_BACKUP_FILE $BACKUP_DIR/ch_backup_$BACKUP_DATE.tgz
