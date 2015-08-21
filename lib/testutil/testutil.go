package testutil

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

// MustTempFile attempts to create a temp file,
// and logs the error to t.Fatalf if it fails.
// The arguments dir and prefix behave as
// documented in ioutil.TempFile.
func MustTempFile(t *testing.T, dir, prefix string) (f *os.File) {
	var err error
	f, err = ioutil.TempFile(dir, prefix)
	must(t, "MustTempFile", err)
	return
}

// Must logs to t.Fatalf if err != nil.
func Must(t *testing.T, err error) {
	must(t, "Must", err)
}

// MustPrefix is like Must, except that if it
// logs to t.Fatalf, the given prefix is prepended
// to the output.
func MustPrefix(t *testing.T, prefix string, err error) {
	must(t, "Must: "+prefix, err)
}

func must(t *testing.T, prefix string, err error) {
	if err != nil {
		nfatalf(t, 3, prefix+": %v", err)
	}
}

// fatalf is equivalent to nfatalf with n = 2;
// it should be called only by top-level functions.
func fatalf(t *testing.T, format string, args ...interface{}) {
	nfatalf(t, 3, format, args...)
}

// nfatalf calls t.Fatalf, but prepends the file and
// line number of the caller's nth ancestor.
func nfatalf(t *testing.T, n int, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(n + 1)
	if !ok {
		t.Fatalf("testutil: unknown file/line: "+format, args...)
	}
	t.Fatalf("testutil: %v:%v: "+format, append([]interface{}{file, line}, args...)...)

}
