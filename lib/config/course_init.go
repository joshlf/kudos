package config

import (
	"os"
	"path/filepath"
)

func dirMaker(dir string, err *error) func() {
	*err = os.Mkdir(dir, os.ModeDir)
	return func() {
		if *err != nil {
			os.Remove(dir)
		}
	}
}

func InitCourse(course, coursePath string) (err error) {
	path := func(s ...string) string {
		return filepath.Join(append([]string{coursePath}, s...)...)
	}

	os.Chdir(coursePath)

	defConfig := DefaultCourseConfig()
	defConfig.Name = course

	defer dirMaker(path(CourseConfigDirName), &err)()
	if err != nil {
		return
	}

	defer dirMaker(path(CourseConfigDirName, "assignments"), &err)()
	if err != nil {
		return
	}

	defer dirMaker(path(CourseConfigDirName, "handin"), &err)()
	if err != nil {
		return
	}

	var file *os.File
	if file, err = os.OpenFile(path(CourseConfigDirName, "config.toml"),
		os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664); err != nil {
		return err
	}

	defConfig.WriteTOML(file)

	return nil
}
