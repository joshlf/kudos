package main

import (
	"github.com/spf13/viper"
)

const (
	DefaultConfigFile = Root + "/etc/kudos/config.toml"
)

// InitConfig initializes the configuration
// with the given configuration file path
func InitConfig(config string) {
	Debug.Printf("using config: %v\n", config)
	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		Error.Printf("could not read config: %v\n", err)
	}
}
