#!/bin/bash

cd "$(dirname "${BASH_SOURCE:-$0}")"

GOOS=linux GOARCH=amd64 go build -o ./example/bin.linux_amd64/geflect
GOOS=windows GOARCH=amd64 go build -o ./example/bin.windows_amd64/geflect.exe
GOOS=darwin GOARCH=amd64 go build -o ./example/bin.darwin_amd64/geflect
