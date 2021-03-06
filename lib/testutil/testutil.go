package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// TB is the interface common to testing.T
// and testing.B. It is used instead of
// testing.TB so that it can be satisfied
// by types other than testing.T and
// testing.B for internal testing purposes,
// but can be treated by users of this
// package as equivalent to testing.TB.
type TB interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
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
func MustTempFile(t TB, dir, prefix string) (f *os.File) {
	f, err := ioutil.TempFile(dir, prefix)
	must(t, err)
	return f
}

// MustTempFile attempts to create a temp directory,
// and logs the error to t.Fatalf if it fails.
// The arguments dir and prefix behave as
// documented in ioutil.TempDir.
func MustTempDir(t TB, dir, prefix string) (name string) {
	name, err := ioutil.TempDir(dir, prefix)
	must(t, err)
	return name
}

// Must logs to t.Fatalf if err != nil.
func Must(t TB, err error) {
	must(t, err)
}

// MustPrefix is like Must, except that if it
// logs to t.Fatalf, the given prefix is prepended
// to the output.
func MustPrefix(t TB, prefix string, err error) {
	if err != nil {
		nfatalf(t, 1, prefix+": %v", err)
	}
}

// MustError logs to t.Fatalf if err == nil
// or if err.Error() != expect.
func MustError(t TB, expect string, err error) {
	if err == nil {
		nfatalf(t, 1, "unexpected nil error")
	}
	if err.Error() != expect {
		nfatalf(t, 1, "unexpected error: got %q; want %q", err, expect)
	}
}

// MustErrorPrefix is like MustError, except that
// if it logs to t.Fatalf, the given prefix is
// prepended to the output.
func MustErrorPrefix(t TB, prefix, expect string, err error) {
	if err == nil {
		nfatalf(t, 1, prefix+": unexpected nil error")
	}
	if err.Error() != expect {
		nfatalf(t, 1, prefix+": unexpected error: got %q; want %q", err, expect)
	}
}

func must(t TB, err error) {
	if err != nil {
		nfatalf(t, 2, "%v", err)
	}
}

// nfatalf calls t.Fatalf, but prepends the file and
// line number of the caller's nth ancestor.
func nfatalf(t TB, n int, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(n + 1)
	if !ok {
		t.Fatalf("unknown file/line: "+format, args...)
	}
	file = filepath.Base(file)
	t.Fatalf("%v:%v: "+format, append([]interface{}{file, line}, args...)...)
}
