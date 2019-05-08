package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var permittedLevels = []string{
	"trace",
	"info",
	"warn",
	"err",
	"access",
	"fatal",
}

var CONFIG = viper.New()

var (
	Auth    AuthConfig
	Db      DbConfig
	VkOAuth OAuthConfig
	Logger  LoggerConfig
)

func init() {
	CONFIG.SetConfigName("config")                                   // name of testConfig file (without extension)
	CONFIG.AddConfigPath("/etc/colors-game/")                        // path to look for the testConfig file in
	CONFIG.AddConfigPath("$HOME/.colors-game")                       // call multiple times to add many search paths
	CONFIG.AddConfigPath(".")                                        // optionally look for testConfig in the working directory
	CONFIG.AddConfigPath("/home/daniknik/go/src/2019_1_OPG_plus_2/") // optionally look for testConfig in the working directory
	CONFIG.AddConfigPath(os.Getenv("COLORS_CONFIG_PATH"))
	err := CONFIG.ReadInConfig() // Find and read the testConfig file
	if err != nil {              // Handle errors reading the testConfig file
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}

	parseAuthConfig(CONFIG, &Auth)

	parseVkConfig(CONFIG, &VkOAuth)

	parseDbConfig(CONFIG, &Db)

	parseLoggerConfig(CONFIG, &Logger)
}

type OAuthConfig struct {
	AppId     string `json:"app_id, string"`
	AppKey    string `json:"app_key, string"`
	AppSecret string `json:"app_secret, string"`
}

func parseVkConfig(v *viper.Viper, conf *OAuthConfig) {
	conf.AppId = v.GetString("oauth.vk.app_id")
	conf.AppKey = v.GetString("oauth.vk.app_key")
	conf.AppSecret = v.GetString("oauth.vk.app_secret")
}

type DbConfig struct {
	AuthDbName        string `json:"auth_db_name"`
	AuthUsersTable    string `json:"auth_users_table"`
	CoreDbName        string `json:"core_db_name"`
	CoreUsersTable    string `json:"core_users_table"`
	ChatDbName        string `json:"chat_db_name"`
	ChatMessagesTable string `json:"chat_messages_table"`
	ChatTypesTable    string `json:"chat_types_table"`
	Host              string `json:"host"`
	Port              string `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	AuthTestDb        string `json:"auth_test_db"`
	CoreTestDb        string `json:"core_test_db"`
	ChatTestDb        string `json:"chat_test_db"`
}

func parseDbConfig(v *viper.Viper, conf *DbConfig) {
	conf.AuthDbName = v.GetString("db.auth_db_name")
	conf.AuthUsersTable = v.GetString("db.auth_users_table")
	conf.CoreDbName = v.GetString("db.core_db_name")
	conf.CoreUsersTable = v.GetString("db.core_users_table")
	conf.ChatDbName = v.GetString("db.chat_db_name")
	conf.ChatMessagesTable = v.GetString("db.chat_messages_table")
	conf.ChatTypesTable = v.GetString("db.chat_types_table")

	conf.AuthTestDb = v.GetString("db.auth_test_db")
	conf.CoreTestDb = v.GetString("db.core_test_db")
	conf.ChatTestDb = v.GetString("db.chat_test_db")

	keyPrefix := "db.envs." + os.Getenv("COLORS_DB")

	confMap := v.GetStringMapString(keyPrefix)

	if len(confMap) == 0 {
		confMap = v.GetStringMapString("db.envs.default")
	}
	if confMap["env"] == "true" {
		conf.Host = os.Getenv(confMap["host"])
		conf.Port = os.Getenv(confMap["port"])
		conf.Username = os.Getenv(confMap["username"])
		conf.Password = os.Getenv(confMap["password"])
	} else {
		conf.Host = confMap["host"]
		conf.Port = confMap["port"]
		conf.Username = confMap["username"]
		conf.Password = confMap["password"]
	}
}

type AuthConfig struct {
	Secret              string `json:"secret"`
	AuthServiceLocation string `json:"auth_service_location"`
	AuthPort            string `json:"auth port"`

	CookieServiceLocation string `json:"cookie_service_location"`
	CookieServicePort     string `json:"cookie_service_port"`
}

func parseAuthConfig(v *viper.Viper, conf *AuthConfig) {
	conf.Secret = v.GetString("auth.secret")

	keyPrefix := "auth.envs." + strings.ToLower(os.Getenv("COLORS_SERVICE_USE_MODE"))
	confMap := v.GetStringMapString(keyPrefix)
	if len(confMap) == 0 {
		confMap = v.GetStringMapString("auth.envs.default")
	}
	conf.AuthPort = confMap["port"]
	conf.AuthServiceLocation = confMap["service_location"]

	conf.CookieServiceLocation = confMap["cookie_service_location"]
	conf.CookieServicePort = confMap["cookie_service_port"]
}

type LoggerConfig struct {
	Levels map[string]io.Writer
	Files  []*os.File
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Levels: map[string]io.Writer{
			"trace":  os.Stdout,
			"info":   os.Stdout,
			"warn":   os.Stdout,
			"err":    os.Stdout,
			"access": os.Stdout,
			"fatal":  os.Stdout,
		},
	}
}

func parseLoggerConfig(v *viper.Viper, conf *LoggerConfig) {
	conf.Levels = NewLoggerConfig().Levels //fucking costyl' but it works only this way. holy crap...

	confMap := v.GetStringMap("logging")
	if len(confMap) == 0 {
		return
	}
	for k, v := range confMap {
		vm := v.(map[string]interface{})
		err := checkLevel(k)
		if err != nil {
			panic(err)
		}

		switch vm["mode"] {
		case "prompt":
			conf.Levels[k] = os.Stdout
		case "suppress":
			conf.Levels[k] = ioutil.Discard
		case "file":
			var filename interface{}
			if filename = vm["file"]; filename == nil {
				filename = "colors.log"
			}
			f, err := os.OpenFile(filename.(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic("error opening file: " + err.Error())
			}
			conf.Levels[k] = f
			conf.Files = append(conf.Files, f)
		default:
			conf.Levels[k] = os.Stdout
		}
	}
}

func checkLevel(l string) error {
	ok := false
	for _, pl := range permittedLevels {
		if l == pl {
			ok = true
		}
	}

	if !ok {
		return fmt.Errorf("logging level is not permited")
	}
	return nil
}
