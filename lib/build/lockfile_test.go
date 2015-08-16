package build

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/synful/kudos/lib/lockfile"
)

var (
	pre26Bin  = "testdata/pre26"
	post26Bin = "testdata/post26"
)

func init() {
	_, parentDir, _, _ := runtime.Caller(0)
	parentDir = filepath.Dir(parentDir)
	pre26Bin = filepath.Join(parentDir, pre26Bin)
	post26Bin = filepath.Join(parentDir, post26Bin)
}

func TestPre26(t *testing.T) {
	defer func() {
		pre26 = false
		versionChecked = false
	}()

	oldpath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldpath)
	err := os.Setenv("PATH", pre26Bin)
	if err != nil {
		t.Fatalf("error setting PATH: %v", err)
	}
	if !getPre26() {
		t.Errorf("getPre26 erroneously returned false")
	}

	// Make sure it uses the old value
	// once it's been set.
	pre26 = false
	if getPre26() {
		t.Errorf("getPre26 erroneously returned true")
	}

	// Reset
	versionChecked = false
	err = os.Setenv("PATH", post26Bin)
	if err != nil {
		t.Fatalf("error setting PATH: %v", err)
	}
	if getPre26() {
		t.Errorf("getPre26 erroneously returned true")
	}

	// Make sure it uses the old value
	// once it's been set.
	pre26 = true
	if !getPre26() {
		t.Errorf("getPre26 erroneously returned false")
	}
}

func TestNewLockfile(t *testing.T) {
	defer func() {
		pre26 = false
		versionChecked = false
	}()

	errstr := new(string)

	pre26 = true
	versionChecked = true
	_, err := NewLockfile("")
	*errstr = "cannot use non-legacy lockfile on Linux kernel " +
		"versions pre-2.6; please recompile with \"lockfile_legacy\" tag"
	expectError(t, errstr, err)

	pre26 = false
	_, err = NewLockfile("")
	*errstr = lockfile.ErrNeedAbsPath.Error()
	expectError(t, errstr, err)

	_, err = NewLockfile("/")
	expectError(t, nil, err)
}

func expectError(t *testing.T, expect *string, err error) {
	switch {
	case expect == nil && err != nil:
		t.Errorf("unexpected error; want <nil>; got %v", err)
	case expect != nil && (err == nil || err.Error() != *expect):
		t.Errorf("unexpected error; want %v; got %v", *expect, err)
	}
}
