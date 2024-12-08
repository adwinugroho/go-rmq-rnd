package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func LoadPath(path string) error {
	fmt.Println("CONFIG LOADED")
	/////////////////////////////
	// INITIALIZED CONFIG FILE //
	/////////////////////////////
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(".")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return err
	}

	// READ FROM ENV
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(&Config)
	if err != nil {
		return err
	}

	return nil
}

func Load() error {
	return LoadPath("config")
}
