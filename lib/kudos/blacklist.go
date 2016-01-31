package kudos

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var blacklistTemplate *template.Template

func init() {
	tmpl := `{{range . }}{{ . }}
{{end}}`
	blacklistTemplate = template.Must(template.New("").Parse(tmpl))
}

func WriteBlacklistFile(path string, uids ...string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = blacklistTemplate.Execute(f, uids)
	if err != nil {
		return err
	}
	return f.Sync()
}

// ParseBlacklistFile reads the file at path and returns
// its contents as a slice of strings, one per line
// (empty strings or all whitespace will be ommitted).
// If the file is a valid blacklist file, then these
// strings will be valid UIDs, but ParseBlacklistFile
// only validates that each line contains a single
// string of numerical characters. It does not validate
// that these are valid UIDs, or even that they represent
// numbers which fit in any particular range (ie, that
// of a 32-bit unsigned integer).
func ParseBlacklistFile(path string) (uids []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(f)
	line := -1 // so we can put line++ at the beginning
	for s.Scan() {
		line++
		parts := strings.Fields(s.Text())
		switch {
		case len(parts) == 0:
			continue
		case len(parts) > 1:
			return nil, fmt.Errorf("line %v: found multiple tokens (expected 1)", line)
		case !isNumerical(parts[0]):
			return nil, fmt.Errorf("line %v: non-numerical uid", line)
		}
		uids = append(uids, parts[0])
	}
	return uids, nil
}

// expects len(s) > 0
func isNumerical(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
