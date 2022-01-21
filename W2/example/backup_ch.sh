#!/bin/bash
now=`date '+%Y%m%d_%H%M%S'`;
echo "andri" | sudo -S clickhouse-backup clean
echo "andri" | sudo -S clickhouse-backup create
echo "andri" | sudo tar -czf /var/lib/clickhouse/backup/ch_backup_$now.tgz -C /var/lib/clickhouse/backup/ .
rsync --delete -r -e "ssh -p 2211" andri@127.0.0.1:/var/lib/clickhouse/backup/* ~/Documents

