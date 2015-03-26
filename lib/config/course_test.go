package config

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
)

var TestConfig = "testdata/sample_course_config.toml"

func init() {
	_, parentDir, _, _ := runtime.Caller(0)
	parentDir = filepath.Dir(parentDir)
	TestConfig = filepath.Join(parentDir, TestConfig)
}

func TestDecodeCourseConfig(t *testing.T) {
	expected := CourseConfig{
		Name:             "cs101",
		TaGroup:          "cs101ta",
		StudentGroup:     "cs101student",
		HandinDir:        HandinDir("handin"),
		ShortDescription: "CS 101",
		LongDescription:  "CS 101 is an introductory course in CS.",
	}

	var testStr []byte
	var config CourseConfig
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
	expected := CourseConfig{
		Name:             "cs101",
		TaGroup:          "cs101tas",
		StudentGroup:     "cs101students",
		HandinDir:        HandinDir("handin"),
		ShortDescription: "CS 101",
		LongDescription:  "This is an introductory course in CS.",
	}
	conf, err := ReadCourseConfig("cs101", "testdata")
	if err != nil {
		t.Errorf("ReadCourseConfig(\"cs101\", \"testdata\"): %v", err)
	} else if !reflect.DeepEqual(expected, conf) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, conf)
	}

	conf, err = ReadCourseConfig("cs102", "testdata")
	if err == nil || err.Error() != "course name in config (cs101) does not match expected name (cs102)" {
		t.Errorf("expected error:\ncourse name in config (cs101) does not match expected name (cs102)\n\ngot:\n%v", err)
	}
}
