package handin

import (
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	acl "github.com/joshlf/go-acl"
	"github.com/joshlf/kudos/lib/perm"
)

// PerformFaclHandin performs a handin of the given
// directory, writing a tar'd and gzip'd version of
// it to target.
func PerformFaclHandin(handin, target string) error {
	tf, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer tf.Close()

	cmd := exec.Command("tar", "c", filepath.Clean(handin))
	gzw := gzip.NewWriter(tf)
	defer gzw.Close()
	cmd.Stdout = gzw
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not write handin archive: %v", err)
	}
	return nil
}

// InitFaclHandin initializes dir by creating it with the
// permissions rwxrwxr-x, and for each given UID, creating
// the file <UID>.tgz with the permissions r--r-----, and
// with an ACL granting write access the user with the given
// UID.
func InitFaclHandin(dir string, uids []string) (err error) {
	// need world r-x so students can cd in
	// and write to their handin files
	err = os.Mkdir(dir, os.ModeDir|perm.Parse("rwxrwxr-x"))
	if err != nil {
		return fmt.Errorf("could not create handin dir: %v", err)
	}
	// put the defer here so that we only remove
	// the directory after we're sure we created
	// it (if we did it earlier, the error could
	// be, for example, that a file already existed
	// there, and then we'd spuriously remove it)
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	for _, uid := range uids {
		path := filepath.Join(dir, uid+".tgz")
		// make sure to use global err
		// (so defered func can check it)
		var f *os.File
		f, err = os.OpenFile(path, os.O_CREATE|os.O_EXCL, perm.Parse("r--r-----"))
		f.Close()
		if err != nil {
			return fmt.Errorf("could not create handin file: %v", err)
		}
		// TODO(joshlf): set group to ta group
		// (maybe just make handin dir g+s at
		// init?)
		a := acl.ACL{
			{acl.TagUserObj, "", os.FileMode(perm.Read)},
			{acl.TagGroupObj, "", os.FileMode(perm.Read)},
			{acl.TagOther, "", 0},
			{acl.TagUser, uid, os.FileMode(perm.Write)},
			{acl.TagMask, "", os.FileMode(perm.Write)},
		}
		err = acl.Set(path, a)
		if err != nil {
			return fmt.Errorf("could not set permissions on handin file: %v", err)
		}
	}
	return nil
}

// InitSetgidHandin initializes dir by creating it with
// the permissions rwxrwx--- (so that the setgid handin
// method is required to write files into it).
func InitSetgidHandin(dir string) (err error) {
	err = os.Mkdir(dir, os.ModeDir|perm.Parse("rwxrwx---"))
	if err != nil {
		return fmt.Errorf("could not create handin dir: %v", err)
	}
	return nil
}
