package main

import (
	"2019_1_OPG_plus_2/internal/app/coreserver"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"os"
)

// @title Colors service API by OPG+2
// @version 1.0
// @description Game based on filling field with color cells
// @license.name Apache 2.0

// @contact.name @DanikNik
// @contact.email nikolsky.dan@gmail.com

// @host localhost:8002
// @BasePath /api

func main() {
	params := coreserver.Params{Port: os.Getenv("PORT")}
	if params.Port == "" {
		params.Port = "8002"
	}

	err := coreserver.StartApp(params)
	if err != nil {
		coreserver.StopApp()
		tsLogger.LogFatal("%s", err)
	}

}
