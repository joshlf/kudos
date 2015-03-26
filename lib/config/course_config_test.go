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

func inti() {
	_, parentDir, _, _ := runtime.Caller(0)
	parentDir = filepath.Dir(parentDir)
	TestRubric = filepath.Join(parentDir, TestRubric)
}

func TestDecodeCourseConfig(t *testing.T) {
	expected := CourseConfig{
		Name:             "cs0",
		TaGroup:          "cs0tas",
		StudentGroup:     "cs0students",
		ShortDescription: "Test CS course",
		LongDescription: `This is a test file providing an example for
how course configs will look`,
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
