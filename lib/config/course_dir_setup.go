package config

import (
	"os"
	"path/filepath"
)

const KudosDirName = ".kudos"

func SetupDir(course, coursePath string) (err error) {
	path := func(s string) string {
		return filepath.Join(coursePath, s)
	}
	os.Chdir(coursePath)

	defConfig := DefaultCourseConfig()
	defConfig.Name = course

	if err = os.Mkdir(path(KudosDirName), os.ModeDir); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.Remove(KudosDirName)
		}
	}()

	return nil
}
