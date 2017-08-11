#!/bin/bash

git pull

go get -v

go build -v

bash -c "go build -ldflags '-w -X main.buildTime=$(date +'%b-%d-%Y-%H:%M:%S') -X main.commit=$(git describe --always --dirty=*)' -v ."

chmod a+x /var/henry.sites/henry.sites

authbind --deep /var/henry.sites/henry.sites
