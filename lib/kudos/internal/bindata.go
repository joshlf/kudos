// Code generated by go-bindata.
// sources:
// example/assignments/assignment.sample
// example/config
// DO NOT EDIT!

package internal

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _exampleAssignmentsAssignmentSample = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x94\x90\x31\x6f\xc3\x20\x14\x84\x67\xf3\x2b\x9e\x98\x5a\x89\x01\x50\x59\xb2\x75\xe8\xd2\xa9\x52\xbb\x55\x19\x70\xa0\x2d\x52\x80\x28\xe0\x29\xf2\x7f\x2f\x86\x94\x92\xd4\x56\x1d\x2f\x88\xbb\xf3\x77\x8f\x77\x42\x1d\xde\x79\xa5\xf1\x06\xb0\x0c\xc1\x7c\x3a\xca\x30\x49\xa2\x93\x36\x8b\x8f\x59\xb4\xda\x45\x38\x3b\x6a\xc8\xc6\xf3\xb0\x87\x07\x02\x9c\x32\x01\x32\x02\xe3\x1b\x4a\xa5\x85\xbb\xa7\xd7\xb7\xfb\x94\x4b\xc1\x2f\xe9\x94\x71\x21\x85\xdf\x51\x37\x7d\xa7\x72\xfc\x36\x7e\x98\x63\x88\x13\xb4\xc8\xeb\xc8\xe7\xf0\xe1\xe8\xfb\xbd\xb6\x19\x9f\x2f\x0c\x6f\x8b\x37\x92\x85\xba\xa0\x77\xde\xa9\xd9\x3e\x71\x7b\x1f\xaf\x7d\xf9\xd8\xe6\x37\xb7\xa9\x85\x29\xca\xac\x95\xfb\xb3\xe8\x97\xf2\x27\x34\xd6\xc1\x1b\x17\x27\x94\xa0\xff\xbc\xac\xcc\xb3\xc8\xe4\xb3\xcc\xaa\x85\xa1\xff\x3b\x77\x77\xd1\x74\xd5\x27\x2b\xf0\x1a\xcb\x45\x6b\x8c\x64\x05\xac\x5f\x0b\x6b\x2e\x97\xab\x47\x23\xfa\x0e\x00\x00\xff\xff\xd5\x69\xce\x59\xc9\x02\x00\x00")

func exampleAssignmentsAssignmentSampleBytes() ([]byte, error) {
	return bindataRead(
		_exampleAssignmentsAssignmentSample,
		"example/assignments/assignment.sample",
	)
}

func exampleAssignmentsAssignmentSample() (*asset, error) {
	bytes, err := exampleAssignmentsAssignmentSampleBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "example/assignments/assignment.sample", size: 713, mode: os.FileMode(420), modTime: time.Unix(1451364025, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _exampleConfig = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xe6\xe2\x54\x4a\xce\x4f\x49\x55\xb2\x52\x50\x4a\x2e\x36\x34\x30\x54\xd2\x01\x8a\xe4\x25\xe6\x82\x45\x9c\x83\x15\x60\x42\x25\x89\xf1\xe9\x45\xf9\xa5\x05\x70\x85\x25\x89\xc5\x60\x89\x94\xd4\xe2\xe4\xa2\xcc\x82\x92\xcc\xfc\x3c\x84\x16\x85\xcc\x62\x85\xc4\x3c\x85\xcc\xbc\x92\xa2\xfc\x94\xd2\xe4\x92\xfc\xa2\x4a\x85\xe4\xfc\xd2\xa2\xe2\x54\xa0\x98\x82\x73\xb0\x9e\x12\x57\x2d\x17\x20\x00\x00\xff\xff\x3b\x3f\x5f\x26\x7c\x00\x00\x00")

func exampleConfigBytes() ([]byte, error) {
	return bindataRead(
		_exampleConfig,
		"example/config",
	)
}

func exampleConfig() (*asset, error) {
	bytes, err := exampleConfigBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "example/config", size: 124, mode: os.FileMode(420), modTime: time.Unix(1451363314, 0)}
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
	"example/assignments/assignment.sample": exampleAssignmentsAssignmentSample,
	"example/config":                        exampleConfig,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"example": &bintree{nil, map[string]*bintree{
		"assignments": &bintree{nil, map[string]*bintree{
			"assignment.sample": &bintree{exampleAssignmentsAssignmentSample, map[string]*bintree{}},
		}},
		"config": &bintree{exampleConfig, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
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

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
