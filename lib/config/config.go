package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/synful/kudos/lib/build"
	"github.com/synful/kudos/lib/log"
)

const (
	DefaultGlobalConfigFile = build.Root + "/etc/kudos/config.toml"
	EnvPrefix               = "KUDOS_"
	CourseEnvVar            = "COURSE"
)

// Config represents the configuration
// of this instance of kudos after all
// configurations (files, flags,
// environment variables, etc) have
// been taken into account.
var Config = struct {
	// Whether the course was set (either in an
	// environment variable or on the command line)
	CourseSet          bool
	Course, CoursePath string
}{}

// The layout of the global config file
type globalConfigFile struct {
	Course_path_prefix string
	Course_path_suffix string
}

// InitConfig initializes the configuration
// by reading the global config file and
// parsing command-line arguments and
// environment variables. If config or course
// are not nil and have been set on the
// command line, their values will be used.
//
// InitConfig will panic if course is not
// nil and course.Name differs from
// CourseEnvVar (case insensitive).
func InitConfig(config *pflag.Flag, course *pflag.Flag) error {
	if course != nil && strings.ToLower(course.Name) != strings.ToLower(CourseEnvVar) {
		panic("config: course.Name differs from CourseEnvVar")
	}

	if config != nil && config.Changed {
		viper.SetConfigFile(config.Value.String())
		log.Debug.Printf("global config set on command line: %v\n", config.Value.String())
	} else {
		viper.SetConfigFile(DefaultGlobalConfigFile)
		log.Debug.Printf("using default global config: %v\n", DefaultGlobalConfigFile)
	}
	// viper.BindPFlag will make this true regardless
	// of whether the flag was specified by the user
	courseSet := viper.IsSet(CourseEnvVar) || course.Changed
	if course != nil {
		viper.BindPFlag(course.Name, course)
		log.Debug.Printf("course set on command line: %v\n", course.Value)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("could not read config: %v", err)
	}

	var cfg globalConfigFile
	err = viper.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("could not parse config: %v", err)
	}

	if courseSet {
		var code code
		if err := code.UnmarshalTOML(viper.GetString(CourseEnvVar)); err != nil {
			return fmt.Errorf("course code: %v", err)
		}
		Config.CourseSet = true
		Config.Course = viper.GetString(CourseEnvVar)
		Config.CoursePath = filepath.Join(cfg.Course_path_prefix, Config.Course, cfg.Course_path_suffix)
	}
	return nil
}

func init() {
	viper.SetEnvPrefix(EnvPrefix)
	// Implementation of BindEnv only returns an error
	// on usage issues (in particular, when 0 args are
	// passed); this call should never return an error
	if err := viper.BindEnv(CourseEnvVar); err != nil {
		panic(fmt.Errorf("config: internal error: %v", err))
	}
}
