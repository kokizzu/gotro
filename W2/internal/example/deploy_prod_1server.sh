#!/usr/bin/env bash

# TODO: --exclude must be added when there are new folder on the server that do not want to be overwritten with local when deploy

alias tcmd='/usr/bin/time -f "CPU: %Us\tReal: %es\tRAM: %MKB"'
alias upx=goupx

echo '> Doing replacements:'

# replace DefaultMaxRequestBodySize with bigger one
sed -i 's/const DefaultMaxRequestBodySize.*/const DefaultMaxRequestBodySize = 512 * 1024 * 1024/' vendor/github.com/valyala/fasthttp/server.go
cat vendor/github.com/valyala/fasthttp/server.go | grep 'const DefaultMaxRequestBodySize'

sed -i 's/case "1", "y", "yes":/case "1", "y", "yes", "true", "on":/' vendor/github.com/valyala/fasthttp/args.go
cat vendor/github.com/valyala/fasthttp/args.go | grep 'case "1", "y", "yes", "true", "on":'

if [ "$GOPATH" == "" ]; then
  echo "> GOPATH not set, auto-set to compile-node's GOPATH.."
  GOPATH=/home/`whoami`/go
fi

echo '> Compiling:'

SSH_USER=SSHUSER_CHANGEME
SERVER=SERVERHOST_CHANGEME
SSH_PORT=22
PROJECT=example
WEB_USER=SERVICEUSER_CHANGEME
WEB_GROUP=SERVICEGROUP_CHANGEME
HOME_DIR=/home/${SERVICEUSER_CHANGEME}/site
SSH_PARAM='-i SERVERPRIVATEKEY_CHANGEME.pem'

# build Luwes
SERVER_DIR=${SSH_USER}@${SERVER}:${HOME_DIR}
BUILD_DATE=`date +.%Y%m%d.%H%M%S`

export GOOS=linux
export GOARCH=amd64

echo "> Building ${PROJECT}.." \
&& go build -ldflags "
" -o /tmp/${PROJECT}.exe \
&& cp /tmp/${PROJECT}.exe /tmp/${PROJECT}.exe-raw & 
PIDS="$PIDS $!"

for pid in $PIDS; do
	wait $pid || let 'fail=1'
done

if [ "$fail" == '1' ]; then
	echo '> One or more build process failed..'
	exit 11
fi

EXECUTABLES=''
REMOTE_CMD='sudo systemctl restart example_rest; ' # 'killall -9 ${PROJECT}.exe'

PIDS=''

echo '> Comparing binaries..' \

b_diff=`cmp -l /tmp/${PROJECT}.exe-raw /tmp/${PROJECT}.exe-${PROJECT} | wc -l`
if [ "$b_diff" -gt 6 ] || [ "$b_diff" == 0 ]; then
	EXECUTABLES="$EXECUTABLES /tmp/${PROJECT}.exe"
#	REMOTE_CMD="sudo systemctl restart CHANGEME; $REMOTE_CMD"
	echo "> ${PROJECT}.exe ${b_diff} bytes differ, compressing.." 
	upx --no-progress /tmp/${PROJECT}.exe &
	PIDS="$PIDS $!"
fi

(cd svelte ; npm run build)

cp svelte/src/pages/api.js svelte/dist/
cp svelte/public/* svelte/dist/ 

# TODO: generate API docs (if gen-route not called)
#if [ "${EXECUTABLES}" != '' ] ; then
#	echo ${BUILD_DATE} > public/last_deploy
#	echo `git log -n 1 | head -n 4` >> public/last_deploy
#	echo '> Generating API docs..'
#	unset GOOS
#	unset GOARCH
#	pushd . &&
#	cd go/apidocs &&
#	go run gen_apidoc.go &&
#	popd & 
#	PIDS="$PIDS $!"
#fi


wait ${PIDS}

echo "> Moving executables and scripts.. ${EXECUTABLES}" \
&& rsync -L -h -t -P -r -e "ssh ${SSH_PARAM} -p ${SSH_PORT}" production/start_*.sh production/Caddyfile production/.env ${EXECUTABLES} ${SERVER_DIR} \
&& echo '> Sychronizing current release..' \
&& rsync -L -h -t -P -r --delete -e "ssh ${SSH_PARAM} -p ${SSH_PORT}" \
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
--exclude 'svelte/dist/upload' \
--exclude 'svelte/node_modules' \
--exclude 'svelte/src' \
--exclude 'svelte/src/vite-plugin-mpa' \
--exclude '_old' \
--exclude '3rdparty' \
--exclude 'handler' \
--exclude 'model' \
--exclude 'deploy' \
--exclude 'conf' \
--exclude 'domain' \
--exclude 'Makefile' \
--exclude 'Caddyfile' \
--exclude 'go*' \
--exclude 'resources' \
--exclude 'git_stats' \
--exclude '.sublime-project' \
--exclude '*.md' \
--exclude '.env' \
--exclude '*.Sql' \
--exclude '*.pem' \
--exclude '*.ppk' \
--exclude '*.yml' \
--exclude 'logs*' \
--exclude "${PROJECT}.exe" \
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
&& ssh ${SSH_PARAM} ${SSH_USER}@${SERVER} -p ${SSH_PORT} "
chown -R ${WEB_USER}:${WEB_GROUP} ${HOME_DIR} ; 
${REMOTE_CMD}" \
&& mv /tmp/${PROJECT}.exe-raw /tmp/${PROJECT}.exe-${PROJECT} 

