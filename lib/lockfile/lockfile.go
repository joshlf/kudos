// Package lockfile implements file-based locking
// that works both between processes on a single
// host, and also between processes sharing an NFS
// file system.
//
// Note that this package relies on guarantees that
// are not provided on NFS pre-version 3 or on Linux
// pre-version 2.6, so it is not safe for on either
// of these systems.
package lockfile

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrNeedAbsPath = errors.New("need absolute path")
)

// Lock is a handle to a file-based lock. The lock
// requires a path to a file (whose existence will
// signal that the lock is locked). In order for
// a Lock to work properly, the lock file should
// not exist ahead of time, unless created by another
// instance of Lock.
//
// A process can use multiple Locks simultaneously,
// but they will behave completely independently,
// and only one can be locked at a time. Locks are safe
// for concurrent access.
type Lock struct {
	file   string
	locked bool
	init   bool
	m      sync.Mutex
}

// New creates a new Lock with the given file, which
// must be an absolute path; if it is not, New will
// return ErrNeedAbsPath. New only initializes the Lock
// datastructure; no filesystem operations are performed
// until a call to TryLock.
func New(path string) (*Lock, error) {
	if !filepath.IsAbs(path) {
		return nil, ErrNeedAbsPath
	}
	return &Lock{
		file: path,
		init: true,
	}, nil
}

// TryLock is equivalent to TryLockN(1, 0).
func (l *Lock) TryLock() (ok bool, err error) {
	return l.TryLockN(1, 0)

}

// TryLockN attempts to acquire the lock up to
// n times, sleeping for the given delay in between
// each attempt. It will panic if the lock is already
// acquired.
func (l *Lock) TryLockN(n int, delay time.Duration) (ok bool, err error) {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.init {
		panic("lockfile: uninitialized lock")
	}
	if l.locked {
		panic("lockfile: tried to lock acquired lock")
	}
	for i := 0; i < n; i++ {
		ok, err = l.tryLock()
		if ok || err != nil {
			return
		}
		// Only sleep if we have tries left
		if i+1 < n {
			time.Sleep(delay)
		}
	}
	return false, nil
}

func (l *Lock) tryLock() (bool, error) {
	f, err := os.OpenFile(l.file, os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return false, nil
		}
		return false, err
	}
	f.Close()
	l.locked = true
	return true, nil
}

// Unlock releases the lock. It will panic
// if the lock is not acquired.
func (l *Lock) Unlock() error {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.init {
		panic("lockfile: uninitialized lock")
	}
	if !l.locked {
		panic("lockfile: tried to unlock unacquired lock")
	}
	err := os.Remove(l.file)
	if err != nil {
		return err
	}
	l.locked = false
	return nil
}
