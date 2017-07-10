#!/usr/bin/env bash

# synchronize service files: auto_backup, auto_restart, caddy

# journalctl -u namaservice
# -r urut terbalik
# -f otomatis follow kalau ada yg baru

# restart service: sudo systemctl restart namaservice

USER=root
HOST=TODO_CHANGE_DOMAIN
PORT=22 #22836 

# to main server
odir=${USER}@${HOST}:/lib/systemd/system/
rsync -h -t -P -r -e "ssh -p ${PORT}" \
 TODO_WEBAPP_SERVICE.service \
 caddy.service \
 auto_backup.service \
 TODO_CRON_SERVICE.service \
 $odir

ssh ${USER}@${HOST} -p ${PORT} 'systemctl daemon-reload'