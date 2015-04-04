package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/synful/kudos/lib/config/internal"
	logpkg "github.com/synful/kudos/lib/log"
	"github.com/synful/kudos/lib/perm"
)

// TODO(synful): this may be broken.
// First, if called like:
// defer dirMaker()()
// it would, on an "already exists"
// error, attempt to remove the
// existant dir (clearly wrong).
//
// Second, the order in which
// the defered funcs are called
// is important (otherwise you
// might try removing /foo before
// /foo/bar, so the first would
// fail).
func dirMaker(dir string, mode os.FileMode, err *error) func() {
	*err = os.Mkdir(dir, mode)
	return func() {
		if *err != nil {
			os.Remove(dir)
		}
	}
}

// InitCourse initializes course in the directory
// coursePath. It panics if coursePath is not an
// absolute path. It returns an error if course
// is not a validly formatted course code.
//
// If log is true, InitCourse will print logging
// information about what it is doing at the
// Verbose level.
func InitCourse(course, coursePath string, log bool) (err error) {
	var code code
	if err := code.UnmarshalTOML(course); err != nil {
		return err
	}
	if !filepath.IsAbs(coursePath) {
		panic("config: non-absolute coursePath")
	}

	printf := logpkg.Verbose.Printf
	if !log {
		printf = func(string, ...interface{}) {}
	}

	defer func() {
		p := recover()
		if p != nil {
			panic(fmt.Errorf("config: internal error: %v", p))
		}
	}()

	conf := internal.MustAsset(filepath.Join("example", CourseConfigFileName))
	assign := internal.MustAsset(filepath.Join("example", CourseAssignmentsDirName, "assignment.toml.sample"))

	c := Course{coursePath, courseConfig{Code: optionalCode{code, true}}}

	dirMode := os.ModeDir | perm.Parse("rwxrwxr-x")
	fileMode := perm.Parse("rw-rw-r--")

	printf("creating %v\n", c.ConfigDir())
	defer dirMaker(c.ConfigDir(), dirMode, &err)()
	if err != nil {
		return
	}

	printf("creating %v\n", c.AssignmentsDir())
	defer dirMaker(c.AssignmentsDir(), dirMode, &err)()
	if err != nil {
		return
	}

	printf("creating %v\n", c.HandinDir())
	defer dirMaker(c.HandinDir(), dirMode, &err)()
	if err != nil {
		return
	}

	printf("creating %v\n", c.ConfigFile())
	file, err := os.OpenFile(c.ConfigFile(), os.O_CREATE|os.O_EXCL|os.O_WRONLY, fileMode)
	if err != nil {
		return
	}

	_, err = file.Write(conf)
	if err != nil {
		return
	}

	path := filepath.Join(c.AssignmentsDir(), "assignment.toml.sample")
	printf("creating %v\n", path)
	file, err = os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fileMode)
	if err != nil {
		return err
	}

	_, err = file.Write(assign)
	if err != nil {
		return
	}

	return
}
