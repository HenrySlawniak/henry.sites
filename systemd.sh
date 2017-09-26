#!/bin/bash

git pull &&\
go get -v &&\
bash -c "go build -ldflags '-w -X main.buildTime=$(date +'%b-%d-%Y-%H:%M:%S') -X main.commit=$(git describe --always --dirty=*)' -v -pkgdir ~/go ." &&\
chmod a+x /var/go/src/github.com/HenrySlawniak/henry.sites/henry.sites &&\
authbind --deep /var/go/src/github.com/HenrySlawniak/henry.sites/henry.sites
