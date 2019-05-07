#!/usr/bin/env bash

go build -o ../cmd/server/main ../cmd/server/main.go;
go build -o ../cmd/auth/main ../cmd/auth/auth_server.go;
go build -o ../cmd/chat/main ../cmd/chat/chat_server.go;
