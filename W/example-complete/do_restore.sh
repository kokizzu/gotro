
fname=`ls -w 1 tmp/*sql.xz | tail -n 1`
xzcat $fname | pg_restore -U geo -d geo
