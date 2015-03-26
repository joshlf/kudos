package config

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
)

var TestConfig = "sample_course_config.toml"

func init() {
	_, parentDir, _, _ := runtime.Caller(0)
	parentDir = filepath.Dir(parentDir)
	TestConfig = filepath.Join(parentDir, TestConfig)
}

func TestDecodeCourseConfig(t *testing.T) {
	expected := CourseConfig{
		Name:             "cs101",
		TaGroup:          "cs101tas",
		StudentGroup:     "cs101students",
		HandinDir:        HandinDir("handin"),
		ShortDescription: "CS 101",
		LongDescription:  "This is an introductory course in CS.",
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
