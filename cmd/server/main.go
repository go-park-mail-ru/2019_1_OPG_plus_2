package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/app/server"
	"os"
)

func main() {
	params := server.Params{Port: os.Getenv("PORT")}
	if params.Port == "" {
		params.Port = "8001"
	}

	err := server.StartApp(params)
	if err != nil {
		fmt.Println()
	}
}
