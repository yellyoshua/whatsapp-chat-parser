#!/bin/sh

go mod tidy && go run /app/cmd/cmd.go
# whatsapp-chat-parser

# docker run --rm --name whatsapp-2 -it -v $(pwd):/app golang:1.16 /bin/bash

# docker build -t whatsapp-chat-parser:v1 . && docker rm -f whatsapp && docker run --name whatsapp --rm --env-file ./.env.production whatsapp-chat-parser:v1