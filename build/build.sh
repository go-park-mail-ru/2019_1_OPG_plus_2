#!/usr/bin/env bash

go build -o ${COLORSBUILD}../cmd/server/main ${COLORSBUILD}../cmd/server/main.go;
go build -o ${COLORSBUILD}../cmd/auth/main ${COLORSBUILD}../cmd/auth/auth_server.go;
go build -o ${COLORSBUILD}../cmd/chat/main ${COLORSBUILD}../cmd/chat/chat_server.go;

