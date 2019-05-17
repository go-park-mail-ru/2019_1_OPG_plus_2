package main

import "2019_1_OPG_plus_2/internal/app/cookieserver"

func main() {
	err := cookieserver.Start()
	if err != nil {
		panic(err)
	}
}
