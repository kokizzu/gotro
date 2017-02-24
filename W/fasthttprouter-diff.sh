#!/usr/bin/env bash

p0='github.com/buaazp/fasthttprouter'
echo "FETCH $p0"
go get -u -v $p0

for x in path router tree; do
	p1="fasthttprouter-$x.go"
	p2="$GOPATH/src/github.com/buaazp/fasthttprouter/$x.go"
	echo
	echo "DIFF $p1 $p2"
	diff --color=always $p1 $p2
done