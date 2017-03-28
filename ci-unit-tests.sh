#!/bin/bash

set -e -u -x

export GOPATH=$(pwd)/project

go get -u github.com/kardianos/govendor

cd project/src/github.com/ONSdigital/dp-apipoc-client/

go test -v -race ./...