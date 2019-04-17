package db

import (
	"log"
)

func testsInitial() {
	// Базы для тестов
	AuthDbName = "colors_auth_test"
	CoreDbName = "colors_core_test"

	if err := Open(); err != nil && err != AlreadyInit {
		log.Fatal(err.Error())
	}
}
