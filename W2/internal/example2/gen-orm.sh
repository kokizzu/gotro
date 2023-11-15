#!/usr/bin/env bash


cd ./model
  cat *.go | grep '//go:generate ' | cut -d ' ' -f 2- | bash -x > /tmp/1.log

for i in ./m*; do
  if [[ ! -d "$i" ]] ; then continue ; fi
  echo $i
  pushd .
  cd "$i"

  # generate ORM
  go test -bench=.

  for j in ./*; do
    echo $j
    if [[ ! -d "$j" ]] ; then continue ; fi

    pushd .
    cd "$j"
    echo `pwd`
    cat *.go | grep '//go:generate ' | cut -d ' ' -f 2- | bash -x >> /tmp/1.log
    popd

  done

  popd

done
