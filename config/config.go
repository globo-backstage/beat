package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

func ReadConfigFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	viper.SetConfigFile(filePath)
	return viper.ReadInConfig()
}

func init() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	envReplacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(envReplacer)

	viper.SetDefault("log.level", "info")
}

func LoadLogSettings() {
	logLevel, err := logrus.ParseLevel(viper.GetString("log.level"))

	if err != nil {
		log.Fatal(err)
	}

	logrus.SetLevel(logLevel)
}
