package util

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func LoadConfig() error {
	viper.AddConfigPath("../../config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("cfg")
	return viper.ReadInConfig()
}

func Hash(pass string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func CheckPassword(pass string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
