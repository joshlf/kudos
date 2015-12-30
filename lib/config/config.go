package config

import (
	"path/filepath"

	"github.com/joshlf/kudos/lib/build"
)

const (
	DefaultGlobalConfigFile = build.Root + "/etc/kudos/config"

	KudosDirName         = ".kudos"
	CourseConfigFileName = "config"
	HandinDirName        = "handin"
	HandinFileName       = "handin.tgz"
	AssignmentDirName    = "assignments"
	HooksDirName         = "hooks"

	PreHandinHookFileName = "pre-handin"

	DBDirName      = "db"
	DBFileName     = "db"
	DBTempFileName = "db.tmp"
	DBLockFileName = "lock"
)

func IgnoreFile(path string) bool {
	base := filepath.Base(path)
	switch {
	case stringSuffix(base, 4) == ".tmp":
		return true
	case stringSuffix(base, 7) == ".sample":
		return true
	default:
		return false
	}
}

func IgnoreFileAndLog(printf func(format string, a ...interface{}), path string) bool {
	ignore := IgnoreFile(path)
	if ignore {
		printf("skipping %v\n", path)
	}
	return ignore
}

func stringSuffix(s string, n int) string {
	if len(s) > n {
		return s[len(s)-n:]
	}
	return s
}
