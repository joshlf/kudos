package kudos

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type UserConfig struct {
	// Blacklist is a slice of UIDs
	Blacklist []string `json:"blacklist",omitempty`
}

func ParseUserConfigFile(path string) (*UserConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := parseUserConfig(f)
	if err != nil {
		return nil, fmt.Errorf("could not parse: %v", err)
	}
	return c, nil
}

func parseUserConfig(r io.Reader) (*UserConfig, error) {
	d := json.NewDecoder(r)
	var config UserConfig
	err := d.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// WriteUserConfigFile writes a json encoding of u
// to path. It does this atomically by first writing
// to a temporary location, and then atomically moving
// the temporary file to the location given by path.
func WriteUserConfigFile(path string, u *UserConfig) error {
	// create a directory (rather than a file)
	// so that we can control the permissions
	// that the file itself is created with
	dir, err := ioutil.TempDir("", "kudos")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	tmppath := filepath.Join(dir, "config")

	buf, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		return fmt.Errorf("could not marshal: %v", err)
	}

	f, err := os.OpenFile(tmppath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}

	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	err = os.Rename(tmppath, path)
	if err != nil {
		return err
	}

	return nil
}
