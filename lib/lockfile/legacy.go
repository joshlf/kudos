package lockfile

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

/*
	ALGORITHM: The algorithm used is a variant on that
	described in the Linux open(2) man page, under the
	O_EXCL heading. When on NFSv3 or later and kernel
	2.6 or later, O_EXCL is guaranteed to be obeyed,
	and can be used for locking - try to create the
	lockfile with O_EXCL, and it is guaranteed that only
	one such concurrent operation will succeed - that
	person has the lock. This is the algorithm used in
	the Lock implementation. However, this atomicity
	property is not guaranteed pre-NFSv3 or pre-kernel
	2.6. The man page describes an alternative - each
	process attempting to acquire the lock creates a file
	which is known ahead of time to be unique (for example,
	by using a combination of hostname and PID - this is
	the approach suggested by the man page). Then, each
	process attempts to create a hard link to that file.
	Only one such operation will succeed, and whoever's
	file has a hard link to it is the one who acquires
	the lock.

	The algorithm used here differs from the one described
	in the man page in the way that it guarantees
	uniqueness among the process-specific files. While
	a PID is trivial to acquire (simply os.Getpid()), the
	hostname is not guaranteed to be unique, so IP address
	must be used instead, which is much tricker. It must
	be made sure that the address used is the one seen by
	the NFS server, which requires figuring out which
	interface is used to route traffic to the NFS server.
	It also is helpless in the face of address translation
	occuring upstream of the host.

	For this reason, the algorithm used here does not
	attempt to construct an identifier guaranteed to be
	unique, but rather uses randomness to make the
	probability of collision negligible. Each process
	generates its own random nonce. The name of the file
	for a given lock (since a single process may create
	multiple independent locks) is the concatenation of
	that process-wide nonce with a unique, atomically
	incremented counter. This guarantees that there can
	never be collisions among locks created by a given
	process.

	Additionally, when a lockfile is created, it is
	created with the O_EXCL flag. This does not guarantee
	atomicity on a NFS, as described above, but it does
	mean that so long as files are created at sufficiently
	different times, collisions will be detected, and the
	process which detects a collision can regenerate its
	nonce and try again. In a normal system, in the
	unimaginably unlikely case in which a collision occurs,
	a nonce regeneration is nearly gauranteed to resolve the
	situation. However, there could be issues like a sytem
	having a bad source of randomness that would cause
	collisions to repeat. For this reason, after a nonce
	is regenerated a number of times, the TryLock call
	gives up and returns ErrCollision. A caller should take
	this as an indication that there is something wrong with
	the system - likely a bad source of randomness.

	There is one further subtlety to this algorithm. Since a
	given process prevents collisions among multiple locks
	that it uses with an atomic counter, forking the process
	is all but guaranteed to cause a collision. However, this
	is acceptable because the collision is guaranteed to be
	detected since the O_EXCL flag will work properly between
	processes on a single machine. The detecting process will
	then regenerate the nonce, and try again, at which point
	the same probabilities discussed above hold again.

	TODO(synful): is this actually true? Maybe O_EXCL on NFS
	mounts is enforced by the NFS server itself, so being on
	the same host doesn't change anything?
*/

const (
	legacyLinkfilename   = "link"
	legacyRandBytes      = 64
	legacyCollisionRetry = 10
)

var (
	legacyNonce     [legacyRandBytes]byte
	legacyHaveNonce bool
	legacyNoncelock sync.RWMutex
)

func invalidateNonce() {
	legacyNoncelock.Lock()
	legacyHaveNonce = false
	legacyNoncelock.Unlock()
}

func getNonce() ([legacyRandBytes]byte, error) {
	legacyNoncelock.RLock()
	if legacyHaveNonce {
		defer legacyNoncelock.RUnlock()
		return legacyNonce, nil
	}
	legacyNoncelock.RUnlock()
	legacyNoncelock.Lock()
	defer legacyNoncelock.Unlock()
	_, err := io.ReadFull(rand.Reader, legacyNonce[:])
	return legacyNonce, err
}

var legacyCtr uint64

func lockname() (string, error) {
	var buf [8 + legacyRandBytes]byte
	var encoded [2 * (8 + legacyRandBytes)]byte
	legacyNonce, err := getNonce()
	if err != nil {
		return "", err
	}
	copy(buf[8:], legacyNonce[:])
	binary.BigEndian.PutUint64(buf[:8], atomic.AddUint64(&legacyCtr, 1))
	hex.Encode(encoded[:], buf[:])
	return string(encoded[:]), nil
}

// LegacyLock is a handle to a file-based lock, and
// implements an algorithm which is safe to use on
// NFS versions prior to 3 and Linux kernel versions
// prior to 2.6. Like Lock, a LegacyLock requires a
// directory in which to operate, which must exist
// ahead of time, and must be empty.
//
// A process can use multiple LegacyLocks simultaneously,
// but they will behave completely independently, and
// only one can be locked at a time. LegacyLocks are safe
// for concurrent access. The zero value of LegacyLock is
// not valid, and any methods called on a zero value
// LegacyLock will panic. To acquire a valid LegacyLock,
// use NewLegacyLock.
type LegacyLock struct {
	dir      string
	lockfile string
	linkfile string
	locked   bool
	init     bool
	m        sync.Mutex
}

// NewLegacyLock creates a new LegacyLock with the given
// directory, which must be an absolute path; if it is not,
// NewLegacyLock will return ErrNeedAbsPath. NewLegacyLock
// only initializes the lock datastructure; no filesystem
// operations are performed until a call to TryLock or
// TryLockN.
func NewLegacyLock(dir string) (*LegacyLock, error) {
	if !filepath.IsAbs(dir) {
		return nil, ErrNeedAbsPath
	}
	lockname, err := lockname()
	if err != nil {
		return nil, err
	}
	return &LegacyLock{
		dir:      dir,
		lockfile: filepath.Join(dir, lockname),
		linkfile: filepath.Join(dir, legacyLinkfilename),
		init:     true,
	}, nil
}

// TryLock is equivalent to TryLockN(1, 0).
func (l *LegacyLock) TryLock() (ok bool, err error) {
	return l.TryLockN(1, 0)
}

// TryLockN attempts to acquire the lock up to
// n times before giving up, pausing for the given
// delay between each attempt. TryLockN will panic
// if l is not initialized or if the lock is already
// acquired.
func (l *LegacyLock) TryLockN(n int, delay time.Duration) (ok bool, err error) {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.init {
		panic("lockfile: uninitialized lock")
	}
	if l.locked {
		panic("lockfile: tried to lock acquired lock")
	}
	f, err := os.OpenFile(l.lockfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0666)
	for i := 0; i < legacyCollisionRetry && os.IsExist(err); i++ {
		invalidateNonce()
		var lname string
		lname, err = lockname()
		if err != nil {
			return false, ErrCollision
		}
		l.lockfile = filepath.Join(l.dir, lname)
		f, err = os.OpenFile(l.lockfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0666)
	}
	if err != nil {
		if os.IsExist(err) {
			return false, ErrCollision
		}
		return false, fmt.Errorf("creating lockfile: %v", err)
	}
	f.Close()
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

func (l *LegacyLock) tryLock() (bool, error) {
	err := os.Link(l.lockfile, l.linkfile)
	if err != nil {
		// The error might be spurious
		// (see the man 2 open section
		// on O_EXCL)
		fi, err2 := os.Stat(l.lockfile)
		if err2 != nil {
			return false, err2
		}
		if fi.Sys().(syscall.Stat_t).Nlink == 2 {
			return true, nil
		}

		// If the error wasn't spurious,
		// this just means we failed to
		// acquire the lock.
		if os.IsExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Unlock releases the lock so that others can
// acquire it. Unlock will panic if l is not
// initialized, or if the lock is not acquired.
func (l *LegacyLock) Unlock() error {
	l.m.Lock()
	defer l.m.Unlock()
	if !l.init {
		panic("lockfile: uninitialized lock")
	}
	if !l.locked {
		panic("lockfile: tried to unlock unacquired lock")
	}
	// Remove the linkfile first so that
	// we minimize the amount of time that
	// we retain the lock unnecessarily
	err := os.Remove(l.linkfile)
	if err != nil {
		return err
	}
	err = os.Remove(l.lockfile)
	if err != nil {
		return err
	}
	l.locked = false
	return nil
}
