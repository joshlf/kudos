package kudos

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joshlf/kudos/lib/config"
)

var (
	ErrNeedAbsPath = errors.New("need absolute path")
)

type Course struct {
	Code        string
	Name        string
	Description string
	TAGroup     string
}

// NOTE: All of the convenience methods to retrieve
// fields of the various parseable* types will either:
//   - check to see if the field is set before dereferencing
//     the pointer if the field is optional
//   - assume that the field has been set and dereference
//     the pointer if the field is mandatory
//
// These methods shouldn't be called except for during
// validation (in a manner that makes sure this is safe)
// or after validation (at which point these invariants
// are guaranteed to hold)

type parseableCourse struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	TAGroup     *string `json:"ta_group"`
}

func (p *parseableCourse) code() string { return *p.Code }

func (p *parseableCourse) name() (s string) {
	if p.Name != nil {
		s = *p.Name
	}
	return
}

func (p *parseableCourse) description() (s string) {
	if p.Description != nil {
		s = *p.Description
	}
	return
}

func (p *parseableCourse) taGroup() string { return *p.TAGroup }

// ParseCourseFileValidateRoot is like ParseCourseFile
// except that it infers the location of the course
// config file from the course root's path, and validates
// that the course code matches the course root.
func ParseCourseFileValidateRoot(courseRoot string) (*Course, error) {
	path := filepath.Join(courseRoot, config.KudosDirName, config.CourseConfigFileName)
	c, err := ParseCourseFile(path)
	if err != nil {
		return nil, err
	}
	base := filepath.Base(courseRoot)
	if c.Code != base {
		return nil, fmt.Errorf("course root name does not match course code (%v)", c.Code)
	}
	return c, nil
}

func ParseCourseFile(path string) (*Course, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := parseCourse(f)
	if err != nil {
		return nil, fmt.Errorf("could not parse: %v", err)
	}
	return c, nil
}

func parseCourse(r io.Reader) (*Course, error) {
	d := json.NewDecoder(r)
	var course parseableCourse
	err := d.Decode(&course)
	if err != nil {
		return nil, err
	}
	if err = validateCourse(course); err != nil {
		return nil, err
	}

	return &Course{
		Code:        course.code(),
		Name:        course.name(),
		Description: course.description(),
		TAGroup:     course.taGroup(),
	}, nil
}

func validateCourse(course parseableCourse) error {
	if course.Code == nil {
		return fmt.Errorf("must have code")
	}
	if err := ValidateCode(*course.Code); err != nil {
		return fmt.Errorf("bad course code %q: %v", *course.Code, err)
	}
	if course.TAGroup == nil {
		return fmt.Errorf("must have TA group")
	}
	// TODO(joshlf): Look up TA group (verify that it exists)
	return nil
}
