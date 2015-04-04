package internal

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _example_assignments_assignment_toml_sample = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x94\x54\x4f\x8f\x1a\x3f\x0c\xbd\xef\xa7\xb0\xe0\xb2\x2b\x01\x02\xf4\xe3\xf2\x93\xf6\xb0\x52\x2b\xb5\x3d\x54\x95\xba\x3d\x21\x0e\x19\xe2\x30\x91\xf2\x67\x9a\x64\x40\xf3\xed\x6b\x27\xc3\xcc\x40\x77\xa5\xed\x0d\x3c\xf6\x7b\x7e\xf6\x73\xe6\xf0\xcb\xe9\xdf\x2d\xc2\xd1\x4b\x84\x36\xa2\x84\xe4\x21\xa0\xc2\xc0\x3f\x52\xad\x23\x88\x18\xf5\xc9\x59\x74\xe9\x61\x5e\xf2\x9e\x61\x56\x82\xeb\xcd\xec\x81\x82\x5f\x5a\x2b\xdc\x32\xa0\x90\xa2\x32\x08\x4e\x58\x5c\xc1\x57\x05\xce\x3b\x04\x42\x88\x0d\x1e\xb5\xd2\x28\x17\xa0\x13\x5c\xb4\x31\x54\x24\x51\x89\xd6\xa4\x42\x83\x13\x96\xcc\xb1\xa2\x0c\xc6\x61\xae\x97\xf1\x53\x4f\xf8\x89\x3a\x96\x22\x11\xcb\x77\x9f\x90\xea\x45\x02\xd9\xc7\x22\xd8\x36\x26\x50\xde\x18\x7f\x29\x0a\x94\x0f\x56\x24\x86\xe4\x24\x42\xfc\xd6\x1a\xf8\x6f\x01\xdb\xf5\x66\x07\x54\xba\xd9\xfe\xbf\x5e\x0b\x0b\x8f\x9f\x7f\xbe\x3e\x65\x82\x91\x32\xc2\x51\x38\xa8\xc5\x19\x09\xd7\x24\xdd\x90\xc0\x5a\x38\xa9\x5d\xcc\x68\xcc\xac\x15\xcd\xab\xcc\x27\x69\x8b\x71\xda\x96\x78\x0b\x09\x35\x49\x0e\x20\x86\xa6\xa9\xd2\x07\x02\x8a\x49\xbb\x63\xba\x12\x2c\xa0\x6a\x13\x4d\x31\x41\xe5\x53\xbd\x82\x17\x29\x75\xd2\xde\x09\x63\xba\x05\x95\xdc\x61\x73\x62\x86\xf7\xce\x74\xc0\xb3\x2f\x38\xb0\x04\x02\x4b\xb4\x1e\x88\xda\x36\xa6\xa3\x52\x5a\xf5\x84\x9e\x47\xb3\xdf\x97\xec\xc3\x61\xb2\x66\xa5\x43\x4c\xb3\x8f\x0e\x6e\x0e\x4d\xf0\xe4\x00\x1b\x29\x77\x3f\xe3\x3f\x9b\xd9\xe1\x3d\xec\x88\x47\xef\xe4\x1d\xf8\xee\x5f\xc0\xb7\x04\xce\xcb\x22\x59\xee\x44\x7b\xe9\x13\x56\xf0\x5a\x0f\x7f\x0a\x5d\xb6\x44\x45\x0e\x2f\x6e\x17\xd6\xbb\x13\x0f\xd0\x18\xf2\x5f\xb3\x34\x78\x46\x33\xe0\x97\xb1\x5b\xd1\x71\x05\x7b\x33\xb2\x11\x05\xf9\xb8\xad\xfa\x9c\x5e\x46\x2c\x83\xeb\x83\x37\xea\x8a\xf8\x89\x8b\x7f\xf4\xfd\xe4\x60\xe3\x35\x2f\xed\x19\x76\xeb\xa2\xe0\xda\xed\x85\x9c\x31\xe1\x89\xf9\x8e\xc4\xf0\xb9\xbe\xe9\x22\xb2\x09\xe8\xa2\xfa\xdd\x97\x33\xeb\x38\x9d\xe1\xe1\x2c\x0c\x69\x7d\xcc\x27\xc0\x47\x37\xb8\xa0\x62\xbf\x69\x47\xa6\x0d\x74\xf0\x2a\x78\x9b\x65\x4e\xaa\x22\x78\x45\xc8\x37\x64\x4f\x24\x96\xea\x46\x6b\x5f\x9b\xa2\xd5\x4e\xf2\xb2\x60\x7a\x35\x02\xf9\xb0\x29\x6e\x5d\x70\x0a\x53\x2a\x18\xfb\xa4\xe7\xa0\x74\x75\xf7\x10\x5c\xaf\xff\xfd\xb1\x6e\xdf\x1a\x6b\x0e\xc2\x58\xb5\x1a\x3b\xca\x00\x30\x3e\x5d\x25\x73\x58\xc1\x76\xf7\xd1\xca\xea\xef\xca\x3f\x01\x00\x00\xff\xff\x90\xed\x95\xf5\x41\x05\x00\x00")

func example_assignments_assignment_toml_sample_bytes() ([]byte, error) {
	return bindata_read(
		_example_assignments_assignment_toml_sample,
		"example/assignments/assignment.toml.sample",
	)
}

func example_assignments_assignment_toml_sample() (*asset, error) {
	bytes, err := example_assignments_assignment_toml_sample_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "example/assignments/assignment.toml.sample", size: 1345, mode: os.FileMode(420), modTime: time.Unix(1428354173, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _example_config_toml = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x92\x31\x6f\xe3\x30\x0c\x85\xf7\xfc\x0a\xc2\x19\x6e\xb9\x33\x2e\x3f\x20\x43\x90\xa5\x9d\x3a\x24\x7b\xa0\x58\x74\x24\x40\x12\x0d\x91\x42\xe1\x7f\x5f\x51\x76\xd2\xa0\x2d\x50\xc0\x93\xf4\xde\xc7\xa7\x47\x6f\xe1\x48\x25\x33\xfe\x61\x18\xc8\x62\x0f\x67\xe7\x19\x62\x61\x81\x68\x64\x70\x20\x0e\xeb\x8d\x4a\xc0\xfa\x8c\x83\x50\x9e\x21\x99\x88\xfd\x66\x7b\xbf\xd8\x43\x37\xf0\xee\xff\xae\xdb\xd4\xb3\x97\x12\x4d\xfa\x97\xd1\x58\x73\x0d\xb8\x48\x17\x6a\xfd\x68\x12\x4f\xc9\x84\xbf\x60\x92\x85\x77\x1f\x02\x58\x1c\x4d\x09\x52\x9d\x42\xcf\xd3\x34\x0e\xf8\x11\x12\x09\x30\x8a\x8e\x53\x96\x0e\x3b\x9e\xe0\x3e\xed\xec\x96\x11\x40\x63\x33\x9f\x0f\x70\xcb\x54\xa6\xe7\x87\x5c\x11\x78\xc2\xc1\x8f\x1e\xad\x62\xc4\x5c\x9a\xe6\x91\x5b\x0c\xff\x08\x63\x29\x16\x93\xfc\x4e\x5c\x85\x5f\xb0\xeb\xe9\x27\xdb\xd5\x47\xfb\x04\x11\xc5\x91\xed\xe1\xad\x95\xc1\x60\x32\x42\x37\x9a\x21\x74\xad\x95\xae\xbe\xf6\xe6\x6d\xa7\xe0\x57\xd1\xd6\x9c\xbf\xb9\x30\x43\x6d\x9f\x62\xc4\x64\xd1\x6a\x57\xa5\x96\xa4\x29\xd5\xb9\x32\xa1\xa4\x80\xcc\xd5\x38\xd7\x12\x81\x67\x16\x8c\x60\x09\x79\xa9\xb1\x4c\x13\x65\x69\x0e\x56\xfc\x12\xe8\xb2\x9a\xf7\x6b\x0a\x8d\x7b\xa8\x7b\xe1\x21\xfb\x16\xf1\xde\xc7\xb2\x99\xef\xdb\x54\x54\x93\x4f\xbe\xc9\x1f\x2b\x52\x91\x49\xe0\x93\x64\xb2\x65\xf9\x77\xd6\xf5\xd6\x1e\x8e\xa7\xbe\xdb\x7c\x04\x00\x00\xff\xff\x73\x39\xb3\x24\x81\x02\x00\x00")

func example_config_toml_bytes() ([]byte, error) {
	return bindata_read(
		_example_config_toml,
		"example/config.toml",
	)
}

func example_config_toml() (*asset, error) {
	bytes, err := example_config_toml_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "example/config.toml", size: 641, mode: os.FileMode(420), modTime: time.Unix(1428354206, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"example/assignments/assignment.toml.sample": example_assignments_assignment_toml_sample,
	"example/config.toml":                        example_config_toml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() (*asset, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"example": &_bintree_t{nil, map[string]*_bintree_t{
		"assignments": &_bintree_t{nil, map[string]*_bintree_t{
			"assignment.toml.sample": &_bintree_t{example_assignments_assignment_toml_sample, map[string]*_bintree_t{}},
		}},
		"config.toml": &_bintree_t{example_config_toml, map[string]*_bintree_t{}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	if err != nil { // File
		return RestoreAsset(dir, name)
	} else { // Dir
		for _, child := range children {
			err = RestoreAssets(dir, path.Join(name, child))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
