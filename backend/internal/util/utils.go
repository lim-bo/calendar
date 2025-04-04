package util

import "github.com/spf13/viper"

func LoadConfig() error {
	viper.AddConfigPath("../../config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
