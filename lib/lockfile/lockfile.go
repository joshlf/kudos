// Package lockfile implements file-based locking
// that works both between processes on a single
// host, and also between processes sharing an NFS
// file system.
//
// This package provides two lock implementations -
// Lock, and LegacyLock. Most users will want Lock.
// However, the algorithm used by Lock is incorrect
// if used on NFS pre-version 3, or if any of the
// machines attempting to acquire the lock are running
// a Linux kernel pre-2.6. If either of these is
// the case, LegacyLock should be used.
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
	// ErrCollision is returned if the randomly generated
	// handle-specific lockfile name used in the LegacyLock
	// algorithm collides with that of another lock handle.
	// It is very likely that if this error is encountered,
	// it indicates an issue with the cryptographic randomness
	// available to this process. ErrCollision can only be
	// returned from LegacyLock's TryLock and TryLockN methods.
	ErrCollision   = errors.New("lockfile name collision")
	ErrNeedAbsPath = errors.New("lockfile needs absolute directory path")
)

// Lock is a handle to a file-based lock. The lock
// requires a directory in which to operate, which
// must exist ahead of time, and must be empty.
//
// A process can use multiple Locks simultaneously,
// but they will behave completely independently,
// and only one can be locked at a time. Locks are
// safe for concurrent access. The zero value of
// Lock is not valid, and any methods called on a
// zero value Lock will panic. To acquire a valid
// Lock, use NewLock.
type Lock struct {
	file   string
	locked bool
	init   bool
	m      sync.Mutex
}

// NewLock creates a new Lock with the given directory,
// which must be an absolute path; if it is not, NewLock
// will return ErrNeedAbsPath. NewLock only initializes
// the lock datastructure; no filesystem operations are
// performed until a call to TryLock or TryLockN.
func NewLock(dir string) (*Lock, error) {
	if !filepath.IsAbs(dir) {
		return nil, ErrNeedAbsPath
	}
	return &Lock{
		file: filepath.Join(dir, lockfilename),
	}, nil
}

// TryLock is equivalent to TryLockN(1, 0).
func (l *Lock) TryLock() (ok bool, err error) {
	return l.TryLockN(1, 0)
}

// TryLockN attempts to acquire the lock up to
// n times before giving up, pausing for the given
// delay between each attempt. TryLockN will panic
// if l is not initialized or if the lock is already
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
		if ok {
			l.locked = true
			return
		}
		if err != nil {
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
	return true, nil
}

// Unlock releases the lock so that others can
// acquire it. Unlock will panic if l is not
// initialized, or if the lock is not acquired.
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
