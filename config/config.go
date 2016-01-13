package config

import (
	"github.com/spf13/viper"
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
}
