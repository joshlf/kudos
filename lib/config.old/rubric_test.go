package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"
)

const TestRubric = "testdata/test_rubric.toml"

func TestWriteRubricFile(t *testing.T) {
	expect := `assignment = "assign01"

[[grade]]
  problem = "prob1"
  comment = ""
  score = # out of 50 points

[[grade]]
  problem = "prob2"
  comment = ""

  [[grade.grade]]
    problem = "a"
    comment = ""
    score = # out of 25 points

  [[grade.grade]]
    problem = "b"
    comment = ""
    score = # out of 25 points
`

	a, err := readAssignConfig("assign01", TestAssignmentConfig1)
	if err != nil {
		t.Fatalf("error reading assignment config: %v",
			err)
	}
	err = WriteRubricFile(Assignment{a, Course{}}, "test_rubric.toml")
	if err != nil {
		t.Fatalf("error writing rubric: %v",
			err)
	}
	defer os.Remove("test_rubric.toml")
	data, err := ioutil.ReadFile("test_rubric.toml")
	if err != nil {
		t.Fatalf("error reading rubric file: %v", err)
	}
	if string(data) != expect {
		t.Errorf("unexpected rubric contents:\n%v", string(data))
	}
}

func TestWriteRubricFileError(t *testing.T) {
	testPanic(t, func() { WriteRubricFileProblems(Assignment{}, nil, "") }, "config: no problems specified")
	testPanic(t, func() { WriteRubricFileProblems(Assignment{}, []string{"foo"}, "") }, "config: unkown problem code")
}

func TestReadRubricFile(t *testing.T) {
	expect := Rubric{
		r: rubricFile{
			AssignmentCode: optionalCode{"assign01", true},
			Grades: []grade{
				grade{
					Problem: optionalCode{"prob1", true},
					Comment: "Almost",
					Score:   optionalNumber{45, true},
					Grades:  nil},
				grade{
					Problem: optionalCode{"prob2", true},
					Grades: []grade{
						grade{
							Problem: optionalCode{"a", true},
							Comment: "Good!",
							Score:   optionalNumber{25, true},
							Grades:  nil},
						grade{
							Problem: optionalCode{"b", true},
							Comment: "You missed a few",
							Score:   optionalNumber{15, true},
							Grades:  nil,
						},
					},
				},
			},
		},
	}

	a, err := readAssignConfig("assign01", TestAssignmentConfig1)
	if err != nil {
		t.Fatalf("error reading assignment config: %v", err)
	}
	r, err := ReadRubricFile(Assignment{a, Course{}}, TestRubric)
	if err != nil {
		t.Fatalf("error reading rubric: %v", err)
	}
	if !reflect.DeepEqual(expect, r) {
		t.Errorf("unexpected rubric: %v", r)
	}
}

func TestReadRubricFileError(t *testing.T) {
	testError(t, func() error { _, err := ReadRubricFile(Assignment{}, ""); return err },
		"could not parse rubric: open : no such file or directory")

	tmp, err := ioutil.TempFile("", "test_kudos_rubric")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
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
	asgn, err := readAssignConfig("assign01", TestAssignmentConfig1)
	if err != nil {
		t.Fatalf("unexpected error reading asignment: %v", err)
	}
	f := func() error { _, err := ReadRubricFile(Assignment{conf: asgn}, tmp.Name()); return err }

	writeContents("bad = ")
	testError(t, f, "could not parse rubric: Near line 1 (last key parsed 'bad'): Expected value but found '\\n' instead.")

	writeContents("")
	testError(t, f, "must have assignment code")

	writeContents(`assignment = "assign02"`)
	testError(t, f, "assignment code in rubric (assign02) does not match expected code (assign01)")

	writeContents(`assignment = "assign01"`)
	testError(t, f, "rubric has no grades")

	contents := `assignment = "assign01"
[[grade]]`
	writeContents(contents)
	testError(t, f, "all grades must specify problem code")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob"`
	writeContents(contents)
	testError(t, f, "unknown problem: prob")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob"`
	writeContents(contents)
	testError(t, f, "unknown problem: prob")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob1"
[[grade.grade]]
problem = "prob"`
	writeContents(contents)
	testError(t, f, "unknown problem: prob1.prob")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob1"
score = 0
[[grade.grade]]
problem = "prob"`
	writeContents(contents)
	testError(t, f, "grade cannot specify score and sub-grades: prob1")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob1"`
	writeContents(contents)
	testError(t, f, "grade must specify score or sub-grades: prob1")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob1"
score = 0
[[grade]]
problem = "prob2"`
	writeContents(contents)
	testError(t, f, "grade must specify score or sub-grades: prob2")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob1"
score = 0
[[grade]]
problem = "prob2"
[[grade.grade]]
problem = "a"`
	writeContents(contents)
	testError(t, f, "grade must specify score or sub-grades: prob2.a")

	contents = `assignment = "assign01"
[[grade]]
problem = "prob1"
score = 0
[[grade]]
problem = "prob2"
[[grade.grade]]
problem = "a"
score = 0
[[grade.grade]]
problem = "a"`
	writeContents(contents)
	testError(t, f, "duplicate grade for problem: prob2.a")
}

func TestRubricMethods(t *testing.T) {
	a, err := readAssignConfig("assign01", TestAssignmentConfig1)
	if err != nil {
		t.Fatalf("error reading assignment config: %v", err)
	}
	r, err := ReadRubricFile(Assignment{a, Course{}}, TestRubric)
	if err != nil {
		t.Fatalf("error reading rubric: %v", err)
	}

	code := r.AssignmentCode()
	if code != "assign01" {
		t.Errorf("unexpected course code: want %v; got %v", "assign01", code)
	}

	expect := []Grade{{"prob1", "Almost", 45, nil}, {"prob2", "", 0, []Grade{{"a", "Good!", 25, nil}, {"b", "You missed a few", 15, nil}}}}
	grades := r.Grades()
	if !reflect.DeepEqual(expect, grades) {
		t.Errorf("unexpected grades slice: want %v; got %v", expect, grades)
	}

	expectScores := []float64{45, 40, 25, 15}
	scores := []float64{grades[0].Score(), grades[1].Score(), grades[1].Grades[0].Score(), grades[1].Grades[1].Score()}
	if !reflect.DeepEqual(expectScores, scores) {
		t.Errorf("unexpected scores slice: want %v; got %v", expect, scores)
	}
}
