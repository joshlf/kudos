package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func testPanic(t *testing.T, f func(), panic string) {
	_, file, line, _ := runtime.Caller(1)
	defer func() {
		p := recover()
		if p == nil || fmt.Sprint(p) != panic {
			t.Errorf("unexpected panic at %v:%v: want %v; got %v", file, line, panic, p)
		}
	}()
	f()
}

func testError(t *testing.T, f func() error, err string) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)
	got := f()
	if got == nil || got.Error() != err {
		t.Fatalf("unexpected error at %v:%v: want %v; got %v", file, line, err, got)
	}
}
