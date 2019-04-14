package db

import "os"

// Docker DB params
var host = ""
var port = "12345"
var username = "root"
var password = "root"

var authDbName = "colors_auth"
var authUsersTable = "users"

var coreDbName = "colors_core"
var coreUsersTable = "users"

func init() {
	// На реальном сервере будет переменная окружения PRODUCTION=on. Меняем параметры на параметры реальной DB.
	if os.Getenv("PRODUCTION") == "on" {
		host = ""
		port = "3306"
		username = "colors"
		password = "A$55Ea~r~|lGvaZ~"
	} else if os.Getenv("TEST_DB") == "on" {
		host = "82.146.59.94"
		port = "12345"
		username = "root"
		password = "root"
	}
	// При сборке бэка в докере будет переменная окружения IN_DOCKER=on.
	// Контейнер бэка работает в одной с базой виртуальной сети под названием opg-net.
	// Эту сеть необходимо предварительно создать командой `docker network create opg-net`.
	if os.Getenv("IN_DOCKER") == "on" {
		host = "colors-db"
		port = "3306"
	}
}
