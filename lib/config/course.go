package config

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	CourseConfigDirName  = ".kudos"
	CourseConfigFileName = "config.toml"
)

type HandinDir string

func (h *HandinDir) UnmarshalTOML(i interface{}) error {
	path, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected string value")
	}
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		return fmt.Errorf("must be relative path")
	}
	*h = HandinDir(path)
	return nil
}

type CourseConfig struct {
	Name             string    `toml:"course_name"`
	TaGroup          string    `toml:"ta_group"`
	StudentGroup     string    `toml:"student_group"`
	HandinDir        HandinDir `toml:"handin_dir"`
	ShortDescription string    `toml:"short_description"`
	LongDescription  string    `toml:"long_description"`
}

func (c CourseConfig) WriteTOML(w io.Writer) (err error) {
	return toml.NewEncoder(w).Encode(&c)
}

func DefaultCourseConfig() CourseConfig {
	return CourseConfig{
		Name:             "cs101",
		TaGroup:          "cs101tas",
		StudentGroup:     "cs101students",
		HandinDir:        HandinDir("handin"),
		ShortDescription: "CS 101",
		LongDescription:  "This is an introductory course in CS.",
	}
}

// func ReadCourseConfig(course, coursePath string) (CourseConfig, error) {
// 	confPath := filepath.Join(coursePath, CourseConfigDirName, CourseConfigFileName)
// 	f, err := os.Open(confPath)
// 	if err != nil {
// 		return CourseConfig{}, fmt.Errorf("could not open course config: %v", err)
// 	}
// 	defer f.Close()
// }
