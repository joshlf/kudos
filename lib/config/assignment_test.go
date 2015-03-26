package config

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

const TESTFILE1 = "sample_spec.toml"

func TestDecoding(t *testing.T) {

	var asgn AssignSpec
	tm, _ := timeparse("Jul 4, 2015 at 12:00am (EST)")
	dur, _ := time.ParseDuration("4m")
	expected := AssignSpec{
		Title: "generic-assignment",
		Problem: []Problem{
			Problem{
				Name:  "funtimes!",
				Files: []string{"file1.c", "file1.readme"},
				Total: GradeNum{40},
			},
			Problem{
				Name:  "oh great!",
				Files: []string{"file2.c"},
				Total: GradeNum{20},
			},
		},
		Handin: Handin{
			Due:   date{tm},
			Grace: duration{dur},
		},
	}
	var testStr []byte
	var err error

	if testStr, err = ioutil.ReadFile(TESTFILE1); err != nil {
		t.Fatalf("Cannot find %v file!", TESTFILE1)
	}

	if _, err = toml.Decode(string(testStr), &asgn); err != nil {
		t.Fatalf("Error decoding file:\n\t %v", err)
	}

	if !reflect.DeepEqual(asgn, expected) {
		t.Fatalf("Expected \n%v\n, Got \n%v\n", expected, asgn)
	}

}

func TestEmitRubric(t *testing.T) {
	tm, _ := timeparse("Jul 4, 2015 at 12:00am (EST)")
	dur, _ := time.ParseDuration("4m")
	var b bytes.Buffer

	orig := &AssignSpec{
		Title: "generic-assignment",
		Problem: []Problem{
			Problem{
				Name:  "funtimes!",
				Files: []string{"file1.c", "file1.readme"},
				Total: GradeNum{40},
			},
			Problem{
				Name:  "oh great!",
				Files: []string{"file2.c"},
				Total: GradeNum{20},
			},
		},
		Handin: Handin{
			Due:   date{tm},
			Grace: duration{dur},
		},
	}

	orig.Rubric().WriteTOML(&b)

	var rubric Rubric
	if _, err := toml.DecodeReader(&b, &rubric); err != nil {
		t.Fatal("Unable to re-decode the rubric template:\n%v", err)
	}

	expected := Rubric{
		Assignment: "generic-assignment",
		Grader:     "",
		Grade: []Grade{
			Grade{
				Problem:  "funtimes!",
				Comment:  "",
				Score:    GradeNum{0},
				Possible: GradeNum{40},
			},
			Grade{
				Problem:  "oh great!",
				Comment:  "",
				Score:    GradeNum{0.},
				Possible: GradeNum{20},
			},
		},
	}
	if !reflect.DeepEqual(expected, rubric) {
		t.Fatalf("Expected \n%v\n, Got \n%v\n", expected, rubric)
	}

}