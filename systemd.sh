#!/bin/bash

git pull

go get -v

go build -v

go build -ldflags '-w -X main.buildTime=$(date +'%b-%d-%Y-%H:%M:%S') -X main.commit=$(git describe --always --dirty=*)' -v .

authbind --deep /var/henry.sites/henry.sites
