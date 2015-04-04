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
