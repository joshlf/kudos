package handin

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"testing"

	acl "github.com/joshlf/go-acl"
	"github.com/joshlf/kudos/lib/perm"
	"github.com/joshlf/kudos/lib/testutil"
)

func TestFaclHandin(t *testing.T) {
	testDir := testutil.MustTempDir(t, "", "kudos")
	defer func() {
		err := os.RemoveAll(testDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not remove temp directory: %v\n", err)
		}
	}()

	/*
		Set up the handin directory
	*/

	usr, err := user.Current()
	testutil.Must(t, err)

	// Unfortunately, there's no way to actually exercise the
	// permissions functionality since we own the file and
	// folder - if we deprive ourselves of write permissions
	// but add a write ACL, we still won't be able to write
	// to the file - the permissions check will fail on the
	// basis of us being the file owner. The same goes for
	// read/execute on the folder. Thus, the best we can do
	// is just manually check that the permissions are what
	// we expect, and expect that everything should work
	// properly in the real world as a result.
	targetDirPath := filepath.Join(testDir, "handin", usr.Uid)
	targetFilePath := filepath.Join(targetDirPath, handinFileName)
	err = InitFaclHandin(filepath.Join(testDir, "handin"), []string{usr.Uid})
	testutil.Must(t, err)

	a, err := acl.Get(targetDirPath)
	testutil.Must(t, err)
	expect := acl.ACL{
		{acl.TagUserObj, "", perm.ParseSingle("rwx")},
		{acl.TagUser, usr.Uid, perm.ParseSingle("r-x")},
		{acl.TagGroupObj, "", perm.ParseSingle("rwx")},
		{acl.TagMask, "", perm.ParseSingle("r-x")},
		{acl.TagOther, "", 0},
	}
	if !reflect.DeepEqual(a, expect) {
		t.Fatalf("directory has wrong permissions: want %v; got %v", expect, a)
	}

	a, err = acl.Get(targetFilePath)
	testutil.Must(t, err)
	expect = acl.ACL{
		{acl.TagUserObj, "", os.FileMode(perm.Read)},
		{acl.TagUser, usr.Uid, os.FileMode(perm.Write)},
		{acl.TagGroupObj, "", os.FileMode(perm.Read)},
		{acl.TagMask, "", os.FileMode(perm.Write)},
		{acl.TagOther, "", 0},
	}
	if !reflect.DeepEqual(a, expect) {
		t.Fatalf("file has wrong permissions: want %v; got %v", expect, a)
	}

	/*
		Set up the handin
	*/

	// remove and recreate with write permissions
	testutil.Must(t, os.Remove(targetFilePath))
	f, err := os.Create(targetFilePath)
	testutil.Must(t, err)
	f.Close()

	handinPath := filepath.Join(testDir, "to_handin")
	err = os.Mkdir(handinPath, 0700)
	testutil.Must(t, err)

	err = ioutil.WriteFile(filepath.Join(testDir, "to_handin", "foo"), []byte("foo\n"), 0600)
	testutil.Must(t, err)

	/*
		Perform the handin
	*/

	pwd, err := os.Getwd()
	testutil.Must(t, err)
	testutil.Must(t, os.Chdir(handinPath))
	defer os.Chdir(pwd)
	err = PerformFaclHandin(targetFilePath)
	testutil.Must(t, err)

	/*
		Verify the handin
	*/

	f, err = os.Open(targetFilePath)
	testutil.Must(t, err)

	gr, err := gzip.NewReader(f)
	testutil.Must(t, err)

	tr := tar.NewReader(gr)

	expected := map[string][]byte{
		"./":    {},
		"./foo": []byte("foo\n"),
	}
	got := make(map[string][]byte)
	for i := 0; i < 2; i++ {
		hdr, err := tr.Next()
		testutil.Must(t, err)
		got[hdr.Name] = make([]byte, hdr.Size)
		_, err = io.ReadFull(tr, got[hdr.Name])
		testutil.Must(t, err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("unexpected tar contents: want:\n%v\n\ngot:\n%v", expected, got)
	}

	/*
		Extract the handin
	*/

	extractDir := filepath.Join(testDir, "extract")
	testutil.Must(t, os.Mkdir(extractDir, 0700))
	testutil.Must(t, ExtractHandin(targetFilePath, extractDir))

	/*
		Verify the extracted handin
	*/

	err = exec.Command("diff", "-rN", handinPath, extractDir).Run()
	testutil.Must(t, err)
}

func TestSetgidHandin(t *testing.T) {
	testDir := testutil.MustTempDir(t, "", "kudos")
	defer os.RemoveAll(testDir)

	dir := filepath.Join(testDir, "handin")
	testutil.Must(t, InitSetgidHandin(dir))

	a, err := acl.Get(dir)
	testutil.Must(t, err)

	if acl.ToUnix(a) != perm.Parse("rwxrwx---") {
		t.Errorf("bad handin directory permissions: want %v; got %v", perm.Parse("rwxrwx---"), acl.ToUnix(a))
	}
}
