package main

import (
	"2019_1_OPG_plus_2/internal/app/chatserver"
)

func main() {
	serv, err := chatserver.Start()
	if err != nil {
		serv.Log.LogFatal("CHAT: %v", err)
	}
}
