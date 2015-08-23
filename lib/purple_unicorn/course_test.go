package purple_unicorn

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCorrectCourseFile(t *testing.T) {
	basic1 := ` {
	"code" : "cs000",
	"name" : "test CS course",
	"description" : "a very simple cs course",
	"ta_group" : "staff",
	"tas" : ["jliebowf", "ezr"],
	"student_group" : "sys",
	"handin_method" : "facl"
}`
	r := bytes.NewBufferString(basic1)
	p, err := ParseCourse(r)
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		grp1 := Group("staff")
		grp2 := Group("sys")
		facl := "facl"
		expect := NewCourse("cs000", "test CS course", "a very simple cs course", &grp1, []User{User("jliebowf"), User("ezr")}, &grp2, &facl)
		if !reflect.DeepEqual(*p, expect) {
			t.Fail()
			t.Logf("Expected \n%v\n Got \n%v\n", *p, expect)
		}
	}
}
