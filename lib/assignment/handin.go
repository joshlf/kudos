package assignment

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/synful/kudos/lib/config"
	"github.com/synful/kudos/lib/perm"
)

const (
	HandinMetadataFileName = ".kudos_metadata"
)

type HandinMetadata struct {
	// TODO(synful)
}

func PerformFaclHandin(metadata HandinMetadata, target string) error {
	// TODO(synful): reimplement using Go's tar package
	// to eliminate the dependancy on tar and the need
	// to first write the metadata file out to the dir.

	mf, err := os.OpenFile(HandinMetadataFileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not create metadata file: %v", err)
	}
	defer mf.Close()

	enc := json.NewEncoder(mf)
	err = enc.Encode(metadata)
	if err != nil {
		return fmt.Errorf("could not write metadata file: %v", err)
	}

	tf, err := os.OpenFile(target, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("could open target: %v", err)
	}
	defer tf.Close()

	cmd := exec.Command("tar", "c", ".")
	gzw := gzip.NewWriter(tf)
	defer gzw.Close()
	cmd.Stdout = gzw
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not write handin archive: %v", err)
	}
	return nil
}

// TODO(synful): have this take a slice of student
// config objects instead of uid strings
func InitFaclHandin(course config.Course, asgn config.Assignment, uids []string) (err error) {
	dir := asgn.HandinDir()

	// need world r-x so students can cd in
	// and write to their handin files
	err = os.Mkdir(dir, os.ModeDir|perm.Parse("rwxrwxr-x"))
	if err != nil {
		return fmt.Errorf("could not create handin dir: %v", err)
	}
	defer func(dir string) {
		if err != nil {
			os.RemoveAll(dir)
		}
	}(dir)

	for _, uid := range uids {
		path := filepath.Join(dir, uid+".tgz")
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, perm.Parse("r--r-----"))
		f.Close()
		if err != nil {
			return fmt.Errorf("could not create handin file: %v", err)
		}
		// TODO(synful): set group to ta group
		// (maybe just make handin dir g+s at
		// init?)

		facl := perm.Facl{perm.User, uid, perm.Write}
		err = perm.SetFacl(uid, facl)
		if err != nil {
			return fmt.Errorf("could not set permissions on handin file: %v", err)
		}
	}
	return nil
}

func InitSetgidHandin(course config.Course, asgn config.Assignment) (err error) {
	dir := asgn.HandinDir()
	err = os.Mkdir(dir, os.ModeDir|perm.Parse("rwxrwx---"))
	if err != nil {
		return fmt.Errorf("could not create handin dir: %v", err)
	}
	return nil
}
