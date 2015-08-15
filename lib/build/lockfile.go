package build

import (
	"github.com/synful/kudos/lib/lockfile"
)

func NewLockfile(dir string) (lockfile.Lock, error) {
	if LockfileLegacy {
		return lockfile.NewLegacy(dir)
	}
	return lockfile.New(dir)
}
