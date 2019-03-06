package main

import (
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/app/server"
	"os"
)

// @title Colors service API by OPG+2
// @version 1.0
// @description Game based on filling field with color cells
// @license.name Apache 2.0

// @contact.name @DanikNik
// @contact.email nikolsky.dan@gmail.com

// @contact.name @severush

// @contact.name @Xatabch

// @host localhost:8001
// @BasePath /api

func main() {
	params := server.Params{Port: os.Getenv("PORT")}
	if params.Port == "" {
		params.Port = "8001"
	}

	err := server.StartApp(params)
	if err != nil {
		panic(err)
	}
}
