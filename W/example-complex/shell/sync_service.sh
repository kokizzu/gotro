#!/usr/bin/env bash

# synchronize service files: auto_backup, auto_restart, caddy, puis2, pups, node_control

# journalctl -u namaservice
# -r urut terbalik
# -f otomatis follow kalau ada yg baru

# restart service: sudo systemctl restart namaservice

HOST=root@CHANGEME
PORT=22 

# to main server
odir=${HOST}:/usr/lib/systemd/system/
rsync -h -t -P -r -e "ssh -p ${PORT}" \
 auto_restart.service \
 CHANGEME.service \
 caddy.service \
 auto_backup.service \
 $odir

ssh ${HOST} -p ${PORT} 'systemctl daemon-reload'