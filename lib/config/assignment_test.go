package config

import (
	"reflect"
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
