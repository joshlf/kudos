package handin

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
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
	// defer os.RemoveAll(testDir)

	usr, err := user.Current()
	testutil.MustPrefix(t, "could not get current user", err)

	// Unfortunately, there's no way to actually exercise the
	// permissions functionality since we own the file - if
	// we deprive ourselves of write permissions but add a
	// write ACL, we still won't be able to write to the file
	// - the permissions check will fail on the basis of us
	// being the file owner. Thus, the best we can do is just
	// manually check that the permissions are what we expect,
	// and expect that everything should work properly in the
	// real world as a result.
	targetPath := filepath.Join(testDir, "handin", usr.Uid+".tgz")
	err = InitFaclHandin(filepath.Join(testDir, "handin"), []string{usr.Uid})
	testutil.MustPrefix(t, "could not init handin directory", err)

	a, err := acl.Get(targetPath)
	expect := acl.ACL{
		{acl.TagUserObj, "", os.FileMode(perm.Read)},
		{acl.TagUser, usr.Uid, os.FileMode(perm.Write)},
		{acl.TagGroupObj, "", os.FileMode(perm.Read)},
		{acl.TagMask, "", os.FileMode(perm.Write)},
		{acl.TagOther, "", 0},
	}
	if !reflect.DeepEqual(a, expect) {
		t.Fatalf("file has wrong permissions: want %v; got %v", expect, a)
	}

	// remove and recreate with write permissions
	testutil.MustPrefix(t, "could not remove handin target", os.Remove(targetPath))
	f, err := os.Create(targetPath)
	testutil.MustPrefix(t, "could not create handin target", err)
	f.Close()

	handinPath := filepath.Join(testDir, "to_handin")
	err = os.Mkdir(handinPath, 0700)
	testutil.MustPrefix(t, "could not create directory to hand in", err)

	err = ioutil.WriteFile(filepath.Join(testDir, "to_handin", "foo"), []byte("foo\n"), 0600)
	testutil.MustPrefix(t, "could not write handin file", err)

	err = PerformFaclHandin(handinPath, targetPath)
	testutil.MustPrefix(t, "could not perform facl handin", err)

	f, err = os.Open(targetPath)
	testutil.MustPrefix(t, "could not open handin archive", err)

	gr, err := gzip.NewReader(f)
	testutil.MustPrefix(t, "could not create gzip reader", err)

	tr := tar.NewReader(gr)

	expected := map[string][]byte{
		// strip leading slash
		handinPath[1:] + "/":    {},
		handinPath[1:] + "/foo": []byte("foo\n"),
	}
	got := make(map[string][]byte)
	for i := 0; i < 2; i++ {
		hdr, err := tr.Next()
		testutil.MustPrefix(t, "could not read handin archive", err)
		got[hdr.Name] = make([]byte, hdr.Size)
		_, err = io.ReadFull(tr, got[hdr.Name])
		testutil.MustPrefix(t, "could not read handin archive", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("unexpected tar contents: want:\n%v\n\ngot:\n%v", expected, got)
	}
}
