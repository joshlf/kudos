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
// will be added as "b/c".
func ParseDir(dir string) (*template.Template, error) {
	var t *template.Template
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
			// tt is the template we're parsing,
			// which might be the root if this
			// is the first template file we've
			// encountered.
			var tt *template.Template
			if t == nil {
				t = template.New(name)
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
		return nil, fmt.Errorf("ParseDir: no template files in directory")
	}
	return t, nil
}
