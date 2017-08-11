#!/bin/bash

git pull

go get -v

go build -v

sudo setcap CAP_NET_BIND_SERVICE=+eip ./henry.sites

./henry.sites
