package assignment

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	HandinMetadataFileName = ".kudos_metadata"
)

type HandinMetadata struct {
	// TODO(synful)
}

func PerformHandin(metadata HandinMetadata, target string) error {
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

	tf, err := ioutil.TempFile("", "kudos")
	if err != nil {
		return fmt.Errorf("could not create temp file: %v", err)
	}
	defer os.Remove(tf.Name())

	cmd := exec.Command("tar", "c", ".")
	gzw := gzip.NewWriter(tf)
	defer gzw.Close()
	cmd.Stdout = gzw
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not create handin archive: %v", err)
	}
	err = os.Rename(tf.Name(), target)
	if err != nil {
		return fmt.Errorf("could not move handin into place: %v", err)
	}
	return nil
}
