package yellow_dingo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/synful/kudos/lib/purple_unicorn"
)

type course struct {
	Code         *string  `json:"code"`
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	TaGroup      *string  `json:"ta_group"`
	TAs          []string `json:"tas"`
	StudentGroup *string  `json:"student_group"`
	HandinMethod *string  `json:"handin_method"`
}

func ParseCourseFile(f string) (*purple_unicorn.Course, error) {
	if r, err := os.Open(f); err == nil {
		if p, err := ParseCourse(r); err != nil {
			return nil, fmt.Errorf("error parsing file %v:\n%v", f, err)
		} else {
			return p, nil
		}
	} else {
		return nil, fmt.Errorf("unable to open course file %v", f)
	}
}

func ParseCourse(r io.Reader) (*purple_unicorn.Course, error) {
	var errs purple_unicorn.ErrList
	d := json.NewDecoder(r)
	var c course
	err := d.Decode(&c)
	if err != nil {
		return nil, err
	}
	course := purple_unicorn.Course{}
	if c.Code == nil {
		errs.Add(fmt.Errorf("must add \"code\" field"))
	} else {
		course.SetCodeNoValidate(purple_unicorn.Code(*c.Code))
	}
	if c.Name == nil {
		errs.Add(fmt.Errorf("must add \"name\" field"))
	} else {
		course.SetNameNoValidate(*c.Name)
	}
	if c.Description != nil {
		course.SetDescriptionNoValidate(*c.Description)
	}
	if c.TaGroup != nil {
		taGroup := purple_unicorn.Group(*c.TaGroup)
		course.SetTaGroupNoValidate(&taGroup)
	}
	for _, ta := range c.TAs {
		course.AddTANoValidate(purple_unicorn.User(ta))
	}
	if c.StudentGroup != nil {
		studGroup := purple_unicorn.Group(*c.StudentGroup)
		course.SetStudentGroupNoValidate(&studGroup)
	}
	course.SetHandinMethodNoValidate(c.HandinMethod)

	if len(errs) > 0 {
		return nil, errs
	}
	return &course, nil
}
