package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TODO(joshlf): Add ExpectError method that expects
// a non-nil error with Error() matching a given string

type testingT interface {
	Fatalf(format string, args ...interface{})
}

// SrcDir attempts to figure out what source
// file it is called from, and returns the
// parent directory of that file. This can
// be useful for tests which have local test
// data, since commands such as
//  go test ./...
// can make it so that the current working
// directory is not necessarily the same as
// the source directory.
func SrcDir() (dir string, ok bool) {
	var f string
	_, f, _, ok = runtime.Caller(1)
	if !ok {
		return
	}
	return filepath.Dir(f), true
}

// MustTempFile attempts to create a temp file,
// and logs the error to t.Fatalf if it fails.
// The arguments dir and prefix behave as
// documented in ioutil.TempFile.
func MustTempFile(t *testing.T, dir, prefix string) (f *os.File) {
	return mustTempFile(t, dir, prefix)
}

func mustTempFile(t testingT, dir, prefix string) *os.File {
	f, err := ioutil.TempFile(dir, prefix)
	must(t, err)
	return f
}

// MustTempFile attempts to create a temp directory,
// and logs the error to t.Fatalf if it fails.
// The arguments dir and prefix behave as
// documented in ioutil.TempDir.
func MustTempDir(t *testing.T, dir, prefix string) (name string) {
	return mustTempDir(t, dir, prefix)
}

func mustTempDir(t testingT, dir, prefix string) string {
	name, err := ioutil.TempDir(dir, prefix)
	must(t, err)
	return name
}

// Must logs to t.Fatalf if err != nil.
func Must(t *testing.T, err error) {
	mustImpl(t, err)
}

// can't name it just "must" because we've
// already got one of those; can't just call
// must directly from Must because a) we
// need this for testing and, b) we need the
// right call stack depth because must assumes
// it's a certain depth from the client caller.
func mustImpl(t testingT, err error) {
	must(t, err)
}

// MustPrefix is like Must, except that if it
// logs to t.Fatalf, the given prefix is prepended
// to the output.
func MustPrefix(t *testing.T, prefix string, err error) {
	mustPrefix(t, prefix, err)
}

func mustPrefix(t testingT, prefix string, err error) {
	if err != nil {
		nfatalf(t, 2, prefix+": %v", err)
	}
}

func must(t testingT, err error) {
	if err != nil {
		nfatalf(t, 3, "%v", err)
	}
}

// fatalf is equivalent to nfatalf with n = 3;
// it should be called only by second-level functions.
func fatalf(t testingT, format string, args ...interface{}) {
	nfatalf(t, 4, format, args...)
}

// nfatalf calls t.Fatalf, but prepends the file and
// line number of the caller's nth ancestor.
func nfatalf(t testingT, n int, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(n + 1)
	if !ok {
		t.Fatalf("unknown file/line: "+format, args...)
	}
	file = filepath.Base(file)
	t.Fatalf("%v:%v: "+format, append([]interface{}{file, line}, args...)...)

}
