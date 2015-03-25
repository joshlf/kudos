package assignment

import (
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
				Total: 40,
			},
			Problem{
				Name:  "oh great!",
				Files: []string{"file2.c"},
				Total: 20,
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
