package assignment

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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

	target := filepath.Join(cwd, "handin_test.tgz")
	_, err = os.Create(target)
	if err != nil {
		t.Fatalf("could not create target file: %v", err)
	}
	defer os.Remove(target)

	err = PerformHandin(HandinMetadata{}, target)
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

	expected := map[string][]byte{
		"./":                {},
		"./.kudos_metadata": []byte("{}\n"),
		"./foo":             []byte("foo\n"),
	}
	got := make(map[string][]byte)
	for i := 0; i < 3; i++ {
		hdr, err := tr.Next()
		if err != nil {
			t.Fatalf("could not read handin archive: %v", err)
		}
		got[hdr.Name] = make([]byte, hdr.Size)
		_, err = io.ReadFull(tr, got[hdr.Name])
		if err != nil {
			t.Fatalf("couldn't read handin archive: %v", err)
		}
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("unexpected tar contents: expectected:\n%v\n\ngot:\n%v", expected, got)
	}
}
