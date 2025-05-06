package notifier

import "github.com/spf13/viper"

func LoadConfig() error {
	viper.AddConfigPath("../../config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("smtp_cfg")
	return viper.ReadInConfig()
}
