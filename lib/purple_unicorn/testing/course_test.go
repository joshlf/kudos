package testing

import (
	"bytes"
	"reflect"
	"testing"

	pu "github.com/synful/kudos/lib/purple_unicorn"
	"github.com/synful/kudos/lib/yellow_dingo"
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
	p, err := yellow_dingo.ParseCourse(r)
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		grp1 := pu.Group("staff")
		grp2 := pu.Group("sys")
		facl := "facl"
		expect := pu.NewCourse("cs000", "test CS course", "a very simple cs course", &grp1, []pu.User{pu.User("jliebowf"), pu.User("ezr")}, &grp2, &facl)
		if !reflect.DeepEqual(*p, expect) {
			t.Fail()
			t.Logf("Expected \n%v\n Got \n%v\n", *p, expect)
		}
	}
}
