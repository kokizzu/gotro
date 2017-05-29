#!/bin/sh

# TODO: --exclude must be added when there are new folder on the server that do not want to be overwritten with local when deploy

alias tcmd='/usr/bin/time \-f\ "CPU: %Us\tReal: %es\tRAM: %MKB"'
alias upx=goupx

echo 'Doing replacements:'

# replace startWithSlate bug
sed -i 's/StartWithSlate bool `json:"startWithSlate,omitempty"`/StartWithSlate bool `json:"startWithSlate"`/' ${GOPATH}/src/google.golang.org/api/youtube/v3/youtube-gen.go
cat ${GOPATH}/src/google.golang.org/api/youtube/v3/youtube-gen.go | grep 'StartWithSlate bool'

# replace DefaultMaxRequestBodySize with bigger one
sed -i 's/const DefaultMaxRequestBodySize.*/const DefaultMaxRequestBodySize = 512 * 1024 * 1024/' ${GOPATH}/src/github.com/valyala/fasthttp/server.go
cat ${GOPATH}/src/github.com/valyala/fasthttp/server.go | grep 'const DefaultMaxRequestBodySize'

sed -i 's/case "1", "y", "yes":/case "1", "y", "yes", "true", "on":/' ${GOPATH}/src/github.com/valyala/fasthttp/args.go
cat ${GOPATH}/src/github.com/valyala/fasthttp/args.go | grep 'case "1", "y", "yes", "true", "on":'

if [ "$GOPATH" == "" ]; then
  echo "GOPATH not set, auto-set to compile-node's GOPATH.."
  GOPATH=/home/`whomai`/go
fi

echo 'Generating API documentations'

pushd . &&
cd go/apidocs &&
go run gen_apidoc.go &&
popd ||
exit 10 & 

echo 'Compiling:'

SERVER=CHANGEME.com
SSH_PORT=22
SVC_PORT=8083
WEB_USER=CHANGEME
WEB_GROUP=users
HOME_DIR=/home/${WEB_USER}/web

# build CHANGEME
SERVER_DIR=root@${SERVER}:${HOME_DIR}
if [ "$(uname)" == 'Darwin' ]; then
  echo 'Error: Unable deploy from MacOSX'
  exit 10
fi
BUILD_DATE=`date +.%Y%m%d.%H%M%S`

#  -X main.VERSION=$compile_date
echo ${BUILD_DATE} > public/last_deploy
echo `git log -n 1 | head -n 4` >> public/last_deploy

echo 'Building CHANGEME service..' \
&& tcmd go build -ldflags "
  -X main.LISTEN_ADDR=:${SVC_PORT}
" -o /tmp/CHANGEME \
&& echo 'Compressing CHANGEME..' \
&& tcmd upx /tmp/CHANGEME \
|| exit 11 

echo 'Moving executables..' \
&& tcmd rsync --progress -e "ssh -p ${SSH_PORT}" /tmp/CHANGEME ${SERVER_DIR} \
&& echo 'Copying scripts..' \
&& tcmd scp -pP ${SSH_PORT} server_stat.sh run_*.sh shell/Caddyfile shell/*.sh shell/*.rb ${SERVER_DIR} \
&& tcmd scp -pP ${SSH_PORT} server_stat.sh run_*.sh shell/Caddyfile shell/*.sh shell/*.rb ${SERVER_DIR}/shell/ \
&& echo 'Sychronizing current release..' \
&& tcmd rsync -L -h -t -P -r --delete --stats -e "ssh -p ${SSH_PORT}" \
--exclude '.*' \
--exclude 'public/js/all.js' \
--exclude 'public/css/all.css' \
--exclude 'public/js/lib.js' \
--exclude 'public/css/lib.css' \
--exclude '_old' \
--exclude 'go' \
--exclude 'resources' \
--exclude 'git_stats' \
--exclude '.sublime-project' \
--exclude '*.sql' \
--exclude 'logs*' \
--exclude 'gin-bin' \
--exclude 'handler' \
--exclude 'ruby*' \
--exclude 'shell' \
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
&& ssh root@${SERVER} -p ${SSH_PORT} "chown -R ${WEB_USER}:${WEB_GROUP} ${HOME_DIR}; systemctl restart CHANGEME"

# TODO: tambahkan di sebelum baris terakhir, kalau SPA sudah selesai
# tcmd scp -rpP ${SSH_PORT} public/SPA/dist/* ${SERVICE_DIR}/public/