package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"os"
)

var CONFIG = viper.New()

var (
	Auth    AuthConfig
	Db      DbConfig
	VkOAuth OAuthConfig
	Logger  LoggerConfig
)

func init() {
	CONFIG.SetConfigName("config")             // name of config file (without extension)
	CONFIG.AddConfigPath("/etc/colors-game/")  // path to look for the config file in
	CONFIG.AddConfigPath("$HOME/.colors-game") // call multiple times to add many search paths
	CONFIG.AddConfigPath(".")                  // optionally look for config in the working directory
	err := CONFIG.ReadInConfig()               // Find and read the config file
	if err != nil {                            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	parseAuthConfig()

	parseVkConfig()

	parseDbConfig()

	parseLoggerConfig()
}

func parseLoggerConfig() {
	//logConf := reflect
}

func parseDbConfig() {
	Db.AuthDbName = CONFIG.GetString("db.auth_db_name")
	Db.AuthUsersTable = CONFIG.GetString("db.auth_users_table")
	Db.CoreDbName = CONFIG.GetString("db.core_db_name")
	Db.CoreUsersTable = CONFIG.GetString("db.core_users_table")

	var keyPrefix string
	switch os.Getenv("COLORS_SERVICE_USE_MODE") {
	case "PRODUCTION":
		keyPrefix = "db.envs.production"
	case "IN_DOCKER":
		keyPrefix = "db.envs.in_docker"
	case "USE_DOCKER_DB":
		keyPrefix = "db.envs.use_docker_db"
	default:
		keyPrefix = "db.envs.default"
	}

	conf := CONFIG.GetStringMapString(keyPrefix)
	if conf["env"] == "true" {
		Db.Host = os.Getenv(conf["host"])
		Db.Port = os.Getenv(conf["port"])
		Db.Username = os.Getenv(conf["username"])
		Db.Password = os.Getenv(conf["password"])
	} else {
		Db.Host = conf["host"]
		Db.Port = conf["port"]
		Db.Username = conf["username"]
		Db.Password = conf["password"]
	}

	fmt.Println(conf)
}

func parseVkConfig() {
	// iterate through fields
	VkOAuth.AppId = CONFIG.GetString("oauth.vk.app_id")
	VkOAuth.AppKey = CONFIG.GetString("oauth.vk.app_key")
	VkOAuth.AppSecret = CONFIG.GetString("oauth.vk.app_secret")
}

func parseAuthConfig() {
	Auth.Secret = CONFIG.GetString("auth.secret")
}

type OAuthConfig struct {
	AppId     string `json:"app_id, string"`
	AppKey    string `json:"app_key, string"`
	AppSecret string `json:"app_secret, string"`
}

type DbConfig struct {
	AuthDbName     string `json:"auth_db_name"`
	AuthUsersTable string `json:"auth_users_table"`
	CoreDbName     string `json:"core_db_name"`
	CoreUsersTable string `json:"core_users_table"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

type AuthConfig struct {
	Secret string `json:"secret"`
}

type LoggerConfig struct {
	Trace  io.Writer
	Info   io.Writer
	Warn   io.Writer
	Err    io.Writer
	Access io.Writer
}
