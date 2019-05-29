package config

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

var TestData = map[string]string{
	"testVkConfig": `
{"oauth": {"vk": {"app_id": "some_id","app_key": "some_key","app_secret": "some_secret"}}}`,
	"testDBConfig": `
{
	"db": {
    	"auth_db_name": "auth_db",
    	"envs": {
      		"default": {
		        "env": false,
		        "host": "host_default",
		        "port": "port_default",
		        "username": "username_default",
		        "password": "password_default"
			},
			"IN_DOCKER_NET": {
	        	"env": false,
	        	"host": "host_in_docker_net",
	        	"port": "port_in_docker_net",
	        	"username": "username_in_docker_net",
	        	"password": "password_in_docker_net"
			},
	      	"USE_DOCKER_DB": {
	        	"env": false,
	        	"host": "host_docker_db",
	        	"port": "port_docker_db",
	        	"username": "username_docker_db",
	        	"password": "password_docker_db"
			},
	      	"PRODUCTION": {
		        "env": true,
		        "host": "DB_HOST",
		        "port": "DB_PORT",
		        "username": "DB_USERNAME",
		        "password": "DB_PASSWORD"
			}
		}
	}
}`,
	"testAuthConfig": `
{
	"auth": {
	    "secret": "some_secret",
	    "envs": {
			"IN_DOCKER_NET": {
		        "service_location": "auth_docker_net",
		        "port": "port_docker_net",
		        "cookie_service_location": "cookie_docker_net",
		        "cookie_service_port": "cookie_port_docker_net"
	      	},
	      	"default": {
	        	"service_location": "auth_default",
	        	"port": "port_default",
	        	"cookie_service_location": "cookie_default",
	        	"cookie_service_port": "cookie_port_default"
	      	}
	    }
  	}
}
`,
	"testGameConfig": `
{
	"game":{
		"envs":{
			"IN_DOCKER_NET": {
		        "service_location": "game_docker_net",
		        "port": "port_docker_net"
	      	},
	      	"default": {
	        	"service_location": "game_default",
	        	"port": "port_default"
	      	}
		}
	}
}
`,

	"testLoggerPrompts": `
{
	"logging": {
	    "trace": {
			"mode": "prompt"
		},
	    "info": {
			"mode": "prompt"
	    },
	    "warn": {
	      	"mode": "prompt"
	    },
	    "err": {
			"mode": "prompt"
	    },
	    "access": {
		    "mode": "prompt"
	    },
	    "fatal": {
	    	"mode": "prompt"
	    }
	}
}`,
	"testLoggerFiles": `
{
	"logging": {
	    "trace": {
			"mode": "file"
	   	},
	    "info": {
			"mode": "file",
	      	"file": "info.log"
	    },
	    "warn": {
	      	"mode": "file",
	      	"file": "warn.log"
	    },
	    "err": {
			"mode": "file",
			"file": "err.log"
	    },
	    "access": {
		    "mode": "file",
		    "file": "access.log"
	    },
	    "fatal": {
	    	"mode": "file",
			"file": "fatal.log"
	    }
  	}
}
`,
	"testLoggerSuppress": `
{
	"logging": {
	    "trace": {
			"mode": "suppress"
	   	},
	    "info": {
			"mode": "suppress"
	    },
	    "warn": {
	      	"mode": "suppress"
	    },
	    "err": {
			"mode": "suppress"
	    },
	    "access": {
		    "mode": "suppress"
	    },
	    "fatal": {
	    	"mode": "suppress"
	    }
  	}
}
`,
}

func Test_parseVkConfig(t *testing.T) {
	tests := []struct {
		name         string
		v            *viper.Viper
		configString string
		conf         *OAuthConfig
		sample       OAuthConfig
	}{
		{
			name:         "VkConfig",
			v:            viper.New(),
			configString: TestData["testVkConfig"],

			conf: &OAuthConfig{},
			sample: OAuthConfig{
				AppSecret: "some_secret",
				AppKey:    "some_key",
				AppId:     "some_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.SetConfigType("json")
			err := tt.v.ReadConfig(strings.NewReader(tt.configString))
			if err != nil {
				t.Fatalf("WTF: %v", err)
			}
			parseVkConfig(tt.v, tt.conf)
			if !reflect.DeepEqual(tt.sample, *tt.conf) {
				t.Errorf("Configs are not equal:\nGOT: %v\nEXP:%v", *tt.conf, tt.sample)
			}
		})

	}
}

func Test_parseDbConfig(t *testing.T) {
	oldValues := map[string]string{
		"COLORS_DB":   os.Getenv("COLORS_DB"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USERNAME": os.Getenv("DB_USERNAME"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
	}

	type args struct {
		v            *viper.Viper
		configString string

		conf   *DbConfig
		sample DbConfig

		envs map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "DbConfigNoEnv",
			args: args{
				v:            viper.New(),
				configString: TestData["testDBConfig"],

				sample: DbConfig{
					Username:   "username_default",
					Port:       "port_default",
					Password:   "password_default",
					Host:       "host_default",
					AuthDbName: "auth_db",
				},
				conf: &DbConfig{},
				envs: map[string]string{
					"COLORS_DB": "", //setting it empty to be sure it is not accidentally filled
				},
			},
		},
		{
			name: "DbConfigInDockerNet",
			args: args{
				v:            viper.New(),
				configString: TestData["testDBConfig"],

				sample: DbConfig{
					Username:   "username_in_docker_net",
					Port:       "port_in_docker_net",
					Password:   "password_in_docker_net",
					Host:       "host_in_docker_net",
					AuthDbName: "auth_db",
				},
				conf: &DbConfig{},
				envs: map[string]string{
					"COLORS_DB": "IN_DOCKER_NET",
				},
			},
		},
		{
			name: "DbConfigUseDockerDb",
			args: args{
				v:            viper.New(),
				configString: TestData["testDBConfig"],

				sample: DbConfig{
					Username:   "username_docker_db",
					Port:       "port_docker_db",
					Password:   "password_docker_db",
					Host:       "host_docker_db",
					AuthDbName: "auth_db",
				},
				conf: &DbConfig{},
				envs: map[string]string{
					"COLORS_DB": "USE_DOCKER_DB",
				},
			},
		},
		{
			name: "DbConfigProduction",
			args: args{
				v:            viper.New(),
				configString: TestData["testDBConfig"],

				sample: DbConfig{
					Username:   "username_production",
					Port:       "port_production",
					Password:   "password_production",
					Host:       "host_production",
					AuthDbName: "auth_db",
				},
				conf: &DbConfig{},
				envs: map[string]string{
					"COLORS_DB":   "PRODUCTION",
					"DB_HOST":     "host_production",
					"DB_PORT":     "port_production",
					"DB_USERNAME": "username_production",
					"DB_PASSWORD": "password_production",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.v.SetConfigType("json")
			err := tt.args.v.ReadConfig(strings.NewReader(tt.args.configString))
			if err != nil {
				t.Errorf("WTF: %v", err)
			}
			for name, val := range tt.args.envs {
				err := os.Setenv(name, val)
				if err != nil || os.Getenv(name) != val {
					t.Errorf("WTF setenv failed")
				}
			}
			parseDbConfig(tt.args.v, tt.args.conf)
			if !reflect.DeepEqual(*tt.args.conf, tt.args.sample) {
				t.Errorf("Configs are not equal:\nGOT: %v\nEXP:%v", *tt.args.conf, tt.args.sample)
			}
		})
	}

	for k, v := range oldValues {
		err := os.Setenv(k, v)
		if err != nil {
			t.Logf("Setting env back to initial values failed: %v: %v", k, v)
		}

	}
}

func Test_parseAuthConfig(t *testing.T) {
	oldValues := map[string]string{
		"COLORS_SERVICE_USE_MODE": os.Getenv("COLORS_SERVICE_USE_MODE"),
	}

	type args struct {
		v    *viper.Viper
		conf *AuthConfig

		configString string
		sample       AuthConfig
		envs         map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "AuthConfigNoEnv",
			args: args{
				v:            viper.New(),
				conf:         &AuthConfig{},
				configString: TestData["testAuthConfig"],

				sample: AuthConfig{
					Secret:                "some_secret",
					AuthServiceLocation:   "auth_default",
					AuthPort:              "port_default",
					CookieServiceLocation: "cookie_default",
					CookieServicePort:     "cookie_port_default",
				},
				envs: map[string]string{
					"COLORS_SERVICE_USE_MODE": "", //setting it empty to be sure it is not accidentally filled
				},
			},
		},
		{
			name: "AuthConfigInDockerNet",
			args: args{
				v:            viper.New(),
				conf:         &AuthConfig{},
				configString: TestData["testAuthConfig"],

				sample: AuthConfig{
					Secret:                "some_secret",
					AuthServiceLocation:   "auth_docker_net",
					AuthPort:              "port_docker_net",
					CookieServiceLocation: "cookie_docker_net",
					CookieServicePort:     "cookie_port_docker_net",
				},
				envs: map[string]string{
					"COLORS_SERVICE_USE_MODE": "IN_DOCKER_NET",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.v.SetConfigType("json")
			err := tt.args.v.ReadConfig(strings.NewReader(tt.args.configString))
			if err != nil {
				t.Fatalf("WTF: %v", err)
			}
			for name, val := range tt.args.envs {
				err := os.Setenv(name, val)
				if err != nil || os.Getenv(name) != val {
					t.Errorf("WTF setenv failed")
				}
			}
			parseAuthConfig(tt.args.v, tt.args.conf)
			if !reflect.DeepEqual(*tt.args.conf, tt.args.sample) {
				t.Errorf("Configs are not equal:\nGOT: %v\nEXP:%v", *tt.args.conf, tt.args.sample)
			}
		})
	}

	for k, v := range oldValues {
		err := os.Setenv(k, v)
		if err != nil {
			t.Logf("Setting env back to initial values failed: %v: %v", k, v)
		}

	}
}

func Test_parseGameConfig(t *testing.T) {
	oldValues := map[string]string{
		"COLORS_SERVICE_USE_MODE": os.Getenv("COLORS_SERVICE_USE_MODE"),
	}

	type args struct {
		v    *viper.Viper
		conf *GameConfig

		configString string
		sample       GameConfig
		envs         map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GameConfigNoEnv",
			args: args{
				v:            viper.New(),
				conf:         &GameConfig{},
				configString: TestData["testGameConfig"],

				sample: GameConfig{
					GameServicePort:     "port_default",
					GameServiceLocation: "game_default",
				},
				envs: map[string]string{
					"COLORS_SERVICE_USE_MODE": "", //setting it empty to be sure it is not accidentally filled
				},
			},
		},
		{
			name: "GameConfigInDockerNet",
			args: args{
				v:            viper.New(),
				conf:         &GameConfig{},
				configString: TestData["testGameConfig"],

				sample: GameConfig{
					GameServicePort:     "port_docker_net",
					GameServiceLocation: "game_docker_net",
				},
				envs: map[string]string{
					"COLORS_SERVICE_USE_MODE": "IN_DOCKER_NET",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.v.SetConfigType("json")
			err := tt.args.v.ReadConfig(strings.NewReader(tt.args.configString))
			if err != nil {
				t.Fatalf("WTF: %v", err)
			}
			for name, val := range tt.args.envs {
				err := os.Setenv(name, val)
				if err != nil || os.Getenv(name) != val {
					t.Errorf("WTF setenv failed")
				}
			}
			parseGameConfig(tt.args.v, tt.args.conf)
			if !reflect.DeepEqual(*tt.args.conf, tt.args.sample) {
				t.Errorf("Configs are not equal:\nGOT: %v\nEXP:%v", *tt.args.conf, tt.args.sample)
			}
		})
	}

	for k, v := range oldValues {
		err := os.Setenv(k, v)
		if err != nil {
			t.Logf("Setting env back to initial values failed: %v: %v", k, v)
		}

	}
}

func Test_parseLoggerConfig(t *testing.T) {
	type args struct {
		v    *viper.Viper
		conf *LoggerConfig

		configString string

		compMode  string
		sample    LoggerConfig
		filenames map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "LoggerNoConfig",
			args: args{
				v:            viper.New(),
				configString: `{"not_logging_key": "not_logging_value"}`,

				conf:     &LoggerConfig{},
				compMode: "sample",
				sample:   *NewLoggerConfig(),
			},
		},
		{
			name: "LoggerPromptConfig",
			args: args{
				v:            viper.New(),
				configString: TestData["testLoggerPrompts"],

				conf:     &LoggerConfig{},
				compMode: "sample",
				sample:   *NewLoggerConfig(),
			},
		},
		{
			name: "LoggerFileConfig",
			args: args{
				v:            viper.New(),
				configString: TestData["testLoggerFiles"],

				conf:     &LoggerConfig{},
				compMode: "file",
				filenames: map[string]string{
					"trace":  "colors.log",
					"info":   "info.log",
					"warn":   "warn.log",
					"err":    "err.log",
					"access": "access.log",
					"fatal":  "fatal.log",
				},
			},
		},
		{
			name: "LoggerSuppressConfig",
			args: args{
				v:            viper.New(),
				configString: TestData["testLoggerSuppress"],

				conf:     &LoggerConfig{},
				compMode: "sample",
				sample: LoggerConfig{
					Levels: map[string]io.Writer{
						"trace":  ioutil.Discard,
						"info":   ioutil.Discard,
						"warn":   ioutil.Discard,
						"err":    ioutil.Discard,
						"access": ioutil.Discard,
						"fatal":  ioutil.Discard,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.v.SetConfigType("json")
			err := tt.args.v.ReadConfig(strings.NewReader(tt.args.configString))
			if err != nil {
				t.Fatalf("WTF: %v", err)
			}
			parseLoggerConfig(tt.args.v, tt.args.conf)
			if tt.args.compMode == "sample" {
				if !reflect.DeepEqual(*tt.args.conf, tt.args.sample) {
					t.Errorf("Configs are not equal:\nGOT: %v\nEXP:%v", *tt.args.conf, tt.args.sample)
				}
			} else if tt.args.compMode == "file" {
				for key, val := range tt.args.conf.Levels {
					if tt.args.filenames[key] != val.(*os.File).Name() {
						t.Errorf("Filenames not equal:\nGOT: %v\nEXP: %v", tt.args.filenames[key], val.(*os.File).Name())
					}
				}
			}

			defer func() {
				for _, f := range tt.args.conf.Files {
					err := f.Close()
					if err != nil {
						t.Fatalf("WTF: %v", err)
					}
				}
			}()
		})
	}
}
