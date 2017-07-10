#!/usr/bin/env bash

# TODO: --exclude must be added when there are new folder on the server that do not want to be overwritten with local when deploy

alias tcmd='/usr/bin/time -f "CPU: %Us\tReal: %es\tRAM: %MKB"'
alias upx=goupx

echo '> Doing replacements:'

# replace DefaultMaxRequestBodySize with bigger one
sed -i 's/const DefaultMaxRequestBodySize.*/const DefaultMaxRequestBodySize = 512 * 1024 * 1024/' ${GOPATH}/src/github.com/valyala/fasthttp/server.go
cat ${GOPATH}/src/github.com/valyala/fasthttp/server.go | grep 'const DefaultMaxRequestBodySize'

sed -i 's/case "1", "y", "yes":/case "1", "y", "yes", "true", "on":/' ${GOPATH}/src/github.com/valyala/fasthttp/args.go
cat ${GOPATH}/src/github.com/valyala/fasthttp/args.go | grep 'case "1", "y", "yes", "true", "on":'

if [ "$GOPATH" == "" ]; then
  echo "> GOPATH not set, auto-set to compile-node's GOPATH.."
  GOPATH=/home/`whomai`/go
fi

echo '> Compiling:'

SSH_USER=root
SERVER=TODO_CHANGE_DOMAIN
SSH_PORT=22 #22836
WEB_PORT=8083
WEB_USER=web
WEB_GROUP=users
HOME_DIR=/home/${WEB_USER}/site

# build example-complete
SERVER_DIR=${SSH_USER}@${SERVER}:${HOME_DIR}
if [ "$(uname)" == 'Darwin' ]; then
  echo '> Error: Unable deploy from MacOSX'
  exit 10
fi
BUILD_DATE=`date +.%Y%m%d.%H%M%S`
PROJECT=example-complete-example-cron
FLAGS="
  -X main.VERSION=${BUILD_DATE}
  -X main.PROJECT_NAME=${PROJECT}
  -X main.DOMAIN=TODO_CHANGE_DOMAIN
"

PIDS=''

echo '> Building example-complete web server..' \
&& go build -ldflags "
  -X main.LISTEN_ADDR=:${WEB_PORT}
	${FLAGS}
" -o /tmp/example-complete \
&& cp /tmp/example-complete /tmp/example-complete-raw &
PIDS="$PIDS $!"

echo '> Building example-cron background service..' \
&& go build -ldflags "
  -X main.LISTEN_ADDR=:${WEB_PORT}
	${FLAGS}
" -o /tmp/example-cron \
&& cp /tmp/example-cron /tmp/example-cron-raw &
PIDS="$PIDS $!"


for pid in $PIDS; do
	wait $pid || let 'fail=1'
done

if [ "$fail" == '1' ]; then
	echo '> One or more build process failed..'
	exit 11
fi

EXECUTABLES=''
REMOTE_CMD='systemctl restart paling_baik;'

PIDS=''

echo '> Comparing binaries..' \

b_diff=`cmp -l /tmp/example-complete-raw /tmp/example-complete-${PROJECT} | wc -l`
if [ "$b_diff" -gt 6 ] || [ "$b_diff" == 0 ]; then
	EXECUTABLES="$EXECUTABLES /tmp/example-complete"
#	REMOTE_CMD="systemctl restart example-complete; $REMOTE_CMD"
	echo "> example-complete service ${b_diff} bytes differ, compressing.."
	upx --no-progress /tmp/example-complete &
	PIDS="$PIDS $!"
fi

b_diff=`cmp -l /tmp/example-cron-raw /tmp/example-cron-${PROJECT} | wc -l`
if [ "$b_diff" -gt 6 ] || [ "$b_diff" == 0 ]; then
	EXECUTABLES="$EXECUTABLES /tmp/example-cron"
	REMOTE_CMD="systemctl restart menolak_lupa; $REMOTE_CMD"
	echo "> example-cron service ${b_diff} bytes differ, compressing.."
	upx --no-progress /tmp/example-cron &
	PIDS="$PIDS $!"
fi

if [ "${EXECUTABLES}" != '' ] ; then
	echo ${BUILD_DATE} > public/last_deploy
	echo `git log -n 1 | head -n 4` >> public/last_deploy
	echo '> Generating API docs..'
	pushd . &&
	cd go/apidocs &&
	go run gen_apidoc.go &&
	popd & 
	PIDS="$PIDS $!"
fi

wait ${PIDS}

echo "> Moving executables and scripts.. ${EXECUTABLES}" \
&& rsync -L -h -t -P -r -e "ssh -p ${SSH_PORT}" run_*.sh shell/Caddyfile shell ${EXECUTABLES} ${SERVER_DIR} \
&& echo '> Sychronizing current release..' \
&& rsync -L -h -t -P -r --delete -e "ssh -p ${SSH_PORT}" \
--exclude '.*' \
--exclude 'public/pictures/*' \
--exclude 'public/videos/*' \
--exclude 'public/profpics/*' \
--exclude 'public/js/all.js' \
--exclude 'public/css/all.css' \
--exclude 'public/js/lib.js' \
--exclude 'public/css/lib.css' \
--exclude 'public/node_modules' \
--exclude 'public/cache/*' \
--exclude '_old' \
--exclude 'go' \
--exclude 'resources' \
--exclude 'git_stats' \
--exclude '.sublime-project' \
--exclude '*.sql' \
--exclude 'logs*' \
--exclude 'example-complete' \
--exclude 'example-cron' \
--exclude 'gin-bin' \
--exclude 'handler' \
--exclude 'ruby*' \
--exclude 'phantomjs*' \
--exclude 'arch' \
--exclude 'sql*' \
--exclude 'db*' \
--exclude 'tmp' \
--exclude '*.go' \
--exclude '*.sh' \
--exclude '*.rb' \
--exclude '*.txt' \
--exclude '*.iml' \
--exclude '*.java' \
--exclude '*.log' \
--exclude '*.service' \
--exclude 'Caddyfile' \
. ${SERVER_DIR} \
&& echo "> Executing remote commands.. ${REMOTE_CMD}" \
&& ssh ${SSH_USER}@${SERVER} -p ${SSH_PORT} "
chown -R ${WEB_USER}:${WEB_GROUP} ${HOME_DIR} ; 
${REMOTE_CMD}" \
&& mv /tmp/example-complete-raw /tmp/example-complete-${PROJECT} \
&& mv /tmp/example-cron-raw /tmp/example-cron-${PROJECT} 

# TODO: tambahkan di sebelum baris terakhir, kalau SPA sudah selesai
# scp -rpP ${SSH_PORT} public/SPA/dist/* ${SERVICE_DIR}/public/