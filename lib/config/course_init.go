package config

import (
	"fmt"
	"os"
	"path/filepath"

	logpkg "github.com/synful/kudos/lib/log"
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
func dirMaker(dir string, err *error) func() {
	*err = os.Mkdir(dir, os.ModeDir)
	return func() {
		if *err != nil {
			os.Remove(dir)
		}
	}
}

// InitCourse initializes course in the directory
// coursePath. It panics if coursePath is not an
// absolute path.
//
// If log is true, InitCourse will print logging
// information about what it is doing at the Info
// level.
func InitCourse(course, coursePath string, log bool) (err error) {
	if !filepath.IsAbs(coursePath) {
		panic("internal error: non-absolute coursePath")
	}

	printf := logpkg.Info.Printf
	if !log {
		printf = func(string, ...interface{}) {}
		fmt.Print(printf)
	}

	path := func(s ...string) string {
		return filepath.Join(append([]string{coursePath}, s...)...)
	}

	conf := DefaultCourseConfig()
	conf.Name = course
	conf.TaGroup = course + "ta"
	conf.StudentGroup = course + "student"

	printf("creating %v\n", path(CourseConfigDirName))
	defer dirMaker(path(CourseConfigDirName), &err)()
	if err != nil {
		return
	}

	printf("creating %v\n", path(CourseConfigDirName, "assignments"))
	defer dirMaker(path(CourseConfigDirName, "assignments"), &err)()
	if err != nil {
		return
	}

	printf("creating %v\n", path(string(conf.HandinDir)))
	defer dirMaker(path(string(conf.HandinDir)), &err)()
	if err != nil {
		return
	}

	printf("creating %v\n", path(CourseConfigDirName, CourseConfigFileName))
	var file *os.File
	if file, err = os.OpenFile(path(CourseConfigDirName, CourseConfigFileName),
		os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664); err != nil {
		return
	}

	err = conf.WriteTOML(file)

	return nil
}
