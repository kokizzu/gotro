#!/usr/bin/env bash

set -e # exit on error
set -x # print all command

TAG=$(ruby -e 't = Time.now; print "v1.#{t.month+(t.year-2021)*12}%02d.#{t.hour}%02d" % [t.day, t.min]')
EXE_NAME="example-${TAG}.exe"
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.VERSION=${TAG}'" -o "${EXE_NAME}" .
echo "${EXE_NAME}"
