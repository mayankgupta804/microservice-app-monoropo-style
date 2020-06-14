package config

import (
	"github.com/spf13/viper"
)

type Server struct {
	Port string
}

type Logger struct {
	Output string
}

type Database struct {
	Name          string
	User          string
	Password      string
	Port          string
	Host          string
	MigrationsDir string
	Dialect       string
}

type Config struct {
	Server   Server
	Database Database
	Logger   Logger
}

var App Config
var DB Database

func Load() {
	viper.SetConfigName("application")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	App = Config{
		Server: Server{
			Port: getStringOrPanic("APP_PORT"),
		},
		Database: Database{
			Name:          getStringOrPanic("DATABASE_NAME"),
			Host:          getStringOrPanic("DATABASE_HOST"),
			Port:          getStringOrPanic("DATABASE_PORT"),
			User:          getStringOrPanic("DATABASE_USER"),
			Password:      getStringOrPanic("DATABASE_PASSWORD"),
			MigrationsDir: getStringOrPanic("DATABASE_MIGRATIONS_DIR"),
			Dialect:       getStringOrPanic("DATABASE_DIALECT"),
		},
		Logger: Logger{
			Output: getStringOrPanic("LOG_OUTPUT"),
		},
	}
}
