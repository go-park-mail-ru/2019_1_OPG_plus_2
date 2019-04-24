package main

import "2019_1_OPG_plus_2/internal/app/authserver"

func main() {
	err := authserver.Start()
	if err != nil {
		panic(err)
	}
}
