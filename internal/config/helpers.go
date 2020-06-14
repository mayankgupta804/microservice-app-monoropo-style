package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func getStringOrPanic(key string) string {
	if !viper.IsSet(key) {
		panic(fmt.Sprintf("%s config is not set", key))
	}

	return viper.GetString(key)
}

func getBoolOrPanic(key string) bool {
	if !viper.IsSet(key) {
		panic(fmt.Sprintf("%s config is not set", key))
	}

	return viper.GetBool(key)
}

func getIntOrPanic(key string) int {
	if !viper.IsSet(key) {
		panic(fmt.Sprintf("%s config is not set", key))
	}

	return viper.GetInt(key)
}
