// Package lockfile implements file-based locking
// that works both between processes on a single
// host, and also between processes sharing an NFS
// file system.
//
// This package provides two lock implementations -
// a default implementation (New), and a legacy
// implementation (NewLegacy). Most users will want
// the default implementation. However, the algorithm
// used is incorrect if used on NFS pre-version 3, or
// if any of the machines attempting to acquire the
// lock are running a Linux kernel pre-2.6. If either
// of these is the case, legacy implementation should
// be used.
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
	// handle-specific lockfile name used in the legacy
	// algorithm collides with that of another lock handle.
	// It is very likely that if this error is encountered,
	// it indicates an issue with the cryptographic randomness
	// available to this process. ErrCollision can only be
	// returned from the legacy implementation's TryLock method.
	ErrCollision   = errors.New("lockfile name collision")
	ErrNeedAbsPath = errors.New("lockfile needs absolute directory path")
)

// Lock is a handle to a file-based lock. The lock
// requires a directory in which to operate, which
// must exist ahead of time, and must be empty.
//
// A process can use multiple Locks simultaneously,
// but they will behave completely independently,
// and only one can locked at a time. Locks are safe
// for concurrent access.
type Lock interface {
	// TryLock attempts to acquire the lock.
	// It will panic if the lock is already acquired.
	TryLock() (ok bool, err error)

	// Unlock releases the lock. It will panic
	// if the lock is not acquired.
	Unlock() error
}

// TryLockN will call l.Lock up to n times,
// sleeping for the given delay in between
// each call. Users should call TryLockN
// rather than implementing the functionality
// themselves, as TryLockN is able to make
// optimizations which rely on structures
// internal to this package.
func TryLockN(l Lock, n int, delay time.Duration) (ok bool, err error) {
	if ll, ok := l.(*legacyLock); ok {
		return ll.tryLockN(n, delay)
	}
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

type lock struct {
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
//
// The Lock returned by New uses an algorithm that is not
// safe for use on NFS shares running NFS versions prior
// to 3, or if any of the machines attempting to acquire
// the lock are running on Linux kernel versions prior to
// 2.6. If either of these conditions hold, NewLegacy
// should be used instead.
//
// Locks returned by New and NewLegacy are incompatible;
// using them together will result in undefined behavior.
func New(dir string) (Lock, error) {
	if !filepath.IsAbs(dir) {
		return nil, ErrNeedAbsPath
	}
	return &lock{
		file: filepath.Join(dir, lockfilename),
	}, nil
}

func (l *lock) TryLock() (ok bool, err error) {
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
	return true, nil
}

func (l *lock) Unlock() error {
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
