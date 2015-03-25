package assignment

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

const TEST_RUBRIC_1 = "sample_rubric.toml"

func TestDecodeRubric(t *testing.T) {
	var rubric Rubric
	expected := Rubric{
		Assignment: "sample assignment!",
		Grader:     "ezr",
		Grades: []Grade{
			Grade{
				Problem: "is-cookie-good?",
				Comment: `Comments for problem 1 here!

wow you did so well!
`,
				Score: GradeNum{40},
				Total: GradeNum{70},
			},
			Grade{
				Problem: "collab policy stuff",
				Comment: `Make sure to read the collaboration policy!
`,
				Score: GradeNum{0.5},
				Total: GradeNum{15},
			},
		},
	}
	var testStr []byte
	var err error

	if testStr, err = ioutil.ReadFile(TEST_RUBRIC_1); err != nil {
		t.Fatalf("Cannot find %v file!", TEST_RUBRIC_1)
	}

	if _, err = toml.Decode(string(testStr), &rubric); err != nil {
		t.Fatalf("Error decoding file:\n\t %v", err)
	}

	if !reflect.DeepEqual(rubric, expected) {
		t.Fatalf("Expected\n%v\n, Got \n%v\n", expected, rubric)
	}
}
