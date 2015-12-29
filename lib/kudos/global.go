package kudos

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type GlobalConfig struct {
	CoursePathPrefix string
}

func CourseCodeToPath(code string, gc *GlobalConfig) string {
	return filepath.Join(gc.CoursePathPrefix, code)
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

type parseableGlobalConfig struct {
	CoursePathPrefix *string `json:"course_path_prefix"`
}

func (p *parseableGlobalConfig) prefix() string { return *p.CoursePathPrefix }

func ParseGlobalConfigFile(path string) (*GlobalConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	gc, err := parseGlobalConfig(f)
	if err != nil {
		return nil, err
	}
	return gc, nil
}

func parseGlobalConfig(r io.Reader) (*GlobalConfig, error) {
	d := json.NewDecoder(r)
	var gc parseableGlobalConfig
	err := d.Decode(&gc)
	if err != nil {
		return nil, err
	}
	if err = validateGlobalConfig(gc); err != nil {
		return nil, err
	}
	return &GlobalConfig{gc.prefix()}, nil
}

func validateGlobalConfig(gc parseableGlobalConfig) error {
	if gc.CoursePathPrefix == nil {
		return fmt.Errorf("must have course_path_prefix")
	}
	return nil
}
