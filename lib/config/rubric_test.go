package config

import (
	"io/ioutil"
	"os"
	"reflect"
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
