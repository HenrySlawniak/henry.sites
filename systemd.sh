#!/bin/bash

git pull

go get -v

go build -v

authbind --deep /var/henry.sites/henry.sites
