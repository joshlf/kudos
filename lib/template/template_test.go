package template

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/synful/kudos/lib/testutil"
)

const templatesDir = "testdata/templates"

func TestParseDir(t *testing.T) {
	dir, ok := testutil.SrcDir()
	if !ok {
		t.Fatalf("could not locate source directory")
	}
	dir = filepath.Join(dir, templatesDir)
	tmpl, err := ParseDir(dir)
	testutil.Must(t, err)

	f := func(w io.Writer) error { return tmpl.ExecuteTemplate(w, "foo", nil) }
	runExpect(t, "foo\nbaz", f)
	f = func(w io.Writer) error { return tmpl.Execute(w, nil) }
	runExpect(t, "baz", f)

	tmpl, err = template.New("root").Parse(`root/{{ template "foo" }}`)
	testutil.Must(t, err)
	tmpl, err = ParseDirAdd(tmpl, dir)
	testutil.Must(t, err)
	f = func(w io.Writer) error { return tmpl.Execute(w, nil) }
	runExpect(t, "root/foo\nbaz", f)
	f = func(w io.Writer) error { return tmpl.ExecuteTemplate(w, "bar/baz", nil) }
	runExpect(t, "baz", f)

	empty := testutil.MustTempDir(t, "", "")
	defer os.Remove(empty)
	_, err = ParseDir(empty)
	want := "template: no files in directory"
	if err == nil || err.Error() != want {
		t.Errorf("unexpected error: got %v; want %v", err, want)
	}
}

func TestCommand(t *testing.T) {
	// We modify them in this function, so restore
	// them after we're done.
	defer func() {
		stdout, stderr = os.Stdout, os.Stderr
	}()

	env := []string{"FOO=BAR"}
	tmpl, err := template.New("").Parse("{{ . }}")
	testutil.Must(t, err)

	var stdoutmock, stderrmock bytes.Buffer
	stdout, stderr = &stdoutmock, &stderrmock
	cmd, err := CommandOut(tmpl, "env", 1, env)
	testutil.Must(t, err)
	testutil.Must(t, cmd.Run())
	if stderr := string(stderrmock.Bytes()); stderr != "" {
		t.Errorf("unexpected stderr: got %v; want empty", stderr)
	}
	fields := bytes.Fields(stdoutmock.Bytes())
	var found bool
	for _, f := range fields {
		if string(f) == env[0] {
			found = true
		}
	}
	if !found {
		t.Errorf("couldn't find expected env var in stdout: %v", env[0])
	}

	stdoutmock, stderrmock = bytes.Buffer{}, bytes.Buffer{}
	cmd, err = CommandOut(tmpl, "cat", 1, nil)
	testutil.Must(t, err)
	testutil.Must(t, cmd.Run())
	if stderr := string(stderrmock.Bytes()); stderr != "" {
		t.Errorf("unexpected stderr: got %v; want empty", stderr)
	}
	if stdout := string(stdoutmock.Bytes()); stdout != "1" {
		t.Errorf("unexpected stdout: got %v; want %v", stdout, "1")
	}
}

func TestRunAll(t *testing.T) {
	var entries []Entry
	var env []string
	for i := 0; i < 10; i++ {
		env = append(env, fmt.Sprintf("TEMPLATE_TEST_%v=foo", i))
		entries = append(entries, Entry{i, env})
	}

	tmpl, err := template.New("").Parse("{{ . }}")
	testutil.Must(t, err)
	var stdoutmock, stderrmock bytes.Buffer
	stdout, stderr = &stdoutmock, &stderrmock
	testutil.Must(t, RunAllOut(tmpl, "env", entries...))

	if stderr := string(stderrmock.Bytes()); stderr != "" {
		t.Errorf("unexpected stderr: got %v; want empty", stderr)
	}
	// Assume env will print variables in order.
	// This is probably a good assumption, but
	// fix this code if the assumption ever breaks.
	var relevantLines []string
	envmap := make(map[string]bool)
	for _, e := range env {
		envmap[e] = true
	}
	// This might not actually be lines (eg, if there are
	// spaces in certain lines), but we don't care about
	// those lines anyway, so it doesn't matter.
	lines := strings.Fields(string(stdoutmock.Bytes()))

	// Filter out only the lines we care about
	for _, l := range lines {
		if envmap[l] {
			relevantLines = append(relevantLines, l)
		}
	}

	// 55 = 10 + 9 + ... + 1
	if len(relevantLines) != 55 {
		t.Errorf("unexpected len(relevantLines): got %v; want %v", len(relevantLines), 55)
	}

	// Assume the first output is 1 line, then 2, then
	// 3, and so on. Keep consuming the next n lines
	// until we run out, and expect it to match env
	// the whole time.
	for i := 0; i < 10; i++ {
		got := relevantLines[:i+1]
		want := env[:i+1]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected output: got %v; want %v", got, want)
		}
		relevantLines = relevantLines[i+1:]
	}

	entries = []Entry{{0, nil}}
	err = RunAll(tmpl, "", entries...)
	enterr, ok := err.(EntryError)
	if !ok {
		t.Errorf("unexpected error type: want %T; got %T", EntryError{}, err)
	}
	if enterr.Index != 0 || !reflect.DeepEqual(enterr.Entry, Entry{0, nil}) ||
		enterr.Error() != "fork/exec : no such file or directory" {
		// Use %#v so Error() method is not used
		t.Errorf("unexpected error: got %#v; want %#v",
			enterr, EntryError{10, Entry{0, nil}, fmt.Errorf("fork/exec : no such file or directory")})
	}
}

// Run f and expect the given string to be
// written to f's argument.
func runExpect(t *testing.T, expect string, f func(io.Writer) error) {
	var buf bytes.Buffer
	err := f(&buf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	got := string(buf.Bytes())
	if got != expect {
		t.Errorf("unexpected output: got %q; want %q", got, expect)
	}
}
