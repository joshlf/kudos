package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// ParseDir walks dir and creates a template containing
// all files it finds. The name of each template will
// be its path relative to dir. For example, if dir is
// "foo" and foo contains:
//  foo/
//      a
//      b/
//        c
// then "foo/a" will be added as "a", and "foo/b/c"
// will be added as "b/c". The returned template itself
// will be that of the first file encountered. There
// must be at least one file in the directory.
func ParseDir(dir string) (*template.Template, error) {
	return ParseDirAdd(nil, dir)
}

// ParseDirAdd is like ParseDir, except that it adds
// the parsed templates to t instead of creating a new
// template.
func ParseDirAdd(t *template.Template, dir string) (*template.Template, error) {
	dir = filepath.Clean(dir)
	prefixlen := len(dir + "/") // filepath.Clean removes trailing slash
	f := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			name := path[prefixlen:]
			var tt *template.Template
			if t == nil {
				t = template.New(name)
			}
			if name == t.Name() {
				tt = t
			} else {
				tt = t.New(name)
			}
			_, err = tt.Parse(string(b))
			return err
		}
		return nil
	}
	err := filepath.Walk(dir, f)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("template: no files in directory")
	}
	return t, nil
}
