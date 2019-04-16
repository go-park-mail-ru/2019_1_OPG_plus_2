package db

import (
	"log"
)

func testsInitial() {
	// Базы для тестов
	authDbName = "colors_auth_test"
	coreDbName = "colors_core_test"

	if err := Open(); err != nil {
		log.Fatal(err.Error())
	}
}

