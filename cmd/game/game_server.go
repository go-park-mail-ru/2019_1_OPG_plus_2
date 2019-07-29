package main

import (
	"2019_1_OPG_plus_2/internal/app/gameserver"
)

func main() {
	err := gameserver.Start()
	if err != nil {
		panic(err)
	}
}
