package config

import (
	"io/ioutil"
	"reflect"
	"runtime"
	"testing"
	"time"
)

const TestFile = "testdata/sample_spec.toml"
const TestAssignmentConfig1 = "testdata/test_assignment_1.toml"
const TestAssignmentConfig2 = "testdata/test_assignment_2.toml"

func TestExample(t *testing.T) {
	tm, _ := timeparse("Jul 4, 2015 at 12:00am (EST)")
	expect := assignConfig{
		Code:    optionalCode{"assign01", true},
		Name:    optionalString{"Assignment 01", true},
		Due:     optionalDate{date(tm), true},
		Handins: nil,
		Problems: []problem{
			problem{
				Code:        optionalCode{"prob1", true},
				Name:        optionalString{"Problem 1", true},
				Points:      optionalNumber{50, true},
				SubProblems: nil,
			},
			problem{
				Code:   optionalCode{"prob2", true},
				Name:   optionalString{"Problem 2", true},
				Points: optionalNumber{0, false},
				SubProblems: []problem{
					problem{
						Code:        optionalCode{"a", true},
						Name:        optionalString{"", false},
						Points:      optionalNumber{25, true},
						SubProblems: nil,
					},
					problem{
						Code:        optionalCode{"b", true},
						Name:        optionalString{"", false},
						Points:      optionalNumber{25, true},
						SubProblems: nil,
					},
				},
			},
		},
	}

	conf, err := readAssignConfig("assign01", TestAssignmentConfig1)
	if err != nil {
		t.Fatalf("unexpected error reading config: %v", err)
	}
	if !reflect.DeepEqual(conf, expect) {
		t.Errorf("got unexpected config: %#v", conf)
	}

	expect.Due = optionalDate{set: false}
	expect.Handins = []handin{
		handin{
			Code:     optionalCode{"first", true},
			Due:      optionalDate{date(tm), true},
			Problems: []code{"prob1"},
		},
		handin{
			Code:     optionalCode{"second", true},
			Due:      optionalDate{date(tm.Add(time.Hour * time.Duration(24))), true},
			Problems: []code{"prob2"},
		},
	}

	conf, err = readAssignConfig("assign01", TestAssignmentConfig2)
	if err != nil {
		t.Fatalf("unexpected error reading config: %v", err)
	}
	if !reflect.DeepEqual(conf, expect) {
		t.Errorf("got unexpected config: %#v", conf)
	}
}

func TestAssignmentMethods(t *testing.T) {
	conf, err := readAssignConfig("assign01", TestAssignmentConfig1)
	if err != nil {
		t.Fatalf("unexpected error reading config: %v", err)
	}
	asgn := Assignment{conf: conf}
	tm, _ := timeparse("Jul 4, 2015 at 12:00am (EST)")
	expect := []interface{}{"assign01", "Assignment 01", tm, false, ".kudos/handin/assign01"}
	got := []interface{}{asgn.Code(), asgn.Name(), asgn.Due(), asgn.HasMultipleHandins(), asgn.HandinDir()}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("unexpected return values: want %v; got %v", expect, got)
	}

	// NOTE: The types here are really tricky,
	// and it's very easy to mess it up (since
	// they're slices of interfaces, and thus
	// there's no static checking). If the
	// test is failing, make sure to check that
	// all types line up properly (for example,
	// you did 50.0 instead of 50, type(nil)
	// instead of an untyped nil, etc).
	probs := asgn.Problems()
	expect = []interface{}{"prob1", "Problem 1", 50.0, []Problem(nil), "prob2", "Problem 2", 50.0, 2, "a", "a", 25.0,
		[]Problem(nil), "b", "b", 25.0, []Problem(nil)}
	got = []interface{}{probs[0].Code, probs[0].Name(), probs[0].Points(), probs[0].SubProblems, probs[1].Code,
		probs[1].Name(), probs[1].Points(), len(probs[1].SubProblems), probs[1].SubProblems[0].Code,
		probs[1].SubProblems[0].Name(), probs[1].SubProblems[0].Points(), probs[1].SubProblems[0].SubProblems,
		probs[1].SubProblems[1].Code, probs[1].SubProblems[1].Name(), probs[1].SubProblems[1].Points(),
		probs[1].SubProblems[1].SubProblems}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("unexpected return values: want \n%#v; \ngot \n%#v", expect, got)
	}

	testPanic(t, func() { asgn.Handins() }, "config: does not have multiple handins")
	asgn.conf.Name.set = false
	if asgn.Name() != "assign01" {
		t.Errorf("unexpected name; want %v; got %v", "assign01", asgn.Name())
	}

	asgn.conf.Due = optionalDate{set: false}
	asgn.conf.Handins = []handin{
		handin{
			Code:     optionalCode{"first", true},
			Due:      optionalDate{date(tm), true},
			Problems: []code{"prob1"},
		},
		handin{
			Code:     optionalCode{"second", true},
			Due:      optionalDate{date(tm.Add(time.Hour * time.Duration(24))), true},
			Problems: []code{"prob2"},
		},
	}
	handins := asgn.Handins()
	expect = []interface{}{"first", tm, 1, "second", tm.Add(time.Hour * time.Duration(24)), 1}
	got = []interface{}{handins[0].Code, handins[0].Due, len(handins[0].Problems),
		handins[1].Code, handins[1].Due, len(handins[1].Problems)}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("unexpected return values: want \n%#v; \ngot \n%#v", expect, got)
	}

	testPanic(t, func() { asgn.Due() }, "config: has multiple handins")
}

func TestReadAssignConfigError(t *testing.T) {
	testError(t, func() error { _, err := readAssignConfig("", "/nonexistant/file"); return err },
		"open /nonexistant/file: no such file or directory")

	tmp, err := ioutil.TempFile("", "test_kudos")
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
	var code string
	f := func() error { _, err := readAssignConfig(code, tmp.Name()); return err }

	writeContents("")
	testError(t, f, "assignment must have code")

	writeContents(`code = "foo"`)
	code = "bar"
	testError(t, f, "assignment code in config (foo) does not match expected code (bar)")

	code = "foo"
	testError(t, f, "assignment has no problems")

	contents := `code = "foo"
due = ""
`
	writeContents(contents)
	testError(t, f, `Type mismatch for 'config.assignConfig.due': parsing time "" as "Jan 2, 2006 at 3:04pm (MST)": cannot parse "" as "Jan"`)

	contents = `code = "foo"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[handin]]
[[problem]]
`
	writeContents(contents)
	testError(t, f, "assignment cannot have due date and handins")

	contents = `code = "foo"
[[handin]]
[[problem]]
`
	writeContents(contents)
	testError(t, f, "assignment cannot have one handin - instead just use a due date")

	contents = `code = "foo"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[problem]]
`
	writeContents(contents)
	testError(t, f, "all problems must have a code")

	contents = `code = "foo"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[problem]]
code = "prob1"
points = 0
[[problem.subproblem]]
`
	writeContents(contents)
	testError(t, f, "problem cannot have points and subproblems: prob1")

	contents = `code = "foo"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[problem]]
code = "prob1"
`
	writeContents(contents)
	testError(t, f, "problem must have points or subproblems: prob1")

	contents = `code = "foo"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[problem]]
code = "prob1"
[[problem.subproblem]]
code = "prob1"
`
	writeContents(contents)
	testError(t, f, "problem must have points or subproblems: prob1.prob1")

	contents = `code = "foo"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[problem]]
code = "prob1"
points = 0
[[problem]]
code = "prob1"
`
	writeContents(contents)
	testError(t, f, "duplicate problem code: prob1")

	contents = `code = "foo"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[problem]]
code = "prob1"
[[problem.subproblem]]
code = "prob1"
points = 0
[[problem.subproblem]]
code = "prob1"
points = 0
`
	writeContents(contents)
	testError(t, f, "duplicate subproblem code: prob1.prob1")

	contents = `code = "foo"
[[handin]]
[[handin]]
[[problem]]
code = "prob1"
points = 0
`
	writeContents(contents)
	testError(t, f, "all handins must have a code")

	contents = `code = "foo"
[[handin]]
code = "handin1"
[[handin]]
code = "handin2"
[[problem]]
code = "prob1"
points = 0
`
	writeContents(contents)
	testError(t, f, "handin must have due date: handin1")

	contents = `code = "foo"
[[handin]]
code = "handin1"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[handin]]
code = "handin2"
due = "Jan 2, 2006 at 3:04pm (MST)"
[[problem]]
code = "prob1"
points = 0
`
	writeContents(contents)
	testError(t, f, "handin must have problems: handin1")

	contents = `code = "foo"
[[handin]]
code = "handin1"
due = "Jan 2, 2006 at 3:04pm (MST)"
problems = ["prob1"]
[[handin]]
code = "handin1"
due = "Jan 2, 2006 at 3:04pm (MST)"
problems = ["prob2"]
[[problem]]
code = "prob1"
points = 0
`
	writeContents(contents)
	testError(t, f, "unknown problem in handin handin1: prob2")

	contents = `code = "foo"
[[handin]]
code = "handin1"
due = "Jan 2, 2006 at 3:04pm (MST)"
problems = ["prob1"]
[[handin]]
code = "handin1"
due = "Jan 2, 2006 at 3:04pm (MST)"
problems = ["prob1"]
[[problem]]
code = "prob1"
points = 0
`
	writeContents(contents)
	testError(t, f, "problem in multiple handins: prob1")

	contents = `code = "foo"
[[handin]]
code = "handin1"
due = "Jan 2, 2006 at 3:04pm (MST)"
problems = ["prob1"]
[[handin]]
code = "handin1"
due = "Jan 2, 2006 at 3:04pm (MST)"
problems = ["prob2"]
[[problem]]
code = "prob1"
points = 0
[[problem]]
code = "prob2"
points = 0
[[problem]]
code = "prob3"
points = 0
`
	writeContents(contents)
	testError(t, f, "problem not in any handins: prob3")
}
