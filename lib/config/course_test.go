package config

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/synful/kudos/lib/perm"
)

var TestConfig = "testdata/test_course_config.toml"

func init() {
	_, parentDir, _, _ := runtime.Caller(0)
	parentDir = filepath.Dir(parentDir)
	TestConfig = filepath.Join(parentDir, TestConfig)
}

func TestCourse(t *testing.T) {
	c := Course{
		path: "/foo/bar/baz",
		config: courseConfig{
			Code:         optionalCode{"foo", true},
			Name:         optionalString{"Foo", true},
			TaGroup:      optionalString{"foota", true},
			StudentGroup: optionalString{"foostudent", true},
			HandinMethod: optionalHandinMethod{handinMethod(FaclMethod), true},
			Description:  optionalString{"An introduction to foo.", true},
		},
	}

	expect := []interface{}{"foo", "Foo", "foota", "foostudent", FaclMethod, true, "An introduction to foo."}
	got := []interface{}{c.Code(), c.Name(), c.TaGroup(), c.StudentGroup(), c.HandinMethod(), c.DescriptionSet(), c.Description()}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("unexpected return values; want %#v; got %#v", expect, got)
	}

	expect = []interface{}{"/foo/bar/baz/.kudos", "/foo/bar/baz/.kudos/config.toml", "/foo/bar/baz/.kudos/handin", "/foo/bar/baz/.kudos/assignments"}
	got = []interface{}{c.ConfigDir(), c.ConfigFile(), c.HandinDir(), c.AssignmentsDir()}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("unexpected return values; want %#v; got %#v", expect, got)
	}

	c.config.Name = optionalString{"", false}
	c.config.Description = optionalString{"", false}
	expect = []interface{}{"foo", ""}
	got = []interface{}{c.Name(), c.Description()}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("unexpected return values; want %#v; got %#v", expect, got)
	}
}

func TestDecodeCourseConfig(t *testing.T) {
	expected := courseConfig{
		Code:         optionalCode{"cs101", true},
		Name:         optionalString{"CS 101", true},
		TaGroup:      optionalString{"cs101ta", true},
		StudentGroup: optionalString{"cs101student", true},
		HandinMethod: optionalHandinMethod{handinMethod(FaclMethod), true},
		Description:  optionalString{"CS 101 is an introductory course in CS.", true},
	}

	var testStr []byte
	var config courseConfig
	var err error

	if testStr, err = ioutil.ReadFile(TestConfig); err != nil {
		t.Fatalf("Cannot find %v file!", TestConfig)
	}
	if _, err = toml.Decode(string(testStr), &config); err != nil {
		t.Fatalf("Error decoding file:\n\t %v", err)
	}

	if !reflect.DeepEqual(expected, config) {
		t.Fatalf("Expected\n%v\n, Got \n%v\n", expected, config)
	}
}

func TestReadCourseConfig(t *testing.T) {
	expected := Course{
		"testdata",
		courseConfig{
			Code:         optionalCode{"cs101", true},
			Name:         optionalString{"CS 101", true},
			TaGroup:      optionalString{"cs101ta", true},
			StudentGroup: optionalString{"cs101student", true},
			HandinMethod: optionalHandinMethod{handinMethod(FaclMethod), true},
			Description:  optionalString{"CS 101 is an introductory course in CS.", true},
		},
	}
	course, err := ReadCourseConfig("cs101", "testdata")
	if err != nil {
		t.Errorf("ReadCourseConfig(\"cs101\", \"testdata\"): %v", err)
	} else if !reflect.DeepEqual(expected, course) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, course)
	}

	course, err = ReadCourseConfig("cs102", "testdata")
	if err == nil || err.Error() != "course code in config (cs101) does not match expected code (cs102)" {
		t.Errorf("expected error:\ncourse code in config (cs101) does not match expected code (cs102)\n\ngot:\n%v", err)
	}
}

func TestReadCourseConfigError(t *testing.T) {
	testError(t, func() error { _, err := ReadCourseConfig("", ""); return err },
		"could not parse course config: open .kudos/config.toml: no such file or directory")

	dir, err := ioutil.TempDir("", "test_kudos_course_dir")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	err = os.Mkdir(filepath.Join(dir, ".kudos"), 0700)
	if err != nil {
		t.Fatalf("could not create kudos dir: %v", err)
	}
	tmp, err := os.Create(filepath.Join(dir, ".kudos/config.toml"))
	if err != nil {
		t.Fatalf("could not create config file: %v", err)
	}

	writeContents := func(contents string) {
		_, _, line, _ := runtime.Caller(1)
		err := tmp.Truncate(0)
		if err != nil {
			t.Fatalf("line %v: could not truncate: %v", line, err)
		}
		_, err = tmp.Write([]byte(contents))
		if err != nil {
			t.Fatalf("line %v: could not write contents: %v", line, err)
		}
		// Otherwise os.Open() will return a file
		// which is already seeked to the end.
		_, err = tmp.Seek(0, 0)
		if err != nil {
			t.Fatalf("line %v: could not seek: %v", line, err)
		}
	}
	var code string
	f := func() error { _, err := ReadCourseConfig(code, dir); return err }

	writeContents("bad = ")
	testError(t, f, "could not parse course config: Near line 1 (last key parsed 'bad'): Expected value but found '\\n' instead.")

	writeContents("")
	testError(t, f, "course code must be set")

	writeContents(`code = "foo"`)
	code = "bar"
	testError(t, f, "course code in config (foo) does not match expected code (bar)")

	code = "foo"
	contents := `code = "foo"`
	writeContents(contents)
	testError(t, f, "ta_group must be set")

	contents = `code = "foo"
ta_group = "foota"`
	writeContents(contents)
	testError(t, f, "student_group must be set")

	contents = `code = "foo"
ta_group = "foota"
student_group = "foostudent"`
	writeContents(contents)
	testError(t, f, "handin_method must be set")

	contents = `code = "foo"
ta_group = "foota"
student_group = "foostudent"
handin_method = "foo"`
	writeContents(contents)
	testError(t, f, "could not parse course config: Type mismatch for 'config.courseConfig.handin_method': allowed methods: facl, setgid")
}

func TestInitCourse(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get pwd: %v", err)
	}

	course := "cs101"
	coursePath := filepath.Join(pwd, course)
	err = os.Mkdir(coursePath, os.ModeDir|perm.Parse("rwx------"))
	if err != nil {
		t.Fatalf("could not create course directory: %v", err)
	}
	defer os.RemoveAll(coursePath)
	err = InitCourse(course, coursePath, false)
	if err != nil {
		t.Fatalf("InitCourse(%v, %v, true): %v", course, coursePath, err)
	}

	err = exec.Command("diff", "-rN", "cs101/.kudos/", "example").Run()
	if err != nil {
		t.Errorf("unexpected error running diff: %v", err)
	}
}
