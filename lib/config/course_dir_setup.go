package config

import (
	"os"
	"path/filepath"
)

const KudosDirName = ".kudos"

func dirMaker(dir string, err *error) func() {
	*err = os.Mkdir(dir, os.ModeDir)
	return func() {
		if *err != nil {
			os.Remove(dir)
		}
	}
}

func SetupDir(course, coursePath string) (err error) {
	path := func(s ...string) string {
		return filepath.Join(append([]string{coursePath}, s...)...)
	}

	os.Chdir(coursePath)

	defConfig := DefaultCourseConfig()
	defConfig.Name = course

	defer dirMaker(path(KudosDirName), &err)()
	if err != nil {
		return
	}

	defer dirMaker(path(KudosDirName, "assignments"), &err)()
	if err != nil {
		return
	}

	defer dirMaker(path(KudosDirName, "handin"), &err)()
	if err != nil {
		return
	}

	var file *os.File
	if file, err = os.OpenFile(path(KudosDirName, "config.toml"),
		os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664); err != nil {
		return err
	}

	defConfig.WriteTOML(file)

	return nil
}
