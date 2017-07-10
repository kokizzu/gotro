
ssh root@TODO_CHANGE_DOMAIN -p 2200 "sudo -u postgres pg_dump -Fc -c geo | xz - -c" \
 | pv -r -b > tmp/db_backup_`date +%Y-%m-%d_%H%M%S`.sql.xz
