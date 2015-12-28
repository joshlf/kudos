package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/joshlf/kudos/lib/testutil"
)

func TestLock(t *testing.T) {
	testDir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(testDir)
	lock, err := New(filepath.Join(testDir, "lock"))
	testutil.Must(t, err)
	ok, err := lock.TryLock()
	testutil.Must(t, err)
	if !ok {
		t.Errorf("failed to acquire lock")
	}
	testutil.Must(t, lock.Unlock())

	// make sure that when we unlock it,
	// it can be successfully locked again
	ok, err = lock.TryLock()
	testutil.Must(t, err)
	if !ok {
		t.Errorf("failed to acquire lock")
	}
}

func TestParallel(t *testing.T) {
	// This test spawns a number of goroutines
	// which each try to acquire a single lock
	// at the same time (they are synchronized
	// by all sleeping until the same time;
	// empirically this appears to bring their
	// executions within 0.1ms of one another).
	//
	// NOTE: This test is basically useless if
	// the goroutines are not run on different
	// CPU cores; be wary of trusting this test
	// if that might not happen (ie, pre-go1.5
	// with GOMAXPROCS=1, or on a single-core
	// machine).

	testDir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(testDir)

	// make sure that every process is running
	// at least one goroutine
	numLocks := 2 * runtime.GOMAXPROCS(0)

	locks := make([]*Lock, numLocks)
	for i := range locks {
		var err error
		locks[i], err = New(filepath.Join(testDir, "lock"))
		testutil.Must(t, err)
	}

	chans := make([]chan bool, numLocks)
	for i := range chans {
		chans[i] = make(chan bool, 1)
	}

	worker := func(l *Lock, target time.Time, c chan bool) {
		time.Sleep(target.Sub(time.Now()))
		ok, err := l.TryLock()
		c <- ok
		testutil.Must(t, err)
	}

	target := time.Now().Add(time.Millisecond)

	for i := 0; i < numLocks; i++ {
		go worker(locks[i], target, chans[i])
	}

	// read fromthe channel here so that
	// we clean up even if testutil.Must
	// calls t.Fatal (and thus panics)
	defer func() {
		numTrue := 0
		bools := make([]bool, numLocks)
		for i := range bools {
			bools[i] = <-chans[i]
			if bools[i] {
				numTrue++
			}
		}

		// don't call t.Fatalf if it's already
		// been called
		r := recover()
		if r == nil && numTrue != 1 {
			str := fmt.Sprintf("%v:%v", 0, bools[0])
			for i, b := range bools[1:] {
				str += fmt.Sprintf(" %v:%v", i+1, b)
			}
			t.Fatalf("bad number of locks acquired (%v): %v", numTrue, str)
		}
	}()
}

func TestTryLockN(t *testing.T) {
	testDir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(testDir)
	lock, err := New(filepath.Join(testDir, "lock"))
	testutil.Must(t, err)
	ok, err := lock.TryLock()
	testutil.Must(t, err)
	if !ok {
		t.Errorf("failed to acquire lock")
	}

	lock2, err := New(filepath.Join(testDir, "lock"))
	testutil.Must(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Done()
	go func() {
		time.Sleep(time.Millisecond)
		lock.Unlock()
	}()
	// wait for 10 * 105us = 1.05ms (just after
	// the goroutine will unlock the lock)
	ok, err = lock2.TryLockN(10, 105*time.Microsecond)
	testutil.Must(t, err)
	if !ok {
		t.Errorf("failed to acquire lock")
	}
}

func TestErr(t *testing.T) {
	lock, err := New("foo")
	testutil.MustError(t, "need absolute path", err)

	// will try to create /dev/null/foo
	// (which should obviously fail)
	lock, err = New("/dev/null/foo")
	testutil.Must(t, err)
	ok, err := lock.TryLock()
	testutil.MustError(t, "open /dev/null/foo: not a directory", err)
	if ok {
		t.Errorf("unexpected return from TryLock: got %v; want false", ok)
	}

	testDir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(testDir)
	lock, err = New(filepath.Join(testDir, "lock"))
	testutil.Must(t, err)
	ok, err = lock.TryLock()
	testutil.Must(t, err)
	if !ok {
		t.Errorf("failed to acquire lock")
	}
	testPanic(t, "lockfile: tried to lock acquired lock", func() { lock.TryLock() })
	testutil.Must(t, lock.Unlock())
	testPanic(t, "lockfile: tried to unlock unacquired lock", func() { lock.Unlock() })

	testPanic(t, "lockfile: uninitialized lock", func() { (&Lock{}).TryLock() })
	testPanic(t, "lockfile: uninitialized lock", func() { (&Lock{}).Unlock() })
}

func testPanic(t *testing.T, err string, f func()) {
	defer func() {
		var prefix string
		_, file, _, ok := runtime.Caller(2)
		if ok {
			prefix = filepath.Base(file) + ": "
		}
		r := recover()
		if r == nil || fmt.Sprint(r) != err {
			t.Errorf(prefix+"wrong panic: want %v; got %v", err, r)
		}
	}()
	f()
}
