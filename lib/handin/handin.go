package handin

import (
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	acl "github.com/joshlf/go-acl"
	"github.com/joshlf/kudos/lib/config"
	"github.com/joshlf/kudos/lib/perm"
)

// PerformFaclHandin performs a handin of the current
// directory, writing a tar'd and gzip'd version of
// it to target.
func PerformFaclHandin(target string) (err error) {
	f, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command("tar", "-c", ".")
	gzw := gzip.NewWriter(f)
	defer gzw.Close()
	cmd.Stdout = gzw
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// ExtractHandin extracts the given handin (which must
// be a tar'd and gzip'd file) to the target directory,
// which must already exist.
func ExtractHandin(handin, target string) (err error) {
	f, err := os.Open(handin)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	cmd := exec.Command("tar", "-x", "-C", target)
	cmd.Stdin = gzr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// InitFaclHandin initializes dir by creating it with the
// permissions rwxrwxr-x, and for each given UID, creating
// the folder <UID> with the permissions rwxrwx--- and with
// an ACL grating read and execute permissions to the user
// with the given UID. Finally, inside this folder, the file
// handin.tgz is created with the permissions r--r-----, and
// with an ACL granting write access to the user. An example
// handin directory structure might look like:
//
//  hw01/                   (u::rwx,g::rwx,o::r-x)
//       1234/              (u::rwx,g::rwx,o::---,u:1234:r-x)
//            handin.tgz    (u::r--,g::r--,o::---,u:1234:-w-)
//       5678/              (u::rwx,g::rwx,o::---,u:5678:r-x)
//            handin.tgz    (u::r--,g::r--,o::---,u:5678:-w-)
//
// The motivation for this design is the following. Granting
// only a write ACL to the student is sufficient to allow them
// to write to the file, but it does not allow them to modify
// the timestamp on the file. The only way for the timestamp
// to be modified is to write to the file, in which case it is
// updated to the current time, which means that the timestamp
// should accurately reflect the time at which they handed in.
//
// However, this alone is not enough. If all files were to be
// in the same directory, students would be able to see the
// metadata about others' handins, which would allow them to
// see how large other students' handins were, and more
// importantly, to see whether and when other students had
// handed in. By placing each handin in its own folder that
// only the student has read/execute permissions on, other
// students are prevented from learning anything about handins
// other than their own.
func InitFaclHandin(dir string, uids []string) (err error) {
	// need world r-x so students can cd in
	// and write to their handin files
	mode := perm.Parse("rwxrwxr-x")
	err = os.Mkdir(dir, mode)
	if err != nil {
		return fmt.Errorf("could not create handin directory: %v", err)
	}
	// set permissions explicitly since original permissions
	// might be masked (by umask)
	err = os.Chmod(dir, mode)
	if err != nil {
		return fmt.Errorf("could not set permissions on handin directory: %v", err)
	}

	// put the defer here so that we only remove
	// the directory after we're sure we created
	// it (if we did it earlier, the error could
	// be, for example, that a file already existed
	// there, and then we'd spuriously remove it)
	//
	// TODO(joshlf): are we sure we want to do this?
	// maybe it'd be better (and even more idiomatic?)
	// to leave it so that the contents can be
	// inspected so that the user can figure out
	// what went wrong?
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	// TODO(joshlf): set group to ta group
	// (or maybe just make global handin dir
	// g+s at course init?)

	for _, uid := range uids {
		path := filepath.Join(dir, uid)
		filepath := filepath.Join(path, config.HandinFileName)

		// make sure to use global err
		// (so defered func can check it)
		err = os.Mkdir(path, perm.Parse("rwxrwx---"))
		if err != nil {
			return fmt.Errorf("could not create handin directory: %v", err)
		}

		// if this code changes, make sure that the
		// permissions on path are still set explicitly
		// (relying on os.Mkdir is not enough - umask
		// might change the permissions)
		a := append(
			acl.FromUnix(perm.Parse("rwxrwx---")),
			acl.Entry{acl.TagUser, uid, perm.ParseSingle("r-x")},
			acl.Entry{acl.TagMask, "", perm.ParseSingle("r-x")},
		)
		err = acl.Set(path, a)
		if err != nil {
			return fmt.Errorf("could not set permissions on handin directory: %v", err)
		}

		var f *os.File
		f, err = os.Create(filepath)
		f.Close()
		if err != nil {
			return fmt.Errorf("could not create handin file: %v", err)
		}

		// if this code changes, make sure that the
		// permissions on filepath are still set explicitly
		// (relying on os.Mkdir is not enough - umask
		// might change the permissions)
		a = append(
			acl.FromUnix(perm.Parse("r--r-----")),
			acl.Entry{acl.TagUser, uid, perm.ParseSingle("-w-")},
			acl.Entry{acl.TagMask, "", perm.ParseSingle("-w-")},
		)
		err = acl.Set(filepath, a)
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
	mode := perm.Parse("rwxrwx---")
	err = os.Mkdir(dir, mode)
	if err != nil {
		return fmt.Errorf("could not create handin directory: %v", err)
	}
	// set permissions explicitly since original permissions
	// might be masked (by umask)
	err = os.Chmod(dir, mode)
	if err != nil {
		return fmt.Errorf("could not set permissions on handin directory: %v", err)
	}
	return nil
}
