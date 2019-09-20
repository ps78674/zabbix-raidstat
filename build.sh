#!/usr/bin/bash

[[ ! -d build ]] && mkdir build
go build -ldflags="-s -w" -buildmode=plugin -o build/adaptec.so plugins/adaptec/main.go
go build -ldflags="-s -w" -buildmode=plugin -o build/hp.so plugins/hp/main.go
go build -ldflags="-s -w" -o build/raidstat main.go
build/raidstat $@
