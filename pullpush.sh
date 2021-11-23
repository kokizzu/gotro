#!/usr/bin/env bash

if [ $# -eq 0 ] ; then
  echo "Usage: 
  ./pullpush.sh 'the commit message'"
  exit
fi

# generate documentation
godocdown A > A/README.md
godocdown B > B/README.md
godocdown C > C/README.md
godocdown F > F/README.md
godocdown I > I/README.md
godocdown L > L/README.md
godocdown M > M/README.md
godocdown S > S/README.md
godocdown T > T/README.md
godocdown X > X/README.md

# format indentation
go fmt ./...
echo "codes formatted.."

# testing if has error
go build loader.go || ( echo 'has error, exiting..' ; kill 0 )
echo "codes tested.."

# testing if has "gokil" included
ag gokil **/*.go && ( echo 'echo should not import previous gokil library..' ; kill 0 )
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

git commit -m "$*" && git pull && git push origin master

git tag -a `ruby -e 't = Time.now; print "v1.#{t.month+(t.year-2021)*12}%02d.#{t.hour}%02d" % [t.day, t.min]'` -m "$*"
git push --tags 

# delete tag: 
# git tag -d v1.mdd.hhmm 
# git push -d origin v1.mdd.hhmm
