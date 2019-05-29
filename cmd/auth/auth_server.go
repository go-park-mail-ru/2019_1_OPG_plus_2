package main

import (
	"2019_1_OPG_plus_2/internal/app/authserver"
	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/config"
)

func main() {
	config.Init()
	auth.Init()
	err := authserver.Start()
	if err != nil {
		panic(err)
	}
}
