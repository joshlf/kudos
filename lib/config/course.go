package config

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/synful/kudos/lib/log"
)

const (
	CourseConfigDirName      = ".kudos"
	CourseConfigFileName     = "config.toml"
	CourseHandinDirName      = "handin"
	CourseAssignmentsDirName = "assignments"
)

// HandinMethod represents a method of
// implementing handing in assignments
// (either using facls or a setgid
// program).
type HandinMethod string

const (
	FaclMethod   HandinMethod = "facl"
	SetgidMethod HandinMethod = "setgid"
)

// So we don't have to export (HandinMethod).UnmarshalTOML
type handinMethod HandinMethod

func (h *handinMethod) UnmarshalTOML(i interface{}) error {
	method, _ := i.(string)
	hmethod := HandinMethod(strings.ToLower(method))
	if hmethod != FaclMethod && hmethod != SetgidMethod {
		return fmt.Errorf("allowed methods: %v, %v", FaclMethod, SetgidMethod)
	}
	*h = handinMethod(hmethod)
	return nil
}

// Course represents the configuration of a course.
type Course struct {
	path   string
	config courseConfig
}

// Code returns c's code.
func (c *Course) Code() string { return string(c.config.Code.code) }

// Name returns the human-readbale
// name of c. If one was not set
// in the config file, it defaults
// to c.Code().
func (c *Course) Name() string {
	if !c.config.Name.set {
		return string(c.config.Code.code)
	}
	return c.config.Name.string
}

// TaGroup returns the name of c's TA group.
func (c *Course) TaGroup() string { return c.config.TaGroup.string }

// StudentGroup returns the name of c's student group.
func (c *Course) StudentGroup() string { return c.config.StudentGroup.string }

// HandinMethod returns c's handin method.
func (c *Course) HandinMethod() HandinMethod { return HandinMethod(c.config.HandinMethod.handinMethod) }

// DescriptionSet returns whether c's config
// specified a description.
func (c *Course) DescriptionSet() bool { return c.config.Description.set }

// Description returns c's description, or
// the empty string if c.DescriptionSet() == false.
func (c *Course) Description() string {
	if !c.config.Description.set {
		// Could probably rely on this
		// being the zero value of
		// c.config.Description.string,
		// but this is safer.
		return ""
	}
	return c.config.Description.string
}

func (c *Course) ConfigDir() string  { return filepath.Join(c.path, CourseConfigDirName) }
func (c *Course) ConfigFile() string { return filepath.Join(c.ConfigDir(), CourseConfigFileName) }
func (c *Course) HandinDir() string  { return filepath.Join(c.ConfigDir(), CourseHandinDirName) }
func (c *Course) AssignmentsDir() string {
	return filepath.Join(c.ConfigDir(), CourseAssignmentsDirName)
}

type courseConfig struct {
	Code         optionalCode         `toml:"code"` // Guaranteed to be set
	Name         optionalString       `toml:"name"`
	TaGroup      optionalString       `toml:"ta_group"`      // Guaranteed to be set
	StudentGroup optionalString       `toml:"student_group"` // Guaranteed to be set
	HandinMethod optionalHandinMethod `toml:"handin_method"` // Guaranteed to be set
	Description  optionalString       `toml:"description"`
}

func (c courseConfig) WriteTOML(w io.Writer) (err error) {
	return toml.NewEncoder(w).Encode(&c)
}

func ReadCourseConfig(course, coursePath string) (Course, error) {
	confPath := filepath.Join(coursePath, CourseConfigDirName, CourseConfigFileName)
	log.Debug.Printf("reading course config file: %v\n", confPath)
	var conf courseConfig
	_, err := toml.DecodeFile(confPath, &conf)
	if err != nil {
		return Course{}, fmt.Errorf("could not parse course config: %v", err)
	}
	if !conf.Code.set {
		return Course{}, fmt.Errorf("course code must be set")
	}
	if course != string(conf.Code.code) {
		return Course{}, fmt.Errorf("course code in config (%v) does not match expected code (%v)", string(conf.Code.code), course)
	}
	if !conf.TaGroup.set {
		return Course{}, fmt.Errorf("ta_group must be set")
	}
	if !conf.StudentGroup.set {
		return Course{}, fmt.Errorf("student_group must be set")
	}
	if !conf.HandinMethod.set {
		return Course{}, fmt.Errorf("handin_method must be set")
	}
	return Course{coursePath, conf}, nil
}
