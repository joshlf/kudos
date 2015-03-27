package db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var pathTestCases = map[string]Path{
	"/":         Path{nil, true},
	"/foo":      Path{[]PathElement{{Name: "foo"}}, true},
	"/foo/bar":  Path{[]PathElement{{Name: "foo"}, {Name: "bar"}}, true},
	"/foo/*":    Path{[]PathElement{{Name: "foo"}, {Any: true}}, true},
	"/foo/*/*":  Path{[]PathElement{{Name: "foo"}, {Any: true}, {Any: true}}, true},
	"/foo/../*": Path{[]PathElement{{Name: "foo"}, {Up: true}, {Any: true}}, true},
	"":          Path{nil, false},
	"foo":       Path{[]PathElement{{Name: "foo"}}, false},
	"foo/bar":   Path{[]PathElement{{Name: "foo"}, {Name: "bar"}}, false},
	"foo/*":     Path{[]PathElement{{Name: "foo"}, {Any: true}}, false},
	"foo/*/*":   Path{[]PathElement{{Name: "foo"}, {Any: true}, {Any: true}}, false},
	"foo/../*":  Path{[]PathElement{{Name: "foo"}, {Up: true}, {Any: true}}, false},
}

func TestParsePath(t *testing.T) {
	for str, path := range pathTestCases {
		got := ParsePath(str)
		if !reflect.DeepEqual(got, path) {
			t.Errorf("parsing %v, expected %v; got %v", str, path, got)
		}
	}
}

func TestUnmarshalPath(t *testing.T) {
	for str, path := range pathTestCases {
		str = "\"" + str + "\""
		var got Path
		json.Unmarshal([]byte(str), &got)
		if !reflect.DeepEqual(got, path) {
			t.Errorf("parsing %v, expected %v; got %v", str, path, got)
		}
	}
}

var constraintEntityTestCases = map[string]ConstraintEntity{
	`{"path": "/foo/bar"}`:  ConstraintEntity{Path{[]PathElement{{Name: "foo"}, {Name: "bar"}}, true}},
	`{"value": "/foo/bar"}`: ConstraintEntity{"/foo/bar"},
	`{"value": 1}`:          ConstraintEntity{float64(1)},
	`{"value": true}`:       ConstraintEntity{true},
}

func TestUnmarshalConstraintEntity(t *testing.T) {
	for str, constraint := range constraintEntityTestCases {
		var got ConstraintEntity
		json.Unmarshal([]byte(str), &got)
		if !reflect.DeepEqual(got, constraint) {
			t.Errorf("parsing %v, expected %v; got %v", str, constraint, got)
		}
	}
}

func TestUnmarshalRelation(t *testing.T) {
	for str, rel := range relationKeywords {
		str = `"` + str + `"`
		var got Relation
		json.Unmarshal([]byte(str), &got)
		if !reflect.DeepEqual(got, rel) {
			t.Errorf("parsing %v, expected %v; got %v", str, rel, got)
		}
	}
	str := `"foo"`
	var got Relation
	err := json.Unmarshal([]byte(str), &got)
	if err == nil || err.Error() != "unknown relation: foo" {
		t.Errorf("unexpected error: expected %v; got %v", "unknown relation: foo", err)
	}
}

var valueTestCases = map[string]Value{
	`"foo"`: Value{"foo"},
	`1`:     Value{int64(1)},
	`1.0`:   Value{float64(1)},
}

func TestUnmarshalValue(t *testing.T) {
	for str, val := range valueTestCases {
		var got Value
		json.Unmarshal([]byte(str), &got)
		if !reflect.DeepEqual(got, val) {
			t.Errorf("parsing %v, expected %v; got %v", str, val, got)
		}
	}
	strs := []string{`true`, `{}`}
	for _, str := range strs {
		var got Value
		err := json.Unmarshal([]byte(str), &got)
		if err == nil || err.Error() != "must be string or number" {
			t.Errorf("unexpected error: expected %v; got %v", "must be string or number", err)
		}
	}
}

func TestPath(t *testing.T) {
	path := "/foo/bar"
	p := ParsePath(path)
	if p.HasWildcards() {
		t.Errorf("ParsePath(%q).HasWildcards() = true", path)
	}
	if p.HasUps() {
		t.Errorf("ParsePath(%q).HasUps() = true", path)
	}

	path = "/foo/*/bar"
	p = ParsePath(path)
	if !p.HasWildcards() {
		t.Errorf("ParsePath(%q).HasWildcards() = false", path)
	}
	if p.HasUps() {
		t.Errorf("ParsePath(%q).HasUps() = true", path)
	}

	path = "/foo/*/bar/.."
	p = ParsePath(path)
	if !p.HasWildcards() {
		t.Errorf("ParsePath(%q).HasWildcards() = false", path)
	}
	if !p.HasUps() {
		t.Errorf("ParsePath(%q).HasUps() = false", path)
	}
}

func ExamplePrint() {
	var c Constraint
	c.A.Entity = ParsePath("/foo/bar/*/baz")
	c.B.Entity = Value{"foo"}
	fmt.Println(c)

	// Output: (path:/foo/bar/*/baz < value:"foo")
}
