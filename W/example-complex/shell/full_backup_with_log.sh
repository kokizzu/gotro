#!/usr/bin/bash

x=/home/CHANGEME/backup/full_backup--`date +%F_%H%M%S`_with_log--$RANDOM.sql.xz
echo "Drop all materialized views.."
( psql -t -U CHANGEME -c "SELECT 'DROP MATERIALIZED VIEW ' || string_agg(oid::regclass::text, ', ') 
FROM   pg_class
WHERE  relkind = 'm';" ) | psql -U CHANGEME 
echo "Dumping data to $x"
( psql -t -U CHANGEME -c "select 'drop table if exists \"' || tablename || '\" cascade;' from pg_tables where schemaname = 'public';
" && pg_dump -U CHANGEME CHANGEME -T 'vx_*' ) | /usr/bin/time -f 'CPU: %Us\tReal: %es\tRAM: %MKB' xz -3 -zf - > $x
echo "Backup file info:"
ls -al $x
echo "Restore with command:
 xzcat $x | psql -U CHANGEME"

## rsync to googledrive
#rc=`which rclone`
#if [ -n "$rc" ]; then
#	# backup database
#	rclone copy --retries 1 --include '*.xz' /backup gdrive:backup
#	# backup web
#	rclone copy --retries 1 /home/web/puis2 gdrive:puis2
#fi
#
