#!/usr/bin/env bash

# preinstall all dependencies
#go build -i -x
echo 'set ownership of $GOROOT..'
old_owner=`ls -ld $GOROOT/pkg | awk '{print $3}'`
[ "$old_owner" != `whoami` ] && sudo chown -Rv `whoami` $GOROOT/pkg
echo 'remove $GOPATH/pkg if go upgraded/downgraded..'
PKG_OS_ARCH=`go version | cut -d ' ' -f 4 | tr '/' '_'`
old_version=`strings $GOPATH/pkg/$PKG_OS_ARCH/github.com/kokizzu/gotro/A.a | grep 'go object' | head -n 1 | cut -f 5 -d ' '`
cur_version=`go version | cut -f 3 -d ' '`
[ "$old_version" != "$cur_version" ] && rm -rf $GOPATH/pkg
if [ "$1" != "skip" ]; then
	echo 'precompile all dependencies..'
	go build -i -v
fi

# -i for auto compile after every changes, default: every request and changes
# killall -9 gin 2 > /dev/null
# killall -9 gin-bin 2 > /dev/null
gin=`which gin`
if [ -z "$gin" ]; then
	gin=$GOPATH/bin/gin
fi
echo 'starting gin..'
$gin -i -a 3001 -p 3000 -b CHANGEME

