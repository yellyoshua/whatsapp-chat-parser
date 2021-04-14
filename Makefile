#!bin/bash

#
# Makefile
# yellyoshua, 2021-04-13 23:06
#

install:
	go get .

test:
	go test -timeout 30s github.com/yellyoshua/whatsapp-chat-parser/whatsapp
	go test -timeout 30s github.com/yellyoshua/whatsapp-chat-parser/utils
	go test -timeout 30s github.com/yellyoshua/whatsapp-chat-parser/paper

clean-dependencies:
	go mod tidy

build:
	go build .

cross-build:
	GOOS=windows GOARCH=386 go build -o whatsapp-chat-parser-win32
	GOOS=linux GOARCH=386 go build -o whatsapp-chat-parser-linux-386
	GOOS=linux GOARCH=amd64 go build -o whatsapp-chat-parser-linux-amd64