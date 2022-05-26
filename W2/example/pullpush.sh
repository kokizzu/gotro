#!/usr/bin/env bash

if [ $# -eq 0 ] ; then
  echo "Usage: 
  ./pullpush.sh 'the commit message'"
  exit
fi

# format indentation
go fmt ./...
echo "codes formatted.."

# update deps
go get -u -v github.com/ory/dockertest@latest
go mod tidy -v

# testing if has "gokil" included
ag gokil **/*.go && ( echo 'echo should not import previous gokil library..' ; kill 0 )
ag luwes **/*.go && ( echo 'echo should not import luwes library..' ; kill 0 )
ag kaimono **/*.go && ( echo 'echo should not import kaimono library..' ; kill 0 )
ag jankenpon **/*.go && ( echo 'echo should not import jankenpon library..' ; kill 0 )
echo "imports checked.."

# add and commit all files
git add .
git status
read -p "Press Ctrl+C to exit, press any enter key to check the diff..
"

# recheck again
git diff --staged
echo 'Going to commit with message: '\"$*\"
read -p "Press Ctrl+C to exit, press any enter key to really commit..
"

git commit -m "$*" && git pull && git push origin main

#git tag -a `ruby -e 't = Time.now; print "v1.#{t.month+(t.year-2021)*12}%02d.#{t.hour}%02d" % [t.day, t.min]'` -m "$*"
#git push --tags 

# delete tag: 
# git tag -d v1.mdd.hhmm 
# git push -d origin v1.mdd.hhmm
