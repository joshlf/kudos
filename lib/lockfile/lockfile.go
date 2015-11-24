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

const (
	lockfilename = "lock"
)

var (
	ErrNeedAbsPath = errors.New("lockfile needs absolute directory path")
)

// Lock is a handle to a file-based lock. The lock
// requires a directory in which to operate, which
// must exist ahead of time, and must be empty.
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

// New creates a new Lock with the given directory, which
// must be an absolute path; if it is not, New will return
// ErrNeedAbsPath. New only initializes the Lock
// datastructure; no filesystem operations are performed
// until a call to TryLock.
func New(dir string) (*Lock, error) {
	if !filepath.IsAbs(dir) {
		return nil, ErrNeedAbsPath
	}
	return &Lock{
		file: filepath.Join(dir, lockfilename),
		init: true,
	}, nil
}

// TryLock attempts to acquire the lock.
// It will panic if the lock is already acquired.
func (l *Lock) TryLock() (ok bool, err error) {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.init {
		panic("lockfile: uninitialized lock")
	}
	if l.locked {
		panic("lockfile: tried to lock acquired lock")
	}
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

// TryLockN is like TryLock, except that it
// will try up to n times to acquire the lock,
// sleeping for the given delay in between
// each attempt.
func (l *Lock) TryLockN(n int, delay time.Duration) (ok bool, err error) {
	for i := 0; i < n; i++ {
		ok, err = l.TryLock()
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
