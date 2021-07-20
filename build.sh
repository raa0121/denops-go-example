#!/bin/sh

which go >/dev/null 2>&1

if [ $? == 1 ];
then
  echo 'go is not installed.'
  exit 1
fi

go env -w GOOS=js GOARCH=wasm
(cd denops/go; go build -o main.wasm .)
go env -u GOOS GOARCH
