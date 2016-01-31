package config

import (
	"path/filepath"

	"github.com/joshlf/kudos/lib/build"
	"github.com/joshlf/kudos/lib/perm"
)

var (
	DefaultGlobalConfigFile = build.Root + "/etc/kudos/config"

	KudosDirName          = ".kudos"
	KudosDirPerms         = perm.Parse("rwxrwxr-x")
	CourseConfigFileName  = "config"
	CourseConfigFilePerms = perm.Parse("rw-rw-r--")
	HandinDirName         = "handin"
	HandinDirPerms        = perm.Parse("rwxrwxr-x")
	SavedHandinsDirName   = "saved_handins"
	SavedHandinsDirPerms  = perm.Parse("rwxrwx---")
	HandinFileName        = "handin.tgz"
	AssignmentDirName     = "assignments"
	AssignmentDirPerms    = perm.Parse("rwxrwx---")
	HooksDirName          = "hooks"
	HooksDirPerms         = perm.Parse("rwxrwxr-x")

	UserConfigFileName    = ".kudosconfig"
	UserConfigFilePerms   = perm.Parse("rw-r--r--")
	UserBlacklistFileName = ".kudosblacklist"
	// no perms specified for blacklist because
	// this is handled by custom logic in the
	// blacklist command

	PreHandinHookFileName  = "pre-handin"
	PreHandinHookFilePerms = perm.Parse("rw-rw-r--")

	DBDirName      = "db"
	DBDirPerms     = perm.Parse("rwxrwx---")
	PubDBDirName   = "pubdb"
	PubDBDirPerms  = perm.Parse("rwxrwxr-x")
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
