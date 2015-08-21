package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

/*
	UTILITY CODE

	If you add or remove lines, most of the tests
	will fail because the line numbers won't match;
	make sure to re-run the tests and update the
	line numbers after any modifications.
*/

// return the "file:line" of the caller's nth ancestor
func fileLinePrefix(n int) string {
	_, file, line, ok := runtime.Caller(n + 1)
	if ok {
		return fmt.Sprintf("%v:%v", filepath.Base(file), line)
	}
	return "unknown file/line"
}

// Use mockError when panicking in tests
// so that when we recover we can tell it
// apart from other panics.
type mockError string

func (m mockError) Error() string { return string(m) }

// mockT implements the testingT interface
type mockT struct{}

func (m mockT) Fatalf(format string, args ...interface{}) {
	panic(mockError(fmt.Sprintf(format, args...)))
}

var mock mockT

func expectFatal(t *testing.T, str string, f func()) {
	// Calculate prefix ahead of time because if there's
	// a panic, the defered function will get called from
	// the call site of the panic, and we can't predict
	// what that stack will look like.
	prefix := fileLinePrefix(1)
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("%v: function failed to call Fatalf", prefix)
		}
		me, ok := r.(mockError)
		if !ok {
			panic(r)
		}
		if string(me) != str {
			t.Errorf("%v: unexpected call to Fatalf: want %v; got %v", prefix, str, me)
		}
	}()
	f()
}

func expectFatalRegex(t *testing.T, re *regexp.Regexp, f func()) {
	// Calculate prefix ahead of time because if there's
	// a panic, the defered function will get called from
	// the call site of the panic, and we can't predict
	// what that stack will look like.
	prefix := fileLinePrefix(1)
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("%v: function failed to call Fatalf", prefix)
		}
		me, ok := r.(mockError)
		if !ok {
			panic(r)
		}
		if !re.MatchString(string(me)) {
			t.Errorf("%v: unexpected call to Fatalf: want regex %v; got %v", prefix, re, me)
		}
	}()
	f()
}

func expectSuccess(t *testing.T, f func()) {
	// Calculate prefix ahead of time because if there's
	// a panic, the defered function will get called from
	// the call site of the panic, and we can't predict
	// what that stack will look like.
	prefix := fileLinePrefix(1)
	defer func() {
		r := recover()
		if me, ok := r.(mockError); ok {
			t.Errorf("%v: function called Fatalf: %v", prefix, me)
		} else if r != nil {
			panic(r)
		}
	}()
	f()
}

/*
	TESTS
*/

func TestMust(t *testing.T) {
	expectSuccess(t, func() { mustImpl(mock, nil) })
	f := func() { mustImpl(mock, errors.New("foo")) }
	expectFatal(t, "testutil_test.go:66: foo", f)
}

func TestMustPrefix(t *testing.T) {
	expectSuccess(t, func() { mustPrefix(mock, "", nil) })
	f := func() { mustPrefix(mock, "foo", errors.New("bar")) }
	expectFatal(t, "testutil_test.go:66: foo: bar", f)
}

func TestMustTempFile(t *testing.T) {
	var name string
	expectSuccess(t, func() { mustTempFile(mock, "", "") })
	defer os.Remove(name)

	// Make a directory we know for a fact is empty
	dir := MustTempDir(t, "", "")
	defer os.RemoveAll(dir)
	nonexistant := filepath.Join(dir, "foo")
	re := regexp.MustCompile("^testutil_test.go:88: open " + nonexistant +
		".*: no such file or directory$")
	expectFatalRegex(t, re, func() { mustTempFile(mock, nonexistant, "") })
}

func TestMustTempDir(t *testing.T) {
	var name string
	expectSuccess(t, func() { mustTempDir(mock, "", "") })
	defer os.Remove(name)

	// Make a directory we know for a fact is empty
	dir := MustTempDir(t, "", "")
	defer os.RemoveAll(dir)
	nonexistant := filepath.Join(dir, "foo")
	re := regexp.MustCompile("^testutil_test.go:88: mkdir " + nonexistant +
		".*: no such file or directory$")
	expectFatalRegex(t, re, func() { mustTempDir(mock, nonexistant, "") })

}
