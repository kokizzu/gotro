#!/usr/bin/env bash

# synchronize service files: auto_backup, auto_restart, caddy

# journalctl -u namaservice
# -r reversed log
# -f follow log

# restart service: sudo systemctl restart example1_rest

SSH_USER=SSHUSER_CHANGEME
SERVER=SERVERHOST_CHANGEME
SSH_PORT=22 #22836
SSH_PARAM='-i ../SERVERPRIVATEKEY_CHANGEME.pem'  

# to main server
odir=${SSH_USER}@${SERVER}:/tmp/
rsync -h -t -P -r -e "ssh ${SSH_PARAM} -p ${SSH_PORT}" \
 example1_rest.service \
 Caddyfile \
 start_example1_rest.sh \
 $odir

ssh ${SSH_PARAM} ${SSH_USER}@${SERVER} -p ${SSH_PORT} 'sudo mv /tmp/*.service /lib/systemd/system/ && 
sudo systemctl daemon-reload && 
sudo systemctl enable example1_rest ; 
sudo mv /tmp/Caddyfile /etc/caddy/ && 
sudo systemctl reload caddy 
'
