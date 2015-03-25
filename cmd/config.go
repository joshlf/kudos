package main

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultConfigFile = Root + "/etc/kudos/config.toml"
	EnvPrefix         = "KUDOS_"
	CourseEnvVar      = "COURSE"
)

type Config struct {
	Course, CoursePath string
}

// InitConfig initializes the configuration
// with the given configuration file path
func InitConfig(config string) {
	Debug.Printf("using config: %v\n", config)

	viper.SetConfigFile(config)
	viper.SetEnvPrefix(EnvPrefix)
	// Implementation of BindEnv only returns an error
	// on usage issues (in particular, when 0 args are
	// passed); this call should never return an error
	if err := viper.BindEnv(CourseEnvVar); err != nil {
		panic(fmt.Errorf("internal error: %v", err))
	}

	if cmdMain.PersistentFlags().Lookup("course").Changed {
		viper.Set(CourseEnvVar, courseFlag)
	}

	err := viper.ReadInConfig()
	if err != nil {
		Error.Printf("could not read config: %v\n", err)
		devFail()
	}

	if err != nil {
		Error.Printf("could not parse config: %v\n", err)
		devFail()
	}
}
