package db

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
)

var host = config.Db.Host
var port = config.Db.Port
var username = config.Db.Username
var password = config.Db.Password

var AuthDbName = config.Db.AuthDbName
var AuthUsersTable = config.Db.AuthUsersTable

var CoreDbName = config.Db.CoreDbName
var CoreUsersTable = config.Db.CoreUsersTable

//func init() {
//	// На реальном сервере будет переменная окружения PRODUCTION=on. Меняем параметры на параметры реальной DB.
//	if os.Getenv("PRODUCTION") == "on" {
//		host = os.Getenv("DB_HOST")
//		port = os.Getenv("DB_PORT")
//		username = os.Getenv("DB_USERNAME")
//		password = os.Getenv("DB_PASSWORD")
//	} else
//	// При сборке бэка в докере будет переменная окружения IN_DOCKER=on.
//	if os.Getenv("IN_DOCKER") == "on" {
//		host = "colors-db"
//		port = "3306"
//		username = "root"
//		password = "root"
//	} else
//	// Если сервер запускается не в докере, но хочется использовать БД из докера USE_DOCKER_DB=on.
//	if os.Getenv("USE_DOCKER_DB") == "on" {
//		host = ""
//		port = "12345"
//		username = "root"
//		password = "root"
//	}
//}
