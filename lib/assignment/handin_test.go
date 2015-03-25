package assignment

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHandin(t *testing.T) {
	cwd, err := os.Getwd()

	testDir, err := ioutil.TempDir(".", "kudos_test")
	if err != nil {
		t.Fatalf("could not create test directory: %v", err)
	}
	defer os.RemoveAll(filepath.Join(cwd, testDir))

	err = ioutil.WriteFile(filepath.Join(testDir, "foo"), []byte("foo\n"), 0666)
	if err != nil {
		t.Fatalf("could not write test file: %v", err)
	}
	if err != nil {
		t.Fatalf("could not get pwd: %v", err)
	}
	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("could not change to test directory: %v", err)
	}
	err = PerformHandin(HandinMetadata{}, filepath.Join(cwd, "handin_test.tgz"))
	defer os.Remove(filepath.Join(cwd, "handin_test.tgz"))
	if err != nil {
		t.Fatalf("could not perform handin: %v", err)
	}
	f, err := os.Open(filepath.Join(cwd, "handin_test.tgz"))
	if err != nil {
		t.Fatalf("could not open handin archive: %v", err)
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("could not create gzip reader: %v", err)
	}
	tr := tar.NewReader(gr)
	hdr, err := tr.Next()
	if err != nil {
		t.Fatalf("could not read handin archive: %v", err)
	}
	if hdr.Name != "./" {
		t.Errorf("unexpected file name in handin archive: expected \"./\"; got \"%v\"", hdr.Name)
	}
	if hdr.Size != 0 {
		t.Errorf("unexpected file size in handin archive: expected 0; got %v", hdr.Size)
	}
	hdr, err = tr.Next()
	if err != nil {
		t.Fatalf("could not read handin archive: %v", err)
	}
	if hdr.Name != "./.kudos_metadata" {
		t.Errorf("unexpected file name in handin archive: expected \"./.kudos_metadata\"; got \"%v\"", hdr.Name)
	}
	hdr, err = tr.Next()
	if err != nil {
		t.Fatalf("could not read handin archive: %v", err)
	}
	if hdr.Name != "./foo" {
		t.Errorf("unexpected file name in handin archive: expected \"./foo\"; got \"%v\"", hdr.Name)
	}
	if hdr.Size != 4 {
		t.Errorf("unexpected file size in handin archive: expected 4; got %v", hdr.Size)
	}
	buf := make([]byte, hdr.Size)
	_, err = io.ReadFull(tr, buf)
	if err != nil {
		t.Errorf("couldn't read from handin archive: %v", err)
	}
	if string(buf) != "foo\n" {
		t.Errorf("unexpected file contents: expected \"foo\\n\"; got \"%s\"", string(buf))
	}
}
